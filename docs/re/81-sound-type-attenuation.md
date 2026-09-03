# 81 — 系統選單的「TYPE 1–4」是四段衰減，不是四種音源

**狀態：confirmed（原始 bytes 交叉解碼）。** 音效那一列的五個選項
（ＯＦＦ／TYPE 1／2／3／4）在 `KI.EXE` 這一側只換一個數字：
`INT 61h AH=0Bh` 的 `AL ＝ (值 − 1) × 4`。TSR 收到之後把它**加進 OPL3
operator 的 Total Level**——也就是**衰減**。四個 TYPE 是四段音量，
每段 4 個 TL 單位（OPL 的 TL 是 0.75 dB/step，約 3 dB）。

- 日期：2026-09-03
- 原始輸入：`workplace/orig/dosv/KI.EXE`，SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`；
  `workplace/orig/dosv/YNSOUND.COM`，SHA-256
  `e2c6a6a8576c4f2a96b7e3f156d7f48c9570ae03539fe9367adb78aebb364fa1`
- 位址：`KI.EXE` 用 IDA 線性位址；`YNSOUND.COM` 用**載入後的記憶體位移**
  （COM 的 org ＝ `0x100`，所以檔案位移 ＝ 記憶體位移 − `0x100`）
- 相關：[`17`](17-dosv-audio-tsr.md)（TSR 的 `INT 61h` 介面）、
  [`55`](55-system-menu-window.md) §4（那一列的五個選項）、
  [`57`](57-opl3-register-map.md)（OPL3 暫存器）

## 1. `KI.EXE` 側：`sub_102D0` 只送一個數字

```asm
sub_102D0:                       ; IDA 000102D0，37 B
        push ax
        mov  al, cs:byte_10CF9   ; 音效設定值 0–4（ds:0CF9h，docs/re/55 §4）
        cmp  al, 1
        jb   short loc_102EF     ; 0 ⇒ 關
        ja   short loc_102DF      ; ≥2 ⇒ 只換型別
        mov  ah, 7 / int 61h     ; ★ 1 ⇒ 先開，再往下換型別
loc_102DF:
        mov  al, cs:byte_10CF9
        dec  al
        shl  al, 1 / shl  al, 1  ; ★ (值 − 1) × 4 ⇒ 0／4／8／12
        mov  ah, 0Bh / int 61h
        jmp  short loc_102F3
loc_102EF:
        mov  ah, 8 / int 61h     ; ★ 0 ⇒ 停止並清空狀態
loc_102F3:
        pop  ax
        retn
```

| 設定值 | 顯示 | 送出去的 |
|--:|---|---|
| 0 | ＯＦＦ | `AH=8`（停止／清空）|
| 1 | TYPE 1 | `AH=7`（初始化）＋ `AH=0Bh, AL=0` |
| 2 | TYPE 2 | `AH=0Bh, AL=4` |
| 3 | TYPE 3 | `AH=0Bh, AL=8` |
| 4 | TYPE 4 | `AH=0Bh, AL=12` |

⭐ **只有從 ＯＦＦ 走到 TYPE 1 才送 `AH=7`。** 選項是**環狀遞增**的
（`sub_16062`：加一、到頂繞回 0，[`55`](55-system-menu-window.md) §4），
所以要到 TYPE 2 一定先經過 TYPE 1——`AH=7` 已經送過了。
**這一條在 remake 沒有對應物，但它解釋了為什麼 `al ≥ 2` 那一支看起來少做一件事。**

## 2. TSR 側：`AH=0Bh` 把那個數字加進 Total Level

`INT 61h` 的分派表在記憶體 `0x0115`（`add bx, 0115h`，
[`17`](17-dosv-audio-tsr.md) §1），**13 筆 word，AH ＝ 0–12**。
第 14 筆起是程式碼，不是表項。`AH=0Bh` ⇒ handler `0x02DE`：

```asm
0x02DE: push ds / push es
        mov  cs:[0996h], al      ; ★ 存下型別參數（0／4／8／12）
        mov  ax, cs / mov ds, ax
        mov  es, ds:[0992h]
        xor  ah, ah
.loop:  call 060Fh               ; 對 ah = 0、1、2 各重算一次
        inc  ah / cmp ah, 3 / jne .loop
        pop  es / pop ds / retn
