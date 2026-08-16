package state

import (
	"fmt"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/capital"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/march"
)

// 軍團表：127 筆 × 64 B，區塊 `+0x22C0`（段內 `2240h`）。
// 出貨的劇本檔裡全零——開局沒有軍團，玩家要自己編成。
// 佈局見 docs/formats/08 §1.7，來源是 `sub_16F26`／`sub_16F86`／`sub_125A3`。
const (
	corpsBase, corpsSize, numCorps = 0x22C0, 64, 127

	// unitSlots 起點：軍團記錄 +0x28 起是六個部隊槽，每槽 4 B
	// （`sub_15285` 的 `add si, 28h` ＋ `add si, 4`，docs/re/09 §3.1）。
	unitSlotBase, unitSlotSize = 0x28, 4

	// aliveFlag 是存在旗標的門檻。編成時寫 0xC0，掃描時比 `cmp byte ptr [si], 80h`。
	aliveFlag = 0x80
	newCorps  = 0xC0
)

// Corps 是一支軍團。
//
// **軍團編號與武將編號一一對應**：`sub_1291A` 直接用
// `(軍團位址 − 0x2240) ÷ 2 + 0x4240` 換算，兩張表同索引平行
// （docs/re/09 §6）。所以這裡不另外存「將領」——索引就是將領。
type Corps struct {
	Alive   bool
	Faction int // +0x01
	Morale  int // +0x06，編成時從勢力的士氣基準複製
	Men     int // +0x04，總兵力（六個槽的和）

	// Units 是六個編成位置的兵種與兵力。空槽的 Men 是 0。
	// 一點兵力 ＝ 10 人，滿編 100 ＝ 1,000 人（說明書 5.5）。
	Units [army.Positions]combat.Unit

	// Heading 是**朝向**（+0x08）：0／1 是東西、2／3 是南北、4 是靜止。
	// `sub_12808` 從「現在座標與下一個路徑點的差」算出來，X 有差就用 X、
	// 沒差才看 Y。野戰要取樣大地圖上的哪兩格由它決定（`sub_14B63`）。
	Heading int

	Direction int // +0x0A，沿路徑表前進的**步進量**（`sub_127F6` 取負再相加）
	Timer     int // +0x0B，每 tick 減 1，歸零走一步
	Interval  int // +0x1E，速度 ＝ 間隔的倒數

	Node, X, Y                   int // +0x0E（÷8）／+0x10／+0x12
	TargetNode, TargetX, TargetY int // +0x14（÷8）／+0x16／+0x18

	// Ordered 是**玩家下令的目標據點**（+0x20），TargetNode 是移動用的
	// 同一個值（+0x14，原版存 ×8）。原版把同一個概念存成兩份：
	// `sub_142AB` 一次寫兩個，`mov [si+14h],bx` 之後 `shr bx,1` ×3 再
	// `mov [si+20h],bl`，所以 +0x20 恆等於 +0x14 ÷ 8。
	//
	// 兩份仍分開留著，因為遷都會讓它們暫時不一致（`sub_14502`：目標是
	// 舊首都的改成新首都，但 +0x14 只在它等於新首都×8 時才改）。
	// 合併成一個欄位會讓那一段自我抵銷。
	Ordered int

	// Delegated 是記錄 +0x00 的位元 2 ＝ **「委任」**（交給電腦指揮）。
	//
	// 玩家在軍團面板下完行軍目標之後，選單「戰鬥指揮／委任／解體」
	// （TALK #76）決定它：選「戰鬥指揮」清掉、選「委任」設起。
	// AI 自己編的軍團預設是委任（`sub_16E8F`），而**君主親征那一支會被清掉**
	// （`sub_1699E` 編完立刻 `and [di], 0FBh`）。
	//
	// 求援只調得動委任中的軍團（`sub_14155` 的 `test byte [di], 4`）——
	// **玩家自己指揮的軍團不會被 AI 搶走**。見 docs/re/45。
	Delegated bool

	// Stage 是記錄 +0x23。下完指令歸 0，選「解體」寫 11；
	// 而求援要求 < 8，所以待解體的軍團調不動（docs/re/45 §3）。
	Stage int

	// Routing 是記錄 +0x00 的旗標 **8**：**敗走中**。
	// 回首都的路要穿過別人的地時進這個狀態（docs/spec/43）。
	//
	// ⚠ 旗標 8 < 0x80，所以這支軍團**不算活著**——地圖上不畫、
	// 勢力的軍團數已經減掉；但它還沒消失，`sub_12A7E` 每 tick 處理它。
	Routing bool
	// RoutTimer 是記錄 +0x03：敗走的倒數，進狀態時寫 48，
	// 每 tick 減 1，歸零時軍團記錄歸零、主將解職。
	RoutTimer int
}

// ⚠ **行軍路線刻意不放在 Corps 裡**，放在 `World.routes`。
//
// 兩個理由。一、原版的軍團記錄沒有這個欄位（原版玩家只能對相鄰據點
// 下令，不需要多段路由），放進來會讓「記錄長什麼樣」與「remake 加了什麼」
// 混在一起。二、**Corps 一旦含有 slice 就不能用 `==` 比**，
// 而存檔的 byte-for-byte round-trip 測試正是靠這個可比性。

// Leader 回傳帶兵的武將編號。軍團與武將同索引，所以就是軍團編號。
func (w *World) Leader(corps int) int { return corps }

