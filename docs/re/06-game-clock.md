# 06 — 遊戲時鐘：`sub_11D8E`

**狀態：整條時間鏈已解，confirmed。**

- 日期：2026-08-08
- 輸入：`workplace/ida/dosv/KI.EXE.i64`
- 推論等級：**confirmed**（機器碼直接讀出來，不是推的）

## 0. 位址對映

這支執行檔是單一段（程式碼與資料同段，程式裡到處是 `cs:` 寫資料）。

```
IDA 線性位址 = ds 偏移 + 0x10000
檔案偏移     = ds 偏移 + 0x200        （MZ 標頭 0x200）
```

驗證：`ds:0D9A` 在檔案 `0x0F9A` 讀到 `SINARIO.DAT\0SAVE...` ✅

## 1. 時鐘的欄位

| ds 偏移 | 型別 | 內容 |
|---|---|---|
| `0CF0h` | u8 | **日** |
| `0CF1h` | u8 | **該月天數**（`0CF0h` 讀成 word 時是高位元組） |
| `0CF2h` | u8 | **子刻**（0–8） |
| `0CF3h` | u8 | **時**（1–24） |
| `0CF4h` | u8 | **月** |
| `0CF6h` | u16 | **年** |

`sub_11E17`（每「時」呼叫一次）把三個欄位畫到畫面上：

```asm
mov ax, ds:0CF6h   / mov bx, 903h / mov di, 2BBh / call sub_1062F   ; 年，3 位數
mov al, ds:0CF4h   / mov bx, 902h / mov di, 2C0h / call sub_1062F   ; 月，2 位數
mov al, ds:0CF0h   / mov bx, 902h / mov di, 2C4h / call sub_1062F   ; 日，2 位數
```

`bx` 的低位元組就是位數（3／2／2），`di` 是畫面位置。
**這支程式把年月日三個欄位釘死了**，不需要任何推測。

## 2. `sub_11D8E`：整條進位鏈

```asm
sub_11D8E:
    cmp  byte ptr ds:0CF2h, 8        ; 子刻 < 8 ?
    jb   loc_11DF4                   ;   → 只加子刻
    cmp  byte ptr ds:0CF3h, 17h      ; 時 < 23 ?
    jb   loc_11DE0                   ;   → 進位到「時」
    mov  ax, ds:0CF0h                ; AL=日  AH=該月天數
    cmp  al, ah
    jb   loc_11DD7                   ; 日 < 該月天數 → 進位到「日」
    cmp  byte ptr ds:0CF4h, 0Ch      ; 月 < 12 ?
    jb   loc_11DC1                   ;   → 進位到「月」
    cmp  word ptr ds:0CF6h, 3E8h     ; 年 < 1000 ?
    jb   loc_11DB8
    mov  word ptr ds:0CF6h, 3E6h     ;   年封頂：設 998
loc_11DB8:
    inc  word ptr ds:0CF6h           ; 年++
    mov  byte ptr ds:0CF4h, 0        ; 月 = 0
loc_11DC1:
    inc  byte ptr ds:0CF4h           ; 月++
    mov  bl, ds:0CF4h
    xor  bh, bh
    mov  ah, [bx-6755h]              ; ← 查每月天數表（ds:98ABh）
    xor  al, al
    mov  ds:0CF0h, ax                ; 日 = 0，同時寫入該月天數
    call sub_15358                   ; ★ 月結
loc_11DD7:
    inc  byte ptr ds:0CF0h           ; 日++
    mov  byte ptr ds:0CF3h, 0        ; 時 = 0
loc_11DE0:
    inc  byte ptr ds:0CF3h           ; 時++
    mov  byte ptr ds:0CF2h, 0        ; 子刻 = 0
    call sub_19377                   ; ★ 季節漸變
    call sub_13E11                   ; ★ 每「時」的世界更新（未解）
    call sub_11E17                   ; 重畫日期
    jmp  loc_11DF8
loc_11DF4:
    inc  byte ptr ds:0CF2h           ; 子刻++
loc_11DF8:
    …速度節流（§4）…
```

**每一層都用「落下去再 fall through」的寫法**，所以進位後所有低位單位
都會被重設成 1（不是 0）——這解釋了為什麼劇本檔的「時」欄位存的是 `1`。

### 時間單位的層級

```
子刻 (0–8，9 階) → 時 (1–24，24 階) → 日 → 月 → 年
```

**一個遊戲日 ＝ 24 「時」× 9 子刻 ＝ 216 個 tick。**

「時」這個單位**說明書從頭到尾沒提過**，但它是實際存在的：
季節漸變、世界更新、畫面重繪都掛在「時」上，不是掛在「日」上。

## 3. 每月天數表（ds:98ABh）

