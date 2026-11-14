# 141 — 全軍退卻之後的 120 拍倒數

## 2026-11-14：撤退傳令與首次出口初始化（CONFORMED，限本片）

接續上一片已完成的停止補兵，核對六隊長傳令、重複全軍退卻拒絕、大將檢查與首次移動。
原版`1A8F6–1A92D`只寫六隊長Next，已有退卻旗標就整次拒絕；
`1A9D0–1A9FA`由新退卻隊長傳給七隊員，`1AAED–1AB38`首次保存舊命令到Next，
並同時初始化Goal、Step及路徑；繼續退卻時只查出口，不每拍改Y或清路。
原版`1ADDA→1AE56`在尋路出隊與移動前先檢查體力門檻，再扣攻城大將體力。
已有[正常1476拍收據](../../workplace/scratch/ui-windows-20261104/cao-assault-original.jsonl)
顯示1472拍Cmd5／Next7，Step已指出口；目前Go仍有一拍與Next差異。
本片先用受控門檻樣本及typed候選審查，不因終點拍數相同就忽略中間機制。
戰略存檔及玩家指令不改，正常GUI與原版常式受控樣本分開驗證。

輸入`KI.EXE`、dosgolem與位址空間身分沿下片。IDA Pro9.4新匯出為
[傳令／退卻](../../workplace/scratch/ui-windows-20261104/retreat-command-20261114/ida-probe.json)
SHA-256 `a08e1ea4b0efebbb382d78b3f021d58cc9c5ed799bcdd73f59c8e5fbe0cf5125`，
DB SHA-256 `fca6e416d46cbb758b1afd93952268ed36bd86f34200adae999d7d67c39df047`；
[隊長傳令](../../workplace/scratch/ui-windows-20261104/retreat-handlers-20261114/ida-probe.json)
SHA-256 `8674162db14236f5c460ee9a36610aaff3379cbbc2ef77a8e55eb66f2e9b73c9`，
DB SHA-256 `bcafb8c7093c2e52f367ed54892b7a79562e48c9096564f9245824659cfecfb9`。
兩份probe均739函式、UID/GID1000，原始檔唯讀，未改正式DB。

[12案24階段受控比較](../../workplace/scratch/ui-windows-20261104/cao-retreat-order-candidate-fixtures-comparison.json)
全等，涵蓋49／50／51體力、雙方低體力、已有退卻旗標、攻守換位、野戰不扣與門檻檢查後扣體力。
fixture從正常241拍的1AE56快照restore，只改RAM；分別停1ADDD與下一次1ADC8，不能當正常玩家擺位。
首階段比較全96槽HP／Cmd／Next及旗標、倒數；次階段只比進入退卻記錄的命令、位置、Goal／Step。
`1AB21`亦須呼叫既有requestPath去重入隊，沿[spec120](120-pathfind-request-queue.md)及re80契約。
typed候選的正常1471–1476拍，玩家48兵共288兵次的11個命令／位置／目標欄位吻合，終點仍1476拍。
既有PathQueued、退場MoveFlag及敵方後段差異不由這11欄代證；停止線維持完整動畫／世界出口之外。
首次初始化保存舊命令、固定出口Y及直接Step，取代spec94每拍重設Goal／按路徑終點清路的舊推論。
原指令、受控候選與正常撤退窗口審查足以列READY，正式接線後仍須冷測、正常滑鼠及負對照。

[正式24階段比較](../../workplace/scratch/ui-windows-20261104/cao-retreat-order-after-fixtures-comparison.json)全等，
含234筆退卻記錄的PathQueued與TargetRefresh；受控Go Rand固定0，原版快照seed1，
門檻／傳令／出口規則本身不取亂數。對話變體、原版繪圖／音效呼叫及後續亂數不在這組常式比較內。
[正式窗口](../../workplace/scratch/ui-windows-20261104/cao-retreat-order-after-window.json)288兵次11欄全等，
原版及typed重播終點仍1476拍，6000／6000人、體力49／220；正常240拍前綴逐byte不變。
41個正式來源套件冷測、全範圍vet、建置與八類突變拒絕通過。
原腳本查詢測試改從已收到全軍令的隊長及退卻階段觀察分支，未改比較方向、常數或查詢公式。

[正常滑鼠](../../workplace/scratch/ui-windows-20261104/cao-retreat-order-remake-gui.json)71個功能／同步點、17圖通過。
第一劇本曹操宣戰呂布、編夏侯淵、行軍途中存讀檔、自然椎陽攻城、改線／登城／突擊／原有快轉後停手，
725拍自然收尾並返回大地圖。GUI輸入另行記錄，不與原版停手樣本混比戰果。
程式`wlgame-cao-retreat-order-r1` SHA-256 `e57c806f9d81c3c796e094ca63024f84d12bca4e7f0a9e478e29d8593f4741c2`。
本片CONFORMED，完整動畫、195拍後的單兵狀態、退場旗標與完整世界出口仍保留；Issue #16／#18未關閉。

## 2026-11-14：撤退停止補兵與場上計數（CONFORMED，限本片）

