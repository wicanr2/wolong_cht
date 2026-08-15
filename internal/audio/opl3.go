// Package audio 是原版音訊的離線重放：把 `*BGM.DAT` 的事件解成 OPL3
// 暫存器寫入序列，再用一顆純 Go 的 OPL3 把序列渲染成 PCM。
//
// 規格見 `docs/spec/29-audio.md`，原版側的出處在 `docs/re/56`（事件）
// 與 `docs/re/57`（OPL3 暫存器、音色、音量、速度）。
package audio

import "math"

// YMF262（OPL3）音源核心。
//
// ## 為什麼是 OPL3 而不是 OPL2
//
// `YNSOUND.COM` 初始化時對「基底埠 +2」寫了 reg `0x105 = 1`（NEW）與
// `0x104 = 0x3F`（六對 4-operator 通道全開），這兩個暫存器 OPL2 沒有
// （`docs/re/57` §1）。六個聲軌各佔一組 4-op 通道，音效走剩下的
// 三個 2-op 通道。
//
// ## 哪些是證據、哪些是模型
//
// **證據**（`docs/re/57`）：暫存器配置、音色欄位版面、音量換算、
// 4-op 通道配對與載波遮罩。那些都由呼叫端（`bgm.go`）餵進來。
//
// **公開規格照做**（不是某份實作的抄本）：對數正弦表與指數表用標準式子
// 產生；包絡走硬體那套「全域計數器 ＋ 0/1 增量樣式表」，不是每個取樣點
// 加固定量——這是 rate 每 +1 快約 1.19 倍的來源；KSL、KSR、feedback、
// 八種波形、四種 4-op 連接型態都照規格。
//
// **仍是模型的部分**（不宣稱逐週期 parity）：顫音／抖音用規格的深度與
// 3.7 Hz／6.1 Hz，波形取三角／正弦近似；相位以浮點算再轉定點，硬體是
// 純整數。**日後對不上原版錄音時，先看是不是落在這半邊。**
//
// 行為對照 [ymfm](https://github.com/aaronsgiles/ymfm)（BSD-3-Clause）
// 與 OPL3 的公開規格；本檔是照規格重寫，沒有搬運上游程式碼。
//
// **rhythm 模式沒有實作**——原版一次都沒開（`docs/re/57` §2：六個 4-op
// ＋ 三個 2-op 已經把資源用完，`0xBD` 的 rhythm 位元從沒被寫過）。
// 要用到的時候再補，不要為了「完整」先寫一段沒有測試的程式碼。
const (
	// OPL3 的原生取樣率。要別的取樣率請在外面重取樣，不要改這裡
	// ——包絡與相位都以晶片時脈為準。
	NativeRate = 49716.0

	opl3Operators = 36
	opl3Channels  = 18

	// 包絡的單位是 0.1875 dB，滿刻度 511（≈ 96 dB）。
	opl3EGMax = 511
)

// 對數正弦表與指數表。OPL 的乘法是在對數域做加法，兩張表就是進出對數域的橋。
var (
	opl3LogSin [256]uint16
	opl3Exp    [256]uint16
)

func init() {
	for i := 0; i < 256; i++ {
		// −log2(sin((i+0.5) × π/512)) × 256
		opl3LogSin[i] = uint16(math.Round(
			-math.Log2(math.Sin((float64(i)+0.5)*math.Pi/512.0)) * 256))
		// (2^(i/256) − 1) × 1024
		opl3Exp[i] = uint16(math.Round((math.Pow(2, float64(i)/256.0) - 1) * 1024))
	}
}

type envPhase uint8

const (
	envOff envPhase = iota
	envAttack
	envDecay
	envSustain
	envRelease
)

type operator struct {
	am, vib, egt, ksr bool
	multiple          uint8
	ksl               uint8
	totalLevel        uint8 // 0..63，每級 0.75 dB
	attack            uint8
	decay             uint8
	sustain           uint8
	release           uint8
	waveform          uint8

	phase     uint32
	env       envPhase
	level     int32 // 目前衰減量，0（最大聲）..511（無聲）
	out, prev int32 // feedback 用的前兩個輸出
}