func (w *World) loadCorps(b []byte) {
	for i := range w.Corps {
		r := b[corpsBase+i*corpsSize:]
		c := Corps{
			Alive:      r[0x00] >= aliveFlag,
			Faction:    int(r[0x01]),
			Men:        u16(r, 0x04),
			Morale:     int(r[0x06]),
			Heading:    int(r[0x08]),
			Direction:  int(r[0x0A]),
			Timer:      int(r[0x0B]),
			Node:       u16(r, 0x0E) / 8,
			X:          u16(r, 0x10),
			Y:          u16(r, 0x12),
			TargetNode: u16(r, 0x14) / 8,
			TargetX:    u16(r, 0x16),
			TargetY:    u16(r, 0x18),
			Interval:   int(r[0x1E]),
			Ordered:   int(r[0x20]),
			Delegated: r[0x00]&0x04 != 0,
			Stage:     int(r[0x23]),
			// 旗標 8 而且不到 0x80 ＝ 敗走中（docs/spec/43）。
			Routing:   r[0x00] < aliveFlag && r[0x00]&0x08 != 0,
			RoutTimer: int(r[0x03]),
		}
		for k := range c.Units {
			s := r[unitSlotBase+k*unitSlotSize:]
			c.Units[k] = combat.Unit{Men: int(s[1]), Kind: kindFromByte(s[2])}
		}
		w.Corps[i] = c
	}
}

func (w *World) saveCorps(b []byte) {
	for i, c := range w.Corps {
		r := b[corpsBase+i*corpsSize:]
		if !c.Alive {
			// 敗走中的軍團要把旗標 8 與倒數寫回去，其餘欄位不動
			// （原版 `sub_12977` 也只改這兩個，docs/spec/43）。
			if c.Routing {
				r[0x00] = 0x08
				r[0x03] = byte(c.RoutTimer)
			}
			// 其餘不存在的軍團**一個 byte 都不動**——重建會抹掉痕跡。
			continue
		}
		r[0x00] = newCorps
		if c.Delegated {
			r[0x00] |= 0x04
		} else {
			r[0x00] &^= 0x04
		}
		r[0x01] = byte(c.Faction)
		r[0x02] = byte(i)
		putU16(r, 0x04, c.Men)
		r[0x06] = byte(c.Morale)
		r[0x08] = byte(c.Heading)
		r[0x0A] = byte(c.Direction)
		r[0x0B] = byte(c.Timer)
		putU16(r, 0x0E, c.Node*8)
		putU16(r, 0x10, c.X)
		putU16(r, 0x12, c.Y)
		putU16(r, 0x14, c.TargetNode*8)
		putU16(r, 0x16, c.TargetX)
		putU16(r, 0x18, c.TargetY)
		r[0x1E] = byte(c.Interval)
		r[0x20] = byte(c.Ordered)
		r[0x23] = byte(c.Stage)
		for k, u := range c.Units {
			s := r[unitSlotBase+k*unitSlotSize:]
			s[1] = byte(u.Men)
			s[2] = byteFromKind(u.Kind)
		}
	}
}

// 兵種在檔案裡是 1-based（`sub_14F8A` 直接寫 3 表示步兵，docs/re/09 §7）。
// EmptySlotKind 是「這個編成位置沒有兵」。原版的兵種欄寫 **4**
// （`sub_14717` 掃到 4 就跳過，`docs/re/30` §4），而 remake 的
// TroopType 是 0-based，所以是 3。
const EmptySlotKind = army.TroopType(3)

func kindFromByte(v byte) army.TroopType {
	if v == 0 {
		return army.Cavalry
	}
	return army.TroopType(v - 1)
}

func byteFromKind(t army.TroopType) byte { return byte(t) + 1 }

// 移動間隔。純騎馬編成走得快（說明書 5.5「騎馬隊のみの軍団は
// 移動速度が速くなります」）。
//
// 數值出自 `sub_16FD2`：掃六個槽，只要有一槽兵種不是騎馬就記一個旗標，
// 最後 `+0x1E` ＝ 全騎馬 2、否則 3（`docs/re/30` §5）。
// **混編一律同速**——多摻一種兵不會更慢。
const (
	IntervalCavalry = 2
	IntervalMixed   = 3
)

// MaxMenPerSlot 是一個編成槽的兵力上限，單位是**點**（一點 10 人）。
// 原版 `sub_14698` 的 `cmp ax, 64h`，而槽位本身也只有 1 byte。
const MaxMenPerSlot = 100

// PreviewFormation 算出「照現在這組兵種按下確定，六個槽會各分到多少兵」，
// **不動任何狀態**——編成畫面每次重畫都要顯示這個結果（`docs/spec/22` §1.2）。
//
// 原版不需要這一支：它每次改兵種就真的把兵退回池再重分，畫面直接讀軍團記錄
// （`docs/re/30` §4.1）。remake 的編成是到按確定才落地，所以預覽要另外算。
func (w *World) PreviewFormation(faction int,
	kinds [army.Positions]army.TroopType, manned [army.Positions]bool) [army.Positions]int {

	if faction < 0 || faction >= len(w.Factions) {
		return [army.Positions]int{}
	}
	pool := w.Factions[faction].Reserves // 值拷貝，分配只動這一份
	return distributeReserves(&pool, kinds, manned)
}

