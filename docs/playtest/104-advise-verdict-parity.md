# 104 — 判決畫面（請求君主出陣）：拿掉事件列之後只剩 remake 的提示

**狀態：通過。** 進言第 4 項「請求君主出陣」的判決畫面
（`sub_13B08` 的三句 ＋ 插圖）第一次逐像素比：
**插圖、上框、選單殘影、下半畫面各 0 px**，殘差只剩 remake 自己加的
「Enter 繼續」。[`35`](35-advise-verdict-screens.md) 先前只做過目視驗收。

- 日期：2026-09-07
- 規格：[`../spec/45`](../spec/45-advise-scene-layout.md) §1.0、
  [`../spec/145`](../spec/145-general-and-faction-cells.md) §1.1（判決台詞就是回報，不寫事件列）
- 原版側：`WOLONG_DOSGOLEM_GAMEDIR=dosgolem/root-noclouds tools/dosgolem.sh …
  "…;sclick:352,15;stap:48,47;steps:300000;move:48,200;steps:80000;
  move:48,232;steps:80000;move:48,264;steps:80000;move:48,296;steps:80000;
  press;steps:1200000;shot:orig-sortie"`
  （**四次 `move` ＝ 第 4 列**，[`../spec/131`](../spec/131-dosgolem-oracle.md) §3.5）
- remake 側：`-save-file workplace/dosgolem/root-noclouds/SAVE.DAT -load-slot 0
  -advise-sortie -cam 0,0
  -fixture-when clock:196/4/17/6 -shot-when clock:196/4/17/6`

## 1. 結果

| 區 | 第一次拍 | 修完 |
|---|---:|---:|
| 插圖 (56,136,288×176) | 0 | **0** |
| 上框 (0,80,256×80) | 0 | **0** |
| 選單殘影 (0,64,120×32) | 3,840（沒畫）| **0** |
| 下半畫面 (0,320,432×80) | 6,000＋（事件列）| **0** |
| `command` 區 | 13,489（沒開命令視窗）| 1,369（＝提示）|

## 2. 兩個缺陷

**一、判決三句之外又寫了一次事件列。** `beginSortie` 寫
「君主親自出陣」／「請求出陣：君主不同意」，而 `sub_13B08` 的三句
已經是回報了——[`../spec/145`](../spec/145-general-and-faction-cells.md) §1.1
的判準（**玩家自己剛下的指令不寫事件列**）套過來，
進言的**八個出口**都拿掉了（遷都成／不成、出陣成／不成、說服失敗、
進言撤回、進言失效、進言成立）。

**二、fixture 沒走真實流程。** `-advise-sortie` 直接 `openAdvise()` ＋
`beginSortie()`，於是命令視窗沒開、「進言」那一格沒反白、選單的框也沒留下。
改成 `pickAdviseCommand(adviseSortieRow)` 走原本那條路之後三件事一起對上。
⭐ 這是 [`../playtest/102`](../playtest/102-help-second-step.md) §2.2 那條
「fixture 造得出來 ≠ 原版走得到」的另一面：**造得出來也可能少了東西**。

## 3. ⭐ 選單上的 `move:` 是「往下一格」

第一次試拍時 `move:48,264` 選到的是第 2 列而不是第 3 列，
`move:48,296` 選到第 1 列——**Y 的絕對值不決定停在哪一列，移動事件的
次數才是**（[`../spec/131`](../spec/131-dosgolem-oracle.md) §3.5）。
既有腳本寫成 `[move:48,200;]*N` 的 `N` 就是列號。

## 4. 未解

| 項目 | 現況 |
|---|---|
| 「君主出征中」那一句 | `openAdvise` 仍寫事件列，而**原版是 TALK #64 訊息框**（[`../re/22`](../re/22-strategy-command-tree.md) §3.4）——不是「該不該寫事件列」，是呈現方式用錯了。還沒拍 |
| 「Enter 繼續」提示 | remake 自己加的操作說明，保留；它讓每一張說服／判決場景差 1,369–4,628 px |
