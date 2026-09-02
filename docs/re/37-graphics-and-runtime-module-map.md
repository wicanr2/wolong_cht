# 37 — 圖庫、繪圖底層與 C runtime 兩個模組的全圖

**狀態：叢集歸屬 confirmed。硬體層的角色由 I/O 埠直接判定（精確），
其餘角色標籤是強證據。**

- 日期：2026-08-14
- 範圍：只驗松崗 DOS/V，位址 `1E800`–`21000` 與 `10000`–`11800`
- 原始 `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 原始 `KI.EXE.i64` SHA-256：`7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`
- 工具：**IDAPython** `tools/ida_module.py`
- 位址空間：IDA DOS/V linear address，segment base `0x10000`

方法與界線同 [`35`](35-strategy-ui-module-map.md) §1。
**這兩個模組有一個好處：它們碰硬體，而 I/O 埠是精確證據**——
不必靠呼叫圖推，`out dx, al` 打到哪個埠就決定了那支在做什麼。

## 1. 硬體層：埠決定角色

| 函式 | B | 埠／中斷 | 角色 |
|---|---:|---|---|
| `sub_1EB11` | 77 | **PPI port 61h**（in／out ×2）| **PC 喇叭發聲**。呼叫者 `sub_10CDE`／`sub_10CE7`（各 9 B，合計 10 個呼叫點）|
| `sub_1EB5E` | 14 | Video status | 等垂直空白 |
| `sub_1EB6C` | 39 | **INT 10h** ＋ EGA palette／overscan | **設畫面模式與調色盤**，開機時由 `sub_1006B` 呼叫一次 |
| `sub_1EB93` | 73 | EGA palette ＋ Graphics Controller ×4 | 調色盤重設；`sub_1006B`、`sub_11A6E` |
| `sub_1EBDC` | 82 | 四處 `out`（IDA 未標註埠）| 呼叫者含 `sub_19946`（戰術主迴圈）→ 疑似音源 |
| `sub_1EC2E`／`sub_1EC6C` | 62／22 | 同上 | 由 `sub_10A65` 進入 |
| `sub_1EC82` | 94 | **INT 1Ah**（BIOS 計時器）| **亂數種子**（[`10`](10-rng.md) 已定案）|
| `sub_1ECE0` | 28 | — | **亂數**（[`10`](10-rng.md)），全庫最常被呼叫的常式之一 |

**`sub_1EB11` 這一組是 PC 喇叭而不是音效卡**——
port 61h 是 PC/XT 的 PPI，控制 speaker gate。DOS/V 側的音源另由
`YNSOUND.COM` 提供（[`17`](17-dosv-audio-tsr.md)）,兩者不是同一條路。

## 2. 圖庫解碼與 blit

| 函式 | B | 呼叫者 | 角色 |
|---|---:|---|---|
| `sub_1E81C` | 325 | `sub_1E77D` | **全庫最大的函式**。圖庫解碼主體 |
| `sub_1F020` | 288 | `sub_1C61F`、`sub_1C6BF`（戰術）| 大圖塊繪製，八處 EGA `out` |
| `sub_1F1A3` | 203 | `sub_10BCD`、`sub_15AFC`、`sub_18607` | [`21`](21-function-census.md) T4 榜首。六處 EGA `out` |
| `sub_1F888` | 176 | `sub_19B6D`、`sub_1C673`、`sub_1C7F…` | 位元對齊 blit（[`18`](18-tactical-button-glyphs.md) §2）|
| `sub_1FA37` | 107 | 15＋ 處 | 圖塊 blit（`ax` ＝ 尺寸），共用層 |
| `sub_1F9B0` | 107 | `sub_10C60`、`sub_10C77`、`sub_16DA…` | — |
| `sub_1FAC2` | 79 | `sub_1950F`、`sub_19796` | — |
| `sub_1FB29`／`sub_1FBA7` | 126／84 | **無直接呼叫者** | 間接分派目標 |
| `sub_1F7A4` | 212 | **無直接呼叫者** | **字型 blitter**——它是被 `loc_1F75E` 用 `call far` 呼叫的（[`29`](29-font-service-int15.md) §3），所以直接 xref 看不到 |

`sub_1F7A4` 是「零呼叫者不等於死碼」最乾淨的例子：
[`21`](21-function-census.md) §6 講的三種可能裡，它屬於「被取址後間接呼叫」，
而那個位址是**開機時由 `INT 15h` 服務填進去的**。

## 3. 檔案 I/O 與資源載入

`sub_1F4DF`（73 B）、`sub_1F593`（84 B）、`sub_1F655`（60 B）分別由
`sub_1E38C`、`sub_1E3A6`、`sub_1F5E7` 呼叫，三支都含 `mov ax, 4200h`＋`INT 21h`
（DOS LSEEK）。`sub_1E38C` 是圖庫載入（[`33`](33-shared-draw-helpers.md) §2 用它讀肖像）。

> ⚠ **這三支一度被工具標成「碰城兵臨時軍團表」**——因為城兵臨時軍團的段內基址
> 正好是 `0x4200`，而 `mov ax, 4200h` 是 DOS 的 LSEEK 功能碼。
> **基址與 DOS 功能碼撞號**，三支工具的 `BASES` 已把 `0x4200` 移除。
> 規則：**立即值當證據前，先問「這個數字還可能是什麼」。**

## 4. `sub_20000`：系統服務分派（已解，見 [`16`](16-idle-clock-event10.md)）

14 B，**呼叫點遍布全庫**。底下六支：

```
sub_20000 → nullsub_4(1) / sub_2002E(65) / sub_20070(42) / sub_2009A(35)
          / sub_200BD(3) / sub_200C0(54)
