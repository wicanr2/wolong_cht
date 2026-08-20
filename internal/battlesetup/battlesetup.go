// Package battlesetup 把戰場資產接到規則層的戰術入口。
//
// **桌面版與手機版共用這一份。** 選哪一張戰場、要不要翻轉、腳本從哪一段取，
// 三件事各有幾個容易錯的細節（野戰的翻轉來自地形配對、攻城的來自
// 「玩家守城」、腳本段編號 ＝ tactic × 4 ＋ 戰場類別）。
// 兩個 UI 各寫一份的話，其中一邊會在特定局面開出不同的戰場。
package battlesetup

import (
	"fmt"
	"os"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/rules/battlefield"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// Provider 是一個世界的戰場來源。
type Provider struct {
	lib   *battle.Library
	sp    *battle.Sprites
	world *state.World
	wmap  *world.Map

	// rotate 是**最後一次建場**用的翻轉旗標。繪圖端要用同一個值，
	// 否則旗與地形會對不上（docs/playtest/40 §11）。
	rotate bool
}

// Options 是 Load 要的東西。
type Options struct {
	// Dir 是原版素材目錄。
	Dir string
	// World 是規則層的世界；戰場的選擇要看據點座標與軍團朝向。
	World *state.World
	// Map 是大地圖，野戰靠它取樣地形。缺席時野戰一律用預設那一張。
	Map *world.Map
	// Warn 收載入過程的警告，nil 就丟掉。**三份資料各自可缺**，
	// 少一份就退一級，不會讓遊戲開不起來。
	Warn func(string)
}

// Load 載入戰場資產並組出戰術來源。
//
// 陣形表載不到就回錯誤——**沒有它玩家的戰鬥只能自動判定**，
// 那是一個要讓呼叫端知道的降級，不是可以靜靜吞掉的細節。
func Load(opt Options) (*Provider, *state.TacticalSetup, error) {
	warn := opt.Warn
	if warn == nil {
		warn = func(string) {}
	}
	forms, err := tactical.LoadFormations(opt.Dir + "/KI.EXE")
	if err != nil {
		return nil, nil, fmt.Errorf("載不到陣形表：%w", err)
	}
	p := &Provider{world: opt.World, wmap: opt.Map}
	p.lib = loadLibrary(opt.Dir, warn)
	if raw, err := os.ReadFile(opt.Dir + "/BATTLE.SCH"); err == nil {
		if sp, err := battle.ParseSprites(raw); err == nil {
			p.sp = sp
		} else {
			warn(fmt.Sprintf("%v；兵會畫成色點", err))
		}
	} else {
		warn(fmt.Sprintf("載不到 BATTLE.SCH（%v）；兵會畫成色點", err))
	}
	setup := &state.TacticalSetup{
		Forms: forms,
		Field: p.BuildField,
		Script: func(node int, siege bool, tactic int) []byte {
			if p.lib == nil {
				return nil
			}
			return p.lib.Script(tactic, battle.Category(p.FieldNumber(node, siege)))
		},
	}
	return p, setup, nil
}

// Library 是戰場地形與腳本，載不到就是 nil。
func (p *Provider) Library() *battle.Library { return p.lib }

// Sprites 是人物圖形，載不到就是 nil（兵畫成色點）。
func (p *Provider) Sprites() *battle.Sprites { return p.sp }

// Rotate 是最後一次建場的翻轉旗標。
func (p *Provider) Rotate() bool { return p.rotate }

func loadLibrary(dir string, warn func(string)) *battle.Library {
	read := func(name string) []byte {
		b, err := os.ReadFile(dir + "/" + name)
		if err != nil {
			warn(fmt.Sprintf("載不到 %s（%v）", name, err))
			return nil
		}
		return b
	}
	m, mdl := read("BATTLE.MAP"), read("BATTLE.MDL")
	if m == nil || mdl == nil {
		warn("沒有戰場地形，會用生成的替代")
		return nil
	}
	lib, err := battle.Parse(m, mdl, read("BATTLE.DAT"))
	if err != nil {
		warn(fmt.Sprintf("%v；會用生成的地形", err))
		return nil
	}
	return lib
}

// BuildField 取一張戰場。
//
//   - **攻城戰**：戰場編號就是據點編號（`docs/re/05`）
//   - **野戰**：從大地圖上即時算（`internal/rules/battlefield`）——
//     取軍團所在格與下方四格的地形類型去配一張 21 筆的表
func (p *Provider) BuildField(node int, siege, rotate bool) *tactical.Field {
	n := p.FieldNumber(node, siege)
	if p.lib == nil || n < 0 || n >= battle.NumFields {
		return SyntheticField(siege)
	}
	// 野戰的翻轉來自地形配對（`sub_14BDD` 換順序才配上），
	// 攻城的來自「玩家守城」（`sub_14ED7`）——兩條路合成同一個旗標。
	if !siege {
		rotate = p.FieldRotates(node)
	}
	p.rotate = rotate
	// 用原始圖塊值建，不只用堆疊高度——城壁與門是從圖塊值認出來的，
	// 而且打壞時要換圖塊再重算高度（docs/re/11 §5.9）。
	// 再加上七層子圖塊表：兩個平面的地面圖與登城都靠它算（docs/spec/36）。
	tiles, gate := p.lib.Tiles(n), p.lib.GateX(n)
	if rotate {
		tiles, gate = battle.Rotate180(tiles), battle.RotateGateX(gate)
	}
	return tactical.NewFieldFromTileLayers(
		tiles, p.lib.Heights(n), p.lib.TileLayers(n), gate)
}

// FieldNumber 回傳這一場用第幾張戰場：攻城就是據點編號，野戰現算。
func (p *Provider) FieldNumber(node int, siege bool) int {
	if siege {
		return node
	}
	return p.fieldForNode(node)
}

// fieldForNode 依大地圖的地形算出野戰要用哪一張戰場。
//
// 取樣的五格與 `sub_14B63` 一致（中心、下、左下、右下、兩格下方），
// 取樣方向用**玩家那一側軍團的朝向**（軍團記錄 `+0x08`）——原版是
// `cmp ah, [si+1] / mov al, [si+8]`，只在勢力等於玩家時才取。
func (p *Provider) fieldForNode(node int) int {
	n, ok := p.neighbours(node)
	if !ok {
		return battlefield.FieldBase + 6
	}
	f, _ := battlefield.Select(p.playerHeading(), n)
	return f
}

// FieldRotates 回報這一張野戰要不要翻轉。
func (p *Provider) FieldRotates(node int) bool {
	n, ok := p.neighbours(node)
	if !ok {
		return false
	}
	_, rot := battlefield.Select(p.playerHeading(), n)
	return rot
}

// neighbours 取據點所在格與下方四格的地形類型。
//
// ⚠ **野外的節點沒有座標**，只有據點有；取不到就讓呼叫端走預設。
func (p *Provider) neighbours(node int) (battlefield.Neighbours, bool) {
	if p.wmap == nil || p.world == nil || node < 0 || node >= len(p.world.Cities) {
		return battlefield.Neighbours{}, false
	}
	cx, cy := p.world.Cities[node].X, p.world.Cities[node].Y
	at := func(dx, dy int) int {
		t, err := p.wmap.Tile(cx+dx, cy+dy)
		if err != nil {
			return 0
		}
		return battlefield.Terrain(t)
	}
	return battlefield.Neighbours{
		Centre:    at(0, 0),
		Down:      at(0, 1),
		DownLeft:  at(-1, 1),
		DownRight: at(1, 1),
		TwoDown:   at(0, 2),
	}, true
}

// playerHeading 回傳玩家目前在動的那一支軍團的朝向。
//
// 原版只在「軍團的勢力 ＝ 玩家的勢力」時才拿 `+0x08`（`sub_14B63` 開頭那兩個
// `cmp ah, [si+1]`），所以戰場的取樣方向是**從玩家的視角**定的。
// 玩家沒有在移動的軍團時退回靜止（4），配對表會走預設那一支。
func (p *Provider) playerHeading() int {
	if p.world == nil {
		return state.HeadingStill
	}
	for i := range p.world.Corps {
		c := &p.world.Corps[i]
		if c.Alive && c.Faction == p.world.Player && c.Heading != state.HeadingStill {
			return c.Heading
		}
	}
	return state.HeadingStill
}

// SyntheticField 是沒有原版戰場資料時的替代品。幾何同尺寸，內容是自己生的。
func SyntheticField(siege bool) *tactical.Field {
	stack := make([][]int, tactical.Height)
	for y := range stack {
		stack[y] = make([]int, tactical.Width)
	}
	if !siege {
		return tactical.NewField(stack, 0)
	}
	const wallX, top, bottom = 40, 8, tactical.Height - 9
	gate := tactical.Height / 2
	for y := top; y <= bottom; y++ {
		if y != gate {
			stack[y][wallX] = 4
		}
	}
	for x := wallX; x < tactical.Width-1; x++ {
		stack[top][x] = 4
		stack[bottom][x] = 4
	}
	return tactical.NewField(stack, wallX)
}

// StageOptions 是驗收捷徑「擺一場戰鬥出來」要的參數。
type StageOptions struct {
	Siege bool
	// Node 是攻城要打哪一座城，−1 表示用守方現在待的那一個。
	Node int
	// Attacker／Defender 是兩支軍團的編號。
	Attacker, Defender int
}

// StageEncounter 把攻方與守方擺到會遭遇的位置，跑到開打為止。
//
// ⚠ 這是**驗收捷徑，不是規則**：它只搬座標與目標，遭遇判定仍然走正常的
// `resolveContact`。兩個 UI 的驗收路徑共用這一支——擺位的規則只留一份，
// 否則兩邊會在不同的局面開打，比出來的畫面沒有意義。
//
// 回傳 true 表示真的開出了一場（或停在遭遇決策上）。
func StageEncounter(w *state.World, rng combat.Rand, opt StageOptions) bool {
	if w == nil || opt.Attacker < 0 || opt.Defender < 0 ||
		opt.Attacker >= len(w.Corps) || opt.Defender >= len(w.Corps) ||
		opt.Attacker == opt.Defender {
		return false
	}
	me, foe := &w.Corps[opt.Attacker], &w.Corps[opt.Defender]
	if !me.Alive || !foe.Alive {
		return false
	}
	if opt.Siege {
		// 攻城：守方待在自己的城裡，攻方從隔壁一格走進去。
		node := foe.Node
		if opt.Node >= 0 && opt.Node < len(w.Cities) {
			node = opt.Node
			// 攻城要打的是**守方的城**，所以捷徑把這座城暫時記到守方名下。
			// 只影響 fixture，不是規則。
			w.Cities[node].Owner = foe.Faction
			foe.Node, foe.TargetNode = node, node
			foe.X, foe.Y = w.Cities[node].X, w.Cities[node].Y
			foe.TargetX, foe.TargetY = foe.X, foe.Y
		}
		me.Node = node
		me.X, me.Y = w.Cities[node].X-1, w.Cities[node].Y
		me.TargetNode = node
		me.TargetX, me.TargetY = w.Cities[node].X, w.Cities[node].Y
		me.Timer = 1
	} else {
		// 野戰：把敵方放在隔壁一格、目標設成我方所在的那一格——
		// 下一次輪到它移動就會撞上（遭遇條件是「同格、不同勢力」）。
		foe.X, foe.Y = me.X-1, me.Y
		foe.TargetX, foe.TargetY = me.X, me.Y
		foe.TargetNode = me.Node
		foe.Timer = 1
	}
	for i := 0; i < 64 && w.PendingBattle() == nil && w.PendingEncounter() == nil; i++ {
		w.Tick(rng)
	}
	return w.PendingBattle() != nil || w.PendingEncounter() != nil
}
