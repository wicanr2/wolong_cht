package state

import (
	"errors"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
)

// ErrIncompleteTick 表示世界仍停在戰術戰鬥中斷的更新，不能寫成存檔。
var ErrIncompleteTick = errors.New("state: incomplete world update")

// start 對應原版批首D18；next是暫態配額進度，不寫進原版區塊。
// tail為尚未跑軍費／倒數的軍團，−1表示沒有未完尾段。
type corpsBatch struct {
	start int
	next  int
	tail  int
	hour  int
}

type tickContinuation struct {
	corps             corpsBatch
	includeMapObjects bool
}

// HasIncompleteTick 回報是否仍有被戰術戰鬥中斷的世界更新（spec272）。
// 地圖命令與存檔需等ResumeTick完成後才開放。
func (w *World) HasIncompleteTick() bool { return w != nil && w.tickResume != nil }

// ResumeTick 續完既有更新，事件交回同一條Event管道。
// 沒有未完更新時不會開始新拍；續跑也可能遇到同批的下一場戰鬥。
func (w *World) ResumeTick(random economy.Rand) Event {
	if !w.HasIncompleteTick() {
		return Event{HourFaction: -1}
	}
	return w.tick(random, w.tickResume.includeMapObjects)
}
