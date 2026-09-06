# 81 — 指令列的反白涵蓋八格，狀態列提示框接上財政與編成

**狀態：通過。** 用 dosgolem 拍原版的「進言」「財政」「編成」三條流程：
**進言那一張五區裡四區 0 px**（只剩 `banner` ＝ 日期），
財政與編成的**狀態列提示框各自對到 0 px／只差一則校訂**。
一次抓到兩個缺口：指令列的反白只接了三格、左下角的提示框只接了行軍。

- 日期：2026-09-06
- 規格：[`../spec/124`](../spec/124-menu-highlight-xor.md) §3.5（反白涵蓋八格）、
  [`../spec/140`](../spec/140-status-message-box.md)（狀態列提示框）
- 出處：`sub_161CA`（`000161CA`）逐行讀完 —— `call sub_10B46` **夾住**
  `call cs:funcs_161FE[bx]`，而那一段不分索引
- 原版側：`WOLONG_DOSGOLEM_GAMEDIR=dosgolem/root-saveb tools/dosgolem.sh
  workplace/dosgolem/advise "wait;click:320,200;click:300,151;until:196/4/20;
  sclick:352,15;save:cmd;stap:48,47;steps:400000;shot:a1-advise;restore:cmd;
  stap:144,47;steps:400000;shot:b-finance;restore:cmd;stap:192,47;
  steps:400000;shot:c-form"`
- remake 側：`tools/parity_shot.sh out.png -direct -scenario 0 -player 0
  -seed 7 {-advise-menu|-open-finance|-open-form} -cam 0,0 -shot-frames 1`

## 1. 抓到的兩個缺口

| # | 症狀 | 成因 | 規格 |
|---|---|---|---|
| 1 | 指令列的「進言」「財政」「編成」**不會反白** | `activeCommandCell()` 只認得三張彈出選單。原版是 `sub_161CA` 在 `call` 前後各 XOR 一次，**八格共用同一段程式碼** | [`../spec/124`](../spec/124-menu-highlight-xor.md) §3.5 |
| 2 | 左下角的**狀態列提示框沒有畫**；財政那一張反而多了一條 remake 自己的提示條 | `sub_18853` 的 48 個呼叫點只接了行軍那一條 | [`../spec/140`](../spec/140-status-message-box.md) |

⭐ **兩個缺口都不是「畫錯」，是「沒畫」**——而沒畫的東西在畫面上不會
留下任何線索。抓到它們的唯一方式是把原版那一張擺在旁邊。

## 2. 結果

### 2.1 進言（五項選單，TALK #77）

| 區 | 修之前 | 修之後 |
|---|---:|---:|
| `banner` | 116 | 116（日期：原版 4月20日、remake 4月1日）|
| `command` | **13,642（98.68%）** | **0** |
| `map` | 0 | 0 |
| `minimap` | 0 | 0 |
| `faction` | 0 | 0 |

⭐ **選單框本身一開始就是 0 px**——`openAdvise` 畫的框、字、反白與原版
逐點相同。差的 98.68% 全部是「指令列根本沒開、那一格也沒反白」。

### 2.2 財政（狀態列 TALK #16）

| 區 | 修之前 | 修之後 | 剩下的是什麼 |
|---|---:|---:|---|
| `command` | 13,643 | **88** | x 144–157 y 47–60 ＝ **原版自己畫的滑鼠游標**（14×14，不可消）|
| `map` | 35,492 | **195** | 見下表 |
| `minimap`／`faction` | 0 | 0 | — |

`map` 那 195 px 分成三群，**三群都有解釋**：

| 群 | px | 是什麼 |
|---|---:|---|
| x 96–126 y 352–365 | 141 | ⭐ **M7 的校訂**：TALK #16「財政**予定**」→「財政**計畫**」（`translations/corrections.json` 的 `id: 16`，`text-error`）|
| x 129–158 y 113–126 | 54 | 資金 73,696（原版 4月20日）vs 74,000（remake 4月1日）|

⭐ **那 141 px 反而是最強的證據**：差異的位置正好落在校訂改過的那兩個字上，
說明 remake 掛的就是 TALK #16 那一則。**校訂是刻意的差異**（`CLAUDE.md` §6），
不是缺陷。

### 2.3 編成（狀態列 TALK #0）

**狀態列提示框逐像素 0 px**（框 `(0,320,256,80)` 單獨比）。

⚠ 其餘區**這一輪不可比**：原版那一張停在武將一覽（4月20日的存檔局面），
remake 的 `-open-form` 走到編成面板（4月1日的新局面），
而且 `-lord-corps` 預設放行讓曹操留在候選裡（[`../spec/76`](../spec/76-lord-not-in-formation.md) 的
remake 差異）。**兩邊不是同一個局面，也不是同一步**——列在 §4。

## 3. ⭐ 三件量出來的事

**一、`-battle-steps` 那一類的坑在這裡也有。** `-advise-menu`／`-open-finance`
／`-open-form` 都**沒有開命令視窗**，於是 `command` 區比的是「指令列 vs 地圖」，
98% 的差異看起來像天大的缺陷，實際上一半是 fixture 沒擺好。
三個 fixture 都補上 `hudSet(hudCommand, true)` ＋ 記下格號。

**二、原版側從新遊戲開局會差一天。** 試過讓原版走 NEW GAME（`click:320,179`
→ `300,158` → `450,128` → `360,280` ×2）好讓兩邊都在 196年4月1日：
`banner` 從 116 降到 48，**但仍不是 0**——原版光是走到主畫面再開指令列就
用掉 5,000 萬道指令 ≈ 一個遊戲日，截圖時已經是 **4月2日**。
`click` 的預設 settle 是六百萬道指令，而 `sclick`／`stap` 沒有 settle 參數。

**三、新遊戲那條路的鏡頭不一樣。** 同一個 `sclick:352,15` 之後，
存檔那條路的鏡頭停在世界原點、新遊戲那條路停在別處——大地圖右半的地形
因此整片不同（`minimap` 11.77%／`faction` 11.84%）。
⭐ **鏡頭是游標拖出來的副產物**（[`../re/84`](../re/84-popup-row-band-and-world-cursor.md) §2），
不是一個可以直接設定的狀態。**要對拍就用同一條路徑取樣。**

## 4. 未解

| 項目 | 現況 |
|---|---|
| 編成那一張的完整對拍 | 要兩邊同一個局面同一步：原版停在武將一覽（4月20日），remake 要載同一份存檔、跑到同一天、停在同一步，而且 `-lord-corps=false` |
| 武將／勢力兩格的反白 | 照公式接了，**沒有原版擷取**。原版那兩格走狀態列提示 ＋ 地圖游標，remake 開的是一覽表（[`../spec/124`](../spec/124-menu-highlight-xor.md) §5）|
| 另外 45 個 `sub_18853` 呼叫點 | 接了行軍、財政、編成三條。其餘要一條一條接、一條一條拍 |
| 原版擷取裡的滑鼠游標 | `command` 那 88 px。原版自己畫的，remake 的截圖模式不畫——與 [`76`](76-battle-talk-parity.md) §4 的 95 px 同一類 |