type channel struct {
	fnum     uint16
	block    uint8
	keyOn    bool
	feedback uint8
	// connect 是 `0xC0` 的 bit 0。2-op 時 true ＝ 並聯；
	// 4-op 時它與配對通道的 connect 一起決定四種連接型態（§algorithm）。
	connect     bool
	left, right bool
}

// opl3SlotToOperator：暫存器位移 → 該 bank 內的 operator 編號。
// `0x20`／`0x40`／`0x60`／`0x80`／`0xE0` 五族共用這張表。
// **位移不是連號的**——每組 6 個之後跳 2，這是 OPL 的固定配線。
var opl3SlotToOperator = map[uint8]int{
	0x00: 0, 0x01: 1, 0x02: 2, 0x03: 3, 0x04: 4, 0x05: 5,
	0x08: 6, 0x09: 7, 0x0A: 8, 0x0B: 9, 0x0C: 10, 0x0D: 11,
	0x10: 12, 0x11: 13, 0x12: 14, 0x13: 15, 0x14: 16, 0x15: 17,
}

// 通道 → 兩個 operator（bank 內編號）。
var opl3OperatorIndex = [9][2]int{
	{0, 3}, {1, 4}, {2, 5},
	{6, 9}, {7, 10}, {8, 11},
	{12, 15}, {13, 16}, {14, 17},
}

var opl3KSLRom = [16]int{0, 32, 40, 44, 48, 50, 52, 54, 56, 57, 58, 59, 60, 61, 62, 63}

var opl3KSLShift = [4]uint{8, 1, 2, 0}

// OPL3 是一顆 YMF262。輸出是立體聲、49,716 Hz。
type OPL3 struct {
	regs [2][256]uint8
	ops  [opl3Operators]operator
	chs  [opl3Channels]channel

	// newMode 是 reg 0x105 bit 0。關著的時候只有前四種波形、沒有 4-op、
	// 第二 bank 不出聲——原版一開機就打開它。
	newMode bool
	// fourOp 是 reg 0x104 的低六位：bit 0–2 ＝ bank 0 的三對，
	// bit 3–5 ＝ bank 1 的三對。
	fourOp uint8

	tremoloDeep bool
	vibratoDeep bool

	sampleRate float64
	samples    uint64
	// egCounter 是包絡的全域計數器。所有 operator 共用它，
	// 所以同一個 rate 的兩個音會同步走——這是聽得出來的特徵。
	egCounter uint32
}

// NewOPL3 建一顆晶片。sampleRate 給 0 就用原生的 49,716 Hz。
func NewOPL3(sampleRate float64) *OPL3 {
	if sampleRate <= 0 {
		sampleRate = NativeRate
	}
	c := &OPL3{sampleRate: sampleRate}
	// 開機預設兩個輸出都關。原版的音色資料自己會寫 0xC0 的 bit 4/5。
	return c
}

