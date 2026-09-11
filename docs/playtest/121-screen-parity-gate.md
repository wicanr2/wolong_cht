# 121 — 十七組畫面逐區重量：十四組符合預期，缺口剩三個

**狀態：通過（帶三個已登記的缺口）。** 把散在二十幾份 `playtest/` 的畫面
對拍組收成一份可重跑的資料（`tools/parity_screens.json`）＋一支閘
（`tools/parity_screens.sh`），**每一組都在當前程式碼上重截重量**，
不引用任何文件裡的舊數字。

- 日期：2026-09-11
- 規格：[`../spec/198`](../spec/198-ivent-scene-frame.md)（本輪新增並修好）、
  [`../spec/197`](../spec/197-shot-mode-draws-no-desktop-cursor.md)、
  [`../spec/90`](../spec/90-same-state-parity.md)、[`../spec/91`](../spec/91-tactical-parity.md)
- 重跑：`tools/parity_screens.sh`（全部）或 `tools/parity_screens.sh main field siege`

## 1. 結果

| 組 | 判定 | 殘差 |
|---|---|---|
| `main` 開局主畫面 | ✅ | **五區全 0 px** |
| `menu-corps`／`menu-city`／`menu-personnel`／`menu-advise` 四張彈出選單 | ✅ | **五區全 0 px** |
| `general-said` 武將一覽 → 曹操台詞 | ✅ | **五區全 0 px** |
| `help7b` 請求協助第二步 | ✅ | **整張 640×400 = 0 px** |
| `form` 編成第一步 | ✅ | `command` 88（原版游標）|
| `finance` 財政視窗 | ✅ | `command` 88、`map` 141（M7 校訂字）|
| `nc-cam` 受控存檔主畫面 | ✅ | `map` 66（原版游標）|
| `picker` 勢力圖例選擇器 | ✅ | `minimap` 82（原版游標）|
| `sortie` 判決畫面 | ✅ | `command` 1369 ＋ `map` 1509（remake 自加的「Enter 繼續」）|
| `video-lcd` 液晶模式 | ✅ | `map` 15,345（remake 多三列系統設定）|
| `quit-menu` 遊戲結束 | ✅ | `map` 14,795（同上）|
| `advise-scene` 說服場景 | ◐ | `map` 3,354 ＝ Enter 提示 2,127 ＋ **進言選單殘影 1,750** |
| `field` 野戰第 52 拍 | ◐ | 七區 0、`field` 95（游標）、**`sb-minimap` 128** |
| `siege` 攻城第 6 拍 | ◐ | 七區 0（含 `bottom`）、**`field` 179**、**`sb-minimap` 432** |

## 2. ⭐ 先前被當成「回歸」的三件事，兩件是取樣點、一件是 fixture

第一輪重量有六組不過。逐一追下來，**只有一個是真缺陷**：

| 症狀 | 看起來像 | 實際是 |
|---|---|---|
| 財政 `map` 141 → **289**，資金 73696 vs 73715 | 經濟規則回歸 | **取樣點只寫到「日」。** `a3dedbb` 把一天改成 23 小時、新的一天從 1 時開始之後，`clock:196/4/20` 停在 1 時，而原版 `until:196/4/20` 停在 2 時。寫成 `clock:196/4/20/2` 就回到 141 |
| 攻城 `bottom` 0 → **930**，六個部隊格的命令圖示全錯 | 戰術規則回歸 | **取樣點落在開場初始化中間。** `-battle-steps` 掃 0/2/4/6/…：只有第 0 拍是 930，第 2 拍起全部 0 |
| 編成 `map` 0 → **49,200**，清單內容全錯 | 武將表回歸 | **fixture 選錯。** `-open-form` 會一路走到第二步的編成面板，停在武將一覽的是 `-open-form-pick` |

三者的共同點：**指標朝「壞掉」的方向動，而原因都在取樣那一側。**
`bisect` 對前兩個各跑了一輪，第一輪還被 [`../spec/197`](../spec/197-shot-mode-draws-no-desktop-cursor.md)
那個桌面游標污染（`236 = 141 + 95`），閾值抬過 95 px 才問得出真正的那一格。

## 3. 真缺陷：事件插圖少了外面那個框

說服場景與判決畫面的 `map` 各差 6,000–8,800 px。原因是原版貼插圖之前
先畫一個 **(48, 128) 304×192** 的底，插圖 (56, 136) 288×176 蓋在中間，
四邊各露 8 px ＝ 畫面上那圈黃邊；remake 的 `drawIventScene` 只貼圖。

