package main

// 戰場畫面。
//
// ⚠ **這裡畫的不是原版的美術。** `BATTLE.MDL` 的像素格式還沒解出來
// （docs/re/11 §6），所以地形用堆疊高度上色、兵用色點。
// **幾何是對的**（64 × 62 的格、立體的層、陣形位置、鎖敵），
// 美術是暫代的——解出像素格式之後換掉這一層就好。

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/battlefield"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// 戰場畫在地圖那一塊（480 × 368），一格 7 px。
const (
	cellPx  = 7
	fieldX  = 0
	fieldY  = bannerH
	visRows = (screenH - bannerH - 20) / cellPx
)

// battleActive 回報現在是不是在戰場畫面。
func (g *game) battleActive() bool { return g.world.PendingBattle() != nil }

// updateBattle 是戰場畫面的輸入。
//
// 說明書 4.1：「**戦闘中は絶対に時間を止められません**」——
// 所以這裡**沒有暫停鍵**，指令要在時間流動中下達。
func (g *game) updateBattle() {
	p := g.world.PendingBattle()
	b := p.Battle
	if !g.scriptsOn {
		g.attachScripts(p)
		g.scriptsOn = true
	}

	if !b.Done {
		// 六個戰術指令。編號與原版一致（docs/re/11 §5.8b）。
		for i, k := range []ebiten.Key{
			ebiten.Key1, ebiten.Key2, ebiten.Key3,
			ebiten.Key4, ebiten.Key5, ebiten.Key6,
		} {
			if pressed(k) {
				side := g.battleSide()
				b.Order(side, -1, tactical.Command(i))
				g.setEvent("下令：" + tactical.Command(i).String())
			}
		}
		// 每個畫面更新推進 speed 幀，與戰略層共用同一個倍率。
		n := g.speed
		if n < 1 {
			n = 1
		}
		for i := 0; i < n && !b.Done; i++ {
			b.Step()
		}
		return
	}
	// 打完了，按 Enter 結算回戰略層。
	if pressed(ebiten.KeyEnter) || pressed(ebiten.KeySpace) {
		if ev := g.world.ResolvePending(g.rng); ev != nil {
			g.setEvent(battleLine(g, *ev))
		}
		g.scriptsOn = false
	}
}

// battleSide 回傳玩家在這一場裡是哪一側。
func (g *game) battleSide() int {
	p := g.world.PendingBattle()
	if g.world.Corps[p.Attacker].Faction == g.world.Player {
		return 0
	}
	return 1
}

func (g *game) drawBattle(screen *ebiten.Image) {
	p := g.world.PendingBattle()
	b := p.Battle

	screen.Fill(color.RGBA{18, 22, 18, 255})
	vector.DrawFilledRect(screen, 0, 0, screenW, bannerH, color.RGBA{32, 24, 16, 255}, false)

	white := color.RGBA{240, 240, 230, 255}
	amber := color.RGBA{240, 200, 120, 255}
	dim := color.RGBA{150, 150, 160, 255}

	// 鏡頭跟著玩家那一側的大將走，否則 62 列塞不進畫面。
	me := b.Sides[g.battleSide()].Soldiers[0]
	top := me.Y - visRows/2
	if top < 0 {
		top = 0
	}
	if top > tactical.Height-visRows {
		top = tactical.Height - visRows
	}

	// 地形：堆疊越高越亮。堆疊 ≥ 4 的在攻城圖上就是城牆（docs/re/11 §4.3）。
	for row := 0; row < visRows; row++ {
		y := top + row
		for x := 0; x < tactical.Width; x++ {
			h := b.Field.StandLevel(x, y)
			if h == 0 {
				continue
			}
			v := uint8(40 + h*24)
			vector.DrawFilledRect(screen,
				float32(fieldX+x*cellPx), float32(fieldY+row*cellPx),
				cellPx, cellPx, color.RGBA{v, v, uint8(30 + h*10), 255}, false)
		}
	}

	// 兵。攻方暖色、守方冷色；大將畫大一格。
	for i := range b.Sides {
		base := color.RGBA{230, 120, 90, 255}
		if i == 1 {
			base = color.RGBA{110, 170, 240, 255}
		}
		for k := range b.Sides[i].Soldiers {
			s := &b.Sides[i].Soldiers[k]
			if !s.Alive || s.Y < top || s.Y >= top+visRows {
				continue
			}
			px := float32(fieldX + s.X*cellPx)
			py := float32(fieldY + (s.Y-top)*cellPx)
			c := base
			if s.Cmd == tactical.Retreat {
				c = color.RGBA{120, 120, 120, 255} // 退卻中畫成灰的
			}
			size := float32(cellPx - 2)
			if s.IsGeneral() {
				size = cellPx
				c = amber
				if i == 1 {
					c = color.RGBA{200, 220, 255, 255}
				}
			}
			vector.DrawFilledRect(screen, px, py, size, size, c, false)
		}
	}

	// 上方橫幅：雙方的兵數與有利／不利。
	g.td.Draw(screen, "戰場", 8, 8, amber)
	for i := range b.Sides {
		s := &b.Sides[i]
		label := "攻方"
		if i == 1 {
			label = "守方"
		}
		g.td.Draw(screen, fmt.Sprintf("%s %4d 兵　%s　大將體力 %3d",
			label, s.Remaining(), advName(b.Advantage[i]), s.Soldiers[0].HP),
			90+i*260, 8, white)
	}

	// 底部：指令列。
	vector.DrawFilledRect(screen, 0, screenH-19, screenW, 19, color.RGBA{0, 0, 0, 210}, false)
	if b.Done {
		g.td.Draw(screen, fmt.Sprintf("%s勝　第 %d 幀　按 Enter 回戰略畫面",
			sideLabel(b.Winner), b.Frame), 4, screenH-17, amber)
	} else {
		g.td.Draw(screen,
			"1 陣形　2 攻擊　3 突擊　4 城壁　5 守陣　6 退卻　（時間不會停）",
			4, screenH-17, dim)
	}
}

