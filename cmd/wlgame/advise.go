package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/persuasion"
	"github.com/wicanr2/wolong_cht/internal/ui/listwin"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 進言（說明書 3.2、3.9）。
//
// **玩家扮演軍師，不是君主**——戰略指令不是直接執行，而是向君主提議。
// 君主可能拒絕，然後玩家挑理由說服他。
//
// 這一層只做畫面與流程；判定全在 internal/rules/persuasion，
// 那邊有 13 條測試把說明書 3.9 的規則釘住。

var adviseCommands = []persuasion.Command{
	persuasion.Hostility, persuasion.CeaseFire, persuasion.Cooperate,
}

// adviseActive 回報進言流程是不是開著。開著時時間會停——
// 它是非常駐視窗（15-realtime.md §2）。
func (g *game) adviseActive() bool { return g.advise != adviseNone }

func (g *game) openAdvise() {
	g.advise = advisePickCommand
	g.sessCur = 0
	g.adviseLog = nil
}

func (g *game) closeAdvise() {
	g.advise = adviseNone
	g.sess = nil
	g.adviseLog = nil
}

// situation 把目前的世界狀態換成說服判定要的局勢。
//
// 每一欄都對應已解出的資料：國力用據點數、疲弊用資金 < 0、
// 侵攻狀態用勢力記錄 +0x19（docs/re/08 §1）。
func (g *game) situation(target int) persuasion.Situation {
	me := g.world.Factions[g.world.Player]
	them := g.world.Factions[target]
	return persuasion.Situation{
		Aggression:  me.Aggression,
		OurCities:   me.Cities,
		TheirCities: them.Cities,
		AllyCities:  them.Cities,
		OurFunds:    me.Funds,
		TheirFunds:  them.Funds,

		// 交友度是**有向的**：這裡要的是自家君主看對方的值。
		Friendship: g.world.Friendship[g.world.Player][target].Value(),

		TheyInvadeThirdParty: them.InvasionTarget != diplomacy.NoTarget &&
			them.InvasionTarget != g.world.Player,
		TheyInvadeUs: them.InvasionTarget == g.world.Player,
	}
}

// updateAdvise 處理進言流程的輸入。回傳 true 表示它吃掉了這一幀的輸入。
func (g *game) updateAdvise() bool {
	if !g.adviseActive() {
		return false
	}
	switch g.advise {
	case advisePickCommand:
		for i := range adviseCommands {
			if pressed(ebiten.Key1 + ebiten.Key(i)) {
				g.adviseCmd = adviseCommands[i]
				g.openTargetList()
				g.advise = advisePickTarget
			}
		}
		if pressed(ebiten.KeyEscape) {
			g.closeAdvise()
		}

	case advisePickTarget:
		// 對象用一覽表選，兩段式選取的規則由 listwin 負責。
		if g.list == nil {
			g.closeAdvise()
			return true
		}
		switch {
		case pressed(ebiten.KeyArrowUp):
			g.list.Move(-1)
		case pressed(ebiten.KeyArrowDown):
			g.list.Move(1)
		case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
			if id, ok := g.list.Confirm(); ok {
				g.target = id
				g.list = nil
				g.beginPersuasion()
			}
		case pressed(ebiten.KeyEscape):
			if g.list.Cancel() {
				g.list = nil
				g.advise = advisePickCommand
			}
		}

	case advisePersuade:
		opts := persuasion.Options(g.adviseCmd)
		switch {
		case pressed(ebiten.KeyArrowUp):
			g.sessCur = (g.sessCur + len(opts) - 1) % len(opts)
		case pressed(ebiten.KeyArrowDown):
			g.sessCur = (g.sessCur + 1) % len(opts)
		case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
			g.offerReason(opts[g.sessCur])
		case pressed(ebiten.KeyEscape):
			g.closeAdvise()
		}
	}
	return true
}

func (g *game) openTargetList() {
	var rows []int
	for _, i := range g.world.AliveFactions() {
		if i != g.world.Player {
			rows = append(rows, i)
		}
	}
	fs := &g.world.Factions
	g.list = listwin.New(listwin.Factions, []listwin.Column{
		{Title: "君主", Less: func(a, b int) bool {
			return g.world.LordName(a) < g.world.LordName(b)
		}},
		{Title: "據點", Less: func(a, b int) bool { return fs[a].Cities > fs[b].Cities }},
		{Title: "武將", Less: func(a, b int) bool { return fs[a].Generals > fs[b].Generals }},
	}, rows, 10, &g.sortMem)
}

