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
	"bytes"
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
	// livingFactionsOffset 是存活勢力數在區塊裡的位置（docs/re/59 §3）。
	livingFactionsOffset = 0x3A

	// 交友度矩陣：列 ＝ 觀察者、欄 ＝ 對象，每列 24 byte（用到前 22 欄）。
	// 位址算法出自 `sub_13119`：`0x600 + 觀察者 × 24 + 對象`（段內偏移），
	// 檔案偏移再加 0x80。
	friendBase, friendStride              = 0x0680, 24
	cityBase, citySize, numCities         = 0x08C0, 32, 192
	generalBase, generalSize, numGenerals = 0x42C0, 32, 127

	// 稅率與三兵種募兵數。原版載到 cs:0D08h，而區塊前 59 B 對映到
	// cs:0CF0h，所以 0D08h − 0CF0h = 0x18 就是區塊內的偏移。
	// 信賴度則是 cs:0D00h（IDA `byte_10D00`），所以是區塊內的 +0x10。
	// Player 的原版 runtime 值是 cs:0CFFh（區塊 +0x0F）；前一個 word
	// cs:0CFDh（區塊 +0x0D）保存同一勢力的記錄表位址 `faction×0x40`。
	playerPtrOffset  = 0x0D
	playerOffset     = 0x0F
	trustOffset      = 0x10

	// titleOffset 是劇本／存檔的標題字串（Big5，NUL 結尾）。
	// 原版四槽選擇視窗的名稱欄畫的就是它（docs/re/52 §4）。
	titleOffset = 0x40
	titleMax    = 0x3B // 到區塊 +0x7B，之後是第 ② 塊的起點 +0x80
	taxOffset        = 0x18
	recruitCapOffset = 0x1A
	nextSettings     = 0x20 // 「來月」的同四項，月結時搬到上面兩處

	// 原版事件佇列位於區塊尾端；事件字的高 byte 可能是勢力／災害
	// 變體，不能只保存低 byte。這裡先保存原始 256 × 4 B，處理時序與
	// handler 效果仍由事件佇列的反組譯切片逐項接入。
	eventQueueOffset    = 0x52C0
	eventQueueEntrySize = 4
	eventQueueEntries   = 0x100
	eventQueueDispatch  = 0x40 // sub_131AE 只處理前 64 筆（0x100 byte）
)

// Faction 是一個勢力的完整狀態。
// NoAdvisor 是「這個勢力沒有軍師」。原版勢力記錄 +0x02 寫 0x7F
// （docs/formats/08），君主選擇畫面看到它就連軍師的頭像都不畫。
const NoAdvisor = 0x7F

type Faction struct {
	Alive    bool
	Lord     int // 君主的武將編號
	Advisor  int // 軍師的武將編號，**0x7F ＝ 無**（記錄 +0x02）
	Capital  int // 首都的據點編號
	Reserves [economy.NumTroopTypes]int

	Generals int // 武將數，月結時不重算（由武將表決定）
	Cities   int // 據點數，月結時重算
	Funds    int // 有號 24 位
	Expense  int // 本月累計支出，月結後歸零

	// MoraleBase 是新編成軍團的初始士氣（記錄 +0x1D，四個劇本都是 200）。
	//
	// 編成軍團時它被複製進軍團的士氣欄位，而說明書編成畫面的士氣值
	// 正好是 200，所以它是士氣基準值（docs/re/08 §4）。
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

	// ReliefSite 是最近一次求援的據點編號（記錄 +0x16）。
	// LostSite 是被別人佔走的據點編號（記錄 +0x17）——
	// 據點的 +0x1A 記的原主是自己、+0x01 已經不是（docs/re/44 §1、§3）。
	ReliefSite int
	LostSite   int

	Diplomat int // 派駐「這個」勢力的外交官（由別人派來），0xFF ＝ 無

	// LowFunds 是記錄 +0x00 的 bit 6：資金低於「取消侵攻」門檻的**一半**
	// 時設起（docs/re/08 §1）。
	//
	// 消費端是 AI 軍團挑目標的 Stage 2：設起時**跳過 LostSite**，
	// 只接 ReliefSite——沒錢就不主動反攻失土（docs/re/65 §4.1）。
	LowFunds bool
}

