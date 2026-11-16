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

// BeginDiplomacyChoice 直接把一份外交提案掛成「擋住世界的決定」。
// 保留給測試與跨平台輔助控制（手機版的金額輸入器要在沒有事件的局面上驗），
// 遊戲流程一律走 beginDiplomacy。
func (w *World) BeginDiplomacyChoice(c DiplomacyChoice) {
	w.diplomacy = &c
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
		choice.InitialResponse, choice.OfferAmount, ok = w.cooperationTerms(choice.Source, choice.Invader, choice.Target)
	case DiplomacyCeasefire:
		choice.InitialResponse, choice.OfferAmount, ok = w.ceasefireTerms(choice.Source, choice.Target)
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
		// sub_13902先以通報者顯示#42，再進場景中的君主提問。
		ev.TalkNotices = append(ev.TalkNotices, TalkNotice{
			Index: 0x2A, City: -1, Faction: -1, General: -1, Amount: -1,
		})
	}
	return true
}

// ResolveDiplomacy保留既有bool入口，與呈現層共用同一份結算。
func (w *World) ResolveDiplomacy(option DiplomacyOption) bool {
	_, applied := w.ResolveDiplomacyWithResult(option)
	return applied
}

// ResolveDiplomacyWithResult按原版13902保留初始回應及一次採納骰。
// 回傳的實際條件供UI報告，不能以玩家建議推定君主最後決定。
func (w *World) ResolveDiplomacyWithResult(option DiplomacyOption) (DiplomacyResult, bool) {
	result, settle := w.PrepareDiplomacy(option)
	if settle == nil {
		return result, false
	}
	return settle()
}

// PrepareDiplomacy決定回覆，但把資金、Trust等寫回留給中央結果確認後。
// 續行只能執行一次；不重擲採納骰，不進存檔。同步caller沿上方入口。
func (w *World) PrepareDiplomacy(option DiplomacyOption) (DiplomacyResult, func() (DiplomacyResult, bool)) {
	result := DiplomacyResult{Response: 2, AdviceOption: option}
	choice := w.diplomacy
	if choice == nil || option > DiplomacyReject {
		return result, nil
	}
	w.diplomacy = nil
	if option == DiplomacyOfferFunds && choice.OfferAmount == 0 {
		result.AdviceOption = DiplomacyAcceptFree
	}
	// 包含拒絕也必須取一次，且不能重新計算外交代表而多取亂數。
	roll := 0
	if w.rng != nil {
		roll = w.rng.Next() & 0xFF
	}
	if choice.InitialAmount < 0 || choice.OfferAmount < 0 {
		return result, nil
	}
	result.Response, result.Amount = choice.InitialResponse, choice.InitialAmount
	result.AdviceAdopted = roll <= clampU8(w.Trust)
	if result.AdviceAdopted {
		result.Response = int(result.AdviceOption)
		if option == DiplomacyOfferFunds {
			result.Amount = choice.OfferAmount
		}
		if option == DiplomacyOfferFunds && choice.OfferAmount > choice.InitialAmount {
			result.Response = 3
		}
	}
	decided := *choice
	settled := false
	return result, func() (DiplomacyResult, bool) {
		if settled {
			return result, false
		}
		settled = true
		return w.applyDiplomacyResult(decided, result)
	}
}

func (w *World) applyDiplomacyResult(choice DiplomacyChoice, result DiplomacyResult) (DiplomacyResult, bool) {
	if result.Response == 3 {
		w.AdjustTrust(diplomacyOverOfferTrustPenalty)
	}
	if result.Response >= 2 {
		return result, false
	}
	declarationNotice := choice.Kind == DiplomacyCooperation && choice.Source == w.Player &&
		choice.Source >= 0 && choice.Source < numFactions && choice.Invader >= 0 && choice.Invader < numFactions &&
		!w.Friendship[choice.Source][choice.Invader].AtWar()
	applied := false
	switch choice.Kind {
	case DiplomacyCooperation:
		applied, result.ReleasedGenerals = w.finishQueuedCooperationWithReturnees(choice.Source, choice.Invader, choice.Target, result.Response, result.Amount)
	case DiplomacyCeasefire:
		applied, result.ReleasedGenerals = w.finishQueuedCeasefireWithReturnees(choice.Source, choice.Target, result.Response, result.Amount)
	}
	if applied && declarationNotice {
		result.TalkNotices = append(result.TalkNotices, TalkNotice{
			Index: 0x1A0, City: -1, Faction: choice.Invader,
			General: w.Factions[choice.Source].Lord, Amount: -1, SpeakerPortrait: true,
		})
	}
	return result, applied
}
