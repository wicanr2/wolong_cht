package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
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
