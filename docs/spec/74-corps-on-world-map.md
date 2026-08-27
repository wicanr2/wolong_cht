# 74 — 軍團要畫在大地圖上

**狀態：CONFORMED。軍團以 `MMAP.MCH` 的原版圖塊疊在大地圖上，桌面與手機共用同一條算式。**

- 日期：2026-08-23
- 出處：`sub_12AF4` → `sub_12B2A` → `sub_1D4C7` → `sub_1D66A`（本輪讀出來，見 §1）
- 推論等級：**confirmed**（座標、圖塊算式與圖庫來源都由機器碼定案）

## 1. 原版做什麼

四支函式串起來。`sub_12AF4` 掃整張軍團表：

```asm
mov     si, 2240h           ; 軍團表
mov     cx, 7Fh             ; 127 筆
loc_12B18:
cmp     byte ptr [si], 0C0h ; ← 存在旗標 ≥ 0xC0 才畫
jb      short loc_12B20
call    sub_12B2A
loc_12B20:
add     si, 40h             ; 每筆 0x40 bytes
loop    loc_12B18
```

`sub_12B2A` 取座標與圖塊編號：

```asm
mov     dx, [si+10h]        ; X（格）
mov     bl, [si+12h]        ; Y（格）
mov     al, [si+9]          ; 勢力編號 × 5
add     al, [si+8]          ; ＋朝向（0/1 ＝ X 減增、2/3 ＝ Y 減增、4 ＝ 靜止）
call    sub_1D4C7
```

`sub_1D4C7` 把圖塊推進**顯示表**，不是直接畫：

```asm
bx -= cs:word_1D856         ; 減鏡頭 Y；借位就不畫
cmp bx, 17h / jnb 不畫      ; 可視 23 列
dx -= cs:word_1D854         ; 減鏡頭 X；借位就不畫
cmp dx, 28h / jnb 不畫      ; 可視 40 行
bx = (Y*40 + X) * 8         ; 每格 8 bytes
ds = cs:word_1D84E          ; 顯示表段
bl = [si+1]                 ; 這格已經疊了幾張
cmp bl, 4 / jnb 不畫        ; ⭐ 推入端最多 4 張（槽位其實有五個，見 ../re/72 §3）
xchg al, [bx+si+3]          ; 疊到第 bl 層
[si+1] = bl + 1
if 值有變 → or [si], 20h    ; 標記這格要重畫
```

## 2. ⭐ 疊圖的圖庫就是 `MMAP.MCH`

`sub_1D66A` 是顯示表的 blitter，走 23×40 格，**地形與疊圖用不同的段**：

```asm
mov     ah, [si+2]              ; 地形圖塊
mov     ds, cs:word_1D84A       ; ← 地形圖庫
call    sub_1D7E7
mov     bl, 3
loc_1D6F4:
mov     ah, [bx+si]             ; 疊圖第 bl−3 層
mov     ds, cs:word_1D84C       ; ← 疊圖圖庫
call    sub_1D804
inc     bx / dec dh / jnz loc_1D6F4
```

而 `sub_1D46A` 把兩個段串起來：

```asm
mov     cs:word_1D84A, ax
add     ax, 800h                ; ⭐ 疊圖庫 ＝ 地形庫 + 0x800 段 ＝ +32,768 B
mov     cs:word_1D84C, ax
```

**32,768 B 正好是 `MMAP.MDL`**（256 張 16×16 4bpp，`internal/assets/world`
的 `TileCount = 32768 / 128`）。所以緊接在後面的疊圖庫就是 **`MMAP.MCH`**——
remake 早就在解它（首都疊圖 `CapitalOverlayTile = 0xFF` 就是 MCH 圖塊，
而 `world_test.go` 的註解已經寫著「固定 `sub_1D804` 的 256×160 物件圖塊區」）。

⭐ **這一格不需要新素材，也不需要新解碼器。** 缺的只是「把軍團接進去」。

## 3. 演算法

```
for each 軍團 c where c.Alive:
    tile = c.Faction*5 + c.Heading      # 每勢力五張：四方向 ＋ 靜止
    在 (c.X, c.Y) 疊上 MCH.Tile(tile)   # 透明像素不蓋地形
```

⚠ **朝向的值域是 5**（`0/1` X 減增、`2/3` Y 減增、`4` 靜止）——
`internal/state.Corps.Heading` 的註解已經寫死這一條，直接用。

⚠ **不要自己另外算「該畫哪一格」。** 座標就是 `Corps.X`／`Corps.Y`，
與原版 `+0x10`／`+0x12` 同一個欄位（`internal/state/corps.go` 已對齊）。

## 4. remake 實作

| 項目 | 位置 |
|---|---|
| 疊圖來源 | `internal/assets/world`（`MCHTile`，已存在）|
| 疊圖清單 | `cmd/wlgame/strategyhud.go` 的 `corpsMarks()`（新增）|
| 繪製 | `internal/assets/library` 的 `RenderWorldMarked` 增一個疊圖參數 |
| 差異 | 無。座標、圖塊編號與圖庫都照原版 |

