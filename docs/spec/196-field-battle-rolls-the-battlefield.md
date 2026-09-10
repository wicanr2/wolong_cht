# 196 — 野戰開打前要先決定戰場，而水域那一格**要擲骰**

**狀態：CONFORMED。** `sub_14A7B`（野戰的自動判定）**第一件事**是
`call sub_14B63`——取樣守方腳下與下方四格的地形、決定用哪一張戰場。
正下方那一格的地形類型是 8（水上木構造）時走 `sub_14C1A`，
而那一支會 **`call sub_1ECE0`**：

```asm
sub_14C1A:                      ; bl ＝ 正下方那一格的地形類型
        cmp bl, 9 / jz .fixed   ; 類型 9 ⇒ 固定 0xD5，不擲骰
        call sub_1ECE0
        and  al, 3
        mov  cl, al / add cl, 0D1h    ; ⇒ 戰場 0xD1 + (亂數 & 3)
        xor  ch, ch / retn
```

規則本身 remake 早就有了（`internal/rules/battlefield`，
`NeedsWaterRoll` ＋ `SelectWith`），但**只接在「進戰術畫面」那條路**
（`internal/battlesetup`）。自動判定的野戰完全沒跑，於是**少取一次亂數**，
之後整條亂數流錯位。

- 日期：2026-09-11
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_14B63`（IDA 線性 00014B63–00014BDC）、
  `sub_14C1A`（00014C1A–00014C4B）、`sub_14BDD`（00014BDD）；
  `sub_14A7B` 是唯一的呼叫端
- 推論等級：confirmed（逐條讀出 ＋ 逐拍取數對上）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §49
- remake 實作：`internal/state/state.go` 的 `SetTerrain`、
  `internal/state/corps.go` 的 `fieldAt`
- 相關：[`121`](121-water-battlefield-selection.md)、[`../re/05`](../re/05-battle-selection.md)、
  [`../re/11`](../re/11-tactical-battle.md) §4.4

## 1. 取樣的五格

```asm
mov dx, [di+1Ch] / mov bx, [di+1Ah]     ; ★ **守方**的佔用圖位址 ＝ 座標
…
mov al, [bx+17Fh] / call sub_14C4C      ; 左下
mov al, [bx+181h] / call sub_14C4C      ; 右下
mov al, [bx]      / call sub_14C4C      ; 中心
mov al, [bx+300h] / call sub_14C4C      ; 兩格下
mov al, [bx+180h] / call sub_14C4C      ; ★ 正下方 → bl
and bl, bl        / jz  → sub_14BDD     ; 平原：查配對表，不擲骰
cmp bl, 8         / jnb → sub_14C1A     ; 水域：可能擲骰
（1–7）           →  戰場 ＝ 0CEh + bl
```

一列是 `0x180` byte，所以 `+0x180` 是正下方、`+0x17F`／`+0x181` 是左右下、
`+0x300` 是兩格下——與 `internal/battlesetup` 既有的
`at(0,0)/at(0,1)/at(-1,1)/at(1,1)/at(0,2)` 一致。

⚠ **座標取的是守方的**（`di`），不是攻方，也不是據點。

## 2. 對拍上長什麼樣

同局面拍 8,889（軍團 50 對軍團 34 的野戰）：

```
原版   ['141D8','141F1','14216','14C22','152F6','152F6','151FC','1521B', …]  18 次
remake ['state.go:1768'×3,        'combat.go:107'×2,     'combat.go:212', …] 17 次
```

`14C22` 就是那一擲。少了它，之後每一場戰鬥的傷亡都換一組亂數——
到 5/31（拍 9,177）差 252 個 byte／160 座據點、三支軍團的兵力欄位。

## 3. remake 的修法

規則層不讀檔案，所以地圖由呼叫端注入（與 `SetRoads` 同一個做法）：

```go
w.SetTerrain(func(x, y int) byte { t, _ := lib.World.Tile(x, y); return t })
```

`fieldAt` 在開打之前照 `sub_14B63` 取樣五格，
`battlefield.NeedsWaterRoll(n)` 為真就 `rng.Next()` 一次。

⚠ **沒注入地圖時不擲骰**：那是降級路徑（缺原版素材也要能跑），
與原版不一致但至少可預測——夾具會把「有沒有接上」印出來。

## 4. 驗證

`tools/parity_ck.sh 196/5/31`：據點 252 B／160 座、勢力 2 B、軍團 21 B／3 支
→ **全 0**。

<!-- 缺口：無 -->
