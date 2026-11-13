# 62 — 被換位的兵，這一幀不動

**狀態：CONFORMED，限既有換位及門圖塊地形旗標片段。** 原版的旗標有兩個作用，兩個都接上了。
補上第二個之後，[`61`](61-soldier-initial-hp-from-morale.md)
（開場體力 ＝ 士氣）才接得上——兵撐得久就會擠在一起，
而擠在一起打不打得到人由這一條決定。

- 日期：2026-11-13
- 出處：`sub_1ADC8`（`0001ADED`）、`sub_1B732`（`0001B732`）、
  `sub_1B240`（`0001B358`）
- 推論等級：**confirmed（靜態）**

## 1. 原版做什麼

2026-11-14射擊略過勘誤：`1AD5C`／`1ADA6`發射成功也設+0x00 bit6，沿既有重畫清除。
因此remake欄位Swapped保留原名，但不只代表被別人換位；射擊當拍同樣略過移動。
原指令及正常重播的射擊硬直見[spec159](159-tactical-command-movement-order.md)，不改本文件的換位資格。

### 2026-11-13：門圖塊的地形旗標排除（CONFORMED，限本片）

本片只修正地形旗標的圖塊條件，不改重畫時機、換位、平面、登城高度、速度或存檔。
正常曹操對呂布重播第152拍槽87位於(45,38)、Z0、低平面，原版bit1為0而Go為1。
原始BATTLE.MAP該格為0xFD，BATTLE.MDL組0定義高度5；目前Go只判高度至少4。

已證實，靜態：IDA Pro9.4、DOS/V linear `1B2D3–1B2DD`用目前XY讀`word_1D2F6`圖塊；
`1B2E1 cmp bl,F0h`及`1B2E4 jnb 1B300`先排除圖塊值≥0xF0，再由`1B2E6–1B2F9`讀定義高度判≥4。
`1B300`清bit1，`1B2FB`才設bit1。輸入KI.EXE SHA-256沿本文件既有版本，
[原指令及bytes](../../workplace/scratch/ui-windows-20261104/structure-cost-20261113/ida-probe.json)
的DB SHA-256為`bad62f79a66cf9fee5976c60f9ade426f4347b44fd3e527d987016d32625eee0`。

候選契約：有原始圖塊時，先排除0xF0–0xFF；低於0xF0才判高度≥4。
沒有原始圖塊的合成戰場保留高度判定。初始化及逐兵重畫共用`Field.HighTerrain`，
鎖敵沿既有消費，不交換或提前更新旗標。原始資料、typed Field、規則、正常GUI沿原路徑，存檔格式不變。
驗收先核對同指令／seed1的182筆原版觀測逐byte不變，再於槽87重畫點補讀實際圖塊與清旗標分支。
候選須消除第152拍差異；後續登城Z差異另列，不藉本片猜補。
回歸涵蓋0xEF／0xF0、0xFD高度5、低高度、合成戰場、初始化及停留重畫，舊程式須被拒絕。
正式接線前須完成候選審查；接線後冷測受影響套件並抽驗正常滑鼠編成、行軍、存讀檔、入戰、登城及退卻。

[審查收據](../../workplace/scratch/ui-windows-20261104/cao-terrain-review.json)確認新原版182筆主觀測逐byte相同，
101筆命令收據也與前片相同。兩筆既有診斷為框外右鍵未改畫面及有界runto未到指定常式；
它們不證明該步被接受，後續自然入戰及固定戰術觀測仍重現，不以診斷作原版無錯宣稱。
槽87於tick151重畫、boundary152之前讀到0xFD／高度5，依序走1B2E1、1B300、1B303，bit1為0。
候選前152拍96槽27欄位及完整RNG相同；第153拍首次差異僅槽87的Z原版1、Go4。
獨立唯讀審查確認差量僅圖塊上限閘，合成戰場與重畫時機保留。靜態、動態及候選證據齊備後列READY。

