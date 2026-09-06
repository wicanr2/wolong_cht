# 149 — 行軍目標是在大地圖上點，不是一覽表

**狀態：CONFORMED。** 原版下行軍指示時，選目標據點走的是**大地圖選點**
（`sub_1703C`：畫一個空心框游標，等玩家點地圖上的據點），
remake 開的是一張 192 列的據點一覽。兩者的畫面、操作與命中規則都不同。

- 日期：2026-09-06
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）
  - `sub_17FDB`（`00017FDB`）：`mov cx, 3 / call sub_18853` → **`call sub_1703C`**
  - `sub_1703C`（`0001703C`）：`sub_11F7F`（畫游標）→ `sub_1E453`（熱區查詢）
    → `sub_1709E`（那一格是不是據點）的迴圈
  - `sub_1709E`（`0001709E`）：兩道閘，見 [`../re/85`](../re/85-march-target-hit-test.md) §1
- 推論等級：**confirmed**（三支逐行讀完 ＋ 原版擷取 `orig-mar3` 五區全 0 px）
- 驗證：[`../playtest/94`](../playtest/94-march-map-picker.md)
- 相關：[`39`](39-march-order-menu.md) §1.1（下一步的三選一）、
  [`../re/85`](../re/85-march-target-hit-test.md)（命中判定與格座標）、
  [`../re/45`](../re/45-corps-command-mode.md) §5（`sub_1703C` 的迴圈）

## 1. 原版做什麼

```
sub_17FDB:
  loop:
    cx = 3 / sub_18853            ; 狀態列 #3「請指示行軍目標之據點。」
    call sub_1703C                ; ★ 大地圖選點；CF=1 ＝ 右鍵取消 → 整條結束
    …                             ; 命中 → TALK #21 ＋ 三選一（docs/spec/39）
    取消三選一 → jmp loop          ; 回去重選據點
```

`sub_1703C` 的迴圈每圈做三件事：

1. `sub_11F7F`：算游標的格座標（寫進 `word_1989A`／`word_1989C`）並畫游標。
2. `sub_1E453`：查游標壓在哪個**熱區**。⭐ **`al ≠ 0` 就走熱區分支，
   這一圈根本不問據點**（[`../re/85`](../re/85-march-target-hit-test.md) §3）。
3. `al == 0`（在地圖上）→ `sub_1709E`：那一格是不是據點。

`sub_1709E` 的兩道閘：

| 閘 | 條件 |
|---|---|
| 一 | 那一格的**圖塊編號落在 `0CBh`–`0D3h`**（9 種據點圖塊）|
| 二 | **據點表登記的 `+0x08` X 與 `+0x0A` Y 與那一格完全相等** |

⭐ **一座城的圖形有 4×4 格，但只有登記的那一格按得到**——
「點圖示中間沒反應」是原版行為，不是 bug
（[`../playtest/67`](../playtest/67-dosgolem-popup-menus.md) 記過這個症狀）。

### 1.1 游標的形狀（原版擷取 `orig-mar3` 逐像素）

大地圖上的游標是 **15×15 的白色圓角空心框 ＋ 往右下偏 1 px 的黑影**，
畫在那一格的左上角（格 (col,row) → 螢幕 `((col−camX)×16, 32+(row−camY)×16)`）。

```
白（色 15，(243,243,243)）：      黑影（同形，(+1,+1)）
 .#############.
 ##...........##
 #.............#
 …（共 11 列）
 #.............#
 ##...........##
 .#############.
```

⚠ **只在選點期間畫。** 原版在別的狀態下畫什麼游標還沒對上
（[`../playtest/91`](../playtest/91-general-boast-parity.md) §3：三個假說都被對照推翻了），
所以這一份只接「選點期間」這一個已經有原版擷取的狀態。

### 1.2 時間停著

`sub_1703C` 的迴圈不推進世界。實測：選點狀態下再跑 900,000 道指令，
遊戲時鐘還是 196年4月17日 6時。

### 1.3 選完軍團先開**軍團情報面板**

`sub_17FDB` 的唯一呼叫端是 `sub_17F90+21`：

```asm
sub_17F90:
    sub_1807B / sub_1812A          ; ★ 兩種情況都先畫軍團情報面板
    al = cs:byte_10CFF
    cmp al, [si+1] / jnz loc_17FBC ; 是玩家的軍團嗎
    call sub_17FDB                 ; ★ 是 → 行軍指示（狀態列換成 #3）
    call sub_14325 / or [si], 2
    jmp  loc_17FC8
loc_17FBC:
    cx = 4 / call sub_18853        ; ★ 不是 → 狀態列 #4，只是看
    call sub_121E7
loc_17FC8:
    call sub_1817D                 ; 擦掉面板
    cx = 0FFFFh / call sub_18853   ; 清狀態列
```

