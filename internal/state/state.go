// Package state 把原版的劇本／存檔載進規則層可以直接跑的結構。
//
// 這一層是**純資料 ＋ 世界迴圈**，不含畫面。所以它可以在無頭環境裡
// 大量重複執行，用長期行為去驗證那些從機器碼讀出來的公式
// （見 cmd/wlsim）。
//
// 檔案格式見 docs/formats/08，各欄位的來源見 docs/re/06 與 docs/re/07。
// 原版資產不隨本專案散布，載入時請自備 SINARIO.DAT。
package state

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/wicanr2/wolong_cht/internal/rules/clock"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/general"
)

// 劇本／存檔的佈局常數（docs/formats/08 §0–§1.6）。
const (
	blockSize = 22208 // 0x56C0，一個劇本區塊
	numBlocks = 4

	factionBase, factionSize, numFactions = 0x0080, 64, 22
	cityBase, citySize, numCities         = 0x08C0, 32, 192
	generalBase, generalSize, numGenerals = 0x42C0, 32, 127

	// 稅率與三兵種募兵數。原版載到 cs:0D08h，而區塊前 59 B 對映到
	// cs:0CF0h，所以 0D08h − 0CF0h = 0x18 就是區塊內的偏移。
	taxOffset        = 0x18
	recruitCapOffset = 0x1A
	nextSettings     = 0x20 // 「來月」的同四項，月結時搬到上面兩處
)

// Faction 是一個勢力的完整狀態。
type Faction struct {
	Alive    bool
	Lord     int // 君主的武將編號
	Advisor  int // 軍師的武將編號，0xFF ＝ 無
	Capital  int // 首都的據點編號
	Reserves [economy.NumTroopTypes]int

	Generals int // 武將數，月結時不重算（由武將表決定）
	Cities   int // 據點數，月結時重算
	Funds    int // 有號 24 位
	Expense  int // 本月累計支出，月結後歸零

	// MoraleBase 是新編成軍團的初始士氣（記錄 +0x1D，四個劇本都是 200）。
	//
	// ⚠ 先前這一欄被標成「疑似信賴度」。編成軍團時它被複製進軍團的
	// 士氣欄位，而說明書編成畫面的士氣值正好是 200 ——
	// 所以它是士氣基準值。**信賴度存在哪還沒解**（docs/re/08 §4）。
	MoraleBase int

	// Aggression 是君主的好戰等級（0–15），**不顯示給玩家**。
	// 呂布 15、曹操 14、劉備 4、劉表 1、劉禪 0。
	// 同一位君主在所有劇本裡是同一個值。
	Aggression int

	// InvasionTarget 是這個勢力正在侵攻的對象（勢力記錄 +0x19），
	// 0xFF ＝ 無。**財政撐不住時原版會自動把它設回 0xFF**
	// （docs/re/08 §1）——這是政略 AI 的一條硬性煞車。
	InvasionTarget int

	Corps int // 軍團數（記錄 +0x14）

	Diplomat int // 派駐「這個」勢力的外交官（由別人派來），0xFF ＝ 無
}

// City 是一個據點的完整狀態。
type City struct {
	// Name 是城市名稱（6 byte Big5，原始位元組）。北京、涿郡、武陵…
	Name string

	// Owner 是**執行期**的所屬勢力（記錄 +0x01）。
	// 月結（sub_153C6）與內政官（sub_15715）讀的都是這一欄。
	Owner int

	// OwnerRecorded 是記錄 +0x1A 的所屬勢力。
	//
	// ⚠ 兩欄在 190/192 個據點上一致，但**劇本 3 的武陵與劇本 4 的南昌
	// 不一致**，而且不一致的那兩個裡 +0x1A 才與勢力記錄的據點數（+0x23）
	// 相符，歷史上也才正確（武陵屬荊南劉備、南昌屬東吳孫權）。
	//
	// 因為執行期讀的是 +0x01，原版實際會把武陵判給劉璋、南昌判給劉禪。
	// **這是原版的資料瑕疵**，remake 預設照抄（用 Owner），
	// 要改成作者意圖就改用這一欄。見 docs/formats/08 §1.6。
	OwnerRecorded int

	X, Y          int
	Production    int
	ProductionCap int
	Growth        int // 實際值（存值 − 100），是成長率不是增量
	Prevention    int // 防災值
	Garrison      int // 城兵數
	GarrisonCap   int
	Governor      int // 派駐的內政官（武將編號），0xFF ＝ 無
}

