package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/cutscene"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// 結局的播放。原版是另一支程式 `D7END.EXE`：十二幕過場 ＋ 一段逐字打出來的
// 結尾文字（docs/re/70、規格 docs/spec/67）。
//
// 節拍照原版換算：`sub_104B5(ax)` 等到 INT 1Ch 的計數器 ≥ ax，而處理常式
// 一次加 4，所以一個單位是四分之一個 BIOS tick（18.2 Hz）。remake 用自己的
// 固定 60 fps，換算 `1 單位 ≈ 60 / (18.2 × 4) ≈ 0.824 幀`。

const (
	// endingFPS 是 remake 的固定幀率；原版沒有固定時間基準（CLAUDE.md §3.1）。
	endingFPS = 60
	// endingUnitNum／endingUnitDen 把原版的延遲單位換成幀：18.2 Hz × 4。
	endingUnitNum = endingFPS * 10
	endingUnitDen = 728
	// 原版的四個延遲常數（docs/re/70 §4）。
	endingFadeStepUnits  = 0x0A
	endingTypeUnits      = 0x12
	endingSceneHoldUnits = 0x32
	endingFinalHoldUnits = 0x258
)

// endingFrames 把原版的延遲單位換成 remake 的幀數，至少一幀。
func endingFrames(units int) int {
	n := units * endingUnitNum / endingUnitDen
	if n < 1 {
		return 1
	}
	return n
}

type endingPhase int

const (
	endingFadeIn endingPhase = iota
	endingType               // 只有第一幕：逐字打出結尾文字
	endingHold
	endingFadeOut
	endingDone
)

type endingState struct {
	art   *cutscene.Ending
	scene int
	phase endingPhase
	step  int // 淡入淡出走到第幾階
	chars int // 已經打出幾個字
	wait  int // 還要等幾幀
}

// beginEnding 準備結局播放。素材載不到就回 nil，呼叫端退回原本的對話框——
// **不要靜靜跳過**：載不到的原因會寫進 g.lastEvent。
func (g *game) beginEnding() *endingState {
	if g == nil || g.origDir == "" {
		return nil
	}
	art, err := cutscene.LoadEnding(g.origDir)
	if err != nil {
		g.lastEvent = "結局過場載入失敗：" + err.Error()
		return nil
	}
	return &endingState{art: art, wait: endingFrames(endingFadeStepUnits)}
}

// update 推一幀。回傳 true 表示放完了。
func (e *endingState) update() bool {
	if e == nil || e.phase == endingDone {
		return true
	}
	if e.wait > 0 {
		e.wait--
		return false
	}
	switch e.phase {
	case endingFadeIn:
		e.step++
		if e.step >= cutscene.FadeSteps {
			e.step = cutscene.FadeSteps - 1
			if e.scene == 0 {
				e.phase, e.wait = endingType, endingFrames(endingTypeUnits)
			} else {
				e.phase, e.wait = endingHold, endingFrames(endingSceneHoldUnits)
			}
			return false
		}
		e.wait = endingFrames(endingFadeStepUnits)
	case endingType:
		e.chars++
		if e.chars >= e.totalChars() {
			e.phase, e.wait = endingHold, endingFrames(endingSceneHoldUnits)
			return false
		}
		e.wait = endingFrames(endingTypeUnits)
	case endingHold:
		if e.scene == cutscene.Scenes-1 {
			// 最後一幕停久一點，然後留在畫面上等玩家。
			e.phase, e.wait = endingDone, 0
			return true
		}
		e.phase, e.wait = endingFadeOut, endingFrames(endingFadeStepUnits)
	case endingFadeOut:
		e.step--
		if e.step < 0 {
			e.scene++
			e.step, e.phase = 0, endingFadeIn
			e.wait = endingFrames(endingFadeStepUnits)
			return false
		}
		e.wait = endingFrames(endingFadeStepUnits)
	}
	return false
}

func (e *endingState) totalChars() int {
	n := 0
	for _, l := range e.art.Lines {
		n += len([]rune(l))
	}
	return n
}