// distributeReserves 照原版 `sub_14698` 把預備兵分給六個槽，並從池裡扣掉
// 實際放進去的量。回傳每個槽分到幾點。
//
// 分配式（docs/spec/21 §2）：同一個兵種佔幾個槽就分成幾份，
// **餘數整個給第一個槽**，之後的槽再對剩下的重分；每槽上限 100 點。
//
// ⚠ 扣掉的量與放進槽裡的量必須是同一個數。原版是
// `sub es:[bx], ax` 之後緊接 `mov [si+1], al`——**同一個 ax**。
// 先前 remake 扣 1000、放 100，等於每編一支軍團就吃掉十倍的池。
func distributeReserves(pool *[economy.NumTroopTypes]int,
	kinds [army.Positions]army.TroopType, manned [army.Positions]bool) [army.Positions]int {

	var left [economy.NumTroopTypes]int
	for k, ok := range manned {
		if ok && int(kinds[k]) >= 0 && int(kinds[k]) < int(economy.NumTroopTypes) {
			left[kinds[k]]++
		}
	}
	var out [army.Positions]int
	for k, ok := range manned {
		if !ok {
			continue
		}
		t := int(kinds[k])
		if t < 0 || t >= int(economy.NumTroopTypes) || left[t] == 0 {
			continue
		}
		n := pool[t]/left[t] + pool[t]%left[t]
		left[t]--
		if n > MaxMenPerSlot {
			n = MaxMenPerSlot
		}
		pool[t] -= n
		out[k] = n
	}
	return out
}

// FormCorps 編成一支軍團（原版 `sub_16F26`）。
//
// leader 是帶兵的武將編號，kinds 是六個位置的兵種，manned 標哪幾個位置要有兵。
// **兵力不由呼叫端決定**：照 `sub_14698` 從勢力的預備兵池分配
// （docs/spec/21 §2），池裡有多少就分多少。
//
// 照原版的順序：武將標成出陣中、軍團繼承勢力的士氣基準、
// 位置設在首都、勢力的軍團數 +1。
func (w *World) FormCorps(leader int, kinds [army.Positions]army.TroopType,
	manned [army.Positions]bool) error {

	if leader < 0 || leader >= numCorps {
		return fmt.Errorf("state: 武將編號 %d 超出 0–%d", leader, numCorps-1)
	}
	g := &w.Generals[leader]
	if !g.Alive {
		return fmt.Errorf("state: 武將 %d 不存在", leader)
	}
	if g.Faction < 0 || g.Faction >= numFactions {
		return fmt.Errorf("state: 武將 %d 在野，不能編成", leader)
	}
	if w.Corps[leader].Alive {
		return fmt.Errorf("state: 武將 %d 已經帶著軍團", leader)
	}
	f := &w.Factions[g.Faction]

	// **大將的位置一定要有兵。** 原版的壞滅判定 `sub_1474A` 直接看
	// `[si+29h]`（第一槽的兵力）是不是 0，是就當軍團已經沒了——
	// 所以一支大將空著的軍團一編出來就會被判掉（docs/re/09 §5）。
	if !manned[0] {
		return fmt.Errorf("state: 大將的位置一定要有兵")
	}

	// 兵力由 `sub_14698` 分配，**不是每槽固定 100 點**（docs/spec/21 §2）。
	// 池裡有多少就分多少，所以「兵不夠」不是錯誤——分完主將槽是 0 才是。
	men := distributeReserves(&f.Reserves, kinds, manned)
	if men[0] == 0 {
		return fmt.Errorf("state: 大將的位置分不到兵（預備兵 %v）", f.Reserves)
	}

	home := w.clampCity(f.Capital)
	c := Corps{
		Alive:   true,
		Faction: g.Faction,
		Morale:  f.MoraleBase,
		Ordered:    f.Capital,
		Node:    home,
		X:       w.Cities[home].X,
		Y:       w.Cities[home].Y,
		// 目標先設成原地，行軍指令下達前不會動。
		TargetNode: home,
		TargetX:    w.Cities[home].X,
		TargetY:    w.Cities[home].Y,
	}
	allCav := false
	for k, ok := range manned {
		if !ok || men[k] == 0 {
			// **空槽在原版是兵種 4**，不是「兵種 0 而人數 0」
			// （`sub_14717` 看到 4 就跳過，docs/re/30 §4）。
			// 寫回存檔與畫面取圖都靠這個值，不能留成騎馬。
			c.Units[k] = combat.Unit{Kind: EmptySlotKind}
			continue
		}
		c.Units[k] = combat.Unit{Men: men[k], Kind: kinds[k]}
		c.Men += men[k]
	}
	allCav = c.rules().AllCavalry()
	c.Interval = IntervalMixed
	if allCav {
		c.Interval = IntervalCavalry
	}
	// 原版的 `sub_16FD2` 每次重算都把 `+0x0B` 寫成 **1**，不是寫成間隔——
	// 所以剛編成的軍團下一個 tick 就會走第一步（`docs/re/30` §5）。
	c.Timer = 1

	w.Corps[leader] = c
	g.Posted = true
	f.Corps++
	return nil
}

// rules 回傳這支軍團在 army 層的視圖。
func (c Corps) rules() army.Corps {
	out := army.Corps{Alive: c.Alive, Faction: c.Faction, Morale: c.Morale,
		Node: c.Node, X: c.X, Y: c.Y,
		TargetNode: c.TargetNode, TargetX: c.TargetX, TargetY: c.TargetY,
		Direction: c.Direction, MoveTimer: c.Timer, MoveInterval: c.Interval}
	for k, u := range c.Units {
		out.Units[k] = u.Kind
		out.Manned[k] = u.Men > 0
	}
	return out
}

