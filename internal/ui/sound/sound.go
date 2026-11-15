// Package sound 是 remake 的音訊播放層。
//
// 音檔**不隨發行包散布**——那是原版衍生物（`CLAUDE.md` §9）。
// 玩家拿自己的原版跑一次 `tools/bgm2ogg.sh` 產生，放進音檔目錄。
// **找不到檔案就靜音跑**，不 fallback 到自製音樂（`docs/spec/29` §5）。
package sound

import (
	"bytes"
	"fmt"
	gameaudio "github.com/wicanr2/wolong_cht/internal/audio"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const sampleRate = 44100

// Bank 是一組音檔。零值可以直接用（等於「沒有音檔」，全部呼叫變成 no-op），
// 所以呼叫端不必到處判斷 nil。
type effectVoice interface {
	Play()
	Pause()
	Close() error
	IsPlaying() bool
	SetVolume(float64)
}
type Bank struct {
	effect          effectVoice
	effectCode      int
	effectUntil     time.Time
	effectDurations map[int]time.Duration
	effectFactory   func([]byte) (effectVoice, error)
	effectNow       func() time.Time

	dir     string
	music   []string // 曲名（不含副檔名），排序過
	effects map[int]string

	mu      sync.Mutex
	ctx     *audio.Context
	player  *audio.Player
	current string
	enabled bool
	// initErr 記下第一次建 player 的失敗。**不重試**——headless 或
	// 沒有音效裝置的環境每一幀重試一次會把主迴圈拖垮。
	initErr error
	// level 是系統選單「音　效」的 TYPE 1–4，存成 0–3
	// （`docs/spec/122`）。原版是加進 OPL 載波的 Total Level，
	// remake 播的是事先算好的 OGG，所以接成播放增益。
	level int
	// silent 是**驗收模式**：狀態全部照舊，只是不碰音效裝置
	// （docs/spec/29 §5.1）。截圖與逐幀錄製跑在沒有音效卡的容器裡，
	// 而把音檔目錄清空會讓 Available() 變成 false，
	// 系統選單那一格就從「TYPE 1」變成「未接入」——
	// **驗收捷徑不可以改到被驗收的畫面**。
	silent bool
}

// Open 掃描音檔目錄。目錄不存在也回傳一個可用的 Bank（沒有音檔而已）。
func Open(dir string) *Bank {
	b := &Bank{dir: dir, effects: map[int]string{}, enabled: true}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return b
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".ogg") {
			continue
		}
		base := strings.TrimSuffix(name, ".ogg")
		var id int
		if n, _ := fmt.Sscanf(base, "sfx-%d", &id); n == 1 {
			b.effects[id] = name
			continue
		}
		b.music = append(b.music, base)
	}
	sort.Strings(b.music)
	return b
}

// SetSilent 開關驗收模式：不出聲，但 Available／Enabled 與選單顯示照舊。
func (b *Bank) SetSilent(on bool) {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.silent = on
	if on {
		b.stopEffectLocked()
	}
	if on && b.player != nil {
		b.player.Pause()
		_ = b.player.Close()
		b.player, b.current = nil, ""
	}
}

// Silent 回報現在是不是驗收模式。
func (b *Bank) Silent() bool { return b != nil && b.silent }

// Available 回報有沒有音檔。系統選單靠它顯示「未接入」。
func (b *Bank) Available() bool { return b != nil && (len(b.music) > 0 || len(b.effects) > 0) }

// Enabled／SetEnabled 是系統選單那一列的開關。
func (b *Bank) Enabled() bool { return b != nil && b.enabled }

func (b *Bank) SetEnabled(on bool) {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.enabled = on
	if !on {
		b.stopEffectLocked()
	}
	if !on && b.player != nil {
		b.player.Pause()
	} else if on && b.player != nil {
		b.player.Play()
	}
}

// levelCount 是 TYPE 的段數（TYPE 1–4）。
const levelCount = 4

// LevelGain 是第 n 段（TYPE n+1）的播放增益。
//
// 原版每段加 4 個 OPL Total Level 單位，而 TL 是 0.75 dB/step，
// 所以一段是 3 dB（`docs/re/81` §2）。
func LevelGain(n int) float64 {
	return math.Pow(10, -float64(n)*4*0.75/20)
}

// Level／SetLevel 是系統選單那一列的 TYPE（0 ＝ TYPE 1）。
func (b *Bank) Level() int {
	if b == nil {
		return 0
	}
	return b.level
}

