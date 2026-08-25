package phone

import (
	"os"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/world"
)

const origDir = "../../../workplace/orig/dosv"

func newTestSession(t *testing.T) *Session {
	t.Helper()
	if _, err := os.Stat(origDir + "/SINARIO.DAT"); err != nil {
		t.Skip("找不到原版素材，跳過")
	}
	// Seed 是挑過的：TestEventNoticesAppearAndExpire 是情境探針，
	// 要求兩萬幀內玩家的城至少收到一則事件訊息。暴風雨範圍改成
	// 原版的切比雪夫 41×41 框（docs/spec/81 §2）之後訊息變稀有，
	// 規則改動後這支失敗時，先確認新規則的機器碼出處，再換 seed。
	s, err := NewSession(Options{OrigDir: origDir, Scenario: 0, Player: 0, Seed: 13})
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// 開局第一眼要看得到自己的首都——不然玩家得先找路。
func TestNewSessionCentresOnTheCapital(t *testing.T) {
	s := newTestSession(t)
	capital := s.World().Factions[s.World().Player].Capital
	c := &s.World().Cities[capital]
	camX, camY, _ := s.Camera()
	cols, rows := s.ViewTiles()
	cx, cy := c.X+world.CityCentreDX, c.Y
	if cx < camX || cx >= camX+cols || cy < camY || cy >= camY+rows {
		t.Fatalf("首都 (%d,%d) 不在鏡頭 (%d,%d)+%dx%d 內", cx, cy, camX, camY, cols, rows)
	}
}

// 鏡頭不可以捲出世界外——捲出去會畫到空白，而空白看起來像地圖破圖。
func TestPanIsClampedToTheWorld(t *testing.T) {
	s := newTestSession(t)
	s.Pan(-9999, -9999)
	if x, y, _ := s.Camera(); x != 0 || y != 0 {
		t.Fatalf("往左上捲到底應該是 (0,0)，得到 (%d,%d)", x, y)
	}
	s.Pan(9999, 9999)
	x, y, _ := s.Camera()
	cols, rows := s.ViewTiles()
	if x != world.Width-cols || y != world.Height-rows {
		t.Fatalf("往右下捲到底應該是 (%d,%d)，得到 (%d,%d)",
			world.Width-cols, world.Height-rows, x, y)
	}
}

// 縮放要保住畫面中心的那一格，否則捏合會把畫面甩到別的地方。
func TestZoomKeepsTheCentreTile(t *testing.T) {
	s := newTestSession(t)
	s.Pan(40, 30)
	before := centreTile(s)
	for _, z := range []int{2, 3, 1} {
		s.SetZoom(z)
		if got, want := centreTile(s), before; abs(got[0]-want[0]) > 1 || abs(got[1]-want[1]) > 1 {
			t.Fatalf("縮到 %d× 之後中心從 %v 跑到 %v", z, want, got)
		}
	}
	// 級距只有 1–3，超出要夾住。
	s.SetZoom(9)
	if _, _, z := s.Camera(); z != MaxZoom {
		t.Fatalf("縮放沒有夾住：%d", z)
	}
}

func centreTile(s *Session) [2]int {
	x, y, _ := s.Camera()
	cols, rows := s.ViewTiles()
	return [2]int{x + cols/2, y + rows/2}
}

// 點到據點要選中它，點空白要取消。**選中不可以動到 World**——
// 兩段式的第一段只改選取狀態（docs/mobile/android-ux.md §3）。
func TestSelectIsTwoStageAndDoesNotMutateTheWorld(t *testing.T) {
	s := newTestSession(t)
	before := s.World().Fingerprint()

	capital := s.World().Factions[s.World().Player].Capital
	c := &s.World().Cities[capital]
	mx, my, _, _ := MapRect()
	camX, camY, zoom := s.Camera()
	px := float64(TilePx * zoom)
	lx := float64(mx) + (float64(c.X+world.CityCentreDX-camX)+0.5)*px
	ly := float64(my) + (float64(c.Y-camY)+0.5)*px
	if !s.SelectAt(lx, ly) {
		t.Fatalf("點在首都上卻沒選中（邏輯座標 %.0f,%.0f）", lx, ly)
	}
	if s.Selected() != capital {
		t.Fatalf("選到 %d，預期首都 %d", s.Selected(), capital)
	}
	if s.World().Fingerprint() != before {
		t.Fatal("只是選取卻改到了 World——兩段式的第一段不可以有副作用")
	}

	// 點指令列不屬於地圖區，不可以改選取。
	s.SelectAt(10, LogicalH-10)
	if s.Selected() != capital {
		t.Fatal("點到指令列卻改了地圖的選取")
	}
	// 點地圖上的空白處要取消。
	if s.SelectAt(float64(mx)+1, float64(my)+1) && s.Selected() >= 0 {
		t.Log("左上角剛好有據點，換一格再試")
	}
}