三支共用 sub_20137(52) → sub_201C6 / sub_201E4 → sub_2020C
                          / sub_20249 → sub_202A0
                          / sub_202BD → sub_202FE
```

`sub_2016B`（50 B）與 `sub_2019D`（41 B）無直接呼叫者——同層的間接目標。

## 5. C runtime 與低階 I/O（`10000`–`11800`）

`start`（67 B）呼叫七支完成初始化：

| 函式 | B | 角色 |
|---|---:|---|
| `sub_1005B` | 15 | — |
| `sub_1006B` | 116 | **硬體初始化**：`sub_1EB6C` 設模式、`sub_1EB93` 調色盤、`sub_1EC82` 亂數種子、`sub_1F720` 字型服務（[`29`](29-font-service-int15.md) §2）|
| `sub_100DF` | 118 | 寫 `word_10D52` 等全域段暫存 |
| `sub_10155` | 95 | 系統服務 4／5／7／8／D／E／F |
| `sub_10210`／`sub_1030F`／`sub_1033B` | 38／40／19 | — |

> ⚠ IDA 把 `start` 的呼叫者標成 `sub_18FC9`（存檔畫面）。
> **那是把跳進 `start` 區段的分支當成呼叫**——`start` 是程式進入點，
> 不會被存檔畫面呼叫。直接 xref 在進入點附近不可信。

其餘可分成四叢：

| 叢 | 函式 | 角色 |
|---|---|---|
| 文字與數字 | `sub_1062F`／`sub_1069A`／`sub_106DE` | 印數字（[`28`](28-text-number-rendering.md) §2）|
| | `sub_106F5`／`sub_106F9`／`sub_106FD` | 各 4 B 的薄包裝（[`33`](33-shared-draw-helpers.md) §1）|
| | `sub_1075B`（119）／`sub_1084A`（90）| 排版繪訊息、×8 變體展開（[`25`](25-message-variants-and-personnel.md) §1）|
| | `sub_107D2`（115）| 肖像四格快取（[`33`](33-shared-draw-helpers.md) §2）|
| 檔案／資源 | `sub_1036F`（84）→ `sub_103C3`／`sub_103E6`／`sub_10414`（161）／`sub_104B5`／`sub_104FF` | `sub_10414` 底下再掛 `sub_1054D`／`sub_1059B`（各 78／82 B，共用 `nullsub_5`）|
| VRAM 搬移 | `sub_109D0`／`sub_10A1C`／`sub_10A65`／`sub_10AAA`／`sub_10AD9`／`sub_10B46`／`sub_10BAF`／`sub_10BCD`／`sub_10C14`／`sub_10C60`／`sub_10C77`／`sub_10CAC`／`sub_10CC3` | 被繪圖層與 EGA 底層共用 |
| 發聲 | `sub_10CDE`／`sub_10CE7`（各 9 B）| 轉呼叫 `sub_1EB11`（PC 喇叭）|

`sub_10337`（4 B，20＋ 個呼叫點）是這個模組裡被呼叫最多的薄包裝。

## 6. 仍要逐行讀的

| 函式 | B | 為什麼 |
|---|---:|---|
| `sub_1E81C` | 325 | **全庫最大**，圖庫解碼主體，只有一個呼叫者 |
| `sub_1F020` | 288 | 大圖塊繪製，戰術側專用 |
| `sub_1F1A3` | 203 | T4 榜首，三個來源不同的呼叫者 |
| `sub_1F7A4` `[DOS/BIOS]` | 212 | 字型 blitter，逐行未解。**同一支函式在 [`29`](29-font-service-int15.md) §9 也列著，那裡是正本**；未解的是「怎麼寫 VRAM」，而 remake 不碰 VGA 平面，所以不擋 remake |
| `sub_1EBDC`／`sub_1EC2E` | 82／62 | 埠沒被 IDA 標註，是不是音源要驗 |
| `sub_10414` 叢 | 161＋ | `sub_103E6` 的 `call cs:word_1036D` 是 [`21`](21-function-census.md) §7 未攤開的間接分派 |
