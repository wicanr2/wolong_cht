package main

import (
	"github.com/wicanr2/wolong_cht/internal/rules/speed"
	"image/color"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 自然策略常駐骨架，數值全部出自機器碼（docs/spec/12 §1、docs/re/47）：
// 640×400、32 px 橫幅、命令 (0,32,432,32)、縮小地圖 (432,32,208,160)、
// 自勢力情報 (432,192,208,208)、地圖 (0,32,640,368) ＝ 40×23 格。
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
		{"map y", strategyMapY, 32},
		{"map width", strategyMapW, 640},
		{"map height", strategyMapH, 368},
		{"map cols", viewCols, 40},
		{"map rows", viewRows, 23},
		{"command width", strategyCommandW, 432},
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
		7*(strategyCommandCellW-strategyCommandTextW) + 24; end > strategyCommandW {
		t.Fatalf("命令列字串尾端 = %d，超出框寬 %d", end, strategyCommandW)
	}
	// 熱區照 sub_1E3D7(al=0x0C, X=24, Y=40, 384×16)。
	if r := strategyCommandCellRect(0); r.Min.X != 24 || r.Min.Y != 40 || r.Dy() != 16 {
		t.Fatalf("命令列第一格 = %v，want x24 y40 h16", r)
	}
	if r := strategyCommandCellRect(7); r.Max.X != 24+384 {
		t.Fatalf("命令列最後一格右緣 = %d，want %d", r.Max.X, 24+384)
	}
	// 自勢力情報視窗內部的絕對座標，出自 docs/re/47 §4.2。
	// 用絕對值而不是位移來驗，是因為機器碼給的就是絕對值——
	// 中間多一層減法就多一個算錯的機會。
	inner := []struct {
		name string
		got  int
		want int
	}{
		{"君主列 y", strategyFactionY + strategyInfoYOffset, 208},
		{"軍師列 y", strategyFactionY + strategyInfoYOffset + 2*strategyInfoRowStep, 240},
		{"值欄 x", strategySidebarX + strategyInfoValueXOffset, 576},
		{"信賴度量條 x", strategySidebarX + strategyTrustXOffset, 456},
		{"信賴度量條 y", strategyFactionY + strategyTrustYOffset, 292},
		{"資金 x", strategySidebarX + strategyFundsXOffset, 560},
		{"資金 y", strategyFactionY + strategyFundsYOffset, 312},
		{"預備兵 x", strategySidebarX + strategyReserveXOffset, 568},
		{"預備兵首列 y", strategyFactionY + strategyReserveYOffset, 328},
		{"預備兵末列 y", strategyFactionY + strategyReserveYOffset + 2*strategyResourceRowStep, 360},
		// 顯示清單場景 0（docs/re/48 §3）：標籤與兩個底層方塊。
		{"「君主」標籤 x", strategySidebarX + strategyInfoLabelXOffset, 512},
		{"「信賴度」標籤 x", strategySidebarX + strategyTrustLabelX, 448},
		{"「信賴度」標籤 y", strategyFactionY + strategyTrustLabelY, 272},
		{"信賴度槽 x", strategySidebarX + strategyTrustSlotX, 448},
		{"信賴度槽 y", strategyFactionY + strategyTrustSlotY, 288},
		{"資源黑底 x", strategySidebarX + strategyResourceBoxX, 448},
		{"資源黑底 y", strategyFactionY + strategyResourceBoxY, 304},
		{"「資金」標籤 x", strategySidebarX + strategyResourceLabelX, 456},
		{"圖形欄 x", strategySidebarX + strategyIconXOffset, 528},
	}
	for _, c := range inner {
		if c.got != c.want {
			t.Errorf("%s = %d，want %d", c.name, c.got, c.want)
		}
	}
	// 資金 7 位與預備兵 6 位在原版都右對齊到 x=616。
	if r := strategySidebarX + strategyFundsXOffset + strategyFundsDigits*textdraw.HalfW; r != 616 {
		t.Errorf("資金右端 = %d，want 616", r)
	}
	if r := strategySidebarX + strategyReserveXOffset + strategyReserveDigits*textdraw.HalfW; r != 616 {
		t.Errorf("預備兵右端 = %d，want 616", r)
	}
}

