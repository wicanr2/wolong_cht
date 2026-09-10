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
	// modelledCorpsBits 是 remake 真的有在維護的那幾個位元：
	// 7／6（存在）與 2（委任）。其餘位元寫回時原樣保留（docs/spec/166）。
	modelledCorpsBits = 0xC0 | 0x04 | 0x02 | 0x01 | standoffBit

	// standoffBit 是 `+0x00` 的位元 5 ＝ 對峙中（docs/spec/175）。
	standoffBit = 0x20

	// standoffTicks 是撞上之後 `+0x03` 的初值（`sub_12831`／`sub_12880`
	// 的 `mov byte ptr [si+3], 0Ch`）。每個軍團巡迴週期減 1，一圈 8 拍，
	// 所以對峙整整 12 × 8 ＝ 96 拍才開打（docs/spec/175）。
	standoffTicks = 0x0C

	// corpsSpriteStride 是一個勢力佔幾張軍團圖塊：四個方向 ＋ 停著。
	// 與 `internal/assets/world` 的 `CorpsHeadings` 是同一個數字
	// （`sub_12B2A`：圖塊 ＝ `[si+9]` ＋ `[si+8]`，docs/spec/74 §3）。
	corpsSpriteStride = 5
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

	// Direction 是 +0x0A：沿路徑表前進的**步進量**（`sub_127F6` 取負再相加）。
	// 存的是原始 byte：正向 `4`、反向 `0xFC`（＝ −4）。
	Direction int
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

	// PathPtr、LinkAddr 是行軍在**原版道路表**裡的位置
	// （記錄 `+0x0C`／`+0x0E`，docs/spec/172）。步進量在 `Direction`（`+0x0A`）。
	//
	// ⭐ `+0x0E` 有兩種語意：停在據點上是**據點編號 × 8**，行軍中是
	// **連結記錄的位址**（≥ 0x800）。`LinkAddr` 只放後者，前者由 `Node` 表示；
	// `LinkAddr == 0` 就代表「現在用據點那一種」。
	//
	// ⚠ 這三個是 remake 沒有拿來做決策的欄位——路由走的是 `routes`。
	// 它們存在的理由是**存檔要與原版互通**：不寫就等於每次存檔都把
	// 原版的行軍狀態抹掉（`CLAUDE.md` §9 的「改寫不是重建」）。
	PathPtr  int
	LinkAddr int

	// OnPath 是記錄 +0x00 的**位元 0 ＝「這一步是從路段中間走出來的」**
	// （docs/spec/173 §1.1）。
	//
	// 設定端是 `sub_126FF`（推進一個路徑點）與 `sub_147BB` 的三個
	// 「已經在路段上」分支；**清除端在 `sub_147BB` 的 `loc_1482F`**
	// （`and byte ptr [si], 0FEh`）——軍團站在據點上重新選路時清掉。
	// 所以它不是一次性的旗標，而是每次移動判定都會重寫。
	//
	// `sub_12662` 靠它分「重算完要不要判 leg 盡頭」，`sub_12708` 靠它
	// 決定要不要跑地形 0CEh–0DDh 那一段。
	OnPath bool

	// Replan 是記錄 +0x00 的**位元 1 ＝「下一步要重算」**（docs/re/34 §2.1）。
	//
	// 改行軍目標的常式只設這個旗標，**方向不當場換**——要等下一次
	// 「輪到移動」時 `sub_12662` 清掉它並呼叫 `sub_147BB`
	// （`0x126A5`），重算完那一拍照樣走一格（docs/spec/177 §1.6）。
	//
	// ⚠ 位元 1 沒設而且軍團走在邊上時，`sub_147BB` **根本不會被呼叫**：
	// 方向就照現有的 `+0x0A` 一路走到端點。
	Replan bool

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

	// Standoff 是記錄 +0x00 的**位元 5 ＝ 對峙中**：要踏進去的那一格被
	// 敵方軍團佔著（`sub_12831`），或那是別人的據點（`sub_12880`）。
	// 對峙中軍團**停在原地**，`Countdown` 倒數完才結算（docs/spec/175）。
	//
	// ⚠ 每次「輪到移動」時原版先清掉它（`sub_125A3` 的 `and [si],0DFh`），
	// 擋著的東西還在才會被再設一次——所以它是**每個週期重算的狀態**，
	// 不是黏著的旗標。
	Standoff bool

	// Countdown 是記錄 `+0x03`。⭐ **原版同一個 byte 兩種用途**，
	// 用 `+0x00` 分辨，而兩種狀態互斥（docs/spec/175 §1.5）：
	//
	//   - `Routing`（旗標 `08h`，已經不 Alive）→ **敗走倒數**，
	//     進狀態時 48，歸零時軍團記錄歸零、主將解職（docs/spec/43）。
	//   - `Standoff`（位元 5，還活著）→ **對峙倒數**，撞上時 12，
	//     每個軍團巡迴週期減 1，減到 1 的那一次才開打。
	//   - 兩者都不成立 → 恆為 **0**（`sub_1264A` 每個週期歸零）。
	//
	// ⚠ **不要拆成兩個 Go 欄位。** 一個 byte 拆成兩個欄位會多出
	// 「兩個欄位、一個 byte」的失步風險；原版本來就是靠 `+0x00` 分辨，
	// 照抄那個結構，載入／寫回／讀取三處都不可能對不上。
	// 分開的是**用途**（哪個狀態讀它），不是儲存。
	Countdown int
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
	w.occupancySeg = 0
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
			PathPtr:   u16(r, 0x0C),
			OnPath:    r[0x00]&0x01 != 0,
			Replan:    r[0x00]&0x02 != 0,
			Delegated: r[0x00]&0x04 != 0,
			Stage:     int(r[0x23]),
			// 旗標 8 而且不到 0x80 ＝ 敗走中（docs/spec/43）。
			Routing:   r[0x00] < aliveFlag && r[0x00]&0x08 != 0,
			// ⭐ 對峙只在**活著**的軍團身上成立；敗走的旗標是 `08h`，
			// 位元 5 本來就不會設（`sub_12977` 整個 byte 寫 8）。
			Standoff:  r[0x00] >= aliveFlag && r[0x00]&standoffBit != 0,
			Countdown: int(r[0x03]),
		}
		for k := range c.Units {
			s := r[unitSlotBase+k*unitSlotSize:]
			c.Units[k] = combat.Unit{Men: int(s[1]), Kind: kindFromByte(s[2])}
		}
		// ⭐ `+0x0E` ≥ 0x800 是**連結記錄的位址**（行軍中），不是據點編號。
		// 照 `/8` 讀會得到 470、488 這種不存在的據點（docs/spec/172 §1）。
		//
		// ⚠ 這時 `Node`（出發據點）**還算不出來**——要有道路圖才知道那條
		// leg 的兩端是誰。先擺目標當佔位，真正的值由 `restoreMarchRoutes`
		// 在 `SetRoads` 時填。在那之前這支軍團算「行軍狀態未還原」，
		// `tickOneCorps` 不會動它（否則第一次移動就會被判成抵達目標，
		// 連帶在錯的地方觸發攻城）。
		if raw := u16(r, 0x0E); raw >= 0x800 {
			c.LinkAddr = raw
			c.Node = c.TargetNode
		}

		// ⚠ **載入不可信的資料要驗範圍。** 未使用的軍團槽裡是垃圾，
		// 而 `Faction` 會被直接拿來索引 `w.Factions`（22 筆）——超範圍就是
		// 執行期 panic，而不是一個看得懂的錯誤。對拍時拿原版記憶體重建
		// 存檔踩過這個：欄位讀成 83，`tickOneCorps` 當場炸掉。
		// 這裡把它降級成「這個槽不存在」，讓壞資料不會變成 crash。
		if c.Faction < 0 || c.Faction >= numFactions {
			c.Alive, c.Faction = false, 0
		}
		// 佔用圖的段基底只要從任何一支活著的軍團反推一次就夠。
		if c.Alive && w.occupancySeg == 0 {
			if seg := u16(r, 0x1C) - c.Y*24; seg > 0 {
				w.occupancySeg = seg
			}
		}
		w.Corps[i] = c
	}
	// 存檔沒有佔用圖，載入端只能從位置重建（docs/spec/183）。
	w.rebuildOccupancy()
}

