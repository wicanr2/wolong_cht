package battlesetup

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/assets/gfx"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/state"
)

const original = "../../workplace/orig/dosv"

func TestBattleAssetsLoadAtFirstUse(t *testing.T) {
	if _, err := os.Stat(filepath.Join(original, "KI.EXE")); err != nil {
		t.Skip("缺少原版素材")
	}
	dir := t.TempDir()
	provide := func(name string) {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(original, name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), b, 0600); err != nil {
			t.Fatal(err)
		}
	}
	provide("KI.EXE")
	var warnings []string
	p, setup, err := Load(Options{Dir: dir, Warn: func(w string) { warnings = append(warnings, w) }})
	if err != nil || p == nil || setup == nil || setup.Forms == nil {
		t.Fatalf("關卡資源尚未到位時仍應能安裝陣形：%v", err)
	}
	if len(warnings) != 0 || p.lib != nil || p.sp != nil {
		t.Fatalf("開局階段不應讀戰場資源：%v", warnings)
	}
	for _, name := range []string{"BATTLE.MAP", "BATTLE.MDL", "BATTLE.DAT"} {
		provide(name)
	}
	field := setup.Field(82, true, false)
	lib := p.Library()
	if lib == nil || !field.HasTiles() || len(setup.Script(82, true, 2, false)) != tactical.ScriptCodeSize {
		t.Fatal("首次開戰應讀得到原版地形與腳本")
	}
	if len(warnings) != 0 || p.sp != nil {
		t.Fatalf("規則層建場不應讀單兵圖形：%v", warnings)
	}
	provide("BATTLE.SCH")
	sprites := p.Sprites()
	if sprites == nil || sprites.SourceTile(battle.CombinedSourceTerrainTiles) == nil {
		t.Fatal("首次建立畫面應讀得到單兵圖形")
	}
	for _, name := range []string{"BATTLE.MAP", "BATTLE.MDL", "BATTLE.DAT", "BATTLE.SCH"} {
		if err := os.Remove(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
	}
	if p.Library() != lib || p.Sprites() != sprites || len(warnings) != 0 {
		t.Fatalf("同一資料來源不應重讀已載入的檔案：%v", warnings)
	}
	missing, missingSetup, err := Load(Options{Dir: dir, Warn: func(w string) { warnings = append(warnings, w) }})
	if err != nil || len(warnings) != 0 {
		t.Fatalf("缺檔到首次使用才回報：%v，%v", err, warnings)
	}
	if missingSetup.Field(82, true, false).HasTiles() || missing.Sprites() != nil || len(warnings) == 0 {
		t.Fatal("首次使用缺圖素材時應啟用地形與人物替代")
	}
}

// 原版節點21的留守讀點是圖庫值0；缺少注入會沿未知近似出擊。
// 同時驗新遊戲共用Load及讀檔重新掛來源，不能只驗孤立setter。
func TestOriginalAIHoldSourceSurvivesWorldReload(t *testing.T) {
	icon, err := os.ReadFile(original + "/ICONGRF.DAT")
	if err != nil {
		t.Skip("缺少原版素材")
	}
	region := gfx.IconRegions[3]
	if len(icon) < region.Offset+region.Length {
		t.Fatal("原圖庫範圍不完整")
	}
	chrome := icon[region.Offset : region.Offset+region.Length]
	data, err := os.ReadFile(original + "/MMAP.MAP")
	if err != nil {
		t.Fatal(err)
	}
	m, err := world.ParseMap(data)
	if err != nil {
		t.Fatal(err)
	}
	for _, present := range []bool{true, false} {
		w, err := state.LoadScenario(original+"/SINARIO.DAT", 0)
		if err != nil {
			t.Fatal(err)
		}
		w.Player = 0
		city := &w.Cities[21]
		city.Owner, city.Threatened, city.Specific, city.Occupancy = 12, true, false, 1
		w.Corps[0] = state.Corps{Alive: true, Faction: 12, Node: 21, X: city.X, Y: city.Y,
			TargetNode: 22, Ordered: 21, Stage: 1, Timer: 1, Interval: 3, Men: 600, Morale: 200}
		source := chrome
		if !present {
			source = nil
		}
		if _, _, err := Load(Options{Dir: original, World: w, Map: m, Chrome: source}); err != nil {
			t.Fatal(err)
		}
		block := w.Bytes()
		restored, err := state.LoadBlock(block)
		if err != nil {
			t.Fatal(err)
		}
		if err := restored.Restore(w.TakeSnapshot()); err != nil {
			t.Fatal(err)
		}
		AttachAIHoldSources(restored, source, m)
		if !bytes.Equal(block, restored.Bytes()) {
			t.Fatal("掛回原素材改寫了存檔")
		}
		for _, current := range []*state.World{w, restored} {
			current.Tick(rng.NewFixed(7))
			want := 2
			if present {
				want = 1
			}
			if got := current.Corps[0].Stage; got != want {
				t.Fatalf("來源存在=%v：Stage=%d，預期%d", present, got, want)
			}
		}
	}
}