func TestDOSVNaturalStrategyTextContainment(t *testing.T) {
	// 這些是 active 8/16 px bitmap 字寬的 containment 測試，不是新增的原版數值證據。
	//
	// 命令列那一項**是**原版數值：`sub_106F5(dx=8, bx=0x28)` → X=8、Y=40，
	// `cs:6181h` 的詞是兩個全形字（32 px），節距 48（詞 32 ＋ 全形空格 16）。
	if strategyCommandTextW != 32 || strategyCommandCellW != 48 ||
		strategyCommandX != 8 || strategyCommandLead != 24 ||
		strategyInfoValueW != 56 {
		t.Fatalf("命令列版面 = x%d lead%d cell%d text%d／info %d，want 8/24/48/32/56",
			strategyCommandX, strategyCommandLead, strategyCommandCellW,
			strategyCommandTextW, strategyInfoValueW)
	}
	// 第八個詞的右緣，以及尾端「 　」之後的位置，都要在框寬 432 內。
	if right := strategyCommandX + strategyCommandLead +
		7*strategyCommandCellW + strategyCommandTextW; right != 400 {
		t.Fatalf("第八個詞右緣 = %d，want 400", right)
	}
	if strategyFundsDigits != 7 || strategyReserveDigits != 6 {
		t.Fatalf("數字槽 = 資金 %d 位、預備兵 %d 位，want 7/6",
			strategyFundsDigits, strategyReserveDigits)
	}

	for _, digits := range []int{strategyFundsDigits, strategyReserveDigits} {
		want := digits * textdraw.HalfW
		for _, value := range []int{0, 1, 9999, 99999, -1, -9999} {
			text := strategyHUDNumber(value, digits)
			if got := textdraw.StringWidth(text); got != want {
				t.Errorf("%d 以 %d 位格式化為 %q、寬度 %d，want %d", value, digits, text, got, want)
			}
		}
	}
	for _, test := range []struct {
		value  int
		digits int
		want   string
	}{
		{10000000, 7, "9999999"}, {-1000000, 7, "-999999"},
		{1000000, 6, "999999"}, {-100000, 6, "-99999"},
	} {
		if got := strategyHUDNumber(test.value, test.digits); got != test.want {
			t.Errorf("超出 %d 槽的 %d 顯示為 %q，want %q", test.digits, test.value, got, test.want)
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

// 財政視窗的版面契約（docs/spec/14）。視窗矩形來自 sub_1895D(cx=0A15h)，
// 數值座標由 sub_16846 的 VRAM 位移換算（一列 80 byte）。
func TestFinanceWindowLayout(t *testing.T) {
	checks := []struct {
		name string
		got  int
		want int
	}{
		{"視窗 x", financeWinX, 16},
		{"視窗 y", financeWinY, 80},
		{"視窗寬", financeWinW, 336},
		{"視窗高", financeWinH, 160},
		{"資金值 x", financeFundsValueX, 104},
		{"資金值 y", financeFundsValueY, 112},
		{"收入／支出值 x", financeIncomeValueX, 288},
		{"今月底欄值 x", financeValueThisX, 120},
		{"次月欄值 x", financeValueNextX, 280},
		{"首列 y", financeRowY, 160},
		{"末列 y", financeRowY + (financeRows-1)*financeRowStep, 208},
		{"綠色圖示欄 x", financeIconNextX, 256},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %d，want %d", c.name, c.got, c.want)
		}
	}
	// 視窗要裝得下所有內容：最後一列的框下緣是 224，視窗到 240。
	if bottom := financeRowY + 64; bottom > financeWinY+financeWinH {
		t.Errorf("兩欄的框到 %d，超出視窗下緣 %d", bottom, financeWinY+financeWinH)
	}
	// 熱區與綠色圖示欄逐格重合——原版 sub_168C3 的四個熱區就是那四格。
	for i := 0; i < financeRows; i++ {
		r := financeRowRect(i)
		if r.Min.X != financeIconNextX || r.Dx() != 24 || r.Dy() != 16 {
			t.Fatalf("第 %d 列熱區 = %v，want x%d 24×16", i, r, financeIconNextX)
		}
		if r.Min.Y != financeRowY+i*financeRowStep {
			t.Errorf("第 %d 列熱區 y = %d，want %d", i, r.Min.Y, financeRowY+i*financeRowStep)
		}
	}
	// 資金 7 位、收入／支出 6 位、兩欄各 5 位（sub_16846 的 bx 低 byte）。
	if financeFundsDigits != 7 || financeAmountDigits != 6 || financeRowDigits != 5 {
		t.Errorf("位數 = %d/%d/%d，want 7/6/5",
			financeFundsDigits, financeAmountDigits, financeRowDigits)
	}
}

// TestCorpsFormationLayout 把編成視窗的版面釘成契約（docs/spec/22）。
//
// 座標全部出自機器碼：視窗來自 `sub_1895D(cx=0C0Fh)`，靜態層是顯示清單
// 場景 5，數值座標由 `sub_16D6F`／`sub_16DA8` 的 VRAM 位移換算。
func TestCorpsFormationLayout(t *testing.T) {
	checks := []struct {
		name string
		got  int
		want int
	}{
		{"視窗 x", formWinX, 144},
		{"視窗 y", formWinY, 112},
		{"視窗寬", formWinW, 240},
		{"視窗高", formWinH, 192},
		{"武將名 x", formNameX, 296},
		{"總兵力值 x", formTotalValueX, 312},
		{"總兵力 y", formTotalY, 152},
		{"士氣值 x", formMoraleValueX, 320},
		{"士氣值 y", formMoraleY, 168},
		{"槽標籤 x", formSlotLabelX, 160},
		{"槽圖示 x", formSlotIconX, 200},
		{"槽數值 x", formSlotValueX, 232},
		{"首槽 y", formSlotY, 192},
		{"末槽 y", formSlotY + (army.Positions-1)*formSlotStep, 272},
		{"預備兵圖示 x", formReserveIconX, 280},
		{"預備兵數值 x", formReserveValueX, 312},
		{"預備兵首列 y", formReserveY, 216},
		{"確定鈕 x", formOKX, 280},
		{"確定鈕 y", formOKY, 272},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %d，want %d", c.name, c.got, c.want)
		}
	}
	// 熱區與兵種圖示逐格重合——原版 sub_16DFD 登記的六個 24×16 就是那六格。
	for k := 0; k < army.Positions; k++ {
		r := formSlotRect(k)
		if r.Min.X != formSlotIconX || r.Dx() != 24 || r.Dy() != 16 {
			t.Fatalf("第 %d 槽熱區 = %v，want x%d 24×16", k, r, formSlotIconX)
		}
		if r.Min.Y != formSlotY+k*formSlotStep {
			t.Errorf("第 %d 槽熱區 y = %d，want %d", k, r.Min.Y, formSlotY+k*formSlotStep)
		}
	}
	// 數字是半形 8 px：預備兵 6 位從 312 到 360，落在 (304,216) 64×48 的框內。
	if right := formReserveValueX + formReserveDigits*textdraw.HalfW; right > 304+64 {
		t.Errorf("預備兵數到 %d，超出值框右緣 %d", right, 304+64)
	}
	// 六槽的框是 (160,192) 112×96，末槽的數字不能溢出去。
	if right := formSlotValueX + formSlotDigits*textdraw.HalfW; right > formSlotLabelX+112 {
		t.Errorf("槽的兵力到 %d，超出六槽框右緣 %d", right, formSlotLabelX+112)
	}
	// 位數：總兵力 4、士氣 3、預備兵 6、槽 4（sub_16D6F/sub_16DA8 的 bx 低 byte）。
	if formTotalDigits != 4 || formMoraleDigits != 3 ||
		formReserveDigits != 6 || formSlotDigits != 4 {
		t.Errorf("位數 = %d/%d/%d/%d，want 4/3/6/4", formTotalDigits,
			formMoraleDigits, formReserveDigits, formSlotDigits)
	}
}

// 據點情報視窗的版面契約（docs/spec/23）。視窗矩形來自 sub_1895D(cx=810h)，
// 數值座標由 sub_17E4A 的 VRAM 位移換算。
func TestCityInfoWindowLayout(t *testing.T) {
	checks := []struct {
		name string
		got  int
		want int
	}{
		{"視窗 x", cityWinX, 0},
		{"視窗 y", cityWinY, 272},
		{"視窗寬", cityWinW, 256},
		{"視窗高", cityWinH, 128},
		{"景觀圖 x", cityViewX, 16},
		{"景觀圖 y", cityViewY, 288},
		{"景觀圖邊長", cityViewSize, 96},
		{"據點名 x", cityNameX, 128},
		{"據點名 y", cityNameY, 288},
		{"類型 x", cityKindX, 192},
		{"城主 y", cityLordY, 304},
		{"標籤 x", cityLabelX, 128},
		{"首列 y", cityRowY, 320},
		{"末列 y", cityRowY + (cityRows-1)*cityRowStep, 368},
		{"城兵數 x", cityGarrisonX, 208},
		{"生產力 x", cityProductionX, 192},
		{"上昇值 x", cityGrowthX, 208},
		{"防災值 x", cityPreventionX, 216},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %d，want %d", c.name, c.got, c.want)
		}
	}
	// 四個值的右端都落在 240——值框是 (128,320) 112×64，右緣正好 240。
	// 這不是巧合，是原版四組（座標，位數）湊出來的，拿來當算術檢查。
	for _, v := range []struct {
		name   string
		x      int
		digits int
	}{
		{"城兵數", cityGarrisonX, cityGarrisonDigits},
		{"生產力", cityProductionX, cityProductionDigits},
		{"上昇值", cityGrowthX, cityGrowthDigits},
		{"防災值", cityPreventionX, cityPreventionDigits},
	} {
		if right := v.x + v.digits*textdraw.HalfW; right != cityLabelX+112 {
			t.Errorf("%s 右端 = %d，want %d", v.name, right, cityLabelX+112)
		}
	}
	// 景觀圖不能壓到右半的標籤欄。
	if cityViewX+cityViewSize > cityLabelX {
		t.Errorf("景觀圖右緣 %d 越過標籤欄 %d", cityViewX+cityViewSize, cityLabelX)
	}
	// 類型字串：六個詞，首都覆寫（docs/re/50 §2.1）。
	for kind, want := range map[int]string{0: "大都市", 1: "中都市", 2: "小都市", 3: "關卡", 4: "戰場"} {
		if got := cityKindLabel(kind, false); got != want {
			t.Errorf("類型 %d = %q，want %q", kind, got, want)
		}
		if got := cityKindLabel(kind, true); got != "首都" {
			t.Errorf("類型 %d 是首都時 = %q，want 首都", kind, got)
		}
	}
}