// City 是一個據點的完整狀態。
type City struct {
	// Name 是據點名稱（6 byte Big5，原始位元組）。北京、涿郡、武陵…
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

	// KindHigh 是 +0x16 的高 4 位：**`KYOGRF.DAT` 的張號**（0–14），
	// 據點情報視窗左半那張 96×96 景觀圖（docs/re/50 §3）。
	//
	// 值域只到 14 是算術逼出來的：原版算檔內位移時把高 16 位丟掉，
	// 張號 15 會溢位指到第一張的中間——15 張正好用滿。
	KindHigh int

	// Adjacency 是記錄 +0x00 的低 4 位：**四個鄰接槽裡哪幾個屬於別的勢力**。
	//
	// ⚠ 不是「哪幾個方向有鄰接」——有沒有鄰居由 +0x1C–+0x1F 的 0xFF
	// 哨兵決定。原版 `sub_1890A` 在據點換手時逐位設／清，並同步加減
	// +0x1B，所以這四位是會動的（docs/re/44 §5）。
	Adjacency int

	// EnemyNeighbours 是記錄 +0x1B ＝ Adjacency 的位元個數
	// （相鄰的敵方據點數）。它是 0 時原版連威脅掃描都不做。
	EnemyNeighbours int

	// ReliefCooldown 是記錄 +0x17：求援冷卻計時器，每輪 −1，非 0 時不再求援。
	// 玩家的據點求援後寫「亂數(0–15) ＋ 24」，AI 的據點寫「離首都的距離」
	// （上限 30）。開局全 0 是因為還沒有人求過援（docs/re/40 §4）。
	ReliefCooldown int

	// Threat 是記錄 +0x14 ＝ 鄰接敵方據點的 Occupancy 總和（周邊威脅量），
	// 每次輪到這個據點時由 `sub_13FA9` 重算。
	Threat int

	// Threatened 是記錄 +0x00 的位元 7、Specific 是位元 6：
	// **受威脅**與**威脅裡有本勢力的侵攻目標**（docs/re/40 §2）。
	//
	// 兩個旗標由威脅掃描每輪重寫，AI 軍團的 Stage 0–2 讀它們決定
	// 要駐守還是出擊（docs/re/65 §2–§4）。
	Threatened bool
	Specific   bool

	// Occupancy 是記錄 +0x18 ＝ 停在這個據點那一格的軍團數。
	//
	// ⚠ **是快取不是狀態。** 原版每次輪到這個據點就從單位佔用圖
	// （98,304 B 的計數器陣列）重抄一次，所以 remake 要重算不要記帳
	// （docs/re/44 §1）。
	Occupancy int

	// Neighbours 是記錄 +0x1C–+0x1F 的四個鄰接據點編號，0xFF ＝ 沒有。
	// 這些槽是政略 AI 建立「相鄰勢力清單」的直接來源
	// （原版 sub_12C52 → sub_12CDF），不能只用由地圖推導的道路圖替代。
	Neighbours [4]int
}