// battle 回傳這支軍團在 combat 層的視圖。
func (w *World) battle(i int) combat.Corps {
	c := w.Corps[i]
	g := w.Generals[i]
	return combat.Corps{
		Faction: c.Faction,
		Leader: combat.Leader{
			Martial: g.Martial, Command: g.Command,
			SiegeAptitude: g.Aptitude[0], FieldAptitude: g.Aptitude[1],
			Rating: g.Rules().Rating(),
		},
		Units: c.Units, Morale: c.Morale, Men: c.Men,
	}
}

func (w *World) applyBattle(i int, b combat.Corps) {
	c := &w.Corps[i]
	c.Units, c.Morale, c.Men = b.Units, b.Morale, b.Men
}

// ---------------------------------------------------------------------------
// 每 tick 的軍團更新（`sub_125A3`）
// ---------------------------------------------------------------------------

// corpsPerTick 是原版每個 tick 處理幾支軍團。
//
// ⭐ **不是全部 127 支**：`sub_125A3` 的 `mov cx, 10h` 只跑 16 筆，
// 從一個游標開始，處理完把游標往前推，`si >= 0x1FC0`（127 × 64）繞回 0。
// 所以軍團是**輪流**被更新的，一輪要 8 個 tick。
const corpsPerTick = 16

// upkeepHour 是收軍費與回士氣的時刻。
//
// ⚠ `sub_12600` 開頭就是 `cmp cs:byte_10CF3, 1 / jz`，而 `ds:0CF3h` 是
// **小時**（`sub_11D8E` 在 `0x17` ＝ 23 進位）。所以軍費不是每 tick 收，
// 是**每天「一時」那個小時收**。docs/re/09 §9 初版寫成「每 tick」，
// 那是只看 `sub_125A3` 的呼叫點、沒往下讀 `sub_12600` 的閘。
const upkeepHour = 1

// CorpsEvent 是一支軍團在這個 tick 發生的事。
type CorpsEvent struct {
	Corps int

	Moved   bool // 這個 tick 走了一步
	Arrived bool // 到達目標

	// Battle 不是 nil 表示打了一場。Enemy 是對手的軍團編號，
	// −1 表示對手是據點的城兵。Mode 是野戰還是攻城。
	Battle *combat.Result
	Enemy  int
	Mode   combat.Mode

	// BattleBefore／BattleAfter 是戰略層的兵力點數（每點 10 人），
	// 給結果視窗與事件記錄使用。它們不是原版存檔欄位，也不參與規則；
	// 只是把戰鬥前後已存在的數值沿事件流帶出來，避免 UI 重新猜測。
	BattleBefore     [2]int
	BattleAfter      [2]int
	BattleCityDamage int

	// Destroyed 是這一戰壞滅的軍團編號（可能兩支都是）。
	Destroyed []int
	// Fate 是壞滅方主將的下場，只在 Destroyed 非空時有意義。
	Fate map[int]combat.Fate

	// Captured 不是 −1 表示這個 tick 佔下了某個據點。
	Captured int

	// Disbanded 表示這支軍團在這個 tick 解體了（`docs/spec/39`）——
	// 兵**回**預備兵池。
	Disbanded bool

	// Routed 表示這支軍團在這個 tick 敗走了（`docs/spec/43`）——
	// 兵**不回**池。兩者都會讓軍團從地圖上消失，但代價完全不同，
	// 所以計數要分開，不然量出來的「軍團損耗」會把回收算成損失。
	Routed bool

	// GovernorReturned 不是 −1 表示**那個據點派駐的內政官被遣回了**
	// （原版 `sub_14D63`，docs/spec/48），值是武將編號。
	GovernorReturned int

	// Relocated 不是 −1 表示**舊主的首都被打下來，遷到了這個據點**
	// （原版 `sub_14DF0`，訊息 30「首都被攻陷了！儘速遷都到\2」）。
	Relocated int
}

// tickCorps 跑一輪軍團更新，回傳這個 tick 發生的事。
func (w *World) tickCorps(hour int, rng combat.Rand) []CorpsEvent {
	var out []CorpsEvent
	for n := 0; n < corpsPerTick; n++ {
		i := w.corpsCursor
		w.corpsCursor = (w.corpsCursor + 1) % numCorps
		if !w.Corps[i].Alive {
			// 敗走中的軍團不算活著，但還有一個倒數要跑
			// （原版 `sub_125A3` 的 `test byte [si+2240h], 8`）。
			if w.Corps[i].Routing {
				w.tickRout(i)
			}
			continue
		}
		if ev := w.tickOneCorps(i, hour, rng); ev != nil {
			out = append(out, *ev)
		}
	}
	return out
}

