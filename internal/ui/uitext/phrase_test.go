package uitext

import "testing"

// 片語層只動明列的那幾條（docs/spec/87 §4）。這一支釘住三件事：
// 命中要換、最長優先、**沒有片語的字串一個 byte 都不動**。
func TestPhraseLayer(t *testing.T) {
	raw := []byte(`{
	  "整句": "whole",
	  "phrases": {" 對 ": " vs ", " 對我方宣戰": " declares war on us"}
	}`)
	tab, err := Parse(En, nil, raw)
	if err != nil {
		t.Fatal(err)
	}

	// 逐句仍然優先。
	if got := tab.Convert("整句"); got != "whole" {
		t.Errorf("逐句表 ＝ %q", got)
	}
	// 拼出來的句子靠片語補接縫。
	if got := tab.Convert("YUAN-YIN 對 CAO-CAO"); got != "YUAN-YIN vs CAO-CAO" {
		t.Errorf("片語沒換：%q", got)
	}
	// ⭐ 最長優先：先換短的會把長的那一條切成「 vs 我方宣戰」。
	if got := tab.Convert("CAO-CAO 對我方宣戰"); got != "CAO-CAO declares war on us" {
		t.Errorf("最長優先失效：%q", got)
	}
	// 沒有片語就原樣回傳——這一層的風險全在「換到不該換的地方」。
	for _, s := range []string{"CAO-CAO", "無關的一句話", "對", "面對面"} {
		if got := tab.Convert(s); got != s {
			t.Errorf("%q 被動到了：%q", s, got)
		}
	}
}

// zh-Hans 的片語要在字級表**之前**做：字形換過之後，繁體的片語鍵就對不上。
func TestPhraseBeforeCharTable(t *testing.T) {
	tab, err := Parse(ZhHans,
		[]byte(`{"對":"对","國":"国"}`),
		[]byte(`{"phrases":{" 對我方宣戰":" 向我方宣战"}}`))
	if err != nil {
		t.Fatal(err)
	}
	if got := tab.Convert("曹操 對我方宣戰"); got != "曹操 向我方宣战" {
		t.Errorf("順序錯了：%q", got)
	}
	// 沒中片語的部分仍走字級表。
	if got := tab.Convert("國"); got != "国" {
		t.Errorf("字級表失效：%q", got)
	}
}
