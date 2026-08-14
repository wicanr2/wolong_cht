// Package savepath 集中管理玩家存檔 overlay 的檔案系統邊界。
// 刻意不依賴 Ebiten，讓路徑安全可以在不啟動圖形 runtime 的情況下測試。
package savepath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// NativeExt 是 remake 原生存檔的副檔名（`docs/spec/20`）。
const NativeExt = ".wlsave"

// NativePath 由原版格式的 overlay 路徑推出某一槽的原生存檔路徑：
//
//	SAVE.DAT + slot 0  →  SAVE-slot1.wlsave
//
// **兩者刻意分開檔案**：原版格式那一份要能直接拿去 DOSBox，
// 而原生檔帶著 remake 專屬狀態。共用一個檔名只會讓其中一種被覆蓋。
func NativePath(overlay string, slot int) (string, error) {
	if overlay == "" {
		return "", errors.New("沒有 overlay 路徑，無法推出原生存檔路徑")
	}
	if slot < 0 {
		return "", fmt.Errorf("槽位 %d 不合法", slot)
	}
	dir, base := filepath.Split(overlay)
	stem := strings.TrimSuffix(base, filepath.Ext(base))
	return filepath.Join(dir, fmt.Sprintf("%s-slot%d%s", stem, slot+1, NativeExt)), nil
}
