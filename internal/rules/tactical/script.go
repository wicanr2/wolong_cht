package tactical

// `BATTLE.DAT` 的腳本直譯器。
//
// 原版每幀執行一個指令（`sub_1A426`，docs/re/11 §3.1）：
//
//	一個指令 ＝ 2 byte：低 5 位是指令碼、高 3 位是參數、第二個 byte 是運算元
//	分支指令（10）後面**還跟一個 word 的跳躍目標**，所以它是 4 byte
//
// 十九個處理常式裡，條件全部寫進同一個暫存器再由分支讀走——
// **這是一台「查詢 → 分支 → 下令 → 等待」的小機器**。
//
// ⚠ 十九個裡有幾個查的是還沒解出來的全域變數（`byte_1D346`、
// `word_1D33C`、`word_1D30A` 那些）。**那幾個在這裡回 0**，
// 並在下面各自標出來——解出來之前不要當成原版行為。

// 指令碼。名稱與作用見 docs/re/11 §3.5。
const (
	opWait     = 0
	opForm     = 1  // 切換陣形
	opLine     = 2  // 陣形線
	opOrder    = 3  // 下命令
	opQD346    = 4  // 未解
	opQD33C    = 5  // 未解
	opQMine    = 6  // 我方六隊命令的最小值
	opQTheirs  = 7  // 敵方六隊命令的最小值
	opQAdv     = 8  // 有利／不利
	opQRand    = 9  // 亂數 mod N
	opBranch   = 10 // 條件分支
	opQA24     = 11 // 未解
	opQA04     = 12 // 未解
	opOrderBy  = 13 // 依兵種下命令
	opQD31C    = 14 // 我方兵數
	opQMin18   = 15 // 門／城壁耐久的最小值
	opMessage  = 16
	opQCmd9    = 17
	opQB03     = 18
	numOpcodes = 19
)

// 分支的五種比較（`sub_1A591` 的 switch）。
const (
	cmpAlways = 0
	cmpEQ     = 1
	cmpNE     = 2
	cmpLT     = 3
	cmpGT     = 4
)

// kindStride 是指令 13 挑隊時乘的倍數。與兵種編碼同源（兵種 × 18）。
const kindStride = 18

// ScriptCodeSize 是一段腳本的長度（`sub_1CBE5` 讀 0x100）。
const ScriptCodeSize = 256

// Script 是一段跑在某一側身上的 AI 腳本。
type Script struct {
	code []byte
	side int

	pc   int // 讀到第幾個 byte
	wait int // 還要等幾幀
	cond int // 條件暫存器（原版 byte_1D315）

	// Halted 為真表示腳本走到底了。原版的腳本是 256 B 的環，
	// 走完就停在那裡；這裡也一樣，不繞回去。
	Halted bool
}

// NewScript 綁一段腳本到某一側。code 是 256 byte 的一段（`Library.Script`）。
func NewScript(code []byte, side int) *Script {
	if len(code) == 0 {
		return &Script{Halted: true}
	}
	return &Script{code: code, side: side}
}

// Step 跑一幀。與原版一樣：**等待中就什麼都不做**。
func (s *Script) Step(b *Battle) {
	if s.Halted {
		return
	}
	if s.wait > 0 {
		s.wait--
		return
	}
	if s.pc+1 >= len(s.code) {
		s.Halted = true
		return
	}
	lo, arg := s.code[s.pc], s.code[s.pc+1]
	op, par := int(lo&0x1F), int((lo&0xE0)>>5)
	s.pc += 2
	if op >= numOpcodes {
		s.Halted = true
		return
	}
	s.exec(b, op, par, int(arg))
}