// Write 寫一個暫存器。bank 0 是基底埠、bank 1 是「基底埠 +2」那一組。
func (c *OPL3) Write(bank int, addr, value uint8) {
	if bank < 0 || bank > 1 {
		return
	}
	c.regs[bank][addr] = value
	base := bank * 9 // 通道編號的 bank 偏移
	opBase := bank * 18

	switch {
	case bank == 1 && addr == 0x04:
		c.fourOp = value & 0x3F
	case bank == 1 && addr == 0x05:
		c.newMode = value&0x01 != 0
	case addr == 0xBD:
		c.tremoloDeep = value&0x80 != 0
		c.vibratoDeep = value&0x40 != 0
	case addr >= 0x20 && addr <= 0x35:
		if i, ok := opl3SlotToOperator[addr-0x20]; ok {
			op := &c.ops[opBase+i]
			op.am = value&0x80 != 0
			op.vib = value&0x40 != 0
			op.egt = value&0x20 != 0
			op.ksr = value&0x10 != 0
			op.multiple = value & 0x0F
		}
	case addr >= 0x40 && addr <= 0x55:
		if i, ok := opl3SlotToOperator[addr-0x40]; ok {
			c.ops[opBase+i].ksl = value >> 6
			c.ops[opBase+i].totalLevel = value & 0x3F
		}
	case addr >= 0x60 && addr <= 0x75:
		if i, ok := opl3SlotToOperator[addr-0x60]; ok {
			c.ops[opBase+i].attack = value >> 4
			c.ops[opBase+i].decay = value & 0x0F
		}
	case addr >= 0x80 && addr <= 0x95:
		if i, ok := opl3SlotToOperator[addr-0x80]; ok {
			c.ops[opBase+i].sustain = value >> 4
			c.ops[opBase+i].release = value & 0x0F
		}
	case addr >= 0xA0 && addr <= 0xA8:
		ch := &c.chs[base+int(addr-0xA0)]
		ch.fnum = (ch.fnum & 0x300) | uint16(value)
	case addr >= 0xB0 && addr <= 0xB8:
		n := base + int(addr-0xB0)
		ch := &c.chs[n]
		ch.fnum = (ch.fnum & 0xFF) | (uint16(value&0x03) << 8)
		ch.block = (value >> 2) & 0x07
		on := value&0x20 != 0
		if on != ch.keyOn {
			ch.keyOn = on
			c.keyChange(n, on)
		}
	case addr >= 0xC0 && addr <= 0xC8:
		ch := &c.chs[base+int(addr-0xC0)]
		ch.feedback = (value >> 1) & 0x07
		ch.connect = value&0x01 != 0
		ch.left = value&0x10 != 0
		ch.right = value&0x20 != 0
	case addr >= 0xE0 && addr <= 0xF5:
		if i, ok := opl3SlotToOperator[addr-0xE0]; ok {
			w := value & 0x07
			if !c.newMode {
				w &= 0x03
			}
			c.ops[opBase+i].waveform = w
		}
	}
}

// fourOpPair 回傳 n 是不是 4-op 組的**前一個**通道，以及配對的後一個通道。
//
// ⚠ 只有 bank 內的 0/1/2 才可能是前一個。判斷寫在一處，
// 因為 key-on 與混音兩邊都要用同一套判準，分開寫遲早會漂。
func (c *OPL3) fourOpPair(n int) (partner int, ok bool) {
	if !c.newMode {
		return 0, false
	}
	bank, idx := n/9, n%9
	if idx > 2 {
		return 0, false
	}
	if c.fourOp&(1<<uint(bank*3+idx)) == 0 {
		return 0, false
	}
	return bank*9 + idx + 3, true
}

// isFourOpSecond 回傳 n 是不是某一組 4-op 的後半（後半不獨立出聲）。
func (c *OPL3) isFourOpSecond(n int) bool {
	bank, idx := n/9, n%9
	if idx < 3 || idx > 5 {
		return false
	}
	return c.newMode && c.fourOp&(1<<uint(bank*3+idx-3)) != 0
}

func (c *OPL3) keyChange(n int, on bool) {
	ops := c.channelOperators(n)
	for _, i := range ops {
		op := &c.ops[i]
		if on {
			op.env = envAttack
			op.phase = 0
			// 觸發時不把 level 歸零：硬體的 attack 是從目前衰減量往上追，
			// 連續彈同一個音才不會有卡頓。
			if op.level >= opl3EGMax {
				op.level = opl3EGMax
			}
		} else if op.env != envOff {
			op.env = envRelease
		}
	}
}

// channelOperators 回傳這個通道實際發聲用到的 operator。
// 4-op 的前半回傳四個（含配對通道的兩個），後半在 4-op 時回傳空的。
func (c *OPL3) channelOperators(n int) []int {
	bank, idx := n/9, n%9
	base := bank * 18
	pair := opl3OperatorIndex[idx]
	if partner, ok := c.fourOpPair(n); ok {
		p := opl3OperatorIndex[partner%9]
		return []int{base + pair[0], base + pair[1], base + p[0], base + p[1]}
	}
	if c.isFourOpSecond(n) {
		return nil
	}
	return []int{base + pair[0], base + pair[1]}
}

