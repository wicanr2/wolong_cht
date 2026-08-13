package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// DOS/V YouTube oracle 的自然策略常駐骨架：640×400、32 px banner、
// 左側 432×336 地圖、右側 208 px；上方 minimap 與下方情報框共用分隔邊。
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
		{"minimap height", strategyMinimapH, 152},
		{"minimap legend y", strategyMinimapLegendY, 168},
		{"minimap swatch y", strategyMinimapSwatchY, 172},
		{"faction y", strategyFactionY, 176},
		{"faction height", strategyFactionH, 224},
	}
	for _, check := range checks {
		if check.got != check.want {
			t.Errorf("%s = %d，want %d", check.name, check.got, check.want)
		}
	}
	if got := strategyFactionY + strategyFactionH; got != screenH {
		t.Fatalf("右欄下方情報框未貼齊畫布底部：y+h=%d，want %d", got, screenH)
	}
	if got := strategyMinimapLegendY + textdraw.GlyphH; got > strategyFactionY+chrome.Tile {
		t.Fatalf("minimap 色標越過共用分隔邊：end=%d，shared edge end=%d", got, strategyFactionY+chrome.Tile)
	}
	if strategySidebarInnerX != 440 || strategyMinimapY != 40 || strategyMinimapW != 192 {
		t.Fatalf("DOS/V minimap 內框 = x%d y%d w%d，want x440 y40 w192",
			strategySidebarInnerX, strategyMinimapY, strategyMinimapW)
	}
	if strategyFactionInnerY != 184 || strategyFactionInnerW != 192 {
		t.Fatalf("DOS/V 情報內框 = y%d w%d，want y184 w192", strategyFactionInnerY, strategyFactionInnerW)
	}
	if strategyInfoYOffset != 24 || strategyInfoRowStep != 17 || strategyInfoDividerXOffset != 120 ||
		strategyInfoDividerH != 48 || strategyTrustYOffset != 88 || strategyResourceDividerY != 120 ||
		strategyResourceBoxY != 128 || strategyResourceBoxH != 88 {
		t.Fatal("DOS/V 情報列／分隔線／資源區幾何契約被改動")
	}
}

func TestDOSVNaturalStrategyTextContainment(t *testing.T) {
	// 這些是 active 8/16 px bitmap 字寬的 containment 測試，不是新增的原版數值證據。
	if strategyCommandTextW != 48 || strategyInfoValueW != 72 {
		t.Fatalf("文字安全寬度 = command %d/info %d，want 48/72", strategyCommandTextW, strategyInfoValueW)
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
