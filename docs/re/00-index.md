# 00 — RE 知識庫入口

**狀態：索引。34 份反組譯筆記按子系統分類，每份一句話說它回答什麼問題。**

- 日期：2026-08-14
- 範圍：松崗 DOS/V `KI.EXE` 為主，PC-98 只在兩版對照處出現
- 目標檔身分：`KI.EXE` SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`

這一份回答「**我想知道 X，該讀哪一份**」。
另外兩張表回答別的問題，動手前一起查（`CLAUDE.md` §10）：

| 想知道 | 查 |
|---|---|
| 某個欄位／常數解了沒 | [`../INDEX.md`](../INDEX.md) 的斷言總表 |
| 某支函式有人讀過嗎 | [`21`](21-function-census.md) 的覆蓋地圖 |
| 某支函式大概在做什麼 | [`24`](24-unread-function-catalogue.md) 的未讀目錄 |
| **我想了解某個子系統** | **本檔** |

---

## 1. 資料表與記憶體佈局

段內基址是所有位址算術的起點，讀任何一份筆記之前先記住這幾個數字：

```
ds:0840h  據點表  192 × 32 B        ds:2240h  軍團表  127 × 64 B
ds:4200h  城兵臨時軍團 1 × 64 B     ds:4240h  武將表  128 × 32 B
勢力表 22 × 64 B（編號 × 64 即記錄位址）
```

| 筆記 | 回答什麼 |
|---|---|
| [`01`](01-first-recon.md) | 兩版有哪些檔、哪 23 個 byte-for-byte 相同 |
| [`04`](04-mmap-entry-points.md) | 大地圖載入的入口點與記憶體佈局 |
| [`14`](14-mmap-mch-objects.md) | `MMAP.MCH` 的戰略地圖物件 |
| [`34`](34-corps-status-bits.md) | 軍團 `+0x00` 的位元圖：哪個位元表示「有行軍指令」「行軍中」 |

> 欄位層級的定義不在 `docs/re/`，在
> [`../formats/08-sinario-save.md`](../formats/08-sinario-save.md)。
> `docs/re/` 回答「程式碼在哪」，`docs/formats/` 回答「資料長什麼樣」。

## 2. 時間與世界更新（即時制的骨架）

| 筆記 | 回答什麼 |
|---|---|
| [`06`](06-game-clock.md) | 五層時間單位、一日 216 tick、速度節流 |
| [`08`](08-hourly-update.md) | 每「時」跑什麼；**每 tick 只更新 16 支軍團** |
| [`07`](07-monthly-settlement.md) | 月結的九支子程式；也是 AI 事件的總表 |
| [`16`](16-idle-clock-event10.md) | 沒有輸入時時鐘怎麼繼續走 |
| [`10`](10-rng.md) | 亂數在哪裡介入 |

## 3. 戰略層：玩家怎麼下指令

| 筆記 | 回答什麼 |
|---|---|
| [`22`](22-strategy-command-tree.md) | **指令列與熱區分派**——整個戰略層的入口 |
| [`25`](25-message-variants-and-personnel.md) | 訊息索引 ≥ `0x196` 的 ×8 變體；人事四個指令 |
| [`30`](30-corps-formation-ui.md) | 軍團編成：兵員池、兵種循環、移動間隔公式 |
| [`31`](31-faction-picker-screen.md) | 勢力一覽（同時是領地圖）；主畫面的重繪分派表 |
| [`32`](32-strategy-detail-panels.md) | 據點詳情欄與軍團詳情欄的欄位對照 |
| [`26`](26-list-window-engine.md) | 一覽表視窗引擎與 23 個 callback |
| [`27`](27-list-row-fields.md) | 一覽表怎麼畫每一列；外交關係等級 |
| [`13`](13-pc98-numeric-window.md) | 數字輸入視窗的量測與 CJK 版面決策 |

**指令列共八項**：進言／人事／財政／編成／軍團／據點／武將／勢力
（字串在 `cs:6181h`，見 [`31`](31-faction-picker-screen.md) §1）。

## 4. 戰術層：戰場

| 筆記 | 回答什麼 |
|---|---|
| [`05`](05-battle-selection.md) | 政略與戰術的接縫：戰場怎麼被選出來 |
| [`11`](11-tactical-battle.md) | 戰術模組結構與戰場資料模型 |
| [`09`](09-combat.md) | 觸發、自動判定、傷亡、武將的下場 |
| [`19`](19-outcome.md) | 敗北 outcome 的接線 |
| [`20`](20-ida-re-coverage-audit.md) | 戰術管線被移植進 remake 的程度（**量的是移植度不是覆蓋度**）|

## 5. 呈現層：圖、字、聲

| 筆記 | 回答什麼 |
|---|---|
| [`02`](02-palette-routine.md) | `.BRG` 的通道順序與亮度縮放 |
| [`03`](03-image-blitter.md) | 圖庫載入器與 VGA 繪製 |
| [`28`](28-text-number-rendering.md) | **文字與數字怎麼畫**：EGA Set/Reset、屬性的 4-bit 配對 |
| [`29`](29-font-service-int15.md) | **原版怎麼顯示中文**：`INT 15h AH=50h` ＋ 倚天字型 |
| [`33`](33-shared-draw-helpers.md) | 共用繪圖層：字串包裝、肖像四格快取、小地圖上色 |
| [`18`](18-tactical-button-glyphs.md) | 戰術底列按鈕的 glyph 資產 |
| [`17`](17-dosv-audio-tsr.md) | DOS/V 音源 TSR 與戰術效果碼 |
| [`23`](23-bgm-resource-format.md) | `*BGM.DAT` 的音樂資源格式 |

## 6. 外交與訊息

| 筆記 | 回答什麼 |
|---|---|
| [`12`](12-diplomacy-dialogue.md) | 停戰說服的訊息索引 `#190`–`#198` |
| [`15`](15-event10-producer.md) | 事件 10 producer 的深度逆向 |

