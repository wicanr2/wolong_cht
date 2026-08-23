package state

const maxDiplomacyOffer = 0x7530

// sub_13C3D 對事件 2／3 的 response=3 以 AL=0x1E 呼叫 sub_13DC9；
// 這是玩家輸入超過原始要求時的信賴度扣除，不是一般拒絕（response=2）。
const diplomacyOverOfferTrustPenalty = -0x1E

// editAmount 是 sub_17C6E 的數值核心：數字鍵把目前值左移一位再加數字，
// 另有 00／退位／最大／清零／結束輸入動作；所有結果都鉗在呼叫端給的上限。
//
// ⚠ **呼叫端只給一個數：上限。** 原版 `[bp+0]` 同時是夾值的上界與
// 「最大」鍵取的值，`si` 開場一律是 0（docs/spec/78 §1）。
func editAmount(current, max int, edit AmountEdit, digit int) (int, bool) {
	if max < 0 {
		return 0, false
	}
	switch edit {
	case AmountAppendDigit:
		if digit < 0 || digit > 9 {
			return current, false
		}
		current = current*10 + digit
	case AmountAppendHundred:
		current *= 100
	case AmountDeleteDigit:
		current /= 10
	case AmountSetMax:
		current = max
	case AmountClear:
		current = 0
	case AmountFinishInput:
		// 原版 sub_17DEA 只設 STC 讓 sub_17C6E 離開操作迴圈，
		// 不改寫目前的 SI；跨平台 UI 以確認鍵完成同一數值效果。
	default:
		return current, false
	}
	if current < 0 {
		current = 0
	}
	if current > max {
		current = max
	}
	return current, true
}

// EditAmountValue 是 `sub_17C6E` 數值核心的公開入口，給沒有自己的 pending
// 狀態的呼叫端用（財政視窗那四個熱區，docs/spec/78）。
//
// **規則只有這一份實作**——不要在呈現層再寫一次乘十與夾上限。
func EditAmountValue(current, max int, edit AmountEdit, digit int) (int, bool) {
	return editAmount(current, max, edit, digit)
}

// PendingDiplomacy 回傳事件 2／3 等待玩家處理的外交視窗副本。
// 這個狀態只存在 runtime，不寫入 SINARIO.DAT／SAVE.DAT。
func (w *World) PendingDiplomacy() *DiplomacyChoice {
	if w.diplomacy == nil {
		return nil
	}
	c := *w.diplomacy
	return &c
}

// SetDiplomacyOfferAmount 保留給測試與跨平台輔助控制的直接設定入口；
// 原版逐位編輯語意由 EditDiplomacyOfferAmount 提供。
func (w *World) SetDiplomacyOfferAmount(amount int) bool {
	if w.diplomacy == nil {
		return false
	}
	if amount < 0 {
		amount = 0
	}
	if amount > maxDiplomacyOffer {
		amount = maxDiplomacyOffer
	}
	w.diplomacy.OfferAmount = amount
	return true
}

// EditDiplomacyOfferAmount 接入 sub_17C6E 的狀態部分。按鍵掃描碼由
// cmd/wlgame 映射成 AmountEdit；這裡只處理已由原始函式證實的數值變化。
func (w *World) EditDiplomacyOfferAmount(edit AmountEdit, digit int) bool {
	if w.diplomacy == nil {
		return false
	}
	amount, ok := editAmount(w.diplomacy.OfferAmount, maxDiplomacyOffer, edit, digit)
	if !ok {
		return false
	}
	w.diplomacy.OfferAmount = amount
	return true
}

// beginDiplomacy 將原版 sub_13220／sub_13262 的玩家路徑停在三選一前。
// 代表與金額仍先走已接入的條件函式；未通過資料有效性時 fail-closed。
func (w *World) beginDiplomacy(ev *Event, choice DiplomacyChoice) bool {
	if w.diplomacy != nil {
		return false
	}
	var ok bool
	switch choice.Kind {
	case DiplomacyCooperation:
		_, choice.OfferAmount, ok = w.cooperationTerms(choice.Source, choice.Invader, choice.Target)
	case DiplomacyCeasefire:
		_, choice.OfferAmount, ok = w.ceasefireTerms(choice.Source, choice.Target)
	default:
		return false
	}
	if !ok {
		return false
	}
	choice.InitialAmount = choice.OfferAmount
	w.diplomacy = &choice
	if ev != nil {
		ev.Diplomacy = w.diplomacy
		// sub_138C7／sub_138E6 在玩家進入三選一前先顯示
		// TALK #360／#373；{3} 是提出請求的勢力君主名。
		index, faction := 0x168, choice.Source
		if choice.Kind == DiplomacyCooperation {
			index, faction = 0x175, choice.Invader
		}
		ev.TalkNotices = append(ev.TalkNotices, TalkNotice{
			Index: index, City: -1, Faction: faction, General: -1, Amount: -1,
		})
	}
	return true
}

// ResolveDiplomacy 消費目前的玩家外交選擇。
//
// 事件 2／3 的三列已按原版 AL=0／1／reject 分開；提供資金使用
// sub_13712／sub_136C4 算出的預設金額與目前已編輯的 OfferAmount。
// Trust advice、其餘亂數建議與原版逐頁訊息仍留在後續 parity 工作；本函式只接入
// 已由 sub_13C3D 證實的超額輸入扣 30 分分支。
func (w *World) ResolveDiplomacy(option DiplomacyOption) bool {
	choice := w.diplomacy
	if choice == nil || option > DiplomacyReject {
		return false
	}
	w.diplomacy = nil
	if option == DiplomacyReject {
		return false
	}

	// sub_138C7／sub_138E6 先把 sub_136C4／sub_13712 的 DX 留在
	// sub_13902 的初始編輯值；玩家關閉視窗前不會再跑一次條件函式。
	// 特別是代表政治值平手時，重算會錯誤多吃一次原版亂數。
	amount := choice.InitialAmount
	if amount < 0 || choice.OfferAmount < 0 {
		return false
	}
	response := 0
	if option == DiplomacyOfferFunds {
		amount = choice.OfferAmount
		// sub_13902 的 `DX > [BP+0Ah]` 會把回傳 AL 設為 3；
		// 外層 handler 以 `cmp al, 2` 視為拒絕，不套用外交收尾。
		if amount > choice.InitialAmount {
			w.AdjustTrust(diplomacyOverOfferTrustPenalty)
			return false
		}
		if amount > 0 {
			response = 1
		}
	}

	switch choice.Kind {
	case DiplomacyCooperation:
		return w.finishQueuedCooperation(choice.Source, choice.Invader, choice.Target, response, amount)
	case DiplomacyCeasefire:
		return w.finishQueuedCeasefire(choice.Source, choice.Target, response, amount)
	default:
		return false
	}
}
