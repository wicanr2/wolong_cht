# 33 — 底列六格是選部隊，不是第二套命令列

**狀態：CONFORMED。** 兩張順序表、選取位元圖、下令時的取用方式、
「沒選 ＝ 全隊」與六張命令圖示都有機器碼出處，已實作並有單測。

- 日期：2026-08-16
- 出處：[`../re/60-tactical-sidebar.md`](../re/60-tactical-sidebar.md) §6.2、§10
  （`sub_1C7F4` `0001C7F4`、`sub_1C6BF` `0001C6BF`、`sub_1C74C` `0001C74C`、
  handler `0001C1B9`、`sub_1A8CC` `0001A8CC`）
- 推論等級：confirmed

## 1. 原版做什麼

畫面底列（y 368–400）是**玩家六個編成位置**各一格，不是命令按鈕。
點一格切換該隊的選取狀態；接著在側欄的指令面板下令，
命令只送給**被選中的隊**，一格都沒選就送給全隊。

說明書 4.3 寫的就是這件事：「點該部隊情報框，出現**黃框**表示已選」／
「**沒有任何黃框時，命令對全軍生效**」。

### 1.1 ⭐ 兩張互為反排列的順序表

| 表 | 內容 | 值 |
|---|---|---|
| `cs:0xD2E4` | 螢幕格 → 隊編號 | `2 4 0 1 5 3` |
| `cs:0xD2EA` | 隊編號 → 螢幕 X | `160 240 0 400 80 320` |

六個編成位置是 0 主將／1 前鋒／2 左翼／3 右翼／4 左備／5 右備，
所以畫面由左到右是：

```
 左翼   左備   主將   前鋒   右備   右翼
   0     80    160    240    320    400
```

⭐ **是空間排列**：左翼與左備在左、主將與前鋒在中、右備與右翼在右。
兩張表互為反排列——`cs:0xD2EA[cs:0xD2E4[i]] = i × 80` 對六格全部成立。

`sub_1C7F4` 用 `cs:0xD2E4` 決定每格的熱區碼（`表值 + 0x15`）；
`sub_1C6BF`（選取框）與 `sub_1C74C`（兵條）用 `cs:0xD2EA` 決定 X。

### 1.2 一格裡有三樣東西

| 東西 | 位置 | 來源 |
|---|---|---|
| 位置名 glyph | 格內 (4, 6)，24×16 | `ICONGRF` 段 1 `0x3900 + 螢幕格 × 0xC0` |
| **目前命令的圖示** | 格內 (54, 6)，24×16 | `ICONGRF` 段 3 的 `命令碼 × 0xC0`（§1.5）|
| **待機兵條** | (格 X ＋ 2, 396)，長上限 `0x4C` ＝ 76，2 px 高，色 12 | `word_1D30A:+0x09 + 4×隊` |

命令圖示由 `sub_1A8CC` 更新，而它開頭就是：

```asm
0001A8CC  cmp si, 600h / jnb → retn      ; ★ 只有側 0（玩家）有這一列
0001A8D4  mov bx, si / mov ah, bh        ; 隊編號 ＝ 單位記錄位移 >> 8
0001A8D8  call sub_1C673
```

六個命令 handler 各自傳自己的碼進來（`sub_1A92E` 傳 0、`sub_1A953` 傳 1、
`sub_1A96D` 傳 2、`sub_1A988` 傳 3、`sub_1A99C` 傳 4、`sub_1A9D0` 的 5 直接跳過），
而且前面都有一道 `cmp al, ah / jz`——**命令沒變就不重畫**。

⭐ `sub_1A96D`（命令 2）在 `byte_1D34B == 0`（攻城）時順手呼叫 `sub_1B7CB`
把門全開，正好是說明書 4.2 的「突擊時守方會開門」。
**這獨立驗證了命令 2 ＝ 突擊**（另一條路是 TALK 台詞，
[`../re/60`](../re/60-tactical-sidebar.md) §6.1）。

### 1.3 選取位元圖與下令

```asm
0001C27D  mov ah, 0Ch / xor cs:byte_1D310, 1     ; 熱區 0x15 → 位元 0
0001C285  test cs:byte_1D310, 1 / jnz .1
0001C28D  xor ah, ah                              ; 沒選中 → 色 0
0001C28F  mov al, 0 / call sub_1C6BF               ; 重畫第 0 隊那一格的框
```

`0x15`–`0x1A` 六個 handler 形狀相同，切位元 0–5、傳 `al` ＝ 隊編號。
下令那一支（`0001C1B9`）：

```asm
0001C1CE  xchg ah, cs:byte_1D310      ; ★ 取出並清空
0001C1D3  and ah, ah / jnz .1
0001C1D7  mov ah, 0FFh                 ; ★ 沒選 ＝ 全隊
.1: cx = 6 / si = 0x1B
    每輪 shr ah,1，進位就寫 [si] = 命令碼；si += 0x100
0001C1F4  call sub_1C6AE               ; 六格的框全部重畫（都變成未選）
```

