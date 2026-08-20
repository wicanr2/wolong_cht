package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// messageDialog 保存已經由 TALK.DAT 展開、並經過 remake 實際字寬換行、
// 但尚未由玩家確認的訊息。原始 TALK 行仍先作為 hard boundary；自動換行
// 只發生在單一原始行內，不把翻譯文字寫回資料檔。
type messageDialog struct {
	lines        []string
	page         int
	portraitPage int // -1 表示這則訊息沒有可回查的 speaker 肖像。

	// lower 表示這一則要畫在**事件場景的下框**（原版 `sub_13CDC`）。
	// 下框寫死是軍師在說話，也就是玩家自己（docs/spec/42 §1）。
	lower bool
	// scene 是要留在背後的 IVENTGRF 頁；−1 表示不畫插圖。
	// 原版的事件對話是疊在插圖上的，插圖不會因為換一句話而消失。
	scene int
}

const (
	// messagePageRows ＝ 訊息框裡放得下幾列。原版的框高 80 px、
	// 內容區 64 px（sub_10BCD 四邊各內縮 8），一列 16 px ⇒ **4 列**。
	// TALK.DAT 尾端的空行是結構終止符，不佔這四列之一。
	//
	// 原文本來就折好了：1,022 則裡只有 4 則有 5 行，而那四則全是
	// 五選一的選單，由選單常式畫、不進這個框（docs/re/66 §6）。
	messagePageRows = talkBoxRows
	// 每列 10 個全形字 ＝ 160 px，出處是框右內緣減文字起點
	// （docs/re/66 §3）。TALK.DAT 有 825 行剛好是這個寬度。
	messageContentWidth = talkTextWidth
)

func layoutMessageLines(lines []string) []string {
	return textdraw.WrapLines(lines, messageContentWidth)
}

func messagePage(lines []string, page int) ([]string, int, bool) {
	if page < 0 {
		page = 0
	}
	pages := (len(lines) + messagePageRows - 1) / messagePageRows
	if pages == 0 {
		return nil, 0, false
	}
	if page >= pages {
		page = pages - 1
	}
	start := page * messagePageRows
	end := start + messagePageRows
	if end > len(lines) {
		end = len(lines)
	}
	return lines[start:end], pages, true
}

func (g *game) messageActive() bool {
	return len(g.messages) > 0
}

// talkLines 取出 TALK.DAT 的原始行並代入目前已證實可用的 marker。
// 缺少 marker 時整則訊息 fail-closed，不顯示半句或把錯誤索引當成文字。
func (g *game) talkLines(index int, vars map[byte]string) ([]string, bool) {
	if g == nil || g.lib == nil {
		return nil, false
	}
	return g.lib.Talk.Lines(index, vars)
}

// textDecodeBig5 讓訊息檔的解碼集中在既有 Big5 呈現路徑；獨立成小函式
// 方便 talkLines 保持「先代 marker、再顯示」的順序。
func textDecodeBig5(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}
	return big5(string(raw))
}

func (g *game) enqueueEventMessages(ev state.Event) {
	for _, notice := range ev.TalkNotices {
		// 事件 2／3 的 base TALK 不是獨立 modal。原版
		// sub_13902 → sub_13C99 會把 base 與 General +0x1E 的變體、
		// 肖像、IVENTGRF 場景及三選一一次畫出；pending choice 的
		// composite renderer 會接手它。state 仍保留 notice 作為事件證據。
		if g.world != nil && g.world.PendingDiplomacy() != nil &&
			isCompositeDiplomacyNotice(notice) {
			continue
		}
		g.enqueueTalkNotice(notice)
	}
	for _, id := range ev.ReleasedGenerals {
		if id < 0 || id >= len(g.world.Generals) ||
			g.world.Player < 0 || g.world.Player >= len(g.world.Factions) ||
			g.world.Generals[id].Faction != g.world.Player {
			// sub_150D7 只有在釋放後的 General +0x1C 等於玩家勢力
			// （byte_10CFF）時，才建立 TALK #37 的玩家通知。
			continue
		}
		g.enqueueTalkNotice(state.TalkNotice{
			Index: 0x25, City: -1, Faction: -1, General: id, Amount: -1,
		})
	}
}

func (g *game) enqueueTalk(index int, vars map[byte]string) {
	g.enqueueTalkWithPortrait(index, vars, -1)
}

func (g *game) enqueueTalkWithPortrait(index int, vars map[byte]string, portraitPage int) {
	lines, ok := g.talkLines(index, vars)
	if !ok || len(lines) == 0 {
		return
	}
	visible := false
	for _, line := range lines {
		if line != "" {
			visible = true
			break
		}
	}
	if !visible {
		// TALK.DAT 的空槽仍可能保留一個資料上的空行；它不應變成
		// 玩家看得到的空白 modal（事件 9 的 #409 就是這種情況）。
		return
	}
	lineWidth := messageContentWidth
	if portraitPage >= 0 {
		lineWidth = talkTextWidth
	}
	g.messages = append(g.messages, messageDialog{
		lines:        textdraw.WrapLines(lines, lineWidth),
		portraitPage: portraitPage,
		scene:        -1,
	})
}

