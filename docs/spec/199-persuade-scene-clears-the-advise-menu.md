# 199 — 說服場景進場要擦掉進言選單

**狀態：CONFORMED。** 進言那五項選單畫在 (0, 64)，而說服場景進場時
`sub_189A4(cx=151Bh)` 蓋掉 **(0, 32, 432, 336)**——選單整個在那個矩形裡，
所以原版的說服畫面上看不到它。remake 把 `adviseMenuStale` 留到
「地圖重畫」才清，於是場景左上角掛著一塊 1,750 px 的殘影。

- 日期：2026-09-11
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_13830` 的進場、`sub_16224` 的 `dx = 400h`
- 推論等級：confirmed（兩個矩形都來自立即值，一個包住另一個）
- 相關：[`140`](140-status-message-box.md) §2.0（同一次進場也清狀態列）、
  [`126`](126-command-popup-menus.md) §1.2（選完不擦的那一條規則）、
  [`45`](45-advise-scene-layout.md) §1

## 1. 原版做什麼

| 東西 | 矩形 | 出處 |
|---|---|---|
| 進言選單 | (0, 64)，大小由內容算 | `sub_16224` 的 `dx = 400h` ⇒ 粗格 (0, 4) |
| 說服場景蓋掉的範圍 | **(0, 32, 432, 336)** | `sub_189A4(cx = 151Bh)`：`cl=1Bh` 27×16 ＝ 432、`ch=15h` 21×16 ＝ 336 |

⇒ 選單的左上角 (0, 64) 在場景矩形內，而場景高度到 y = 368，
**整個選單被蓋掉**。這不是另一道清除指令，是同一次重畫的副作用。

⭐ [`140`](140-status-message-box.md) §2.0 已經記過這個矩形，
當時用它解釋**狀態列為什麼要另外清**（狀態列在 (0, 320, 256, 80)，
下緣 32 px 蓋不到）。同一句話反過來讀就是：**蓋得到的東西不必另外清**，
而進言選單正是蓋得到的那一種。

## 2. remake 錯在哪

`adviseMenuStale` 只在 `moveCamTo` 清（`cmd/wlgame/main.go:1361`），
而說服場景不動鏡頭。`drawAdviseMenu` 因此在 `advisePersuade` 底下照畫。

⚠ **不能改成「選完就不畫」**——那會踩掉 [`126`](126-command-popup-menus.md) §1.2：
原版把目標一覽表直接畫在選單上面，露出來的那一列一直掛在清單上緣。
選單要留到**說服場景進場**才消失，不是選完就消失。

## 3. remake 實作

`beginPersuasion` 在清狀態列的同一處把 `adviseMenuStale` 也清掉
（`cmd/wlgame/advise.go`）。兩件事出自同一次進場，放在一起。

## 4. 驗證

`tools/parity_screens.sh advise-scene` 的 `map` 區：**3,354 → 1,604 px**，
消掉的 1,750 px 正是選單那一塊
（[`../playtest/121`](../playtest/121-screen-parity-gate.md) §1）。
剩下的 1,604 是 remake 自加的「Enter 繼續」提示，已登記在
`REMAKE-PLAN.md` 的 Intentional differences。

## 5. 未解

- 判決畫面（遷都／請求君主出陣，[`49`](49-advise-relocate-and-sortie.md)）走的不是
  `sub_13830`，它的進場蓋不蓋這個矩形還沒對過。`sortie` 那一組的
  「選單殘影」rect 現在量到 0 px，但那是 `-advise-sortie` 這個 fixture
  沒有經過選單造成的，**不算證據**。
