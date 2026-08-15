package audio

import (
	"encoding/binary"
	"fmt"
)

// `*BGM.DAT` 的重放：把事件解成 OPL3 暫存器寫入。
//
// 這一份是 `YNSOUND.COM` 那顆 INT 8 ISR 的重寫（`docs/re/56` §2、§5）。
// **刻意照原版的順序做**——先 key-off 再設音高、控制事件不佔時間、
// 換音色會把「已 key-off」旗標立起來——因為那些順序聽得出來。
//
// 三張查表在 TSR 裡不在資料裡（`docs/re/56` §3），所以 `Tables` 一定要
// 從玩家自己的 `YNSOUND.COM` 讀進來。**不准把表寫死在程式碼裡**：
// 那是原版資料（`CLAUDE.md` §9），而且兩版的 TSR 不見得一樣。

const (
	// ISR 的頻率：PIT divisor 256（`docs/re/57` §5）。
	isrHz = 1193182.0 / 256.0

	trackCount = 6
	// 未使用的聲軌一律指到這個 stub。
	emptyTrack = 0x22

	tblLength = 0x0AB0 // COM offset
	tblB0     = 0x0AD0
	tblA0     = 0x0B50
	comBase   = 0x0100

	instrumentSize = 32
	effectSize     = 16
)

// Tables 是 TSR 裡的三張查表。
type Tables struct {
	Length [32]byte // 長度索引 → tick 數
	// B0 用**整個**音符低 byte 當索引，所以是 128 筆（八度 0–7），
	// 不是 32 筆。表的邊界是 `0x0AD0` 到 `0x0B50`，中間沒有別的東西。
	B0 [128]byte // 音符低 byte → block ＋ F-Number 高 2 位
	A0 [16]byte  // 音名 → F-Number 低 8 位
}

// LoadTables 從 `YNSOUND.COM` 的原始 bytes 取出三張表。
func LoadTables(com []byte) (*Tables, error) {
	if len(com) < tblA0-comBase+16 {
		return nil, fmt.Errorf("YNSOUND.COM 只有 %d bytes，讀不到查表", len(com))
	}
	t := &Tables{}
	copy(t.Length[:], com[tblLength-comBase:])
	copy(t.B0[:], com[tblB0-comBase:])
	copy(t.A0[:], com[tblA0-comBase:])
	return t, nil
}

// Song 是一個曲塊。
type Song struct {
	Data        []byte
	Instruments int // 音色表在曲塊內的位移
	Tracks      [trackCount]int
}

// Songs 把 `*BGM.DAT` 拆成曲塊。
//
// `BGM.DAT` 前 0x100 bytes 是 11 組 (offset, length)；`OPENBGM.DAT` 這種
// 單曲檔沒有索引，整個檔就是一個曲塊（`docs/re/23` §2）。
func Songs(data []byte) ([]Song, error) {
	if len(data) < 0x40 {
		return nil, fmt.Errorf("檔案只有 %d bytes", len(data))
	}
	var spans [][2]int
	if binary.LittleEndian.Uint32(data) == 0x0100 {
		for i := 0; i < 11; i++ {
			off := int(binary.LittleEndian.Uint32(data[i*8:]))
			ln := int(binary.LittleEndian.Uint32(data[i*8+4:]))
			if off < 0 || ln <= 0 || off+ln > len(data) {
				return nil, fmt.Errorf("第 %d 曲的索引 (0x%X, %d) 超出檔案", i, off, ln)
			}
			spans = append(spans, [2]int{off, ln})
		}
	} else {
		spans = append(spans, [2]int{0, len(data)})
	}
	out := make([]Song, 0, len(spans))
	for _, s := range spans {
		b := data[s[0] : s[0]+s[1]]
		if len(b) < 0x1C {
			return nil, fmt.Errorf("曲塊只有 %d bytes", len(b))
		}
		song := Song{Data: b, Instruments: int(binary.LittleEndian.Uint16(b[2:]))}
		for i := range song.Tracks {
			song.Tracks[i] = int(binary.LittleEndian.Uint16(b[0x10+i*2:]))
		}
		out = append(out, song)
	}
	return out, nil
}

