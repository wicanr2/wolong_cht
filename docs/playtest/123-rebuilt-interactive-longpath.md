# 123：重建 EXE 互動長程實跑（matching spec/262）

**狀態：通過。** 冷重建 `KI.EXE` 在 dosgolem 下走完一條玩家路徑：
開局 → 關事件訊息 → 指令列 → 財政 → 人事清單 → 推進到 196 年
4 月 20 日。五張截圖、時鐘與累計指令數與原版側全同。

- 日期：2026-10-13。
- 規格：[`../spec/262`](../spec/262-matching-interactive-longpath.md)。
- 原版側：`tools/dosgolem.sh /tmp/lp2-orig "<時間軸>"`（`GAMEDIR=orig/dosv`）。
- 重建側：`WOLONG_DOSGOLEM_GAMEDIR=longpath-rebuilt`，同時間軸。
  `longpath-rebuilt/KI.EXE` 與原版 SHA-256 同為 `fffeba98…3868`。
- 時間軸：playtest/122 開局段（至 `shot:map0`）→ `steps:2000000` →
  右鍵×4 關事件訊息 → 指令列（滑鼠 672,365）→ 財政（470,397）→
  `shot:fin` → 右鍵退回 → 人事（430,397）→ 內政官任命（430,422）→
  `shot:list` → 右鍵×2 → `until:4/12；4/15；4/18；4/20` →
  `shot:later；clock`。

## 比對結果

| 取樣點 | 內容 | 原版／重建 SHA-256 |
|---|---|---|
| boot | NEW GAME 對話框 | `4233ec31…809d84` 全同 |
| map0 | 開局大地圖 | `9ab876ac…51a14a8` 全同（與 playtest/122 的 map 同值） |
| fin | 財政面板＋軍師指示 | `15d4a72e…daa2` 全同 |
| list | 人事據點清單 | `a452b718…6ea` 全同 |
| later | 4 月 20 日大地圖 | `faadaf51…812` 全同 |

- 雙側最終累計指令數同為 191,906,561；`clock` 同為 196 年 4 月 20 日 1 時。
- 完整截圖在 `/tmp/lp2-orig/`、`/tmp/lp2-rebuilt/`（`cmp` 逐張比過）。

## 過程記錄（後人不用重踩）

- 先把 `until:4/6；4/7；4/8` 插在互動段之前會改變事件時序：到 4/8 後
  沒有待關的事件訊息，第二記 `rpress` 因「畫面完全沒變」被 dosgolem
  拒絕。正解是沿用 spec/261 的 proven 序列（`steps` 推到事件出現、
  右鍵×4），日期推進段全部放在互動之後。
- 財政面板一記右鍵就退回指令列（不是直接回大地圖）；人事清單兩記
  右鍵退回後 `until` 才能一路跑到 4/20，中途沒有事件訊息擋路。
