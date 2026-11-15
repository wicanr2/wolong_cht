package sound

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type testEffectVoice struct {
	playing, paused, closed bool
	volume                  float64
}

func (v *testEffectVoice) Play()               { v.playing = true }
func (v *testEffectVoice) Pause()              { v.playing, v.paused = false, true }
func (v *testEffectVoice) Close() error        { v.closed = true; return nil }
func (v *testEffectVoice) IsPlaying() bool     { return v.playing }
func (v *testEffectVoice) SetVolume(x float64) { v.volume = x }

func policyTestBank(t *testing.T) (*Bank, *time.Time, *[]*testEffectVoice) {
	t.Helper()
	dir := t.TempDir()
	for _, id := range []int{10, 11, 12} {
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("sfx-%02d.ogg", id)), []byte("fake factory input"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	b := Open(dir)
	now := time.Unix(1000, 0)
	voices := []*testEffectVoice{}
	b.effectNow = func() time.Time { return now }
	b.effectFactory = func([]byte) (effectVoice, error) {
		v := &testEffectVoice{}
		voices = append(voices, v)
		return v, nil
	}
	b.effectDurations = map[int]time.Duration{10: time.Second, 11: time.Second, 12: time.Second}
	return b, &now, &voices
}

func TestEffectPriorityUsesSourceDeadlineAndSingleVoice(t *testing.T) {
	b, now, voices := policyTestBank(t)
	b.PlayEffect(12)
	b.PlayEffect(11)
	if len(*voices) != 2 || !(*voices)[0].closed || !(*voices)[0].paused {
		t.Fatal("較高優先度應關閉舊聲部")
	}
	(*voices)[1].playing = false // OGG狀態不能縮短原始busy。
	b.PlayEffect(12)
	if len(*voices) != 2 {
		t.Fatal("原始busy未到期時，應拒絕較低優先度")
	}
	b.PlayEffect(10)
	*now = now.Add(100 * time.Millisecond)
	b.PlayEffect(10)
	if len(*voices) != 4 || !(*voices)[2].closed {
		t.Fatal("相同效果也應關聲並重新起音")
	}
	*now = now.Add(time.Second - time.Nanosecond)
	b.PlayEffect(12)
	if len(*voices) != 4 {
		t.Fatal("到期前一奈秒仍應busy")
	}
	*now = now.Add(time.Nanosecond)
	b.PlayEffect(12)
	if len(*voices) != 5 || !(*voices)[3].closed {
		t.Fatal("原始busy到期後，OGG仍播放也應接受新效果")
	}
	for _, v := range (*voices)[:4] {
		if !v.closed || v.playing {
			t.Fatal("同時只能保留一個效果聲部")
		}
	}
}

func TestEffectVolumeMuteAndOffCloseOwnedVoice(t *testing.T) {
	b, _, voices := policyTestBank(t)
	b.PlayEffect(10)
	b.SetLevel(2)
	if (*voices)[0].volume != LevelGain(2) {
		t.Fatal("音量應套用現有效果聲部")
	}
	b.SetEnabled(false)
	if !(*voices)[0].closed {
		t.Fatal("OFF應關閉現有效果")
	}
	b.SetEnabled(true)
	b.PlayEffect(12)
	b.SetSilent(true)
	if !(*voices)[1].closed || !b.Available() || !b.Enabled() {
		t.Fatal("驗收靜音應關閉聲部並保留設定狀態")
	}
	b.PlayEffect(10)
	if len(*voices) != 2 {
		t.Fatal("靜音不能建立player")
	}
	b.SetSilent(false)
	b.PlayEffect(11)
	b.PlayEffect(0)
	if !(*voices)[2].closed || b.effect != nil || b.effectCode != 0 {
		t.Fatal("效果0應停止現有效果，不需要sfx-00檔案")
	}
}

func TestEffectMissingAssetAndTimingFallback(t *testing.T) {
	b, _, voices := policyTestBank(t)
	b.PlayEffect(10)
	b.PlayEffect(1) // 缺檔不改現有聲部或優先度。
	if len(*voices) != 1 || (*voices)[0].closed || b.effectCode != 10 {
		t.Fatal("缺少音檔應保留目前聲音")
	}
	b.effectDurations = nil
	b.PlayEffect(12)
	if len(*voices) != 1 {
		t.Fatal("缺原始時序時，以播放狀態fallback")
	}
	(*voices)[0].playing = false
	b.PlayEffect(12)
	if len(*voices) != 2 {
		t.Fatal("fallback播完後應接受新效果")
	}
}

func TestEffectBusyStartsAfterPlayerCreation(t *testing.T) {
	b, now, voices := policyTestBank(t)
	factory := b.effectFactory
	b.effectFactory = func(raw []byte) (effectVoice, error) {
		*now = now.Add(10 * time.Second)
		return factory(raw)
	}
	b.PlayEffect(10)
	b.PlayEffect(12)
	if len(*voices) != 1 || b.effectCode != 10 {
		t.Fatal("首次解碼／裝置準備的耗時不能吃掉聲音開始後的busy")
	}
}

func TestEffectTimingConfigurationRejectsBrokenChain(t *testing.T) {
	b, _, _ := policyTestBank(t)
	data := make([]byte, 13*16)
	data[10*16+14] = 6
	if err := b.SetEffectTimings(data); err != nil {
		t.Fatal(err)
	}
	got := b.effectDurations[10]
	if got <= 0 || got >= time.Second {
		t.Fatalf("原始六拍忙碌時間不合理：%s", got)
	}
	data[10*16+13] = 10
	if err := b.SetEffectTimings(data); err == nil || b.effectDurations[10] != got {
		t.Fatal("壞鏈應回錯誤，不能蓋掉有效時序")
	}
}
