# 153 — 「遊戲結束」的確認是兩項選單，位置寫死在 (352, 272)

**狀態：CONFORMED。** 系統選單第 6 列跳的是「　終　　了　」／「　取　　消　」
兩項選單，不是 ＹＥＳ／ＮＯ 對話框。⭐ **它的框是六個 handler 裡唯一
把座標寫死在指令裡的**——行軍三選一與大地圖兩項都跟著游標走。

- 日期：2026-09-06
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）
  - 系統選單分派表 `cs:6056h` 第 6 筆 ＝ `60B4`（[`../re/31`](../re/31-faction-picker-screen.md) §表、
    [`../re/55`](../re/55-system-menu-window.md) §5）
  - handler `0x160B4`–`0x160CB`（IDA 沒建成函式，落在 `sub_160A5` 之後）
  - TALK #81：`['　終　　了　', '　取　　消　', '']`
- 推論等級：**confirmed**（機器碼逐條 ＋ 原版擷取 `orig-sys5`）
- 驗證：[`../playtest/101`](../playtest/101-quit-menu-parity.md)
- 相關：[`13`](13-main-window-toggles.md)（系統選單六列）、[`26`](26-yes-no-dialog.md)（F10 那條，remake 差異）、
  [`30`](30-victory.md) §1（離開碼 0 的觸發點就是這一段的 `0x160C8`）、
  [`151`](151-map-click.md)（同一支選單引擎 `sub_193E9`）

## 1. 原版做什麼

```asm
                mov     ax, 102h        ; ah = 1（模式）、al = 2（兩項）
                mov     cx, 51h ; 'Q'   ; ★ TALK #81
                mov     dx, 1116h       ; ★ dl = 16h → X = 352、dh = 11h → Y = 272
                call    sub_193E9       ; 選單引擎（與 sub_11F0E／sub_1804E 同一支）
                jb      short locret_160CB      ; CF ＝ 右鍵取消
                and     al, al
                jnz     short locret_160CB      ; 選第 1 列「取消」也離開
                xor     al, al
                call    sub_11CB1               ; 第 0 列「終了」→ 離開碼 0
locret_160CB:   retn
```

三條出路裡**兩條什麼都不做**：右鍵、以及選中第 1 列。
只有第 0 列會走到 `sub_11CB1`，離開碼 `al = 0` ＝ 回 `YNVSHELL.COM`
（[`30`](30-victory.md) §1 那張表的第一列）。**沒有第二層確認、沒有自動存檔。**

### 1.1 ⭐ 位置是寫死的，不跟著游標

| 選單 | 位置怎麼來 | 夾制 |
|---|---|---|
| 行軍三選一 `sub_1804E` | 游標所在格 | X ≤ `21h`、Y ≤ `18h −` 項數 |
| 大地圖兩項 `sub_11F0E` | 游標所在格 | X ≤ `23h`、Y ≤ `16h` |
| **這一張 `0x160B4`** | **`dx = 1116h` 常數** | 不需要——常數本身就在畫面內 |

`dl = 0x16`、`dh = 0x11`，格 × 16 ⇒ **(352, 272)**。
系統視窗自己是 `(208, 112, 208, 192)`（[`12`](12-strategy-chrome.md)），
第 6 列的值格中心大約在 (376, 280)——**框正好貼在被點的那一格右下**，
看起來像跟著游標，其實是常數剛好落在那裡。

⚠ 所以**不能拿這一張去驗「選單跟著游標」的規則**，
也不能反過來用它推 `sub_193E9` 對超出畫面的處理
（[`151`](151-map-click.md) §未解那條仍然沒解）。

### 1.2 框的尺寸

字串是六個全形字一列、兩列，與大地圖那兩項同一個形狀
⇒ 框寬 112、高 48（[`125`](125-menu-box-width-from-padding.md) 的算式）。
右緣 352 + 112 ＝ 464，下緣 272 + 48 ＝ 320，都在畫面內。

## 2. remake 要怎麼改

| 項目 | 作法 |
|---|---|
| 狀態 | `game.quitMenu`（`cmd/wlgame/quitmenu.go`）——**零值 ＝ 沒開** |
| 開啟 | `dispatchSystemRow` 的 `sysRowQuit` 改呼叫 `g.openQuitMenu()` |
| 字串 | `talkmenu.MenuLabels(lib.Talk, 0x51, …)`——⛔ 不能走 `talkLines`，它會 `TrimRight` 掉行尾全形空白而**框寬由第一列的字數決定**（[`125`](125-menu-box-width-from-padding.md)）|
| 繪製 | `drawLegacyChoiceBox(screen, 352, 272, rows, sel)`，與大地圖兩項共用 |
| 點擊 | `talkChoiceClick(352, 272, rows)`；右鍵／ESC 取消 |
| 離開 | 第 0 列 → `ebiten.Termination` |

`updateQuitMenu` 排在 `Update` 最前面（截圖與錄影的收尾之後），
理由是**它是模態的**：原版在 `sub_193E9` 的迴圈裡，時鐘與 AI 都不動。

### 2.1 remake 差異：F10 那條保留 ＹＥＳ／ＮＯ

`CLAUDE.md` §9 明訂「ESC 只取消，F10 才離開，離開前跳 Ｙ／Ｎ 確認」——
那是 remake 自己加的快捷鍵，原版沒有 F10。
**兩條路各自保留**：F10 走 [`26`](26-yes-no-dialog.md) 的置中對話框，
系統選單第 6 列走這一份的兩項選單。

## 3. 驗證

| 方式 | 內容 |
|---|---|
| 對原版 ✅ | [`../playtest/101`](../playtest/101-quit-menu-parity.md) |
| 單元測試 | `TestQuitMenuOpensFromSystemRow`：點第 6 列會開，且不動 `quitting` |
| 單元測試 | `TestQuitMenuRowSemantics`：第 0 列回 true、第 1 列與取消回 false |
| 單元測試 | `TestQuitMenuAnchorIsFixed`：座標就是 `(0x16*16, 0x11*16)` |

## 4. 未解

| 項目 | 現況 |
|---|---|
| 離開時自動存檔 | `CLAUDE.md` §9 要求「離開前自動存檔，存檔失敗就不離開」，remake 兩條路目前都只是 `ebiten.Termination`。**原版這一段沒有存檔**（`xor al,al` 直接走），所以那是 remake 差異，還沒實作 |
