package main

// 戰場畫面。
//
// `BATTLE.MDL` 的子圖塊與 `BATTLE.SCH` 的人物圖形都已經解出來了
// （docs/formats/07 §8–§10），等角繪圖走 battleview.go；
// **沒有原版素材時退回高度圖 ＋ 色點**，幾何一樣是對的
// （64 × 62 的格、立體的層、陣形位置、鎖敵）。

import (
	"strings"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/battlesetup"
	"github.com/wicanr2/wolong_cht/internal/rules/speed"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// battleActive 回報現在是不是在戰場畫面。
func (g *game) battleActive() bool { return g.world.PendingBattle() != nil }

// updateBattleChoice 是行軍遭遇後的「戰鬥指揮／委任」選單。
// 這個選單不是除錯捷徑：它只消費正常行軍觸發的 EncounterChoice，
// 在玩家決定前不會讓戰略時鐘繼續走。
func (g *game) updateBattleChoice() {
	x, y := ebiten.CursorPosition()
	if row, ok := battleChoiceRowAt(x, y); ok {
		// Hover 只改變反白列；只有 JustPressed 才能確認。
		g.battleChoiceRow = row
	}
	if pressed(ebiten.KeyArrowUp) || pressed(ebiten.KeyArrowDown) {
		g.battleChoiceRow = 1 - g.battleChoiceRow
	}
	if pressed(ebiten.Key1) {
		g.battleChoiceRow = 0
	}
	if pressed(ebiten.Key2) {
		g.battleChoiceRow = 1
	}
	keyboardConfirm := pressed(ebiten.KeyEnter) || pressed(ebiten.KeySpace) ||
		pressed(ebiten.Key1) || pressed(ebiten.Key2)
	mouseConfirm := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	if !keyboardConfirm && !mouseConfirm {
		return
	}
	if mouseConfirm {
		// 畫面外、框線與列距空白都不能觸發既有確認路徑。
		if _, ok := battleChoiceRowAt(x, y); !ok {
			return
		}
	}

	if g.battleChoiceRow == 0 {
		if err := g.world.ChooseBattleCommand(); err != nil {
			g.setEvent(err.Error())
			return
		}
		g.view = nil
		g.setEvent("戰鬥指揮")
		return
	}
	ev := g.world.ChooseBattleDelegate(g.rng)
	if ev == nil {
		g.setEvent("沒有待處理的遭遇")
		return
	}
	g.setEvent("委任：" + battleLine(g, *ev))
}

// drawBattleChoice 畫出原版行軍抵達後的兩路選擇。
func (g *game) drawBattleChoice(screen *ebiten.Image, c *state.EncounterChoice) {
	l := battleChoiceLayoutFor()
	g.chrome.Window(screen, l.Window.X, l.Window.Y, l.Window.W, l.Window.H, chrome.Menu)
	amber := color.RGBA{240, 200, 120, 255}
	white := chrome.Paper
	dim := color.RGBA{170, 170, 180, 255}
	selected := color.RGBA{140, 230, 140, 255}
	x, y, h := l.Window.X, l.Window.Y, l.Window.H

	att := big5(g.world.Generals[c.Attacker].Name)
	def := "城兵"
	if c.Defender >= 0 {
		def = big5(g.world.Generals[c.Defender].Name)
	}
	g.td.Draw(screen, "遭遇戰", x+chrome.Tile+4, y+chrome.Tile+2, amber)
	g.td.Draw(screen, att+" 對 "+def, x+chrome.Tile+4,
		y+chrome.Tile+textdraw.GlyphH+6, white)
	mode := "野戰"
	if c.Mode == combat.Siege {
		mode = "攻城"
	}
	g.td.Draw(screen, mode+"　請選擇處理方式", x+chrome.Tile+4,
		y+chrome.Tile+2*(textdraw.GlyphH+4), dim)

	rows := []string{"戰鬥指揮", "委任"}
	for i, row := range rows {
		col := white
		mark := "　"
		if i == g.battleChoiceRow {
			col, mark = selected, "●"
		}
		g.td.Draw(screen, fmt.Sprintf("%s%d　%s", mark, i+1, row),
			l.Rows[i].X, l.Rows[i].Y, col)
	}
	g.td.Draw(screen, "↑↓ 選擇　1-2 直選　Enter 確定",
		x+chrome.Tile+4, y+h-chrome.Tile-textdraw.GlyphH, dim)
}

// speedToastFrames 是調速度之後那行提示顯示幾幀（約 1.5 秒 @60fps）。
//
// **常駐顯示會破壞版面 parity**（原版戰場沒有速度指示），所以只在剛調過
// 的時候浮一下。這是 remake 差異。
const speedToastFrames = 90

// drawSpeedToast 在戰場左下角浮一行「戰術速度 N」。
func (g *game) drawSpeedToast(dst *ebiten.Image, l dosvBattleLayout) {
	if g.speedToast <= 0 {
		return
	}
	text := "戰術速度" + speed.Labels[clamp(g.tacticalSpeed, 0, speed.Levels-1)]
	w := g.td.Width(text) + 16
	x := l.Field.X + 8
	y := l.BottomCommands.Y - textdraw.GlyphH - 12
	vector.DrawFilledRect(dst, float32(x), float32(y-4), float32(w),
		float32(textdraw.GlyphH+8), color.RGBA{0, 0, 0, 200}, false)
	g.td.Draw(dst, text, x+8, y, chrome.Paper)
}

// updateBattle 是戰場畫面的輸入。
//
// 說明書 4.1：「**戦闘中は絶対に時間を止められません**」——
// 所以這裡**沒有暫停鍵**，指令要在時間流動中下達。
func (g *game) updateBattle() {
	p := g.world.PendingBattle()
	b := p.Battle
	l := dosvBattleLayoutFor(screenW, screenH)
	if g.view == nil {
		g.view = g.newBattleView(g.battle.FieldNumber(p.Node, p.Mode == combat.Siege))
	}

	// ＋／− 調**戰術速度**。原版的戰術速度是獨立設定（系統選單第 5 列），
	// 而戰場畫面獨佔輸入，所以這裡要自己接一次——不然進了戰場就調不到。
	for i, k := range []ebiten.Key{ebiten.KeyMinus, ebiten.KeyEqual} {
		if pressed(k) {
			g.adjustSpeed(true, []int{-1, 1}[i])
			g.speedToast = speedToastFrames
		}
	}
	if g.speedToast > 0 {
		g.speedToast--
	}

	if !b.Done {
		g.startBattleTalk(p)
		// 右鍵先過熱區表（docs/spec/32 §2.1）。放在對白推進之前，
		// 因為原版的右鍵分派（0x1C01D）不看對白狀態。
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			if x, y := ebiten.CursorPosition(); handleBattleRightClick(b, x, y) {
				return
			}
		}
		talkAdvanced := g.advanceBattleTalkInput()
		// 開場 50 tick ＋ 整段單挑是**不可打斷的 modal**（原版全在
		// `sub_1A1C5` 裡，主迴圈的輸入分派還沒開始，docs/spec/80 §3）。
		if b.OpeningActive() {
			talkAdvanced = true
		}
		// 六個戰術指令。編號與原版一致（docs/re/11 §5.8b）。
		if !talkAdvanced {
			for i, k := range []ebiten.Key{
				ebiten.Key1, ebiten.Key2, ebiten.Key3,
				ebiten.Key4, ebiten.Key5, ebiten.Key6,
			} {
				if pressed(k) {
					g.issueBattleCommand(b, i)
				}
			}
			if x, y := ebiten.CursorPosition(); inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if g.view != nil && l.SideMiniMap.containsPoint(x, y) {
					g.view.SetCameraFromMiniMap(x, y)
					return
				}
				// 底列六格是**選部隊**，不是命令（docs/spec/33）。
				if slot, ok := splitBattleCommandIndexAt(l.BottomCommands, x, y); ok {
					g.toggleBattleSquad(b, slot)
					return
				}
				// 陣形選單：8 欄 × 2 列 ＝ 十六個陣形（docs/spec/37 §2.1）。
				if idx, ok := battleFormationIndexAt(l.SideFormation, x, y); ok {
					g.setBattleFormation(b, idx)
					return
				}
				// 陣形線三選一（同上 §2.2）：由上而下 ＝ 敵軍側／中央／自軍側。
				if k, ok := battleLineIndexAt(l.SideLines, x, y); ok {
					// 三格改的是**玩家自己**的陣形線，所以永遠取原版
					// 「玩家」那一組（`LineFor` 的 side 0），與玩家站
					// 攻方還是守方無關（docs/spec/56 §1）。
					b.Sides[g.battleSide()].Line = tactical.LineFor(0, 2-k)
					return
				}
				// 側欄面板：畫面列序不是命令碼順序，要查表。
				if row, ok := battleSideCommandIndexAt(l.SideCommands, x, y); ok {
					g.issueBattleCommand(b, battleSideCommandRowCode[row])
				}
			}
		}
		// 這一個畫面推進幾個戰場幀。**戰術速度是獨立設定，而且值要 ×16**
		// ——原版第 5 列的 handler `sub_160A5` 就做這一件事
		// （docs/re/61 §4、docs/spec/34）。
		n := g.tacticalThrottle.Steps(g.tacticalSpeed, speed.TacticalMul, speed.HighSpeedTactical)
		for i := 0; i < n && !b.Done; i++ {
			b.Step()
		}
		// 規則層把原版的效果碼排隊，播放在這裡（docs/re/17 §3）。
		// 碼就是 `SOUND.DAT` 的記錄編號，不必另外對照表。
		for _, code := range b.TakeSoundEffects() {
			g.sound.PlayEffect(int(code))
		}
		// 這一幀新產生的單挑喊話當幀就掛框，不等下一幀。
		g.pumpDuelTalks(p)
		g.tickBattleTalk(n)
		return
	}
	clearBattleTalkSession(g, b)
	// 打完了，按 Enter 或**點一下滑鼠**結算回戰略層（docs/spec/89 §3）。
	//
	// ⚠ 用 `JustPressed` 不是 `Pressed`：戰鬥是用滑鼠下令的，
	// 最後一個指令那一下如果還按著，用「持續按著」判定會讓結果畫面
	// 在出現的同一幀就被關掉——玩家根本看不到它。
	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle)
	if pressed(ebiten.KeyEnter) || pressed(ebiten.KeySpace) || clicked {
		if ev := g.world.ResolvePending(g.rng); ev != nil {
			g.setEvent(battleLine(g, *ev))
		}
		g.view = nil
	}
}

