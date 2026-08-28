package state

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
)

// 世界指紋。規格 docs/spec/69。
//
// ⚠ **這不是原版機制**，是 remake 的驗證設施：把一個 `World` 縮成一個
// 可比較的值，讓「同一個 seed 跑兩次結果一樣」變成一行斷言。
//
// 涵蓋存檔位元組 ＋ 不序列化的 runtime 欄位；**不涵蓋**道路圖（常量）、
// 戰術戰鬥、待決狀態的內容與任何畫面／音訊（docs/spec/69 §3）。
// **指紋相同只代表規則層的狀態相同。**

// rngState 是亂數來源可選的自我描述。`internal/rules/rng` 的實作有它，
// 測試用的假亂數沒有——取不到就記一個明確的標記，不要靜靜跳過。
type rngState interface{ State() (byte, byte) }

// Fingerprint 回傳這個世界的決定性指紋。
func (w *World) Fingerprint() [32]byte {
	if w == nil {
		return sha256.Sum256(nil)
	}
	h := sha256.New()
	h.Write(w.Bytes())

	// ⚠ 以下一律固定寬度、固定順序寫入。**不得走 map**——
	// 走 map 的話同一個世界每次算出不同的指紋，而那種錯誤會以
	// 「Android 與桌面不一致」的形式出現，查起來會很久。
	u32 := func(v int) {
		var b [4]byte
		binary.BigEndian.PutUint32(b[:], uint32(int32(v)))
		h.Write(b[:])
	}
	u8 := func(v uint8) { h.Write([]byte{v}) }
	flag := func(b bool) {
		if b {
			h.Write([]byte{1})
			return
		}
		h.Write([]byte{0})
	}

	u32(w.cityCursor)
	u32(w.eventCursor)
	u8(w.eventDelay)
	u32(w.corpsCursor)
	for _, v := range w.cityBias {
		u32(v)
	}
	for _, m := range w.disasterMarkers {
		u8(uint8(m))
	}
	for _, v := range w.disasterMarkerLevels {
		u8(v)
	}
	for _, o := range w.disasterObjects {
		flag(o.active)
		u8(uint8(o.kind))
		u8(o.typeCode)
		u32(o.city)
		u32(o.x)
		u32(o.y)
		u8(o.timer)
		u8(o.interval)
		u8(o.phase)
	}
	flag(w.stormArea != nil)

	// 待決狀態只記「有沒有」，內容是 UI 層的暫態（docs/spec/69 §3）。
	flag(w.pending != nil)
	flag(w.diplomacy != nil)
	flag(w.funding != nil)

	if r, ok := w.rng.(rngState); ok {
		c, s := r.State()
		h.Write([]byte{1, c, s})
	} else {
		// 取不到就明講，不要與「c=0 s=0」混在一起。
		h.Write([]byte{0, 0, 0})
	}

	var out [32]byte
	copy(out[:], h.Sum(nil))
	return out
}

// FingerprintHex 是指紋的前 16 個十六進位字，給 log 與測試訊息用。
func (w *World) FingerprintHex() string {
	f := w.Fingerprint()
	return hex.EncodeToString(f[:])[:16]
}
