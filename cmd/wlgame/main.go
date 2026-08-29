// wlgame 是戰略主畫面的可玩垂直切片。
//
// 它把三層接起來跑：
//
//	internal/state    ← 從 SINARIO.DAT 載入真實的世界
//	internal/rules    ← 時鐘、月結、災害
//	internal/assets   ← 大地圖的圖塊與四季調色盤
//
//	tools/go.sh run ./cmd/wlgame -orig workplace/orig/dosv
//
// 它已接上命令、軍團、行軍、戰術戰鬥與四槽存檔 overlay；仍不是完整
// 發行版，尚缺 M7 全量文意校訂、目標平台實機與部分原版素材語意。
// 這裡把「時間在跑的世界」呈現出來，讓已定案的規格
// （docs/mechanics/15-realtime.md）能用眼睛驗收：
//
//   - 時間是連續的，不是回合制
//   - **開啟非常駐視窗會讓時間停止**（§2 的暫停規則）
//   - 季節在 3/6/9/12 月的前 16 天漸變，不是瞬間切換
//   - 月結時資金、預備兵、生產力會跳動
//
// 中文用倚天 16×15 點陣字（`-font` 指到自備的字型目錄）。
// **字型檔不隨本專案散布**，與原版資料同一個處理方式。
// 沒帶 `-font` 也跑得起來，只是中文會畫成空心方框——
// 缺字要看得出來，不能靜靜吃掉。
package main

