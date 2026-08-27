# 98 — 爬不上去的那一下要打門：未破的門是這樣被打開的

**狀態：CONFORMED。** remake 的 `tryClimb` 碰到未破的門就直接回 false，
於是登城的兵站在門格上**永遠卡住**——門一點傷都沒有（耐久 80 全程不變）。
補上「爬不上去就打它」之後，站上門格的兵會把門打開。
⚠ 據點 82 那條 fixture 的行為沒有因此改變：那一隊走不到門格
（[`../playtest/52`](../playtest/52-siege-timeseries-parity.md) §6）。

- 日期：2026-08-27
- 出處：`sub_1B0D3`（`0001B0D3`–`0001B115`）＋ `sub_1B186`
  （`0001B186`–`0001B1B0`），本輪 dump
- 推論等級：**confirmed**

## 1. 原版做什麼

```asm
sub_1B0D3:                          ; 往上爬
  di = bx & 0FFFh
  al = 圖塊[di]                     ; 腳下這一格
  cmp al, 0F0h / jb 失敗            ; ★ 不是門 → 根本不能爬
  [si+5] = 面向（跟著門的奇偶）
  bx += 1000h
  call sub_1B186                    ; 上一層站得住嗎
  jb  loc_1B109                     ; 站不住
  … 成功：[si+0Ah]++、[si+0Ch]=bx、[si+1Eh]=10h
loc_1B109:
  and al, al
  jz  失敗                          ; ★ al == 0 → 就是失敗
  call loc_1B533                    ; ★ al != 0 → 走**碰撞處理**
```

```asm
sub_1B186:                          ; 上一層站不站得住
  al = es:[bx+1000h]
  test al, 7Fh / jnz loc_1B1AD      ; ★ 上一層有**實體** → al = 實體編號，失敗
  dl = 圖塊[bx & 0FFFh]
  xor al, al
  cmp dl, 0F8h / jb loc_1B1AD       ; 圖塊 < 0xF8（未破的門）→ al = 0，失敗
  clc / retn                        ; ≥ 0xF8 → 成功
loc_1B1AD:
  al &= 7Fh / stc / retn
```

⭐ **兩種失敗要分開看**，`sub_1B0D3` 的 `and al, al` 就是在分它們：

| 上一層 | `al` | `sub_1B0D3` 的反應 |
|---|---:|---|
| 有實體（門／城壁本體）| 實體編號 ≠ 0 | **`loc_1B533` 碰撞處理 → 打它** |
| 空的，但腳下的圖塊是未破的門 | 0 | 什麼都不做 |

**所以未破的門是這樣被打開的**：登城的兵走到門格上 → 試著爬 →
上一層就是那道門 → 走碰撞處理，每一幀扣一點耐久。門的耐久只有 **80**
（城壁是（城兵數＋50）×10，動輒上千），很快就破，接著就爬得上去。

⭐ **這也解釋了原版的門強度條為什麼一直不出現**：那個條的呼叫端有兩道閘，
其中一道是 `cmp [di+1], 1` ＝ **只有城壁（類型 1）才顯示**
（[`32`](32-gate-strength-bar.md) §2）。打門不會亮條——
與原版七張影格量到的一致（[`../playtest/52`](../playtest/52-siege-timeseries-parity.md) §2）。

## 2. remake 現況

```go
if !s.CanClimb() || !b.Field.IsGateCell(s.X, s.Y) ||
    b.Field.GateBlocksHighPlane(s.X, s.Y) {
    return false
}
```

三個條件擠在同一個 `return false` 裡，**未破的門與「不是門」被當成同一件事**。
量到的後果（據點 82、`SAVE-L.DAT`、呂布 35 攻）：

| | 值 |
|---|---|
| 四道門的耐久 | 80 → **80**（全程沒被打過）|
| 城壁 | 1,660 → **0**（被磨穿）|
| 站上牆頂的人數 | **0** |
| 門強度條 | 一路顯示（原版是一次都沒有）|

## 3. 演算法

```
tryClimb：
    不能爬的兵種 或 腳下不是門格   → 失敗
    上一層有實體                   → **打它**（hitStructure），失敗
    上一層空著但門還沒破           → 失敗（什麼都不做）
    其餘                           → 爬上去
```

remake 的地形模型裡「上一層有實體」＝ 那一格有還沒破的結構物
（`breakableAt`）——未破的門正是這一種。

## 4. remake 實作

| 項目 | 位置 |
|---|---|
| 規則層 | `internal/rules/tactical/soldier.go` 的 `tryClimb` |
| 差異 | 無 |

## 5. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestClimbIntoUnbrokenGateHitsIt`（`internal/rules/tactical`）。做過正對照——**拿掉那一下會紅** |
| 對原版 | [`../playtest/52`](../playtest/52-siege-timeseries-parity.md)：門強度條要跟原版一樣不出現 |

## 6. 未解

| 項目 | 現況 |
|---|---|
| `loc_1B533` 的完整分流 | 這裡只用到「撞到結構物」那一條。它同時也是敵我碰撞的入口（[`../re/11`](../re/11-tactical-battle.md) §5.16），另外兩條在水平移動那邊已經接了 |
