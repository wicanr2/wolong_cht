# 91 — 武將自陳那一張：五區全 0 px，事件列是唯一的差異

**狀態：通過。** 指令列「武將」那格選走一位之後（自陳台詞 ＋ 清單留著
＋ 那一列反白 ＋ 狀態列 #24），逐區比對**五個分區全部 0 px**——
這是第一張連滑鼠游標都 0 的完整畫面。唯一的缺陷是 remake 在
事件列多寫了一句「選擇了 曹操」，4,015 px 全在那條列上。

- 日期：2026-09-06
- 規格：[`../spec/145`](../spec/145-general-and-faction-cells.md) §1／§1.1
- 原版側：`WOLONG_DOSGOLEM_GAMEDIR=dosgolem/root-noclouds tools/dosgolem.sh
  workplace/parity/cells "wait;click:320,200;click:300,151,1000;runto:11CD0;
  sclick:352,15;stap:336,47;steps:400000;stap:200,112;steps:400000;
  stap:200,112;steps:900000;shot:orig-gen1;clock"` → 196年4月17日 6時
- remake 側：`-save-file workplace/dosgolem/root-noclouds/SAVE.DAT -load-slot 0
  -open-list -list-pick-row 0 -cam 0,0
  -fixture-when clock:196/4/17/6 -shot-when clock:196/4/17/6`
  ⚠ **受控存檔不能省**（[`../spec/147`](../spec/147-controlled-parity-save.md)）：原版側跑的是 `root-noclouds`，remake 只用 `-direct` 從劇本跑到同一時刻會有雲、也會有軌跡分歧。這一行本輪補上，補之前照著跑對不出文中的數字
  （`-list-pick-row` 本輪新增）

## 1. 結果

| 區 | 修之前 | 修之後 |
|---|---:|---:|
| `banner` | 0 | **0** |
| `command` | 0 | **0** |
| `map` | 4,015 | **0** |
| `minimap` | 0 | **0** |
| `faction` | 0 | **0** |

原版那一張畫的是：清單第 0 列「曹操　9 13 13　曹操　**君王**」整列反白
（底 `519241`、字 `f3e300`），訊息框裡曹操的肖像配
**「哼！城塞在我的武力之前就如同紙片！」**，左下角狀態列
「確認指示之武將的能力。」＝ #24。三件事逐像素相同。

⭐ 台詞驗算：曹操的攻城／野戰／水戰適性最高值落在**攻城**（`+0x0E`）
⇒ 組 `0x1A9`，配他的變體 `+0x1E` ⇒ 城塞戰那八句其中一句。
與 [`86`](86-general-faction-cells-parity.md) §2 的夏侯淵（野戰組 `0x1AA`）
一起，四組裡的兩組各有一個原版錨點。

## 2. 抓到的缺陷：事件列寫了「選擇了 ○○○」

那一句是台詞還沒接上以前的暫時產物，接上之後變成同一件事講兩遍。
拿掉之後 `map` 從 4,015 掉到 0。

**判準寫進了 [`../spec/145`](../spec/145-general-and-faction-cells.md) §1.1**：
事件列登記的是「世界狀態變了」，不是「玩家點了什麼」。純瀏覽的流程不寫。

## 3. ⭐ 這一張原版**沒畫**滑鼠游標

前面每一張帶原版擷取的畫面都留著 66–95 px 的游標殘差
（[`82`](82-strategy-date-aligned-parity.md)、[`90`](90-formation-second-step.md)），
而這一張是 0——放大 `(196,108)` 那一塊，原版就是清單的反白列，沒有紅箭頭。

**成因已解**：畫不畫由 `seg002` 的一個顯示旗標 `byte_20100` 決定
（[`../re/88`](../re/88-mouse-cursor-visibility.md)），不是由「畫面上有什麼」決定。
這一張的最終狀態是 `byte_20100 = 0`，[`92`](92-personnel-assign-parity.md)
那一族是 1——兩組都用 `ipeek:20100:1` 讀出來，與截圖逐格對上。

⭐ **所以「有沒有游標」是可觀測的，不必再從畫面推。**
對拍時多送一個 `ipeek:20100:1`，殘差就分得出
「remake 少畫了游標」與「原版根本沒畫」。

⚠ remake 這一側還沒同步：54 個 hide／show 呼叫點要逐一對應到 remake
的繪圖流程（[`../re/88`](../re/88-mouse-cursor-visibility.md) §5）。
在那之前 remake 只在地圖選點時自繪。

## 4. 未解

| 項目 | 現況 |
|---|---|
| remake 要把游標的 hide／show 接在哪 | 原版的規則已解（[`../re/88`](../re/88-mouse-cursor-visibility.md)），但 54 個呼叫點還沒對應到 remake 的繪圖流程 |
| `0x1A8`（俘虜）與 `0x1AB`（海戰）兩組 | 沒有原版擷取；需要一個有俘虜、或水戰適性最高的武將的局面 |
