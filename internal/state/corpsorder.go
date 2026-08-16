package state

// 行軍指示的三選一與軍團抵達時的狀態機。
//
// 規格 `docs/spec/39`，機器碼出處 `docs/re/45`（流程與寫入值）與
// `docs/re/64`（`sub_14325` 的 12 筆分派表、解散的五個動作）。
//
// ⚠ **委任只影響遭遇，不影響行軍**：原版 `sub_14E5C` 在委任位元有設時
// 退回自動判定，而抵達處理完全不看那個位元。

import (
	"fmt"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
)

// 軍團記錄 `+0x23` 的 Stage（`docs/re/64` §2）。八個值全部實作：
// 0–3 分玩家與非玩家兩套 handler，8–11 不分（`docs/re/65`）。
const (
	// StageNormal 是一般行軍。戰鬥指揮與委任都寫這個值。
	StageNormal = 0
	// StageDone 是補完兵之後的值（原版 `sub_14499` 的 `[si+23h] = 3`）。
	// 對 AI 來說它是「體檢」：六槽不足額就解體，否則等士氣。
	StageDone = 3
	// StageWaitMorale 是「等士氣」：士氣回到勢力基準才繼續行動。
	StageWaitMorale = 8
	// StageResupply 是「在首都補兵」：抵達處理會退回再重分，然後轉 StageDone。
	StageResupply = 9
	// StageHomeResupply 是「回首都補兵的路上」：先把目標校正成首都。
	StageHomeResupply = 10
	// StageDisband 是「解體」：先回首都，到了就解散。
	StageDisband = 11
)

// MarchMode 是行軍指示的第二段——選完目的地之後那個三選一
// （TALK #76「　戰鬥指揮　」「　委　　任　」「　解　　體　」）。
type MarchMode int

const (
	// MarchCommand 是戰鬥指揮：遭遇時跳選單問玩家。
	MarchCommand MarchMode = iota
	// MarchDelegate 是委任：交給電腦指揮，遭遇時直接自動判定。
	MarchDelegate
	// MarchDisband 是解體：走回首都後把兵還回預備兵池。
	MarchDisband
)

func (m MarchMode) String() string {
	switch m {
	case MarchDelegate:
		return "委任"
	case MarchDisband:
		return "解體"
	default:
		return "戰鬥指揮"
	}
}

// DisbandAllowed 回報這支軍團的三選一要不要出現第三項。
//
// **原版只在目標據點就是自己的首都時才給「解體」**（`sub_17FDB` 的
// `cmp [bx+3], [bp+0]`）——解散要把兵還回預備兵池，而池子在首都。
func (w *World) DisbandAllowed(corps int) bool {
	if corps < 0 || corps >= numCorps || !w.Corps[corps].Alive {
		return false
	}
	c := w.Corps[corps]
	if c.Faction < 0 || c.Faction >= numFactions {
		return false
	}
	return c.Ordered == w.Factions[c.Faction].Capital
}

// SetMarchMode 寫下行軍指示的第二段。**要在 March 之後呼叫**——
// 「解體」能不能選是看 `Ordered`（原版也是先選據點再開選單）。
func (w *World) SetMarchMode(corps int, mode MarchMode) error {
	if corps < 0 || corps >= numCorps || !w.Corps[corps].Alive {
		return fmt.Errorf("state: 軍團 %d 不存在", corps)
	}
	c := &w.Corps[corps]
	switch mode {
	case MarchCommand:
		c.Delegated, c.Stage = false, StageNormal
	case MarchDelegate:
		c.Delegated, c.Stage = true, StageNormal
	case MarchDisband:
		if !w.DisbandAllowed(corps) {
			return fmt.Errorf("state: 解體只能在目標是首都時下達")
		}
		c.Stage = StageDisband
	default:
		return fmt.Errorf("state: 沒有這個行軍指示 %d", mode)
	}
	// 三條路的共同尾巴：計時器寫 1，下一個 tick 就動（`sub_17FDB`）。
	c.Timer = 1
	return nil
}