// level 是這一幀的亮度（0–1）。原版是 17 階換色盤；remake 用同樣的階數
// 疊一層黑，**階數照抄、算式是 remake 的**（色階算式還沒讀，docs/spec/67 §6）。
func (e *endingState) level() float64 {
	return float64(e.step) / float64(cutscene.FadeSteps-1)
}

// finished 回報最後一幕是不是已經停住了。
func (e *endingState) finished() bool { return e != nil && e.phase == endingDone }

// updateEnding 是結局播放期間唯一的輸入：放完之後任意鍵或滑鼠左鍵收掉。
//
// ⚠ 原版只認滑鼠鍵（`sub_104C7` 只開了 INT 33h）。remake 兩種都收——
// 只認滑鼠會讓人以為當掉（docs/spec/67 §3）。
func (g *game) updateEnding() error {
	if g.ending == nil {
		return nil
	}
	done := g.ending.update()
	if !done {
		return nil
	}
	if pressed(ebiten.KeyEnter) || pressed(ebiten.KeySpace) || pressed(ebiten.KeyEscape) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.ending = nil
	}
	return nil
}

func (g *game) drawEnding(screen *ebiten.Image) {
	if g.ending == nil {
		return
	}
	e := g.ending
	screen.Fill(color.RGBA{0, 0, 0, 255})
	img := e.art.Frames[e.scene]
	key := e.scene
	if e.scene == 0 && e.phase != endingType && e.art.Final != nil {
		// 「終」那一頁是打完字之後才換上去的（docs/formats/09 §6）。
		img, key = e.art.Final, -1
	}
	if img != nil {
		if g.endingCache == nil {
			g.endingCache = map[int]*ebiten.Image{}
		}
		tex, ok := g.endingCache[key]
		if !ok {
			tex = ebiten.NewImageFromImage(img)
			g.endingCache[key] = tex
		}
		screen.DrawImage(tex, nil)
	}
	if e.scene == 0 && e.chars > 0 {
		g.drawEndingText(screen, e)
	}
	// 淡入淡出：階數照原版的 17 階，疊一層黑。
	if a := 1 - e.level(); a > 0 {
		vector.DrawFilledRect(screen, 0, 0, screenW, screenH,
			color.RGBA{0, 0, 0, uint8(a * 255)}, false)
	}
}

func (g *game) drawEndingText(screen *ebiten.Image, e *endingState) {
	if g.td == nil || !g.td.Available() {
		return
	}
	left := e.chars
	ink := color.RGBA{255, 255, 255, 255}
	for row, line := range e.art.Lines {
		if left <= 0 {
			return
		}
		runes := []rune(line)
		if len(runes) > left {
			runes = runes[:left]
		}
		left -= len(runes)
		y := cutscene.TextY + row*cutscene.TextLeading
		for i, ch := range runes {
			x := cutscene.TextX + i*cutscene.TextAdvance
			g.td.Draw(screen, string(ch), x, y, ink)
		}
	}
}

// endingActive 回報現在是不是在放結局。
func (g *game) endingActive() bool {
	return g != nil && g.ending != nil
}

// maybeBeginEnding 在勝利 latch 之後起一次結局播放。
func (g *game) maybeBeginEnding() {
	if g == nil || g.world == nil || g.ending != nil || g.endingShown {
		return
	}
	if g.world.Outcome() != state.Victory {
		return
	}
	g.endingShown = true
	g.ending = g.beginEnding()
}

// openEndingFixture 直接跳到結局的第 n 幕，只給驗收用。
//
// 淡入淡出的階數推到最亮、文字全部打完——截圖要的是**穩定**的畫面，
// 而動畫中的任何一幀都會隨啟動時機變。
func openEndingFixture(g *game, scene int) {
	if g == nil {
		return
	}
	e := g.beginEnding()
	if e == nil {
		return
	}
	if scene < 0 {
		scene = 0
	}
	if scene >= cutscene.Scenes {
		scene = cutscene.Scenes - 1
	}
	e.scene, e.step, e.phase = scene, cutscene.FadeSteps-1, endingHold
	if scene == 0 {
		// 第一幕停在「字打完、還沒換成『終』那一頁」——那一格才看得到文字。
		e.phase = endingType
	}
	e.chars = e.totalChars()
	e.wait = 1 << 30
	g.ending, g.endingShown = e, true
}
