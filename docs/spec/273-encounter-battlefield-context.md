# 273 — 在接戰位置選場一次，供戰術與畫面共用

**狀態：CONFORMED。** 日期：2026-11-04。範圍是 Issue #31 的選場接線，
不宣稱全部水域分支或整款 remake 已完成。

- 日期：2026-11-04

## 證據與限制

- 輸入是松崗 DOS/V `KI.EXE`，SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 使用 IDA Pro 9.4 的線性位址。沿用 [spec251](251-c-engagement-tactical.md) 的固定匯出與
  `tools/c_recovery/engagement_generated.inc`。

- 資料庫 SHA-256：
  `cc4bf0da038a769fc7c12f8dcc47e884c1eb25ce959ceed0f86c3fb44182308f`。
  匯出 `workplace/matching-decompilation/c-engagement/ida-closed-v5/ida-probe.json`，
  SHA-256 `32ebc39f59e53f16d18f2f4d5bb16ae8f6971ef2dd4e733976a93dd8f296a5d9`。
  五支選場函式的每條指令均與原始檔相同，檔頭512 bytes、IDA基址10000h。
  收據：`workplace/scratch/field-selection-20261104/ida-bytes-proof.json`。
- 已證實：`14A8E` 先呼叫 `14B63`，`14A9E` 才呼叫 `14C72` 挑應戰者。
  因此選場使用最初接觸守軍座標，以及最初接觸雙方中的玩家朝向，
  不用重選代表軍團的朝向。
- 已證實：`14B71` 讀玩家攻方 `+08`，`14B79` 再由玩家守方覆寫；
  `14B7D/14B80` 讀守方 `+1C/+1A`，取以守方所在格為中心的十字五格。
  平原配對、類型1–7與類型8抽一次亂數沿用 [spec121](121-water-battlefield-selection.md)
  及 [spec196](196-field-battle-rolls-the-battlefield.md)。
- 原版兩方皆非玩家時保留入口 AL，沒有明設預設朝向。此輪不還原該暫存器生命週期，
  維持靜止方向4的 remake 近似；自動戰鬥仍消費類型8的一次亂數。
- 類型9的213編號已證實，其翻轉涉及 `14C33` 的遠指標與後續 DS 讀取。
  此輪維持既有不翻轉近似，不以「鄰格CA」或「攻方勢力」猜補。
- [playtest128](../playtest/128-player-field-collision.md) 的正常道路野戰為平原198，
  可作正常流程回歸，不能證明非平原與水域全部對拍。

## 實作契約

2026-11-04 證據審查：獨立唯讀審查確認地址及朝向來源，修正上述兩處敘述後進入READY。

1. 規則層在 `fieldAt` 挑代表守軍之前，以接觸雙方選出 `{Number, Rotate}`。
   地形不存在時採198、不抽骰；有地形時只在中心類型8取一次傳入亂數。
2. 選場結果以參數傳給 `fight`、`beginTactical`。攻城編號仍是據點編號，
   玩家守城時翻轉。驗收捷徑 `StageBattle` 也選一次，不另設正式規則。
3. `Pending.Node` 保留據點與戰果用途。新增執行期 `FieldNumber/FieldRotate`，
   戰場、適性、AI腳本、等角圖、喊話、單挑與音樂均讀同一結果。
4. `TacticalSetup` 函式簽名不變，第一個參數實際傳選定的戰場編號。
   `battlesetup.Provider` 只載資產，不讀世界、不抽骰、不按據點永久快取。
   建場回 nil 而退回自動判定時，不再選場或多抽骰。
5. `battlesetup.AttachTerrain` 把已解碼的大地圖掛到目前 World；越界回0。
   載不到陣形表時仍掛地形，供自動戰鬥使用。桌面讀檔重新掛地形，
   手機讀檔重新掛戰術及地形，不引用舊 World。
6. 原始資產 → 既有 parser → 地形分類／選場 → 規則與兩端UI → 存讀檔。
   原版及原生存檔格式均不改；待戰中仍依spec272禁止存檔。

## 驗收

