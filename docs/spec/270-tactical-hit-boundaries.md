# 270：戰術命令與選隊的完整點擊邊界

**狀態：CONFORMED。原版熱區表與remake在全畫布256,000個座標一致，舊內縮實作被拒絕；44項正常戰術GUI與51項行軍回歸通過。範圍限於側欄命令與底列選隊的點擊邊界。**

- 日期：2026-11-03。
- 範圍：桌面戰術側欄六命令、底列六選隊格的命中區。戰鬥規則、速度、存檔格式不變。
- 原版：松崗 DOS/V `KI.EXE`，SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 靜態來源：[re60](../re/60-tactical-sidebar.md) §1.2／§6。

- IDA Pro9.4 database linear，segment base `0x10000`；該筆記的 DB SHA-256 `c3abeca0548ef7fd5ecf9fe531ea5971c27b143f1077291372171f809e14ab51`。
- 執行期來源：[playtest69](../playtest/69-dosgolem-tactical-controls.md)；選隊排列見 [spec33](33-squad-selection.md)，命令／陣線見 [spec37](37-tactical-player-controls.md)。推論等級為已證實。

## 原版契約

`sub_1C863` 對側欄註冊下表完整矩形，`sub_1C7F4` 註冊底列六格。
原版熱區圖為8×8格，這些矩形都對齊格界，沒有四邊各2px的內縮。
命令由 `0x1C1B9`／退卻 `0x1C21A` 分派，單擊即生效；不套用文字列帶的兩段選取。

| 區域 | 半開範圍 | 資料分派 |
|---|---|---|
| 側欄六列 | `496 ≤ x < 544`、`280 ≤ y < 376`，每列16px | 由上而下命令碼 `2,1,0,3,4,5` |
| 底列六格 | `0 ≤ x < 480`、`368 ≤ y < 400`，每格80px | 由左而右隊編號 `2,4,0,1,5,3` |

側欄 `544..551` 不屬於命令，`x552` 起是陣線格；`y376` 起是快轉列。
底列 `x480` 起是側欄柱，不能當第六格。相鄰列／格在分界處直接切換，沒有點不到的空隙。

## 實作與驗收

- `battleHitRect` 使用完整正尺寸格子，移除內縮；繪圖座標與命令／隊編號換位表不變。
- 舊測試要求框邊不命中，與原版證據衝突；改用 `TestBattleCommandHitTestsMatchOriginalFullHotspots`，依上表檢查全畫布及外部邊界。
- 正常首頁、透過既有設定調整戰術速度、正常編成守軍，等自然遭遇進戰場；在原版邊界點選隊、下令並回到戰略存讀檔。不注入遭遇或勝負。
- 原版 dosgolem 熱區圖與邊界查詢已重生，收據在 `workplace/shots/battle-hit-20261103/resume/`。原版 `siege:` 是輔助取樣，不作正常攻城路徑；正常remake守軍遭遇另行驗證。
- 正式套件冷測、靜態檢查與既有滑鼠行軍回歸已通過；新熱區表另直接餵給桌面測試，不只依賴手寫矩形。

## remake驗證結果

| 範圍 | 結果 |
|---|---|
| 全畫布命中 | `TestBattleCommandHitTestsMatchOriginalFullHotspots` 以原版矩形檢查258,084個畫布內外座標，各檢查側欄及底列；修正前首先在 `(496,280)` 失敗，修正後通過 |
| 原版熱區表 | `TestBattleHitTestsAgainstOriginalHotmap` 直接讀4000-byte原版表，逐點檢查640×400畫布；側欄4,608、底列15,360個像素的命中及分派全同，其他座標均不命中這兩組控制 |
| 原版表反例 | 僅將 `battleHitRect` 換回2px內縮，以同一原版表重跑，首先在 `(496,280)` 原版碼09、remake未命中而失敗 |
| 正常修正前 | 正常首頁編成守軍、自然遭遇進戰場，點 `(0,368)` 後選取位元仍0，收據在 `before/` |
| 正常修正後 | 44項通過：第一格及末格的左上／右下邊界選隊與解除、守陣及攻擊的邊界下令、柱與空隙不命中、鄰接陣線、退卻、返回戰略及存讀檔 |
| 主要行軍路徑 | 同一二進位51項通過，保留文字上下選列、編成、行軍中存讀檔與和平折返 |
| 正式程式 | 39個有測試套件的正式範圍冷測與vet通過；最後只讀觀測補充後，另冷跑及vet桌面套件通過 |

正常戰術驗證固定remake seed1。曹操軍團37留守許昌82，196/6/6 5時由軍團81進攻，
玩家為守方。退卻後失守許昌、首都移至洛陽58，存讀檔保留這次世界結果。
退卻不改 `Standing`，用大將的 `Cmd/Next = 5/5` 驗證，沒有為測試改動規則。
正常截圖已目視，沒有由direct-entry、強制勝負或座標注入進戰場。

受測程式為本機 `workplace/scratch/wlgame-battle-hit`，SHA-256
`bcd7904992f486eb64fe4513b0a02c78a3d0594a8a1ca303ab53dc647a2e64c0`。
收據根目錄內的 `mouse-toolbar-gui.json` 與 `regression/mouse-toolbar-gui.json` 綁定同一雜湊。
來源封存及逐檔收據由同目錄的 `source-manifest.json` 索引。

