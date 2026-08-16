package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/persuasion"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
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
	g.ally = -1
	g.target = -1
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
	// 「有別的勢力正在侵攻我方」——停戰提案的「我正在防禦戰」用這個，
	// 原版是 `sub_16577` 裡掃 22 個勢力的迴圈，**而且跳過交涉對象本身**。
	anyone := false
	for i := range g.world.Factions {
		f := &g.world.Factions[i]
		if f.Alive && i != target && f.InvasionTarget == g.world.Player {
			anyone = true
			break
		}
	}
	s := persuasion.Situation{
		Trust:       g.world.Trust,
		Aggression:  me.Aggression,
		OurCities:   me.Cities,
		TheirCities: them.Cities,
		OurFunds:    me.Funds,
		TheirFunds:  them.Funds,

		// 交友度是**有向的**，而且判定式要的是**含和平位元的原始值**
		// （`0x80` ＝ 和平）——門檻常數本身就帶著它。
		Friendship: g.world.Friendship[g.world.Player][target].Raw(),

		TheyInvadeThirdParty: them.InvasionTarget != diplomacy.NoTarget &&
			them.InvasionTarget != g.world.Player,
		TheyInvadeUs:    them.InvasionTarget == g.world.Player,
		AnyoneInvadesUs: anyone,
	}
	if g.adviseCmd == persuasion.Cooperate && g.ally >= 0 &&
		g.ally < len(g.world.Factions) {
		s.AllyCities = g.world.Factions[g.ally].Cities
		s.AllyFriendship = g.world.Friendship[g.world.Player][g.ally].Raw()
		s.SameFactionPicked = g.ally == target
	}
	return s
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
				if g.adviseCmd == persuasion.Cooperate {
					g.advise = advisePickAlly
				} else {
					g.advise = advisePickTarget
				}
			}
		}
		if pressed(ebiten.KeyEscape) {
			g.closeAdvise()
		}

	case advisePickAlly, advisePickTarget:
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
				g.list = nil
				if g.advise == advisePickAlly {
					g.ally = id
					g.openTargetList()
					g.advise = advisePickTarget
				} else {
					g.target = id
					g.beginPersuasion()
				}
			}
		case pressed(ebiten.KeyEscape):
			if g.list.Cancel() {
				g.list = nil
				if g.advise == advisePickTarget && g.adviseCmd == persuasion.Cooperate {
					g.openTargetList()
					g.advise = advisePickAlly
				} else {
					g.advise = advisePickCommand
				}
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
	g.openFactionPicker(rows, "↑↓ 移動　Enter 選取／決定　1-6 排序　ESC 取消", nil)
}

func (g *game) beginPersuasion() {
	s := g.situation(g.target)
	g.sessCur = 0
	g.advise = advisePersuade
	// ① 君主開場、② 軍師的進言——兩則都是原版的原文（docs/spec/44 §2）。
	base := adviseTalkBase(g.adviseCmd)
	g.adviseLog = nil
	for _, l := range []string{
		g.adviseLine("君主", base+g.playerTalkVariant()),
		g.adviseLine("軍師", base+3),
	} {
		if l != "" {
			g.adviseLog = append(g.adviseLog, l)
		}
	}
	queued := false
	if g.adviseCmd == persuasion.Hostility {
		queued = g.world.HasQueuedDeclaration(g.world.Player, g.target)
	}
	switch reaction := persuasion.FirstReaction(g.adviseCmd, s, queued); reaction {
	case persuasion.AskReason:
		g.sess = persuasion.Begin(g.adviseCmd, s)
	default:
		g.sess = nil
		g.adjustTrust(persuasion.ReactionTrustDelta(reaction))
		base := adviseTalkBase(g.adviseCmd)
		if line := g.adviseLine("君主",
			adviseReplyIndex(base, reaction, g.playerTalkVariant())); line != "" {
			g.adviseLog = append(g.adviseLog, line)
		}
		if reaction == persuasion.Agree {
			g.commitAdvice()
		}
	}
}

func (g *game) offerReason(r persuasion.Reason) {
	if g.sess == nil {
		return
	}
	out, dt := g.sess.Offer(r)
	g.adjustTrust(dt)
	g.adviseLog = append(g.adviseLog, "軍師：「"+r.String()+"。」")

	switch out {
	case persuasion.Agreed:
		g.commitAdvice()
	case persuasion.Failed:
		g.adviseLog = append(g.adviseLog,
			fmt.Sprintf("君主：「此言不實。」　信賴度 %d", dt))
		g.lastEvent = "說服失敗"
		g.sess = nil
	case persuasion.Withdrawn:
		g.adviseLog = append(g.adviseLog, "軍師：「……此事容後再議。」")
		g.lastEvent = "進言撤回"
		g.sess = nil
	default:
		g.adviseLog = append(g.adviseLog, "君主：「唔……還有呢？」")
	}
}

// adjustTrust 對應原版 sub_13D91／sub_13DC9 的 byte 飽和行為。
func (g *game) adjustTrust(delta int) {
	if g.world != nil {
		g.world.AdjustTrust(delta)
	}
}