- 以固定原版函式或固定 IDA bytes 複核選場資料流，保留輸入身分、位址及等級。
- 測試非平原接戰座標、參戰玩家方向含靜止、重選守軍前的方向、同據點連續兩場、
  水域抽骰次序、Provider無副作用、建場失敗降級，以及桌面／手機讀檔後仍可開戰。
  這些是 remake 測試，不混稱正常玩家對拍。
- 重跑 playtest128 的正常滑鼠野戰，抽驗攻守城共用入口；
  使用有界Docker、冷測、來源與收據雜湊。錯座標或多抽一次的錯版應被測試拒絕。
- 完整水域可見結果及類型9翻轉仍留 Issue #31，不關閉 Issue。

## 實測結果

本機收據根目錄：`workplace/scratch/field-selection-20261104/`。
程式 `wlgame` SHA-256：`eb13ad8822e3f99eb07caebeb8b6c2c58c597e819f7628dd08a4fa20c59766f2`。
工具鏈沿用 `san1-go:1.26.7-ebiten2.9.9-r1`，Go1.26.7／Ebiten2.9.9；
GUI用既有llvmpipe及LP_NUM_THREADS=0契約。

| 驗證 | 結果與界線 |
|---|---|
| 原始bytes | 五支固定IDA函式逐條與KI.EXE一致，`ida-bytes-proof.json` |
| 原版局部地形 | 沿用spec251固定O2二進位，以dosgolem重跑terrain，60案／120階段完整RAM與plane通過；受控原始RNG在執行前固定，無重擲，僅原版對固定C，不冒稱Go正常海戰 |
| remake測試 | 40套件冷測、11套件無測試，vet通過；移除Provider旋轉快取後另重跑該套件。新測試驗接觸座標、朝向4、代表守軍順序、水域單次抽骰、同據點連續兩場、建場失敗及桌面／手機讀檔 |
| 錯版反例 | Go overlay把座標改為據點時選出198而非193；多抽一次在建場前觀測到2次而非1次。兩案實際退出1，沒有編譯失敗或panic |
| 道路同狀態回歸 | 沿用playtest128原始快照與258-byte raw RNG，兩軍96點等待全通過；入場編號198、不翻轉、子刻5及游標0相同。MCH、renderer bit4、一格原版暫扣及入戰RNG階段限制沿用128 |
| 正常玩家路徑 | seed1由首頁、設定、編成、宣戰、地圖出兵、行軍存讀檔至接戰，野戰50項、攻城67項、守城68項通過，含滑鼠命令、結算及戰後存讀檔。GUI與原版快照亂數不同，不作像素或完整戰役聲明 |

守城第一次停在「行軍中開存檔」：此時有自然訊息，工具列正確停用。
驗收器改為在原8秒預算內以滑鼠推進訊息，再開存檔；同一程式、seed及400秒總預算重跑通過。
`gui-defend/` 保留失敗，`gui-defend-retry/` 保存有效結果，未因失敗改遊戲規則。
三條路徑的戰場截圖已目視，攻守城的幾何翻轉與側欄可見；完整畫面同狀態對拍仍另驗。

`verification.json`、來源／輸入／收據manifest及本機 `source.tar.gz` 保存重現資訊。
原版資料及衍生物均維持gitignore。沒有commit、push或Issue寫入。

## 未解

| 項目 | 等級 | 是否阻塞本規格 |
|---|---|---|
| 非玩家雙方的入口AL | 未知 | 否，保留已標示近似，不影響類型8抽骰次數 |
| 類型9完整翻轉語意 | 未知 | 否，維持既有限制；Issue #31仍開啟 |
| 正常非平原／海戰的完整同狀態可見結果 | 類型8選場及能力已由playtest129抽驗；完整畫面／結算待驗 | 否，完成聲明限接線與已列抽驗 |

2026-11-04 後續正常水域證據見 [playtest129](../playtest/129-normal-water-field-battle.md)：
18攻72守，CA中心選212，96點及選場後完整RNG通過，原版兵記錄與Go海戰戰力相同；
另有46項正常滑鼠水戰。遊戲程式沿用本規格版本，沒有重新實作已完成的選場接線。