func battleCommandIndexForKey(k ebiten.Key) (int, bool) {
	for i, key := range [...]ebiten.Key{
		ebiten.Key1, ebiten.Key2, ebiten.Key3,
		ebiten.Key4, ebiten.Key5, ebiten.Key6,
	} {
		if k == key {
			return i, true
		}
	}
	return 0, false
}

// issueBattleCommand 走玩家那一條下令路（docs/spec/33 §1.3）：
// 命令只送給被選中的隊，一隊都沒選就送給全隊，送完把選取清空。
func (g *game) issueBattleCommand(b *tactical.Battle, command int) {
	if b == nil || b.Done || command < 0 || command >= len(battleCommandLabels) {
		return
	}
	side := g.battleSide()
	order := tactical.Command(command)
	if !b.OrderSelected(side, order) {
		// 原版跳 TALK 582「這哪裡有城壁啊！！」並且不下令。
		g.setEvent("這哪裡有城壁啊！！")
		return
	}
	g.setEvent("下令：" + order.String())
}

// toggleBattleSquad 切換底列第 slot 格對應那一隊的選取（熱區 0x15–0x1A）。
func (g *game) toggleBattleSquad(b *tactical.Battle, slot int) {
	if b == nil || b.Done || slot < 0 || slot >= len(battleBottomSlotSquad) {
		return
	}
	b.ToggleSquadSelection(g.battleSide(), battleBottomSlotSquad[slot])
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
	l := dosvBattleLayoutFor(screenW, screenH)

	screen.Fill(color.RGBA{18, 22, 18, 255})

	// 鏡頭跟著玩家那一側的大將走。
	me := b.Sides[g.battleSide()].Soldiers[0]

	if g.view != nil {
		g.drawBattleIso(screen, b, &me)
		g.drawBattleChrome(screen, b, p, l)
		g.drawSpeedToast(screen, l)
		g.drawBattleResult(screen, b, p)
		return
	}

	// 沒有原版美術時的退路：由上往下的高度圖。仍遵守同一個完整
	// 戰術版面，避免資產缺失時又退回右側大片空黑。
	amber := color.RGBA{240, 200, 120, 255}
	cellPx := 7
	fieldX, fieldY := l.Field.X, l.Field.Y
	visRows := l.Field.H / cellPx
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
			if fieldX+x*cellPx >= l.Field.right() {
				break
			}
			h := b.Field.StandLevel(x, y)
			if h == 0 {
				continue
			}
			v := uint8(40 + h*24)
			vector.DrawFilledRect(screen,
				float32(fieldX+x*cellPx), float32(fieldY+row*cellPx),
				float32(cellPx), float32(cellPx), color.RGBA{v, v, uint8(30 + h*10), 255}, false)
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
			if fieldX+st.X*cellPx >= l.Field.right() {
				continue
			}
			vector.DrawFilledRect(screen,
				float32(fieldX+st.X*cellPx), float32(fieldY+(y-top)*cellPx),
				float32(cellPx), float32(cellPx), c, false)
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
			if fieldX+s.X*cellPx >= l.Field.right() {
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
				size = float32(cellPx)
				c = amber
				if i == 1 {
					c = color.RGBA{200, 220, 255, 255}
				}
			}
			vector.DrawFilledRect(screen, px, py, size, size, c, false)
		}
	}

	g.drawBattleChrome(screen, b, p, l)
	g.drawBattleResult(screen, b, p)
}

