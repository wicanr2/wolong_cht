package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEditedAbilitiesRoundTripIntoRetainedSlotBytes(t *testing.T) {
	raw := make([]byte, blockSize*numBlocks)
	raw[generalBase+0x0E], raw[generalBase+0x0F], raw[generalBase+0x10] = 0x12, 0x37, 0x49
	w := loadBlock(raw[:blockSize])
	before := w.Generals[0].Abilities()
	bad := before
	bad.Tactic, bad.Martial = 8, 10
	if w.ApplyGeneralAbilities(0, bad) || w.Generals[0].Abilities() != before {
		t.Fatal("invalid draft must leave all stats untouched")
	}
	want := GeneralAbilities{Martial: 11, Command: 15, Politics: 2, Tactic: 7, Aptitude: [3]int{10, 0, 3}}
	if !w.ApplyGeneralAbilities(0, want) {
		t.Fatal("valid draft rejected")
	}
	dir := t.TempDir()
	src, dst := filepath.Join(dir, "source.dat"), filepath.Join(dir, "save.dat")
	if err := os.WriteFile(src, raw, 0600); err != nil {
		t.Fatal(err)
	}
	if err := w.SaveInto(src, dst, 1); err != nil {
		t.Fatal(err)
	}
	readback, err := LoadScenario(dst, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got := readback.Generals[0].Abilities(); got != want {
		t.Fatalf("legacy slot decoded %+v, want %+v", got, want)
	}
	r := readback.raw[generalBase:]
	if r[0x0E] != 0xA2 || r[0x0F] != 0x07 || r[0x10] != 0x39 {
		t.Fatalf("aptitude bits changed outside high nibble: %x", r[0x0E:0x11])
	}
}