// General 是一名武將的完整狀態。
type General struct {
	Alive    bool
	Name     string
	Alias    string // 呼び名。多數與 Name 相同
	Aptitude [3]int // 已經 >>4 的小值
	Martial  int
	Command  int
	Politics int
	Timer    int // 每月遞減，歸零才行動
	Faction  int // 所屬勢力，0xFF ＝ 在野
	Posting  int // 派駐狀態，0xFF ＝ 未派駐
}

// Rules 回傳這名武將在規則層的視圖。
func (g General) Rules() general.General {
	return general.General{
		Name: g.Name, Aptitude: g.Aptitude,
		Martial: g.Martial, Command: g.Command, Politics: g.Politics,
	}
}

// World 是一整個遊戲狀態。
type World struct {
	// raw 是載入時那個劇本區塊的完整位元組。
	//
	// **存檔採「改寫」而不是「重建」**（CLAUDE.md §10）：
	// Save 會從這份原始位元組出發，只覆寫已經解出來的欄位，
	// 其餘一個 byte 都不動。這樣還沒解的區域（事件佇列、軍團表、
	// 那 69 byte 不載入的空隙…）能原封不動保留，
	// 存檔也不會因為我們理解不完整而損毀。
	raw []byte

	Clock    clock.Clock
	Factions [numFactions]Faction
	Cities   [numCities]City
	Generals [numGenerals]General

	// Player 是玩家所仕的勢力編號。原版存在 cs:0CFFh，
	// 但劇本檔裡沒有（開新遊戲時才選），所以預設 −1。
	Player int

	// Trust 是信賴度：君主對軍師（＝玩家）的評價。
	// 歸 0 → 被逐出勢力 → Game Over（docs/mechanics/80-victory.md §1）。
	//
	// ⚠ **它在存檔裡的位置還沒找到。** 勢力記錄 +0x1D 曾被誤判成信賴度，
	// 實際是士氣基準（docs/re/08 §4）。所以這一欄目前是**執行期狀態**，
	// 存檔不會保留 —— 找到欄位之前不要假裝存得起來。
	Trust int

	// 稅率與募兵數是**玩家專屬**的設定（AI 不用，見 docs/re/07 §8）。
	// Next 那一組是玩家在財政視窗改的值，月結時才搬到生效那一組。
	TaxRate        int
	RecruitCap     [economy.NumTroopTypes]int
	NextTaxRate    int
	NextRecruitCap [economy.NumTroopTypes]int
}

func u16(b []byte, off int) int { return int(binary.LittleEndian.Uint16(b[off:])) }

// i24 讀一個有號 24 位元的值。原版的資金就是這樣存的，
// 而且判斷正負是看第 23 位（docs/re/07 §2）。
func i24(b []byte, off int) int {
	v := int(b[off]) | int(b[off+1])<<8 | int(b[off+2])<<16
	if v&0x800000 != 0 {
		v -= 1 << 24
	}
	return v
}

func decodeName(b []byte) string {
	// 姓名是 6 byte 的 Big5（3 個全形字），不足的補全形空白。
	// 這裡不做編碼轉換——呼叫端要顯示時再處理，
	// 規則層本身不需要看得懂字。
	s := make([]byte, 0, 6)
	for i := 0; i < len(b); i++ {
		if b[i] == 0 {
			break
		}
		s = append(s, b[i])
	}
	return string(s)
}

