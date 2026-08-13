package state

import (
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/strategyai"
)

// runtimeCityBase 是原版事件 Param 指向城市記錄的段內位址；它與
// SAVE.DAT 的 cityBase（檔案偏移 0x08C0）不同，不能混用。
const runtimeCityBase = 0x0840

// queueEvent 是 sub_12FB1 → sub_12FBF 的可重播轉接。
//
// 原版的 sub_12FB1 把目前勢力的編號放在事件字高 byte，低 byte 才是
// dispatch code；sub_12FBF 以亂數決定前 64 格的搜尋起點，再向尾端找空格。
// 這裡保留同一個搜尋邊界與「不回繞」行為。slotHint 目前只供未來接入
// sub_1301C 類型的指定起點；0xFF 代表原版 sub_12FBF 的亂數起點。
func (w *World) queueEvent(rng economy.Rand, source int, code byte, param uint16, slotHint byte) bool {
	if source < 0 || source > 0xFF || code == 0 {
		return false
	}
	start := int(w.eventCursor)
	if slotHint == 0xFF {
		start += rng.Next() & 0x7C
	} else {
		start += int(slotHint) * eventQueueEntrySize
	}
	for off := start; off < eventQueueDispatch*eventQueueEntrySize; off += eventQueueEntrySize {
		i := off / eventQueueEntrySize
		if byte(w.events[i].Code) != 0 {
			continue
		}
		w.events[i] = QueuedEvent{
			Code:  uint16(source)<<8 | uint16(code),
			Param: param,
		}
		return true
	}
	return false
}

// queueFullEvent 是 sub_1301C 的指定提示位置版本。slotHint 不是 AI 的
// `sub_12FBF` 邊界；這條路徑從提示位置一路找完整 256 格。
func (w *World) queueFullEvent(source int, code byte, param uint16, slotHint byte) bool {
	if source < 0 || source > 0xFF || code == 0 {
		return false
	}
	start := int(w.eventCursor) + int(slotHint)*eventQueueEntrySize
	for off := start; off < eventQueueEntries*eventQueueEntrySize; off += eventQueueEntrySize {
		i := off / eventQueueEntrySize
		if byte(w.events[i].Code) != 0 {
			continue
		}
		w.events[i] = QueuedEvent{
			Code:  uint16(source)<<8 | uint16(code),
			Param: param,
		}
		return true
	}
	return false
}

// queuePlayerEvent 是 sub_1301C 的完整 256 格寫入路徑。
//
// AI 產生端 sub_12FBF 只從目前游標開始找前 64 格；玩家進言的
// sub_164F1／sub_16623 則傳入 BL=14h 後呼叫 sub_1301C，會從同一個
// 游標加上第 20 格的提示位置找完整 256 格。兩條路徑不能共用
// eventQueueDispatch 的上限，否則玩家事件在前 64 格滿時會被錯誤丟掉。
func (w *World) queuePlayerEvent(source int, code byte, param uint16) bool {
	return w.queueFullEvent(source, code, param, 0x14)
}

// HasQueuedDeclaration 是 sub_16475 對 sub_1304E 的玩家路徑查詢。
// 原版只比對事件 Code 與 Param 低 byte；Param 高 byte 的 FF 附帶值
// 不參與「已經排過宣戰」判定。
func (w *World) HasQueuedDeclaration(source, target int) bool {
	if source < 0 || source >= numFactions || target < 0 || target >= numFactions {
		return false
	}
	want := uint16(source)<<8 | 1
	for _, e := range w.events {
		if e.Code == want && int(byte(e.Param)) == target {
			return true
		}
	}
	return false
}

// ApplyPlayerHostility 是 sub_16405 同意敵對提案後的直接
// sub_13526 收尾。玩家宣戰不是事件 1 的延遲 handler；事件 1 是政略
// AI 使用的佇列路徑。這裡仍在寫入前重驗證和平／存活條件，避免 UI
// 進言期間世界狀態改變後把過期提案套用下去。
func (w *World) ApplyPlayerHostility(target int) bool {
	if w.Player < 0 || w.Player >= numFactions || target < 0 || target >= numFactions ||
		target == w.Player || !w.Factions[w.Player].Alive || !w.Factions[target].Alive ||
		w.Friendship[w.Player][target].AtWar() {
		return false
	}
	return w.applyQueuedDeclaration(w.Player, target)
}

// QueueEvent10 是 remake 的事件 10 raw producer。
//
// 原版 `sub_13496` 的 consumer 已證實把事件字高當成 General index、把
// Param 原樣當成 TALK.DAT index；但 DOS/V `.i64` 對 `sub_12FBF`／
// `sub_1301C` 的所有直接 caller 都沒有送入低碼 0x0A，因此找不到可重播的
// 自然觸發時序。這個入口只建立已證實的四 byte queue payload，採完整
// 256 格搜尋供劇本／測試／外部事件源注入；不宣稱這個搜尋位置就是原版
// 未定位 producer 的時序。
func (w *World) QueueEvent10(general, talkIndex int) bool {
	if general < 0 || general >= numGenerals || talkIndex < 0 || talkIndex > 0xFFFF {
		return false
	}
	return w.queueFullEvent(general, 0x0A, uint16(talkIndex), 0)
}

