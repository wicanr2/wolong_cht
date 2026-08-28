# 103 — 手機版外交提案的「提示金額」：數字鍵盤

**狀態：CONFORMED（2026-08-29 實作並有單測）。**

- 日期：2026-08-29
- 出處：原版三選一是 `sub_13902`（[`78`](78-amount-input-editor.md) §1 的呼叫端之一），
  「提示金額」走 `sub_17C6E` 的數值輸入器；規則層已照原版做成
  `state.EditDiplomacyOfferAmount`（逐位輸入、百、刪、最大、清除，上限 `0x7530`）。
  **鍵盤的版面是 remake 差異**（手機沒有原版那張 640×400 的數字視窗可放）。
- 推論等級：規則 confirmed（沿用 `78`）；版面不涉及原版事實。

## 1. 為什麼

手機版的外交提案先前只給「接受／拒絕」兩項，第二項「提示金額」因為沒有數值
輸入器而整個拿掉（`android-ux.md` §7）。原版三選一的第二項是玩家談判的主要工具
——少了它，停戰與協力都只能白給或拒絕。

## 2. 做法

| 項目 | 內容 |
|---|---|
| 三選一 | 用 `TALK.DAT` #283 的原文：「答應／提示金額／拒絕」，順序＝`state.DiplomacyOption` |
| 點「提示金額」 | 主區換成數字鍵盤：上方一列標題（誰提什麼）＋ 目前金額，下方 4×4 鍵 |
| 鍵 | `7 8 9 ⌫`／`4 5 6 百`／`1 2 3 最大`／`清除 0 取消 確定`——與原版數值視窗的動作一一對應（`78` §1：數字、百、刪一位、最大、清除、確定），「取消」回三選一 |
| 每一鍵 | 直接呼叫 `world.EditDiplomacyOfferAmount`，上限與清零規則都在規則層 |
| 確定 | `ResolveDiplomacy(DiplomacyOfferFunds)`，用當下的 `OfferAmount` |
| 返回鍵 | 先關鍵盤回三選一，再照原本的順序 |
| 熱區 | 每鍵 ≥ 48 dp（主區 960×420 分 4×4，每格 240×~90）|

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 鍵盤與熱區 | `internal/ui/phone/modal.go`：`amountPad`、`amountKeys`、`amountKeyRect`、`drawAmountPad` |
| 規則層入口 | `state.BeginDiplomacyChoice`（測試與輔助控制用，遊戲流程仍走事件）|
| 差異 | 版面是 remake 的；輸入動作與上限照原版 |

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestDiplomacyAmountKeypad`（`internal/ui/phone`）：三選一有三項、按鍵改到 `OfferAmount`、確定後決定消失且 `OfferAmount` 進入結算 |
| 截圖 | `tools/phone_shot.sh`（`WOLONG_DIPLOMACY=1` 停在鍵盤）|

## 5. 未解

| 項目 | 現況 |
|---|---|
| 撥款請求的「指定金額」 | 原版 `sub_17C6E` 的另一個呼叫端；手機版仍只給「照要求撥款／拒絕」，同一套鍵盤可以接，先不做 |
