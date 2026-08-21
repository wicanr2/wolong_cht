# 28 — 文字與數字的繪製層

**狀態：數字繪製、兩支字串繪製與 EGA 平面寫入方式 confirmed。
單字元 blitter `loc_1F75E` 與字型查表 `sub_1F878` 的內部未逐行讀。**

- 日期：2026-08-13
- 範圍：只驗松崗 DOS/V
- 原始 `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 原始 `KI.EXE.i64` SHA-256：`7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`
- 工具：IDA Pro 9.4；自我修改區另以 `KI.EXE` 原始 bytes 交叉解碼
- `KI.EXE` 檔案位移 ＝ 段內偏移 ＋ `0x200`
- 位址空間：IDA DOS/V linear address，segment base `0x10000`

## 1. 底層：EGA 的 Set/Reset

所有文字與數字都不是逐像素寫顏色，而是靠 EGA 的 Set/Reset 暫存器：
先把顏色寫進埠 `0x3CF`，再把字型的點陣 byte 寫進 VRAM，
硬體就用那個顏色把有設位元的像素填上。

```asm
mov dx, 3CFh
mov al, <顏色>
out dx, al        ; EGA Graphics Controller data register
mov al, es:[di]   ; dummy read，載入 latch
movsb             ; 寫字型 byte
add di, 4Fh       ; 下一列
```

- VRAM 段是 `0A0C8h`。相對標準的 `0A000h` 多了 `0xC8` 個 paragraph
  ＝ `0xC80` bytes ＝ 40 列，所以這一組常式的 y 原點在螢幕第 40 列。
- **列距是 `0x50` ＝ 80 bytes**（`movsb` 已經 +1，所以只加 `0x4F`）。
  80 bytes × 8 像素 ＝ 640，與熱區圖的定址式
  （[`22`](22-strategy-command-tree.md) §2）算出來的寬度一致。
- 每個字高 **16 列**，畫完 `di -= 0x500`（16 × `0x50`）回到起點。

`movsb` 之前的 `mov al, es:[di]` 是**必要的**：EGA 寫入前要先讀一次把
四個平面載進 latch，否則沒被字型位元覆蓋的平面會被清掉。
移植時漏掉這一步會得到「顏色只剩一個平面」的畫面。

## 2. `sub_1062F`：數字

```
輸入  dx:ax = 32 位有號值
      bl    = 位數
      bh    = 背景色（高 4 位）: 前景色（低 4 位）
      di    = VRAM 位置（最左一位）
      ds    = cs:word_10D54（字型段）
```

```asm
cl = bl / dec cx / add di, cx      ; 移到最右一位，從個位往左印
bp = dx                            ; 留著判斷正負
cmp dh, 80h / jb →                 ; 負數
  not dx / not ax / add ax,1 / adc dx,0   ; 32 位兩補數取絕對值
  dec bl                           ; 空出一位給負號
loop:
  cx = 0Ah / div cx                ; dx = 餘數 ＝ 這一位
  sub_1069A                        ; 印一位
  dec di / dec bl / jz done
  and ax, ax / jnz loop            ; 還有高位就繼續
  ; 值印完了，剩下的位置補背景
  dh = bh >> 4 / ah = dh
  sub_106DE                        ; 用背景色填滿一格
  dec di / dec bl / jnz 補迴圈
done:
  cmp bp, 0 / jge →
    dx = 0Ah / sub_1069A           ; 字型第 10 格 ＝ 負號
