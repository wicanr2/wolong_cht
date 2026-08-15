package main

// 軍團的三個指令：編成、行軍、軍團一覽。
//
// 規則全部在 internal/state 與 internal/rules/army，這裡只做操作介面。
// 說明書 3.2 的指令選單裡「軍隊編成」與「行軍指示」是分開的兩項，
// 而且訊息也分開（`translations` 的 #0 是「進行軍隊編組。請選擇武將。」、
// #2 是「請選擇進行行軍指示之軍團。」）——所以這裡也做成兩條流程。

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
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

// 編成視窗的版面**全部出自原版**（docs/spec/22）：視窗矩形來自
// `sub_1895D(cx=0C0Fh)`，靜態層是顯示清單場景 5，數值座標由
// `sub_16D6F`／`sub_16DA8` 的 VRAM 位移換算（一列 80 byte）。
const (
	formWinX, formWinY = 144, 112
	formWinW, formWinH = 240, 192

	formPortraitX, formPortraitY = 152, 120
	formNameX, formNameY         = 296, 128

	formHeadLabelX               = 248
	formTitleY                   = 128
	formTotalY, formMoraleY      = 152, 168
	formTotalValueX              = 312
	formMoraleValueX             = 320
	formTotalDigits              = 4
	formMoraleDigits             = 3

	// 六個槽：標籤 → 兵種圖示 → 兵力，一列三段。
	formSlotLabelX = 160
	formSlotIconX  = 200
	formSlotValueX = 232
	formSlotY      = 192
	formSlotStep   = 16
	formSlotDigits = 4

	// 右側預備兵欄：圖示在 280，數字在 312，三列。
	formReserveLabelX, formReserveLabelY = 280, 192
	formReserveIconX                     = 280
	formReserveValueX                    = 312
	formReserveY                         = 216
	formReserveDigits                    = 6

	formOKX, formOKY = 280, 272
	formOKW, formOKH = 88, 16
	formOKTextX      = 304

	formIconW, formIconH = 24, 16

	// remake 差異：操作提示自己一個框，接在原版視窗下面。
	// 四列：錯誤訊息一列、操作三列——一列塞不下 240 的框寬。
	formHintY = formWinY + formWinH
	formHintH = 80
)

// formSlotLabels 是六個槽的標籤，取自顯示清單場景 5 的字串。
//
// **不用 `army.Position.String()`**：那一組是規則層的用語，第一個是「大將」
// （原版 TALK #62 也這樣說），而編成視窗上印的是「主將」。
var formSlotLabels = [army.Positions]string{"主將", "前鋒", "左翼", "右翼", "左備", "右備"}

// formSlotRect 是第 k 個槽的可點矩形（原版熱區 0x3E+k），
// **與兵種圖示逐格重合**。
func formSlotRect(k int) image.Rectangle {
	if k < 0 || k >= army.Positions {
		return image.Rectangle{}
	}
	y := formSlotY + k*formSlotStep
	return image.Rect(formSlotIconX, y, formSlotIconX+formIconW, y+formIconH)
}

