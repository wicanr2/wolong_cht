# 141 — 全軍退卻之後的 120 拍倒數

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

## 5. 未解

| 項目 | 現況 |
|---|---|
| 退卻中要不要補兵 | 原版下令之後場上八拍歸零、**沒有補兵進場**；remake 的 `reinforce()` 會補（補進來的兵下一幀被 `applySquadLeaderGone` 改成退卻）。兩邊最後都會結束，但**中途的場上人數不同**，沒有逐拍對過 |
| `word_1D31C` 的兩個 byte | 量到開場是 48／48（＝六隊 × 八人），[`../re/11`](../re/11-tactical-battle.md) §5.9 寫的是「含畫面外待機的」。**兩種讀法都還沒有直接證據**，這一份只用到「它歸零時結束」這一點 |

<!-- 缺口：兩項，見上表 -->
