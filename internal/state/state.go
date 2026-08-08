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

	"github.com/wicanr2/wolong_cht/internal/rules/capital"
	"github.com/wicanr2/wolong_cht/internal/rules/clock"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/general"
	"github.com/wicanr2/wolong_cht/internal/rules/governor"
	"github.com/wicanr2/wolong_cht/internal/rules/march"
)

// 劇本／存檔的佈局常數（docs/formats/08 §0–§1.6）。
const (
	blockSize = 22208 // 0x56C0，一個劇本區塊
	numBlocks = 4

	factionBase, factionSize, numFactions = 0x0080, 64, 22

	// 交友度矩陣：列 ＝ 觀察者、欄 ＝ 對象，每列 24 byte（用到前 22 欄）。
	// 位址算法出自 `sub_13119`：`0x600 + 觀察者 × 24 + 對象`（段內偏移），
	// 檔案偏移再加 0x80。
	friendBase, friendStride              = 0x0680, 24
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

	// LowFunds 是記錄 +0x00 的 bit 6：資金低於「取消侵攻」門檻的**一半**
	// 時設起。用途還沒解——原版設了它但還沒找到誰讀（docs/re/08 §1）。
	LowFunds bool
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

	// Kind 是據點類型（記錄 +0x16 的低 4 位）：
	// 0 大城／1 中城／2 小城／3 關／4 戰場。**數字越小城越大。**
	//
	// 選首都時拿它當第一 key（`internal/rules/capital`），
	// 而類型本身也解釋了兩個先前看起來突兀的數字：
	// 關的城兵上限最高（194–254），戰場是 0（不能駐兵）。
	// 見 docs/formats/08 §1.6。
	Kind int

	// KindHigh 是 +0x16 的高 4 位（0–14）。嚴格巢狀在 Kind 內，
	// 疑似大地圖上的外觀編號，**未解**——存著只為了寫回時不失真。
	KindHigh int

	// Adjacency 是記錄 +0x00 的低 4 位：四個方向哪幾個有鄰接
	// （對應 +0x1C–+0x1F 的四個據點編號，docs/formats/08 §4.1）。
	Adjacency int
}