func (w *World) tickOneCorps(i, hour int, rng combat.Rand) *CorpsEvent {
	c := &w.Corps[i]
	ev := CorpsEvent{Corps: i, Enemy: -1, Captured: -1,
		Relocated: capital.None, GovernorReturned: noGovernor}

	// ① 移動的節拍。原版先減再判斷：間隔 N 表示每 N 個 tick 走一步。
	c.Timer--
	if c.Timer <= 0 {
		c.Timer = c.Interval
		// ⚠ **停在目標上也要跑抵達處理**：原版 `sub_12662` 一開頭就比
		// 「現在節點 ＝ 目標節點」，相同就直接呼叫 `sub_14325` 分派，
		// 不需要移動（`docs/re/64` §1）。解體下在「已經在首都」時就靠這條。
		// ⚠ 判「到了」要連座標一起看：remake 的 `Node` 在**踩到據點座標**
		// 時就更新（中繼據點也算），單看它會把「還在路上但經過目標據點」
		// 誤判成抵達。
		if c.Node == c.TargetNode && c.X == c.TargetX && c.Y == c.TargetY {
			w.arriveCorps(i, rng)
			if !c.Alive {
				ev.Disbanded, ev.Routed = !c.Routing, c.Routing
				return &ev
			}
		} else if w.step(i) {
			ev.Moved = true
			ev.Arrived = c.Node == c.TargetNode
			if ev.Arrived {
				w.arriveCorps(i, rng)
				if !c.Alive {
					ev.Disbanded, ev.Routed = !c.Routing, c.Routing
					return &ev
				}
			}
			w.resolveContact(i, &ev, rng)
		}
	}

	// ② 軍費與士氣。每天「一時」那個小時才收。
	// 這一步壞滅的軍團不收——它已經不在了。
	if hour == upkeepHour && c.Alive {
		inField := army.KindOf(c.Node) == army.FieldNode
		f := &w.Factions[c.Faction]
		f.Expense = economy.ClampFunds(f.Expense + combat.Upkeep(c.Men, inField))
		cc := combat.Corps{Morale: c.Morale}
		combat.Recover(&cc, f.MoraleBase, inField)
		c.Morale = cc.Morale
	}

	if !ev.Moved && ev.Battle == nil {
		return nil
	}
	return &ev
}

// step 把軍團往目標推進一格，回傳有沒有真的動。
//
// **逐格走在原版的道路上。** 路徑是 `internal/assets/world` 用
// 原版的走訪常式（`sub_1E81C`／`sub_1E961`）算出來的，
// 不是最短路——原版照著畫出來的路走，會繞。
//
// 沒有道路圖時（缺原版素材）`routes` 是空的，退回直線逼近：
// **缺素材要能降級跑，不是整個動不了。**
func (w *World) step(i int) bool {
	c := &w.Corps[i]

	// ① 有格子路徑就逐格走 —— **這條路徑每一格都踩在道路圖塊上**
	//    （`internal/assets/world` 有逐格檢查的測試）。
	if cells := w.routes[i]; len(cells) > 0 {
		next := cells[0]
		w.routes[i] = cells[1:]
		c.Heading = headingTo(c.X, c.Y, next[0], next[1])
		c.X, c.Y = next[0], next[1]
		// 踩到某個據點的座標就算抵達那個據點。中繼據點也要更新，
		// 不然攻城、遭遇這些判定會在錯的地方觸發。
		if n := w.cityAt(next[0], next[1]); n >= 0 {
			c.Node = n
		}
		if len(w.routes[i]) == 0 {
			c.Node = c.TargetNode
			c.Heading = HeadingStill
		}
		return true
	}

	// ② 沒有路徑（缺素材、或圖裡這一段沒有格子序列）→ 退回直線。
	if c.X == c.TargetX && c.Y == c.TargetY {
		if c.Node != c.TargetNode {
			c.Node = c.TargetNode
			c.Heading = HeadingStill
			return true
		}
		c.Heading = HeadingStill
		return false
	}
	c.Heading = headingTo(c.X, c.Y, c.TargetX, c.TargetY)
	c.X += sign(c.TargetX - c.X)
	c.Y += sign(c.TargetY - c.Y)
	if c.X == c.TargetX && c.Y == c.TargetY {
		c.Node = c.TargetNode
		c.Heading = HeadingStill
	}
	return true
}

// cityAt 回傳座標上的據點編號，沒有回 −1。
//
// 用線性掃描是刻意的：192 筆而已，而建索引就要面對「兩個據點同座標」
// 這種原版資料可能有的狀況。線性掃描回第一個，行為明確。
func (w *World) cityAt(x, y int) int {
	for i := range w.Cities {
		if w.Cities[i].X == x && w.Cities[i].Y == y {
			return i
		}
	}
	return -1
}

// 朝向的四個值加上「靜止」。原版寫進軍團記錄 `+0x08`
// （`sub_12808`；到站時 `sub_12662`／`sub_127A2` 改寫成 4）。
//
// 編碼是**符號位元**來的：`ax = 現在 − 目標`，取 bit 15 轉成 0／1，
// 南北那一組再加 2。所以 0／1 是「X 減少／增加」，2／3 是「Y 減少／增加」。
const (
	HeadingXMinus = 0
	HeadingXPlus  = 1
	HeadingYMinus = 2
	HeadingYPlus  = 3
	HeadingStill  = 4
)

// headingTo 重現 `sub_12808`：**X 有差就只看 X，X 相同才看 Y**。
func headingTo(x, y, tx, ty int) int {
	if d := tx - x; d != 0 {
		if d < 0 {
			return HeadingXMinus
		}
		return HeadingXPlus
	}
	if d := ty - y; d != 0 {
		if d < 0 {
			return HeadingYMinus
		}
		return HeadingYPlus
	}
	return HeadingStill
}

func sign(v int) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	}
	return 0
}

