# 193 — 士氣或大將槽歸零時，戰後的階段**不設**

**狀態：READY（已實作，但沒有對拍樣本能分開）。** `sub_1474A` 在重算（`sub_16FD2`）之後、
分流之前有兩道提早返回：

```asm
call sub_16FD2                       ; 總兵力／間隔／`+0x0B` ← 1
cmp  byte ptr [si+6], 0   / jz .stc  ; ★ 士氣 0
cmp  byte ptr [si+29h], 0 / jz .stc  ; ★ 大將槽（第 0 槽）兵力 0
and  cl, cl / jz .stay               ; cl ＝ 0（勝方）→ `[si+23h] ← 8`
…
.stay:  mov byte ptr [si+23h], 8
.stc:   stc / retn                   ; ⇒ **`+0x23` 一個字都沒寫**
```

⇒ **打完之後士氣歸零的軍團，`+0x23` 保持原值**。攻城的勝方不判壞滅
（[`187`](187-only-the-loser-is-judged.md)），所以它活著並帶著戰前的
階段——階段 8（等士氣）只在軍團**站到目標上**時才被分派器讀到，
所以這個差別要等它走到目的地才顯現。

remake 的 `retreatOrPerish` 在 `won` 時無條件寫 `StageWaitMorale`。

- 日期：2026-09-11
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_1474A`（IDA 線性 0001474A–000147BA）、
  `sub_14483`（00014483，等士氣）、`sub_14325`（00014325，階段分派）
- 推論等級：**機器碼 confirmed，行為未驗**——自動判定的大將槽保底 1，
  士氣也只在「戰前 < 100」時歸零（`scaleMorale`），所以同局面 7,935 拍裡
  **一場都沒走到這兩條**。改動照機器碼補上，對拍結果不變。
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §48
- remake 實作：`internal/state/aimarch.go` 的 `retreatOrPerish`
- 相關：[`187`](187-only-the-loser-is-judged.md)、[`46`](46-post-battle-retreat.md)

## 1. 為什麼現在看不見

`combat.Destroyed` 已經是這兩個檢查（`docs/re/09` §5），而
`retreatOrPerish` 的 `won` 分支**在它成立時仍然寫 `+0x23`**——
原版不寫。差別只在「壞滅的勝方」這種組合上，而：

- 自動判定的大將槽保底 1（`shrink` 的 `isGeneralSlot`），
- 士氣只在戰前不足 100 時歸零，

所以要湊出這個組合得先有一支士氣 < 100 還去攻城的軍團。
**同局面到 5/25 為止一場都沒有。**

⚠ 這一條留著的理由是「機器碼就是這樣寫的」，不是「觀測要求這樣改」。
真的要驗它，得另外造一個士氣 < 100 的局面。

## 2. remake 的修法

```go
w.recalcCorps(i)
// ⭐ 原版在這裡就可能提早返回（`cmp [si+6],0` / `cmp [si+29h],0`）：
// 士氣或大將槽兵力歸零 ⇒ STC，而且 **`+0x23` 一個字都沒寫**。
if c.Morale == 0 || c.Units[0].Men == 0 {
	return true
}
if won { c.Stage = StageWaitMorale; return false }
```

回傳 `true` ＝ 原版的 STC。攻城的勝方那一側本來就不消費這個旗標
（[`187`](187-only-the-loser-is-judged.md)），所以「回 true」不會
讓勝方壞滅——`fightGarrison`／`resolveCorpsBattle` 都會把它蓋掉。

## 3. 驗證

`tools/parity_ck.sh` 全部取樣點**與改動前相同**（5/20 之前全 0）。
這一條沒有讓任何一個 byte 變好，也沒有變壞。

<!-- 缺口：無 -->
