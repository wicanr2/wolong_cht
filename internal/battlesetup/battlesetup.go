// Package battlesetup 把戰場資產接到規則層的戰術入口。
//
// 桌面與手機共用這份資產來源；規則層先決定戰場編號及翻轉，這裡只載入結果。
package battlesetup

import (
	"fmt"
	"os"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// Provider 是不依賴目前世界的戰場資產來源。
type Provider struct {
	lib *battle.Library
	sp  *battle.Sprites
}

// Options 是 Load 要的東西。
type Options struct {
	// Dir 是原版素材目錄。
	Dir string
	// World 是要掛上地形查詢的目前世界，Provider不保留它。
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
	AttachTerrain(opt.World, opt.Map)
	warn := opt.Warn
	if warn == nil {
		warn = func(string) {}
	}
	forms, err := tactical.LoadFormations(opt.Dir + "/KI.EXE")
	if err != nil {
		return nil, nil, fmt.Errorf("載不到陣形表：%w", err)
	}
	p := &Provider{}
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
		// 能力分類為0／1／2；腳本另依玩家攻守選四槽（spec59）。
		// ⛔ 不可以拿大地圖圖塊代替（`state.TacticalSetup.Category` 的說明）。
		Category: func(field int, siege bool) int {
			return battle.Category(field)
		},
		Script: func(field int, siege bool, tactic int, playerDefends bool) []byte {
			if p.lib == nil {
				return nil
			}
			return p.lib.Script(tactic, battle.ScriptSlot(field, playerDefends))
		},
	}
	return p, setup, nil
}

// Library 是戰場地形與腳本，載不到就是 nil。
func (p *Provider) Library() *battle.Library { return p.lib }

// Sprites 是人物圖形，載不到就是 nil（兵畫成色點）。
func (p *Provider) Sprites() *battle.Sprites { return p.sp }

// AttachTerrain 把已解碼地圖掛到目前世界，也供讀檔及自動戰鬥使用。
func AttachTerrain(w *state.World, m *world.Map) {
	if w == nil {
		return
	}
	if m == nil {
		w.SetTerrain(nil)
		return
	}
	w.SetTerrain(func(x, y int) byte { tile, _ := m.Tile(x, y); return tile })
}

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

// BuildField 載入規則層已選定的戰場，不重新選場或消費亂數（spec273）。
func (p *Provider) BuildField(n int, siege, rotate bool) *tactical.Field {
	if p.lib == nil || n < 0 || n >= battle.NumFields {
		return SyntheticField(siege)
	}
	// 用原始圖塊值建，不只用堆疊高度——城壁與門是從圖塊值認出來的，
	// 而且打壞時要換圖塊再重算高度（docs/re/11 §5.9）。
	// 再加上七層子圖塊表：兩個平面的地面圖與登城都靠它算（docs/spec/36）。
	// 規則格號包含表頭與尾段，與士兵的原版Y一致（spec36）。
	tiles, gate := p.lib.MinimapRows(n), p.lib.GateX(n)
	if rotate {
		copy(tiles[1:len(tiles)-1], battle.Rotate180(tiles[1:len(tiles)-1]))
		gate = battle.RotateGateX(gate)
	}
	return tactical.NewFieldFromTileLayers(
		tiles, p.lib.Heights(n), p.lib.TileLayers(n), gate)
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
			w.ClearMarchRoute(opt.Defender)
		}
		w.ClearMarchRoute(opt.Attacker)
		me.Node = node
		me.X, me.Y = w.Cities[node].X-1, w.Cities[node].Y
		me.TargetNode = node
		me.TargetX, me.TargetY = w.Cities[node].X, w.Cities[node].Y
		me.Timer = 1
		skipStandoff(me)
	} else {
		// 野戰：把敵方放在隔壁一格、目標設成我方所在的那一格——
		// 下一次輪到它移動就會撞上（遭遇條件是「同格、不同勢力」）。
		// ⚠ 座標被外力改掉，存檔記的那條 leg 就作廢了——不丟掉的話
		// 這支軍團會被行軍還原守衛凍住，遭遇一次都不會發生
		// （docs/spec/172 §4.6）。
		w.ClearMarchRoute(opt.Defender)
		foe.X, foe.Y = me.X-1, me.Y
		foe.TargetX, foe.TargetY = me.X, me.Y
		foe.TargetNode = me.Node
		foe.Timer = 1
		skipStandoff(foe)
	}
	for i := 0; i < 64 && w.PendingBattle() == nil; i++ {
		w.Tick(rng)
	}
	return w.PendingBattle() != nil
}

// skipStandoff 把對峙倒數推到最後一格，讓下一個移動拍當場結算。
//
// ⚠ **這是 fixture 的捷徑，不是規則。** 原版撞上之後要對峙 12 個軍團
// 巡迴週期（＝ 96 拍）才開打（docs/spec/175）；驗收捷徑只跑 64 拍，
// 而且多推那 96 拍會把世界帶到別的狀態——守方還可能在那段空檔裡被
// AI 判成無事可做而解散，於是「同狀態對拍」比的就不是同一個狀態了。
// 對峙倒數本身由 `internal/state` 的單元測試驗。
func skipStandoff(c *state.Corps) {
	c.Standoff, c.Countdown = true, 1
}