import (
	"github.com/wicanr2/wolong_cht/internal/battlesetup"
	"github.com/wicanr2/wolong_cht/internal/ui/isoview"
	"github.com/wicanr2/wolong_cht/internal/rules/speed"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/assets/cjk"
	"github.com/wicanr2/wolong_cht/internal/assets/gfx"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/rules/battlefield"
	"github.com/wicanr2/wolong_cht/internal/rules/bgm"
	"github.com/wicanr2/wolong_cht/internal/rules/clock"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/march"
	"github.com/wicanr2/wolong_cht/internal/rules/persuasion"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/savepath"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/langpack"
	"github.com/wicanr2/wolong_cht/internal/ui/listwin"
	"github.com/wicanr2/wolong_cht/internal/ui/uitext"
	"github.com/wicanr2/wolong_cht/internal/ui/sound"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 原版 DOS/V 是 640×400：最上方 32 px 是橫幅，接著 32 px 命令列；
// 左側是 27×21 格的大地圖，右側 208 px 是縮小地圖／自勢力情報常駐 HUD。
// 命令、列表、財政與事件選擇等暫存視窗仍浮在這個自然策略畫面上。
//
// 這三個數字沿用 640×400 DOS/V 畫布與原版 16 px 格座標；歷史 PC-98 截圖只作
// 交叉驗證，不作本輪 DOS/V 畫面 oracle：
// 地圖是**對齊格線**畫的、一格 16 px，可視範圍 40×23 格 ＝ 640×368。
// 這個格數是從原版讀出來的：`sub_1D615` 的迴圈 `cx=0x28 / dx=0x17`
// （docs/re/47 §3.2），捲動上限 `cmp cx, 27Fh` 也照 640 寬算。
const (
	screenW, screenH = 640, 400
	bannerH          = 32
	viewCols         = strategyMapW / 16 // 40
	viewRows         = strategyMapH / 16 // 23
)

// 把某一格移到畫面上的哪一格（原版 `sub_12151` 的兩個入口參數，
// docs/spec/52）。開新遊戲、點縮小地圖、跳到某個據點都走這一條。
//
// `sub_12151` 的立即值就是 `ax=14h`／`cx=0Ch` ＝ (20, 12)：
// 水平放在 40 欄的正中央、垂直比 23 列的中央高半列。
//
// ⚠ 這裡曾經寫 16——那是為了抵銷「地圖整張左移四格」而來的，
// 而左移的成因是 `MMAP.MAP` 解壓後開頭有 4 byte 長度欄位
// （`world.MapHeader`）。兩邊一起修好之後，這裡照機器碼寫 20。
const (
	centreCol = 20
	centreRow = 12
)

// 橫幅右段那三個數字欄的右緣。橫幅本身（ICONGRF 段 0）已經印好
// 「年 月 日」三個字，數字要填在每個字**左邊**那塊黑底上。
// 三個值是拿松崗實機的主畫面量出來的（docs/spec/52 §2），
// 不是量橫幅圖檔——圖檔只給得出「年」在哪，給不出數字靠右靠到哪。
const (
	// 三個右緣剛好等距 32 px —— 這是「量對了」的算術檢查。
	bannerYearRight  = 496
	bannerMonthRight = 528
	bannerDayRight   = 560
	// 字模的 16 列裡上下各留一列空白，所以墨水落在 y 9–22。
	bannerTextY = 8
	// 日期的字色是調色盤第 9 色（#F3D392），與橫幅上「年 月 日」同色。
	bannerDateInk = 9

	// eventLineFrames 是事件列顯示幾幀。**六秒**——手機版的事件通知
	// 已經在用這個值（`internal/ui/phone`），兩端一致比另外挑一個數字好。
	eventLineFrames = 6 * 60
)

type game struct {
	lib    *library.Library
	// talkBase 是母本繁中的訊息表。切語言要換 lib.Talk，換回來時
	// 得有原本那一份——重讀一次 TALK.DAT 太慢，也會把校訂弄丟。
	talkBase *text.Table
	// talkPinned ＝ 使用者用 -talk-json 指定了訊息表，切語言不要蓋掉它。
	talkPinned bool
	fontDir    string
	// langNotice 是切換語言後短暫顯示的語系名，langNoticeAt 是已經畫了幾幀。
	langNotice   string
	langNoticeAt int
	world  *state.World
	rng    *rng.Rand
	td     *textdraw.Drawer
	chrome *chrome.Set // 原版視窗外框（ICONGRF 段 3）
	// amountFrame 是 DOS/V sub_17D0D 的 96×64 數值視窗內框；
	// 一般 chrome.Set 只負責其他視窗的 8×8 邊框。
	amountFrame *ebiten.Image
	// cursorImage 是 DOS/V KI.EXE seg002:031B 的白框／紅填 16×16 游標。
	// 數值 modal 內用它取代向量高亮，鍵盤 fallback 才會把箭頭定位到目前格位。
	cursorImage *ebiten.Image
	// battleCommandBase／battleCommandGlyphs／battleSideCommands 是松崗
	// DOS/V ICONGRF 段 1 的原版戰術指令素材；nil 時才用文字 fallback。
	battleCommandBase   *ebiten.Image
	battleCommandGlyphs [6]*ebiten.Image
	// battleOrderIcons 是底列每格右半那張「目前命令」的圖示
	// （ICONGRF 段 3 的 `碼 × 0xC0`，docs/spec/33 §1.2）。
	battleOrderIcons [6]*ebiten.Image
	// battleArmIcons 是底列每格中間那張兵種圖示
	// （ICONGRF 段 3 的 `0x480 + (兵種−1) × 0xC0`，docs/spec/33 §1.6）。
	battleArmIcons [3]*ebiten.Image
	battleSideCommands  *ebiten.Image
	battleCommandSelect color.RGBA

	// 側欄另外四塊美術（docs/re/60 §1.2）。battleSideFlags[0] 是下格
	// （我方，段 1 0x1000）、[1] 是上格（對方，0x0800）。
	battleSideFlags      [2]*ebiten.Image
	battleFormationStrip *ebiten.Image
	battleSideFooter     *ebiten.Image

	// 側欄外框的四塊（docs/spec/31 §1.1）：橫帶、角、左柱、右柱。
	// 索引與 library.BattleFramePart 相同。
	battleFrame [4]*ebiten.Image

	// battle 是戰場來源（`internal/battlesetup`）。翻轉旗標也在它身上——
	// 繪圖層要與建場用同一個值，否則畫面與規則層會不一致（docs/spec/56）。
	battle *battlesetup.Provider

	// battleCamAt 是 `-battle-cam` 的覆寫值（驗收用；nil ＝ 照原版初值）。
	battleCamAt *[2]int

	// 標題兩行與兩條計量條的顏色，一律查調色盤（docs/spec/54）。
	battleTitlePlace color.RGBA // 索引 9：地名與「作戰」
	battleTitleLord  color.RGBA // 索引 11：君主名與「對」
	battleMenBar     color.RGBA // 索引 12：兵力條
	battleHealthBar  color.RGBA // 索引 11：大將體力條

	// battleFormation 是陣形選單目前選中的那一格（原版 byte_1D346），
	// 值域 0–15；−1 表示還沒選。
	battleFormation int

	// 小地圖上的部隊點：sub_1B240 用調色盤索引 10（側 0 ＝ 我方）
	// 與 3（側 1 ＝ 對方）。
	battleUnitDotAlly color.RGBA
	battleUnitDotFoe  color.RGBA

	// 門強度條的顏色：sub_1C4D2 的 ah=0x0B。
	battleGateBarColor color.RGBA

	// minimapInk 是縮小地圖標記用的六個色號（docs/re/62 §2）：
	// 依序 0（黑）、15（白）、10（紅）、12（黃）、3（藍）、8（深藍）。
	minimapInk [6]color.RGBA
	// minimapFaction 是圖例第二格盯著的勢力（原版 `cs:byte_198A7`，
	// 靜態初值 0）。它決定哪個勢力的據點畫成白框藍心。
	minimapFaction int
	// factionPicker：22 勢力的選擇視窗開著沒有（原版熱區 0x17 →
	// sub_15AFC，docs/spec/35 §2.5.2）。
	factionPicker bool

	// cancelFn 讓測試注入「取消被按下」（docs/spec/73）。nil ＝ 讀真的輸入。
	cancelFn func() bool

	// lordCorps ＝ 允許把君主編成軍團長（docs/spec/76）。
	// ⚠ **預設 true，與原版不同**——使用者裁定的 remake 差異。
	// 想要原版行為就在系統選單把那一列切成「不可」。
	lordCorps bool
	// damageReport 是戰後結果畫面要不要多印一行「攻城損害」。
	// **原版沒有這個報告**，所以預設 false（docs/spec/89）。
	damageReport bool

	// roads 與 tactical 是掛在 World 上的執行期來源，不屬於存檔本體。
	// 讀取另一個槽位後要重新掛回，否則數值雖然恢復，行軍／戰鬥會悄悄退回
	// 沒有道路圖與戰術資料的降級路徑。
	roads    *march.Graph
	tactical *state.TacticalSetup

	// 存檔採明確指定的可寫 overlay；空字串代表這次執行沒有開啟持久化。
	origDir    string
	// 結局過場（docs/spec/67）。ending 是正在播的那一次，
	// endingShown 擋住「收掉之後又自己跳出來」。
	ending      *endingState
	endingShown bool
	endingCache map[int]*ebiten.Image
	sourceFile string
	saveFile   string
	saveBase   string
	saveUI     saveUIState
	// battleFastForward 是戰場 `▶▶`（快轉）開著；battleFFTouched 記錄按過
	// ——原版按過才描那一圈框（docs/spec/102）。
	battleFastForward, battleFFTouched bool
	launcher   *launcherModel
	// naming 是「自定」軍師命名視窗開著（docs/spec/104）；
	// launcherPreviewWorld 是選君主那一頁預覽的劇本世界（命名要看武將的肖像）；
	// customAdvisor 是命名完成、等開局套進世界的結果。
	naming               *namingModel
	launcherPreviewWorld *state.World
	customAdvisor        *customAdvisor

	camX, camY int

	// hud 是主畫面四個常駐視窗的開關集合，對應原版 `byte_198A6` 的
	// bit 0–3（docs/spec/13）。**初值四個全開是 remake 差異**——
	// 原版新遊戲時的初值還沒讀出來。
	hud hudWindow

	// speed／tacticalSpeed 是**原版的檔位 0–4**（0 ＝ 最高速、4 ＝ 最低速），
	// 不是「每畫面幾 tick」。原版是兩個獨立設定：戰略速度在 `ds:0CFAh`
	// 直接當等待量，戰術速度在 `ds:0CFBh` 要先 ×16（`sub_160A5`）。
	// 每一檔實際多快由 speed.Throttle 換算（docs/spec/34、docs/re/61）。
	speed            int // 戰略速度檔位
	tacticalSpeed    int // 戰術速度檔位
	strategyThrottle speed.Throttle
	tacticalThrottle speed.Throttle
	speedToast       int // 剛調過速度時在戰場浮一行提示，剩幾幀

	// sound 是音訊播放層。**音檔不隨發行包散布**——玩家自己跑
	// `tools/bgm2ogg.sh` 從原版產生（docs/spec/29 §5、§6）。
	// 沒有音檔時 Bank 的每個方法都是 no-op，所以呼叫端不必判斷。
	sound *sound.Bank

	// idleGate 對應松崗繁中版 sub_11F7F 的「游標座標未變」判定。
	// 它是讓自然世界迴圈開始的 UI 閘門；事件 queue 的 consumer 仍在
	// World.TickMap 的每時更新內，不能把兩者混成同一件事。
	idleGate idleClockGate

	// list 是開著的一覽表（武將／據點…）。它是**非常駐視窗**，
	// 所以開著的時候時間會停（15-realtime.md §2）。
	battleLib *battle.Library
	// view 是戰場的等角繪圖資源。沒有原版美術時是 nil，畫面退回高度圖。
	view          *isoview.View
	battleSprites *battle.Sprites

	list    *listwin.List
	sortMem listwin.Memory

	// listRow 把一列資料格式化成**逐欄**的字串，欄數與該家族的欄位定義
	// 相同（docs/spec/38）。一覽表本身只管狀態機（listwin），
	// 顯示什麼由開啟它的人決定。
	listRow func(id int) []string
	// listTitle 是那一組欄位的標題字串（原版是一整條，不是逐欄拼的）。
	listTitle string
	// listCellInk 讓某一格換色：士氣 < 100、外交「交戰」都要換
	// （docs/re/27 §2、§4）。回 false 就用預設墨色。
	listCellInk func(id, col int) (color.RGBA, bool)
	// listPick 是決定一列之後要做的事。回傳 true 表示關掉一覽表。
	listPick func(id int) bool
	listHint string
	// listTouched：游標列的淡色提示只在玩家操作過清單之後才畫——
	// 原版剛開窗沒有任何反白（playtest/42 §4 的 q0），游標由滑鼠本體表示。
	// 與編成視窗的 f.keyboard 同一個先例。
	listTouched bool

	// form 是編成流程的狀態。
	form formState
	// marchMode 是行軍指示的第二段（戰鬥指揮／委任／解體，docs/spec/39）。
	marchMode marchModeState

	// finance 是財政畫面的狀態。與 form 一樣是**非常駐視窗**，
	// 開著的時候時間會停（15-realtime.md §2）。
	finance financeState

	// cityInfo 是據點情報視窗的狀態（docs/spec/23）。原版純顯示，
	// 按右鍵關掉。
	cityInfo cityInfoState

	// corpsInfo 是軍團情報視窗的狀態（docs/spec/24）。它蓋在自勢力情報
	// 那一格上——原版就是同一個矩形。
	corpsInfo corpsInfoState

	// 進言的狀態機：選指令 → 選對象 → 說服。
	advise    adviseStage
	adviseCmd persuasion.Command
	ally      int
	target    int
	sess      *persuasion.Session
	sessCur   int
	// adviseCmdRow 是進言五項選單的游標。
	adviseCmdRow int

	// adviseNode 是遷都選到的據點（進言第四項）。
	adviseNode int

	// adviseQueue 是還沒演到的句子（docs/spec/45 §1.1）。
	adviseQueue []adviseStep

	// 進言畫面上兩個框各自的最新一句（docs/spec/45 §1）。
	adviseLordSaid    []string // 上框：君主
	adviseAdvisorSaid []string // 下框：軍師，也就是玩家自己

	lastEvent string
	// lastEventAt 是事件列**內容變成現在這一句**的那一幀，
	// lastEventShown 是上一次看到的那一句。原版的大地圖沒有這個框
	//（docs/spec/88 §2），remake 加的提示至少要像個提示。
	//
	// ⭐ 用「內容變了就重計時」而不是在每個設定點記時間：`lastEvent`
	// 有三十幾個設定點，其中幾個是 `+=`。**要靠每個地方都記得的機制
	// 遲早會漏掉一個**，而漏掉的症狀是那一種事件的框永遠不消失。
	lastEventAt    int
	lastEventShown string
	messages    []messageDialog
	// scenarioTitles 是四個劇本的標題（區塊 +0x40），啟動時從 SINARIO.DAT 讀。
	scenarioTitles map[int]string

	quitting bool
	quitYes  bool // ＹＥＳ／ＮＯ 對話框的選取（remake 的鍵盤操作）

	// battleTalkSession 是 presentation-only 的戰術開場 TALK queue；不進 state／存檔。
	battleTalkSession battleTalkSession

	// diplomacyRow 是事件 2／3 外交三選一視窗目前反白的選項。
	diplomacyRow int
	// diplomacyEditingAmount 對應原版 sub_13902 選到「資金提供」後
	// 進入 sub_17C6E 的第二個 modal；未進入數值器前 Enter 仍只選列。
	diplomacyEditingAmount bool

	// fundingRow 是事件 4／5 撥款三選一視窗目前反白的選項。
	fundingRow int
	// fundingEditingAmount 對應事件 4／5 的 sub_139E8 → sub_17C6E。
	fundingEditingAmount bool

	// amountCursor 是原版滑鼠格選取的跨平台呈現狀態。row／col 直接
	// 對應 CS:7D93h 的 3×6 格；鍵盤只是把同一個格動作映射出來。
	amountCursorRow, amountCursorCol int
	amountCursorActive               bool
	// amountKeyboard：鍵盤 fallback 游標只在用過鍵盤後畫（playtest/42 §3）。
	amountKeyboard bool
	// hideAmountCursor：對拍 fixture 用——headless 的指標位置不可控，
	// 會把游標畫進比對框（playtest/42 §3）。
	hideAmountCursor bool
	amountCursorOwner                int
	// amountAnchorX／Y 是 `sub_17C6E` 的 `dx`／`bx`——**錨點由呼叫端給**，
	// 事件 2／3／4／5 是 (88,184)，財政的四個熱區是 (296,184)
	// （docs/spec/78 §1.2）。0 表示還沒設，一律當成事件那一組。
	amountAnchorX, amountAnchorY int

	// 截圖模式：跑 shotAt 幀之後把畫面存成 PNG 然後結束。
	// 這是這支程式**唯一**的自動驗收方式——Ebiten 要顯示器，
	// 所以 CI／容器裡要靠 Xvfb ＋ 這個旗標才驗得到畫面。
	shotPath string
	shotAt   int
	frame    int
	shotDone bool

	// 逐幀錄製（docs/spec/71）。nil ＝ 沒開。
	rec     *recorder
	recDone bool

	// 災害物件圖像快取。key 內含季節、原版 object type 與 8 相位；
	// Level 不進 key，因為它是 marker 強度，不是動畫幀。
	disasterImages map[int]*ebiten.Image
}

// adviseStage 是進言流程的階段。
type adviseStage int

const (
	adviseNone        adviseStage = iota
	advisePickCommand             // 選五項進言
	advisePickAlly                // 協力要請先選協力方
	advisePickTarget              // 選對象勢力
	advisePersuade                // 君主拒絕了，開始說服
	advisePickCapital             // 遷都：選要搬去的據點
	adviseVerdict                 // 遷都／出陣：君主一句話定案，沒有說服迴圈
)

// timeRuns 是暫停規則。
//
//	時間推進 ⟺ 開啟中的視窗集合 ⊆ {命令, 自勢力情報, 縮小地圖}
//
// 刻意寫成一個函式而不是散在各視窗的開關程式碼裡——
// 這樣「哪些視窗會停時間」只有一個地方可以改。
func (g *game) timeRuns() bool {
	if g.world != nil && g.world.Outcome() != state.InProgress {
		return false
	}
	// 一覽表、進言、編成都是非常駐視窗 —— 開著就停時間。
	if g.list != nil || g.adviseActive() || g.form.active || g.marchMode.active ||
		g.finance.active ||
		g.saveUI.active || g.messageActive() {
		return false
	}
	// 四個常駐視窗裡**只有系統視窗會停時間**（說明書 3.1、
	// docs/mechanics/15-realtime.md §2）。
	if g.hudOpen(hudSystem) {
		return false
	}
	return true
}

// openGeneralList 開武將一覽。欄位照原版的六欄（docs/spec/38 §1.2）。
func (g *game) openGeneralList() {
	var rows []int
	for i, gen := range g.world.Generals {
		if gen.Alive && gen.Faction == g.world.Player {
			rows = append(rows, i)
		}
	}
	g.openGeneralPicker(rows, "↑↓ 移動　Enter 選取／決定　1-6 排序　ESC 取消",
		func(i int) bool {
			g.lastEvent = "選擇了 " + big5(g.world.Generals[i].Name)
			return true
		})
}

// openGeneralPicker 開一張武將清單（看或選都用這一張，docs/re/26 §4.2）。
func (g *game) openGeneralPicker(rows []int, hint string, pick func(int) bool) {
	g.list = listwin.New(listwin.Generals, g.listColumnsGenerals(), rows,
		listRowsPerPage, &g.sortMem)
	g.listTouched = false
	g.listTitle = listFamilyGenerals.Title
	g.listRow = g.listRowGeneral
	g.listCellInk = nil
	g.listHint = hint
	g.listPick = pick
}

// openCityPicker 開一張據點清單。**上昇率 0 換色**（docs/re/27 §5）。
func (g *game) openCityPicker(rows []int, hint string, pick func(int) bool) {
	g.list = listwin.New(listwin.Cities, g.listColumnsCities(), rows,
		listRowsPerPage, &g.sortMem)
	g.listTouched = false
	g.listTitle = listFamilyCities.Title
	g.listRow = g.listRowCity
	g.listCellInk = func(id, col int) (color.RGBA, bool) {
		// 原版比的是存值 100（存值 ＝ 實際成長 ＋ 100），remake 的
		// `Growth` 已經是實際成長，所以門檻是 0。
		if col == 2 && g.world.Cities[id].Growth == 0 {
			return listWarnInk, true
		}
		return color.RGBA{}, false
	}
	g.listHint = hint
	g.listPick = pick
}

// openFactionPicker 開一張勢力清單。「交戰」那一格換色（docs/re/27 §4）。
func (g *game) openFactionPicker(rows []int, hint string, pick func(int) bool) {
	g.list = listwin.New(listwin.Factions, g.listColumnsFactions(), rows,
		listRowsPerPage, &g.sortMem)
	g.listTouched = false
	g.listTitle = listFamilyFactions.Title
	g.listRow = g.listRowFaction
	g.listCellInk = func(id, col int) (color.RGBA, bool) {
		if col == 4 && g.factionDiplomacy(id) == 0 {
			return listWarnInk, true
		}
		return color.RGBA{}, false
	}
	g.listHint = hint
	g.listPick = pick
}

// drawList 畫一覽表。
//
// 原版的清單視窗是**米色底、黑字**，選取的那一列是**綠色反白條**；
// 顏色常數沿用共用 ICONGRF／palette 證據，PC-98 實機只作歷史交叉驗證。
// 兩段式選取的第一下就是把那一列變綠，第二下才決定——
// 這在原版的君主選擇畫面上實際看得到。
func (g *game) drawList(screen *ebiten.Image) {
	l := g.list
	fields := listFieldsFor(l)
	// 外框畫在客戶區**外面**：原版的 (24,88,384,176) 是內容區，
	// 16（標題）＋ 10 × 16 剛好填滿（docs/spec/38 §1.1）。
	g.chrome.Window(screen, listWinX-chrome.Tile, listWinY-chrome.Tile,
		listWinW+2*chrome.Tile, listWinH+2*chrome.Tile, chrome.Sheet)

	ink := chrome.Ink
	dim := color.RGBA{90, 80, 70, 255}

	// 標題列是**黑底白字**（影片影格看得到），而標題本身是原版的
	// 一整條字串（docs/re/26 §4.1），不是逐欄拼出來的。
	vector.DrawFilledRect(screen, float32(listWinX), float32(listWinY),
		float32(listWinW), float32(listRowH), color.RGBA{0, 0, 0, 255}, false)
	// 標題整條與文字欄一樣有半格縮排（實機量測：「武將名」ink 從 48 起，
	// docs/spec/38 §1.5）。
	g.td.Draw(screen, g.listTitle, listBodyX()+listTextInset, listWinY,
		g.paletteInk(strategyInkNormal, chrome.Paper))

	rows, first := l.Visible()
	// ⭐ **一頁永遠畫滿十列**：原版沒有資料的那幾列印的是分隔線那一行
	// （全形欄印 `－`、半形欄印 `-`），不是留白。證據是 DOS/V 實錄影格
	// 的軍團一覽——兩支軍團，下面八列全是破折號（docs/spec/38 §1.4）。
	for i := len(rows); i < listRowsPerPage; i++ {
		y := listRowY(i)
		for _, f := range fields {
			// 數字欄的破折號**貼欄本體、不加半格縮排**（parity-menus7 的
			// m1：「----」左緣 120＝欄起點）；文字欄照 §1.5 在欄起點＋8。
			x := listBodyX() + f.X
			if !f.Numeric {
				x += listTextInset
			}
			g.td.Draw(screen, listDashes(f), x, y, ink)
		}
	}
	for i, r := range rows {
		y := listRowY(i)
		if first+i == l.Cursor &&
			(g.listTouched || l.Phase() == listwin.Selected) {
			hl := color.RGBA{200, 210, 170, 255}
			if l.Phase() == listwin.Selected {
				hl = chrome.Select // 反白：原版就是這個綠
			}
			// 反白列從清單本體的左緣起（原版 sub_184BC 是
			// `word_181AE + 0x10`、寬 `word_181B2 − 0x10`）。
			vector.DrawFilledRect(screen, float32(listBodyX()), float32(y),
				float32(listBodyW()), float32(listRowH), hl, false)
		}
		cells := g.listRow(r)
		for col, cell := range cells {
			if col >= len(fields) {
				break
			}
			c := ink
			if g.listCellInk != nil {
				if got, ok := g.listCellInk(r, col); ok {
					c = got
				}
			}
			// 數字欄右靠到分隔線右緣、用 8×16 原版字模（sub_1062F 那一套，
			// docs/spec/38 §1.5）；文字欄在欄起點＋8。
			if fields[col].Numeric {
				if v, err := strconv.Atoi(strings.TrimSpace(cell)); err == nil {
					digits := fields[col].W / textdraw.HalfW
					g.drawOriginalNumber(screen, v,
						listFieldRight(fields, col)-digits*gfx.DigitWidth, y, digits, c)
					continue
				}
				g.td.Draw(screen, cell,
					listFieldRight(fields, col)-textdraw.StringWidth(cell), y, c)
				continue
			}
			g.td.Draw(screen, fitCell(g.td.Text(cell), listCellRoom(fields, col)),
				listFieldX(fields, col), y, c)
		}
	}

	g.drawListScrollbar(screen, l, dim, ink)
}

// drawListScrollbar 照實機量測畫（docs/spec/38 §1.6）：標題列那格純黑、
// ▲▼ 是 3D 綠鈕（面＝色 13、高光＝色 2、影＝黑）、槽純黑、
// **滑塊是比例式的 3D 綠塊**（白頂／白左、色 1 右／底），
// 高度＝槽高×可見/總列、位置＝槽高×Top/總列（orig-w3-target 實測
// 21 列時滑塊 60px／槽 128px）。整頁剛好裝滿時滑塊佔滿全槽，
// 看起來像「沒有滑塊」——武將一覽第一輪就是這樣誤判的。
func (g *game) drawListScrollbar(screen *ebiten.Image, l *listwin.List,
	dim, ink color.RGBA) {

	face := g.paletteInk(13, color.RGBA{130, 162, 97, 255})
	hi := g.paletteInk(strategyInkNormal, chrome.Paper) // 高光是白（色 15），量測 m15
	grey := g.paletteInk(2, color.RGBA{162, 178, 178, 255})
	shade := g.paletteInk(1, color.RGBA{97, 113, 113, 255})
	black := color.RGBA{0, 0, 0, 255}
	fill := func(r image.Rectangle, c color.RGBA) {
		vector.DrawFilledRect(screen, float32(r.Min.X), float32(r.Min.Y),
			float32(r.Dx()), float32(r.Dy()), c, false)
	}
	// 標題列左邊那格：純黑。
	fill(image.Rect(listWinX, listWinY, listWinX+listScrollW, listWinY+listRowH), black)
	// 槽：純黑、左右黑邊之間放**比例式滑塊**（orig-w3-target 逐像素）。
	white := hi
	track := listScrollTrackRect()
	fill(track, black)
	slotH := track.Dy()
	thumbY, thumbH := track.Min.Y, slotH
	if l != nil && len(l.Rows) > l.Height && l.Height > 0 {
		thumbH = slotH * l.Height / len(l.Rows)
		thumbY = track.Min.Y + slotH*l.Top/len(l.Rows)
		if thumbY+thumbH > track.Max.Y {
			thumbY = track.Max.Y - thumbH
		}
	}
	// 滑塊的 3D：白頂列＋白左欄、色 1 右欄＋底列、面綠（色 13）。
	fill(image.Rect(track.Min.X+1, thumbY, track.Max.X-2, thumbY+thumbH), face)
	fill(image.Rect(track.Min.X+1, thumbY, track.Min.X+2, thumbY+thumbH), white)
	fill(image.Rect(track.Min.X+2, thumbY, track.Max.X-2, thumbY+1), white)
	fill(image.Rect(track.Max.X-3, thumbY+1, track.Max.X-2, thumbY+thumbH), shade)
	fill(image.Rect(track.Min.X+2, thumbY+thumbH-1, track.Max.X-3, thumbY+thumbH), shade)
	// ▲▼ 鈕（m10/m11 逐像素）：高光 2px 厚（頂兩列＋左兩欄）、
	// 右緣黑從第 2 列起、底列只留左角一點高光；箭頭十列
	// 寬 2,2,4,4,6,6,8,8,10,10、以 x＝Min+8 為中線；
	// ▲ 的箭頭從第 3 列起、▼ 從第 5 列起。
	button := func(r image.Rectangle, up bool) {
		fill(r, face)
		fill(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+1), hi)
		fill(image.Rect(r.Min.X, r.Min.Y, r.Min.X+1, r.Max.Y-1), hi)
		fill(image.Rect(r.Min.X+1, r.Min.Y+1, r.Max.X-1, r.Min.Y+2), grey)
		fill(image.Rect(r.Min.X+1, r.Min.Y+1, r.Min.X+2, r.Max.Y-1), grey)
		fill(image.Rect(r.Max.X-2, r.Min.Y+2, r.Max.X-1, r.Max.Y-2), shade)
		fill(image.Rect(r.Min.X+2, r.Max.Y-2, r.Max.X-1, r.Max.Y-1), shade)
		fill(image.Rect(r.Max.X-1, r.Min.Y+1, r.Max.X, r.Max.Y), black)
		fill(image.Rect(r.Min.X+1, r.Max.Y-1, r.Max.X, r.Max.Y), black)
		fill(image.Rect(r.Min.X, r.Max.Y-1, r.Min.X+1, r.Max.Y), hi)
		for i := 0; i < 10; i++ {
			w := (i/2)*2 + 2
			y := r.Min.Y + 2 + i
			if !up {
				w = 10 - (i/2)*2
				y = r.Min.Y + 4 + i
			}
			fill(image.Rect(r.Min.X+8-w/2, y, r.Min.X+8-w/2+w, y+1), black)
		}
	}
	button(listScrollUpRect(), true)
	button(listScrollDownRect(), false)
}

