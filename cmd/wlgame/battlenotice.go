package main

// 使用者選定的三秒提示與按需戰報，屬於remake介面（spec89續片）。
import (
	"encoding/json"
	"fmt"
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
	"github.com/wicanr2/wolong_cht/translations"
)

const battleToastDuration = 3 * time.Second

var battleReportLabels = func() map[string]map[string]string {
	b, err := translations.Files.ReadFile("desktop-battle-report.json")
	if err != nil {
		panic(err)
	}
	var labels map[string]map[string]string
	if err := json.Unmarshal(b, &labels); err != nil {
		panic(err)
	}
	for _, lang := range []string{"zh-hant", "zh-hans", "ja", "en"} {
		for _, key := range []string{"button", "view", "close", "title", "attacker", "defender", "commander", "before", "after", "result", "attackerWins", "defenderWins", "garrison", "yourSide", "same", "changed", "general"} {
			if labels[lang][key] == "" {
				panic("incomplete battle report labels: " + lang + "/" + key)
			}
		}
	}
	return labels
}()

func battleReportText(key string) string { return battleReportLabels[desktopLocale()][key] }

type battleReportSnapshot struct {
	rawNames      [2]string
	before, after [2]int
	defenderWins  bool
	playerSide    int
}

type battleReportState struct {
	world                  *state.World
	latest                 *battleReportSnapshot
	open, dismissed, drawn bool
	shownAt                time.Time
}

func (r *battleReportState) showAt(now time.Time) {
	if r.shownAt.IsZero() {
		r.shownAt = now
	}
}

func (r *battleReportState) pendingAt(now time.Time) bool {
	return r.latest != nil && !r.open && !r.dismissed &&
		(r.shownAt.IsZero() || now.Before(r.shownAt.Add(battleToastDuration)))
}

func (g *game) hasBattleReport() bool {
	return g != nil && g.world != nil && g.battleReport.world == g.world && g.battleReport.latest != nil
}

func (g *game) battleReportOpen() bool { return g.hasBattleReport() && g.battleReport.open }

// 一次複製結果，不從稍後已移動／易主的軍團重算戰果。原名保留到繪製時才轉語系。
func (g *game) recordBattleReport(e state.CorpsEvent) {
	if g.world == nil || e.Battle == nil || e.Corps < 0 || e.Corps >= len(g.world.Generals) {
		return
	}
	r := &battleReportSnapshot{before: e.BattleBefore, after: e.BattleAfter,
		defenderWins: e.Battle.DefenderWins, playerSide: -1}
	r.rawNames[0] = g.world.Generals[e.Corps].Name
	if g.world.Corps[e.Corps].Faction == g.world.Player {
		r.playerSide = 0
	}
	if e.Enemy >= 0 && e.Enemy < len(g.world.Generals) {
		r.rawNames[1] = g.world.Generals[e.Enemy].Name
		if g.world.Corps[e.Enemy].Faction == g.world.Player {
			r.playerSide = 1
		}
	}
	g.battleReport = battleReportState{world: g.world, latest: r}
}

func (r *battleReportSnapshot) name(side int) string {
	if r.rawNames[side] == "" {
		return battleReportText("garrison")
	}
	return big5(r.rawNames[side])
}

func (g *game) battleReportSummary() string {
	if !g.hasBattleReport() {
		return ""
	}
	r := g.battleReport.latest
	if r.playerSide < 0 {
		return battleReportText("general")
	}
	i := r.playerSide
	if r.before[i] == r.after[i] {
		return fmt.Sprintf(battleReportText("same"), r.name(i))
	}
	return fmt.Sprintf(battleReportText("changed"), r.name(i), r.before[i]*10, r.after[i]*10)
}

func (g *game) openBattleReport() {
	if g.hasBattleReport() {
		g.battleReport.open, g.battleReport.dismissed = true, true
		g.idleGate.Pause()
	}
}

func (g *game) battleToastVisibleAt(now time.Time) bool {
	return g.hasBattleReport() && g.battleReport.pendingAt(now) && g.desktopHostReady() &&
		!g.hudOpen(hudSystem) && !g.factionPicker && !(g.lastEvent != "" && g.frame-g.lastEventAt < eventLineFrames)
}

func (g *game) reportWidth(s string) int {
	if g.td != nil {
		return g.td.Width(s)
	}
	return textdraw.StringWidth(s)
}