// General 是一名武將的完整狀態。
type General struct {
	Alive    bool
	Name     string
	Alias    string // 呼び名。多數與 Name 相同

	// Portrait 是 KAOGRF 的頁碼（記錄 +0x01）。
	//
	// **不是武將編號。** 曹操是第 16 個武將但頭像在第 50 頁、
	// 荀彧是第 62 個但在第 121 頁——拿武將編號當頁碼會畫出別人的臉。
	// 這個欄位原本只知道「127 筆各不相同」，是拿 PC-98 實機的君主確認畫面
	// 比對出來的：畫面上的曹操與 KAOGRF 第 50 頁是同一張圖
	// （docs/playtest/07）。KAOGRF 有 150 張 > 127 人，也對得上。
	Portrait int    // 見上（宣告順序照記錄偏移）
	Aptitude [3]int // 已經 >>4 的小值
	Martial  int
	Command  int
	Politics int
	Timer    int // 每月遞減，歸零才行動

	// Posted 是「出陣中」（記錄 +0x17）。編成軍團時原版寫 1
	// （`sub_16F26`），武將被俘時寫 4（`sub_129C3`）、釋放時歸零。
	Posted bool

	// Tactic 是戰場行動腳本編號（記錄 +0x16，值域 0–7）。
	// `BATTLE.DAT` 的段編號 ＝ 本值 × 4 ＋ 戰場類別（docs/re/11 §3.3）。
	// 它與能力值單調對應：0 是呂布那型（平均武力 13.6），7 是純文官（1.5）。
	Tactic int

	// LoyalToDeath 是記錄 +0x00 的 bit 4：**舊主已滅時寧可自刎也不改事二主**
	// （`sub_129C3` → 訊息 0x43，docs/re/09 §6）。
	// 旗標那個 byte 有 7 種值，目前只解出這一個位元。
	LoyalToDeath bool

	// Budget 是官員手上的經費餘額（記錄 +0x1A）。外交官每次工作
	// 扣 23 − 政治，歸零就停擺，要再向君主開口要錢。
	Budget int

	Faction int // 所屬勢力，0xFF ＝ 在野

	// Captor 是**舊主的勢力編號**（記錄 +0x1D），0xFF ＝ 非捕虜。
	//
	// ⚠ 這一欄先前叫 Posting，記成「派駐狀態」。四個劇本開局全是 0xFF，
	// 兩種解讀都說得通——是用法把它定下來的：戰敗被擒時寫入舊主
	// （`sub_129C3`）、釋放時清掉並通知舊主（`sub_150D7`）、
	// 月結時判歸降（`sub_1585F`）。見 docs/re/09 §6。
	Captor int
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
	// cityCursor 是據點整備的輪轉游標（原版 `word_10D1E`）。
	// 每 tick 前進一格，192 個據點輪一圈 ≈ 一天。
	cityCursor int

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

	// Corps 是軍團表。**索引與武將表平行**——軍團 i 由武將 i 帶
	// （`sub_1291A` 直接換算兩張表的位址，docs/re/09 §6）。
	// 出貨的劇本檔裡全零：開局沒有軍團，玩家要自己編成。
	Corps [numCorps]Corps

	// tactical 是戰術戰鬥的戰場來源；nil 表示全部走自動判定。
	tactical *TacticalSetup

	// cityBias 是載入時「勢力記錄的據點數 − 實際數出來的」，
	// 給不變量檢查當基準線（見 invariant.go）。
	cityBias [numFactions]int
	// pending 是一場等著被跑完的戰術戰鬥。它還在的時候世界不前進。
	pending *Pending
	// rng 是給戰術層用的亂數源，開戰時記下來。
	rng combat.Rand

	// corpsCursor 是下一支要更新的軍團（原版 cs:0D18h）。
	// 每 tick 只推進 16 支，所以掃完一輪要 8 個 tick。
	corpsCursor int

	// roads 是據點道路圖，從 MMAP 推導後由呼叫端注入（SetRoads）。
	// nil 表示沒有素材 —— 行軍退回直線移動。
	roads *march.Graph

	// routes[i] 是軍團 i 還沒走完的**地圖格**序列（沿真正的道路）。
	// 見 corps.go 的說明：它刻意不放進 Corps，
	// 那樣 Corps 才留得住 `==` 可比性（存檔 round-trip 測試靠它）。
	routes [numCorps][][2]int

	// hourFaction 是下一個輪到的勢力（原版 cs:0D1Ch，以 si 步進 0x40）。
	// 不匯出——它是迴圈的內部游標，不是遊戲狀態的一部分。
	hourFaction int

	// Player 是玩家所仕的勢力編號。原版存在 cs:0CFFh，
	// 但劇本檔裡沒有（開新遊戲時才選），所以預設 −1。
	Player int

	// Friendship 是交友度矩陣：Friendship[觀察者][對象]。
	//
	// **這是有向的**——A 對 B 與 B 對 A 是兩個獨立的值，
	// 畫面上只看得到自家君主那一側（docs/mechanics/50-diplomacy.md §1）。
	//
	// 對角線是 0xFF（自己），值 127 超過一般上限 100。
	Friendship [numFactions][numFactions]diplomacy.Friendship

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

// 原版的哨兵值。0xFF 表示「沒有」——**不等於 Go 的零值**（CLAUDE.md §8）。
const (
	noFaction = 0xFF
	noCity    = 0xFF
)

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

	for i := range w.Friendship {
		row := b[friendBase+i*friendStride:]
		for j := range w.Friendship[i] {
			w.Friendship[i][j] = diplomacy.Friendship(row[j])
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
			Kind:          int(r[0x16]) & 0x0F,
			KindHigh:      int(r[0x16]) >> 4,
			Adjacency:     int(r[0x00]) & 0x0F,
		}
	}

	for i := range w.Generals {
		r := b[generalBase+i*generalSize:]
		w.Generals[i] = General{
			Alive:    r[0x00] >= 0x80,
			Name:     decodeName(r[0x02:0x08]),
			Alias:    decodeName(r[0x08:0x0E]),
			Portrait: int(r[0x01]),
			Aptitude: [3]int{
				int(r[0x0E]) >> 4, int(r[0x0F]) >> 4, int(r[0x10]) >> 4,
			},
			Martial:      int(r[0x11]),
			Command:      int(r[0x12]),
			Politics:     int(r[0x13]),
			Tactic:       int(r[0x16]),
			Timer:        int(r[0x18]),
			Posted:       r[0x17] != 0,
			LoyalToDeath: r[0x00]&0x10 != 0,
			Budget:       int(r[0x1A]),
			Faction:      int(r[0x1C]),
			Captor:       int(r[0x1D]),
		}
	}
	w.loadCorps(b)
	// 記下據點數的基準差額，給不變量檢查用（見 invariant.go）。
	w.snapshotBias()
	return w, nil
}

