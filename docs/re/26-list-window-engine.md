# 26 — 一覽表視窗引擎

**狀態：視窗幾何、五個一覽表家族的描述子、選取迴圈與持久化排序狀態 confirmed。
四個 callback 各自的內部行為（怎麼填列、怎麼畫列）未讀。**

- 日期：2026-08-13
- 範圍：只驗松崗 DOS/V
- 原始 `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 原始 `KI.EXE.i64` SHA-256：`7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`
- 工具：IDA Pro 9.4，`tools/ida_listwin.idc`
- 位址空間：IDA DOS/V linear address，segment base `0x10000`

## 1. 一個引擎，五種清單

戰略層的每個指令都要先選東西——選據點、選武將、選勢力、選軍團。
這些一覽表**共用同一個引擎**，差別只在四個 callback 指標。

```
sub_181C0(dx=x, bx=y, cx=尺寸)      開視窗
sub_1820E(ax, bx, si, di, cl=狀態)  跑選取迴圈，回傳 bx＝選中項；CF=1 取消
```

八個呼叫端全部用同一組視窗幾何 `dx=0x18, bx=0x58, cx=0x0B18`。

## 2. `sub_181C0`：視窗幾何

```asm
mov word_181B6, cx
mov ax, 4 / xchg al, cl / shl ax, cl    ; ax = cl_原值 << 4
mov word_181B2, ax                      ; 寬（像素）
xor ax, ax / mov al, ch / shl ax, cl
mov word_181B4, ax                      ; 高（像素）
mov word_181AE, dx / mov word_181B0, bx ; 左上角
sub dx, 8 / sub bx, 8 / shr dx, cl / shr bx, cl
```

`cx` 把尺寸打包成 **16 像素為單位的格數**：`cl` ＝ 寬幾格、`ch` ＝ 高幾格。
一覽表用的 `cx = 0x0B18` 就是 **24 格寬 × 11 格高 ＝ 384 × 176 像素**，
左上角在 `(0x18, 0x58)` ＝ `(24, 88)`。

640×400 的畫面上，這個視窗蓋住 x 24–408、y 88–264。
x 的範圍與指令列（[`22`](22-strategy-command-tree.md) §3，24–408）**完全相同**。

## 3. `sub_1820E`：選取迴圈

```asm
sub sp, 200h / mov bp, sp          ; 512 B ＝ 256 個 word 的結果陣列
mov cs:word_181A8, ax              ; ← ax 是函式指標
mov cs:word_181A6, bx              ; ← bx 是函式指標
mov cs:word_183D3, si
mov cs:word_18320, di / mov cs:word_18591, di
mov cs:byte_181BC, 0
xor ch, ch / shl cx, 1
mov cs:word_181BA, cx              ; 進來的狀態 ×2 ＝ 目前選取列的位移
es=ss / di=bp / cx=100h / ax=FFFFh / rep stosw   ; 陣列填 0xFFFF
call cs:word_181A8                 ; ← 建清單：把項目編號寫進 SS:[bp] 陣列
mov cs:byte_181BD, ah
mov al, 3Dh / call loc_1828F
mov bx, cs:word_181BA / call loc_1857F
call sub_18412                     ; ← 熱區分派迴圈
jb → 取消（CF=1）
mov bl, al / bx = al × 2 / bp += bx
mov bx, [bp+0]                     ; ← 從陣列取出選中項
clc
```

三件事定案：

1. **`ax` 與 `bx` 是函式指標，不是資料。** 呼叫端把自己的 callback 位址傳進來，
   引擎存進 `word_181A8`／`word_181A6` 再間接呼叫。
2. **清單最多 256 項**（512 B ÷ 2），未使用的格子是 `0xFFFF`。
3. **回傳的是陣列裡的值**，不是列號——所以 callback 寫進去的就是呼叫端要的編號。

`sub_18412`（[`22`](22-strategy-command-tree.md) §2 的熱區分派）在這裡當選取迴圈用：
它讀游標下的熱區編號減 `0x3D` 再查 `funcs_18450`，五個 handler 對應
捲動與選取，超過的走 `loc_1857F`。

### 3.1 訂正：`word_181A6`／`word_181A8` 靜態可解

[`21-function-census.md`](21-function-census.md) §7 把
`call cs:word_181A6` 這一組列為「開機時指向 return stub，實際目標由 runtime 填，
要靠動態 trace」。**靜態影像確實是 `07 1F C3 90`（`pop es; pop ds; retn; nop`），
但目標完全靜態可解**——八個呼叫端都在 `sub_1820E` 進入時用立即值傳進來。

教訓寫成規則：**看到「全域函式指標」不要先判定要動態 trace。
先找誰寫它**——寫入端如果是立即值，候選集合就是有限且已知的。

## 4. 五個一覽表家族

八個呼叫端，五組 `(bx, si, di)`。同一組的兩個變體只有 `ax`（建清單的 callback）不同，
也就是**同一種清單的兩種取法**。

| 家族 | `bx` 畫列 | `si` | `di` | 呼叫端 | `ax` 建清單 | 用在哪 |
|---|---|---|---|---|---|---|
| 軍團 | `0x724D` | `0x70EC` | `0x713D` | `sub_1716D` | `0x71A8` | 軍團指令：位置確認／行軍指示 |
| | | | | `sub_171D3` | `0x7217` | `sub_11E46` |
| 據點 | `0x745F` | `0x7378` | `0x73CE` | `sub_17400` | `0x743B` | 據點一覽、遷都、內政官任免 |
| 武將 | `0x76DC` | `0x7550` | `0x75C8` | `sub_17663` | `0x76A0` | 選武將（任命、編成）|
| | | | | `sub_175FA` | `0x763C` | 指令列 #6 武將：看能力 |
| 勢力 | `0x796C` | `0x77EE` | `0x7875` | `sub_17906` | `0x7944` | 外交、外交官任免 |
| | | | | `sub_178A7` | `0x78E5` | 指令列 #7 勢力：看交友度 |
| 開局選勢力 | `0x7B90` | `0x7AC6` | `0x7B12` | `sub_17B3C` | `0x7B6F` | `sub_11AC3` 新遊戲選玩家勢力 |

「同一種清單、兩種取法」與
[`../mechanics/10-strategy.md`](../mechanics/10-strategy.md) §3 的說明對得上：
指令列的「武將」「勢力」是**看**（顯示表上沒列的參數、交友度），
人事與外交的一覽表是**選**。兩者欄位不同，所以建清單的 callback 不同、
畫列的 callback 相同。

## 5. 排序狀態是持久的

每個家族有自己的狀態 byte，進 `sub_1820E` 前讀出、回來後寫回：

| 家族 | 狀態 byte |
|---|---|
| 軍團 | `cs:word_198A8` 低位元組 |
| 據點 | `cs:word_198A8` 高位元組 |
| 武將 | `cs:word_198AA` 低位元組 |
| 勢力 | `cs:word_198AA` 高位元組 |

這正是 [`../mechanics/10-strategy.md`](../mechanics/10-strategy.md) §5 由說明書記下的
「排序狀態以**視窗種類**為單位記住，不必每次重排」——現在有機器碼側的對照，
而且知道它存在哪、只有一個 byte。

`sub_17663`（選武將）進來前先 `xor cl, cl` 再讀，所以**它每次都從第一列開始**；
其餘七個直接沿用上次的位置。

## 6. 解任：`sub_16B4F` 與 `sub_16C2A`

兩支完全同構，只差目標欄位（據點 `+0x19` ／ 勢力 `+0x2A`）：

```asm
mov bh, 0FFh
xchg bh, [si+19h]        ; 欄位 ← 0xFF，bh ← 原本的武將編號
cmp bh, 0FFh / jnz +2
  stc / retn             ; 本來就沒人 → 失敗
