# 164 — 求援的三道閘：`sub_14028` / `sub_140C9` / `sub_140B3`

**狀態：CONFORMED。** 三處都已改（`internal/state/strategy.go`），
`internal/state` 冷跑全綠；同局面對拍沒有回歸。
⚠ 個案取樣（「鄰居敵對、非侵攻目標、佔用 0」的據點看到求援）還沒做，
而對拍顯示求援機率路徑的次數仍與原版差很多（§7）。

- 日期：2026-09-09
- 出處：`KI.EXE`（松崗 DOS/V）`sub_14028`（IDA 線性 00014028）、
  `sub_140C9`（000140C9）、`sub_140B3`（000140B3）、`sub_14575`（00014575）
- 相關：[`../re/40`](../re/40-garrison-relief-request.md) §2／§4、
  [`../re/44`](../re/44-threat-and-reinforcement-ai.md) §2、
  [`163`](163-city-tick-order.md)、[`../playtest/119`](../playtest/119-rng-pace-comparison.md)

## 1. 三個差異

`internal/state/strategy.go` 的 `relieve`／`requestRelief` 與原版有三處不一致。
三處都在同一組常式裡，都由組語直接讀出來。

| # | 原版 | remake 現況 |
|---|---|---|
| ① | 受威脅而且**這一格沒有軍團**就立刻求援，**與「威脅是不是具體目標」無關** | 多了一個 `r.Specific &&` |
| ② | `sub_140B3`（記下求援據點）是**獨立的一次呼叫**，冷卻中照樣執行 | 冷卻中直接 return，`ReliefSite` 沒寫 |
| ③ | **編出至少一支軍團**才寫冷卻（`sub_14575` 回 CF=1 就跳過） | 無條件寫冷卻 |

## 2. ① 立刻求援的條件裡沒有「具體目標」

```asm
sub_14028:
        and     byte ptr [si+840h], 3Fh      ; 清 bit 7／6
        mov     al, 80h
        cmp     byte ptr [bp+0], 0FEh
        jb      short loc_1403E              ; < 0FE：威脅有具體目標
        jz      short loc_14040              ; = 0FE：一般威脅
        mov     byte ptr [si+857h], 0        ; = 0FF：沒有威脅 → 清冷卻
        jmp     short loc_14053              ;          clc（不呼叫 sub_14057）
loc_1403E:
        or      al, 40h                      ; ← 具體與否**只影響 +0x00 的 bit 6**
loc_14040:
        or      [si+840h], al
        cmp     byte ptr [si+858h], 1
        jnb     short loc_14055              ; 佔用 ≥ 1 → stc → 走 sub_14057
        mov     al, 1
        call    sub_140C9                    ; 佔用 = 0 → 立刻求援一支
        call    sub_140B3
loc_14053:
        clc
        retn
```

⭐ **`jb` 與 `jz` 匯流到同一個 `loc_14040`。** 具體目標只決定要不要多設 bit 6，
分支條件從頭到尾只有 `[si+858h]`（停在這一格的軍團數）。

remake 寫成 `if r.Specific && c.Occupancy == 0`，於是
「鄰居敵對但不是本勢力的侵攻目標、而且沒有守軍」這一格**永遠不會求援**——
這正是「電腦被打了也不調兵」的形狀。

## 3. ② 冷卻中仍然記下求援據點

```asm
sub_140C9:
        cmp     byte ptr [si+857h], 0
        jz      short loc_140D1
        retn                                 ; ← 冷卻中，這一支結束
…
sub_140B3:                                   ; ← 但呼叫端是**分開的兩次 call**
        mov     bh, [si+841h]                ; 勢力記錄 +0x16 ＝ 求援的據點編號
        …
        mov     [bx+16h], ah
```

兩個呼叫端（`sub_14028+25`、`sub_14057+3A`）都是
`call sub_140C9 / call sub_140B3` 兩行。所以**冷卻只擋住「發出請求」，
擋不住「記下是哪個據點在求援」**——而勢力記錄 `+0x16` 是 AI 軍團挑目標時讀的。

## 4. ③ 沒編出軍團就不寫冷卻

```asm
sub_140C9  loc_140FF（AI 的據點）:
        …
        call    sub_14575
        jb      short loc_14153              ; ← CF=1：一支都沒編出來 → 不寫冷卻
        …計算離首都的距離…
loc_1414F:
        mov     [di+857h], al                ; 寫冷卻
loc_14153:
        pop     si
        retn

sub_14575:
        bl = max(5, 資金 >> 13)
        sub     bl, [si+14h]                 ; 減掉現有軍團數
        ja      short loc_14594
        stc / jmp 結束                       ; 額度用完 → CF=1
loc_145A5:
        call    sub_145C1                    ; 編一支
        jb      short loc_145B7              ; 沒有可用的武將 → 跳出
        inc     ah
        …
loc_145B7:
        clc
        and     ah, ah
        jnz     short loc_145BD
        stc                                  ; 一支都沒編出來 → CF=1
```

冷卻的語意是「剛剛真的叫到援軍了，先別再叫」。**沒叫到就不該進入冷卻**，
否則據點會在最需要援軍的時候（額度用完、沒有閒置武將）反而閉嘴 30 拍。

## 5. remake 要改什麼

`internal/state/strategy.go`：

```go
// relieve
if !r.Threatened { c.ReliefCooldown = 0; return nil }
if c.Occupancy == 0 {                      // ← 拿掉 r.Specific &&
    return w.requestRelief(site, 1, rng)
}
if len(r.Targets) == 0 { return nil }
```

```go
// requestRelief
f := &w.Factions[c.Owner]
defer func() { f.ReliefSite = site }()     // ← sub_140B3：冷卻中照樣寫
if c.ReliefCooldown != 0 { return nil }
…
formed := 0
for n := threat.Budget(f.Funds, f.Corps); n > 0 && want > 0; n-- {
    if w.formAICorpsTo(c.Owner, site) == nil { break }
    formed, want = formed+1, want-1
}
if formed == 0 { return nil }              // ← sub_14575 CF=1：不寫冷卻
c.ReliefCooldown = threat.AICooldown(…)
```

## 6. 怎麼驗

1. `internal/state` 既有測試全綠（`ai_probe_test.go` 直接呼叫 `refreshCityThreat`）。
2. 取數節拍對拍：`tools/parity_pace_diff.py` 的不一致比例要下降，
   而且**前段的完全一致區間不能縮短**（現況：前 1,977 拍逐拍相同）。
3. ① 的直接證據：找一個「鄰居敵對、非侵攻目標、佔用 0」的據點，
   改後要出現求援（AI 側 ＝ 編軍團，玩家側 ＝ TALK #38）。

## 7. 未解

| 項目 | 現況 |
|---|---|
| AI 自動判定戰鬥的觸發時機 | 對拍第一個分歧在拍 1978（remake）／2119（原版），差 141 拍。成因未解 |
| `sub_145C1` 挑武將的規則 | 取 `+0x11` 最大且 `+0x17 == 0` 的那一位；remake 的 `formAICorpsTo` 是否同序未驗 |