// drawBattleResult 是戰術戰鬥完成後、回寫戰略層之前的明確結果畫面。
// 這個停留點對應原版「戰鬥結束後仍先顯示結果，玩家確認才返回」的
// 玩家路徑；按 Enter 仍由 updateBattle 呼叫 ResolvePending。
func (g *game) drawBattleResult(screen *ebiten.Image, b *tactical.Battle, p *state.Pending) {
	if !b.Done {
		return
	}
	o := b.Result()
	const x, y, w, h = 72, 104, 496, 160
	g.chrome.Window(screen, x, y, w, h, chrome.Menu)
	amber := color.RGBA{240, 200, 120, 255}
	white := chrome.Paper
	dim := color.RGBA{170, 170, 180, 255}
	winner := "攻方勝"
	if !o.AttackerWins {
		winner = "守方勝"
	}
	g.td.Draw(screen, "戰鬥結果　"+winner, x+chrome.Tile+4,
		y+chrome.Tile+2, amber)
	attBefore := g.world.Corps[p.Attacker].Men * tactical.MenPerSoldier
	defName := "守方"
	defBefore := 0
	if p.Defender >= 0 {
		defName = "守方"
		defBefore = g.world.Corps[p.Defender].Men * tactical.MenPerSoldier
	} else {
		defName = "城兵"
		defBefore = p.Garrison * tactical.MenPerSoldier
	}
	g.td.Draw(screen, fmt.Sprintf("攻方兵力　%d → %d", attBefore, o.Men[0]),
		x+chrome.Tile+4, y+chrome.Tile+textdraw.GlyphH+8, white)
	g.td.Draw(screen, fmt.Sprintf("%s兵力　%d → %d", defName, defBefore, o.Men[1]),
		x+chrome.Tile+4, y+chrome.Tile+2*(textdraw.GlyphH+8), white)
	// ⭐ **原版戰術結束後沒有損害報告**（使用者裁定 2026-08-26，
	// docs/spec/89）。這一行是 remake 的驗收資訊，預設關著，
	// 要看的人在系統選單第 8 列打開。
	if p.Mode == combat.Siege && g.damageReport {
		g.td.Draw(screen, fmt.Sprintf("攻城損害　%d", b.CityDamage(p.CityWall)),
			x+chrome.Tile+4, y+chrome.Tile+3*(textdraw.GlyphH+8), dim)
	}
	g.td.Draw(screen, "Enter 或點一下　回到戰略畫面", x+chrome.Tile+4,
		y+h-chrome.Tile-textdraw.GlyphH, dim)
}

// drawBattleChrome 畫戰術畫面本身的完整 DOS/V 編排。
//
// 原版不是「上方一條文字 banner ＋ 左邊一張圖」；主將 TALK、雙方狀態、
// 縮圖、命令區與底部 glyph row 都是戰場畫面的一部分。TALK 的 producer
// 不在本批修改範圍內，因此沒有 payload 時只顯示可診斷的版面 fallback。
func (g *game) drawBattleChrome(screen *ebiten.Image, b *tactical.Battle, p *state.Pending, l dosvBattleLayout) {
	g.drawBattleSidebar(screen, b, p, l)
	g.drawBattleGateBar(screen, b)
	talks := g.battleTalkState(b, p)
	if talks.visible(0) {
		g.drawBattleTalk(screen, talks.text(0), talks.speaker(0), l.TopTalk, 0, talks.portrait(0), p)
	}
	if talks.visible(1) {
		g.drawBattleTalk(screen, talks.text(1), talks.speaker(1), l.BottomTalk, 1, talks.portrait(1), p)
	}
	g.drawBattleKeys(screen, b, l.BottomCommands)
}

// battleTalkState 只呈現目前 queue 的一筆。未知 marker 已在 enqueue 時
// fail-closed；此處仍保留同一個安全閘，避免錯誤資料進入 renderer。
func (g *game) battleTalkState(b *tactical.Battle, p *state.Pending) battleTalkState {
	if g == nil || b == nil || p == nil {
		return battleTalkState{}
	}
	g.startBattleTalk(p)
	s := &g.battleTalkSession
	if s.battle != b || !s.queue.active() {
		return battleTalkState{}
	}
	// **兩側各一個框，可以同時掛著**（原版 word_1D322／word_1D324，
	// docs/spec/60）——上框是側 0、下框是側 1。
	result := battleTalkState{}
	for side := 0; side < 2; side++ {
		entry, ok := s.queue.current(side)
		if !ok {
			continue
		}
		text, ok := g.battleTalkText(entry)
		if !ok {
			s.queue.clear(side)
			continue
		}
		if side == 0 {
			result.Top, result.TopPortrait, result.TopSpeaker = text, entry.Portrait, entry.Speaker
		} else {
			result.Bottom, result.BottomPortrait, result.BottomSpeaker = text, entry.Portrait, entry.Speaker
		}
	}
	return result
}

func (g *game) drawBattleTalk(screen *ebiten.Image, text, speaker string, r battleRect, side, portraitPage int, p *state.Pending) {
	if text == "" {
		return
	}
	// 框內是**純黑**、沒有龍紋，也沒有獨立的姓名列——說話者的名字由
	// \1 直接代在句中，畫成色 9（playtest/43 的實機逐像素）。
	g.chrome.Window(screen, r.X, r.Y, r.W, r.H, chrome.Blank)
	white := g.paletteInk(strategyInkNormal, chrome.Paper)
	nameInk := g.paletteInk(strategyInkDim, color.RGBA{243, 211, 146, 255})

	commander := g.battleCommander(p, side)
	if commander >= 0 && commander < len(g.world.Generals) && g.lib != nil && portraitPage >= 0 {
		if portrait, err := g.lib.Portrait(portraitPage, int(g.world.Clock.Season())); err == nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(r.X+chrome.Tile), float64(r.Y+chrome.Tile))
			screen.DrawImage(ebiten.NewImageFromImage(portrait), op)
		}
	}

	tx := r.X + chrome.Tile + 72
	// 文字從框頂下 16 px 起、**行距 16**（⚠ textdraw.GlyphH 是 15，
	// 拿它當行距會整行高 1px——實機量測 m23：行1 墨水 16 起、行2 32 起）。
	ty := r.Y + 16
	maxWidth := r.W - (tx - r.X) - chrome.Tile - 4
	for i, line := range textdraw.WrapLine(text, maxWidth) {
		if i >= 3 {
			break
		}
		g.drawTalkLineWithName(screen, line, speaker, tx, ty+i*16, white, nameInk)
	}
}

