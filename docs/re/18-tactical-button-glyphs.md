# 18 — DOS/V 戰術底列按鈕 glyph 候選資產研究

**狀態：六個命令 glyph、底板、右欄複合面板與選取矩形已解出並接入 remake。**

> 2026-08-12 勘誤：本文早期把 `sub_1F888` 解讀為每次重複
> `segment1+0x3900` 的同一張圖，因而把六 glyph 來源標為 unknown。
> 新證據直接追蹤 `sub_1F938` 對 `SI` 的間接前進：每列加 3 bytes，
> 16 列、四平面共前進 `0xC0`；`sub_1F888` 本身沒有保留 `SI`。
> 因此 `sub_1C7F4` 六次呼叫實際連續消費 `0x3900..0x3D7F`。

- 日期：2026-08-12
- 平台：只研究松崗 DOS/V；PC-98 不在本輪範圍
- 研究工具：IDA Pro 9.4，映像 `ida-pro-9.4-ver2:uidfix-v1`
- IDA 位址基準：DOS/V 線性位址；所有 `segment1+0xNNNN` 是 `ICONGRF.DAT` 第 1 段內相對位移
- 原始 `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 原始 `KI.EXE.i64` SHA-256：`7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`
- 原始 `ICONGRF.DAT` SHA-256：`2154782c045b898aa5fafa74a4ff0c3745771ec85799b7882e1c4c009b3f1c1d`
- 原始 `GAMEPAL.BRG` SHA-256：`1f0119c75ea5cd333bd3ac75ef92030f93924f011728edc9b1f727c483263708`
- 原版參考圖：`yt-wolong-natural-320s.png` SHA-256 `b9816eed9502aab54639207a3c3f3ffc978011cb614505584d6e24b67d85da88`；`yt-wolong-natural-400s.png` SHA-256 `3de7ceca52329b89dee72f08066b0021c0e82d91c392d411aacb11107372d386`

## 1. 本輪可重現產物

接觸表由唯讀研究工具產生，沒有把任何候選接到遊戲：

- [完整 blit 候選接觸表](../images/icongrf-seg1-blit-candidates.png)
- [尾端 16×8 滑動窗口](../images/icongrf-seg1-tail-16x8-windows.png)
- [尾端 16×16 候選](../images/icongrf-seg1-tail-16x16-candidates.png)
- 工具：[icon_segment1_contact.py](../../tools/icon_segment1_contact.py)

產物 SHA-256：完整候選 `609bebcc52f215ed7921bf1320f95b373ef1df52094016869e381a406e108491`；
16×8 窗口 `4c63fce708b7f9522af3526e14d70e5ec07ef81bc65b61cd0e50679a0591247a`；
16×16 候選 `129f110553a86f0d502ebe883ac1e338c2ae4c78e950f0560486de0500f4fcfc`。

工具從 `ICONGRF.DAT` file offset `0x2800`、length `0x3F00` 擷取 segment 1，依
IDA blit 的 planar stride／列數解碼；每個圖塊上方直接標示 segment-relative raw offset。
它不是 runtime decoder，也沒有賦予圖塊 UI 語意。

## 2. `sub_1C7A9` 呼叫鏈與候選尺寸

### 2.1 `sub_1C7F4`：不是 25 次迴圈

IDA Pro 9.4 的 `sub_1C7A9` 範圍為 `0001C7A9–0001C7F4`，在
`0001C7DB` 呼叫 `sub_1C7F4`。`sub_1C7F4` 範圍為 `0001C7F4–0001C83E`：

```asm
0001C7F4  xor dx, dx
0001C7F6  mov bp, 0D2E4h
0001C7F9  mov di, 7300h
0001C7FC  push dx
0001C7FE  mov bx, 170h
0001C800  mov cx, 40Ah
0001C803  mov al, cs:[bp+0]
0001C807  add al, 15h
0001C809  call sub_1E3D7
0001C80D  mov ax, 2005h
0001C810  mov bx, di
0001C812  mov si, 3000h
0001C815  call sub_1FA37
0001C818  inc bp
0001C819  add di, 0Ah
0001C81C  add dx, 50h
0001C823  jb  0001C7FC
0001C825  mov si, 3900h
0001C828  mov ax, 6
0001C82B  mov dx, 4
0001C82E  mov bx, 176h
0001C831  mov cx, 1018h
0001C834  call sub_1F888
0001C837  add dx, 50h
0001C83A  dec ax
0001C83B  jnz 0001C831
0001C83D  retn
```

因此可證實的是：

| 原始位址 | 來源 | blit 參數 | 直接尺寸／消費量 | 判定 |
|---|---|---|---:|---|
| `0001C815` | `segment1+0x3000` | `AX=0x2005` → `sub_1FA37` | 80×32、`0x500` B | 已證實為同一來源的底板候選 |
| `0001C834` | `segment1+0x3900` | `sub_1F888`、`CX=0x1018` | 3 bytes/plane × 16 rows × 4 planes = `0xC0` B；可視為 24×16 候選 | 尺寸為強推論，UI 語意 unknown |

第一段的 `DX=0,0x50,0xA0,0xF0,0x140,0x190`，下一次成為 `0x1E0` 後離開，
所以是 **6 次**，不是 25 次。第二段 `AX=6` 的 `sub_1F888` 也是 6 次。
`sub_1FA37` 會保存／還原 `SI`，所以不能把這 6 次解讀成 6 個連續來源圖塊；
它們重複使用同一個 `0x3000` 或 `0x3900` 來源。

`sub_1E3D7` 在每個 slot 前另以 `AL=cs:[bp]+0x15` 寫入其 buffer，這是獨立的
runtime buffer 寫入端，並非讀取 `ICONGRF` segment 1 的六個 glyph。此點只證實
「可見 slot 差異不必來自 `0x3000` 的六份資產」，不替該 buffer 猜命名。

### 2.2 25 次固定迴圈實際位於 `sub_1CA3B`

IDA `sub_1CA3B` 範圍 `0001CA3B–0001CAA8`：

```asm
0001CA40  mov bp, 3Ch
0001CA43  mov cx, 19h
0001CA49  mov bx, bp
0001CA4B  mov si, 3E00h
0001CA4E  call sub_1FA37
0001CA51  add bx, 12h
0001CA54  mov si, 3E80h
0001CA57  call sub_1FA37
0001CA5A  add bp, 500h
0001CA5E  loop 0001CA49
```

這是 25 次迴圈，每次使用同一對來源 `0x3E00`／`0x3E80`，每個來源由
`AX=0x1001` 定位為 16×16、`0x80` B。它不是 `sub_1C7F4` 的 25 步進，
也不能由迴圈次數直接推成六個底列命令 glyph。

同一支後段另直接使用：

| 原始位址／常式 | segment-relative offset | `AX`／來源尺寸 | 證據等級 |
|---|---:|---|---|
| `0001CA69`、`0001CA72`、`0001CA7B`、`0001CA84` | `0x3DC0` | `0x0801` → 16×8 | 已證實直接 blit 尺寸 |
| `sub_1CAA8`，由 `0001CA8D` 等四處呼叫 | `0x3D80` | `0x0801` → 16×8；每組 8 次重複來源 | 已證實直接 blit 尺寸 |

### 2.3 `sub_1C863` 的其他 segment-1 候選

IDA `sub_1C863` 使用下列 segment-relative source。接觸表已全部納入：

| 原始位址 | offset | `AX` | 尺寸／消費量 |
|---|---:|---:|---:|
| `0001C86C` | `0x0800` | `0x2008` | 128×32、`0x800` B |
| `0001C886` | `0x1000` | `0x2008` | 128×32、`0x800` B |
| `0001C89F` | `0x0000` | `0x2008` | 128×32、`0x800` B |
| `0001C8B9` | `0x1800` | `0x6008` | 128×96、`0x1800` B |
| `0001C943` | `0x3500` | `0x1008` | 128×16、`0x400` B |

這些是直接被同一初始化流程 blit 的資產候選；接觸表只呈現形狀，不把它們命名成
面板、旗幟、兵種或命令。

## 3. 與原版 320s／400s 底列比對

原版參考圖的來源是使用者提供的 YouTube 擷取，解析度 478×360 且經過有損縮放，
所以只能作視覺結構 oracle，不能作同狀態逐像素證據。

兩張圖的戰術底列都清楚呈現：

- 戰場左側下方有 6 個等寬、相鄰、紅黑框的命令 slot。
- 每個 slot 可見繁中命令文字／字形與一個彩色圖示區。
- 320s 與 400s 的 slot 幾何一致，但彩色狀態／選取內容不同；這支持「同一底板
  加上 runtime 狀態或其他繪製來源」的模型。

候選逐項結果：

| 候選 | 接觸表形狀 | 與底列可見內容的比對 | 結論 |
|---|---|---|---|
| `0x3000` 80×32 | 寬紅框、黑色內部、下方金色線；沒有六份不同內容 | 幾何上像可重複 slot 背板，但不像完整「文字＋圖示」glyph | **底板候選已證實；六 glyph 來源 unknown** |
| `0x3500` 128×16 | 黃色箭頭／指示形狀 | 不對應六個等寬命令 slot | 候選資產已定位；用途 unknown |
| `0x3900` 24×16 | 小型人物／圖形狀 planar 候選 | 不是六格文字列；且經 `sub_1F888` 位元對齊 | 尺寸強推論；用途 unknown |
| `0x0000`、`0x0800`、`0x1000`、`0x1800` | 含大型框、文字或戰術側欄樣式的複合畫面 | 不是單一 16×16 命令 glyph 陣列 | 直接 blit 候選；用途不可由畫面單獨命名 |
| `0x3D80`、`0x3DC0`、`0x3E00`、`0x3E80` | 小型紅框、色條、直向 motif | 沒有六個底列命令文字的順序／尺寸關係 | 直接來源已定位；命令 glyph 語意 unknown |

## 4. 選取 overlay 不等於 glyph 本體

IDA `sub_1C6BF`（`0001C6BF–0001C6F6`）以 `AL=0..5` 查
`CS:[BX-0x2D16]` 的六個位置，再呼叫 `sub_1F020` 兩次；它沒有讀取
`ICONGRF`。`AH=0`／`AH=0x0C` 的狀態由呼叫端提供，`sub_1C6AE` 以 `CX=6`
重畫六個 slot。這證實：

- 六個 slot 的 selected／unselected overlay 是一條獨立路徑。
- `byte_1D310` 的選取 bitfield 與 glyph 資產不能混為一談。
- `sub_1F020` 是 EGA 矩形／遮罩繪製，不是 glyph loader。

因此本輪不會把 `sub_1C6BF` 的位置表或 overlay mask 當作六個 glyph offset。

## 5. 結論與停止條件

## 5A. 2026-08-12 深入勘誤後的定案

### 已證實

1. `sub_1F938` 在每列後以 `add si,dx`（`DX=3`）前進來源；
   `CH=16`，`sub_1F888` 又依次處理四平面，所以每次呼叫消費 `0xC0` bytes。
2. `sub_1C7F4` 的六次呼叫並未重設 `SI`；六張 24×16 glyph 偏移依序是
   `0x3900、0x39C0、0x3A80、0x3B40、0x3C00、0x3CC0`，末端恰為 `0x3D80`。
3. `segment1+0x3000` 為 80×32 底板，`sub_1C7F4` 將它貼在 `(0,368)`
   起每 80 px 一格；glyph 目的點為每格 `(4,6)`。
4. `sub_1C863` 以 `segment1+0x1800`、`AX=0x6008` 把 128×96 複合命令面板
   直接貼在 `(496,280)`。面板內是單欄六列，不是 remake 先前的 2×3 文字格。
5. `sub_1C6BF` 以 `sub_1F020` 畫兩層矩形；啟用時 `AH=0x0C`，
   取消時 `AH=0`。remake 現直接使用當季 `GAMEPAL.BRG` palette index 12。
6. `sub_1E3D7` 這條路徑寫的是 RAM 命中表，不是 framebuffer。
   `CS:D2E4..D2E9`（IDA 線性 `1D2E4..1D2E9`）raw bytes 為
   `02 04 00 01 05 03`，加 `0x15` 後寫入六個 80×32 slot 的 hit code。
   這個順序與可視 glyph 順序是分離資料，不得把 raw hit code 當 glyph index。

六圖接觸表：[icongrf-seg1-battle-command-glyphs.png](../images/icongrf-seg1-battle-command-glyphs.png)；
實機畫面：[wlgame-tactical-original-command-glyphs.png](../images/wlgame-tactical-original-command-glyphs.png)。

前述早期「六 glyph 來源 unknown」與「右欄需用 2×3 文字 fallback」結論均已推翻；
保留下方原文是為了讓錯誤形成原因可追溯。

### 已證實

1. DOS/V `sub_1C7F4` 的直接來源為 `segment1+0x3000` 與 `segment1+0x3900`。
2. `0x3000` 的直接 blit 尺寸為 80×32，且同一來源重複 6 次。
3. `0x3900` 由 `sub_1F888` 的 3-byte stride／16-row／4-plane 消費模型得到 0xC0 B；
   source 的 24×16 解讀是強推論，不能升級為 UI 語意。
4. 25 次固定迴圈在 `sub_1CA3B`，使用 `0x3E00`／`0x3E80` 的 16×16 資源。
5. `sub_1C6BF` 只處理 selected overlay，不提供 glyph 本體。

### 強推論

- `0x3000` 是底列 slot 的可重複框／底板候選，因其 80×32 尺寸與六次 `DX += 0x50`
  排列相符；但「底板」以外的具體 UI 語意仍應視為強推論。
- `0x3900` 是 24×16 的位元對齊小圖候選；尚不能命名為旗、人物、游標或按鈕。

### 未知，刻意不實作

- 原版底列六個繁中命令文字／圖示的精確來源、offset、順序與完整尺寸。
- `sub_1E3D7` 寫入 buffer 的資料結構與後續顯示語意。
- `0x3000` 背板上可見文字／圖示是否由字型、另一段繪製、或 runtime buffer 合成。
- `0x3900` 與 `0x3500` 的 UI 用途。

本輪沒有修改 `internal/assets/gfx`、`cmd/wlgame/battle.go`、`WORKLIST.md`、
`CONTEXT.md`，也沒有建立 `HANDOFF.md`。在取得新的原版固定狀態畫面或更多
`sub_1E3D7` consumer／font blit 證據前，研究在此收斂；不可用候選接觸表猜接 UI。
