# 170 — AI 挑到新目標要**兩拍**才會出發：`+0x20` 是意圖、`+0x14` 是行軍目標

**狀態：CONFORMED。** `aiStage2` 改成只寫意圖、`aiStage0` 與玩家的到站
處理都先 `retarget`（`sub_14548`）；同局面對拍的軍團 39 **逐拍與原版相同**，
逐拍不一致 1,730 → 1,279。

- 日期：2026-09-09
- 出處：`KI.EXE`（松崗 DOS/V）`sub_1440F`（IDA 線性 0001440F）、
  `sub_1439D`（0001439D）、`sub_14548`（00014548）、`sub_14325`（00014325）、
  `sub_12662`（00012662）、`sub_14300`（00014300）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §21
- 相關：[`../re/64`](../re/64-corps-arrival-state-machine.md)、
  [`../re/65`](../re/65-ai-march-decision-chain.md)、[`169`](169-march-has-no-city-stub.md)

## 1. 兩個欄位不是同一件事

| 欄位 | 意思 | 誰寫 |
|---|---|---|
| `+0x20` | **意圖**：AI 決定要去哪個據點 | `sub_1440F`（AI 挑待辦）、`sub_17FDB`（玩家下令）|
| `+0x14`／`+0x16`／`+0x18` | **行軍目標**：節點 × 8、X、Y | `sub_14548`，由各 Stage handler 呼叫 |

`sub_1440F` 只寫意圖，**不動行軍目標**：

```asm
loc_14446:
        cmp     [si+20h], al
        jz      short loc_14451
        mov     [si+20h], al          ; ⭐ 只寫 +0x20
        or      byte ptr [si], 2      ; 位元 1 ＝ 下一步要重算（docs/re/34 §2.1）
loc_14451:
        mov     byte ptr [si+23h], 0  ; Stage → 0
        dec     byte ptr [di+18h]
```

落實是下一拍的事。`sub_14325` 依 Stage 分派，AI 的 Stage 0 是 `sub_1439D`：

```asm
sub_1439D:
        mov     bx, ax                ; ax ＝ +0x20 指的據點記錄位址
        call    sub_14548             ; ⭐ 把 +0x14／+0x16／+0x18 設成那個據點
        jb      short loc_143A5       ; CF=1 ＝ 已經站在那裡
        retn
loc_143A5:
        test    byte ptr [bx], 40h    ; 據點 +0x00 位元 6（威脅有具體目標）
        jnz     short locret_143AE
        mov     byte ptr [si+23h], 1  ; Stage → 1
```

```asm
sub_14548:                            ; bx ＝ 據點記錄位址
        mov     ax, [bx+8]  / mov dx, [bx+0Ah]      ; 據點的 X、Y
        sub     bx, 840h / shr bx,1 / shr bx,1      ; 據點編號 × 8
        cmp     [si+10h], ax / jnz .set             ; 已經在那裡？
        cmp     [si+12h], dx / jnz .set
        cmp     [si+0Eh], bx / jnz .set
        stc / jmp .write
.set:   clc
.write: mov [si+16h], ax / mov [si+18h], dx / mov [si+14h], bx
```

⚠ `sub_14370`（**玩家**的 Stage 0–3）第一件事也是 `sub_14548`——
兩半張分派表都靠它落實。

## 2. 實測：三個週期

同一份快照，軍團 39（勢力 0，據點 82 → 88）。`sub_125A3` 每拍只更新 16 支、
間隔 3，所以這支每 **24 拍**才輪到一次：

| 拍 | 旗標 | 朝向 | 現節點 | `+0x14` 目標 | `+0x20` 意圖 | Stage | 座標 |
|---:|---|---:|---:|---:|---:|---:|---|
| 31 | `C0` | 4 | 82 | 82 | 82 | 2 | (206,114) |
| **32** | **`C2`** | 4 | 82 | **82** | **88** | **0** | (206,114) |
| 56 | `C2` | 4 | 82 | **88** | 88 | 0 | (206,114) |
| **80** | `C0` | 3 | 488 | 88 | 88 | 0 | **(206,116)** |

- 拍 32：`sub_1440F` 寫意圖 88、設位元 1、Stage → 0。**行軍目標還是 82。**
- 拍 56：`+0x0E == +0x14`（都是 82）⇒ 走「到站」分支 ⇒ `sub_14325` ⇒
  `sub_1439D` ⇒ `sub_14548` 把行軍目標設成 88。**仍然不動。**
- 拍 80：`+0x0E ≠ +0x14` ⇒ 移動分支 ⇒ 清位元 1、`sub_147BB` 選路 ⇒
  `sub_12708` 走出第一格。

⇒ **從「決定去哪」到「踏出第一步」跨三個更新週期。**

## 3. remake 現況

`internal/state/aimarch.go`：

```go
// aiStage2 ＝ sub_1440F
if c.Ordered != dest {
    _ = w.March(i, dest)      // ⚠ March 同時寫 Ordered 與 TargetNode
}

// aiStage0 ＝ sub_1439D
func (w *World) aiStage0(i int) {
    c := &w.Corps[i]
    if w.citySpecific(c.TargetNode) { return }
    c.Stage = 1               // ⚠ 少了 sub_14548
}
```

兩處合起來讓 remake **少一拍**：拍 32 決定、拍 56 就走。
軍團 39 的格子序列因此整條早一個週期
（[`169`](169-march-has-no-city-stub.md) §6）。

## 4. remake 要改什麼

**一、`aiStage2` 只寫意圖。**

```go
if c.Ordered != dest {
    c.Ordered = dest          // ⭐ 不碰 TargetNode——那是下一拍 sub_14548 的事
}
c.Stage = StageNormal
```

**二、`aiStage0` 先落實再判斷。** `sub_14548` 的語意是
「把行軍目標設成 `Ordered` 指的據點，並回報現在是不是已經站在那裡」：

```go
func (w *World) aiStage0(i int) {
    c := &w.Corps[i]
    here := w.retarget(i, c.Ordered)   // ← sub_14548
    if !here {
        return                          // 還沒到，這一拍只設目標
    }
    if w.citySpecific(c.Node) { return }
    c.Stage = 1
}
```

**三、玩家的到站處理（`sub_14370`）也要先 `retarget`。**
`arriveCorps` 走玩家那一半時第一件事同樣是 `sub_14548`。

`retarget` 是新的共用小函式（`sub_14548`）：設 `TargetNode`／`TargetX`／`TargetY`、
重算 `routes`，回傳「軍團的 `Node`／`X`／`Y` 是不是已經等於那個據點」。

## 5. 怎麼驗

1. 同局面對拍：軍團 39 在拍 32／56 都不動、拍 80 才走到 (206,116)，
   之後逐拍與原版相同。
2. `internal/state` 既有測試維持綠燈。
3. 逐拍取數的第一個分歧點要往後推。

## 6. 未解

| 項目 | 現況 |
|---|---|
| 軍團 `+0x00` 位元 1 | remake 沒建模。它的效果（下一次移動前重查道路表）被 remake 的「`March` 當場算好 routes」涵蓋，但存檔寫回是靠 `modelledCorpsBits` 原樣保留的，不是真的維護 |
| `sub_14325` 分派表的 16 項 | 玩家半張與 AI 半張各 8 項，remake 只對到 0–3 與 8/10/11；其餘未逐支對過 |