func TestProviderLoadsExplicitSelectionForEveryEncounter(t *testing.T) {
	if _, err := os.Stat(original + "/KI.EXE"); err != nil {
		t.Skip("缺少原版素材")
	}
	w, err := state.LoadScenario(original+"/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	p, s, err := Load(Options{Dir: original, World: w})
	if err != nil {
		t.Fatal(err)
	}
	lib := p.Library()
	if lib == nil {
		t.Fatal("missing battle library")
	}
	for _, tc := range []struct {
		number int
		rotate bool
	}{{193, true}, {210, false}, {82, true}, {211, false}, {193, false}} {
		tiles, gate := lib.MinimapRows(tc.number), lib.GateX(tc.number)
		if tc.rotate {
			copy(tiles[1:len(tiles)-1], battle.Rotate180(tiles[1:len(tiles)-1]))
			gate = battle.RotateGateX(gate)
		}
		want := tactical.NewFieldFromTileLayers(tiles, lib.Heights(tc.number), lib.TileLayers(tc.number), gate)
		got := s.Field(tc.number, tc.number < 192, tc.rotate)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("wrong geometry for explicit field %d rotate %v", tc.number, tc.rotate)
		}
		if got := s.Category(tc.number, tc.number < 192); got != battle.Category(tc.number) {
			t.Fatalf("category %d", got)
		}
		if !reflect.DeepEqual(s.Script(tc.number, false, 2, false), lib.Script(2, battle.ScriptSlot(tc.number, false))) {
			t.Fatalf("script field %d", tc.number)
		}
	}
}

func TestProviderSelectsFourOriginalScriptSlots(t *testing.T) {
	if _, err := os.Stat(original + "/KI.EXE"); err != nil {
		t.Skip("缺少原版素材")
	}
	p, setup, err := Load(Options{Dir: original})
	if err != nil {
		t.Fatal(err)
	}
	lib := p.Library()
	if lib == nil {
		t.Fatal("缺少真實戰場腳本")
	}
	for _, tc := range []struct {
		field, slot int
		defends     bool
	}{{76, 1, false}, {76, 0, true}, {191, 1, false}, {191, 0, true},
		{192, 2, false}, {208, 2, true}, {209, 3, false}, {213, 3, true}} {
		for _, tactic := range []int{0, 7} {
			got := setup.Script(tc.field, tc.field < 192, tactic, tc.defends)
			if len(got) != tactical.ScriptCodeSize || !reflect.DeepEqual(got, lib.Script(tactic, tc.slot)) {
				t.Fatalf("field=%d playerDefends=%v tactic=%d：腳本應取槽%d", tc.field, tc.defends, tactic, tc.slot)
			}
		}
	}
}

func TestProviderKeepsOriginalRuleRows(t *testing.T) {
	if _, err := os.Stat(original + "/KI.EXE"); err != nil {
		t.Skip("缺少原版素材")
	}
	p, _, err := Load(Options{Dir: original})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(original + "/BATTLE.MAP")
	if err != nil {
		t.Fatal(err)
	}
	const field = 76
	block := raw[battle.IndexSize+field*battle.FieldSize : battle.IndexSize+(field+1)*battle.FieldSize]
	for _, rotate := range []bool{false, true} {
		f := p.BuildField(field, true, rotate)
		for y := 0; y < 64; y++ {
			for x := 0; x < 64; x++ {
				want := block[y*64+x]
				if rotate && y > 0 && y < 63 {
					want = battle.RotateTile(block[(63-y)*64+63-x])
				}
				if got := f.Tile(x, y); got != want {
					t.Fatalf("原版列號 field=%d rotate=%v (%d,%d)：%#x，應為%#x", field, rotate, x, y, got, want)
				}
			}
		}
	}
	f := p.BuildField(field, true, false)
	for _, tc := range []struct{ x, y, z int }{{55, 36, 0}, {55, 26, 1}, {55, 24, 1}} {
		if z, ok := f.GroundLevel(tc.x, tc.y, tactical.PlaneLow); !ok || z != tc.z {
			t.Fatalf("原版低平面(%d,%d)=%d/%v，應為%d", tc.x, tc.y, z, ok, tc.z)
		}
	}
}

func TestMissingFormationsStillAttachTerrainForAutomaticBattles(t *testing.T) {
	if _, err := os.Stat(original + "/SINARIO.DAT"); err != nil {
		t.Skip("缺少原版素材")
	}
	w, err := state.LoadScenario(original+"/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	for _, f := range []int{0, 12} {
		w.Factions[f].Reserves = [economy.NumTroopTypes]int{6000, 6000, 6000}
		if err := w.FormCorps(w.Factions[f].Lord, [army.Positions]army.TroopType{}, [army.Positions]bool{true, true, true, true, true, true}); err != nil {
			t.Fatal(err)
		}
	}
	a, d := w.Factions[0].Lord, w.Factions[12].Lord
	w.Corps[d].X, w.Corps[d].Y = 100, 110
	m := &world.Map{Tiles: make([]byte, world.Width*world.Height)}
	m.Tiles[110*world.Width+100] = 0xCA
	if _, _, err := Load(Options{Dir: t.TempDir(), World: w, Map: m}); err == nil {
		t.Fatal("missing forms should report error")
	}
	// Supply a simple field afterwards. StageBattle selects from the terrain
	// attached even though asset loading failed, without touching city coords.
	w.SetTactical(&state.TacticalSetup{Forms: tactical.SyntheticFormations(), Field: func(int, bool, bool) *tactical.Field { return SyntheticField(false) }})
	if err := w.StageBattle(a, d, combat.Field, rng.NewFixed(3)); err != nil {
		t.Fatal(err)
	}
	if n := w.PendingBattle().FieldNumber; n < 209 || n > 212 {
		t.Fatalf("terrain was lost after asset error: %d", n)
	}
}
