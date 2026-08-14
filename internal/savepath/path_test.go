package savepath

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitialLoadPath(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "SINARIO.DAT")
	overlay := filepath.Join(dir, "save", "SAVE.DAT")
	if err := os.WriteFile(source, []byte("source"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := InitialLoadPath(source, overlay)
	if err != nil {
		t.Fatal(err)
	}
	if got != source {
		t.Fatalf("missing overlay path = %q, want source %q", got, source)
	}

	if err := os.MkdirAll(filepath.Dir(overlay), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(overlay, []byte("save"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err = InitialLoadPath(source, overlay)
	if err != nil {
		t.Fatal(err)
	}
	if got != overlay {
		t.Fatalf("existing overlay path = %q, want %q", got, overlay)
	}
}

func TestInitialLoadPathRejectsOriginalAndDirectory(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "SINARIO.DAT")
	if err := os.WriteFile(source, []byte("source"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := InitialLoadPath(source, source); err == nil {
		t.Fatal("same source and overlay should be rejected")
	}

	overlayDir := filepath.Join(dir, "save")
	if err := os.Mkdir(overlayDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := InitialLoadPath(source, overlayDir); err == nil {
		t.Fatal("directory overlay should be rejected")
	}
}

func TestNativePath(t *testing.T) {
	got, err := NativePath("/tmp/x/SAVE.DAT", 0)
	if err != nil || got != "/tmp/x/SAVE-slot1.wlsave" {
		t.Fatalf("NativePath = %q, %v；want /tmp/x/SAVE-slot1.wlsave", got, err)
	}
	if got, _ := NativePath("/tmp/x/SAVE.DAT", 3); got != "/tmp/x/SAVE-slot4.wlsave" {
		t.Fatalf("第四槽 = %q", got)
	}
	// 原生檔與原版 overlay 一定是不同檔案，否則其中一種會被覆蓋。
	if got, _ := NativePath("/tmp/x/SAVE.DAT", 0); SamePath(got, "/tmp/x/SAVE.DAT") {
		t.Fatal("原生存檔路徑與原版 overlay 相同")
	}
	if _, err := NativePath("", 0); err == nil {
		t.Error("沒有 overlay 路徑時應該報錯")
	}
	if _, err := NativePath("/tmp/x/SAVE.DAT", -1); err == nil {
		t.Error("負數槽位應該報錯")
	}
}