// LoadScenario 從 SINARIO.DAT（或 SAVE.DAT）讀出第 index 個區塊。
func LoadScenario(path string, index int) (*World, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(raw) != blockSize*numBlocks {
		return nil, fmt.Errorf("%s 大小 %d，預期 %d", path, len(raw), blockSize*numBlocks)
	}
	if index < 0 || index >= numBlocks {
		return nil, fmt.Errorf("劇本編號 %d 超出 0–%d", index, numBlocks-1)
	}
	b := raw[index*blockSize : (index+1)*blockSize]

	// 信賴度的欄位還沒找到，開局先給說明書截圖上的值。
	w := &World{Player: -1, Trust: 200, raw: append([]byte(nil), b...)}

	// 遊戲時鐘（docs/formats/08 §1.1）。+0x01 的該月天數不另存，
	// clock 套件用 DaysInMonth(Month) 算得出來。
	w.Clock = clock.Clock{
		Day:     int(b[0x00]),
		Subtick: int(b[0x02]),
		Hour:    int(b[0x03]),
		Month:   u16(b, 0x04),
		Year:    u16(b, 0x06),
	}
	w.TaxRate = int(b[taxOffset])
	w.NextTaxRate = int(b[nextSettings])
	for i := 0; i < int(economy.NumTroopTypes); i++ {
		w.RecruitCap[i] = u16(b, recruitCapOffset+i*2)
		w.NextRecruitCap[i] = u16(b, nextSettings+2+i*2)
	}

	for i := range w.Factions {
		r := b[factionBase+i*factionSize:]
		w.Factions[i] = Faction{
			Alive:          r[0x00] >= 0x80,
			Lord:           int(r[0x01]),
			Advisor:        int(r[0x02]),
			Capital:        int(r[0x03]),
			Generals:       int(r[0x18]),
			Expense:        i24(r, 0x1A),
			Corps:          int(r[0x14]),
			InvasionTarget: int(r[0x19]),

			// ⚠ +0x1D 是**士氣基準值**，不是信賴度 ——
			// 編成軍團時它被複製進軍團的士氣欄位（docs/re/08 §4）。
			// 信賴度存在哪還沒解。
			MoraleBase: int(r[0x1D]),
			Funds:      i24(r, 0x20),
			Cities:     int(r[0x23]),
			Aggression: int(r[0x28]),
			Diplomat:   int(r[0x2A]),
		}
		for t := 0; t < int(economy.NumTroopTypes); t++ {
			w.Factions[i].Reserves[t] = u16(r, 0x04+t*2)
		}
	}

	for i := range w.Cities {
		r := b[cityBase+i*citySize:]
		w.Cities[i] = City{
			Name:          decodeName(r[0x02:0x08]),
			Owner:         int(r[0x01]),
			OwnerRecorded: int(r[0x1A]),
			X:             u16(r, 0x08),
			Y:             u16(r, 0x0A),
			ProductionCap: u16(r, 0x0C),
			Production:    u16(r, 0x0E),
			Growth:        int(r[0x10]) - 100, // 存值帶 +100 偏移
			Prevention:    int(r[0x11]),
			GarrisonCap:   int(r[0x12]),
			Garrison:      int(r[0x13]),
			Governor:      int(r[0x19]),
		}
	}

	for i := range w.Generals {
		r := b[generalBase+i*generalSize:]
		w.Generals[i] = General{
			Alive: r[0x00] >= 0x80,
			Name:  decodeName(r[0x02:0x08]),
			Alias: decodeName(r[0x08:0x0E]),
			Aptitude: [3]int{
				int(r[0x0E]) >> 4, int(r[0x0F]) >> 4, int(r[0x10]) >> 4,
			},
			Martial:  int(r[0x11]),
			Command:  int(r[0x12]),
			Politics: int(r[0x13]),
			Timer:    int(r[0x18]),
			Faction:  int(r[0x1C]),
			Posting:  int(r[0x1D]),
		}
	}
	return w, nil
}

// Event 是一個 tick 裡發生的事，供呼叫端記錄或呈現。
type Event struct {
	Clock      clock.Event
	Settled    bool                     // 這個 tick 跑了月結
	Disaster   map[int]economy.Disaster // 據點編號 → 災害
	Storm      *economy.StormArea
	Eliminated []int // 這個 tick 被判定滅亡的勢力
}