// drawTalkLineWithName 畫一行對白，句中等於 speaker 的那一段換色
// （原版把 \1 代入的名字畫成色 9）。
func (g *game) drawTalkLineWithName(screen *ebiten.Image, line, speaker string, x, y int, ink, nameInk color.RGBA) {
	if speaker == "" {
		g.td.Draw(screen, line, x, y, ink)
		return
	}
	rest := line
	for {
		idx := strings.Index(rest, speaker)
		if idx < 0 {
			g.td.Draw(screen, rest, x, y, ink)
			return
		}
		if idx > 0 {
			g.td.Draw(screen, rest[:idx], x, y, ink)
			x += textdraw.StringWidth(rest[:idx])
		}
		g.td.Draw(screen, speaker, x, y, nameInk)
		x += textdraw.StringWidth(speaker)
		rest = rest[idx+len(speaker):]
	}
}

func (g *game) battleCommander(p *state.Pending, side int) int {
	corps := p.Attacker
	if side == 1 {
		corps = p.Defender
	}
	if corps < 0 || corps >= len(g.world.Corps) || !g.world.Corps[corps].Alive {
		return -1
	}
	return g.world.Leader(corps)
}

func (g *game) drawBattleSidebar(screen *ebiten.Image, b *tactical.Battle, p *state.Pending, l dosvBattleLayout) {
	vector.DrawFilledRect(screen, float32(l.Sidebar.X), float32(l.Sidebar.Y),
		float32(l.Sidebar.W), float32(l.Sidebar.H), chrome.Menu, false)
	// 上格是對方、下格是我方（docs/spec/31 §1）。原版靠互換
	// word_10D2E／word_10D30 做到這件事；這裡直接用玩家所在的側別，
	// 效果相同而且不必動 tactical 的側別編號。
	ally := g.battleSide()
	foe := 1 - ally
	g.drawBattleSideTitle(screen, p, l.SideTitle, ally, foe)
	g.drawBattleSideCell(screen, b, p, l.SideFoe, foe, true)
	g.drawBattleMiniMap(screen, b, l.SideMiniMap)
	g.drawBattleSideCell(screen, b, p, l.SideAlly, ally, false)
	g.drawBattleFormationStrip(screen, b, l.SideFormation)
	g.drawBattleSideCommands(screen, b, l.SideCommands)
	g.drawBattleSideFooter(screen, l.SideFooter)
	g.drawBattleSideFrame(screen, l.Sidebar)
}

// drawBattleSideFrame 貼側欄的兩根柱子、四個角與四條橫帶（docs/spec/31 §1.1）。
//
// 原版是 `sub_1CA3B` 一塊一塊貼圖，不是畫線；少了它側欄兩側是空的。
// 畫在最後：橫帶壓在七格的邊界上，與原版的繪製順序一致。
func (g *game) drawBattleSideFrame(screen *ebiten.Image, r battleRect) {
	blit := func(img *ebiten.Image, x, y int) {
		if img == nil {
			return
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(img, op)
	}
	left, right := r.X, r.right()-16
	for y := 0; y < r.H; y += 16 { // 25 × 16 ＝ 400
		blit(g.battleFrame[library.BattleFrameLeft], left, r.Y+y)
		blit(g.battleFrame[library.BattleFrameRight], right, r.Y+y)
	}
	for _, y := range battleSideBandY {
		for i := 0; i < 8; i++ { // 8 × 16 ＝ 128
			blit(g.battleFrame[library.BattleFrameBand], r.X+16+i*16, r.Y+y)
		}
	}
	for _, y := range [2]int{0, r.H - 8} {
		blit(g.battleFrame[library.BattleFrameCorner], left, r.Y+y)
		blit(g.battleFrame[library.BattleFrameCorner], right, r.Y+y)
	}
}

// battleSideBandY 是四條橫帶的 y（docs/spec/31 §1.1）——七格的分隔線。
var battleSideBandY = [4]int{0, 40, 240, 392}

// drawBattleSideTitle 畫標題兩行（docs/re/60 §3）：
//
//	(512,8)  戰場地名     (560,8)  「作戰」
//	(512,24) 我方君主     (560,24) 「對」    (576,24) 對方君主
//
// 原版每一段都是固定三個全形字（48 px），所以三個 x 是 512／560／576。
func (g *game) drawBattleSideTitle(screen *ebiten.Image,
	p *state.Pending, r battleRect, ally, foe int) {
	// ⭐ 標題格的底是**色 0（黑）**：原版先 sub_1F1A3 填 (496,8)–(623,39)
	// 再寫字（docs/spec/31 §1.1）。畫成視窗底色會讓這一格整片不同。
	vector.DrawFilledRect(screen, float32(r.X), float32(r.Y),
		float32(r.W), float32(r.H), chrome.Blank, false)
	place := g.battleTitlePlace // 調色盤索引 9
	lord := g.battleTitleLord   // 調色盤索引 11
	x0 := r.X + 16              // 512 − 496
	if uiLang.Lang().Latin() {
		// ⚠ **三段擠不進 2×16 個半形格**（地名 12、君主名各 11），
		// 照原版的三個 x 畫會互相蓋掉（`SHOUCHBNttle`）。
		// 會撞名的那兩段優先：君主名裁掉會讓兩個不同的人長得一樣，
		// 地名不顯示只是少一個資訊，而戰場就畫在眼前（docs/spec/87 §5）。
		g.td.Draw(screen, g.battleLordName(p, ally), x0, r.Y, lord)
		g.td.Draw(screen, "vs "+g.battleLordName(p, foe), x0, r.Y+16, lord)
		return
	}
	g.td.Draw(screen, g.battleFieldName(p), x0, r.Y, place)
	g.td.Draw(screen, "作戰", x0+48, r.Y, place)
	g.td.Draw(screen, g.battleLordName(p, ally), x0, r.Y+16, lord)
	g.td.Draw(screen, "對", x0+48, r.Y+16, lord)
	g.td.Draw(screen, g.battleLordName(p, foe), x0+64, r.Y+16, lord)
}

// drawBattleSideCell 畫一格將旗：主將名 ＋ 兵力條 ＋ 大將體力條。
// top 為真時用上格（對方）的排列，否則用下格（我方）的（docs/re/60 §5）。
func (g *game) drawBattleSideCell(screen *ebiten.Image, b *tactical.Battle,
	p *state.Pending, r battleRect, side int, top bool) {
	slot := 0
	if top {
		slot = 1
	}
	if img := g.battleSideFlags[slot]; img != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(r.X), float64(r.Y))
		screen.DrawImage(img, op)
	} else {
		g.chrome.Window(screen, r.X, r.Y, r.W, r.H, chrome.Menu)
	}
	cell := battleSideCellLayoutFor(r, top)
	name := sideLabel(side)
	if commander := g.battleCommander(p, side); commander >= 0 && commander < len(g.world.Generals) {
		name = big5(g.world.Generals[commander].Name)
	}
	g.td.Draw(screen, name, cell.Name.X, cell.Name.Y, chrome.Paper)

	s := &b.Sides[side]
	health := 0
	if len(s.Soldiers) > 0 {
		health = s.Soldiers[0].HP
	}
	menLen, healthLen := battleSideBarLengths(s.Remaining(), health)
	// 索引 12／11 一律查調色盤（docs/spec/54）；未填的部分原版畫成色 0。
	g.drawBattleBar(screen, cell.MenBar, menLen, g.battleMenBar)
	g.drawBattleBar(screen, cell.HealthBar, healthLen, g.battleHealthBar)
}

