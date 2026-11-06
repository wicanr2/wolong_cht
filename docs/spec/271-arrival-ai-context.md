# 271：抵達分派前的 AI 上下文

**狀態：CONFORMED。正式實作通過81點比較、錯版拒絕、冷測與67項正常攻方滑鼠驗收。**

- 日期：2026-11-03。
- 來源：[spec195](195-ai-stage1-asks-about-the-intent.md)、[playtest126](../playtest/126-player-march-collision-snapshot.md)。
- 原版：松崗 DOS/V `KI.EXE`，SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 位址：IDA Pro 9.4 database linear，base `0x10000`。
- 固定 IDA 匯出：`workplace/matching-decompilation/c-tick/census/ida-probe.json`，來源與資料庫身分見 [re132](../re/132-c-army-update-restoration.md)。
- 新原版正常操作軌跡：`workplace/scratch/ai-context-20261103/original/original.log`。

## 原因與證據

已證實，軍團74在首次城門碰撞後第10拍進入 `sub_143AF` 前，
`sub_12662` 已在 `0x1266E` 呼叫 `sub_128F4`。
該常式的 `0x128F4: mov di,bx`、`0x128F6/0x128F8: shl di,1`
把目前節點的原始偏移乘4。節點偏移為據點編號乘8，因此留下 `據點編號×32`。
這裡沒有加據點表基址 `0x0840`。

正常原版軌跡直接觀測到 `DI: 7D00→0350→06A0→0D40`。
軍團74腳下據點106，`106×32=0x0D40`；實際讀取的 `DS:0D58` 是0。
此值決定留守，不消費出擊延遲的亂數。舊 Go 模型仍為未知值−1，因而出擊。

已證實，本切片的 DI 在同一次抵達流程被寫入。spec195 的兩個 writer 不能涵蓋此路徑，
「只來自上一支軍團」不適用本例。歷史收據保留，不改寫成當時已通過。

## 實作契約

1. 範圍限 `sub_12662` 的目標相同分支。Go 在此處呼叫抵達分派前，
   以腳下據點寫入上下文。途中留守分支維持既有行為，不共用此算式。
2. 本次只修正腳下據點屬於軍團勢力的已證實路徑。
   敵方據點的 `sub_1291A` 後續副作用不以此切片推定。
3. 寫入值為 `node×citySize`，不加 `runtimeCityBase`，不使用意圖或目標代替腳下據點。
4. `DI+0x18` 若落在勢力記錄後半或勢力表與據點表之間，
   不得向下取整當成某勢力的武將數，也不得因不屬於勢力表而猜成3。
   由目前世界的原版區塊映射讀取該 byte，保留已解欄位與未解 bytes 的既有權威分工。
5. 直接呼叫 `arriveCorps` 仍代表 `sub_14325`，不假裝已經執行上游常式。
   未知上下文的既有近似保留。原生存檔格式不在本切片範圍。

## 驗證閘門

- 固定原版 gate 快照、258-byte raw RNG、事件節流與道路，連跑80拍。
  必要欄位與完整 raw RNG 應在81個檢查點一致。
- 正常原版重跑取得的 gate 與前16拍必須和 playtest126 的原始收據逐檔一致。
- 比較抵達分支，覆蓋相鄰 byte 門檻2／3、未知與錯誤舊上下文、
  偶數／奇數的勢力記錄位移及表間區域；檢查 Stage、Timer及完整 RNG。
  途中留守走既有回歸測試，不套入未驗證的共同算式。
- 用移除本次寫入的實際錯版確認原先第10拍分岔會復發。
- 回歸正式 Go 套件與正常玩家攻方滑鼠路徑；存讀檔後應由正式抵達路徑重新產生上下文。

## 未解範圍

完整 DI 生命週期、敵方據點副作用、守方碰撞矩陣、完整戰鬥及戰役不由此片代證。
只修正已造成玩家行軍對拍分岔的可重播路徑，不追整個 executable 的暫存器歷史。