// arriveCorps 是原版 `sub_14325`：軍團停在目標據點上時，
// 每個 tick 依 Stage 分派一次。
//
// ⚠ **不是「剛抵達的那一刻」**：原版 `sub_12662` 一開頭就比較現在節點與
// 目標節點，相同就直接走分派，所以停著不動的軍團每個 tick 都會跑一次。
// 解體如果下在「軍團已經在首都」的情況，靠的正是這一點。
//
// 分派表 12 筆：**Stage ≥ 8 不分玩家**，0–3 才分玩家與非玩家兩套。
func (w *World) arriveCorps(i int, rng rander) {
	c := &w.Corps[i]
	switch c.Stage {
	case StageDisband:
		w.arriveDisband(i)
		return
	case StageResupply:
		w.resupplyCorps(i)
		return
	case StageWaitMorale:
		w.waitMorale(i)
		return
	case StageHomeResupply:
		w.headHomeResupply(i)
		return
	}
	if c.Faction != w.Player {
		w.aiArrive(i, rng)
		return
	}
	// 玩家的 Stage 0–3 走同一支（`sub_14370`）：歸零，然後看要不要補兵。
	c.Stage = StageNormal
	if c.Faction < 0 || c.Faction >= numFactions {
		return
	}
	// **兵力滿編就不補**：0x258 ＝ 600 點 ＝ 六槽各 100。
	if c.Men >= army60000Points {
		return
	}
	if c.Node == w.Factions[c.Faction].Capital {
		c.Stage = StageResupply
	}
}

// army60000Points 是六槽滿編的點數（原版 `cmp word [si+4], 258h`）。
// 一點 10 人，600 點 ＝ 6,000 人。
const army60000Points = 600

// arriveDisband 是 `sub_144D6`：目標不是首都就先校正成首都，
// 已經在首都就解散。
func (w *World) arriveDisband(i int) {
	c := &w.Corps[i]
	if c.Faction < 0 || c.Faction >= numFactions {
		return
	}
	capital := w.clampCity(w.Factions[c.Faction].Capital)
	if c.Ordered != capital {
		// 原版還會設 `+0x00` 位元 1 ＝「下一步要重算」，由 `sub_12662`
		// 清掉並呼叫 `sub_147BB` 重查道路表（`docs/re/64` §6）。
		// remake 的 `March` 一次算完整條路徑，重下一次就等價。
		_ = w.March(i, capital)
		return
	}
	if c.Node != capital {
		return // 還在路上
	}
	w.disbandCorps(i)
}

// resupplyCorps 是 `sub_14499`：六槽退回池 → 重新分配 → 轉 StageDone。
//
// **走的是編成畫面同一條路**（`sub_1461D` → `sub_14717` ＋ `sub_14698`），
// 所以補兵的分配規則與編成完全一致。
func (w *World) resupplyCorps(i int) {
	c := &w.Corps[i]
	if c.Faction < 0 || c.Faction >= numFactions {
		return
	}
	f := &w.Factions[c.Faction]
	kinds, manned := c.slots()
	w.poolBack(i)
	men := distributeReserves(&f.Reserves, kinds, manned)
	c.Men = 0
	for k := range c.Units {
		if !manned[k] || men[k] == 0 {
			c.Units[k] = combat.Unit{Kind: EmptySlotKind}
			continue
		}
		c.Units[k] = combat.Unit{Men: men[k], Kind: kinds[k]}
		c.Men += men[k]
	}
	c.Stage = StageDone
}

// slots 把六個槽拆成 distributeReserves 要的兩個陣列。
// **空槽看的是兵種欄**（原版的 4），不是兵力 0——打到剩 0 的槽仍要補。
func (c Corps) slots() (kinds [army.Positions]army.TroopType,
	manned [army.Positions]bool) {

	for k, u := range c.Units {
		kinds[k] = u.Kind
		manned[k] = u.Kind != EmptySlotKind
	}
	return kinds, manned
}

// poolBack 是 `sub_14717`：六個槽的兵員全部退回勢力的預備兵池。
func (w *World) poolBack(i int) {
	c := &w.Corps[i]
	f := &w.Factions[c.Faction]
	for k, u := range c.Units {
		if u.Kind == EmptySlotKind || u.Men == 0 {
			continue
		}
		if t := int(u.Kind); t >= 0 && t < len(f.Reserves) {
			f.Reserves[t] += u.Men
		}
		c.Units[k] = combat.Unit{Kind: u.Kind}
	}
	c.Men = 0
}

// disbandCorps 是 `sub_14651` 的五個動作（`docs/re/64` §3）。
func (w *World) disbandCorps(i int) {
	c := &w.Corps[i]
	if !c.Alive {
		return
	}
	w.poolBack(i)                  // ② 兵回預備兵池
	if c.Faction >= 0 && c.Faction < numFactions {
		f := &w.Factions[c.Faction]
		if f.Corps > 0 {
			f.Corps-- // ① 勢力的軍團數 −1
		}
	}
	if i >= 0 && i < len(w.Generals) {
		w.Generals[i].Posted = false // ④ 主將解職（軍團編號 ＝ 主將編號）
	}
	c.Alive = false // ③ 軍團記錄歸零
	c.Stage = StageNormal
	// ⑤ 大地圖佔用圖 −1：remake 的佔用是每 tick 由位置推導的，
	// 沒有要維護的計數器（`docs/re/44` §1 記的那張表是快取）。
}
