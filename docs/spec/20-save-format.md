# 20 — remake 原生存檔格式

**狀態：CONFORMED。編解碼、路徑與遊戲接線都實作並驗過。
存檔一次寫兩份（原版格式 ＋ 原生檔），讀檔優先原生檔。**

- 日期：2026-08-14
- 出處：[`docs/formats/08`](../formats/08-sinario-save.md)（原版區塊版面）、
  `internal/state/state.go`（現行載入／寫回）、`internal/savepath`（路徑邊界）
- 推論等級：原版格式 **confirmed**；本規格是 remake 的設計，不是原版行為

## 1. 為什麼要有自己的格式

現況：remake **只用原版格式**。載入 `SINARIO.DAT`／`SAVE.DAT` 的第 N 個
22,208 B 區塊，存檔時走「改寫」——從原始 bytes 出發，只蓋已解欄位，
未解區域一個 byte 都不動（`CLAUDE.md` §9）。byte-for-byte round-trip 有測試守著。

這個策略保住了保存價值，但有三個代價：

| 代價 | 說明 |
|---|---|
| **裝不下 remake 的狀態** | 行軍路線（`routes`）、災害物件、風暴區、事件游標、據點輪轉游標都不在原版記錄裡。目前的作法是「不序列化」——**存檔再讀回來就掉了** |
| **讀寫都要整條解碼鏈** | 想看一份存檔裡誰在哪，得先跑 `LoadScenario`。外部工具（測試 fixture、對拍腳本、bug 回報）沒有便宜的入口 |
| **不可 diff** | 兩份存檔差在哪，只能比 22,208 個 byte |

## 2. 設計

### 2.1 一條不可退讓的約束

> **原生存檔必須能無損匯出回原版格式。**

所以原生檔**要帶著原始區塊的 bytes**。少了它，那 7 KB 未解區
（`+0x1EC0`–`+0x42C0`）與所有還沒解的欄位在匯出時就只能填 0——
那等於把「改寫不是重建」這條硬規則作廢。

### 2.2 檔案結構

單一檔案，JSON（理由見 §6 決策 A）：

```jsonc
{
  "format": "wolong-save/1",
  "origin": {
    "source": "SINARIO.DAT",        // 或 SAVE.DAT
    "block": 0,                      // 劇本／存檔槽
    "block_sha256": "…"              // raw 的雜湊，驗檔案自己有沒有壞
  },
  "raw": "<base64 的 22,208 bytes>",  // ← 存檔當時的狀態寫成原版格式，§2.1
  "state": { …已解欄位，人看得懂的形狀… },
  "runtime": { …原版記錄裡沒有的 remake 狀態… }
}
```

### 2.3 權威來源怎麼分

**`state` 對已解欄位有權威，`raw` 供應其餘一切。**

⚠ `raw` 是 **`w.Bytes()`**：從載入時的原始位元組出發、蓋上目前的已解欄位。
它同時是保存錨點（未解區域原樣帶著走）與進度（已解欄位是現在的值）。
**不可以存 `w.RawBlock()`**——那是載入當時的區塊，存進去等於把開檔之後
的進度全部丟掉，而檔案的格式、雜湊、載入流程全部看起來正常。
擋這件事的是 `TestEncodeCapturesTheCurrentState`：它**先推進 600 tick 再存**，
而只在開局狀態存檔的測試對這個錯誤完全無感。

- 載入：以 `raw` 建出 World（走現行的 `LoadScenario` 解碼），
  再把 `state` 的欄位覆蓋上去。
- 匯出：拿 `raw` 當底，覆寫已解欄位（就是現行的 `SaveInto`），
  **產出必須與「同一份 World 走原版路徑存檔」byte-for-byte 相同**。
- 一致性：載入時若 `state` 與 `raw` 解出來的值不合，**大聲失敗**，
  不要靜靜挑一邊。那代表檔案被手改過或版本不合。

> `raw` 與 `state` 重複儲存同一份資訊是**刻意的**。
> 重複可以互相檢查；只留一份就沒有東西可以對。

### 2.4 `runtime` 放什麼

現在「不序列化」而應該存起來的：

| 欄位 | 現況 |
|---|---|
| `routes`（每支軍團的行軍路線）| 讀檔後歸零，軍團會改走重算的路 |
| `disasterObjects`／`disasterMarkers`／`stormArea` | 災害的執行期物件，讀檔後消失 |
| `eventCursor`／`eventDelay` | 事件佇列游標，讀檔後重設 |
| `cityCursor` | 據點輪轉游標——**它決定 AI 下一個處理哪個據點**（`docs/spec/10`）|
| `Player` | 目前由呼叫端設定，不在存檔裡 |

