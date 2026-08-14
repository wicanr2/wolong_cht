# 43 — 未解缺口總表（生成的）

**狀態：生成的清單，跑 `tools/py.sh tools/re_open_questions.py` 重出。
這一份不下結論，只把各文件的「未解」表集中到一處。**

- 日期：2026-08-14
- 產生工具：`tools/re_open_questions.py`
- 來源：`docs/` 底下所有文件的未解小節、表格裡標未解的列，與收尾是「…未解」的散句

既有的三張表回答別的問題（`CLAUDE.md` §10）：`docs/INDEX.md` 是**已解**的斷言、
[`21`](21-function-census.md) 是函式有沒有人寫過、[`24`](24-unread-function-catalogue.md) 是未讀函式在做什麼。
**這一份是唯一回答「還有什麼沒解」的。**

> 「擋住什麼」由來源目錄決定，「怎麼裁決」由關鍵字決定——兩欄都是機械算出來的，
> **不是逐條判斷過的優先序**。要排優先序請自己讀那一列指到的小節。

## 1. 總量

| 擋住什麼 | 缺口數 | 靜態可解 | 要實測 | 兩版對照 |
|---|---:|---:|---:|---:|
| 規則正確性 | 7 | 7 | 0 | 0 |
| 資料保存 | 28 | 28 | 0 | 0 |
| 程式碼理解 | 82 | 80 | 2 | 0 |
| 外部資料 | 9 | 9 | 0 | 0 |
| **合計** | **126** | 124 | 2 | 0 |

其中 **1 條明講被防拷擋著**——那條路沒通之前，它們不會因為多讀組語而前進。

## 2.1 規則正確性（7 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`mechanics/15-realtime.md`](../mechanics/15-realtime.md) | **世界更新** | 每「時」 / `sub_13E11`（內容未解） | 靜態 |
| [`mechanics/20-military.md`](../mechanics/20-military.md) | 玩家六個位置如何完整影響戰力仍未解 | （散句） | 靜態 |
| [`mechanics/40-economy.md`](../mechanics/40-economy.md) | 那兩個歸零的欄位很可能就是**本月累計的收入與支出 → 假說，待驗 | （散句） | 靜態 |
| [`mechanics/60-personnel.md`](../mechanics/60-personnel.md) | 政治如何影響內政效果與外交官要價仍未解。 | （散句） | 靜態 |
| [`mechanics/70-ai.md`](../mechanics/70-ai.md) | 交友關係的實際範圍、每月增量、協力成功的門檻值。 | （散句） | 靜態 |
| [`mechanics/70-ai.md`](../mechanics/70-ai.md) | cx` 在 `≥ 0x100` 時多半是另一個編號空間或帶旗標，未解。 | （散句） | 靜態 |
| [`mechanics/80-victory.md`](../mechanics/80-victory.md) | 四個劇本的結束條件是否不同 | `OPEN_S*` / `END_S*` 檔（未解） | 靜態 |

## 2.2 資料保存（28 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`formats/02-brg-palette.md`](../formats/02-brg-palette.md) | 4–7 | **未解**，高飽和純色。假說見 `docs/re/02` §6.2 | 靜態 |
| [`formats/03-grf-images.md`](../formats/03-grf-images.md) | 3 | `0x9700` / `0x23A0` (9,120) / `sub_1006B` / 走 `sub_1F888`（**位元對齊**的繪製常式，可放在非 8 倍數的 x）。未解 | 靜態 |
| [`formats/04-map-sch-container.md`](../formats/04-map-sch-container.md) | 狀態：容器格式的索引層 READY，壓縮演算法未解。 | （散句） | 靜態 |
| [`formats/05-mmap-worldmap.md`](../formats/05-mmap-worldmap.md) | 比對的正是方向碼那個欄位）。 強證據，未定案。 | （散句） | 靜態 |
| [`formats/07-battle.md`](../formats/07-battle.md) | 推論等級：**分段、筆數、緩衝區大小 confirmed；圖塊尺寸未解 | （散句） | 靜態 |
| [`formats/07-battle.md`](../formats/07-battle.md) | mov     al, es:[bx+1]       ; → byte_1AB4F（用途未解） | （散句） | 靜態 |
| [`formats/07-battle.md`](../formats/07-battle.md) | 表頭與尾段各 64 byte，內容未解。 | （散句） | 靜態 |
| [`formats/07-battle.md`](../formats/07-battle.md) | ——表頭筆數多於資料組數，原因未解。 | （散句） | 靜態 |
| [`formats/07-battle.md`](../formats/07-battle.md) | 其餘 61,440 B 是像素資料，格式未解。 | （散句） | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0x0008` | 51 / 未解的全域狀態（一起載入 `cs:0CF0h`） | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0x0080` | 2,112 / **勢力表：22 筆 × 64 B**（`docs/re/06` §5）＋ 其後未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0x1EC0` | 7,168 / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0x3AC0`…`+0x42C0` | — / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0x16` | u16 / 恆 `FFFF` / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0x19` | u8 / 恆 `FF` / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0` | 1 / 旗標（7 種值）。**bit 4 ＝ 不事二主**：舊主已滅時被俘會自刎（`sub_129C3` → 訊息 0x43） / 其餘位元未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+16`（`0x10`） | 1 / 同格式，戰略層沒有呼叫點取得到（要模式 2）。**計入評價** / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+20`（`0x14`） | 1 / 值域 0–7 / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+25`,`+27` | 2 / 含 `0xFF` 哨兵 / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0` | 1 / 值域 0–15，14 種。**據點類型？** / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+12` | 2 / uint16，北京 21000、代 4000、涿郡 4300 / 未解（人口／資金？） | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+14` | 2 / uint16，北京 20714、代 3000、涿郡 3142 / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+16` | 1 / 值域 10–179 / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+17` | 1 / **全部 192 筆都是 100** / 未解（防災值初值？） | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+18` | 2 / uint16 / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+22` | 1 / 值域 0–226，22 種 / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+27`–`+31` | 5 / 含 `0xFF` 哨兵 / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0`／`+3` | 未解 | 靜態 |

