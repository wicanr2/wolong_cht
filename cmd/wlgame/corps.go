package main

// 軍團的三個指令：編成、行軍、軍團一覽。
//
// 規則全部在 internal/state 與 internal/rules/army，這裡只做操作介面。
// 說明書 3.2 的指令選單裡「軍隊編成」與「行軍指示」是分開的兩項，
// 而且訊息也分開（`translations` 的 #0 是「進行軍隊編組。請選擇武將。」、
// #2 是「請選擇進行行軍指示之軍團。」）——所以這裡也做成兩條流程。

import (
	"fmt"
	"sort"
	"image"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/listwin"
)

// formState 是編成畫面的狀態：選好武將之後，逐槽指定兵種。
type formState struct {
	active bool
	leader int
	slot   int
	kinds  [army.Positions]army.TroopType
	manned [army.Positions]bool
	// keyboard 記「玩家用過鍵盤沒有」。原版沒有選取狀態——六個槽都是
	// 滑鼠熱區，點一下就循環（docs/re/49 §5）。remake 的鍵盤是額外加的，
	// 所以選取框也只在用過鍵盤之後才畫，滑鼠操作時畫面與原版一致。
	keyboard bool
}

// ---------------------------------------------------------------------------
// 編成
// ---------------------------------------------------------------------------

// formCandidates 是還可以帶兵的武將。
//
// ⚠ **君主在不在裡面由系統選單那一列決定**（docs/spec/76）。
// 原版讓君主帶兵走的是**另一條路**：請求出陣時 `sub_16E8F` 由君主本人
// 帶一支軍團，而且 `sub_16EC9` 專門擋「君主已經帶著軍團」（docs/spec/11）。
// 如果一般編成本來就選得到君主，那條專用路徑與那道擋都沒有存在的必要。
//
// ⚠ **remake 預設放行**（`g.lordCorps` 初值 true），與原版不同——
// 使用者裁定的差異。切成「不可」才是原版行為。
//
// ⚠ **這道擋只在玩家的編成指令上**，`autoFormCorps` 不受影響——
// 出陣那條本來就是要讓君主帶兵。
func (g *game) formCandidates() []int {
	lord := -1
	if !g.lordCorps {
		if p := g.world.Player; p >= 0 && p < len(g.world.Factions) {
			lord = g.world.Factions[p].Lord
		}
	}
	var rows []int
	for i, gen := range g.world.Generals {
		if gen.Alive && gen.Faction == g.world.Player &&
			!gen.Posted && gen.Captor == 0xFF && i != lord {
			rows = append(rows, i)
		}
	}
	return rows
}

// beginForm 開始編成：先選一名還沒帶兵的武將。
//
// **原版是「選武將 → 編成 → 回到選武將」的迴圈**（`sub_16C5E`，docs/re/30 §1），
// 編成視窗畫在武將一覽**上面**，一覽表不會被擦掉。所以這裡選完不關一覽表，
// 編成完或取消就回到它繼續選下一位，右鍵／ESC 才離開整條流程。
func (g *game) beginForm() {
	rows := g.formCandidates()
	if len(rows) == 0 {
		g.lastEvent = "沒有可以帶兵的武將"
		return
	}
	g.openGeneralPicker(rows, "選帶兵的武將　Enter 選取／決定　1-6 排序　ESC 取消", nil)
	g.listPick = func(i int) bool {
		g.form = formState{active: true, leader: i}
		// 預設六槽全滿的騎馬編成——玩家再逐槽改。
		for k := range g.form.manned {
			g.form.manned[k] = true
		}
		return false // 一覽表留在背景（原版不擦掉它）
	}
}

// refreshFormCandidates 在編成完成之後更新背景那張一覽表——
// 剛帶兵的那位已經 `Posted`，不該還留在候選裡。
func (g *game) refreshFormCandidates() {
	if g.list == nil {
		return
	}
	rows := g.formCandidates()
	if len(rows) == 0 {
		g.list = nil
		return
	}
	g.list.Rows = rows
	if g.list.Cursor >= len(rows) {
		g.list.Cursor = len(rows) - 1
	}
	if g.list.Top > g.list.Cursor {
		g.list.Top = g.list.Cursor
	}
}

// formCycleKind 是原版點一下槽的動作：**兵種 +1，`1 → 2 → 3 → 1` 循環**
// （docs/re/30 §3）。空槽的兵種是 4，`inc` 之後也落回 1，
// 所以**點過的槽不會再變回空槽**——這與鍵盤的空白鍵不同。
func (f *formState) cycleKind(k int) {
	if k < 0 || k >= army.Positions {
		return
	}
	if !f.manned[k] {
		f.kinds[k], f.manned[k] = army.Cavalry, true
		return
	}
	f.kinds[k] = (f.kinds[k] + 1) % 3
}