// drawBattleBar 照 sub_10AAA：先畫 filled px 的實色，其餘畫色 0。
func (g *game) drawBattleBar(screen *ebiten.Image, r battleRect, filled int, col color.RGBA) {
	if filled > 0 {
		vector.DrawFilledRect(screen, float32(r.X), float32(r.Y),
			float32(filled), float32(r.H), col, false)
	}
	if rest := r.W - filled; rest > 0 {
		vector.DrawFilledRect(screen, float32(r.X+filled), float32(r.Y),
			float32(rest), float32(r.H), color.RGBA{0, 0, 0, 255}, false)
	}
}

// drawBattleFormationStrip 畫十六個陣形那一格（docs/re/60 §8）。
// 選取框是格內縮 1 px 的 14×14；原版選中色 12、取消色 10。
func (g *game) drawBattleFormationStrip(screen *ebiten.Image, b *tactical.Battle, r battleRect) {
	if img := g.battleFormationStrip; img != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(r.X), float64(r.Y))
		screen.DrawImage(img, op)
	} else {
		g.chrome.Window(screen, r.X, r.Y, r.W, r.H, chrome.Menu)
	}
	idx := g.battleFormation
	if idx < 0 || idx >= battleFormationCols*battleFormationRows {
		return
	}
	cell := battleFormationCellRect(r, idx)
	vector.StrokeRect(screen, float32(cell.X+1), float32(cell.Y+1), 14, 14,
		1, g.battleCommandSelect, false)
}

// setBattleFormation 選一個陣形（原版 handler `0x1C11A`）。
//
// 原版只把 `byte_1D346` 與 `word_1D342`（＝ 格 × 96）寫下去，
// **沒有立刻重排部隊**——要等下一次「陣形」命令生效才走過去。
func (g *game) setBattleFormation(b *tactical.Battle, idx int) {
	if idx < 0 || idx >= tactical.NumFormations {
		return
	}
	g.battleFormation = idx
	b.Sides[g.battleSide()].Formation = idx
}

// 門強度條的版面（docs/spec/32 §2，全部是原版直接座標）。
const (
	gateBarLabelX  = 264 // sub_1C407 的 dx=0x108
	gateBarLabelY  = 8
	gateBarClearX1 = 471 // sub_10BCD 算出的右緣
	gateBarX       = 320 // sub_1C4D2 的 dx=0x140
	gateBarY       = 16  // sub_1C4D2 的 bx=0x10
	gateBarLen     = 0x97
	gateBarH       = 2
	gateBarShift   = 4 // 耐久 >> 4

	// 熱區 0x1D 的矩形（0001C43C 的註冊）。⚠ **比條本身大得多**——
	// 條只有 (320,16) 起 151×2 px，熱區涵蓋 (256,0)–(479,31) 整塊。
	// 照抄，不要縮到條的大小：原版讓玩家點那一帶就收掉。
	gateBarHotX = 256
	gateBarHotY = 0
	gateBarHotW = 224
	gateBarHotH = 32
)

// battleRightClickHotspots 是戰場畫面的**右鍵**熱區表。
//
// 原版的熱區碼有兩張分派表（docs/re/60 §10）：左鍵 `cs:0xC048`、
// 右鍵 `cs:0xC086`。⭐ **`0x1B`／`0x1C`／`0x1D` 在左鍵表裡全是預設值**——
// 點門強度條左鍵沒有反應，要右鍵才有（TALK.DAT #4「以滑鼠的右鍵回復。」）。
//
// 右鍵表只有四筆有 handler，其中三筆在 remake 不適用：`0x1B`／`0x1C` 是
// 大將名牌（remake 沒做），`0x1F` 沒有註冊點（docs/re/60 §12）。
// **做成表而不是寫死成「右鍵就收掉」**——原版本來就是按熱區碼分派的。
// battleRightClickTarget 是熱區 handler 碰得到的戰場狀態。
// 用介面而不是 *tactical.Battle，測試才能只餵它真正需要的那幾個方法。
type battleRightClickTarget interface {
	StructureBar() (durability int, shown bool)
	DismissStructureBar()
}

var battleRightClickHotspots = []struct {
	code int // 原版的熱區碼
	rect battleRect
	// live 回報這一格現在存不存在。原版的熱區是**建立版面時才註冊**、
	// 收掉時清除，所以不能一直掛著。
	live func(battleRightClickTarget) bool
	do   func(battleRightClickTarget)
}{
	{
		code: 0x1D,
		rect: battleRect{X: gateBarHotX, Y: gateBarHotY, W: gateBarHotW, H: gateBarHotH},
		live: func(b battleRightClickTarget) bool { _, shown := b.StructureBar(); return shown },
		do:   func(b battleRightClickTarget) { b.DismissStructureBar() },
	},
}

// handleBattleRightClick 送一次右鍵到熱區表。回傳有沒有人接走。
func handleBattleRightClick(b battleRightClickTarget, x, y int) bool {
	if b == nil {
		return false
	}
	for _, h := range battleRightClickHotspots {
		if h.live != nil && !h.live(b) {
			continue
		}
		if h.rect.containsPoint(x, y) {
			h.do(b)
			return true
		}
	}
	return false
}