⭐ **狀態列 #4 只在「別人的軍團」時掛**——自己的軍團直接進行軍指示，
狀態列是 #3。[`140`](140-status-message-box.md) §1.1 先前把 #4 記成
「自己的軍團」，反了。

⭐ **面板在選點期間一直留著**，而且那時候它**不吃輸入**：右鍵是選點的取消
（`sub_1703C` 回 CF=1 → `loc_18046` → `sub_1817D` 擦面板 ＋ 清狀態列）。

### 1.4 面板上的數字用原版的 8×16 字模

原版擷取逐像素量到：`1` 只有 4 px 寬、`0` 6 px，與一覽表的數字欄同一套
（`sub_1062F`）。remake 先前用文字字型的數字（`1` 是 6 px），整串位置就差開。

## 2. remake 要怎麼改

| 項目 | 作法 |
|---|---|
| 選點狀態 | `cmd/wlgame/mappick.go` 的 `mapPickState`：`pick func(city int)`；**`pick` 是 nil 就是沒在選**——零值安全 |
| 進入 | `cmd/wlgame/corps.go` 的 `pickDestination` 改成開選點，不開一覽表 |
| 游標 | `drawMapPickCursor`：§1.1 的兩層遮罩，畫在游標所在格的左上角 |
| 格座標 | `camX + x/16`、`camY + (y − strategyMapY)/16`。⚠ remake 的鏡頭本來就是**以格為單位**，沒有原版那個「兩次捨去各吃一格」的問題（[`../re/85`](../re/85-march-target-hit-test.md) §2）|
| 命中 | `cityAtTile(col, row)`：據點的 `X`／`Y` **完全相等**才算。閘一（圖塊編號）不必另外查——圖塊是據點就一定在表裡，反之亦然 |
| 熱區優先 | 游標在指令列／小地圖／勢力欄上時**不問據點**，與原版同序（`sub_1E453` 先問） |
| 取消 | 右鍵 ＝ 整條流程結束（`sub_1703C` 回 CF=1 → `sub_17FDB` 的 `loc_18046`）|
| 時間 | `timeRuns()` 加一條：選點期間停時間（§1.2）|
| 軍團情報面板 | `showCorpsPanel(corps)` 只畫面板不碰狀態列；`openCorpsInfo` ＝ 面板 ＋ #4（別人的軍團）。選點期間 `updateCorpsInfo` 不吃輸入，「ESC 關閉」那行 remake 提示也不畫 |
| 面板的數字 | `drawOriginalNumber`（§1.4）|
| 指令格 | `commandFlowRunning` 要把選點算成「流程還在跑」，否則 `syncCommandFlow` 會在選點那一刻把反白與狀態列一起收掉 |
| 殘留的選單框 | 進選點時 `closePopupMenu()`——原版是地圖一重畫就沒了（[`126`](126-command-popup-menus.md) §1.2）|
| 驗收旗標 | `-open-march-pick`：編一支軍團並停在選點狀態。配 `-pick-tile X,Y` 把游標釘在指定格——**headless 沒有指標**，同 `hideAmountCursor` 的理由 |

### 2.1 拿掉的 remake 便利

原本的一覽表**照切比雪夫距離排序**（`pickDestination`），
[`../playtest/56`](../playtest/56-lubu-flow-parity.md) 標記過那是 remake 便利。
改成地圖選點之後那張表就不存在了，連帶那條差異一起消失。

⚠ **這是把 remake 拉回原版，不是新增差異。** 一覽表看得到全部 192 個據點、
地圖選點只看得到畫面內的——這是原版的操作成本，照抄。

## 3. 驗證

| 方式 | 內容 |
|---|---|
| 對原版 ✅ | [`../playtest/94`](../playtest/94-march-map-picker.md)：`orig-mar3` **五區全 0 px**，含原版自己畫的游標 |
| 單元測試 | `TestCityAtTileNeedsExactRegisteredTile`：登記格命中、旁邊三格落空 |
| 單元測試 | `TestMapPickCursorMask`：兩層遮罩的尺寸，與「黑影不是白框位移一格」的那兩格 |
| 單元測試 | `TestMapPickTileAt`：螢幕座標換格，地圖區以外回 false |
| 單元測試 | `TestMapPickKeepsCommandCellLit`：選點期間指令格還亮著 |
| 單元測試 | `TestMapPickStopsTime`：選點期間 `timeRuns()` 回 false |
| 突變測試 | 把命中改成「4×4 範圍內都算」，第一支要變紅 |

## 4. 未解

| 項目 | 現況 |
|---|---|
| 游標推到畫面邊緣時鏡頭跟過去 | 原版的滑鼠座標是**世界座標**，推到視野外鏡頭會捲（[`../re/84`](../re/84-popup-row-band-and-world-cursor.md) §2）。remake 的滑鼠被視窗框住，還沒接這個行為 |
| 別的狀態下的游標 | [`../playtest/91`](../playtest/91-general-boast-parity.md) §3 還沒對上，所以只接選點這一個 |
