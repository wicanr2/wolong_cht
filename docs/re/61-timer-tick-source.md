# 61 — 計時中斷是誰發的：節流的頻率終於有數字了

**狀態：✅ 解出來了。** `ds:0D2Dh` 那個計數器由**音效驅動**加，
頻率 **291.3 Hz**（PIT 重設成 4660.9 Hz，驅動自己分頻 16）。
`docs/re/06` §4 只知道「等 N 個計時中斷」，現在 N 有了單位，
兩層速度設定的實際幀率／日長都算得出來（§4）。
順帶解出**戰術速度存在 `ds:0CFBh`**。

- 日期：2026-08-16
- 範圍：松崗 DOS/V。`KI.EXE` ＋ `YNSOUND.COM`
- 原始檔 SHA-256：
  `KI.EXE` `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`／
  `YNSOUND.COM` `e2c6a6a8576c4f2a96b7e3f156d7f48c9570ae03539fe9367adb78aebb364fa1`
- 資料庫 SHA-256：`KI.EXE.i64` `8a8fd7d528e0498000fd04282300588b637fbcb8aa48deb09242f1f41f532691`／
  `YNSOUND.COM.i64` `8590d94aded4b5eb149b44f23ad5c41788e4b155534b2d8bfa3fa7e276949eff`
- 工具：`tools/ida_dump.py`、`tools/ida_decode.py`（TSR 大半是 `db`，要強制解碼）
- 位址空間：IDA DOS/V linear address（`DS` ＝ `CS`，`ds:0CFAh` ＝ linear `0x10CFA`）

## 1. 遊戲這一側：回呼只有 6 條指令

`KI.EXE` 註冊一支 far 回呼給音效驅動，回呼本身什麼都不做，
只留下「有沒有跳過」與「跳了幾次」：

```asm
0001033B  push ds / push ax / push dx
0001033E  ax = cs / ds = ax
00010342  dx = 356h                  ; ★ 回呼位址 ＝ cs:0356h
00010345  ax = 0C00h / int 61h       ; ★ 向音效驅動登記
0001034A  pop dx / pop ax / pop ds / retn

0001034E  ax = 0C01h / int 61h       ; 取消登記
00010355  retn

00010356  cli / pushf                ; ★ 回呼本體
00010358  mov cs:byte_10D2C, 1       ;   旗標：這一輪有中斷來過
0001035E  cmp cs:byte_10D2D, 0FFh
00010364  jnb  +5                    ;   計數器飽和在 0FFh，不回繞
00010366  inc cs:byte_10D2D
0001036B  popf / retf
```

`byte_10D2C`／`byte_10D2D` 就是 `docs/re/06` §4 等的那兩個
（`ds:0D2Ch`／`ds:0D2Dh`——同一段，`DS` ＝ `CS`）。

⚠ **計數器飽和不回繞**，所以「等到計數器 ≥ N」在極端情況會早退不會漏等。

## 2. 驅動這一側：`INT 61h AH=0Ch` 就是「登記時鐘回呼」

`YNSOUND.COM` 的 `INT 61h` 服務裡（COM offset `0x2FE`）：

```asm
   cmp al, 1 / jz .取消
   mov cs:[093Eh], dx        ; ★ 把 DS:DX 直接寫進 call 的運算元
   mov cs:[0940h], ds
   retn
.取消:
   ax = 097Ah                ; 指回驅動自己的 retf
   cs:[093Eh] = ax / cs:[0940h] = cs
   retn
```

`cs:093Eh` 是 §3 那條 `call far ptr 0:0` 的運算元——**靜態影像裡的 `0:0`
是「還沒登記」，不是分析失敗**（`CLAUDE.md` §7 第 26 條的同一個形狀）。
取消時填回 `0x097A`，那裡就是一個 `retf`。

## 3. 時鐘鏈：PIT 4660.9 Hz → 回呼 291.3 Hz → BIOS 18.2 Hz

安裝時（COM offset `0x8B6`）：

```asm
   ax = 3508h / int 21h            ; 取原本的 INT 8
   cs:[094Dh] = bx / cs:[094Fh] = es
   ds = cs / dx = 913h
   ax = 2508h / int 21h            ; ★ 把 INT 8 換成自己的
   al = 36h / out 43h, al          ; ★ 重設 PIT ch0
   ax = 0100h / out 40h, al / al = ah / out 40h, al   ; ★ 除數 ＝ 0x0100
   cs:[0B6Ah] = 10h                ; 回呼分頻：16
   cs:[0B6Bh] = 0                  ; BIOS 分頻：0 → 第一次 dec 變 0FFh
```

處理常式（COM offset `0x813`）：

```asm
   cli / pushf / push ax
   test cs:[099Eh], 2 / jz .跳過音樂
     dec cs:[0B69h] / jnz .跳過音樂
     cs:[0B69h] = cs:[0B68h]       ; 音樂 tempo 分頻器，值由 0x859 依曲速算
     call 播放引擎
.跳過音樂:
   dec cs:[0B6Ah] / jnz .下一段
     cs:[0B6Ah] = 10h
     call far ptr 0:0              ; ★ 遊戲登記的回呼（§2 patch 進來的）
     cli
.下一段:
   dec cs:[0B6Bh] / jnz .EOI
     pop ax / popf / jmp far ptr 原本的 INT 8    ; ★ 每 256 次鏈一次
.EOI:
   al = 20h / out 20h, al / pop ax / popf / sti / iret
```