**這些一律標成 remake 差異**（`CLAUDE.md` §1）：原版沒有它們，
所以匯出回原版格式時會遺失，而遺失的後果要寫進文件。

## 3. remake 實作（打算改哪裡）

| 項目 | 位置 | 動作 |
|---|---|---|
| 新套件 | `internal/savefile` ✅ | `Encode`／`Decode`／`ExportOriginal`／`VerifyExport`。純編解碼，不依賴 Ebiten |
| 匯出 | 同上 | `ExportOriginal(*state.World, raw []byte) ([]byte, error)`，內部就是現行的 `SaveInto` 邏輯 |
| 狀態層 | `internal/state/snapshot.go` ✅ | `Snapshot` 型別 ＋ `RawBlock`／`TakeSnapshot`／`Restore`／`LoadBlock`。JSON tag 只在 `Snapshot` 上，沒有灑進 World |
| 路徑 | `internal/savepath` | `NativePath(overlay, slot)`：`SAVE.DAT` → `SAVE-slot1.wlsave`。原版 overlay 與原生檔各自獨立，互不覆蓋 |
| 遊戲 | `cmd/wlgame/save.go` | **一次存兩份**：原生檔（遊戲讀的）＋ 原版格式 overlay（拿去 DOSBox 的）。讀檔優先原生檔，沒有就退回原版格式 |

> **為什麼不是「選單多一項匯出」**（原本 §3 這樣寫）：多一個動作就多一個
> 「忘記匯出」的狀態，而原版格式正是這個專案的保存價值所在。
> 兩份一起寫沒有額外 UI、也沒有不一致的中間態，代價只是每次多寫 22 KB。

## 4. 驗證

| 方式 | 內容 | 狀態 |
|---|---|---|
| 單元測試 | **來回一致**：原版 → 原生 → 匯出，byte-for-byte 相同 | ✅ `TestRoundTripToOriginalIsByteForByte` |
| 單元測試 | **`runtime` 真的活著**：`routes`／`cityCursor`／`eventCursor`／AI 開關 | ✅ `TestRuntimeStateSurvivesTheRoundTrip` |
| 單元測試 | **改壞要炸**：動過 `raw`、版本不合、`runtime` 索引超界 | ✅ `TestDecodeRejectsTamperedFile` |
| 單元測試 | 沒有原版檔案也載得動（決策 C）| ✅ `TestDecodeDoesNotNeedTheOriginalFiles` |
| 單元測試 | 存檔寫出兩份、讀檔優先原生檔、原生檔壞掉要炸、不存在要退回原版格式 | ✅ `TestSaveWritesBothFormatsAndReadsBackNative` |
| 單元測試 | 路徑推導與「原生檔不會蓋掉原版 overlay」 | ✅ `TestNativePath` |
| 對原版 | 匯出的檔案放回 DOSBox 能讀 | ⬜ **未做**（工具見 [`../playtest/21`](../playtest/21-dosboxx-bridge-sampling.md)）|

## 5. 未解

| 項目 | 現況 |
|---|---|
| 存檔區塊的 7 KB 未解區 | `+0x1EC0`–`+0x42C0`，靠 `raw` 原樣保存，但**內容仍不知道**（`docs/formats/08`）|
| 原版 `SAVE.DAT` 的槽位語意 | 四個槽與 `SINARIO.DAT` 的四個劇本是不是同一個編號空間，未確認 |

## 6. ⚠ 要裁定的三件事

| # | 決策 | 我的建議 |
|---|---|---|
| **A** | JSON 還是二進位？ | **JSON**。可 diff、可手改來造測試 fixture、bug 回報看得懂。代價是檔案大（22 KB 的 base64 ≈ 30 KB，加上 state 約 60–80 KB）。存檔不是效能瓶頸 |
| **B** | `raw` 要不要壓縮？ | **先不要**。壓了就不能直接 diff，而這個格式的主要價值就是看得見。要省空間再說 |
| **C** | 原生檔要不要能在**沒有原版檔案**時載入？ | **要**。`raw` 已經自帶那 22,208 B，所以技術上可行。但這與「啟動要驗證原版檔案」的政策衝突——**驗證應該在啟動時做一次，不是每次讀檔做**，否則玩家換機器就讀不了自己的存檔 |
