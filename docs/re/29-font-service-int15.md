# 29 — 原版怎麼顯示中文：INT 15h 字型服務與 `END_S13/S14.DAT`

**狀態：整條鏈 confirmed（靜態）。`KI.EXE` 側走 DOS/V 的 `INT 15h AH=50h`
向常駐服務要字模，`STR.EXE` 是那個服務，字模來自 16×15 的 Big5 點陣字檔。
唯一未解的是 `STR.EXE` 寫死的檔名與封裝內的資料檔不同步（§6）。**

- 日期：2026-08-13
- 範圍：松崗 DOS/V。PC-98 側不適用（有字型 ROM，`int 15h` 出現 0 次）
- `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- `STR.EXE` 916 B；解碼方式：直接讀 `KI.EXE`／`STR.EXE` 的原始 bytes
  （`STR.EXE` 的常駐段只有 404 B 載入映像，IDA 建不出有意義的函式）
- 位址：`KI.EXE` 用 IDA DOS/V linear address；`STR.EXE` 用**常駐段內偏移**

這一份結掉 CLAUDE.md §2.1／`docs/reference/04` 長期掛著的
「原版怎麼顯示中文仍未解」。

## 1. 全鏈一句話

```
KI.EXE 開機  →  sub_1F720 發兩次 INT 15h AX=5000h（BH=0 半形／BH=1 全形）
              ←  常駐服務回 ES:BX ＝ 取字模常式的遠位址
              →  把兩個位址 patch 進 loc_1F75E 裡的兩條 call far 0000:0000

畫一個字    →  loc_1F75E 在堆疊開 32 B 緩衝，call far <向量>
              →  常式 open/lseek/read/close 字型檔，把 30 B（全形）
                 或 15 B（半形）搬進緩衝，第 16 列補 0
              →  回到 loc_1F75E，用 EGA Set/Reset 把緩衝畫上 VRAM（docs/re/28 §1）
