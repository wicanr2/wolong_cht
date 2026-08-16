# 45 — 進言的畫面：插圖 ＋ 兩個框輪流講話 ＋ 五列選單

**狀態：CONFORMED。** 三個框的位置、誰的肖像、選單的矩形與列數都有機器碼出處。
⭐ **選單框是 (80, 176) 160×96、固定五列**——外交、撥款、說服三處
共用同一支 `sub_13B7E`，那組座標與尺寸全是寫死的。

- 日期：2026-08-17
- 原始 `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 分析時的 `KI.EXE.i64` SHA-256：`8a8fd7d528e0498000fd04282300588b637fbcb8aa48deb09242f1f41f532691`
- 工具：**IDAPython** `tools/ida_dump.py`、`tools/ida_range.py`
- 出處：[`../re/66`](../re/66-message-box-geometry.md) §5.1／§5.2、
  [`44`](44-advise-original-text.md) §3
- 推論等級：**confirmed**（座標與尺寸都是立即值，兩條無關的路徑各自算出同一個高）

## 1. 三個框

`sub_13830` 每一步都指名一支繪製常式，所以「誰在哪個框說話」不必推：

| 步 | 常式 | 框 | 講話者 | 肖像 |
|---|---|---|---|---|
| 插圖 | `sub_13D09(al=0)` | (64, 144) 288×176 | — | `IVENTGRF` 第 0 頁 |
| ① 君主開場 | `sub_13C99` | **上框** (0, 80, 256, 80) | 玩家的君主 | 武將 `+0x01` |
| ② 軍師進言 | `sub_13CDC` | **下框** (128, 288, 256, 80) | **一定是軍師**（＝玩家）| 勢力 `+0x02` → 武將 `+0x01` |
| ③ 君主回答 | `sub_13C99` | 上框 | 同 ① | 同 ① |

說服迴圈把 ②③ 重複下去：`sub_13B5A` 每一輪先 `sub_13CDC`（軍師說出挑的理由），
再由 `sub_13BA9` 收尾 `sub_13C99`（君主的反應）。
**兩個框各自只顯示「最新的那一句」**，不是一份會往下捲的對話紀錄。

三個講話框的寬高完全一樣（`sub_1075B` 把 `cx = 0510h` 寫死），只有位置不同。

## 2. 選單框：(80, 176) 160×96，五列

`sub_13B7E` 是三個地方共用的選單：外交三選一（`sub_13902`）、
撥款（`sub_139E8`）、說服五選一（`sub_13B5A`）。

```asm
00013B7E  mov dx, 50h / mov bx, 0B0h / mov cx, 600Ah
00013B87  call sub_19796        ; 先把這塊畫面存起來
00013B8C  mov al, bl / xor ah, ah / mov cx, [bp+0]
00013B93  mov dl, 5 / mov dh, 0Bh
00013B97  call sub_193E9        ; 選單本體，回傳選到第幾列
00013B9A  jb short loc_13B8C    ; 取消就重來
00013B9C  … call sub_197C3      ; 還原
```

`sub_19796` 把 `dx`／`bx` 換算成 VRAM 位址，**單位是像素**：

```asm
dx >>= 3            ; x ÷ 8 ＝ 每列 80 byte 裡的第幾個 byte
bx <<= 4            ; y × 16
dx += bx            ; y × 16 + x/8
bx <<= 2 / bx += dx ; ★ y × 64 再加上去 ⇒ y × 80 + x/8
```

`cx = 600Ah` 一路傳到 `sub_1FAC2` → `sub_1FB11`：後者 `shl al, 1` 之後
拿 `cl` 當每列 byte 數、`dh` 當列數、`add bp, 50h` 換行。

```
每列 0Ah × 2 ＝ 14h ＝ 20 byte ⇒ 160 px
列數 60h                       ⇒  96 px
```

⭐ **選單框 ＝ (80, 176, 160, 96)。**

### 2.1 96 是「五列」算出來的，不是巧合

框內上下各內縮 8 px、一列 16 px：`8 + 5×16 + 8 ＝ 96`。
而 `sub_193E9` 收到的 `dl = 5` 正是那個 5——**兩條互不相干的路徑
（存畫面的區塊高、選單本體的列數立即值）算出同一個數**。

`dl`／`dh` 是被 patch 進選單迴圈本體的立即值：

```asm
00019440  db 0B2h, 09h    ; mov dl, 9   ← sub_193E9 執行前 xchg 成 5
00019442  db 0B6h, 0Fh    ; mov dh, 0Fh ← 換成 0Bh
```

IDA 把這四個 byte render 成 `db`，因為 `sub_193E9` 開頭
`xchg dl, cs:byte_19441` 會寫進自己的程式碼區段
（`CLAUDE.md` §7 第 28 條的判準）。

### 2.2 選項有幾個是傳進去的，框不跟著縮

`sub_13B5A` 傳 `al = 5`、`sub_13902` 傳 `al = 3`，而框的矩形三處相同。
**選項少的時候下面就是空的**，框不會縮。

## 3. remake 實作

| 項目 | 作法 |
|---|---|
| 說服畫面 | `drawAdvise` 的 `advisePersuade` 畫 composite：`drawIventScene(0)` ＋ 兩個 `drawLegacyTalkBox` ＋ `drawLegacyChoiceBox` |
| 兩個框的內容 | `game.adviseLordSaid`／`adviseAdvisorSaid` 各存**一句**（已折好的行）；`adviseSay(who, index)` 換掉其中一個 |
| 說話者 | 靠框的位置與肖像，句子裡沒有標記。上框 `playerLordPortrait()`、下框 `playerAdvisorPortrait()` |
| 選單框 | `talkChoiceW/H = 160 × 96`、`talkChoiceRows = 5`；`drawLegacyChoiceBox` 與 `talkChoiceClick` 都用 `talkChoiceRows` |
| 能選幾列 | `talkChoiceClick(rows)` 由呼叫端傳——說服 5、外交與撥款 3，對應原版的 `al`（§2.2）|

選單框的尺寸不只影響進言：外交（`drawDiplomacy`）與撥款（`drawFunding`）
畫的是同一個框，三處一起對齊。

### 3.1 保留的 remake 差異

- **選指令那一階段**（敵對／停戰／協力）仍是自己的小視窗。
  原版走的是戰略指令樹（`docs/re/22`），那一層還沒重做。
- **原版一次只顯示一句、按鍵才往下走**；remake 沒有這個逐句節拍，
  兩個框直接顯示最新一句。
- **按鍵提示**是 remake 加的細框，放在橫幅與上框之間——
  畫面底部 360–392 是事件視窗的位置，放那裡會疊在一起。

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestChoiceBoxMatchesOriginalRect`：(80, 176, 160, 96) 與五列，且 `8 + 列數×16 + 8 == 高` |
| 單元測試 | `TestChoiceClickCoversRequestedRowsOnly`：五列剛好填滿框內緣，最後一列不越界 |
| 單元測試 | `TestAdviseSceneLinesGoToTheRightBox`：君主的句子進上框、軍師的進下框，各自帶對的肖像 |
| 截圖 | [`../playtest/34`](../playtest/34-advise-scene-screens.md)：五選一選單與提出理由之後各一張 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| `dh = 0Bh` | 選單本體的第二個立即值，語意未定（不是框寬——框寬由 `cx = 600Ah` 決定）|
| 逐句節拍 | 原版每句要等玩家按鍵，`sub_10241`／`sub_102C2` 那一段還沒讀 |
| 選單的反白樣式 | 原版怎麼畫游標列沒解，remake 用自己的反白條 |