## 2.3 程式碼理解（82 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`re/01-first-recon.md`](../re/01-first-recon.md) | 是加了新過場、還是把原本的長段拆開，未解。 | （散句） | 靜態 |
| [`re/01-first-recon.md`](../re/01-first-recon.md) | `PASS.MAP`／`PASS.SCH` | dosv / **PC-98 沒有**。關隘資料，移植時新增或改名。未解 | 靜態 |
| [`re/03-image-blitter.md`](../re/03-image-blitter.md) | sub_1FAC2` 是另一支繪製常式（`shl al, 1` 後才 `mov cx, ax`），用途未解。 | （散句） | 靜態 |
| [`re/04-mmap-entry-points.md`](../re/04-mmap-entry-points.md) | call sub_1E3C0                    ; ← 未解 | （散句） | 靜態 |
| [`re/05-battle-selection.md`](../re/05-battle-selection.md) | 推論等級：**兩條路徑的存在與參數 confirmed**；地形對映表的內容未解 | （散句） | 靜態 |
| [`re/05-battle-selection.md`](../re/05-battle-selection.md) | `sub_14C4C` 是「地圖圖塊值 → 地形類型」的對映 | **未解** | 靜態 |
| [`re/05-battle-selection.md`](../re/05-battle-selection.md) | `+0x08` | 勢力相符時取用的值 / 未解 | 靜態 |
| [`re/06-game-clock.md`](../re/06-game-clock.md) | 勢力記錄 64 B 的欄位表 | `+0`、`+1Ah`、`+1Ch` 已知，其餘未解 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | 123 個武將裡有 80 人只有一欄非零、32 人兩欄、9 人三欄—— → 強證據，待驗 | （散句） | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | 連結表每筆的實際佈局（`+8` 為何落在下一筆） | 佈局仍未解，但**容量的矛盾解掉了**（見 §7.1） | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `+0x00` bit 6 | 資金低於門檻一半時設起（用途未解） | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `+0x08` | 建立時寫 `4` / 未解 | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `+0x0C` | 行軍中的暫存（`sub_12708` 寫） / 未解 | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `+0x23` | 建立時寫 `1` / 未解 | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `+0x04` | byte / `sub_1E81C` 的回傳 `ah` / 未解 | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `+0x08` | word / 另一張表的索引（`bx << 3`） / 未解 | 靜態 |
| [`re/09-combat.md`](../re/09-combat.md) | `sub_129C3` 的 `[bx+17h] = 4` | 武將 `+0x17` 未解 | 靜態 |
| [`re/11-tactical-battle.md`](../re/11-tactical-battle.md) | `ds:0D2FC` | 8,192 / 未解 | 靜態 |
| [`re/11-tactical-battle.md`](../re/11-tactical-battle.md) | `ds:0D2F8` | 4,096 / 未解（第二份戰場？） | 靜態 |
| [`re/11-tactical-battle.md`](../re/11-tactical-battle.md) | `ds:0D306` | 30,720 / 未解 | 靜態 |
| [`re/11-tactical-battle.md`](../re/11-tactical-battle.md) | `ds:0D42` | 4,000 / 未解（與大地圖的連結表同大小） | 靜態 |
| [`re/22-strategy-command-tree.md`](../re/22-strategy-command-tree.md) | 指令列圖形來源與熱區圖的登記方式未解。 | （散句） | 靜態 |
| [`re/22-strategy-command-tree.md`](../re/22-strategy-command-tree.md) | 指令列的圖形來源 | 八個 40×16 按鈕的圖從哪來未讀。`ICONGRF` 段 1 仍未解（[`../formats/03`](../formats/03-grf-images.md) §5.4），尺寸上是候選，但**沒有證據**，不要當結論 | 靜態 |
| [`re/22-strategy-command-tree.md`](../re/22-strategy-command-tree.md) | 熱區圖怎麼登記 | 寫入端 `sub_1E41B` 沒有任何文件提過 | 靜態 |
| [`re/22-strategy-command-tree.md`](../re/22-strategy-command-tree.md) | `off_159D2` 的其餘槽位 | 16 筆裡只有 [1] `0x1614A`、[2] `0x15E1E`、[3] `0x15A3A`、[4] `sub_15FAA`、[13] `sub_161CA` 非 `nullsub_1`。這是頂層模式分派表，[1]–[3] 未讀 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | 聲軌資料本身的事件編碼、`+0x02`／`+0x04` 兩張表的內容與音色語意仍未解。 | （散句） | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | `+0x00` | 2 B / 未解 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | `+0x02` | 2 B / 指向靠近曲尾的一張表；載入時存進 `cs:099Ah`。**內容未解** | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | `+0x04` | 2 B / 指向另一張表；存進 `cs:099Ch`。它是曲塊的**最後一段**，長 `0x30` B。**內容未解** | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | `+0x06`–`+0x0F` | 10 B / 未解 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | 聲軌資料的事件編碼 | 完全未解。音高、長度、迴圈、包絡都還沒讀 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | `+0x00`、`+0x06`–`+0x0F` | 12 B 未解 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | DOS/V 的實際音源晶片 | `17` §4 有 register path，卡種仍未定案 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | 曲號 ↔ 場景的對應 | `KI.EXE` 呼叫端傳哪個索引還沒對過 | 靜態 |
| [`re/25-message-variants-and-personnel.md`](../re/25-message-variants-and-personnel.md) | `sub_16B4F`／`sub_16C2A` | 解任的實際動作未讀 | 靜態 |
| [`re/25-message-variants-and-personnel.md`](../re/25-message-variants-and-personnel.md) | `sub_17400`／`sub_17906`／`sub_17663` | 選據點／選勢力／選武將三支一覽表未讀 | 靜態 |
| [`re/25-message-variants-and-personnel.md`](../re/25-message-variants-and-personnel.md) | `sub_11D46` | 四支人事指令離開前都呼叫它，未讀 | 靜態 |
| [`re/25-message-variants-and-personnel.md`](../re/25-message-variants-and-personnel.md) | `sub_1075B` 尾端的 `sub_106F9` | 收到 `ax = 0xFF00 \ / ah`；`al=[武將+1]`（頭像編號）在這條路徑上似乎沒被用到，要再讀 | 靜態 |
| [`re/26-list-window-engine.md`](../re/26-list-window-engine.md) | `sub_1748F` | 畫一列的實作（193 B），未讀 | 靜態 |
| [`re/26-list-window-engine.md`](../re/26-list-window-engine.md) | `word_183D3` | `si` 也存進這裡，讀取端還沒找 | 靜態 |
| [`re/26-list-window-engine.md`](../re/26-list-window-engine.md) | `funcs_18450` 五個 handler | 捲動與選取的實際行為（`nullsub_2`／`sub_18463`／`sub_184DD`／`sub_1851A`／`sub_18546`） | 靜態 |
| [`re/26-list-window-engine.md`](../re/26-list-window-engine.md) | `sub_11D46` | 17 個呼叫點，人事四支離開前都呼叫，未讀 | 靜態 |
| [`re/27-list-row-fields.md`](../re/27-list-row-fields.md) | 開局選勢力的逐列 `sub_17BC0` | 未逐欄對照（欄位與勢力一覽重疊，但少了外交兩欄） | 靜態 |
| [`re/28-text-number-rendering.md`](../re/28-text-number-rendering.md) | `sub_1F7A4` | 把 32 B 字模緩衝畫上 VRAM 的實際迴圈，未逐行讀 | 靜態 |
| [`re/28-text-number-rendering.md`](../re/28-text-number-rendering.md) | 屬性的其餘位元 | bit 2 是陰影已證實；`0x9001`／`0x9000` 的 bit 0 差在哪未讀 | 靜態 |
| [`re/28-text-number-rendering.md`](../re/28-text-number-rendering.md) | `word_10D54` | 數字字型段的來源未追。**數字不走 `INT 15h`**（`sub_1069A` 直接從 `ds` 取點陣），所以與 §29 是兩條路 | 靜態 |
| [`re/29-font-service-int15.md`](../re/29-font-service-int15.md) | `END_S10/S11` 與 `STR.EXE` 檔名不同步 | §6，要實跑裁決 | 實測 |
| [`re/29-font-service-int15.md`](../re/29-font-service-int15.md) | `END_S13.DAT` 前 408 格的來源 | 不是 `stdfont.15` 的任何一段，也不是 `usrfont.15m`（256 B） | 靜態 |
| [`re/29-font-service-int15.md`](../re/29-font-service-int15.md) | `END_S15.DAT`（5,242 B） | `KI.EXE` 的字串表引用它，但**不是字型**（大小不是 30／15 的倍數），也不是過場圖（沒有 `00 F4 01` 檔頭），Big5 解碼是亂碼 → 疑似壓縮 | 靜態 |
| [`re/29-font-service-int15.md`](../re/29-font-service-int15.md) | `sub_1F7A4` | 把 32 B 緩衝畫上 VRAM 的實際迴圈，未逐行讀 | 靜態 |
| [`re/29-font-service-int15.md`](../re/29-font-service-int15.md) | `YNFONT.EXE` 怎麼顯示中文 | 它不走 INT 15h（0 次），防拷畫面的中文是它自己畫的。與本鏈無關，仍未解 | 靜態・**防拷擋著** |
| [`re/30-corps-formation-ui.md`](../re/30-corps-formation-ui.md) | 軍團 `+0x23` | `sub_16F26` 寫 1，用途未解 | 靜態 |
| [`re/30-corps-formation-ui.md`](../re/30-corps-formation-ui.md) | 軍團 `+0x00` 的位元 | 建立時 `0xC0`，AI 路徑（`sub_16E8F`）另外 `or 4`。各位元語意未解 | 靜態 |
| [`re/30-corps-formation-ui.md`](../re/30-corps-formation-ui.md) | `sub_1D4C7` | 大地圖上實際畫圖塊的常式，未讀 | 靜態 |
| [`re/30-corps-formation-ui.md`](../re/30-corps-formation-ui.md) | `sub_16D6F`／`sub_16DA8` | 兩支印數字的常式，畫的是哪幾個欄位未逐一對過 | 靜態 |
| [`re/31-faction-picker-screen.md`](../re/31-faction-picker-screen.md) | `cs:6056` 表的長度 | 前六筆是一組小 handler，後五筆疑似越過表尾（§1.2） | 靜態 |
| [`re/31-faction-picker-screen.md`](../re/31-faction-picker-screen.md) | `cs:byte_198A6` | 位元 1 對應勢力一覽；其餘位元未解 | 靜態 |
| [`re/31-faction-picker-screen.md`](../re/31-faction-picker-screen.md) | `sub_15AD1 → sub_15AFC` 的進入路徑 | 不在 `off_159D2` 表裡，未定位 | 靜態 |
| [`re/32-strategy-detail-panels.md`](../re/32-strategy-detail-panels.md) | `sub_107D2` | `sub_1807B` 用它畫主將（`bx=3EB7h`），疑似肖像，未讀 | 靜態 |
| [`re/32-strategy-detail-panels.md`](../re/32-strategy-detail-panels.md) | `sub_1812A`／`sub_1817D` | 軍團面板的前後收尾，未讀 | 靜態 |
| [`re/32-strategy-detail-panels.md`](../re/32-strategy-detail-panels.md) | 軍團 `+0x00` 的位元怎麼清 | 三處設定都找到了，清除點未找到 | 靜態 |
| [`re/32-strategy-detail-panels.md`](../re/32-strategy-detail-panels.md) | `sub_14325` | 下完行軍令之後呼叫，未讀 | 靜態 |
| [`re/32-strategy-detail-panels.md`](../re/32-strategy-detail-panels.md) | `+0x16` 高 4 位 | 仍未解（`0`–`14`），這兩支面板都沒用到它 | 靜態 |
| [`re/33-shared-draw-helpers.md`](../re/33-shared-draw-helpers.md) | `sub_1E38C` | 從圖庫段載入 `di` bytes，內部未讀 | 靜態 |
| [`re/33-shared-draw-helpers.md`](../re/33-shared-draw-helpers.md) | `cs:word_10D40` | 肖像圖庫所在的段，誰載入它未追 | 靜態 |
| [`re/33-shared-draw-helpers.md`](../re/33-shared-draw-helpers.md) | `sub_1FA37` 的 `ax = 4004h` | 圖塊尺寸的編碼方式未逐位對過（[`18`](18-tactical-button-glyphs.md) §2 只定了「ax ＝ 尺寸」） | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | 0 | `sub_12459`／`sub_126FF`（候選） / `sub_12533`（候選） / 未定 | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | 4 | `sub_12B3C`（候選） / `sub_12BA8`（候選） / 未定。`sub_12B3C` 呼叫 `sub_15D19`（小地圖畫點） | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | 位元 0／4／5 的語意 | 有成對的設定與清除點，但那幾支函式的操作對象尚未逐支確認是軍團 | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | 位元 3 | 掃描裡沒出現。間接寫入抓不到，不能據此說它不存在 | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | `sub_12977` 的 `mov byte [si], 8` | 該函式同時碰武將表與軍團表，`si` 指哪一張未確認 | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | `sub_144A9`／`sub_144D6`／`sub_12880` | 表歸屬指向據點表，`or [si], 2`／`or [si], 20h` 的語意要另外讀 | 靜態 |
| [`re/37-graphics-and-runtime-module-map.md`](../re/37-graphics-and-runtime-module-map.md) | `sub_1F7A4` | 212 / 字型 blitter，[`29`](29-font-service-int15.md) §9 已列為未解 | 靜態 |
| [`re/40-garrison-relief-request.md`](../re/40-garrison-relief-request.md) | `sub_13FA9`（127 B） | 填那 16 B 威脅緩衝的，未讀 | 靜態 |
| [`re/40-garrison-relief-request.md`](../re/40-garrison-relief-request.md) | `sub_140B3`（22 B） | 求援後的收尾，未讀 | 靜態 |
| [`re/40-garrison-relief-request.md`](../re/40-garrison-relief-request.md) | `sub_14575`（76 B） | AI 側取「該勢力的某個目標」，未讀 | 靜態 |
| [`re/40-garrison-relief-request.md`](../re/40-garrison-relief-request.md) | 據點 `+0x14`／`+0x18` | 一個被 `sub_13FA9` 寫、一個當門檻，語意未定 | 靜態 |
| [`re/40-garrison-relief-request.md`](../re/40-garrison-relief-request.md) | `+0x20` 與 `+0x14` 的關係 | §5 的張力，要實測 | 實測 |
| [`re/42-leaf-functions.md`](../re/42-leaf-functions.md) | `INT 61h` 的四個服務號（`ah=4`／`7`／`8`、`ax=09F2h`／`0C01h`） | 對應什麼音效動作要看 `YNSOUND.COM`（[`17`](17-dosv-audio-tsr.md)） | 靜態 |
| [`re/42-leaf-functions.md`](../re/42-leaf-functions.md) | `cs:byte_198A6` 位元 3 | 對應 `sub_15FAA`，設定與清除端都沒找到 | 靜態 |
| [`re/42-leaf-functions.md`](../re/42-leaf-functions.md) | `sub_1E9A7` 的 8 bytes 參數表 | 表本身沒讀 | 靜態 |
| [`re/42-leaf-functions.md`](../re/42-leaf-functions.md) | `byte_1020E`／`byte_10CF9` | 音源相關的兩個旗標 | 靜態 |

## 2.5 外部資料（9 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`reference/03-baked-japanese.md`](../reference/03-baked-japanese.md) | `END_S3`–`END_S11` | 九張全不同 / 有燒字，**松崗重繪過** | 靜態 |
| [`reference/03-baked-japanese.md`](../reference/03-baked-japanese.md) | `END_S1`／`S2`／`S12` | 相同 / 沒有燒字 | 靜態 |
| [`reference/03-baked-japanese.md`](../reference/03-baked-japanese.md) | `OPEN_S1`–`S6` | 六張全相同 / 開場沒有燒字，文字是疊繪的 | 靜態 |
| [`reference/03-baked-japanese.md`](../reference/03-baked-japanese.md) | **`ICONGRF.DAT`** | **相同** / **沒重繪 → 裡面的日文留著** | 靜態 |
| [`reference/03-baked-japanese.md`](../reference/03-baked-japanese.md) | `KAOGRF`／`KYOGRF`／`IVENTGRF` | 相同 / 純圖像，目前看過的都沒有文字 | 靜態 |
| [`reference/04-first-survey.md`](../reference/04-first-survey.md) | `SHOW.O` | 57,148 B / 被 `INSTALL.EXE` 與 `LOGO.EXE` 引用。開頭 `3c df 00 00 11 af 01 00 50 00 80 07`。**未解** | 靜態 |
| [`reference/04-first-survey.md`](../reference/04-first-survey.md) | `KYOGRF.DAT` | 69,120 / 未解 | 靜態 |
| [`reference/04-first-survey.md`](../reference/04-first-survey.md) | FM 3 聲 ＋ SSG 3 聲，埠 `0x188`／`0x18A`。 DOS/V 側未解。 | （散句） | 靜態 |
| [`reference/05-eten-font-provenance.md`](../reference/05-eten-font-provenance.md) | `END_S13/S14/S15` 是中文版加的結局段 | S13／S14 是字型。**`END_S15` 仍未解** | 靜態 |

## 3. 這支工具的盲區

抽取只認四種結構（專門的未解小節、表格最後一欄標未解的列、
`**未解**：…`、收尾是「…未解」的句子）。**寫在段落中段、
或用別的詞說「這個還不知道」的缺口抽不到**——下列檔案提到未解
卻一列都沒抽出來，要嘛缺口寫成別的句式，要嘛那些字樣只是在講別的事：

- `docs/formats/01-talk-dat.md`
- `docs/formats/06-mmap-rle.md`
- `docs/mechanics/00-index.md`
- `docs/mechanics/10-strategy.md`
- `docs/mechanics/30-combat.md`
- `docs/mechanics/50-diplomacy.md`
- `docs/playtest/10-event-message-modal.md`
- `docs/promo/dosv-adlib-and-tactical-review.md`
- `docs/promo/yt-remake-pixel-review.md`
- `docs/re/00-index.md`
- `docs/re/02-palette-routine.md`
- `docs/re/12-diplomacy-dialogue.md`
- `docs/re/13-pc98-numeric-window.md`
- `docs/re/15-event10-producer.md`
- `docs/re/17-dosv-audio-tsr.md`
- `docs/re/19-outcome.md`
- `docs/re/20-ida-re-coverage-audit.md`
- `docs/re/21-function-census.md`
- `docs/re/24-unread-function-catalogue.md`
- `docs/re/35-strategy-ui-module-map.md`
- `docs/re/39-remaining-unread.md`
- `docs/reference/02-jp-cht-diff.md`

只印抽得到的部分，會讓解析失敗長得像「那份文件沒有缺口」。
這一節就是為了讓那個差別看得見。
