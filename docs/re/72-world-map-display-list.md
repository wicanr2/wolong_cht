# 72 — 大地圖的顯示表：地形一層 ＋ 最多四層疊圖

**狀態：已解。** 大地圖不是直接畫的，中間隔著一張 40×23 的顯示表；
軍團就是疊在那張表上的一層，圖庫是 `MMAP.MCH`。

- 日期：2026-08-23
- 輸入檔：`workplace/orig/dosv/KI.EXE`
  SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 資料庫：`KI.EXE.i64` SHA-256 `7885fbce1fd54f318fa04ff355755514e76030af70acb5e11645b02ed7483bed`（739 支函式）
- 推論等級：**confirmed**
- 規格：[`../spec/74`](../spec/74-corps-on-world-map.md)

## 1. 四支串起來

| 位址 | 做什麼 |
|---|---|
| `sub_12AF4` | 掃整張軍團表：`si=2240h`、每筆 `0x40`、127 筆 |
| `sub_12B2A` | 取一支軍團的座標與圖塊編號 |
| `sub_1D4C7` | 把圖塊**推進顯示表**（不是畫）|
| `sub_1D66A` | 顯示表的 blitter，逐格把地形與疊圖畫出來 |

## 2. `sub_12AF4`：存在旗標是 `>= 0xC0`

```asm
mov     si, 2240h
mov     cx, 7Fh             ; 127 筆
loc_12B18:
cmp     byte ptr [si], 0C0h ; ← 軍團記錄 +0x00
jb      short loc_12B20     ; 小於就不畫
call    sub_12B2A
loc_12B20:
add     si, 40h
loop    loc_12B18
```

前面還有一輪掃 `test byte [si], 20h` → `sub_12B3C`（未讀，推測是擦除舊位置）。

## 3. `sub_1D4C7`：顯示表的結構

```asm
bx -= cs:word_1D856         ; 減鏡頭 Y，借位就不畫
cmp bx, 17h / jnb 不畫      ; 可視 23 列
dx -= cs:word_1D854         ; 減鏡頭 X，借位就不畫
cmp dx, 28h / jnb 不畫      ; 可視 40 行
bx = (Y*40 + X) * 8         ; 每格 8 bytes
ds = cs:word_1D84E          ; 顯示表段
bl = [si+1]                 ; 這格疊了幾張
cmp bl, 4 / jnb 不畫        ; ⭐ 每格最多 4 層
xchg al, [bx+si+3]          ; 疊到第 bl 層
[si+1] = bl + 1
若值有變 → or [si], 20h     ; 標記重畫
```

每格 8 bytes 的版面：`+0` 旗標、`+1` 疊了幾層、`+2` 地形圖塊、`+3..+6` 疊圖四層。

## 4. ⭐ 兩個圖庫段差 `0x800`

`sub_1D66A` 畫的時候地形與疊圖用**不同的段**：

```asm
mov     ah, [si+2]              ; 地形
mov     ds, cs:word_1D84A       ; ← 地形圖庫
call    sub_1D7E7
mov     bl, 3
loc_1D6F4:
mov     ah, [bx+si]             ; 疊圖
mov     ds, cs:word_1D84C       ; ← 疊圖圖庫
call    sub_1D804
inc bx / dec dh / jnz loc_1D6F4
```

而 `sub_1D46A` 把兩者串起來：

```asm
mov     cs:word_1D84A, ax
add     ax, 800h                ; ＝ +32,768 byte
mov     cs:word_1D84C, ax
```

**32,768 byte 正好是 `MMAP.MDL`**（256 張 16×16 4bpp）。
所以疊圖庫就是緊接其後的 **`MMAP.MCH`**——`internal/assets/world`
早就在解它，首都疊圖 `CapitalOverlayTile = 0xFF` 用的就是這一組。

## 5. 找法值得記

一開始只知道 `sub_12B2A` 呼叫 `sub_1D4C7`，而 `sub_1D4C7` 標「未讀」。
**沒有直接去猜圖塊從哪來**，而是照著資料流走：

1. dump `sub_1D4C7` → 看到它只寫顯示表，不畫 → **畫的人是別人**
2. xref `word_1D84E`（顯示表段）→ 找到消費端 `sub_1D66A`
3. `sub_1D66A` 用了兩個不同的段 → xref 那兩個 → `sub_1D46A` 一行 `add ax, 800h` 定案

⭐ **三次 xref，零次猜測。** 對照組是先前那條走不通的路：
數了 `ICONGRF` 段 1 的長度 16,128 ÷ 128 ＝ 126 塊，而軍團要 110 張，
覺得「數字對得上」——**那是統計特徵不是語意**（`CLAUDE.md` §7 第 17 條），
真正的答案在資料流裡。

## 6. 未解

| 項目 | 現況 |
|---|---|
| `sub_12B3C` | 軍團旗標 `0x20` 成立時呼叫，推測是擦除舊位置，未讀 |
| `sub_1D782`／`sub_1D7E7`／`sub_1D804` | 三支實際搬像素的常式，未讀 |
| 那 110 張軍團圖的逐張外觀 | 算式定案、抽驗過勢力 0 靜止那一張（紅色軍旗），**22 × 5 沒有逐張看過** |
