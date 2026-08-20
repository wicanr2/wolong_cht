// Package wolongmobile 是手機版的平台殼。
//
// 它只做三件事：把觸控與生命週期轉成輸入、把畫布交給 internal/ui/phone、
// 在 Android 側註冊 Ebiten 的 mobile 入口。**規則層不會知道自己跑在手機上。**
package wolongmobile

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/wicanr2/wolong_cht/internal/assets/cjk"
	"github.com/wicanr2/wolong_cht/internal/ui/phone"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

type game struct {
	sess *phone.Session
	td   *textdraw.Drawer
	err  string

	// opt／fontDir 是**還沒開局**時記著的設定。手機上這兩個一開始是空的，
	// 要等 Java 側把資料根目錄交過來（見 SetDataRoot）。
	opt     phone.Options
	fontDir string
	tried   bool

	// drag 是單指拖曳的狀態。按下的那一刻先記位置，
	// **移動超過門檻才算拖曳**——否則每一次點擊都會順便平移一格。
	dragging   bool
	dragMoved  bool
	dragX      int
	dragY      int
	pinchStart float64

	frame int
	// shotAt／shotPath 是驗收路徑：跑到第 N 幀存一張圖就離開。
	shotAt   int
	shotPath string
}

// dragThreshold 是「這是拖曳不是點擊」的門檻（螢幕像素）。
const dragThreshold = 12

func newGame(opt phone.Options, fontDir string, shotPath string, shotAt int) *game {
	g := &game{opt: opt, fontDir: fontDir, shotPath: shotPath, shotAt: shotAt}
	g.ensure()
	return g
}

// ensure 在資料就位之後才開局。
//
// ⭐ **不能在 init() 就開局**：Android 的資料路徑要由 Java 側算出來
// （`getExternalFilesDir`／`getFilesDir` 都不是常數），而 `mobile.SetGame`
// 又必須在 init() 呼叫。所以開局延後到第一次 Update——
// 這也順便讓「還沒匯入資料」變成一個**畫得出來的狀態**，不是崩潰。
func (g *game) ensure() {
	if g.sess != nil || g.tried {
		return
	}
	opt := g.opt
	fontDir := g.fontDir
	if opt.OrigDir == "" {
		root := DataRoot()
		if root == "" {
			return // 還在等 Java 交路徑，這一幀先畫等待畫面。
		}
		opt.OrigDir = filepath.Join(root, "orig")
		if fontDir == "" {
			fontDir = filepath.Join(root, "eten")
		}
	}
	g.tried = true
	// ⭐ 字型要在開局**之前**載入。載入失敗的訊息本身是中文，
	// 字型跟著開局一起載的話，開局失敗時螢幕上只剩一片底色——
	// 最需要訊息的那一刻反而什麼都印不出來。
	if f, err := cjk.LoadDir(fontDir, cjk.Options{}); err == nil {
		a, _ := cjk.LoadASCIIDir(fontDir)
		g.td = textdraw.New(f, a)
	} else {
		log.Printf("WOLONG_INIT 字型 %q 載不起來：%v", fontDir, err)
	}
	log.Printf("WOLONG_INIT root=%q orig=%q", DataRoot(), opt.OrigDir)
	sess, err := phone.NewSession(opt)
	if err != nil {
		// ⚠ 缺檔要**指名**：Android 端的資料是使用者自己匯入的，
		// 「載入失敗」四個字對他毫無用處。
		g.err = err.Error()
		log.Printf("WOLONG_INIT 開局失敗：%v", err)
		return
	}
	g.sess = sess
	// 驗收鉤子：直接把畫面推到要看的狀態。手機上這幾個環境變數都是空的。
	if v, err := strconv.Atoi(os.Getenv("WOLONG_ZOOM")); err == nil && v > 0 {
		sess.SetZoom(v)
	}
	if v, err := strconv.Atoi(os.Getenv("WOLONG_SELECT")); err == nil && v >= 0 {
		sess.Select(v)
	}
	if os.Getenv("WOLONG_PAUSED") != "" {
		sess.SetPaused(true)
	}
	if v, err := strconv.Atoi(os.Getenv("WOLONG_SHEET")); err == nil && v >= 0 {
		sess.OpenSheet(phone.Command(v))
		if t, err := strconv.Atoi(os.Getenv("WOLONG_TAB")); err == nil {
			sess.SetSheetTab(t)
		}
	}
	if v, err := strconv.Atoi(os.Getenv("WOLONG_ADVISE")); err == nil && v >= 0 {
		sess.PickAdvise(v)
	}
	if v, err := strconv.Atoi(os.Getenv("WOLONG_SIEGE")); err == nil && v >= 0 {
		if err := sess.OpenDemoBattle(v); err != nil {
			log.Printf("WOLONG_SIEGE：%v", err)
		}
	}
}

