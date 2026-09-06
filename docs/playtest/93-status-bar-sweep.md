# 93 — 狀態列逐條補原版擷取：財政三格全部 95 px

**狀態：進行中。** [`../spec/140`](../spec/140-status-message-box.md) §5 的
「其餘 20 幾條只有單元測試釘住索引與流程，沒有逐張擷取」開始補。
這一份收攏每一條的擷取結果，補到哪一條就在 §1 加一列。

- 日期：2026-09-06
- 規格：[`../spec/140`](../spec/140-status-message-box.md) §1.1（27 個帶索引的呼叫點）
- 受控存檔：`dosgolem/root-noclouds`（[`88`](88-controlled-save-parity.md)）
- 取樣點：196年4月17日 6時，兩邊都用 `clock:196/4/17/6`

## 1. 逐條

| TALK | 內容 | 原版怎麼走到 | remake fixture | `map` | 出處 |
|---:|---|---|---|---:|---|
| #1 | 編成第二層 | `stap:192,47` → 選武將 | `-open-form -form-pick-row 0 -lord-corps=false` | 78 | [`90`](90-formation-second-step.md) |
| #2 | 軍團 → 行軍：選軍團 | `stap:240,47` → 選單第 1 列 | `-open-command-menu corps:1` | **0** | [`84`](84-personnel-corps-list-parity.md) |
| #5／#6／#8 | 進言：敵對／停戰／請求協助 | `stap:48,47` → 選單第 N 列 | `-advise-target -advise-pick-row N` | **0** | [`89`](89-advise-status-parity.md) |
| #9 | 人事任命：選武將 | 人事 → 選城／勢力 → 選將 | `-open-command-menu personnel:0:0:0` | 95 | [`92`](92-personnel-assign-parity.md) |
| #11／#13 | 內政官任命／解任 | `stap:96,47` → 選單第 0／1 列 | `-open-command-menu personnel:0|1` | 278 | [`84`](84-personnel-corps-list-parity.md) |
| #12／#14 | 外交官任命／解任 | 同上第 2／3 列 | `-open-command-menu personnel:2|3` | **0** | [`84`](84-personnel-corps-list-parity.md) |
| #15 | 進言 → 遷都 | `stap:48,47` → 選單第 3 列 | `-advise-target -advise-pick-row 3` | 805 | [`89`](89-advise-status-parity.md)（內政吃亂數）|
| **#18** | **財政：騎兵募集人數** | `stap:144,47` → `stap:270,184` | `-open-finance -finance-amount 1` | **95** | 本份 |
| **#19** | **財政：弓兵募集人數** | 同上 `stap:270,200` | `-open-finance -finance-amount 2` | **95** | 本份 |
| **#20** | **財政：步兵募集人數** | 同上 `stap:270,216` | `-open-finance -finance-amount 3` | **95** | 本份 |
| **#0** | **編成：選武將** | `stap:192,47` | `-open-form-pick -lord-corps=false` | **0** | 本份 |
| **#16** | **財政主畫面** | `stap:144,47` | `-open-finance` | 141 | 本份（§3：那是校訂）|
| #22 | 軍團 → 位置確認 | `stap:240,47` → 選單第 0 列 | `-open-command-menu corps:0` | **0** | [`84`](84-personnel-corps-list-parity.md) |
| #23 | 據點 → 據點一覽 | `stap:288,47` → 選單第 1 列 | `-open-cities` | 278 | [`83`](83-city-list-parity.md) |
| #24／#25 | 武將／勢力一覽 | `stap:336,47`／`stap:384,47` | `-open-list`／`-open-factions` | **0** | [`86`](86-general-faction-cells-parity.md) |

95 px ＝ 原版自己畫的 14×14 滑鼠游標，停在最後一次點擊的位置。
278 px 是據點一覽每小時會動的兩欄（已裁定，[`83`](83-city-list-parity.md) §4.1）。

## 2. 還沒拍的

| TALK | 內容 | 為什麼還沒拍 |
|---:|---|---|
| #3 | 行軍：指示目標據點 | 已拍過（[`../spec/140`](../spec/140-status-message-box.md) §5），紀錄沒收進本表 |
| #4 | 軍團情報面板（`sub_17F90`）| 要在地圖上點軍團才開得到 |
| #4 | 勢力一覽常駐視窗（`sub_15AD1`）| **remake 還沒接那個視窗**（[`../re/31`](../re/31-faction-picker-screen.md)）|
| #7 | 請求協助的第二步（協同進攻對象）| 要再走一層：先選協助勢力 |
| #21 | 行軍三選一（帶 `\2` 目標據點名）| 要走完軍團 → 目的地兩層 |
| #10 | 「請選擇解任之武將。」| ⭐ **原版一個呼叫點都沒有**，是死文字，拍不到也不該拍 |

## 3. ⭐ 第一次在逐像素上看到「校訂」造成的差異

財政主畫面（#16）那 141 px 裡有 87 px 是文字本身不同：

| | 文字 |
|---|---|
| 原版 | 請指示下個月以後的財政**予定**。 |
| remake | 請指示下個月以後的財政**計畫**。 |

那是 `translations/corrections.json` 第 16 則的**刻意校訂**
（「繁中句子殘留日文『予定』，改為同義繁中詞」），不是缺陷。

⭐ **這一類差異之後只會變多**：60 筆校訂裡任何一則出現在對拍畫面上，
那一張就永遠不會是 0 px。判準是**先查 `corrections.json` 再判斷**——
「文字不一樣」在這個專案裡有兩個完全不同的成因，
而畫面上分不出來。
