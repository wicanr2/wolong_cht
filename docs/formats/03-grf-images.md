# 03 — `*GRF.DAT` 圖庫格式

**狀態：`KAOGRF`／`KYOGRF`／`IVENTGRF` READY，`ICONGRF` 部分解。**

- 日期：2026-08-07
- 出處：`docs/re/03-image-blitter.md`（`KI.EXE` 的載入器與繪製常式）
- 工具：`tools/grf.py`
- 驗收：150 張全部解出、**餘 0 byte**、視覺正確（`docs/images/kaogrf-sheet.png`）

## 1. 像素格式

**4 bpp planar，plane-major。** 一張圖：

```
plane 0  整張      width/8 × height byte
plane 1  整張      同上
plane 2  整張
plane 3  整張
```

像素值 = `p0 | p1<<1 | p2<<2 | p3<<3`，每個 byte 的**最高位在最左**。

沒有檔頭、沒有壓縮（熵 3.4–5.2 bit/byte），檔案就是一張接一張。

> **為什麼是 plane-major 而不是逐列交錯**：繪製常式 `sub_1FA37` 對四個平面
> 各呼叫一次 `sub_1FAA2`，而 `sub_1FAA2` 不重設來源指標 `si`——
> 四次呼叫連續吃掉四段資料。見 `docs/re/03` §2。

## 2. `KAOGRF.DAT`：武將頭像（confirmed）

| 項目 | 值 | 依據 |
|---|---|---|
| 尺寸 | **64 × 64** | 繪製常式跑 64 列 × 每列 8 byte |
| 每張 | **2,048 byte** | 載入器 `mov di, 800h` |
| 定位 | **編號 × 2048** | 載入器的 `shl ax,1 / rcl cx,1` ×3 |
| 張數 | **150** | 307,200 ÷ 2,048，餘 0 |
| 調色盤 | `GAMEPAL.BRG` 第 0 組 | 同一支常式先載入 `GAMEPAL` |

**150 張對上遊戲的 146 名武將**（社群資料），多出來的四格是空白或備用。

![KAOGRF 150 張頭像](../images/kaogrf-sheet.png)

### 載入器有 4 格快取

`sub_107D2` 維護一張 4 格的快取表（`[bx+846h]`），未命中時以 round-robin
挑一格（`byte_10845`）覆寫。**remake 不需要照做**——那是為了省 1994 年的記憶體，
不影響行為。但要記錄，那是保存的一部分。

## 3. `KYOGRF.DAT`：15 張 96×96（confirmed）

`sub_17F1A`：

| 依據 | 值 |
|---|---|
| `mov dx, 1200h` ＋ `mul dx` | 定位 ＝ 編號 × 4,608 |
| `mov di, 1200h` | 每張 4,608 byte |
| `mov ax, 6006h` | `ah`＝96 列、`al`＝6 → 每列 12 byte ＝ **96 px 寬** |

驗算 12 × 96 × 4 ＝ 4,608 ✓。69,120 ÷ 4,608 ＝ **15 張**，餘 0。
目的 VRAM 偏移 `bx = 5A02h`。

![KYOGRF 15 張據點景觀](../images/kyogrf-sheet.png)

畫的是**據點景觀**：城門、關隘、宮殿、山道、雪山、水岸、營寨。
`KYO` 應該是「拠点（kyoten）」的頭一個音（**假說**，未驗）。

索引來源是 `[si+16h] >> 4`——某筆記錄的 `+0x16` 欄位取高四位。
**`si` 指向哪一種記錄還沒確認**；圖既然是據點景觀，假說是據點記錄，
而 `>> 4` 取高四位、15 張圖，值域正好吻合 1–15。
確認之後這是一條直通據點資料結構的線。

## 4. `IVENTGRF.DAT`：3 張 288×176（confirmed）

`sub_13D09` 載入、`sub_13D68` 繪製：

| 依據 | 值 |
|---|---|
| `mov dx, 6300h` ＋ `mul dx` | 定位 ＝ 編號 × 25,344 |
| `mov di, 6300h` | 每張 25,344 byte |
| `mov ax, 0B012h` | `ah`＝176 列、`al`＝0x12＝18 → 每列 36 byte ＝ **288 px 寬** |