// QueuePlayerCeasefire 是 sub_164F1 在停戰提案成立後寫入事件 6 的
// producer。事件字高是回報方，而不是玩家；事件 6 handler 之後才以
// sub_136C4 的 SI=回報方、DI=玩家方向完成付款與停戰。
func (w *World) QueuePlayerCeasefire(target int) bool {
	if w.Player < 0 || w.Player >= numFactions || target < 0 || target >= numFactions ||
		target == w.Player || !w.Factions[w.Player].Alive || !w.Factions[target].Alive ||
		w.Factions[target].Diplomat != noFaction || !w.Friendship[w.Player][target].AtWar() {
		return false
	}
	for _, e := range w.events {
		if e.Code == uint16(target)<<8|6 {
			return false
		}
	}
	return w.queuePlayerEvent(target, 6, 0)
}

// QueuePlayerCooperation 是 sub_16623 在協力提案成立後寫入事件 7 的
// producer。Param 的高低 byte 都是侵攻目標；handler 實際只取低 byte，
// 但原始格式仍完整保留。
func (w *World) QueuePlayerCooperation(ally, invader int) bool {
	if w.Player < 0 || w.Player >= numFactions || ally < 0 || ally >= numFactions ||
		invader < 0 || invader >= numFactions || ally == w.Player || invader == w.Player ||
		ally == invader || !w.Factions[w.Player].Alive || !w.Factions[ally].Alive ||
		!w.Factions[invader].Alive || w.Factions[ally].Diplomat != noFaction ||
		!w.Friendship[w.Player][invader].AtWar() {
		return false
	}
	for _, e := range w.events {
		if e.Code == uint16(ally)<<8|7 {
			return false
		}
	}
	param := uint16(invader)<<8 | uint16(invader)
	return w.queuePlayerEvent(ally, 7, param)
}

// compactEventQueue 是原版 sub_12BD9 的月度事件佇列壓縮。
//
// 原版不是把已處理事件逐筆從頭搬走，而是每月直接丟掉前 64 格，將第
// 64 格之後的 192 格搬到開頭，再把尾端 64 格清零。游標與每十次一筆的
// 節流計數也在同一支重設。這裡先接入資料／時序，不假設尚未完成的 UI
// handler 已經存在。
func (w *World) compactEventQueue() {
	copy(w.events[:eventQueueEntries-eventQueueDispatch], w.events[eventQueueDispatch:])
	for i := eventQueueEntries - eventQueueDispatch; i < eventQueueEntries; i++ {
		w.events[i] = QueuedEvent{}
	}
	w.eventCursor = 0
	w.eventDelay = 7
}

// takeNextQueuedEvent 重現 sub_131AE 的節流與取格順序。
//
// 呼叫一次代表原版的一次「每時」事件處理嘗試。回傳 false 可能表示尚未
// 到第十次，也可能是該格為空；呼叫端日後接 handler 時仍要依 Code 的
// 低 byte dispatch，不能把高 byte 當成另一個事件代碼。
func (w *World) takeNextQueuedEvent() (QueuedEvent, bool) {
	if w.eventDelay != 1 {
		w.eventDelay--
		return QueuedEvent{}, false
	}
	w.eventDelay = 0
	if w.eventCursor >= eventQueueDispatch*eventQueueEntrySize {
		return QueuedEvent{}, false
	}

	w.eventDelay = 10
	e := w.events[w.eventCursor/eventQueueEntrySize]
	w.eventCursor += eventQueueEntrySize
	if byte(e.Code) == 0 {
		return QueuedEvent{}, false
	}
	return e, true
}

// dispatchQueuedEvent 接入目前已有獨立證據的 dispatch handler。
//
// 事件 1（sub_1320C → sub_13526）與事件 8（sub_133EA → sub_16A3D）
// 會改變政略狀態；事件 2（sub_13220 → sub_13712）完成合作的狀態部分；
// 事件 3（sub_13262 → sub_136C4）完成停戰的狀態部分；
// 事件 4（sub_132A9 → sub_139E8）與事件 5（sub_132E9 → sub_139E8）
// 會停在玩家撥款視窗；事件 6（sub_13327 → sub_136C4）與事件 7
// （sub_13388 → sub_13712）會自動完成外交官回報的狀態收尾；
// 事件 9（sub_13485 → sub_150D7）釋放指定俘虜武將的狀態部分；事件 10
// （sub_13496）保留 raw producer／訊息 consumer 邊界；事件 11／12 接入暴風雨／火災／暴動
// 的 runtime marker、sub_14269 持久效果與事件 12 延遲清除；事件 13
// （sub_13507）會扣玩家信賴度。尚未解出的完整訊息、物件動畫與其他
// dispatch code 仍只取出，不把未證實資料流當成已知效果。
func appendDiplomacyReportNotice(ev *Event, index, faction, general, amount int) {
	if ev == nil {
		return
	}
	ev.TalkNotices = append(ev.TalkNotices, TalkNotice{
		Index: index, City: -1, Faction: faction, General: general, Amount: amount,
	})
}