正式接線通過[七套件冷測／vet／建置](../../workplace/scratch/ui-windows-20261104/cao-terrain-tests.log)，
0xEF／0xF0、0xFD高度5、低高度、合成高度及初始化／重畫皆有回歸；舊程式實際被兩條斷言拒絕。
[正式證據](../../workplace/scratch/ui-windows-20261104/cao-terrain-evidence.json)前152拍96槽27欄位、完整RNG及腳本相同；
成本8192 bytes及96槽剩餘點數仍以前153拍為界，154拍首次不同。153拍槽87 Z1／Z4保留，未調整登城動作。
[正常滑鼠](../../workplace/scratch/ui-windows-20261104/cao-terrain-remake-gui.json)69個檢查點、16圖通過，
曹操第一劇本宣戰呂布、編夏侯淵、停手行軍、途中存讀檔、自然遇敵、改線、登城、突擊及退卻返回。
登城210→216拍指標不動，主將X28→34，Goal保持(42,32,4)；已目視城牆、單兵及側欄指令。
[來源收據](../../workplace/scratch/ui-windows-20261104/cao-terrain-integration.json)保存兩檔來源、binary、輸入及證據雜湊。
不宣稱完整戰果、所有地形、許昌留守、架梯／投射物或手機GUI完成。

2026-11-13純Z後續：`1B0FD`逐層上升已由[spec36續片](36-ground-planes-and-climbing.md)接線，
前156拍逐兵／RNG相同。上述153拍Z1／Z4是旗標片當時的界線，目前首次差異為157拍受擊／命令及RNG。

### 2026-11-12：換位的中繼高度與地形旗標續片

**狀態：CONFORMED，限本片。** 舊文的被換位跳過機制保留。
本片從第85拍槽73／79的StepZ差異開始，先以可丟棄覆蓋檔驗證，
不修改換位資格、碰撞格表、傷害、AI命令或存檔格式。

- 工具為IDA Pro9.4，位址均為DOS/V linear；輸入KI.EXE SHA-256為
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
  一次性DB SHA-256為`47537903862d87add7996cc733e57e3e1ed0182771cf1c1946bc5122bb36b16e`。
  [原始指令與bytes](../../workplace/scratch/ui-windows-20261104/swap-20261112/ida-probe.json)
  保留原始函式名、739函式探針、工具與UID/GID。
- **已證實，靜態**：`sub_1B732`交換位置、碰撞索引及`+1E`平面後，
  `1B778–1B78D`以`1B783 or al,ah`先檢查雙方平面OR為0，再檢查雙方目前Z的OR不為0。
  兩條件皆成立，`1B78F–1B795`才交換兩人的`+12h`中繼高度。
  平地兩人Z均0或任一方在高平面時，中繼高度各自保留。
- **已證實，靜態與實測**：此函式不交換`+00 bit1`地形旗標。
  它只以`1B732 or [di],40h`標記被換者。舊Go交換HighTerrain的行為沒有原版依據。
  原版保留旗標直到各自重畫，不在換位函式依新位置重新計算。
  `sub_1B240`的`1B2D3–1B300`於每名兵重畫時，先排除圖塊≥0xF0，再以定義高度是否至少4重設bit1。
  被換者輪到時也會重畫；已經輪過的被換者則保留原旗標直到下一拍。
  [重畫指令與bytes](../../workplace/scratch/ui-windows-20261104/swap-redraw-20261112/ida-probe.json)
  的DB SHA-256為`1a217870ca1a7ed5cba51c4d04110c4ba5bb83246af2db4453bbc9c72c41f25a`。
  此旗標由鎖敵的平面懲罰消費，不能只當繪圖狀態。
- dosgolem來源`f79eab6db8c64a524a06535e53e5d32b0bbcf7c8`，原版與remake使用同一初態及受控RNG。
  固定seed1檔SHA-256為`2624479abfc3b6b8b9fc4b12961433dfaa55c08e301381b5167068b9a2654fb7`。
  [正常玩家輸入](../../workplace/scratch/ui-windows-20261104/cao-swap-original.commands)
  自第一劇本曹操開局、向呂布宣戰、編成行軍後自然入戰；沒有直接開戰或強制勝負。
  [重播快照](../../workplace/scratch/ui-windows-20261104/cao-swap-original.jsonl)117筆與前輪逐筆相同。