func advName(a tactical.Advantage) string {
	switch a {
	case tactical.Advantaged:
		return "有利"
	case tactical.Disadvantaged:
		return "不利"
	}
	return "普通"
}

func sideLabel(i int) string {
	if i == 0 {
		return "攻方"
	}
	return "守方"
}

// ---------------------------------------------------------------------------
// 戰場來源
// ---------------------------------------------------------------------------

// installTactical 把戰場來源裝到世界上。
//
// 三份資料各自可缺：陣形表在 `KI.EXE`、地形在 `BATTLE.MAP` ＋ `BATTLE.MDL`、
// AI 腳本在 `BATTLE.DAT`。**少一份就退一級**，不會讓遊戲開不起來——
// 最差的情況是玩家的戰鬥也走自動判定。
func (g *game) installTactical(dir string) {
	forms, err := tactical.LoadFormations(dir + "/KI.EXE")
	if err != nil {
		log.Printf("⚠ 載不到陣形表（%v）；玩家的戰鬥會走自動判定", err)
		return
	}
	lib := loadBattleLibrary(dir)
	g.battleLib = lib
	g.world.SetTactical(&state.TacticalSetup{
		Forms: forms,
		Field: func(node int, siege bool) *tactical.Field {
			return g.buildField(node, siege)
		},
	})
}

func loadBattleLibrary(dir string) *battle.Library {
	read := func(name string) []byte {
		b, err := os.ReadFile(dir + "/" + name)
		if err != nil {
			log.Printf("⚠ 載不到 %s（%v）", name, err)
			return nil
		}
		return b
	}
	m, mdl := read("BATTLE.MAP"), read("BATTLE.MDL")
	if m == nil || mdl == nil {
		log.Print("⚠ 沒有戰場地形，會用生成的替代")
		return nil
	}
	lib, err := battle.Parse(m, mdl, read("BATTLE.DAT"))
	if err != nil {
		log.Printf("⚠ %v；會用生成的地形", err)
		return nil
	}
	return lib
}

// buildField 取一張戰場。
//
//   - **攻城戰**：戰場編號就是據點編號（`docs/re/05`）
//   - **野戰**：從大地圖上即時算（`internal/rules/battlefield`）——
//     取軍團所在格與下方四格的地形類型去配一張 21 筆的表
func (g *game) buildField(node int, siege bool) *tactical.Field {
	n := node
	if !siege {
		n = g.fieldForNode(node)
	}
	if g.battleLib == nil || n < 0 || n >= battle.NumFields {
		return syntheticField(siege)
	}
	return tactical.NewField(g.battleLib.Stacks(n), g.battleLib.GateX(n))
}

