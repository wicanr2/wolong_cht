package main

// 桌面偏好與遊戲存檔分離，契約見 docs/spec/156。
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/wicanr2/wolong_cht/internal/ui/uitext"
)

type desktopPreferences struct {
	Version             int    `json:"version"`
	Language            string `json:"language"`
	VideoLCD            bool   `json:"video_lcd"`
	Sound               int    `json:"sound"`
	StrategySpeed       int    `json:"strategy_speed"`
	TacticalSpeed       int    `json:"tactical_speed"`
	LordCorps           bool   `json:"lord_corps"`
	DamageReport        bool   `json:"damage_report"`
	BattleResultSeconds int    `json:"battle_result_seconds"`
}

func (g *game) desktopPreferences() desktopPreferences {
	return desktopPreferences{1, string(uiLang.Lang()), g.videoLCD, g.soundOption(),
		g.speed, g.tacticalSpeed, g.lordCorps, g.damageReport, g.battleResultSeconds}
}

func (p desktopPreferences) validate() error {
	if p.Version != 1 {
		return fmt.Errorf("不支援的偏好版本 %d", p.Version)
	}
	switch p.Language {
	case "zh-hant", "zh-hans", "ja", "en":
	default:
		return fmt.Errorf("無效語言 %q", p.Language)
	}
	if p.Sound < 0 || p.Sound > 4 || p.StrategySpeed < 0 || p.StrategySpeed > 4 || p.TacticalSpeed < 0 || p.TacticalSpeed > 4 {
		return fmt.Errorf("音效或速度檔位超出 0–4")
	}
	for _, seconds := range battleResultDurations {
		if p.BattleResultSeconds == seconds {
			return nil
		}
	}
	return fmt.Errorf("無效的戰後結果秒數 %d", p.BattleResultSeconds)
}

func readDesktopPreferences(path string, defaults desktopPreferences) (desktopPreferences, error) {
	f, err := os.Open(path)
	if err != nil {
		return defaults, err
	}
	defer f.Close()
	const maxBytes = 64 * 1024
	raw, err := io.ReadAll(io.LimitReader(f, maxBytes+1))
	if err != nil {
		return defaults, err
	}
	if len(raw) > maxBytes {
		return defaults, fmt.Errorf("偏好檔案超過 64 KiB")
	}
	p := defaults
	p.Version = 0 // 必須明示版本；缺少其他欄位沿用預設。
	if err := json.Unmarshal(raw, &p); err != nil {
		return defaults, err
	}
	if err := p.validate(); err != nil {
		return defaults, err
	}
	return p, nil
}

func writeDesktopPreferences(path string, p desktopPreferences) error {
	if err := p.validate(); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".preferences-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(append(raw, '\n')); err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(f.Name(), path)
}

func (g *game) applyDesktopPreferences(p desktopPreferences, explicit map[string]bool) error {
	if !explicit["lang"] {
		if err := g.setLanguage(uitext.Language(p.Language)); err != nil {
			return err
		}
	}
	if !explicit["video-lcd"] {
		g.videoLCD = p.VideoLCD
	}
	if !explicit["speed"] {
		g.speed = p.StrategySpeed
	}
	if !explicit["tactical-speed"] {
		g.tacticalSpeed = p.TacticalSpeed
	}
	if !explicit["lord-corps"] {
		g.lordCorps = p.LordCorps
	}
	if !explicit["siege-damage"] {
		g.damageReport = p.DamageReport
	}
	g.setSoundOption(p.Sound)
	g.battleResultSeconds = p.BattleResultSeconds
	return nil
}

func (g *game) initDesktopPreferences() {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("設定記憶未啟用：%v", err)
		return
	}
	path := filepath.Join(dir, "wolong-remake", "preferences.json")
	p, err := readDesktopPreferences(path, g.desktopPreferences())
	if err == nil {
		explicit := map[string]bool{}
		flagVisit(func(name string) { explicit[name] = true })
		err = g.applyDesktopPreferences(p, explicit)
	}
	if err != nil && !os.IsNotExist(err) {
		log.Printf("未載入桌面偏好，保留原檔並沿用預設：%v", err)
	}
	// 啟動與套用偏好都不寫檔；只有之後的玩家操作才可保存。
	g.preferencesPath = path
}

func (g *game) saveDesktopPreferences() {
	if g.preferencesPath == "" {
		return
	}
	if err := writeDesktopPreferences(g.preferencesPath, g.desktopPreferences()); err != nil {
		log.Printf("偏好保存失敗：%v", err)
		g.setEvent("設定本次有效，但無法保存到下次啟動")
	}
}