```

三件事定案：

1. **右對齊**，從個位往左印，前導位置用**背景色**填滿（不是空白字元）。
2. **負數走 32 位兩補數取絕對值**，並少印一位留給負號；
   負號是字型的第 10 格（0–9 是數字）。
3. `bl == 1` 的保護：`dec bl` 之後若歸零會補回 1，避免位數變成 0。

### 2.1 換色警示改的是前景色

`sub_1069A` 分兩趟畫一個數字：

```asm
al = bh & 0Fh / out 3CFh, al   ; ← 前景色
… 用字型 byte 畫 …
al = bh >> 4   / out 3CFh, al  ; ← 背景色
… 用字型 byte 的反相畫 …
```

所以 `bh` 是一對 4-bit 顏色。一覽表正常用 `bh = 0x90`（背景 9、前景 0），
警示時改成 `0x9A`（背景 9、前景 **A**）——
[`27`](27-list-row-fields.md) §5 的三條換色規則，機制都是**改前景色**。

## 3. `loc_10701`：固定三個全形字

名稱欄（武將名、據點名、勢力名）走這一支。它**只印三個全形字**，
不看結束符：

```asm
mov cs:[071Fh], ax / mov cs:[0736h], ax / mov cs:[074Dh], ax   ; ← 自我修改
di = dx / bp = bx
ax = [si]   / xchg al, ah / sub_1F878 / cx = ax / mov ax, <屬性> / loc_1F75E
dx = di+10h / ax = [si+2] / xchg al, ah / sub_1F878 / cx = ax / mov ax, <屬性> / loc_1F75E
dx = di+20h / ax = [si+4] / xchg al, ah / sub_1F878 / cx = ax / mov ax, <屬性> / loc_1F75E
```

- **三個字，間距 16 像素**，整欄寬 48 像素。
  記錄裡的名稱欄正好是 6 bytes（3 個 Big5 字），兩邊剛好對上，
  所以顯示層沒有截字邏輯。
- `xchg al, ah`：Big5 是高位在前，x86 讀進 `ax` 之後要交換才是字碼。
- **自我修改**：屬性被寫進三個 inline `mov ax, imm16` 的運算元
  （`0x071F`／`0x0736`／`0x074D`）。原因是 `ax` 同時要當字碼與屬性用，
  patch 立即值比每次 push/pop 省。IDA 因此在這一區建不出乾淨的函式，
  上面的解碼是拿 `KI.EXE` 的原始 bytes 對出來的。

> 這是這份程式碼第二處自我修改（另一處是 `loc_1A065`，
> [`20`](20-ida-re-coverage-audit.md) §2.2）。**看到 IDA 解出不合理的指令，
> 先懷疑自我修改，不要先懷疑分析壞掉。**

## 4. `sub_1F6DC`：變長字串

欄位標題、分隔線與空列的「－－－」走這一支：

```asm
di = 10h                      ; 預設步進 16（全形）
loop:
  ch = [si] / si++
  if ch == 0 → done           ; NUL 結尾
  cl = [si] / si++
  xchg ax, cx / sub_1F878     ; 查字型；CF 回報寬度
  jnb → di = 8 / dec si       ; ← 半形：步進 8，退回多讀的那個 byte
  xchg ax, cx
  test al, 4 / jz →           ; 屬性 bit 2 ＝ 陰影
    inc bx / inc dx / ah >>= 4 / loc_1F75E    ; 先在 (x+1,y+1) 用高 4 位顏色畫
  loc_1F75E                   ; 再畫本體
  dx += di
```

- **NUL 結尾，長度不限。**
- **全形／半形自動判斷**：`sub_1F878` 用 CF 回報這個碼是不是雙位元組，
  是就步進 16、不是就退回一個 byte 並步進 8。
  這解釋了分隔線 `－－－　 ----` 為什麼能全形半形混排。
- **屬性 bit 2 ＝ 陰影**：先在右下一像素用高 4 位的顏色畫一次再畫本體。

## 5. 兩支的分工

| | `loc_10701` | `sub_1F6DC` |
|---|---|---|
| 長度 | 固定 3 個全形字 | NUL 結尾，不限 |
| 寬度判斷 | 一律全形 | 自動全／半形 |
| 陰影 | 無 | 屬性 bit 2 |
| 用在 | 名稱欄（武將／據點／勢力）| 標題、分隔線、空列 |

名稱欄用固定版是因為記錄裡的名稱就是 6 bytes，不必檢查結束符——
**資料的固定長度直接換成了程式碼的簡化**。

## 6. 未解

| 項目 | 現況 |
|---|---|
| `sub_1F7A4` | 把 32 B 字模緩衝畫上 VRAM 的實際迴圈，未逐行讀 |
| 屬性的其餘位元 | bit 2 是陰影已證實；`0x9001`／`0x9000` 的 bit 0 差在哪未讀 |
| `word_10D4C` 那一組 | 來源已解——`sub_100DF` 開機把 `ICONGRF` 段 3 切五塊，`word_10D54` 是 `+0x0840` 的 11 格 × 16 列數字字模（[`../spec/52`](../spec/52-main-screen-camera-and-banner-date.md) §4）；緊接在後的 `+0x08F0` 另有一組 11 格，用途未解 |
