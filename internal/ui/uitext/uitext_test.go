package uitext

import (
	"os"
	"path/filepath"
	"testing"
)

func writeJSON(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

// 查找順序：覆寫詞表 → 字級表 → 原文（docs/spec/84 §2）。
func TestConvertPriority(t *testing.T) {
	over := writeJSON(t, "ui.json", `{"乾坤":"乾坤（覆寫）"}`)
	chars := writeJSON(t, "chars.json", `{"乾":"干","國":"国"}`)
	tb, err := Load(ZhHans, chars, over)
	if err != nil {
		t.Fatal(err)
	}
	if got := tb.Convert("乾坤"); got != "乾坤（覆寫）" {
		t.Fatalf("覆寫詞表沒有優先：%q", got)
	}
	if got := tb.Convert("國乾"); got != "国干" {
		t.Fatalf("字級表轉換錯：%q", got)
	}
	if got := tb.Convert("兵"); got != "兵" {
		t.Fatalf("表外字應原樣：%q", got)
	}
}

// 母本語系與 nil 表都是恆等——缺檔不能讓文字消失。
func TestZhHantAndNilAreIdentity(t *testing.T) {
	tb, err := Load(ZhHant, "")
	if err != nil {
		t.Fatal(err)
	}
	if got := tb.Convert("信賴度"); got != "信賴度" {
		t.Fatalf("母本語系被改寫：%q", got)
	}
	var nilTable *Table
	if got := nilTable.Convert("信賴度"); got != "信賴度" {
		t.Fatalf("nil 表被改寫：%q", got)
	}
	if nilTable.RuneMap() != nil {
		t.Fatal("nil 表的 RuneMap 應為 nil")
	}
}

func TestParseLanguage(t *testing.T) {
	if l, err := ParseLanguage("ZH-Hans"); err != nil || l != ZhHans {
		t.Fatalf("大小寫不敏感解析失敗：%v %v", l, err)
	}
	if _, err := ParseLanguage("jp"); err == nil {
		t.Fatal("不支援的語系應報錯")
	}
}

// 真實字級表（版控內的產出）要載得起來，且把常用繁體字轉掉。
func TestRealCharTableLoads(t *testing.T) {
	p := "../../../translations/t2s-chars.json"
	if _, err := os.Stat(p); err != nil {
		t.Skip("找不到 t2s-chars.json")
	}
	tb, err := Load(ZhHans, p)
	if err != nil {
		t.Fatal(err)
	}
	if got := tb.Convert("信賴度"); got != "信赖度" {
		t.Fatalf("信賴度 → %q，預期 信赖度", got)
	}
	if got := tb.Convert("孫策攻擊劉繇"); got != "孙策攻击刘繇" {
		t.Fatalf("人名轉換錯：%q", got)
	}
}

// UI 的數值欄位是先 Sprintf 才畫的，畫到語系層時數字已經填進去了。
// 樣板查找要能把它換回去（docs/spec/84 §2）。
func TestNumericTemplateLookup(t *testing.T) {
	over := writeJSON(t, "ui.json", `{"兵 %d":"Troops %d","士氣 %d／%d":"Morale %d/%d"}`)
	tb, err := Load(En, "", over)
	if err != nil {
		t.Fatal(err)
	}
	if got := tb.Convert("兵 300"); got != "Troops 300" {
		t.Fatalf("單一數值樣板 = %q", got)
	}
	if got := tb.Convert("士氣 80／120"); got != "Morale 80/120" {
		t.Fatalf("兩個數值要依序填回 = %q", got)
	}
	// 樣板查不到就原樣——寧可顯示原文，也不要把數字擺錯位置。
	if got := tb.Convert("金 300"); got != "金 300" {
		t.Fatalf("查不到的樣板被改了 = %q", got)
	}
}

// 名表包了一層 {"note":..., "names":{...}}，要載得進來。
func TestNestedNamesTable(t *testing.T) {
	p := writeJSON(t, "names.json", `{"note":"說明","names":{"曹操":"CAO-CAO"}}`)
	tb, err := Load(En, "", p)
	if err != nil {
		t.Fatal(err)
	}
	if got := tb.Convert("曹操"); got != "CAO-CAO" {
		t.Fatalf("名表沒載進來 = %q", got)
	}
}

// 三個語系的名表都要載得起來，而且抽樣的名字要對（docs/spec/84 §4）。
func TestRealNameTables(t *testing.T) {
	for _, tc := range []struct {
		lang Language
		file string
		want map[string]string
	}{
		{Ja, "names-ja.json", map[string]string{
			"曹操": "曹操", "孫權": "孫権", "郭汜": "郭汜", "廣武": "広武"}},
		{En, "names-en.json", map[string]string{
			"曹操": "CAO-CAO", "諸葛亮": "ZHUGE-LIANG", "劉禪": "LIU-SHAN",
			"許昌": "XUCHANG", "長安": "CHANG'AN", "會稽": "KUAIJI"}},
	} {
		p := "../../../translations/" + tc.file
		if _, err := os.Stat(p); err != nil {
			t.Skipf("找不到 %s", tc.file)
		}
		tb, err := Load(tc.lang, "", p)
		if err != nil {
			t.Fatalf("%s: %v", tc.file, err)
		}
		for zh, want := range tc.want {
			if got := tb.Convert(zh); got != want {
				t.Errorf("%s：%s → %q，預期 %q", tc.file, zh, got, want)
			}
		}
	}
}