// drawBattleGateBar 畫「門強度」＋進度條（docs/spec/32）。
//
// 它不是常駐 HUD：城壁或城門挨打才出現，20 幀後自己收掉，
// 而且**只在更小的耐久出現時才更新**——所以這條只會往下掉。
func (g *game) drawBattleGateBar(screen *ebiten.Image, b *tactical.Battle) {
	if b == nil {
		return
	}
	durability, shown := b.StructureBar()
	if !shown {
		return
	}
	// 原版先把 (264,8)-(471,23) 填色 0 再畫字（sub_10BCD）。
	vector.DrawFilledRect(screen, gateBarLabelX, gateBarLabelY,
		gateBarClearX1-gateBarLabelX+1, textdraw.GlyphH,
		color.RGBA{0, 0, 0, 255}, false)
	g.td.Draw(screen, "門強度", gateBarLabelX, gateBarLabelY, chrome.Paper)

	filled := durability >> gateBarShift
	if filled > gateBarLen {
		filled = gateBarLen
	}
	if filled < 0 {
		filled = 0
	}
	g.drawBattleBar(screen,
		battleRect{X: gateBarX, Y: gateBarY, W: gateBarLen, H: gateBarH},
		filled, g.battleGateBarColor)
}

// drawBattleSideFooter 只畫 segment1+0x3500 那一條。
// 原版按下去會切 loc_1A065 的自我修改碼，語意未解，所以不接行為
// （docs/spec/31 §5）。
func (g *game) drawBattleSideFooter(screen *ebiten.Image, r battleRect) {
	if img := g.battleSideFooter; img != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(r.X), float64(r.Y))
		screen.DrawImage(img, op)
		return
	}
	vector.DrawFilledRect(screen, float32(r.X), float32(r.Y),
		float32(r.W), float32(r.H), chrome.Menu, false)
}

// battleFieldName 對應 sub_1C955：攻城戰用據點名，野戰依戰場類別
// 取「陸上」或「海上」。
//
// 戰場類別由 sub_19A33 依戰場編號分三段（docs/re/58 §4）：
// < 0xC0 攻城、0xC0–0xD0 類別 1、≥ 0xD1 類別 2。
func (g *game) battleFieldName(p *state.Pending) string {
	if p == nil {
		return ""
	}
	siege := p.Mode == combat.Siege
	field := g.battle.FieldNumber(p.Node, siege)
	if siege {
		if field >= 0 && field < len(g.world.Cities) {
			return big5(g.world.Cities[field].Name)
		}
		return ""
	}
	if field >= 0xD1 {
		return "海上"
	}
	return "陸上"
}

// battleLordName 對應 sub_1C9AB：主將 → 所屬勢力 → 君主 → 武將名。
// 原版沒有「勢力名」欄位，畫面上的勢力名就是君主的姓名
// （docs/re/33 §1 的 sub_188B0）。
func (g *game) battleLordName(p *state.Pending, side int) string {
	commander := g.battleCommander(p, side)
	if commander < 0 || commander >= len(g.world.Generals) {
		return ""
	}
	f := g.world.Generals[commander].Faction
	if f < 0 || f >= len(g.world.Factions) {
		return ""
	}
	lord := g.world.Factions[f].Lord
	if lord < 0 || lord >= len(g.world.Generals) {
		return ""
	}
	return big5(g.world.Generals[lord].Name)
}

// drawBattleMiniMap 畫出戰術初始化時建立的 DOS/V 128×128 base minimap，
// 再疊上部隊點。
//
// 底圖：`sub_1C83E`／`sub_1C4FA`／`sub_1C51E`，MAP tile → MDL attribute
// → 2×2 palette block。
//
// 部隊點：`sub_1B240` 在 `0001B284` 依單位記錄的位址分色——
// `si < 0x600` 用調色盤索引 10、否則用 3。那個 `0x600` 就是兩側單位陣列
// 的分界（`word_1D30E` 每側 0x600 B），也就是**側 0 一色、側 1 一色**。
// 原版的側 0 恆為玩家（docs/re/60 §5），所以這裡用 g.battleSide() 對齊。
func (g *game) drawBattleMiniMap(screen *ebiten.Image, b *tactical.Battle, r battleRect) {
	if g.view == nil || g.view.Minimap() == nil || r.W < battle.TacticalMinimapWidth ||
		r.H < battle.TacticalMinimapHeight {
		// 缺少 raw MAP/MDL 或 palette 時明確留空；不能重新引入高度／
		// 每兩格取樣的假 minimap。
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(r.X), float64(r.Y))
	screen.DrawImage(g.view.Minimap(), op)
	// 順序照原版的抹除規則反推（docs/re/60 §7）：抹掉一個部隊點時，
	// 落在陣形線上還原成 11、落在十字上還原成 0——**所以部隊點在最上面，
	// 陣形線壓在十字上面**。
	g.drawBattleMiniMapCross(screen, r)
	g.drawBattleMiniMapLine(screen, b, r)
	g.drawBattleMiniMapUnits(screen, b, r)
	// ⚠ **不畫外框。** 原版的小地圖直接貼在側欄的黑底上，四周沒有線；
	// 先前那一圈白框是 remake 自己加的，逐像素對拍一眼就看得出來
	// （docs/playtest/40 §3.1）。側欄的分隔是四條橫帶（docs/spec/31 §1.1）。
}

// drawBattleMiniMapCross 畫小地圖上的游標十字（原版 `sub_1C577`，色 0）。
//
// ⭐ **十字不是鏡頭。** 它自己有一組變數（`word_1D32C`／`word_1D32E`，
// 初值 `0x20`／`0x21`），只是縮圖點選會把兩者一起改
// （docs/re/60 §7）。拿十字的位置去反推鏡頭會得到一個看似合理的錯值。
//
// 換算與地形同一條：`x = 496 + 2×Y`、`y = 80 + 2×(63 − X)`，
// 所以「固定 X」是一條水平線、「固定 Y」是一條垂直線。
func (g *game) drawBattleMiniMapCross(screen *ebiten.Image, r battleRect) {
	if g.view == nil {
		return
	}
	ink := chrome.Ink
	if x := r.X + 2*g.viewCursorY(); x >= r.X && x < r.X+battle.TacticalMinimapWidth {
		vector.DrawFilledRect(screen, float32(x), float32(r.Y), 2,
			float32(battle.TacticalMinimapHeight), ink, false)
	}
	if y := r.Y + 2*(battle.Width-1-g.viewCursorX()); y >= r.Y &&
		y < r.Y+battle.TacticalMinimapHeight {
		vector.DrawFilledRect(screen, float32(r.X), float32(y),
			float32(battle.TacticalMinimapWidth), 2, ink, false)
	}
}

