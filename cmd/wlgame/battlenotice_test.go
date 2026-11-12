package main

import (
	"fmt"
	"image/color"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

func TestBattleNoticeVisibleClockAndReopen(t *testing.T) {
	g := &game{world: &state.World{}, battleNotices: true}
	g.recordBattleReport(state.CorpsEvent{Corps: 0, Enemy: 1, Battle: &combat.Result{}})
	now := time.Unix(100, 0)
	if !g.battleReport.pendingAt(now.Add(time.Hour)) {
		t.Fatal("首次繪製前已逾時")
	}
	g.battleReport.showAt(now)
	g.battleReport.showAt(now.Add(time.Second))
	if !g.battleReport.pendingAt(now.Add(3*time.Second-time.Nanosecond)) || g.battleReport.pendingAt(now.Add(3*time.Second)) {
		t.Fatal("三秒邊界或重繪計時錯誤")
	}
	if !g.desktopEnabled(desktopReport) || g.mapClickBusy() {
		t.Fatal("短提示到期後不能重開，或阻擋地圖")
	}
	g.openBattleReport()
	if !g.mapClickBusy() || g.desktopHostReady() || !g.desktopEnabled(desktopBack) || !g.desktopEnabled(desktopConfirm) {
		t.Fatal("明細未取得模態輸入")
	}
	g.pointerAction = desktopBack
	g.updateBattleReport()
	if g.battleReportOpen() || g.battleReport.pendingAt(now) || !g.hasBattleReport() {
		t.Fatal("關閉後重播提示或清掉明細")
	}
	g.openBattleReport()
	g.pointerAction = desktopConfirm
	g.updateBattleReport()
	if g.battleReportOpen() {
		t.Fatal("工具列確認未關閉明細")
	}
	g.recordBattleReport(state.CorpsEvent{Corps: 0, Enemy: 1, Battle: &combat.Result{}})
	if !g.battleReport.pendingAt(now) || !g.battleReport.shownAt.IsZero() {
		t.Fatal("下一場未重設提示")
	}
	g.world = &state.World{}
	if g.hasBattleReport() || g.battleReportOpen() || g.desktopEnabled(desktopReport) {
		t.Fatal("新世界沿用舊戰報")
	}
}

func TestBattleNoticesDefaultOffAndOnlyNewResultsNotify(t *testing.T) {
	g := &game{world: &state.World{}}
	e := state.CorpsEvent{Corps: 0, Enemy: 1, Battle: &combat.Result{}}
	for i := 0; i < 3; i++ {
		e.BattleAfter[0] = 600 - i
		g.recordBattleReport(e)
		if g.battleNotices || g.battleToastVisibleAt(time.Now()) || g.battleReportOpen() || g.mapClickBusy() {
			t.Fatal("預設或連續戰報打擾地圖")
		}
	}
	if !g.hasBattleReport() || g.battleReport.latest.after[0] != 598 || !g.desktopEnabled(desktopReport) {
		t.Fatal("關閉自動提示後無法手動查看最近一場")
	}
	g.openBattleReport()
	if !g.battleReportOpen() {
		t.Fatal("手動明細未開啟")
	}
	g.pointerAction = desktopBack
	g.updateBattleReport()
	g.dispatchSystemRow(sysRowBattleNotices, true)
	if !g.battleNotices || g.battleToastVisibleAt(time.Now()) {
		t.Fatal("啟用時補播舊戰報")
	}
	g.recordBattleReport(e)
	if !g.battleToastVisibleAt(time.Now()) {
		t.Fatal("啟用後的新戰鬥未提示")
	}
	g.dispatchSystemRow(sysRowBattleNotices, false)
	if g.battleNotices || g.battleToastVisibleAt(time.Now()) {
		t.Fatal("停用後仍顯示提示")
	}
	g.dispatchSystemRow(sysRowBattleNotices, true)
	if g.battleToastVisibleAt(time.Now()) {
		t.Fatal("重開選項後補播已取消的提示")
	}
	if sysRowBattleNotices != 9 || sysWinY+sysWinH > screenH {
		t.Fatal("選項改動既有列或超出畫布")
	}
	for lang, labels := range battleReportLabels {
		if textdraw.StringWidth(labels["setting"]) > sysValueX-sysLabelX-8 {
			t.Fatalf("%s 戰報選項標籤溢出", lang)
		}
	}
}

// 獨立RunGame讀回實際GPU像素，避免把只有模態狀態的GUI收據當成可讀文字的證據。
func TestBattleReportGPUKeepsTextAboveNativeMap(t *testing.T) {
	if os.Getenv("WOLONG_BATTLE_NOTICE_GPU") == "" {
		t.Skip("separate RunGame process required")
	}
	lib, err := library.Load(os.Getenv("WOLONG_BATTLE_NOTICE_ORIG"))
	if err != nil {
		t.Fatal(err)
	}
	p := &battleReportGPUProbe{g: &game{world: &state.World{}, chrome: chrome.Load(lib, 0), td: textdraw.New(nil, nil)}}
	if err := ebiten.RunGame(p); err != nil && err != ebiten.Termination {
		t.Fatal(err)
	}
	if p.err != nil {
		t.Fatal(p.err)
	}
}

type battleReportGPUProbe struct {
	g    *game
	done bool
	err  error
}

func (p *battleReportGPUProbe) Layout(int, int) (int, int) { return screenW, screenH }
func (p *battleReportGPUProbe) Update() error {
	if p.done {
		return ebiten.Termination
	}
	return nil
}
func (p *battleReportGPUProbe) Draw(_ *ebiten.Image) {
	if p.done {
		return
	}
	p.done = true
	g := p.g
	g.recordBattleReport(state.CorpsEvent{Corps: 0, Enemy: 1, Battle: &combat.Result{}, BattleBefore: [2]int{600, 600}, BattleAfter: [2]int{550, 600}})
	g.openBattleReport()
	classic, legacy := ebiten.NewImage(screenW, screenH), ebiten.NewImage(screenW, screenH)
	defer classic.Dispose()
	defer legacy.Dispose()
	classic.Fill(color.Black)
	legacy.Fill(color.Black)
	g.drawBattleReport(classic)
	frame := &hdComposite{}
	frame.begin(2)
	defer frame.surface.Dispose()
	g.hdFrame = frame
	g.drawBattleReport(legacy)
	frame.flush(legacy)
	ink := 0
	for y := battleReportRect.Min.Y; y < battleReportRect.Max.Y; y++ {
		for x := battleReportRect.Min.X; x < battleReportRect.Max.X; x++ {
			want := classic.At(x, y)
			if color.RGBAModel.Convert(want).(color.RGBA) == chrome.Paper {
				ink++
			}
			for dy := 0; dy < 2; dy++ {
				for dx := 0; dx < 2; dx++ {
					if frame.surface.At(x*2+dx, y*2+dy) != want {
						p.err = fmt.Errorf("HD戰報像素不同: %d,%d", x, y)
						return
					}
				}
			}
		}
	}
	if ink < 100 {
		p.err = fmt.Errorf("戰報缺少文字像素: %d", ink)
	}
}

func TestQuickBattleReportSnapshotAndInputPriority(t *testing.T) {
	w := &state.World{Player: 0}
	w.Generals[1].Name, w.Generals[2].Name = "ATTACKER", "DEFENDER"
	w.Corps[1].Faction, w.Corps[2].Faction = 13, 0
	g := &game{world: w, lastEvent: "existing event", battleNotices: true}
	e := state.CorpsEvent{Corps: 1, Enemy: 2, GovernorReturned: -1, Captured: -1,
		Battle: &combat.Result{DefenderWins: true}, BattleBefore: [2]int{600, 500}, BattleAfter: [2]int{300, 450}}
	g.reportCorps(state.Event{Corps: []state.CorpsEvent{e}})
	if !g.hasBattleReport() || g.battleReport.latest.playerSide != 1 || !g.battleReport.latest.defenderWins || g.lastEvent != "existing event" {
		t.Fatal("快速戰報未記錄守方，或漏進原事件列")
	}
	e.BattleAfter[1], e.Battle.DefenderWins = 1, false
	w.Generals[2].Name = "CHANGED"
	r := g.battleReport.latest
	if r.after != [2]int{300, 450} || r.name(1) != "DEFENDER" || !r.defenderWins {
		t.Fatal("戰後變動改寫已記錄快照")
	}
	g.lastEvent = ""
	if width := g.reportWidth(g.battleReportSummary()); width <= 336 && width > g.battleToastRect().Dx()-80 {
		t.Fatal("一般摘要被查看按鈕擠掉")
	}
	if !g.battleToastVisibleAt(time.Now()) {
		t.Fatal("可操作地圖未顯示提示")
	}
	g.saveUI.active = true
	if g.battleToastVisibleAt(time.Now()) || g.desktopEnabled(desktopReport) || !g.battleReport.shownAt.IsZero() {
		t.Fatal("提示蓋過存檔視窗或提早起算")
	}
	g.saveUI.active = false
	g.factionPicker = true
	if g.battleToastVisibleAt(time.Now()) || g.desktopEnabled(desktopReport) {
		t.Fatal("提示搶走勢力選擇輸入")
	}
}

func TestBattleReportToolbarAndLocalizedFieldsFit(t *testing.T) {
	for _, lang := range []string{"zh-hant", "zh-hans", "ja", "en"} {
		labels := desktopLabels[lang]
		end := 0
		for _, action := range desktopActions {
			r := desktopButton(action)
			if r.Min.X < end || r.Max.X > screenW || r.Dy() != desktopBarHeight-4 {
				t.Fatalf("%s 工具列超界或重疊: %v", lang, r)
			}
			end = r.Max.X
			label := ""
			switch {
			case action == desktopReport:
				label = battleReportLabels[lang]["button"]
			case action <= desktopQuit:
				label = labels[int(action)-1]
			default:
				continue
			}
			if textdraw.StringWidth(label) > r.Dx()-8 {
				t.Fatalf("%s 按鈕溢出: %q", lang, label)
			}
		}
		texts := battleReportLabels[lang]
		for _, key := range []string{"same", "changed"} {
			var line string
			if key == "changed" {
				line = fmt.Sprintf(texts[key], strings.Repeat("將", 10), 60000, 59990)
			} else {
				line = fmt.Sprintf(texts[key], strings.Repeat("將", 10))
			}
			got := fitBattleReportText(line, 336, textdraw.StringWidth)
			if got == "" || textdraw.StringWidth(got) > 336 {
				t.Fatalf("摘要超界: %s %s", lang, got)
			}
		}
		for _, key := range []string{"commander", "before", "after"} {
			got := fitBattleReportText(texts[key], 104, textdraw.StringWidth)
			if textdraw.StringWidth(got) > 104 {
				t.Fatalf("表頭超界: %s/%s", lang, key)
			}
		}
	}
	if fitBattleReportText("中文字", 4, textdraw.StringWidth) != "" {
		t.Fatal("極窄欄位溢出")
	}
}