**下完令選取就清空**——不是保留。

### 1.4 選取框的幾何（`sub_1C6BF`）

兩個同心矩形，`ah` 選中時 `0x0C`、取消時 `0`：

```
外框 (X + 2, 372) – (X + 77, 392)
內框 (X + 3, 373) – (X + 76, 391)
```

### 1.5 ⭐ 六張命令圖示就在段 3 的最前面，而且圖畫的是什麼可以逐張核對

`sub_1C673` 從 `word_10D48` 段取 `命令碼 × 0xC0`。`sub_1006B` 把 `ICONGRF`
的檔案位移 `0x9700`（＝段 3 的起點）讀進 `word_10D48`，`sub_100DF` 只再往後
配置 `word_10D4A = word_10D48 + 0x6C` 段等指標（[`../re/48`](../re/48-window-display-list.md) §6.1），
**`word_10D48` 本身就是段 3 的位移 0**。所以段內位移就是 `碼 × 0xC0`。

解出來的六張與命令碼**逐張相符**：

| 碼 | 段 3 位移 | 圖 | 命令 |
|---:|---|---|---|
| 0 | `0x0000` | 紅色的陣地 | 陣形 |
| 1 | `0x00C0` | 長槍 | 攻擊 |
| 2 | `0x0180` | 紅黃軍旗 | 突擊 |
| 3 | `0x0240` | 磚牆 | 城壁 |
| 4 | `0x0300` | 盾牌 | 守陣 |
| 5 | `0x03C0` | 白旗 | 退卻 |

⭐ **這是內容檢查，不是推論**：碼 3 過去只能由「六個命令扣掉已知的四個」
反推，磚牆那張圖把它獨立驗了一次。六張的尾端 `0x480` 落在外框 motif
（`0x6C0`）之前，中間第 7–9 張是兵種圖示的橘色版（馬／弓／步），
與 `0x1BA0` 的紅版、`0x1EA0` 的綠版同一批。

## 2. 演算法

```
點螢幕第 i 格  →  隊 = cs:D2E4[i]  →  選取位元圖 ^= 1 << 隊

下令(命令碼 c):
    if c == 城壁 and 這張圖沒有城:            # byte_1AB4F == 0
        跳 TALK 582「這哪裡有城壁啊！！」，不下令
        return
    位元圖 = 取出並清空(選取位元圖)
    if 位元圖 == 0: 位元圖 = 0xFF             # 全隊
    for 隊 in 0..5:
        if 位元圖 & (1 << 隊): 那一隊的每個兵 +0x1B ← c
    跳 TALK 0x1B1 + c
```

⚠ **城壁令對玩家是「拒絕」，對腳本是「降級成攻擊」**——
兩條路不同（[`../re/60`](../re/60-tactical-sidebar.md) §6.1）。
`Battle.Order` 現有的 `ScaleWal → Attack` 是腳本那一條，**保留**；
玩家那一條要另外走。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 選取位元圖 | `internal/rules/tactical`：`Side.Selected`、`Battle.ToggleSquadSelection`、`Battle.TakeSquadSelection` |
| 玩家下令 | `Battle.OrderSelected`（城壁在非攻城戰回 `false` ＝ 拒絕）|
| 兩張順序表 | `cmd/wlgame/battlelayout.go`：`battleBottomSlotSquad`、`battleSquadSlotX` |
| 底列繪製 | `cmd/wlgame/battle.go` `drawBattleKeys`：位置名 glyph ＋ **命令圖示** ＋ 待機兵條 ＋ 選取框 |
| 命令圖示 | `library.DOSVOrderIcon(碼, 季)`；圖示畫的是**該隊隊長的 `Cmd`**（每隊第一個兵）|
| 輸入 | 底列點擊 → 切換選取；側欄面板點擊 → `battleSideCommandRowCode[列]` |
| 差異 | **鍵盤 1–6 直接下命令**是 remake 加的（原版這一層只有滑鼠）；它一樣走 `OrderSelected`，所以選取的行為一致 |

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestDOSVOrderIcons`（`internal/assets/gfx`）：六張互不相同、都有內容、尾端不撞到外框圖塊 |
| 單元測試 | `TestSquadSelectionMatchesRawBitfield`、`TestOrderSelectedRejectsScaleWallOffSiege`（`internal/rules/tactical`）；`TestBottomSlotSquadTablesAreInverse`、`TestBattleSideCommandClickUsesRowCode`（`cmd/wlgame`）|
| 對原版 | 兩張順序表的值直接取自 `cs:0xD2E4`／`cs:0xD2EA`；畫面證據見 `docs/images/wlgame-tactical-squad-select.png` |

## 5. 未解

| 項目 | 現況 |
|---|---|
| 待機兵條的欄位語意 | `word_1D30A:+0x09 + 4k` 在 [`../re/11`](../re/11-tactical-battle.md) §3.9 記成「第 k 隊的待機兵數」；條的上限 76 遠小於一隊 100 兵，所以開局會頂在上限 |
