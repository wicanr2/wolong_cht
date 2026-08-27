package phone

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/ui/sound"
)

// 手機端的場景查詢要餵得出正確的曲子（docs/spec/92 §2.2）。
//
// ⭐ 規則本身在 `internal/rules/bgm`，那邊有逐條的表；這一支只驗
// **Session 的狀態有沒有正確翻進 Scene**——翻錯的話兩端會放不同的曲子，
// 而那是「一條規則兩份實作」想避免的正是這件事。
func TestMusicSceneFollowsSessionState(t *testing.T) {
	s := newTestSession(t)

	scene := s.MusicScene()
	if scene.Month != s.world.Clock.Month {
		t.Errorf("月份沒帶進去：%d vs %d", scene.Month, s.world.Clock.Month)
	}
	if scene.Battle != nil {
		t.Error("開局不該在戰術畫面")
	}
	if got := s.MusicTrack(); got == "" {
		t.Error("開局應該有一首季節曲")
	}

	// 面板打開＝事件與對話那一組。
	s.sheet.open = true
	if got := s.MusicScene(); !got.Message {
		t.Error("面板開著時 Message 應該為真")
	}
	if got := s.MusicTrack(); got != "bgm-6" {
		t.Errorf("面板開著時的曲子 ＝ %q，預期 bgm-6", got)
	}
	s.sheet.open = false
}

// 「未接入」與「關」是兩件事：關是玩家的選擇，未接入是這一台沒有音檔。
func TestSoundRowSeparatesMissingFromOff(t *testing.T) {
	s := newTestSession(t)
	if got := s.soundValue(); got != "未接入" {
		t.Errorf("沒有音庫時 ＝ %q，預期「未接入」", got)
	}
	s.ToggleSound() // 沒有音庫時不該做事
	if got := s.soundValue(); got != "未接入" {
		t.Errorf("點過之後 ＝ %q，未接入不該變成「關」", got)
	}

	// 空目錄的 Bank：Available() 為假，仍然是「未接入」。
	s.AttachSound(sound.Open(t.TempDir()))
	if got := s.soundValue(); got != "未接入" {
		t.Errorf("空音庫 ＝ %q，預期「未接入」", got)
	}
}
