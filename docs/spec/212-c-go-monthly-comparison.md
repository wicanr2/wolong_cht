# 212：原版／C／Go 完整月結比較

**狀態：CONFORMED。四劇本 36 固定向量的完整原版區塊及最終 RNG 相同，限局部月結規則。**

- 日期：2026-10-08
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 原版入口：IDA Pro 9.4 database linear `sub_15358`，完整 C 呼叫鏈見 [`re/97`](../re/97-c-monthly-politics-restoration.md)
- Go 入口：`internal/state/state.go` 的月結規則區段，原本嵌在 `tick`，每時更新在其後

## 1. 研究契約

從同一原版劇本區塊載入原版/C 與 Go，以同一玩家、設定、queue、raw RNG 初值執行一個月結。
原版每個規則 callee 保留真實指令；C 用已審查來源；只有通知／音效／重畫為 RET fixture。
Go 純抽出原有月結區段，正常 `tick` 仍以同樣順序呼叫，不藉抽出改規則。
研究適配只在 `matching` build tag 下開放，不能成為正式遊戲 shortcut。

先比較載入／匯出前置基準，再比月結後的 globals、勢力、據點、武將、corps、queue 與完整 RNG。
每個差異逐欄定位，初始差異與規則差異分開。近似 event 10 以現有明示開關關閉，其他正式規則保留。
不能透過移除規則、放寬 mask 或丟棄 queue／RNG 差異取得通過。

## 2. 實作閘門

本規格允許可丟棄研究 probe 與不改行為的月結區段抽取。發現會改變正式規則的差異時，
先以原始 bytes／原版動態向量審查為 READY 的個別契約，再進正式規則修正。
正常 `tick` 的既有月結次序冷測必須通過；局部直接入口不宣稱 GUI／玩家垂直鏈完成。

比較、載入基準、source identity 及六個逐 byte／RNG 拒絕 selftest 均通過，見
[`re/98`](../re/98-c-go-monthly-comparison.md)。

## 3. 未解範圍

| 項目 | 邊界 |
|---|---|
| 其他後期狀態與全部政治條件 | 未驗證，36 固定原版向量不外推任意世界 |
| 正常 UI、音效、存檔及長程流程 | 未驗證，純規則比較不取代正常玩家驗收 |
| C 機器碼匹配 | 未驗證，整檔逐位元組匹配仍為組語基準 |