// 軍團情報視窗的版面契約（docs/spec/24）。視窗矩形來自 sub_1895D(cx=0D0Dh)，
// 數值座標由 sub_1807B／sub_1812A 的 VRAM 位移換算。
func TestCorpsInfoWindowLayout(t *testing.T) {
	checks := []struct {
		name string
		got  int
		want int
	}{
		{"視窗 x", corpsWinX, 432},
		{"視窗 y", corpsWinY, 192},
		{"視窗寬", corpsWinW, 208},
		{"視窗高", corpsWinH, 208},
		{"頭像 x", corpsPortraitX, 440},
		{"頭像 y", corpsPortraitY, 200},
		{"標籤 x", corpsHeadLabelX, 512},
		{"名字 x", corpsHeadValueX, 576},
		{"垂直線 x", corpsDividerX, 560},
		{"總兵力 y", corpsTotalY, 272},
		{"總兵力值 x", corpsTotalX, 536},
		{"斜線 x", corpsSlashX, 568},
		{"士氣 x", corpsMoraleX, 584},
		{"槽標籤 x", corpsSlotLabelX, 464},
		{"槽圖示 x", corpsSlotIconX, 520},
		{"槽數值 x", corpsSlotValueX, 576},
		{"首槽 y", corpsSlotY, 288},
		{"末槽 y", corpsSlotY + (army.Positions-1)*corpsSlotStep, 368},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %d，want %d", c.name, c.got, c.want)
		}
	}
	// 「總兵力 6000／200」是一列三段，接續但不重疊：
	// 4 位半形數字到 568 剛好接上全形斜線，斜線 16 px 之後是士氣。
	if right := corpsTotalX + corpsTotalDigits*textdraw.HalfW; right != corpsSlashX {
		t.Errorf("總兵力右端 = %d，want 接在斜線 %d", right, corpsSlashX)
	}
	if corpsSlashX+textdraw.GlyphW != corpsMoraleX {
		t.Errorf("斜線右端 = %d，want 接在士氣 %d", corpsSlashX+textdraw.GlyphW, corpsMoraleX)
	}
	// 六槽的數字 4 位到 608，落在視窗（右緣 640）內。
	if right := corpsSlotValueX + corpsSlotDigits*textdraw.HalfW; right > corpsWinX+corpsWinW {
		t.Errorf("槽的兵力到 %d，超出視窗右緣 %d", right, corpsWinX+corpsWinW)
	}
	// 這個視窗與自勢力情報是同一格——原版就是蓋上去的。
	if corpsWinX != strategySidebarX || corpsWinY != strategyFactionY ||
		corpsWinW != strategySidebarW || corpsWinH != strategyFactionH {
		t.Errorf("軍團情報 (%d,%d,%d,%d) 與自勢力情報 (%d,%d,%d,%d) 不同格",
			corpsWinX, corpsWinY, corpsWinW, corpsWinH,
			strategySidebarX, strategyFactionY, strategySidebarW, strategyFactionH)
	}
}