// Tick 推進一個 tick。月結、季節、災害都掛在對應的進位事件上，
// 順序照原版（docs/re/06 §5、docs/re/07 §1）。
func (w *World) Tick(rng economy.Rand) Event {
	ev := Event{Clock: w.Clock.Advance()}
	if !ev.Clock.Month {
		return ev
	}
	ev.Settled = true

	// ① 各勢力的月結。
	cities := w.economyCities()
	for i := range w.Factions {
		f := &w.Factions[i]
		if !f.Alive {
			continue
		}
		ef := economy.Faction{
			Funds:      f.Funds,
			Reserves:   f.Reserves,
			Capital:    cities[w.clampCity(f.Capital)],
			TaxRate:    w.TaxRate,
			RecruitCap: w.RecruitCap,
			Expense:    f.Expense,
			AI:         i != w.Player,
		}
		res := economy.Settle(&ef, cities, i, rng)
		f.Funds, f.Reserves, f.Expense = ef.Funds, ef.Reserves, ef.Expense
		f.Cities = res.Cities
	}

	// ② 生產力與上昇值。原版掃全部 192 個據點，不分勢力；
	//    稅率修正只套用在玩家的據點。
	for i := range w.Cities {
		c := &w.Cities[i]
		cs := economy.CityState{
			Production:    c.Production,
			ProductionCap: c.ProductionCap,
			Growth:        c.Growth,
			Owner:         c.Owner,
		}
		economy.GrowCity(&cs, w.TaxRate, c.Owner == w.Player, rng)
		c.Production, c.Growth = cs.Production, cs.Growth
	}

	// ③ 災害。
	ev.Disaster = map[int]economy.Disaster{}
	for i := range w.Cities {
		c := &w.Cities[i]
		if d := economy.RollCityDisaster(c.Prevention, c.Growth+100, rng); d != economy.NoDisaster {
			ev.Disaster[i] = d
		}
	}
	ev.Storm = economy.RollStorm(cities, rng)

	// ④ 「來月」的設定生效（原版是一次 4 個 word 的複製）。
	w.TaxRate, w.RecruitCap = w.NextTaxRate, w.NextRecruitCap

	// ⑤ 勢力滅亡判定。據點數歸零就出局。
	//    ⚠ 原版「滅亡」的精確條件還沒反組譯出來（docs/mechanics/80-victory.md §3），
	//    這裡先用最直覺的一條，並標成 remake 的暫定規則。
	for i := range w.Factions {
		if w.Factions[i].Alive && w.Factions[i].Cities == 0 {
			w.Factions[i].Alive = false
			ev.Eliminated = append(ev.Eliminated, i)
		}
	}
	return ev
}

func (w *World) clampCity(i int) int {
	if i < 0 || i >= numCities {
		return 0
	}
	return i
}

func (w *World) economyCities() []economy.City {
	out := make([]economy.City, numCities)
	for i, c := range w.Cities {
		out[i] = economy.City{X: c.X, Y: c.Y, Production: c.Production, Owner: c.Owner}
	}
	return out
}

// AliveFactions 回傳還存在的勢力編號。
func (w *World) AliveFactions() []int {
	var out []int
	for i, f := range w.Factions {
		if f.Alive {
			out = append(out, i)
		}
	}
	return out
}

// LordName 回傳某個勢力的君主姓名（原始 Big5 byte）。
func (w *World) LordName(faction int) string {
	f := w.Factions[faction]
	if f.Lord < 0 || f.Lord >= numGenerals {
		return ""
	}
	return w.Generals[f.Lord].Name
}

// ---------------------------------------------------------------------------
// 存檔
// ---------------------------------------------------------------------------

func putU16(b []byte, off, v int) {
	binary.LittleEndian.PutUint16(b[off:], uint16(v))
}

// putI24 寫一個有號 24 位元的值，並照原版的上下限鉗住。
func putI24(b []byte, off, v int) {
	if v > economy.MaxFunds {
		v = economy.MaxFunds
	}
	if v < economy.MinFunds {
		v = economy.MinFunds
	}
	u := uint32(v) & 0xFFFFFF
	b[off] = byte(u)
	b[off+1] = byte(u >> 8)
	b[off+2] = byte(u >> 16)
}

