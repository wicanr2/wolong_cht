// Package capital 是首都的選擇與遷都。
//
// 出處是 `sub_16A3D`（選新首都）與兩個呼叫端：
// `sub_14DF0`（首都被打下來就遷）與事件 8（沒仗打時主動遷）。
// 完整推導見 `docs/mechanics/70-ai.md`。
//
// ⚠ **這一支曾經被記成「挑侵攻目標」，錯了。**
// `cmp al, [si+1]` 比的是據點的**擁有者**——它掃的是自己的地盤，
// 不是候選目標。留這句話是因為那個誤讀很自然：
// 一支「掃 192 個據點挑一個」的常式，掛在 AI 決策鏈上，
// 看起來就像在挑目標。
package capital

// Site 是選首都要看的據點欄位。
type Site struct {
	Owner int
	// Kind 是據點類型（記錄 +0x16 的低 4 位）：
	// 0 大城／1 中城／2 小城／3 關／4 戰場。**數字越小城越大。**
	Kind int
	// Production 是目前的生產力（+0x0E），不是上限。
	Production int
	// Adjacency 是鄰接遮罩（+0x00 的低 4 位）。
	Adjacency int
}

// None 是「一個據點都沒有」。原版用 CF ＝ 1 表示。
const None = -1

// Pick 選出 faction 應該定都的據點。
//
// 原版是一次線性掃描，兩個條件都拿**目前最佳值**當門檻：
//
//	類型 ≤ 目前最佳（取最小 ＝ 城最大）
//	生產力 ≥ 目前最佳（取最大）
//
// ⚠ **兩個條件是「且」，而且門檻邊掃邊更新**，所以結果與掃描順序有關。
// 這不是可以整理成「先比類型再比生產力」的排序——照抄掃描。
//
// 還有一個 `test byte ptr [si], 1Fh` 的分支：一旦掃到鄰接遮罩為 0 的據點，
// 之後只有同樣遮罩為 0 的才能再取代它。**在四個劇本的初始資料上，
// 這個分支不改變任何結果**，但照抄——它是原版的行為，
// 而能不能觸發取決於執行期的地盤形狀，不是初始資料說了算。
func Pick(sites []Site, faction int) int {
	bestKind, bestProd, best := 0xFF, 0, None
	latched := false
	for i := range sites {
		s := &sites[i]
		if s.Owner != faction {
			continue
		}
		kind := s.Kind & 0x0F
		if bestKind < kind {
			continue
		}
		if bestProd > s.Production {
			continue
		}
		if !latched {
			bestKind, bestProd, best = kind, s.Production, i
		}
		if s.Adjacency&0x1F != 0 {
			continue
		}
		latched = true
		bestKind, bestProd, best = kind, s.Production, i
	}
	return best
}
