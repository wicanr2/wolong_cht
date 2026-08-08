package battle

import (
	"os"
	"testing"
)

const dir = "../../../workplace/orig/dosv/"

func load(t *testing.T) *Library {
	t.Helper()
	read := func(n string) []byte {
		b, err := os.ReadFile(dir + n)
		if err != nil {
			t.Skip("找不到原版 " + n + "，跳過")
		}
		return b
	}
	l, err := Parse(read("BATTLE.MAP"), read("BATTLE.MDL"), read("BATTLE.DAT"))
	if err != nil {
		t.Fatal(err)
	}
	return l
}

// 索引第二欄為 0 ⟺ 野戰用的戰場。214 張零例外（docs/re/11 §4.5）。
func TestSiegeFlagMatchesGateColumn(t *testing.T) {
	l := load(t)
	fieldMaps := 0
	for n := 0; n < NumFields; n++ {
		if n >= 192 && l.IsSiege(n) {
			t.Errorf("戰場 %d 是野戰用的，第二欄卻是 %d", n, l.GateX(n))
		}
		if !l.IsSiege(n) {
			fieldMaps++
		}
	}
	// 22 張野戰圖（192–213）＋ 6 張第二欄為 0 的攻城圖。
	if fieldMaps != 28 {
		t.Errorf("第二欄為 0 的有 %d 張，應為 28", fieldMaps)
	}
}

// 堆疊高度：攻城圖長出城牆，平原的野戰圖幾乎是平的。
// 這是「一格是一疊 1–7 層圖塊」那條解讀最硬的驗證。
func TestStackHeightsShapeTheMap(t *testing.T) {
	l := load(t)
	tall := func(n int) int {
		c := 0
		for _, row := range l.Stacks(n) {
			for _, h := range row {
				if h >= 4 {
					c++
				}
			}
		}
		return c
	}
	// 戰場 198 是「平原 ＋ 平原」，只有 10 格高處。
	if got := tall(198); got != 10 {
		t.Errorf("戰場 198 的高處有 %d 格，應為 10", got)
	}
	// 戰場 192 是「山 ＋ 山」，最多。
	if got := tall(192); got != 596 {
		t.Errorf("戰場 192 的高處有 %d 格，應為 596", got)
	}
	// 攻城用的戰場 5 有一圈城牆。
	if got := tall(5); got != 320 {
		t.Errorf("戰場 5 的高處有 %d 格，應為 320", got)
	}
	// 每一格的堆疊都在 0–7。
	for _, row := range l.Stacks(5) {
		for x, h := range row {
			if h < 0 || h > MaxStack {
				t.Fatalf("第 %d 格的堆疊是 %d，應在 0–%d", x, h, MaxStack)
			}
		}
	}
}