// Event 是一個 tick 裡發生的事，供呼叫端記錄或呈現。
type Event struct {
	Clock      clock.Event
	Settled    bool                     // 這個 tick 跑了月結
	Disaster   map[int]economy.Disaster // 據點編號 → 災害
	Storm      *economy.StormArea
	Eliminated []int // 這個 tick 被判定滅亡的勢力

	// HourFaction 是這個 tick 輪到的勢力編號，−1 表示這個 tick 沒有輪到誰。
	// 原版每「時」只處理一個勢力，22 個勢力輪一圈（docs/re/08 §1）。
	HourFaction int

	// InvasionCancelled 為真表示上面那個勢力因為財政撐不住，
	// 侵攻目標被自動清掉了。
	InvasionCancelled bool

	// FriendshipUp 為真表示外交官這次做出了成果（交友度 +1）。
	FriendshipUp bool

	// Relocated 記錄這個 tick 遷都的勢力：勢力編號 → 新首都的據點編號。
	// 兩個觸發都會寫進來：首都被打下來（`sub_14DF0`），
	// 以及沒仗打時的主動遷都（事件 8）。
	Relocated map[int]int

	// Corps 是這個 tick 裡動過或打過的軍團。原版每 tick 只更新 16 支
	// （`sub_125A3` 的 `mov cx, 10h`），所以這裡通常是空的或很短。
	Corps []CorpsEvent
}

// Tick 推進一個 tick。月結、季節、災害都掛在對應的進位事件上，
// 順序照原版（docs/re/06 §5、docs/re/07 §1）。
func (w *World) Tick(rng economy.Rand) Event {
	// 有戰術戰鬥還沒打完，世界就停在那裡——原版進戰術畫面時戰略時間也停了。
	if w.pending != nil {
		return Event{HourFaction: -1}
	}
	w.rng = rng
	ev := Event{Clock: w.Clock.Advance(), HourFaction: -1}

	// 軍團先動。原版的主迴圈是「先 sub_125A3 再 sub_11D8E（時鐘）」，
	// 不過時鐘已經在上面推進了，所以這裡用推進後的小時去判軍費。
	ev.Corps = w.tickCorps(w.Clock.Hour, rng)

	// 據點整備：**每 tick 一個**，游標輪轉（原版 `sub_13EFD` 的
	// `mov si, word_10D1E` … `add si, 20h`）。192 個據點輪一圈，
	// 而一天是 216 tick，所以每個據點大約每天一次。
	w.tickCity(rng)

	if ev.Clock.Hour {
		w.hourly(&ev, rng)
	}
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
		// 月結會照實際的據點重算，開局那個差額到此消失
		// （劇本 3、4 的資料瑕疵，見 invariant.go）。
		w.cityBias[i] = 0
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

	// ⑤ 勢力滅亡判定。**據點與軍團都沒了**才出局。
	//
	//    ⚠ 原版「滅亡」的精確條件還沒反組譯出來
	//    （docs/mechanics/80-victory.md §3），這裡是 remake 的暫定規則。
	//
	//    初版只看據點數歸零，於是**還有軍團在野外的勢力會被判死**，
	//    留下一支沒有主人的軍團——不變量層的「已滅勢力還有軍團」抓到的
	//    就是這個。兩種修法裡選了這一種而不是「判死時順手刪掉軍團」：
	//    還有軍團在外的勢力仍然可能打下城來，直接刪軍團等於替原版
	//    決定了一條我們還沒讀出來的規則。**暫定規則要往保守的方向定。**
	//
	//    ⚠ 這個 bug 是**實作內政官之後才炸出來的**：城兵數變了 → 戰況變了
	//    → 才走到「最後一城被佔但軍團還在外面」那個組合。
	//    加一個會改變長期軌跡的機制，等於幫舊程式做了一次隨機測試。
	for i := range w.Factions {
		if w.Factions[i].Alive && w.Factions[i].Cities == 0 &&
			w.Factions[i].Corps == 0 {
			w.Factions[i].Alive = false
			ev.Eliminated = append(ev.Eliminated, i)
		}
	}
	return ev
}

