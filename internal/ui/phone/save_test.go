package phone

import (
	"os"
	"path/filepath"
	"testing"
)

// 手機版的存讀檔要走 overlay，**原版素材一個 byte 都不能動**。
func TestSaveWritesOverlayAndLeavesTheOriginalAlone(t *testing.T) {
	s := newSessionInTempRoot(t)
	src := s.sourceSave()
	before, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}

	if err := s.SaveSlot(0); err != nil {
		t.Fatalf("存檔失敗：%v", err)
	}
	after, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("存檔改到了原版 SAVE.DAT")
	}
	if _, err := os.Stat(s.overlaySave()); err != nil {
		t.Fatalf("overlay 沒有被建出來：%v", err)
	}
}

// 存了再讀要拿回同一個世界。
//
// ⚠ 判準是**原版格式那 22,208 B 一模一樣**加上執行期游標回到原位，
// 不是指紋相同：原生存檔不帶亂數相位與災害的執行期覆蓋層
//（`docs/spec/20` §2.4 列了它帶哪些欄位），那幾項讀檔後本來就會不同。
func TestSaveLoadRoundTripRestoresTheBlockAndCursors(t *testing.T) {
	s := newSessionInTempRoot(t)
	// 先跑一段，讓世界離開開局狀態；不然「讀回來一樣」可能只是因為沒動過。
	for i := 0; i < 300; i++ {
		s.Tick()
	}
	want := append([]byte(nil), s.World().Bytes()...)
	wantSnap := s.World().TakeSnapshot()

	if err := s.SaveSlot(1); err != nil {
		t.Fatalf("存檔失敗：%v", err)
	}
	for i := 0; i < 300; i++ {
		s.Tick()
	}
	if string(s.World().Bytes()) == string(want) {
		t.Fatal("又跑了 300 tick 區塊卻沒變，這個測試證明不了東西")
	}
	if err := s.LoadSlot(1); err != nil {
		t.Fatalf("讀檔失敗：%v", err)
	}
	if got := s.World().Bytes(); string(got) != string(want) {
		t.Fatalf("讀回來的區塊與存的時候不同（%d vs %d byte）", len(got), len(want))
	}
	if got := s.World().TakeSnapshot(); got.CityCursor != wantSnap.CityCursor ||
		got.EventCursor != wantSnap.EventCursor {
		t.Fatalf("執行期游標沒有回到原位：%+v，預期 %+v", got, wantSnap)
	}
}

// 空槽要說「空」，寫過的槽要說出日期。
func TestSlotLabelDistinguishesEmptyFromWritten(t *testing.T) {
	s := newSessionInTempRoot(t)
	if got := s.slotLabel(2); got != "空" {
		t.Fatalf("沒寫過的槽標成 %q，預期「空」", got)
	}
	if err := s.SaveSlot(2); err != nil {
		t.Fatal(err)
	}
	if got := s.slotLabel(2); got == "空" || got == "壞檔" {
		t.Fatalf("寫過的槽標成 %q", got)
	}
}

// newSessionInTempRoot 把原版素材連結到一個可寫的暫存根目錄底下。
//
// ⚠ 不能直接拿 `workplace/orig/dosv` 當根：存檔會寫到它的**兄弟目錄**，
// 而那個目錄在 repo 裡（`workplace/orig/save`）。測試不該在那裡留東西。
func newSessionInTempRoot(t *testing.T) *Session {
	t.Helper()
	if _, err := os.Stat(origDir + "/SINARIO.DAT"); err != nil {
		t.Skip("找不到原版素材，跳過")
	}
	root := t.TempDir()
	orig := filepath.Join(root, "orig")
	if err := os.MkdirAll(orig, 0o755); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(origDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		abs, err := filepath.Abs(filepath.Join(origDir, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		// 連結而不是複製：4.4 MB 的素材每個測試複製一次很浪費，
		// 而且**連結讀不到就是路徑錯了**，比複製更早暴露問題。
		if err := os.Symlink(abs, filepath.Join(orig, e.Name())); err != nil {
			t.Fatal(err)
		}
	}
	s, err := NewSession(Options{OrigDir: orig, Scenario: 0, Player: 0, Seed: 7})
	if err != nil {
		t.Fatal(err)
	}
	return s
}
