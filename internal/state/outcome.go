package state

// OutcomeKind 是目前已證實、會離開 DOS/V 主遊戲循環的玩家結果。
//
// 原版證據：sub_13DC9 在信賴度歸零後以 selector 0x019E 顯示訊息並呼叫
// sub_11CB1；sub_14DF0 找不到替代據點時寫 capital=0xFF／清 alive bit，
// sub_14FCE 對玩家勢力同樣呼叫 sub_11CB1。sub_11CB1 的後續去向仍未知，
// 因此這裡只表達敗北原因，不把它命名成返回標題或結束程式。
type OutcomeKind uint8

const (
	InProgress OutcomeKind = iota
	DefeatTrustZero
	DefeatFactionEliminated
	// Victory 是原版離開碼 2（`D7END.EXE`）。閘門只有一個：
	// 存活勢力數減到 1（`sub_11CD0` 的 `cmp cs:byte_10D2A, 1`）。
	// ⚠ 「佔領所有城池」在機器碼裡**不是**另一條規則（docs/re/59 §5）。
	Victory
)

// Outcome 回傳目前結果。它是唯讀 accessor；結果一旦離開 InProgress 就不會
// 被另一個原因覆蓋。
func (w *World) Outcome() OutcomeKind {
	if w == nil {
		return InProgress
	}
	return w.outcome
}

// OutcomeMessageSelector 回傳已證實的原版訊息 selector。勢力滅亡的原版
// selector 尚未定位，所以刻意回傳 ok=false。
func (w *World) OutcomeMessageSelector() (uint16, bool) {
	if w == nil {
		return 0, false
	}
	switch w.outcome {
	case DefeatTrustZero:
		return 0x019E, true
	case Victory:
		// 原版結局送 TALK #75 與 #407（docs/re/59 §2）；
		// 內容還沒對出來，所以這裡只回前者。
		return 0x004B, true
	default:
		return 0, false
	}
}

// latchOutcome 只在真實 mutation 邊界呼叫；不可用每幀掃描補設結果。
func (w *World) latchOutcome(kind OutcomeKind) bool {
	if w == nil || kind == InProgress || w.outcome != InProgress {
		return false
	}
	w.outcome = kind
	return true
}

// AdjustTrust 是所有執行期信賴度扣減／增加的共同 mutation API。
// 由 1 降至 0 的同一個寫入邊界會 latch DefeatTrustZero；已經 latch 後不再
// 改寫結果。這裡保留 u8 飽和，對應 sub_13D91／sub_13DC9。
func (w *World) AdjustTrust(delta int) int {
	if w == nil || w.outcome != InProgress {
		if w == nil {
			return 0
		}
		return w.Trust
	}
	before := clampU8(w.Trust)
	w.Trust = clampU8(before + delta)
	if before > 0 && w.Trust == 0 {
		w.latchOutcome(DefeatTrustZero)
	}
	return w.Trust
}

// eliminateFaction 是**唯一**讓一個勢力滅亡的入口。
//
// 原版的順序有意義（`sub_14FCE`，docs/re/59 §4）：先判斷滅亡的是不是
// 玩家所仕的勢力（→ 敗北），**才**減存活勢力數。反過來寫的話，
// 玩家自己滅亡時會因為「剩一個」被誤判成結局。
//
// winner 是打下最後一個據點的勢力；沒有對象時傳 noFaction。
func (w *World) eliminateFaction(i, winner int) {
	if w == nil || i < 0 || i >= numFactions || !w.Factions[i].Alive {
		return
	}
	w.Factions[i].Alive = false
	if i == w.Player {
		w.latchOutcome(DefeatFactionEliminated)
	}
	if w.LivingFactions > 0 {
		w.LivingFactions--
	}
	w.disperseFaction(i, winner)
	if w.LivingFactions == 1 {
		w.latchOutcome(Victory)
	}
}

// disperseFaction 是 `sub_14FCE` 的武將處置迴圈：滅亡的勢力**不會留下
// 在外面的軍團**，127 名武將逐一按三條路處置。
//
//	俘虜（+0x1D ≠ 0xFF）      → sub_150D7 釋放，回原主或在野
//	非君主且有職務            → sub_150B4 解散軍團、變成在野武將
//	君主本人，或無職          → sub_129C3 軍團歸零，成為勝方的俘虜
//
// **這條規則先前是「還沒讀出來」的那一條**：舊寫法只清 Alive，
// 於是首都被打下來的勢力會留下沒有主人的軍團。
func (w *World) disperseFaction(i, winner int) {
	lord := w.Factions[i].Lord
	for gi := range w.Generals {
		g := &w.Generals[gi]
		if !g.Alive || g.Faction != i {
			continue
		}
		if g.Captor != noFaction {
			// 這一位是被 i 關著的俘虜，主人沒了就放出來。
			w.releaseGeneral(gi)
			continue
		}
		hadPost := g.Posted
		g.Posted = false
		if gi != lord && hadPost {
			g.Faction = noFaction // 在野
			continue
		}
		if winner >= 0 && winner < numFactions {
			g.Captor, g.Faction = i, winner
		} else {
			g.Faction = noFaction
		}
	}
	for ci := range w.Corps {
		c := &w.Corps[ci]
		if !c.Alive || c.Faction != i {
			continue
		}
		c.Alive = false
		c.Stage = StageNormal
	}
	w.Factions[i].Corps = 0
	w.Factions[i].Generals = 0
}

// DebugLatchOutcomeForShot 僅供 wlgame 的 -open-outcome 截圖 fixture 使用。
// 正常玩家路徑不呼叫它；production outcome 仍必須由 AdjustTrust 或 capture
// 的真實 mutation 邊界產生。
func (w *World) DebugLatchOutcomeForShot(kind OutcomeKind) {
	w.latchOutcome(kind)
}
