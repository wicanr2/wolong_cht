# 110 — 指令列的「軍團」是兩項彈出選單

**狀態：CONFORMED。原版點「軍團」跳的是**位置確認／行軍指示**兩項選單；
remake 先前直接開軍團一覽，行軍只掛在鍵盤的 `M`。**

- 日期：2026-09-01
- 出處：[`../re/22`](../re/22-strategy-command-tree.md) §3.3
  （`sub_1628F`、`sub_193E9(ax=2, cx=4Fh, dx=40Ch)`）；
  實機 [`../playtest/54`](../playtest/54-menu-second-row-tap.md)、
  [`../playtest/56`](../playtest/56-lubu-flow-parity.md) §4.6
- 推論等級：**confirmed**——兩個分支的呼叫鏈都在反組譯裡，選單也拍到了

## 1. 原版做什麼

```asm
mov ax, 2 / mov cx, 4Fh / mov dx, 40Ch / call sub_193E9   ; 選單 TALK #79，2 項
al = 0 → 位置確認：sub_18853(cx=16h) → sub_1716D → sub_12151
al = 1 → 行軍指示：sub_18853(cx=2)   → sub_1716D → sub_17F90
```

| 項 | 狀態列 | 之後 |
|---|---|---|
| 位置確認 | TALK #22「將游標移動至軍團的現在位置。」 | 選一支軍團 → **鏡頭移過去** |
| 行軍指示 | TALK #2「請選擇進行行軍指示之軍團。」 | 選一支軍團 → 目的地一覽 → 三選一 |

兩項用的是**同一張軍團一覽**（`sub_1716D`，回傳選中的軍團記錄位址，
CF=1 為取消）。

### 1.1 位置確認的鏡頭

```asm
mov dx, [bx+10h]      ; 軍團記錄 +0x10 = X
mov bx, [bx+12h]      ; 軍團記錄 +0x12 = Y
mov ax, 14h / mov cx, 0Ch
call sub_12151
```

`sub_12151(ax, cx)` 把鏡頭放到 `(dx − ax, bx − cx)`，
立即值 `0x14`／`0x0C` ＝ **(20, 12)**，與開局鏡頭是首都 −(20,12)
（[`52`](52-main-screen-camera-and-banner-date.md)）是同一組常數。

## 2. 演算法

```
軍團指令：
  選單（TALK #79，兩項）
    第 0 項 位置確認：軍團一覽 → 選中 c → 鏡頭 = (c.X − 20, c.Y − 12)，夾邊界
    第 1 項 行軍指示：軍團一覽 → 選中 c → 目的地一覽 → 三選一
```

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 呈現層 | `cmd/wlgame/corps.go` 的 `openCorpsCommandMenu`（新）、`beginLocateCorps`（新）；`naturalCommandActions` 的第 5 格改指向選單 |
| 差異 | **鍵盤 `C` 開選單、`M` 直接跳行軍指示**。`M` 是 remake 既有的捷徑（原版沒有鍵盤），保留 |

⭐ 鏡頭移動沿用既有的寫法（`c.X-centreCol, c.Y-centreRow` ＋ `clampCam()`），
`centreCol`／`centreRow` 就是 20／12，與 `sub_12151` 的立即值同一組。

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestCorpsCommandMenuRows`、`TestLocateCorpsMovesCamera`（`cmd/wlgame`）|
| 對原版 | [`../playtest/56`](../playtest/56-lubu-flow-parity.md) §4.6 的原版選單影格 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| 選單本身的逐像素對拍 | 原版選單拍到了（`parity-tap5/menu.png`），**但沒有做逐像素比對**——remake 的彈出選單走 `drawLegacyChoiceBox`，位置與底色未量 |
| 「據點」指令同樣是兩項選單 | `../re/22` §3.4：TALK #82「首都確認／據點一覽」。**這一輪沒動**，remake 的「據點」目前直接開一覽 ＋ 情報卡 |