// drawBattleMiniMapLine 在小地圖上畫玩家的陣形線（原版 `sub_1C5AE`，色 11）。
//
// 換算與部隊點同一條：戰場 X 對應小地圖的 y ＝ 2 × (63 − X)，
// 所以一條「固定 X」的陣形線在小地圖上是**水平線**。
func (g *game) drawBattleMiniMapLine(screen *ebiten.Image, b *tactical.Battle, r battleRect) {
	if b == nil {
		return
	}
	y := r.Y + 2*(battle.Width-1-b.Sides[g.battleSide()].Line)
	if y < r.Y || y >= r.bottom() {
		return
	}
	// ⭐ **高 2 px**：小地圖上一格戰場 ＝ 2 px（`sub_1C5AE` 畫的是一整行格子）。
	// 畫 1 px 的話原版那條線會少掉一半——實測原版佔 y 與 y+1 兩列。
	vector.DrawFilledRect(screen, float32(r.X), float32(y),
		float32(battle.TacticalMinimapWidth), 2, g.battleGateBarColor, false)
}

// drawBattleMiniMapUnits 把每個活著的兵畫成一個 2×2 點。
//
// 座標換算與底圖同一條（sub_1C51E）：
//
//	x = 2 × mapY        y = 2 × (63 − mapX)
func (g *game) drawBattleMiniMapUnits(screen *ebiten.Image, b *tactical.Battle, r battleRect) {
	if b == nil {
		return
	}
	ally := g.battleSide()
	for side := range b.Sides {
		col := g.battleUnitDotFoe
		if side == ally {
			col = g.battleUnitDotAlly
		}
		for i := range b.Sides[side].Soldiers {
			s := &b.Sides[side].Soldiers[i]
			if !s.Alive {
				continue
			}
			x := r.X + 2*s.Y
			y := r.Y + 2*(battle.Width-1-s.X)
			if x < r.X || y < r.Y || x+2 > r.right() || y+2 > r.bottom() {
				continue
			}
			vector.DrawFilledRect(screen, float32(x), float32(y), 2, 2, col, false)
		}
	}
}

func (g *game) drawBattleSideCommands(screen *ebiten.Image, b *tactical.Battle, r battleRect) {
	selected := -1
	if len(b.Sides) > 0 && len(b.Sides[g.battleSide()].Soldiers) > 0 {
		selected = int(b.Sides[g.battleSide()].Soldiers[0].Cmd)
	}
	if g.battleSideCommands != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(r.X), float64(r.Y))
		screen.DrawImage(g.battleSideCommands, op)
		// ⭐ **指令面板不畫選取框。** 原版的選取框只有兩處：底列六個 slot
		// （`sub_1C6BF`）與陣形選單（`sub_1C61F`），指令那六列一處都沒有
		// （docs/re/60 §6、§9）。實機截圖裡雙方都在交戰、命令一定不是空的，
		// 那六列上仍然沒有框（docs/playtest/40 §3.1）。
		// **目前命令要看底列**——每一格右半就是那一隊的命令圖示，
		// 那是原版放這個資訊的地方（docs/spec/33 §1.2）。
		return
	}
	g.chrome.Window(screen, r.X, r.Y, r.W, r.H, chrome.Menu)
	for i, cell := range battleSideCommandCells(r) {
		label := battleSideCommandRowLabel(i)
		if battleSideCommandRowCode[i] == selected {
			vector.DrawFilledRect(screen, float32(cell.X), float32(cell.Y),
				float32(cell.W), float32(cell.H), chrome.Select, false)
		}
		x := cell.X + (cell.W-battleCommandTextWidth(label))/2
		y := cell.Y + (cell.H-textdraw.GlyphH)/2
		g.td.Draw(screen, label, x, y, chrome.Paper)
	}
}

// drawBattleKeys 畫底列六格。**那是玩家的六個編成位置，不是命令按鈕**
// （docs/spec/33）：`sub_1C7F4` 在 (0,368) 起每 80 px 重貼 80×32 底板，
// 把位置名 glyph 放在各格 (4,6)，`sub_1C74C` 在 y=396 畫待機兵條，
// `sub_1C6BF` 在選中的格上畫兩個同心框。
//
// 由左到右是「左翼 左備 主將 前鋒 右備 右翼」——順序來自 `cs:0xD2E4`。
func (g *game) drawBattleKeys(screen *ebiten.Image, b *tactical.Battle, r battleRect) {
	amber := color.RGBA{240, 200, 120, 255}
	dim := color.RGBA{190, 190, 200, 255}
	side := &b.Sides[g.battleSide()]
	for i, cell := range splitBattleCommandCells(r) {
		squad := battleBottomSlotSquad[i]
		picked := side.SquadSelected(squad)
		if g.battleCommandBase != nil && g.battleCommandGlyphs[i] != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(cell.X), float64(cell.Y))
			screen.DrawImage(g.battleCommandBase, op)
			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(cell.X+battleSlotGlyphX),
				float64(cell.Y+battleSlotGlyphY))
			screen.DrawImage(g.battleCommandGlyphs[i], op)
			// 格子中間是**這一隊的兵種**（sub_19B6D，docs/spec/33 §1.6）。
			if arm := battleSlotArmIcon(side.Kinds[squad]); arm >= 0 &&
				arm < len(g.battleArmIcons) && g.battleArmIcons[arm] != nil {
				op = &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(cell.X+battleSlotArmX),
					float64(cell.Y+battleSlotArmY))
				screen.DrawImage(g.battleArmIcons[arm], op)
			}
			// 格子右半是**這一隊目前的命令**（sub_1C673，docs/spec/33 §1.2）。
			// 每隊的第一個兵是隊長，命令看它的 Cmd；就位（7）畫成陣形。
			if code := battleSlotOrderIcon(side.Soldiers[squad*tactical.PerSquad].Cmd); code >= 0 &&
				code < len(g.battleOrderIcons) && g.battleOrderIcons[code] != nil {
				op = &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(cell.X+battleSlotOrderX),
					float64(cell.Y+battleSlotOrderY))
				screen.DrawImage(g.battleOrderIcons[code], op)
			}
		} else {
			vector.StrokeRect(screen, float32(cell.X+1), float32(cell.Y+1),
				float32(cell.W-2), float32(cell.H-2), 1,
				color.RGBA{90, 170, 90, 255}, false)
			col := dim
			if picked {
				col = amber
			}
			label := battleSquadLabels[squad]
			x := cell.X + (cell.W-battleCommandTextWidth(label))/2
			g.td.Draw(screen, label, x, r.Y+5, col)
		}
		// 待機兵條：sub_1C74C 的長度上限 0x4C，色 12。
		g.drawBattleBar(screen,
			battleRect{X: cell.X + battleSlotBarX, Y: battleSlotBarY,
				W: battleSlotBarLen, H: battleSlotBarH},
			clampBarLen(side.Reserve[squad], battleSlotBarLen),
			g.battleCommandSelect)
		if picked {
			outer, inner := battleSlotSelectRects(cell.X)
			vector.StrokeRect(screen, float32(outer.X), float32(outer.Y),
				float32(outer.W), float32(outer.H), 1, g.battleCommandSelect, false)
			vector.StrokeRect(screen, float32(inner.X), float32(inner.Y),
				float32(inner.W), float32(inner.H), 1, g.battleCommandSelect, false)
		}
	}
	if b.Done {
		g.td.Draw(screen, "Enter 回戰略", r.X+r.W-112, r.Y+5, amber)
	}
}

