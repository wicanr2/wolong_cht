//go:build android

package wolongmobile

import "github.com/hajimehoshi/ebiten/v2/mobile"

// defaultDirs 在 Android 上刻意回空：資料根目錄由 Java 側算出來，
// 經 SetDataRoot 交進來（`getExternalFilesDir`／`getFilesDir` 都不是常數）。
//
// ⚠ Android 11 以上 `os.ReadFile` 拿不到 SAF 的 `content://`，
// 所以資料一定要先被複製到 app 的私有目錄（docs/mobile/android-plan.md §3）。
func defaultDirs() (orig, font string) { return "", "" }

// current 是這個行程的那一局。Android 的生命週期回呼（返回鍵）
// 要摸到它，而 `mobile.SetGame` 之後 Ebiten 不再交出 Game。
var current *game

func init() {
	opt, font := optionsFromEnv()
	current = newGame(opt, font, "", 0)
	mobile.SetGame(current)
}

// Back 由 Java 的 onBackPressed 呼叫。回傳 false 表示沒東西可關，
// Java 端才把返回鍵交回系統。
//
// ⚠ **不要讓返回鍵直接結束遊戲**：即時制沒有暫停，
// 誤觸一次就是進度沒了（docs/mobile/android-ux.md §4）。
func Back() bool {
	if current == nil {
		return false
	}
	return current.back()
}

// Initialize 給產生出來的 Java wrapper 一個穩定的入口。
func Initialize() {}