xor bl, bl               ; bx = 編號 × 256
shr bx,1 / shr bx,1 / shr bx,1
add bx, 4240h            ; ÷8 → 編號 × 32，加基址 = 武將記錄
mov byte [bx+1Ah], 0     ; 經費餘額 ← 0
mov byte [bx+17h], 0     ; 職務 ← 無
clc
```

`xchg` 一步同時「清空欄位」與「取出原值」，所以**沒有中間狀態**。

**解任會沒收經費餘額。** `+0x1A` 是官員的經費餘額
（[`../formats/08-sinario-save.md`](../formats/08-sinario-save.md) §3，撥款金額 ÷ 128），
解任時直接歸零——先撥款再解任，那筆錢就消失了。
這是可玩性上的實質規則，不是清理程式碼。

## 7. 未解

| 項目 | 現況 |
|---|---|
| 五組 callback 的內部 | `0x70EC`–`0x7B90` 這 17 個位址都沒有函式邊界，要先讓 IDA 認出它們 |
| `si`／`di` 兩個指標的用途 | 存進 `word_183D3`／`word_18320`／`word_18591`，呼叫點還沒找 |
| `funcs_18450` 五個 handler | 捲動與選取的實際行為（`nullsub_2`／`sub_18463`／`sub_184DD`／`sub_1851A`／`sub_18546`）|
| `byte_181BD` | `sub_1820E` 存 `ah`（疑似筆數），讀取端未找 |
| `sub_11D46` | 17 個呼叫點，人事四支離開前都呼叫，未讀 |
