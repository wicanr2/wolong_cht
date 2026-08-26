package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 半形語系的戰場標題兩行都不准超過標題格（docs/spec/87 §5）。
//
// ⚠ 量的是**最長的那一組**，不是手上這一場：`GONGSUN-ZAN` 與
// `YANGPINGGUAN` 平常碰不到，但碰到時畫面就爛了，而那時沒有人在看測試。
func TestLatinBattleTitleFits(t *testing.T) {
	raw, err := os.ReadFile("../../translations/names-en.json")
	if err != nil {
		t.Skipf("讀不到英文名表：%v", err)
	}
	var doc struct {
		Names map[string]string `json:"names"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	longest := ""
	for _, v := range doc.Names {
		if len(v) > len(longest) {
			longest = v
		}
	}
	if longest == "" {
		t.Fatal("名表是空的")
	}
	for _, line := range []string{longest, "vs " + longest} {
		if w := textdraw.StringWidth(line); w > battleSideTitleW {
			t.Errorf("%q 寬 %d px，標題格只有 %d", line, w, battleSideTitleW)
		}
	}
}

// 將旗上的名字欄：半形語系用 8 個字母，中日文維持原版的三個全形字。
func TestLatinSideNameWidth(t *testing.T) {
	if battleSideNameLatinW/textdraw.HalfW != 8 {
		t.Errorf("半形名字欄 ＝ %d px（%d 個字母），規格是 8 個",
			battleSideNameLatinW, battleSideNameLatinW/textdraw.HalfW)
	}
	// 往左不可以蓋到旗子的 ▶，往右不可以蓋到「軍」。
	if battleSideNameLatinX+battleSideNameLatinW != battleSideNameX+battleSideNameW {
		t.Errorf("右緣對不齊：半形 %d，原版 %d",
			battleSideNameLatinX+battleSideNameLatinW, battleSideNameX+battleSideNameW)
	}
}