// hourly 跑原版 `sub_13E11`：**每「時」只處理一個勢力**，
// 22 個勢力輪一圈 ＝ 22 小時，所以每個勢力大約每天被處理一次。
//
// 順序照原版：① 侵攻的財政檢查 → ② 預備兵維持費 → ③ 外交官。
// 完整反組譯見 docs/re/08 §1–§3。
func (w *World) hourly(ev *Event, rng economy.Rand) {
	i := w.hourFaction
	w.hourFaction = (i + 1) % numFactions
	ev.HourFaction = i

	f := &w.Factions[i]
	if !f.Alive {
		return
	}

	// ① 財政撐不住就自動取消侵攻。門檻與**據點數**掛鉤——
	//    地盤越大，維持一場侵攻需要的最低資金越高。
	if !diplomacy.CanSustainInvasion(f.Funds, f.Cities) {
		if f.InvasionTarget != diplomacy.NoTarget {
			ev.InvasionCancelled = true
		}
		f.InvasionTarget = diplomacy.NoTarget
		// 再低一半就設 bit 6。原版設了它，但還沒找到誰讀。
		f.LowFunds = !diplomacy.CanSustainInvasion(f.Funds*2, f.Cities)
	} else {
		f.LowFunds = false
	}

	// ② 預備兵維持費：三個兵種相加除以 32，累進本月支出。
	//    這是說明書 5.2「予備兵にも月単位で経費がかかり」的實際來源——
	//    它其實是**每小時**扣，只是月末才結算。
	total := 0
	for _, n := range f.Reserves {
		total += n
	}
	f.Expense = economy.ClampFunds(f.Expense + total>>5)

	// ③ 外交官。派駐在這個勢力的外交官是**別人派來的**，
	//    所以要改的是「派遣方 → 這個勢力」那一格交友度。
	ev.FriendshipUp = w.runDiplomat(i, rng)

	// ④ 主動遷都（原版的事件 8，`sub_12D3A`）。
	//    **沒有侵攻目標時**才會發，機率 rand(0..255) < 0x40 ＝ 25%。
	//    「閒著沒仗打就把首都搬到最好的城」——與侵攻互斥。
	if f.InvasionTarget == diplomacy.NoTarget && rng.Next()&0xFF < 0x40 {
		if next := w.relocateCapital(i); next != capital.None {
			if ev.Relocated == nil {
				ev.Relocated = map[int]int{}
			}
			ev.Relocated[i] = next
		}
	}
}

// relocateCapital 把 faction 的首都搬到 `internal/rules/capital` 選出的據點。
//
// 原版寫回時用 `xchg ah, [si+3]` 再 `cmp al, ah`——**位置沒變就什麼都不做**，
// 連通知都不發。這裡照抄那個相等判斷。
//
// 回傳新首都的據點編號，沒動或無據點可遷回 capital.None。
// **事件記錄交給呼叫端**——兩個觸發點的事件型別不同。
func (w *World) relocateCapital(faction int) int {
	sites := make([]capital.Site, len(w.Cities))
	for i := range w.Cities {
		c := &w.Cities[i]
		sites[i] = capital.Site{
			Owner: c.Owner, Kind: c.Kind,
			Production: c.Production, Adjacency: c.Adjacency,
		}
	}
	next := capital.Pick(sites, faction)
	if next == capital.None || next == w.Factions[faction].Capital {
		return capital.None
	}
	w.Factions[faction].Capital = next
	return next
}