// updateForm 是編成畫面的輸入。滑鼠照原版的熱區（docs/spec/22 §2），
// 鍵盤是 remake 加的。
func (g *game) updateForm() {
	f := &g.form
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		f.active = false // 原版 sub_121E7 回 CF=1 ＝ 右鍵取消
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if k, ok := formSlotAt(x, y); ok {
			f.keyboard, f.slot = false, k
			f.cycleKind(k)
			return
		}
		if image.Pt(x, y).In(formOKRect()) {
			f.keyboard = false
			g.commitForm()
			return
		}
		return
	}
	switch {
	case pressed(ebiten.KeyEscape):
		f.active = false
	case pressed(ebiten.KeyArrowUp):
		f.keyboard = true
		f.slot = (f.slot + army.Positions - 1) % army.Positions
	case pressed(ebiten.KeyArrowDown):
		f.keyboard = true
		f.slot = (f.slot + 1) % army.Positions
	case pressed(ebiten.KeySpace):
		f.keyboard = true
		f.manned[f.slot] = !f.manned[f.slot]
	case pressed(ebiten.Key1):
		f.keyboard = true
		f.kinds[f.slot], f.manned[f.slot] = army.Cavalry, true
	case pressed(ebiten.Key2):
		f.keyboard = true
		f.kinds[f.slot], f.manned[f.slot] = army.Archer, true
	case pressed(ebiten.Key3):
		f.keyboard = true
		f.kinds[f.slot], f.manned[f.slot] = army.Infantry, true
	case pressed(ebiten.KeyEnter):
		f.keyboard = true
		g.commitForm()
	}
}

// commitForm 是「確定」：編成成功就關掉編成視窗**回到武將一覽**
// （原版 sub_16C5E 的迴圈，docs/re/30 §1）。主將槽沒兵時原版走
// TALK #62，remake 照走同一則訊息。
func (g *game) commitForm() {
	f := &g.form
	// 原版的確定鈕只有一種拒絕理由：主將槽沒有兵（docs/re/30 §3.1），
	// 所以 state 回什麼錯誤都走同一則訊息。
	if err := g.world.FormCorps(f.leader, f.kinds, f.manned); err != nil {
		g.enqueueTalk(formNoLeaderTroopsTalk, nil)
		return
	}
	g.lastEvent = big5(g.world.Generals[f.leader].Name) + " 編成完畢"
	f.active = false
	g.refreshFormCandidates()
}

// formNoLeaderTroopsTalk 是 TALK #62「大將的部隊的兵員不足啊。」，
// 原版在確定鈕上主將槽為空時走 `sub_18810(cx=3Eh)`（docs/re/30 §3.1）。
const formNoLeaderTroopsTalk = 0x3E

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

// formSlotAt 反查點到哪一個槽。
func formSlotAt(x, y int) (int, bool) {
	p := image.Pt(x, y)
	for k := 0; k < army.Positions; k++ {
		if p.In(formSlotRect(k)) {
			return k, true
		}
	}
	return 0, false
}

// formOKRect 是確定鈕（原版熱區 0x44），與它的底框逐格重合。
func formOKRect() image.Rectangle {
	return image.Rect(formOKX, formOKY, formOKX+formOKW, formOKY+formOKH)
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
	season := int(g.world.Clock.Season())

	// 靜態層（顯示清單場景 5）。
	for _, box := range []struct{ x, y, w, h int }{
		{240, 152, 112, 32}, {formSlotLabelX, formSlotY, 112, 96},
		{304, formReserveY, 64, 48},
	} {
		vector.DrawFilledRect(screen, float32(box.x), float32(box.y),
			float32(box.w), float32(box.h), color.Black, false)
	}
	// 確定鈕是**按鈕**：底色 7 ＋ 一圈 9／6（docs/re/48 §2.1）。
	g.dlButton(screen, formOKX, formOKY, formOKW, formOKH)
	g.td.Draw(screen, "將軍", formHeadLabelX, formTitleY, ink)
	g.td.Draw(screen, "總兵力", formHeadLabelX, formTotalY, ink)
	g.td.Draw(screen, "士氣值", formHeadLabelX, formMoraleY, ink)
	g.td.Draw(screen, "預備兵數", formReserveLabelX, formReserveLabelY, labelInk)
	g.td.Draw(screen, "確 定", formOKTextX, formOKY, g.dlButtonInk())
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
	//
	// 頭像走 `sub_107D2`（docs/re/33 §2），頁碼取武將記錄的 +0x01
	// 不是武將編號。原版那裡是四格 round-robin 快取，remake 一次全載入——
	// **只差在載入時機，畫出來的東西相同**。
	gen := g.world.Generals[f.leader]
	if img, err := g.lib.Portrait(gen.Portrait, season); err == nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(formPortraitX), float64(formPortraitY))
		screen.DrawImage(ebiten.NewImageFromImage(img), op)
	}
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

	// ↓ **remake 差異**：原版沒有選取狀態——六個槽都是滑鼠熱區，
	// 點一下兵種就 +1。所以選取框只在玩家用過鍵盤之後才畫。
	if f.keyboard {
		sel := formSlotRect(f.slot)
		vector.StrokeRect(screen, float32(sel.Min.X-1), float32(sel.Min.Y-1),
			float32(sel.Dx()+2), float32(sel.Dy()+2), 1, ink, false)
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
	// 預設照距離排：一張 192 列的表按編號排，玩家要翻半天才找得到隔壁。
	// **原版的欄位裡沒有「距離」**（docs/re/26 §4.1），所以只排順序、
	// 不加欄——這是 remake 的便利，不是原版行為。
	sort.SliceStable(rows, func(a, b int) bool { return dist(rows[a]) < dist(rows[b]) })
	g.openCityPicker(rows, "選擇目的地　Enter 選取／決定　1-6 排序　ESC 取消", nil)
	g.listPick = func(i int) bool {
		if err := g.world.March(corps, i); err != nil {
			g.setEvent(err.Error())
			return true
		}
		// 原版選完據點還有第二段：戰鬥指揮／委任／解體（docs/spec/39）。
		g.beginMarchMode(corps, i)
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
		func(corps int) bool {
			g.openCorpsInfo(corps) // 原版的軍團情報視窗（docs/spec/24）
			return true
		})
}