// listFieldsFor 取這張一覽表的欄位定義；沒設過就用武將那一組。
func listFieldsFor(l *listwin.List) []listField {
	if l == nil {
		return nil
	}
	latin := uiLang.Lang().Latin()
	switch l.Kind {
	case listwin.Cities:
		if latin {
			return latinFieldsCities
		}
		return listFamilyCities.fields()
	case listwin.Corps:
		if latin {
			return latinFieldsCorps
		}
		return listFamilyCorps.fields()
	case listwin.Factions:
		if latin {
			return latinFieldsFactions
		}
		return listFamilyFactions.fields()
	default:
		if latin {
			return latinFieldsGenerals
		}
		return listFamilyGenerals.fields()
	}
}

// setLanguage 換語言。**換的是呈現，不是遊戲**——世界狀態、時鐘、
// 存檔一律不動（docs/spec/86 §4）。
func (g *game) setLanguage(lang uitext.Language) error {
	p, err := langpack.Load(lang, g.fontDir)
	if err != nil {
		return err
	}
	uiLang = p.Text
	if !g.talkPinned && g.lib != nil {
		// 母本語系回到原版解出來的那一份（含校訂）。
		g.lib.Talk = g.talkBase
		if p.Talk != nil {
			g.lib.Talk = p.Talk
		}
	}
	if p.Font == nil && lang != uitext.ZhHant {
		log.Printf("⚠ %s 的全形字型載不到；缺的字會畫成方框", lang)
	}
	p.Apply(g.td)
	// 標題是依欄界生成的，要在新的詞表上重新登記（docs/spec/85 §4）。
	registerLatinListTitles()
	g.langNotice = languageName(lang)
	g.langNoticeAt = 0
	return nil
}

// languageName 是切換時畫在畫面上的提示，**用該語言自己的寫法**。
// 名字取自 `langpack.Choices`——與啟動殼層、手機面板同一份。
func languageName(lang uitext.Language) string {
	for _, c := range langpack.Choices {
		if c.Lang == lang {
			return c.Name
		}
	}
	return string(lang)
}

// registerLatinListTitles 把四個家族的標題換成依半形欄界生成的那一條
// （docs/spec/85 §4）。登記進語系詞表，畫的時候由 Drawer 的翻譯鉤子換掉，
// 呼叫端仍然只認得原本那條中文標題。
func registerLatinListTitles() {
	if !uiLang.Lang().Latin() {
		return
	}
	for _, t := range []struct {
		family listFamily
		fields []listField
		labels []string
	}{
		{listFamilyGenerals, latinFieldsGenerals, latinLabelsGenerals},
		{listFamilyCities, latinFieldsCities, latinLabelsCities},
		{listFamilyCorps, latinFieldsCorps, latinLabelsCorps},
		{listFamilyFactions, latinFieldsFactions, latinLabelsFactions},
	} {
		uiLang.Add(t.family.Title, latinTitle(t.fields, t.labels))
	}
}

func pressed(k ebiten.Key) bool { return inpututil.IsKeyJustPressed(k) }

// 大地圖的四季配樂：月 → 曲號。原版的表在 `cs:9309h`，
// 由 `sub_19321` 以「月 − 1」查（docs/re/58 §2）。
// **這不是照聽感排的**——同一個月份索引還有第二張表決定季節調色盤，
// 兩張表逐月吻合，那是「曲 2–5 是四季」的交叉驗證。
// musicTrack 把目前的狀態交給 `internal/rules/bgm`。
//
// ⭐ **規則不在這裡**：手機端要一樣的行為，抄第二份會長出差異
//（CLAUDE.md §7 第 6 條）。這一支只負責把 `game` 的狀態翻成 `bgm.Scene`。
func (g *game) musicTrack() string {
	scene := bgm.Scene{
		Launcher: g.launcher != nil,
		Ending:   g.endingActive(),
	}
	if g.world == nil {
		return bgm.Track(scene)
	}
	scene.GameOver = g.world.Outcome() != state.InProgress
	scene.Message = g.messageActive() || g.adviseActive()
	scene.Month = g.world.Clock.Month
	if g.battleActive() {
		b := bgm.Battle{Field: battlefield.FieldBase}
		if p := g.world.PendingBattle(); p != nil {
			b.Field = g.battle.FieldNumber(p.Node, p.Mode == combat.Siege)
			b.PlayerAttacks = p.Attacker >= 0 && p.Attacker < len(g.world.Corps) &&
				g.world.Corps[p.Attacker].Faction == g.world.Player
		}
		scene.Battle = &b
	}
	return bgm.Track(scene)
}

func (g *game) updateMusic() {
	if name := g.musicTrack(); name != "" {
		g.sound.PlayMusic(name)
	}
}

