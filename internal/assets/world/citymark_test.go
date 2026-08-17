package world

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCityCentreTileByOwnership 釘住 docs/spec/53 §2：地圖檔存的是無所屬那一張。
func TestCityCentreTileByOwnership(t *testing.T) {
	for _, base := range []byte{205, 208, 211} {
		for _, c := range []struct {
			own  Ownership
			want byte
		}{
			{OwnedBySelf, base - 2},
			{OwnedByOther, base - 1},
			{Unowned, base},
		} {
			if got := CityCentreTile(base, c.own); got != c.want {
				t.Errorf("CityCentreTile(%d, %v) = %d，預期 %d", base, c.own, got, c.want)
			}
		}
	}
	// 不是據點中心就原樣回傳——寧可少畫也不要換錯格子。
	for _, other := range []byte{0, 204, 206, 212, 255} {
		if got := CityCentreTile(other, OwnedBySelf); got != other {
			t.Errorf("CityCentreTile(%d, 自勢力) 動了非中心格 → %d", other, got)
		}
	}
}

// TestOwnershipOf 釘住 0x18 ＝ 無所屬（docs/re/62 §2）。
func TestOwnershipOf(t *testing.T) {
	if got := OwnershipOf(NoOwner, 3); got != Unowned {
		t.Errorf("0x18 應該是無所屬，得到 %v", got)
	}
	if got := OwnershipOf(3, 3); got != OwnedBySelf {
		t.Errorf("同編號應該是自勢力，得到 %v", got)
	}
	if got := OwnershipOf(4, 3); got != OwnedByOther {
		t.Errorf("不同編號應該是他勢力，得到 %v", got)
	}
}

// TestCityCentreIsRecordPlusFour 用原版素材驗 docs/spec/53 §5：
// **192 座據點的 (X+4, Y) 全部是據點中心圖塊，而且全圖只有這 192 格。**
//
// 這一條同時擋兩件事：+4 這個位移漂掉，以及據點表與地圖對不上。
func TestCityCentreIsRecordPlusFour(t *testing.T) {
	dir := filepath.Join("..", "..", "..", "workplace", "orig", "dosv")
	raw, err := os.ReadFile(filepath.Join(dir, "MMAP.MAP"))
	if err != nil {
		t.Skipf("找不到原版素材，跳過：%v", err)
	}
	m, err := ParseMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, v := range m.Tiles {
		if IsCityCentre(v) {
			total++
		}
	}
	if total != 192 {
		t.Errorf("全圖的據點中心格有 %d 個，預期 192（＝據點數）", total)
	}
}
