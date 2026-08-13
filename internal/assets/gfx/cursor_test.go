package gfx

import "testing"

func TestDOSVCursorResource(t *testing.T) {
	raw := read(t, "KI.EXE")
	idx, err := DecodeDOSVCursor(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(idx), DOSVCursorWidth*DOSVCursorHeight; got != want {
		t.Fatalf("DOS/V cursor 像素數 = %d，want %d", got, want)
	}
	counts := map[byte]int{}
	for _, v := range idx {
		counts[v]++
	}
	if counts[DOSVCursorOutline] != 39 || counts[DOSVCursorFill] != 56 ||
		counts[DOSVCursorTransparent] != 161 {
		t.Fatalf("DOS/V cursor mask 統計 = %#v，want white=39 red=56 transparent=161", counts)
	}
	wantRows := []string{
		"WW..............",
		"WRWW............",
		".WRRWW..........",
		".WRRRRWW........",
		"..WRRRRRWW......",
		"..WRRRRRRRWW....",
		"...WRRRRRRRRWW..",
		"...WRRRRRRRW....",
		"....WRRRRRW.....",
		"....WRRRRRRW....",
		".....WRRWRRRW...",
		".....WRW.WRRRW..",
		"......W...WRRW..",
		"......W....WW...",
		"................",
		"................",
	}
	for y, want := range wantRows {
		row := make([]byte, DOSVCursorWidth)
		for x := range row {
			switch idx[y*DOSVCursorWidth+x] {
			case DOSVCursorOutline:
				row[x] = 'W'
			case DOSVCursorFill:
				row[x] = 'R'
			default:
				row[x] = '.'
			}
		}
		if got := string(row); got != want {
			t.Fatalf("cursor row %d = %q，want %q", y, got, want)
		}
	}
}
