package text

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// origPath 回傳原版素材的路徑。素材不進版控，沒有就跳過測試
// ——CI 上不該因為缺原版檔而紅。
func origPath(t *testing.T, ver, name string) string {
	t.Helper()
	p := filepath.Join("..", "..", "..", "workplace", "orig", ver, name)
	if _, err := os.Stat(p); err != nil {
		t.Skipf("找不到原版素材 %s，跳過", p)
	}
	return p
}

// TestRoundTrip 是本套件唯一有意義的驗收：解出來再組回去必須
// byte-for-byte 相同。寫不回去的中文化工具是不能用的。
func TestRoundTrip(t *testing.T) {
	for _, ver := range []string{"dosv", "pc98"} {
		t.Run(ver, func(t *testing.T) {
			raw, err := os.ReadFile(origPath(t, ver, "TALK.DAT"))
			if err != nil {
				t.Fatal(err)
			}
			tbl, err := Parse(raw)
			if err != nil {
				t.Fatal(err)
			}
			got := tbl.Bytes()
			if !bytes.Equal(got, raw) {
				for i := range got {
					if i >= len(raw) || got[i] != raw[i] {
						t.Fatalf("第一個差異 @ 0x%04x（原 %d B，重建 %d B）",
							i, len(raw), len(got))
					}
				}
				t.Fatalf("長度不同：原 %d B，重建 %d B", len(raw), len(got))
			}
		})
	}
}

// TestMarkerSet 釘住 docs/formats/01 §3 的統計：六種標記，沒有 \5。
func TestMarkerSet(t *testing.T) {
	raw, err := os.ReadFile(origPath(t, "dosv", "TALK.DAT"))
	if err != nil {
		t.Fatal(err)
	}
	tbl, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	count := map[byte]int{}
	for _, m := range tbl.Messages {
		for _, mk := range m.Markers() {
			count[mk]++
		}
	}
	want := map[byte]int{'1': 102, '2': 37, '3': 85, '4': 34, '6': 72, '7': 18}
	if len(count) != len(want) {
		t.Fatalf("標記種類 %d 種，預期 %d 種：%v", len(count), len(want), count)
	}
	for k, v := range want {
		if count[k] != v {
			t.Errorf("標記 \\%c 出現 %d 次，預期 %d", k, count[k], v)
		}
	}
	if _, ok := count['5']; ok {
		t.Error("出現了 \\5 —— 原版沒有這個標記")
	}
}

// TestTrailingBlankLinesArePreserved 釘住最容易踩的那個坑：
// 訊息結尾的空行是資料，正規化掉就寫不回去。
func TestTrailingBlankLinesArePreserved(t *testing.T) {
	raw, err := os.ReadFile(origPath(t, "dosv", "TALK.DAT"))
	if err != nil {
		t.Fatal(err)
	}
	tbl, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	var single, blank int
	for _, m := range tbl.Messages {
		if len(m.Lines) == 1 && len(m.Lines[0].Parts) == 0 {
			blank++ // 空訊息：只有一個 NUL
		}
		if n := len(m.Bytes()); n == 1 {
			single++
		}
	}
	if blank != 78 || single != 78 {
		t.Errorf("空訊息 %d 則（單 NUL %d 則），預期各 78 則", blank, single)
	}
}

// 同一個標記出現多次時要依序取值（docs/spec/106）：原版的 formatter 是共用的
// 堆疊游標，「{1}大人的兵馬，遇上{1}的兵馬了」兩個 `{1}` 是兩個不同的武將。
func TestLinesSeqGivesEachRepeatedMarkerItsOwnValue(t *testing.T) {
	tb := &Table{enc: UTF8}
	tb.Messages[7] = Message{Lines: []Line{{Parts: []Part{
		{Marker: '1'}, {Raw: []byte("大人遇上")}, {Marker: '1'}, {Raw: []byte("了")},
	}}}}
	vars := map[byte]string{'1': "甲"}
	got, ok := tb.LinesSeq(7, vars, map[byte][]string{'1': {"甲", "乙"}})
	if !ok || len(got) != 1 || got[0] != "甲大人遇上乙了" {
		t.Fatalf("LinesSeq = %q, %v", got, ok)
	}
	// seq 用完就退回 vars；沒給 seq 時與 Lines 相同。
	plain, ok := tb.Lines(7, vars)
	if !ok || plain[0] != "甲大人遇上甲了" {
		t.Fatalf("Lines = %q, %v", plain, ok)
	}
	short, ok := tb.LinesSeq(7, vars, map[byte][]string{'1': {"丙"}})
	if !ok || short[0] != "丙大人遇上甲了" {
		t.Fatalf("seq 用完沒有退回 vars：%q", short)
	}
}

// ⭐ 共用游標：`{6}` 也吃一個參數，所以兩種版面的 `{1}` 拿到不同的人
// （`docs/spec/136`）。
func TestLinesStreamSharesOneCursorAcrossMarkers(t *testing.T) {
	tb := &Table{enc: UTF8}
	// #7：{1}啊　　#8：{6}我是{1}
	tb.Messages[7] = Message{Lines: []Line{
		{Parts: []Part{{Marker: '1'}, {Raw: []byte("啊")}}},
		{},
	}}
	tb.Messages[8] = Message{Lines: []Line{
		{Parts: []Part{{Marker: '6'}, {Raw: []byte("我是")}, {Marker: '1'}}},
		{},
	}}
	params := []string{"對手", "說話者"}

	got, names, ok := tb.LinesStream(7, params)
	if !ok || len(got) != 1 || got[0] != "對手啊" {
		t.Fatalf("沒有 {6} 時 {1} 應該拿到對手，得到 %v（ok=%v）", got, ok)
	}
	if len(names) != 1 || names[0] != "對手" {
		t.Errorf("要畫色 9 的名字回報成 %v，應為 [對手]", names)
	}
	got, names, ok = tb.LinesStream(8, params)
	if !ok || len(got) != 1 || got[0] != "我是說話者" {
		t.Fatalf("{6} 吃掉一個之後 {1} 應該拿到說話者，得到 %v（ok=%v）", got, ok)
	}
	if len(names) != 1 || names[0] != "說話者" {
		t.Errorf("要畫色 9 的名字回報成 %v，應為 [說話者]", names)
	}
}

// 參數不夠就整則丟棄，不畫半句（fail-closed，同 Lines）。
func TestLinesStreamFailsClosedWhenParamsRunOut(t *testing.T) {
	tb := &Table{enc: UTF8}
	tb.Messages[7] = Message{Lines: []Line{
		{Parts: []Part{{Marker: '6'}, {Marker: '1'}, {Marker: '1'}}},
		{},
	}}
	if got, _, ok := tb.LinesStream(7, []string{"甲", "乙"}); ok {
		t.Fatalf("參數只有兩個卻代出 %v，應該整則丟棄", got)
	}
}