// resolveContact 檢查這一步之後有沒有撞上敵人或走進別人的據點。
//
// 兩條路對應原版的 `sub_12831`（野戰遭遇）與 `sub_12880`（攻城），
// 判定條件照原版：**同格且不同勢力**才打，走進自家據點直接通過。
func (w *World) resolveContact(i int, ev *CorpsEvent, rng combat.Rand) {
	c := w.Corps[i]

	// ⭐ **據點要先問，野戰後問。**
	//
	// 原版的順序反過來（`sub_12708` 先看佔用圖 `cmp byte ptr [di], 0`，
	// 有人才叫 `sub_12831` 打野戰，沒人而且是據點圖塊才叫 `sub_12880`
	// 打攻城），但那是建立在**一個據點佔好幾格地圖**上的：
	// 據點的圖塊值是 `0xCE`–`0xDD` 一整段，守軍站在其中一格，
	// 攻方通常踏進的是**別的那幾格**——那幾格佔用圖是 0，所以走攻城，
	// 接著 `sub_14C72` 再用**據點自己的座標**把守軍找出來。
	//
	// 本專案的據點是**一個點**，攻方必然踏在守軍那一格上，
	// 照抄順序的話永遠打成野戰、攻城那條路永遠走不到。
	// 所以這裡把順序倒過來——**這是為了補上地圖模型的差異，
	// 不是規則不同**（docs/re/09 §2）。
	if army.KindOf(c.Node) == army.CityNode {
		city := &w.Cities[c.Node]
		switch {
		case city.Owner == combat.NeutralFaction:
			// 中立據點沒有主人，所以沒有「首都失守」這回事。
			city.Owner = c.Faction
			w.Factions[c.Faction].Cities++
			ev.Captured = c.Node
			return
		case city.Owner != c.Faction:
			// 城裡有守軍就打守軍，沒有就打城兵。
			for j := range w.Corps {
				d := w.Corps[j]
				if d.Alive && d.Faction == city.Owner && d.Node == c.Node {
					w.fight(i, j, ev, combat.Siege, city.Garrison, rng)
					return
				}
			}
			w.fightGarrison(i, ev, rng)
			return
		}
		// 自家的據點：不打，繼續往下看有沒有敵軍團同格（不該發生，但不擋）。
	}

	// 野戰：同一格上有別的勢力的軍團。
	for j := range w.Corps {
		d := w.Corps[j]
		if j == i || !d.Alive || d.Faction == c.Faction {
			continue
		}
		if d.X == c.X && d.Y == c.Y {
			w.fight(i, j, ev, combat.Field, 0, rng)
			return
		}
	}
}

func (w *World) fight(att, def int, ev *CorpsEvent, m combat.Mode, garrison int, rng combat.Rand) {
	// ⭐ 玩家的勢力捲進去先問戰鬥指揮／委任，其餘自動判定（原版 `sub_14E5C`）。
	ev.Mode = m
	if w.wantsTactical(att, def) {
		// 先確認這場的地形來源確實能建出來；沒有對應戰場時照原版的
		// 委任退路，不讓玩家卡在一個永遠選不了的空選單。
		node := w.Corps[att].Node
		if w.tactical.Field(node, m == combat.Siege) != nil {
			w.encounter = &EncounterChoice{
				Attacker: att, Defender: def, Node: node,
				Mode: m, Garrison: garrison,
			}
			return
		}
	}
	w.resolveCorpsBattle(ev, att, def, m, garrison, rng)
}

// resolveCorpsBattle 執行一場已決定委任的軍團對軍團戰鬥。
// 戰鬥指揮的出口走 ResolvePending；兩者最後共用同一組戰後處理。
func (w *World) resolveCorpsBattle(ev *CorpsEvent, att, def int, m combat.Mode, garrison int, rng combat.Rand) {
	a, d := w.battle(att), w.battle(def)
	ev.BattleBefore = [2]int{a.Men, d.Men}
	r := combat.Resolve(&a, &d, m, garrison, rng)
	w.applyBattle(att, a)
	w.applyBattle(def, d)
	ev.Battle, ev.Enemy, ev.Mode = &r, def, m
	ev.BattleAfter = [2]int{a.Men, d.Men}
	ev.BattleCityDamage = r.CityDamage
	w.damageCity(w.Corps[att].Node, m, r)

	// 原版兩邊各跑一次 `sub_1474A`：士氣判定之外，**敗方退不了也算壞滅**
	// （docs/spec/46 §1）。守方站在自家城裡走「不退」那一支，
	// 所以攻城的易主判定不受影響。
	attDead := r.AttackerDestroyed || w.retreatOrPerish(att, !r.DefenderWins)
	defDead := r.DefenderDestroyed || w.retreatOrPerish(def, r.DefenderWins)
	w.afterBattle(ev, att, attDead, def, rng)
	w.afterBattle(ev, def, defDead, att, rng)

	if defDead && !attDead && m == combat.Siege {
		w.capture(att, ev, rng)
	}
}