// 四槽選擇視窗的版面契約（docs/spec/25）。視窗矩形來自 sub_1895D(cx=0F13h)，
// 數值座標由 sub_18C20 的 VRAM 位移換算。
func TestSlotSelectWindowLayout(t *testing.T) {
	checks := []struct {
		name string
		got  int
		want int
	}{
		{"視窗 x", savePanelX, 96},
		{"視窗 y", savePanelY, 80},
		{"視窗寬", savePanelW, 304},
		{"視窗高", savePanelH, 240},
		{"標題 x", saveTitleX, 184},
		{"標題 y", saveTitleY, 91},
		{"水平線 y", saveRuleY, 111},
		{"水平線長", saveRuleW, 287},
		{"名稱欄 x", saveNameBoxX, 120},
		{"名稱欄 y", saveNameBoxY, 118},
		{"名稱 y", saveNameY, 120},
		{"日期欄 x", saveSlotX, 256},
		{"日期欄 y", saveSlotY, 144},
		{"列距", saveSlotStep, 48},
		{"年 x", saveYearX, 264},
		{"月 x", saveMonthX, 304},
		{"日 x", saveDayX, 336},
		{"「年月日」x", saveDateLabelX, 288},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %d，want %d", c.name, c.got, c.want)
		}
	}
	// 三個數字各自嵌在「年　月　日」那三個字的前面：
	// 年 3 位到 288（「年」），月 2 位到 320（「月」），日 2 位到 352（「日」）。
	if right := saveYearX + saveYearDigits*textdraw.HalfW; right != saveDateLabelX {
		t.Errorf("年的右端 = %d，want 接在「年」%d", right, saveDateLabelX)
	}
	if right := saveMonthX + saveMonthDigits*textdraw.HalfW; right != saveDateLabelX+2*textdraw.GlyphW {
		t.Errorf("月的右端 = %d，want %d", right, saveDateLabelX+2*textdraw.GlyphW)
	}
	if right := saveDayX + saveMonthDigits*textdraw.HalfW; right != saveDateLabelX+4*textdraw.GlyphW {
		t.Errorf("日的右端 = %d，want %d", right, saveDateLabelX+4*textdraw.GlyphW)
	}
	// 四個槽都要落在視窗內（末槽的日期欄下緣 304 < 320）。
	if bottom := saveSlotY + 3*saveSlotStep + saveSlotH; bottom > savePanelY+savePanelH {
		t.Errorf("末槽下緣 %d 超出視窗 %d", bottom, savePanelY+savePanelH)
	}
	// 熱區與日期欄逐格重合。
	for i := 0; i < 4; i++ {
		r := saveSlotRect(i)
		if r.Min.X != saveSlotX || r.Dx() != saveSlotW || r.Dy() != saveSlotH {
			t.Fatalf("第 %d 槽熱區 = %v，want x%d %d×%d", i, r, saveSlotX, saveSlotW, saveSlotH)
		}
	}
}

// ＹＥＳ／ＮＯ 對話框的版面契約（docs/spec/26）。
func TestYesNoDialogLayout(t *testing.T) {
	checks := []struct {
		name string
		got  int
		want int
	}{
		{"框寬", yesNoW, 208},
		{"框高", yesNoH, 96},
		{"水平線 dy", yesNoRuleDY, 31},
		{"水平線長", yesNoRuleW, 191},
		{"選項框 dx", yesNoBoxDX, 40},
		{"ＹＥＳ dy", yesNoYesDY, 40},
		{"ＮＯ dy", yesNoNoDY, 64},
		{"選項框寬", yesNoBoxW, 128},
		{"選項框高", yesNoBoxH, 16},
		{"文字 dx", yesNoTextDX, 80},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %d，want %d", c.name, c.got, c.want)
		}
	}
	// 兩個選項框都要落在框內，且中間留得下那條 8 px 的縫。
	if yesNoNoDY+yesNoBoxH > yesNoH {
		t.Errorf("ＮＯ 的下緣 %d 超出框高 %d", yesNoNoDY+yesNoBoxH, yesNoH)
	}
	if gap := yesNoNoDY - (yesNoYesDY + yesNoBoxH); gap != 8 {
		t.Errorf("兩個選項之間的縫 = %d，want 8", gap)
	}
	if r := yesNoRect(216, 152); r.Dx() != yesNoW || r.Dy() != yesNoH {
		t.Errorf("yesNoRect = %v", r)
	}
}

