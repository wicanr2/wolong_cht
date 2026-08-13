package state

import (
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
)

// 原版 sub_123FF 配置 32 筆、每筆 0x10 byte 的 runtime 物件記錄。
// 這些欄位不在 SINARIO／SAVE 區塊裡；只保存原版 map loop 需要的動畫
// 時序，不把它誤寫成持久化遊戲規則。
const (
	disasterObjectSlots    = 32
	disasterObjectInterval = 0x10
	disasterMovingSlot     = 0x10 // sub_12459 calls sub_1248A for SI >= 2140h.

	// sub_1248A 的固定回繞界線不是目前畫布尺寸；它直接寫入 raw
	// word_10D22/24/26/28 失界時的 16-bit 常數。X/Y 物件座標仍以
	// 有號 16-bit word 保存，避免把 FFF0h 誤當成極大的正座標。
	disasterWrapMinX = -0x10
	disasterWrapMaxX = 0x190
	disasterWrapMinY = -0x10
	disasterWrapMaxY = 0x110
)

type disasterObject struct {
	active   bool
	kind     economy.Disaster
	typeCode uint8
	city     int
	x, y     int
	timer    uint8
	interval uint8
	phase    uint8
	dirty    bool
	// xDrift/yDrift 對應原版 [si+08]/[si+0A]，各自是低 byte 在前的
	// 16-bit 固定點 word；[si+09]/[si+0B] 是同一 word 的高 byte，
	// 不要拆成 Go 的浮點速度，否則 sub_124FF 的 signed byte clamp
	// 與 byte wrap 會改變邊界行為。
	xDrift uint16
	yDrift uint16
}

// DisasterObjectSnapshot 是給呈現層的原版 runtime 物件快照。
// Phase 是這次繪製要使用的相位；原版在 dirty frame 先畫舊相位，
// 再把記錄內的 +0Fh 遞增，所以 RenderDisasterObjects 會在回傳後更新
// 內部相位，而不是提前一格。
type DisasterObjectSnapshot struct {
	Slot     int
	City     int
	Kind     economy.Disaster
	TypeCode uint8
	X, Y     int
	Phase    uint8
}

func disasterObjectType(kind economy.Disaster) (uint8, bool) {
	switch kind {
	case economy.Fire:
		return 1, true // sub_134B1 的 AH=1 → sub_123FF AL=1
	case economy.Riot:
		return 2, true // sub_134B1 的 AH=2 → sub_123FF AL=2
	default:
		return 0, false
	}
}

// createDisasterObject 是 sub_123FF 的最小已證實轉接。
// 原版掃到第一筆 [si] < 80h 的空槽；不對同一城市做去重。
func (w *World) createDisasterObject(city int, kind economy.Disaster) bool {
	if city < 0 || city >= len(w.Cities) {
		return false
	}
	typeCode, ok := disasterObjectType(kind)
	if !ok {
		return false
	}
	for slot := range w.disasterObjects {
		if w.disasterObjects[slot].active {
			continue
		}
		c := w.Cities[city]
		w.disasterObjects[slot] = disasterObject{
			active:   true,
			kind:     kind,
			typeCode: typeCode,
			city:     city,
			x:        c.X,
			y:        c.Y,
			timer:    1, // sub_123FF [si+0Ch] = 1
			interval: disasterObjectInterval,
			phase:    1, // sub_123FF [si+0Fh] = 1
		}
		return true
	}
	return false
}

// clearDisasterObjects 是 sub_12438 的座標比對清除。
// 原版會清掉所有相同城市座標的 active record，而不是只清第一筆。
func (w *World) clearDisasterObjects(city int) {
	for i := range w.disasterObjects {
		if w.disasterObjects[i].active && w.disasterObjects[i].city == city {
			w.disasterObjects[i] = disasterObject{}
		}
	}
}

// AdvanceMapObjects 在 remake 的明示暫停（speed=0）時單獨跑可見物件動畫；
// 正常完整順序由 World.TickMap 以「據點／軍團／物件／時鐘」完成。rng 必須
// 在物件更新前綁定，讓 sub_1248A 的 moving slot 使用本次 map-loop 的亂數流。
func (w *World) AdvanceMapObjects(rng economy.Rand) {
	w.rng = rng
	w.AdvanceDisasterObjects()
}