func (b *Bank) SetLevel(n int) {
	if b == nil {
		return
	}
	if n < 0 {
		n = 0
	}
	if n >= levelCount {
		n = levelCount - 1
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.level = n
	if b.effect != nil {
		b.effect.SetVolume(LevelGain(n))
	}
	if b.player != nil {
		b.player.SetVolume(LevelGain(n))
	}
}

// Music 回傳可用的曲名，給系統選單或除錯用。
func (b *Bank) Music() []string {
	if b == nil {
		return nil
	}
	return b.music
}

func (b *Bank) context() *audio.Context {
	if b.ctx == nil {
		b.ctx = audio.NewContext(sampleRate)
	}
	return b.ctx
}

// PlayMusic 換一首背景音樂並無限循環。同一首重複呼叫不會重頭播。
func (b *Bank) PlayMusic(name string) {
	if b == nil || !b.enabled || b.initErr != nil || b.Silent() {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if name == b.current && b.player != nil {
		return
	}
	path := filepath.Join(b.dir, name+".ogg")
	raw, err := os.ReadFile(path)
	if err != nil {
		return // 沒有這首就靜音，不是錯誤
	}
	stream, err := vorbis.DecodeF32(bytes.NewReader(raw))
	if err != nil {
		b.initErr = err
		return
	}
	loop := audio.NewInfiniteLoopF32(stream, stream.Length())
	player, err := b.context().NewPlayerF32(loop)
	if err != nil {
		b.initErr = err
		return
	}
	if b.player != nil {
		b.player.Pause()
		_ = b.player.Close()
	}
	b.player, b.current = player, name
	player.SetVolume(LevelGain(b.level))
	player.Play()
}

// StopMusic 停掉目前這首。
func (b *Bank) StopMusic() {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.player != nil {
		b.player.Pause()
		_ = b.player.Close()
		b.player, b.current = nil, ""
	}
}

// SetEffectTimings從原始SOUND.DAT取得鏈長，不以OGG殘響推定busy。
// 遊戲adapter在啟動時提供；不寫入設定、存檔或音檔metadata。
func (b *Bank) SetEffectTimings(data []byte) error {
	if b == nil {
		return nil
	}
	effects, err := gameaudio.Effects(data)
	if err != nil {
		return err
	}
	durations := make(map[int]time.Duration, len(effects))
	for id := range effects {
		d, err := gameaudio.EffectDuration(effects, id)
		if err != nil {
			return err
		}
		durations[id] = d
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.effectDurations = durations
	return nil
}

// allowEffect對應KI102F5：相同或較高優先度可打斷，忙碌時拒絕較低優先度。
func allowEffect(last, code int, busy bool) bool { return code <= last || !busy }
func (b *Bank) stopEffectLocked() {
	if b.effect != nil {
		b.effect.Pause()
		_ = b.effect.Close()
		b.effect = nil
	}
	b.effectUntil = time.Time{}
}
func (b *Bank) createEffectPlayer(raw []byte) (effectVoice, error) {
	stream, err := vorbis.DecodeF32(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	player, err := b.context().NewPlayerF32(stream)
	if err != nil {
		b.initErr = err
	}
	return player, err
}

// PlayEffect依原版優先度選單一效果聲部，接受時停止舊聲並重新起音。
// 未提供原始時序的獨立Bank以OGG播放狀態fallback，不作原版時序聲明。
func (b *Bank) PlayEffect(id int) {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.enabled || b.initErr != nil || b.silent || id < 0 || id > 255 {
		return
	}
	if id == 0 {
		b.stopEffectLocked()
		b.effectCode = 0
		return
	}
	name, ok := b.effects[id]
	if !ok {
		return
	}
	clock := b.effectNow
	if clock == nil {
		clock = time.Now
	}
	now := clock()
	busy := false
	if b.effect != nil {
		if _, ok := b.effectDurations[b.effectCode]; ok {
			busy = now.Before(b.effectUntil)
		} else {
			busy = b.effect.IsPlaying()
		}
	}
	if !allowEffect(b.effectCode, id, busy) {
		return
	}
	raw, err := os.ReadFile(filepath.Join(b.dir, name))
	if err != nil {
		return
	}
	factory := b.effectFactory
	if factory == nil {
		factory = b.createEffectPlayer
	}
	player, err := factory(raw)
	if err != nil {
		return
	}
	b.stopEffectLocked()
	b.effect = player
	b.effectCode = id
	if d, ok := b.effectDurations[id]; ok {
		b.effectUntil = clock().Add(d)
	}
	player.SetVolume(LevelGain(b.level))
	player.Play()
}

// Err 回報播放層第一次失敗的原因（沒有音效裝置時會有值）。
func (b *Bank) Err() error {
	if b == nil {
		return nil
	}
	return b.initErr
}
