package main

import (
	"os"
	"path/filepath"
	"testing"
)

// ⭐ 這一支擋的是「解開的包裡預設路徑不成立」那個安靜的失敗，
// 同時擋住它的三個副作用（docs/spec/72 §3）。
func TestResolveDataDir(t *testing.T) {
	// exe 旁邊的三種版面各造一份，測試才知道自己在驗哪一條。
	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("取不到執行檔路徑：%v", err)
	}
	exeDir := filepath.Dir(exe)

	const def = "workplace/orig/dosv"
	const bundled = "gamedata"

	t.Run("明講的旗標一律尊重", func(t *testing.T) {
		// ⭐ 即使那個目錄不存在也要照用——對拍與驗收明講路徑，
		// 我們不可以「幫他找一個比較好的」。
		if got := resolveDataDir("/tmp/沒有這個目錄", def, bundled); got != "/tmp/沒有這個目錄" {
			t.Errorf("明講的路徑被換掉了：%s", got)
		}
	})

	t.Run("repo 內行為不變", func(t *testing.T) {
		// 造一個當下工作目錄底下的 def，模擬 repo 內開發。
		wd := t.TempDir()
		old, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		if err := os.Chdir(wd); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { os.Chdir(old) })
		if err := os.MkdirAll(filepath.Join(wd, def), 0o755); err != nil {
			t.Fatal(err)
		}
		if got := resolveDataDir(def, def, bundled); got != def {
			t.Errorf("repo 路徑存在時應該原樣回傳，得到 %s", got)
		}
	})

	t.Run("exe 旁邊的完整包", func(t *testing.T) {
		// 工作目錄換到一個沒有 def 的地方，第二條才不會先命中。
		wd := t.TempDir()
		old, _ := os.Getwd()
		if err := os.Chdir(wd); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { os.Chdir(old) })

		want := filepath.Join(exeDir, bundled)
		if err := os.MkdirAll(want, 0o755); err != nil {
			t.Skipf("造不出 %s（唯讀的測試環境）：%v", want, err)
		}
		t.Cleanup(func() { os.RemoveAll(want) })
		if got := resolveDataDir(def, def, bundled); got != want {
			t.Errorf("應該找到 exe 旁邊的 %s，得到 %s", want, got)
		}
	})

	t.Run("都找不到就回預設值", func(t *testing.T) {
		wd := t.TempDir()
		old, _ := os.Getwd()
		if err := os.Chdir(wd); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { os.Chdir(old) })
		// ⚠ 回預設值而不是空字串：載入器要噴一個帶路徑的錯，
		// 空字串會變成「讀當下目錄」，症狀完全不同。
		if got := resolveDataDir(def, def, "這個名字不會存在"); got != def {
			t.Errorf("都找不到時應該回預設值，得到 %s", got)
		}
	})
}

// 兩個包內目錄名不能互撞，也不能是空字串——空字串會讓
// filepath.Join(dir, "") 退化成 exe 目錄本身。
func TestBundledDirNames(t *testing.T) {
	if bundledOrigDir == "" || bundledFontDir == "" {
		t.Fatal("包內目錄名不可為空")
	}
	if bundledOrigDir == bundledFontDir {
		t.Fatal("資料與字型不能放同一個目錄：ImportActivity 也是分開的")
	}
}