```

`cs:[0996h]` **全檔只有一個寫入端（這裡）與一個讀取端**（`0x065D`）：

```asm
0x060F: ...
        mov  cl, 0Fh
        sub  cl, [bx+0A56h]      ; 這一聲部自己的音量
        shl  cl, 1 / shl cl, 1   ; ×4
        add  cl, cs:[0996h]      ; ★ 加上型別參數
        cmp  cl, 3Fh / jb  $+4 / mov cl, 3Fh
        ...
        mov  al, ah / add al, 40h ; ★ OPL 暫存器 40h+n ＝ KSL/TL
        ...
.op:    mov  ah, es:[si]         ; 音色資料的那一個 byte
        mov  bh, ah
        and  ah, 3Fh             ; TL（bit 0–5）
        and  bh, 0C0h            ; KSL（bit 6–7）
        add  ah, cl              ; ★ 衰減加上去
        cmp  ah, 3Fh / jb  $+4 / mov ah, 3Fh
        or   ah, bh
        call 0890h               ; 寫進晶片
        inc  si
        add  al, [bx+0A52h]      ; 下一個 operator 的暫存器位移
        inc  bx / cmp bl, 4 / jb .op
```

三件事同時對上，所以這不是「聽起來像」的推論：

1. **暫存器位址是 `40h + n`**——OPL 的 `40h`–`55h` 就是 KSL/TL
   （[`57`](57-opl3-register-map.md)）。
2. **欄位切法是 `& 3Fh` 與 `& 0C0h`**——TL 六位元、KSL 兩位元。
3. **上限夾在 `3Fh`**——TL 的最大值，也就是最安靜。

⇒ `AL` 越大，加得越多，聲音越小。**TYPE 1 最大聲，TYPE 4 最小聲，
每段差 4 個 TL 單位。**

⚠ `0x060F` 只被 `AH=0Bh` 的三次迴圈與別處呼叫；**它重算的是既有聲部的
TL，不是把值記進音色表**——所以型別要在每次載入新音色之後再送一次，
或由播放路徑自己帶上。這一段沒逐行讀（§4）。

## 3. 順帶定案的三個服務號

[`17`](17-dosv-audio-tsr.md) §7 掛著「`ah=4`／`7`／`8` 對應什麼動作要看
`YNSOUND.COM`」。三個都讀了：

| AH | handler | 作用 |
|--:|---|---|
| `04` | `0x01AC` | **卸載**：`INT 21h AX=251Ch` 從 `cs:[0B62h]` 還原 INT 1Ch 向量，清狀態，六個聲部逐一靜音 |
| `07` | `0x01F6` | 初始化（資源 layout parser，[`17`](17-dosv-audio-tsr.md) §2 已證實）|
| `08` | `0x0269` | **停止並清空**：`cs:[099Eh]` bit 1 沒設就直接回；設了就把 `099Ah` 起 0x23 words、`0A56h` 起 0x2D words 清零，六個聲部靜音，最後清掉 bit 1 |
| `0C` | `0x02FB` | 設回呼的 far pointer；`AL=1` 走內建的 `cs:097Ah`。`KI.EXE` 的 `ax=0C01h`（[`42`](42-leaf-functions.md) §7）就是**恢復內建回呼** |

⇒ `AH=4` 與 `AH=8` 都會靜音，差別是 **`AH=4` 連中斷向量一起還原**（要退出程式），
`AH=8` 只停播放（要換曲子或關音效）。

## 4. remake 怎麼接

remake 的音訊是純 Go 的 OPL3 渲染（[`../spec/29`](../spec/29-audio.md)），
不經過 DOS，也沒有 TSR。**四個 TYPE 對 remake 是一個「主衰減」參數**，
規格與實作見 [`../spec/122`](../spec/122-sound-type-levels.md)。

## 5. 未解

| 項目 | 現況 |
|---|---|
| `0x060F` 的音色來源 | `es:[si]` 從 `ds:[0992h]` 指的段讀，`si` 由 `[bx+0A5Ch] << 5 + ds:[099Ah]` 算出。**這一段的資料結構沒逐行讀**，所以「換音色之後型別要不要重送」還不確定 |
| `AH=0Bh` 只重算三個聲部 | 迴圈是 `ah = 0、1、2`，而靜音路徑跑的是六個（`ah = 5..0`）。剩下三個什麼時候拿到新的衰減沒查 |
| `AH=09h`（`ax=09F2h`）| `cs:[099Eh]` bit 1 沒設就直接回；設了就對 `ds:[0A4Ch]` 個聲部逐一呼叫 `0x049E`，參數 `al=91h`／`ah=0F2h`。`0x049E` 沒讀 |