func (g *game) battleToastRect() image.Rectangle {
	w := max(280, min(416, g.reportWidth(g.battleReportSummary())+80))
	w = (w + 7) / 8 * 8
	return image.Rect(8, screenH-40, 8+w, screenH-8)
}

var battleReportRect = image.Rect(72, 96, 568, 304)
var battleReportCloseRect = image.Rect(440, 264, 552, 288)

// 三秒提示沒有輸入閘，只有玩家點到它才主動開明細並吃掉該次地圖點擊。
func (g *game) updateBattleToast() bool {
	if !g.battleToastVisibleAt(time.Now()) || !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return false
	}
	x, y := cursorPosition()
	if !image.Pt(x, y).In(g.battleToastRect()) {
		return false
	}
	g.openBattleReport()
	return true
}

func (g *game) updateBattleReport() {
	g.idleGate.Pause()
	x, y := cursorPosition()
	if g.cancelled() || g.confirmed() ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && image.Pt(x, y).In(battleReportCloseRect) {
		g.battleReport.open = false
	}
}

func fitBattleReportText(s string, room int, width func(string) int) string {
	if width(s) <= room {
		return s
	}
	rs := []rune(s)
	for len(rs) > 0 {
		rs = rs[:len(rs)-1]
		if candidate := string(rs) + "..."; width(candidate) <= room {
			return candidate
		}
	}
	return ""
}

func (g *game) drawReportText(dst *ebiten.Image, s string, x, y, w int) {
	g.td.Draw(dst, fitBattleReportText(s, w, g.td.Width), x, y, chrome.Paper)
}

func (g *game) drawBattleToast(dst *ebiten.Image) {
	g.battleReport.drawn = false
	if g.chrome == nil || g.td == nil || !g.battleToastVisibleAt(time.Now()) {
		return
	}
	if g.hdFrame != nil {
		g.hdFrame.flush(dst)
	}
	g.battleReport.showAt(time.Now())
	g.battleReport.drawn = true
	r := g.battleToastRect()
	g.chrome.Window(dst, r.Min.X, r.Min.Y, r.Dx(), r.Dy(), chrome.Blank)
	g.drawReportText(dst, g.battleReportSummary(), r.Min.X+8, r.Min.Y+8, r.Dx()-80)
	g.drawReportText(dst, battleReportText("view"), r.Max.X-56, r.Min.Y+8, 48)
}

func (g *game) drawBattleReport(dst *ebiten.Image) {
	if g.chrome == nil || g.td == nil || !g.battleReportOpen() {
		return
	}
	if g.hdFrame != nil {
		g.hdFrame.flush(dst)
	}
	r := battleReportRect
	g.chrome.Window(dst, r.Min.X, r.Min.Y, r.Dx(), r.Dy(), chrome.Menu)
	g.drawReportText(dst, battleReportText("title"), 88, 112, 464)
	b := g.battleReport.latest
	if b.playerSide >= 0 {
		role := []string{"attacker", "defender"}[b.playerSide]
		g.drawReportText(dst, fmt.Sprintf(battleReportText("yourSide"), battleReportText(role)), 88, 136, 464)
	}
	g.drawReportText(dst, battleReportText("commander"), 168, 160, 144)
	g.drawReportText(dst, battleReportText("before"), 328, 160, 104)
	g.drawReportText(dst, battleReportText("after"), 448, 160, 104)
	for i, role := range []string{"attacker", "defender"} {
		y := 184 + i*24
		g.drawReportText(dst, battleReportText(role), 88, y, 72)
		g.drawReportText(dst, b.name(i), 168, y, 144)
		g.drawReportText(dst, fmt.Sprint(b.before[i]*10), 328, y, 104)
		g.drawReportText(dst, fmt.Sprint(b.after[i]*10), 448, y, 104)
	}
	winner := "attackerWins"
	if b.defenderWins {
		winner = "defenderWins"
	}
	g.drawReportText(dst, fmt.Sprintf(battleReportText("result"), battleReportText(winner)), 88, 240, 464)
	c := battleReportCloseRect
	g.chrome.Window(dst, c.Min.X, c.Min.Y, c.Dx(), c.Dy(), chrome.Blank)
	label := battleReportText("close")
	g.drawReportText(dst, label, c.Min.X+(c.Dx()-g.reportWidth(label))/2, c.Min.Y+4, c.Dx()-16)
}