// General 是一名武將的完整狀態。
type General struct {
	Alive bool
	Name  string
	Alias string // 呼び名。多數與 Name 相同

	// Portrait 是 KAOGRF 的頁碼（記錄 +0x01）。
	//
	// **不是武將編號。** 曹操是第 16 個武將但頭像在第 50 頁、
	// 荀彧是第 62 個但在第 121 頁——拿武將編號當頁碼會畫出別人的臉。
	// 這個欄位原本只知道「127 筆各不相同」，是拿 PC-98 實機的君主確認畫面
	// 比對出來的：畫面上的曹操與 KAOGRF 第 50 頁是同一張圖
	// （docs/playtest/07）。KAOGRF 有 150 張 > 127 人，也對得上。
	Portrait int // 見上（宣告順序照記錄偏移）

	// TalkVariant 是記錄 +0x1E 的原始值；sub_13C99 以它取 0–2 的
	// TALK 變體（值 >= 3 時先減 3）。它不是武將編號，也不應拿來
	// 推測人物身分；保留原始欄位是為了讓事件 2／3 的 composite TALK
	// 可以回到正確的 prompt。語意等級：已證實（IDA 00013C99）。
	TalkVariant int

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
	// 它與能力值單調對應：0 是呂布那型（平均武術 13.6），7 是純文官（1.5）。
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

	// events 是原版區塊 +0x52C0 的 256 筆事件佇列。事件字完整保留
	// u16（包含高 byte 的勢力／災害分流），參數也是 u16；目前接上
	// 存檔 round-trip、時鐘／壓縮，以及已有獨立證據的 1／2／3／4／5／6／7／8／9／10／11／12／13 handler（10 為訊息邊界、11／12 為 runtime marker 與延遲效果）。
	events [eventQueueEntries]QueuedEvent

	// disasterMarkers 是事件 11／12 在執行期寫入據點記錄 +0x15 的
	// 災害 marker。原版 sub_14269 會在該據點輪到時消耗這個 marker，
	// 並把防災值／上昇值／生產力／城兵寫回 City；marker 本身是 runtime
	// 欄位，事件 12 的清除事件到達後才歸零。這兩個陣列不序列化。
	disasterMarkers      [numCities]economy.Disaster
	disasterMarkerLevels [numCities]uint8
	// disasterObjects 是 sub_123FF／sub_12459／sub_12533 的 32 筆非存檔
	// runtime 物件；事件 12 的火災／暴動才會建立，清除事件會移除。
	disasterObjects [disasterObjectSlots]disasterObject
	stormArea       *economy.StormArea

	// eventCursor／eventDelay 是原版的 runtime 游標（`word_10D20`）與
	// 節流計數（`byte_131AD`），不在存檔區塊內；載入新狀態與月結都重設。
	eventCursor int
	eventDelay  uint8

	// tactical 是戰術戰鬥的戰場來源；nil 表示全部走自動判定。
	tactical *TacticalSetup

	// cityBias 是載入時「勢力記錄的據點數 − 實際數出來的」，
	// 給不變量檢查當基準線（見 invariant.go）。
	cityBias [numFactions]int
	// pending 是一場等著被跑完的戰術戰鬥。它還在的時候世界不前進。
	pending *Pending
	// encounter 是一場玩家尚未選擇「戰鬥指揮／委任」的遭遇。
	// 它和 pending 一樣會凍結戰略時間，但還沒有建立戰術戰場。
	encounter *EncounterChoice
	// diplomacy 是事件 2／3 的玩家互動；它是 runtime 狀態，不進存檔。
	diplomacy *DiplomacyChoice
	// funding 是事件 4／5 的玩家撥款互動；它是 runtime 狀態，不進存檔。
	funding *FundingChoice
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

	// Title 是劇本／存檔的標題（區塊 +0x40，Big5 原始 byte）：
	// 「第一章」「呂布歸天」之類。四槽選擇視窗的名稱欄印的就是它。
	Title string

	// Player 是玩家所仕的勢力編號。原版存於 cs:0CFFh（區塊 +0x0F），
	// 同一全域區段的 cs:0CFDh（+0x0D）保存勢力表位址；新劇本兩者都是
	// 0xFF／0xFFFF，只有玩家選定或有效存檔才有值。
	Player int

	// strategicAI 是執行期開關。載入／存檔本身不包含「誰是玩家」這個
	// 啟動參數；wlgame／wlsim 在設定 Player 後明確啟用，讓純格式／規則
	// 測試可以仍然只跑已驗證的時鐘與月結，不被長期 AI 軌跡混入。
	strategicAI bool

	// approximateEvent10 是事件 10 的 remake 近似 producer 開關。原版
	// `sub_13496` 的 consumer 已知，但自然 `0x0A` writer 尚未定位；載入
	// 劇本時預設開啟，讓正常遊戲有可玩的俘虜消息；raw fixture 可關閉它，
	// 避免把替代規則混入原版 queue 邊界測試。這是 runtime 設定，不存檔。
	approximateEvent10 bool

	// outcome 是 runtime 的單次結果 latch，不寫入存檔。只由已證實的
	// 信賴度／據點 mutation 邊界設定；不能由 UI 每幀掃描推導。
	outcome OutcomeKind

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
	// 原版載入／存檔把 cs:0CF0h 起的 0x3B byte 整段搬入／搬出；
	// 說服流程的 `byte_10D00`（IDA `seg000:10D00`）在這段的 +0x10，
	// 因此它是可持久化的 u8。勢力記錄 +0x1D 是士氣基準，不是信賴度。
	Trust int

	// LivingFactions 是還沒滅亡的勢力數（區塊 +0x3A）。
	// **只在 eliminateFaction 裡減**——原版也只有一個 `dec`
	// （docs/re/59 §3）。減到 1 就是結局。
	LivingFactions int

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

// QueuedEvent 是原版事件佇列的一筆原始記錄。
//
// Code 保留完整 u16：低 byte 是 dispatch code，高 byte 在部分路徑承載
// 勢力或災害變體。這個型別用於存檔 round-trip 與已接入的事件 handler；
// 不要把高 byte 丟掉。
type QueuedEvent struct {
	Code  uint16
	Param uint16
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
	return loadBlock(raw[index*blockSize : (index+1)*blockSize]), nil
}

// blockTitle 取出區塊 +0x40 的標題（Big5，NUL 結尾）。
func blockTitle(b []byte) string {
	raw := b[titleOffset : titleOffset+titleMax]
	if i := bytes.IndexByte(raw, 0); i >= 0 {
		raw = raw[:i]
	}
	return string(raw)
}

// loadBlock 從一個劇本／存檔區塊建出 World。
// LoadScenario 與 LoadBlock（`snapshot.go`）共用這一支——
// **解碼只留一份實作**（`CLAUDE.md` §7 第 6 條）。
func loadBlock(b []byte) *World {
	// Trust 是原始全域區段的 byte_10D00（區塊 +0x10）。
	player := -1
	if p := int(b[playerOffset]); p >= 0 && p < numFactions &&
		u16(b, playerPtrOffset) == p*factionSize {
		player = p
	}
	w := &World{
		Player:      player,
		Title:       blockTitle(b),
		Trust:       int(b[trustOffset]),
		raw:         append([]byte(nil), b...),
		eventDelay:  7,
		eventCursor: 0,
		// 事件 10 的原版 producer unknown；remake 預設使用明確標示的
		// 近似 producer，仍可由 SetApproximateEvent10(false) 關閉。
		approximateEvent10: true,
	}
	for i := range w.events {
		off := eventQueueOffset + i*eventQueueEntrySize
		w.events[i] = QueuedEvent{
			Code:  binary.LittleEndian.Uint16(b[off:]),
			Param: binary.LittleEndian.Uint16(b[off+2:]),
		}
	}

	// 遊戲時鐘（docs/formats/08 §1.1）。+0x01 的該月天數不另存，
	// clock 套件用 DaysInMonth(Month) 算得出來。
	w.Clock = clock.Clock{
		Day:     int(b[0x00]),
		Subtick: int(b[0x02]),
		Hour:    int(b[0x03]),
		Month:   u16(b, 0x04),
		Year:    u16(b, 0x06),
	}
	// 存活勢力數（區塊 +0x3A，59 byte 全域區塊的最後一格）。
	// 原版 `cs:0D2Ah` 全庫只有一個 `dec`，靠這個欄位載入初值；
	// 減到 1 就是結局（docs/re/59 §3）。
	w.LivingFactions = int(b[livingFactionsOffset])
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
			ReliefSite:     int(r[0x16]),
			LostSite:       int(r[0x17]),

			// +0x1D 是**士氣基準值**，不是信賴度 —— 編成軍團時它被
			// 複製進軍團的士氣欄位（docs/re/08 §4）。信賴度在全域 +0x10。
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
			Threatened:    r[0x00]&0x80 != 0,
			Specific:      r[0x00]&0x40 != 0,
			EnemyNeighbours: int(r[0x1B]),
			Threat:          int(r[0x14]),
			Occupancy:       int(r[0x18]),
			ReliefCooldown:  int(r[0x17]),
			Neighbours:    [4]int{int(r[0x1C]), int(r[0x1D]), int(r[0x1E]), int(r[0x1F])},
		}
	}

	for i := range w.Generals {
		r := b[generalBase+i*generalSize:]
		w.Generals[i] = General{
			Alive:       r[0x00] >= 0x80,
			Name:        decodeName(r[0x02:0x08]),
			Alias:       decodeName(r[0x08:0x0E]),
			Portrait:    int(r[0x01]),
			TalkVariant: int(r[0x1E]),
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
	return w
}

// ResolveTalkFormatter2 重現 DOS/V TALK formatter `\\2` 取字串的 raw 規則。
//
// 證據：DOS/V KI.EXE `seg000:000108DB` 先從 SS:[DI] 取一個 word；高 byte
// 為 FF 時，把低 byte 轉成 `0x0840 + city×0x20`，否則直接把 word 當成
// DS 位移，兩條路徑最後都再加 2。這裡的 DS 動態區對應所載劇本區塊的
// `+0x80` 至事件佇列前；不把 queue 或區塊外資料當成可顯示字串，解析
// 不到時回傳 false，讓呈現層整則 fail-closed。
func (w *World) ResolveTalkFormatter2(word uint16) ([]byte, bool) {
	const dynamicSize = eventQueueOffset - factionBase // 0x5240
	if w == nil || len(w.raw) < factionBase+dynamicSize {
		return nil, false
	}

	off := int(word)
	if byte(word>>8) == 0xFF {
		city := int(byte(word))
		if city >= numCities {
			return nil, false
		}
		off = runtimeCityBase + city*citySize
	}
	off += 2
	if off < 0 || off >= dynamicSize {
		return nil, false
	}

	data := w.raw[factionBase+off : factionBase+dynamicSize]
	for i, b := range data {
		if b == 0 {
			data = data[:i]
			break
		}
	}
	return append([]byte(nil), data...), true
}

// FundingKind 是原版事件 4／5 的玩家撥款互動類型。
type FundingKind byte

const (
	FundingGovernor FundingKind = 4
	FundingDiplomat FundingKind = 5
)

// FundingOption 是原版 sub_139E8 的三列選項。
type FundingOption byte

const (
	FundingFullAmount FundingOption = iota
	FundingSetAmount
	FundingReject
)

// FundingChoice 是暫存在 runtime 的玩家撥款視窗，不寫入存檔。
// Subject 是事件 4 的據點編號或事件 5 的勢力編號；Officer 是實際收到
// 經費的武將編號。RequestedAmount 是事件佇列帶來的原始要求，OfferAmount
// 是玩家在「指定金額」路徑裡目前輸入的值。
type FundingChoice struct {
	Kind            FundingKind
	Subject         int
	Officer         int
	RequestedAmount int
	OfferAmount     int
}

// AmountEdit 是原版 sub_17C6E 數值輸入器已證實的編輯動作。
//
// 原始函式的 DOS/V 輸入 API／平台按鍵仍由呈現層映射；這裡只保存它呼叫的
// 數值語意，讓事件 2／3 與 4／5 共用同一組邊界與上限。
type AmountEdit byte

const (
	AmountAppendDigit AmountEdit = iota
	AmountAppendHundred
	AmountDeleteDigit
	// AmountSetMax 是原版 `sub_17DEC` 的 `mov si, [bp+0]`——那一格是
	// **上限**，不是呼叫端給的初值（`si` 開場就是 0，呼叫端沒有機會給初值）。
	// 圖庫上那顆鍵的字樣本來就寫著「最大」（docs/spec/78 §1）。
	AmountSetMax
	AmountClear
	AmountFinishInput
)

// DiplomacyKind 是原版事件 2／3 的玩家外交互動類型。
type DiplomacyKind byte

const (
	DiplomacyCooperation DiplomacyKind = 2
	DiplomacyCeasefire   DiplomacyKind = 3
)

// DiplomacyOption 是原版 sub_13902 的三列選項。
// OfferFunds 使用狀態層已解出的預設說服金額；數值編輯器的狀態語意由
// AmountEdit 保存，PC-98 掃描碼與逐頁畫面仍屬呈現層工作。
type DiplomacyOption byte

const (
	DiplomacyAcceptFree DiplomacyOption = iota
	DiplomacyOfferFunds
	DiplomacyReject
)

// DiplomacyChoice 是暫存在 runtime 的玩家外交視窗，不寫入存檔。
type DiplomacyChoice struct {
	Kind          DiplomacyKind
	Source        int
	Invader       int
	Target        int
	InitialAmount int
	OfferAmount   int
}

// Event 是一個 tick 裡發生的事，供呼叫端記錄或呈現。
type Event struct {
	Clock      clock.Event
	Settled    bool                     // 這個 tick 跑了月結
	Disaster   map[int]economy.Disaster // 據點編號 → 災害
	Storm      *economy.StormArea
	Eliminated []int // 這個 tick 被判定滅亡的勢力

	// TalkNotices 是事件處理器已確認要交給呈現層的 TALK.DAT 訊息。
	// 它只保存原版索引與可回查的 state 目標，不把文字、編碼或排版塞進
	// state；索引與 marker 來源見 docs/re/07 §18、docs/formats/01。
	TalkNotices []TalkNotice

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

	// Diplomacy 是這個 tick 剛掛起的玩家外交選擇；選擇完成前世界停止。
	Diplomacy *DiplomacyChoice

	// Funding 是這個 tick 剛掛起的事件 4／5 撥款選擇；選擇完成前世界停止。
	Funding *FundingChoice

	// Strategy 記錄政略 AI 在這個 tick 做出的「宣戰／編成」動作。
	// Corps 或 Destination 為 −1 時，表示該欄位在這筆事件沒有動作。
	// 這是觀測用事件，不寫回存檔。
	Strategy []StrategyEvent

	// Corps 是這個 tick 裡動過或打過的軍團。原版每 tick 只更新 16 支
	// （`sub_125A3` 的 `mov cx, 10h`），所以這裡通常是空的或很短。
	Corps []CorpsEvent

	// ReleasedGenerals 記錄事件 9（或其同一狀態收尾）釋放的武將索引，
	// 讓 GUI 能取用原版 TALK.DAT 句子並顯示可回查的事件通知；對話框排版
	// 仍由呈現層決定，不寫入存檔。
	ReleasedGenerals []int
}

// DisasterMarker 是事件 11／12 留在執行期據點記錄上的災害標記。
//
// 這不是存檔欄位，也不代表已解出原版物件動畫；呈現層只能把它當成
// 「目前有一個待套用／仍在顯示中的災害狀態」來讀。Level 保留原版
// marker 的強度，讓不同呈現層可以做一致的可視化，但不得把它解讀成
// 動畫幀數或剩餘時間。
type DisasterMarker struct {
	Kind  economy.Disaster
	Level uint8
}

// DisasterMarkerAt 回傳指定據點目前的 runtime 災害 marker。
//
// 這個唯讀方法刻意不暴露內部陣列，也不把 runtime marker 寫進存檔；
// 無效據點與 NoDisaster 都回傳 ok=false。它是 wlgame 視覺層的接縫，
// 不是新的規則來源。
func (w *World) DisasterMarkerAt(cityID int) (DisasterMarker, bool) {
	if cityID < 0 || cityID >= len(w.disasterMarkers) {
		return DisasterMarker{}, false
	}
	kind := w.disasterMarkers[cityID]
	if kind == economy.NoDisaster {
		return DisasterMarker{}, false
	}
	return DisasterMarker{Kind: kind, Level: w.disasterMarkerLevels[cityID]}, true
}

// StormAreaSnapshot 回傳目前 runtime 暴風雨範圍的副本。
//
// StormArea 由事件 11 的 handler 暫存，沒有值時回傳 ok=false；回傳副本
// 可避免 UI 意外修改規則層的指標內容。
func (w *World) StormAreaSnapshot() (economy.StormArea, bool) {
	if w.stormArea == nil {
		return economy.StormArea{}, false
	}
	return *w.stormArea, true
}

// TalkNotice 是一則原版訊息的結構化呈現要求。
//
// Index 是 TALK.DAT 的原始槽位。City、Faction 與 General 是可選的 state 索引，
// 未使用時為 -1；Amount 是 `\\7` 數值 marker 的原始十進位值，未使用時為 -1。
// 目前災害訊息使用 City，外交官回報使用 Faction／General，事件 10
// 使用 General 保存 formatter 的高位元組，事件 13 不帶目標。
// RawFormatterWordValid／RawFormatterWord 是原版直接呼叫 TALK 時，從
// `SS:[DI]` 取出的原始 word；它不是 City ID。Valid 必須顯式設起，避免
// Go 結構的零值把「未提供」誤當成原版 word 0。
// 這裡不直接保存展開後文字，讓 Big5 round-trip 與 UI 排版仍由資產／呈現層
// 負責，也避免把尚未解出的 formatter 參數誤升格成語意。
// CapitalMovedTalkBase 是遷都之後那一句的**組編號**（不是索引）。
//
// `sub_133FD` 對兩條路都傳 `cx = 0x1A4`，展開後是 518–525
// （docs/spec/64）：0–2 是自國君主下令，3–7 是他國遷都的情報回報。
// 說話者不同，取到的變體就不同——展開要用**原始** `+0x1E`。
const CapitalMovedTalkBase = 0x1A4

type TalkNotice struct {
	Index                 int
	City                  int
	Faction               int
	General               int
	Amount                int
	RawFormatterWord      int
	RawFormatterWordValid bool
	Secondary             bool // sub_13C3D 的第二次 TALK 呼叫
	NoPortrait            bool // 原版直接 sub_18810 的文字／選單沒有肖像 blit
}

// StrategyEvent 是政略 AI 的可觀測動作。
type StrategyEvent struct {
	Faction     int
	Target      int
	Corps       int
	Destination int
}

// Tick 推進一個 tick。月結、季節、災害都掛在對應的進位事件上。
//
// 原版 sub_11CD0 的可觀測順序是：據點 sub_13EFD → 軍團 sub_125A3
// → MCH 物件 sub_12459 → 時鐘 sub_11D8E。Tick 是規則 tick 的
// 據點／軍團／時鐘部分；一次可見 map-loop 的完整順序由 TickMap 提供，
// 不把 UI 的 g.speed 倍速誤套到物件動畫。
func (w *World) Tick(rng economy.Rand) Event {
	return w.tick(rng, false)
}

// TickMap 跑一次完整的原版 map-loop：據點、軍團、MCH 物件，最後才是時鐘。
// wlgame 每個可見畫面的第一個規則 tick 使用它；同畫面的額外 speed tick
// 使用 Tick，避免自訂快轉改變物件 16-update cadence。
func (w *World) TickMap(rng economy.Rand) Event {
	return w.tick(rng, true)
}

func (w *World) tick(rng economy.Rand, includeMapObjects bool) Event {
	// sub_11CB1 離開主遊戲循環後，remake 的 state 不再讓 Tick／AI／時鐘
	// 產生任何副作用；已經排入的訊息仍可由呈現層讀取。
	if w.outcome != InProgress {
		return Event{HourFaction: -1}
	}
	// 有戰術戰鬥、玩家尚未選擇的遭遇、外交提案或撥款請求，世界都停在那裡。
	// 原版進戰術畫面前也會先問「戰鬥指揮／委任」，這個選單同樣不能讓
	// 下一個軍團或時鐘在背景偷偷前進。
	if w.pending != nil || w.encounter != nil || w.diplomacy != nil || w.funding != nil {
		return Event{HourFaction: -1}
	}
	w.rng = rng
	ev := Event{HourFaction: -1}

	// ① 據點整備：**每 tick 一個**，游標輪轉（原版 `sub_13EFD` 的
	// `mov si, word_10D1E` … `add si, 20h`）。192 個據點輪一圈，
	// 而一天是 216 tick，所以每個據點大約每天一次。
	cityStrategy, cityNotices := w.tickCity(rng)
	ev.Strategy = append(ev.Strategy, cityStrategy...)
	ev.TalkNotices = append(ev.TalkNotices, cityNotices...)

	// ② 軍團在時鐘進位前更新，使用尚未 Advance 的小時。
	ev.Corps = w.tickCorps(w.Clock.Hour, rng)

	// ③ 完整 map-loop 才更新一次 MCH 物件；額外的 speed tick 跳過這裡。
	if includeMapObjects {
		w.AdvanceDisasterObjects()
	}

	// ④ 最後才進入 sub_11D8E。
	ev.Clock = w.Clock.Advance()

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

	// 原版 sub_12BD9 緊接月結經濟處理後壓縮事件佇列，並重設
	// `word_10D20`／`byte_131AD`。已證實的 queue 邊界先照原版保存，避免
	// 積壓事件跨月時留下錯誤資料；事件 handler 由每小時流程逐筆取出。
	w.compactEventQueue()
	// 原版 sub_15358 在月結壓縮後先跑 sub_15715／sub_1578F，將玩家
	// 內政官／外交官的撥款請求放進事件佇列，再進入其他政略評估。
	if w.Player >= 0 && w.Player < numFactions {
		w.queueFundingRequests(rng)
	}
	// 暴風雨與火災／暴動也是月結產生的事件：sub_122DB 先寫事件 11，
	// sub_12286 再按據點順序寫事件 0x010C／0x020C。事件 12 的 Param
	// 保存原版 runtime city record 位址，不把檔案偏移或 city ID 偷換進去。
	w.stormArea = ev.Storm
	if ev.Storm != nil {
		w.queueEvent(rng, 0, 11, 0, 0xFF)
	}
	for i := range w.Cities {
		d, ok := ev.Disaster[i]
		if !ok {
			continue
		}
		variant := 0
		if d == economy.Fire {
			variant = 1
		} else if d == economy.Riot {
			variant = 2
		} else {
			continue
		}
		param := uint16(runtimeCityBase + i*citySize)
		w.queueEvent(rng, variant, 12, param, 0xFF)
	}

	// 原版 sub_12BD9 在月結的經濟處理之後跑政略評估。宣戰／遷都決策先
	// 寫入 queue，再由每小時的 sub_131AE 邊界逐筆處理；其餘尚未解出的
	// handler 仍不在這個轉接層裡。每一筆決策仍只使用已由機器碼確認的
	// 欄位與比較式。
	if w.strategicAI {
		strategyEvents, relocated := w.runStrategicAI(rng)
		ev.Strategy = append(ev.Strategy, strategyEvents...)
		if len(relocated) > 0 {
			ev.Relocated = relocated
		}
	}

	// 原版 event 10 的自然 writer 仍未知。remake 在所有已證實的月結／
	// queue producer 之後，使用獨立、可關閉的近似 producer；它只產生
	// raw `(general<<8)|0x0A`，實際 TALK 仍等下一次每時 dispatcher。
	w.produceApproximateEvent10(rng)

	// ④ 「來月」的設定生效（原版是一次 4 個 word 的複製）。
	w.TaxRate, w.RecruitCap = w.NextTaxRate, w.NextRecruitCap

	// ⑤ 勢力滅亡判定：**據點數歸零就出局**（原版 `sub_14CF3` 的
	//    `dec [bx+23h]` → `sub_14DF0` 回 CF ＝ 1 → `sub_14FCE`，
	//    docs/re/59 §4）。
	//
	//    原版是在據點易主的當下判的，remake 那條路在 `capture`；
	//    這個月結掃描是後備，接住不是由攻城造成的據點歸零。
	//    軍團不必另外處理——`eliminateFaction` 會把在外的軍團一起收掉
	//    （原版 `sub_14FCE` 的武將處置迴圈，docs/re/59 §4.1）。
	for i := range w.Factions {
		if w.Factions[i].Alive && w.Factions[i].Cities == 0 {
			w.eliminateFaction(i, noFaction)
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
	// 原版 sub_13E11 的第一個呼叫就是 sub_131AE；它與當小時輪到的
	// 勢力財政檢查分開，不能把事件延後到月結或直接同步套用。
	w.dispatchQueuedEvent(ev)
	if w.diplomacy != nil || w.funding != nil {
		return
	}

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
		// 再低一半就設 bit 6：AI 軍團挑目標時會跳過「奪回失土」。
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
	old := w.Factions[faction].Capital
	w.Factions[faction].Capital = next
	// 原版 `sub_133FD`／`sub_14DF0` 在首都真的變更後都呼叫
	// `sub_14502`：同勢力、以舊首都為 Ordered 的活軍團改掛新首都；
	// 若目標正好是新首都，目標據點只改回舊首都，X/Y 保留原值。
	// 這裡只接入已證實的三個欄位效果，不虛構原版的路徑重算。
	w.syncCorpsAfterCapitalChange(faction, old, next)
	return next
}

// syncCorpsAfterCapitalChange 是 `sub_14502` 的非破壞性欄位轉接。
//
// 原版掃描 127 筆軍團，條件是存在旗標 >= 0x80、勢力等於來源，且
// `Corps+0x20` 等於舊首都；命中後寫 `+0x20 = newCapital`。若
// `Corps+0x14` 等於新首都×8，另寫成舊首都×8，並保留 `+0x16/+0x18`。
func (w *World) syncCorpsAfterCapitalChange(faction, oldCapital, newCapital int) {
	if faction < 0 || faction >= len(w.Factions) || oldCapital < 0 || newCapital < 0 {
		return
	}
	for i := range w.Corps {
		c := &w.Corps[i]
		if !c.Alive || c.Faction != faction || c.Ordered != oldCapital {
			continue
		}
		c.Ordered = newCapital
		if c.TargetNode == newCapital {
			c.TargetNode = oldCapital
		}
	}
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

// EnableStrategicAI 啟用月結政略評估與敵方軍團編成。
//
// 這是執行期設定，不會寫入 SINARIO.DAT／SAVE.DAT；呼叫端應在設定 Player
// 後呼叫。未啟用時，World 仍可作為純經濟／格式驗證模型使用。
func (w *World) EnableStrategicAI() { w.strategicAI = true }

// SetApproximateEvent10 控制事件 10 的 remake 近似自然 producer。
// false 適合只驗證原始 queue／consumer 的 fixture；正常遊戲預設為 true。
func (w *World) SetApproximateEvent10(enabled bool) {
	if w != nil {
		w.approximateEvent10 = enabled
	}
}

// ApproximateEvent10Enabled 回傳目前是否啟用事件 10 近似 producer。
func (w *World) ApproximateEvent10Enabled() bool {
	return w != nil && w.approximateEvent10
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

func clampU8(v int) int {
	if v < 0 {
		return 0
	}
	if v > 0xFF {
		return 0xFF
	}
	return v
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
	b[livingFactionsOffset] = byte(w.LivingFactions)
	b[0x00] = byte(w.Clock.Day)
	b[0x01] = byte(clock.DaysInMonth(w.Clock.Month))
	b[0x02] = byte(w.Clock.Subtick)
	b[0x03] = byte(w.Clock.Hour)
	putU16(b, 0x04, w.Clock.Month)
	putU16(b, 0x06, w.Clock.Year)
	if w.Player >= 0 && w.Player < numFactions {
		putU16(b, playerPtrOffset, w.Player*factionSize)
		b[playerOffset] = byte(w.Player)
	}
	b[trustOffset] = byte(clampU8(w.Trust))
	for i, e := range w.events {
		off := eventQueueOffset + i*eventQueueEntrySize
		putU16(b, off, int(e.Code))
		putU16(b, off+2, int(e.Param))
	}

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
		r[0x16] = byte(f.ReliefSite)
		r[0x17] = byte(f.LostSite)
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
		// 這四個欄位是執行期會動的（威脅掃描與據點換手），
		// 所以存檔要帶著走；沒跑過 tick 的話寫回來的就是讀進去的值。
		// 位元 4／5 沒解，原樣保留（改寫不是重建）。
		flags := byte(0)
		if c.Threatened {
			flags |= 0x80
		}
		if c.Specific {
			flags |= 0x40
		}
		r[0x00] = r[0x00]&0x30 | flags | byte(c.Adjacency&0x0F)
		r[0x14] = byte(c.Threat)
		r[0x18] = byte(c.Occupancy)
		r[0x17] = byte(c.ReliefCooldown)
		r[0x1B] = byte(c.EnemyNeighbours)
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
func (w *World) tickCity(rng economy.Rand) ([]StrategyEvent, []TalkNotice) {
	if len(w.Cities) == 0 {
		return nil, nil
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
	// 原版 sub_13EFD 在 sub_14194 之後無條件呼叫 sub_14269；
	// 事件 11／12 寫入的 +0x15 marker 不是只有畫面效果。
	w.applyCityDisasterEffect(w.cityCursor)
	// 威脅偵測（原版 sub_13EFD 的佔用圖抄寫 ＋ sub_13F74 → sub_13FA9）。
	// 中立據點只更新佔用數與鄰接遮罩，不做威脅判斷——
	// 原版的 `cmp byte ptr [si+841h], 18h / jz` 只跳過 sub_13F74。
	notices := w.refreshCityThreat(w.cityCursor, rng)
	if w.strategicAI && c.Owner >= 0 && c.Owner < numFactions && c.Owner != w.Player {
		if ev := w.formAICorps(c.Owner); ev != nil {
			return []StrategyEvent{*ev}, notices
		}
	}
	return nil, notices
}

// applyCityDisasterEffect 重現原版 sub_14269（IDA 線性位址 00014269）。
//
// marker 是事件 11／12 寫進據點 runtime record +0x15 的 byte：先從 +0x11
// 防災值扣 marker；若不足，差額再依原版的 byte／word 算術扣 +0x10 的
// 上昇值存值、+0x0E 的生產力與 +0x13 的城兵。這裡保留生產力的 16 位元
// 減法，不把未證實的「飽和為零」套到原版沒有夾住的那一欄。
func (w *World) applyCityDisasterEffect(cityID int) {
	if cityID < 0 || cityID >= len(w.Cities) {
		return
	}
	marker := int(w.disasterMarkerLevels[cityID])
	if marker == 0 {
		return
	}
	c := &w.Cities[cityID]
	if c.Prevention >= marker {
		c.Prevention -= marker
		return
	}

	deficit := marker - c.Prevention
	c.Prevention = 0

	// 原版 +0x10 是存值（實際上昇值 + 100），sub 指令不足時寫 0。
	storedGrowth := c.Growth + 100
	if storedGrowth < deficit {
		storedGrowth = 0
	} else {
		storedGrowth -= deficit
	}
	c.Growth = storedGrowth - 100

	// +0x0F 是生產力 u16 的高 byte；原版 `mul byte ptr [si+84Fh]`
	// 後右移兩位，再對 +0x0E 做 16 位元 sub。
	productionLoss := (deficit * (c.Production >> 8)) >> 2
	c.Production = int(uint16(uint16(c.Production) - uint16(productionLoss)))

	// 原版把差額右移一位後對 +0x13 做 byte sub，不足即寫 0。
	garrisonLoss := deficit >> 1
	if c.Garrison < garrisonLoss {
		c.Garrison = 0
	} else {
		c.Garrison -= garrisonLoss
	}
}
