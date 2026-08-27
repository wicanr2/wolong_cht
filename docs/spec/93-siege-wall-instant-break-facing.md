# 93 — 攻城「一撞歸零」的面向常數要跟著戰場翻轉

**狀態：CONFORMED。** 面向改成跟著戰場擺法走之後，城壁從第 148 幀開始
每幀掉 1–3 點、門強度量表跟著亮起，攻城對拍的 `field` 從 4.32% 降到 0.86%。

- 日期：2026-08-27
- 出處：[`../re/11`](../re/11-tactical-battle.md) §5.11 的
  `seg000:B5EB` 那一段（`cmp cs:byte_10D35, 80h / jnb .flipped`）
- 推論等級：confirmed（兩個分支的面向常數直接寫在機器碼裡）

## 1. 原版做什麼

兵撞上城壁／門時走 `seg000:B5B7`–`B5FE`。攻城戰有一條「耐久直接歸零」的捷徑，
而它**按玩家站哪一側分成兩個分支**：

```asm
cmp cs:byte_1D34B, 0 / jnz .shared        ; 只有攻城戰
cmp cs:byte_10D35, 80h / jnb .flipped     ; bit 7 ＝ 玩家是守方
  cmp si, 600h / jnb .out                 ; 守方碰不壞城壁
  cmp byte ptr [si+5], 0 / jnz .shared    ; ★ 面向 0
  mov word ptr [di+18h], 0
.flipped:
  cmp si, 600h / jb .out                  ; 守方碰不壞城壁
  cmp byte ptr [si+5], 2 / jnz .shared    ; ★ 面向 2
  mov word ptr [di+18h], 0
```

面向的四個值由四支移動常式寫定（`sub_1B047`／`1B069`／`1B08B`／`1B0AF`）：

| 值 | 同時做的事 | 方向 |
|---:|---|---|
| 0 | `dec [si+6]` | West（X−）|
| 1 | `dec [si+8]` | North（Y−）|
| 2 | `inc [si+6]` | East（X+）|
| 3 | `inc [si+8]` | South（Y+）|

⭐ **兩個分支合起來就是「背對城」。** 攻城圖的城壁一律在 X 33–46，
玩家的陣形原點固定在 X 5、腳本側在 X 58（[`56`](56-battlefield-rotation.md) §1）：

| 局面 | 戰場 | 攻方出發 | 城在 | 朝城 | 背對城 | 原版比的值 |
|---|---|---:|---|---|---|---:|
| 玩家攻城 | 不翻 | X 5 | X 33–46 | East | **West** | 0 ＝ West |
| 玩家守城 | 轉 180 度 | X 58 | X 17–30 | West | **East** | 2 ＝ East |

也就是說這個常數不是「某個固定方向」，是**戰場擺法的函數**。

## 2. 演算法

```
背對城的面向 away = (玩家是攻方) ? West : East

hitStructure(side, facing, x, y):
    ...
    if 攻城戰 and facing == away:
        耐久 = 0        # 接著落到「耐久 0 → 整排一起垮」
```

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 規則層 | `internal/rules/tactical/structure.go`：`awayFromCastle` 常數改成看 `b.PlayerSide` 的函式 |
| 差異 | 無——這一條是把原版的第二個分支補上，不是 remake 差異 |

remake 的 `Sides[0]` 永遠是攻方，而戰場翻不翻綁的是**玩家站哪一側**
（`battlesetup.BuildField` 的 `rotate` ＝ 玩家守城），所以判準用
`b.PlayerSide`：等於 `DefenderSide` 就是翻過的圖。

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestInstantBreakFacingFollowsRotation`（`internal/rules/tactical`）|
| 對原版 | [`../playtest/49`](../playtest/49-parity-retest-20260827.md) §3.1：原版 `e10.png` 那一刻**城壁完好、量表亮著**，而修正前的 remake 在第 148 幀就把城壁打垮 |

## 5. 未解

| 項目 | 現況 |
|---|---|

<!-- 缺口：無 -->
