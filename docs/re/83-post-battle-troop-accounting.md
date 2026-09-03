# 83 — 打完之後兵力怎麼結算：三項相加，而待機數從來不會被清成 0

**狀態：confirmed（靜態）。** 戰後寫回軍團的每槽兵數是
**退場生還 ＋ 場上存活 ＋ 待機**三項相加（`sub_19F58` 的
`al = [si+3] + [si+1]`，前者在戰鬥中就累加了）。⭐ 順帶封閉一個負證據：
碰結算緩衝區的**只有 11 支函式**，逐一看完之後可以說
**原版沒有任何地方把某一隊的待機數清成 0**——只有 `dec`（開場取用、補兵）。
remake 的 `squadLeaderGone` 清 `Reserve[squad] = 0`，那是原版沒有的
（[`../spec/128`](../spec/128-squad-leader-gone-keeps-reserve.md)）。

- 日期：2026-09-03
- 原始輸入：`workplace/orig/dosv/KI.EXE`，SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- IDA database：`workplace/ida/dosv/KI.EXE.i64` SHA-256
  `65736f11b0b28a5b3a6db9e1a3d205cc24f0eaebc82b508ee0d7d283f6240572`（739 支函式）
- 工具：`tools/ida_var_writers.py`（`word_1D30A` 的 11 筆參考）、`tools/ida_dump.py`
- 位址：DOS/V 線性位址
- 相關：[`11`](11-tactical-battle.md) §5.6（1 隊長 ＋ 7 隊員）、§5.9（勝負條件）、
  [`../spec/65`](../spec/65-retreated-soldiers-survive.md)（退卻是保命不是損失）

## 1. 結算緩衝區的版面（`word_1D30A` 段）

`sub_19E97` 開場把軍團記錄搬進來，一側 0x20 byte：

```
+0        （搬進來之前先寫 0）
+1 … +7   從軍團記錄 +1 起的 7 byte
+8+4k     每隊 4 byte × 6 隊：
  +8+4k     ← 從軍團記錄搬來
  +9+4k     ⭐ **待機兵數**
  +10+4k    ← 從軍團記錄搬來
  +11+4k    ⭐ **生還累計**（`stosb` 寫 0 ⇒ 開場歸零）
```

side 1 的基址是 `+0x20`（`sub_19C45` 的 `mov di, 20h`）。

三支互相印證同一格：`sub_19C45` 的 `es:[bx+di+9]`、`sub_1B413` 的
`es:[bx+1]`（`bx` ＝ 隊×4＋8）、`sub_19F58` 的 `[si+1]`（`si` 從 8 起每次 +4）
——**位址都是「側基址 ＋ 隊×4 ＋ 9」**。

## 2. ⭐ `sub_19F2C` 數的是「還站在場上的兵」

```asm
sub_19F2C:                          ; 00019F2C，44 B
        mov  cx, 30h                ; 48 個兵記錄
.next:  cmp  byte ptr es:[di], 80h  ; ★ bit 7 ＝ 在場
        jb   .skip                  ;   不在場就不數
        mov  bx, di
        cmp  di, 600h / jb $+2 / sub bx, 600h    ; 換算成側內位移
        shl  bx,1 / shl bx,1 / and bx, 0FC00h / xchg bl,bh   ; ⇒ 隊 × 4
        inc  byte ptr [bx+si+0Bh]   ; ★ 累加到「生還累計」
.skip:  add  di, 20h / loop .next
```

⚠ **它沒有先歸零**，因為那一格開場就是 0（§1），而且戰鬥中
`sub_1B4B8` 已經在同一格累加了——所以 `sub_19F2C` 是**補上最後一批**，
不是重數。

| 誰 | 什麼時候 | 加什麼 |
|---|---|---|
| `sub_1B4B8`（`ah = 0`）| 戰鬥中，兵退到畫面外 | 那一隊 +1（`ah = 1` ＝ 倒地，**不加**）|
| `sub_19F2C` | 戰鬥結束 | 還站在場上的每個兵 +1 |

## 3. 寫回軍團：三項相加

```asm
sub_19F58:                          ; 00019F58，72 B
        add  si, 8 / add di, 28h    ; si → 緩衝區第 0 隊、di → 軍團 +0x28
        xor  ax, ax / mov dx, ax / mov cx, 6
.next:  mov  al, [si+3]             ; ★ 生還累計（退場 ＋ 場上）
        add  al, [si+1]             ; ★ ＋ 待機
        add  dx, ax                 ; 總和
        mov  es:[di+1], al          ; ★ 寫回這一槽的兵數
        add  si, 4 / add di, 4 / loop .next
        ; 總和拿去按比例縮士氣：
        mov  cx, dx / xchg cx, es:[di+4] / jz +
        mov  al, es:[di+6] / mul word ptr es:[di+4] / div cx
+       mov  es:[di+6], al
```

⭐ **最後三行是士氣按兵力比例縮**：新士氣 ＝ 舊士氣 × 新總兵力 ÷ 舊總兵力。
`es:[di+4]` 被 `xchg` 換成新總兵力，所以同一條指令既讀舊值又寫新值。

## 4. ⭐ 隊長不在場是「每幀重新施加」，不是一次性事件

`sub_1A754`／`sub_1A785`（每幀，由 `sub_1A6FA` 呼叫）逐隊檢查：

```asm
cmp byte ptr [si], 80h / jb .noLeader   ; ★ 隊長的 bit 7
call sub_1A85B / call sub_1A7B7          ; 隊長在 → 正常更新
jmp .members
.noLeader: call sub_1A83F                ; ⇒ 七名隊員的命令一律改成 5（退卻）
```

`sub_1A83F` 只做那一件事（[`11`](11-tactical-battle.md) §5.6），
**不動待機數**。所以原版的行為鏈是：

```
隊長倒下
  → 每一幀都把在場的七名改成退卻
  → 他們退到畫面外（sub_1B4B8, ah=0）⇒ 生還累計 +1
  → 那一格空出來，而待機數還是非 0
  → sub_1B413 把待機兵補進場（補兵路徑不看隊長在不在）
  → 下一幀又被改成退卻 …
```

**待機兵最後全部經由「補進場 → 退卻 → 走出畫面」被算成生還。**
帳目上等價於「待機兵直接生還」，但走的是實際的離場常式。

## 5. 未解

| 項目 | 現況 |
|---|---|
| 補兵會不會補到隊長那一格 | `sub_1B413` 只在**待機數為 0** 那條路才判 `test si, 0FFh`（隊長格）；待機非 0 的路徑沒有排除它。remake 的 `reinforce` 是**一律不補隊長格**。要對得先在原版上量「隊長死後那一格會不會出現新的兵」 |
| `+8+4k` 與 `+10+4k` 的語意 | 從軍團記錄 `+0x28+4k` 起搬進來的兩個 byte，`sub_19F58` 只讀 `+1`／`+3`，這兩格戰鬥中沒有讀取端被找到 |
| `word_1D31A` 的語意 | `sub_1AEA9` 數「在場且 `+0x19` 非零」的兵，`sub_1ADC8` 尾段拿它算優勢（`byte_1D31E`）。`+0x19` 是一個倒數計時器（`sub_1AB7C` 寫 `28h`、`sub_1ADC8` 每幀 `dec`），**它代表什麼還沒解** |
