package phone

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/ui/isoview"
)

// 原版的 480×368 視野要塞得進戰場區。
//
// ⚠ **裁掉一截不是小事**：看得到多少戰場會影響決策。這一條紅了就要改
// 版面（例如再收掉一條），不是改縮放。
func TestBattleFieldFitsTheOriginalViewport(t *testing.T) {
	_, _, w, h := BattleFieldRect()
	if w < isoview.NativeW || h < isoview.NativeH {
		t.Fatalf("戰場區 %d×%d，塞不下原版的 %d×%d",
			w, h, isoview.NativeW, isoview.NativeH)
	}
}

// 上排六格是**空間排列**（左翼 左備 主將 前鋒 右備 右翼），不是記錄順序。
func TestSquadStripUsesTheOriginalSpatialOrder(t *testing.T) {
	want := []string{"左翼", "左備", "主將", "前鋒", "右備", "右翼"}
	for i := range want {
		if got := unitLabel(squadSlot(i)); got != want[i] {
			t.Errorf("第 %d 格是 %q，原版是 %q", i, got, want[i])
		}
	}
}

// 命令列的第 i 個送的是原版第 i 列的命令碼——**畫面順序不是命令碼順序**。
func TestBattleCommandRowSendsTheOriginalCode(t *testing.T) {
	if got := battleCommandRow(0); got != tactical.Charge {
		t.Fatalf("第 0 個命令是 %v，原版第 0 列送的是突擊", got)
	}
	for i := range battle.SideCommandRowCode {
		if int(battleCommandRow(i)) != battle.SideCommandRowCode[i] {
			t.Errorf("第 %d 個命令碼 %d，表裡是 %d",
				i, battleCommandRow(i), battle.SideCommandRowCode[i])
		}
	}
}

// 六格與六命令都要比 48 dp 的觸控下限大，而且不重疊、不超出畫面。
func TestBattleRowsAreTappable(t *testing.T) {
	const minTouch = 48
	for _, f := range []func(int) (int, int, int, int){SquadRect, BattleCommandRect} {
		var prevRight int
		for i := 0; i < army.Positions; i++ {
			x, y, w, h := f(i)
			if w < minTouch || h < minTouch {
				t.Errorf("第 %d 格 %d×%d 比觸控下限小", i, w, h)
			}
			if x < prevRight {
				t.Errorf("第 %d 格與前一格重疊", i)
			}
			if y+h > LogicalH {
				t.Errorf("第 %d 格超出畫面底部", i)
			}
			prevRight = x + w
		}
	}
}

// 開了戰場之後：主區換成戰場、點擊不再走到大地圖，而且**時間仍在走**
//（戰術層沒有暫停）。
func TestBattleTakesOverTheScreen(t *testing.T) {
	s := newTestSession(t)
	if err := s.OpenDemoBattle(82); err != nil {
		t.Skipf("擺不出戰鬥：%v", err)
	}
	if !s.BattleActive() {
		t.Fatal("開了戰場卻不在戰場狀態")
	}
	capital := s.World().Factions[s.World().Player].Capital
	s.Select(capital)
	// 點在原本是指令列的位置：戰場上那一條不畫，點下去只能是命令列。
	x, y, w, h := BattleCommandRect(0)
	s.Tap(float64(x+w/2), float64(y+h/2))
	if s.SheetOpen() {
		t.Fatal("戰場上點到底列卻開了面板")
	}

	before := s.Battle().Frame
	for i := 0; i < 60; i++ {
		s.Tick()
	}
	if s.Battle() != nil && s.Battle().Frame == before {
		t.Fatal("戰場的時鐘沒有走——戰術層沒有暫停")
	}
}

// 長跑：連續推進二十萬幀（約六年遊戲時間），每次有東西擋著就選一個，看時鐘是不是一直在走。
//
// ⭐ 這一條擋的是「玩到某一天畫面不動了」——**擋住世界的四種狀態少接一個
// 就會這樣**，而且從畫面上看不出原因。委任與拒絕是這裡的選擇，
// 因為它們不需要再開任何畫面。
func TestLongRunNeverStalls(t *testing.T) {
	if testing.Short() {
		t.Skip("長跑，-short 時跳過")
	}
	s := newTestSession(t)
	start := s.World().Clock
	blocked := 0
	for i := 0; i < 200000; i++ {
		switch {
		case s.BattleActive():
			// 戰場自己會走完，不必介入。
		case s.ModalKind() != modalNone:
			s.PickModal(len(s.ModalOptions()) - 1) // 委任／拒絕
			blocked++
		}
		before := s.World().Clock
		s.Tick()
		if s.World().Clock == before && !s.BattleActive() &&
			s.ModalKind() == modalNone && s.Speed() == DefaultSpeed {
			// 節流器不是每一幀都推進，所以時鐘不動很正常；
			// 這裡只是不讓迴圈永遠空轉，實際判準在迴圈外。
			continue
		}
	}
	end := s.World().Clock
	if end == start {
		t.Fatalf("跑了二十萬幀時鐘完全沒動（擋住 %d 次）", blocked)
	}
	t.Logf("%d年%d月%d日 → %d年%d月%d日，中途擋住 %d 次",
		start.Year, start.Month, start.Day, end.Year, end.Month, end.Day, blocked)
}
