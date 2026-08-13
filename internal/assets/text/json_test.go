package text

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadJSONPreservesMarkersAndTrailingBlankLines(t *testing.T) {
	messages := make([][]string, MessageCount)
	for i := range messages {
		messages[i] = []string{""}
	}
	messages[7] = []string{"{1}外交費{6}{7}。", ""}
	data, err := json.Marshal(jsonTable{Encoding: "cp950", Messages: messages})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "talk.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	table, err := LoadJSON(path, Big5)
	if err != nil {
		t.Fatal(err)
	}
	got := table.Messages[7]
	if len(got.Lines) != 2 || len(got.Lines[1].Parts) != 0 {
		t.Fatalf("尾端空行未保留：%#v", got.Lines)
	}
	wantMarkers := []byte{'1', '6', '7'}
	markers := got.Markers()
	if string(markers) != string(wantMarkers) {
		t.Fatalf("marker 錯誤：%q，預期 %q", markers, wantMarkers)
	}
	if got.Lines[0].Parts[0].Marker != '1' || got.Lines[0].Parts[2].Marker != '6' ||
		got.Lines[0].Parts[3].Marker != '7' {
		t.Fatalf("marker 分段錯誤：%#v", got.Lines[0].Parts)
	}
}

func TestLoadJSONRejectsWrongMessageCount(t *testing.T) {
	data, err := json.Marshal(jsonTable{Encoding: "cp950", Messages: [][]string{{"只有一則"}}})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "short.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadJSON(path, Big5); err == nil {
		t.Fatal("訊息數量錯誤卻成功載入")
	}
}