⭐ **那一行組語 [`../spec/45`](../spec/45-advise-scene-layout.md) §1.0 早就抄進來了**
（`dx=3 / bx=8 / cx=0C13h / call sub_10C14`），當時只拿它反推插圖座標，
沒注意到同一段還畫了一個底。修法與出處寫成 [`../spec/198`](../spec/198-ivent-scene-frame.md)。

改一支 `drawIventScene` 涵蓋五個呼叫端（說服、外交、判決、事件通知、撥款）：

| | 修之前 | 修之後 |
|---|---:|---:|
| `sortie` `map` | 6,998 | **1,509**（只剩 Enter 提示）|
| `advise-scene` `map` | 8,843 | **3,354** |

## 4. 縮圖那幾格：先錯兩次，最後定位到繪製層

這一格我連續給了**兩個錯的歸因**，記在這裡因為兩次的形狀一樣。

**第一次**：照 [`../spec/133`](../spec/133-opening-deployment.md) §3.6 的
「Y 是亂數、兩邊不同源」寫成「結構性、永遠不會是 0」。
錯在沒查那條斷言的前提——`wlgame` 的 `-battle-exact`
（[`../spec/90`](../spec/90-same-state-parity.md) §2.5）就是為了同步亂數而存在的，
[`114`](114-focus-and-same-battle.md) 還逐槽核對過 96 槽全同，**而它比那條斷言早三天**。

**第二次**：接上 `-battle-exact` ＋ 原版擺位前的 RNG（`0x19C45`）之後，
`sb-minimap` 恆為 624 且不隨 `-shot-frames` 變，於是我寫成「擺位結果本身不同」。
錯在 **`-battle-steps` 的預設是 120**——我掃的是 `-shot-frames`，
而真正推進戰場的那個旗標一直掛在預設值上，兵早就走進陣形了。

### 4.1 把兩個旗標都明寫之後：逐兵全等

```
tools/go.sh run ./cmd/wlgame -direct -scenario 0 -player 0 \
  -save-file workplace/parity/exact2/orig.DAT -load-slot 0 \
  -rng-state workplace/parity/exact2/spawn-rng.bin \
  -battle-exact -open-siege -siege-corps 35,39 -siege-node 82 \
  -shot /tmp/x.png -shot-frames 0 -battle-steps 0 -list-units
```

| 槽 | 原版 | remake | |
|---|---|---|---|
| 側 1 / 隊 0 | (62,26) 體 220、(62,21)、(62,35)、(62,38)、(62,44)… | **側 0** 同值 | ✅ |
| 側 0 / 隊 0 | ( 1,25) 體 188、( 1,41)、( 1,26)、( 1,44)、( 1,22)… | **側 1** 同值 | ✅ |

**96 槽逐槽相同**，側號對調是已知的（remake 的 `Sides[0]` 恆為攻方、
原版的側 0 恆為玩家，[`../spec/133`](../spec/133-opening-deployment.md) §3.5）。

⇒ 擺位沒有問題，亂數也同步了。

### 4.2 第三個錯的歸因：取樣單位

把取樣點也對齊（`-battle-steps 6`）之後 `field` 收到 95 px、`bottom` 0，
只剩 `sb-minimap` 428。我看了放大的並排圖——地形逐像素相同、只有部隊點的
落點不同——於是寫成「縮圖的座標換算」。**那是第三個錯的歸因。**

錯在原版側用的是 `steps:6000000`（**指令**），而
[`tools/dosgolem.sh`](../../tools/dosgolem.sh) 的說明就寫著
「`ticks:N` **以戰術節拍為單位**推進 N 拍 ← 走位期取樣要用這個，不要用 `steps:`」。

改用 `ticks` 重取原版：

| | 原版 `ticks:6` | remake `-battle-steps 6` |
|---|---|---|
| 側 0 / 隊 0 | (3,27)(2,41)(2,25)(3,45)(3,21)(3,29)(2,31)(3,37) | (3,28)(2,40)(1,25)(3,44)(3,21)(3,29)(3,31)(3,36) |

**8 個兵裡 2 個完全相同、6 個差 1 格。** 用 `steps` 取的那一張差到 4–7 格
（原版的 Y 已經收斂到 28–38，remake 還散在 21–44），整個結論因此建立在
一個錯的取樣單位上。

### 4.3 現在的判準：逐兵 ±1 格