// ⭐ 命中算式：原版把 Y 除以 8，**第 2 條是縫**——兩個選項都不算，
// 不是四捨五入到最近的按鈕（docs/spec/26 §2）。
func TestYesNoHitGap(t *testing.T) {
	const x, y = 216, 152
	yesY := y + yesNoYesDY
	for _, c := range []struct {
		name    string
		px, py  int
		wantHit bool
		wantYes bool
	}{
		{"ＹＥＳ 上緣", x + yesNoBoxDX, yesY, true, true},
		{"ＹＥＳ 下緣", x + yesNoBoxDX + 127, yesY + 15, true, true},
		{"縫的上緣", x + yesNoBoxDX, yesY + 16, false, false},
		{"縫的下緣", x + yesNoBoxDX, yesY + 23, false, false},
		{"ＮＯ 上緣", x + yesNoBoxDX, yesY + 24, true, false},
		{"ＮＯ 下緣", x + yesNoBoxDX, yesY + 71, true, false},
		{"左邊界外", x + yesNoBoxDX - 1, yesY, false, false},
		{"右邊界外", x + yesNoBoxDX + 128, yesY, false, false},
		{"上邊界外", x + yesNoBoxDX, yesY - 1, false, false},
		{"下邊界外", x + yesNoBoxDX, yesY + 72, false, false},
	} {
		hit, yes := hitTestYesNo(x, y, c.px, c.py)
		if hit != c.wantHit || (hit && yes != c.wantYes) {
			t.Errorf("%s (%d,%d) = hit%v yes%v，want hit%v yes%v",
				c.name, c.px, c.py, hit, yes, c.wantHit, c.wantYes)
		}
	}
}

// 君主選擇卡的版面契約（docs/spec/27）。視窗矩形來自 sub_1895D(cx=0C0Fh)，
// 數值座標由 sub_18EA0 的 VRAM 位移換算。
func TestLordCardLayout(t *testing.T) {
	checks := []struct {
		name string
		got  int
		want int
	}{
		{"視窗 x", lordCardX, 160},
		{"視窗 y", lordCardY, 112},
		{"視窗寬", lordCardW, 240},
		{"視窗高", lordCardH, 192},
		{"君主頭像 x", lordPortraitX, 184},
		{"君主頭像 y", lordPortraitY, 128},
		{"君主名 x", lordNameX, 200},
		{"君主名 y", lordNameY, 216},
		{"軍師頭像 x", lordAdvPortraitX, 312},
		{"軍師頭像 y", lordAdvPortraitY, 168},
		{"軍師名 x", lordAdvNameX, 328},
		{"軍師名 y", lordAdvNameY, 144},
		{"標籤 x", lordLabelX, 184},
		{"首都 y", lordCapitalY, 240},
		{"武將數 y", lordGeneralsY, 256},
		{"據點數 y", lordCitiesY, 272},
		{"首都名 x", lordCapitalNameX, 264},
		{"數字 x", lordCountX, 272},
		{"垂直線 x", lordDividerX, 247},
		{"自定 y", lordCustomY, 248},
		{"確定 y", lordOKY, 272},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %d，want %d", c.name, c.got, c.want)
		}
	}
	// ⚠ 熱區編號與畫面上下顛倒：0x20 是下面的「確定」、0x21 是上面的「自定」。
	// 這裡至少要保證兩顆按鈕的上下關係沒被寫反。
	if lordCustomY >= lordOKY {
		t.Errorf("自定 (%d) 應該在確定 (%d) 上面", lordCustomY, lordOKY)
	}
	// 兩顆按鈕與左下的資訊框都要落在視窗內。
	if right := lordOKX + lordButtonW; right > lordCardX+lordCardW {
		t.Errorf("按鈕右緣 %d 超出視窗 %d", right, lordCardX+lordCardW)
	}
	if bottom := lordOKY + lordButtonH; bottom > lordCardY+lordCardH {
		t.Errorf("按鈕下緣 %d 超出視窗 %d", bottom, lordCardY+lordCardH)
	}
	// 三位數的右端不能穿過「自定／確定」那一欄（左緣 328）。
	if right := lordCountX + lordCountDigits*textdraw.HalfW; right > lordCustomX {
		t.Errorf("數字右端 %d 穿到按鈕欄 %d", right, lordCustomX)
	}
}