func (g *game) Update() error {
	g.frame++
	// 截圖模式要等 Draw 真正取到像素後才結束；只用 `frame > shotAt`
	// 會在高更新速率下跳過那一幀，讓 packaged smoke 沒有 PNG 卻仍 exit 0。
	if g.shotPath != "" && g.shotDone {
		return ebiten.Termination
	}
	// 錄滿了才結束，而且與截圖模式同一個理由放在 Update：
	// Draw 取到像素之後要等下一次 Update 才收工，否則最後一張可能還沒寫完。
	if g.recDone {
		return ebiten.Termination
	}
	// [HARD] ESC 只取消／關視窗，F10 才離開（CLAUDE.md §10）。
	if g.quitting {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			px, py := ebiten.CursorPosition()
			if hit, yes := hitTestYesNo(quitDialogX, quitDialogY, px, py); hit {
				if yes {
					return ebiten.Termination
				}
				g.quitting = false
			}
			return nil
		}
		switch {
		case pressed(ebiten.KeyArrowUp), pressed(ebiten.KeyArrowDown):
			g.quitYes = !g.quitYes
		case pressed(ebiten.KeyY):
			return ebiten.Termination
		case pressed(ebiten.KeyEnter):
			if g.quitYes {
				return ebiten.Termination
			}
			g.quitting = false
		// ⚠ 右鍵也要能取消。這個對話框先前只認 N／ESC，是 2026-08-23
		// 那一輪掃七個面板時漏掉的第八個（docs/spec/73 §2）——它不在
		// `docs/spec/73` 的清單裡，因為它不是「面板」而是離開確認。
		case pressed(ebiten.KeyN), g.cancelled():
			g.quitting = false
		}
		return nil
	}
	// F9 循環切語言（remake 差異，docs/spec/86 §4）。原版沒有這個鍵，
	// 而系統選單是逐像素對過的，不能加第五列。
	if pressed(ebiten.KeyF9) {
		if err := g.setLanguage(langpack.Next(uiLang.Lang())); err != nil {
			log.Printf("切換語言失敗：%v", err)
		}
		return nil
	}
	if pressed(ebiten.KeyF10) {
		g.quitting, g.quitYes = true, false // 預設停在 ＮＯ
		return nil
	}
	g.updateMusic()
	if g.launcher != nil {
		return g.updateLauncher()
	}
	if g.world != nil && g.world.Outcome() != state.InProgress {
		// 勝利先放結局過場，放完（或載不到素材）才回到結果對話框。
		g.maybeBeginEnding()
	}
	if g.endingActive() {
		return g.updateEnding()
	}
	if g.world != nil && g.world.Outcome() != state.InProgress {
		return g.updateOutcome()
	}
	// ⭐ **訊息先於戰場**：原版遭遇時 `sub_14EB9`（訊息框）在 `sub_11B5A`
	// （開戰術畫面）之前，玩家按掉才看到戰場（docs/spec/105）。
	// 戰略層在戰鬥中不跑，所以這裡不會有別的訊息插隊。
	if g.messageActive() && g.battleActive() {
		g.updateMessageOnly()
		return nil
	}
	// 戰場畫面獨佔一切——戰鬥中不能停時間，也不能開別的視窗。
	if g.battleActive() {
		g.updateBattle()
		return nil
	}
	// 事件前置報告可能與外交／撥款 pending 同一個 tick 產生；原版先
	// 顯示 TALK，再讓玩家進入選擇，因此通知 modal 優先於這兩個選單。
	if g.messageActive() {
		g.updateMessageOnly()
		return nil
	}
	// 事件 2／3 外交三選一也會凍結戰略時間。
	if g.world.PendingDiplomacy() != nil {
		g.updateDiplomacy()
		return nil
	}
	// 事件 4／5 內政官／外交官撥款也會凍結戰略時間。
	if g.world.PendingFunding() != nil {
		g.updateFunding()
		return nil
	}
	// 勢力選擇視窗是模態的（原版 sub_15AFC 自己跑一個等待迴圈，
	// 只認自己的熱區 0x20 與右鍵，docs/spec/35 §2.5.2）。
	if g.updateFactionPicker() {
		return nil
	}
	// 進言流程是模態的，優先吃輸入。
	if g.adviseActive() {
		g.updateAdvise()
		return nil
	}
	// 存檔／讀取是模態視窗，不能讓背景的命令鍵穿透。
	if g.saveUI.active {
		g.updateSaveUI()
		return nil
	}
	// 系統視窗內才接受存檔快捷鍵；避免在地圖上誤觸而改變狀態。
	if g.hudOpen(hudSystem) {
		switch {
		case pressed(ebiten.KeyS):
			g.beginSaveUI(saveWrite)
			return nil
		case pressed(ebiten.KeyL):
			g.beginSaveUI(saveRead)
			return nil
		}
	}
	// 一覽表開著時吃掉所有輸入 —— 它是模態的（說明書 3.8 的兩段式操作
	// 只有在獨佔輸入時才成立）。
	// 行軍指示的三選一是模態的，排在最前面（它會蓋在一覽表上）。
	if g.marchMode.active {
		g.updateMarchMode()
		return nil
	}
	// 編成畫面排在一覽表**之前**：原版的編成視窗是畫在武將一覽上面的
	// （docs/re/30 §1），一覽表留在背景但不吃輸入。
	if g.form.active {
		g.updateForm()
		return nil
	}
	if g.list != nil {
		g.updateListUI()
		return nil
	}
	// 財政畫面是模態的，優先吃輸入。
	if g.finance.active {
		g.updateFinance()
		return nil
	}
	// 據點情報視窗是模態的，優先吃輸入。
	if g.cityInfo.active {
		g.updateCityInfo()
		return nil
	}
	// 軍團情報視窗是模態的，優先吃輸入。
	if g.corpsInfo.active {
		g.updateCorpsInfo()
		return nil
	}
	// 自然策略頂端八格只在沒有 active modal／戰鬥／啟動器時接收滑鼠。
	// 上面的 return 順序是輸入隔離閘；winSystem 是唯一非 resident 的原生
	// 視窗，也必須阻止命令列點擊穿透。游標 hover 不改狀態，因為目前沒有
	// 足夠 DOS/V 證據證明原版頂端八格的 hover highlight。
	// 橫幅右側五格開關：**左鍵開、右鍵關**（docs/spec/13 §2.3）。
	// 它排在系統視窗的閘**之前**——不然系統視窗一開就再也關不掉。
	// 位置上不會衝突：開關在 y<32，系統視窗在 y 112–304。
	if left := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft); left ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		if i, ok := hitTestHUDSwitch(x, y); ok {
			if w := hudSwitchWindow(i); w != 0 {
				g.hudSet(w, left)
			}
			return nil
		}
		// 縮小地圖圖例的右半格（熱區 0x17 ＝ (536,168,96,16)）開勢力選擇視窗。
		if left && g.hudOpen(hudMinimap) && hitTestMinimapLegend(x, y) {
			g.factionPicker = true
			return nil
		}
		// 熱區 0x16：點地圖區把鏡頭捲過去（docs/spec/35 §2.5.1）。
		if left && g.hudOpen(hudMinimap) {
			if col, row, ok := minimapWorldAt(x, y); ok {
				g.centreCamOn(col, row)
				return nil
			}
		}
	}
	// 系統選單開著時，那六列吃滑鼠。**原版的六個 handler 沒讀**
	// （docs/re/55 §4），所以這裡的接法是照標籤字面意思的 remake 差異：
	// 左鍵 +1／右鍵 −1 調速度，兩個 ＯＫ 列接既有的存讀與離開確認。
	if g.hudOpen(hudSystem) {
		if left := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft); left ||
			inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			x, y := ebiten.CursorPosition()
			if row, ok := hitTestSystemRow(x, y); ok {
				g.dispatchSystemRow(row, left)
				return nil
			}
		}
	}
	if !g.hudOpen(hudSystem) {
		if g.hudOpen(hudCommand) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if command, ok := hitTestNaturalCommand(x, y); ok {
				g.dispatchNaturalCommand(command)
				return nil
			}
		}
	}
	for _, key := range []ebiten.Key{
		ebiten.KeyP, ebiten.KeyJ, ebiten.KeyF, ebiten.KeyA,
		ebiten.KeyC, ebiten.KeyT, ebiten.KeyG, ebiten.KeyK,
	} {
		if pressed(key) {
			command, ok := strategyCommandForShortcut(key)
			if ok {
				g.dispatchNaturalCommand(command)
				return nil
			}
		}
	}
	switch {
	case pressed(ebiten.KeyM):
		g.beginMarch()
		return nil
	}
	// 松崗繁中版會在游標未移動且沒有新輸入時才走 sub_11CD0 的自然世界
	// 迴圈。上方會開啟 modal 的命令已經 return；這裡保留仍可穿透到
	// timeRuns 的 UI 操作，確保同一 frame 不會一邊接受命令一邊推進日期。
	inputActive := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) ||
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) ||
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	if pressed(ebiten.KeyEscape) {
		inputActive = true
		// 由上而下關掉最上面那個開著的視窗。
		for i := hudSwitchN - 1; i >= 0; i-- {
			if w := hudSwitchWindow(i); w != 0 && g.hudOpen(w) {
				g.hudSet(w, false)
				break
			}
		}
	}
	// remake 差異：原版只有滑鼠點橫幅那五格，鍵盤 1–4 是自己加的。
	for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4} {
		if pressed(k) {
			inputActive = true
			w := hudSwitchWindow(i)
			g.hudSet(w, !g.hudOpen(w))
		}
	}
	// ＋／− 調速度：戰術畫面調戰術速度，其餘調戰略速度。
	for i, k := range []ebiten.Key{ebiten.KeyMinus, ebiten.KeyEqual} {
		if pressed(k) {
			inputActive = true
			g.adjustSpeed(g.battleActive(), []int{-1, 1}[i])
		}
	}

	step := 1
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		step = 8
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		inputActive = true
		g.camX += step
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		inputActive = true
		g.camX -= step
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		inputActive = true
		g.camY += step
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		inputActive = true
		g.camY -= step
	}
	g.clampCam()

	if !g.timeRuns() {
		return nil
	}
	// 第一個觀測 frame 也視為尚未 idle；原版必須先得到一筆穩定座標，
	// 才會設 byte_198A3 的 bit 7。任何游標移動或命令都會停住這次
	// 據點／軍團／物件／時鐘更新，下一個靜止 frame 才可恢復。
	cursorX, cursorY := ebiten.CursorPosition()
	if !g.idleGate.Allows(cursorX, cursorY, inputActive) {
		return nil
	}
	// 系統視窗開著時時間停止（說明書 3.1，docs/spec/13 §2.4）。
	// 另外三個視窗開著時時間照走——這是原版明講的差別。
	// **戰略層沒有別的暫停鍵**：檔位 0 在原版是最快不是暫停。
	if g.hudOpen(hudSystem) {
		g.world.AdvanceMapObjects(g.rng)
		return nil
	}
	// 這一個畫面推進幾個子刻，由節流累加器決定（docs/spec/34）。
	steps := g.strategyThrottle.Steps(g.speed, 1, speed.HighSpeedStrategy)
	if steps == 0 {
		g.world.AdvanceMapObjects(g.rng)
		return nil
	}
	for i := 0; i < steps; i++ {
		// 第一個規則 tick 跑完整的原版順序「據點／軍團／物件／時鐘」；
		// 同一畫面的額外 speed tick 不重跑物件。
		var ev state.Event
		if i == 0 {
			ev = g.world.TickMap(g.rng)
		} else {
			ev = g.world.Tick(g.rng)
		}
		if ev.Settled {
			g.lastEvent = "月結"
			if n := len(ev.Disaster); n > 0 {
				g.lastEvent += fmt.Sprintf("　災害%d", n)
			}
			if ev.Storm != nil {
				g.lastEvent += "　暴風雨"
			}
		}
		for _, i := range ev.Eliminated {
			g.lastEvent = big5(g.world.LordName(i)) + " 滅亡"
		}
		g.reportCorps(ev)
		g.reportStrategy(ev)
		g.enqueueEventMessages(ev)
		if g.world.Outcome() != state.InProgress {
			break
		}
		if g.messageActive() {
			break
		}
	}
	return nil
}