玩家可見問題：曹操第一劇本對呂布、椎陽攻城停手到底，原版1476拍收尾，
remake卻補兵到1579拍才由倒數結束。兩側兵力與大將體力相同，不能因此忽略收尾延遲。
沿用原版seed1與正常宣戰、編成、行軍入口，不用直接開戰或修改戰場RAM。
[原版逐拍收據](../../workplace/scratch/ui-windows-20261104/cao-assault-original.jsonl)
與[舊實作重播](../../workplace/scratch/ui-windows-20261104/cao-assault-remake.jsonl)分開保存。

原版輸入`KI.EXE` SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`；
dosgolem `f79eab6db8c64a524a06535e53e5d32b0bbcf7c8`。
新[IDA Pro9.4匯出](../../workplace/scratch/ui-windows-20261104/retreat-finish-20261114/ida-probe.json)
位址空間為DOS/V linear，SHA-256 `f1c09578cbb7cb1d1e81e025bbb5c78e10dc508f224accc3966bd304cee923de`；
一次性DB SHA-256 `0a1fcb27f9d0d48efbc7c476226f2d50ad45dd05123ebe91bb305e8f03da3f47`。
輸出probe確認739函式、UID/GID1000，原始執行檔唯讀。

已證實的窄契約：

- `1B43E–1B443`比較原版角色1／2與全軍撤退旗標；該側撤退就不補兵。另一側仍能補兵。
  個別受傷退卻、隊長死亡本身不等於全軍撤退，不清除待機兵，也不把此閘擴成雙方停補。
- `1ADCF`清場上計數；`1AE29`只在原先在場的單兵分支呼叫`1AEA9`，
  補兵分支`1AE15–1AE18`當拍不計入。待機與已離場生還不進這個計數。
  計數保留依序處理時的狀態，不改成所有移動結束後重新數Alive。
- `1A727–1A746`依玩家、對方順序檢查計數，先歸零者輸；120拍倒數出口仍排在前面。
  remake的攻守編號須先經PlayerSide映射，兩側同時為0時亦如此。

typed輸入沿用`endPhase`與`PlayerSide`；移動批次建立兩側場上計數，再交勝負檢查。
來源軍團／待機資料→Battle→Step→Result→既有世界結算及UI，戰略存檔格式不變。
本片不改大將體力檢查位置、全軍退卻命令初始化、補兵落點與逐槽排程。
先前195拍之後的狀態差異保留，不以整場勝方相同宣稱全程逐拍一致。
驗收須覆蓋雙PlayerSide、單側與全軍退卻、待機保留、雙方同時空場、倒數兜底，
再跑相同原版初始狀態到底及正常滑鼠自然收尾。中間動畫與完整世界出口仍分開驗證。

[候選審查](../../workplace/scratch/ui-windows-20261104/cao-retreat-candidate3-comparison.json)確認
1476拍原版自然出口、兩側各6000人、每隊100兵及大將體力49／220，typed候選四項全等。
原242筆觀測的共通欄位與103筆命令前綴不變，保留第64／88步既有診斷。
初步只在移動後重數Alive的候選不足以保留原版計數時點，改用移動批次逐槽累計。
另行移動大將檢查的實驗得到1477拍，未採用；該命令初始化仍須獨立查證，不能為湊拍數一併猜補。
依原指令、正常原版出口及候選審查，以上窄契約列READY。

[正式比較](../../workplace/scratch/ui-windows-20261104/cao-retreat-after-comparison.json)與候選逐byte相同，
原版及remake皆1476拍結束，勝方、兩側6000人、每隊100兵與大將體力49／220全等。
正常240拍前綴與上一輪逐byte相同；195拍後的已知狀態差異仍保留。
正式來源套件41個冷測、全範圍vet及建置通過；五種突變分別拒絕撤退仍補兵、待機進勝負計數、
移動後重數Alive、攻守順序取代角色順序與移動計數漏接。
原倒數測試曾靠持續補兵撐120拍，已改成硬直阻止退場，正確建立兜底前提。

[正常滑鼠抽驗](../../workplace/scratch/ui-windows-20261104/cao-retreat-finish-remake-gui.json)71個功能／同步點、17圖通過。
第一劇本曹操宣戰呂布、編夏侯淵、停手行軍、途中存讀檔，自然在椎陽進攻城戰；
改陣形線、登城、突擊並使用原有快轉後停手，725拍自然敗退並返回大地圖。
GUI這條有下令，戰果為5600／5660人、體力38／214；它與上方停手到底的對拍是不同輸入，分開記錄。
追蹤僅增加在場、生還、體力及結束時Outcome的觀測，不新增正式遊戲旗標或修改戰果。
戰報維持預設關閉。本片CONFORMED，不關閉涵蓋完整世界出口的Issue #18與戰術全程的Issue #16。

**狀態：CONFORMED。** 原版下全軍退卻令之後起一道 120 拍的倒數，
走完就結束戰鬥、由**沒退卻的那一側**獲勝。remake 先前只有
「一側補不出兵」這一條出口，退卻的兵走不出去就永遠打不完。

- 日期：2026-09-06
- 出處：[`../re/11`](../re/11-tactical-battle.md) §5.9／§5.9.1
  （`sub_1A6FA` ＝ `0001A6FA`、`sub_1A8F6` ＝ `0001A8F6`、
  `sub_19A33` 的初始化、`sub_19FDC` ＝ `00019FDC`）
- 推論等級：**confirmed**（機器碼 ＋ dosgolem 實測一拍減一次）
- 相關：[`135`](135-script-message-command.md)（`byte_1D349` 的另一個用途：
  腳本訊息的閘）、[`128`](128-squad-leader-gone-keeps-reserve.md)（隊長不在場的連鎖退卻）

## 1. 原版做什麼

```asm
sub_19A33（戰鬥初始化）
    mov cs:byte_1D349, al      ; al ＝ 0
    mov cs:byte_1D34A, 78h     ; ★ 120，在這裡設，不是下令時設

