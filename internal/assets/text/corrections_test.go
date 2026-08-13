package text

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApplyCorrectionsAppliesOnlyVerifiedPublicOverlay(t *testing.T) {
	messages := make([][]string, MessageCount)
	for i := range messages {
		messages[i] = []string{""}
	}
	messages[7] = []string{"原文{1}"}
	raw := tableFromLines(t, messages).Bytes()
	table, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	fix := "修正{1}"
	manifest, err := json.Marshal(correctionManifest{Corrections: []correctionItem{{
		ID: 7, CHT: "原文{1}", Fix: &fix,
	}}})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "corrections.json")
	if err := os.WriteFile(path, manifest, 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := ApplyCorrections(table, path, Big5)
	if err != nil {
		t.Fatal(err)
	}
	text, err := messageText(got.Messages[7], Big5)
	if err != nil {
		t.Fatal(err)
	}
	if text != "修正{1}" {
		t.Fatalf("校訂後文字 = %q，預期修正{1}", text)
	}
	original, err := messageText(table.Messages[7], Big5)
	if err != nil {
		t.Fatal(err)
	}
	if original != "原文{1}" {
		t.Fatalf("輸入表被改寫：%q", original)
	}
}

func TestApplyCorrectionsFailsClosedOnUnexpectedSource(t *testing.T) {
	messages := make([][]string, MessageCount)
	for i := range messages {
		messages[i] = []string{""}
	}
	messages[7] = []string{"不是預期原文"}
	raw := tableFromLines(t, messages).Bytes()
	table, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	fix := "修正"
	manifest, err := json.Marshal(correctionManifest{Corrections: []correctionItem{{
		ID: 7, CHT: "原文", Fix: &fix,
	}}})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "corrections.json")
	if err := os.WriteFile(path, manifest, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ApplyCorrections(table, path, Big5); err == nil {
		t.Fatal("來源文字不符仍成功套用校訂")
	}
}

func TestApplyCorrectionsExpandsSameTextReference(t *testing.T) {
	messages := make([][]string, MessageCount)
	for i := range messages {
		messages[i] = []string{""}
	}
	messages[7] = []string{"原文"}
	messages[8] = []string{"原文"}
	table, err := Parse(tableFromLines(t, messages).Bytes())
	if err != nil {
		t.Fatal(err)
	}
	fix := "修正"
	manifest, err := json.Marshal(correctionManifest{Corrections: []correctionItem{
		{ID: 7, CHT: "原文", Fix: &fix},
		{ID: 8, CHT: "同 7", Fix: &fix},
	}})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "corrections.json")
	if err := os.WriteFile(path, manifest, 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := ApplyCorrections(table, path, Big5)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []int{7, 8} {
		text, err := messageText(got.Messages[id], Big5)
		if err != nil {
			t.Fatal(err)
		}
		if text != "修正" {
			t.Fatalf("#%d = %q，預期修正", id, text)
		}
	}
}

func TestApplyCorrectionsPreservesValidatedWrapAndTrailingLines(t *testing.T) {
	messages := make([][]string, MessageCount)
	for i := range messages {
		messages[i] = []string{""}
	}
	original := strings.Repeat("原", 10)
	messages[7] = []string{original, "", ""}
	table, err := Parse(tableFromLines(t, messages).Bytes())
	if err != nil {
		t.Fatal(err)
	}
	fix := strings.Repeat("甲", 12)
	manifest, err := json.Marshal(correctionManifest{Corrections: []correctionItem{{
		ID: 7, CHT: original, Fix: &fix,
	}}})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "corrections.json")
	if err := os.WriteFile(path, manifest, 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := ApplyCorrections(table, path, Big5)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Messages[7].Lines) != 4 {
		t.Fatalf("行數 = %d，預期 4（兩行加兩個尾端空行）", len(got.Messages[7].Lines))
	}
	for i, want := range []string{strings.Repeat("甲", 10), strings.Repeat("甲", 2), "", ""} {
		actual, err := messageText(Message{Lines: []Line{got.Messages[7].Lines[i]}}, Big5)
		if err != nil {
			t.Fatal(err)
		}
		if actual != want {
			t.Fatalf("第 %d 行 = %q，預期 %q", i, actual, want)
		}
	}
}

// TestApplyCorrectionsMatchesGeneratedTableWhenOriginalAvailable 把公開的最小
// 覆蓋表與既有 1,022 則產生表逐則比對。完整表僅作本機／CI 的 oracle，不能放進
// 發行包；少了合法原版測試輸入時則跳過，而不是以替代資料宣稱相同。
func TestApplyCorrectionsMatchesGeneratedTableWhenOriginalAvailable(t *testing.T) {
	originalPath := filepath.Join("..", "..", "..", "workplace", "orig", "dosv", "TALK.DAT")
	original, err := os.ReadFile(originalPath)
	if errors.Is(err, os.ErrNotExist) {
		t.Skipf("未提供松崗 TALK.DAT：%s", originalPath)
	}
	if err != nil {
		t.Fatal(err)
	}
	raw, err := Parse(original)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ApplyCorrections(raw, filepath.Join("..", "..", "..", "translations", "corrections.json"), Big5)
	if err != nil {
		t.Fatal(err)
	}
	want, err := LoadJSON(filepath.Join("..", "..", "..", "translations", "talk-dosv-corrected.json"), Big5)
	if err != nil {
		t.Fatal(err)
	}
	for id := range got.Messages {
		gotBytes := got.Messages[id].Bytes()
		wantBytes := want.Messages[id].Bytes()
		if !bytes.Equal(gotBytes, wantBytes) {
			t.Fatalf("#%d 套用公開校訂後與既有產生表不同：得到 %q，預期 %q", id,
				messageLines(got.Messages[id]), messageLines(want.Messages[id]))
		}
	}
}

func messageLines(msg Message) []string {
	lines := make([]string, len(msg.Lines))
	for i, line := range msg.Lines {
		lines[i], _ = messageText(Message{Lines: []Line{line}}, Big5)
	}
	return lines
}

func tableFromLines(t *testing.T, messages [][]string) *Table {
	t.Helper()
	table := &Table{}
	for i, lines := range messages {
		message, err := messageFromJSONLines(lines, Big5)
		if err != nil {
			t.Fatalf("第 %d 則測試文字無法轉換：%v", i, err)
		}
		table.Messages[i] = message
	}
	return table
}
