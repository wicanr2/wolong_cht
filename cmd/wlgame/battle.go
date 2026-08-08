package main

// 戰場畫面。
//
// `BATTLE.MDL` 的子圖塊與 `BATTLE.SCH` 的人物圖形都已經解出來了
// （docs/formats/07 §8–§10），等角繪圖走 battleview.go；
// **沒有原版素材時退回高度圖 ＋ 色點**，幾何一樣是對的
// （64 × 62 的格、立體的層、陣形位置、鎖敵）。

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
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
	if g.view == nil {
		g.view = g.newBattleView(g.fieldNumber(p.Node, p.Mode == combat.Siege))
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
		g.view = nil
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

	amber := color.RGBA{240, 200, 120, 255}

	// 鏡頭跟著玩家那一側的大將走。
	me := b.Sides[g.battleSide()].Soldiers[0]

	if g.view != nil {
		g.drawBattleIso(screen, b, &me)
		g.drawBattleBanner(screen, b, p)
		g.drawBattleKeys(screen, b)
		return
	}

	// 沒有原版美術時的退路：由上往下的高度圖。
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

	// 城壁與門。它們不是地形而是**實體**（docs/re/11 §5.11），
	// 所以另外畫：完好的城壁畫石灰色、門畫木色，打壞的不畫（那一格已經變平地）。
	for _, st := range b.Structures {
		if st.Broken {
			continue
		}
		c := color.RGBA{200, 200, 190, 255} // 城壁
		if st.Kind == tactical.KindGate {
			c = color.RGBA{170, 120, 60, 255} // 門
		}
		for r := 0; r < st.Run; r++ {
			y := st.Y + r
			if y < top || y >= top+visRows {
				continue
			}
			vector.DrawFilledRect(screen,
				float32(fieldX+st.X*cellPx), float32(fieldY+(y-top)*cellPx),
				cellPx, cellPx, c, false)
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

	g.drawBattleBanner(screen, b, p)

	g.drawBattleKeys(screen, b)
}

// drawBattleBanner 畫上方橫幅與右側的城壁狀態。兩條繪製路徑共用。
func (g *game) drawBattleBanner(screen *ebiten.Image, b *tactical.Battle, p *state.Pending) {
	white := color.RGBA{240, 240, 230, 255}
	amber := color.RGBA{240, 200, 120, 255}
	dim := color.RGBA{150, 150, 160, 255}
	_ = p
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
	// 右側的空白欄放城壁的狀態。攻城戰打的是城壁耐久，
	// 而腳本指令 15 看的就是這個值（docs/re/11 §5.11）。
	if b.Field.IsSiege() && len(b.Structures) > 0 {
		const px = isoNativeW*isoScale + 6
		py := bannerH + 8
		g.chrome.Window(screen, px-chrome.Tile-2, py-chrome.Tile-2,
			screenW-(px-chrome.Tile-2)-2, 72, chrome.Menu)
		g.td.Draw(screen, "城壁", px, py, amber)
		min, intact := b.MinWallDurability()
		txt := fmt.Sprintf("耐久 %d", min)
		if !intact {
			txt = "已破"
		}
		g.td.Draw(screen, txt, px, py+18, white)
		broken := 0
		for _, st := range b.Structures {
			if st.Broken {
				broken++
			}
		}
		g.td.Draw(screen, fmt.Sprintf("%d／%d 段", len(b.Structures)-broken,
			len(b.Structures)), px, py+36, dim)
	}

}

// drawBattleKeys 畫底部的指令列。兩條繪製路徑共用。
func (g *game) drawBattleKeys(screen *ebiten.Image, b *tactical.Battle) {
	amber := color.RGBA{240, 200, 120, 255}
	dim := color.RGBA{150, 150, 160, 255}
	// 底部：指令列。與戰略畫面一樣用原版外框，不要黑底長條。
	const h = 32
	g.chrome.Window(screen, 0, screenH-h, screenW, h, chrome.Menu)
	if b.Done {
		g.td.Draw(screen, fmt.Sprintf("%s勝　第 %d 幀　按 Enter 回戰略畫面",
			sideLabel(b.Winner), b.Frame), chrome.Tile+4, screenH-h+9, amber)
	} else {
		g.td.Draw(screen,
			"1 陣形　2 攻擊　3 突擊　4 城壁　5 守陣　6 退卻",
			chrome.Tile+4, screenH-h+9, dim)
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
	if raw, err := os.ReadFile(dir + "/BATTLE.SCH"); err == nil {
		if sp, err := battle.ParseSprites(raw); err == nil {
			g.battleSprites = sp
		} else {
			log.Printf("⚠ %v；兵會畫成色點", err)
		}
	} else {
		log.Printf("⚠ 載不到 BATTLE.SCH（%v）；兵會畫成色點", err)
	}
	g.world.SetTactical(&state.TacticalSetup{
		Forms: forms,
		Field: func(node int, siege bool) *tactical.Field {
			return g.buildField(node, siege)
		},
		Script: func(node int, siege bool, tactic int) []byte {
			if g.battleLib == nil {
				return nil
			}
			return g.battleLib.Script(tactic, battle.Category(g.fieldNumber(node, siege)))
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
	n := g.fieldNumber(node, siege)
	if g.battleLib == nil || n < 0 || n >= battle.NumFields {
		return syntheticField(siege)
	}
	// 用原始圖塊值建，不只用堆疊高度——城壁與門是從圖塊值認出來的，
	// 而且打壞時要換圖塊再重算高度（docs/re/11 §5.9）。
	return tactical.NewFieldFromTiles(
		g.battleLib.Tiles(n), g.battleLib.Heights(n), g.battleLib.GateX(n))
}

// fieldNumber 回傳這一場用第幾張戰場：攻城就是據點編號，野戰現算。
func (g *game) fieldNumber(node int, siege bool) int {
	if siege {
		return node
	}
	return g.fieldForNode(node)
}

// fieldForNode 依大地圖的地形算出野戰要用哪一張戰場。
//
// 取樣的五格與 `sub_14B63` 一致（中心、下、左下、右下、兩格下方），
// 取樣方向用**玩家那一側軍團的朝向**（軍團記錄 `+0x08`）——原版是
// `cmp ah, [si+1] / mov al, [si+8]`，只在勢力等於玩家時才取。
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
	f, _ := battlefield.Select(g.playerHeading(), n)
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

// playerHeading 回傳玩家目前在動的那一支軍團的朝向。
//
// 原版只在「軍團的勢力 ＝ 玩家的勢力」時才拿 `+0x08`（`sub_14B63` 開頭那兩個
// `cmp ah, [si+1]`），所以戰場的取樣方向是**從玩家的視角**定的。
// 玩家沒有在移動的軍團時退回靜止（4），配對表會走預設那一支。
func (g *game) playerHeading() int {
	for i := range g.world.Corps {
		c := &g.world.Corps[i]
		if c.Alive && c.Faction == g.world.Player && c.Heading != state.HeadingStill {
			return c.Heading
		}
	}
	return state.HeadingStill
}

// demoBattle 是**驗收用**的捷徑：直接擺一場戰術戰鬥出來。
func (g *game) demoBattle(siege bool) {
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
	me := &g.world.Corps[mine]
	foe := &g.world.Corps[theirs]
	if siege {
		// 攻城：守方待在自己的城裡，攻方從隔壁一格走進去。
		// 走的是**正常的遭遇判定**——`resolveContact` 先問據點再問野戰。
		node := foe.Node
		me.Node = node
		me.X, me.Y = g.world.Cities[node].X-1, g.world.Cities[node].Y
		me.TargetNode = node
		me.TargetX, me.TargetY = g.world.Cities[node].X, g.world.Cities[node].Y
		me.Timer = 1
	} else {
		// 野戰：把敵方放在隔壁一格、目標設成我方所在的那一格——
		// 下一次輪到它移動就會撞上（遭遇條件是「同格、不同勢力」）。
		foe.X, foe.Y = me.X-1, me.Y
		foe.TargetX, foe.TargetY = me.X, me.Y
		foe.TargetNode = me.Node
		foe.Timer = 1
	}
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
