# 166 — 存檔寫回：軍團旗標的未解位元與據點游標

**狀態：READY**

- 日期：2026-09-09
- 出處：`KI.EXE`（松崗 DOS/V）`sub_16F26`（軍團建立寫 `0xC0`）、
  `sub_13EFD`（`word_10D1E`）；位元圖見 [`../re/34`](../re/34-corps-status-bits.md) §2
- 相關：[`162`](162-city-cursor-from-save.md)、`CLAUDE.md` §9
  （**存檔寫回是「改寫」不是「重建」**：只蓋已解欄位，未解區域一個 byte 都不動）

## 1. 兩個缺口

拿原版記憶體重建的存檔跑 **0 拍**再寫回，逐 byte 比對來源：

| 表 | 差異 |
|---|---|
| 勢力／據點／武將 | 完全相同 |
| **軍團** | `+0x00` 有 6 支不同（`C5` → `C4`）|

而跑 200 拍之後，寫回的全域欄位 `+0x2E`（據點巡迴游標）**還停在起始值**。

## 2. 軍團 `+0x00`：整個 byte 被覆寫

```go
r[0x00] = newCorps          // 0xC0
if c.Delegated { r[0x00] |= 0x04 } else { r[0x00] &^= 0x04 }
```

remake 建模的只有位元 7／6（存在）與位元 2（委任）。
位元 0／1／4／5 有設定端與清除端、語意尚未定案
（[`../re/34`](../re/34-corps-status-bits.md) §2），**但它們是原版的狀態**——
整個 byte 寫死等於每次存檔都把它們抹掉。

原版只在**建立軍團**時整個寫 `0xC0`（`sub_16F26`）；既有軍團的其他位元
由各自的設定／清除端維護。所以寫回要分兩種：

- 來源那一格本來就有軍團（`src[0x00] >= 0x80`）→ **保留未建模的位元**
- 來源那一格是空的（remake 這一局新編的）→ 照 `sub_16F26` 寫 `0xC0`

## 3. 據點游標 `+0x2E`

[`162`](162-city-cursor-from-save.md) 讓載入端讀了它，寫回端沒跟上。
症狀是「存檔再讀回來，AI 的時間軸會跳回原點」——而且**存檔比對會失真**：
拿原版某一刻的記憶體跟 remake 同一刻的世界比，游標永遠對不上。

單位是**段內偏移**（據點編號 × 32），與載入端同一個換算。

## 4. remake 要改什麼

`internal/state/corps.go` 的 `saveCorps`：

```go
const modelledBits = 0xC0 | 0x04 // 存在（位元 7／6）＋ 委任（位元 2）
keep := byte(0)
if r[0x00] >= aliveFlag {
    keep = r[0x00] &^ modelledBits   // 未建模的位元原樣留著
}
r[0x00] = newCorps | keep
if c.Delegated { r[0x00] |= 0x04 }
```

`internal/state/state.go` 的 `Bytes()`：寫 `putU16(b, cityCursorOffset, w.cityCursor*citySize)`。

## 5. 怎麼驗

1. 拿原版記憶體重建的存檔跑 0 拍再寫回，**四張表逐 byte 相同**。
2. 跑 N 拍後寫回，`+0x2E ÷ 32` 等於當下的游標。
3. 既有的 round-trip 測試（`TestRoutingSurvivesSaveRoundTrip` 等）維持綠燈。
