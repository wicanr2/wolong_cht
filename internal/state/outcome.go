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
		// 原版結局送兩則：先 #0x4B（無肖像的捷報），再由君主說一句
		// （組編號 `0x197`，見 VictoryLordTalkIndex）。docs/re/59 §2。
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
			demoteCapturedSovereign(g)
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

// VictoryLordTalkBase 是結局那一句的**組編號**。`sub_11CD0` 送的是
// `cx = 197h`，而 `TALK.DAT` 索引 ≥ `0x196` 是八格一組——展開成
// `0x196 + (0x197 − 0x196) × 8 + 說話型` ＝ **414 ＋ 變體**。
//
// 君主是主公型，說話型落在 0–2，對到的三則正好都是君主對軍師的褒獎；
// 3–7 那五格是空的。**不是 `#407`**——那一則是財政赤字的催促。
const VictoryLordTalkBase = 0x197

// VictoryLordTalkIndex 回傳結局第二則要用的組編號與說話者。
//
// 說話者是玩家所仕勢力的**君主**（`sub_11CD0` 取 `[bx+425Eh]` 的說話型
// 與 `[bx+4241h]` 的肖像，bx 就是君主的武將記錄）。
func (w *World) VictoryLordTalkIndex() (index, general int, ok bool) {
	if w == nil || w.outcome != Victory {
		return 0, 0, false
	}
	if w.Player < 0 || w.Player >= len(w.Factions) {
		return 0, 0, false
	}
	lord := w.Factions[w.Player].Lord
	if lord < 0 || lord >= len(w.Generals) {
		return 0, 0, false
	}
	return VictoryLordTalkBase, lord, true
}

// demoteCapturedSovereign 是 `sub_129C3` 的 `loc_12A12`：被俘的人只要
// 旗標 bit 6（主公型）成立，就清掉它並把說話類型加 3——0／1／2 正好
// 搬到 3／4／5，也就是**主公型換成臣下型**（docs/spec/127）。
//
// ⚠ 這一段**不以「舊主已滅」為條件**。舊主還在時 `jnb loc_12A12` 直接
// 跳過來，舊主已滅但沒有 bit 4（不事二主）也落到這裡；只有走自刎那一條
// 才跳過，而那條的下一步是整筆歸零。
//
// 先測 bit 再清，所以被俘兩次不會加兩次——第二次 bit 6 已經是 0 了。
// 釋放（`sub_150D7`）**不會把它加回去**，這是不可逆的搬移。
func demoteCapturedSovereign(g *General) {
	if g == nil || !g.Sovereign {
		return
	}
	g.Sovereign = false
	g.TalkVariant += 3
}