## 7. 方法論與量測（讀之前先看這幾條）

| 筆記 | 回答什麼 |
|---|---|
| [`21`](21-function-census.md) | 全函式覆蓋普查；**§3.1 為什麼要排除目錄型文件** |
| [`24`](24-unread-function-catalogue.md) | 未讀函式的證據與下手順序 |
| [`39`](39-remaining-unread.md) | **剩餘 90 支的逐支歸屬**（生成的，可重跑）|
| [`20`](20-ida-re-coverage-audit.md) | 手挑取樣的問題，與全量量測的差別 |
| [`34`](34-corps-status-bits.md) §3 | IDAPython 的用法與兩個仍然成立的坑 |

### 這個專案反覆踩到的四種錯

寫進規則是為了下一輪不必重踩。完整清單在 `CLAUDE.md` §7，這裡只列與 RE 直接相關的：

1. **零命中要先問掃描條件對不對。** IDC 缺函式在 headless 底下是靜默中止，
   輸出看起來像「跑完了，沒命中」（[`30`](30-corps-formation-ui.md) §8）。
   **掃描腳本一律帶正對照。**
2. **表歸屬是函式層級不是指令層級。** 暫存器會跨函式帶著走
   （[`34`](34-corps-status-bits.md) §2.4）。
3. **判斷暫存器語意要看被呼叫的那一支**，不是看數值像什麼
   （[`22`](22-strategy-command-tree.md) §3.3）。
4. **有辦法直接讀的東西不要用呼叫關係去推。** 印一張分派表只要一行程式碼
   （`CONTEXT.md` 的推翻表）。

## 8. 現況與缺口

完整清單見 [`39`](39-remaining-unread.md)。

| 模組 | 未讀 | T4 bytes |
|---|---:|---:|
| 戰術戰鬥 | 37 | 1,144 |
| 訊息格式化與 TALK 輸出 | 13 | 664 |
| 月結、經濟、AI 決策 | 12 | 339 |
| 圖庫解碼與繪圖底層 | 11 | 360 |
| 啟動、C runtime 與低階 I/O | 9 | 200 |
| 戰略資料存取 | 6 | 185 |
| 戰略 UI | 2 | 26 |
| 戰術呈現 | 0 | 0 |

**合計 90 支、2,918 bytes（5%）。** 模組級全圖見
[`35`](35-strategy-ui-module-map.md)（戰略 UI）、
[`36`](36-tactical-module-map.md)（戰術戰鬥）、
[`37`](37-graphics-and-runtime-module-map.md)（圖庫與 C runtime）、
[`38`](38-strategy-core-module-map.md)（指令樹、到站處理、月結）。

**數字是 2026-08-14 的快照，要現況一律重跑**
（`tools/py.sh tools/re_coverage.py workplace/ida/dosv/census/census.tsv`）。

## 9. 怎麼加一份新筆記

1. 編號流水，一個發現一份。標題寫「回答什麼問題」不是「讀了哪支函式」。
2. 檔頭要有：狀態行、日期、範圍（哪一版）、輸入檔 SHA-256、`.i64` SHA-256、
   工具、位址空間。**沒有這幾項的筆記無法被後續驗證。**
3. 每條結論標推論等級：confirmed／強證據／假說／未知。
4. **正文只寫現況**，推翻紀錄集中到 `CONTEXT.md` 的「已被推翻的斷言」。
   `tools/index.py` 有一道閘檢查這一點（它比對的是措辭，所以引用這條規則本身
   也會被攔——換個說法就好）。
5. 提交前跑 `tools/check.sh`（六道閘）。