// 每個 operator 的實際包絡速率：rate = 4 × 欄位值 ＋ key-scale。
func rateOf(op *operator, field uint8, ch *channel) int {
	if field == 0 {
		return 0
	}
	ks := int(ch.block) >> 1
	if op.ksr {
		ks = (int(ch.block) << 1) | int((ch.fnum>>9)&1)
	}
	rate := 4*int(field) + ks
	if rate > 63 {
		rate = 63
	}
	return rate
}

// opl3EGIncrement 是包絡的增量樣式表。
//
// OPL 的包絡不是「每個取樣點固定加多少」，而是**一個全域計數器配一張
// 0/1 樣式表**：速率的低兩位選一列，計數器的高位選一欄，選到 1 才走一步。
// 這樣才做得出「rate 每 +1 大約快 1.19 倍」的階梯，而不是只有 2 的冪次。
var opl3EGIncrement = [4][8]uint8{
	{0, 1, 0, 1, 0, 1, 0, 1},
	{0, 1, 0, 1, 1, 1, 0, 1},
	{0, 1, 1, 1, 0, 1, 1, 1},
	{0, 1, 1, 1, 1, 1, 1, 1},
}

// egStep 回傳這個取樣點該走幾個包絡單位（0 表示不動）。
//
// `shift = 13 − rate/4` 是分頻量：rate 每大 4，計數器就少除一次 2，
// 時間剛好減半；rate ≥ 52 之後 shift 觸底，改成把增量往左移，
// 於是最快的幾檔仍然能繼續加速。
func (c *OPL3) egStep(rate int) int32 {
	if rate < 4 {
		return 0
	}
	shift := 13 - (rate >> 2)
	sel := rate & 3
	if shift > 0 {
		mask := uint32(1)<<uint(shift) - 1
		if c.egCounter&mask != 0 {
			return 0
		}
		return int32(opl3EGIncrement[sel][(c.egCounter>>uint(shift))&7])
	}
	return int32(opl3EGIncrement[sel][c.egCounter&7]) << uint(-shift)
}

func (c *OPL3) stepEnvelope(op *operator, ch *channel) {
	switch op.env {
	case envAttack:
		r := rateOf(op, op.attack, ch)
		if r >= 60 {
			// 最快的幾檔在硬體上等於「立刻到頂」。
			op.level = 0
			op.env = envDecay
			return
		}
		if inc := c.egStep(r); inc > 0 {
			// Attack 是往 0 逼近的指數曲線：離目標越近走得越慢。
			op.level += (^op.level * inc) >> 3
		}
		if op.level <= 0 {
			op.level = 0
			op.env = envDecay
		}
	case envDecay:
		op.level += c.egStep(rateOf(op, op.decay, ch))
		if sustain := int32(op.sustain) * 32; op.level >= sustain {
			op.level = sustain
			op.env = envSustain
		}
	case envSustain:
		if !op.egt {
			// EGT = 0（percussive）時 sustain 階段仍以 release 速率繼續衰減。
			op.level += c.egStep(rateOf(op, op.release, ch))
			if op.level >= opl3EGMax {
				op.level = opl3EGMax
				op.env = envOff
			}
		}
	case envRelease:
		op.level += c.egStep(rateOf(op, op.release, ch))
		if op.level >= opl3EGMax {
			op.level = opl3EGMax
			op.env = envOff
		}
	}
	if op.level > opl3EGMax {
		op.level = opl3EGMax
	}
}

// kslAttenuation 依 block／fnum 算出高音衰減，單位與包絡相同。
func kslAttenuation(ksl uint8, block uint8, fnum uint16) float64 {
	if ksl == 0 {
		return 0
	}
	v := (opl3KSLRom[fnum>>6] << 2) - int(8-block)<<5
	if v < 0 {
		v = 0
	}
	return float64(v >> opl3KSLShift[ksl])
}

// expOut 把「對數域的總衰減」換回線性輸出。
func expOut(level int) int32 {
	if level < 0 {
		level = 0
	}
	if level > 0x1FFF {
		level = 0x1FFF
	}
	return int32(((int(opl3Exp[(level&0xFF)^0xFF]) | 0x400) << 1) >> uint(level>>8))
}