### 4.05 ⭐ 軍團站在首都中心格時，**首都疊圖不畫**

原版 `probe-march/e5.png` 的許昌（曹操剛編成一支軍團，軍團站在首都裡）
逐格窮舉過：對那一格試「首都疊圖 開／關」× MCH 圖塊 0–255 共 512 種組合，
**只有一種是 0 px**——

| 組合 | 不同像素 |
|---|---:|
| **據點中心（自勢力）＋ MCH 4，不疊首都** | **0** |
| 據點中心 ＋ 首都疊圖 ＋ MCH 4 | 40 |
| 據點中心 ＋ 首都疊圖，無軍團 | 195 |
| 據點中心，無首都疊圖、無軍團 | 162 |

MCH 4 ＝ `CorpsTile(0, HeadingStill)`，也就是曹操那一支靜止的軍團。
差的那 40 px 全在圖塊的外圈——軍團圖塊在那一圈是透明的，
所以底下是誰就露誰：原版露的是據點中心，remake 露的是首都疊圖。

⚠ **首都疊圖本身沒有問題**：開局主畫面（沒有軍團）五區逐像素相同
（[`../playtest/37`](../playtest/37-main-screen-parity.md)），2026-08-27 重驗仍是 0 px。
兩個都對，錯的是**同一格同時有兩者**時的取捨。

**機制未讀**：`sub_1D4C7` 的顯示表是後推的疊在前面的（[`../re/72`](../re/72-world-map-display-list.md) §3），
照那個模型兩層都該畫。是誰決定不推首都那一層，沒有讀出來——
推論等級因此是**強證據（實機逐像素 0 差）**，不是 confirmed。

### 4.1 驗收路徑：要看得到「軍團在路上」

`-open-form`／`-open-corps`／`-open-march-list` 三個 fixture 都停在**視窗裡**，
而視窗蓋住大地圖——用它們截圖看不到軍團疊在哪。要驗這一條得有一條
**停在大地圖**的路徑：

| 旗標 | 作用 |
|---|---|
| `-corps-on-map` | 編一支軍團（走 `formCandidates` 的真實資格判定），不開任何視窗 |
| `-march-to N` | 配上面：對那支軍團下行軍指示到據點 N。`-shot-frames` 推進的 tick 會讓它沿原版道路走出城 |

⚠ **軍團待在自己城裡時看不出來是正常的。** 軍團圖塊疊在據點中心格上，
而據點中心徽記與軍團圖塊都是紅色系——原版也一樣（`probe-march/e5.png`
的許昌與同一張圖上沒有軍團的城比較，差別只在中心那個凸字形）。
所以驗收要截**兩張**：城裡那一張證明疊得上去，路上那一張證明它會動。

⚠ **推入上限那條先不做。** remake 的疊圖只有據點與軍團兩種來源，
還碰不到上限；照抄一個沒有機會觸發的限制只會讓程式更難讀。
**這是刻意的省略，不是漏掉**——真的要疊到四層時再補。

## 5. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestCorpsTile`：釘死圖塊算式（勢力 × 5 ＋ 朝向）、越界退回靜止、兩個勢力的區間不重疊 |
| 單元測試 | `TestCorpsMarksSkipDeadCorps`（`cmd/wlgame`）：`Alive=false` 不出現在 `corpsMarks()` |
| 實跑 | `-corps-on-map`（城裡）與 `-corps-on-map -march-to N`（路上）兩張截圖（§4.1）|
| 對原版 | `workplace/promo-live/probe-march/e5.png` 的許昌：原版在城中心多一個凸字形（[`../playtest/50`](../playtest/50-corps-on-map.md)）|

## 6. 未解

| 項目 | 現況 |
|---|---|
| 那 110 張圖在 MCH 裡的實際外觀 | 算式定案，但**沒有逐張看過** 22 勢力 × 5 方向長什麼樣。已看過的：勢力 0 的靜止與行進（[`../playtest/50`](../playtest/50-corps-on-map.md)）|
| 每格 4 層上限 | 刻意沒做（§4）|
| `sub_12B3C` | 軍團旗標 `0x20` 成立時呼叫的那一支（推測是擦除舊位置），未讀 |
| 首都疊圖為什麼不畫 | §4.05 的行為由實機逐像素定案，**機制沒讀出來**。可能是同一個疊圖槽被覆寫，也可能是據點那一趟看到軍團就跳過。下手點：`sub_1D4C7` 的呼叫者裡，推首都疊圖的那一支 |
| 別的疊圖組合 | 只驗過「首都 ＋ 軍團」。**災害物件 ＋ 軍團**、**非首都據點 ＋ 軍團**都沒有樣本 |
