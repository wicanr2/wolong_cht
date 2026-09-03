package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/ui/sound"
)

// soundGame 給一個有音檔的播放層。內容不必是真的 ogg——
// 這幾支測試只走選單那一側，不碰解碼。
func soundGame(t *testing.T) *game {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "bgm-0.ogg"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	b := sound.Open(dir)
	b.SetSilent(true)
	return &game{sound: b}
}

// ⚠ 預設要是 TYPE 1：docs/playtest/39 的逐像素對拍靠這一格顯示「TYPE 1」。
func TestSoundValueDefaultsToType1(t *testing.T) {
	g := soundGame(t)
	if got := g.soundValue(); got != "TYPE 1" {
		t.Errorf("預設值 ＝ %q，應為 TYPE 1（docs/spec/122 §3）", got)
	}
	if got := (&game{sound: sound.Open(t.TempDir())}).soundValue(); got != "未接入" {
		t.Errorf("沒有音檔時 ＝ %q，應為 未接入", got)
	}
}

// 原版是環狀遞增（sub_16062），五個選項 ＯＦＦ／TYPE 1–4。
func TestCycleSoundWalksOriginalFiveOptions(t *testing.T) {
	g := soundGame(t)
	want := []string{"TYPE 2", "TYPE 3", "TYPE 4", "ＯＦＦ", "TYPE 1"}
	for i, w := range want {
		g.cycleSound(true)
		if got := g.soundValue(); got != w {
			t.Fatalf("左鍵第 %d 次 ＝ %q，應為 %q", i+1, got, w)
		}
	}
	back := []string{"ＯＦＦ", "TYPE 4", "TYPE 3", "TYPE 2", "TYPE 1"}
	for i, w := range back {
		g.cycleSound(false)
		if got := g.soundValue(); got != w {
			t.Fatalf("右鍵第 %d 次 ＝ %q，應為 %q", i+1, got, w)
		}
	}
}

// 沒有音檔時那一列不做事——「未接入」不是可以被切掉的選項。
func TestCycleSoundIgnoredWithoutAudio(t *testing.T) {
	g := &game{sound: sound.Open(t.TempDir())}
	g.cycleSound(true)
	if got := g.soundValue(); got != "未接入" {
		t.Errorf("沒有音檔時被切成 %q", got)
	}
}
