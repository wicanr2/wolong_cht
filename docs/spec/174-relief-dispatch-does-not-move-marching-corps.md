# 174 — 求援調兵只調**站著**的軍團，而且只寫意圖

**狀態：CONFORMED。** `sub_14155` 比的是軍團 `+0x0E`，行軍中那是連結記錄位址，
所以調不到；而且它只寫 `+0x20`／`+0x23`，不碰行軍目標。remake 兩件都做反了，
結果是 AI 軍團**走八格就被彈回起點，永遠到不了**。

- 日期：2026-09-09
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_14155`（IDA 線性 00014155）、
  `sub_14057`（00014057）、`sub_1439D`（0001439D）、`sub_14548`（00014548）
- 推論等級：confirmed（逐條反組譯）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §29
- 相關：[`170`](170-ai-target-takes-two-ticks.md)（同一種病的第一個案例）、
  [`164`](164-relief-request-gates.md)、[`172`](172-corps-march-fields.md) §1

## 1. 原版做什麼

```asm
sub_14155:                       ; cl ＝ 目標據點、ch ＝ Stage、
                                 ; dl ＝ 要派幾支、dh ＝ 這一格有幾支
        sub     dh, dl           ; dh ← 可以略過幾支
        mov     di, 2240h        ; 軍團表
        mov     bx, si
        shr     bx, 1
        shr     bx, 1            ; bx ← 據點編號 × 8
loc_14160:
        cmp     bx, [di+0Eh]     ; ⭐ 比的是 +0x0E
        jnz     short loc_1418E
        cmp     byte ptr [di], 80h
        jb      short loc_1418E
        and     dh, dh
        jz      short loc_14179
        call    sub_1ECE0        ; 略過額度還在 → 抽一次
        cmp     al, 40h
        jnb     short loc_14179
        dec     dh
        jmp     short loc_1418E  ; 25% 機率略過這一支
loc_14179:
        test    byte ptr [di], 4 ; 委任
        jz      short loc_1418A
        cmp     byte ptr [di+23h], 8
        jnb     short loc_1418A
        mov     [di+20h], cl     ; ⭐ 只寫意圖
        mov     [di+23h], ch     ; ⭐ 與 Stage
loc_1418A:
        dec     dl
        jz      short locret_14193
loc_1418E:
        add     di, 40h
        jmp     short loc_14160
```

呼叫端 `sub_14057` 傳的是 `ch = 0`：

```asm
loc_14099:
        and     al, al
        jnz     short loc_1409F
        mov     al, 1
loc_1409F:
        mov     cl, ss:[di]      ; 目標據點
        mov     ch, 0            ; ⭐ Stage ← 0
        mov     dl, al           ; 要派幾支
        mov     dh, [si+858h]    ; 這一格有幾支
        call    sub_14155
```

### 1.1 `+0x0E` 的兩種語意，這裡只有一種對得上

軍團 `+0x0E` 停在據點上是**據點編號 × 8**（< 0x800），行軍中是
**連結記錄的位址**（≥ 0x800，[`172`](172-corps-march-fields.md) §1）。
`cmp bx, [di+0Eh]` 的 `bx` 恆 < 0x800，所以**行軍中的軍團永遠不相等**。

⇒ 求援調得動的只有「此刻站在那一格」的軍團。這不是最佳化，是語意：
援軍是把守軍調去別處，不是把半路上的軍團拉回來重新出發。

### 1.2 只寫意圖，落實是下一拍的事

`mov [di+20h], cl` ＋ `mov [di+23h], ch` 之後就沒有別的了——
`+0x14`／`+0x16`／`+0x18` 要等這支軍團下一次被 `sub_125A3` 更新時，
走「現在節點 ＝ 目標節點」那一支到 `sub_14325` → `sub_1439D` →
`sub_14548` 才設。與 [`170`](170-ai-target-takes-two-ticks.md) 是同一條規矩。

## 2. remake 差在哪

```go
// internal/state/strategy.go：dispatchGarrison
gs[i] = threat.Garrison{
    At:    cp.Node,                    // ⚠ Node 是**出發據點**，不是所在地
    ...
}
for _, i := range threat.Dispatch(...) {
    _ = w.March(i, target)             // ⚠ March 會寫 TargetNode 並重算路徑
}
```

兩個錯疊起來的後果是一個**閉環**：

1. 軍團 19 從據點 129 出發往 122，`Node` 一路留在 129。
2. 192 拍後據點 129 又輪到求援機率路徑，`dispatchGarrison` 把它當成
   「還站在 129 的守軍」。
3. `w.March(19, 122)` 從 `c.Node`（129）重算整條路徑 ⇒ 下一步跳回路徑的第一格。

實測：軍團 19 在拍 1,687 出發，走到第 8 格（298,160），
拍 1,879 回到 (293,165) 從頭走，拍 2,071 再一次。**永遠到不了。**

⭐ 這個缺陷**不會讓任何測試變紅**：軍團活著、座標合法、走在道路圖上、
存檔欄位也都寫得出來。只有同局面逐拍對拍看得見。

## 3. remake 要改什麼

**一、`At` 要是「此刻站在哪」。**

```go
at := cp.Node
if cp.LinkAddr != 0 {
    at = -1     // 行軍中：原版的 +0x0E 是連結記錄位址，比不上任何據點
}
```

**二、只寫意圖與 Stage，不呼叫 `March`。**

```go
for _, i := range threat.Dispatch(gs, site, want, skip, next) {
    c := &w.Corps[i]
    c.Ordered = target      // sub_14155：mov [di+20h], cl
    c.Stage = StageNormal   // sub_14155：mov [di+23h], ch，呼叫端傳 0
}
```

| 項目 | 位置 |
|---|---|
| 狀態層 | `internal/state/strategy.go`：`dispatchGarrison` |
| 差異 | 無（照原版）|

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 軍團 19 的軌跡 | `rng_pace.go -corps 19`：出發之後**單調前進**，拍 2,023 抵達 122，不再每 192 拍彈回 |
| 同局面對拍的軍團表 | 拍 1,900 從 12 個 byte／2 支降到 **4 個 byte／1 支**（其餘由 [`169`](169-march-has-no-city-stub.md) §3.1.1 與 [`173`](173-corps-flag-bit0-and-sprite-fields.md) 補完到 0）|
| 單元測試 | `internal/state` 全綠 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| `sub_14057` 的 `dl` | 用的是「威脅目標的索引」（亂數 & 3，0 當 1），不是威脅量算出來的支數。看起來像原版的怪癖，照抄 |
