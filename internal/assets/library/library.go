// Package library 把原版素材目錄載成一組可檢視的項目。
//
// **這一層刻意不 import Ebiten。** Ebiten 在 init 期就要求顯示器，
// 混進來會讓無頭環境連截圖工具都跑不起來 —— 這個坑實際踩過一次：
// 第一版把截圖功能和 Ebiten 檢視器放同一個 package，結果
// `-shot` 模式在容器裡直接 panic 在 ebiten 的 init。
package library

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/wicanr2/wolong_cht/internal/assets/gfx"
	"github.com/wicanr2/wolong_cht/internal/assets/palette"
	"github.com/wicanr2/wolong_cht/internal/assets/text"
)

// Entry 是一種可檢視的素材。
type Entry struct {
	Label string
	Spec  gfx.Spec
	Data  []byte
	Count int
}

// Library 是一個原版素材目錄解出來的全部東西。
type Library struct {
	Palette *palette.Palette
	Talk    *text.Table
	Entries []Entry
	Warns   []string // 餘數不是 0 之類的警告，呼叫端要顯示出來
}

func read(dir, name string) ([]byte, error) {
	b, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return nil, fmt.Errorf("讀不到 %s：%w\n"+
			"原版素材不隨本專案散布，請用 -orig 指向自備的原版目錄", name, err)
	}
	return b, nil
}

// Load 載入一個原版素材目錄（dosv 或 pc98 都可以）。
func Load(dir string) (*Library, error) {
	palRaw, err := read(dir, "GAMEPAL.BRG")
	if err != nil {
		return nil, err
	}
	pal, err := palette.Parse(palRaw)
	if err != nil {
		return nil, err
	}
	talkRaw, err := read(dir, "TALK.DAT")
	if err != nil {
		return nil, err
	}
	talk, err := text.Parse(talkRaw)
	if err != nil {
		return nil, err
	}

	lib := &Library{Palette: pal, Talk: talk}
	add := func(label string, spec gfx.Spec, data []byte) {
		n, rem := spec.Count(data)
		if rem != 0 {
			// 餘數不是 0 就代表尺寸錯了。不靜悄悄吞掉——
			// 那正是「看起來能跑但其實解錯」的來源。
			lib.Warns = append(lib.Warns, fmt.Sprintf(
				"%s 餘 %d byte，尺寸 %dx%d 可能是錯的",
				label, rem, spec.Width, spec.Height))
		}
		lib.Entries = append(lib.Entries, Entry{label, spec, data, n})
	}

	for _, c := range []struct {
		file string
		spec gfx.Spec
		name string
	}{
		{"KAOGRF.DAT", gfx.Kao, "KAOGRF portraits"},
		{"KYOGRF.DAT", gfx.Kyo, "KYOGRF locations"},
		{"IVENTGRF.DAT", gfx.Ivent, "IVENTGRF cutscenes"},
	} {
		data, err := read(dir, c.file)
		if err != nil {
			return nil, err
		}
		add(c.name, c.spec, data)
	}

	icon, err := read(dir, "ICONGRF.DAT")
	if err != nil {
		return nil, err
	}
	for _, r := range gfx.IconRegions {
		if r.Spec.Width == 0 {
			continue // 段 1、段 3 的尺寸還沒解，不假裝畫得出來
		}
		add("ICONGRF/"+r.Name, r.Spec, icon[r.Offset:r.Offset+r.Length])
	}
	return lib, nil
}

// Render 畫出第 asset 種素材的第 page 張。
func (l *Library) Render(asset, page, bank int) (*image.RGBA, error) {
	if asset < 0 || asset >= len(l.Entries) {
		return nil, fmt.Errorf("素材編號 %d 超出範圍（共 %d 種）",
			asset, len(l.Entries))
	}
	e := l.Entries[asset]
	if e.Count > 0 {
		page = ((page % e.Count) + e.Count) % e.Count
	}
	return e.Spec.RenderRGBA(e.Data, page, l.Palette, bank)
}
