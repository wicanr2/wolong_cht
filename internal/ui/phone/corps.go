package phone

import (
	"fmt"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
)

// 軍團：編成與行軍。
//
// 規則全在 `internal/state`（`FormCorps` 是原版 `sub_16F26`、
// `March` 走道路圖）。這一層只決定「怎麼點」。

// corpsForm 是編成中的那一支。leader < 0 表示還沒選武將。
type corpsForm struct {
	leader int
	kinds  [army.Positions]army.TroopType
	manned [army.Positions]bool
}

// newCorpsForm 是預設編成：六個位置都有兵、全部步兵。
//
// ⚠ 預設值會被**直接送進規則層**，所以要挑一個講得出理由的：
// 步兵是三個兵種裡最便宜也最不挑地形的一種，當底比較不會讓玩家
// 在沒注意到的情況下拿到一支昂貴的騎兵團。
func newCorpsForm() corpsForm {
	f := corpsForm{leader: -1}
	for i := range f.kinds {
		f.kinds[i] = army.Infantry
		f.manned[i] = true
	}
	return f
}

// unitLabel 把六個編成位置翻成中文。
//
// ⚠ 順序是**記錄裡的順序**（`docs/spec/21`），不是畫面上的排列。
// 原版戰場底列的空間排列是「左翼 左備 主將 前鋒 右備 右翼」，
// 兩者不同——照畫面排會把兵種配錯位置。
func unitLabel(i int) string {
	return [army.Positions]string{"主將", "前鋒", "左翼", "右翼", "左備", "右備"}[i]
}

// kindLabel 是兵種在畫面上的字，空槽畫成「空」。
func kindLabel(k army.TroopType, manned bool) string {
	if !manned {
		return "空"
	}
	return k.String()
}

// corpsFormRows 是編成頁的內容：還沒選武將時列候選人，
// 選了之後列六個位置與一列「編成」。
func (s *Session) corpsFormRows() []sheetRow {
	if s.form.leader < 0 {
		return s.corpsCandidateRows()
	}
	rows := make([]sheetRow, 0, army.Positions+2)
	rows = append(rows, sheetRow{
		name: s.Localise(s.world.Generals[s.form.leader].Name),
		cols: []string{"換人"},
	})
	for i := 0; i < army.Positions; i++ {
		rows = append(rows, sheetRow{
			name: unitLabel(i),
			cols: []string{kindLabel(s.form.kinds[i], s.form.manned[i])},
		})
	}
	rows = append(rows, sheetRow{name: "編成", cols: []string{s.reserveSummary()}})
	return rows
}

// reserveSummary 是預備兵的三個數字。編成分得到多少兵由它決定
//（`distributeReserves`，docs/spec/21 §2），所以要擺在按鈕旁邊。
func (s *Session) reserveSummary() string {
	p := s.world.Player
	if p < 0 || p >= len(s.world.Factions) {
		return ""
	}
	r := s.world.Factions[p].Reserves
	return fmt.Sprintf("預備兵 騎 %d／弓 %d／步 %d",
		r[army.Cavalry]*MenPerPoint, r[army.Archer]*MenPerPoint, r[army.Infantry]*MenPerPoint)
}

// corpsCandidateRows 是可以帶兵的武將。
//
// ⚠ 已經帶著軍團的人不列——`FormCorps` 會擋，但**讓人點得到一個必定失敗的
// 選項**是介面的錯，不是規則層的錯。
func (s *Session) corpsCandidateRows() []sheetRow {
	rows := make([]sheetRow, 0, 32)
	for _, i := range s.corpsCandidates() {
		g := &s.world.Generals[i]
		rows = append(rows, sheetRow{
			name: s.Localise(g.Name),
			cols: []string{
				fmt.Sprintf("武 %d", g.Martial),
				fmt.Sprintf("統 %d", g.Command),
				s.generalPost(i),
			},
		})
	}
	if len(rows) == 0 {
		rows = append(rows, sheetRow{name: "沒有可以帶兵的武將", dim: true})
	}
	return rows
}

func (s *Session) corpsCandidates() []int {
	out := make([]int, 0, 32)
	for i := range s.world.Generals {
		g := &s.world.Generals[i]
		if !g.Alive || g.Faction != s.world.Player {
			continue
		}
		if i < len(s.world.Corps) && s.world.Corps[i].Alive {
			continue
		}
		out = append(out, i)
	}
	return out
}

// tapCorpsFormRow 點編成頁的一列。
func (s *Session) tapCorpsFormRow(row int) {
	if s.form.leader < 0 {
		ids := s.corpsCandidates()
		if row >= 0 && row < len(ids) {
			f := newCorpsForm()
			f.leader = ids[row]
			s.form = f
		}
		return
	}
	switch {
	case row == 0:
		s.form.leader = -1 // 換人
	case row >= 1 && row <= army.Positions:
		s.cycleSlot(row - 1)
	case row == army.Positions+1:
		s.lastErr = s.world.FormCorps(s.form.leader, s.form.kinds, s.form.manned)
		if s.lastErr == nil {
			s.form = newCorpsForm()
			s.sheet.tab = 0
		}
	}
}

// cycleSlot 把一個位置在「騎馬 → 弓兵 → 步兵 → 空」之間換。
//
// ⚠ **主將那一格不准空**：原版的壞滅判定直接看第一槽的兵力是不是 0，
// 空著的話軍團一編出來就會被判掉（docs/re/09 §5）。
func (s *Session) cycleSlot(i int) {
	if !s.form.manned[i] {
		s.form.manned[i] = true
		s.form.kinds[i] = army.Cavalry
		return
	}
	switch s.form.kinds[i] {
	case army.Cavalry:
		s.form.kinds[i] = army.Archer
	case army.Archer:
		s.form.kinds[i] = army.Infantry
	default:
		if i == 0 {
			s.form.kinds[i] = army.Cavalry // 主將那一格轉回頭，不給空
			return
		}
		s.form.manned[i] = false
	}
}

// MarchSelected 把一支軍團派往目前選中的據點。
//
// ⭐ 與遷都同一個手勢：**先在地圖上點目的地，再開軍團**。
// 手機上面板一開就蓋住地圖，反過來的順序做不到。
func (s *Session) MarchSelected(corps int) {
	if s.selected < 0 {
		s.lastErr = fmt.Errorf("先在地圖上點目的地")
		return
	}
	s.lastErr = s.world.March(corps, s.selected)
}