// diplomacyCaptiveFlags 是 sub_137D8 → sub_13138 的兩個方向 bit。
// sub_13138 不檢查 General 的 active byte，只比 +1Dh（Captor）與 +1Ch
// （Faction）；因此這裡也保留 raw record 的欄位語意，不擅自加 Alive gate。
func (w *World) diplomacyCaptiveFlags(a, b int) uint8 {
	var flags uint8
	for _, g := range w.Generals {
		if g.Captor == a && g.Faction == b {
			flags |= 1
		}
		if g.Captor == b && g.Faction == a {
			flags |= 2
		}
		if flags == 3 {
			break
		}
	}
	return flags
}

func appendDiplomacySecondaryNotice(ev *Event, index int, noPortrait bool, rawWord int, rawWordValid bool) {
	if ev == nil {
		return
	}
	ev.TalkNotices = append(ev.TalkNotices, TalkNotice{
		Index: index, City: -1, Faction: -1, General: -1, Amount: -1,
		RawFormatterWord: rawWord, RawFormatterWordValid: rawWordValid,
		Secondary: true, NoPortrait: noPortrait,
	})
}

func (w *World) dispatchQueuedEvent(ev *Event) {
	e, ok := w.takeNextQueuedEvent()
	if !ok {
		return
	}
	source := int(e.Code >> 8)
	switch byte(e.Code) {
	case 1:
		if w.applyQueuedDeclaration(source, int(byte(e.Param))) {
			// 宣戰的可觀測紀錄在入佇列的月結事件已產生；此處不重複計數。
		}
	case 2:
		// 事件字高 byte 是玩家／合作方；參數低 byte 是侵攻方，
		// 高 byte 是侵攻方目前的目標。三者都不能由低 byte 推導替代。
		invader, target := int(byte(e.Param)), int(byte(e.Param>>8))
		if source == w.Player {
			w.beginDiplomacy(ev, DiplomacyChoice{
				Kind: DiplomacyCooperation, Source: source, Invader: invader, Target: target,
			})
			return
		}
		w.applyQueuedCooperation(source, invader, target)
	case 3:
		// sub_13262 只把參數低 byte 當成停戰對象；高 byte 是
		// sub_12FB1 的原始附帶值，不能拿來覆蓋對象。
		target := int(byte(e.Param))
		if target == w.Player {
			w.beginDiplomacy(ev, DiplomacyChoice{
				Kind: DiplomacyCeasefire, Source: source, Target: target,
			})
			return
		}
		w.applyQueuedCeasefire(source, target)
	case 4:
		// sub_132A9 用事件字高還原據點索引，Param 是
		// sub_15715 算出的要求金額；處理時重新確認 +19h 仍指向同一名
		// 內政官，不能把事件發送時的舊指標直接沿用。
		cityID := source
		if cityID < 0 || cityID >= numCities {
			return
		}
		officer := w.Cities[cityID].Governor
		if w.beginFunding(ev, FundingChoice{
			Kind: FundingGovernor, Subject: cityID, Officer: officer,
			RequestedAmount: int(e.Param),
		}) && ev != nil {
			// sub_132A9 在進入 sub_139E8 前先以 CX=38h 顯示
			// TALK #56；DI 指向城市、堆疊參數指向內政官。
			ev.TalkNotices = append(ev.TalkNotices, TalkNotice{
				Index: 0x38, City: cityID, Faction: -1, General: officer, Amount: -1,
			})
		}
	case 5:
		// sub_132E9 的字高是勢力索引；+2Ah 是派駐該勢力的外交官。
		factionID := source
		if factionID < 0 || factionID >= numFactions {
			return
		}
		officer := w.Factions[factionID].Diplomat
		if w.beginFunding(ev, FundingChoice{
			Kind: FundingDiplomat, Subject: factionID, Officer: officer,
			RequestedAmount: int(e.Param),
		}) && ev != nil {
			// sub_132E9 在進入 sub_139E8 前先以 CX=39h 顯示
			// TALK #57；\3 是勢力君主名，\1 是外交官名。
			ev.TalkNotices = append(ev.TalkNotices, TalkNotice{
				Index: 0x39, City: -1, Faction: factionID, General: officer, Amount: -1,
			})
		}
	case 6:
		// sub_13327 的事件字高是回報停戰的對方勢力；它先檢查
		// 該勢力 +0x2A 仍有外交官，再以 sub_136C4(SI=對方、
		// DI=玩家) 結算。換成目前 state 的參數方向，就是玩家付給
		// 對方：applyQueuedCeasefire(Player, source)。
		if source < 0 || source >= numFactions || source == w.Player ||
			w.Player < 0 || w.Player >= numFactions ||
			!w.Factions[source].Alive || w.Factions[source].Diplomat == noFaction {
			return
		}
		// sub_13327 先以 #39（TALK #57）報告外交官與回報勢力，
		// 再由 sub_136C4 決定 #2B／#2C／#2D（TALK #43／#44／#45）。
		// `DX` 是 sub_136C4 乘 0x3E8 後的金額，#2C 的 \\7 讀它。
		appendDiplomacyReportNotice(ev, 0x39, source, w.Factions[source].Diplomat, -1)
		response, amount, ok := w.ceasefireTerms(w.Player, source)
		if !ok {
			appendDiplomacyReportNotice(ev, 0x3A, -1, -1, -1)
			return
		}
		noticeAmount := -1
		if response == 1 {
			noticeAmount = amount
		}
		appendDiplomacyReportNotice(ev, 0x2B+response, source, -1, noticeAmount)
		// sub_13C3D 的第二次直接呼叫條件：AH（雙向俘虜關係）非零，
		// 且 AL 不是 2／3。索引是 CX+1Dh = #72。第一個 formatter 呼叫
		// 讀 SS:[SP] 的保存 DI；返回後第二個呼叫沿用恢復的 DI，`\\2`
		// 會讀 SS:[DI]。那是當次 DOS stack 的暫存 word，不能由 faction
		// record 或持久 World state 重建；保留 raw index，但沒有動態 trace
		// 捕捉 payload 時一律明確 invalid，讓呈現層 fail-closed。
		if response < 2 && w.diplomacyCaptiveFlags(source, w.Player) != 0 {
			appendDiplomacySecondaryNotice(ev, 0x48, false, -1, false)
		}
		if response < 2 {
			w.finishQueuedCeasefire(w.Player, source, response, amount)
		}
	case 7:
		// sub_13388 的字高是協力方，Param 低 byte 是協力方要攻擊的
		// 侵攻目標；Param 高 byte 只是原始呼叫端留下的同值欄位，
		// handler 只用 DL。協力方必須仍有外交官，否則不套用效果。
		invader := int(byte(e.Param))
		if source < 0 || source >= numFactions || invader < 0 || invader >= numFactions ||
			source == w.Player || invader == w.Player || source == invader ||
			w.Player < 0 || w.Player >= numFactions ||
			!w.Factions[source].Alive || !w.Factions[invader].Alive ||
			w.Factions[source].Diplomat == noFaction {
			return
		}
		// sub_13388 與事件 6 共用 #39，成功／金額／失敗分別是
		// #2F／#30／#31（TALK #47／#48／#49）。
		appendDiplomacyReportNotice(ev, 0x39, source, w.Factions[source].Diplomat, -1)
		response, amount, ok := w.cooperationTerms(source, invader, w.Player)
		if !ok {
			appendDiplomacyReportNotice(ev, 0x3A, -1, -1, -1)
			return
		}
		noticeAmount := -1
		if response == 1 {
			noticeAmount = amount
		}
		appendDiplomacyReportNotice(ev, 0x2F+response, source, -1, noticeAmount)
		// 同一個 sub_13C3D 條件；事件 7 的 CX+1Dh 是 #76。
		// #76 本身沒有 formatter marker，原版是直接文字／選單樣式
		// TALK，因此不附帶玩家君主肖像。
		if response < 2 && w.diplomacyCaptiveFlags(source, w.Player) != 0 {
			appendDiplomacySecondaryNotice(ev, 0x4C, true, -1, false)
		}
		if response < 2 {
			w.finishQueuedCooperation(source, invader, w.Player, response, amount)
		}
	case 8:
		if source < 0 || source >= numFactions || source == w.Player || !w.Factions[source].Alive {
			return
		}
		if next := w.relocateCapital(source); next != capitalNone {
			if ev.Relocated == nil {
				ev.Relocated = map[int]int{}
			}
			ev.Relocated[source] = next
		}
	case 9:
		// sub_13485 先把 AX（完整事件 Code）右移三次後加上
		// General 基址；因此高 byte 是武將索引，低 byte 9 不是勢力。
		// sub_150D7 的訊息索引已由呈現層取用；state 只把已證實欄位與
		// ReleasedGenerals 觀測事件寫出，不在這裡混入對話框排版。
		if w.releaseGeneral(source) && ev != nil {
			ev.ReleasedGenerals = append(ev.ReleasedGenerals, source)
		}
	case 10:
		// sub_13496：AL=事件字高、CX=Param，先組成 FFxx 的 formatter
		// word，再以 CX 呼叫 sub_18810。\1 handler 會把 FFxx 解成
		// General index；這是事件 10 唯一已證實的呈現資料流。
		general := -1
		if source >= 0 && source < numGenerals {
			general = source
		}
		if ev != nil {
			ev.TalkNotices = append(ev.TalkNotices, TalkNotice{
				Index: int(e.Param), City: -1, Faction: -1,
				General: general, Amount: -1,
			})
		}
	case 11:
		// sub_134A6 更新的是城市 +0x15 的暴風雨動畫標記。這個欄位
		// 不被月結規則讀取，因此只保留在 runtime，不寫回 City 持久欄位。
		w.applyQueuedStormMarker(ev)
	case 12:
		// sub_134B1 的高 byte 1／2 是火災／暴動，0 是延遲清除；
		// Param 是 runtimeCityBase + city×0x20 的段內位址，不是 city ID。
		w.applyQueuedDisasterMarker(ev, source, e.Param)
	case 13:
		// sub_13507 固定把 AL 設成 0x32，再進 sub_13DC9；DX
		// 只供訊息 formatter 使用，不是扣除量。
		w.AdjustTrust(-50)
		if ev != nil {
			// sub_13507 的固定 CX=33h 對應 TALK.DAT #51；
			// 這是玩家深度赤字事件的君主警告，不是另一個信賴度欄位。
			ev.TalkNotices = append(ev.TalkNotices, TalkNotice{
				Index: 0x33, City: -1, Faction: -1, General: -1, Amount: -1,
			})
		}
	}
}

