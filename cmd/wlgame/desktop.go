package main

// 桌面外殼與原版畫面分開。規格與使用者裁定見 docs/spec/266。
import (
	"encoding/json"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/langpack"
	"github.com/wicanr2/wolong_cht/translations"
)

const desktopBarHeight = 28

// 對拍與錄影仍直接使用 game，座標偏移保持零。
var desktopGameOffset int

type desktopAction int

type mouseHintID int

const (
	mouseSelect mouseHintID = iota
	mouseContinue
	mouseReturnTitle
	mouseReturnMap
	mouseChooseRuler
)

var desktopLabels = loadDesktopLabels()

func loadDesktopLabels() map[string][]string {
	data, err := translations.Files.ReadFile("desktop-controls.json")
	var labels map[string][]string
	if err == nil {
		err = json.Unmarshal(data, &labels)
	}
	if err != nil {
		log.Panicf("讀取桌面工具列文字：%v", err)
	}
	for _, lang := range []string{"zh-hant", "zh-hans", "ja", "en"} {
		if len(labels[lang]) != 7 || len(labels["hints-"+lang]) != 5 {
			log.Panicf("桌面工具列文字不完整：%s", lang)
		}
	}
	return labels
}

func mouseHint(h mouseHintID) string {
	lang := desktopLocale()
	rows := desktopLabels["hints-"+lang]
	if len(rows) == 0 {
		rows = desktopLabels["hints-zh-hant"]
	}
	return rows[int(h)]
}

func desktopLocale() string {
	if uiLang != nil {
		return string(uiLang.Lang())
	}
	return "zh-hant"
}

const (
	desktopNone desktopAction = iota
	desktopLanguage
	desktopSettings
	desktopSave
	desktopLoad
	desktopBack
	desktopConfirm
	desktopQuit
	desktopDisplay
	desktopZoomOut
	desktopZoomReset
	desktopZoomIn
	desktopReport
)

var desktopActions = []desktopAction{desktopLanguage, desktopDisplay, desktopZoomOut, desktopZoomReset, desktopZoomIn, desktopSettings, desktopReport, desktopSave, desktopLoad, desktopBack, desktopConfirm, desktopQuit}
var desktopWidths = []int{48, 72, 24, 44, 24, 48, 56, 48, 48, 48, 64, 48}
var desktopDisplayLabels = func() map[string][]string {
	b, err := translations.Files.ReadFile("desktop-display.json")
	var m map[string][]string
	if err == nil {
		err = json.Unmarshal(b, &m)
	}
	if err != nil {
		log.Panicf("desktop display labels: %v", err)
	}
	for _, lang := range []string{"zh-hant", "zh-hans", "ja", "en"} {
		if len(m[lang]) != 2 {
			log.Panicf("desktop display labels: %s", lang)
		}
	}
	return m
}()

func desktopButton(a desktopAction) image.Rectangle {
	x := 12
	for i, action := range desktopActions {
		if a == action {
			return image.Rect(x, 2, x+desktopWidths[i], desktopBarHeight-2)
		}
		x += desktopWidths[i] + 4
	}
	return image.Rectangle{}
}

func desktopHit(x, y int) desktopAction {
	for _, a := range desktopActions {
		if image.Pt(x, y).In(desktopButton(a)) {
			return a
		}
	}
	return desktopNone
}

type desktopGesture struct {
	previous, owned [3]bool
	armed           desktopAction
}

// 完整隔離三個按鍵，包括拖出工具列後放開的那一幀。
func (b *desktopGesture) pointer(x, y int, down [3]bool, enabled func(desktopAction) bool) (desktopAction, bool) {
	return b.pointerWithHeight(x, y, down, enabled, desktopBarHeight)
}

