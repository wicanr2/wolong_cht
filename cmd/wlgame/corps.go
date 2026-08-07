package main

// 軍團的三個指令：編成、行軍、軍團一覽。
//
// 規則全部在 internal/state 與 internal/rules/army，這裡只做操作介面。
// 說明書 3.2 的指令選單裡「軍隊編成」與「行軍指示」是分開的兩項，
// 而且訊息也分開（`translations` 的 #0 是「進行軍隊編組。請選擇武將。」、
// #2 是「請選擇進行行軍指示之軍團。」）——所以這裡也做成兩條流程。

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/listwin"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// formState 是編成畫面的狀態：選好武將之後，逐槽指定兵種。
type formState struct {
	active bool
	leader int
	slot   int
	kinds  [army.Positions]army.TroopType
	manned [army.Positions]bool
	err    string
}

// ---------------------------------------------------------------------------
// 編成
// ---------------------------------------------------------------------------

// beginForm 開始編成：先選一名還沒帶兵的武將。
func (g *game) beginForm() {
	var rows []int
	for i, gen := range g.world.Generals {
		if gen.Alive && gen.Faction == g.world.Player &&
			!gen.Posted && gen.Captor == 0xFF {
			rows = append(rows, i)
		}
	}
	if len(rows) == 0 {
		g.lastEvent = "沒有可以帶兵的武將"
		return
	}
	gs := g.world.Generals
	g.list = listwin.New(listwin.Generals, []listwin.Column{
		{Title: "武將名", Less: func(a, b int) bool { return gs[a].Name < gs[b].Name }},
		{Title: "武術", Less: func(a, b int) bool { return gs[a].Martial > gs[b].Martial }},
		{Title: "統率", Less: func(a, b int) bool { return gs[a].Command > gs[b].Command }},
		{Title: "攻城", Less: func(a, b int) bool { return gs[a].Aptitude[0] > gs[b].Aptitude[0] }},
		{Title: "野戰", Less: func(a, b int) bool { return gs[a].Aptitude[1] > gs[b].Aptitude[1] }},
	}, rows, 12, &g.sortMem)
	g.listRow = func(i int) (string, string) {
		gen := gs[i]
		// 適性是**場合適性**不是兵種適性（docs/re/09 §3.3），
		// 所以這裡的欄位名是「攻城／野戰」。
		return big5(gen.Name), fmt.Sprintf("%4d%8d%8d%8d",
			gen.Martial, gen.Command, gen.Aptitude[0], gen.Aptitude[1])
	}
	g.listPick = func(i int) bool {
		g.form = formState{active: true, leader: i}
		// 預設六槽全滿的騎馬編成——玩家再逐槽改。
		for k := range g.form.manned {
			g.form.manned[k] = true
		}
		return true
	}
	g.listHint = "選擇帶兵的武將　Enter 選取／決定　1-5 排序　ESC 取消"
}

// updateForm 是編成畫面的輸入。
func (g *game) updateForm() {
	f := &g.form
	switch {
	case pressed(ebiten.KeyEscape):
		f.active = false
	case pressed(ebiten.KeyArrowUp):
		f.slot = (f.slot + army.Positions - 1) % army.Positions
	case pressed(ebiten.KeyArrowDown):
		f.slot = (f.slot + 1) % army.Positions
	case pressed(ebiten.KeySpace):
		f.manned[f.slot] = !f.manned[f.slot]
	case pressed(ebiten.Key1):
		f.kinds[f.slot], f.manned[f.slot] = army.Cavalry, true
	case pressed(ebiten.Key2):
		f.kinds[f.slot], f.manned[f.slot] = army.Archer, true
	case pressed(ebiten.Key3):
		f.kinds[f.slot], f.manned[f.slot] = army.Infantry, true
	case pressed(ebiten.KeyEnter):
		if err := g.world.FormCorps(f.leader, f.kinds, f.manned); err != nil {
			f.err = plainErr(err)
			return
		}
		g.lastEvent = big5(g.world.Generals[f.leader].Name) + " 編成完畢"
		f.active = false
	}
}