func (w *World) saveCorps(b []byte) {
	for i, c := range w.Corps {
		r := b[corpsBase+i*corpsSize:]
		if !c.Alive {
			// 敗走中的軍團要把旗標 8 與倒數寫回去，其餘欄位不動
			// （原版 `sub_12977` 也只改這兩個，docs/spec/43）。
			if c.Routing {
				r[0x00] = 0x08
				r[0x03] = byte(c.Countdown)
			}
			// 其餘不存在的軍團**一個 byte 都不動**——重建會抹掉痕跡。
			continue
		}
		// ⚠ **不要整個 byte 寫死。** remake 只建模位元 7／6（存在）與
		// 位元 2（委任）；位元 0／1／4／5 有設定端與清除端、語意還沒定案
		// （docs/re/34 §2），但它們是原版的狀態——覆寫等於每次存檔都抹掉。
		// 原版只在**建立軍團**時整個寫 0xC0（`sub_16F26`），既有軍團的
		// 其他位元由各自的維護端負責（docs/spec/166）。
		keep := byte(0)
		if r[0x00] >= aliveFlag {
			keep = r[0x00] &^ modelledCorpsBits
		}
		r[0x00] = newCorps | keep
		if c.Delegated {
			r[0x00] |= 0x04
		} else {
			r[0x00] &^= 0x04
		}
		if c.OnPath {
			r[0x00] |= 0x01
		} else {
			r[0x00] &^= 0x01
		}
		if c.Replan {
			r[0x00] |= 0x02
		} else {
			r[0x00] &^= 0x02
		}
		if c.Standoff {
			r[0x00] |= standoffBit
		} else {
			r[0x00] &^= standoffBit
		}
		// ⚠ **活著的軍團也要寫 `+0x03`。** 不在對峙時原版恆為 0
		// （`sub_1264A` 每個週期歸零）——不寫就會留著存檔裡的舊值。
		r[0x03] = byte(c.Countdown)
		r[0x01] = byte(c.Faction)
		r[0x02] = byte(i)
		putU16(r, 0x04, c.Men)
		r[0x06] = byte(c.Morale)
		r[0x08] = byte(c.Heading)
		// `+0x09` 是導出值：**勢力編號 × 5**（`sub_16FD2`，docs/spec/173 §1.3）。
		// 每個勢力五張圖塊（四方向 ＋ 停著），與 `world.CorpsHeadings` 同一個
		// 常數——狀態層不依賴資產層，所以這裡另寫一份並互指。
		r[0x09] = byte(c.Faction * corpsSpriteStride)
		r[0x0A] = byte(c.Direction)
		r[0x0B] = byte(c.Timer)
		putU16(r, 0x0C, c.PathPtr)
		// `+0x0E` 的兩種語意：行軍中是連結記錄位址、停著是據點 × 8。
		if c.LinkAddr != 0 {
			putU16(r, 0x0E, c.LinkAddr)
		} else {
			putU16(r, 0x0E, c.Node*8)
		}
		putU16(r, 0x10, c.X)
		putU16(r, 0x12, c.Y)
		// 佔用圖的兩欄是導出值：偏移 ＝ X、段 ＝ 基底 ＋ Y × 24。
		if w.occupancySeg > 0 {
			putU16(r, 0x1A, c.X)
			putU16(r, 0x1C, w.occupancySeg+c.Y*24)
		}
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

// recalcCorps 是 `sub_16FD2`：軍團內容變了就重算的統一入口
// （`docs/re/30` §5）。四個呼叫者，戰後的 `sub_1474A` 第一行就是它。
//
// 三件事：
//
//   - `+0x04` 總兵力 ＝ 六槽 `+1` 的和（16 位無號）
//   - `+0x1E` 移動間隔 ＝ 全騎馬 2、否則 3
//   - **`+0x0B` 移動計時寫 1**（`0001701D`）——不是寫成間隔，
//     所以重算過的軍團**下一次輪到就走**
//
// ⭐ 第三件事對同局面對拍是看得見的：一場戰鬥之後兩邊的計時器會同步
// 歸到 1，少了它，攻守雙方的移動節拍會各自漂掉（docs/spec/179）。
func (w *World) recalcCorps(i int) {
	c := &w.Corps[i]
	men, allCav := 0, true
	for _, u := range c.Units {
		if u.Kind != army.Cavalry {
			allCav = false
		}
		men += u.Men
	}
	c.Men = men & 0xFFFF
	c.Interval = IntervalMixed
	if allCav {
		c.Interval = IntervalCavalry
	}
	c.Timer = 1
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

// newCorpsRecord 是**兩條編成路徑共用的初值**（原版 `sub_16F26` 建記錄、
// `sub_16FD2` 收尾）。玩家編成走 `FormCorps`、AI 編成走 `autoFormCorps`，
// 兩邊各抄一份的結果是抄漏的欄位用 Go 的零值頂上——
// `+0x08` 變成「往 X 減」、`+0x0B` 變成間隔（docs/spec/173 §2）。
//
// ⭐ **靜止是 4 不是 0**（docs/spec/74 §3）。Go 的零值 0 是「朝 X 減的
// 方向走」，大地圖上會畫成側面行進的圖塊——剛編成的軍團站在城裡，
// 該畫靜止那一張（CLAUDE.md §7 第 11 條）。
//
// ⭐ **計時是 1 不是間隔**：`sub_16FD2` 收尾寫 `[si+0Bh] = 1`，
// 新編的軍團**下一拍就輪得到**。
func (w *World) newCorpsRecord(faction, capital, morale int) Corps {
	home := w.clampCity(capital)
	return Corps{
		Alive:   true,
		Faction: faction,
		Morale:  morale,
		Ordered: capital,
		Node:    home,
		Heading: HeadingStill,
		Timer:   1,
		X:       w.Cities[home].X,
		Y:       w.Cities[home].Y,
		// 目標先設成原地，行軍指令下達前不會動。
		TargetNode: home,
		TargetX:    w.Cities[home].X,
		TargetY:    w.Cities[home].Y,
	}
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

	c := w.newCorpsRecord(g.Faction, f.Capital, f.MoraleBase)
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
	// ⭐ 編成**會**進佔用圖：`sub_16F86` 設完 `+0x1A`／`+0x1C` 之後
	// `inc byte ptr [bx]`（docs/spec/183）。
	w.enterCell(c.X, c.Y)
	g.Duty = DutyCorpsLeader
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
	// FateSides 是那一刻的勝方與敗方勢力，key 同 Fate。
	//
	// ⭐ **不能事後從武將記錄反推**：被擒之後 `Generals[i].Faction`
	// 已經換成勝方，而訊息要比的是**舊主**（原版靠武將 `+0x1D` 保存，
	// `docs/spec/123`）。
	FateSides map[int]FateSide

	// Captured 不是 −1 表示這個 tick 佔下了某個據點。
	Captured int

	// TalkNotices 是這一支軍團在這個 tick 要跳的訊息框（目前只有進戰術畫面
	// 前那一則，`docs/spec/105`）。`World.tick` 會把它併進 `Event.TalkNotices`。
	TalkNotices []TalkNotice

	// Disbanded 表示這支軍團在這個 tick 解體了（`docs/spec/39`）——
	// 兵**回**預備兵池。
	Disbanded bool

	// Routed 表示這支軍團在這個 tick 敗走了（`docs/spec/43`）——
	// 兵**不回**池。兩者都會讓軍團從地圖上消失，但代價完全不同，
	// 所以計數要分開，不然量出來的「軍團損耗」會把回收算成損失。
	Routed bool

	// RoutEnded 表示這支軍團的敗走倒數在這個 tick 歸零（原版 `sub_12A7E`）：
	// 軍團記錄清掉、主將解職。**訊息掛在這一刻**，不是敗走的當下
	// （`docs/spec/77`）。
	RoutEnded bool

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
		// ⭐ **一圈是 128 格不是 127。** 原版 `sub_125A3` 在 16 次
		// `add si, 40h` 之後才 `cmp si, 1FC0h`，所以檢查點上的 si 只會是
		// 0x400 的倍數，一圈實際走過 si = 0…0x1FC0 ＝ 128 格、8 拍。
		// 軍團表本身也是 128 格（`0x22C0`–`0x42C0` ＝ 0x2000 ÷ 64）；
		// remake 只建模前 127 格（跟著武將數），第 128 格空轉一拍
		// ——用 127 取模會每 8 拍就比原版多轉一格（docs/spec/168 §2.1）。
		w.corpsCursor = (w.corpsCursor + 1) % corpsSlots
		if i >= numCorps {
			continue
		}
		if !w.Corps[i].Alive {
			// 敗走中的軍團不算活著，但還有一個倒數要跑
			// （原版 `sub_125A3` 的 `test byte [si+2240h], 8`）。
			if w.Corps[i].Routing && w.tickRout(i) {
				out = append(out, CorpsEvent{Corps: i, Enemy: -1, Captured: -1,
					Relocated: capital.None, GovernorReturned: noGovernor,
					RoutEnded: true})
			}
			continue
		}
		if ev := w.tickOneCorps(i, hour, rng); ev != nil {
			out = append(out, *ev)
		}
		// ⭐ `sub_1264A` 由 `sub_125A3` 在**每次巡到**時呼叫，
		// 不論這一拍有沒有輪到移動——倒數因此是「每 8 拍減 1」，
		// 不是「每 24 拍減 1」（docs/spec/175 §1.3）。
		w.tickStandoff(i)
	}
	return out
}

// tickStandoff 是 `sub_1264A`：對峙倒數。
//
//	test byte ptr [si], 20h
//	jnz  .hold
//	mov  byte ptr [si+3], 0      ; 沒卡著 → 歸零
//	mov  byte ptr [si+21h], 0
//	retn
//	.hold:
//	dec  byte ptr [si+3]
//	jnz  retn
//	mov  byte ptr [si+3], 1      ; 減到 0 就停在 1，等下一次撞上開打
//
// ⚠ `dec` 是 byte 運算：`+0x03` 是 0 時會變成 255，不是負的。
// 正常流程碰不到（設位元 5 的同一支常式當場寫 12），但**存檔可以**，
// 所以照抄比「看起來比較合理」的寫法安全。
//
// ⚠ `+0x21` 也一起歸零，remake 還沒建模那一欄（docs/spec/175 §5）。
func (w *World) tickStandoff(i int) {
	c := &w.Corps[i]
	if !c.Standoff {
		c.Countdown = 0
		return
	}
	c.Countdown = int(byte(c.Countdown - 1))
	if c.Countdown == 0 {
		c.Countdown = 1
	}
}

// atTargetNode 是原版 `sub_12662` 開頭的 `cmp bx, [si+14h]`
// （`+0x0E` 對 `+0x14`）：**軍團的現在節點就是行軍目標**。
//
// ⚠ 不能只比 `Corps.Node`。它是 remake 自己的欄位，行軍中留著出發那一站，
// 所以「從 X 出發、目標也是 X」（退卻回首都、掉頭）會**每走一格都判成
// 抵達**——連帶每一格都跑一次 `arriveCorps`，把目標欄位重寫一遍。
// 原版沒有這個問題：`+0x0E` 在路上是連結記錄位址（≥ `800h`），
// 一定不等於 `+0x14`（docs/spec/179 §2）。
func (w *World) atTargetNode(i int) bool {
	c := &w.Corps[i]
	if c.LinkAddr != 0 || c.Node != c.TargetNode {
		return false
	}
	// ⚠ 缺道路圖時 `LinkAddr` 恆為 0（`step` 退回直線逼近），上面那道閘
	// 就失效了，所以再問一次座標。比的是**目標據點的座標**，不是
	// `TargetX`／`TargetY`——後兩格在戰後退卻時留著舊目標的值
	// （`sub_1474A` 不寫它們，docs/spec/177 §1.4）。
	if !validCity(c.TargetNode) {
		return true
	}
	return c.X == w.Cities[c.TargetNode].X && c.Y == w.Cities[c.TargetNode].Y
}

func (w *World) tickOneCorps(i, hour int, rng combat.Rand) *CorpsEvent {
	c := &w.Corps[i]
	// 第二道保險：載入端已經擋過（loadCorps），這裡再擋一次，
	// 因為 Faction 在下面被當索引用，而快照/測試也可能塞進別的值。
	if c.Faction < 0 || c.Faction >= numFactions {
		return nil
	}
	// ⚠ **行軍狀態還沒還原就不要動它。** 存檔記著「走在哪一條路的第幾格」
	// （`+0x0E`／`+0x0C`／`+0x0A`），但那要有道路圖才解得開；在
	// `SetRoads` 之前 `Node` 只是暫時擺著目標，這時候跑一步會被判成抵達，
	// 連帶在錯的地方觸發攻城（docs/spec/172 §4.5）。
	// 不動比走錯好——不動看得出來，走錯看不出來。
	if c.LinkAddr != 0 && len(w.routes[i]) == 0 {
		return nil
	}
	ev := CorpsEvent{Corps: i, Enemy: -1, Captured: -1,
		Relocated: capital.None, GovernorReturned: noGovernor}

	// ① 移動的節拍。原版先減再判斷：間隔 N 表示每 N 個 tick 走一步。
	c.Timer--
	if c.Timer <= 0 {
		c.Timer = c.Interval
		// ⭐ **輪到移動就先清對峙旗標**（`sub_125A3` 的 `and [si],0DFh`）。
		// 擋著的東西還在的話，下面的移動會把它再設回來——所以它是
		// 每個週期重算的狀態，不是黏著的旗標（docs/spec/175）。
		c.Standoff = false
		// ⚠ **停在目標上也要跑抵達處理**：原版 `sub_12662` 一開頭就比
		// 「現在節點 ＝ 目標節點」，相同就直接呼叫 `sub_14325` 分派，
		// 不需要移動（`docs/re/64` §1）。解體下在「已經在首都」時就靠這條。
		// ⚠ 原版比的是 `+0x0E` 與 `+0x14`，而 `+0x0E` 在行軍中是**連結
		// 記錄位址**（≥ `800h`），所以「還在路上」永遠不會相等。
		// remake 的 `Node` 語意不同（行軍中留著出發那一站），等價寫法是
		// **先問還在不在路段上**——`LinkAddr` 只在走完整條路線時歸零
		// （docs/spec/179 §2）。
		if w.atTargetNode(i) {
			w.arriveCorps(i, rng)
			if !c.Alive {
				ev.Disbanded, ev.Routed = !c.Routing, c.Routing
				return &ev
			}
		} else {
			// ⭐ 位元 1 ＝「下一步要重算」：`sub_12662` 在 `0x126A5`
			// 清掉它並呼叫 `sub_147BB`，**然後照樣走一格**（沒有出口）。
			// 位元 1 沒設而且走在邊上時，方向根本不重算——就照現有的
			// `+0x0A` 一路走到端點（docs/spec/177 §1.6）。
			if c.Replan {
				c.Replan = false
				w.replanOnLeg(i)
			}
			// ⭐ **位元 0 在「問擋不擋」之前就寫好了。**
			// `sub_12662` 的 `loc_126F2` 走 `sub_126FF`（`or [si], 1`
			// ＋ 路徑指標前進），而 `sub_126FF` **尾端落空**到
			// `sub_12708`——也就是先設旗標、再問下一格擋不擋
			// （docs/re/34 §2.05）。剛從據點出發那一步走的是
			// `loc_126C0`，直接進 `sub_12708`，位元 0 維持 0。
			//
			// ⇒ 判準是**移動前**的 `+0x0E`：`≥ 800h`（在邊上）設、
			// `< 800h`（站在據點上）清。把它擺到 `standoffBlocks`
			// 之後就會漏掉「對峙那一拍」——原版那時已經設好了
			// （同局面拍 4,900 的軍團 73）。
			c.OnPath = c.LinkAddr != 0
			if w.standoffBlocks(i, &ev, rng) {
				// ⭐ **踏進去之前先問**（`sub_12708`）：下一格被敵方軍團
				// 佔著或是別人的據點，這一拍就**不動**——設位元 5、
				// `+0x03` 從 12 倒數，減到 1 的那一次才結算（docs/spec/175）。
				// `ev.Moved` 維持 false，對峙的 96 拍在事件層是靜的。
			} else if w.step(i) {
				ev.Moved = true
				// ⚠ **走完最後一格的那一拍不分派。** 原版 `sub_12708`
				// 寫完座標就落到 `loc_126F5`（佔用圖 +1）然後 `retn`；
				// `sub_14325` 只在**下一次輪到移動**時跑——那時
				// `cmp bx, [si+14h]` 才相等（docs/spec/179 §2.1）。
				// 少了這一拍的間隔，Stage 機會整條早一個移動週期。
				ev.Arrived = w.atTargetNode(i)
			}
		}
	}

	// ② 軍費與士氣。每天「一時」那個小時才收。
	// 這一步壞滅的軍團不收——它已經不在了。
	//
	// ⭐ **軍費是當場從資金扣的**（原版 `sub_12600` → `sub_1562B`），
	// 不進「本月支出」那個累加器——那一格只有預備兵維持費在寫
	// （`sub_13E65`，docs/spec/50）。總額一個月下來一樣，
	// 差在月中：出兵中的勢力資金會即時往下掉，而每小時的侵攻財政閘
	// 讀的正是資金。
	if hour == upkeepHour && c.Alive {
		// ⚠ 判準是 `cmp word ptr [si+0Eh], 800h` ＝ **有沒有走在路段上**，
		// 不是節點的種類（docs/spec/178）。`Corps.Node` 在行軍中留著
		// 出發那一站，拿它去問「在不在野外」永遠答「在城裡」——
		// 於是行軍中的軍團收便宜的軍費、而且每小時回 10 點士氣。
		onLeg := c.LinkAddr != 0
		f := &w.Factions[c.Faction]
		f.Funds = economy.ClampFunds(f.Funds - combat.Upkeep(c.Men, onLeg))
		cc := combat.Corps{Morale: c.Morale}
		combat.Recover(&cc, f.MoraleBase, onLeg)
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

	// ⓪ 邊界掉頭：前方那一端的據點屬於**和平**的別勢力就折返
	//    （`docs/spec/132`）。原版每一步都問一次，而且只在野外問。
	if w.turnBackAtBorder(i) {
		return true
	}

	// ① 有格子路徑就逐格走 —— **這條路徑每一格都踩在道路圖塊上**
	//    （`internal/assets/world` 有逐格檢查的測試）。
	if cells := w.routes[i]; len(cells) > 0 {
		next := cells[0]
		w.routes[i] = cells[1:]
		// ⚠ 位元 0 **不在這裡維護**——原版寫它的是 `sub_147BB`，
		// 而那一支不是每拍都跑（docs/re/34 §2.05）。維護點在呼叫端。
		// 同步吃掉一格標記：原版每走一步就 `bx += [si+0Ah]` 再寫回 `+0x0C`。
		var head = -1
		if mk := w.routeMarks[i]; len(mk) > 0 {
			c.PathPtr, c.LinkAddr, c.Direction = mk[0].PathPtr, mk[0].LinkAddr,
				byteStep(mk[0].Step)
			// ⭐ 朝向看的是**下一個路徑點**，不是剛走過的那一格
			// （`sub_12804`，docs/spec/173 §1.2）。
			if n := mk[0].Next; n != ([2]int{}) {
				head = headingTo(next[0], next[1], n[0], n[1])
			}
			w.routeMarks[i] = mk[1:]
		} else if c.LinkAddr != 0 {
			// 掉頭走的反向段沒有道路表的標記，但原版照樣每走一步就
			// `bx += [si+0Ah]` 寫回 `+0x0C`——步進是有號的，反向就是
			// 往回數（docs/spec/177 §1.1）。少了這一步，路徑點位址會
			// 凍在掉頭那一格，與原版逐格拉開。
			c.PathPtr += int8Step(c.Direction)
		}
		if head < 0 {
			// 沒有標記（掉頭走的反向段、或圖裡沒有格子序列）。
			// ⭐ 剩下的格子序列就是路徑點序列，所以**下一個路徑點**照樣
			// 問得到——原版 `sub_12804` 看的是它，不是剛走過的那一格
			// （docs/spec/173 §1.2）。
			if rest := w.routes[i]; len(rest) > 0 {
				head = headingTo(next[0], next[1], rest[0][0], rest[0][1])
			}
		}
		if head < 0 {
			// 真的沒有下一格（最後一步）→ 退回「這一步走的方向」。
			head = headingTo(c.X, c.Y, next[0], next[1])
		}
		c.Heading = head
		// 原版 `sub_12662` 的 `dec byte [di]`（`loc_12697`，走之前）與
		// `inc byte [di]`（`loc_126F5`，走完之後）——**只有真的走了一步
		// 才動圖**（docs/spec/183）。
		w.exitCell(c.X, c.Y)
		c.X, c.Y = next[0], next[1]
		w.enterCell(c.X, c.Y)
		// 踩到某個據點的座標就算抵達那個據點。中繼據點也要更新，
		// 不然攻城、遭遇這些判定會在錯的地方觸發。
		//
		// ⭐ **中繼據點也要把 `+0x0E` 換回據點編號 × 8。** 原版
		// `sub_127A2` 走完一條連結就 `mov [si+0Eh], bx`，而 `bx < 600h`
		// 時那就是據點節點——它不分「中途經過」與「終點」。下一拍
		// `sub_147BB` 的 `bx < 800h` 那一半再重新選路、寫新的連結位址
		// 並 `and [si], 0FEh` 清掉位元 0（`loc_14869`）。
		//
		// 少了這一步，remake 一路記著連結位址走完全程，於是
		// `+0x0E` 與位元 0 在每個中繼據點都與原版差一拍
		// （同局面拍 4,750：軍團 4 的 `+0x0E`、軍團 89 的位元 0）。
		// 連帶影響軍費與士氣——站在據點上那一拍不算在野外。
		if n := w.cityAt(next[0], next[1]); n >= 0 {
			c.Node = n
			c.LinkAddr = 0
		}
		if len(w.routes[i]) == 0 {
			c.Node = c.TargetNode
			c.Heading = HeadingStill
			// 到站：`+0x0E` 換回據點編號 × 8（`sub_127A2`），
			// 而 `+0x0C`／`+0x0A` **留著最後一筆**——原版沒有清它們。
			c.LinkAddr = 0
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
	w.exitCell(c.X, c.Y) // 同 §① 的佔用圖維護（docs/spec/183）
	c.X += sign(c.TargetX - c.X)
	c.Y += sign(c.TargetY - c.Y)
	w.enterCell(c.X, c.Y)
	if c.X == c.TargetX && c.Y == c.TargetY {
		c.Node = c.TargetNode
		c.Heading = HeadingStill
	}
	return true
}

// turnBackAtBorder 是原版的邊界掉頭（`sub_142AB`，`docs/spec/132`）。
//
// **要打誰，得先跟誰交戰。** 軍團走在野外路徑上時，每一步都問一次
// 「下一個要踏進的據點是誰的」：自己的、中立的、正在交戰的都放行，
// **和平的別勢力就把目標改成路的另一端**，也就是掉頭走回剛離開的據點。
//
// ⚠ **原版不在下令時檢查。** 指令下得成、訊息也跳「向{2}移動下」，
// 折返發生在路上，畫面上沒有任何錯誤提示——玩家只看到軍團走出去又走回來。
// 改成下令時拒絕會多出原版沒有的訊息，也會遮掉這條玩法規則。
//
// 回傳有沒有掉頭（掉頭那一 tick 不再往前走，與原版一致：
// `sub_142AB` 只改目標，移動由下一次重算負責）。
func (w *World) turnBackAtBorder(i int) bool {
	c := &w.Corps[i]
	if c.Node < 0 || c.Node >= len(w.Cities) {
		return false
	}
	// 站在據點上不問——原版 `sub_12662` 的 `cmp bx, 800h` 那一條。
	if c.X == w.Cities[c.Node].X && c.Y == w.Cities[c.Node].Y {
		return false
	}
	next := w.nextCityOnRoute(i)
	if next < 0 || next == c.Node {
		return false
	}
	if !w.borderIsClosed(c.Faction, next) {
		return false
	}
	// 掉頭：目標換成路的另一端 ＝ 剛離開的據點，路線用反向那一段。
	back := c.Node
	c.TargetNode, c.Ordered = back, back
	c.TargetX, c.TargetY = w.Cities[back].X, w.Cities[back].Y
	// ⚠ 掉頭走的是自己算的反向段，不是道路表的 leg，所以沒有標記可對。
	// 清掉，讓 `+0x0C`／`+0x0E` 留在掉頭前的值——原版 `sub_142AB` 也是
	// 直接改 `+0x14`／`+0x0E`，不重排路徑點。
	w.routes[i], w.routeMarks[i] = w.reverseLeg(next, back, c.X, c.Y), nil
	c.Heading = headingTo(c.X, c.Y, c.TargetX, c.TargetY)
	return true
}

// borderIsClosed 回「軍團屬於 faction 時，能不能踏進 node 這個據點」。
//
// 三個放行條件照抄 `sub_142AB`：同勢力、中立（`army.NeutralFaction`）、
// 或**交戰中**（交友度的最高位元 0）。
func (w *World) borderIsClosed(faction, node int) bool {
	owner := w.Cities[node].Owner
	if owner == faction || owner == combat.NeutralFaction {
		return false
	}
	if faction < 0 || faction >= len(w.Friendship) ||
		owner < 0 || owner >= len(w.Friendship[faction]) {
		// 既不是自己的、也不是中立、又不是合法勢力編號 ⇒ 資料壞了。
		// **當成擋住**：軍團停下來看得見，走進一個不存在的勢力看不見。
		return true
	}
	// ⚠ 判準是**含和平位元的原始值**，不是低 7 位的交友度。
	return !w.Friendship[faction][owner].AtWar()
}

// nextCityOnRoute 回「照現在的路線，下一個會踏進的據點」，沒有回 −1。
//
// 原版是一段一段走，所以「前方端點」就是這一段路的終點；remake 的路線
// 是多段串起來的，等價的問法是「序列裡第一個踩到據點座標的格子」。
func (w *World) nextCityOnRoute(i int) int {
	for _, cell := range w.routes[i] {
		if n := w.cityAt(cell[0], cell[1]); n >= 0 {
			return n
		}
	}
	return w.Corps[i].TargetNode
}

// replanOnLeg 是 `sub_147BB` 的「`+0x0E` ≥ `800h`」那一半：軍團走在
// 某條邊上時**只決定往這條邊的哪一端**，不從出發據點重算整條路
// （docs/spec/177 §1.1）。
//
// 原版三個分支：目標就是端點 B ⇒ 步進 `4`（正向）、目標就是端點 A ⇒
// `0FCh`（反向）、都不是 ⇒ `loc_1491B` 算成本挑一端。remake 用「離目標
// 較近的那一端」代替第三支——成本函數是自我修改碼，還沒讀出來
// （docs/spec/177 §5）。
//
// 回傳有沒有換方向。⚠ **重算完那一拍照樣走一格**——原版
// `sub_12662` 在 `0x126A8` 呼叫完 `sub_147BB` 之後就落到 `sub_12708`
// 寫座標，中間沒有出口。
func (w *World) replanOnLeg(i int) bool {
	c := &w.Corps[i]
	if c.LinkAddr == 0 || w.roads == nil {
		return false
	}
	a, b, ok := w.roads.EdgeByLink(c.LinkAddr)
	if !ok {
		return false
	}
	// 站在端點上就不算「在邊上」——那是 `< 800h` 那一半的事。
	for _, n := range []int{a, b} {
		if n >= 0 && n < len(w.Cities) &&
			c.X == w.Cities[n].X && c.Y == w.Cities[n].Y {
			return false
		}
	}
	want := b
	switch {
	case c.TargetNode == b:
	case c.TargetNode == a:
		want = a
	default:
		da, db := w.roads.Distance(a, c.TargetNode), w.roads.Distance(b, c.TargetNode)
		if da >= 0 && (db < 0 || da < db) {
			want = a
		}
	}
	ahead := w.nextCityOnRoute(i)
	if ahead == want {
		return false // 已經朝著那一端走，什麼都不必動
	}
	cells := w.reverseLeg(ahead, want, c.X, c.Y)
	if len(cells) == 0 {
		return false // 切不到就不要把路徑清空——那會讓軍團整支凍住
	}
	// ⚠ 反向段沒有道路表的標記可對，`+0x0C`／`+0x0E` 留在換向前的值：
	// 原版這一半也只寫 `+0x0A`，不碰那兩格。
	w.routes[i], w.routeMarks[i] = cells, nil
	c.Direction = byteStep(-int8Step(c.Direction))
	c.Heading = headingTo(c.X, c.Y, cells[0][0], cells[0][1])
	return true
}

// int8Step 把 `+0x0A` 的 byte 讀成有號步進（`4` 或 `0FCh` ＝ −4）。
func int8Step(v int) int { return int(int8(v)) }

// reverseLeg 回「從目前這一格走回 back」的格子序列。
//
// 反向那一段路的格子序列與正向是同一批，所以取 `CellRoute(from, back)`
// 再從目前這一格切開就好。切不到（缺道路圖、或這一段沒有格子序列）
// 就回 nil，讓 `step` 退回直線逼近。
func (w *World) reverseLeg(from, back, x, y int) [][2]int {
	if w.roads == nil {
		return nil
	}
	rev := w.roads.CellRoute(from, back)
	for k, cell := range rev {
		if cell[0] == x && cell[1] == y {
			return append([][2]int(nil), rev[k+1:]...)
		}
	}
	return nil
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

// onCity 回「這支軍團現在是不是**站在據點上**」。
//
// 原版的判準是節點欄 `+0x0E` 小於 `800h`（`sub_12662` 的
// `cmp bx, 800h`、`sub_1487B` 的「現在在野外」，docs/spec/46 §2）：
// 行軍中那一欄放的是**連結記錄的位址**，一定 ≥ `800h`。
//
// ⚠ **remake 的 `Node` 在行軍中留著出發／中繼據點的編號**，光看它會把
// 「走在路上」讀成「站在城裡」。座標一起看才等價——`turnBackAtBorder`
// 早就是這樣問的，這裡把同一個判準抽出來共用。
func (w *World) onCity(i int) bool {
	// ⚠ `army.KindOf` 對**負數**回 `CityNode`（它只比上界），
	// 所以範圍檢查要靠 `standsOn`，不能只看 KindOf。
	return army.KindOf(w.Corps[i].Node) == army.CityNode &&
		w.standsOn(i, w.Corps[i].Node)
}

// standsOn 回「這支軍團的**座標**是不是就在 node 那個據點上」。
//
// 原版 `sub_14C72` 收「同一格上有誰」的名單用的是座標比對
// （`cmp ax, [bx+12h]` ＋ `cmp dx, [bx+10h]`），**不是節點欄**——
// 走在路上的軍團節點欄放的是連結記錄位址，本來就對不上任何據點。
func (w *World) standsOn(i, node int) bool {
	if node < 0 || node >= len(w.Cities) || i < 0 || i >= len(w.Corps) {
		return false
	}
	c := &w.Corps[i]
	return c.X == w.Cities[node].X && c.Y == w.Cities[node].Y
}

// nextCell 回「這一拍要踏進去的那一格」，沒有下一步就回 false。
//
// 對應原版 `sub_12708` 進來時 `es:[bx]`／`es:[bx+2]` 指的那一筆路徑點：
// 指標已經在 `sub_126FF` 加過步進，所以是**下一筆**，不是腳下這一筆
// （docs/spec/173 §1.1）。
func (w *World) nextCell(i int) (int, int, bool) {
	if cells := w.routes[i]; len(cells) > 0 {
		// ⭐ **問的是原版的路徑點，不是走到的格子。** 兩者只差終點
		// 那一筆：格子序列的最後一格是據點中心，路徑點序列的最後一筆是
		// **那一端的城門格**（`march.Edge.BGate`）。守軍站在城門格上時
		// 原版撞到的是軍團（野戰），拿據點中心去問會撞到據點（攻城）——
		// 同局面拍 4,910 就差在這裡（docs/spec/186）。
		if mk := w.routeMarks[i]; len(mk) > 0 && mk[0].Point != ([2]int{}) {
			return mk[0].Point[0], mk[0].Point[1], true
		}
		return cells[0][0], cells[0][1], true
	}
	// 沒有道路圖時 `step` 走直線退路（缺素材要能降級跑），
	// 下一格就是逼近目標的那一步。
	c := &w.Corps[i]
	if c.X == c.TargetX && c.Y == c.TargetY {
		return 0, 0, false
	}
	return c.X + sign(c.TargetX-c.X), c.Y + sign(c.TargetY-c.Y), true
}

// blockerAt 回「(x,y) 這一格擋不擋得住 i 這支軍團」：
// 回傳擋路的據點編號（攻城那條）或軍團編號（野戰那條），都沒有就回 −1／−1。
//
// 兩條路對應原版 `sub_12708` 裡**串聯**的兩個閘：
//
//	cmp byte ptr [di], 0   ; 佔用圖有人 → sub_12831（是敵人就擋）
//	test byte ptr [si], 1  ; 已經走上路徑，而且
//	cmp al, 0CEh / 0DDh    ; 下一格的地形是據點圖塊 → sub_12880（不是自己的就擋）
//
// ⭐ **據點要先問，佔用圖後問**——順序與原版相反，理由是地圖模型不同：
// 原版一個據點佔 `0CEh`–`0DDh` 一整段圖塊，攻方踏進的通常是**別的那幾格**，
// 那幾格佔用圖是 0 所以走攻城，再由 `sub_14C72` 用據點座標把守軍找出來。
// 本專案的據點是**一個點**，攻方必然踏在守軍那一格上，照抄順序會永遠
// 打成野戰、攻城那條路永遠走不到（docs/re/09 §2）。
func (w *World) blockerAt(i, x, y int) (int, int) {
	c := &w.Corps[i]
	// ⚠ **缺道路圖時沒有「城門格」這回事**（`routes` 是空的，`nextCell`
	// 回的是逼近目標的那一步 ＝ 據點中心）。那時守軍必然站在同一格上，
	// 照原版順序問會把每一場攻城都變成野戰——所以降級路徑**據點先問**。
	// 這是明示的 remake 差異，只在缺原版素材時走到（docs/spec/186 §3）。
	if len(w.routes[i]) == 0 {
		if n := w.cityAt(x, y); n >= 0 && w.Cities[n].Owner != c.Faction {
			return n, -1
		}
	}
	// ⭐ **① 先問軍團**（`sub_12831`），沒被擋才問據點（`sub_12880`）。
	// 掃軍團表找**第一支**站在那一格上而且活著的；找到就停，
	// 是自己人就放行（軍團可以疊同格）而且**繼續往下問據點**
	// （原版 `sub_12831` 回 STC，`sub_12708` 的 `jb loc_1273C`）。
	//
	// ⚠ 這裡的 `(x, y)` 是**路徑點**，不是走到的格子——最後一步的
	// 路徑點是城門格，據點中心是格子（docs/spec/186）。守軍站在
	// 城門格上時撞到的是軍團，站在據點中心時撞不到、往下走攻城。
	for j := range w.Corps {
		d := &w.Corps[j]
		if j == i || !d.Alive {
			continue
		}
		if d.X == x && d.Y == y {
			if d.Faction != c.Faction {
				return -1, j
			}
			break
		}
	}
	// ⭐ **② 再問據點**（`sub_12880`）：這條連結通往的據點是別人的就擋。
	//
	// 原版看的是連結的端點據點，而 `sub_12708` 用「下一格的圖塊落在
	// `0CEh`–`0DDh`」＋ `test byte ptr [si], 1` 把它限制在**快踏進據點**
	// 的那幾步。remake 的等價條件是「這一步走到的格子是段的終點
	// （據點中心），而路徑點已經是城門格」——兩者只在那一步不相等。
	cell := [2]int{x, y}
	if cells, mk := w.routes[i], w.routeMarks[i]; len(cells) > 0 {
		cell = cells[0]
		// 路徑點就是這一格 ⇒ 還沒走到城門格，據點那一條不成立。
		// ⚠ `Point` 是零值表示這一段沒有道路表的標記（`cellRoute` 補的
		// 零值 CellMark），那時沒有城門格可比，退回舊行為逐格問。
		if len(mk) > 0 && mk[0].Point != ([2]int{}) && mk[0].Point == cells[0] {
			return -1, -1
		}
	}
	if n := w.cityAt(cell[0], cell[1]); n >= 0 && w.Cities[n].Owner != c.Faction {
		return n, -1
	}
	return -1, -1
}

// standoffBlocks 是 `sub_12708` 的「踏進去之前先問」：擋住就**這一拍不動**，
// 設 `+0x00` 位元 5、`+0x03` ← 12，減到 1 的那一次才結算（docs/spec/175）。
//
// 回傳 true 表示這一拍不移動。倒數由 `tickStandoff`（＝ `sub_1264A`）
// 在**每次巡到**時減 1，一圈 8 拍 ⇒ 對峙整整 96 拍。
//
// ⚠ 擋路的對象**每個移動拍重新找一次**，不記在軍團記錄裡——原版也是
// 這樣（`sub_12831` 每次都重掃 127 支）。敵人先走掉，對峙就自己散了。
func (w *World) standoffBlocks(i int, ev *CorpsEvent, rng combat.Rand) bool {
	c := &w.Corps[i]
	x, y, ok := w.nextCell(i)
	if !ok {
		return false
	}
	node, enemy := w.blockerAt(i, x, y)
	if node < 0 && enemy < 0 {
		return false
	}
	c.Standoff = true
	switch {
	case c.Countdown > 1:
		// 還在倒數。原版每個週期播一次音效（`sub_102F5(al=3)`），
		// 而且 `+0x03 & 3` 是對峙動畫的相位——兩者 remake 都還沒接
		// （docs/spec/175 §5）。
		return true
	case c.Countdown == 1:
		// ⭐ 減到 1 的那一次才結算。`sub_14A7B`／`sub_14ADE` 一進去
		// 就把雙方的位元 5 與 `+0x03` 清掉。
		c.Standoff, c.Countdown = false, 0
		if enemy >= 0 {
			w.Corps[enemy].Standoff, w.Corps[enemy].Countdown = false, 0
		}
		if node >= 0 {
			w.siegeAt(i, node, ev, rng)
		} else {
			w.fieldAt(i, enemy, ev, rng)
		}
		return true
	default:
		// 第一次撞上：`mov byte ptr [si+3], 0Ch`。
		c.Countdown = standoffTicks
		return true
	}
}

// siegeAt 是對峙倒數走完之後**打據點**那一條（`sub_12880` 的
// `call sub_14ADE`）。
//
// ⚠ 目標據點由呼叫端指定，不能用 `Node` 去找：軍團**還沒踏進去**，
// `Node` 留在前一站（docs/spec/175 §2）。
func (w *World) siegeAt(i, node int, ev *CorpsEvent, rng combat.Rand) {
	c := &w.Corps[i]
	if node < 0 || node >= len(w.Cities) {
		return
	}
	city := &w.Cities[node]
	if city.Owner == combat.NeutralFaction {
		// 中立據點沒有主人，所以沒有「首都失守」這回事。
		// ⚠ 它一樣要對峙滿 12 個週期——`sub_12880` 的歸屬檢查是
		// `cmp [di+841h], al`，無主（0x18）與別人的據點走同一條路。
		city.Owner = c.Faction
		w.Factions[c.Faction].Cities++
		ev.Captured = node
		return
	}
	if city.Owner == c.Faction {
		return
	}
	// 城裡有守軍就打守軍（多支疊同格時照 `sub_14C72` 計分挑應戰者，
	// docs/spec/82），沒有就打城兵。
	//
	// ⚠ **判準是「腳踩在據點的座標上」**（`sub_14C72` 收的是同一格的名單，
	// 攻城時那一格就是據點座標）。拿 `Node` 比會把「從這座城出發、
	// 正走在路上」的軍團算成守軍——`Node` 在行軍中留著出發那一站。
	cx, cy := city.X, city.Y
	if j := w.pickDefender(i, city.Owner, func(d *Corps) bool {
		return d.X == cx && d.Y == cy
	}); j >= 0 {
		w.fight(i, j, node, ev, combat.Siege, city.Garrison, rng)
		return
	}
	w.fightGarrison(i, node, ev, rng)
}

// fieldAt 是對峙倒數走完之後**打軍團**那一條（`sub_12831` 的
// `call sub_14A7B`）。撞到的那一支由 `blockerAt` 指定；同一格上疊了
// 好幾支時照 `sub_14C72` 計分挑應戰者（docs/spec/82）。
//
// 戰場沿用攻方的 `Node`——原版的野戰戰場是 `sub_14B63` 從**守方那一格
// 周圍的五格地形**算出來的，remake 這一層還是近似（docs/re/78）。
func (w *World) fieldAt(i, enemy int, ev *CorpsEvent, rng combat.Rand) {
	if enemy < 0 || enemy >= len(w.Corps) {
		return
	}
	x, y, f := w.Corps[enemy].X, w.Corps[enemy].Y, w.Corps[enemy].Faction
	if k := w.pickDefender(i, f, func(d *Corps) bool {
		return d.X == x && d.Y == y
	}); k >= 0 {
		enemy = k
	}
	w.fight(i, enemy, w.Corps[i].Node, ev, combat.Field, 0, rng)
}

// pickDefender 是 `sub_14C72` 的挑選：faction 勢力裡通過 at 條件
// （同格／同據點）的軍團逐支計分，**嚴格最大**者應戰（同分取先者）：
//
//	分數 ＝ ((兵數 >> 4) 的低 byte) × (士氣 >> 4) × ((大將評價 >> 4) ＋ 1)
//
// 評價是 `sub_155A6` 的衍生值（`general.Rating()`）。沒有候選回 −1
// （原版 CF=1，攻城走打城兵那條）。attacker 只用來排除自己。
func (w *World) pickDefender(attacker, faction int, at func(*Corps) bool) int {
	best, bestScore := -1, -1
	for j := range w.Corps {
		d := &w.Corps[j]
		if j == attacker || !d.Alive || d.Faction != faction || !at(d) {
			continue
		}
		rating := 0
		if lead := w.Leader(j); lead >= 0 && lead < len(w.Generals) {
			rating = w.Generals[lead].Rules().Rating()
		}
		score := ((d.Men >> 4) & 0xFF) * (d.Morale >> 4) * (rating>>4 + 1)
		if score > bestScore {
			best, bestScore = j, score
		}
	}
	return best
}

func (w *World) fight(att, def, node int, ev *CorpsEvent, m combat.Mode, garrison int, rng combat.Rand) {
	// ⭐ 玩家的勢力捲進去而且那一方**沒有委任**，就直接開戰術畫面
	// （原版 `sub_14E5C`／`sub_14ED7`：`sub_14EB9` → `sub_11B5A`，**中間沒有選單**，
	// 實機 docs/playtest/55）。「戰鬥指揮／委任」是行軍指示時就決定的
	// （docs/spec/39），遭遇當下只看委任位元。其餘自動判定。
	ev.Mode = m
	if w.wantsTactical(att, def) && w.beginTactical(att, def, node, m, garrison) {
		// 原版進戰術畫面前先跳一則訊息（`sub_14EB9`／`sub_14F58`，docs/spec/105）。
		ev.TalkNotices = append(ev.TalkNotices, w.encounterNotice(att, def, m))
		return
	}
	w.resolveCorpsBattle(ev, att, def, node, m, garrison, rng)
}

// 進戰術畫面前那一則訊息的 TALK 索引（原版 `sub_14E5C`／`sub_14ED7` 的 `cx`）。
const (
	// talkSiegeCityFallen ＝ #26「{2}受到{1}兵馬的攻擊，被攻陷了！！」
	// （空城自動判定後，`sub_14ED7` 的 `cx = 1Ah`）。
	talkSiegeCityFallen = 0x1A
	// talkSiegeIncoming ＝ #27「{1}的兵馬，向{2}進攻過來了！！」（玩家守城）。
	talkSiegeIncoming = 0x1B
	// talkSiegeOutgoing ＝ #28「{1}大人的兵馬，向{2}進攻了！！」（玩家攻城）。
	talkSiegeOutgoing = 0x1C
	// talkFieldEncounter ＝ #29「{1}大人的兵馬，遇上{1}的兵馬了！！」（野戰）。
	// ⭐ **兩個 `{1}`**：前者攻方主將、後者守方主將（`sub_14EB9` 依序推兩個參數）。
	talkFieldEncounter = 0x1D
)

// encounterNotice 是進戰術畫面前的那一則訊息（docs/spec/105 §1）。
//
// 攻城兩則的 `{1}` 都是**攻方**主將、`{2}` 都是據點，差別只在玩家站哪一邊；
// 野戰那一則的兩個 `{1}` 是攻守兩個主將。
func (w *World) encounterNotice(att, def int, m combat.Mode) TalkNotice {
	n := TalkNotice{City: -1, Faction: -1, General: w.Leader(att), Amount: -1}
	if m == combat.Siege {
		n.Index = talkSiegeIncoming
		if w.Corps[att].Faction == w.Player {
			n.Index = talkSiegeOutgoing
		}
		n.City = w.Corps[att].Node
		return n
	}
	// ⚠ **第一個 `{1}` 是玩家那一方的主將**，不是攻方：原版玩家守方那條路
	// 在呼叫 `sub_14EB9` 之前先 `xchg si, di`（`0x14E9F`），所以兩個參數
	// 的順序跟著玩家走。攻城那兩則沒有這一步，`{1}` 一律是攻方。
	first, second := att, def
	if def >= 0 && def < len(w.Corps) && w.Corps[def].Faction == w.Player {
		first, second = def, att
	}
	n.Index = talkFieldEncounter
	n.General = w.Leader(first)
	n.SeqGenerals = []int{w.Leader(first), w.Leader(second)}
	return n
}

// resolveCorpsBattle 執行一場已決定委任的軍團對軍團戰鬥。
// 戰鬥指揮的出口走 ResolvePending；兩者最後共用同一組戰後處理。
func (w *World) resolveCorpsBattle(ev *CorpsEvent, att, def, node int, m combat.Mode, garrison int, rng combat.Rand) {
	a, d := w.battle(att), w.battle(def)
	ev.BattleBefore = [2]int{a.Men, d.Men}
	r := combat.Resolve(&a, &d, m, garrison, rng)
	w.applyBattle(att, a)
	w.applyBattle(def, d)
	ev.Battle, ev.Enemy, ev.Mode = &r, def, m
	ev.BattleAfter = [2]int{a.Men, d.Men}
	ev.BattleCityDamage = r.CityDamage
	w.damageCity(node, m, r)

	// 原版兩邊各跑一次 `sub_1474A`：士氣判定之外，**敗方退不了也算壞滅**
	// （docs/spec/46 §1）。守方站在自家城裡走「不退」那一支，
	// 所以攻城的易主判定不受影響。
	attDead := r.AttackerDestroyed || w.retreatOrPerish(att, !r.DefenderWins)
	defDead := r.DefenderDestroyed || w.retreatOrPerish(def, r.DefenderWins)
	w.afterBattle(ev, att, node, attDead, def, rng)
	w.afterBattle(ev, def, node, defDead, att, rng)

	if defDead && !attDead && m == combat.Siege {
		w.capture(att, node, ev, rng)
	}
}

// fightGarrison 打的是據點的城兵——原版在 `ds:4200h` 現搭一支臨時軍團
// （`sub_14F8A`，docs/re/09 §7）。守方不是軍團，所以不會有壞滅或被擒。
func (w *World) fightGarrison(att, node int, ev *CorpsEvent, rng combat.Rand) {
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
	w.afterBattle(ev, att, node, attDead, -1, rng)
	if !r.DefenderWins && !attDead {
		// ⭐ **敵軍攻下玩家的空城** → 原版跳 #26（`sub_14ED7` 的 `loc_14EF1`：
		// 玩家是守方而 `bx == 4200h`（城裡沒有駐守軍團）→ `sub_15130`
		// 自動判定，`al == 0`（攻方贏）才 `sub_14F71`）。
		// **玩家自己攻空城那條路是靜的**——`loc_14F2B` 的
		// `cmp bx, 4200h` 直接 `jz loc_14EE7`，判定完就回，沒有訊息。
		// 變數順序照 `sub_14F71` 推堆疊的次序：先據點（`di`）後主將
		// （`ax = [si+2]`），對應 #26「{2}受到{1}兵馬的攻擊」。
		fellForPlayer := city.Owner == w.Player
		w.capture(att, node, ev, rng)
		if fellForPlayer && ev.Captured == node {
			ev.TalkNotices = append(ev.TalkNotices, TalkNotice{
				Index: talkSiegeCityFallen, City: node, Faction: -1,
				General: w.Leader(att), Amount: -1,
			})
		}
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
func (w *World) afterBattle(ev *CorpsEvent, i, node int, destroyed bool, victor int, rng combat.Rand) {
	if !destroyed {
		return
	}
	winner := combat.NeutralFaction
	if victor >= 0 {
		winner = w.Corps[victor].Faction
	} else if army.KindOf(node) == army.CityNode && node >= 0 && node < len(w.Cities) {
		// 勝方是城兵：那一場的據點由呼叫端指定——軍團對峙時還沒踏進去，
		// `Node` 找不到它（docs/spec/175 §2）。
		winner = w.Cities[node].Owner
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
	w.exitCell(c.X, c.Y) // 佔用圖 −1（原版 `sub_12977`，docs/spec/183）
	g.Duty = DutyNone
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
	if ev.FateSides == nil {
		ev.FateSides = map[int]FateSide{}
	}
	ev.FateSides[i] = FateSide{Winner: winner, Loser: loser}
}

// FateSide 是一次下場判定的兩個勢力（`docs/spec/123`）。
type FateSide struct{ Winner, Loser int }

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
	// `sub_14D63` 只清職務，**不清 +0x1A**（那是手動解任才做的，docs/spec/143 §2）。
	w.Generals[id].Duty = DutyNone
	return id
}

// noGovernorSlot 是據點記錄 `+0x19` 的哨兵（原版寫 `0xFF`）。
// 存檔要 byte-for-byte 寫得回去，所以這裡存的是原始值不是 −1。
const noGovernorSlot = 0xFF

// capture 把據點換手（`sub_14CF3`）。
func (w *World) capture(att, node int, ev *CorpsEvent, rng combat.Rand) {
	if node < 0 || node >= len(w.Cities) || army.KindOf(node) != army.CityNode {
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
		// ⚠ **名單比的是座標，不是 `Node`**（`sub_14C72` 的
		// `cmp ax, [bx+12h]` ＋ `cmp dx, [bx+10h]`）。remake 的 `Node`
		// 在行軍中留著出發／中繼那一站，拿它當條件會把**還走在半路上**
		// 的軍團也調頭——原版只動站在那一格上的守軍（docs/spec/47 §4.1）。
		if !c.Alive || c.Faction != old || !w.standsOn(i, node) {
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
	w.routes[corps], w.routeMarks[corps] = nil, nil
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
	// `marks` 與格子逐格對應，是軍團 `+0x0A`／`+0x0C`／`+0x0E` 的來源
	// （docs/spec/172）。
	w.routes[corps], w.routeMarks[corps] = w.roads.CellRouteMarked(c.Node, node)
	return nil
}

// SetRoads 掛上道路圖。**規則層不讀檔案**，所以圖由呼叫端
// （`cmd/wlgame` 等）從 `MMAP` 推導後注入。
//
// 沒有掛的話行軍退回直線移動——缺原版素材時要能降級跑，
// 不是整個動不了。
func (w *World) SetRoads(g *march.Graph) {
	w.roads = g
	w.restoredMarches, w.unresolvedMarches = w.restoreMarchRoutes()
}

// ClearMarchRoute 丟掉一支軍團的格子路徑與存檔裡的路徑指標。
//
// ⭐ 給「把軍團直接擺到某一格」的驗收捷徑用（`internal/battlesetup`）：
// 座標被外力改掉之後，`+0x0C`／`+0x0E` 描述的那條 leg 就不再是它走的路。
// 留著有兩個後果——沿舊路徑走回去，或者被 `tickOneCorps` 的還原守衛
// 當成「有 `LinkAddr` 卻沒有路徑」而**整支凍住**（docs/spec/172 §4.6）。
func (w *World) ClearMarchRoute(i int) {
	if i < 0 || i >= len(w.Corps) {
		return
	}
	w.routes[i], w.routeMarks[i] = nil, nil
	w.Corps[i].LinkAddr, w.Corps[i].PathPtr = 0, 0
}

// UnresolvedMarches 回傳「載入之後行軍狀態還原不了」的軍團數。
//
// ⭐ 給對拍夾具與長跑用：那些軍團**完全不動**（`tickOneCorps` 早退），
// 所以少了這個計數，症狀會長得像「AI 什麼都沒做」而不是「載入沒還原」。
func (w *World) UnresolvedMarches() int { return w.unresolvedMarches }

// RestoredMarches 回傳「載入之後**接回**行軍路徑」的軍團數。
//
// ⭐ 它是 `UnresolvedMarches() == 0` 的**正對照**。沒有它，
// 「一支都沒還原不了」與「根本沒有軍團需要還原」在斷言上長得一樣——
// 那是假零（`~/diagnosis-notes` 02：查詢回空的四種形狀）。
func (w *World) RestoredMarches() int { return w.restoredMarches }

// restoreMarchRoutes 把「載入存檔時正在行軍」的軍團接回格子路徑。
//
// ⭐ 存檔裡的 `+0x0E`（連結記錄位址）＋ `+0x0C`（路徑點位址）＋ `+0x0A`
// （步進方向）**完整描述了軍團走在哪一條路的第幾格**，所以路徑重建得回來。
// 少了這一步，載入之後那些軍團走的是直線退路——而畫面與存檔欄位看起來
// 都正常，只有逐格對拍才看得見（docs/spec/172）。
//
// ⚠ 道路圖是呼叫端注入的，所以這件事只能掛在 `SetRoads` 上，不能放在
// `loadCorps` 裡——那時還沒有圖。
func (w *World) restoreMarchRoutes() (restored, unresolved int) {
	if w.roads == nil {
		return 0, 0
	}
	for i := range w.Corps {
		c := &w.Corps[i]
		if !c.Alive || len(w.routes[i]) > 0 {
			continue
		}
		switch {
		case c.LinkAddr != 0:
			// ① 走在某條 leg 上：`+0x0E` 說是哪一條、`+0x0A` 說往哪一邊、
			//    `+0x0C` 說走到第幾格。
			a, b, ok := w.roads.EdgeByLink(c.LinkAddr)
			if !ok {
				unresolved++
				continue
			}
			// ⚠ **`+0x0A` 可能是 0**：軍團已經在野外、但路徑指標還沒設
			// （`+0x00` 位元 1 ＝「下一步要重算」，原版下一拍才跑
			// `sub_147BB` 選路）。這時方向要從別處推——先看目標是不是
			// 這條邊的一端，再兩端都試一次，能對齊的才算數。
			cands := []int{a, b}
			switch {
			case int8(c.Direction) < 0:
				cands = []int{b, a} // 0xFC ＝ −4，沿路徑表倒著走
			case a == c.TargetNode:
				cands = []int{b, a}
			}
			done := false
			for _, from := range cands {
				cells, marks := w.roads.CellRouteMarked(from, c.TargetNode)
				k := alignRoute(cells, marks, c.PathPtr, c.X, c.Y)
				if k < 0 {
					continue
				}
				c.Node = from
				w.routes[i], w.routeMarks[i] = cells[k+1:], marks[k+1:]
				done = true
				break
			}
			if done {
				restored++
			} else {
				unresolved++
			}
		case c.Node != c.TargetNode:
			// ② 目標已經定了但還沒踏出去（原版是下一拍 `sub_147BB` 才選路）。
			//    整條重算就好，軍團還站在出發據點上。
			w.routes[i], w.routeMarks[i] = w.roads.CellRouteMarked(c.Node, c.TargetNode)
			restored++
		}
	}
	return restored, unresolved
}

// alignRoute 找出軍團現在站在整條路線的第幾格。
//
// ⭐ 兩層：先比**路徑點位址**（`+0x0C`，精確），對不上再比**座標**——
// 目標被改過時位址會落在別條 leg 上，但軍團的座標一定還在路線上。
// 兩層都對不上就回 −1，呼叫端把它算成「還原不了」而不是默默從頭走。
func alignRoute(cells [][2]int, marks []march.CellMark, pathPtr, x, y int) int {
	for k := range marks {
		if marks[k].PathPtr == pathPtr {
			return k
		}
	}
	for k := range cells {
		if cells[k][0] == x && cells[k][1] == y {
			return k
		}
	}
	return -1
}

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

// byteStep 把 ±4 的步進量換成記錄裡的 byte（`4` 與 `0xFC`）。
func byteStep(step int) int { return int(byte(int8(step))) }