func (b *desktopGesture) pointerWithHeight(x, y int, down [3]bool, enabled func(desktopAction) bool, height int) (desktopAction, bool) {
	inside := image.Pt(x, y).In(image.Rect(0, 0, screenW, height))
	blocked := inside
	action := desktopNone
	for i := range down {
		if inside && down[i] && !b.previous[i] {
			b.owned[i] = true
			if i == 0 {
				b.armed = desktopHit(x, y)
			}
		}
		if b.owned[i] {
			blocked = true
			if !down[i] {
				if i == 0 && b.armed != desktopNone && desktopHit(x, y) == b.armed && enabled(b.armed) {
					action = b.armed
				}
				b.owned[i] = false
				if i == 0 {
					b.armed = desktopNone
				}
			}
		}
	}
	b.previous = down
	return action, blocked
}

func (g *game) desktopHostReady() bool {
	return g.world != nil && g.launcher == nil && g.world.Outcome() == state.InProgress &&
		!g.world.HasIncompleteTick() &&
		!g.battleActive() && !g.endingActive() && !g.quitting && !g.quitMenu.active &&
		!g.mapClickBusy() && g.world.PendingDiplomacy() == nil && g.world.PendingFunding() == nil
}

func (g *game) desktopEnabled(a desktopAction) bool {
	if g == nil {
		return false
	}
	if g.choiceOwnsInput() {
		return false
	}
	switch a {
	case desktopLanguage:
		return !g.quitting && (g.launcher != nil && g.naming == nil || g.desktopHostReady() && !g.hudOpen(hudSystem))
	case desktopDisplay:
		return g.hdAssetsAvailable() && !g.videoLCD && !g.quitting && (g.launcher != nil && g.naming == nil || g.desktopHostReady())
	case desktopZoomOut:
		return g.strategyZoomReady() && g.mapZoomPercent() > 50
	case desktopZoomReset:
		return g.strategyZoomReady() && g.mapZoomPercent() != 100
	case desktopZoomIn:
		return g.strategyZoomReady() && g.mapZoomPercent() < 200
	case desktopSettings, desktopSave, desktopLoad:
		return g.desktopHostReady()
	case desktopReport:
		return g.hasBattleReport() && g.desktopHostReady() && !g.hudOpen(hudSystem) && !g.factionPicker
	case desktopBack:
		return g.quitting || g.quitMenu.active || g.naming != nil ||
			g.launcher != nil && g.launcher.phase != launcherTitle ||
			g.launcher == nil && !g.battleActive() && (g.mapClickBusy() || g.hudOpen(hudSystem) ||
				g.world != nil && (g.world.PendingDiplomacy() != nil || g.world.PendingFunding() != nil))
	case desktopConfirm:
		return g.quitting || g.quitMenu.active || g.launcher != nil || g.messageActive() ||
			g.world != nil && (g.world.Outcome() != state.InProgress ||
				!g.battleActive() && (g.mapClickBusy() ||
					g.world.PendingDiplomacy() != nil || g.world.PendingFunding() != nil))
	case desktopQuit:
		return !g.quitting && !g.quitMenu.active
	}
	return false
}

func (g *game) confirmed() bool { return g.pointerAction == desktopConfirm }

type desktopGame struct {
	g                *game
	bar              desktopGesture
	canvas           *ebiten.Image
	labels           map[string][]string
	trace            *os.File
	traceLast        string
	traceAction      desktopAction
	hd               hdComposite
	finalScale       float64
	hdFrameDrawn     bool
	barVisible       bool
	finalGeometry    ebiten.GeoM
	finalBounds      image.Rectangle
	hasFinalGeometry bool
}

func newDesktopGame(g *game) *desktopGame {
	s := &desktopGame{g: g, labels: desktopLabels}
	if path := os.Getenv("WOLONG_DESKTOP_TRACE"); path != "" {
		var err error
		s.trace, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
		if err != nil {
			log.Fatalf("開啟桌面操作紀錄：%v", err)
		}
	}
	desktopGameOffset = desktopRevealHeight
	return s
}

func (s *desktopGame) Layout(int, int) (int, int) {
	return screenW, screenH + s.barHeight()
}

