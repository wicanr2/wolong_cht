# 88 — 受控存檔對拍：跑滿一個遊戲日，四區 0 px、`map` 只剩游標框

**狀態：通過。** 用 `tools/parity_save.py --no-clouds` 做一份把雲關掉的
受控存檔，兩邊都載它。跑滿一個遊戲日之後 `banner`／`command`／`minimap`／
`faction` 全 0 px，`map` **66 px**——而那 66 px 是原版自己畫的地圖游標框，
剛好落在鏡頭跳過去的那一格。

- 日期：2026-09-06
- 規格：[`../spec/147`](../spec/147-controlled-parity-save.md)
- 原版側：`WOLONG_DOSGOLEM_GAMEDIR=dosgolem/root-noclouds tools/dosgolem.sh
  workplace/parity/status "wait;click:320,200;click:300,151;until:196/4/20;
  sclick:352,15;stap:384,47;steps:400000;stap:200,128;steps:400000;
  stap:200,128;steps:1600000;rpress;steps:1200000;shot:orig-nc-cam;clock"`
- remake 側：`tools/parity_shot.sh out.png -direct -scenario 0 -player 0 -seed 7
  -save-file workplace/dosgolem/root-noclouds/SAVE.DAT -load-slot 0
  -cam 271,153 -open-window 0 -shot-when clock:196/4/21/11`

## 1. 結果

| 存檔 | 取樣點 | `banner` | `command` | `map` | `minimap` | `faction` |
|---|---|---:|---:|---:|---:|---:|
| 原本的（remake 沒畫雲）| 196/4/21 8 時 | 0 | 0 | 11,420 | 0 | 0 |
| 原本的（remake 畫了雲，各自漂）| 同上 | 0 | 0 | 22,210 | 0 | 0 |
| **`--no-clouds`** | 196/4/21 11 時 | **0** | **0** | **66** | **0** | **0** |

⭐ **中間那一列是這件事的重點。** 把雲接上去之後，同時鐘的殘差**變大**——
因為畫面上同時有原版的雲與 remake 的雲。**接對了反而更差**，
是即時制對拍最容易被誤讀成回歸的一種形狀。受控存檔把它整個繞開。

## 2. 那 66 px

差異落在**一格**：螢幕 (320–335, 224–239) ＝ 地圖格 (291,165) ＝ 濡須口本身。
原版在那一格的邊上多了一圈白（`f3f3f3`），內部的據點徽記兩邊逐點相同。

那是**原版畫在地圖上的游標框**：白色空心 16×16，貼著游標所在的格。
鏡頭由 `sub_12151` 跳到目標據點之後，游標就落在那一格上。
remake 沒有這個框（[`../spec/147`](../spec/147-controlled-parity-save.md) §5）。

⚠ **這與先前那個「88 px」是同一個東西的另一種樣子**：在選單／指令列上
原版畫的是 14×14 的箭頭，在大地圖上畫的是 16×16 的空心框。
五個分區把 640×400 鋪滿，所以游標停在哪都會落進某一區。

## 3. 驗證那個開關真的有效

載入受控存檔跑到 196/4/21，`peek:2754:2140:48`：

```
00 00 F3 FF 0E 01 00 00 F7 FC F2 00 0C 10 00 00
└ 旗標 0    └ 座標與計時器與檔案裡一模一樣
```

**旗標清掉了，座標與計時器一動也沒動**——`sub_12459` 真的整筆跳過，
不是「畫不出來但還在動」。畫面上也確認一朵雲都沒有。

## 4. 未解

| 項目 | 現況 |
|---|---|
| 原版的地圖游標框 | remake 沒畫；要接得先讀出顏色、線寬與更新時機 |
| 還能關掉什麼 | 目前只有雲。天災、AI 出兵、募兵都吃亂數，各自需要自己的欄位 |