// Instrument 取第 n 號音色的 32 bytes。超出範圍回傳 nil。
func (s *Song) Instrument(n int) []byte {
	off := s.Instruments + n*instrumentSize
	if n < 0 || off < 0 || off+instrumentSize > len(s.Data) {
		return nil
	}
	return s.Data[off : off+instrumentSize]
}

// track 是一個聲軌的執行期狀態，欄位對應 TSR 裡的那幾塊記憶體。
type track struct {
	ptr       int   // cs:09F6 — 目前的事件指標
	countdown int   // cs:0A02 — 還剩幾個 tick
	volume    uint8 // cs:0A56
	fadeBias  int8  // cs:0A38 — 漸變累積的偏移
	fadeStep  int8  // cs:09B2
	fadeLeft  uint8 // cs:09AC — 還要漸幾次
	fadePeriod uint8 // cs:09A0 — 每次隔幾個 tick
	fadeTimer uint8  // cs:09A6
	instrument uint8 // cs:0A5C
	pitch     uint16 // cs:0A0E — 待寫的 A0/B0 值
	live      bool

	markJump   int // cs:09B8
	markLoop   int // cs:09C4
	loopCount  uint8
	markSub    int // cs:09D0
	subReturn  int // cs:09DC
}

// Player 重現 TSR 的播放引擎。
type Player struct {
	song   *Song
	tbl    *Tables
	chip   *OPL3
	tracks [trackCount]track

	tempoDiv int // cs:0B68 — ISR 的分頻值
	isrLeft  int
	Flag     uint8 // cs:0999 — 給遊戲讀的同步旗標

	// Done 在所有聲軌都跑到資料尾端時變 true。原版靠控制事件無限循環，
	// 這裡只用來擋住「資料讀爆」的情形，不是正常的結束條件。
	Done bool
}

// NewPlayer 準備一首曲子。晶片會照原版的初始化順序設好。
func NewPlayer(song *Song, tbl *Tables, chip *OPL3) *Player {
	p := &Player{song: song, tbl: tbl, chip: chip, tempoDiv: 1, isrLeft: 1}
	// 原版初始化（`docs/re/57` §1）：先開 NEW 再開六對 4-operator。
	chip.Write(1, 0x05, 0x01)
	chip.Write(1, 0x04, 0x3F)
	for i := range p.tracks {
		t := &p.tracks[i]
		t.ptr = song.Tracks[i]
		t.countdown = 1
		t.live = song.Tracks[i] != emptyTrack
		p.silence(i)
	}
	return p
}

// regs 回傳這一軌四個 operator 的暫存器位移與所屬 bank。
//
// 聲軌 t 佔通道 (t%3) 與 (t%3)+3，bank ＝ t/3（`docs/re/57` §2）。
func (p *Player) regs(t int) (bank int, ch uint8, ops [4]uint8) {
	bank = t / 3
	ch = uint8(t % 3)
	ops = [4]uint8{ch, ch + 3, ch + 8, ch + 11}
	return
}

// silence 把一軌的四個 operator 壓到最小音量並關掉包絡，等同原版的
// 「歸零一軌」（`0x1035D`）。
func (p *Player) silence(t int) {
	bank, ch, ops := p.regs(t)
	p.keyOff(t)
	for _, o := range ops {
		p.chip.Write(bank, 0x40+o, 0x3F)
		p.chip.Write(bank, 0x80+o, 0xFF)
	}
	p.chip.Write(bank, 0xC0+ch, 0x30)
	p.chip.Write(bank, 0xC0+ch+3, 0x30)
}

