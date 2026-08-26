package langpack

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/ui/uitext"
)

// 四個語系都要載得起來，而且內容是對的。
//
// 這一支同時驗**內嵌**那條路：測試不從 repo 根目錄跑，
// SearchPaths 的相對路徑一定落空，載得到就是內嵌生效了。
func TestLoadAllLanguages(t *testing.T) {
	for _, tc := range []struct {
		lang     uitext.Language
		wantTalk bool
		name     string // 人名轉換的抽樣
		want     string
	}{
		{uitext.ZhHant, false, "曹操", "曹操"},
		{uitext.ZhHans, true, "曹操", "曹操"}, // 簡體靠字級表，這兩個字同形
		{uitext.Ja, true, "孫權", "孫権"},
		{uitext.En, true, "曹操", "CAO-CAO"},
	} {
		p, err := Load(tc.lang, "")
		if err != nil {
			t.Fatalf("%s: %v", tc.lang, err)
		}
		if (p.Talk != nil) != tc.wantTalk {
			t.Errorf("%s: Talk 存在 = %v，預期 %v", tc.lang, p.Talk != nil, tc.wantTalk)
		}
		if got := p.Convert(tc.name); got != tc.want {
			t.Errorf("%s: %s → %q，預期 %q", tc.lang, tc.name, got, tc.want)
		}
		if !Available(tc.lang) {
			t.Errorf("%s: Available 回 false", tc.lang)
		}
	}
}

// 簡體的字級轉換要真的生效（同形字看不出來，換一個會變的字）。
func TestSimplifiedConvertsCharacters(t *testing.T) {
	p, err := Load(uitext.ZhHans, "")
	if err != nil {
		t.Fatal(err)
	}
	if got := p.Convert("劉備"); got != "刘备" {
		t.Fatalf("劉備 → %q，預期 刘备", got)
	}
}

// F9 的循環順序：繞一圈回到母本。
func TestNextCyclesThroughAvailable(t *testing.T) {
	seen := map[uitext.Language]bool{}
	cur := uitext.ZhHant
	for i := 0; i < 4; i++ {
		cur = Next(cur)
		seen[cur] = true
	}
	if cur != uitext.ZhHant {
		t.Errorf("四步之後停在 %s，預期回到母本", cur)
	}
	for _, l := range []uitext.Language{uitext.ZhHans, uitext.Ja, uitext.En} {
		if !seen[l] {
			t.Errorf("循環漏掉 %s", l)
		}
	}
}