// battleSquadLabels 以**編成位置編號**為索引（0 主將…5 右備），
// 只在取不到原版 glyph 時當 fallback。
var battleSquadLabels = [...]string{"主將", "前鋒", "左翼", "右翼", "左備", "右備"}

func clampBarLen(v, max int) int {
	if v < 0 {
		return 0
	}
	if v > max {
		return max
	}
	return v
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
	p, setup, err := battlesetup.Load(battlesetup.Options{
		Dir: dir, World: g.world, Map: worldMapOf(g.lib),
		Warn: func(m string) { log.Print("⚠ " + m) },
	})
	if err != nil {
		log.Printf("⚠ %v；玩家的戰鬥會走自動判定", err)
		return
	}
	g.battle = p
	g.battleLib = p.Library()
	g.battleSprites = p.Sprites()
	g.tactical = setup
	g.world.SetTactical(setup)
}

// worldMapOf 取大地圖，缺席時回 nil（野戰會走預設戰場）。
func worldMapOf(l *library.Library) *world.Map {
	if l == nil {
		return nil
	}
	return l.World
}

// demoBattle 是**驗收用**的捷徑：直接擺一場戰術戰鬥出來。
// choose=false 時停在正常遭遇決策選單，供畫面驗收。
// demoBattle 是驗收捷徑。siegeNode ≥ 0 時指定攻城的戰場（＝據點編號），
// **為的是能在 remake 上開出與原版影格同一張戰場**——不同戰場的地形與
// 圖塊組都不同，拿兩張不同的戰場比顏色不算數。
func (g *game) demoBattle(siege, choose bool, f siegeFixture) {
	siegeNode, defend := f.node, f.defend
	if att, def := f.corps[0], f.corps[1]; att >= 0 && def >= 0 {
		g.demoBattleWithCorps(siege, choose, siegeNode, att, def, f.steps)
		return
	}
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
	if defend {
		// 玩家守城：兩個角色對調。原版這一支會 `or byte_10D35, 0C0h`，
		// 也就是**側欄換邊 ＋ 戰場轉 180 度**（docs/spec/56 §1）。
		me, foe = foe, me
	}
	g.stageEncounter(siege, choose, siegeNode, f.steps, me, foe)
}

// demoBattleWithCorps 拿**存檔裡現成的**兩支軍團開一場，不現編
// （`-siege-corps`，docs/spec/90 §2.3）。兵種、人數、主將都照存檔，
// 所以側欄的名字與計量條才有機會與原版一致。
//
// 攻守由參數直接指定，**不再看 `-siege-defend`**——玩家站哪一邊
// 由那兩支軍團的勢力決定，戰場要不要轉 180 度也跟著（docs/spec/56 §1）。
func (g *game) demoBattleWithCorps(siege, choose bool, siegeNode, att, def, steps int) {
	if att >= len(g.world.Corps) || def >= len(g.world.Corps) || att == def {
		g.setEvent("-siege-corps 的編號超出範圍")
		return
	}
	me, foe := &g.world.Corps[att], &g.world.Corps[def]
	if !me.Alive || !foe.Alive {
		g.setEvent("-siege-corps 指到的軍團不存在（用 -list-corps 看有哪些）")
		return
	}
	g.stageEncounter(siege, choose, siegeNode, steps, me, foe)
}

// stageEncounter 把攻方與守方擺到會遭遇的位置，跑到開打，
// 再推進 steps 個戰術 tick。兩條 fixture 路徑共用這一支——
// **擺位與遭遇的規則只留一份實作**（`CLAUDE.md` §7 第 6 條）。
func (g *game) stageEncounter(siege, choose bool, siegeNode, steps int, me, foe *state.Corps) {
	// 擺位與遭遇的規則在 `internal/battlesetup`，手機版的驗收路徑用同一支。
	att, def := corpsIndex(g.world, me), corpsIndex(g.world, foe)
	battlesetup.StageEncounter(g.world, g.rng, battlesetup.StageOptions{
		Siege: siege, Node: siegeNode, Attacker: att, Defender: def,
	})
	// demoBattle 是驗收捷徑，但仍經過與正常行軍相同的遭遇決策狀態；
	// 這個旗標的語意是「最後開戰場」，所以在這裡選擇戰鬥指揮。
	if choose && g.world.PendingEncounter() != nil {
		if err := g.world.ChooseBattleCommand(); err != nil {
			g.setEvent(err.Error())
			return
		}
	}

	if g.battleActive() {
		// 只跑到部隊展開。900 tick 會使野戰 fixture 在第一幀 GUI
		// 前就結束，造成「攻城有戰場、兩軍遭遇只有戰果」的假差異。
		// steps ＝ 0 是**原版開場對白那一幀**：一個兵都還沒動。
		for i := 0; i < steps; i++ {
			g.world.PendingBattle().Battle.Step()
		}
		// ⭐ 對白框的時鐘要跟著推：正常迴圈是 `updateBattle` 推 n 格戰場
		// **同時** `tickBattleTalk(n)`，驗收路徑漏掉的話框永遠不消失
		// （原版 60 tick 到期，docs/spec/60）。
		g.startBattleTalk(g.world.PendingBattle())
		g.tickBattleTalk(steps)
		g.tacticalSpeed = 1
	}
}

// corpsIndex 反查一支軍團在表裡的編號。
func corpsIndex(w *state.World, c *state.Corps) int {
	for i := range w.Corps {
		if &w.Corps[i] == c {
			return i
		}
	}
	return -1
}