func (p *Player) keyOff(t int) {
	bank, ch, _ := p.regs(t)
	p.chip.Write(bank, 0xA0+ch, 0)
	p.chip.Write(bank, 0xB0+ch, 0)
}

func (p *Player) keyOn(t int) {
	bank, ch, _ := p.regs(t)
	v := p.tracks[t].pitch
	p.chip.Write(bank, 0xA0+ch, uint8(v))
	p.chip.Write(bank, 0xB0+ch, uint8(v>>8))
}

// setPitch 只算出待寫值，不寫晶片——與原版的 `0x1069E` 一致。
func (p *Player) setPitch(t int, note uint8) {
	b0 := p.tbl.B0[note&0x7F] | 0x20 // key-on 位元
	a0 := p.tbl.A0[note&0x0F]
	p.tracks[t].pitch = uint16(b0)<<8 | uint16(a0)
}

// carrierMask 是四個 operator 裡哪幾個是載波，由兩個通道的 CNT 位元決定。
// 這張表在 TSR 裡是 `cs:0A4Eh`（`docs/re/57` §1）。
var carrierMask = [4]uint8{0x08, 0x0A, 0x09, 0x0D}

func (p *Player) carriers(inst []byte) uint8 {
	idx := (inst[0x14]&1)<<1 | (inst[0x15] & 1)
	return carrierMask[idx]
}

// program 把整組 operator 參數寫進晶片（原版的 `0x1053E`）。
func (p *Player) program(t int) {
	inst := p.song.Instrument(int(p.tracks[t].instrument))
	if inst == nil {
		return
	}
	bank, ch, ops := p.regs(t)
	p.keyOff(t)
	mask := p.carriers(inst)
	for i, o := range ops {
		if mask&(1<<uint(i)) != 0 {
			p.chip.Write(bank, 0x40+o, 0x3F) // 先壓掉載波，免得換音色時爆音
		}
	}
	p.chip.Write(bank, 0xC0+ch, inst[0x14])
	p.chip.Write(bank, 0xC0+ch+3, inst[0x15])
	for i, o := range ops {
		p.chip.Write(bank, 0x20+o, inst[0x00+i])
		p.chip.Write(bank, 0x60+o, inst[0x08+i])
		p.chip.Write(bank, 0x80+o, inst[0x0C+i])
		p.chip.Write(bank, 0xE0+o, inst[0x10+i])
	}
	p.applyVolume(t)
}

// applyVolume 把 0–15 的音量換算成載波的 TL（原版的 `0x1060F`）。
//
// **只有載波改**，調變器直接寫音色原值——改調變器會變音色不是變音量。
func (p *Player) applyVolume(t int) {
	inst := p.song.Instrument(int(p.tracks[t].instrument))
	if inst == nil {
		return
	}
	bank, _, ops := p.regs(t)
	mask := p.carriers(inst)
	atten := int(15-p.tracks[t].volume) * 4
	if atten > 0x3F {
		atten = 0x3F
	}
	for i, o := range ops {
		v := inst[0x04+i]
		if mask&(1<<uint(i)) != 0 {
			lvl := int(v&0x3F) + atten
			if lvl > 0x3F {
				lvl = 0x3F
			}
			v = v&0xC0 | uint8(lvl)
		}
		p.chip.Write(bank, 0x40+o, v)
	}
}

func clampVolume(v int) uint8 {
	if v < 0 {
		return 0
	}
	if v > 15 {
		return 15
	}
	return uint8(v)
}

// ISR 走一步。回傳 true 表示這一步真的推進了音樂（不是被分頻擋掉）。
func (p *Player) isr() bool {
	p.isrLeft--
	if p.isrLeft > 0 {
		return false
	}
	p.isrLeft = p.tempoDiv
	p.tick()
	return true
}

