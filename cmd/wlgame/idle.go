package main

import "github.com/wicanr2/wolong_cht/internal/rules/speed"

// idleClockGate 是自然世界迴圈的輸入閘門。
//
// 松崗繁中版 KI.EXE 有兩層：`sub_11F7F` 比較本次與前次滑鼠座標，**只有座標
// 變了**才把倒數 `ds:98A5h` 重設成 12；`sub_11CD0` 只要看到那個倒數非零就
// 提早 retn，跳過 `sub_13EFD`（據點的每 tick 處理）。於是游標一直動就一直
// 被重設，**移動期間世界完全不前進**；停下來之後還要等 `ds:0D2Dh` 累積到
// `0A0h` ＝ 160 個回呼才恢復。
//
// ⭐ **這是即時制底下給玩家的反應時間**（使用者裁定 2026-09-02）：沒有它，
// 光是把游標移到指令列、再移到目標據點，時間就已經走掉，一條命令下不完。
//
// 規格：docs/spec/112-cursor-idle-resume-delay.md。
type idleClockGate struct {
	seen bool
	x    int
	y    int
	// acc 是中斷單位的累加器，與速度層共用同一個時間基準
	// （speed.UnitsPerFrame／speed.Scale），不要換成寫死的幀數。
	acc int
}

// idleResumeCallbacks 是游標停下之後要等幾個時鐘回呼才恢復。
// `sub_11CD0` 的 `cmp byte ptr ds:0D2Dh, 0A0h`（docs/re/61 §5.1）。
const idleResumeCallbacks = 0xA0

// idleResumeUnits 是同一個門檻換成中斷單位：160 × 600 ＝ 96,000，
// 每幀累加 2,913 ⇒ 33 幀，60 fps 下 0.550 秒（原版 0.5493 秒）。
const idleResumeUnits = idleResumeCallbacks * speed.Scale

// Allows 記錄這一 frame 的游標位置，並回報自然世界迴圈能否執行。
// 第一次呼叫一定回傳 false，避免未初始化座標被誤判成 idle。
//
// 兩種輸入的效果不同，而且**這個差別有出處**：
//
//   - scrolling ＝ 游標移動類（方向鍵捲鏡頭）。原版把游標推到畫面邊緣捲地圖
//     走的是 `sub_11F7F` 的同一支（`add ds:9882h, cx` 之後落到 `loc_11FD0`），
//     所以**重新等滿**。
//   - command ＝ 滑鼠鍵與 remake 自己加的鍵盤捷徑（ESC／1–4／＋−）。
//     原版按鍵不寫 `ds:98A5h`，鍵盤捷徑更是 remake 才有的
//     （`docs/re/47` §3）——所以**只擋這一 frame，不重新計時**。
func (g *idleClockGate) Allows(x, y int, scrolling, command bool) bool {
	moved := !g.seen || g.x != x || g.y != y
	g.seen = true
	g.x = x
	g.y = y
	if moved || scrolling {
		g.acc = 0
		return false
	}
	if g.acc < idleResumeUnits {
		g.acc += speed.UnitsPerFrame
	}
	return g.acc >= idleResumeUnits && !command
}

// Pause 重新開始等待，效果與游標移動一次相同。
// 原版在擦掉訊息框之後也設一次倒數（`sub_18810` 的 `mov cs:byte_198A5, 8`），
// 讓玩家看完訊息、手還沒回到滑鼠時世界不會先跑掉。
func (g *idleClockGate) Pause() { g.acc = 0 }
