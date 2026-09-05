# 134 — 尋路是波數佇列，不是先進先出：成本高的格子要等波數追上才展開

**狀態：CONFORMED。** 原版 `loc_1BD46` 的擴散把「成本 ≥ 目前波數」的格子
**推回佇列**，等波數追上才展開；remake 用純 FIFO 直接展開，於是
**繞路成本 8 完全不起作用**——路一律是直線，撞到大將就永久卡住。
改成波數佇列之後，同一場第 70 拍**96 個兵逐槽全等**，小地圖 8 px → 0 px。

- 日期：2026-09-05
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/orig/dosv/KI.EXE`）
  `loc_1BD46` 擴散主迴圈 `0001BDF1`–`0001BE49`、四個方向的
  `sub_1BF2A`／`sub_1BF4D`／`sub_1BF70`／`sub_1BF97`；
  記錄在 [`../re/11`](../re/11-tactical-battle.md) §5.15
- 推論等級：confirmed（機器碼 ＋ 原版執行時的軌跡，見 §4）

## 1. 原版做什麼

擴散用一條環狀佇列，另外帶一個**波數** `dx`（起始 2）。
每彈出一格先看它的成本：

```asm
loc_1BDF1:  mov bx, [di] / inc di / inc di / and di, 47FFh   ; 彈出
            cmp bx, cs:word_1BD44 / jz loc_1BE4A             ; ★ 彈到終點就收工
loc_1BE00:  cmp dx, [bx]
            ja  short loc_1BE0E                              ; 波數 > 成本 → 展開
            mov [si], bx / inc si / inc si / and si, 47FFh   ; ★ 否則推回佇列
            jmp short loc_1BE38
```

四個方向各一支，以 −X（`sub_1BF2A`）為例：

```asm
cmp dx, [bx-2] / jnb locret          ; 波數 ≥ 鄰格成本 → 不動
mov al, es:[bx+2000h]                ; 那一格的佔用成本（有兵 ＝ 8）
xor ah, ah / add ax, dx              ; ★ 新成本 ＝ **波數** ＋ 佔用成本
mov [bx], ax
mov [si], bx / inc si / inc si       ; 推進佇列
```

波數的推進在迴圈尾端：

```asm
loc_1BE38:  cmp di, cx / jnz loc_1BDF1   ; 這一波還沒消化完
            inc dx                       ; ★ 波數 +1
            mov cx, si                   ; 新的本波結尾 ＝ 目前的寫入位置
            cmp si, di / jnz loc_1BDF1   ; 佇列還有東西
            → 失敗（stc）
```

⭐ **這是 dial queue（桶佇列）版的 Dijkstra**，不是 BFS：
一格的成本是 `波數 + 佔用成本`，而它要等到 `波數 > 自己的成本` 才輪得到展開。
所以踩過一個有兵的格子，等於**在佇列裡多躺八波**——繞路真的比較快。

⚠ **早退條件（彈出終點就收工）本身沒問題**，錯的是排序：
在純 FIFO 之下，直線路徑一定比繞路先碰到終點。

## 2. 演算法

```
成本[起點] = 1；佇列 = [起點]；波數 = 2；本波結尾 = 佇列尾
迴圈：
    格 = 彈出()
    若 格 == 終點 → 回溯
    若 波數 <= 成本[格] → 推回佇列（這一波不展開）
    否則 對四個方向的鄰格 nb（通行位元允許者）：
            若 波數 >= 成本[nb] → 跳過
            成本[nb] = 波數 + 佔用成本(nb)
            前驅[nb] = 格
            推進佇列(nb)
    若 讀指標 == 本波結尾：波數 += 1；本波結尾 = 佇列尾
    佇列空 → 走不到
```

回溯不變（只有轉彎才記一個點，上限 64 點）。

**佔用成本**照舊：有兵的格子 8，其餘 0（`sub_1B240` 維護，`docs/re/63` §1）。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 規則層 | `internal/rules/tactical/path.go` 的 `FindPathForcing`（擴散那一段） |
| 差異 | `force`（可破壞地形當成通行）仍是 remake 加的，見 §1 的既有註記；波數與早退照原版 |

改動只在擴散的排序，回溯、`MaxWaypoints`、平面切換、`force` 都不動。

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 原版執行時軌跡 | dosgolem `ticks:4` ＋ `units:1/0` 逐拍取樣同一場（呂布軍 #35 攻許昌 82）：第 0 隊的位 2 在第 12→20 拍走 `(61,29) → (60,31) → (59,33)`，**繞出那一欄**再回到 `(60,34)`；位 3 走到 `(62,31)`。見 [`../playtest/75`](../playtest/75-pathfind-detour.md) |
| 單元測試 | `internal/rules/tactical/path_test.go` 的 `TestFindPathGoesAroundAnOccupiedCell`（大將擋路 → 路徑不經過那一格） |
| 逐兵對拍 | 同一場第 70 拍，96 個兵**逐槽全等**（先前 1 個差 3 格），小地圖 8 px → 0 px |