// fightGarrison 打的是據點的城兵——原版在 `ds:4200h` 現搭一支臨時軍團
// （`sub_14F8A`，docs/re/09 §7）。守方不是軍團，所以不會有壞滅或被擒。
func (w *World) fightGarrison(att int, ev *CorpsEvent, rng combat.Rand) {
	node := w.Corps[att].Node
	city := &w.Cities[node]
	a := w.battle(att)
	g := combat.Garrison(city.Owner, city.Garrison)
	ev.BattleBefore = [2]int{a.Men, g.Men}
	r := combat.Resolve(&a, &g, combat.Siege, city.Garrison, rng)
	w.applyBattle(att, a)
	ev.Battle, ev.Enemy, ev.Mode = &r, -1, combat.Siege
	ev.BattleAfter = [2]int{a.Men, g.Men}
	ev.BattleCityDamage = r.CityDamage
	w.damageCity(node, combat.Siege, r)

	// 守方是城兵不是軍團，所以只有攻方要跑 `sub_1474A`。
	attDead := r.AttackerDestroyed || w.retreatOrPerish(att, !r.DefenderWins)
	w.afterBattle(ev, att, attDead, -1, rng)
	if !r.DefenderWins && !attDead {
		w.capture(att, ev, rng)
	}
}

// damageCity 套用攻城戰對據點的損傷。城兵、上昇值、防災值各扣同一個量，
// **不分勝敗**（`sub_151B3`，docs/re/09 §4.1）。
func (w *World) damageCity(node int, m combat.Mode, r combat.Result) {
	if m != combat.Siege || army.KindOf(node) != army.CityNode {
		return
	}
	c := &w.Cities[node]
	c.Garrison = clampDown(c.Garrison, r.CityDamage)
	c.Prevention = clampDown(c.Prevention, r.CityDamage)
	// 上昇值在記憶體裡是「實際值 ＋ 100」的存值，原版扣的是存值。
	c.Growth = clampDown(c.Growth+100, r.CityDamage) - 100
}

func clampDown(v, d int) int {
	if v -= d; v < 0 {
		return 0
	}
	return v
}

// afterBattle 處理壞滅：軍團消失、主將擲一次下場（`sub_1291A`）。
//
// victor 是勝方的軍團編號，−1 表示勝方是據點的城兵
// （那時勝方勢力就是該據點的所屬）。
func (w *World) afterBattle(ev *CorpsEvent, i int, destroyed bool, victor int, rng combat.Rand) {
	if !destroyed {
		return
	}
	winner := combat.NeutralFaction
	if victor >= 0 {
		winner = w.Corps[victor].Faction
	} else if n := w.Corps[i].Node; army.KindOf(n) == army.CityNode {
		winner = w.Cities[n].Owner
	}
	w.corpsPerishes(ev, i, winner, rng)
}

// corpsPerishes 是「這支軍團沒了」的共同出口：軍團消失、勢力軍團數 −1、
// 主將擲一次下場（原版 `sub_1291A`）。
//
// 兩個入口：戰敗壞滅（`sub_1474A`）與**據點失守後無處可退**
// （`sub_14DA4` 的 `jb` 分支，[`docs/spec/47`](../../docs/spec/47-city-fall-corps-redirect.md)）。
func (w *World) corpsPerishes(ev *CorpsEvent, i, winner int, rng combat.Rand) {
	c := &w.Corps[i]
	loser := c.Faction
	if loser < 0 || loser >= numFactions || i >= len(w.Generals) {
		return
	}

	f := &w.Factions[loser]
	g := &w.Generals[i]
	fate := combat.RollFate(combat.Captive{
		Rating:       g.Rules().Rating(),
		IsRuler:      f.Lord == i,
		HasCapital:   f.Capital != noCity,
		LoyalToDeath: g.LoyalToDeath,
		LordSurvives: f.Alive,
	}, winner, loser, rng)

	c.Alive = false
	g.Posted = false
	if f.Corps > 0 {
		f.Corps--
	}
	switch fate {
	case combat.Captured:
		g.Captor = loser
		g.Faction = winner
		// 原版 `sub_12AD2(al=0FFh, ah=舊勢力)`：只從舊勢力扣，
		// 不加給俘虜方——被俘期間這名武將不算在任何一方的武將數裡。
		w.dropGeneralCount(loser)
	case combat.Suicide:
		g.Alive = false
		g.Faction = noFaction
		w.dropGeneralCount(loser)
	}

	ev.Destroyed = append(ev.Destroyed, i)
	if ev.Fate == nil {
		ev.Fate = map[int]combat.Fate{}
	}
	ev.Fate[i] = fate
}

// noGovernor 是「這個據點沒有派駐內政官」。原版的哨兵是 0xFF，
// 而 remake 的事件層用 −1——**兩個都不是 0**，因為 0 是合法的武將編號
// （`CLAUDE.md` §7 第 11 條）。
const noGovernor = -1

// returnGovernor 是 `sub_14D63`：據點易主時把派駐的內政官遣回。
// 回傳被遣回的武將編號，沒有內政官就回 `noGovernor`。
//
// 只動兩個欄位：據點的內政官槽清成 0xFF、那名武將的「出陣中」歸零。
// **不降職、不處分**——原版就只有這兩行。
func (w *World) returnGovernor(node int) int {
	if node < 0 || node >= len(w.Cities) {
		return noGovernor
	}
	city := &w.Cities[node]
	id := city.Governor
	if id < 0 || id >= len(w.Generals) {
		return noGovernor
	}
	city.Governor = noGovernorSlot
	w.Generals[id].Posted = false
	return id
}

// noGovernorSlot 是據點記錄 `+0x19` 的哨兵（原版寫 `0xFF`）。
// 存檔要 byte-for-byte 寫得回去，所以這裡存的是原始值不是 −1。
const noGovernorSlot = 0xFF

