package phone

import (
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/text"

	"github.com/wicanr2/wolong_cht/internal/ui/uitext"
)

// roster 是整份武將名冊串起來，當語系有沒有換的探針。
func roster(s *Session) string {
	var b strings.Builder
	for i := range s.world.Generals {
		b.WriteString(s.Localise(string(s.world.Generals[i].Name[:])))
		b.WriteByte('|')
	}
	return b.String()
}

// plainTalk 找一則不含變數的訊息當探針。
//
// ⚠ 不能寫死索引：`Lines` 對缺變數的訊息是 fail-closed（回 false），
// 而哪幾則帶變數是原版資料決定的，寫死等於把測試綁在某一則上。
func plainTalk(t *testing.T, s *Session) int {
	t.Helper()
	for i := 0; i < text.MessageCount; i++ {
		if lines, ok := s.lib.Talk.Lines(i, nil); ok && len(lines) > 0 && lines[0] != "" {
			return i
		}
	}
	t.Fatal("整份 TALK 找不到一則不含變數的訊息")
	return -1
}

// talkLine 取一則訊息的第一行。
func talkLine(t *testing.T, s *Session, index int) string {
	t.Helper()
	lines, ok := s.lib.Talk.Lines(index, nil)
	if !ok || len(lines) == 0 {
		t.Fatalf("TALK #%d 取不到（ok=%v）", index, ok)
	}
	return lines[0]
}

// 手機沒有命令列，所以「換語言」這條路只有系統面板一個入口
// （docs/spec/86 §4）。這一支釘住換過去之後**三個出口都跟著換**。
func TestSetLanguageSwitchesTalkAndNames(t *testing.T) {
	s := newTestSession(t)

	// ⚠ 探針不能只看一個名字：日文版有 271/343 個名字與繁中**寫法相同**
	//（同樣是漢字），挑到其中一個就會誤判成「沒換」。整份名冊才有鑑別力。
	zhRoster := roster(s)
	rawName := string(s.world.Generals[0].Name[:])
	zhName := s.Localise(rawName)
	if zhName == "" {
		t.Fatal("母本的武將名是空的")
	}
	probe := plainTalk(t, s)
	zhTalk := talkLine(t, s, probe)

	for _, lang := range []uitext.Language{uitext.ZhHans, uitext.Ja, uitext.En} {
		if err := s.SetLanguage(lang, ""); err != nil {
			t.Fatalf("%s：%v", lang, err)
		}
		if got := s.Language(); got != lang {
			t.Fatalf("%s：Language() ＝ %q", lang, got)
		}
		if s.LangPack() == nil {
			t.Fatalf("%s：LangPack() 是 nil，呈現層掛不上字型與詞表", lang)
		}
		if got := roster(s); got == zhRoster {
			t.Errorf("%s：整份武將名冊與母本一字不差", lang)
		}
		if got := talkLine(t, s, probe); got == zhTalk {
			t.Errorf("%s：TALK 探針沒換（仍是 %q）", lang, got)
		}
	}

	// 英文的人名要是羅馬拼音，不是還沒轉的漢字。
	if err := s.SetLanguage(uitext.En, ""); err != nil {
		t.Fatal(err)
	}
	if got := s.Localise(rawName); strings.ContainsFunc(got, func(r rune) bool { return r > 0x2000 }) {
		t.Errorf("英文版的武將名仍含全形字：%q", got)
	}

	// ⭐ 切回母本要**完全還原**：語系層殘留會讓「換回來」與「沒換過」
	// 長得不一樣，而那種差異在畫面上幾乎看不出來。
	if err := s.SetLanguage(uitext.ZhHant, ""); err != nil {
		t.Fatal(err)
	}
	if got := roster(s); got != zhRoster {
		t.Error("切回繁中後武將名冊與母本不同")
	}
	if got := talkLine(t, s, probe); got != zhTalk {
		t.Errorf("切回繁中後 TALK 探針 ＝ %q，母本是 %q", got, zhTalk)
	}
}

// 系統面板的「語言」頁要列得出四個語系，而且目前這一個要有記號——
// 換過去之後畫面全是那個語言，沒有記號就認不出自己選了什麼。
func TestLanguageSheetRowsMarkCurrent(t *testing.T) {
	s := newTestSession(t)
	rows := s.languageRows()
	if len(rows) != len(LanguageChoices) {
		t.Fatalf("語言頁列數 ＝ %d，語系有 %d 個", len(rows), len(LanguageChoices))
	}
	marked := 0
	for i, r := range rows {
		if r.name != LanguageChoices[i].Name {
			t.Errorf("第 %d 列 ＝ %q，預期 %q", i, r.name, LanguageChoices[i].Name)
		}
		if len(r.cols) > 0 && r.cols[0] != "" {
			marked++
			if LanguageChoices[i].Lang != s.Language() {
				t.Errorf("第 %d 列打了勾，但它是 %s 不是目前的 %s", i, LanguageChoices[i].Lang, s.Language())
			}
		}
	}
	if marked != 1 {
		t.Errorf("打勾的列有 %d 個，應該剛好 1 個", marked)
	}
}
