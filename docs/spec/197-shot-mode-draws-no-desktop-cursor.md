# 197 — 截圖對拍不畫桌面模式的游標

**狀態：CONFORMED。** `-shot` 是**對拍用的路徑**，而 2026-09-08 的
`ef666c5` 讓桌面捕捉模式（`desktopPointer.active`）在畫面上常駐兩個東西：

- `drawMouseCursor`：14×14 的紅箭頭（跟著 `cursorPosition()`）
- `drawMapPickCursor`：游標所在那一格的 **16×16 白黑框**

原版同狀態的畫面沒有它們，於是每一張對拍截圖都多出一塊差異。

- 日期：2026-09-11
- 出處：這是 **remake 差異的修正**，不是原版規則——原版的游標旗標見
  [`154`](154-mouse-cursor.md)
- 推論等級：confirmed（二分到 `ef666c5`，改掉之後主畫面回到五區全 PASS）
- 實測：[`../playtest/120`](../playtest/120-screen-parity-retest-20260911.md)
- remake 實作：`cmd/wlgame/cursor.go` 的 `drawMouseCursor`、
  `desktopMapCursorVisible`

## 1. 症狀

| 組 | 修之前 | 修之後 |
|---|---:|---:|
| 主畫面 `map` | 101 px（0.07%）| **0 px** ⇒ 五區全 PASS |
| 野戰 `field` | 190 px（0.11%）| **95 px** ＝ [`58`](../playtest/58-parity-retest-20260902.md) 的值（原版錄影裡的滑鼠游標）|

主畫面那 101 px 落在 x 320–335、y 208–223——**正好一個 16×16 的格子**，
顏色只有白 `(243,243,243)` 與黑，與 `drawMapPickCursor` 畫的框一致。

## 2. 修法

兩處都加「截圖模式就不畫」：

```go
// drawMouseCursor
if at == nil && g.shotPath != "" {
	return
}

// desktopMapCursorVisible
if g == nil || g.shotPath != "" || !desktopPointer.active || … {
	return false
}
```

⚠ **`-cursor X,Y` 不受影響**：那一條走 `g.cursorAt != nil`，是對拍
**明確要求**畫游標的路徑（[`154`](154-mouse-cursor.md)）。

## 3. 為什麼這種回歸會活這麼久

`ef666c5` 是一個 111 個檔案的大 commit（「完成桌面操作修正、呂布 AI
根因與 v.1.0.18 發行準備」）。桌面操作是**遊玩端**的改良，而它順手
改掉的是**對拍端**的行為——兩者共用同一個 `Draw`。

⇒ 之後每一輪都在改規則層、用同狀態對拍驗規則，**沒有一輪回頭截圖**，
所以這塊 101 px 一路活到打包前的重量。

⭐ 缺的是**一支可重跑的畫面對拍閘**：規則層有
`tools/parity_ck.sh`（六張表逐 byte），畫面層沒有對應的東西。

<!-- 缺口：無 -->