func (g *game) Update() error {
	g.frame++
	g.ensure()
	if g.sess == nil {
		return nil
	}
	g.handleInput()
	g.sess.Tick()
	g.logFingerprint()
	return nil
}

// fingerprintFrames 是要印指紋的幀。
//
// ⭐ 這是 **Android 里程碑 A 的驗收訊號**（docs/mobile/android-plan.md §5）：
// 同一個 seed、同樣的幀數，手機與桌面要算出同一個指紋（docs/spec/69）。
// 幀數是規則層的 tick 數，與機器快慢無關，所以兩邊可比。
// ⚠ 只有在**完全沒有輸入**的情況下可比——smoke 期間不要點畫面。
var fingerprintFrames = parseFingerprintFrames(os.Getenv("WOLONG_FP_FRAMES"))

// SetFingerprintFrames 由 Java 側轉交 Intent 的 `fp_frames`。
//
// ⚠ **Android 的 app 不繼承 adb 的環境變數**，所以桌面那條 `WOLONG_*`
// 的路在手機上一律無效；要從外面帶參數只能走 Intent 或檔案。
func SetFingerprintFrames(s string) {
	if s == "" {
		return
	}
	fingerprintFrames = parseFingerprintFrames(s)
}

// parseFingerprintFrames 讀 `WOLONG_FP_FRAMES`（逗號分隔的幀號）。
//
// ⚠ 模擬器在忙碌的機器上只跑到個位數 fps，而 Ebiten 在 GL context
// 遺失時會結束整個 app——跑到第 600 幀要好幾分鐘，中途被砍的機率不低。
// 幀號可設定是為了讓驗收能挑一組跑得完的，**不是為了放寬判準**：
// 同一組幀號兩邊仍然要一模一樣。
func parseFingerprintFrames(s string) map[int]bool {
	out := map[int]bool{}
	if s == "" {
		for _, n := range []int{1, 60, 600} {
			out[n] = true
		}
		return out
	}
	for _, f := range strings.Split(s, ",") {
		if n, err := strconv.Atoi(strings.TrimSpace(f)); err == nil && n > 0 {
			out[n] = true
		}
	}
	return out
}

func (g *game) logFingerprint() {
	if !fingerprintFrames[g.frame] {
		return
	}
	// Android 上 Go 的 log 會進 logcat（tag GoLog），桌面上進 stderr——
	// 兩邊都抓得到，不必為了驗收另外開通道。
	log.Printf("WOLONG_FP frame=%d fp=%s", g.frame, g.sess.World().FingerprintHex())
}

func (g *game) handleInput() {
	// 桌面上用 Escape 驗返回鍵的行為。
	// ⚠ **Android 的實體返回鍵不走這裡**：Ebiten 沒有對應的鍵碼，
	// 由 Java 的 onBackPressed 呼叫 Back()（android.go）。
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.back()
		return
	}
	// 觸控與滑鼠走同一條路——桌面上才驗得動（docs/mobile/android-plan.md §4）。
	ids := ebiten.AppendTouchIDs(nil)
	switch {
	case len(ids) >= 2:
		g.handlePinch(ids)
		return
	case len(ids) == 1:
		x, y := ebiten.TouchPosition(ids[0])
		g.handlePointer(x, y, inpututil.IsTouchJustReleased(ids[0]),
			len(inpututil.AppendJustPressedTouchIDs(nil)) > 0)
		return
	}
	g.pinchStart = 0
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) ||
		inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.handlePointer(x, y, inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft),
			inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft))
		return
	}
	g.dragging = false
}

// back 是返回鍵的行為：關進言 → 關面板 → 收小卡。
// 回傳 false 表示沒東西可關，呼叫端自己決定要不要離開。
func (g *game) back() bool {
	if g.sess == nil {
		return false
	}
	if g.sess.AdviseStage() != 0 {
		g.sess.CloseAdvise()
		return true
	}
	return g.sess.Back()
}

func (g *game) handlePointer(x, y int, released, pressed bool) {
	if pressed {
		g.dragging, g.dragMoved, g.dragX, g.dragY = true, false, x, y
		return
	}
	if released {
		if g.dragging && !g.dragMoved {
			g.sess.Tap(float64(x), float64(y))
		}
		g.dragging = false
		return
	}
	if !g.dragging {
		return
	}
	dx, dy := x-g.dragX, y-g.dragY
	if !g.dragMoved && abs(dx) < dragThreshold && abs(dy) < dragThreshold {
		return
	}
	g.dragMoved = true
	// 面板開著的時候拖曳是**捲列表**，不是平移地圖——
	// 底下的地圖看不見，平移它等於把操作丟進看不到的地方。
	if g.sess.SheetOpen() || g.sess.AdviseStage() != 0 {
		if sy := dy / rowScrollPx; sy != 0 {
			g.sess.ScrollRows(-sy)
			g.dragY = y
		}
		return
	}
	// 一格一格地捲：邏輯座標是整數格，半格的位移沒有意義。
	px := phone.TilePx
	if sx := dx / px; sx != 0 {
		g.sess.Pan(-sx, 0)
		g.dragX = x
	}
	if sy := dy / px; sy != 0 {
		g.sess.Pan(0, -sy)
		g.dragY = y
	}
}