// cityIDFromRuntimeParam 把事件 12 的原版段內指標轉成目前 state 的索引。
func cityIDFromRuntimeParam(param uint16) (int, bool) {
	off := int(param) - runtimeCityBase
	if off < 0 || off >= numCities*citySize || off%citySize != 0 {
		return 0, false
	}
	return off / citySize, true
}

// applyQueuedStormMarker 是 sub_134A6 → sub_1237E 的 runtime marker
// 轉接。原版以暴風雨中心的曼哈頓距離圈選城市，並把強度減去距離的
// 一半；這個 +0x15 marker 由 sub_14269 在據點輪轉時消耗。
func (w *World) applyQueuedStormMarker(ev *Event) {
	if w.stormArea == nil || w.rng == nil {
		return
	}
	strength := byte(w.rng.Next()&0x0F + 0x18)
	centerX := w.stormArea.MinX + 5
	centerY := w.stormArea.MinY + 5
	for i, c := range w.Cities {
		dx := absInt(centerX - c.X)
		dy := absInt(centerY - c.Y)
		if dx > dy {
			dx = dy
		}
		if dx > 0x14 {
			continue
		}
		level := int(strength) - dx/2
		if level <= 0 {
			continue
		}
		w.disasterMarkers[i] = economy.Storm
		w.disasterMarkerLevels[i] = byte(level)
		if ev != nil && c.Owner == w.Player {
			// sub_1237E 以城市記錄作為 TALK.DAT \\2 的 formatter
			// 游標，CX=46h 對應 #70「發生了暴風雨」。
			ev.TalkNotices = append(ev.TalkNotices, TalkNotice{
				Index: 0x46, City: i, Faction: -1, General: -1, Amount: -1,
			})
		}
	}
	if ev != nil && ev.Storm == nil {
		area := *w.stormArea
		ev.Storm = &area
	}
}