// 系統選單的版面契約（docs/spec/13 §2.6、docs/re/55）。
// 視窗矩形來自 sub_1895D(cx=0C0Dh)，六列由顯示清單場景 2 給。
// ⚠ remake 在那六列**後面**多加一列「主君編成」（docs/spec/76），
// 所以視窗比原版高一個列距；原版那六列的座標一格都沒動。
func TestSystemMenuLayout(t *testing.T) {
	checks := []struct {
		name string
		got  int
		want int
	}{
		{"視窗 x", sysWinX, 208},
		{"視窗 y", sysWinY, 112},
		{"視窗寬", sysWinW, 208},
		// ⚠ 原版是 192（六列）。remake 多了「主君編成」（docs/spec/76）
		// 與「損害報告」（docs/spec/89）兩列，所以是 192 ＋ 兩個列距。
		// **原版值仍然釘在算式裡**，改列距或列數都會被這一條抓到。
		{"視窗高", sysWinH, 192 + 2*sysRowStep},
		{"標題 x", sysTitleX, 228},
		{"標題 y", sysTitleY, 124},
		{"水平線 y", sysRuleY, 142},
		{"水平線長", sysRuleW, 191},
		{"標籤底 x", sysLabelBoxX, 222},
		{"標籤底 y", sysLabelBoxY, 150},
		{"標籤 x", sysLabelX, 232},
		{"標籤 y", sysLabelY, 152},
		{"值格 x", sysValueX, 352},
		{"值格 y", sysValueY, 152},
		{"列距", sysRowStep, 24},
		// ⚠ 原版六列（熱區 0x20–0x25）；第 7、8 列是 remake 加的。
		{"列數", sysRows, 8},
		// ⭐ **原版六列一格都沒動**：遊戲結束仍是索引 5，新列加在後面。
		{"遊戲結束仍在原版的第 6 列", sysRowQuit, 5},
		{"主君編成是加在最後的第 7 列", sysRowLordCorps, 6},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %d，want %d", c.name, c.got, c.want)
		}
	}
	// 每一列都要落在視窗內。
	if bottom := sysLabelBoxY + (sysRows-1)*sysRowStep + sysLabelBoxH; bottom > sysWinY+sysWinH {
		t.Errorf("末列下緣 %d 超出視窗 %d", bottom, sysWinY+sysWinH)
	}
	// 值格不能疊到左邊的標籤底（128 寬，右緣 350）。
	if sysValueX < sysLabelBoxX+sysLabelBoxW {
		t.Errorf("值格 x=%d 疊到標籤底右緣 %d", sysValueX, sysLabelBoxX+sysLabelBoxW)
	}
	// 熱區與值格逐格重合、列距 24（原版 0x20–0x25 ＋ remake 多的那一格）。
	for k := 0; k < sysRows; k++ {
		r := sysRowRect(k)
		if r.Min.X != sysValueX || r.Dx() != sysValueW || r.Dy() != sysValueH {
			t.Fatalf("第 %d 列熱區 = %v，want x%d %d×%d", k, r, sysValueX, sysValueW, sysValueH)
		}
		if r.Min.Y != sysValueY+k*sysRowStep {
			t.Errorf("第 %d 列熱區 y = %d，want %d", k, r.Min.Y, sysValueY+k*sysRowStep)
		}
	}
	// ⭐ **原版那六列的標籤與順序一格都不能動**——沒接功能的也要留著。
	// remake 加的那一列只能排在它們後面（docs/spec/76）。
	if sysMenuLabels[0] != "資料儲存" || sysMenuLabels[sysRowQuit] != "遊戲結束" {
		t.Errorf("原版六列的標籤被動到了：%v", sysMenuLabels)
	}
}

// 戰略速度與戰術速度是**兩個獨立設定**（原版說明書 3.5、系統選單第 4／5 列）。
// remake 先前共用一個變數，選單那兩列因此永遠一樣。
func TestSpeedsAreIndependent(t *testing.T) {
	g := &game{speed: 2, tacticalSpeed: 2}
	// ＋ ＝ 更快 ＝ 檔位往 0 走（原版的檔位越大越慢）。
	g.adjustSpeed(false, 1)
	if g.speed != 1 || g.tacticalSpeed != 2 {
		t.Fatalf("調戰略速度動到了戰術：%d／%d", g.speed, g.tacticalSpeed)
	}
	g.adjustSpeed(true, -2)
	if g.speed != 1 || g.tacticalSpeed != 4 {
		t.Fatalf("調戰術速度動到了戰略：%d／%d", g.speed, g.tacticalSpeed)
	}
	// 兩端都要夾在原版的 0–4，沒有第六檔也沒有負檔。
	g.adjustSpeed(false, -100)
	g.adjustSpeed(true, 1000)
	if g.speed != speed.Levels-1 || g.tacticalSpeed != 0 {
		t.Fatalf("沒有夾在 0–%d：%d／%d", speed.Levels-1, g.speed, g.tacticalSpeed)
	}
}

// 系統選單的速度列是**五檔循環**（原版 `sub_16062`：+1，到頂繞回 0），
// 不是 ±1 的數值微調。標籤在 `ds:6033h`，選項數 5 出自 `ds:5FF4h`。
func TestSystemRowCyclesSpeedThroughFiveOriginalSteps(t *testing.T) {
	g := &game{}
	for i := 1; i <= speed.Levels; i++ {
		g.dispatchSystemRow(sysRowStrategySpeed, true)
		if want := i % speed.Levels; g.speed != want {
			t.Fatalf("第 %d 次左鍵：檔位 ＝ %d，預期 %d", i, g.speed, want)
		}
	}
	if g.speed != 0 {
		t.Errorf("繞一圈之後應該回到第 0 檔，得到 %d", g.speed)
	}
	// 右鍵往回一檔（remake 的便利，不是原版行為）。
	g.dispatchSystemRow(sysRowStrategySpeed, false)
	if g.speed != speed.Levels-1 {
		t.Errorf("右鍵應該退到最後一檔 %d，得到 %d", speed.Levels-1, g.speed)
	}
	// 戰術速度是另一個設定，不能被戰略速度帶著動。
	if g.tacticalSpeed != 0 {
		t.Errorf("戰術速度被連動了：%d", g.tacticalSpeed)
	}
	g.dispatchSystemRow(sysRowTacticalSpeed, true)
	if g.tacticalSpeed != 1 {
		t.Errorf("戰術速度 ＝ %d，預期 1", g.tacticalSpeed)
	}
}

