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
	tb, err := Load(ZhHans, over, chars)
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
	tb, err := Load(ZhHant, "", "")
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
	tb, err := Load(ZhHans, "", p)
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
