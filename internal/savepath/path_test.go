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
