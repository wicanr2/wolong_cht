# 69 — 世界指紋：同一個 seed 兩次跑出同一個值

**狀態：CONFORMED。⚠ 這一份沒有原版出處——它是 remake 的驗證設施，不是原版機制。**
規格照樣要寫，因為它會改到 `internal/`，而且**它涵蓋什麼、不涵蓋什麼**
決定了拿它下的結論可不可信。

- 日期：2026-08-20
- 出處：**無**。原版沒有這個東西
- 推論等級：不適用（不是對原版的斷言）
- 為什麼需要：Android 里程碑 A 的判準是「同一個 seed 在手機與桌面算出相同結果」，
  而現在沒有辦法把一個 `World` 縮成一個可比較的值
  （[`../mobile/android-plan.md`](../mobile/android-plan.md) §5）

## 1. 要解決的問題

長跑之後怎麼知道兩次跑出來的世界一樣？現在只能逐欄位比，
而 `World` 有 22 個勢力 × 24 B、192 個據點、127 名武將、256 筆事件佇列，
外加十幾個不序列化的 runtime 欄位。**逐欄位比對寫一次就沒人願意維護。**

⭐ 這件事的價值不只在 Android：桌面端也需要一個**決定性迴歸**——
改動規則層之後，同一個 seed 跑同樣的 tick 數應該得到同樣的世界，
不然就是有人不小心把亂數、走訪順序或 map 迭代帶進了規則路徑。

## 2. 涵蓋什麼

指紋 ＝ SHA-256（存檔位元組 ‖ runtime 欄位的正規編碼）。

| 進指紋的 | 從哪來 |
|---|---|
| 整個劇本／存檔區塊 | `World.Bytes()`——時鐘、勢力、據點、武將、軍團、事件佇列、存活勢力數都在裡面 |
| 據點整備游標 | `cityCursor` |
| 事件游標與延遲 | `eventCursor`／`eventDelay` |
| 軍團游標 | `corpsCursor` |
| 災害 marker 與等級 | `disasterMarkers`／`disasterMarkerLevels`（192 格）|
| 災害物件 32 槽 | 每槽的種類、位置與計時 |
| 亂數狀態 | `rng.Rand` 的 `c`／`s` 兩個 byte（表在播種後就不再變）|
| 待決狀態的**種類** | `pending`／`encounter`／`diplomacy`／`funding` 是不是 nil |

## 3. ⚠ 不涵蓋什麼

**這一節比上一節重要**——不寫清楚的話，指紋相同會被讀成「兩邊完全一樣」。

| 不進指紋的 | 為什麼 |
|---|---|
| 道路圖 `roads` | 從 `MMAP.MAP` 算出來的常量，跑再久都不會變 |
| 戰術戰鬥 `tactical` | 那是一場子戰鬥的內部狀態，有自己的驗收（`internal/rules/tactical`）|
| 待決狀態的**內容** | 只記「有沒有」，不記選項細節——那些是 UI 層的暫態 |
| 畫面、音訊、輸入 | 指紋是規則層的，**與畫面無關**。兩邊指紋相同不代表畫面相同 |

⭐ **亂數表不進指紋是刻意的**：`rng.Rand` 的 256 B 表在播種之後就不再改，
把它算進去只是讓每個指紋多算 256 B，不會多抓到任何一種分歧。
c／s 兩個 byte 就完全決定了接下來的輸出。

## 4. remake 實作

| 項目 | 位置 |
|---|---|
| `World.Fingerprint() [32]byte` | `internal/state/fingerprint.go` |
| 十六進位短碼 | `World.FingerprintHex() string`（前 16 個字，給 log 與測試訊息用）|
| 亂數狀態 | `rng.Rand.State() (c, s byte)`；`state` 端以 `interface{ State() (byte, byte) }` 取，取不到就記一個明確的「沒有亂數」標記 |

⚠ **編碼要固定**：所有整數以固定寬度大端寫入，陣列照索引順序走，
**不得走 map**。走 map 的話同一個世界每次算出不同的指紋，
而那種錯誤會以「Android 與桌面不一致」的形式出現，查起來會很久。

## 5. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 ✅ | `TestFingerprintIsDeterministic`：同一個 seed 跑 500 tick 兩次 → 指紋相同；換 seed 或少跑 100 tick → 指紋不同 |
| 單元測試 ✅ | `TestFingerprintIsStableForTheSameWorld`：同一個世界連算 9 次 → 值相同（抓 map 走訪那一類的不決定性）|
| 單元測試（正對照）✅ | `TestFingerprintCoversEveryRecordedField`：**15 個進指紋的欄位各改一格，指紋一定要變**。這一條擋的是「漏掉某個欄位」——漏掉不會讓任何測試變紅，只會讓指紋安靜地失去偵測力 |
| 單元測試 ✅ | `TestFingerprintDistinguishesMissingRandomSource`：「問不到亂數狀態」與「狀態剛好是 0, 0」不可以算出同一個指紋 |

## 6. 未解

| 項目 | 現況 |
|---|---|
| 跨平台實測 | Android 端還沒有東西可以跑（里程碑 A 本身）|
| 戰術戰鬥要不要進指紋 | 目前不進。要驗戰場的決定性得另外做一個，`tactical.Battle` 的欄位更多 |
