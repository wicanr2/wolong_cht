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

## 2026-11-16：求援編成的候選與失敗寫入

**狀態：CONFORMED，限本節編成鏈。** 原版 `140C9 → 14575 → 145C1 → 16E8F` 的求援編成
不要求勢力已設定侵攻目標。`145CE–145EA` 依勢力、職務 0、武術嚴格大於目前最高值
選將；初值為 0，因此全為 0 就失敗，同值保留先掃到者。此處沒有另查存活或舊主欄。
`16EA0` 在試配兵源之前寫軍團 `+2`，兵源不足時仍保留該 byte；不建立活軍團或扣兵。
`16F0A` 逐槽寫入 `+0x2A + 槽號 × 4` 的兵種，後面的槽失敗時不還原已寫的兵種。
試配的兵源只在堆疊副本扣除。追加 192 個不足六槽的固定案例，舊 Go 有 64 案漏寫，
原版完整區塊、佔用圖與 RNG 證據見
[部分配兵原版閉包](../../workplace/scratch/ai-relief-20261116/partial-original-proof.json)與
[修正前結果](../../workplace/scratch/ai-relief-20261116/partial-old-result.json)。此項先達 READY 再實作。

- 輸入：松崗 DOS/V `KI.EXE`，SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 工具與定位：沿 [re/128](../re/128-c-city-tick-restoration.md) 及
  [re/117](../re/117-c-ruler-sortie-restoration.md) 的 IDA Pro 9.4 匯出，12 函式逐筆核對
  原始 bytes 與 MZ relocation 後的載入指令；位址為 IDA database linear。
- 已證實：dosgolem 固定每案初始狀態與完整 258-byte RNG，執行求援與獨立記錄求援地點的 caller。
  四劇本、侵攻目標有／無、容量、挑將、零武術、部分兵源、冷卻與需求共 768 案。
  比較完整 22,208-byte 區塊、128 個原始軍團槽、98,304-byte 佔用圖與完整 RNG。
  舊 Go 有 204 案不同；僅依上述原指令訂正的可丟棄候選實作全等。
- 正式範圍：`requestRelief` 的原版編成鏈、部分失敗寫回與正常存讀檔，不改資料格式、
  工具列或強化 AI。`dest < 0` 的測試探針仍需有侵攻目標，避免把它當成正式求援入口。
- 兩個原始欄位邊界案例用於辨別額外過濾，不宣稱正常遊戲會留下「存活旗標清除但仍有武術」
  或「無職務但有舊主欄」的記錄。正常玩家與長程驗證另列。

本機[原版閉包](../../workplace/scratch/ai-relief-20261116/original-proof.json)、
[舊 Go 結果](../../workplace/scratch/ai-relief-20261116/old-result.json)與
[候選結果](../../workplace/scratch/ai-relief-20261116/candidate-result.json)保存輸入、原始定位與限制。
證據足以進入正式實作；判斷時不擲新骰或挑選通過樣本。

正式來源的 768＋192 案全等。四個額外過濾／失敗寫回的反例皆拒絕，追加配兵案例
另拒絕修正前版本。被俘者依原版保持職務 4，避免挑將移除額外過濾後再次編入俘虜，見
[spec143](143-general-duty-field.md) 的同日補驗。全專案冷測、vet、建置與存讀檔檢查通過；
最終程式的正常曹操滑鼠 62 項／9 圖通過編成、停手行軍、途中存讀檔、自然野戰、守陣、
退卻與戰後存讀檔。不把局部狀態注入當成正常玩家驗收。
來源、受測程式、正常流程與限制保存於
[整合收據](../../workplace/scratch/ai-relief-20261116/integration.json)。

長程沿 Issue #10：原版已固定來源推進至 196/7/16；正式 Go 的六表在 5/1、5/16 全等，
6/1 起仍不同，見[六日期結果](../../workplace/scratch/ai-relief-20261116/long-parity-result.json)。
全區塊診斷另保留敗走與失效軍團記錄差異，不因六表的活軍團篩選而宣稱全區塊一致。

## 7. 未解

| 項目 | 現況 |
|---|---|
| AI 自動判定戰鬥的觸發時機 | 對拍第一個分歧在拍 1978（remake）／2119（原版），差 141 拍。成因未解 |
| ~~`sub_145C1` 挑武將的規則~~ | 已由本頁 2026-11-16 固定原版矩陣核對；正值最高武術、同值先掃到者與失敗寫回均已閉合 |
