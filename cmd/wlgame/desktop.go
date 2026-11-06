package main

// 桌面外殼與原版畫面分開。規格與使用者裁定見 docs/spec/266。
import (
	"encoding/json"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/langpack"
	"github.com/wicanr2/wolong_cht/translations"
)

const desktopBarHeight = 40

// 對拍與錄影仍直接使用 game，座標偏移保持零。
var desktopGameOffset int

type desktopAction int

type mouseHintID int

const (
	mouseSelect mouseHintID = iota
	mouseClose
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
		if len(labels[lang]) != 7 || len(labels["hints-"+lang]) != 6 {
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
)

func desktopButton(a desktopAction) image.Rectangle {
	if a < desktopLanguage || a > desktopQuit {
		return image.Rectangle{}
	}
	x := 6 + (int(a)-1)*90
	return image.Rect(x, 4, x+88, desktopBarHeight-4)
}

func desktopHit(x, y int) desktopAction {
	for a := desktopLanguage; a <= desktopQuit; a++ {
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
	inside := image.Pt(x, y).In(image.Rect(0, 0, screenW, desktopBarHeight))
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
	case desktopSettings, desktopSave, desktopLoad:
		return g.desktopHostReady()
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
	g           *game
	bar         desktopGesture
	canvas      *ebiten.Image
	labels      map[string][]string
	trace       *os.File
	traceLast   string
	traceAction desktopAction
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
	desktopGameOffset = desktopBarHeight
	return s
}

func (s *desktopGame) Layout(int, int) (int, int) {
	return screenW, screenH + desktopBarHeight
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
	wasOwned := s.bar.owned[0] || s.bar.owned[1] || s.bar.owned[2]
	action, blocked := s.bar.pointer(x, y+desktopBarHeight, down, g.desktopEnabled)
	s.traceAction = action
	if !blocked {
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
				"leaderOrder": [2]int{int(b.Sides[side].Soldiers[0].Cmd), int(b.Sides[side].Soldiers[0].Next)},
				"fastForward": g.battleFastForward,
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
	s.g.Draw(s.canvas)
	if s.g.launcher != nil && s.g.quitting {
		s.g.drawYesNo(s.canvas, quitDialogX, quitDialogY, "確定離開？", s.g.quitYes)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, desktopBarHeight)
	dst.DrawImage(s.canvas, op)
	vector.DrawFilledRect(dst, 0, 0, screenW, desktopBarHeight, color.RGBA{20, 27, 37, 255}, false)
	labels := s.labels[desktopLocale()]
	if len(labels) == 0 {
		labels = s.labels["zh-hant"]
	}
	x, y := cursorPosition()
	for a := desktopLanguage; a <= desktopQuit; a++ {
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
			label := labels[int(a)-1]
			s.g.td.Draw(dst, label, r.Min.X+(r.Dx()-s.g.td.Width(label))/2, 12, ink)
		}
	}
	if desktopPointer.active && y < 0 && s.g.choiceBand.owner == nil && !s.g.listScrollOwnsInput() {
		drawArrowCursor(dst, image.Pt(x, y+desktopBarHeight), color.RGBA{255, 255, 255, 255}, color.RGBA{240, 210, 149, 255})
	}
}
