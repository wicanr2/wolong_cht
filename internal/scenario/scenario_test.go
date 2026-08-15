package scenario

import (
	"encoding/json"
	"os"
	"testing"
	"unicode/utf8"

	"github.com/wicanr2/wolong_cht/internal/state"
)

const origPath = "../../workplace/orig/dosv/SINARIO.DAT"

func load(t *testing.T, block int) *state.World {
	t.Helper()
	if _, err := os.Stat(origPath); err != nil {
		t.Skip("沒有原版 SINARIO.DAT，跳過")
	}
	w, err := state.LoadScenario(origPath, block)
	if err != nil {
		t.Fatalf("載入區塊 %d：%v", block, err)
	}
	return w
}

// ⭐ 驗收標準是 byte-for-byte：匯出成 JSON 再套回一份新載入的 World，
// `Bytes()` 要與原本完全相同。這同時擋住兩件事——欄位漏掉，
// 以及名字在 Big5 ↔ UTF-8 之間來回時被改動。
func TestScenarioRoundTrip(t *testing.T) {
	for block := 0; block < 4; block++ {
		w := load(t, block)
		before := w.Bytes()

		data, err := json.Marshal(FromWorld(w, Meta{Block: block}))
		if err != nil {
			t.Fatalf("區塊 %d 匯出：%v", block, err)
		}
		var s Scenario
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatalf("區塊 %d 讀回：%v", block, err)
		}
		w2 := load(t, block)
		if err := s.ApplyTo(w2); err != nil {
			t.Fatalf("區塊 %d 套用：%v", block, err)
		}
		after := w2.Bytes()
		if len(before) != len(after) {
			t.Fatalf("區塊 %d 長度 %d → %d", block, len(before), len(after))
		}
		for i := range before {
			if before[i] != after[i] {
				t.Fatalf("區塊 %d 在 +0x%X 不一致：%02X → %02X",
					block, i, before[i], after[i])
			}
		}
	}
}

// 匯出的名字要是合法 UTF-8——JSON 才有意義，編輯器才看得懂。
func TestScenarioNamesAreUTF8(t *testing.T) {
	s := FromWorld(load(t, 0), Meta{})
	if !utf8.ValidString(s.Meta.Title) || s.Meta.Title == "" {
		t.Errorf("標題不是合法 UTF-8：%q", s.Meta.Title)
	}
	for i, c := range s.Cities {
		if !utf8.ValidString(c.Name) {
			t.Errorf("據點 %d 的名字不是合法 UTF-8：%q", i, c.Name)
		}
	}
	named := 0
	for i, g := range s.Generals {
		if !utf8.ValidString(g.Name) {
			t.Errorf("武將 %d 的名字不是合法 UTF-8：%q", i, g.Name)
		}
		if g.Alive && g.Name != "" {
			named++
		}
	}
	if named < 100 {
		t.Errorf("只有 %d 個在世武將有名字——解碼八成壞了", named)
	}
}

// 改名字要真的寫得進去，而且**變短時尾巴不能留著舊字**。
func TestScenarioRenameWritesBack(t *testing.T) {
	w := load(t, 0)
	s := FromWorld(w, Meta{})
	// 找一個名字比三個字短的據點來改，才驗得到補空白那條路。
	s.Cities[0].Name = "新"
	w2 := load(t, 0)
	if err := s.ApplyTo(w2); err != nil {
		t.Fatalf("套用：%v", err)
	}
	if got := []byte(w2.Cities[0].Name); len(got) != nameBytes {
		t.Fatalf("改名後的欄位長度 = %d，want %d（要補全形空白）", len(got), nameBytes)
	}
	back := FromWorld(w2, Meta{})
	if back.Cities[0].Name != "新" {
		t.Errorf("改名後再匯出 = %q，want 新", back.Cities[0].Name)
	}
	// 其他據點不受影響。
	if back.Cities[1].Name != s.Cities[1].Name {
		t.Errorf("改一個名字動到了別的：%q → %q", s.Cities[1].Name, back.Cities[1].Name)
	}
}

// 太長的名字要**回錯誤**，不能截斷後靜靜寫進去。
func TestScenarioRejectsOverlongName(t *testing.T) {
	s := FromWorld(load(t, 0), Meta{})
	s.Cities[0].Name = "這個名字太長了"
	if err := s.ApplyTo(load(t, 0)); err == nil {
		t.Fatal("超長的名字卻套用成功了")
	}
}

// 筆數不對要擋下來——編輯器少寫一筆會讓後面整排錯位。
func TestScenarioRejectsWrongCounts(t *testing.T) {
	s := FromWorld(load(t, 0), Meta{})
	s.Cities = s.Cities[:len(s.Cities)-1]
	if err := s.ApplyTo(load(t, 0)); err == nil {
		t.Fatal("據點少一筆卻套用成功了")
	}
}