// waveOut 依波形選擇把相位轉成「對數域振幅 ＋ 正負號」。
// 波形 4–7 只有 NEW 位元開著才選得到（Write 已經擋掉）。
func waveOut(phase uint32, waveform uint8) (logAmp int, sign int) {
	p := phase & 0x3FF
	neg := false
	switch waveform {
	case 0:
		neg = p&0x200 != 0
	case 1:
		if p&0x200 != 0 {
			return -1, 0 // 靜音半週
		}
	case 2:
		p &= 0x1FF
	case 3:
		if p&0x100 != 0 {
			return -1, 0
		}
		p &= 0x0FF
	case 4, 5:
		// 前半週壓成兩倍頻，後半週靜音。
		if p&0x200 != 0 {
			return -1, 0
		}
		p = (p << 1) & 0x3FF
		if waveform == 4 {
			neg = p&0x200 != 0
		}
	case 6:
		// 方波：振幅固定，只有正負。
		return 0, sgn(p&0x200 == 0)
	case 7:
		// 對數鋸齒。
		if p&0x200 != 0 {
			return int((0x200 - (p & 0x1FF)) << 3), -1
		}
		return int(p << 3), 1
	}
	idx := p & 0xFF
	if p&0x100 != 0 {
		idx ^= 0xFF
	}
	return int(opl3LogSin[idx]), sgn(!neg)
}

func sgn(positive bool) int {
	if positive {
		return 1
	}
	return -1
}

func (c *OPL3) tremolo() float64 {
	depth := 1.0 // dB
	if c.tremoloDeep {
		depth = 4.8
	}
	ph := math.Mod(float64(c.samples)*3.7/c.sampleRate, 1)
	tri := ph * 2
	if tri > 1 {
		tri = 2 - tri
	}
	return tri * depth / 0.1875
}

func (c *OPL3) vibratoCents(op *operator) float64 {
	if !op.vib {
		return 0
	}
	depth := 7.0
	if c.vibratoDeep {
		depth = 14.0
	}
	return depth * math.Sin(2*math.Pi*6.1*float64(c.samples)/c.sampleRate)
}

// multipleFactor 是 MULT 欄位的倍率。0 是 0.5，10/12/14 有重複值——
// **這張表不是 `max(1, n)`**，照抄公式會在四個值上錯掉。
var multipleFactor = [16]float64{0.5, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 12, 12, 15, 15}

// phaseIncrement 是每個取樣點的相位增量（20 位元定點，高 10 位是波表索引）。
func (c *OPL3) phaseIncrement(op *operator, ch *channel) uint32 {
	f := float64(ch.fnum) * math.Exp2(float64(ch.block)) * math.Exp2(c.vibratoCents(op)/1200)
	inc := f * multipleFactor[op.multiple] * 1024 / (1 << 20) * (NativeRate / c.sampleRate)
	return uint32(inc * 1024) // 低 10 位是小數
}

func (c *OPL3) operatorOutput(op *operator, ch *channel, modulation int32) int32 {
	logAmp, sign := waveOut(uint32(int32(op.phase>>10)+modulation), op.waveform)
	if sign == 0 {
		return 0
	}
	level := float64(op.level) +
		float64(op.totalLevel)*4 +
		kslAttenuation(op.ksl, ch.block, ch.fnum)
	if op.am {
		level += c.tremolo()
	}
	return int32(sign) * expOut(logAmp+int(level)*8)
}

// Generate 產生一個取樣點（左、右），範圍約 ±8 千。
func (c *OPL3) Generate() (float64, float64) {
	var left, right float64
	for n := 0; n < opl3Channels; n++ {
		if c.isFourOpSecond(n) {
			continue // 由前半一起處理
		}
		ch := &c.chs[n]
		var mix float64
		if partner, ok := c.fourOpPair(n); ok {
			mix = c.generateFourOp(n, partner)
		} else {
			mix = c.generateTwoOp(n)
		}
		if mix == 0 {
			continue
		}
		// OPL3 沒開輸出位元就真的不出聲。原版的音色資料一律兩邊都開。
		if ch.left {
			left += mix
		}
		if ch.right {
			right += mix
		}
	}
	c.samples++
	c.egCounter++
	return left, right
}