func (s *desktopGame) Update() error {
	g := s.g
	s.traceAction = desktopNone
	defer s.writeTrace()
	pushX, pushY := g.updateDesktopPointer()
	if g.focusManaged && !ebiten.IsFocused() {
		s.bar = desktopGesture{}
		return nil
	}
	x, y := cursorPosition()
	down := [3]bool{
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft),
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight),
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle),
	}
	if g.choiceOwnsInput() || g.listScrollOwnsInput() {
		// 原版文字選單獨占滑鼠，不能被工具列消費或回傳排隊的手勢。
		s.bar = desktopGesture{previous: down}
		return g.updateWithPointer(0, 0)
	}
	if visible := autoBarVisible(s.barVisible, x, y+s.barHeight(), down, s.bar.previous, s.bar.owned, false); visible != s.barVisible {
		s.setBarVisible(visible)
		g.idleGate.Pause()
		return nil
	}
	wasOwned := s.bar.owned[0] || s.bar.owned[1] || s.bar.owned[2]
	enabled := g.desktopEnabled
	if !s.barVisible {
		enabled = func(desktopAction) bool { return false }
	}
	action, blocked := s.bar.pointerWithHeight(x, y+s.barHeight(), down, enabled, s.barHeight())
	s.traceAction = action
	if !blocked {
		_, wheelY := ebiten.Wheel()
		if g.zoomMapWheel(wheelY, x, y, down) {
			g.idleGate.Pause()
			return nil
		}
		return g.updateWithPointer(pushX, pushY)
	}
	// 停手後沿原版 idle gate 恢復；游標留在工具列不應永久凍結軍團。
	if !wasOwned && down == [3]bool{} && action == desktopNone {
		return g.updateWithPointer(pushX, pushY)
	}
	g.idleGate.Pause()
	if action == desktopNone || !g.desktopEnabled(action) {
		return nil
	}
	switch action {
	case desktopLanguage:
		if err := g.setLanguage(langpack.Next(uiLang.Lang())); err != nil {
			log.Printf("切換語言失敗：%v", err)
		}
	case desktopSettings:
		g.hudSet(hudSystem, !g.hudOpen(hudSystem))
	case desktopReport:
		g.openBattleReport()
	case desktopDisplay:
		g.hdDisplay = !g.hdDisplay
		g.saveDesktopPreferences()
	case desktopZoomOut:
		g.zoomMapTo(g.mapZoomPercent()-25, strategyMapW/2, strategyMapH/2)
	case desktopZoomReset:
		g.zoomMapTo(100, strategyMapW/2, strategyMapH/2)
	case desktopZoomIn:
		g.zoomMapTo(g.mapZoomPercent()+25, strategyMapW/2, strategyMapH/2)
	case desktopSave:
		g.beginSaveUI(saveWrite)
	case desktopLoad:
		g.beginSaveUI(saveRead)
	case desktopQuit:
		g.quitting, g.quitYes = true, false
	case desktopBack, desktopConfirm:
		g.pointerAction = action
		defer func() { g.pointerAction = desktopNone }()
		return g.updateWithPointer(0, 0)
	}
	return nil
}