// AdvanceDisasterObjects 對應 sub_12459，一次代表一次 map-loop 更新。
// timer 歸零時重載 16、設 dirty bit；相位不在這裡遞增，因為原版
// sub_12533 是後面的 render loop 才先畫舊相位再遞增 +0Fh。
func (w *World) AdvanceDisasterObjects() {
	for i := range w.disasterObjects {
		o := &w.disasterObjects[i]
		if !o.active {
			continue
		}
		if o.timer > 0 {
			o.timer--
		}
		if o.timer == 0 {
			o.timer = o.interval
			o.dirty = true
			// 原版只讓後半段的 16 筆 runtime record 走
			// sub_1248A；前半段只換動畫 phase。沒有 Tick 提供的
			// 亂數源時，仍保留 timer／dirty 行為，但不猜測移動結果。
			if i >= disasterMovingSlot && w.rng != nil {
				w.advanceMovingDisasterObject(o)
			}
		}
	}
}

// advanceMovingDisasterObject 是 DOS/V sub_1248A（0001248A）的 typed
// 轉接。sub_1248A 呼叫兩次 sub_124FF（000124FF），再以 storm bounds
// word_10D22/24/26/28 決定方向 byte，最後把失界座標回繞到固定 raw
// 常數。這段只服務 runtime object，不進存檔。
func (w *World) advanceMovingDisasterObject(o *disasterObject) {
	if o == nil || w == nil || w.rng == nil {
		return
	}
	minX, minY := disasterWrapMinX, disasterWrapMinY
	maxX, maxY := disasterWrapMaxX, disasterWrapMaxY
	if w.stormArea != nil {
		minX, minY = w.stormArea.MinX, w.stormArea.MinY
		maxX, maxY = w.stormArea.MaxX, w.stormArea.MaxY
	}

	var step int
	o.xDrift, step = sub124FF(o.xDrift, w.rng.Next())
	o.x = addSignedWord(o.x, step)
	o.yDrift, step = sub124FF(o.yDrift, w.rng.Next())
	o.y = addSignedWord(o.y, step)

	dx := 0
	if o.x < minX {
		dx = 1
	} else if o.x > maxX {
		dx = -1
	}
	if o.x < disasterWrapMinX {
		o.x = disasterWrapMaxX
	} else if o.x > disasterWrapMaxX {
		o.x = disasterWrapMinX
	}

	dy := 0
	if o.y < minY {
		dy = 1
	} else if o.y > maxY {
		dy = -1
	}
	if o.y < disasterWrapMinY {
		o.y = disasterWrapMaxY
	} else if o.y > disasterWrapMaxY {
		o.y = disasterWrapMinY
	}

	o.xDrift = addHighByte(o.xDrift, dx)
	o.yDrift = addHighByte(o.yDrift, dy)
}

// sub124FF 重現原版的 signed-byte 固定點正規化。回傳值是 AX 的
// 有號位移（0、1 或 -1），newDrift 是原始 DX word 的 byte-wrap 結果。
func sub124FF(drift uint16, randomByte int) (newDrift uint16, whole int) {
	lo := int(drift & 0xFF)
	hi := int(int8(byte(drift >> 8)))
	r := randomByte & 7
	if r <= 2 {
		hi += r - 1
	}
	if hi < -15 {
		hi = -15
	}
	if hi > 15 {
		hi = 15
	}

	lo = (lo + hi) & 0xFF
	signedLo := int(int8(byte(lo)))
	switch {
	case signedLo >= 15:
		lo = (lo - 15) & 0xFF
		whole = 1
	case signedLo <= -15:
		lo = (lo + 15) & 0xFF
		whole = -1
	}
	return uint16(byte(lo)) | uint16(byte(hi))<<8, whole
}

func addSignedWord(value, delta int) int {
	return int(int16(uint16(value) + uint16(int16(delta))))
}

func addHighByte(word uint16, delta int) uint16 {
	return (word & 0x00FF) | uint16(byte(int(word>>8)+delta))<<8
}

// RenderDisasterObjects 對應 sub_12533 的一次物件繪製掃描。
// 每筆 dirty record 以目前 phase 畫出，再清 dirty 並把 phase 遞增八循環。
// 回傳的新 slice 不含任何可修改內部狀態的指標。
func (w *World) RenderDisasterObjects() []DisasterObjectSnapshot {
	out := make([]DisasterObjectSnapshot, 0, disasterObjectSlots)
	for slot := range w.disasterObjects {
		o := &w.disasterObjects[slot]
		if !o.active {
			continue
		}
		out = append(out, DisasterObjectSnapshot{
			Slot: slot, City: o.city, Kind: o.kind, TypeCode: o.typeCode,
			X: o.x, Y: o.y, Phase: o.phase,
		})
		if o.dirty {
			o.dirty = false
			o.phase = (o.phase + 1) & 7
		}
	}
	return out
}