// drawForm 畫編成畫面。六個位置的名稱照說明書 5.5 的編成畫面。
func (g *game) drawForm(screen *ebiten.Image) {
	f := &g.form
	if !f.active {
		return
	}
	const x, y, w, h = 40, 52, 400, 214
	vector.DrawFilledRect(screen, x, y, w, h, color.RGBA{0, 0, 0, 230}, false)
	vector.StrokeRect(screen, x, y, w, h, 1, color.RGBA{240, 200, 120, 255}, false)

	white := color.RGBA{240, 240, 230, 255}
	amber := color.RGBA{240, 200, 120, 255}
	dim := color.RGBA{150, 150, 160, 255}
	red := color.RGBA{240, 140, 140, 255}

	g.td.Draw(screen, "軍隊編成　"+big5(g.world.Generals[f.leader].Name), x+8, y+6, amber)

	// 剩餘預備兵。編成一個位置固定扣 1,000 人（說明書 5.5）。
	res := g.world.Factions[g.world.Player].Reserves
	g.td.Draw(screen, fmt.Sprintf("預備兵　騎馬%5d　弓兵%5d　步兵%5d",
		res[0], res[1], res[2]), x+8, y+26, dim)

	ry := y + 50
	men := 0
	for k := 0; k < army.Positions; k++ {
		col := white
		if k == f.slot {
			vector.DrawFilledRect(screen, x+4, float32(ry-1), w-8,
				float32(textdraw.GlyphH+2), color.RGBA{70, 60, 30, 255}, false)
			col = amber
		}
		kind := "（空）"
		if f.manned[k] {
			kind = f.kinds[k].String()
			men += army.MenPerUnit
		}
		g.td.Draw(screen, army.Position(k).String()+"　"+kind, x+12, ry, col)
		ry += textdraw.GlyphH + 2
	}
	g.td.Draw(screen, fmt.Sprintf("合計 %d 人", men), x+240, y+50, white)

	if f.err != "" {
		g.td.Draw(screen, f.err, x+8, y+h-3*textdraw.GlyphH-12, red)
	}
	// 提示分兩列——一列塞不下就會畫到視窗外面去。
	g.td.Draw(screen, "↑↓ 選位置　1 騎馬　2 弓兵　3 步兵　空白 空位",
		x+8, y+h-2*textdraw.GlyphH-8, dim)
	g.td.Draw(screen, "Enter 編成　ESC 取消",
		x+8, y+h-textdraw.GlyphH-4, dim)
}

// ---------------------------------------------------------------------------
// 行軍
// ---------------------------------------------------------------------------

// beginMarch 開始行軍：先選軍團，再選目的地。
func (g *game) beginMarch() {
	rows := g.playerCorps()
	if len(rows) == 0 {
		g.lastEvent = "沒有軍團可以行軍"
		return
	}
	g.openCorpsListWith(rows, "選擇行軍的軍團　Enter 選取／決定　ESC 取消",
		func(i int) bool {
			g.pickDestination(i)
			return false // 直接換成下一張一覽表，不關視窗
		})
}

// pickDestination 選行軍的目的地。
//
// 全部 192 個據點都列出來，但**預設照距離排序**——一張 192 列的表
// 若按編號排，玩家要翻半天才找得到隔壁那座城。
func (g *game) pickDestination(corps int) {
	cs := g.world.Cities
	from := g.world.Corps[corps]
	dist := func(i int) int {
		dx, dy := cs[i].X-from.X, cs[i].Y-from.Y
		if dx < 0 {
			dx = -dx
		}
		if dy < 0 {
			dy = -dy
		}
		if dx > dy {
			return dx
		}
		return dy // 切比雪夫距離，與月結收入用的同一種
	}
	var rows []int
	for i := range cs {
		if i != from.Node {
			rows = append(rows, i)
		}
	}
	g.list = listwin.New(listwin.Cities, []listwin.Column{
		{Title: "據點", Less: func(a, b int) bool { return cs[a].Name < cs[b].Name }},
		{Title: "距離", Less: func(a, b int) bool { return dist(a) < dist(b) }},
		{Title: "所屬", Less: func(a, b int) bool { return cs[a].Owner < cs[b].Owner }},
		{Title: "城兵", Less: func(a, b int) bool { return cs[a].Garrison > cs[b].Garrison }},
	}, rows, 12, &g.sortMem)
	g.list.SortBy(1) // 預設照距離排，最近的在最上面
	g.listRow = func(i int) (string, string) {
		owner := "無主"
		if o := cs[i].Owner; o >= 0 && o < 22 && g.world.Factions[o].Alive {
			owner = big5(g.world.LordName(o))
		}
		return big5(cs[i].Name), fmt.Sprintf("%4d  %s  %4d", dist(i), owner, cs[i].Garrison)
	}
	g.listPick = func(i int) bool {
		if err := g.world.March(corps, i); err != nil {
			g.setEvent(err.Error())
			return true
		}
		g.lastEvent = big5(g.world.Generals[corps].Name) + " 向 " +
			big5(cs[i].Name) + " 行軍"
		return true
	}
	g.listHint = "選擇目的地　Enter 選取／決定　1-4 排序　ESC 取消"
}

// ---------------------------------------------------------------------------
// 軍團一覽
// ---------------------------------------------------------------------------

func (g *game) playerCorps() []int {
	var rows []int
	for _, i := range g.world.AliveCorps() {
		if g.world.Corps[i].Faction == g.world.Player {
			rows = append(rows, i)
		}
	}
	return rows
}

func (g *game) openCorpsList() {
	rows := g.playerCorps()
	if len(rows) == 0 {
		g.lastEvent = "還沒有軍團"
		return
	}
	g.openCorpsListWith(rows, "↑↓ 移動　Enter 選取／決定　1-5 排序　ESC 取消",
		func(int) bool { return true })
}