- [換位前後記錄](../../workplace/scratch/ui-windows-20261104/cao-swap-events.jsonl)
  在`1B732`及`1B798`各取一次，截至第115拍共467對，全部符合上述條件。
  17對位於低平面坡地，3對的StepZ不同，4對的地形旗標不同。
  第84拍槽73／79換位前Z為1／0、StepZ為1／0，換位後均為0／1；
  這解釋了第85拍命令後、移動前的首次差異。此重播沒有高平面換位，高平面排除條件只由靜態指令支持。

實作契約：`swapWith`按上述條件交換StepZ，移除HighTerrain對調；
`drawUnitCollision`在每名活兵的重畫點，依目前圖塊及高度更新HighTerrain，不等待移動成功。
既有位置、平面、碰撞格及Swapped處理保留。此為原版規則勘誤，不是新玩法。
`TestSwapWaypointHeightAndTerrainFlags`須分開測平地、單側坡地、雙側坡地與高平面；
正式修正再以同初態逐拍比較、正常滑鼠入戰／改線／退卻及相關冷測驗證。
存檔欄位不變，正常途中存讀檔沿既有GUI路徑抽測。
尚未驗證的高平面導航、投射物、架梯及完整結算不由本片推定完成。

2026-09-08 大將例外（READY）：本輪 IDA 9.4 匯出的
`ida-collision-lifecycle/collision-entry.json` 保留未歸屬函式區段原始 bytes。
IDA DOS/V linear `0001B548–0001B557`：兵種 0 且命令 0 或 5 時，
直接跳到 `0001B584` 的平面檢查，繞過敵我分派與一般換位四閘。
`0001B584–0001B598` 比較的是 +0x1E 平面，並非兩兵實際 Z。
此分支為已證實；原輸入 EXE 身分同下列 SHA-256。
允許在碰撞及 swapWith 補上此例外，測同平面大將遇到移動中同伴及敵人，
仍保留跨平面雙方須可爬的限制。不可用任意放寬全部換位閘代替。

2026-09-08 訂正（窄修正 READY）：`[di] & 0x61` 也擋 bit 5。
IDA Pro 9.4 DOS/V linear `0001AF6B` 設 0x20；`0001AFE4` 等成功出口
直接返回，只有 `0001AFBF`、`0001AFCC`、`0001AFD0` 清除。
因此舊「重入保護、離開即清」解釋不成立，成功後仍保留到後續清除分支。
原始指令見 `workplace/delegated-battle-20260908/ida-script/battle-cause.json`，
碰撞遮罩見 re/11 §5.16；輸入 KI.EXE SHA-256
`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
此旗標生命週期為已證實；實作需保存 bit 5 並在 swapWith 擋住目標，
以第二拍全兵追蹤及野戰重播驗證，不把旗標當局部函式變數。

單位記錄 `+0x00` 的 bit 6 是「**這一幀被別人換走了**」。
`sub_1B732`（對調位置）在最後一行標記被換的那一個：

```asm
0001B732  or byte ptr [di], 40h      ; di ＝ 被推開的那個
```

它有**兩個**作用：

### 1.1 擋掉同一幀的第二次對調

`loc_1B56D` 的四道檢查裡有 `[di] & 0x61`，bit 6 就在其中
（[`../re/11`](../re/11-tactical-battle.md) §5.16）。

### 1.2 ⭐ 讓那個兵**跳過自己這一幀的更新**

每幀的單位迴圈 `sub_1ADC8` 從 `si = 0` 每 `0x20` 走到 `0xC00`
（96 個單位 ＝ 兩側各 48），開頭就測這個位元：

```asm
0001ADE7  mov  ax, [si]
0001ADE9  cmp  al, 80h
0001ADEB  jb   loc_1AE00
0001ADED  test al, 40h               ; ★ 這一幀被換走了
0001ADEF  jnz  loc_1AE26             ; ★ 直接跳到重畫——不移動、不攻擊
…
0001AE1A  call sub_1AF69             ; 移動（撞到人就在裡面打）
…
loc_1AE26:
0001AE26  call sub_1B240             ; 重畫；尾端 and [si],0BFh 清掉 bit 6
0001AE29  call sub_1AEA9             ; 計數
```

所以旗標的壽命是「**從被換走那一刻，到自己這一格輪到為止**」，
而輪到時做的事是**跳過**，不是照常更新。

⚠ 這一點與 [`../re/11`](../re/11-tactical-battle.md) §5.16 原本只寫的
「擋同一幀第三個人再來換」不同——`0001ADED` 那兩行是第二個作用。

## 2. 為什麼少了它會卡死

攻擊在原版是**移動撞出來的**：走進敵人佔的格子就打
（`loc_1B5A1` → `sub_1B618`），走進自己人佔的格子就對調。
一個兵一幀只走一步，所以「打得到」要靠**輪到自己時就站在敵人旁邊**。

少了 §1.2 的話，被推開的兵下一幀馬上又往前擠，把剛換到前排的同伴
再換回去。前排那一格一直換人，而換進去的那個永遠不是
「下一個輪到更新」的那個——量到的樣子是 48 個兵圍著 2 個打，
95 萬 tick 一次傷害都沒有，逐 tick 追一支圍攻的兵，
它在自己更新時**永遠**站在距離 2，而同一時刻的全域快照裡站在距離 1。

有了 §1.2：被推開的兵閒置一幀，剛換到前排的那個旗標是乾淨的，
下一幀輪到它就從前排踏進敵人那一格 → 打。同一場戰鬥
從「95 萬 tick 不動」變成 **5,000 tick 內分勝負**。

## 3. 演算法

```
每個兵的一幀：
    if 旗標.被換走:
        清掉旗標；只做重畫與姿勢翻面        # 不移動、不攻擊
        return
    …鎖敵、換令、命令分派、移動…
    清掉旗標；姿勢翻面