// 五個標籤照 `ds:6033h` 的順序，值就是檔位本身。
func TestSpeedLabelsMatchOriginalFiveSteps(t *testing.T) {
	if len(speed.Labels) != speed.Levels {
		t.Fatalf("原版是 %d 檔，標籤有 %d 個", speed.Levels, len(speed.Labels))
	}
	want := [speed.Levels]string{"最高速", " 高速 ", " 普通 ", " 低速 ", "最低速"}
	if speed.Labels != want {
		t.Errorf("標籤 = %v，預期 ds:6033h 的 %v", speed.Labels, want)
	}
}

// 節流換算要對得上原版的實際速率：戰略每個子刻等「檔位」個 291.3 Hz 中斷。
// 出處 docs/re/61 §4、docs/spec/34 §3。
func TestSpeedThrottleMatchesOriginalRates(t *testing.T) {
	// 每檔十秒（600 個畫面更新）該推進幾個子刻 ＝ 2913 ÷ 檔位。
	// 量十秒而不是一秒，是為了讓「累加器剩下不滿一步」的量化誤差
	// 小於要驗的 1%——一秒的殘量在最低速那一檔就佔 1.1%。
	for level := 1; level < speed.Levels; level++ {
		var th speed.Throttle
		got := 0
		for f := 0; f < 600; f++ {
			got += th.Steps(level, 1, speed.HighSpeedStrategy)
		}
		want := 2913.0 / float64(level)
		if d := float64(got) - want; d > want/100 || d < -want/100 {
			t.Errorf("檔位 %d：十秒 %d 個子刻，原版 %.1f", level, got, want)
		}
	}
	// 檔位 0 是 remake 差異：原版不等待，這裡給固定上限。
	var th speed.Throttle
	if n := th.Steps(0, 1, speed.HighSpeedStrategy); n != speed.HighSpeedStrategy {
		t.Errorf("最高速 = %d，預期 %d", n, speed.HighSpeedStrategy)
	}
}

// 戰術層的值要先 ×16 才當等待量（`sub_160A5`），所以同檔位下
// 戰場幀數是戰略子刻數的 1/16——「高速」正好是 18.2 fps。
func TestTacticalThrottleIsSixteenTimesCoarser(t *testing.T) {
	for level := 1; level < speed.Levels; level++ {
		var strat, tact speed.Throttle
		s, t2 := 0, 0
		for f := 0; f < 600; f++ {
			s += strat.Steps(level, 1, speed.HighSpeedStrategy)
			t2 += tact.Steps(level, speed.TacticalMul, speed.HighSpeedTactical)
		}
		if d := s - t2*speed.TacticalMul; d > speed.TacticalMul || d < -speed.TacticalMul {
			t.Errorf("檔位 %d：戰略 %d 子刻／戰術 %d 幀，差得太多", level, s, t2)
		}
	}
	// 高速（檔位 1）＝ 18.2 fps ＝ 標準 BIOS tick。
	var th speed.Throttle
	n := 0
	for f := 0; f < 60; f++ {
		n += th.Steps(1, speed.TacticalMul, speed.HighSpeedTactical)
	}
	if n != 18 {
		t.Errorf("高速一秒 %d 幀，預期 18（18.2 fps）", n)
	}
}

// 縮小地圖標記的四種顏色出自 `sub_15CE0`（docs/re/62 §2）：
// 無所屬 0x0F、自勢力 0xAC、盯著的勢力 0xF3、其餘 0x83。
func TestMinimapMarkerColoursFollowOwnership(t *testing.T) {
	w := &state.World{Player: 1}
	for i := range w.Factions {
		w.Factions[i].Alive = true
	}
	g := &game{world: w, minimapFaction: 3}
	// 六個色號各給一個可分辨的假色，才看得出挑錯。
	for i := range g.minimapInk {
		g.minimapInk[i] = color.RGBA{uint8(i + 1), 0, 0, 255}
	}
	cases := []struct {
		name           string
		owner          int
		border, centre int
	}{
		{"無所屬", 24, 0, 1},
		{"自勢力", 1, 2, 3},
		{"盯著的勢力", 3, 1, 4},
		{"其餘勢力", 5, 5, 4},
	}
	for _, c := range cases {
		border, centre := g.minimapMarkerColours(c.owner)
		if border != g.minimapInk[c.border] || centre != g.minimapInk[c.centre] {
			t.Errorf("%s：外框 %v／中心 %v，預期 minimapInk[%d]／[%d]",
				c.name, border, centre, c.border, c.centre)
		}
	}
}

// 標記座標 ＝ 原點 + 格/2 − 2（docs/re/62 §2 的 `1B6h`／`26h`）。
func TestMinimapMarkerGeometryMatchesRaw(t *testing.T) {
	// 原版的 x 原點是 0x1B6 ＝ 438 ＝ 地圖區 440 減 2；y 是 0x26 ＝ 38 ＝ 40 減 2。
	if x, y := minimapMarkerPos(0, 0); x != strategyMinimapX-2 || y != strategyMinimapY-2 {
		t.Errorf("(0,0) → (%d,%d)，預期 (%d,%d)", x, y, strategyMinimapX-2, strategyMinimapY-2)
	}
	// 世界 384×256 格對地圖區 192×128 px：右下角那一格要落在區內。
	x, y := minimapMarkerPos(383, 255)
	if x != strategyMinimapX+191-2 || y != strategyMinimapY+127-2 {
		t.Errorf("(383,255) → (%d,%d)，預期 (%d,%d)",
			x, y, strategyMinimapX+189, strategyMinimapY+125)
	}
}

