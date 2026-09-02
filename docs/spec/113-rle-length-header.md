# 113 — RLE 資料檔的 4 byte 長度頭

**狀態：CONFORMED。** 原版的載入器在解壓之前先 `LSEEK` 過檔頭的 4 byte，
那 4 byte 是解壓後的長度。remake 走 `rle.DecodeFile`，
19 個過場檔加 `MMAP.MAP` 逐檔解到宣告值。

- 日期：2026-09-02
- 出處：`D7OPEN.EXE` 的 `sub_10E04`、`KI.EXE`（SHA-256 `fffeba98…d43868`）的
  `sub_1F655`、`D7END.EXE` 的同一段——三個執行檔逐指令相同
  （[`../re/76`](../re/76-d7open-opening-player.md) §5、
  [`../formats/06`](../formats/06-mmap-rle.md) §3、
  [`../formats/09`](../formats/09-cutscene-images.md) §1.1）
- 推論等級：**confirmed（靜態 ＋ 逐檔長度核對 ＋ 畫面）**

## 1. 原版做什麼

```asm
sub_1F655:                      ; KI.EXE，sub_1F5E7 的開檔子程序
        mov     ax, 4200h       ; LSEEK，從檔頭起算
        xor     cx, cx
        mov     dx, 4           ; ⭐ 讀第一個 byte 之前先跳過 4
        int     21h
        mov     bx, 800h        ; 32 KB 讀取緩衝
```

被跳過的 4 byte 是**小端 u32 ＝ 解壓後的長度**。兩版各 20 個檔走這條路
（`MMAP.MAP`、`OPEN_S1`–`S6`、`END_S1`–`S12`、`GAMEOVER.DAT`），
宣告值與跳頭解出來的長度**逐檔相等**。

## 2. 算式

```
want = u32le(src[0:4])
out  = rle.Decode(src[4:])
len(out) == want        ← 不成立就是解錯，不是「尾巴沒編進去」
```

⭐ **長度是驗收條件，不是參考值。** 從 0 開始解會在某一處掉相位，
差幾十個 byte、內容整體位移。`MMAP.MAP` 是唯一的例外：它的頭
`00 80 01 00` 沒有相鄰重複，RLE 原樣吐出而且相位不變，所以
`Decode(src)[4:] == Decode(src[4:])`——現行的 `world.ParseMap` 因此一直是對的。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| `DecodeFile(src) ([]byte, error)` | `internal/assets/rle`：讀長度頭、解 `src[4:]`、核對長度 |
| 過場圖 | `internal/assets/cutscene` 的 `Decode` 改呼叫它 |
| 世界地圖 | `internal/assets/world` 的 `ParseMap` 改呼叫它；`Map.Header` 改成**壓縮檔**的前 4 byte |
| Python | `tools/rle.py` 加對應的 `decode_file`，供 `tools/` 底下的一次性量測共用 |

⚠ **`rle.Decode` 本身不動。** 它是演算法，長度頭是容器——
混在一起之後，任何「拿一段裸資料試解」的呼叫端都得先湊一個假檔頭。

⚠ **`Map.Tiles`／`CityCentreDX`／`nodeDX`／`centreCol` 的座標框不變。**
`Tiles` 是 `out[:98304]`——長度頭已經在解壓之前被跳掉，
所以圖塊的第 0 格仍然是地圖的第 0 格，那三個常數維持 0／0／20。

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestAllCutscenesDecodeToDeclaredLength`：19 個過場檔逐檔 `len(out) == want` |
| 單元測試 | `TestOpenS2IsTwelveFrames`：12 × 32,000，且 `f0`≈`f4`、`f0`≠`f1` |
| 單元測試 | `TestDecodeFileSkipsHeaderAndChecksLength`／`TestDecodeFileRejectsWrongLength` |
| 單元測試 | `TestMMapHeaderPassesThroughUnchanged`：釘住 §2 的例外 |
| 單元測試 | `TestMapSize`（`internal/assets/world`）：`Map.Extra` 必須是 0 B |
| 畫面 | `END_S12` 的右半從雜訊變成封面兩側的龍紋；`OPEN_S2` 第 0 幀的 x ＝ 160 接縫消失 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| `BATTLE.MAP`／`MMAP.MCH`／`BATTLE.MDL` 走哪一支載入器 | 沒查（[`../formats/06`](../formats/06-mmap-rle.md) §6）。它們的前 4 byte 不是長度，所以**至少不是這一族** |