// tick 是音樂的一步：每軌各自倒數，歸零才處理事件。
func (p *Player) tick() {
	for i := 0; i < trackCount; i++ {
		p.advanceFade(i)
		t := &p.tracks[i]
		if !t.live {
			continue
		}
		t.countdown--
		if t.countdown > 0 {
			continue
		}
		p.step(i)
	}
}

// step 處理一軌的下一個事件。控制事件不佔時間，所以會一直讀到音符為止。
func (p *Player) step(i int) {
	t := &p.tracks[i]
	data := p.song.Data
	// 控制事件可以連續出現，但不該無限——遇到相互跳轉的壞資料要停得下來。
	for guard := 0; guard < 4096; guard++ {
		if t.ptr < 0 || t.ptr+2 > len(data) {
			t.live = false
			p.checkDone()
			return
		}
		lo, hi := data[t.ptr], data[t.ptr+1]
		if lo >= 0x80 {
			t.ptr += 2
			p.control(i, lo, hi)
			continue
		}
		keyedOff := false
		if hi&0x80 == 0 { // bit 7 ＝ 連音：不先 key-off
			p.keyOff(i)
			keyedOff = true
		}
		t.countdown = int(p.tbl.Length[hi&0x7F])
		if lo&0x0F == 0 {
			p.keyOff(i) // 休止
		} else {
			p.setPitch(i, lo)
			if keyedOff {
				p.keyOn(i)
			}
		}
		t.ptr += 2
		return
	}
	t.live = false
	p.checkDone()
}

func (p *Player) checkDone() {
	for i := range p.tracks {
		if p.tracks[i].live {
			return
		}
	}
	p.Done = true
}

// control 是八個 handler 的分派（`docs/re/56` §5）。
func (p *Player) control(i int, lo, hi uint8) {
	t := &p.tracks[i]
	switch (lo >> 4) & 7 {
	case 0: // 音量
		t.volume = clampVolume(int(hi) + int(t.fadeBias))
		p.applyVolume(i)
	case 1: // 音量漸變
		t.fadeStep = 1
		if lo&1 != 0 {
			t.fadeStep = -1
		}
		t.fadeLeft = hi >> 4
		t.fadePeriod = (hi & 0x0F) * 4
		t.fadeTimer = t.fadePeriod
	case 2: // 音色
		t.instrument = hi
		p.program(i)
	case 3: // 速度
		n := int(0xFF-int(hi)) * 11 / 8
		if n < 1 {
			n = 1
		}
		p.tempoDiv = n
	case 4: // 跳轉
		switch lo {
		case 0xC1:
			if t.loopCount > 0 {
				t.loopCount--
			}
			if t.loopCount != 0 {
				t.ptr = t.markLoop
			}
		case 0xC2:
			t.subReturn = t.ptr
			t.ptr = t.markSub
		case 0xC3:
			if t.subReturn != 0 {
				t.ptr = t.subReturn
			}
		default:
			t.ptr = t.markJump
		}
	case 5: // 設記號
		switch lo {
		case 0xD1:
			t.loopCount = hi
			t.markLoop = t.ptr
		case 0xD2:
			t.subReturn = 0
			t.markSub = t.ptr
		default:
			t.markJump = t.ptr
		}
	case 6: // 原版只有一個 retn
	case 7: // 給遊戲讀的同步旗標
		p.Flag = hi
	}
}

// advanceFade 是原版 `0x1075E` 每 tick 做的音量漸變。
func (p *Player) advanceFade(i int) {
	t := &p.tracks[i]
	if t.fadePeriod == 0 || t.fadeStep == 0 {
		return
	}
	if t.fadeTimer > 0 {
		t.fadeTimer--
		return
	}
	if t.fadeLeft == 0 {
		t.fadeStep = 0
		t.fadeBias = 0
		return
	}
	t.fadeLeft--
	t.fadeTimer = t.fadePeriod
	t.fadeBias += t.fadeStep
	t.volume = clampVolume(int(t.volume) + int(t.fadeStep))
	p.applyVolume(i)
}