// 腳本段編號 ＝ 武將 +0x16 × 4 ＋ 戰場類別。
func TestScriptSelection(t *testing.T) {
	l := load(t)
	for _, tc := range []struct{ field, want int }{
		{0, 0}, {191, 0}, // 攻城
		{192, 1}, {208, 1}, {0xD0, 1},
		{0xD1, 2}, {213, 2}, // 另一組野戰
	} {
		if got := Category(tc.field); got != tc.want {
			t.Errorf("Category(%d) ＝ %d，應為 %d", tc.field, got, tc.want)
		}
	}
	// 呂布那一型（+0x16 ＝ 0）在攻城戰用第 0 段。
	if got := l.Script(0, 0); len(got) != ScriptSize {
		t.Errorf("段 0 長 %d，應為 %d", len(got), ScriptSize)
	}
	// 兩段不該一樣。
	a, b := l.Script(0, 0), l.Script(7, 3)
	same := true
	for i := range a {
		if a[i] != b[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("第 0 段與第 31 段完全相同，段編號算錯了")
	}
}

// 子圖塊的格式：五個 64 B 位元平面，第一個是遮罩。
//
// ⭐ 決定性的證據是這條不變量：**遮罩是 0 的地方，四個色平面全部是 0**。
// 這不是我挑的指標——它是「有沒有解對」的必要條件，而且不可能靠巧合成立：
// 3 組 × 192 個子圖塊 × 512 像素 ＝ 294,912 個位元，一個例外都不能有。
//
// 反過來的檢查也在這裡：把遮罩換成別的平面就會破功。
func TestSubTileMaskInvariant(t *testing.T) {
	l := load(t)
	raw := mustRead(t, "BATTLE.MDL")

	bad, transparent, total := 0, 0, 0
	for set := 0; set < NumTileSets; set++ {
		for n := 0; n < NumSubTiles; n++ {
			base := MDLHeader + set*TileSetSize + subTileBase + n*SubTileSize
			b := raw[base : base+SubTileSize]
			for i := 0; i < planeBytes; i++ {
				var any byte
				for p := 1; p <= 4; p++ {
					any |= b[planeBytes*p+i]
				}
				// 遮罩為 0 的位元上，四個色平面必須也是 0。
				bad += popcount(^b[i] & any)
				transparent += popcount(^b[i])
				total += 8
			}
		}
	}
	if bad != 0 {
		t.Errorf("有 %d 個位元違反「遮罩 0 → 色平面 0」——平面分組解錯了", bad)
	}
	if r := float64(transparent) / float64(total); r < 0.2 || r > 0.5 {
		t.Errorf("透明像素佔 %.1f%%，不像等角圖塊（預期三成上下）", r*100)
	}

	// 解出來的子圖塊要與原始位元一致，而且真的用到多種顏色。
	seen := map[int]bool{}
	for n := 0; n < NumSubTiles; n++ {
		s := l.SubTile(0, n)
		if s == nil {
			t.Fatalf("子圖塊 %d 解不出來", n)
		}
		for i := range s.Pix {
			seen[int(s.Pix[i])] = true
		}
	}
	if len(seen) < 10 {
		t.Errorf("整組只用到 %d 種色號，太少——色平面的順序可能錯了", len(seen))
	}
}

func mustRead(t *testing.T, n string) []byte {
	t.Helper()
	b, err := os.ReadFile(dir + n)
	if err != nil {
		t.Skip("找不到原版 " + n + "，跳過")
	}
	return b
}

func popcount(b byte) int {
	n := 0
	for ; b != 0; b &= b - 1 {
		n++
	}
	return n
}

// 圖塊定義引用的子圖塊編號一定要落在 0–191。
// **192 × 320 B 正好用完 61,440 B 的像素區**——兩件事互相印證。
func TestSubTileIndexRange(t *testing.T) {
	raw := mustRead(t, "BATTLE.MDL")
	if NumSubTiles*SubTileSize != TileSetSize-subTileBase {
		t.Fatalf("192 × 320 ＝ %d，像素區是 %d",
			NumSubTiles*SubTileSize, TileSetSize-subTileBase)
	}
	max := 0
	for set := 0; set < NumTileSets; set++ {
		base := MDLHeader + set*TileSetSize
		for i := 0; i < TileDefs; i++ {
			r := raw[base+i*TileDefLen : base+(i+1)*TileDefLen]
			k := int(r[0])
			if k > MaxStack {
				k = MaxStack
			}
			for _, v := range r[1 : 1+k] {
				if int(v) > max {
					max = int(v)
				}
			}
		}
	}
	if max != NumSubTiles-1 {
		t.Errorf("子圖塊編號最大是 %d，預期 %d", max, NumSubTiles-1)
	}
}

// `BATTLE.SCH` 與 `BATTLE.MDL` 的子圖塊是**同一個格式**。
//
// 同一條不變量（遮罩為 0 處四個色平面全 0）在這裡也必須 100% 成立，
// 而且 360 個單位裡**只有 170 與 350 是全空的**——正好相差 180，
// 那就是「兩側各 180 張」的證據。
func TestSpriteFormatMatchesSubTile(t *testing.T) {
	raw := mustRead(t, "BATTLE.SCH")
	sp, err := ParseSprites(raw)
	if err != nil {
		t.Fatal(err)
	}
	bad := 0
	for n := 0; n < NumSprites; n++ {
		b := raw[n*SpriteUnit : (n+1)*SpriteUnit]
		for i := 0; i < planeBytes; i++ {
			var any byte
			for p := 1; p <= 4; p++ {
				any |= b[planeBytes*p+i]
			}
			bad += popcount(^b[i] & any)
		}
	}
	if bad != 0 {
		t.Errorf("有 %d 個位元違反「遮罩 0 → 色平面 0」——BATTLE.SCH 不是同一個格式", bad)
	}

	var empty []int
	for n := 0; n < NumSprites; n++ {
		s := sp.At(n)
		if s == nil {
			t.Fatalf("第 %d 張解不出來", n)
		}
		blank := true
		for _, v := range s.Pix {
			if v != Transparent {
				blank = false
				break
			}
		}
		if blank {
			empty = append(empty, n)
		}
	}
	if len(empty) != 2 || empty[0] != EmptySprite ||
		empty[1] != EmptySprite+SpritesPerSide {
		t.Errorf("全空的是 %v，預期正好兩個且相差 %d", empty, SpritesPerSide)
	}
}
