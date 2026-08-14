package savefile

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

const scenario = "../../workplace/orig/dosv/SINARIO.DAT"

func load(t *testing.T) *state.World {
	t.Helper()
	if _, err := os.Stat(scenario); err != nil {
		t.Skip("找不到原版 SINARIO.DAT，跳過")
	}
	w, err := state.LoadScenario(scenario, 0)
	if err != nil {
		t.Fatal(err)
	}
	return w
}

// 原版 → 原生 → 匯出，必須與直接走原版路徑的輸出 byte-for-byte 相同。
// **這是這個格式唯一的硬約束**（docs/spec/20 §2.1）。
func TestRoundTripToOriginalIsByteForByte(t *testing.T) {
	w := load(t)
	want := ExportOriginal(w)

	data, err := Encode(w, Origin{Source: "SINARIO.DAT", Block: 0})
	if err != nil {
		t.Fatal(err)
	}
	back, _, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyExport(back, want); err != nil {
		t.Fatal(err)
	}
}

// 原版格式裝不下的執行期狀態要活過存檔。改之前這一條是紅的——
// routes 與 cityCursor 讀檔後會歸零。
func TestRuntimeStateSurvivesTheRoundTrip(t *testing.T) {
	w := load(t)
	w.EnableStrategicAI()
	snap := w.TakeSnapshot()
	snap.CityCursor = 137
	snap.EventCursor = 42
	snap.Routes = map[int][][2]int{3: {{10, 20}, {11, 21}}}
	if err := w.Restore(snap); err != nil {
		t.Fatal(err)
	}

	data, err := Encode(w, Origin{Source: "SAVE.DAT", Block: 1})
	if err != nil {
		t.Fatal(err)
	}
	back, f, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	got := back.TakeSnapshot()
	if got.CityCursor != 137 || got.EventCursor != 42 {
		t.Fatalf("游標 = 據點 %d／事件 %d，want 137/42", got.CityCursor, got.EventCursor)
	}
	if len(got.Routes[3]) != 2 || got.Routes[3][1] != [2]int{11, 21} {
		t.Fatalf("行軍路線 = %v，want [[10 20] [11 21]]", got.Routes[3])
	}
	if !got.StrategicAI {
		t.Error("政略 AI 的開關沒有活過存檔")
	}
	if f.Origin.Block != 1 || f.Origin.Source != "SAVE.DAT" {
		t.Errorf("來源紀錄 = %+v", f.Origin)
	}
}

// 檔案被改壞要大聲失敗，不能靜靜採用。
func TestDecodeRejectsTamperedFile(t *testing.T) {
	w := load(t)
	data, err := Encode(w, Origin{Source: "SINARIO.DAT"})
	if err != nil {
		t.Fatal(err)
	}
	var f File
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatal(err)
	}

	// ① 動了 raw，雜湊就對不上。
	tampered := f
	tampered.Raw = append([]byte(nil), f.Raw...)
	tampered.Raw[0x1234] ^= 0xFF
	if b, _ := json.Marshal(tampered); func() bool { _, _, err := Decode(b); return err == nil }() {
		t.Error("改過 raw 的檔案被接受了")
	}

	// ② 版本不合。
	bad := f
	bad.Format = "wolong-save/999"
	if b, _ := json.Marshal(bad); func() bool { _, _, err := Decode(b); return err == nil }() {
		t.Error("未知格式版本被接受了")
	}

	// ③ Runtime 索引超界——夾住會產生「能玩但行為不對」的世界。
	oob := f
	oob.Runtime.CityCursor = 9999
	if b, _ := json.Marshal(oob); func() bool { _, _, err := Decode(b); return err == nil }() {
		t.Error("超界的據點游標被接受了")
	}
}

// 沒有原版檔案也要載得動（決策 C）：Decode 只吃 bytes，不碰檔案系統。
func TestDecodeDoesNotNeedTheOriginalFiles(t *testing.T) {
	w := load(t)
	data, err := Encode(w, Origin{Source: "SINARIO.DAT"})
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := Decode(data); err != nil {
		t.Fatalf("Decode 應該只靠檔案自己的 raw：%v", err)
	}
}
