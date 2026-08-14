package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 自然策略常駐骨架，數值全部出自機器碼（docs/spec/12 §1、docs/re/47）：
// 640×400、32 px 橫幅、命令 (0,32,432,32)、縮小地圖 (432,32,208,160)、
// 自勢力情報 (432,192,208,208)、地圖 (0,64,432,336)。
func TestDOSVNaturalStrategySkeleton(t *testing.T) {
	checks := []struct {
		name string
		got  int
		want int
	}{
		{"screen width", screenW, 640},
		{"screen height", screenH, 400},
		{"banner height", bannerH, 32},
		{"command y", strategyCommandY, 32},
		{"command height", strategyCommandH, 32},
		{"map y", strategyMapY, 64},
		{"map width", strategyMapW, 432},
		{"map height", strategyMapH, 336},
		{"sidebar x", strategySidebarX, 432},
		{"sidebar width", strategySidebarW, 208},
		{"minimap height", strategyMinimapH, 160},
		{"minimap legend y", strategyMinimapLegendY, 168},
		{"minimap swatch y", strategyMinimapSwatchY, 172},
		{"faction y", strategyFactionY, 192},
		{"faction height", strategyFactionH, 208},
		{"command text x", strategyCommandX + strategyCommandLead, 32},
		{"command hit x", strategyCommandHitX, 24},
	}
	for _, check := range checks {
		if check.got != check.want {
			t.Errorf("%s = %d，want %d", check.name, check.got, check.want)
		}
	}
	// 右欄三段相加剛好鋪滿畫面高度，是「矩形讀對了」的算術檢查。
	if got := bannerH + strategyMinimapH + strategyFactionH; got != screenH {
		t.Fatalf("右欄三段相加 = %d，want %d", got, screenH)
	}
	if got := strategyMinimapLegendY + textdraw.GlyphH; got > strategyFactionY-chrome.Tile {
		t.Fatalf("minimap 色標越過視窗下緣：end=%d，內框下緣=%d", got, strategyFactionY-chrome.Tile)
	}
	if strategySidebarInnerX != 440 || strategyMinimapY != 40 || strategyMinimapW != 192 {
		t.Fatalf("DOS/V minimap 內框 = x%d y%d w%d，want x440 y40 w192",
			strategySidebarInnerX, strategyMinimapY, strategyMinimapW)
	}
	if strategyFactionInnerY != 200 || strategyFactionInnerW != 192 {
		t.Fatalf("DOS/V 情報內框 = y%d w%d，want y200 w192", strategyFactionInnerY, strategyFactionInnerW)
	}
	// 命令列的八個字與尾端空白要落在 432 的框寬內（docs/re/47 §4.1）。
	if end := strategyCommandX + strategyCommandLead + 8*strategyCommandTextW +
		7*(strategyCommandCellW-strategyCommandTextW) + 24; end > strategyMapW {
		t.Fatalf("命令列字串尾端 = %d，超出框寬 %d", end, strategyMapW)
	}
	// 熱區照 sub_1E3D7(al=0x0C, X=24, Y=40, 384×16)。
	if r := strategyCommandCellRect(0); r.Min.X != 24 || r.Min.Y != 40 || r.Dy() != 16 {
		t.Fatalf("命令列第一格 = %v，want x24 y40 h16", r)
	}
	if r := strategyCommandCellRect(7); r.Max.X != 24+384 {
		t.Fatalf("命令列最後一格右緣 = %d，want %d", r.Max.X, 24+384)
	}
	if strategyInfoYOffset != 24 || strategyInfoRowStep != 17 || strategyInfoDividerXOffset != 120 ||
		strategyInfoDividerH != 48 || strategyTrustYOffset != 88 || strategyResourceDividerY != 120 ||
		strategyResourceBoxY != 128 || strategyResourceBoxH != 88 {
		t.Fatal("DOS/V 情報列／分隔線／資源區幾何契約被改動")
	}
}

func TestDOSVNaturalStrategyTextContainment(t *testing.T) {
	// 這些是 active 8/16 px bitmap 字寬的 containment 測試，不是新增的原版數值證據。
	//
	// 命令列那一項**是**原版數值：`sub_106F5(dx=8, bx=0x28)` → X=8、Y=40，
	// `cs:6181h` 的詞是兩個全形字（32 px），節距 48（詞 32 ＋ 全形空格 16）。
	if strategyCommandTextW != 32 || strategyCommandCellW != 48 ||
		strategyCommandX != 8 || strategyCommandLead != 24 ||
		strategyInfoValueW != 72 {
		t.Fatalf("命令列版面 = x%d lead%d cell%d text%d／info %d，want 8/24/48/32/72",
			strategyCommandX, strategyCommandLead, strategyCommandCellW,
			strategyCommandTextW, strategyInfoValueW)
	}
	// 第八個詞的右緣，以及尾端「 　」之後的位置，都要在框寬 432 內。
	if right := strategyCommandX + strategyCommandLead +
		7*strategyCommandCellW + strategyCommandTextW; right != 400 {
		t.Fatalf("第八個詞右緣 = %d，want 400", right)
	}
	if strategyNumberSlots != 5 || strategyNumberW != 40 || strategyNumberXOffset != 160 {
		t.Fatalf("數字欄 = slots %d width %d x-offset %d，want 5/40/160",
			strategyNumberSlots, strategyNumberW, strategyNumberXOffset)
	}

	for _, value := range []int{0, 1, 9999, 99999, -1, -9999} {
		text := strategyHUDNumber(value)
		if got := textdraw.StringWidth(text); got != strategyNumberW {
			t.Errorf("%d 格式化為 %q、寬度 %d，want %d", value, text, got, strategyNumberW)
		}
	}
	for _, test := range []struct {
		value int
		want  string
	}{{100000, "99999"}, {655000, "99999"}, {-10000, "-9999"}, {-655000, "-9999"}} {
		if got := strategyHUDNumber(test.value); got != test.want {
			t.Errorf("超出五槽的 %d 顯示為 %q，want %q", test.value, got, test.want)
		}
	}

	for _, test := range []struct {
		name string
		text string
		max  int
		want string
	}{
		{name: "halfwidth", text: "12345", max: 40, want: "12345"},
		{name: "fullwidth", text: "曹操許昌", max: 48, want: "曹操許"},
		{name: "mixed", text: "ABC曹操", max: 48, want: "ABC曹"},
	} {
		bounded := strategyHUDSingleLine(test.text, test.max)
		if bounded != test.want {
			t.Errorf("%s = %q，want %q", test.name, bounded, test.want)
		}
		if got := textdraw.StringWidth(bounded); got > test.max {
			t.Errorf("%s 寬度 %d 越過安全寬度 %d", test.name, got, test.max)
		}
	}
	for _, label := range naturalCommandLabels {
		if got := textdraw.StringWidth(strategyHUDSingleLine(label, strategyCommandTextW)); got > strategyCommandTextW {
			t.Errorf("命令 %q 寬度 %d 越過單格 %d", label, got, strategyCommandTextW)
		}
	}
}
