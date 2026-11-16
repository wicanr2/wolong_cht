# DOS/V 已證實敗北 outcome 接線

**狀態：READY（只涵蓋兩種敗北；不涵蓋勝利、君主死亡或原版返回標題）。**

- 日期：2026-08-12

## 證據邊界

本切片只使用松崗 DOS/V `KI.EXE` 的既有 IDA Pro 9.4 證據：

- `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- IDA `.i64` SHA-256：`7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`
- 位址空間：IDA DOS/V 線性位址。
- `sub_13DC9`：信賴度扣減借位後歸零並選 `0x019E` 群組，再呼叫 `sub_11CB1`。
  恰好歸零的selector與退出判定依caller的CX，見下方2026-11-16抽驗。
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

信賴度由 `World.AdjustTrust` 集中寫入，預設在正值降至零時latch。事件13另依13507的
CX=0 caller契約，扣50後為零即latch，包含輸入已為零；不能由這個caller外推其他扣減入口。
`OutcomeMessageSelector` 回傳的是019E群組，實際TALK槽為470加玩家君主原始說話類型，
並以 `TalkNoticeVars` 代入軍師姓名。素材、索引或參數不可用時沿既有原因句fallback。
勢力滅亡selector未定位，仍顯示remake原因句。2026-08-12的初稿誤把019E當414槽；
414實為褒獎，2026-11-16原版抽驗已推翻並修正，詳見下節。

據點易主仍由 `internal/state/corps.go:capture` 處理。當失守的是首都且
`relocateCapital` 找不到替代據點時，該勢力的 `Capital` 設為 `0xFF`、`Alive` 清除；
若該勢力等於 `World.Player`，同一 mutation 邊界 latch `DefeatFactionEliminated`。
有替代據點、非玩家勢力消滅都不會設定玩家 outcome。

Outcome 後 `Tick`、`TickMap` 與未完成戰術結算不再產生副作用。訊息／世界狀態仍可讀，
呈現層先等待前置TALK確認，再顯示原因modal。工具列確認或modal內左鍵回到remake launcher；這是 remake presentation policy，不是
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

## 未解

| 項目 | 現況 |
|---|---|
| 勢力滅亡 selector | 未定位。remake 只顯示克制的 fallback 句，不冒充原版文字 |

## 2026-11-16：事件13歸零抽驗與勘誤

**狀態：本節CONFORMED，限明示五個邊界與正常讀檔正對照。** 規格、版本、READY快照及入口見
[spec219](../spec/219-c-event-handlers.md#2026-11-16赤字警告歸零邊界)，
[整合收據](../../workplace/scratch/closeout-warning-zero-20261116/integration.json)保存完整來源與執行證據。

- 已證實：原版Trust0／1／49扣50發生借位，先通報51，再以019E群組、raw1及肖像50說471逐出，
  兩次右鍵後抵達11CB1。Trust50扣至零則第二句仍是無關編組0，再抵達11CB1；51降至1並返回。
- 原版處理期間未用callee stub或注入結果。正常標題／四槽讀檔的49案也有相同兩則訊息與退出入口，
  時鐘7時與直接取樣6時各自完整比較。22208-byte持久區塊及258-byte RNG六案全等。
- 舊Go漏判輸入0；將群組414錯讀為褒獎；Update與Draw又讓敗北框蓋過前置報告，報告無法確認。
  三處接線均已修正。舊狀態／舊文字回歸及只恢復舊優先序的GUI錯版均被拒絕。
- 五條正常滑鼠讀檔共29項／14圖通過，含停時、右鍵報告與左鍵敗北回標題；另冷測、vet與野戰62項／9圖通過。
  已抽看49案原版兩框及Go兩框。原版逐出用君主肖像框，Go沿既有較大的敗北框、標題、按鈕及暗背景，
  不稱視窗像素一致。Go仍在訊息入列前扣Trust，原版在通報確認後扣；本片只對收尾狀態與確認順序。
- 使用者A只已驗非致命51／80／255。50案不顯示無關編組，改顯示敗北原因，維持既有敗北外殼政策，
  不以A替代這個差異的說明。11CB1後不在本次原版執行範圍，回launcher不稱原版後續去向。