// 圖例第二格不能停在自勢力，也不能停在已滅亡的勢力
// （原版 `sub_15AFC` 明文擋掉自勢力，docs/re/62 §4.2）。
func TestMinimapFactionSkipsSelfAndDead(t *testing.T) {
	w := &state.World{Player: 0}
	w.Factions[0].Alive = true
	w.Factions[2].Alive = true
	g := &game{world: w, minimapFaction: 0}
	// ⭐ **開局盯的就是勢力 0，即使那是自己。** 原版沒有初始化
	// `cs:byte_198A7`，資料段開機是 0——實機截圖上圖例兩格都是「曹操」
	// （docs/playtest/38）。擋自勢力的是選單，不是初值。
	if got := g.watchedFaction(); got != 0 {
		t.Errorf("watchedFaction = %d，預期 0（開局盯的就是勢力 0）", got)
	}
	// 換一個要走選單，而選單擋掉自己（原版 `sub_15AFC` 的
	// `cmp al, cs:byte_10CFF`，見 factionpicker_test.go）。
	if g.pickerSelectable(0) {
		t.Error("自勢力不該可選")
	}
	if !g.pickerSelectable(2) {
		t.Error("勢力 2 活著又不是自己，應該可選")
	}
	// 盯著的勢力滅亡就往後找，找到還活著的 0。
	w.Factions[2].Alive = false
	if got := g.watchedFaction(); got != 0 {
		t.Errorf("盯著的滅亡後應往後找到 0，得到 %d", got)
	}
	// 全部滅亡時不能無限迴圈。
	w.Factions[0].Alive = false
	if got := g.watchedFaction(); got != -1 {
		t.Errorf("沒有活著的勢力時應回 −1，得到 %d", got)
	}
}

// 圖例的右半格是原版熱區 0x17 ＝ (536, 168, 96, 16)。
func TestMinimapLegendHitBoxMatchesRawHotzone(t *testing.T) {
	if strategyMinimapLegendY != 168 {
		t.Fatalf("圖例 y = %d，原版是 168", strategyMinimapLegendY)
	}
	inside := [][2]int{{536, 168}, {631, 183}}
	outside := [][2]int{{535, 168}, {536, 167}, {632, 168}, {536, 184}}
	for _, p := range inside {
		if !hitTestMinimapLegend(p[0], p[1]) {
			t.Errorf("(%d,%d) 應該在熱區內", p[0], p[1])
		}
	}
	for _, p := range outside {
		if hitTestMinimapLegend(p[0], p[1]) {
			t.Errorf("(%d,%d) 不該在熱區內", p[0], p[1])
		}
	}
}

// TestFormSlotHitBoxesMatchRaw 釘住編成畫面的滑鼠熱區（docs/re/49 §5）：
// 六個槽 (200, 192+16k) 24×16、確定鈕 (280, 272) 88×16。
func TestFormSlotHitBoxesMatchRaw(t *testing.T) {
	for k := 0; k < army.Positions; k++ {
		y := 192 + k*16
		for _, p := range [][2]int{{200, y}, {223, y + 15}} {
			got, ok := formSlotAt(p[0], p[1])
			if !ok || got != k {
				t.Errorf("(%d,%d) → 槽 %d ok=%v，want 槽 %d", p[0], p[1], got, ok, k)
			}
		}
		// 圖示左邊與右邊各一格外不算，否則相鄰的標籤與數字也會被吃掉。
		for _, p := range [][2]int{{199, y}, {224, y}} {
			if _, ok := formSlotAt(p[0], p[1]); ok {
				t.Errorf("(%d,%d) 不該落在第 %d 槽", p[0], p[1], k)
			}
		}
	}
	ok := formOKRect()
	if ok.Min.X != 280 || ok.Min.Y != 272 || ok.Dx() != 88 || ok.Dy() != 16 {
		t.Errorf("確定鈕熱區 = %v，want (280,272) 88×16", ok)
	}
}

// TestFormSlotClickCyclesKind 釘住原版點一下槽的動作（docs/re/30 §3）：
// 兵種 +1，`1 → 2 → 3 → 1`；**空槽（兵種 4）也落回 1**，
// 所以點過的槽不會再變回空槽。
func TestFormSlotClickCyclesKind(t *testing.T) {
	var f formState
	f.manned[0] = false // 空槽
	f.cycleKind(0)
	if !f.manned[0] || f.kinds[0] != army.Cavalry {
		t.Fatalf("空槽點一下 = manned %v／兵種 %v，want 騎馬", f.manned[0], f.kinds[0])
	}
	want := []army.TroopType{army.Archer, army.Infantry, army.Cavalry}
	for i, w := range want {
		f.cycleKind(0)
		if f.kinds[0] != w {
			t.Errorf("第 %d 次循環 = %v，want %v", i+1, f.kinds[0], w)
		}
		if !f.manned[0] {
			t.Errorf("第 %d 次循環之後變成空槽了", i+1)
		}
	}
}