// applyQueuedDisasterMarker 是 sub_134B1 的兩階段 runtime marker 狀態。
// marker 寫入後的持久數值效果由 sub_14269（applyCityDisasterEffect）在
// 據點輪轉時執行；這裡仍只負責事件 12 的 marker／清除排程與通知。
func (w *World) applyQueuedDisasterMarker(ev *Event, source int, param uint16) {
	city, ok := cityIDFromRuntimeParam(param)
	if !ok {
		return
	}
	if source == 0 {
		w.disasterMarkers[city] = economy.NoDisaster
		w.disasterMarkerLevels[city] = 0
		w.clearDisasterObjects(city)
		return
	}
	var kind economy.Disaster
	switch source {
	case 1:
		kind = economy.Fire
	case 2:
		kind = economy.Riot
	default:
		return
	}
	w.disasterMarkers[city] = kind
	// sub_134B1 先以城市座標建立 MCH runtime object，再由後續 map
	// loop 依 +0Ch／+0Fh 驅動動畫；物件滿 32 筆時保留 marker，但不
	// 虛構一個不存在的動畫槽。
	w.createDisasterObject(city, kind)
	if w.rng != nil {
		w.disasterMarkerLevels[city] = byte(w.rng.Next()&7 + 4)
		// sub_134B1 第二次亂數決定動畫延遲 6..13，然後用
		// sub_1301C 的完整 256 格 producer 排入清除事件。
		delay := byte(w.rng.Next()&7 + 6)
		w.queueFullEvent(0, 12, param, delay)
	}
	if ev != nil {
		if ev.Disaster == nil {
			ev.Disaster = map[int]economy.Disaster{}
		}
		ev.Disaster[city] = kind
		if w.Cities[city].Owner == w.Player {
			// sub_134B1 的 CX=46h + AH：AH=1／2 分別得到
			// #71／#72；DI 指向同一筆城市記錄，供 TALK.DAT \\2 使用。
			ev.TalkNotices = append(ev.TalkNotices, TalkNotice{
				Index: 0x46 + source, City: city, Faction: -1, General: -1, Amount: -1,
			})
		}
	}
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// applyQueuedCooperation 是 sub_13220 的狀態部分。
//
// payload 是 Code 高 byte＝合作方、Param 低 byte＝侵攻方、Param 高 byte＝
// 被侵攻方。sub_13712 通過時，sub_135ED(AL=1) 由被侵攻方付出說服金額給
// 合作方，再呼叫 sub_13526 讓合作方對侵攻方進入既有宣戰收尾。
func (w *World) applyQueuedCooperation(source, invader, target int) bool {
	response, amount, ok := w.cooperationTerms(source, invader, target)
	if !ok || response >= 2 {
		return false
	}
	return w.finishQueuedCooperation(source, invader, target, response, amount)
}

// cooperationTerms 是 sub_13712 的條件／金額部分。response 仍保留原版
// AL：0＝無條件、1＝付費、2＝拒絕；玩家互動可在已驗證條件上選擇 0／1。
func (w *World) cooperationTerms(source, invader, target int) (response, amount int, ok bool) {
	if source < 0 || source >= numFactions || invader < 0 || invader >= numFactions ||
		target < 0 || target >= numFactions || source == invader || source == target || invader == target {
		return 0, 0, false
	}
	if !w.Factions[source].Alive || !w.Factions[invader].Alive || !w.Factions[target].Alive {
		return 0, 0, false
	}

	if _, ok := w.diplomacyRepresentative(source, target); !ok {
		return 0, 0, false
	}

	// sub_13712：先比較合作方→被侵攻方與合作方→侵攻方的原始值；
	// 再以合作方好戰度的 `0x28 + 2×+0x28` 作第二道 gate。
	friendToInvader := w.Friendship[source][invader].Raw()
	friendToTarget := w.Friendship[source][target].Raw()
	response = 1
	if friendToTarget < friendToInvader {
		response = 2
	}
	value := friendToTarget & 0x7F
	if value < 0x28+w.Factions[source].Aggression*2 {
		response = 2
	}

	// `0x5A − value`，鉗到 0..0x3C，除二後乘 1000。
	amount = 0x5A - value
	if amount < 0 {
		amount = 0
	}
	if amount > 0x3C {
		amount = 0x3C
	}
	amount = amount / 2 * 1000
	if amount == 0 {
		response = 0
	}
	return response, amount, true
}

// finishQueuedCooperation 是 sub_135ED(AL=response) 與後續
// sub_13526 的狀態部分。它不再重算條件，方便玩家三選一沿用同一個
// 原版 payload；numeric input 尚未在此猜測。
func (w *World) finishQueuedCooperation(source, invader, target, response, amount int) bool {
	if response < 0 || response >= 2 {
		return false
	}
	// sub_135ED(AL=1) 的資金方向是 SI（合作方）增加、DI（被侵攻方）
	// 減少；AL=0 仍會釋放兩方俘虜但不轉帳。
	if response == 1 {
		w.Factions[source].Funds = economy.ClampFunds(w.Factions[source].Funds + amount)
		w.Factions[target].Funds = economy.ClampFunds(w.Factions[target].Funds - amount)
	}
	w.releaseDiplomaticCaptives(source, target)

	// sub_13220 最後以 DL=Param 低 byte 呼叫 sub_13526；這裡重用已接入
	// 的事件 1 狀態收尾，保留玩家不直接寫 +0x19 的原版邊界。
	return w.applyQueuedDeclaration(source, invader)
}

// applyQueuedCooperationReport 是事件 7（sub_13388）的狀態轉接。
//
// 原始呼叫是 sub_13712(SI=協力方、BX=侵攻目標、DI=玩家)，所以玩家是
// 被協力方要求付款的一方；完成後協力方再以 sub_13526 對侵攻目標宣戰。
// 這正好對應 finishQueuedCooperation(協力方、侵攻目標、玩家)。
func (w *World) applyQueuedCooperationReport(ally, invader int) bool {
	if w.Player < 0 || w.Player >= numFactions || ally < 0 || ally >= numFactions ||
		invader < 0 || invader >= numFactions || ally == w.Player || invader == w.Player ||
		ally == invader || w.Factions[ally].Diplomat == noFaction {
		return false
	}
	response, amount, ok := w.cooperationTerms(ally, invader, w.Player)
	if !ok || response >= 2 {
		return false
	}
	return w.finishQueuedCooperation(ally, invader, w.Player, response, amount)
}

// highestUnpostedPolitics 是 sub_137F5 的資料掃描：同勢力、未出陣的
// 武將中取政治最高者；同值保留較早的表格位置。這個 helper 不把「軍師」
// 或「外交官」當成新欄位，因為原版只讀 General +0x13／+0x17／+0x1C。
func (w *World) highestUnpostedPolitics(faction int) (int, bool) {
	best := -1
	for i, g := range w.Generals {
		if !g.Alive || g.Faction != faction || g.Posted {
			continue
		}
		if best < 0 || g.Politics > w.Generals[best].Politics {
			best = i
		}
	}
	if best < 0 {
		return 0, false
	}
	return w.Generals[best].Politics, true
}

// diplomacyRepresentative 是 sub_13771 的可回查資料流。
//
// 回傳值是原版 AL：依兩名代表的政治值與平手亂數，結果為
// `政治×2` 或 `(16−政治)×2`。它不是直接的接受旗標；sub_136C4
// 再把這個值放進停戰金額公式。未找到原版會讀到的有效武將時，採
// fail-closed，不從鄰近欄位猜代表。
func (w *World) diplomacyRepresentative(si, di int) (int, bool) {
	if si < 0 || si >= numFactions || di < 0 || di >= numFactions {
		return 0, false
	}

	// `dh`：有派駐外交官就讀 Faction +0x2A 指向的武將，否則讀
	// sub_137F5(di) 選出的對方最高政治未出陣武將。
	dh := 0
	if id := w.Factions[si].Diplomat; id != noFaction {
		if id < 0 || id >= numGenerals || !w.Generals[id].Alive {
			return 0, false
		}
		dh = w.Generals[id].Politics
	} else {
		var ok bool
		dh, ok = w.highestUnpostedPolitics(di)
		if !ok {
			return 0, false
		}
	}

	// `dl`：君主未出陣時改取本勢力最高政治未出陣武將；否則使用
	// Faction +0x01 指向的君主本人。
	lord := w.Factions[si].Lord
	if lord < 0 || lord >= numGenerals || !w.Generals[lord].Alive {
		return 0, false
	}
	dl := w.Generals[lord].Politics
	if w.Generals[lord].Posted {
		// sub_13771 first reads the parallel corps record at ds:2240h
		// (+0x00). A posted lord without an existing corps fails the same
		// validity gate; do not let the General flag alone invent a representative.
		if lord >= len(w.Corps) || !w.Corps[lord].Alive {
			return 0, false
		}
	} else {
		var ok bool
		dl, ok = w.highestUnpostedPolitics(si)
		if !ok {
			return 0, false
		}
	}

	if dl > dh {
		return dl * 2, true
	}
	if dl < dh {
		return (16 - dh) * 2, true
	}
	// sub_13771 的平手分支呼叫 sub_1ECE0 並測試 bit 0。直接由
	// dispatchQueuedEvent 呼叫的單元測試沒有 Tick 提供亂數時，保留
	// 一個固定、可重播的分支；正常 Tick 會使用 World.rng。
	if w.rng != nil && w.rng.Next()&1 != 0 {
		return (16 - dh) * 2, true
	}
	return dl * 2, true
}

// releaseGeneral 是 sub_150D7 的單一武將狀態部分。
//
// 原版先清 General +0x17（出陣中）與 +0x1D（俘虜方），再依原俘虜方
// 是否仍存在，把 +0x1C（所屬勢力）寫回俘虜方或在野。事件 9 直接指定
// General；合作／停戰則由下方配對掃描呼叫同一條寫入路徑。
func (w *World) releaseGeneral(generalID int) bool {
	if generalID < 0 || generalID >= numGenerals || !w.Generals[generalID].Alive {
		return false
	}

	g := &w.Generals[generalID]
	oldCaptor := g.Captor
	g.Posted = false
	g.Captor = noFaction
	if oldCaptor >= 0 && oldCaptor < numFactions && w.Factions[oldCaptor].Alive {
		g.Faction = oldCaptor
	} else {
		g.Faction = noFaction
	}
	return true
}

// releaseDiplomaticCaptives 是 sub_135ED → sub_150D7 的配對狀態部分。
// 和平成立後，兩方互相俘虜的武將回到原 Captor 勢力；若原勢力已滅，
// 則回到在野。畫面通知仍屬尚未接入的 UI 層，這裡只寫回已證實欄位。
func (w *World) releaseDiplomaticCaptives(a, b int) {
	for i := range w.Generals {
		g := &w.Generals[i]
		if !g.Alive || !((g.Faction == a && g.Captor == b) || (g.Faction == b && g.Captor == a)) {
			continue
		}
		w.releaseGeneral(i)
	}
}

// applyQueuedCeasefire 是 sub_13262 的狀態部分。
//
// 事件字高 byte 是提出停戰的勢力，參數低 byte 是對方。sub_136C4 先算
// 說服金額：成功且金額不為零時由提出方付給對方；若對方正在侵攻提出方，
// AL 會變成 2 而 handler 不套用效果。金額為零時原版仍會走停戰收尾。
func (w *World) applyQueuedCeasefire(source, target int) bool {
	response, amount, ok := w.ceasefireTerms(source, target)
	if !ok || response >= 2 {
		return false
	}
	return w.finishQueuedCeasefire(source, target, response, amount)
}

// applyQueuedCeasefireReport 是事件 6（sub_13327）的狀態轉接。
//
// 事件 6 的原始暫存器方向與事件 3 不同：sub_136C4 以回報方放在 SI、
// 玩家放在 DI，表示玩家是提出停戰並支付的一方；因此不能把事件字高
// 直接當成 applyQueuedCeasefire 的 source。這裡只做方向轉換，條件、
// 代表政治值、資金與停戰收尾全部沿用同一條已驗證路徑。
func (w *World) applyQueuedCeasefireReport(other int) bool {
	if w.Player < 0 || w.Player >= numFactions || other < 0 || other >= numFactions ||
		other == w.Player || w.Factions[other].Diplomat == noFaction {
		return false
	}
	return w.applyQueuedCeasefire(w.Player, other)
}

// ceasefireTerms 是 sub_136C4 的代表政治／交友度／侵攻目標計算。
func (w *World) ceasefireTerms(source, target int) (response, amount int, ok bool) {
	if source < 0 || source >= numFactions || target < 0 || target >= numFactions || source == target {
		return 0, 0, false
	}
	if !w.Factions[source].Alive || !w.Factions[target].Alive {
		return 0, 0, false
	}

	politics, ok := w.diplomacyRepresentative(target, source)
	if !ok {
		return 0, 0, false
	}

	// sub_136C4：`AL = friendship(target→source) − (aggression+2)`，
	// 負值歸零；再以 `30−AL` 加到政治結果，除二後乘 1000。
	friendship := w.Friendship[target][source].Value()
	adjusted := friendship - (w.Factions[target].Aggression + 2)
	if adjusted < 0 {
		adjusted = 0
	}
	bonus := 30 - adjusted
	if bonus < 0 {
		bonus = 0
	}
	amount = politics + bonus
	if amount < 0 {
		amount = 0
	}
	amount = amount / 2 * 1000

	response = 1
	if w.Factions[target].InvasionTarget == source {
		response = 2
	}
	if amount == 0 {
		response = 0
	}
	return response, amount, true
}

// finishQueuedCeasefire 是 sub_135ED 與 sub_145F8／sub_14236／
// sub_13669 的狀態收尾；response 由 AI 計算或玩家三選一提供。
func (w *World) finishQueuedCeasefire(source, target, response, amount int) bool {
	if response < 0 || response >= 2 {
		return false
	}
	// sub_135ED(AL=1) 的兩個 24 位資金更新；AL=0 仍會釋放俘虜，
	// 但不轉帳。World 的 Funds 已是有號整數，ClampFunds 對應原版
	// sub_15609／sub_1563B 的上下限。
	if response == 1 {
		w.Factions[target].Funds = economy.ClampFunds(w.Factions[target].Funds + amount)
		w.Factions[source].Funds = economy.ClampFunds(w.Factions[source].Funds - amount)
	}
	w.releaseDiplomaticCaptives(source, target)

	// sub_145F8：只清掉彼此互指的侵攻目標。
	if w.Factions[source].InvasionTarget == target {
		w.Factions[source].InvasionTarget = diplomacy.NoTarget
	}
	if w.Factions[target].InvasionTarget == source {
		w.Factions[target].InvasionTarget = diplomacy.NoTarget
	}

	// sub_14236：只有執行期 Owner 與記錄 OwnerRecorded 都落在兩方
	// 時才同步記錄欄位，其他歷史／資料瑕疵保留不動。
	for i := range w.Cities {
		c := &w.Cities[i]
		ownerInPair := c.Owner == source || c.Owner == target
		recordedInPair := c.OwnerRecorded == source || c.OwnerRecorded == target
		if ownerInPair && recordedInPair {
			c.OwnerRecorded = c.Owner
		}
	}

	// sub_13669：取兩個方向的 7-bit 交友度較小值，設回和平位元。
	value := w.Friendship[source][target].Value()
	if other := w.Friendship[target][source].Value(); other < value {
		value = other
	}
	w.Friendship[source][target] = diplomacy.Peace(value)
	w.Friendship[target][source] = diplomacy.Peace(value)
	return true
}

// applyQueuedDeclaration 是 sub_1320C → sub_13526 的狀態部分。
//
// 事件字高 byte 是發起勢力，參數低 byte 是目標；參數高 byte 在事件 1
// 的 handler 不參與判定。玩家勢力不會被自動填入侵攻目標，但仍會走
// 原版的回頭宣戰與雙向交戰值更新路徑。
func (w *World) applyQueuedDeclaration(source, target int) bool {
	if source < 0 || source >= numFactions || target < 0 || target > diplomacy.NoTarget {
		return false
	}
	if !w.Factions[source].Alive {
		return false
	}
	if target != diplomacy.NoTarget && target != combat.NeutralFaction &&
		(target < 0 || target >= numFactions || !w.Factions[target].Alive) {
		return false
	}

	// sub_1351A 以存在旗標重新驗證兩邊；中立（24）沒有勢力記錄，
	// 原版事件 1 仍可把它寫成目標，但不進交友度矩陣。
	if w.Factions[source].InvasionTarget < combat.NeutralFaction {
		return false
	}
	if source != w.Player {
		w.Factions[source].InvasionTarget = target
	}
	if target == diplomacy.NoTarget || target == combat.NeutralFaction {
		return true
	}

	// sub_135AB：被宣戰者若手上的目標比宣戰者更強，就先吃原本的目標；
	// 否則回頭把宣戰者設成侵攻目標。中立沒有可讀的勢力記錄，這裡採
	// fail-closed 的 remake 邊界，不重現原版對交友度表起點的越界讀。
	defender := &w.Factions[target]
	current := defender.InvasionTarget
	retaliate := current == diplomacy.NoTarget || current == combat.NeutralFaction ||
		current < 0 || current >= numFactions
	if !retaliate {
		retaliate = strategyai.Power(strategyFactionValue(w, source)) >
			strategyai.Power(strategyFactionValue(w, current))
	}
	if retaliate && target != w.Player {
		defender.InvasionTarget = source
	}

	war := strategyai.AtWarValue(w.Friendship[source][target], w.Friendship[target][source])
	w.Friendship[source][target] = war
	w.Friendship[target][source] = war
	return true
}

// strategyFactionValue 只把 sub_13091 需要的資料投影給 strategyai.Power。
func strategyFactionValue(w *World, faction int) strategyai.Faction {
	f := w.Factions[faction]
	return strategyai.Faction{
		Alive: f.Alive, Cities: f.Cities, Funds: f.Funds,
		Reserves: f.Reserves, Aggression: f.Aggression,
		InvasionTarget: f.InvasionTarget,
	}
}