func (s *Script) exec(b *Battle, op, par, arg int) {
	side := s.side
	switch op {
	case opWait:
		s.wait = arg

	case opForm:
		// 原版是 `word_1D344 ← 運算元 × 96`，而 96 ＝ 48 個兵 × 2 byte，
		// 所以運算元就是陣形編號（docs/re/11 §5.8d）。
		if arg < NumFormations {
			b.Sides[side].Formation = arg
		}

	case opLine:
		// 三個常數 58／36／16，由**運算元**選（不是參數）。
		b.Sides[side].Line = LineFor(side, arg)

	case opOrder:
		// 參數 7 ＝ 全軍，0–5 ＝ 指定一隊。
		squad := par
		if par == 7 {
			squad = -1
		}
		b.Order(side, squad, Command(arg))

	case opOrderBy:
		// 只對某個兵種的隊下命令（`cmp ah, [si+24h]`，參數 × 18）。
		want := Kind(par * kindStride)
		for k := 0; k < Squads; k++ {
			if b.Sides[side].Soldiers[k*PerSquad].Kind == want {
				b.Order(side, k, Command(arg))
			}
		}

	case opQMine:
		s.cond = minCommand(&b.Sides[side])
	case opQTheirs:
		s.cond = minCommand(&b.Sides[1-side])
	case opQAdv:
		s.cond = int(b.Advantage[side])
	case opQRand:
		n := arg
		if n == 0 {
			n = Squads
		}
		s.cond = b.rng.Next() % n
	case opQD31C:
		s.cond = clampByte(b.Sides[side].Remaining())
	case opQMin18:
		// 門與城壁的耐久。**還沒實作**（那張 16 筆的表沒接），
		// 回 0 等於「已經被打爛」，腳本會走比較積極的分支。
		s.cond = 0
	case opQCmd9:
		s.cond = 0
		if b.Sides[side].Soldiers[0].Cmd >= 9 {
			s.cond = 1
		}
	case opQB03:
		s.cond = b.Sides[side].Soldiers[0].HP

	case opQD346, opQD33C, opQA24, opQA04:
		// ⚠ 這四個查的全域變數還沒解（docs/re/11 §3.5）。回 0。
		s.cond = 0

	case opMessage:
		// 訊息 0x1CE + N。戰場的訊息還沒接，先記進 Log。
		b.Log = append(b.Log, "（訊息 "+itoa(0x1CE+arg)+"）")

	case opBranch:
		s.branch(par, arg)
	}
}

// branch 重現 `sub_1A591`：目標是**下一個 word 的高位元組 × 2**，
// 而那個 word 的低位元組是 0。不成立就跳過它。
func (s *Script) branch(par, arg int) {
	take := false
	switch par {
	case cmpAlways:
		take = true
	case cmpEQ:
		take = s.cond == arg
	case cmpNE:
		take = s.cond != arg
	case cmpLT:
		take = s.cond < arg
	case cmpGT:
		take = s.cond > arg
	}
	if s.pc+1 >= len(s.code) || s.code[s.pc] != 0 {
		// 後面那個 word 不是目標（低位元組非 0）——原版也是這樣就不跳。
		return
	}
	if take {
		s.pc = int(s.code[s.pc+1]) * 2
		return
	}
	s.pc += 2 // 跳過目標
}

// minCommand 回傳一側六隊「生效中的命令」的最小值。
// 原版把 ≥ 4 的當成 0（`sub_1A52E`）。
func minCommand(sd *Side) int {
	m := 255
	for k := 0; k < Squads; k++ {
		c := int(sd.Soldiers[k*PerSquad].Cmd)
		if c >= 4 {
			c = 0
		}
		if c < m {
			m = c
		}
	}
	if m == 255 {
		return 0
	}
	return m
}

func clampByte(v int) int {
	if v > 255 {
		return 255
	}
	if v < 0 {
		return 0
	}
	return v
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	var b [8]byte
	i := len(b)
	for v > 0 {
		i--
		b[i] = byte('0' + v%10)
		v /= 10
	}
	return string(b[i:])
}

// SetScript 讓某一側由腳本驅動。傳 nil 就改回手動（玩家操作）。
func (b *Battle) SetScript(side int, s *Script) { b.scripts[side] = s }