| 原版圖 | `field` | `sb-minimap` |
|---|---:|---:|
| `steps:6000000` | 95 | 428 |
| `ticks:6` | 294 | 388 |
| **`ticks:12`** | **317** | **264** |

⚠ `field` 在 `steps` 那一張反而最小（95）——那是**假訊號**：
`field` 區只看得到戰場的一部分，開場的兵多半在可視範圍外，
所以這一區對逐兵差異不敏感。**判準要用 `sb-minimap`**。

`sb-minimap` 的殘差對得上算術：縮圖上每格 2 px，一個兵錯 1 格貢獻 8 px
（4 個消失 ＋ 4 個新增），48 兵 × 8 ＝ 384 ≈ `ticks:6` 量到的 388。

⭐ **差的方向不一致**（+1、−1、−1、−1、+1、−1 混著），所以不是單純的
取樣相位。這是**戰術層逐兵移動的差異**，要用規則層的逐拍比對去追
（像 [`tools/parity_ck.sh`](../../tools/parity_ck.sh) 對六張表做的那樣），
不是畫面對拍能解的——畫面只告訴你「有差」，告訴不了「第幾拍開始差」。

### 4.4 這一格連錯三次，形狀都一樣

| # | 歸因 | 錯在哪 |
|---|---|---|
| 1 | 「亂數不同源，永遠不會是 0」 | 沒查那條斷言的前提；`-battle-exact` 比它早三天 |
| 2 | 「擺位結果本身不同」 | 掃 `-shot-frames` 時 `-battle-steps` 掛在預設 120 |
| 3 | 「縮圖的座標換算」 | 原版側用 `steps`（指令）當取樣單位，該用 `ticks`（節拍）|

三次都是**拿一個沒對齊的比較去解釋差異**，而每一次的數字看起來都夠穩定、
夠像一個結論。教訓 `default-flag-advances-the-fixture` 已經工具化
（`tools/parity_screens.py` 擋戰場組漏寫旗標）；取樣單位那一項還只有規則。

## 5. 攻城到底能不能對拍

**能**，就是上表的 `siege` 那一組（dosgolem 的 `siege:35,82`，
[`72`](72-same-battle-parity.md)）。

不能重跑的是另一組：[`58`](58-parity-retest-20260902.md) §1.2 的
`workplace/promo-live/parity-battle4/`。那個目錄有 `t1`–`t16` 十六張圖
（[`120`](120-screen-parity-retest-20260911.md) §3 先前寫成「只有一張」，那是看漏了），
但它們的取樣點是 `capture-metadata.txt` 裡的 **`wait:6`** ——**牆上秒數**。
即時制底下每次停在不同的遊戲日期（`CLAUDE.md` §3.1），
所以那十六張對不回任何一個遊戲時刻。**不是檔案不見，是取樣點不可重建。**

## 6. 為什麼要有這支閘

規則層的 [`tools/parity_ck.sh`](../../tools/parity_ck.sh) 每一輪都跑得動，
是因為「哪些檢查點、拿什麼比」是資料。畫面層先前沒有這一層：

- 二十幾份 `playtest/` 各記各的，`89`–`104` 整批**全文沒有出現 `parity_shot`**
  ——remake 側只有旗標片段，照抄跑不起來
- 六組的 remake 側漏了 `-save-file`，而原版側跑的是受控存檔
- `58` §1.2 的原版側只寫目錄不寫哪一張圖

⇒ **寫不齊就重跑不出來，而重跑不出來的對拍在下一輪規則改動之後完全不會開口。**
[`../spec/197`](../spec/197-shot-mode-draws-no-desktop-cursor.md) 那個桌面游標
在畫面上待了三天沒人發現，就是這個原因。

## 7. 未解

| 項目 | 現況 |
|---|---|
| 進言選單在說服場景上沒關掉 | remake 殘影 1,750 px；原版進場時清掉。`advise-scene` 的 `gap` |
| 「Enter 繼續」提示要不要登記成 remake 差異 | 目前只在 `cmd/wlgame/advise.go:678` 與 [`104`](104-advise-verdict-parity.md) §1，沒有 spec |
| 插圖框的填色 | 中間被插圖全蓋住，這個畫面量不到（[`../spec/198`](../spec/198-ivent-scene-frame.md) §5）|
| 閘還沒收進去的組 | `38`／`39`／`42`／`83`／`84`／`92`／`94`／`96`／`98`／`99`／`105`／`106` 等；`37` 的孫策第二樣本與 `40` 的攻城第二輪原版側不是同狀態 |