// capture 把據點換手（`sub_14CF3`）。
func (w *World) capture(att int, ev *CorpsEvent, rng combat.Rand) {
	node := w.Corps[att].Node
	if army.KindOf(node) != army.CityNode {
		return
	}
	city := &w.Cities[node]
	old := city.Owner
	next := w.Corps[att].Faction
	if old == next {
		return
	}
	// 原本無主（0x18）就沒有「奪取」，只有新主的據點數 +1。
	if old != combat.NeutralFaction && old >= 0 && old < numFactions {
		if w.Factions[old].Cities > 0 {
			w.Factions[old].Cities--
		}
	}
	city.Owner = next
	city.OwnerRecorded = next
	ev.Captured = node
	// 換旗之後第一件事是把派駐的內政官遣回（`sub_14D63`，docs/spec/48）。
	// **舊主是無主時整段跳過**（原版 `cmp bh, 18h / jz`）。
	if old != combat.NeutralFaction {
		ev.GovernorReturned = w.returnGovernor(node)
	}
	// 原版的順序是**遷都 → 調頭 → 滅亡判定 → 新主據點數 +1**
	// （`sub_14CF3` 逐行）。調頭排在遷都之後不是細節——
	// `sub_1487B` 找的是**新首都**的方向。
	finished := false
	if old >= 0 && old < numFactions && w.Factions[old].Capital == node {
		ev.Relocated = w.relocateCapital(old)
		if ev.Relocated == capital.None {
			// sub_14DF0：首都失守且找不到替代據點時，capital=0xFF
			// 並清除勢力 alive bit；sub_14FCE 隨後對玩家離開主循環。
			w.Factions[old].Capital = noCity
			finished = true
		}
	}
	w.redirectFallenCityCorps(ev, node, old, next, rng)
	if finished {
		w.eliminateFaction(old, next)
	}
	w.Factions[next].Cities++
}

// redirectFallenCityCorps 是 `sub_14DA4`：據點易主之後，**舊主留在那一格上
// 的軍團**逐一改成「回家的下一站」；退不了的走 `sub_1291A`（主將擲下場）。
//
// 名單是 `sub_14C72` 在開打前收的——同一格、同一勢力、還活著的軍團，
// 最多 127 支。所以這一條處理的是**疊在同一格上、沒被捲進那一場的守軍**。
func (w *World) redirectFallenCityCorps(ev *CorpsEvent, node, old, winner int, rng combat.Rand) {
	if old < 0 || old >= numFactions || node < 0 || node >= len(w.Cities) {
		return
	}
	for i := range w.Corps {
		c := &w.Corps[i]
		if !c.Alive || c.Faction != old || c.Node != node {
			continue
		}
		hop := w.nextHopHome(i)
		if hop < 0 {
			w.corpsPerishes(ev, i, winner, rng)
			continue
		}
		_ = w.March(i, hop)
		// 原版還會 `mov byte ptr [si+0Bh], 1` ＋ `or byte ptr [si], 2`：
		// 下一個 tick 就重算並起步。
		c.Timer = 1
	}
}

// March 給軍團下行軍指令：往 node 那個據點走。
//
// 原版的目標是一組三元組（據點編號 × 8、X、Y），行軍就是把現在的那組
// 往目標推。這裡只接受據點——說明書 3.2 的行軍指令也是選據點，
// 野外座標是原版內部推路徑時才用到的。
func (w *World) March(corps, node int) error {
	if corps < 0 || corps >= numCorps || !w.Corps[corps].Alive {
		return fmt.Errorf("state: 軍團 %d 不存在", corps)
	}
	if node < 0 || node >= numCities {
		return fmt.Errorf("state: 據點編號 %d 超出 0–%d", node, numCities-1)
	}
	c := &w.Corps[corps]
	c.TargetNode = node
	// 原版下行軍指令時兩個欄位都寫：`sub_17FDB` 的 `mov [si+20h], al`
	// 寫的就是玩家剛選的目的地（`docs/re/27` §7）。少寫這一個，
	// 存檔的 +0x20 會停在編成時的首都，與原版分歧。
	c.Ordered = node
	c.TargetX, c.TargetY = w.Cities[node].X, w.Cities[node].Y
	w.routes[corps] = nil
	if w.roads == nil || node == c.Node {
		return nil
	}
	path := w.roads.Route(c.Node, node)
	if path == nil {
		// **走不到要說走不到**，不要默默走直線穿過山河。
		c.TargetNode = c.Node
		c.Ordered = c.Node
		c.TargetX, c.TargetY = w.Cities[w.clampCity(c.Node)].X, w.Cities[w.clampCity(c.Node)].Y
		return fmt.Errorf("state: 從 %s 沒有路可以到 %s",
			w.Cities[w.clampCity(c.Node)].Name, w.Cities[node].Name)
	}
	// 有格子序列就用格子序列（沿真正的道路走）；沒有就留空，退回直線。
	w.routes[corps] = w.roads.CellRoute(c.Node, node)
	return nil
}

// SetRoads 掛上道路圖。**規則層不讀檔案**，所以圖由呼叫端
// （`cmd/wlgame` 等）從 `MMAP` 推導後注入。
//
// 沒有掛的話行軍退回直線移動——缺原版素材時要能降級跑，
// 不是整個動不了。
func (w *World) SetRoads(g *march.Graph) { w.roads = g }

// AliveCorps 回傳還在的軍團編號。
func (w *World) AliveCorps() []int {
	var out []int
	for i, c := range w.Corps {
		if c.Alive {
			out = append(out, i)
		}
	}
	return out
}
