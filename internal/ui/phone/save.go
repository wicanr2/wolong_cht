package phone

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wicanr2/wolong_cht/internal/savefile"
	"github.com/wicanr2/wolong_cht/internal/savepath"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// SaveSlots 是存檔槽數。原版的 `SAVE.DAT` 就是四個 22,208 B 區塊。
const SaveSlots = 4

// 存檔的檔案位置。
//
// ⚠ **原版資產唯讀**（CLAUDE.md §9）：讀的是使用者匯入的
// `orig/SAVE.DAT`，寫的一律在 `save/` 底下。手機上這兩個目錄都在
// app 的私有空間裡，但這條界線與桌面版一樣不能鬆。
func (s *Session) sourceSave() string { return filepath.Join(s.origDir, "SAVE.DAT") }
func (s *Session) overlaySave() string {
	return filepath.Join(filepath.Dir(s.origDir), "save", "SAVE.DAT")
}

// slotLabel 說這一槽現在有什麼。
func (s *Session) slotLabel(slot int) string {
	path, err := savepath.NativePath(s.overlaySave(), slot)
	if err != nil {
		return "－"
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "空"
	}
	if err != nil {
		return "讀不到"
	}
	w, _, err := savefile.Decode(data)
	if err != nil {
		// ⚠ 壞掉的檔要**講出來**。顯示成「空」的話，玩家會直接存過去，
		// 而那一槽原本可能是他唯一的進度。
		return "壞檔"
	}
	return fmt.Sprintf("%d年%d月%d日", w.Clock.Year, w.Clock.Month, w.Clock.Day)
}

// SaveSlot 存進第 slot 槽。
//
// ⭐ **原版格式與原生檔一起寫**。原版格式是這個專案的保存價值所在，
// 原生檔才裝得下執行期狀態（游標、路徑，docs/spec/20 §2.4）。
// 分成兩個動作就會多出一個「忘記匯出」的狀態。
func (s *Session) SaveSlot(slot int) error {
	if slot < 0 || slot >= SaveSlots {
		return fmt.Errorf("槽位 %d 超出 1–%d", slot+1, SaveSlots)
	}
	dst := s.overlaySave()
	if savepath.SamePath(dst, s.sourceSave()) {
		return fmt.Errorf("拒絕覆寫原始素材：%s", dst)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("建立存檔目錄失敗：%w", err)
	}
	// 已經有 overlay 就從它出發，否則從唯讀來源出發——
	// 這樣其他三槽的內容會原封不動留著。
	base := dst
	if _, err := os.Stat(dst); errors.Is(err, os.ErrNotExist) {
		base = s.sourceSave()
	}
	if err := s.world.SaveInto(base, dst, slot); err != nil {
		return fmt.Errorf("儲存第 %d 槽失敗：%w", slot+1, err)
	}
	data, err := savefile.Encode(s.world, savefile.Origin{
		Source: filepath.Base(dst), Block: slot,
	})
	if err != nil {
		return fmt.Errorf("編碼原生存檔失敗：%w", err)
	}
	path, err := savepath.NativePath(dst, slot)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadSlot 讀第 slot 槽。
//
// 優先讀原生檔——原版格式裝不下執行期狀態，從它讀回來的世界會少一部分。
// **原生檔存在但壞掉一律報錯**，不靜靜退回原版格式：
// 那會讓玩家以為讀到了最新進度。
func (s *Session) LoadSlot(slot int) error {
	if slot < 0 || slot >= SaveSlots {
		return fmt.Errorf("槽位 %d 超出 1–%d", slot+1, SaveSlots)
	}
	dst := s.overlaySave()
	var w *state.World
	path, err := savepath.NativePath(dst, slot)
	if err == nil {
		data, rerr := os.ReadFile(path)
		switch {
		case rerr == nil:
			if w, _, err = savefile.Decode(data); err != nil {
				return fmt.Errorf("第 %d 槽的原生存檔有問題：%w", slot+1, err)
			}
		case !errors.Is(rerr, os.ErrNotExist):
			return fmt.Errorf("讀取原生存檔失敗：%w", rerr)
		}
	}
	if w == nil {
		src := dst
		if _, serr := os.Stat(dst); errors.Is(serr, os.ErrNotExist) {
			src = s.sourceSave()
		}
		if w, err = state.LoadScenario(src, slot); err != nil {
			return fmt.Errorf("讀取第 %d 槽失敗：%w", slot+1, err)
		}
	}
	// 空槽沒有玩家欄位，沿用目前的勢力。
	if w.Player < 0 {
		w.Player = s.world.Player
	}
	s.world = w
	s.selected = -1
	// 讀檔換掉了整個 World，道路圖要重新掛上——**它不在存檔裡**，
	// 是從 `MMAP.MAP` 推出來的常量（`internal/assets/world.RoadEdges`）。
	s.attachRoads()
	s.centreOnCapital()
	return nil
}
