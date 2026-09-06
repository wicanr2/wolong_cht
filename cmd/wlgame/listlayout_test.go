package main

import (
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 羅馬化人名比原版的三個全形字寬，欄寬要裁——但要裁得讀得出來
// （docs/spec/84 §6）。
func TestFitCellTruncatesRomanisedNames(t *testing.T) {
	const room = 56 // 7 個半形格
	for _, tc := range []struct{ in, want string }{
		{"CAO-CAO", "CAO-CAO"},    // 塞得下就原樣
		{"XIAHOU-YUAN", "XIAHOU"}, // 切在連字號上 → 連字號也去掉
		{"ZHUGE-LIANG", "ZHUGE"},  // 名只剩一個字母 → 只留姓
		{"XU-HUANG", "XU-HUAN"},   // 名還剩得夠多 → 照切
		{"曹操", "曹操"},              // 中文照原樣（3 字 48px 本來就塞得下）
	} {
		if got := fitCell(tc.in, room); got != tc.want {
			t.Errorf("fitCell(%q, %d) = %q，預期 %q", tc.in, room, got, tc.want)
		}
	}
	if got := fitCell("CAO-CAO", 0); got != "CAO-CAO" {
		t.Errorf("欄寬 0（未知）時不該裁：%q", got)
	}
}

// 半形語系的欄界（docs/spec/85）：欄不可重疊、不可超出本體、
// 欄數要與中文那一套相同（不同會與 listRow 的格數對不上）。
func TestLatinListFieldsFitTheBody(t *testing.T) {
	// room 是各文字欄該放得下幾個半形字（docs/spec/85 §3 的表）。
	for _, tc := range []struct {
		name   string
		latin  []listField
		labels []string
		zh     []listField
		room   map[int]int
	}{
		{"武將", latinFieldsGenerals, latinLabelsGenerals, listFamilyGenerals.fields(),
			map[int]int{0: 12, 4: 9, 5: 7}},
		{"據點", latinFieldsCities, latinLabelsCities, listFamilyCities.fields(),
			map[int]int{0: 12, 5: 10}},
		{"軍團", latinFieldsCorps, latinLabelsCorps, listFamilyCorps.fields(),
			map[int]int{0: 10, 3: 8, 4: 8, 5: 4}},
		{"勢力", latinFieldsFactions, latinLabelsFactions, listFamilyFactions.fields(),
			map[int]int{0: 9, 3: 8, 4: 4, 5: 8}},
	} {
		if len(tc.latin) != len(tc.zh) {
			t.Errorf("%s：半形 %d 欄、中文 %d 欄", tc.name, len(tc.latin), len(tc.zh))
		}
		if len(tc.labels) != len(tc.latin) {
			t.Errorf("%s：標籤 %d 個、欄位 %d 個", tc.name, len(tc.labels), len(tc.latin))
		}
		for i, f := range tc.latin {
			if i > 0 {
				prev := tc.latin[i-1]
				if f.X < prev.X+prev.W {
					t.Errorf("%s 第 %d 欄 X=%d 疊到前一欄（到 %d）",
						tc.name, i, f.X, prev.X+prev.W)
				}
			}
			if f.X+f.W > listBodyW() {
				t.Errorf("%s 第 %d 欄右緣 %d 超出本體 %d",
					tc.name, i, f.X+f.W, listBodyW())
			}
			// 文字欄的可用寬要與規格逐欄相符——這一格是 spec/85 §3
			// 的覆蓋率算出來的，改動欄界就要回頭改那張表。
			if want, ok := tc.room[i]; ok {
				if got := listCellRoom(tc.latin, i) / textdraw.HalfW; got != want {
					t.Errorf("%s 第 %d 欄放得下 %d 字，規格是 %d 字",
						tc.name, i, got, want)
				}
			}
		}
	}
}

// 標題列的標籤要落在自己欄位的文字起點上，而且不可以互相疊字。
func TestLatinTitleAlignsToColumns(t *testing.T) {
	title := latinTitle(latinFieldsGenerals, latinLabelsGenerals)
	for i, label := range latinLabelsGenerals {
		want := (latinFieldsGenerals[i].X + listTextInset) / textdraw.HalfW
		if got := strings.Index(title, label); got != want {
			t.Errorf("標籤 %q 在第 %d 格，預期第 %d 格（%q）", label, got, want, title)
		}
	}
	if len(title)*textdraw.HalfW > listBodyW() {
		t.Errorf("標題 %d px 超出本體 %d px：%q",
			len(title)*textdraw.HalfW, listBodyW(), title)
	}
}

// 標題的標籤彼此不可相黏——`MarCmdPolFaction` 讀不出來是三欄。
func TestLatinTitlesKeepLabelsApart(t *testing.T) {
	for _, tc := range []struct {
		name   string
		fields []listField
		labels []string
	}{
		{"武將", latinFieldsGenerals, latinLabelsGenerals},
		{"據點", latinFieldsCities, latinLabelsCities},
		{"軍團", latinFieldsCorps, latinLabelsCorps},
		{"勢力", latinFieldsFactions, latinLabelsFactions},
	} {
		title := latinTitle(tc.fields, tc.labels)
		if len(title)*textdraw.HalfW > listBodyW() {
			t.Errorf("%s 標題 %d px 超出本體：%q", tc.name, len(title)*textdraw.HalfW, title)
		}
		for i, label := range tc.labels {
			if label == "" {
				continue
			}
			at := (tc.fields[i].X + listTextInset) / textdraw.HalfW
			if got := strings.Index(title, label); got != at {
				t.Errorf("%s 的 %q 落在第 %d 格，欄位在第 %d 格：%q",
					tc.name, label, got, at, title)
			}
			if at > 0 && title[at-1] != ' ' {
				t.Errorf("%s 的 %q 與前一個標籤相黏：%q", tc.name, label, title)
			}
		}
	}
}

// 空列的數字欄破折號：武將一覽在**欄起點 ＋ 8**（docs/spec/38）。
//
// ⭐ 原版的分隔線與欄界是兩份資料——軍團一覽兩者重合、武將一覽差一個
// 半格。同一條 `Sep` 推不出兩者，所以這裡把量到的絕對座標釘住。
func TestGeneralListEmptyRowDashesMatchOriginal(t *testing.T) {
	f := listFamilyGenerals.fields()
	// 三個數字欄（武術／統率／政治）。
	var numeric []listField
	for _, c := range f {
		if c.Numeric {
			numeric = append(numeric, c)
		}
	}
	if len(numeric) != 3 {
		t.Fatalf("武將一覽有 %d 個數字欄，want 3", len(numeric))
	}
	// 量到的破折號左緣（絕對，body ＝ 40）。
	for i, want := range []int{120, 160, 200} {
		got := listBodyX() + numeric[i].X + listFamilyGenerals.NumericDashInset
		if got != want {
			t.Errorf("第 %d 個數字欄的破折號左緣 = %d，want %d", i, got, want)
		}
	}
	// 數值靠右對齊的右緣也一起釘住——它與破折號**不是同一個位置**。
	for i, want := range []int{128, 168, 208} {
		if got := listBodyX() + numeric[i].X + numeric[i].W; got != want {
			t.Errorf("第 %d 個數字欄的數值右緣 = %d，want %d", i, got, want)
		}
	}
	// 軍團一覽沒有這個位移：兩者重合。
	if listFamilyCorps.NumericDashInset != 0 {
		t.Error("軍團一覽的破折號是貼欄起點的（parity-menus7 m1），不該有位移")
	}
}