// 只讀取實際玩家狀態，供正常 GUI 驗收同步；不注入輸入、時間或亂數。
func (s *desktopGame) writeTrace() {
	if s.trace == nil {
		return
	}
	g := s.g
	x, y := cursorPosition()
	phase := -1
	if g.launcher != nil {
		phase = int(g.launcher.phase)
	}
	v := map[string]any{
		"hdLocationsAvailable": g.hdLocation.pack.Count(), "hdLocationDrawn": g.hdLocation.drawn, "hdLocationPage": g.hdLocation.page, "hdLocationCity": g.hdLocation.city, "hdLocationCacheCount": len(g.hdLocation.cache),
		"barVisible": s.barVisible, "barHeight": s.barHeight(), "barFullHeight": desktopBarHeight, "barRevealHeight": desktopRevealHeight,
		"windowSize":        [2]int{s.finalBounds.Dx(), s.finalBounds.Dy()},
		"hdBannerAvailable": g.hdBanner.pack != nil, "hdBannerDrawn": g.hdBanner.drawn,
		"hdWorldAvailable": g.hdWorld.pack.Count(), "hdWorldDrawn": g.hdWorld.drawn, "hdWorldCacheBytes": g.hdWorld.bytes,
		"hdWorldOrigin": g.hdWorld.drawOrigin, "hdWorldZoom": g.hdWorld.drawZoom,
		"zoomPercent": g.mapZoomPercent(), "zoomReady": g.strategyZoomReady(), "camSub": [2]int{g.camSubX, g.camSubY},
		"hudMinimap": g.hudOpen(hudMinimap), "hudCommand": g.hudOpen(hudCommand), "cityInfo": g.cityInfo.active,
		"hdDisplay": g.hdDisplay, "hdAvailable": g.hdPortraits.Count(), "hdNativeDraws": g.hdNativeDraws, "hdScale": s.hd.scale,
		"hdDrawPhase": g.hdDrawPhase, "hudFaction": g.hudOpen(hudFaction),
		"videoLCD": g.videoLCD, "battleNotices": g.battleNotices,
		"factionPicker": g.factionPicker, "statusBox": g.statusBoxActive(),
		"x": x, "y": y, "captured": desktopPointer.active, "phase": phase,
		"cursorMode":   int(ebiten.CursorMode()),
		"pointerReady": desktopPointer.active && desktopPointer.rebase == 0,
		"pointerRaw":   [2]int{desktopPointer.rawX, desktopPointer.rawY},
		"world":        g.world != nil, "form": g.form.active, "list": g.list != nil,
		"save": g.saveUI.active, "system": g.hudOpen(hudSystem), "quit": g.quitting,
		"messages": len(g.messages), "marchMode": g.marchMode.active, "mapPick": g.mapPickActive(),
		"advise": int(g.advise), "adviseQueue": len(g.adviseQueue),
		"adviseNode": g.adviseNode, "adviseTarget": g.target, "adviseAlly": g.ally,
		"speed": g.speed, "tacticalSpeed": g.tacticalSpeed, "lordCorps": g.lordCorps,
		"action": int(s.traceAction), "cam": [2]int{g.camX, g.camY},
		"buttons": [3]bool{ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft), ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight), ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)},
		"canBack": g.desktopEnabled(desktopBack), "hostReady": g.desktopHostReady(),
	}
	buttons := map[string][4]int{}
	for _, row := range []struct {
		key    string
		action desktopAction
	}{{"zoomOut", desktopZoomOut}, {"zoomReset", desktopZoomReset}, {"zoomIn", desktopZoomIn}, {"display", desktopDisplay}, {"language", desktopLanguage}, {"settings", desktopSettings}, {"report", desktopReport}, {"save", desktopSave}, {"load", desktopLoad}, {"back", desktopBack}, {"confirm", desktopConfirm}, {"quit", desktopQuit}} {
		r := desktopButton(row.action)
		buttons[row.key] = [4]int{r.Min.X, r.Min.Y, r.Max.X, r.Max.Y}
	}
	v["toolbar"] = buttons
	if g.hasBattleReport() {
		r := g.battleReport.latest
		v["battleReport"] = map[string]any{
			"open": g.battleReportOpen(), "toastVisible": g.battleReport.drawn && g.battleToastVisibleAt(time.Now()),
			"summary": g.battleReportSummary(), "before": r.before, "after": r.after, "playerSide": r.playerSide,
			"names": [2]string{r.name(0), r.name(1)}, "defenderWins": r.defenderWins,
			"toastRect": g.battleToastRect(), "closeRect": battleReportCloseRect,
		}
	}
	v["cityInfoCity"] = g.cityInfo.city
	if g.list != nil {
		v["listPhase"], v["listCursor"], v["listPick"] = int(g.list.Phase()), g.list.Cursor, g.listPick != nil
		v["listRows"], v["listTop"] = g.list.Rows, g.list.Top
	}
	v["listDragging"] = g.listScrollOwnsInput()
	if g.listScrollOwnsInput() {
		v["listPreviewTop"] = g.listDrag.top
	}
	if g.launcher != nil {
		v["factionTop"], v["factionCount"], v["factionSelected"] = g.launcher.factionListTop(), len(g.launcher.players), g.launcher.factionSelected
	}
	v["language"] = desktopLocale()
	v["choiceActive"] = g.choiceBand.owner != nil
	if g.choiceBand.owner != nil {
		v["choiceRow"], v["choiceY"], v["choiceCount"] = *g.choiceBand.owner, desktopPointer.y, g.choiceBand.count
	}
	v["popupActive"], v["popupRow"] = g.popupMenuActive(), g.cmdMenu.row
	if g.world != nil {
		v["capital"] = g.world.Factions[g.world.Player].Capital
		v["player"] = g.world.Player
		wars := []int{}
		for i, f := range g.world.Friendship[g.world.Player] {
			if i != g.world.Player && f.AtWar() {
				wars = append(wars, i)
			}
		}
		v["wars"] = wars
		owners := make([]int, len(g.world.Cities))
		for i, c := range g.world.Cities {
			owners[i] = c.Owner
		}
		v["cityOwners"] = owners
		if p := g.world.PendingBattle(); p != nil {
			b := p.Battle
			side := g.battleSide()
			v["battle"] = map[string]any{
				"attacker": p.Attacker, "defender": p.Defender, "node": p.Node,
				"mode": int(p.Mode), "side": side, "frame": b.Frame, "done": b.Done,
				"opening": b.OpeningActive(), "talk": g.battleTalkSession.queue.active(),
				"selected": b.Sides[side].Selected, "formation": b.Sides[side].Formation,
				"line": b.Sides[side].Line, "standing": int(b.Sides[side].Standing),
				"leaderOrder":    [2]int{int(b.Sides[side].Soldiers[0].Cmd), int(b.Sides[side].Soldiers[0].Next)},
				"leaderPosition": [3]int{b.Sides[side].Soldiers[0].X, b.Sides[side].Soldiers[0].Y, b.Sides[side].Soldiers[0].Z},
				"leaderGoal":     [3]int{b.Sides[side].Soldiers[0].GoalX, b.Sides[side].Soldiers[0].GoalY, b.Sides[side].Soldiers[0].GoalZ},
				"leaderStamina":  b.Sides[side].Soldiers[0].Stamina,
				"fastForward":    g.battleFastForward,
				"alive":          [2]int{b.Sides[0].Alive(), b.Sides[1].Alive()},
				"survivors":      [2]int{b.Sides[0].Survivors(), b.Sides[1].Survivors()},
				"generalHP":      [2]int{b.Sides[0].Soldiers[0].HP, b.Sides[1].Soldiers[0].HP},
			}
			if b.Done {
				v["battleOutcome"] = b.Result()
			}
		}
		c := g.world.Clock
		v["date"] = [4]int{c.Year, c.Month, c.Day, c.Hour}
		corps := [][7]int{}
		routes := [][8]int{}
		for i, c := range g.world.Corps {
			if c.Alive {
				corps = append(corps, [7]int{i, c.X, c.Y, c.Heading, c.Stage, c.TargetNode, c.Timer})
				routes = append(routes, [8]int{i, c.Node, c.LinkAddr, c.PathPtr, c.Direction, c.Ordered, c.TargetX, c.TargetY})
			}
		}
		v["corps"] = corps
		v["corpsRoutes"] = routes
	}
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("桌面操作紀錄：%v", err)
		return
	}
	if string(data) != s.traceLast {
		if _, err := s.trace.Write(append(data, '\n')); err != nil {
			log.Printf("桌面操作紀錄：%v", err)
		}
		s.traceLast = string(data)
	}
}

