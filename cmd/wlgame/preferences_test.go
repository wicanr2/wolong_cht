package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/ui/sound"
)

func TestDesktopPreferencesFileLifecycle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings", "preferences.json")
	defaults := desktopPreferences{Version: 1, Language: "zh-hant", Sound: 1, StrategySpeed: 2, TacticalSpeed: 2, LordCorps: true}
	if got, err := readDesktopPreferences(path, defaults); !os.IsNotExist(err) || got != defaults {
		t.Fatalf("缺檔應保留預設：%+v %v", got, err)
	}
	p := defaults
	p.BattleResultSeconds, p.Sound, p.Language = 15, 4, "en"
	for _, seconds := range []int{15, 3, 0} {
		p.BattleResultSeconds = seconds
		if err := writeDesktopPreferences(path, p); err != nil {
			t.Fatal(err)
		}
		if got, err := readDesktopPreferences(path, defaults); err != nil || got != p {
			t.Fatalf("重啟讀回錯誤：%+v %v", got, err)
		}
	}
	files, _ := os.ReadDir(filepath.Dir(path))
	if len(files) != 1 {
		t.Fatalf("留下暫存檔：%v", files)
	}
	for _, bad := range []string{`{`, `{"version":2}`, `{"version":1,"sound":8}`, `{"version":1,"battle_result_seconds":-1}`, `{"version":1,"language":"unknown"}`, strings.Repeat(" ", 65537)} {
		if err := os.WriteFile(path, []byte(bad), 0600); err != nil {
			t.Fatal(err)
		}
		if got, err := readDesktopPreferences(path, defaults); err == nil || got != defaults {
			t.Fatalf("無效檔案不得污染預設：%+v %v", got, err)
		}
		raw, _ := os.ReadFile(path)
		if string(raw) != bad {
			t.Fatal("讀取失敗不應覆寫原檔")
		}
	}
}

func TestDesktopPreferencesReplaceFailurePreservesTarget(t *testing.T) {
	path := filepath.Join(t.TempDir(), "preferences.json")
	if err := os.Mkdir(path, 0700); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(path, "keep")
	if err := os.WriteFile(marker, []byte("原內容"), 0600); err != nil {
		t.Fatal(err)
	}
	p := desktopPreferences{Version: 1, Language: "zh-hant"}
	if err := writeDesktopPreferences(path, p); err == nil {
		t.Fatal("不可替換目錄")
	}
	if raw, err := os.ReadFile(marker); err != nil || string(raw) != "原內容" {
		t.Fatalf("替換失敗損壞既有內容：%q %v", raw, err)
	}
	files, _ := os.ReadDir(filepath.Dir(path))
	if len(files) != 1 {
		t.Fatal("替換失敗留下暫存檔")
	}
}

func TestDesktopPreferencesExplicitFlagsAndPlayerChange(t *testing.T) {
	g := &game{speed: 1, tacticalSpeed: 4, lordCorps: true, sound: sound.Open("")}
	p := desktopPreferences{Version: 1, Language: "en", StrategySpeed: 3, TacticalSpeed: 2, Sound: 4, BattleResultSeconds: 10}
	if err := g.applyDesktopPreferences(p, map[string]bool{"lang": true, "speed": true, "lord-corps": true}); err != nil {
		t.Fatal(err)
	}
	if g.speed != 1 || !g.lordCorps || g.tacticalSpeed != 2 || g.soundOption() != 4 || g.battleResultSeconds != 10 {
		t.Fatalf("明示參數未優先或偏好未接上：%+v", g.desktopPreferences())
	}
	g.preferencesPath = filepath.Join(t.TempDir(), "preferences.json")
	g.dispatchSystemRow(sysRowBattleResult, true)
	p, err := readDesktopPreferences(g.preferencesPath, desktopPreferences{})
	if err != nil || p.BattleResultSeconds != 15 {
		t.Fatalf("玩家變更未保存：%+v %v", p, err)
	}
	g.adjustSpeed(false, 1)
	p, err = readDesktopPreferences(g.preferencesPath, desktopPreferences{})
	if err != nil || p.StrategySpeed != 0 {
		t.Fatalf("快捷鍵未保存：%+v %v", p, err)
	}
}