sub_1A8F6（全軍退卻）
    cmp cs:byte_1D349, 0 / jnz → stc / retn   ; ★ 已經有一側在退卻就不受理
    al = (di == 0) ? 1 : 2
    mov cs:byte_1D349, al
    bx = di + 1Bh；寫六次 [bx] = 5、每次 bx += 100h   ; 六個隊長的命令欄
    sub_1C315(dl = al − 1)                            ; 播「全軍撤退！！」

sub_1A6FA（每拍，主迴圈）
    cmp cs:byte_1D349, 0 / jz .1
    dec cs:byte_1D34A / jnz .1
    xor al, al / xchg al, cs:byte_1D349
    cmp al, 1 / jnz .0 / mov cs:byte_1D349, 1
.0: call sub_19FDC             ; 結束戰鬥，回傳 al = byte_1D349
```

⭐ **倒數排在「補不出兵」那兩條出口之前**（[`../re/11`](../re/11-tactical-battle.md) §5.9）。

## 2. 演算法

```
戰鬥開始：退卻旗標 ← 0、倒數 ← 120
全軍退卻（玩家按退却、或大將體力不支）：
    退卻旗標 ≠ 0 → 不受理（原版回 CF=1）
    退卻旗標 ← 退卻的那一側 + 1
每拍（在勝負判定的最前面）：
    退卻旗標 ≠ 0 →
        倒數 −= 1
        倒數 == 0 → 結束，勝方 ＝ 沒退卻的那一側
```

⚠ **倒數是在初始化時設的，不是在下令時設的**——照抄。兩者在原版等價
（第二次退卻不受理），寫成「下令時設 120」看起來一樣，但那是我們的簡化。

⚠ **正常打完不會用到它。** 實測退卻令下去之後 **8 拍**場上就清空了
（六隊 × 八人都站在自己那一側的邊緣欄），走的是「補不出兵」那條出口。
倒數是兜底：退不出去時才由它收尾。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 規則層 | `internal/rules/tactical/battle.go`：`RetreatCountdown` ＝ 120、`Battle.retreat`（倒數）、`checkVictory` 開頭那一段 |
| 起倒數 | 既有的 `Battle.endPhase`（`Order(side, −1, Retreat)` 寫 `side + 1`，[`135`](135-script-message-command.md)）——**與原版的 `byte_1D349` 是同一個值**，不另外開一個欄位 |
| 勝方 | `Winner = 2 − endPhase`（`endPhase` 1 ⇒ 勝方 1、2 ⇒ 勝方 0），與原版 `xchg` 那兩行等價 |
| 差異 | 無 |

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestRetreatEndsBattleAfter120Ticks`（`internal/rules/tactical`）：下全軍退卻之後**第 120 拍**結束，勝方是沒退卻的那一側 |
| 單元測試 | `TestRetreatCountdownOnlyRunsWhileRetreating`：沒有人退卻時倒數不動 |
| 單元測試 | `TestSecondRetreatDoesNotRestartCountdown`：第二側也下退卻令，倒數不重來、勝方仍是第一側的對手 |
| 突變測試 | 把 120 改成別的值、把「退卻中才遞減」的閘拿掉，各要有測試變紅 |
| 對原版 | [`../playtest/80`](../playtest/80-retreat-countdown.md)：`word_1D318` 走 8 拍的同時 `byte_1D34A` 從 `0x78` 走到 `0x70` |
| 沒有回歸 | 同一場攻城戰第 70 拍的截圖**改動前後逐 byte 相同**（[`../playtest/80`](../playtest/80-retreat-countdown.md) §4）。那個取樣點沒有人退卻，所以「一個 byte 都不變」才是對的結果——比逐區 0 px 更強 |

## 5. 已解項目

| 項目 | 現況 |
|---|---|
| 退卻中要不要補兵 | 已證實：2026-11-14原指令與正常自然收尾確認，全軍撤退的那一側停止補兵。實作與抽驗範圍見本檔續片 |
| `word_1D31C` 的兩個 byte | 已證實：2026-11-14確認只數移動批次的在場分支。正常收尾時場上0而待機仍有552兵，推翻「含待機」舊解釋；原指令與收據見本檔續片 |

<!-- 缺口：無 -->
