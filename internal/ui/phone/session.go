package phone

import (
	"fmt"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/battlesetup"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/rules/clock"
	"github.com/wicanr2/wolong_cht/internal/rules/march"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/rules/speed"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// Session 是手機版的一局。它擁有**真的**規則層——與桌面版共用
// `internal/state` 與 `internal/rules`，不是另一套邏輯
// （docs/mobile/android-plan.md §4）。
type Session struct {
	lib   *library.Library
	world *state.World
	rand  *rng.Rand

	// origDir 是原版素材的目錄。存讀檔要用它推出來源與 overlay 的位置。
	origDir string

	// camX／camY 是鏡頭左上角的世界格；zoom 只允許整數倍。
	camX, camY, zoom int

	// selected 是點中的據點編號，-1 表示沒有選。
	// ⭐ **選中不等於下令**：第一段只改這個欄位，不動 World
	//（docs/mobile/android-ux.md §3 的兩段式）。
	selected int

	// paused 讓驗收路徑可以停住時鐘。⚠ 這不是遊戲內的暫停——
	// 原版即時制沒有那個東西（docs/mobile/android-ux.md §5）。
	paused bool

	// sheet 是指令列打開的面板（sheet.go）。
	sheet sheet

	// advise 是進言流程（advise.go）。
	advise advise

	// form 是編成中的軍團（corps.go）。
	form corpsForm

	// battle 是進行中的戰術戰鬥（battle.go）。
	battle battleState

	// setup 是戰場來源。載不到就是 nil——玩家的戰鬥會走自動判定，
	// **那是降級不是錯誤**（`internal/battlesetup`）。
	setup *battlesetup.Provider

	// lastErr 是最後一次存讀檔的結果。**失敗一定要看得到**——
	// 手機上沒有終端機，靜靜失敗等於資料不見了卻沒人知道。
	lastErr error

	// speed 是**原版的戰略速度檔位 0–4**（0 ＝ 最高速、4 ＝ 最低速），
	// 不是「每畫面幾 tick」。節流器與桌面版共用同一份
	//（`internal/rules/speed`，docs/spec/34）。
	speed    int
	throttle speed.Throttle
}

// Options 是開一局要的東西。字型缺席不是錯誤，中文會顯示成方框
// （與桌面版同一個行為）。
type Options struct {
	OrigDir  string
	Scenario int
	Player   int
	Seed     int
}

// NewSession 載入原版素材與劇本。
//
// ⚠ **缺檔要指名**：Android 端的資料是使用者自己匯入的，
// 「載入失敗」四個字對他毫無用處（docs/mobile/android-plan.md §3）。
func NewSession(opt Options) (*Session, error) {
	lib, err := library.Load(opt.OrigDir)
	if err != nil {
		return nil, fmt.Errorf("載入原版素材：%w", err)
	}
	w, err := state.LoadScenario(opt.OrigDir+"/SINARIO.DAT", opt.Scenario)
	if err != nil {
		return nil, fmt.Errorf("載入劇本 %d：%w", opt.Scenario, err)
	}
	w.Player = opt.Player
	s := &Session{
		lib: lib, world: w, rand: rng.NewFixed(opt.Seed), origDir: opt.OrigDir,
		zoom: 1, selected: -1, speed: DefaultSpeed, form: newCorpsForm(),
	}
	s.attachRoads()
	s.attachBattlefield()
	s.centreOnCapital()
	return s, nil
}

// attachBattlefield 掛上戰場來源。
//
// ⚠ 載不到就維持 nil：**玩家的戰鬥會走自動判定**，遊戲照樣跑得起來
//（`internal/battlesetup`）。這是刻意的降級——手機上使用者匯入的資料
// 可能不齊，為此開不了遊戲比少一個畫面糟得多。
func (s *Session) attachBattlefield() {
	p, setup, err := battlesetup.Load(battlesetup.Options{
		Dir: s.origDir, World: s.world, Map: s.lib.World,
	})
	if err != nil {
		return
	}
	s.setup = p
	s.world.SetTactical(setup)
}

// attachRoads 掛上道路圖。
//
// ⚠ **少了這一步行軍會走直線**，而規則層不會報錯——軍團照樣會動，
// 只是路徑與原版不同。桌面版在同一個位置做同一件事；
// 兩邊都要掛，不然「手機版與桌面版跑出不同結果」會被誤判成規則層的問題。
func (s *Session) attachRoads() {
	if s.lib == nil || s.lib.World == nil || s.world == nil {
		return
	}
	xy := make([][2]int, len(s.world.Cities))
	for i := range s.world.Cities {
		xy[i] = [2]int{s.world.Cities[i].X, s.world.Cities[i].Y}
	}
	edges, err := world.RoadEdges(s.lib.World, xy)
	if err != nil {
		return
	}
	s.world.SetRoads(march.New(len(s.world.Cities), world.MarchEdges(edges, xy)))
}

// World 讓驗收路徑取得規則層狀態（例如比對指紋，docs/spec/69）。
func (s *Session) World() *state.World { return s.world }

// SetPaused 停住或恢復時鐘，只給驗收路徑用。
func (s *Session) SetPaused(p bool) { s.paused = p }

// Selected 回傳目前選中的據點編號，-1 表示沒有。
func (s *Session) Selected() int { return s.selected }

// Camera 回傳鏡頭左上角的世界格與縮放。
func (s *Session) Camera() (x, y, zoom int) { return s.camX, s.camY, s.zoom }

// Clock 回傳遊戲時鐘。
func (s *Session) Clock() clock.Clock { return s.world.Clock }

// DefaultSpeed 是開局的戰略速度檔位。取「普通」，與桌面版的預設一致。
const DefaultSpeed = 2

// Tick 推進這一個畫面該推進的規則層步數。
//
// ⭐ **不是一畫面一 tick**：原版的速度是「推進一步之後等 N 個計時中斷」，
// 檔位不同步數就不同。節流器與桌面版共用（docs/spec/34），
// 兩邊的時間流速因此對得起來。
func (s *Session) Tick() {
	if s.world == nil {
		return
	}
	// ⭐ 戰場開著時**戰略時間停住**——原版進戰術畫面時就是這樣
	//（docs/mechanics/15-realtime.md §2）。戰術層自己有一組獨立的節流。
	if s.BattleActive() {
		s.tickBattle()
		return
	}
	if s.paused {
		return
	}
	for n := s.throttle.Steps(s.speed, 1, speed.HighSpeedStrategy); n > 0; n-- {
		s.world.Tick(s.rand)
	}
}

// Speed 回報戰略速度檔位。
func (s *Session) Speed() int { return s.speed }

// SetSpeed 換檔位。0 ＝ 最高速、4 ＝ 最低速，超出範圍夾住。
func (s *Session) SetSpeed(v int) {
	if v < 0 {
		v = 0
	}
	if v >= speed.Levels {
		v = speed.Levels - 1
	}
	s.speed = v
}

// centreOnCapital 把鏡頭移到玩家的首都。開局第一眼要看得到自己。
func (s *Session) centreOnCapital() {
	p := s.world.Player
	if p < 0 || p >= len(s.world.Factions) {
		return
	}
	cap := s.world.Factions[p].Capital
	if cap < 0 || cap >= len(s.world.Cities) {
		return
	}
	c := &s.world.Cities[cap]
	cols, rows := s.viewTiles()
	s.setCamera(c.X+world.CityCentreDX-cols/2, c.Y-rows/2)
}

// viewTiles 是目前縮放下地圖區看得到幾格。
func (s *Session) viewTiles() (cols, rows int) {
	_, _, w, h := MapRect()
	px := TilePx * s.zoom
	return w / px, h / px
}

// setCamera 夾住鏡頭，不讓它捲出世界外。
func (s *Session) setCamera(x, y int) {
	cols, rows := s.viewTiles()
	if m := world.Width - cols; x > m {
		x = m
	}
	if m := world.Height - rows; y > m {
		y = m
	}
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	s.camX, s.camY = x, y
}

// Pan 以**世界格**為單位平移鏡頭。
func (s *Session) Pan(dx, dy int) { s.setCamera(s.camX+dx, s.camY+dy) }

// SetZoom 換縮放級距，並讓畫面中心維持在同一個世界格上。
// 不做這件事的話捏合會把畫面甩到別的地方去。
func (s *Session) SetZoom(z int) {
	z = ClampZoom(z)
	if z == s.zoom {
		return
	}
	oldCols, oldRows := s.viewTiles()
	cx, cy := s.camX+oldCols/2, s.camY+oldRows/2
	s.zoom = z
	newCols, newRows := s.viewTiles()
	s.setCamera(cx-newCols/2, cy-newRows/2)
}

// SelectAt 把地圖區的邏輯座標換成世界格，選中那一格上的據點。
// 點空白會取消選取。回傳是不是選到了東西。
func (s *Session) SelectAt(lx, ly float64) bool {
	mx, my, mw, mh := MapRect()
	if lx < float64(mx) || lx >= float64(mx+mw) ||
		ly < float64(my) || ly >= float64(my+mh) {
		return false
	}
	px := float64(TilePx * s.zoom)
	tx := s.camX + int((lx-float64(mx))/px)
	ty := s.camY + int((ly-float64(my))/px)
	if n := s.cityAt(tx, ty); n >= 0 {
		s.selected = n
		return true
	}
	s.selected = -1
	return false
}

// cityHalf 是據點在地圖上的命中半徑（世界格）。
//
// ⭐ 據點的圖形不是同一個尺寸：大城 5×5 格、小城 3×3 格
//（`world.applyDecor` 的裝飾磚 222–225 對 5×5、226–229 對 3×3）。
// 半徑取 2 讓最大的城整塊都可以點，小城則多出一格餘裕——
// 手指的落點比滑鼠散，寧可寬一格（docs/mobile/android-ux.md §1）。
const cityHalf = 2

func (s *Session) cityAt(tx, ty int) int {
	best, bestD := -1, cityHalf*2+1
	for i := range s.world.Cities {
		c := &s.world.Cities[i]
		cx, cy := c.X+world.CityCentreDX, c.Y
		dx, dy := abs(cx-tx), abs(cy-ty)
		if dx > cityHalf || dy > cityHalf {
			continue
		}
		// 兩個據點的命中框重疊時取近的，不要看誰先出現在表裡。
		if d := dx + dy; d < bestD {
			best, bestD = i, d
		}
	}
	return best
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// Select 直接選中一個據點，給驗收路徑用。
// ⭐ 與 SelectAt 走同一個欄位——驗收看到的畫面與手指點出來的是同一種狀態。
func (s *Session) Select(city int) {
	if city < 0 || city >= len(s.world.Cities) {
		s.selected = -1
		return
	}
	s.selected = city
	// 選中之後要看得到它，否則截圖上只有小卡沒有目標。
	c := &s.world.Cities[city]
	cols, rows := s.viewTiles()
	s.setCamera(c.X+world.CityCentreDX-cols/2, c.Y-rows/2)
}

// CityAtTile 讓測試不必經過畫面座標就能驗命中判定。
func (s *Session) CityAtTile(tx, ty int) int { return s.cityAt(tx, ty) }

// ViewTiles 回傳目前縮放下看得到幾格，給測試與版面驗收用。
func (s *Session) ViewTiles() (cols, rows int) { return s.viewTiles() }