// adviseTalkBase 是三個進言指令的 TALK 起點——`sub_16405`／`sub_164F1`／
// `sub_16623` 傳給 `sub_13830` 的 `cx`（docs/spec/44 §1）。
//
// ⚠ **三組措辭各不相同**，不能共用一組。
func adviseTalkBase(c persuasion.Command) int {
	switch c {
	case persuasion.CeaseFire:
		return 0x96 // 150
	case persuasion.Cooperate:
		return 0xD6 // 214
	default:
		return 0x56 // 86，敵對（進兵）
	}
}

// adviseReplyIndex 是君主回答的 TALK 索引：`base + 4 + 結果碼×3`
// （`sub_13830` 的 `cx = base + al×3 + 4`）。碼 ≥ 4 一律用 83。
//
// 每個位置佔三則，因為 `sub_13C99` 會把君主的**說話型**加進索引。
func adviseReplyIndex(base int, r persuasion.Reaction, variant int) int {
	if r >= persuasion.SameFaction {
		return 0x53 + variant // #83「我想軍師並不是來談笑的。」
	}
	return base + 4 + 3*int(r) + variant
}

// adviseLine 把一則 TALK 展開成一行，前面加說話者標記。
//
// 標記是 **remake 差異**：原版用兩個講話框 ＋ 肖像分辨誰在說
// （docs/re/66 §5.2），而這裡的清單視窗沒有肖像可以放。
func (g *game) adviseLine(speaker string, index int) string {
	lines, ok := g.talkLines(index, g.adviseTalkVars())
	if !ok || len(lines) == 0 {
		return ""
	}
	joined := ""
	for _, l := range lines {
		if l == "" {
			continue
		}
		joined += l
	}
	if joined == "" {
		return ""
	}
	return speaker + "：「" + joined + "」"
}

// adviseTalkVars 是進言那幾則的變數：`{3}` 交涉對象、`{4}` 軍師（玩家）、
// `{6}` 排版標記（空字串，原版 handler 只調 X 不輸出字元）。
func (g *game) adviseTalkVars() map[byte]string {
	vars := map[byte]string{'6': ""}
	if g.target >= 0 && g.target < len(g.world.Factions) {
		vars['3'] = big5(g.world.LordName(g.target))
	}
	if p := g.world.Player; p >= 0 && p < len(g.world.Factions) {
		if a := g.world.Factions[p].Advisor; a >= 0 && a < len(g.world.Generals) {
			vars['4'] = big5(g.world.Generals[a].Name)
		}
	}
	return vars
}

// commitAdvice 將第一反應直接同意或說服成功接到原版 producer：敵對
// 立即走 sub_13526，停戰／協力分別寫入事件 6／7，讓後續每小時 handler
// 依原版節拍處理外交結果。
func (g *game) commitAdvice() bool {
	var ok bool
	switch g.adviseCmd {
	case persuasion.Hostility:
		ok = g.world.ApplyPlayerHostility(g.target)
	case persuasion.CeaseFire:
		ok = g.world.QueuePlayerCeasefire(g.target)
	case persuasion.Cooperate:
		ok = g.world.QueuePlayerCooperation(g.ally, g.target)
	}
	if !ok {
		g.adviseLog = append(g.adviseLog, "君主：「局勢已變，這項進言沒有成立。」")
		g.lastEvent = "進言失效"
		g.sess = nil
		return false
	}
	g.adviseLog = append(g.adviseLog, "君主：「好，就依你所言。」")
	g.lastEvent = g.adviseCmd.String() + " 成立"
	if g.adviseCmd == persuasion.CeaseFire || g.adviseCmd == persuasion.Cooperate {
		g.lastEvent += "（外交事件已排入）"
	}
	g.sess = nil
	return true
}

// drawAdvise 畫進言流程。
func (g *game) drawAdvise(screen *ebiten.Image) {
	white := color.RGBA{240, 240, 230, 255}
	amber := color.RGBA{240, 200, 120, 255}
	dim := color.RGBA{150, 150, 160, 255}
	red := color.RGBA{240, 140, 140, 255}

	// 進言的視窗也走原版外框（ICONGRF 段 3）。尺寸先進位到 8 的倍數，
	// 不然邊框會切在半塊上。
	box := func(x, y, w, h int) {
		up := func(v int) int { return (v + chrome.Tile - 1) / chrome.Tile * chrome.Tile }
		g.chrome.Window(screen, x, y, up(w), up(h), chrome.Menu)
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
		box(40, 44, 496, h)
		// 說得是在跟**自己的君主**對話（玩家是軍師），所以右上角放君主頭像。
		// 頁碼取武將記錄的 +0x01，不是武將編號（見 state.General.Portrait）。
		if lord := g.world.Factions[g.world.Player].Lord; lord >= 0 &&
			lord < len(g.world.Generals) {
			if img, err := g.lib.Portrait(g.world.Generals[lord].Portrait,
				int(g.world.Clock.Season())); err == nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(40+496-8-64), 52)
				screen.DrawImage(ebiten.NewImageFromImage(img), op)
			}
		}
		subject := big5(g.world.LordName(g.target))
		if g.adviseCmd == persuasion.Cooperate {
			subject = big5(g.world.LordName(g.ally)) + " → " + subject
		}
		g.td.Draw(screen, "說　得　　對象 "+subject,
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