```
偏移 98ABh:  00 1f 1c 1f 1e 1f 1e 1f 1f 1e 1f 1e 1f
索引          0   1  2  3  4  5  6  7  8  9 10 11 12
```

| 月 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| 天數 | 31 | **28** | 31 | 30 | 31 | 30 | 31 | 31 | 30 | 31 | 30 | 31 |

**二月固定 28 天，沒有閏年判斷。**

這張表獨立驗證了劇本的起始月份：四個劇本的「該月天數」欄位
（30／30／30／31）分別對應 4／9／6／5 月，全部吻合。
**劇本 2 是 9 月不是說明書寫的 10 月**——10 月是 31 天，對不上。

## 4. 速度節流

```asm
    mov  al, ds:0CFAh          ; 速度設定值
    and  al, al
    jz   locret                ; 0 → 不等待（最速）
loc_11DFF:
    cmp  byte ptr ds:0D2Ch, 0
    jz   loc_11DFF             ; 等計時中斷把旗標拉起來
loc_11E06:
    cmp  ds:0D2Dh, al
    jb   loc_11E06             ; 等計數器 ≥ 速度值
    mov  byte ptr ds:0D2Dh, 0
    mov  byte ptr ds:0D2Ch, 0
```

| 位址 | 內容 |
|---|---|
| `ds:0CFAh` | **戰略速度設定值**。`0` ＝ 不等待 |
| `ds:0D2Ch` | 計時中斷設的旗標 |
| `ds:0D2Dh` | 計時中斷累加的計數器 |

**這正是說明書 3.5 節那句話的機器碼證據**：

> 選択されている速度が機械の速度を超えるといくら速い速度を設定しても
> それ以上速くなりません

速度設定只是**每個 tick 之後至少等幾個計時中斷**。設 0 就完全不等，
於是速度上限由機器跑一輪主迴圈要多久決定。

⭐ **那個「計時中斷」有頻率**：`YNSOUND.COM` 把 PIT 重設成 4660.9 Hz，
自己分頻 16 之後回呼遊戲，也就是 **291.3 Hz**（[`61`](61-timer-tick-source.md)）。
所以「沒有固定 tick rate」只對**最高速**成立——
其餘四檔的等待量子是固定的 3.43 ms，一個遊戲日 0.74 ～ 2.97 秒，
數字見 [`61`](61-timer-tick-source.md) §4。

## 5. `sub_15358`：月結

換月時（`loc_11DC1`）呼叫。

```asm
mov  ds, word_10D52          ; 勢力表所在的段
mov  si, 0 / xor cl, cl
loc_15362:
    cmp  byte ptr [si], 80h  ; 該勢力存在？
    jb   loc_15381
    mov  ax, [si+1Ah]        ; 讀兩個累計欄位
    mov  dl, [si+1Ch]
    call sub_1563B / sub_153C6 / sub_15609
    xor  ax, ax
    mov  [si+1Ah], ax        ; 歸零
    mov  [si+1Ch], al
    call sub_15828
loc_15381:
    add  si, 40h             ; ← 每筆 64 byte
    inc  cl
    cmp  cl, 16h             ; ← 共 22 筆
    jb   loc_15362
    call sub_15695 / sub_1585F / sub_155A6 / sub_12BD9 / sub_15715
    call sub_1578F / sub_122DB / sub_12286 / sub_157FE
    mov  si, 0D10h / mov di, 0D08h / mov cx, 4
loc_153AF:
    mov  ax, cs:[si] / mov cs:[di], ax     ; ★ 4 個 word：來月 → 今月
    …
    loop loc_153AF
    mov  al, 0Eh / call sub_15E80
```

### 兩個直接可用的結論

| 結論 | 證據 |
|---|---|
| **勢力表 ＝ 22 筆 × 64 byte** | 迴圈的 `add si,40h` 與 `cmp cl,16h` |
| 勢力記錄 `+0` ≥ `0x80` 表示該勢力存在 | `cmp byte ptr [si], 80h` |
| 勢力記錄 `+1Ah`(word)、`+1Ch`(byte) 是**月結時歸零的累計值** | 讀完就寫 0 |

### ⭐ 「次月末才生效」的實作就是一次 4 word 的複製

```
cs:0D10h ──複製 4 個 word──> cs:0D08h
```

正好對上說明書 3.2 節：

> ここにセットされた値は**来月末**より使用されます

**4 個 word ＝ 稅率、騎馬募兵數、弓兵募兵數、步兵募兵數。**
`0D10h` 是「來月」緩衝區，`0D08h` 是「今月」生效值。
玩家在財政視窗改的是 `0D10h`，月結時才搬到 `0D08h`。

## 6. `sub_19377`：季節是**漸變**的

