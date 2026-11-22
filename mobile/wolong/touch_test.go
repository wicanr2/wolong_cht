package wolongmobile

import (
	"image"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/ui/phone"
)

func TestPinchReleaseDoesNotIssueBattleCommand(t *testing.T) {
	const orig = "../../workplace/orig/dosv"
	if _, err := os.Stat(orig + "/SINARIO.DAT"); os.IsNotExist(err) {
		t.Skip("找不到原版素材")
	}
	s, err := phone.NewSession(phone.Options{OrigDir: orig, Seed: 5})
	if err != nil {
		t.Fatal(err)
	}
	if err := s.OpenDemoBattle(82); err != nil {
		t.Fatal(err)
	}
	s.Tick()
	g := &game{sess: s}
	b := s.Battle()
	soldier := &b.Sides[b.PlayerSide].Soldiers[0]
	before := soldier.Next
	x, y, w, h := phone.BattleCommandRect(0)
	command := image.Pt(x+w/2, y+h/2)
	px, py, pw, ph := phone.BattleFieldRect()
	center := image.Pt(px+pw/2, py+ph/2)
	pos := map[ebiten.TouchID]image.Point{1: center.Sub(image.Pt(20, 0)), 2: center.Add(image.Pt(20, 0))}
	at := func(id ebiten.TouchID) (int, int) { return pos[id].X, pos[id].Y }

	g.handleTouch([]ebiten.TouchID{1}, at)
	g.handleTouch([]ebiten.TouchID{1, 2}, at)
	pos[2] = center.Add(image.Pt(55, 0))
	g.handleTouch([]ebiten.TouchID{1, 2}, at)
	pos[1] = command
	g.handleTouch([]ebiten.TouchID{1}, at)
	g.handleTouch(nil, at)
	if soldier.Next != before {
		t.Fatal("releasing a pinch over the command strip issued an order")
	}
	g.handleTouch([]ebiten.TouchID{1}, at)
	g.handleTouch(nil, at)
	if soldier.Next != tactical.Charge {
		t.Fatal("a new tap should issue the selected battle command")
	}
}