func (g *game) openCorpsListWith(rows []int, hint string, pick func(int) bool) {
	g.list = listwin.New(listwin.Corps, g.listColumnsCorps(), rows,
		listRowsPerPage, &g.sortMem)
	g.listTitle = listFamilyCorps.Title
	g.listRow = g.listRowCorps
	// 兩條換色（docs/re/27 §5）：**總兵數 < 300 點**（＝ 3,000 人，半編）
	// 與**士氣 < 100**。兩者都是「值得注意」不是錯誤。
	g.listCellInk = func(id, col int) (color.RGBA, bool) {
		c := g.world.Corps[id]
		switch {
		case col == 1 && c.Men < corpsHalfStrength:
			return listWarnInk, true
		case col == 2 && c.Morale < 100:
			return listWarnInk, true
		}
		return color.RGBA{}, false
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
		g.reportGovernorReturn(e)
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

// reportGovernorReturn 是 `sub_14D63` 的兩則訊息（docs/spec/48）：
// 據點被攻陷時派駐的內政官回來了，先跳一般通知，再由他自己說一句。
//
// 第二則的索引在八格變體的範圍裡，要走 `resolveBattleTalkIndex`——
// 直接拿 `0x1A6` 當索引會落到 422「．．．．」那一組去。
func (g *game) reportGovernorReturn(e state.CorpsEvent) {
	id := e.GovernorReturned
	if id < 0 || id >= len(g.world.Generals) {
		return
	}
	gen := &g.world.Generals[id]
	if gen.Faction != g.world.Player {
		return // 只報自己人的，二十二個勢力天天在丟城
	}
	city := ""
	if e.Captured >= 0 && e.Captured < len(g.world.Cities) {
		city = big5(g.world.Cities[e.Captured].Name)
	}
	vars := map[byte]string{'1': big5(gen.Name), '2': city, '6': ""}
	g.enqueueTalk(governorReturnTalk, vars)
	g.enqueueTalkWithPortrait(
		resolveBattleTalkIndex(governorRegretTalkBase, gen.TalkVariant),
		vars, gen.Portrait)
}

const (
	// governorReturnTalk 是「{2}內政官的{1}大人因為據點被攻陷而歸來了。」
	governorReturnTalk = 0x44
	// governorRegretTalkBase 是內政官自己那一句的**組編號**（不是索引）。
	// 實際落在 534–541（docs/spec/48 §2）。
	governorRegretTalkBase = 0x1A6
)

// demoCorps 是**驗收用**的捷徑：直接把畫面帶到編成或軍團一覽，
// 免得截圖前要按一長串鍵。正常玩不會走到這裡。
//
// 編成的內容照**實際的預備兵**湊——開局的勢力大多湊不滿六個位置，
// 驗收畫面就該長成那個樣子，不然截出來的是假的。
func (g *game) demoCorps(list bool) {
	// ⚠ **不要在這裡自己抄一份資格判定。** 先前這支自己掃 Generals，
	// 於是 `formCandidates()` 把君主擋掉之後，驗收 fixture 仍然把曹操
	// 編成軍團長——修好的規則沒有套到驗收路徑上，而截圖看起來像沒修。
	// 一條規則只留一份實作（CLAUDE.md §7 第 6 條）。
	leaders := g.formCandidates()
	if len(leaders) == 0 {
		return
	}
	if len(leaders) > 2 {
		leaders = leaders[:2]
	}
	kinds, manned := g.affordable()

	if !list {
		// 編成畫面：走**真實流程**（選武將 → 編成），這樣截圖裡的
		// 武將一覽會像原版一樣留在編成視窗底下。
		// **不要先幫他編好**，否則按確定會撞到「已經帶著軍團」。
		g.beginForm()
		if g.list == nil {
			return
		}
		g.confirmListSelection() // 反白
		g.confirmListSelection() // 決定
		if g.form.active {
			g.form.slot, g.form.kinds, g.form.manned = 2, kinds, manned
		}
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