func (g *game) beginPersuasion() {
	g.sess = persuasion.Begin(g.adviseCmd, g.situation(g.target))
	g.sessCur = 0
	g.advise = advisePersuade
	g.adviseLog = []string{
		big5(g.world.LordName(g.world.Player)) + "：「" +
			big5(g.world.LordName(g.target)) + "？說來聽聽。」",
	}
}

func (g *game) offerReason(r persuasion.Reason) {
	out, dt := g.sess.Offer(r)
	g.world.Trust += dt
	if g.world.Trust < 0 {
		g.world.Trust = 0
	}
	g.adviseLog = append(g.adviseLog, "軍師：「"+r.String()+"。」")

	switch out {
	case persuasion.Agreed:
		g.adviseLog = append(g.adviseLog, "君主：「好，就依你所言。」")
		g.lastEvent = g.adviseCmd.String() + " 成立"
		g.advise = adviseNone
		g.sess = nil
	case persuasion.Failed:
		g.adviseLog = append(g.adviseLog,
			fmt.Sprintf("君主：「此言不實。」　信賴度 %d", dt))
		g.lastEvent = "說服失敗"
		g.advise = adviseNone
		g.sess = nil
	case persuasion.Withdrawn:
		g.adviseLog = append(g.adviseLog, "軍師：「……此事容後再議。」")
		g.lastEvent = "進言撤回"
		g.advise = adviseNone
		g.sess = nil
	default:
		g.adviseLog = append(g.adviseLog, "君主：「唔……還有呢？」")
	}
	// 說服結束後把對話留在畫面上，等玩家按 ESC 關掉。
	if g.advise == adviseNone && len(g.adviseLog) > 0 {
		g.advise = advisePersuade
		g.sess = nil
	}
}

// drawAdvise 畫進言流程。
func (g *game) drawAdvise(screen *ebiten.Image) {
	white := color.RGBA{240, 240, 230, 255}
	amber := color.RGBA{240, 200, 120, 255}
	dim := color.RGBA{150, 150, 160, 255}
	red := color.RGBA{240, 140, 140, 255}

	box := func(x, y, w, h int) {
		vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h),
			color.RGBA{0, 0, 0, 225}, false)
		vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h),
			1, amber, false)
	}
	lh := textdraw.GlyphH + 2

	switch g.advise {
	case advisePickCommand:
		box(40, 60, 260, 30+len(adviseCommands)*lh)
		g.td.Draw(screen, "進　言", 48, 66, amber)
		for i, c := range adviseCommands {
			g.td.Draw(screen, fmt.Sprintf("%d　%s", i+1, c.String()),
				48, 86+i*lh, white)
		}

	case advisePersuade:
		// 高度要含最後那一行提示，否則會被框邊切掉。
		h := 46 + len(g.adviseLog)*lh
		if g.sess != nil {
			h += 8 + 7*lh
		}
		box(40, 44, 420, h)
		g.td.Draw(screen, "說　得　　對象 "+big5(g.world.LordName(g.target)),
			48, 50, amber)
		y := 70
		for _, ln := range g.adviseLog {
			col := white
			if len(ln) > 3 && ln[:3] == "軍師" {
				col = dim
			}
			g.td.Draw(screen, ln, 48, y, col)
			y += lh
		}
		if g.sess == nil {
			g.td.Draw(screen, "ESC 關閉", 48, y+6, dim)
			return
		}
		y += 8
		for i, r := range persuasion.Options(g.adviseCmd) {
			col := white
			mark := "　"
			if i == g.sessCur {
				mark = "●"
				col = amber
			}
			if r == persuasion.Withdraw {
				col = dim
				if i == g.sessCur {
					col = red
				}
			}
			g.td.Draw(screen, mark+r.String(), 48, y, col)
			y += lh
		}
		g.td.Draw(screen, "↑↓ 選擇　Enter 提出　ESC 放棄", 48, y+4, dim)
	}
}
