# 44 — 新遊戲的信賴度初始值：實機 0xFF 滿格

**狀態：confirmed。** 新遊戲（劇本 1）進主畫面後開自勢力情報視窗，
信賴度量條**滿格**（實測 x 456–615 正好 160 px、高 2 px）。
以 `sub_15F27` 的長度式反推，`byte_10D00 ＝ 0xFF`。

- 日期：2026-08-25
- 對象：松崗 DOS/V（`wolong-dosboxx`，唯讀掛載）
- 擷取物：`workplace/promo-live/trust-init7/p0-info.png`
- timeline：`wait:130;click:320,215;wait:3;click:300,190;wait:4;click:450,154;wait:2;click:352,336;wait:3;click:352,336;wait:6;click:0,0;wait:1;click:40,0;press;wait:2;shot:p0-info`
  （NEW GAME 四步 → 歸零 → 自勢力情報圖示；主畫面歸零後主機 x ＝ 遊戲 x ÷ 9.6，
  自勢力情報熱區 (368,0,32,32) 中心 384 ⇒ 40）

## 1. 量到什麼

| 項目 | 值 |
|---|---|
| 量條紅色像素 | x 456–615（**160 px 滿長**）、y ＝ 遊戲 292、高 2 |
| 長度式（`sub_15F27`） | `(信賴度×100 + 0x9F) ÷ 0xA0`，滿長 160 |
| 反推 | 160 ⇔ 信賴度 ≥ 0xFE；配合區塊 `+0x10 = 0xFF`，取 **0xFF** |

## 2. `sub_18B12` 的時序（靜態，與實測相符）

新遊戲選完劇本後（`sub_11AC3` 的 `loc_11AE0`）：

1. `sub_18CAE`（ah=0）把劇本區塊的 59-byte 全域段載入——
   `byte_10D00 ← 區塊 +0x10（0xFF）`、`word_10CFD ← 區塊 +0x0D（0xFFFF，未選勢力）`。
2. 緊接著 `mov bx, cs:word_10CFD / mov al, [bx+2Bh] / mov cs:byte_10D00, al`。
   `bx = 0xFFFF` 時有效位址 16-bit 迴繞成 `ds:0x2A`；執行期 `ds`
   （`word_10D52`）的 offset 0 ＝ 檔內 `+0x80`（勢力表基底），
   所以讀到的是**勢力 0 記錄的 `+0x2A`（駐外交官欄，四劇本開局全 0xFF）**。
3. 之後選君主（`loc_11AF8` 寫 `word_10CFD`）**不再重寫** `byte_10D00`。

兩條路徑殊途同歸 `0xFF`：這行 `[bx+2Bh]` 是**實際無害的越界讀**
（每個劇本的勢力 `+0x2B` 檔內全 0，若 `word_10CFD` 已指向真勢力，
它會把信賴度歸零——但時序上它只會在 `word_10CFD=0xFFFF` 時跑）。
remake 以區塊 `+0x10` 初始化 `World.Trust`，行為一致，不移植這行。

## 3. 未解

<!-- 缺口：無 -->
