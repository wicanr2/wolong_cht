# 187 — 戰後判誰：攻城只判敗方，野戰誰壞滅就判誰

**狀態：CONFORMED。** 兩個入口的規則**不一樣**：

| 入口 | 規則 |
|---|---|
| `sub_14ADE`（攻城）| 拿 `al`（誰贏）分流——**攻方贏只看 `ah` 位元 1（守方壞滅），守方贏只看位元 0**。勝方即使 `sub_1474A` 回 STC（士氣歸零）也不判 |
| `sub_14A7B`（野戰）| **不看誰贏**，`ah` 哪一位設起來就判誰；**兩邊都壞滅（`ah` ＝ 3）時只判攻方**（`cmp ah,2 / jb` 與 `jz` 都不成立，落到 `loc_14AC1`）|

remake 的 `resolveCorpsBattle` 兩條共用一段「兩邊都判」，
於是攻城時勝方士氣掉到 0 會被多消滅一次。

- 日期：2026-09-10
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_14ADE`（IDA 線性 00014ADE–00014B62）、
  `sub_15130`（00015130–000151B2）、`sub_1474A`（0001474A–000147BA）
- 推論等級：confirmed（三支逐條讀出）
- 相關：[`46`](46-post-battle-retreat.md)（`sub_1474A` 的三個壞滅入口）、
  [`../re/09`](../re/09-combat.md) §4.3、§5

## 1. 原版做什麼

```asm
sub_14ADE:
        …
        call sub_14C72 / jb loc_14B34         ; 找守軍；沒有 → 打城兵
        call sub_14ED7                        ; 有守軍：軍團 vs 軍團
        and  al, al / jnz loc_14B4B           ; ★ al ≠ 0 ⇒ 守方贏
        test ah, 2  / jz  loc_14B41           ; ★ 攻方贏 → 只看守方壞滅
        mov  al, [si+1] / xchg bx, si
        call sub_1291A                        ; 守將的下場
        xchg bx, si
loc_14B41:
        mov  al, [si+1] / pop si / push si
        call sub_14CF3                        ; 據點易主
        jmp  loc_14B56
loc_14B4B:                                    ; ── 守方贏 ──
        test ah, 1 / jz loc_14B56             ; ★ 只看攻方壞滅
        mov  al, [di+1] / call sub_1291A      ; 攻將的下場
loc_14B56: …
```

`ah` 的兩位來自 `sub_15130`：

```asm
        mov cl, al / call sub_1474A / jnb .1 / or ah, 1   ; si ＝ 攻方 → 位元 0
.1:     xchg si, di
        mov cl, al / xor cl, 1 / call sub_1474A / jnb .2 / or ah, 2   ; 守方 → 位元 1
```

⇒ **兩邊都跑 `sub_1474A`**（總兵力重算、移動計時寫 1、Stage 設定都照跑），
但**只有敗方那一位會被消費**。

⭐ 這也修正了 [`../re/09`](../re/09-combat.md) §4.3 的一句話：
「士氣不足 100 的軍團打贏也會散掉」——士氣**確實**歸零
（`sub_151B3` 對勝方一樣清），`sub_1474A` **確實**回 STC，
但那個 STC 在勝方那一側沒有消費端。散不掉。

## 2. remake 錯在哪

```go
attDead := r.AttackerDestroyed || w.retreatOrPerish(att, !r.DefenderWins)
defDead := r.DefenderDestroyed || w.retreatOrPerish(def, r.DefenderWins)
w.afterBattle(ev, att, node, attDead, def, rng)
w.afterBattle(ev, def, node, defDead, att, rng)
```

`retreatOrPerish` 對勝方回 false，但 `r.AttackerDestroyed`／
`r.DefenderDestroyed`（＝ 士氣 0 或大將槽 0）對勝方一樣成立，
於是勝方也被送進 `corpsPerishes`，多擲一次 `RollFate`。

`fightGarrison`（打城兵那條）同樣要改：原版 `loc_14B34` 之後也是
`and al, al / jnz loc_14B4B`，攻方贏就直接易主，不判攻將。

## 3. 驗證

| | 修正前 | 修正後 |
|---|---:|---:|
| 逐拍取數不一致 | 10 / 5,480 | **3 / 5,480（0.1%）**|
| 第一個分歧 | 5,104 | **5,176** |
| 拍 5,176 的據點表 | 103 個 byte／57 座 | **14 個 byte／4 座** |

拍 5,104（軍團 44 攻據點 2，守軍軍團 73）：原版 17 次、remake 18 次
——多的那一次是勝方的 `RollFate`。改完一致。

### 1.1 野戰那一半

```asm
sub_14A7B:
        …
        call sub_14C72 / jb 結束              ; 沒有對手就不打
        mov  di, bx
        and  byte ptr [si], 0DFh / mov byte ptr [si+3], 0    ; 兩邊都清對峙
        and  byte ptr [di], 0DFh / mov byte ptr [di+3], 0
        call sub_14E5C
        and  ah, ah / jz 結束                 ; ★ 沒有人壞滅
        cmp  ah, 2 / jb loc_14AC1             ; ah ＝ 1 → 判 si（攻方）
                     jz loc_14AC9             ; ah ＝ 2 → 判 di（守方）
                                              ; ah ＝ 3 → 落下，仍判 si
loc_14AC1: mov al, [di+1] / call sub_1291A
        jmp 結束
loc_14AC9: xchg si, di / mov al, [di+1] / call sub_1291A / xchg si, di
```

⚠ **`al` 傳的是勝方勢力**（`sub_1291A` 的第四個例外「勝方 ＝ 敗方勢力」
要靠它），所以 `loc_14AC9` 先 `xchg` 再取 `[di+1]`。

## 4. 未解

- 野戰在 `ah` ＝ 3 時為什麼只判攻方，還沒有解釋。
  照字面是 `cmp/jb/jz` 三分支漏掉了第四種情形——**照抄，不要補**。
