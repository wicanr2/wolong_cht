package state

import "github.com/wicanr2/wolong_cht/internal/rules/economy"

const (
	maxFundingOffer = 0x7530 // sub_17C6E 的輸入上限 30,000
	fundingMinOffer = 0x1F4  // sub_139E8 對非零初始要求的下限 500
)

// PendingFunding 回傳事件 4／5 等待玩家處理的撥款視窗副本。
// 這個狀態只存在 runtime，不寫入 SINARIO.DAT／SAVE.DAT。
func (w *World) PendingFunding() *FundingChoice {
	if w.funding == nil {
		return nil
	}
	c := *w.funding
	return &c
}

// SetFundingAmount 是 sub_139E8 → sub_17C6E 的數值輸入狀態部分。
// 原版輸入允許玩家把數字改到 0；初始要求的 500 下限在 beginFunding
// 套用，不能在這裡把玩家的自訂值再偷偷拉回 500。
func (w *World) SetFundingAmount(amount int) bool {
	if w.funding == nil {
		return false
	}
	if amount < 0 {
		amount = 0
	}
	if amount > maxFundingOffer {
		amount = maxFundingOffer
	}
	w.funding.OfferAmount = amount
	return true
}

// EditFundingAmount 接入事件 4／5 共用的 sub_17C6E 數值核心。
// 上限是 `maxFundingOffer`，畫面只負責把按鍵轉成 AmountEdit。
func (w *World) EditFundingAmount(edit AmountEdit, digit int) bool {
	if w.funding == nil {
		return false
	}
	amount, ok := editAmount(w.funding.OfferAmount, maxFundingOffer, edit, digit)
	if !ok {
		return false
	}
	w.funding.OfferAmount = amount
	return true
}

// beginFunding 將原版 sub_132A9／sub_132E9 停在 sub_139E8 的玩家選擇前。
// 事件 4／5 的處理端會重新確認派駐關係；事件排入後若官員已被換走，
// 這筆事件直接作廢，不能把舊指標當成仍有效。
func (w *World) beginFunding(ev *Event, choice FundingChoice) bool {
	if ev == nil || w.funding != nil || w.diplomacy != nil {
		return false
	}
	if _, ok := w.fundingOfficer(&choice); !ok {
		return false
	}
	if choice.RequestedAmount < 0 {
		choice.RequestedAmount = 0
	}
	if choice.RequestedAmount != 0 && choice.RequestedAmount < fundingMinOffer {
		choice.RequestedAmount = fundingMinOffer
	}
	if choice.RequestedAmount > maxFundingOffer {
		choice.RequestedAmount = maxFundingOffer
	}
	choice.OfferAmount = choice.RequestedAmount
	w.funding = &choice
	ev.Funding = w.funding
	return true
}

func (w *World) fundingOfficer(choice *FundingChoice) (int, bool) {
	if choice == nil || choice.Officer < 0 || choice.Officer >= numGenerals {
		return 0, false
	}
	switch choice.Kind {
	case FundingGovernor:
		if choice.Subject < 0 || choice.Subject >= numCities ||
			w.Cities[choice.Subject].Governor != choice.Officer {
			return 0, false
		}
	case FundingDiplomat:
		if choice.Subject < 0 || choice.Subject >= numFactions ||
			w.Factions[choice.Subject].Diplomat != choice.Officer {
			return 0, false
		}
	default:
		return 0, false
	}
	return choice.Officer, true
}

// ResolveFunding 消費目前的玩家撥款選擇。
//
// sub_139E8 的已證實尾段是：拒絕完全沒有副作用；其他兩列把金額乘二後
// 取高位寫入官員 +0x1A（等於 amount／128 的 byte），再從玩家勢力扣掉
// 原始金額。結果／收尾 TALK 的索引映射在 cmd/wlgame 接入；PC-98 對話
// 排版仍是呈現層的獨立邊界。
func (w *World) ResolveFunding(option FundingOption) bool {
	choice := w.funding
	if choice == nil || option > FundingReject {
		return false
	}
	w.funding = nil
	if option == FundingReject {
		return false
	}
	if w.Player < 0 || w.Player >= numFactions {
		return false
	}
	officer, ok := w.fundingOfficer(choice)
	if !ok {
		return false
	}
	amount := choice.RequestedAmount
	if option == FundingSetAmount {
		amount = choice.OfferAmount
	}
	if amount < 0 {
		amount = 0
	}
	if amount > maxFundingOffer {
		amount = maxFundingOffer
	}
	// sub_139E8 的指定金額編輯在輸入 0 時回傳選項碼 2；尾端只
	// 顯示訊息，不寫官員經費也不扣玩家資金。超過原始要求則是
	// 另一個已證實的選項碼 3，仍會照輸入金額完成，不能套用外交
	// 事件的「超過即拒絕」規則。
	if option == FundingSetAmount && amount == 0 {
		return false
	}
	w.Generals[officer].Budget = clampU8(amount / 0x80)
	w.Factions[w.Player].Funds = economy.ClampFunds(w.Factions[w.Player].Funds - amount)
	return true
}