## 原版新收據

2026-11-03續查使用唯讀dosgolem `f79eab6db8c64a524a06535e53e5d32b0bbcf7c8`，
正常第一章曹操開局及編成軍團37，再於196/4/8 4時用 `siege:37,122` 輔助開戰，守軍為72。
這與正常remake抽驗的戰鬥不同，只比較固定介面的點擊區域，不比較亂數、結算或整幅像素。

取樣先到原版 `sub_1E453`，檢查實際CS、世界表與 `IDAIn` 目標。
本次CS為 `0110`，絕對與段內換算均指向實體 `0F553`，沒有發生段錯置。
`resume/attempt1/` 停在遭遇訊息；[re133](../re/133-c-startup-restoration.md) 已證實
`sub_1222B` 的這條等待路徑只讀右鍵。送出右鍵後進入戰場，呼叫觀測如下。

| 原版入口 | 累計指令數 | 意義 |
|---|---:|---|
| `0x14ADE` | 222,850,962 | 原版攻城入口，世界表段 `2754` |
| `0x11B5A` | 225,856,152 | 遭遇訊息右鍵返回後進入戰術分支 |
| `0x19946` | 225,875,726 | 完整戰場初始化 |
| `0x1C7A9` | 227,818,650 | 介面建構，DS／ES均為 `0110` |
| `0x1C7E4` | 228,257,067 | 完成控制區註冊 |

位址採IDA Pro9.4 database linear，segment base `0x10000`。
沿用的 `workplace/matching-decompilation/c-engagement/ida-closed-v5/ida-probe.json`
SHA-256為 `32ebc39f59e53f16d18f2f4d5bb16ae8f6971ef2dd4e733976a93dd8f296a5d9`；
相關指令的檔案偏移與bytes另對照原版EXE，不修改IDA資料庫或名稱。

`resume/attempt2/` 已進戰場，卻在右鍵處理越過入口後再等同一入口，因而誤報逾時。
最終時間軸移除此重複等待，從節拍6推進至56，在累計298,798,791道指令讀表並截圖。
表內有20個固定戰術熱區，加上開場對話碼27；對話範圍不覆蓋本規格的12格。
完整表長4000 bytes，SHA-256為
`17b08e0147010a49e94f046564533872777aaa2f2c1423026c88164fc7fdc4bf`。
另查詢12格的兩個對角及四個鄰接區域，共28點，均與原版表相符。

收據入口為 `resume/verification.json` 與 `resume/source-manifest.json`；
原版時間軸為 `resume/script.txt`，擷取程式在 `resume/probe/`，重播命令在 `resume/replay.sh`。
原版 `original-battle.png` 是未縮放640×400擷取，已目視確認。
以保存的 `replay.sh` 從乾淨空目錄重跑，熱區表、三張截圖及完整指令紀錄逐檔相同；
重播收據在 `resume/replay-check/`。
本次沒有重新更動遊戲邏輯，受測二進位與前述44／51項GUI相同。

直接比較原版表的命令如下，路徑為容器內路徑；未提供環境變數時明確略過此私人素材測試，
已提供但缺檔、長度錯誤或沒有完整控制區時必須失敗。

```sh
WOLONG_GO_IMAGE=san1-go:1.26.7-ebiten2.9.9-r1 WOLONG_GO_TIMEOUT=120 \
WOLONG_TACTICAL_HOTMAP=/src/workplace/shots/battle-hit-20261103/resume/original-hotzones.bin \
tools/go.sh test -count=1 -v \
  -run 'TestBattleHitTestsAgainstOriginalHotmap|TestBattleCommandHitTestsMatchOriginalFullHotspots' ./cmd/wlgame
```

### 先前失敗的取樣紀錄

本輪沿用re60的已證實矩形及playtest69的歷史執行期紀錄作實作契約。
歷史文件引用的本機原始擷取已不存在，不補造或把remake圖登錄成原版圖。

- `oracle-attempt1/`、`oracle-attempt2/` 的日期條件取樣未達成。有限profile與日期探針另確認原版仍在推進，不能據此稱遊戲或oracle整體受阻。
- `oracle-attempt3/` 的開局只有軍團19、37、72，沒有舊腳本預期的軍團35；改用實際駐122的72後，`oracle-attempt4/` 仍未達戰場入口，原因未裁決。
- `oracle-probe/` 保存可丟棄的原始介面建構探針，沒有修改golem或遊戲EXE；未達預定的註冊完成點，亦不作成功熱區收據。

上述失敗收據保留。續查後已取得有效戰場與新熱區表，不再將此項列為未知；
先前直接呼叫介面建構的探針缺少完整戰場前置，也不改列成功收據。

## 未解範圍

| 項目 | 限制 |
|---|---|
| Issue #30 完整攻守行軍 | 本片只修戰術控制邊界；正常守軍遭遇不代證玩家攻方與逐拍碰撞矩陣 |
| 戰鬥規則與畫面 | 不改陣形、單挑、結算或戰場動畫，不宣稱整幅同狀態 parity |
