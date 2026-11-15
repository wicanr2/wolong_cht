package state

// 松崗 DOS/V 標準配置的相對位置，出處為 sub_100DF 與 spec271。
// 這是原版到站錯位讀取的資料映射，不是宿主或模擬器的實體位址。
const (
	dosvWorldToIcons = 0x5240
	dosvWorldToMap   = 0x75f0
	dosvIconBankSize = 0x23a0
	aiHoldStride     = 1024
)

// SetAIHoldSources 注入原版途中留守會讀到的九個圖庫值及34個地圖值。
// iconBank3 是原 ICONGRF 第3段；tile 只查原版384×256地圖。
// 不載入或依賴HD圖片。來源缺席時保留既有未知近似，不能硬編素材值。
func (w *World) SetAIHoldSources(iconBank3 []byte, tile func(x, y int) byte) {
	w.aiHoldIconsReady = len(iconBank3) == dosvIconBankSize
	if w.aiHoldIconsReady {
		for i := range w.aiHoldIcons {
			w.aiHoldIcons[i] = iconBank3[(i+21)*aiHoldStride+24-dosvWorldToIcons]
		}
	}
	w.aiHoldMapReady = tile != nil
	if tile != nil {
		for i := range w.aiHoldMap {
			off := (i+30)*aiHoldStride + 24 - dosvWorldToMap
			w.aiHoldMap[i] = tile(off%384, off/384)
		}
	}
}

// residualField18 讀本次到站的 DI+18h。世界範圍沿既有typed寫回契約；
// 留守的存檔外讀點使用原素材。正常呼叫前均已重建 DI，不需跨存檔保存。
func (w *World) residualField18() int {
	d := w.stage1DI
	if d < 0 {
		return aiCrowded + 1
	}
	at := (d + 24) & 0xffff
	if at < dosvWorldToIcons && len(w.raw) == blockSize {
		return int(w.Bytes()[factionBase+at])
	}
	if d%aiHoldStride == 0 {
		n := d / aiHoldStride
		if n >= 21 && n <= 29 && w.aiHoldIconsReady {
			return int(w.aiHoldIcons[n-21])
		}
		if n >= 30 && n <= 63 && w.aiHoldMapReady {
			return int(w.aiHoldMap[n-30])
		}
	}
	// 沒有原始區塊的純規則fixture，保留既有欄位查詢。
	if d >= runtimeCityBase {
		if city := (d - runtimeCityBase) / citySize; validCity(city) {
			return w.Cities[city].Occupancy
		}
	} else if f := d / factionSize; f >= 0 && f < numFactions {
		return w.Factions[f].Generals
	}
	return aiCrowded + 1
}