三個頻率：

| 層 | 除數 | 頻率 | 用途 |
|---|--:|--:|---|
| PIT ch0 | `0x0100` ＝ 256 | **4660.87 Hz** | 音樂引擎的最小刻度 |
| 回呼分頻 `cs:0B6Ah` | 16 | **291.30 Hz** | ★ 遊戲的節流單位 |
| BIOS 分頻 `cs:0B6Bh` | 256 | **18.206 Hz** | 鏈回原本的 `INT 8`，系統時鐘不會走鐘 |

1,193,182 ÷ 256 ÷ 16 ＝ 291.30；÷ 256 ＝ 18.206 ＝ 標準 BIOS tick。
**驅動沒有動系統時間**，只是在中間插了一層自己的快時鐘。

⇒ 沒有音效驅動時回呼不會發生，`ds:0D2Dh` 恆為 0，
**兩層節流的等待迴圈都會卡死**。這解釋了為什麼 `YNSOUND.COM`
是啟動流程的必要條件而不是選配。

## 4. ⭐ 於是兩層速度都有數字了

### 戰略：`ds:0CFAh` 直接當「等幾個回呼」

`sub_11D8E` 每跑一次推進一個**子刻**，之後等 `ds:0CFAh` 個回呼
（`docs/re/06` §4）。一個遊戲日 ＝ 24 時 × 9 子刻 ＝ **216 子刻**：

| 設定 | 值 | 子刻／秒 | 一個遊戲日 |
|---|--:|--:|--:|
| 最高速 | 0 | 不等待 | 機器上限 |
| 　高速 | 1 | 291.3 | **0.74 秒** |
| 　普通 | 2 | 145.7 | 1.48 秒 |
| 　低速 | 3 | 97.1 | 2.22 秒 |
| 最低速 | 4 | 72.8 | 2.97 秒 |

### 戰術：`ds:0CFBh` 先 ×16 再當「等幾個回呼」

系統選單第 5 列的 handler 只做一件事（`0x160A5`）：

```asm
000160A5  al = ds:0CFBh
000160A8  shl al,1 ×4          ; ★ ×16
000160B0  ds:0CFCh = al
```

戰場主迴圈（`0x1A0F2`）等的是 `ds:0CFCh`，形狀與戰略層一模一樣：

```asm
0001A0F2  al = ds:0CFCh / and al,al / jz 略過
0001A0F9  cmp byte ds:0D2Ch, 0 / jz $      ; 等旗標
0001A100  cmp ds:0D2Dh, al / jb $          ; 等計數器 ≥ al
          ds:0D2Dh = 0 / ds:0D2Ch = 0
```

| 設定 | 值 | 等幾個回呼 | 戰場幀率 |
|---|--:|--:|--:|
| 最高速 | 0 | 0 | 機器上限 |
| 　高速 | 1 | 16 | **18.2 fps** |
| 　普通 | 2 | 32 | 9.1 fps |
| 　低速 | 3 | 48 | 6.1 fps |
| 最低速 | 4 | 64 | 4.6 fps |

⭐ **那個 ×16 不是隨便挑的**：16 正好是驅動自己的回呼分頻值，
所以戰術「高速」＝ 每個 PIT‑tick 群一幀 ＝ 標準 BIOS tick 的 18.2 Hz。
戰場動畫是照 18.2 fps 設計的，戰略層則細 16 倍。

⚠ **兩層的設定值語意不同**（一個直接用、一個 ×16），
但選單上是同一組標籤、同一張表（[`55`](55-system-menu-window.md) §4）。
remake 若把兩個速度接成同一條換算會錯 16 倍。

## 5. 另一個用同一組計數器的地方

`sub_11CD0`（戰略主迴圈的一步）開頭有一段**等 0xA0 個回呼**的狀態機：

```asm
00011CED  cmp byte ptr ds:0D2Dh, 0A0h / jb retn    ; 0xA0 ＝ 160 ≈ 0.55 秒
00011CF4  ax = 2 / call sub_20000                  ; 滑鼠服務
          ds:0D2Ch = 0 / ds:0D2Dh = 0 / ds:98A5h = 0
```

`ds:98A5h` 是一個小狀態：`> 1` 時每輪只清計數器並遞減，
`== 1` 時等滿 0xA0 個回呼（**約 0.55 秒**）才做事並歸 0，
`== 0` 時走正常流程。也就是一個**以回呼為單位的延時器**，
0.55 秒這個量級像是提示訊息的停留時間。內容沒讀完，記在這裡是因為
**它證明 `ds:0D2Dh` 不是節流專用**——改動計數器語意會影響到這裡。

## 6. 未解

| 項目 | 現況 |
|---|---|
| `ds:98A5h` 那個延時器實際在等什麼 | §5。呼叫端 `sub_11BE0` 沒讀 |
| 音樂 tempo 分頻器 `cs:0B68h` 的算式 | `0x859` 那 20 條指令：`al = ((0FFh − ah) × 13) >> 3`，`ah` 從哪來沒讀 |
| `cs:099Eh` 的 bit 1 | 「音樂啟用」是從用法推的，寫入端沒讀 |
| 無音效驅動時的行為 | §3 推論「會卡死」，**沒有實測**——DOSBox 拿掉 `YNSOUND.COM` 跑一次就能驗 |
