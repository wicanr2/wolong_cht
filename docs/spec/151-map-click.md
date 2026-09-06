# 151 — 在大地圖上點一格：據點開情報卡、軍團開情報面板

**狀態：CONFORMED。** 原版的大地圖左鍵是一條完整的分派（`sub_11E46`）：
那一格是據點就開情報卡、有軍團就開軍團面板、兩者都有就先跳一張兩項的選單。
remake 的大地圖點擊**什麼都不做**。

- 日期：2026-09-06
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）
  - `sub_11E46`（`00011E46`）：分派本體
  - `sub_11F0E`（`00011F0E`）：兩項選單（TALK **#80**「　據　　點　」「　軍　團　」）
  - `sub_171D3`（`000171D3`）＋ 建表 callback（線性 `0x17217`）：那一格上的軍團一覽
  - `sub_17E1F`（據點情報卡）、`sub_17F90`（軍團情報面板，[`149`](149-march-target-map-picker.md) §1.3）
- 推論等級：**confirmed**（分派逐行讀完；`0x17217` 由原始 bytes 交叉解碼）
- 驗證：[`../playtest/96`](../playtest/96-map-click.md)
- 相關：[`23`](23-city-info-window.md)、[`24`](24-corps-info-window.md)、
  [`149`](149-march-target-map-picker.md)（同一套格座標與命中判定）

## 1. 原版做什麼

```
sub_11E46（大地圖左鍵）:
    格 ← 游標（word_1989A 行／word_1989C 列）
    圖塊 ∈ 0CBh–0D3h 且據點表的 X／Y 完全相等 → [bp+0] = 據點記錄
    佔用圖（cs:word_19872 段）那一格的值 → [bp+2]
    有據點 且 佔用 ≠ 0 → sub_11F0E（兩項選單）
        右鍵取消 → 收工
        選 1（軍團）→ 走軍團分支
        選 0（據點）→ 走據點分支
    只有據點 → sub_17E1F（情報卡）
    只有軍團（佔用 ≥ 1）→ 迴圈：
        sub_171D3   ; 那一格上的軍團一覽；右鍵 → 收工
        si = 軍團記錄
        sub_17F90   ; 軍團情報面板（自己的接著進行軍指示）
        回迴圈開頭   ; ★ 可以連著看好幾支
```

⭐ **命中判定與行軍目標選點是同一套**（[`149`](149-march-target-map-picker.md) §1）：
圖塊 `0CBh`–`0D3h` ＋ 據點表 X／Y 完全相等，**一座城只有登記的那一格按得到**。

### 1.1 那一格上的軍團怎麼列（線性 `0x17217`）

```asm
    si = 2240h                     ; 軍團表起點
    xor di, di / xor ah, ah
    al = cs:byte_10CFF             ; ⚠ 讀了玩家勢力，但**這一支從頭到尾沒比對它**
.loop:
    cmp byte [si], 80h   / jb  next   ; ① 軍團存在
    cmp <游標列>, [si+10h] / jnz next  ; ② +0x10 ＝ Y
    cmp <游標行>, [si+12h] / jnz next  ; ③ +0x12 ＝ X
    [bp+di] = si / ah++ / di += 2
next:
    si += 40h / cmp si, 41C0h / jnz .loop
```

⭐ **不分敵我**——`al` 那一行是留下來沒用到的。所以點得到別人的軍團，
而「別人的」由 `sub_17F90` 用狀態列 #4 處理
（[`149`](149-march-target-map-picker.md) §1.3）。

軍團表：起點 `0x2240`、stride `0x40`、終點 `0x41C0` ⇒ **126 支**。

### 1.2 兩項選單開在游標旁邊

```asm
sub_11F0E:
    ax = cs:word_19898 / al > 16h → al = 16h    ; Y 夾住
    dx = cs:word_19896 / dl > 23h → dl = 23h    ; X 夾住
    dh = al / ax = 102h / cx = 50h / call sub_193E9
```

與行軍三選一同一支 `sub_193E9`，字串同樣是六個全形字一列（框寬 112），
**只有夾制上限不同**：

| | X 上限 | Y 上限 | 夾住之後的右／下緣 |
|---|---:|---:|---|
| `sub_1804E`（行軍三選一，2 項）| `21h` | `18h − 項數` ＝ `16h` | 528＋112 ＝ **640**／352＋48 ＝ **400** |
| `sub_11F0E`（這一張，2 項）| **`23h`** | `16h`（寫死）| 560＋112 ＝ **672**／352＋48 ＝ **400** |

⚠ **直的那道兩邊都剛好貼齊畫面，橫的那道只有 `21h` 貼齊。**
`23h` 讓框超出右緣 32 px。兩個立即值都是機器碼直接讀到的，
所以這是**原版就這樣**，不是誰抄錯——但「夾制是為了不讓框出畫面」
這個說法對這一支不成立（§4）。**照抄常數，不要替它補一個它沒寫的規則。**

## 2. remake 要怎麼改

| 項目 | 位置 |
|---|---|
| 分派 | `cmd/wlgame/mapclick.go` 的 `updateMapClick()`：排在熱區與各視窗之後，與地圖選點同一個位置 |
| 命中 | 沿用 `cityAtTile`（[`149`](149-march-target-map-picker.md)）＋ 新的 `corpsAtTile` |
| 據點 | 既有的 `openCityInfo` |
| 軍團一覽 | `openCorpsListWith` 帶「那一格上的」列，選完走 `enterCorpsFromMap` |
| 進入軍團 | `enterCorpsFromMap`：自己的 → `showCorpsPanel` ＋ `pickDestination`（[`149`](149-march-target-map-picker.md) §1.3）；別人的 → `openCorpsInfo`（狀態列 #4）|
| 兩項選單 | `mapChoiceState` ＋ TALK #80，沿用 `drawLegacyChoiceBox`／`talkChoiceClick`；夾制照 §1.2 |
| 差異 | 原版的軍團分支是**迴圈**（看完回一覽表繼續看）。remake 照抄：`enterCorpsFromMap` 回 false |

## 3. 驗證

| 方式 | 內容 |
|---|---|
| 對原版 ✅ | [`../playtest/96`](../playtest/96-map-click.md)：兩項選單 112×48 **逐像素 0 px**；整張畫面比不了（鏡頭是格級的）|
| 單元測試 | `TestCorpsAtTileListsBothSides`：不分敵我、只認完全相等的格 |
| 單元測試 | `TestMapChoiceAnchorClamps`：§1.2 的兩個夾制上限 |
| 單元測試 | `TestMapClickDispatch`：只有據點／只有軍團／兩者都有的三條分支 |

## 4. 未解

| 項目 | 現況 |
|---|---|
| `sub_11F0E` 的 X 夾制 `23h` | 夾住之後框的右緣落在 672，超出畫面 32 px（§1.2）。要嘛 `sub_193E9` 對超出的部分另有處理、要嘛這一張的框比 112 窄——**兩個都還沒驗**。目前照抄常數 |
