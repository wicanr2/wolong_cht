# 134 — 戰術尋路：波數佇列、反向擴散與成本回溯

**狀態：CONFORMED，限本片12筆尋路與前40拍抽驗。2026-11-12補回反向擴散與成本回溯；原波數排序的歷史收據保留。**

- 日期：2026-11-12

## 2026-11-12 證據審查與實作契約

輸入為松崗DOS/V `KI.EXE`，SHA-256
`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
工具為IDA Pro9.4，位址均為IDA DOS/V linear；
[原指令及bytes](../../workplace/scratch/ui-windows-20261104/path-20261112/ida-probe.json)與
[資料庫身分](../../workplace/scratch/ui-windows-20261104/path-20261112/database.sha256)由
[區間工具](../../tools/ida_range.py)經`tools/ida.sh probe`匯出，不改正式資料庫或符號。

| 原始定位 | 已證實契約 |
|---|---|
| `0001AEF6–1AF1C` | AX為士兵目前XY，DX為GoalXY |
| `0001BD54–1BD6E`、`1BDD3–1BDD5` | 由DX算擴散根，AX存入停止位置1BD44；成本1放在目標，搜尋朝士兵位置展開 |
| `0001BDF9–1BDFE`、`1BE4A` | 停止格一出隊就返回，不等其佔用成本追上；回溯預算為當時波數減2，不能改成停止格成本減1 |
| `0001BE5A–1BE71` | 回溯先試−X、+X，再試−Y、+Y，鄰格成本必須等於預算；不讀擴散時的前驅 |
| `0001BE75–1BEAE` | 優先沿目前軸；換軸時才寫當前XY中繼點，連同終點最多64個 |
| `0001BEE0–1BF13` | 沒有相符XY成本時可經門格換平面，否則遞減預算再找；獨立Z點保留，下降的負高度由1B00D夾零 |

原版執行器dosgolem `f79eab6db8c64a524a06535e53e5d32b0bbcf7c8`，正常第一劇本曹操對呂布、
固定seed1；只在首次RND前設258-byte RNG，沒有直接開戰、注入部隊或強制結果。
新增1BD46／1AF44／1AF4C唯讀觀測，重播至16拍取樣點，17個共同快照含指令數、96筆記錄及RNG與上輪相同。
[動態證據](../../workplace/scratch/ui-windows-20261104/cao-path-evidence.json)核對12筆尋路，
11筆成功的完整路徑與1筆失敗均與可丟棄候選一致。
第一筆由(1,44)往(5,43)，原版與候選均回傳(2,44)、(2,43)、(5,43)，
舊Go為(5,44)、(5,43)。成本圖的根在(5,43)，同時有動態與靜態證據，不從圖片反推。

實作只改正式單兵使用的明確終點平面路徑：反向擴散後依上述成本與軸順序取點。
保留既有導航、佔用成本、終點平面選擇、每拍兩筆FIFO及force fallback；不猜補本片發現的其他表格差異。
`goalPlane=-1`只供既有通用helper及測試，保留原先前向搜尋／前驅回溯的相容語意；
正式`computePath`明示終點平面，不走此分支。兩種模式共用既有波數擴散，不複製另一套搜尋器。

驗收包括加權小地形的三個精確轉角、上／下平面的獨立Z點、不可攀登與無路情境，
舊錯版須被回歸測試拒絕。候選前40拍、每拍96槽27欄位及完整RNG均與原版相同；
正式程式另需冷測戰術、世界、桌面、手機，並由正常滑鼠入戰、調陣形線、下令與返回。
存檔格式不變。不宣稱完整戰鬥、所有地形或導航表已對齊。

## 實作與驗收結果

正式規則的[逐兵比較](../../workplace/scratch/ui-windows-20261104/cao-path-after-comparison.json)
通過前40拍、每拍96槽27欄位及完整RNG；[尋路證據](../../workplace/scratch/ui-windows-20261104/cao-path-evidence.json)
核對12筆請求，11筆成功完整點列與1筆失敗皆相同。第4拍槽29中繼點已修正。
[冷測／vet／建置](../../workplace/scratch/ui-windows-20261104/cao-path-tests.log)通過戰術、世界、桌面及手機四套件，
[負對照](../../workplace/scratch/ui-windows-20261104/cao-path-negative.log)確認舊前驅回溯會被精確轉角測試拒絕。
獨立Z點、下降門格、不可攀登及無路情境亦通過。

正式程式[正常滑鼠流程](../../workplace/scratch/ui-windows-20261104/cao-path-remake-gui.json)
通過66個功能／同步檢查點、15圖，涵蓋曹操開局、宣戰呂布、編夏侯淵、停手行軍、途中存讀檔、
自然入戰、改陣形線、突擊、主動退卻返回。沒有直接開戰或強制勝負。
[接線收據](../../workplace/scratch/ui-windows-20261104/cao-path-integration.json)
封存三檔來源、正式binary及工具雜湊；[本機頁](../../workplace/scratch/ui-windows-20261104/cao-path-game-review.html)
四張流程畫面通過390寬離線抽驗，沒有外部請求。自然GUI不當作同狀態整場對拍。

## 未解與限制

| 項目 | 目前證據與範圍 |
|---|---|
| 完整導航與佔用成本的語意映射 | 第一筆全低平面表有586個導航旗標及26個成本值不同，原始成本含0、8、50、100；尚未逐項追寫入端與建立語意映射。本片只驗已到達的前40拍及12筆尋路，不由此開啟全表重寫或宣稱全場parity |

## 歷史記錄：波數排序修正

原波數排序片段曾標為CONFORMED。原版 `loc_1BD46` 的擴散把「成本 ≥ 目前波數」的格子
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
