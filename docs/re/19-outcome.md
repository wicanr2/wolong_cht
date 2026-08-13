# DOS/V 已證實敗北 outcome 接線

**狀態：READY（只涵蓋兩種敗北；不涵蓋勝利、君主死亡或原版返回標題）。**

- 日期：2026-08-12

## 證據邊界

本切片只使用松崗 DOS/V `KI.EXE` 的既有 IDA Pro 9.4 證據：

- `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- IDA `.i64` SHA-256：`7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`
- 位址空間：IDA DOS/V 線性位址。
- `sub_13DC9`：信賴度扣減造成 `byte_10D00=0` 後顯示 selector `0x019E`，呼叫
  `sub_11CB1`。
- `sub_14DF0`：首都被攻陷且 `sub_16A3D` 找不到同勢力替代據點時，寫
  `capital=0xFF` 並清除勢力 alive bit。
- `sub_14FCE`：若被清除的是玩家勢力，呼叫 `sub_11CB1`。

`sub_11CB1` 只證實離開主遊戲循環；它之後是返回標題、結束程式或其他流程仍為
unknown。因此 remake 的「確認後回啟動畫面」是明示的 presentation policy，不是原版
行為宣稱。

## Remake 映射

`internal/state/outcome.go` 提供 runtime、不可存檔的單次 latch：

- `InProgress`
- `DefeatTrustZero`
- `DefeatFactionEliminated`

信賴度由 `World.AdjustTrust` 集中寫入。已知 production caller 是進言 UI、事件 13、
外交超額提案；由 1 降至 0 的同一寫入邊界才 latch `DefeatTrustZero`，已經為 0 不會
重新觸發。原版 selector 以 `OutcomeMessageSelector` 回傳 `0x019E`；DOS/V
`TALK.DAT` 第 414 則可安全解碼，沒有 marker／payload 依賴，因此 GUI 優先顯示這則
原版訊息。若 TALK 資產不可用、索引失效或代入不安全，才顯示「信賴度歸零，已被逐出
勢力。」這個 remake fallback。勢力滅亡 selector 尚未定位，故只顯示克制的 remake
原因句，不把它冒充原版訊息。

據點易主仍由 `internal/state/corps.go:capture` 處理。當失守的是首都且
`relocateCapital` 找不到替代據點時，該勢力的 `Capital` 設為 `0xFF`、`Alive` 清除；
若該勢力等於 `World.Player`，同一 mutation 邊界 latch `DefeatFactionEliminated`。
有替代據點、非玩家勢力消滅都不會設定玩家 outcome。

Outcome 後 `Tick`、`TickMap` 與未完成戰術結算不再產生副作用。訊息／世界狀態仍可讀，
呈現層顯示原因 modal，不把反組譯位址、研究狀態或未知流程放進玩家畫面。Enter、Space
或 modal 內左鍵確認後回到 remake launcher；這是 remake presentation policy，不是
原版 `sub_11CB1` 後續去向的宣稱，也不會自動重開同一局。

## 驗證

deterministic tests 涵蓋：

- Trust 1→0、Trust 已為 0 不觸發；
- 失去最後首都、仍有替代首都；
- 非玩家勢力消滅；
- outcome latch 不覆蓋；
- modal 凍結與確認後 launcher。

`-open-outcome trust`／`-open-outcome faction` 僅為截圖 fixture；前者仍透過
`AdjustTrust`，後者只是 UI 取樣用的受控 latch，正常玩家路徑不依賴它。

實機畫面（Docker＋Xvfb，640×400）：

- [信賴度歸零 modal](../images/wlgame-outcome-trust.png)
- [玩家勢力消滅 modal](../images/wlgame-outcome-faction.png)
