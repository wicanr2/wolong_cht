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

// TestCityCentreIsRecordCoordinate 用原版素材驗 docs/spec/53 §5：
// **192 座據點的記錄座標本身就是據點中心圖塊，而且全圖只有這 192 格。**
//
// ⭐ 這一條的第二個迴圈是**位移的正對照**：把地圖多讀／少讀四格
// （`MMAP.MAP` 解壓後開頭那 4 byte 是長度欄位，見 world.MapHeader）
// 命中數會從 192 掉到 0。少了它，「位移錯了」與「據點表對不上」
// 在測試裡長得一模一樣。
func TestCityCentreIsRecordCoordinate(t *testing.T) {
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

	// 據點表：劇本 0 的區塊 +0x08C0 起，192 筆 × 32 B，
	// +0x08／+0x0A 是 X／Y（docs/formats/08 §1.6）。
	scen, err := os.ReadFile(filepath.Join(dir, "SINARIO.DAT"))
	if err != nil {
		t.Skipf("找不到 SINARIO.DAT，跳過座標比對：%v", err)
	}
	const cityTable, recSize, count = 0x08C0, 32, 192
	if len(scen) < cityTable+recSize*count {
		t.Fatalf("SINARIO.DAT 只有 %d B，讀不到據點表", len(scen))
	}
	for _, shift := range []int{0, -4, 4} {
		hit := 0
		for i := 0; i < count; i++ {
			r := scen[cityTable+recSize*i:]
			x := (int(r[8]) | int(r[9])<<8) + shift
			y := int(r[10]) | int(r[11])<<8
			if x < 0 || x >= Width || y < 0 || y >= Height {
				continue
			}
			if IsCityCentre(m.Tiles[y*Width+x]) {
				hit++
			}
		}
		switch shift {
		case 0:
			if hit != count {
				t.Errorf("記錄座標命中中心格 %d／%d——地圖的起始位移錯了", hit, count)
			}
		default:
			if hit != 0 {
				t.Errorf("位移 %+d 也命中了 %d 座，這個檢查沒有鑑別力", shift, hit)
			}
		}
	}
}
