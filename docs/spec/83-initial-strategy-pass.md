# 83 — 新遊戲的開局政略評估（sub_12BD9 的第二個呼叫點）

**狀態：READY → CONFORMED（2026-08-26）。**

- 日期：2026-08-26
- 出處：`KI.EXE`（DOS/V，SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`）。
  `sub_12BD9`（政略評估：建 22 個排序緩衝、佇列遷都／協力／停戰／宣戰事件）
  的呼叫者**全庫恰好兩個**：
  1. `sub_15358+3C`（月結尾段，remake 已接）
  2. **`sub_11AC3+66`（新遊戲流程：選定君主、寫入 `word_10CFD` 之後立即）**
  讀檔路徑（`sub_18B40` → `loc_11B2C`）**跳過**這個呼叫——
  佇列狀態隨存檔還原，不重評。
- 推論等級：**confirmed**（呼叫者窮舉；下方行為差由探針與實機對照證實）

## 1. 為什麼這一個呼叫點是行為級的

宣戰的財政閘是嚴格比較 `資金>>8 > 據點數×16+64`（`sub_12EFB`）。
孫策開局資金 21000（word 82）、1 城（門檻 80）——**只贏 2**；
第一次月結後資金掉到 word 76，之後逐月下滑，閘永遠關上。

| 時點 | 資金 word | 門檻 | 對劉繇宣戰三道閘 |
|---|---:|---:|---|
| 開局（原版的初評時點） | 82 | 80 | **全過**（交友 167≤167、國力 2000 vs 350） |
| 196/5 月結（remake 原本唯一的評估時點） | 76 | 80 | 財政閘失敗，此後每月更遠 |

這正是 [`../playtest/45`](../playtest/45-ai-longrun-comparison.md) 抓到的
「孫策不渡江攻劉繇」分歧的根因——**不是跨江路徑**（道路邊 122–129
存在，15 步；城 129 的鄰接槽含 122、位元開啟），**也不是水戰模式**。

## 2. 改動

| 項目 | 內容 |
|---|---|
| `internal/state` | 新增 `World.RunInitialStrategyPass(rng)`：政略 AI 開著時跑一次 `runStrategicAI`（同月結那次；宣戰仍走事件 1 佇列，不直接改 `+0x19`）。報表捨棄——原版的訊息也是佇列 dispatch 時才顯示 |
| `cmd/wlgame` | `startWorld` 增加 `newGame` 參數：launcher 新遊戲＝true、讀檔＝false、direct start＝「載入路徑等於劇本檔」 |
| `internal/ui/phone` | `NewSession`（新遊戲）呼叫；`save.go` 讀檔路徑不呼叫 |
| `cmd/wlsim` | 模擬一律從劇本開局＝新遊戲，呼叫 |

## 3. 驗證

- `TestInitialStrategyPassQueuesOpeningDeclarations`（`internal/state`）：
  劇本 1 初評後佇列裡有孫策的事件 1；dispatch 後孫策→劉繇交戰。
- 長程對照重跑：remake 孫策 1→3（吃劉繇），與原版半年軌跡一致
  （[`../playtest/45`](../playtest/45-ai-longrun-comparison.md) §5 更新）。

## 4. 未解

<!-- 缺口：無 -->
