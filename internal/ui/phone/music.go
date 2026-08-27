package phone

import (
	"github.com/wicanr2/wolong_cht/internal/rules/battlefield"
	"github.com/wicanr2/wolong_cht/internal/rules/bgm"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/sound"
)

// MusicScene 把目前的狀態翻成選曲條件（docs/spec/92 §2.2）。
//
// ⭐ **規則不在這裡**：曲目與場景的對應在 `internal/rules/bgm`，
// 桌面與手機共用同一份。這一支只負責把 Session 的狀態填進去——
// 抄第二份規則會長出兩邊不一樣的行為（CLAUDE.md §7 第 6 條）。
func (s *Session) MusicScene() bgm.Scene {
	if s == nil || s.world == nil {
		return bgm.Scene{}
	}
	scene := bgm.Scene{
		GameOver: s.world.Outcome() != state.InProgress,
		Message:  len(s.notices) > 0 || s.sheet.open,
		Month:    s.world.Clock.Month,
	}
	if p := s.world.PendingBattle(); p != nil {
		b := bgm.Battle{Field: battlefield.FieldBase}
		if s.setup != nil {
			b.Field = s.setup.FieldNumber(p.Node, p.Mode == combat.Siege)
		}
		b.PlayerAttacks = p.Attacker >= 0 && p.Attacker < len(s.world.Corps) &&
			s.world.Corps[p.Attacker].Faction == s.world.Player
		scene.Battle = &b
	}
	return scene
}

// MusicTrack 是「現在該放哪一首」，空字串表示不換曲。
func (s *Session) MusicTrack() string { return bgm.Track(s.MusicScene()) }

// AttachSound 讓呈現層把音庫交進來，系統面板才畫得出「音效」那一列。
//
// ⚠ **Session 不播音樂**，播的是呈現層（`mobile/wolong`）——
// 規則層與版面層都不該認識 Ebiten 的 audio。這裡只握著它好回答
// 「有沒有音檔」與「開著沒有」。
func (s *Session) AttachSound(b *sound.Bank) {
	if s == nil {
		return
	}
	s.music = b
}

// ToggleSound 翻轉音效開關。沒有音檔時不做事——那一列顯示「未接入」，
// 點它不該把狀態變成「關」（那會讓缺口從畫面上消失）。
func (s *Session) ToggleSound() {
	if s == nil || s.music == nil || !s.music.Available() {
		return
	}
	s.music.SetEnabled(!s.music.Enabled())
}
