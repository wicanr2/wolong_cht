# 141 — 全軍退卻之後的 120 拍倒數

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