審查已確認途中留守成功時，`sub_14300` 的 `0x14318/0x1431A/0x1431C`
會再將BX左移三次，呼叫 `sub_128F4` 時的BX不同，不能共用目標相同分支的算式。
此路徑尚未取得新正常原版收據，不進本次正式修正。

## 證據審查

2026-11-03：固定IDA匯出的78條指令與KI原始bytes相符；新原版gate及前16拍的
68份狀態資料重現舊收據。限縮後的可丟棄實作通過81個檢查點與16個門檻案例。
完整區塊的2,144筆差異只在 `0x21C0..0x22BF` 的MCH物件區，不宣稱整區塊逐byte相同。
審查拒絕了最初把途中留守當成相同算式的版本，該版保存在本機 `prototype-initial/`。
READY時另列錯版拒絕、冷回歸與正常GUI為正式驗收閘門，結果見下節。

## 正式實作與驗證

正式變更在 `internal/state/corps.go` 的目標相同分支，以及
`internal/state/aimarch.go` 的 `residualField18`。
`internal/state/arrival_context_test.go` 的16個案例覆蓋據點106、42、43、64，
門檻2／3與未知／錯誤舊上下文；檢查Stage、Timer與完整258-byte RNG。

| 項目 | 結果 |
|---|---|
| 正式80拍比較 | 81點的必要欄位、完整RNG及98,304格佔用低7位全同 |
| 完整區塊額外檢查 | 81點在MCH物件區以外零byte差異；原始差異仍保存，不宣稱整塊完全相同 |
| 移除抵達前寫入的錯版 | 真實行為測試有8個案例失敗；80拍比較從第10拍失敗，與原先分岔一致 |
| 正式Go冷測 | `go test -count=1 ./cmd/... ./internal/... ./mobile/... ./translations/...`，39個套件通過，12個無測試 |
| 正常攻方滑鼠路徑 | 67項通過；首頁開局、宣戰、編成、等待自然守軍、出兵、途中存讀檔、玩家主攻、滑鼠突擊、戰後存讀檔 |
| 目視 | `attack-charge.png` 為汝南戰術畫面；`attack-post-load.png` 為戰後讀檔的戰略畫面 |

受測程式為本機 `workplace/scratch/ai-context-20261103/wlgame`，SHA-256
`e1d987e089e47724998d7b29feff7d8a776b3def5e5933c43a9a91fd6d660e44`。
GUI仍用seed1，沒有注入軍團、座標、交戰或勝負；正常輸入的牆鐘逾時只作腳本上限。
該GUI與原版80拍快照不是相同戰鬥，不互相代證整場戰鬥結果或像素一致。

本機收據根目錄為 `workplace/scratch/ai-context-20261103/`：
`original/` 保存新原版軌跡，`ida-evidence.json` 保存逐指令bytes核對，
`remake-final/`、`final-comparison.json` 保存正式重播，
`negative-tests.log`、`remake-negative/`、`negative-comparison.json` 保存真實錯版。
`cold-tests.log` 與 `gui-attack/` 保存正式回歸及正常玩家路徑。
來源封存與重播命令由同目錄 `source-manifest.json`、`verification.json` 索引。

原版重播沿用playtest126的正常輸入，取樣依遊戲tick與函式入口，不依等待秒數。
比較器仍為 `tools/march_snapshot_compare.py`；命令使用既有 `tools/go.sh` 與 `tools/py.sh`
在Docker內執行。正式來源基線為 `bd09fd40e344efee9d01e286a70f2a2303890caf` 加工作樹，
工具鏈沿用 `san1-go:1.26.7-ebiten2.9.9-r1`。

本輪未改存檔格式。快照沒有保存DI的既有缺口仍屬限制，
本片的正常抵達流程會重新產生已證實上下文，不把它換稱為完整DI保存。
Issue #30仍開啟；途中留守、守方碰撞、完整戰鬥及戰役不由本片代證。
