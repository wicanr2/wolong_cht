// Package savepath 集中管理玩家存檔 overlay 的檔案系統邊界。
// 刻意不依賴 Ebiten，讓路徑安全可以在不啟動圖形 runtime 的情況下測試。
package savepath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// InitialLoadPath 優先選取既有 overlay，否則退回唯讀來源。
// 缺少的 overlay 要等第一次明確儲存時才建立。
func InitialLoadPath(source, overlay string) (string, error) {
	if overlay == "" {
		return source, nil
	}
	if SamePath(source, overlay) {
		return "", fmt.Errorf("存檔 overlay 不可與原始素材相同：%s", overlay)
	}
	info, err := os.Stat(overlay)
	switch {
	case err == nil:
		if info.IsDir() {
			return "", fmt.Errorf("存檔 overlay 是目錄：%s", overlay)
		}
		return overlay, nil
	case errors.Is(err, os.ErrNotExist):
		return source, nil
	default:
		return "", fmt.Errorf("檢查存檔 overlay 失敗：%w", err)
	}
}

// SamePath 同時攔截路徑文字別名，以及兩個路徑都存在時的硬連結／符號連結別名。
// 它是防止就地寫入原版資料的失敗即關閉（fail-closed）護欄。
func SamePath(a, b string) bool {
	aa, errA := filepath.Abs(a)
	bb, errB := filepath.Abs(b)
	if errA == nil && errB == nil && filepath.Clean(aa) == filepath.Clean(bb) {
		return true
	}
	ai, errA := os.Stat(a)
	bi, errB := os.Stat(b)
	return errA == nil && errB == nil && os.SameFile(ai, bi)
}