func (g *game) clampCam() {
	if g.camX < 0 {
		g.camX = 0
	}
	if g.camY < 0 {
		g.camY = 0
	}
	if m := 384 - viewCols; g.camX > m {
		g.camX = m
	}
	if m := 256 - viewRows; g.camY > m {
		g.camY = m
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	if g.launcher != nil {
		g.drawLauncher(screen)
		g.maybeSaveShot(screen)
		return
	}
	if g.endingActive() {
		g.drawEnding(screen)
		g.maybeSaveShot(screen)
		return
	}
	if g.world == nil {
		return
	}
	// ⭐ **訊息還在時不畫戰場**：原版遭遇的訊息框跳在大地圖上，
	// 按掉才換成戰術畫面（`sub_14EB9` → `sub_11B5A`，docs/spec/105）。
	// remake 兩件事在同一個 tick 發生，所以由這一格決定先看到哪一個。
	if g.world.Outcome() == state.InProgress && g.battleActive() && !g.messageActive() {
		g.drawBattle(screen)
		g.maybeSaveShot(screen)
		return
	}
	season := int(g.world.Clock.Season())

	// 大地圖鋪滿橫幅以下的全部畫面。四季調色盤直接吃時鐘算出來的季節——
	// 所以畫面會隨遊戲時間換季，不需要另外驅動。
	if img, err := g.lib.RenderWorldMarked(g.camX, g.camY, viewCols, viewRows,
		season, g.cityMarks(), g.corpsMarks()); err == nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, strategyMapY)
		screen.DrawImage(ebiten.NewImageFromImage(img), op)
	}
	g.drawDisasterOverlay(screen)

	// 橫幅是原版美術（ICONGRF 段 0），不是自己畫的長條。
	// 上面已經印好「臥竜伝」與「年 月 日」，這裡只填數字。
	// ⚠ 橫幅寫的是日文「臥竜伝」——松崗版沒有重繪這張圖
	// （docs/reference/03），中文化要補的缺口之一。
	if b, err := g.lib.Banner(season); err == nil {
		screen.DrawImage(ebiten.NewImageFromImage(b), &ebiten.DrawImageOptions{})
	} else {
		vector.DrawFilledRect(screen, 0, 0, screenW, bannerH,
			color.RGBA{32, 24, 16, 255}, false)
		g.td.Draw(screen, "臥龍傳", 10, 8, color.RGBA{240, 200, 120, 255})
	}
	c := g.world.Clock
	// 原版的日期與橫幅上的「年 月 日」同色；硬寫的白色在實機對拍上差得出來。
	ink := color.RGBA{240, 240, 230, 255}
	if col, err := g.lib.PaletteColor(season, bannerDateInk); err == nil {
		ink = col
	}
	g.drawBannerNumber(screen, c.Year, bannerYearRight, ink)
	g.drawBannerNumber(screen, c.Month, bannerMonthRight, ink)
	g.drawBannerNumber(screen, c.Day, bannerDayRight, ink)

	// 自然策略 HUD 是原版主畫面的固定骨架；四個視窗仍是可獨立切換的暫存層。
	g.drawNaturalStrategyHUD(screen)

	// 事件列只在有事的時候出現，而且畫成一個視窗而不是貼在畫面底部的黑條。
	// ⭐ **六秒後自己消失**：與手機版同一個行為（docs/spec/88 §2）。
	if g.lastEvent != g.lastEventShown {
		g.lastEventShown, g.lastEventAt = g.lastEvent, g.frame
	}
	if g.lastEvent != "" && g.frame-g.lastEventAt < eventLineFrames {
		w := g.td.Width(g.lastEvent) + 4*chrome.Tile
		w = (w/chrome.Tile + 1) * chrome.Tile
		x := (screenW - w) / 2 / chrome.Tile * chrome.Tile
		// 與 TALK 框同一族：黑底金框（docs/spec/88 §1）。
		g.chrome.Window(screen, x, screenH-40, w, 32, chrome.Blank)
		g.td.Draw(screen, g.lastEvent, x+2*chrome.Tile, screenH-32, chrome.Paper)
	}
	if !g.td.Available() {
		g.td.Draw(screen, "（未載入字型）", 8, screenH-20,
			color.RGBA{240, 140, 140, 255})
	}

	if g.list != nil {
		g.drawList(screen)
	}
	g.drawForm(screen)
	g.drawMarchMode(screen)
	g.drawFinance(screen)
	g.drawCityInfo(screen)
	g.drawCorpsInfo(screen)
	g.drawAdvise(screen)
	g.drawSaveUI(screen)
	if choice := g.world.PendingDiplomacy(); choice != nil {
		g.drawDiplomacy(screen, choice)
	}
	if choice := g.world.PendingFunding(); choice != nil {
		g.drawFunding(screen, choice)
	}
	if g.messageActive() {
		g.drawMessage(screen)
	}
	if g.world.Outcome() != state.InProgress {
		g.drawOutcome(screen)
	}

	if g.quitting {
		// 原版版面的 ＹＥＳ／ＮＯ 對話框（docs/spec/26），居中。
		g.drawYesNo(screen, quitDialogX, quitDialogY, "確定離開？", g.quitYes)
	}
	g.maybeSaveShot(screen)
}

// drawDisasterOverlay 把 state 的 runtime 災害 marker 接到戰略地圖。
//
// 事件 12 的火災／暴動物件現在使用 MMAP.MCH 原版矩陣與 16×16 平面圖塊：
// sub_123FF 寫入 object type 1／2，sub_12533 以 CS:985Ah 查表取 8 相位，
// 再由 loc_1D51F 置中寫入 40×23 的戰略地圖 cell buffer。remake 在同一個
// 地圖格座標上合成圖像；相位時鐘是呈現層的固定 frame clock，因為目前
// runtime object 的 [si+0F]／[si+0C] 已由 state 保存為非存檔欄位，
// 由 map-loop timer 驅動；這裡只負責把本次 render 取得的相位交給 MCH。
//
// 暴風雨事件 11 的已證實 handler 只更新據點 +0x15 與範圍，不會呼叫
// sub_123FF；所以它仍畫範圍輪廓。若 MCH 缺檔，才退回低干擾向量 marker，
// 不把 fallback 說成原版 parity。
func (g *game) drawDisasterOverlay(screen *ebiten.Image) {
	objects := g.world.RenderDisasterObjects()
	objectByCity := make(map[int]state.DisasterObjectSnapshot, len(objects))
	for _, object := range objects {
		// 原版掃描順序是 slot 0→31；同城重複時第一筆先被呈現層取用。
		if _, exists := objectByCity[object.City]; !exists {
			objectByCity[object.City] = object
		}
	}
	for cityID, city := range g.world.Cities {
		marker, ok := g.world.DisasterMarkerAt(cityID)
		if !ok {
			continue
		}
		x := (city.X - g.camX) * world.TileSize
		y := strategyMapY + (city.Y-g.camY)*world.TileSize
		phase := uint8(1)
		if object, exists := objectByCity[cityID]; exists {
			phase = object.Phase
		}
		if img := g.disasterImage(marker.Kind, int(phase)); img != nil {
			w, h := img.Bounds().Dx(), img.Bounds().Dy()
			drawX := x + world.TileSize/2 - w/2
			drawY := y + world.TileSize/2 - h/2
			if drawX+w <= 0 || drawX >= strategyMapW || drawY+h <= strategyMapY || drawY >= screenH {
				continue
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(drawX), float64(drawY))
			screen.DrawImage(img, op)
			continue
		}
		if x+world.TileSize <= 0 || x >= strategyMapW || y+world.TileSize <= strategyMapY || y >= screenH {
			continue
		}
		g.drawDisasterMarker(screen, x, y, marker)
	}
	if area, ok := g.world.StormAreaSnapshot(); ok {
		g.drawStormArea(screen, area)
	}
}

func disasterObjectType(kind economy.Disaster) (int, bool) {
	switch kind {
	case economy.Fire:
		return 1, true // sub_134B1 的 AH=1
	case economy.Riot:
		return 2, true // sub_134B1 的 AH=2
	default:
		return 0, false
	}
}

func (g *game) disasterImage(kind economy.Disaster, phase int) *ebiten.Image {
	objectType, ok := disasterObjectType(kind)
	if !ok || g.lib == nil || g.lib.MCH == nil || phase < 0 || phase >= 8 {
		return nil
	}
	season := 0
	if g.world != nil {
		season = int(g.world.Clock.Season())
	}
	key := (((season*3)+objectType)*8 + phase)
	if g.disasterImages == nil {
		g.disasterImages = make(map[int]*ebiten.Image)
	}
	if img, exists := g.disasterImages[key]; exists {
		return img
	}
	pattern, ok := g.lib.MCH.PatternFor(objectType, phase)
	if !ok {
		return nil
	}
	img, err := g.lib.MCH.RenderPattern(pattern, g.lib.Palette, season)
	if err != nil {
		return nil
	}
	result := ebiten.NewImageFromImage(img)
	g.disasterImages[key] = result
	return result
}

func (g *game) drawDisasterMarker(screen *ebiten.Image, x, y int, marker state.DisasterMarker) {
	// 深色底讓標記在四季地圖上都可辨識；Level 只控制亮度，不被解讀為幀數。
	alpha := uint8(180 + int(marker.Level&3)*18)
	dark := color.RGBA{24, 18, 20, 220}
	vector.DrawFilledRect(screen, float32(x+9), float32(y+1), 7, 7, dark, false)
	switch marker.Kind {
	case economy.Fire:
		vector.DrawFilledRect(screen, float32(x+10), float32(y+2), 5, 5,
			color.RGBA{246, 92, 32, alpha}, false)
		vector.DrawFilledRect(screen, float32(x+12), float32(y), 2, 4,
			color.RGBA{255, 210, 54, alpha}, false)
	case economy.Riot:
		vector.DrawFilledRect(screen, float32(x+10), float32(y+2), 5, 5,
			color.RGBA{196, 56, 170, alpha}, false)
		vector.DrawFilledRect(screen, float32(x+11), float32(y+1), 3, 7,
			color.RGBA{255, 224, 92, alpha}, false)
	case economy.Storm:
		vector.DrawFilledRect(screen, float32(x+10), float32(y+2), 5, 5,
			color.RGBA{64, 180, 226, alpha}, false)
		vector.DrawFilledRect(screen, float32(x+12), float32(y+1), 2, 7,
			color.RGBA{190, 242, 255, alpha}, false)
	}
}

func (g *game) drawStormArea(screen *ebiten.Image, area economy.StormArea) {
	// 範圍本身只畫輪廓，避免把地圖地形整片蓋掉；據點 marker 才是目前
	// 受影響據點的狀態提示。座標與 StormArea 同樣以地圖格為單位。
	x := (area.MinX - g.camX) * world.TileSize
	y := strategyMapY + (area.MinY-g.camY)*world.TileSize
	w := (area.MaxX - area.MinX + 1) * world.TileSize
	h := (area.MaxY - area.MinY + 1) * world.TileSize
	if w <= 0 || h <= 0 || x+w <= 0 || x >= strategyMapW || y+h <= strategyMapY || y >= screenH {
		return
	}
	c := color.RGBA{80, 194, 232, 150}
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), 2, c, false)
	vector.DrawFilledRect(screen, float32(x), float32(y+h-2), float32(w), 2, c, false)
	vector.DrawFilledRect(screen, float32(x), float32(y), 2, float32(h), c, false)
	vector.DrawFilledRect(screen, float32(x+w-2), float32(y), 2, float32(h), c, false)
}

// maybeSaveShot 在達到目標幀後的第一個 Draw 取像素，下一次 Update 才結束。
// 它不依賴 Draw 與 Update 恰好一對一，對 Xvfb 與不同平台的繪製節奏都安全。
//
// 逐幀錄製（docs/spec/71）也掛在這裡：`Draw` 有四個 return 點，
// 每一個都呼叫這一支，掛在別處會漏掉戰場或結局那幾條路徑。
func (g *game) maybeSaveShot(screen *ebiten.Image) {
	if g.rec != nil && !g.recDone {
		g.recDone = g.rec.shot(screen)
	}
	if g.shotPath == "" || g.shotDone || g.frame < g.shotAt {
		return
	}
	g.shotDone = g.saveShot(screen)
}

// saveShot 把目前畫面寫成 PNG。只在截圖模式用。
func (g *game) saveShot(screen *ebiten.Image) bool {
	b := screen.Bounds()
	img := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			img.Set(x, y, screen.At(x, y))
		}
	}
	f, err := os.Create(g.shotPath)
	if err != nil {
		log.Printf("⚠ 截圖失敗：%v", err)
		return false
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		log.Printf("⚠ 截圖編碼失敗：%v", err)
		return false
	}
	if g.world == nil {
		// 啟動殼層的截圖（`-open-naming`）：還沒有世界，沒有日期可印。
		log.Printf("截圖 → %s（第 %d 幀，啟動殼層）", g.shotPath, g.frame)
		return true
	}
	log.Printf("截圖 → %s（第 %d 幀，%d年%d月%d日）",
		g.shotPath, g.frame, g.world.Clock.Year, g.world.Clock.Month, g.world.Clock.Day)
	return true
}

func (g *game) Layout(int, int) (int, int) { return screenW, screenH }

// seasonName 用中文的季節名。clock.Season 的 String() 已經是中文，
// 這裡只是把它包成一個明確的名字，讓畫面程式讀起來清楚。
func seasonName(s clock.Season) string { return s.String() }

// big5 把 internal/state 保留的原始位元組轉成 UTF-8。
// 解析層刻意保留原始位元組才能 round-trip，所以轉換發生在最外層。
// uiLang 是目前語系的轉換表（-lang，docs/spec/84）。
//
// 放成套件變數是刻意的：`big5` 是**原版資料轉成畫面文字的唯一入口**
// （人名、地名、君主名都走它），語系轉換掛在這裡就一次涵蓋全部，
// 不必去追 69 個呼叫點。zh-hant 時是 nil，Convert 是恆等。
var uiLang *uitext.Table

func big5(s string) string {
	if s == "" {
		return "－"
	}
	return uiLang.Convert(text.Decode([]byte(s), text.Big5))
}

// talkMessage 取一則原版 TALK.DAT 訊息並代入已知變數。
// TALK.DAT 的行尾是原版對話框硬斷行；事件列是 remake 的單行觀測列，
// 因此這裡只合併行，不宣稱已重現原版逐頁對話框。
func (g *game) talkMessage(index int, vars map[byte]string) string {
	lines, ok := g.talkLines(index, vars)
	if !ok {
		return ""
	}
	var b strings.Builder
	for _, line := range lines {
		b.WriteString(line)
	}
	return b.String()
}

// advisorName 回傳玩家所仕勢力的軍師名。0xFF 表示沒有軍師——
// 開新遊戲時玩家本人就是要填進那一格的人。
func (g *game) advisorName() string {
	f := g.world.Factions[g.world.Player]
	if f.Advisor < 0 || f.Advisor >= len(g.world.Generals) {
		return ""
	}
	return g.world.Generals[f.Advisor].Name
}

const (
	defaultTalkCorrections = "translations/corrections.json"
	defaultOrigDir         = "workplace/orig/dosv"
	defaultFontDir         = "workplace/eten"
)