// runDiplomat 讓派駐在 target 的外交官工作一次，回報交友度有沒有提升。
func (w *World) runDiplomat(target int, rng economy.Rand) bool {
	id := w.Factions[target].Diplomat
	if id < 0 || id >= numGenerals {
		return false // 0xFF ＝ 沒有外交官派駐
	}
	g := &w.Generals[id]
	sender := g.Faction
	if sender < 0 || sender >= numFactions || sender == target {
		return false
	}

	d := diplomacy.Diplomat{Politics: g.Politics, Budget: g.Budget}
	fr := w.Friendship[sender][target]
	up := d.Tick(&fr, rng)
	w.Friendship[sender][target] = fr
	g.Budget = d.Budget
	return up
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

	for i := range w.Friendship {
		row := b[friendBase+i*friendStride:]
		for j := range w.Friendship[i] {
			row[j] = byte(w.Friendship[i][j])
		}
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
		r[0x16] = byte(c.KindHigh<<4 | c.Kind&0x0F)
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
		r[0x16] = byte(g.Tactic)
		r[0x18] = byte(g.Timer)
		if g.Posted {
			r[0x17] = 1
		} else {
			r[0x17] = 0
		}
		r[0x1A] = byte(g.Budget)
		r[0x1C] = byte(g.Faction)
		r[0x1D] = byte(g.Captor)
	}
	w.saveCorps(b)
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

// tickCity 跑游標指到的那一個據點的整備（原版 `sub_13EFD` → `sub_14194`）。
//
// ⭐ **少了這一層，AI 的據點會單調掉到暴動。** 月結每月扣上昇值
// `rand(0..15)`（期望 −7.5），補回來的就是這裡：AI 的據點每天有 9/16
// 的機率 +1，月期望 +16.9。實作這一層之前，模擬跑 120 個月會出現
// 1872 次暴動（`docs/re/07` §19）。
func (w *World) tickCity(rng economy.Rand) {
	if len(w.Cities) == 0 {
		return
	}
	w.cityCursor = (w.cityCursor + 1) % len(w.Cities)
	c := &w.Cities[w.cityCursor]
	// ⚠ **中立據點也要整備。** 原版 `sub_13EFD` 的
	// `cmp byte ptr [si+841h], 18h / jz` 只跳過 `sub_13F74`，
	// **`sub_14194` 是無條件呼叫的**。
	//
	// 初版在這裡加了一個原版沒有的 owner 檢查，結果中立據點的上昇值
	// 只有每月 −rand(0..15) 而沒有回補，十年累積 681 次暴動——
	// 而玩家與 AI 的上昇值統計都是 +100，**看起來完全正常**。
	// 症狀出現在統計欄位看不到的那一群身上。

	gc := governor.City{
		Growth: c.Growth, Prevention: c.Prevention,
		Garrison: c.Garrison, GarrisonCap: c.GarrisonCap,
	}
	var gov *governor.Official
	// 只有玩家的據點會走內政官那一支（原版 `cmp al, [si+841h]`）。
	isPlayer := c.Owner == w.Player
	if isPlayer && c.Governor >= 0 && c.Governor < len(w.Generals) {
		g := &w.Generals[c.Governor]
		gov = &governor.Official{
			Politics: g.Politics, Martial: g.Martial, Budget: g.Budget,
		}
	}
	governor.Tick(&gc, gov, isPlayer, func() int { return rng.Next() & 0xFF })
	c.Growth, c.Prevention = gc.Growth, gc.Prevention
	c.Garrison = gc.Garrison
	if gov != nil {
		w.Generals[c.Governor].Budget = gov.Budget
	}
}