// rowScrollPx 是列表捲一列要拖多少邏輯像素。取一列的高度，
// 手指移動的距離與列表跑的距離才會一致。
const rowScrollPx = 30

func (g *game) handlePinch(ids []ebiten.TouchID) {
	x0, y0 := ebiten.TouchPosition(ids[0])
	x1, y1 := ebiten.TouchPosition(ids[1])
	d := dist(x0, y0, x1, y1)
	if g.pinchStart == 0 {
		g.pinchStart = d
		return
	}
	_, _, z := g.sess.Camera()
	switch {
	case d > g.pinchStart*1.4:
		g.sess.SetZoom(z + 1)
		g.pinchStart = d
	case d < g.pinchStart/1.4:
		g.sess.SetZoom(z - 1)
		g.pinchStart = d
	}
	g.dragging = false
}

func (g *game) Draw(screen *ebiten.Image) {
	if g.sess == nil {
		g.drawNotReady(screen)
		g.maybeShot(screen)
		return
	}
	g.sess.Draw(screen, g.td)
	g.maybeShot(screen)
}

// drawNotReady 畫「還沒開局」。⚠ 這裡**不能只印英文**——訊息要指名缺什麼，
// 而字型本身也可能還沒到，所以文字缺席時仍然要有可辨識的底色。
func (g *game) drawNotReady(screen *ebiten.Image) {
	msg := "等待匯入原版資料"
	if g.err != "" {
		msg = g.err
	}
	screen.Fill(color.RGBA{40, 12, 12, 255})
	if g.td != nil && g.td.Available() {
		g.td.Draw(screen, "臥龍傳 Remake", 24, 24, color.RGBA{255, 200, 200, 255})
		g.td.Draw(screen, msg, 24, 56, color.RGBA{255, 200, 200, 255})
	}
}

func (g *game) Layout(_, _ int) (int, int) { return phone.LogicalW, phone.LogicalH }

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func dist(x0, y0, x1, y1 int) float64 {
	return math.Hypot(float64(x1-x0), float64(y1-y0))
}

// maybeShot 是驗收路徑：跑到第 N 幀存一張圖就離開。
//
// ⚠ 手機上不會走到這裡（`shotPath` 是空的）。它存在的理由是
// **讓手機版的畫面能用桌面那條 30 秒的迴圈驗**，
// 而不是每次都起模擬器（docs/mobile/android-plan.md §6）。
func (g *game) maybeShot(screen *ebiten.Image) {
	if g.shotPath == "" || g.frame < g.shotAt {
		return
	}
	f, err := os.Create(g.shotPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	screen.ReadPixels(img.Pix)
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
	os.Exit(0)
}

// optionsFromEnv 讓桌面驗收不必改程式就能換劇本、勢力與種子。
func optionsFromEnv() (phone.Options, string) {
	atoi := func(k string, def int) int {
		if v, err := strconv.Atoi(os.Getenv(k)); err == nil {
			return v
		}
		return def
	}
	// ⚠ 預設值**與平台有關**：桌面上是 repo 內的路徑，Android 上是空字串
	// ——手機的資料路徑要等 Java 側交過來（SetDataRoot），
	// 在這裡填一個 repo 相對路徑會讓手機端去讀一個永遠不存在的目錄。
	defOrig, defFont := defaultDirs()
	dir := os.Getenv("WOLONG_ORIG")
	if dir == "" {
		dir = defOrig
	}
	font := os.Getenv("WOLONG_FONT")
	if font == "" {
		font = defFont
	}
	return phone.Options{
		OrigDir:  dir,
		Scenario: atoi("WOLONG_SCENARIO", 0),
		Player:   atoi("WOLONG_PLAYER", 0),
		Seed:     atoi("WOLONG_SEED", 7),
	}, font
}

// dataRoot 是原版資料與字型的根目錄：`<root>/orig` 放 69 個原版檔，
// `<root>/eten` 放點陣字。
//
// ⚠ **兩者都由使用者自備**，不進 APK（deny-list 的邊界對手機版一樣成立）。
var (
	dataMu   sync.RWMutex
	dataRoot string
)

// SetDataRoot 由 Java 側在 `setContentView` 之前呼叫。
// gomobile 會把它綁成 `Wolongmobile.setDataRoot(String)`。
func SetDataRoot(p string) {
	dataMu.Lock()
	dataRoot = p
	dataMu.Unlock()
}

// DataRoot 回傳目前的資料根目錄，還沒設定時是空字串。
func DataRoot() string {
	dataMu.RLock()
	defer dataMu.RUnlock()
	return dataRoot
}