// enqueueAdvisorTalk 把一則放進**事件場景的下框**：軍師的肖像、
// 下框的位置，背後留著 IVENTGRF 的那一頁（docs/spec/42 §3）。
func (g *game) enqueueAdvisorTalk(index int, vars map[byte]string, scene int) {
	before := len(g.messages)
	g.enqueueTalkWithPortrait(index, vars, g.playerAdvisorPortrait())
	if len(g.messages) == before {
		return // 空槽，enqueueTalkWithPortrait 已經擋掉了
	}
	d := &g.messages[len(g.messages)-1]
	d.lower, d.scene = true, scene
}

// playerAdvisorPortrait 是軍師的頭像（勢力 +0x02 → 武將 +0x01）。
// 沒有軍師（原版寫 0x7F）就退回一般通知的那一張，不要畫錯人。
func (g *game) playerAdvisorPortrait() int {
	if g == nil || g.world == nil || g.world.Player < 0 ||
		g.world.Player >= len(g.world.Factions) {
		return defaultPortraitPage
	}
	advisor := g.world.Factions[g.world.Player].Advisor
	if advisor < 0 || advisor >= len(g.world.Generals) ||
		!g.world.Generals[advisor].Alive {
		return defaultPortraitPage
	}
	return g.world.Generals[advisor].Portrait
}

func isCompositeDiplomacyNotice(notice state.TalkNotice) bool {
	return notice.Index == diplomacyTalkBase(state.DiplomacyCeasefire) ||
		notice.Index == diplomacyTalkBase(state.DiplomacyCooperation)
}

// noticePortraitPage 以 state notice 可直接回查的 General／Faction 取肖像；
// 沒有 speaker 指標時退回玩家君主。事件 3 的玩家君主路徑已由原版
// sub_13902 → sub_187FF／sub_13C99 與 fixture 截圖證實；其他 generic notice
// 的退回順序是呈現層的強推論，不把它升格成原版 speaker 語意。
func (g *game) noticePortraitPage(notice state.TalkNotice) int {
	if g == nil || g.world == nil {
		return -1
	}
	if notice.General >= 0 && notice.General < len(g.world.Generals) {
		return g.world.Generals[notice.General].Portrait
	}
	if notice.Faction >= 0 && notice.Faction < len(g.world.Factions) {
		lord := g.world.Factions[notice.Faction].Lord
		if lord >= 0 && lord < len(g.world.Generals) {
			return g.world.Generals[lord].Portrait
		}
	}
	return g.playerLordPortrait()
}

// noticeTalkIndex 把**組編號**展開成實際索引。
//
// `TALK.DAT` 索引 ≥ `0x196` 是八格一組，`sub_18810` 以說話者的
// `+0x1E` 選組內第幾個（docs/re/25 §1）。
//
// ⚠ 展開用的是**原始** `+0x1E`，不是 `talkVariant()` 收斂過的 0–2：
// 那個「≥3 減 3」只適用 `sub_13C99` 的君主路徑，而臣下型的變體就在 3–7。
// 收斂過會讓外交官的回報變成君主的命令句。
func (g *game) noticeTalkIndex(notice state.TalkNotice) int {
	if notice.Index < talkVariantGroupBase {
		return notice.Index
	}
	variant := 0
	if g != nil && g.world != nil &&
		notice.General >= 0 && notice.General < len(g.world.Generals) {
		variant = g.world.Generals[notice.General].TalkVariant
	}
	return resolveBattleTalkIndex(notice.Index, variant)
}

func (g *game) enqueueTalkNotice(notice state.TalkNotice) {
	// 變數的組法在 `internal/state`，手機版用同一支。
	vars, ok := g.world.TalkNoticeVars(notice, big5)
	if !ok {
		return
	}
	portrait := g.noticePortraitPage(notice)
	if notice.NoPortrait {
		portrait = -1
	}
	g.enqueueTalkWithPortrait(g.noticeTalkIndex(notice), vars, portrait)
}

func (g *game) drawMessage(screen *ebiten.Image) {
	d := g.messages[0]
	lines, pages, ok := messagePage(d.lines, d.page)
	if !ok {
		return
	}
	// **原版只有一種訊息框**（docs/spec/41）：寬高固定、有沒有講話的人
	// 不改變版面，變的只有位置。沒有指定肖像時用一般通知的那一張
	// （KAOGRF 第 147 張）。
	portrait := d.portraitPage
	if portrait < 0 {
		portrait = defaultPortraitPage
	}
	x, y := talkBoxX, talkBoxY
	if d.lower {
		// 事件場景的下框：軍師在說話（docs/spec/42）。
		x, y = talkLowerBoxX, talkLowerBoxY
	}
	if d.scene >= 0 {
		g.drawIventScene(screen, d.scene)
	}
	g.drawLegacyTalkBox(screen, x, y, talkBoxW, talkBoxH, lines, portrait)
	// 翻頁提示是 remake 加的：原版靠等待輸入，沒有頁碼。
	if pages > 1 {
		g.td.Draw(screen, fmt.Sprintf("%d／%d", d.page+1, pages),
			x+talkBoxW-48, y+talkBoxH-24, chrome.Paper)
	}
}
