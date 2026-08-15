# 28 — 劇本的 JSON 匯出與匯入

**狀態：CONFORMED。四個區塊 round-trip 全過。**

- 日期：2026-08-15
- 出處：[`docs/formats/08`](../formats/08-sinario-save.md)（區塊佈局與每一欄的語意）；
  解碼／編碼沿用 `internal/state` 既有的 `loadBlock`／`World.Bytes`
- 推論等級：**沿用**——這一份不新增任何對原版的斷言，只換一種呈現

## 1. 為什麼要它

編輯器要能讀、能改、能存回去。目前劇本只有兩種形態：22,208 bytes 的
二進位區塊，以及 `state.World` 這個記憶體結構。**中間缺一層人看得懂、
diff 得出來、又能無損寫回的格式。**

## 2. `[HARD]` 設計約束

### 2.1 JSON **不是**完整劇本，是「已解欄位的投影」

`CLAUDE.md` §9 的存檔策略是**改寫不是重建**：從原始 bytes 出發，
只蓋已解欄位，未解區域一個 byte 都不動。JSON 沿用同一條：

```
匯出： SINARIO.DAT ──loadBlock──> World ──> JSON（已解欄位）
匯入： SINARIO.DAT ──loadBlock──> World ──套用 JSON──> World.Bytes()（改寫）──> 新檔
```

⭐ **匯入一定要有原始檔**。JSON 少了 `+0x1EC0` 那 7 KB 未解區與其他還沒
解出來的欄位；用 JSON 從零重建會產生一個「看起來對、實際少東西」的檔案。
這正是 `CLAUDE.md` §7 第 12 條說的那種缺口——**從零建立的路徑會藏 bug**。

### 2.2 名字轉成 UTF-8

`City.Name`／`General.Name` 在 `state` 裡是**原始 Big5 bytes 塞在 string**。
直接 `json.Marshal` 會被替換成 U+FFFD，**不可逆**。所以 JSON 一律存 UTF-8，
匯入時用 `text.Encode` 轉回 Big5，並檢查長度不超過原欄位。

### 2.3 不進版控

匯出的 JSON 含原版資料，`.gitignore` 已經擋掉 `workplace/`；
**不得 commit**（`CLAUDE.md` §9）。

## 3. 範圍

| 進 JSON | 內容 |
|---|---|
| `meta` | 劇本標題（區塊 `+0x40`）、來源檔名、區塊編號、來源 SHA-256 |
| `clock` | 年／月／日／時／子刻 |
| `player` | 玩家勢力編號、信賴度 |
| `economy` | 稅率、次月稅率、徵兵上限 ×2 |
| `factions` | 22 筆：君主、軍師、首都、預備兵 ×3、資金、武將數、據點數、士氣基準、交友度、攻擊性… |
| `cities` | 192 筆：名字、擁有者、座標、生產力／上限、上昇值、防災值、城兵／上限、類型、景觀圖號、內政官、鄰接 |
| `generals` | 127 筆：名字、所屬、能力值、適性、忠誠、職務、頭像、捕虜狀態… |

**不進 JSON**：事件佇列（`+0x52C0`）、未解區域。它們由改寫策略原樣保留。

## 4. 工具

`cmd/wlscen`（新增）：

```sh
tools/go.sh run ./cmd/wlscen export -in SINARIO.DAT -block 0 -out s0.json
tools/go.sh run ./cmd/wlscen import -in SINARIO.DAT -block 0 -json s0.json -out new.DAT
tools/go.sh run ./cmd/wlscen roundtrip -in SINARIO.DAT     # 四個區塊各驗一次
```

`roundtrip` 是驗收指令：匯出再匯入，**輸出必須與輸入 byte-for-byte 相同**。

## 5. 實作

| 檔案 | 內容 |
|---|---|
| `internal/scenario/scenario.go` | DTO 定義 ＋ `FromWorld`／`ApplyTo` |
| `internal/scenario/scenario_test.go` | round-trip 測試（四個區塊）、名字轉碼、越界拒絕 |
| `cmd/wlscen/main.go` | 三個子指令 |

DTO 與 `state` 的結構平行但**名字是 UTF-8、欄位有 json 標籤**。
不直接 `json.Marshal(World)`：那會連私有狀態的投影一起洩出去，
而且名字會壞掉（§2.2）。

## 6. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestScenarioRoundTrip`（四個區塊逐 byte 相同）、`TestScenarioNamesAreUTF8`、`TestScenarioRenameWritesBack`（改名要寫得進去，變短時補全形空白）、`TestScenarioRejectsOverlongName`、`TestScenarioRejectsWrongCounts` |
| 指令 | `wlscen roundtrip -in workplace/orig/dosv/SINARIO.DAT` → 四個區塊各 22,208 bytes 完全相同 |

### 6.1 名字為什麼「沒改就不重編」

`text.Decode` 會把尾端的全形空白與 NUL 修掉，再編碼回去補的是全形空白，
而原始檔那兩個 byte 可能是 NUL——**一來一回就不是同一份 bytes 了**。
所以 `keepOrEncode` 先比對 UTF-8 是否相同，相同就原樣保留，
真的改了才編碼並補到 6 bytes（不補的話新名字比舊的短時會留下舊字的後半，
因為 `state` 的寫回是 `copy` 不清欄位）。

## 7. 未解

| 項目 | 現況 |
|---|---|
| 事件佇列 | 這一輪不進 JSON。編輯器要動它得先有 UI 語意 |
| 未解區域 | `+0x1EC0` 那 7 KB 仍是黑盒，只能靠改寫保留 |
| 編輯器 | 這一份只做資料層。UI 是另一份規格 |