驗算 36 × 176 × 4 ＝ 25,344 ✓。76,032 ÷ 25,344 ＝ **3 張**，餘 0。
目的 VRAM 偏移 `bx = 2A87h`。

![IVENTGRF 3 張事件圖](../images/iventgrf-sheet.png)

三張都是**劇情過場**：朝堂、密談、燭光夜話。`IVENT` 是 `EVENT` 的拼法。

## 5. `ICONGRF.DAT`：四段組合檔（部分解）

**不是圖庫，是四段不同用途的資料拼起來的。** 四段是分別載入的，
長度加起來 **剛好等於檔案大小 47,776**：

| 段 | 位移 | 長度 | 載入處 | 內容 |
|---|---|---:|---|---|
| 0 | `0x0000` | `0x2800` (10,240) | `sub_1006B` | **640×32 標題橫幅**（`ax=2028h`，`bx=0` → 畫面最頂端）✅ |
| 1 | `0x2800` | `0x3F00` (16,128) | `sub_1C7A9` | 混合尺寸的圖塊。已知有 16×16（`ax=1001h`，128 B／塊）與 16×8（`ax=801h`，64 B／塊）兩種 |
| 2 | `0x6700` | `0x3000` (12,288) | `sub_1006B` | **192×128 縮小地圖底圖**（`ax=800Ch`）✅ |
| 3 | `0x9700` | `0x23A0` (9,120) | `sub_1006B` | 走 `sub_1F888`（**位元對齊**的繪製常式，可放在非 8 倍數的 x）。未解 |

### 5.1 段 0：標題橫幅上是**日文**

![ICONGRF 段 0 標題橫幅](../images/icongrf-r0-banner.png)

橫幅上寫的是「臥竜伝」，不是「臥龍傳」。而 `ICONGRF.DAT`
**兩版 byte-for-byte 完全相同**（`docs/re/01`）——
**松崗版根本沒有重繪這張橫幅。**

這是一個實質的中文化缺口，remake 要補（見 `docs/reference/03`）。

### 5.2 段 2：縮小地圖

![ICONGRF 段 2 縮小地圖](../images/icongrf-r2.png)

與日文說明書第 3.4 節「縮小マップウインドウ」對得上：
可直接點擊移動視角，紅＝自勢力、藍＝敵勢力、黑白＝中立。

`sub_19541` 以像素座標從這一段取任意矩形（`si = bx × 24` → 列跨距 24 byte
＝ 192 寬，與段 2 的尺寸一致），代表它是被當成**來源圖集**用的。

### 5.3 段 1、段 3 未解

段 1 的圖塊不是從段首開始密排——`sub_1CA3B` 用的 `si` 是 `3E00h`／`3E80h`／`3DC0h`
這類段內偏移，中間還有別的東西。**要逐一追每個 `si` 常數的來源**，
不能整段當成同尺寸的圖塊陣列（試過，出來是雜訊）。

段 3 走的 `sub_1F888` 是另一支繪製常式，還沒讀。

## 6. 工具

```sh
# 全部排成總覽
python3 tools/grf.py sheet workplace/orig/dosv/KAOGRF.DAT 64 64 \
        workplace/orig/dosv/GAMEPAL.BRG 0 out.png 15

# 單張放大
python3 tools/grf.py one workplace/orig/dosv/KAOGRF.DAT 64 64 \
        workplace/orig/dosv/GAMEPAL.BRG 0 42 one.png 6
```

`ICONGRF` 這種組合檔用 `region`：

```sh
python3 tools/grf.py region workplace/orig/dosv/ICONGRF.DAT 0x6700 192 128 \
        workplace/orig/dosv/GAMEPAL.BRG 0 out.png 1
```

`sheet`／`region` 都會印出「餘 N byte」。**餘數不是 0 就代表尺寸猜錯了**，
是最便宜的檢核。四個圖庫的尺寸全部餘 0。

## 7. 兩版共用

`KAOGRF`／`KYOGRF`／`ICONGRF`／`IVENTGRF` **四個檔在 PC-98 與松崗版 byte-for-byte
完全相同**（`docs/re/01`）。解碼器兩版通用。