// Bytes 把世界狀態寫回一個 22,208 byte 的劇本／存檔區塊。
//
// **策略是「改寫」不是「重建」**：從載入時的原始位元組出發，
// 只覆寫下面列出的已解欄位。還沒解的區域一個 byte 都不動。
//
// 這條規則不是潔癖——原版區塊裡至少還有事件佇列（+0x52C0 起 1,024 B，
// 256 筆 × 4 B）、軍團表（+0x22C0 起 127 筆 × 32 B）、
// 以及 +0x3B–+0x7F 那 69 byte 不載入的空隙。
// 重建會把它們全部歸零，等於損毀存檔。
func (w *World) Bytes() []byte {
	b := append([]byte(nil), w.raw...)

	// 遊戲時鐘。+0x01 的該月天數是快取值，原版在換月時一起寫，
	// 這裡也一起寫回去，否則進位判斷會用到舊的天數。
	b[0x00] = byte(w.Clock.Day)
	b[0x01] = byte(clock.DaysInMonth(w.Clock.Month))
	b[0x02] = byte(w.Clock.Subtick)
	b[0x03] = byte(w.Clock.Hour)
	putU16(b, 0x04, w.Clock.Month)
	putU16(b, 0x06, w.Clock.Year)

	b[taxOffset] = byte(w.TaxRate)
	b[nextSettings] = byte(w.NextTaxRate)
	for i := 0; i < int(economy.NumTroopTypes); i++ {
		putU16(b, recruitCapOffset+i*2, w.RecruitCap[i])
		putU16(b, nextSettings+2+i*2, w.NextRecruitCap[i])
	}

	for i, f := range w.Factions {
		r := b[factionBase+i*factionSize:]
		if f.Alive {
			r[0x00] |= 0x80
		} else {
			r[0x00] &^= 0x80
		}
		r[0x01] = byte(f.Lord)
		r[0x02] = byte(f.Advisor)
		r[0x03] = byte(f.Capital)
		for t := 0; t < int(economy.NumTroopTypes); t++ {
			putU16(r, 0x04+t*2, f.Reserves[t])
		}
		r[0x14] = byte(f.Corps)
		r[0x18] = byte(f.Generals)
		r[0x19] = byte(f.InvasionTarget)
		putI24(r, 0x1A, f.Expense)
		r[0x1D] = byte(f.MoraleBase)
		putI24(r, 0x20, f.Funds)
		r[0x23] = byte(f.Cities)
		r[0x28] = byte(f.Aggression)
		r[0x2A] = byte(f.Diplomat)
	}

	for i, c := range w.Cities {
		r := b[cityBase+i*citySize:]
		r[0x01] = byte(c.Owner)
		copy(r[0x02:0x08], []byte(c.Name))
		putU16(r, 0x08, c.X)
		putU16(r, 0x0A, c.Y)
		putU16(r, 0x0C, c.ProductionCap)
		putU16(r, 0x0E, c.Production)
		r[0x10] = byte(c.Growth + 100)
		r[0x11] = byte(c.Prevention)
		r[0x12] = byte(c.GarrisonCap)
		r[0x13] = byte(c.Garrison)
		r[0x19] = byte(c.Governor)
		r[0x1A] = byte(c.OwnerRecorded)
	}

	for i, g := range w.Generals {
		r := b[generalBase+i*generalSize:]
		if g.Alive {
			r[0x00] |= 0x80
		} else {
			r[0x00] &^= 0x80
		}
		copy(r[0x02:0x08], []byte(g.Name))
		copy(r[0x08:0x0E], []byte(g.Alias))
		// 適性存的是 ×16 的值（讀進來時 >>4）。
		for k := 0; k < 3; k++ {
			r[0x0E+k] = byte(g.Aptitude[k] << 4)
		}
		r[0x11] = byte(g.Martial)
		r[0x12] = byte(g.Command)
		r[0x13] = byte(g.Politics)
		r[0x18] = byte(g.Timer)
		r[0x1C] = byte(g.Faction)
		r[0x1D] = byte(g.Posting)
	}
	return b
}

// SaveInto 把這個世界寫回檔案的第 index 個區塊，其餘三個區塊原封不動。
//
// ⚠ 原版資產是唯讀的（CLAUDE.md §10）。這個函式**只寫呼叫端指定的
// 輸出路徑**，不會就地改原始檔——要覆寫原版存檔是呼叫端的決定。
func (w *World) SaveInto(srcPath, dstPath string, index int) error {
	raw, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	if len(raw) != blockSize*numBlocks {
		return fmt.Errorf("%s 大小 %d，預期 %d", srcPath, len(raw), blockSize*numBlocks)
	}
	if index < 0 || index >= numBlocks {
		return fmt.Errorf("槽位 %d 超出 0–%d", index, numBlocks-1)
	}
	out := append([]byte(nil), raw...)
	copy(out[index*blockSize:(index+1)*blockSize], w.Bytes())
	return os.WriteFile(dstPath, out, 0o644)
}