func (g *game) openCorpsListWith(rows []int, hint string, pick func(int) bool) {
	cp := g.world.Corps
	gs := g.world.Generals
	g.list = listwin.New(listwin.Corps, []listwin.Column{
		{Title: "大將", Less: func(a, b int) bool { return gs[a].Name < gs[b].Name }},
		{Title: "兵力", Less: func(a, b int) bool { return cp[a].Men > cp[b].Men }},
		{Title: "士氣", Less: func(a, b int) bool { return cp[a].Morale > cp[b].Morale }},
		{Title: "所在", Less: func(a, b int) bool { return cp[a].Node < cp[b].Node }},
		{Title: "目標", Less: func(a, b int) bool { return cp[a].TargetNode < cp[b].TargetNode }},
	}, rows, 12, &g.sortMem)
	g.listRow = func(i int) (string, string) {
		c := cp[i]
		// 一點兵力 ＝ 10 人。士氣低於 100 的軍團再戰必壞滅
		// （docs/re/09 §4.4），所以這個數字要看得到。
		return big5(gs[i].Name), fmt.Sprintf("%5d %4d  %s → %s",
			c.Men*10, c.Morale,
			big5(g.world.Cities[c.Node].Name),
			big5(g.world.Cities[c.TargetNode].Name))
	}
	g.listPick = pick
	g.listHint = hint
}

// ---------------------------------------------------------------------------
// 事件回報
// ---------------------------------------------------------------------------

// reportCorps 把這個 tick 的軍團事件寫進狀態列。
// 只報**與玩家有關**的——二十二個勢力天天在打，全報會刷屏。
func (g *game) reportCorps(ev state.Event) {
	for _, e := range ev.Corps {
		if e.Battle == nil {
			continue
		}
		mine := g.world.Corps[e.Corps].Faction == g.world.Player ||
			(e.Enemy >= 0 && g.world.Corps[e.Enemy].Faction == g.world.Player)
		for _, d := range e.Destroyed {
			if g.world.Generals[d].Faction == g.world.Player {
				mine = true
			}
		}
		if !mine {
			continue
		}
		g.lastEvent = battleLine(g, e)
	}
}

func battleLine(g *game, e state.CorpsEvent) string {
	who := big5(g.world.Generals[e.Corps].Name)
	against := "城兵"
	if e.Enemy >= 0 {
		against = big5(g.world.Generals[e.Enemy].Name)
	}
	line := who + " 對 " + against
	for _, d := range e.Destroyed {
		name := big5(g.world.Generals[d].Name)
		switch e.Fate[d] {
		case combat.Escaped:
			line += "　" + name + " 部隊壞滅（脫身）"
		case combat.Captured:
			line += "　" + name + " 被擒"
		case combat.Suicide:
			line += "　" + name + " 自刎"
		}
	}
	if e.Captured >= 0 {
		line += "　攻下 " + big5(g.world.Cities[e.Captured].Name)
	}
	return line
}

// demoCorps 是**驗收用**的捷徑：直接把畫面帶到編成或軍團一覽，
// 免得截圖前要按一長串鍵。正常玩不會走到這裡。
//
// 編成的內容照**實際的預備兵**湊——開局的勢力大多湊不滿六個位置，
// 驗收畫面就該長成那個樣子，不然截出來的是假的。
func (g *game) demoCorps(list bool) {
	p := g.world.Player
	var leaders []int
	for i, gen := range g.world.Generals {
		if gen.Alive && gen.Faction == p && !gen.Posted {
			leaders = append(leaders, i)
		}
		if len(leaders) == 2 {
			break
		}
	}
	if len(leaders) == 0 {
		return
	}
	kinds, manned := g.affordable()

	if !list {
		// 編成畫面：**不要先幫他編好**，否則按 Enter 會撞到
		// 「已經帶著軍團」而看不到正常流程。
		g.form = formState{active: true, leader: leaders[0], slot: 2,
			kinds: kinds, manned: manned}
		return
	}
	for _, l := range leaders {
		if !manned[0] {
			break // 預備兵湊不出大將那一隊了，不是錯誤
		}
		if err := g.world.FormCorps(l, kinds, manned); err != nil {
			g.setEvent(err.Error())
			break
		}
		kinds, manned = g.affordable()
	}
	g.openCorpsList()
	if g.list != nil {
		g.list.Move(1)
	}
}

// affordable 依現有的預備兵湊一個編成：騎馬優先，其次弓兵、步兵。
func (g *game) affordable() (kinds [army.Positions]army.TroopType, manned [army.Positions]bool) {
	res := g.world.Factions[g.world.Player].Reserves
	slot := 0
	for t := army.Cavalry; t <= army.Infantry && slot < army.Positions; t++ {
		for res[t] >= army.MenPerUnit && slot < army.Positions {
			kinds[slot], manned[slot] = t, true
			res[t] -= army.MenPerUnit
			slot++
		}
	}
	return kinds, manned
}

// setEvent 把訊息放到事件列。錯誤訊息帶著 Go 的套件前綴（`state: …`），
// 那是給 log 看的，不是給玩家看的。
func (g *game) setEvent(msg string) { g.lastEvent = strings.TrimPrefix(msg, "state: ") }

func plainErr(err error) string { return strings.TrimPrefix(err.Error(), "state: ") }