// bundledTalkCorrectionsPath 保持「發行包不含完整原版文字表」的同時，讓
// corrections.json 可隨可執行檔安裝。先尊重明確環境設定與目前工作目錄，
// 再尋找一般 tar/zip 與 AppImage 的同包路徑；若都不存在則回傳預設相對路徑，
// 由 LoadWithOptions 顯示可診斷的失敗，而不是靜默跳過校訂。
// bundledTranslationPath 依 bundledTalkCorrectionsPath 的候選順序找
// translations/ 底下的語系檔；找不到回空字串（呼叫端 fallback 母本）。
func bundledTranslationPath(rel string) string {
	if fileExists(rel) {
		return rel
	}
	executable, err := os.Executable()
	if err != nil {
		return ""
	}
	dir := filepath.Dir(executable)
	for _, candidate := range []string{
		filepath.Join(dir, rel),
		filepath.Join(dir, "..", rel),
		filepath.Join(dir, "..", "share", "wolong-remake", rel),
	} {
		if fileExists(candidate) {
			return candidate
		}
	}
	return ""
}

func bundledTalkCorrectionsPath() string {
	if configured := os.Getenv("WOLONG_TALK_CORRECTIONS"); configured != "" {
		return configured
	}
	if fileExists(defaultTalkCorrections) {
		return defaultTalkCorrections
	}
	executable, err := os.Executable()
	if err != nil {
		return defaultTalkCorrections
	}
	dir := filepath.Dir(executable)
	for _, candidate := range []string{
		filepath.Join(dir, defaultTalkCorrections),
		filepath.Join(dir, "..", defaultTalkCorrections),
		filepath.Join(dir, "..", "share", "wolong-remake", defaultTalkCorrections),
	} {
		if fileExists(candidate) {
			return candidate
		}
	}
	return defaultTalkCorrections
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// 完整包裡原版資料與點陣字的目錄名（docs/spec/72 §2）。
const (
	bundledOrigDir  = "gamedata"
	bundledFontDir  = "fonts"
	bundledAudioDir = "audio"
)

// resolveDataDir 決定實際要用的資料目錄。
//
// `-orig`／`-font` 的預設值是 **repo 相對路徑**，解開的發行包裡不成立；
// 使用者跑 `./wlgame` 會靜靜地少掉字型，中文變成方框。這一支照
// bundledTalkCorrectionsPath 的形狀補上同一組退路。
//
// ⭐ 三條性質缺一不可（docs/spec/72 §3）：
//
//  1. **不覆蓋明講的旗標**——`value != def` 就直接回傳。對拍與驗收
//     全部明講 `-orig`，一個字都不受影響。
//  2. **repo 內行為不變**——`workplace/orig/dosv` 在的時候第二條就命中。
//  3. **都找不到就回預設值**，讓既有的載入器噴可診斷的錯。
//     ⚠ 不要靜默跳過——沉默的成功比失敗難發現。
// resolveAudioDir 決定音檔目錄（docs/spec/75 §2）。
//
// ⚠ **判準與 resolveDataDir 不同，不要湊成一支。** `-orig` 的預設值是一個
// repo 相對路徑，判準是「那個路徑不存在才去找」；`-audio` 的預設值是
// **空字串**，判準是「**沒給才去找**」。兩者湊在一起會讓其中一邊在某個
// 情況下錯。
//
// 找不到就回空字串，維持既有的靜音行為——完整包以外的批次不含音檔。
func resolveAudioDir(value string) string {
	if value != "" {
		return value
	}
	executable, err := os.Executable()
	if err != nil {
		return ""
	}
	dir := filepath.Dir(executable)
	for _, candidate := range []string{
		filepath.Join(dir, bundledAudioDir),
		filepath.Join(dir, "..", bundledAudioDir),
		filepath.Join(dir, "..", "share", "wolong-remake", bundledAudioDir),
	} {
		if dirExists(candidate) {
			return candidate
		}
	}
	return ""
}

func resolveDataDir(value, def, bundled string) string {
	if value != def {
		return value
	}
	if dirExists(def) {
		return def
	}
	executable, err := os.Executable()
	if err != nil {
		return def
	}
	dir := filepath.Dir(executable)
	for _, candidate := range []string{
		filepath.Join(dir, bundled),
		filepath.Join(dir, "..", bundled),
		filepath.Join(dir, "..", "share", "wolong-remake", bundled),
	} {
		if dirExists(candidate) {
			return candidate
		}
	}
	return def
}

func main() {
	dir := flag.String("orig", defaultOrigDir, "原版素材目錄（請自備）")
	scenPath := flag.String("scenario-file", "", "劇本檔路徑（預設 <orig>/SINARIO.DAT）")
	scenario := flag.Int("scenario", 0, "劇本編號 0–3（直接啟動／驗收用）")
	player := flag.Int("player", 0, "玩家所仕的勢力編號（直接啟動／驗收用）")
	directStart := flag.Bool("direct", false, "跳過一般玩家啟動殼層，直接啟動指定劇本／玩家（驗收用）")
	fontDir := flag.String("font", defaultFontDir, "倚天點陣字目錄（請自備）")
	// ⚠ 預設是空的（靜音）。Ebiten 的音訊錯誤沒有可查詢的 API，
	// 沒有音效裝置時 `RunGame` 會直接帶著 ALSA 的錯誤結束——
	// 無頭驗收與 CI 全部會掛。
	//
	// ⭐ 完整包會自己找執行檔旁邊的 `audio/`（docs/spec/75），
	// 所以「留白＝靜音」這道防護在**驗收模式**另外補（見 flag.Parse 之後）。
	audioDir := flag.String("audio", "", "ogg 音檔目錄（先跑 tools/bgm2ogg.sh 產生；留白＝靜音；完整包會自己找旁邊的 audio/）")
	speed := flag.Int("speed", 2, "戰略速度檔位 0–4（0 ＝ 最高速、4 ＝ 最低速）")
	tacticalSpeed := flag.Int("tactical-speed", 2, "戰術速度檔位 0–4（0 ＝ 最高速、4 ＝ 最低速）")
	seed := flag.Int("seed", -1, "驗收用固定亂數種子；負值時照原版以時鐘播種")
	shot := flag.String("shot", "", "跑 N 幀之後截圖到這個路徑就結束（驗收用）")
	framesDir := flag.String("frames-dir", "", "把每一張畫出來的圖寫成 fNNNNN.png（推廣片素材，docs/spec/71）")
	framesN := flag.Int("frames", 300, "配 -frames-dir：錄幾張就結束")
	shotFrames := flag.Int("shot-frames", 120, "截圖前先跑幾幀")
	saveFile := flag.String("save-file", "", "可寫的四槽存檔 overlay 路徑；一般啟動可選讀檔")
	loadSlot := flag.Int("load-slot", -1, "直接啟動時先從 -save-file 的第 N 槽（0–3）載入（驗收用；原版存檔也讀得動）")
	openWin := flag.Int("open-window", -1, "截圖前先打開第幾個視窗（0–3；−2 ＝ 三個常駐視窗；−3 ＝ 再加系統選單。對拍用）")
	camAt := flag.String("cam", "", "把大地圖鏡頭移到指定格 `X,Y`（對拍用；原版點過視窗開關之後鏡頭就不在開局位置了）")
	openList := flag.Bool("open-list", false, "截圖前先開武將一覽（驗收用；開著、無選取）")
	openFormPick := flag.Bool("open-form-pick", false, "截圖前停在編成的武將一覽（對拍用，與原版指令列 #3 剛開的狀態相同）")
	formPickRow := flag.Int("form-pick-row", 0, "配 -open-form：選候選清單的第 N 列當主將（對拍用）")
	lordCorpsFlag := flag.Bool("lord-corps", true, "允許把君主編成軍團長（docs/spec/76；對拍原版行為時給 false）")
	damageReportFlag := flag.Bool("siege-damage", false, "戰後結果多印一行攻城損害（remake 的驗收資訊，原版沒有；docs/spec/89）")
	openAdvise := flag.Bool("open-advise", false, "截圖前先跑到說服畫面（驗收用）")
	adviseMenu := flag.Bool("advise-menu", false, "單獨用：停在進言的五項選單；配 -open-advise：停在五選一的理由選單（驗收用）")
	adviseSortie := flag.Bool("advise-sortie", false, "截圖前跑「請求君主出陣」的三句定案畫面（驗收用）")
	adviseTarget := flag.Bool("advise-target", false, "截圖前停在進言→交戰的目標勢力清單（對拍用，docs/spec/90 §5.1）")
	openCities := flag.Bool("open-cities", false, "截圖前開據點一覽（對拍用，docs/spec/90 §5.1）")
	openFactions := flag.Bool("open-factions", false, "截圖前開勢力一覽（對拍用，docs/spec/90 §5.1）")
	openCityInfo := flag.Int("open-cityinfo", -2, "截圖前開第 N 個據點的情報卡（−1＝玩家首都；對拍用，docs/spec/90 §5.1）")
	openForm := flag.Bool("open-form", false, "截圖前先編一支軍團並開編成畫面（驗收用）")
	openCorps := flag.Bool("open-corps", false, "截圖前先編兩支軍團並開軍團一覽（驗收用）")
	corpsOnMap := flag.Bool("corps-on-map", false, "截圖前編一支軍團，**不開任何視窗**停在大地圖（docs/spec/74 §4.1）")
	marchTo := flag.Int("march-to", -1, "配 -corps-on-map：對那支軍團下行軍指示到據點 N；`-shot-frames` 推進的 tick 會讓它上路")
	openBattle := flag.Bool("open-battle", false, "截圖前先開一場野戰的戰術戰鬥（驗收用）")
	openSiege := flag.Bool("open-siege", false, "截圖前先開一場攻城的戰術戰鬥（驗收用）")
	openEnding := flag.Int("open-ending", -1, "直接跳到結局的第幾幕（0–11，驗收用）")
	openMarchMode := flag.Bool("open-march-mode", false, "截圖前停在行軍指示的三選一（驗收用）")
	openMarchList := flag.Bool("open-march-list", false, "截圖前編一支軍團並停在行軍目的地一覽（驗收用）")
	openNaming := flag.Bool("open-naming", false, "停在啟動殼層選君主那一頁並打開「自定」命名視窗（驗收用，docs/spec/104）")
	battleFF := flag.Bool("battle-ff", false, "配 -open-battle／-open-siege：截圖前先按下 `▶▶` 快轉（驗收用，docs/spec/102）")
	siegeNode := flag.Int("siege-node", -1, "指定攻城的戰場＝據點編號（驗收用，配 -open-siege）")
	siegeDefend := flag.Bool("siege-defend", false, "攻城時玩家當守方（原版會把戰場轉 180 度，docs/spec/56）")
	siegeCorps := flag.String("siege-corps", "", "拿**存檔裡現成的**兩支軍團開攻城戰：`攻,守`（編號用 -list-corps 看，docs/spec/90 §2.3）")
	battleSteps := flag.Int("battle-steps", 120, "截圖前推進幾個戰術 tick；0 ＝ 原版開場對白那一幀")
	listCorps := flag.Bool("list-corps", false, "把載入後還活著的軍團印出來（編號、勢力、主將、兵力）")
	battleCam := flag.String("battle-cam", "", "覆寫戰術鏡頭的世界格 `X,Y`（驗收用；原版初值是 36,14）")
	openFinance := flag.Bool("open-finance", false, "截圖前先開財政視窗（對拍用，docs/spec/14 §4）")
	financeAmount := flag.Int("finance-amount", -1, "配 -open-finance：再開第 N 列（0–3）的數值輸入器（docs/spec/78）")
	openMessage := flag.Bool("open-message", false, "截圖前先開玩家首都的暴風雨 TALK #70 通知（驗收用）")
	openTalkIndex := flag.Int("open-talk-index", -1, "截圖前直接開指定 TALK.DAT 槽位（驗收用）")
	openOutcome := flag.String("open-outcome", "", "只供截圖的敗北 modal fixture：trust 或 faction")
	talkJSON := flag.String("talk-json", "", "完整繁中 TALK JSON（研究用）")
	langFlag := flag.String("lang", "zh-hant", "語系：zh-hant／zh-hans／ja／en（docs/spec/84）")
	talkCorrections := flag.String("talk-corrections", bundledTalkCorrectionsPath(), "繁中 TALK 校訂覆蓋")
	flag.Parse()

	// 解開的完整包裡，預設的 repo 相對路徑不成立（docs/spec/72 §3）。
	*dir = resolveDataDir(*dir, defaultOrigDir, bundledOrigDir)
	*fontDir = resolveDataDir(*fontDir, defaultFontDir, bundledFontDir)
	*audioDir = resolveAudioDir(*audioDir)
	// ⭐ **驗收模式不出聲，但保留音效狀態**（docs/spec/29 §5.1）。
	// 沒有音效裝置時 Ebiten 的 `RunGame` 會直接帶著 ALSA 的錯誤結束，
	// 而截圖與逐幀錄製都跑在沒有音效卡的容器裡。
	//
	// ⚠ **不要清空音檔目錄。** 那樣 `Bank.Available()` 會變成 false，
	// 系統選單那一格就從「TYPE 1」變成「未接入」——驗收捷徑改到了
	// 被驗收的畫面（playtest/49 §2 的 272 px）。碰音效裝置的只有
	// `PlayMusic`／`PlayEffect`，關掉它們就夠了。
	//
	// 判準是**用途**不是環境：不是「有沒有顯示器」（Xvfb 有 DISPLAY
	// 卻沒有音效卡），也不是猜 `/dev/snd`（那是平台細節）。
	silentAudio := *shot != "" || *framesDir != ""
	flagVisit = func(fn func(string)) { flag.CommandLine.Visit(func(f *flag.Flag) { fn(f.Name) }) }

	lang, err := uitext.ParseLanguage(*langFlag)
	if err != nil {
		log.Fatal(err)
	}
	// 語系檔優先找執行檔旁邊的（完整包），再落回內嵌那一份（docs/spec/86 §2）。
	if p := bundledTranslationPath("translations/talk-en.json"); p != "" {
		langpack.SearchPaths = append([]string{filepath.Dir(p)}, langpack.SearchPaths...)
	}

	lib, err := library.LoadWithOptions(*dir, library.LoadOptions{TalkJSON: *talkJSON, TalkCorrections: *talkCorrections})
	if err != nil {
		log.Fatal(err)
	}
	for _, warning := range lib.Warns {
		log.Printf("⚠ %s", warning)
	}
	path := *scenPath
	if path == "" {
		path = filepath.Join(*dir, "SINARIO.DAT")
	}
	loadPath, err := savepath.InitialLoadPath(path, *saveFile)
	if err != nil {
		log.Fatal(err)
	}

	var ascii *cjk.ASCIIFont
	if a, err := cjk.LoadASCIIDir(*fontDir); err != nil {
		log.Printf("⚠ 載不到倚天半形字型（%v）", err)
	} else {
		ascii = a
	}
	gameRNG := rng.Now()
	if *seed >= 0 {
		gameRNG = rng.NewFixed(*seed)
		log.Printf("驗收固定亂數種子：%d", *seed)
	}
	g := &game{lib: lib, rng: gameRNG, speed: *speed, tacticalSpeed: *tacticalSpeed,
		lordCorps: true, // docs/spec/76：預設放行（remake 差異）
		td:       textdraw.New(nil, ascii),
		shotPath: *shot, shotAt: *shotFrames, origDir: *dir, sourceFile: path,
		rec: newRecorder(*framesDir, *framesN),
		saveFile: *saveFile, saveBase: path, sound: sound.Open(*audioDir)}
	g.sound.SetSilent(silentAudio)
	if *audioDir != "" && !g.sound.Available() {
		log.Printf("音檔目錄 %s 沒有 ogg，靜音跑。要有音樂請跑 tools/bgm2ogg.sh", *audioDir)
	}
	g.talkBase, g.talkPinned, g.fontDir = lib.Talk, *talkJSON != "", *fontDir
	if err := g.setLanguage(lang); err != nil {
		log.Fatal(err)
	}
	// 四個常駐視窗**預設全關**，這是原版數值：新遊戲流程的最後一行是
	// `sub_11A6E` 的 `mov cs:byte_198A6, 0`（docs/re/47 §3.3），
	// PC-98 實跑進到主畫面看到的也正是滿版地圖。玩家自己點橫幅右側
	// 那幾格把要看的視窗叫出來。
	g.hud = 0
	g.chrome = chrome.Load(lib, 0)
	if !g.chrome.Available() {
		log.Printf("⚠ 取不到 ICONGRF 段 3 的視窗外框，改畫純色框")
	}

	// `-open-naming` 要的是啟動殼層本身，即使帶了 `-shot` 也不走直啟。
	if (*directStart || directStartFlagWasPassed()) && !*openNaming {
		if err := g.startWorld(loadPath, *scenario, *player, true, loadPath == path); err != nil {
			log.Fatal(err)
		}
		// -load-slot 讓驗收路徑直接從存檔開局，不必走讀檔選單。
		// **原版存檔也讀得動**（`readSave` 找不到原生檔就退回
		// `state.LoadScenario`），所以可以拿原版存的那一份做同狀態對拍
		// （docs/spec/90 §2）。
		if *loadSlot >= 0 {
			if err := g.readSave(*loadSlot); err != nil {
				log.Fatalf("-load-slot %d 讀不起來：%v", *loadSlot, err)
			}
			// **空的槽要當錯誤**：載進來的是一塊沒有內容的區塊，
			// 解出來是一堆看似合理的垃圾（武將名字都在，據點編號卻超出範圍），
			// 而對拍會把那些垃圾算成「remake 畫錯了」。年份 0 只可能是空槽——
			// 四個劇本都從 190 年以後開始。
			if g.world.Clock.Year == 0 {
				log.Fatalf("-load-slot %d 是空的（年份 0）——那一槽沒有存檔", *loadSlot)
			}
			log.Printf("從第 %d 槽載入：%d年%d月%d日，玩家勢力 %d",
				*loadSlot+1, g.world.Clock.Year, g.world.Clock.Month,
				g.world.Clock.Day, g.world.Player)
			// readSave 會在事件列留一句「已讀取第 N 槽」。那是給玩家看的，
			// **對拍時它會蓋掉畫面下緣**（`map` 與 `faction` 兩區各差一成以上），
			// 而原版載入完什麼都不留。驗收路徑清掉它。
			g.lastEvent = ""
		}
		if *listCorps {
			logAliveCorps(g)
		}
		if *openEnding >= 0 {
			openEndingFixture(g, *openEnding)
		}
		g.lordCorps = *lordCorpsFlag
		g.damageReport = *damageReportFlag
		configureDirectFixtures(g, *openWin, *openList, *openAdvise, *adviseMenu, *adviseSortie, *adviseTarget, *openCities, *openFactions, *openCityInfo, *openForm, *openCorps, *openMarchList,
			*openMarchMode, *openBattle, *openSiege, *openMessage, *openFinance, *financeAmount, *openFormPick, *formPickRow,
			*openTalkIndex, *openOutcome, parseSiegeFixture(*siegeNode, *siegeDefend, *siegeCorps, *battleSteps),
			corpsMapFixture{enabled: *corpsOnMap, marchTo: *marchTo},
			*camAt, *battleCam)
		if *battleFF {
			g.toggleBattleFastForward()
		}
	} else {
		slots := inspectLauncherSlots(*saveFile)
		// 劇本標題從檔案讀，不硬編（docs/spec/25 §1.2）。
		g.scenarioTitles = map[int]string{}
		for i := 0; i < 4; i++ {
			if w, err := state.LoadScenario(loadPath, i); err == nil {
				g.scenarioTitles[i] = big5(w.Title)
			}
		}
		g.launcher = newLauncher(hasAvailableLauncherSlot(slots), slots)
		if *openNaming {
			// 驗收用：跳到劇本 0 的選君主頁並打開命名視窗。
			// setScenarioPlayers 只在「選劇本」那一頁收資料，先把狀態機擺過去。
			g.launcher.phase = launcherScenario
			if err := g.applyLauncherResult(launcherResult{kind: launcherPreviewScenario, scenario: 0}); err != nil {
				log.Printf("-open-naming：%v", err)
			}
			g.launcher.phase = launcherSelectPlayer
			if err := g.openNaming(); err != nil {
				log.Printf("-open-naming：%v", err)
			}
		}
	}

	ebiten.SetWindowSize(screenW*2, screenH*2)
	ebiten.SetWindowTitle("臥龍傳－三國制霸之計")
	if err := ebiten.RunGame(g); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}

func (g *game) buildRoads(w *state.World) *march.Graph {
	if g.lib == nil || g.lib.World == nil || w == nil {
		return nil
	}
	xy := make([][2]int, len(w.Cities))
	for i := range w.Cities {
		xy[i] = [2]int{w.Cities[i].X, w.Cities[i].Y}
	}
	edges, err := world.RoadEdges(g.lib.World, xy)
	if err != nil {
		log.Printf("⚠ 推不出道路圖（%v）；行軍會走直線", err)
		return nil
	}
	log.Printf("道路圖：%d 條路", len(edges))
	return march.New(len(w.Cities), world.MarchEdges(edges, xy))
}

// startWorld 是唯一的正式 World 建立入口。一般 launcher 在確認新局／
// 讀檔後才呼叫；direct fixture 也走同一條路，避免兩套初始化語意漂移。
func (g *game) startWorld(path string, slot int, player int, overridePlayer, newGame bool) error {
	w, err := state.LoadScenario(path, slot)
	if err != nil {
		return err
	}
	if overridePlayer {
		if !validLauncherPlayer(w, player) {
			return fmt.Errorf("玩家勢力 %d 在劇本 %d 不合法", player, slot+1)
		}
		w.Player = player
	} else if !validLauncherPlayer(w, w.Player) {
		return fmt.Errorf("第 %d 槽沒有合法玩家資料", slot+1)
	}
	w.EnableStrategicAI()
	if newGame {
		// 原版新遊戲在選定君主後立即跑一次政略評估（sub_11AC3+66 →
		// sub_12BD9）；讀檔路徑跳過——佇列隨存檔還原（docs/spec/83）。
		w.RunInitialStrategyPass(g.rng)
	}

	g.world = w
	g.roads = g.buildRoads(w)
	if g.roads != nil {
		w.SetRoads(g.roads)
	}
	g.tactical = nil
	g.battleLib = nil
	g.battleSprites = nil
	g.view = nil
	g.battleCommandBase = nil
	g.battleCommandGlyphs = [6]*ebiten.Image{}
	g.battleOrderIcons = [6]*ebiten.Image{}
	g.battleSideCommands = nil
	g.battleCommandSelect = color.RGBA{240, 0, 0, 255}
	g.battleSideFlags = [2]*ebiten.Image{}
	g.battleFormationStrip = nil
	g.battleSideFooter = nil
	g.battleFormation = -1
	g.battleUnitDotAlly = color.RGBA{110, 235, 110, 255}
	g.battleUnitDotFoe = color.RGBA{110, 200, 235, 255}
	g.battleGateBarColor = color.RGBA{140, 225, 225, 255}
	g.installTactical(g.origDir)
	g.saveBase = path
	g.hud = 0
	g.list = nil
	g.form = formState{}
	g.marchMode = marchModeState{}
	g.finance = financeState{}
	g.advise = adviseNone
	g.messages = nil
	g.lastEvent = ""
	g.quitting = false
	g.idleGate = idleClockGate{}

	season := int(w.Clock.Season())
	g.chrome = chrome.Load(g.lib, season)
	if !g.chrome.Available() {
		log.Printf("⚠ 取不到 ICONGRF 段 3 的視窗外框，改畫純色框")
	}
	if frame, err := g.lib.DOSVAmountPanel(season); err != nil {
		log.Printf("⚠ 取不到 DOS/V 數值視窗內框，指定金額畫面改用通用框：%v", err)
	} else {
		g.amountFrame = ebiten.NewImageFromImage(frame)
	}
	if cursor, err := g.lib.DOSVCursor(season); err != nil {
		log.Printf("⚠ 取不到 DOS/V 內建硬體游標，指定金額畫面不疊游標：%v", err)
	} else {
		g.cursorImage = ebiten.NewImageFromImage(cursor)
	}
	if base, err := g.lib.DOSVBattleCommandBase(season); err != nil {
		log.Printf("⚠ 取不到 DOS/V 戰術底列底板，改用文字 fallback：%v", err)
	} else {
		g.battleCommandBase = ebiten.NewImageFromImage(base)
	}
	for i := range g.battleCommandGlyphs {
		glyph, err := g.lib.DOSVBattleCommandGlyph(i, season)
		if err != nil {
			log.Printf("⚠ 取不到 DOS/V 戰術指令 glyph %d，改用文字 fallback：%v", i, err)
			g.battleCommandGlyphs = [6]*ebiten.Image{}
			break
		}
		g.battleCommandGlyphs[i] = ebiten.NewImageFromImage(glyph)
	}
	for i := range g.battleOrderIcons {
		icon, err := g.lib.DOSVOrderIcon(i, season)
		if err != nil {
			log.Printf("⚠ 取不到 DOS/V 戰術命令圖示 %d，底列只畫位置名：%v", i, err)
			g.battleOrderIcons = [6]*ebiten.Image{}
			break
		}
		g.battleOrderIcons[i] = ebiten.NewImageFromImage(icon)
	}
	for i := range g.battleArmIcons {
		icon, err := g.lib.DOSVSquadArmIcon(i, season)
		if err != nil {
			log.Printf("⚠ 取不到 DOS/V 戰術兵種圖示 %d，底列不畫兵種：%v", i, err)
			g.battleArmIcons = [3]*ebiten.Image{}
			break
		}
		g.battleArmIcons[i] = ebiten.NewImageFromImage(icon)
	}
	if panel, err := g.lib.DOSVBattleSideCommands(season); err != nil {
		log.Printf("⚠ 取不到 DOS/V 戰術右欄命令面板，改用文字 fallback：%v", err)
	} else {
		g.battleSideCommands = ebiten.NewImageFromImage(panel)
	}
	if selected, err := g.lib.PaletteColor(season, 0x0C); err != nil {
		log.Printf("⚠ 取不到 DOS/V 戰術指令選取色 0x0C，使用紅色 fallback：%v", err)
	} else {
		g.battleCommandSelect = selected
	}
	for i, foe := range [2]bool{false, true} {
		flag, err := g.lib.DOSVBattleFlag(foe, season)
		if err != nil {
			log.Printf("⚠ 取不到 DOS/V 戰術將旗 %d，那一格改畫通用框：%v", i, err)
			continue
		}
		g.battleSideFlags[i] = ebiten.NewImageFromImage(flag)
	}
	if strip, err := g.lib.DOSVBattleFormationStrip(season); err != nil {
		log.Printf("⚠ 取不到 DOS/V 戰術陣形列，那一格改畫通用框：%v", err)
	} else {
		g.battleFormationStrip = ebiten.NewImageFromImage(strip)
	}
	if footer, err := g.lib.DOSVBattleSideFooter(season); err != nil {
		log.Printf("⚠ 取不到 DOS/V 戰術側欄底列：%v", err)
	} else {
		g.battleSideFooter = ebiten.NewImageFromImage(footer)
	}
	// 取不到調色盤時的替代色（無素材的測試環境）。真正的值在下面那個迴圈裡覆蓋。
	g.battleTitlePlace = color.RGBA{150, 190, 235, 255}
	g.battleTitleLord = color.RGBA{140, 225, 225, 255}
	g.battleMenBar = color.RGBA{235, 120, 110, 255}
	g.battleHealthBar = color.RGBA{140, 225, 225, 255}
	for part := range g.battleFrame {
		img, err := g.lib.DOSVBattleFrame(library.BattleFramePart(part), season)
		if err != nil {
			log.Printf("⚠ 取不到 DOS/V 戰術側欄外框第 %d 塊：%v", part, err)
			g.battleFrame = [4]*ebiten.Image{}
			break
		}
		g.battleFrame[part] = ebiten.NewImageFromImage(img)
	}
	for _, c := range []struct {
		idx int
		dst *color.RGBA
	}{{10, &g.battleUnitDotAlly}, {3, &g.battleUnitDotFoe}, {11, &g.battleGateBarColor},
		{9, &g.battleTitlePlace}, {11, &g.battleTitleLord},
		{12, &g.battleMenBar}, {11, &g.battleHealthBar}} {
		col, err := g.lib.PaletteColor(season, c.idx)
		if err != nil {
			log.Printf("⚠ 取不到小地圖部隊點色 %d，使用 fallback：%v", c.idx, err)
			continue
		}
		*c.dst = col
	}
	// 縮小地圖標記的六個色號（docs/re/62 §2）。取不到就留 fallback。
	g.minimapInk = [6]color.RGBA{
		{0, 0, 0, 255}, {255, 255, 255, 255}, {221, 0, 0, 255},
		{255, 238, 0, 255}, {51, 68, 221, 255}, {0, 34, 102, 255}}
	for i, idx := range [6]int{0, 15, 10, 12, 3, 8} {
		col, err := g.lib.PaletteColor(season, idx)
		if err != nil {
			log.Printf("⚠ 取不到縮小地圖標記色 %d，使用 fallback：%v", idx, err)
			continue
		}
		g.minimapInk[i] = col
	}
	g.minimapFaction = 0
	if cap := w.Factions[w.Player].Capital; cap >= 0 && cap < len(w.Cities) {
		g.camX = w.Cities[cap].X - centreCol
		g.camY = w.Cities[cap].Y - centreRow
	}
	g.clampCam()
	g.checkCityCentres()
	log.Printf("劇本 %d：%d年%d月%d日，勢力 %d 個，玩家所仕 %d（君主 %s）",
		slot+1, w.Clock.Year, w.Clock.Month, w.Clock.Day,
		len(w.AliveFactions()), w.Player, text.Decode([]byte(w.LordName(w.Player)), text.Big5))
	return nil
}

// siegeFixture 是 `-siege-*`／`-battle-steps` 那一組驗收旗標。
//
// corps 兩格是「攻方軍團、守方軍團」的編號，**−1 ＝ 現編**
// （`demoBattle` 的舊行為）。指定既有軍團時兵種與人數照存檔，
// 側欄的名字與計量條才有機會與原版一致（docs/spec/90 §2.3）。
// corpsMapFixture 是「停在大地圖看軍團」的驗收 fixture（docs/spec/74 §4.1）。
//
// 既有的 -open-form／-open-corps／-open-march-list 都停在視窗裡，
// 而視窗蓋住大地圖——用它們截圖看不到軍團疊在哪一格。
type corpsMapFixture struct {
	enabled bool
	marchTo int // −1 ＝ 只編成，不下行軍指示
}

type siegeFixture struct {
	node   int
	defend bool
	corps  [2]int
	steps  int
}

// parseSiegeFixture 把 `-siege-corps 攻,守` 解成兩個編號。
// 格式不對就退回「現編」，並在 log 說明——**驗收旗標不要靜默失敗**，
// 否則對出來的差異會被算到別的成因頭上。
func parseSiegeFixture(node int, defend bool, corps string, steps int) siegeFixture {
	f := siegeFixture{node: node, defend: defend, corps: [2]int{-1, -1}, steps: steps}
	if corps == "" {
		return f
	}
	var att, def int
	if _, err := fmt.Sscanf(corps, "%d,%d", &att, &def); err != nil || att < 0 || def < 0 {
		log.Printf("⚠ -siege-corps 要 `攻,守` 兩個非負整數，收到 %q：改用現編的軍團", corps)
		return f
	}
	f.corps = [2]int{att, def}
	return f
}

// logAliveCorps 印出載入後還活著的軍團，給 `-siege-corps` 挑編號用。
// 軍團與武將同索引（`state.World.Leader`），所以編號就是主將的編號。
func logAliveCorps(g *game) {
	if g.world == nil {
		return
	}
	n := 0
	for i := range g.world.Corps {
		c := &g.world.Corps[i]
		if !c.Alive {
			continue
		}
		n++
		log.Printf("軍團 %3d：勢力 %2d　主將 %-8s　兵 %5d　士氣 %3d　據點 %3d　位置 (%3d,%3d)",
			i, c.Faction, big5(g.world.Generals[g.world.Leader(i)].Name),
			c.Men, c.Morale, c.Node, c.X, c.Y)
	}
	log.Printf("還活著的軍團共 %d 支", n)
}


func configureDirectFixtures(g *game, openWin int, openList, openAdvise, adviseMenu, adviseSortie, adviseTarget, openCities, openFactions bool, openCityInfo int, openForm, openCorps, openMarchList, openMarchMode,
	openBattle, openSiege, openMessage, openFinance bool, financeAmount int, openFormPick bool, formPickRow, openTalkIndex int,
	openOutcome string, siege siegeFixture, corpsMap corpsMapFixture, camAt, battleCam string) {
	w := g.world
	if w == nil {
		return
	}
	player := w.Player
	if openMessage && player >= 0 && player < len(w.Factions) {
		capital := w.Factions[player].Capital
		if capital >= 0 && capital < len(w.Cities) {
			g.enqueueTalkNotice(state.TalkNotice{Index: 0x46, City: capital, Faction: -1, General: -1, Amount: -1})
		}
	}
	if openTalkIndex >= 0 {
		vars := map[byte]string{'1': "武將", '2': "據點", '3': "君主", '4': "軍師", '5': "目標", '6': "", '7': "1234"}
		if len(w.Generals) > 0 {
			vars['1'] = big5(w.Generals[0].Name)
		}
		if player >= 0 && player < len(w.Factions) {
			vars['3'] = big5(w.LordName(player))
			advisor := w.Factions[player].Advisor
			if advisor >= 0 && advisor < len(w.Generals) {
				vars['4'] = big5(w.Generals[advisor].Name)
			}
			city := w.Factions[player].Capital
			if city >= 0 && city < len(w.Cities) {
				vars['2'] = big5(w.Cities[city].Name)
			}
		}
		if openTalkIndex < len(g.lib.Talk.Messages) {
			g.enqueueTalk(openTalkIndex, vars)
		} else {
			log.Printf("⚠ TALK 槽位超出範圍：%d", openTalkIndex)
		}
	}
	// 只供截圖的 fixture；正常 state path 仍由 AdjustTrust／capture 產生。
	switch openOutcome {
	case "trust":
		w.AdjustTrust(-w.Trust)
	case "faction":
		w.DebugLatchOutcomeForShot(state.DefeatFactionEliminated)
	}
	// -open-window N 開第 N 個視窗（0–3），驗收用。
	// −2 ＝ 四窗全開。原版實錄的主畫面就是這個狀態
	// （docs/playtest/27），沒有它就沒辦法對拍。
	if openWin == -2 {
		g.hudSet(hudCommand|hudFaction|hudMinimap, true)
	}
	// −3 連系統選單一起開。原版開著它時**時間停止**，所以這一組也是
	// 「時間凍住」那個對拍狀態（docs/spec/90 §2.1）。
	if openWin == -3 {
		g.hudSet(hudCommand|hudFaction|hudMinimap|hudSystem, true)
	}
	if camAt != "" {
		var cx, cy int
		if _, err := fmt.Sscanf(camAt, "%d,%d", &cx, &cy); err != nil {
			log.Printf("⚠ -cam 要 `X,Y` 兩個整數，收到 %q：%v", camAt, err)
		} else {
			g.camX, g.camY = cx, cy
			g.clampCam()
		}
	}
	if w := hudSwitchWindow(openWin); openWin >= 0 && w != 0 {
		g.hudSet(w, true)
	}
	if openList {
		// 開著、無選取——原版剛開窗沒有反白列（playtest/42 §4）。
		g.openGeneralList()
	}
	if openFormPick {
		// 編成的武將一覽（原版指令列 #3 剛開的狀態：候選已濾、無選取）。
		g.beginForm()
	}
	if openBattle || openSiege {
		g.demoBattle(openSiege, siege)
	}
	if battleCam != "" {
		var bx, by int
		if _, err := fmt.Sscanf(battleCam, "%d,%d", &bx, &by); err != nil {
			log.Printf("⚠ -battle-cam 要 `X,Y` 兩個整數，收到 %q：%v", battleCam, err)
		} else {
			// 戰術 view 是進戰場那一幀才建的，所以存起來由 newBattleView 套用。
			g.battleCamAt = &[2]int{bx, by}
			if g.view != nil {
				g.view.SetCameraWorld(bx, by)
			}
		}
	}
	if openForm || openCorps {
		g.demoCorps(openCorps, formPickRow)
	}
	if corpsMap.enabled {
		g.demoCorpsOnMap(corpsMap.marchTo)
	}
	// 財政視窗（對拍用，docs/spec/14 §4）。原版是命令列 #2 直接開視窗，
	// -finance-amount N 再開第 N 列的數值輸入器（docs/spec/78）。
	if openFinance || financeAmount >= 0 {
		g.beginFinance()
		if financeAmount >= 0 && financeAmount < financeRows {
			g.beginFinanceAmount(financeAmount)
			// 對拍截圖不畫游標：原版那一刻的游標在存還原區外
			// （p6-amount.png），headless 的指標位置又不可控。
			g.hideAmountCursor = true
		}
	}
	if openMarchMode {
		g.demoMarchMode()
	}
	if openMarchList {
		g.demoMarchList()
	}
	if adviseMenu && !openAdvise {
		g.openAdvise() // 停在五項選單
		return
	}
	if adviseTarget {
		// 進言 → 交戰（第 0 列）：目標勢力清單剛開的狀態（docs/spec/90 §5.1）。
		g.openAdvise()
		g.pickAdviseCommand(0)
		return
	}
	if openCities {
		g.openCityList()
		return
	}
	if openFactions {
		g.openFactionList()
		return
	}
	if openCityInfo > -2 {
		n := openCityInfo
		if n == -1 {
			n = g.world.Factions[g.world.Player].Capital
		}
		g.openCityInfo(n)
		return
	}
	if adviseSortie {
		g.openAdvise()
		g.beginSortie()
		return
	}
	if openAdvise {
		g.openAdvise()
		g.adviseCmd = persuasion.Hostility
		g.target = 13
		g.beginPersuasion()
		if g.sess != nil && !adviseMenu {
			g.offerReason(persuasion.WeAreStronger)
		}
		// 驗收要看的是講完之後的畫面，把逐句節拍跑完（docs/spec/45 §1.1）。
		for g.adviseAdvance() {
		}
	}
}