```

**字型不在 `KI.EXE` 裡。** 這解釋了 2026-08-08 那次「把 `YNFONT.EXE`
整個當 1bpp 點陣圖畫出來沒有任何字形」——找錯檔案了。

## 2. `sub_1F720`：向服務要兩個入口

`0001F720`，被 `sub_1006B`（開機初始化）呼叫一次。

```asm
mov ax, 5000h / mov bx, 0     / mov dx, 810h  / xor bp, bp / int 15h
jb  → 跳過                     ; CF=1 ＝ 沒有服務
mov cs:word_1F791, bx / mov cs:word_1F793, es   ; ← 半形向量
mov ax, 5000h / mov bx, 100h  / mov dx, 1010h / xor bp, bp / int 15h
jb  → 跳過
mov cs:word_1F78A, bx / mov cs:word_1F78C, es   ; ← 全形向量
```

- `AH=50h` 是 DOS/V 的字型服務號，`AL=0` ＝取得字型存取入口。
- `BH` 選字型：**0 ＝ 半形、1 ＝ 全形**。
- `DX` 是字型尺寸：`0x0810` ＝ **8×16**、`0x1010` ＝ **16×16**。
  ⚠ 這裡宣告的是 16 列，但實際檔案每個字只有 **15 列**，
  第 16 列由取字模常式自己補 0（§4）。
- **`jb` 之後不設任何預設值。** 沒有服務時兩個向量停在 `0000:0000`，
  下一次畫字就是 `call far 0000:0000`。原版沒有 fallback。

## 3. `loc_1F75E`：兩條被 patch 的遠呼叫

`loc_1F75E` 是單字元 blitter（[`28`](28-text-number-rendering.md) §6 列為未解）。
它的骨架：

```asm
sub sp, 20h              ; 32 B 字模緩衝（16×16 1bpp）
… 裁切：dx > 270h（624）或 bx > 180h（384）就不畫 …
and ch, ch
jz  → 半形
9A 00 00 00 00           ; 0001F789: call far 0000:0000  ← 全形，運算元在 1F78A
jmp →
9A 00 00 00 00           ; 0001F790: call far 0000:0000  ← 半形，運算元在 1F791
call sub_1F7A4           ; 把緩衝畫上 VRAM
```

`ch` 是 Big5 高位元組（`sub_1F878` 已把寬度分類完，見
[`28`](28-text-number-rendering.md) §4），非 0 就走全形。

> **⚠ 這是第三處執行期改寫程式碼。** 前兩處是 `loc_1A065`
> （[`20`](20-ida-re-coverage-audit.md) §2.2）與 `loc_10701`
> （[`28`](28-text-number-rendering.md) §3）。差別是前兩處**自己改自己**，
> 這一處是**開機時被 `sub_1F720` 改**，兩者相距 0x70 個 byte。
> IDA 在 `1F789` 只會顯示 `call far 0:0`——看到「呼叫零位址」不要
> 判定分析壞掉，先找誰在寫那個位址。

## 4. `STR.EXE`：常駐的字型服務

916 B ＝ 512 B MZ 檔頭 ＋ **404 B 載入映像**。常駐段的配置：

| 偏移 | 內容 |
|---|---|
| `+00` | `"DOSV"` 簽章 |
| `+04` | 舊的 INT 15h 向量（offset／segment）|
| `+10` | `start` 用的暫存 |
| `+14` | 全形字型檔名字串 |
| `+20` | 半形字型檔名字串 |
| `+2C` | INT 15h handler |
| `+5D` | **全形取字模常式**（服務回傳的入口）|
| `+DF` | **半形取字模常式** |
| `+11F` | `start`（安裝／解除）|

### 4.1 安裝：同一支執行檔負責裝與卸

`start` 先 `INT 21h AX=3515h` 取現行 INT 15h 向量，檢查該段開頭是不是
`"DOSV"`：**是 → 這次執行等於解除常駐**（還原向量、釋放記憶體）；
**否 → 存舊向量、`AX=2515h` 把自己的 handler 掛上、`INT 21h AH=31h` 常駐**。

`YNVSHELL.COM` 的程式表裡 `STR.EXE` 排在 `YNFONT.EXE` 之後、`KI.EXE` 之前，
而且 `YNVSHELL.COM` 內有一個 `" /R"` 字串——與「跑第二次就解除」的設計對得上。

### 4.2 handler：先讓路，再自己回答

```asm
cmp ah, 50h / jnz → jmp far cs:[04]      ; 不是字型服務就轉給舊 handler
pushf / call far cs:[04]                 ; 是的話先讓舊 handler 處理
cmp al, 0 / jnz → iret                   ; 別人處理掉了就不插手
stc / mov ah, 1                          ; 預設回失敗
cmp bh, 1 / ja → iret                    ; 只認得 BH = 0 或 1
mov ax, cs / mov es, ax
or bh, bh / jnz → bx = 5Dh               ; BH=1 全形
bx = 0DFh                                ; BH=0 半形
sub ah, ah / clc / iret                  ; 回 ES:BX
```

**這是 fallback 設計**：真正的 DOS/V 環境（`$FONT`／`FONTX` 之類）先回答，
沒人回答時才由 `STR.EXE` 頂上。所以松崗版在 DOS/V 上跑用系統字型、
在普通 DOS 上跑用自己帶的字型，同一份 `KI.EXE` 兩邊都能跑。

### 4.3 取字模常式：每畫一個字就開關檔一次

全形（`+5D`）：

```asm
push ds/dx/bx/cx
ds = cs / dx = 14h              ; 全形字型檔名
ax = 3D00h / int 21h            ; open read-only
bx = ax                         ; handle
; ── Big5 → 序號 ──
cmp ch, 0A4h / jb  → ax = (ch-0A1h) * 9Dh
cmp ch, 0C9h / jb  → ax = (ch-0A4h) * 9Dh + 198h     ; +408
cmp ch, 0FAh / jb  → ax = (ch-0C9h) * 9Dh + 16B1h    ; +5809
else                 ax = 56h                         ; 越界 → 固定第 0x56 格
cmp cl, 7Eh / ja   → cl -= 22h
cl -= 40h / ch = 0 / ax += cx
; ── 取 30 B ──
mul 1Eh                          ; dx:ax = 序號 × 30
ax = 4200h / int 21h             ; lseek
ds = es / dx = si                ; 呼叫端的緩衝區
cx = 1Eh / ah = 3Fh / int 21h    ; read 30 B
mov word ptr es:[si+1Eh], 0      ; ← 第 16 列補 0
ah = 3Eh / int 21h               ; close
ax = 0 / pop / retf
```

半形（`+DF`）同形，只是 `ax = cl × 0Fh`、讀 15 B、`mov byte ptr es:[si+0Fh], 0`。

四件事定案：

1. **全形是 16×15、30 B/字；半形是 8×15、15 B/字。** 宣告 16 列、
   實存 15 列，最後一列固定空白。這是**倚天點陣字的規格**，不是 DOS/V 的
   16×16 標準——松崗把倚天字型塞進 DOS/V 的介面裡。
2. **`0x9D` ＝ 157**，是 Big5 每個高位元組底下的字碼數
   （`0x40`–`0x7E` 共 63 ＋ `0xA1`–`0xFE` 共 94）。低位元組 `> 0x7E`
   先減 `0x22` 就把兩段接成連續的 0–156。
3. **越界字碼回第 `0x56` 格**，不是不畫。所以缺字會顯示成某個固定字形。
4. **每畫一個字就 open／lseek／read／close 一次。** 一行 20 個字就是
   80 次 DOS 呼叫。這對「戰略速度受機器效能限制」（CLAUDE.md §2.1）
   是個具體註腳——原版的文字繪製成本主要在磁碟 I/O，
   當年靠 SMARTDRV 之類的快取撐著。
   **remake 不要照抄這個結構**（`internal/assets/cjk` 一次載入整份字型）。

## 5. 字型檔：`END_S13.DAT` 與 `END_S14.DAT`

封裝內只有兩個檔符合上面的規格，而且**只在松崗版存在**（PC-98 沒有）：

| 檔案 | 大小 | 內容 |
|---|---|---|
| `END_S14.DAT` | 3,840 | **與倚天 `ascfont.15` byte-for-byte 相同**（256 字 × 15 B）|
| `END_S13.DAT` | 404,992 | 408 格 Big5 符號區字模 ＋ 倚天 `stdfont.15` 全文（尾端截 68 B）|

`END_S13.DAT` 的組成是逐 byte 驗過的：

```
END_S13[0     : 12240] = 408 格符號區字模（非全零，來源不是 stdfont.15）
END_S13[12240 :      ] = stdfont.15[0 : 392752]        ← 完全相同
stdfont.15 尾端 68 B（2.3 格）未收進來
```

`12240 = 408 × 30`，而 408 正是取字模常式第二段索引的基底
（`+0x198`）——**兩個獨立來源給出同一個數字**，這是最硬的那種對上。
倚天 `stdfont.15` 本來就從 `0xA440`（漢字區起點）開始，
松崗補上前面的符號區就直接對齊了 Big5。

驗收（全部零例外）：

- 「一」「國」「龍」「虎」「牢」「關」「君」「主」「軍」「師」「資」「金」
  「稅」「率」「兵」十五個字照索引式取出來的字模，在倚天 `stdfont.15`
  裡都找得到，**序號一律差 −408**。
- 把 `TALK.DAT`／`SAVE.DAT`／`SINARIO.DAT`／`KI.EXE` 掃出的
  **3,200 個 Big5 候選字**全部照索引式查 `END_S13.DAT`：
  **越界 0 個、字模全空白 0 個。**

> ⚠ 所以 `END_S13.DAT` 是**倚天字型的衍生物**。§10 的「不得散布倚天字型」
> 對它一樣適用；deny-list 已經擋掉所有 `.DAT`，不必另外加規則。

## 6. 未解：檔名對不上

`STR.EXE` 寫死的兩個檔名是 **`END_S10.DAT`／`END_S11.DAT`**，
不是 `END_S13`／`END_S14`。而封裝內的 `END_S10/S11`：

- 大小 40,806／10,637，**都不是 30 或 15 的整數倍**；
- 開頭六個 byte 是 `00 F4 01 00 00 00`，與 `END_S3`／`END_S9`／`END_S12`
  這些結局過場圖**完全一致**；
- `D7END.EXE` 的字串表裡 `END_S1..S12.DAT` 一應俱全（8.3 補空格格式），
  **包含 S10 與 S11**，也就是結局播放器把它們當圖用；
- 照全形索引式取「一」得到亂碼，取「國」直接越界。

**所以封裝內的 `END_S10/S11` 是結局過場圖，`STR.EXE` 若照字面執行會讀到亂碼。**
兩個方向都說得通，目前沒有證據裁決：

1. 安裝流程會把字型檔搬成 `END_S10/S11`（`INSTALL.EXE` 引用
   `STDFONT.24`／`ASCFONT.24`／`ASCFONT.15` 三個檔名，可能是轉檔來源），
   而我們手上的封裝是**安裝媒體**（`DISK1`–`DISK4` 標記檔還在），不是安裝後的目錄；
2. 封裝內的 `STR.EXE` 與資料檔不同步（開發期殘留），實際出貨版檔名已改成 S13/S14。

裁決要靠實跑，而 DOS/V 側的 oracle 被防拷擋著
（[`../playtest/01`](../playtest/01-dosbox-dosv.md)）。
**這一項不擋 remake**——字型的格式、索引式與資料檔都已經確定。

## 7. 兩個可驗證的預測

留在這裡，之後拿到可跑的 DOS/V oracle 時當檢查點：

1. **全形標點應該顯示得出來。** `END_S13.DAT` 前 408 格不是空白。
   若實跑看到標點是空的，代表符號區的來源判斷有誤。
2. **缺字會顯示成同一個固定字形**（第 `0x56` 格），不是空白也不是方框。

## 8. 對 remake 的影響

| 項目 | 處置 |
|---|---|
| 字高 | 原版全形是 **16×15**，第 16 列固定空白。`internal/assets/cjk` 若用 16×16 字型，行距與原版差一列 |
| 半形 | 8×15，與全形同高，所以中英混排不需要基線調整 |
| 缺字 | 原版回第 `0x56` 格；remake 應明確標記為 remake 差異，不要靜默沿用 |
| 載入 | 原版逐字開檔，remake 一次載入。**這是 remake 差異，記在這裡** |
| 授權 | 倚天字型不隨本專案散布（§10）。remake 走 `-font` 指定玩家自備的字庫 |

## 9. 未解

| 項目 | 現況 |
|---|---|
| `END_S10/S11` 與 `STR.EXE` 檔名不同步 | §6，要實跑裁決 |
| `END_S13.DAT` 前 408 格的來源 | 不是 `stdfont.15` 的任何一段，也不是 `usrfont.15m`（256 B）|
| `END_S15.DAT`（5,242 B）| `KI.EXE` 的字串表引用它，但**不是字型**（大小不是 30／15 的倍數），也不是過場圖（沒有 `00 F4 01` 檔頭），Big5 解碼是亂碼 → 疑似壓縮 |
| `sub_1F7A4` | 把 32 B 緩衝畫上 VRAM 的實際迴圈，未逐行讀 |
| `YNFONT.EXE` 怎麼顯示中文 | 它不走 INT 15h（0 次），防拷畫面的中文是它自己畫的。與本鏈無關，仍未解 |
