# 文件索引（自動產生，不要手改）

由 `tools/index.py generate` 從 `docs/**/*.md` 生出來。
**狀態與日期是從各文件的內文讀的，不是另外維護的一份摘要**——
手寫的索引會過時，而過時的索引比沒有索引更糟（`docs/playtest/04` 的「四件事被擋住」整張表在寫下的當下就是錯的）。

提交前跑 `tools/index.py check`，它會擋下：狀態行與內文矛盾、
同一斷言在兩份文件等級不同、連結壞掉、`CONTEXT.md` 漏登記。

## 文件

| 文件 | 標題 | 狀態 | 日期 |
|---|---|---|---|
| [`docs/formats/01-talk-dat.md`](formats/01-talk-dat.md) | 01 — TALK.DAT 訊息表格式 | READY | 2026-08-07 |
| [`docs/formats/02-brg-palette.md`](formats/02-brg-palette.md) | 02 — .BRG 調色盤格式 | READY | 2026-08-07 |
| [`docs/formats/03-grf-images.md`](formats/03-grf-images.md) | 03 — GRF.DAT 圖庫格式 | KAOGRF／KYOGRF／IVENTGRF READY，ICONGRF 部分解。 | 2026-08-07 |
| [`docs/formats/04-map-sch-container.md`](formats/04-map-sch-container.md) | 04 — .MAP／.SCH／.MCH：兩種完全不同的東西 | 容器格式的索引層 READY，壓縮演算法未解。 | 2026-08-07 |
| [`docs/formats/05-mmap-worldmap.md`](formats/05-mmap-worldmap.md) | 05 — MMAP. 大地圖 | READY。MMAP.MDL、地圖尺寸、自動連接都 confirmed； MMAP.MAP 的 RLE 見 [doc… | 2026-08-07 |
| [`docs/formats/06-mmap-rle.md`](formats/06-mmap-rle.md) | 06 — MMAP.MAP 的 RLE 壓縮 | READY。 | 2026-08-07 |
| [`docs/formats/07-battle.md`](formats/07-battle.md) | 07 — BATTLE. 戰場資料 | 分段結構、圖塊定義與像素格式都 confirmed。 | 2026-08-07 |
| [`docs/formats/08-sinario-save.md`](formats/08-sinario-save.md) | 08 — SINARIO.DAT / SAVE.DAT：劇本與存檔 | 整體結構 confirmed，武將能力值 confirmed，其餘欄位進行中。 | 2026-08-07 |
| [`docs/mechanics/00-index.md`](mechanics/00-index.md) | 00 — 遊戲機制索引 | 索引與推論等級定義，長期有效。 | 2026-08-08 |
| [`docs/mechanics/10-strategy.md`](mechanics/10-strategy.md) | 10 — 大地圖政略 | 指令清單完整（說明書），數值判定未解。 | 2026-08-08 |
| [`docs/mechanics/15-realtime.md`](mechanics/15-realtime.md) | 15 — 即時制的時間模型 | ✅ READY。整條時間鏈已在機器碼裡讀出來。 | 2026-08-08 |
| [`docs/mechanics/20-military.md`](mechanics/20-military.md) | 20 — 行軍與軍團 | 道路網與行軍已解並實作；軍團編成的數值判定未解。 | 2026-08-08 |
| [`docs/mechanics/30-combat.md`](mechanics/30-combat.md) | 30 — 戰場（戰術） | 進場規則與戰略層的自動判定全解；戰術層已解並實作 | 2026-08-08 |
| [`docs/mechanics/40-economy.md`](mechanics/40-economy.md) | 40 — 經濟：資金與預備兵 | 機制與公式都已解 | 2026-08-08 |
| [`docs/mechanics/50-diplomacy.md`](mechanics/50-diplomacy.md) | 50 — 外交 | 成立條件與外交官的數值都已解 | 2026-08-08 |
| [`docs/mechanics/60-personnel.md`](mechanics/60-personnel.md) | 60 — 武將 | 三個能力值的作用已定案（說明書＋機器碼），數值公式部分已知。 | 2026-08-08 |
| [`docs/mechanics/70-ai.md`](mechanics/70-ai.md) | 70 — 電腦 AI 的判斷邏輯 | 侵攻的起意與取消已定案（機器碼）；決策鏈的細節蒐集中。 | 2026-08-08 |
| [`docs/mechanics/80-victory.md`](mechanics/80-victory.md) | 80 — 勝負判定 | 勝負條件全部定案（說明書 1.1、1.2 節）。 | 2026-08-08 |
| [`docs/playtest/01-dosbox-dosv.md`](playtest/01-dosbox-dosv.md) | 01 — 松崗 DOS/V 版首次實跑（DOSBox） | DOS/V 版可開機；字型懸案結案， | 2026-08-07 |
| [`docs/playtest/02-dosboxx-pc98.md`](playtest/02-dosboxx-pc98.md) | 02 — PC-98 日文原版實跑（DOSBox-X）：oracle 建立完成 | PC-98 oracle 已建立，沒有防拷。 | 2026-08-07 |
| [`docs/playtest/03-verification-log.md`](playtest/03-verification-log.md) | 03 — 補驗紀錄：哪些結論撐得住，哪些撐不住 | 持續累加。專記驗證與失敗的驗證嘗試。 | 2026-08-07 |
| [`docs/playtest/04-mouse-automation-blocked.md`](playtest/04-mouse-automation-blocked.md) | 04 — 受阻：PC-98 oracle 動不了遊戲內的滑鼠游標 | ⛔ 已作廢——本文的兩個核心判斷都被推翻，解法見 docs/playtest/06。 | 2026-08-07 |
| [`docs/playtest/05-viewport-and-city-coords.md`](playtest/05-viewport-and-city-coords.md) | 05 — 一張靜止的截圖就把據點座標定案了 | SINARIO.DAT 的據點 X／Y 座標升到 confirmed。 | 2026-08-08 |
| [`docs/playtest/06-mouse-solved.md`](playtest/06-mouse-solved.md) | 06 — PC-98 oracle 的滑鼠通了 | 可以自動操作遊戲。 | 2026-08-08 |
| [`docs/playtest/07-in-game-oracle.md`](playtest/07-in-game-oracle.md) | 07 — 進到遊戲本體：三個欄位靠實機畫面定案 | 可自動玩進遊戲本體；首都與頭像編號兩欄位升 confirmed。 | 2026-08-08 |
| [`docs/re/01-first-recon.md`](re/01-first-recon.md) | 01 — 首輪偵查：兩版檔案清單與比對 | READY。兩版檔案清單、執行結構、逐檔比對都完成。 | 2026-08-07 |
| [`docs/re/02-palette-routine.md`](re/02-palette-routine.md) | 02 — 調色盤常式：.BRG 的通道順序與亮度縮放 | READY。兩版調色盤常式互證，.BRG 格式定案。 | 2026-08-07 |
| [`docs/re/03-image-blitter.md`](re/03-image-blitter.md) | 03 — 圖庫載入器與 VGA 繪製常式 | READY。圖庫載入器與四平面繪製常式都讀完。 | 2026-08-07 |
| [`docs/re/04-mmap-entry-points.md`](re/04-mmap-entry-points.md) | 04 — 大地圖：入口點與記憶體佈局 | 入口全部讀完，尺寸與自動連接已定案。 | 2026-08-07 |
| [`docs/re/05-battle-selection.md`](re/05-battle-selection.md) | 05 — 政略與戰術的接縫：戰場怎麼被選出來 | 政略↔戰術的接縫已解；戰場記錄的部分欄位未解。 | 2026-08-07 |
| [`docs/re/06-game-clock.md`](re/06-game-clock.md) | 06 — 遊戲時鐘：sub_11D8E | 整條時間鏈已解，confirmed。 | 2026-08-08 |
| [`docs/re/07-monthly-settlement.md`](re/07-monthly-settlement.md) | 07 — 月結：sub_15358 與它的九支子程式 | 經濟公式全部讀出來，confirmed。 | 2026-08-08 |
| [`docs/re/08-hourly-update.md`](re/08-hourly-update.md) | 08 — 每「時」的世界更新：sub_13E11 | 主結構已解。軍團編成、外交官效果、事件分派都在這裡。 | 2026-08-08 |
| [`docs/re/09-combat.md`](re/09-combat.md) | 09 — 戰鬥：觸發、自動判定、傷亡與武將的下場 | 戰略層的戰鬥全解。戰術畫面（sub_19FA0）還沒碰。 | 2026-08-08 |
| [`docs/re/10-rng.md`](re/10-rng.md) | 10 — 亂數產生器：sub_1ECE0 與 sub_1EC82 | 全解，已實作成 internal/rules/rng。 | 2026-08-08 |
| [`docs/re/11-tactical-battle.md`](re/11-tactical-battle.md) | 11 — 戰術戰鬥：模組結構與戰場資料模型 | 模組骨架與戰場資料模型已解。戰鬥規則本身還沒解。 | 2026-08-08 |
| [`docs/reference/01-jp-manual.md`](reference/01-jp-manual.md) | 01 — 日文原版說明書判讀紀錄 | 有實質機制的頁都讀完了，剩 p.6 啟動操作與 p.36–38 附錄。 | 2026-08-08 |
| [`docs/reference/02-jp-cht-diff.md`](reference/02-jp-cht-diff.md) | 02 — 日中對照：TALK.DAT 第一批發現 | 全量 1,022 則的 | 2026-08-07 |
| [`docs/reference/03-baked-japanese.md`](reference/03-baked-japanese.md) | 03 — 燒進美術裡的日文：松崗版沒重繪的部分 | 已確認的缺口：標題橫幅「臥竜伝」兩版相同（松崗沒重繪）。 | 2026-08-07 |

## 斷言（欄位／常數 → 推論等級 → 出處）

共 136 條。**要查「這件事解了沒」先看這裡**，
不要重讀整份文件，更不要重推一次。

### confirmed（65 條）

| 鍵 | 出處 |
|---|---|
| 1. ✅ 成功：地形類型的空間關係檢查 ▸ 類型 1 鄰接水 | `docs/playtest/03-verification-log.md` |
| 1. ✅ 成功：地形類型的空間關係檢查 ▸ 類型 2 鄰接水 | `docs/playtest/03-verification-log.md` |
| 1. ✅ 成功：地形類型的空間關係檢查 ▸ 類型 8 鄰接水 | `docs/playtest/03-verification-log.md` |
| 1. 戰場怎麼被選出來 ▸ 攻城戰 | `docs/mechanics/30-combat.md` |
| 1. 戰場怎麼被選出來 ▸ 野戰 | `docs/mechanics/30-combat.md` |
| 1.4 ⭐⭐ 兩個可直接寫成程式的定義 ▸ 國力 | `docs/mechanics/70-ai.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x00 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x01 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x02 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x14 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x18 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x1A | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x1C | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x1D | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x20 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x23 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x3E | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x00 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x01 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x02 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x08 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x0A | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x0C | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x0E | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x10 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x13 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x16 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x19 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x1A | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x1C–+0x1F | `docs/formats/08-sinario-save.md` |
| 2. 地形類型對映表（confirmed） ▸ 1 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 2 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 6 | `docs/mechanics/30-combat.md` |
| 3. 武將記錄（32 byte） ▸ +1 | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +14（0x0E） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +15（0x0F） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +17 | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +18 | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +19 | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +2 | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +23（0x17） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +24（0x18） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +28（0x1C） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +29（0x1D） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +31（0x1F） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +8 | `docs/formats/08-sinario-save.md` |
| 3. 軍團記錄的欄位（部分） ▸ +0x1A | `docs/re/05-battle-selection.md` |
| 3. 軍團記錄的欄位（部分） ▸ +0x1C | `docs/re/05-battle-selection.md` |
| 4. 據點記錄（32 byte） ▸ +2 | `docs/formats/08-sinario-save.md` |
| 據點記錄再解出四個欄位 ▸ +19h | `docs/re/07-monthly-settlement.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x00 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x01 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x02 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x06 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x0A | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x0B | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x0E | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x10 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x12 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x14 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x16 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x18 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x1A | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x1C | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x1E | `docs/re/08-hourly-update.md` |

### 強證據（30 條）

| 鍵 | 出處 |
|---|---|
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x03 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x04 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x06 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x08 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x11 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x12 | `docs/formats/08-sinario-save.md` |
| 2. 地形類型對映表（confirmed） ▸ 0 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 3 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 4 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 5 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 7 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 8 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 9 | `docs/mechanics/30-combat.md` |
| 3. 軍團記錄的欄位（部分） ▸ +0x01 | `docs/re/05-battle-selection.md` |
| 4. 據點記錄（32 byte） ▸ +1 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +10 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +26 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +8 | `docs/formats/08-sinario-save.md` |
| 據點記錄再解出四個欄位 ▸ +11h | `docs/re/07-monthly-settlement.md` |
| 據點記錄再解出四個欄位 ▸ +12h | `docs/re/07-monthly-settlement.md` |
| 據點記錄再解出四個欄位 ▸ +13h | `docs/re/07-monthly-settlement.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x04 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x20 | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x00 | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x02 | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x06 | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x0A | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x0C | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x0E | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x0F | `docs/re/08-hourly-update.md` |

### 說明書（8 條）

| 鍵 | 出處 |
|---|---|
| 1. 三個能力值 ▸ +17 | `docs/mechanics/60-personnel.md` |
| 1. 三個能力值 ▸ +18 | `docs/mechanics/60-personnel.md` |
| 1. 三個能力值 ▸ +19 | `docs/mechanics/60-personnel.md` |
| 3.4 其他戰術判定（說明書） ▸ 中央突破戰法 | `docs/mechanics/70-ai.md` |
| 3.4 其他戰術判定（說明書） ▸ 擊破狀態 | `docs/mechanics/70-ai.md` |
| 3.4 其他戰術判定（說明書） ▸ 突擊時機 | `docs/mechanics/70-ai.md` |
| 3.4 其他戰術判定（說明書） ▸ 陣形有利不利 | `docs/mechanics/70-ai.md` |
| 5.8 ⭐ 十一個命令處理常式，隊長與隊員各一套 ▸ 命令 | `docs/re/11-tactical-battle.md` |

### 未解（33 條）

| 鍵 | 出處 |
|---|---|
| 1. 檔案 ＝ 4 個劇本區塊 × 22,208 B ▸ +0x0008 | `docs/formats/08-sinario-save.md` |
| 1. 檔案 ＝ 4 個劇本區塊 × 22,208 B ▸ +0x1EC0 | `docs/formats/08-sinario-save.md` |
| 1. 檔案 ＝ 4 個劇本區塊 × 22,208 B ▸ +0x3AC0…+0x42C0 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x16 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x19 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x3F | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x17 | `docs/formats/08-sinario-save.md` |
| 2. 記憶體佈局：sub_1CC31 ▸ ds:0D2F8 | `docs/re/11-tactical-battle.md` |
| 2. 記憶體佈局：sub_1CC31 ▸ ds:0D2FC | `docs/re/11-tactical-battle.md` |
| 2. 記憶體佈局：sub_1CC31 ▸ ds:0D306 | `docs/re/11-tactical-battle.md` |
| 2. 記憶體佈局：sub_1CC31 ▸ ds:0D42 | `docs/re/11-tactical-battle.md` |
| 3. 武將記錄（32 byte） ▸ +0 | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +16（0x10） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +20（0x14） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +25,+27 | `docs/formats/08-sinario-save.md` |
| 3. 軍團記錄的欄位（部分） ▸ +0x08 | `docs/re/05-battle-selection.md` |
| 4. 只在單邊存在的檔（confirmed） ▸ PASS.MAP／PASS.SCH | `docs/re/01-first-recon.md` |
| 4. 據點記錄（32 byte） ▸ +0 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +12 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +14 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +16 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +17 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +18 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +22 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +27–+31 | `docs/formats/08-sinario-save.md` |
| 6. 月結與季節掛在時鐘的哪裡（confirmed） ▸ 世界更新 | `docs/mechanics/15-realtime.md` |
| 檔案 ▸ 10-strategy.md | `docs/mechanics/00-index.md` |
| 檔案 ▸ 40-economy.md | `docs/mechanics/00-index.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x08 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x0C | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x23 | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x04 | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x08 | `docs/re/08-hourly-update.md` |

