# 76 — 編成的候選武將不含君主

**狀態：CONFORMED。編成候選排除君主；出陣那條照舊讓君主帶兵。**

- 日期：2026-08-23
- 出處：見 §2（**強證據，不是 confirmed**）
- 推論等級：**強證據**

## 1. 問題

`cmd/wlgame/corps.go` 的 `formCandidates()`：

```go
if gen.Alive && gen.Faction == g.world.Player &&
    !gen.Posted && gen.Captor == 0xFF {
```

**沒有排除君主**，所以 `-open-corps` 截出來的軍團一覽是
「曹操 6000／夏侯淵 6000」——玩家把自己的君主編成了軍團長。

## 2. 憑什麼說原版不行

⚠ **這一條是強證據不是 confirmed**，四項一起看：

| # | 證據 | 強度 |
|---|---|---|
| 1 | 使用者實機經驗：「原版主君不能夠編成」 | 使用者是玩過原版的人 |
| 2 | 受控 DOSBox-X 實機：走指令列「編成」編出來的軍團長是**夏侯惇**，不是曹操（`workplace/promo-live/probe-march/e1-corpslist.png`，196年4月7日）| 一次觀察 |
| 3 | ⭐ **原版讓君主帶兵走的是另一條路**：請求出陣時 `sub_16E8F` 由君主本人帶一支軍團，而且 `sub_16EC9` 專門擋「君主已經帶著軍團」（[`11`](11-ai-sortie.md)）| **結構論證** |
| 4 | 未讀：`sub_1820E`（`sub_17663` 用的清單引擎）到底怎麼濾候選 | 缺口 |

⭐ 第 3 項是這裡最有份量的一條。**如果一般編成本來就選得到君主，
那條專用的出陣路徑與那道「已經帶著軍團」的擋都沒有存在的必要。**
一個機制存在，通常是因為別的路走不到它。

⚠ 第 2 項只是一次觀察——那一次的點擊選的是清單第一列，
**不排除曹操只是排在別的位置**。所以它是佐證不是定案。

## 3. 演算法

```
formCandidates():
    for 每個武將 g:
        g 活著 且 屬於玩家勢力 且 未帶兵 且 不是俘虜
        且 **g 不是本勢力的君主**       ← 這一行是新的
```

⚠ **只擋玩家的編成指令，不要擋 `autoFormCorps`。**
請求出陣那條（`internal/state/advise.go:85`）
`autoFormCorps(w.Player, w.Factions[w.Player].Lord, false)`
**本來就是要讓君主帶兵**，照 `docs/spec/11`。兩條路不同，判定不共用。

## 4. remake 實作

| 項目 | 位置 |
|---|---|
| 候選過濾 | `cmd/wlgame/corps.go` 的 `formCandidates()` |
| 不受影響 | `internal/state/strategy.go` 的 `autoFormCorps`（出陣路徑）|
| 差異 | 無（照原版行為）|

## 5. 驗證

| 方式 | 結果 |
|---|---|
| 單元測試 | ✅ `TestFormCandidatesExcludeLord`（含反向對照：其他兩位仍在候選）|
| 既有條件沒被弄壞 | ✅ `TestFormCandidatesKeepsExistingFilters`：已帶兵／俘虜／別勢力／死亡照樣排除 |
| 實跑 | ✅ `-open-corps` 的軍團一覽從「曹操／夏侯淵」變成「**夏侯淵／夏侯惇**」 |
| 與原版一致 | ✅ 原版實機那次編出來的也是**夏侯惇** |

### 5.1 ⚠ 驗收 fixture 自己抄了一份規則

第一次改完 `formCandidates()`，截圖裡**曹操還在**——因為 `demoCorps`
（`-open-corps` 用的 fixture）自己掃了一遍 `Generals`，沒有走那支。

⭐ **修好的規則沒有套到驗收路徑上，而截圖看起來像沒修。**
現在 `demoCorps` 直接用 `formCandidates()`。這是同一輪裡第二次踩到
「一條規則兩份實作」（另一次是右鍵取消散成十二份）。

## 6. 未解

| 項目 | 現況 |
|---|---|
| `sub_1820E` 的候選過濾條件 | 未讀。定案要靠它——現在是強證據 |
| 君主被編成之後原版會怎樣 | 沒試過。若原版其實允許、只是清單排序讓人以為不行，這一條要推翻 |