```

## 4. remake 實作

| 項目 | 位置 |
|---|---|
| 旗標 | `internal/rules/tactical/tactical.go` 的 `Soldier.Swapped` |
| 設定 | `soldier.go` `swapWith`：`other.Swapped = true` |
| 擋第二次對調 | `swapWith` 開頭的 `other.Swapped` 檢查（§1.1，已有）|
| **跳過一幀** | `soldier.go` `updateSoldier` 開頭：旗標立著就清掉、翻姿勢、`return`（§1.2）|

## 5. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestSwappedSoldierSkipsItsTurn`（`internal/rules/tactical`）：旗標立著那一幀不動、旗標要清掉、下一幀恢復移動 |
| 迴歸 | `TestTacticalBattleAlwaysResolves`（`internal/state`）：開場體力改成士氣之後，一場 48 對 48 的野戰在 5,000 tick 內分出勝負（之前 95 萬 tick 不動）|
| 對原版 | 側欄的兩格計量條（[`../playtest/40`](../playtest/40-tactical-parity.md) §10）|

## 6. 相鄰的那一格

同一個迴圈緊接著測 `[si+1]`——挨打之後的硬直計時。
兩者接力擋掉三幀，見 [`63`](63-hit-stun.md)。

## 2026-11-12續片驗證

[正式115拍比較](../../workplace/scratch/ui-windows-20261104/cao-swap-after-comparison.json)
及[獨立核對](../../workplace/scratch/ui-windows-20261104/cao-swap-evidence.json)通過96槽27欄位、完整RNG及腳本PC／wait／cond。
上限由本輪原版115個實際入口決定，不把計數器116當成完整第116拍。
七套件冷測／vet／建置、四錯版負對照及正常滑鼠66點、15圖通過，含途中存讀檔、自然入戰、改線、停手移動及正常退卻。
[來源收據](../../workplace/scratch/ui-windows-20261104/cao-swap-integration.json)保存正式程式、受測三檔來源及工具雜湊；
[本機頁](../../workplace/scratch/ui-windows-20261104/cao-swap-game-review.html)展示四張實際流程畫面，不作整場同狀態或像素對拍聲明。
舊spec／RE勘誤與不受影響引用見[回填表](../../workplace/scratch/ui-windows-20261104/cao-swap-backlinks.json)，收尾腳本檢查原始位址與勘誤標記。

本片沒有新增未解項。高平面動態樣本不足為證據限制，後段戰鬥仍由spec159與既有Issues管理。
<!-- 缺口：無 -->
