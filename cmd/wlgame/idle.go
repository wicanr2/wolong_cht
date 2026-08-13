package main

// idleClockGate 是自然世界迴圈的輸入閘門。
//
// 松崗繁中版 KI.EXE 的 sub_11F7F 比較本次與前次滑鼠座標；相同時才在
// byte_198A3 設 idle bit，主迴圈隨後才能進 sub_11CD0。remake 不重用該
// 暫存 byte，而以這個小型、可單測的狀態保留同一個可觀測契約：必須先
// 觀測到游標穩定，且當前 frame 沒有命令／按鈕輸入，才允許世界更新。
type idleClockGate struct {
	seen bool
	x    int
	y    int
}

// Allows 記錄這一 frame 的游標位置，並回報自然世界迴圈能否執行。
// 第一次呼叫一定回傳 false，避免未初始化座標被誤判成 idle。
func (g *idleClockGate) Allows(x, y int, inputActive bool) bool {
	moved := !g.seen || g.x != x || g.y != y
	g.seen = true
	g.x = x
	g.y = y
	return !moved && !inputActive
}