// drawForm 畫編成畫面（docs/spec/22）。
func (g *game) drawForm(screen *ebiten.Image) {
	f := &g.form
	if !f.active {
		return
	}
	g.chrome.Window(screen, formWinX, formWinY, formWinW, formWinH, chrome.Menu)

	ink := g.paletteInk(strategyInkNormal, chrome.Paper)
	labelInk := g.paletteInk(strategyInkDim, color.RGBA{255, 223, 154, 255})
	warnInk := g.paletteInk(strategyInkGauge, color.RGBA{210, 48, 40, 255})
	season := int(g.world.Clock.Season())

	// 靜態層（顯示清單場景 5）。
	for _, box := range []struct{ x, y, w, h int }{
		{240, 152, 112, 32}, {formSlotLabelX, formSlotY, 112, 96},
		{304, formReserveY, 64, 48}, {formOKX, formOKY, formOKW, formOKH},
	} {
		vector.DrawFilledRect(screen, float32(box.x), float32(box.y),
			float32(box.w), float32(box.h), color.Black, false)
	}
	g.td.Draw(screen, "將軍", formHeadLabelX, formTitleY, ink)
	g.td.Draw(screen, "總兵力", formHeadLabelX, formTotalY, ink)
	g.td.Draw(screen, "士氣值", formHeadLabelX, formMoraleY, ink)
	g.td.Draw(screen, "預備兵數", formReserveLabelX, formReserveLabelY, labelInk)
	g.td.Draw(screen, "確 定", formOKTextX, formOKY, ink)
	for k := 0; k < army.Positions; k++ {
		g.td.Draw(screen, formSlotLabels[k],
			formSlotLabelX, formSlotY+k*formSlotStep, ink)
	}

	// 兵力照 `sub_14698` 的分配式預覽，**不是「槽數 × 1000」**——
	// 兵不夠時原版就是分多少算多少（docs/spec/21 §2）。
	men := g.world.PreviewFormation(g.world.Player, f.kinds, f.manned)

	// 右側預備兵欄：三張紅色圖示（財政那一組的第 2–4 張）＋ 六位數。
	//
	// 顯示的是**扣掉這次編成之後**的餘額。原版每改一次兵種就真的退回池再
	// 重分，所以它畫面上那個數字已經扣過了；remake 到按確定才落地，
	// 這裡要自己減，否則玩家會看到一筆同時算在兩個地方的兵。
	res := g.world.Factions[g.world.Player].Reserves
	for k := 0; k < army.Positions; k++ {
		if t := int(f.kinds[k]); f.manned[k] && t >= 0 && t < len(res) {
			res[t] -= men[k]
		}
	}
	for i := 0; i < 3; i++ {
		y := formReserveY + i*formSlotStep
		if img, err := g.lib.DOSVResourceIcon(i+1, false, season); err == nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(formReserveIconX), float64(y))
			screen.DrawImage(ebiten.NewImageFromImage(img), op)
		}
		g.td.Draw(screen, strategyHUDNumber(res[i]*strategyReserveMenPerPoint,
			formReserveDigits), formReserveValueX, y, labelInk)
	}

	// 動態層：主將、士氣、六個槽與總兵力。
	gen := g.world.Generals[f.leader]
	g.td.Draw(screen, big5(gen.Name), formNameX, formNameY, ink)
	g.td.Draw(screen, strategyHUDNumber(
		g.world.Factions[g.world.Player].MoraleBase, formMoraleDigits),
		formMoraleValueX, formMoraleY, ink)

	total := 0
	for k := 0; k < army.Positions; k++ {
		y := formSlotY + k*formSlotStep
		kind := 4 // 原版的兵種 4 ＝ 空槽
		if f.manned[k] {
			kind = int(f.kinds[k]) + 1
		}
		if img, err := g.lib.DOSVTroopIcon(kind, season); err == nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(formSlotIconX), float64(y))
			screen.DrawImage(ebiten.NewImageFromImage(img), op)
		}
		g.td.Draw(screen, strategyHUDNumber(men[k]*strategyReserveMenPerPoint,
			formSlotDigits), formSlotValueX, y, labelInk)
		total += men[k]
	}
	g.td.Draw(screen, strategyHUDNumber(total*strategyReserveMenPerPoint,
		formTotalDigits), formTotalValueX, formTotalY, ink)

	// ↓ 以下是 **remake 差異**，原版沒有：選取標記與操作提示。
	// 原版點一下槽就是兵種 +1 循環，改完立刻重算，所以不需要選取狀態。
	sel := formSlotRect(f.slot)
	vector.StrokeRect(screen, float32(sel.Min.X-1), float32(sel.Min.Y-1),
		float32(sel.Dx()+2), float32(sel.Dy()+2), 1, ink, false)

	g.chrome.Window(screen, formWinX, formHintY, formWinW, formHintH, chrome.Menu)
	hy := formHintY + 8
	if f.err != "" {
		g.td.Draw(screen, f.err, formWinX+8, hy, warnInk)
	}
	for i, line := range []string{
		"↑↓ 選位置　1 騎馬",
		"2 弓兵　3 步兵　空白 空位",
		"Enter 編成　ESC 取消",
	} {
		g.td.Draw(screen, line, formWinX+8,
			hy+(i+1)*(textdraw.GlyphH+2), labelInk)
	}
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

// reportStrategy 只把涉及玩家的政略動作放進狀態列；其餘勢力的月度宣戰
// 不應每幀刷屏，但敵人正式把玩家列為目標時必須讓正常遊戲看得見。
func (g *game) reportStrategy(ev state.Event) {
	for _, s := range ev.Strategy {
		if s.Target != g.world.Player || s.Faction < 0 || s.Faction >= len(g.world.Factions) {
			continue
		}
		lord := big5(g.world.LordName(s.Faction))
		if s.Corps < 0 {
			g.lastEvent = lord + " 對我方宣戰"
			continue
		}
		if s.Destination >= 0 && s.Destination < len(g.world.Cities) {
			g.lastEvent = lord + " 軍團向 " + big5(g.world.Cities[s.Destination].Name) + " 行軍"
		}
	}
}

func battleLine(g *game, e state.CorpsEvent) string {
	who := big5(g.world.Generals[e.Corps].Name)
	against := "城兵"
	if e.Enemy >= 0 {
		against = big5(g.world.Generals[e.Enemy].Name)
	}
	line := who + " 對 " + against
	if e.Battle != nil {
		winner := "攻方勝"
		if e.Battle.DefenderWins {
			winner = "守方勝"
		}
		line += fmt.Sprintf("　%s　兵力 %d→%d／%d→%d",
			winner,
			e.BattleBefore[0]*10, e.BattleAfter[0]*10,
			e.BattleBefore[1]*10, e.BattleAfter[1]*10)
		if e.BattleCityDamage > 0 {
			line += fmt.Sprintf("　據點損害 %d", e.BattleCityDamage)
		}
	}
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
	// **池的單位是點，不是人**（docs/spec/21 §1）：一個滿編槽是
	// `state.MaxMenPerSlot` ＝ 100 點 ＝ 1,000 人。拿 `army.MenPerUnit`
	// 來比會把門檻抬高十倍，六個槽裡只填得出一個。
	res := g.world.Factions[g.world.Player].Reserves
	slot := 0
	for t := army.Cavalry; t <= army.Infantry && slot < army.Positions; t++ {
		for res[t] >= state.MaxMenPerSlot && slot < army.Positions {
			kinds[slot], manned[slot] = t, true
			res[t] -= state.MaxMenPerSlot
			slot++
		}
	}
	return kinds, manned
}

// setEvent 把訊息放到事件列。錯誤訊息帶著 Go 的套件前綴（`state: …`），
// 那是給 log 看的，不是給玩家看的。
func (g *game) setEvent(msg string) { g.lastEvent = strings.TrimPrefix(msg, "state: ") }

func plainErr(err error) string { return strings.TrimPrefix(err.Error(), "state: ") }
