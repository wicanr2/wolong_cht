//go:build !android

package wolongmobile

import (
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

// RunDesktop 在桌面跑手機版。
//
// ⚠ 這**不是**給玩家的桌面版（那是 `cmd/wlgame`，也是唯一的對拍基準）。
// 它是手機 UI 的開發與驗收入口：同一份 `internal/ui/phone`，
// 用桌面的快迴圈驗，最後才進 APK。
//
// 環境變數：`WOLONG_ORIG`／`WOLONG_FONT`／`WOLONG_SCENARIO`／
// `WOLONG_PLAYER`／`WOLONG_SEED`／`WOLONG_SHOT`／`WOLONG_SHOT_FRAME`。
func RunDesktop() error {
	opt, font := optionsFromEnv()
	shot := os.Getenv("WOLONG_SHOT")
	at, err := strconv.Atoi(os.Getenv("WOLONG_SHOT_FRAME"))
	if err != nil || at <= 0 {
		at = 60
	}
	ebiten.SetWindowSize(960, 540)
	ebiten.SetWindowTitle("臥龍傳 Remake — 手機版")
	g := newGame(opt, font, shot, at)
	// 推廣片的逐幀輸出（`WOLONG_FRAMES_DIR`），手機上不存在。
	g.rec = newDemoRecorder()
	return ebiten.RunGame(g)
}

// defaultDirs 是桌面驗收的預設路徑（repo 相對）。
func defaultDirs() (orig, font string) { return "workplace/orig/dosv", "workplace/eten" }