// fieldForNode 依大地圖的地形算出野戰要用哪一張戰場。
//
// 取樣的五格與 `sub_14B63` 一致（中心、下、左下、右下、兩格下方）。
// ⚠ **行進方向先固定用 2**——原版讀的是軍團記錄 `+0x08`，
// 而那個欄位在本專案的軍團結構裡還沒接出來。
func (g *game) fieldForNode(node int) int {
	if g.lib == nil || g.lib.World == nil {
		return battlefield.FieldBase + 6
	}
	// 野外的節點編號沒有座標，只有據點有；用據點的格座標取樣。
	if node < 0 || node >= len(g.world.Cities) {
		return battlefield.FieldBase + 6
	}
	cx, cy := g.world.Cities[node].X, g.world.Cities[node].Y
	at := func(dx, dy int) int {
		t, err := g.lib.World.Tile(cx+dx, cy+dy)
		if err != nil {
			return 0
		}
		return battlefield.Terrain(t)
	}
	n := battlefield.Neighbours{
		Centre:    at(0, 0),
		Down:      at(0, 1),
		DownLeft:  at(-1, 1),
		DownRight: at(1, 1),
		TwoDown:   at(0, 2),
	}
	f, _ := battlefield.Select(2, n)
	if f >= 209 && f <= 213 {
		f = battlefield.SelectWater(f-battlefield.TerrainBase, g.rng.Next())
	}
	return f
}

// syntheticField 是沒有原版戰場資料時的替代品。幾何同尺寸，內容是自己生的。
func syntheticField(siege bool) *tactical.Field {
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

// attachScripts 讓**不是玩家的那一側**由原版的 AI 腳本驅動。
//
// 段編號 ＝ 帶兵武將的 `+0x16` × 4 ＋ 戰場類別（`sub_1CBE5`）。
func (g *game) attachScripts(p *state.Pending) {
	if g.battleLib == nil {
		return
	}
	field := p.Node
	if p.Mode != combat.Siege {
		field = g.fieldForNode(p.Node)
	}
	cat := battle.Category(field)
	for side, corps := range [2]int{p.Attacker, p.Defender} {
		if corps < 0 || g.world.Corps[corps].Faction == g.world.Player {
			continue // 玩家那一側留給玩家下命令
		}
		code := g.battleLib.Script(g.world.Generals[corps].Tactic, cat)
		if code != nil {
			p.Battle.SetScript(side, tactical.NewScript(code, side))
		}
	}
}

// demoBattle 是**驗收用**的捷徑：直接擺一場戰術戰鬥出來。
func (g *game) demoBattle() {
	p := g.world.Player
	var mine, theirs int = -1, -1
	for i, gen := range g.world.Generals {
		if !gen.Alive || gen.Posted {
			continue
		}
		if gen.Faction == p && mine < 0 {
			mine = i
		}
		if gen.Faction != p && gen.Faction < 22 && theirs < 0 {
			theirs = i
		}
	}
	if mine < 0 || theirs < 0 {
		return
	}
	kinds := [6]army.TroopType{
		army.Cavalry, army.Cavalry, army.Archer,
		army.Archer, army.Infantry, army.Infantry,
	}
	manned := [6]bool{true, true, true, true, true, true}
	for _, l := range []int{mine, theirs} {
		f := g.world.Generals[l].Faction
		g.world.Factions[f].Reserves = [3]int{9000, 9000, 9000}
		if err := g.world.FormCorps(l, kinds, manned); err != nil {
			g.setEvent(err.Error())
			return
		}
	}
	// 把敵方放在隔壁一格、目標設成我方所在的那一格——
	// 下一次輪到它移動就會撞上（原版的野戰遭遇條件是「同格、不同勢力」）。
	me := &g.world.Corps[mine]
	foe := &g.world.Corps[theirs]
	foe.X, foe.Y = me.X-1, me.Y
	foe.TargetX, foe.TargetY = me.X, me.Y
	foe.TargetNode = me.Node
	foe.Timer = 1
	for i := 0; i < 64 && !g.battleActive(); i++ {
		g.world.Tick(g.rng)
	}
	if g.battleActive() {
		// 先跑一段，截圖時才看得到陣線接觸；速度調 1 免得截圖前就打完。
		for i := 0; i < 900; i++ {
			g.world.PendingBattle().Battle.Step()
		}
		g.speed = 1
	}
}
