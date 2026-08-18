# 63 — 被打中的兵有硬直：當幀 ＋ 之後兩幀都不動

**狀態：CONFORMED。** 計時器、遞減點與清旗標的時機都在機器碼裡，
remake 三格都接上了。

- 日期：2026-08-18
- 出處：`sub_1B618`（`0001B66E`–`0001B682`、`0001B69D`）、
  `sub_1B6BC`（`0001B6F0`）、`sub_1ADC8`（`0001ADF1`–`0001ADFE`）
- 推論等級：**confirmed（靜態）**

## 1. 原版做什麼

命中之後，`sub_1B618` 對**被打的那一個**寫五個欄位：

```asm
0001B66E  mov byte ptr [di+1],  2      ; ★ 硬直計時 ← 2
0001B672  and byte ptr [di+2],  1      ; 清掉 +0x02 其餘位元
0001B676  or  byte ptr [di+2], 10h     ; bit 4 ＝ 剛被打中
0001B67A  mov byte ptr [di+5],  0      ; 面向歸零（0 ＝ 西）
0001B67E  mov byte ptr [di+13h], 8     ; +0x13 ← 8
0001B682  or  byte ptr [di],   40h     ; bit 6 ＝ 這一幀不動
```

打到敵方大將的 `sub_1B6BC` 寫同一組（`0001B6F0` 起，少了 `+0x13`）。

### 1.1 ⭐ 四幀，兩個機制接力

每幀的單位迴圈 `sub_1ADC8` 對活著的兵依序測兩件事：

```asm
0001ADED  test al, 40h                ; bit 6（被換走／剛被打中）
0001ADEF  jnz  loc_1AE26              ;   → 只重畫（bit 6 在那裡清掉）
0001ADF1  and  ah, ah                 ; ah ＝ [si+1]
0001ADF3  jz   loc_1AE1A              ;   0 → 正常更新（移動、攻擊）
0001ADF5  dec  byte ptr [si+1]
0001ADF8  jnz  loc_1AE26              ;   還沒歸零 → 只重畫
0001ADFA  and  byte ptr [si+2], 1     ; ★ 歸零那一幀才清掉受擊旗標
0001ADFE  jmp  loc_1AE26              ;   這一幀也還是只重畫
```

所以挨一下之後：

| 幀 | 擋住它的 | 做什麼 |
|---:|---|---|
| n | bit 6（[`62`](62-swapped-unit-skips-its-turn.md)）| 只重畫，bit 6 清掉 |
| n+1 | `[si+1]` 2 → 1 | 只重畫 |
| n+2 | `[si+1]` 1 → 0，同時清掉「剛被打中」| 只重畫 |
| n+3 | — | 恢復移動與攻擊 |

⭐ **「剛被打中」的旗標不是每幀清的**，是**硬直結束那一幀**才清。
它同時是換位的擋條件之一（`[di+02] & 0x10`，
[`../re/11`](../re/11-tactical-battle.md) §5.16）——所以挨打的兵
在硬直期間也**不會被同伴換走**，站著挨完這三幀。

### 1.2 陣亡走同一個計時器

血扣到 0 時改寫成另一組：

```asm
0001B697  and byte ptr [di],   10h     ; 只留 bit 4
0001B69A  or  byte ptr [di],    1      ; bit 0 ＝ 陣亡
0001B69D  mov byte ptr [di+1],  4      ; 計時 ← 4
0001B6A1  mov byte ptr [di+3],  0
```

bit 7 被清掉，於是迴圈走另一條（`loc_1AE00`）——那是倒地動畫，
數完 4 幀由 `sub_1B4B8` 收掉。**本規格不含這一段。**

## 2. 演算法

```
命中(受害者):
    受害者.硬直 ← 2
    受害者.剛被打中 ← 真
    受害者.面向 ← 西
    受害者.這一幀不動 ← 真

每個兵的一幀:
    if 這一幀不動:  清掉它；只重畫；return          # 62
    if 硬直 > 0:
        硬直 −= 1
        if 硬直 == 0: 剛被打中 ← 假
        只重畫；return
    …鎖敵、換令、命令分派、移動…
```

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 硬直計時 | `internal/rules/tactical/tactical.go` 的 `Soldier.Stun`（原版 `+0x01`）|
| 幀數常數 | `damage.go` 的 `HitStunFrames = 2` |
| 設定 | `damage.go` `applyHit`／`hitGeneral`：`e.Stun = HitStunFrames` |
| 遞減與清旗標 | `soldier.go` `updateSoldier`，接在 `Swapped` 那一段後面 |
| 面向歸零 | `e.Facing = West`（`West == 0`，與 `[di+5] = 0` 相同）|

⚠ `Hurt` **不能**在每幀開頭無條件清掉——它要撐到硬直結束，
換位的擋條件才擋得住。

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestHitStunSkipsThreeFrames`（`internal/rules/tactical`）：挨打後三幀不動、第四幀恢復、`Hurt` 撐到最後一幀才清 |
| 迴歸 | `TestTacticalBattleAlwaysResolves`（`internal/state`）：48 對 48 的野戰在 5,000 tick 內分出勝負 |
| 對原版 | 第 61 步還沒接觸，所以逐區對拍的數字不動（[`../playtest/40`](../playtest/40-tactical-parity.md) §10）|

## 5. 未解

| 項目 | 現況 |
|---|---|
| `+0x13` ← 8 | `sub_1B618` 寫、`sub_1B6BC` 不寫。那個欄位誰讀還沒查 |
| 倒地動畫（§1.2）| 4 幀之後 `sub_1B4B8` 收掉，remake 直接把 `Alive` 設成 false |
