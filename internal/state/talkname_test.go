package state_test

import (
	"os"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// TestTalkNameUsesAlias 釘住 docs/spec/119：`\1`／`\4` 代入的是**呼び名**
// （武將記錄 `+0x08`），不是姓名（`+0x02`）。
//
// ⭐ 這一支同時斷言兩件事，缺一不可：
//   - 那四個人的 `TalkName()` 是呼び名（改動有效）
//   - **其餘每一個武將兩欄相同**（改動的影響面就是這四個人，不是全庫）
//
// 少了第二條，一支「永遠回 Alias」與一支「永遠回 Name」都可能矇混過關——
// 因為 508 筆裡有 504 筆兩者本來就一樣。
func TestTalkNameUsesAlias(t *testing.T) {
	const scen = "../../workplace/orig/dosv/SINARIO.DAT"
	if _, err := os.Stat(scen); err != nil {
		t.Skipf("找不到 %s", scen)
	}
	// 松崗版是 Big5，比對前先解碼（`text.Decode` 會把定長補的全形空白去掉）。
	want := map[string]string{
		"諸葛亮": "孔明",
		"司馬懿": "仲達",
		"龐統":  "鳳雛",
		// ⚠ 劇本一的李暹：呼び名整格是全形空白，解出來是空字串——
		// **原版就是畫一片空白**，不是資料缺漏。
		"李暹": "",
	}
	seen := map[string]bool{}
	for s := 0; s < 4; s++ {
		w, err := state.LoadScenario(scen, s)
		if err != nil {
			t.Fatalf("劇本 %d: %v", s, err)
		}
		for i := range w.Generals {
			g := &w.Generals[i]
			if g.Name == "" {
				continue
			}
			name := text.Decode([]byte(g.Name), text.Big5)
			got := text.Decode([]byte(g.TalkName()), text.Big5)
			if exp, ok := want[name]; ok {
				seen[name] = true
				if got != exp {
					t.Errorf("劇本 %d 的 %q：TalkName() = %q，want %q",
						s, name, got, exp)
				}
				continue
			}
			// 其餘武將兩欄必須相同——改動只碰得到上面那四個人。
			if g.Alias != "" && g.Alias != g.Name {
				t.Errorf("劇本 %d 武將 %d（%q）的呼び名 %q 不在預期清單裡——"+
					"影響面比 docs/spec/119 §2 寫的大",
					s, i, name, text.Decode([]byte(g.Alias), text.Big5))
			}
		}
	}
	for name := range want {
		if !seen[name] {
			t.Errorf("四個劇本裡一次都沒找到 %q——這一支在驗空氣", name)
		}
	}
}

// 正對照：`Alias` 空白時退回 `Name`（remake 差異，docs/spec/119 §3）。
func TestTalkNameFallsBackToName(t *testing.T) {
	g := state.General{Name: "曹操　"}
	if got := g.TalkName(); got != "曹操　" {
		t.Errorf("沒有呼び名時 TalkName() = %q，want 姓名", got)
	}
	g.Alias = "孟德　"
	if got := g.TalkName(); got != "孟德　" {
		t.Errorf("有呼び名時 TalkName() = %q，want 呼び名", got)
	}
}