func (s *desktopGame) Draw(dst *ebiten.Image) {
	if s.canvas == nil {
		s.canvas = ebiten.NewImage(screenW, screenH)
	}
	s.canvas.Clear()
	s.hdFrameDrawn = s.g.hdDisplay && !s.g.videoLCD && s.g.hdAssetsAvailable()
	s.g.hdWorld.drawn = 0
	s.g.hdBanner.drawn = false
	s.g.hdLocation.drawn = false
	s.g.hdLocation.page, s.g.hdLocation.city = -1, -1
	s.g.hdNativeDraws = 0
	s.g.hdDrawPhase = -1
	if s.g.launcher != nil {
		s.g.hdDrawPhase = int(s.g.launcher.phase)
	}
	if s.hdFrameDrawn {
		s.hd.begin(s.finalScale)
		s.g.hdFrame = &s.hd
	}
	s.g.Draw(s.canvas)
	s.g.hdFrame = nil
	if s.g.launcher != nil && s.g.quitting {
		s.g.drawYesNo(s.canvas, quitDialogX, quitDialogY, "確定離開？", s.g.quitYes)
	}
	if s.hdFrameDrawn {
		s.hd.flush(s.canvas)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1/s.hd.scale, 1/s.hd.scale)
		s.canvas.DrawImage(s.hd.surface, op)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, float64(s.barHeight()))
	dst.DrawImage(s.canvas, op)
	if !s.barVisible {
		vector.DrawFilledRect(dst, 0, 0, screenW, desktopRevealHeight, color.RGBA{20, 27, 37, 255}, false)
		vector.DrawFilledRect(dst, screenW/2-18, 2, 36, 2, color.RGBA{190, 170, 126, 255}, false)
		if desktopPointer.active {
			x, y := cursorPosition()
			if y < 0 {
				drawArrowCursor(dst, image.Pt(x, y+s.barHeight()), color.RGBA{255, 255, 255, 255}, color.RGBA{240, 210, 149, 255})
			}
		}
		return
	}
	vector.DrawFilledRect(dst, 0, 0, screenW, desktopBarHeight, color.RGBA{20, 27, 37, 255}, false)
	labels := s.labels[desktopLocale()]
	if len(labels) == 0 {
		labels = s.labels["zh-hant"]
	}
	x, y := cursorPosition()
	for _, a := range desktopActions {
		r := desktopButton(a)
		ink := color.RGBA{240, 210, 149, 255}
		fill := color.RGBA{39, 48, 62, 255}
		if !s.g.desktopEnabled(a) {
			ink = color.RGBA{111, 117, 126, 255}
		} else if image.Pt(x, y+desktopBarHeight).In(r) {
			fill = color.RGBA{68, 78, 94, 255}
		}
		vector.DrawFilledRect(dst, float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy()), ink, false)
		vector.DrawFilledRect(dst, float32(r.Min.X+1), float32(r.Min.Y+1), float32(r.Dx()-2), float32(r.Dy()-2), fill, false)
		if s.g.td != nil {
			label := ""
			if a == desktopDisplay {
				index := 0
				if s.g.hdDisplay && !s.g.videoLCD && s.g.hdAssetsAvailable() {
					index = 1
				}
				label = desktopDisplayLabels[desktopLocale()][index]
			} else if a >= desktopZoomOut && a <= desktopZoomIn {
				label = s.g.desktopZoomLabel(a)
			} else if a == desktopReport {
				label = battleReportText("button")
			} else {
				label = labels[int(a)-1]
			}
			s.g.td.Draw(dst, label, r.Min.X+(r.Dx()-s.g.td.Width(label))/2, (desktopBarHeight-16)/2, ink)
		}
	}
	if desktopPointer.active && y < 0 && s.g.choiceBand.owner == nil && !s.g.listScrollOwnsInput() {
		drawArrowCursor(dst, image.Pt(x, y+desktopBarHeight), color.RGBA{255, 255, 255, 255}, color.RGBA{240, 210, 149, 255})
	}
}