```asm
cmp  byte ptr ds:0CF3h, 1     ; 只在「時 == 1」執行
jnz  ret
mov  ah, ds:0CF0h
cmp  ah, 10h
ja   ret                      ; 只在「日 ≤ 16」執行
mov  al, ds:0CF4h
cmp  al, 3  → loc_1939B
cmp  al, 6  → loc_193AA
cmp  al, 9  → loc_193B9
cmp  al, 0Ch → loc_193C8
```

四個分支都是同一個形狀：

| 月 | 來源盤 | 目標盤 | 傳給 `sub_193D7` ＝ **曲號**（[`58`](58-bgm-scene-mapping.md)）|
|---:|---|---|---:|
| 3 | `1934h` | `18A4h` | 2 |
| 6 | `18A4h` | `18D4h` | 3 |
| 9 | `18D4h` | `1904h` | 4 |
| 12 | `1904h` | `1934h` | 5 |

四張調色盤間隔 `0x30` ＝ 48 byte ＝ **16 色 × 3 通道**，
與 `.BRG` 的格式一致（`docs/formats/04`）。

| 位址 | 季節 |
|---|---|
| `18A4h` | 春（3 月轉入） |
| `18D4h` | 夏（6 月轉入） |
| `1904h` | 秋（9 月轉入） |
| `1934h` | 冬（12 月轉入） |

### 漸變是 16 天

條件是「日 ≤ 16 且時 == 1」→ **每天執行一次，連續 16 天**。
`sub_10A65(來源, 目標)` 每次走一步。

**季節不是在某一天突然切換的，是在 3、6、9、12 月的前 16 天慢慢褪過去。**
remake 若做成瞬間切換，畫面感受會完全不同。

⭐ **音樂不跟著漸變**：`sub_193D7` 的 `ah` 是日，日 ＝ 1 時停、
**日 ＝ 2 時整首換掉**（[`58`](58-bgm-scene-mapping.md) §2）。
畫面褪 16 天、音樂第 2 天就換——**兩者不同步是原版行為，不是缺口。**

## 7. 劇本／存檔區塊的完整佈局

`sub_18CAE`（載入）與 `sub_18CFF`（存檔）用同一組位移：

```asm
mul  56C0h                    ; 劇本編號 × 22208 → 檔案位移
mov  bx, cs / mov si, 0CF0h / mov di, 3Bh     ; ① 59 B ← cs:0CF0h
add  ax, 80h
mov  bx, ds:0D52h / xor si,si / mov di, 5240h ; ② 21056 B ← 勢力／據點／武將段
add  ax, 5240h
mov  bx, ds:0D56h / xor si,si / mov di, 400h  ; ③ 1024 B
```

| 檔案位移 | 長度 | 去向 |
|---|---:|---|
| `+0x0000` | 59 | `cs:0CF0h`——**時鐘 ＋ 全域狀態** |
| `+0x0080` | 21,056 | `ds:0D52h` 段——勢力表 ＋ 據點表 ＋ 武將表 |
| `+0x52C0` | 1,024 | `ds:0D56h` 段 |

```
0x80 + 0x5240 + 0x400 = 0x56C0 = 22,208 ✅ 正好是一個劇本區塊
```

**這解開了 `docs/formats/08` 的整體結構。** 而且對得上兩件既有結論：

- `docs/formats/08` 的「`+0x52C0` 起 1,024 B 全零」＝ 第 ③ 塊
- `docs/re/05` 的「據點記錄基址 `0x840`」＝ 第 ② 塊內的偏移
  （檔案 `+0x08C0` − `0x80` ＝ `0x840`）✅

`sub_18CAE` 進入點帶 `bx = 0DA6h`（＝ `SAVE.DAT` 的檔名位址，
緊接在 `SINARIO.DAT` 之後），所以同一支程式靠檔名參數決定讀哪個檔。

## 8. 還沒解的

| 缺口 | 下手點 |
|---|---|
| `sub_13E11`——每「時」的世界更新（AI？行軍？） | 直接讀 |
| `sub_15358` 呼叫的 9 支子程式各做什麼 | 逐一讀，經濟公式都在裡面 |
| 勢力記錄 64 B 的欄位表 | `+0`、`+1Ah`、`+1Ch` 已知，其餘未解 |
| `cs:0CF0h` 那 59 byte 裡除了時鐘的其餘部分 | 對照存檔 diff |
| `sub_10A65` 的內插演算法 | 直接讀 |
| ~~`sub_193D7(al)` 的 2/3/4/5 是什麼索引~~ | **已解**：`BGM.DAT` 的曲號，四季配樂（[`58`](58-bgm-scene-mapping.md) §2）|
| ~~戰術速度存在哪（`0CFAh` 是戰略速度）~~ | **已解**：`ds:0CFBh`，用之前先 ×16 存進 `ds:0CFCh`（[`61`](61-timer-tick-source.md) §4）|
