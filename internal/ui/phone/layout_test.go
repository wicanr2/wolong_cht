package phone

import "testing"

// 四個指令鈕不重疊、都在指令列裡，而且**每一個都比 48 dp 的觸控下限寬**。
func TestCommandRectsAreTappableAndDisjoint(t *testing.T) {
	// 邏輯畫布 960 寬對應到常見手機的短邊；48 dp 在這個尺度上約 48 邏輯 px。
	const minTouch = 48
	var prevRight int
	for i := 0; i < int(numCommands); i++ {
		x, y, w, h := CommandRect(i)
		if w < minTouch || h < minTouch {
			t.Errorf("第 %d 個指令鈕 %d×%d，比觸控下限 %d 小", i, w, h, minTouch)
		}
		if y < LogicalH-CommandH || y+h > LogicalH {
			t.Errorf("第 %d 個指令鈕超出指令列：y=%d h=%d", i, y, h)
		}
		if x < prevRight {
			t.Errorf("第 %d 個指令鈕與前一個重疊：x=%d，前一個右緣 %d", i, x, prevRight)
		}
		prevRight = x + w
	}
	if prevRight > LogicalW {
		t.Errorf("最後一個指令鈕右緣 %d 超出畫面寬 %d", prevRight, LogicalW)
	}
}

// 地圖區要把狀態列與指令列之間的空間用滿，不留縫。
func TestMapRectFillsTheMiddle(t *testing.T) {
	x, y, w, h := MapRect()
	if x != 0 || w != LogicalW {
		t.Errorf("地圖區沒有用滿寬度：x=%d w=%d", x, w)
	}
	if y != StatusH || y+h != LogicalH-CommandH {
		t.Errorf("地圖區與上下兩條之間有縫：y=%d h=%d", y, h)
	}
}

func TestZoomIsIntegerOnly(t *testing.T) {
	for _, c := range []struct{ in, want int }{{0, 1}, {1, 1}, {2, 2}, {3, 3}, {4, 3}, {-5, 1}} {
		if got := ClampZoom(c.in); got != c.want {
			t.Errorf("ClampZoom(%d) = %d，預期 %d", c.in, got, c.want)
		}
	}
}