func (c *OPL3) generateTwoOp(n int) float64 {
	ch := &c.chs[n]
	ops := c.channelOperators(n)
	if len(ops) != 2 {
		return 0
	}
	mod, car := &c.ops[ops[0]], &c.ops[ops[1]]
	c.stepEnvelope(mod, ch)
	c.stepEnvelope(car, ch)
	if mod.env == envOff && car.env == envOff {
		return 0
	}
	modOut := c.operatorOutput(mod, ch, c.feedbackOf(ch, mod))
	mod.prev, mod.out = mod.out, modOut
	var mix float64
	if ch.connect {
		mix = float64(modOut) + float64(c.operatorOutput(car, ch, 0))
	} else {
		mix = float64(c.operatorOutput(car, ch, modOut>>1))
	}
	mod.phase += c.phaseIncrement(mod, ch)
	car.phase += c.phaseIncrement(car, ch)
	return mix
}

// generateFourOp 走四種 4-op 連接型態。
//
// 型態由「前一個通道的 CNT ＋ 後一個通道的 CNT」決定，載波集合是：
//
//	0 0  op1→op2→op3→op4          載波 op4
//	0 1  op1→op2 ＋ op3→op4        載波 op2、op4
//	1 0  op1 ＋ op2→op3→op4        載波 op1、op4
//	1 1  op1 ＋ op2→op3 ＋ op4      載波 op1、op3、op4
//
// ⭐ 這四組載波正好等於 `YNSOUND.COM` 那張遮罩表 `08 0A 09 0D`
// （`docs/re/57` §1）——原版的音量常式只對載波加衰減，
// 所以這個拓樸不是「照規格猜」，是被原版程式碼交叉驗證過的。
func (c *OPL3) generateFourOp(n, partner int) float64 {
	ch, ch2 := &c.chs[n], &c.chs[partner]
	ops := c.channelOperators(n)
	if len(ops) != 4 {
		return 0
	}
	op1, op2 := &c.ops[ops[0]], &c.ops[ops[1]]
	op3, op4 := &c.ops[ops[2]], &c.ops[ops[3]]
	for _, o := range []*operator{op1, op2, op3, op4} {
		c.stepEnvelope(o, ch)
	}
	if op1.env == envOff && op2.env == envOff && op3.env == envOff && op4.env == envOff {
		return 0
	}

	// 頻率一律用前一個通道的（4-op 模式下後一個的 A0/B0 不作用），
	// 所以四個 operator 的相位增量都拿 ch 算。
	out1 := c.operatorOutput(op1, ch, c.feedbackOf(ch, op1))
	op1.prev, op1.out = op1.out, out1

	var mix float64
	switch {
	case !ch.connect && !ch2.connect:
		out2 := c.operatorOutput(op2, ch, out1>>1)
		out3 := c.operatorOutput(op3, ch, out2>>1)
		mix = float64(c.operatorOutput(op4, ch, out3>>1))
	case !ch.connect && ch2.connect:
		mix = float64(c.operatorOutput(op2, ch, out1>>1))
		out3 := c.operatorOutput(op3, ch, 0)
		mix += float64(c.operatorOutput(op4, ch, out3>>1))
	case ch.connect && !ch2.connect:
		mix = float64(out1)
		out2 := c.operatorOutput(op2, ch, 0)
		out3 := c.operatorOutput(op3, ch, out2>>1)
		mix += float64(c.operatorOutput(op4, ch, out3>>1))
	default:
		mix = float64(out1)
		out2 := c.operatorOutput(op2, ch, 0)
		mix += float64(c.operatorOutput(op3, ch, out2>>1))
		mix += float64(c.operatorOutput(op4, ch, 0))
	}
	for _, o := range []*operator{op1, op2, op3, op4} {
		o.phase += c.phaseIncrement(o, ch)
	}
	return mix
}

// feedbackOf 是 operator 1 的自我回授：取前兩次輸出的平均。
func (c *OPL3) feedbackOf(ch *channel, op *operator) int32 {
	if ch.feedback == 0 {
		return 0
	}
	return ((op.out + op.prev) >> (9 - ch.feedback)) >> 1
}
