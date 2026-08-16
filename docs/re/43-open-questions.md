# 43 — 未解缺口總表（生成的）

**狀態：生成的清單，跑 `tools/py.sh tools/re_open_questions.py` 重出。
這一份不下結論，只把各文件的「未解」表集中到一處。**

- 日期：2026-08-16
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
| 規則正確性 | 54 | 52 | 2 | 0 |
| 資料保存 | 37 | 37 | 0 | 0 |
| 程式碼理解 | 215 | 208 | 7 | 0 |
| 驗收 | 54 | 48 | 6 | 0 |
| 外部資料 | 17 | 16 | 0 | 1 |
| 其他 | 91 | 86 | 5 | 0 |
| **合計** | **468** | 447 | 20 | 1 |

## 2.1 規則正確性（54 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`mechanics/10-strategy.md`](../mechanics/10-strategy.md) | **防災值**怎麼成長 | 欄位語意已知（據點 `+0x11`），成長規則未解 | 靜態 |
| [`mechanics/10-strategy.md`](../mechanics/10-strategy.md) | 行軍費用 | **說明書 10.6「行軍の費用」有專節，還沒讀** | 靜態 |
| [`mechanics/15-realtime.md`](../mechanics/15-realtime.md) | **世界更新** | 每「時」 / `sub_13E11`（內容未解） | 靜態 |
| [`mechanics/15-realtime.md`](../mechanics/15-realtime.md) | `sub_13E11` 每「時」做什麼 | 行軍與 AI 的節拍，寫到那一層時再解 | 靜態 |
| [`mechanics/15-realtime.md`](../mechanics/15-realtime.md) | `sub_10A65` 的內插演算法 | 只影響畫面 | 靜態 |
| [`mechanics/15-realtime.md`](../mechanics/15-realtime.md) | 最高速那一檔在原版實機上是多少 | 機器相依，要實測才有數字；只影響手感調校 | 實測 |
| [`mechanics/20-military.md`](../mechanics/20-military.md) | 玩家六個位置如何完整影響戰力仍未解 | （散句） | 靜態 |
| [`mechanics/20-military.md`](../mechanics/20-military.md) | 軍團編成的數值判定（六個位置怎麼影響戰力，見 `30-combat.md`） | （未解小節內文） | 靜態 |
| [`mechanics/20-military.md`](../mechanics/20-military.md) | 路上（192–255）與野外（≥256）節點：remake 目前只用據點層的圖。 | （未解小節內文） | 靜態 |
| [`mechanics/20-military.md`](../mechanics/20-military.md) | sub_14502` 第二段的方向（`docs/re/08` §7.7） | （未解小節內文） | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | 戰術層不再是空白：`sub_19FA0` 的入口、腳本、戰場資料模型、核心移動／一般近戰／ | （未解小節內文） | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | 0 | **陣形** / 走到指定座標，到位就轉成「已就位」狀態 | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | 1 | **攻擊** / 大將（`+0x04 == 0`）不攻擊，只移動 | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | 2 | **突擊** / 攻城戰時多跑一支開門的常式 | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | 3 | **城壁移動** / 野戰時自動變成「攻擊」；腳本只在攻城段下它 | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | 4 | **守陣** / 敵人曼哈頓距離 ≤ **16 格**就換行為 | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | 5 | **退卻** / 受傷自動改成它、隊長倒下全隊改成它、不可打斷 | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | **疲勞度** | 兵記錄 `+0x19` / 走到陣形位置**補滿 128**、攻擊時上限壓到 **40**、移動每幀 **−1** | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | **士氣值** | 待查 / — | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | **將軍體力** | 大將的 `+0x03` / **低於 50 自動退卻**；攻城方**每 10 幀 −1** | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | **有利／不利** | `byte_1D31E` / 兩軍「還有餘力的兵數」相減，**差距 ≤ 8 判成普通** | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | 大將／騎馬 | 目標座標 ← 敵人的座標，**追過去打** | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | **弓兵** | 目標 ← **自己的**座標（不動），依與敵人的**高度差**分支 | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | 步兵 | 近戰，而且**挨箭只吃四分之一傷害** | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | 繞路點是誰算出來的（真正的尋路演算法），以及士氣值存在哪。 | （未解小節內文） | 靜態 |
| [`mechanics/40-economy.md`](../mechanics/40-economy.md) | 預備兵維持費的單價（三兵種是否不同） | `sub_15358` 尾段的批次呼叫 | 靜態 |
| [`mechanics/40-economy.md`](../mechanics/40-economy.md) | 支出（`+0x1A`／`+0x1C`）是**誰累加**的 | 找寫入點 | 靜態 |
| [`mechanics/40-economy.md`](../mechanics/40-economy.md) | 防災值、城兵數的月度變化 | `sub_15358` 尾段剩下的六支 | 靜態 |
| [`mechanics/40-economy.md`](../mechanics/40-economy.md) | 「北方／南方」的判定邊界（座標？據點旗標？） | `SINARIO.DAT` 據點表 ＋ 反組譯募兵 | 靜態 |
| [`mechanics/40-economy.md`](../mechanics/40-economy.md) | 行軍啟動費用的計算基準 | 反組譯 | 靜態 |
| [`mechanics/40-economy.md`](../mechanics/40-economy.md) | 災害的實際傷害量 | `sub_12FBF` 的事件表（`010Ch` 火災／`020Ch` 暴動／`0Bh` 暴風雨） | 靜態 |
| [`mechanics/40-economy.md`](../mechanics/40-economy.md) | 那兩個歸零的欄位很可能就是**本月累計的收入與支出 → 假說，待驗 | （散句） | 靜態 |
| [`mechanics/50-diplomacy.md`](../mechanics/50-diplomacy.md) | **六階的門檻** | 反組譯顯示程式 | 靜態 |
| [`mechanics/50-diplomacy.md`](../mechanics/50-diplomacy.md) | `g(對方君主好戰等級)` 的實際修正量 | 反組譯外交結算 | 靜態 |
| [`mechanics/50-diplomacy.md`](../mechanics/50-diplomacy.md) | 「良好」對應的實際交友值 | 反組譯 ＋ `SINARIO.DAT` | 靜態 |
| [`mechanics/50-diplomacy.md`](../mechanics/50-diplomacy.md) | 外交官每月提升的量與金額的關係 | 反組譯月結 | 靜態 |
| [`mechanics/60-personnel.md`](../mechanics/60-personnel.md) | 政治如何影響內政效果與外交官要價仍未解。 | （散句） | 靜態 |
| [`mechanics/60-personnel.md`](../mechanics/60-personnel.md) | 武術的勝敗判定分布 | 反組譯一騎打ち | 靜態 |
| [`mechanics/60-personnel.md`](../mechanics/60-personnel.md) | 評價實際餵進哪些判定 | 反組譯戰鬥結算 | 靜態 |
| [`mechanics/60-personnel.md`](../mechanics/60-personnel.md) | `+0Eh`／`+0Fh`／`+10h` 是不是兵種適性 | 與戰鬥程式對照 | 靜態 |
| [`mechanics/60-personnel.md`](../mechanics/60-personnel.md) | 政治如何影響內政效果與外交官要價 | 反組譯 | 靜態 |
| [`mechanics/60-personnel.md`](../mechanics/60-personnel.md) | `+0` 旗標的 7 種值代表什麼 | 身分已確定在 `+0x17`，`+0` 是別的東西。只知道 bit 4 ＝ 不事二主 | 靜態 |
| [`mechanics/60-personnel.md`](../mechanics/60-personnel.md) | `+0x1E` 為什麼 0–2 恆空 | 掃 127 筆武將的 `+0x1E` 值分佈 | 靜態 |
| [`mechanics/60-personnel.md`](../mechanics/60-personnel.md) | 被俘之後的處理 | 未知 | 靜態 |
| [`mechanics/60-personnel.md`](../mechanics/60-personnel.md) | 信賴度存在哪 | `SINARIO.DAT`／`SAVE.DAT` 每區塊 `+0x10`，對應 `cs:0D00h`／IDA `byte_10D00`；勢力 `+0x1D` 是士氣基準 | 靜態 |
| [`mechanics/70-ai.md`](../mechanics/70-ai.md) | 交友關係的實際範圍、每月增量、協力成功的門檻值。 | （散句） | 靜態 |
| [`mechanics/70-ai.md`](../mechanics/70-ai.md) | 10 | `sub_13496` / 訊息-only 邊界；保留事件取出，不虛構持久欄位 / handler 只把 `AH`／`DX` 組成 `sub_18810` formatter 參數，尚未完成 TALK.DAT 逐句顯示 | 靜態 |
| [`mechanics/70-ai.md`](../mechanics/70-ai.md) | cx` 在 `≥ 0x100` 時多半是另一個編號空間或帶旗標，未解。 | （散句） | 靜態 |
| [`mechanics/80-victory.md`](../mechanics/80-victory.md) | 信賴度的值域 | **0…255，byte 飽和**（`seg000:10D00`、`sub_13D91`／`13DC9`）→ confirmed | 靜態 |
| [`mechanics/80-victory.md`](../mechanics/80-victory.md) | 初始／新遊戲初始化值 | 原始劇本 `+0x10` 可直接讀；第 1 劇本目前為 `0xFF`；`sub_18B12` 完整時序 → 強證據／待 oracle | 實測 |
| [`mechanics/80-victory.md`](../mechanics/80-victory.md) | 玩家進言／說得的信賴度增減量 | **已證實**：`sub_13830` 的第一反應為 `+20`／`−20`，多理由完成為 `+10`，錯選理由為 `−20`；事件 13 的 `−50` 另由 `sub_13507` 定案。事件 2／3 等外交回報的其他增減仍需逐分支對拍 | 靜態 |
| [`mechanics/80-victory.md`](../mechanics/80-victory.md) | 四個劇本的結束條件是否不同 | `OPEN_S*` / `END_S*` 檔（未解）。**觸發條件本身四劇本共用**——差別只在初始勢力數 | 靜態 |
| [`mechanics/80-victory.md`](../mechanics/80-victory.md) | 君主陣亡時軍師怎麼辦 | 未知 | 靜態 |
| [`mechanics/80-victory.md`](../mechanics/80-victory.md) | END_S1`–`END_S4`（結局動畫？）與四個劇本的結局有關，格式還沒碰。 | （未解小節內文） | 靜態 |

## 2.2 資料保存（37 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`formats/01-talk-dat.md`](../formats/01-talk-dat.md) | 六個變數標記各自代表什麼（要從 `KI.EXE` 反追）。 | （未解小節內文） | 靜態 |
| [`formats/01-talk-dat.md`](../formats/01-talk-dat.md) | 偏移表 `[1023] = 0` 是保留還是有意義。 | （未解小節內文） | 靜態 |
| [`formats/01-talk-dat.md`](../formats/01-talk-dat.md) | 78 則空訊息的索引位置有沒有規律。 | （未解小節內文） | 靜態 |
| [`formats/01-talk-dat.md`](../formats/01-talk-dat.md) | 訊息索引與遊戲事件的對應（哪一則在什麼時候顯示）。 | （未解小節內文） | 靜態 |
| [`formats/02-brg-palette.md`](../formats/02-brg-palette.md) | 誰載入、誰選組（`docs/re/02` §7）。 | （未解小節內文） | 靜態 |
| [`formats/02-brg-palette.md`](../formats/02-brg-palette.md) | OPENPAL` 6 組與 `ENDPAL` 12 組各自對應哪些畫面。 | （未解小節內文） | 靜態 |
| [`formats/03-grf-images.md`](../formats/03-grf-images.md) | 3 | `0x9700` / `0x23A0` (9,120) / `sub_1006B` / 走 `sub_1F888`（**位元對齊**的繪製常式，可放在非 8 倍數的 x）。未解 | 靜態 |
| [`formats/03-grf-images.md`](../formats/03-grf-images.md) | `0x0480` | 24×16 × 3 / 兵種圖示的**橘色版**：馬／弓／步 / 尚未找到取用端 | 靜態 |
| [`formats/04-map-sch-container.md`](../formats/04-map-sch-container.md) | 狀態：容器格式的索引層 READY，壓縮演算法未解。 | （散句） | 靜態 |
| [`formats/05-mmap-worldmap.md`](../formats/05-mmap-worldmap.md) | MMAP.MAP`（80,716 B）的編碼**。 | （未解小節內文） | 靜態 |
| [`formats/05-mmap-worldmap.md`](../formats/05-mmap-worldmap.md) | MMAP.MCH` 的 256×160-byte MCH 圖塊、0xA000 metadata、事件 12 火災／暴動 | （未解小節內文） | 靜態 |
| [`formats/05-mmap-worldmap.md`](../formats/05-mmap-worldmap.md) | sub_1E717` 建出來的記錄實際被誰用（§3）。 | （未解小節內文） | 靜態 |
| [`formats/05-mmap-worldmap.md`](../formats/05-mmap-worldmap.md) | byte_1E47E` 與 `byte_1E47F` 兩個流水號的分工（`sub_1E567` 寫前者、 | （未解小節內文） | 靜態 |
| [`formats/05-mmap-worldmap.md`](../formats/05-mmap-worldmap.md) | 比對的正是方向碼那個欄位）。 強證據，未定案。 | （散句） | 靜態 |
| [`formats/07-battle.md`](../formats/07-battle.md) | 戰場的 `64 × 62` 哪一邊是寬 | 格數確定（§2.1），方向沒定 | 靜態 |
| [`formats/07-battle.md`](../formats/07-battle.md) | 每個戰場 4,096 B 的表頭與尾段各 64 byte | 內容未解（§2.1） | 靜態 |
| [`formats/07-battle.md`](../formats/07-battle.md) | 表頭 16 筆但資料只有 3 組 | 多出來的 13 筆是什麼未解（§3） | 靜態 |
| [`formats/07-battle.md`](../formats/07-battle.md) | mov     al, es:[bx+1]       ; → byte_1AB4F（用途未解） | （散句） | 靜態 |
| [`formats/07-battle.md`](../formats/07-battle.md) | 表頭與尾段各 64 byte，內容未解。 | （散句） | 靜態 |
| [`formats/07-battle.md`](../formats/07-battle.md) | ——表頭筆數多於資料組數，原因未解。 | （散句） | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0x0008` | 51 / 未解的全域狀態（一起載入 `cs:0CF0h`） | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0x0080` | 2,112 / **勢力表：22 筆 × 64 B**（`docs/re/06` §5）＋ 其後未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0x1EC0` | 7,168 / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0x3AC0`…`+0x42C0` | — / 未解 | 靜態 |
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

## 2.3 程式碼理解（215 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`re/00-index.md`](../re/00-index.md) | `43` | **缺口總表**：各文件「未解」表的集中版（生成的，`check.sh` 每次重出） | 靜態 |
| [`re/00-index.md`](../re/00-index.md) | 739 支函式的 T4 全部歸零**——每一支都有 `docs/re/` 層級的記錄。 | （未解小節內文） | 靜態 |
| [`re/01-first-recon.md`](../re/01-first-recon.md) | 是加了新過場、還是把原本的長段拆開，未解。 | （散句） | 靜態 |
| [`re/01-first-recon.md`](../re/01-first-recon.md) | `YNMOUSE.COM` | pc98 / 滑鼠驅動。dosv 版把它併進 `KI.EXE`？未驗 | 靜態 |
| [`re/01-first-recon.md`](../re/01-first-recon.md) | `PASS.MAP`／`PASS.SCH` | dosv / **PC-98 沒有**。關隘資料，移植時新增或改名。未解 | 靜態 |
| [`re/02-palette-routine.md`](../re/02-palette-routine.md) | 誰把 `.BRG` 的內容載到 `cs:1964h`（`sub_109AF` 載到 `cs:18A4h`， | （未解小節內文） | 靜態 |
| [`re/02-palette-routine.md`](../re/02-palette-routine.md) | OPENPAL`（6 組）、`ENDPAL`（12 組）的分組對應哪些畫面。 | （未解小節內文） | 靜態 |
| [`re/02-palette-routine.md`](../re/02-palette-routine.md) | 設定表 `cs:0x5FF4` 每筆後三個 byte 是什麼（第 4 筆的第二個 word | （未解小節內文） | 靜態 |
| [`re/03-image-blitter.md`](../re/03-image-blitter.md) | ICONGRF`／`KYOGRF`／`IVENTGRF` 的尺寸。 | （未解小節內文） | 靜態 |
| [`re/03-image-blitter.md`](../re/03-image-blitter.md) | sub_1E38C`（帶位移的讀檔）與 `sub_1F4A2`（開檔／讀檔／關檔）的完整介面。 | （未解小節內文） | 靜態 |
| [`re/03-image-blitter.md`](../re/03-image-blitter.md) | sub_1FAC2` 是另一支繪製常式（`shl al, 1` 後才 `mov cx, ax`），用途未解。 | （未解小節內文） | 靜態 |
| [`re/04-mmap-entry-points.md`](../re/04-mmap-entry-points.md) | call sub_1E3C0                    ; ← 未解 | （散句） | 靜態 |
| [`re/04-mmap-entry-points.md`](../re/04-mmap-entry-points.md) | MMAP.MAP`（80,716 B）的編碼**——展開成 98,304 B（＝ 384 × 256 格）。 | （未解小節內文） | 靜態 |
| [`re/04-mmap-entry-points.md`](../re/04-mmap-entry-points.md) | sub_20000` 的三次呼叫在設定什麼（參數像是範圍：`17FFh`、`101Fh`）。 | （未解小節內文） | 靜態 |
| [`re/04-mmap-entry-points.md`](../re/04-mmap-entry-points.md) | word_19872` 的 98,304 B 畫布怎麼對應到螢幕 | （未解小節內文） | 靜態 |
| [`re/05-battle-selection.md`](../re/05-battle-selection.md) | 推論等級：**兩條路徑的存在與參數 confirmed**；地形對映表的內容未解 | （散句） | 靜態 |
| [`re/05-battle-selection.md`](../re/05-battle-selection.md) | `sub_14C4C` 是「地圖圖塊值 → 地形類型」的對映 | **未解** | 靜態 |
| [`re/05-battle-selection.md`](../re/05-battle-selection.md) | `+0x08` | 勢力相符時取用的值 / 未解 | 靜態 |
| [`re/06-game-clock.md`](../re/06-game-clock.md) | `sub_13E11`——每「時」的世界更新（AI？行軍？） | 直接讀 | 靜態 |
| [`re/06-game-clock.md`](../re/06-game-clock.md) | `sub_15358` 呼叫的 9 支子程式各做什麼 | 逐一讀，經濟公式都在裡面 | 靜態 |
| [`re/06-game-clock.md`](../re/06-game-clock.md) | 勢力記錄 64 B 的欄位表 | `+0`、`+1Ah`、`+1Ch` 已知，其餘未解 | 靜態 |
| [`re/06-game-clock.md`](../re/06-game-clock.md) | `cs:0CF0h` 那 59 byte 裡除了時鐘的其餘部分 | 對照存檔 diff | 靜態 |
| [`re/06-game-clock.md`](../re/06-game-clock.md) | `sub_10A65` 的內插演算法 | 直接讀 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | 10 | `sub_13496` / 訊息-only：建立武將／參數 formatter 游標；持久狀態尚未找到 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | `sub_15940` 的兩個分支 | 已派駐武將的每月行動，會發訊息 `0x41`／`0x42`。分支 2 有一行 `mov byte ptr [si+1Ch], 18h`（把所屬勢力寫成 24）**與「+1Ch 是勢力編號、只有 0–21」矛盾**，還沒解釋 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | `sub_12BD9` | 已讀：對 22 個勢力各建 0x30 的候選緩衝區，串起交友度排序、協力／停戰／宣戰產生器與遷都事件 8 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | `sub_14269`／`sub_13EFD` | 事件 11／12 寫入的城市 `+0x15` marker 在據點輪轉時先扣防災值；不足時再扣上昇值、生產力與城兵，已接入 `World.applyCityDisasterEffect`；物件動畫仍未完 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | 武將 `+1Ah` | 官員「要錢中」的旗標／金額，`sub_12FBF` 的事件會寫它 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | 武將 `+0Eh`／`+0Fh`／`+10h` 是不是兵種適性 | 需要與戰鬥程式對照 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | `sub_157FE` 觸發的事件內容 | `sub_12FBF(ax=0Dh, dx=196h)` | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | 軍團記錄 32 B 的欄位表 | 段內 `2240h`，開局全零（沒有軍團） | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | 段內 `2040h`／`2140h` 的兩張 16 × 16 B 表 | `sub_123FF` 會在 `2040h` 那張找空位配置；`2140h` 那張開局已有 16 筆，`+2`／`+4` 看起來是地圖座標 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | 123 個武將裡有 80 人只有一欄非零、32 人兩欄、9 人三欄—— → 強證據，待驗 | （散句） | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | 連結表每筆的實際佈局（`+8` 為何落在下一筆） | 佈局仍未解，但**容量的矛盾解掉了**（見 §7.1） | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `sub_15456` 用 stride 32 掃軍團表 | 與 64 矛盾，疑似原版 bug | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `sub_14300`（AI 行軍判斷）、`sub_142AB`（野外移動） | 直接讀 | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | 軍團 `+0x0A` ＝ 行進方向的值域 | 直接讀 | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | 軍團 `+0x08`／`+0x23` | 建立時寫 4 與 1 | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | 信賴度存在哪 | **區塊 `+0x10`，對應 `cs:0D00h`／IDA `byte_10D00`**；勢力 `+0x1D` 是士氣基準、武將 `+0x1D` 是捕虜狀態。初始／新遊戲初始化的完整時序與所有事件增減量仍待 oracle | 實測 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `+0x00` bit 6 | 資金低於門檻一半時設起（用途未解） | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `+0x08` | 建立時寫 `4` / 未解 | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `+0x0C` | 行軍中的暫存（`sub_12708` 寫） / 未解 | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `+0x23` | 建立時寫 `1` / 未解 | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `+0x04` | byte / `sub_1E81C` 的回傳 `ah` / 未解 | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `+0x08` | word / 另一張表的索引（`bx << 3`） / 未解 | 靜態 |
| [`re/09-combat.md`](../re/09-combat.md) | **戰術完整結算** | `TestNormalScenarioTacticalBattleTerminates` 已證實真實正常攻城的狀態層勝負／傷亡回寫，`wlgame-ai-postbattle.png` 證明正常 GUI 回戰略；GUI 戰後訊息、完整狀態對拍與少數分支仍未完 | 靜態 |
| [`re/09-combat.md`](../re/09-combat.md) | `sub_1AD7F` 攻擊分支 | `shootSpecial` 已接入 `CH=0x20` 的相鄰格／垂直效果；`+0x1E` 的初始化／上移／下移／交換來源與 `sub_1AC55` 的 raw 比較已確認並接成 `PlaneHigh`，普通箭原版 SCH 單幀圖形已接回，完整投射物動畫／同狀態對拍仍待確認 | 靜態 |
| [`re/09-combat.md`](../re/09-combat.md) | 原版／remake 同狀態對拍 | 需有效時序原版存檔或可重建的同狀態 oracle | 實測 |
| [`re/09-combat.md`](../re/09-combat.md) | `sub_14C72`：怎麼挑出對手軍團 | 野戰與攻城共用，回傳 `bx` | 靜態 |
| [`re/09-combat.md`](../re/09-combat.md) | 地形係數表的列 2、武將適性 `+0x10` | 兩者都要 `al = 2` 才取得到，戰略層沒有呼叫點 | 靜態 |
| [`re/09-combat.md`](../re/09-combat.md) | 武將旗標 `+0x00` 的 bit 4（自刎） | 值域已知有 7 種，只解出這一個位元 | 靜態 |
| [`re/09-combat.md`](../re/09-combat.md) | 據點 `+0x10`／`+0x11` 被攻城扣減 | 欄位語意已知（上昇值／防災值），但「被打過的城成長變慢」還沒在數值上驗過 | 靜態 |
| [`re/09-combat.md`](../re/09-combat.md) | `[si+3]` 的 0／1／≥2 是誰設的 | 決定哪一支軍團會進戰鬥畫面 | 靜態 |
| [`re/09-combat.md`](../re/09-combat.md) | `sub_129C3` 的 `[bx+17h] = 4` | 武將 `+0x17` 未解 | 靜態 |
| [`re/11-tactical-battle.md`](../re/11-tactical-battle.md) | `ds:0D2FC` | 8,192 / 未解 | 靜態 |
| [`re/11-tactical-battle.md`](../re/11-tactical-battle.md) | `ds:0D2F8` | 4,096 / 未解（第二份戰場？） | 靜態 |
| [`re/11-tactical-battle.md`](../re/11-tactical-battle.md) | `ds:0D306` | 30,720 / 未解 | 靜態 |
| [`re/11-tactical-battle.md`](../re/11-tactical-battle.md) | `ds:0D42` | 4,000 / 未解（與大地圖的連結表同大小） | 靜態 |
| [`re/11-tactical-battle.md`](../re/11-tactical-battle.md) | 按「解開之後能換到什麼」排序： | （未解小節內文） | 靜態 |
| [`re/11-tactical-battle.md`](../re/11-tactical-battle.md) | `ds:0D2FC`／`0D306` | 通行圖與未解 / — | 靜態 |
| [`re/12-diplomacy-dialogue.md`](../re/12-diplomacy-dialogue.md) | `AH` 的完整欄位名稱 | 語意由日中原文並列確認，欄位名本身未定（§3） | 靜態 |
| [`re/12-diplomacy-dialogue.md`](../re/12-diplomacy-dialogue.md) | #367–#372／#380–#385 的 AH／信賴度次要回覆 | 未解，不可當成完整的原版對話流程（§8） | 靜態 |
| [`re/12-diplomacy-dialogue.md`](../re/12-diplomacy-dialogue.md) | #73／#77 | 未定位，不得拿來補接事件 6／7（§9） | 靜態 |
| [`re/12-diplomacy-dialogue.md`](../re/12-diplomacy-dialogue.md) | 事件 6／7 次要 TALK 的 formatter 參數契約 | 缺參數且語意未知，維持 fail-closed（§10） | 靜態 |
| [`re/13-pc98-numeric-window.md`](../re/13-pc98-numeric-window.md) | sub_17D0D`（`00017D0D`）使用 `DS:SI=word_10D50:0600h`、`AX=4006h`； | （未解小節內文） | 靜態 |
| [`re/13-pc98-numeric-window.md`](../re/13-pc98-numeric-window.md) | 依 `sub_100DF` 的段指標配置，`word_10D50:0600h` 換算為 | （未解小節內文） | 靜態 |
| [`re/13-pc98-numeric-window.md`](../re/13-pc98-numeric-window.md) | sub_17C6E` caller 的保存區仍從 `(80,176)`、112×80 開始；96×64 內框目的地是 | （未解小節內文） | 靜態 |
| [`re/13-pc98-numeric-window.md`](../re/13-pc98-numeric-window.md) | cmd/wlgame` 有 DOS/V 資源時直接繪製這張內框，缺資源才降級通用框；數值格的 raw | （未解小節內文） | 靜態 |
| [`re/15-event10-producer.md`](../re/15-event10-producer.md) | 以下來源沒有證據，不能補成事實：未被 IDA 建成函式的 far code、以暫存器或指標 | （未解小節內文） | 靜態 |
| [`re/17-dosv-audio-tsr.md`](../re/17-dosv-audio-tsr.md) | `0x330` 的用途 | MPU-401 的標準埠，沒找到讀它的地方 | 靜態 |
| [`re/17-dosv-audio-tsr.md`](../re/17-dosv-audio-tsr.md) | `INT 61h` 的四個服務號 | `ah=4`／`7`／`8` 與 `ax=09F2h`／`0C01h`，對應什麼動作要看 `YNSOUND.COM`（`42` §7） | 靜態 |
| [`re/19-outcome.md`](../re/19-outcome.md) | 勢力滅亡 selector | 未定位。remake 只顯示克制的 fallback 句，不冒充原版文字 | 靜態 |
| [`re/20-ida-re-coverage-audit.md`](../re/20-ida-re-coverage-audit.md) | §5 的四個 gate 裡，前兩個已由 §7／§8 兩個切片收掉；剩下的是這些： | （未解小節內文） | 靜態 |
| [`re/20-ida-re-coverage-audit.md`](../re/20-ida-re-coverage-audit.md) | 同狀態動態 oracle | 沒有可重放的存檔／輸入序列，所以「原版等價」目前無法驗。**這是還沒做，不是做不了**——DOS/V 的密碼頁空白確認就會過（`../playtest/18`），PC-98 側連除錯器都接好了（`../playtest/21`） | 實測 |
| [`re/20-ida-re-coverage-audit.md`](../re/20-ida-re-coverage-audit.md) | 逐幀執行順序 | 顯示串列與相機已重建，但整幀的呼叫順序沒有逐幀對過 | 靜態 |
| [`re/20-ida-re-coverage-audit.md`](../re/20-ida-re-coverage-audit.md) | `loc_1A065` 的 runtime bytes | 自我修改碼，靜態影像看不到每輪的實際內容（§2.2） | 靜態 |
| [`re/20-ida-re-coverage-audit.md`](../re/20-ida-re-coverage-audit.md) | 四層差分（terrain／display list／composited／HUD） | 沒有 machine-readable diff，目前只有 layout-only 比較 | 靜態 |
| [`re/22-strategy-command-tree.md`](../re/22-strategy-command-tree.md) | 指令列圖形來源與熱區圖的登記方式未解。 | （散句） | 靜態 |
| [`re/22-strategy-command-tree.md`](../re/22-strategy-command-tree.md) | 熱區圖怎麼登記 | 寫入端 `sub_1E41B` 沒有任何文件提過 | 靜態 |
| [`re/22-strategy-command-tree.md`](../re/22-strategy-command-tree.md) | `off_159D2` 的其餘槽位 | 16 筆裡只有 [1] `0x1614A`、[2] `0x15E1E`、[3] `0x15A3A`、[4] `sub_15FAA`、[13] `sub_161CA` 非 `nullsub_1`。這是頂層模式分派表，[1]–[3] 未讀 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | `+0x00` | 2 B / 未解 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | `+0x06`–`+0x0F` | 10 B / 未解 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | `+0x04` 那張表 | 大小與音效記錄相同（3 × 16 B），但驅動沒讀它。見 `57` §8 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | `+0x00`、`+0x06`–`+0x0F` | 12 B 未解 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | 曲號 ↔ 場景的對應 | `KI.EXE` 呼叫端傳哪個索引還沒對過 | 靜態 |
| [`re/25-message-variants-and-personnel.md`](../re/25-message-variants-and-personnel.md) | `sub_16B4F`／`sub_16C2A` | 解任的實際動作未讀 | 靜態 |
| [`re/25-message-variants-and-personnel.md`](../re/25-message-variants-and-personnel.md) | `sub_17400`／`sub_17906`／`sub_17663` | 選據點／選勢力／選武將三支一覽表未讀 | 靜態 |
| [`re/25-message-variants-and-personnel.md`](../re/25-message-variants-and-personnel.md) | `sub_11D46` | 四支人事指令離開前都呼叫它，未讀 | 靜態 |
| [`re/25-message-variants-and-personnel.md`](../re/25-message-variants-and-personnel.md) | `sub_1075B` 尾端的 `sub_106F9` | 收到 `ax = 0xFF00 \ / ah`；`al=[武將+1]`（頭像編號）在這條路徑上似乎沒被用到，要再讀 | 靜態 |
| [`re/26-list-window-engine.md`](../re/26-list-window-engine.md) | `sub_1748F` | 畫一列的實作（193 B），未讀 | 靜態 |
| [`re/26-list-window-engine.md`](../re/26-list-window-engine.md) | `word_183D3` | `si` 也存進這裡，讀取端還沒找 | 靜態 |
| [`re/26-list-window-engine.md`](../re/26-list-window-engine.md) | `sub_11D46` | 17 個呼叫點，人事四支離開前都呼叫，未讀 | 靜態 |
| [`re/27-list-row-fields.md`](../re/27-list-row-fields.md) | 開局選勢力的逐列 `sub_17BC0` | 未逐欄對照（欄位與勢力一覽重疊，但少了外交兩欄） | 靜態 |
| [`re/28-text-number-rendering.md`](../re/28-text-number-rendering.md) | `sub_1F7A4` | 把 32 B 字模緩衝畫上 VRAM 的實際迴圈，未逐行讀 | 靜態 |
| [`re/28-text-number-rendering.md`](../re/28-text-number-rendering.md) | 屬性的其餘位元 | bit 2 是陰影已證實；`0x9001`／`0x9000` 的 bit 0 差在哪未讀 | 靜態 |
| [`re/28-text-number-rendering.md`](../re/28-text-number-rendering.md) | `word_10D54` | 數字字型段的來源未追。**數字不走 `INT 15h`**（`sub_1069A` 直接從 `ds` 取點陣），所以與 §29 是兩條路 | 靜態 |
| [`re/29-font-service-int15.md`](../re/29-font-service-int15.md) | `END_S10/S11` 與 `STR.EXE` 檔名不同步 | §6，要實跑裁決 | 實測 |
| [`re/29-font-service-int15.md`](../re/29-font-service-int15.md) | `END_S13.DAT` 前 408 格的來源 | 不是 `stdfont.15` 的任何一段，也不是 `usrfont.15m`（256 B） | 靜態 |
| [`re/29-font-service-int15.md`](../re/29-font-service-int15.md) | `END_S15.DAT`（5,242 B） | `KI.EXE` 的字串表引用它，但**不是字型**（大小不是 30／15 的倍數），也不是過場圖（沒有 `00 F4 01` 檔頭），Big5 解碼是亂碼 → 疑似壓縮 | 靜態 |
| [`re/29-font-service-int15.md`](../re/29-font-service-int15.md) | `sub_1F7A4` | 把 32 B 緩衝畫上 VRAM 的實際迴圈，未逐行讀 | 靜態 |
| [`re/29-font-service-int15.md`](../re/29-font-service-int15.md) | `YNFONT.EXE` 怎麼顯示中文 | 它不走 INT 15h（0 次），密碼輸入畫面的中文是它自己畫的。與本鏈無關，仍未解 | 靜態 |
| [`re/30-corps-formation-ui.md`](../re/30-corps-formation-ui.md) | 軍團 `+0x23` | `sub_16F26` 寫 1，用途未解 | 靜態 |
| [`re/30-corps-formation-ui.md`](../re/30-corps-formation-ui.md) | `sub_1D4C7` | 大地圖上實際畫圖塊的常式，未讀 | 靜態 |
| [`re/31-faction-picker-screen.md`](../re/31-faction-picker-screen.md) | 分派表已印出，但 `sub_15AD1 → sub_15AFC` 的進入路徑仍未定位。 | （散句） | 靜態 |
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
| [`re/33-shared-draw-helpers.md`](../re/33-shared-draw-helpers.md) | `sub_1FA37` 的 `ax = 4004h` | 圖塊尺寸的編碼方式未逐位對過（`18` §2 只定了「ax ＝ 尺寸」） | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | 0 | `sub_12459`／`sub_126FF`（候選） / `sub_12533`（候選） / 未定 | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | 4 | `sub_12B3C`（候選） / `sub_12BA8`（候選） / 未定。`sub_12B3C` 呼叫 `sub_15D19`（小地圖畫點） | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | 位元 0／4／5 的語意 | 有成對的設定與清除點，但那幾支函式的操作對象尚未逐支確認是軍團 | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | 位元 3 | 掃描裡沒出現。間接寫入抓不到，不能據此說它不存在 | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | `sub_12977` 的 `mov byte [si], 8` | 該函式同時碰武將表與軍團表，`si` 指哪一張未確認 | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | `sub_144A9`／`sub_144D6`／`sub_12880` | 表歸屬指向據點表，`or [si], 2`／`or [si], 20h` 的語意要另外讀 | 靜態 |
| [`re/35-strategy-ui-module-map.md`](../re/35-strategy-ui-module-map.md) | `sub_18FC9` 叢 | — / 存檔畫面的槽位與按鈕對應未驗（§2.8） | 靜態 |
| [`re/37-graphics-and-runtime-module-map.md`](../re/37-graphics-and-runtime-module-map.md) | `sub_1F7A4` | 212 / 字型 blitter，`29` §9 已列為未解 | 靜態 |
| [`re/40-garrison-relief-request.md`](../re/40-garrison-relief-request.md) | 三支後續函式與兩個欄位由 `44` 收掉： | （未解小節內文） | 靜態 |
| [`re/40-garrison-relief-request.md`](../re/40-garrison-relief-request.md) | `+0x20` 與 `+0x14` 的關係 | §5 的張力，要實測 | 實測 |
| [`re/42-leaf-functions.md`](../re/42-leaf-functions.md) | `INT 61h` 的四個服務號（`ah=4`／`7`／`8`、`ax=09F2h`／`0C01h`） | 對應什麼音效動作要看 `YNSOUND.COM`（`17`） | 靜態 |
| [`re/42-leaf-functions.md`](../re/42-leaf-functions.md) | `cs:byte_198A6` 位元 3 | 對應 `sub_15FAA`，設定與清除端都沒找到 | 靜態 |
| [`re/42-leaf-functions.md`](../re/42-leaf-functions.md) | `sub_1E9A7` 的 8 bytes 參數表 | 表本身沒讀 | 靜態 |
| [`re/42-leaf-functions.md`](../re/42-leaf-functions.md) | `byte_1020E`／`byte_10CF9` | 音源相關的兩個旗標 | 靜態 |
| [`re/44-threat-and-reinforcement-ai.md`](../re/44-threat-and-reinforcement-ai.md) | `sub_14194`／`sub_14269` | `sub_13EFD` 每 tick 還呼叫這兩支，未讀 | 靜態 |
| [`re/44-threat-and-reinforcement-ai.md`](../re/44-threat-and-reinforcement-ai.md) | 軍團 `+0x23` | 兩條派兵路徑都寫 0、`sub_16F26` 寫 1、`sub_14155` 要求 `< 8`。像是階段計數，語意未定 | 靜態 |
| [`re/44-threat-and-reinforcement-ai.md`](../re/44-threat-and-reinforcement-ai.md) | 據點 `+0x00` 的 bit 4／5 | bit 6／7 是威脅旗標、低 4 位是敵方鄰居，中間兩位仍未見寫入端 | 靜態 |
| [`re/44-threat-and-reinforcement-ai.md`](../re/44-threat-and-reinforcement-ai.md) | 勢力 `+0x17` 的讀取端 | 寫入端在 §1，誰讀它未找 | 靜態 |
| [`re/44-threat-and-reinforcement-ai.md`](../re/44-threat-and-reinforcement-ai.md) | `+0x20` 與 `+0x14` 的張力 | `sub_14575` 與 `sub_14155` 都只寫 `+0x20`，`40` §5 的張力還在 | 靜態 |
| [`re/45-corps-command-mode.md`](../re/45-corps-command-mode.md) | `sub_193E9` 的選單協定 | `ah = 1`、`cx` ＝ 首項索引、`dx`／`bx` ＝ 位置；回傳值怎麼編碼未逐位對過 | 靜態 |
| [`re/45-corps-command-mode.md`](../re/45-corps-command-mode.md) | `sub_1703C` | 選據點的那一支，未讀 | 靜態 |
| [`re/45-corps-command-mode.md`](../re/45-corps-command-mode.md) | `sub_13E11`（每「時」） | 已讀：它每次只處理**一支**軍團（`word_10D1C` 每次 `+0x40`，繞到 `0x580` 歸零 ＝ 22 支輪替），條件是 `+0x00` 位元 7 設起來；接著算 `[si+23h] × 8 + 0x18` 拿去和 `+0x21`（u16）比，過了就寫 `[si+19h] = 0xFF`，再比一… | 靜態 |
| [`re/46-strategy-chrome-cell-layer.md`](../re/46-strategy-chrome-cell-layer.md) | 樣式碼 | 只確定 `0` ＝ 擦除、`0x0B` ＝ 指令列、`0x0C`／`0x0F` 出現在別處；完整值域未列 | 靜態 |
| [`re/46-strategy-chrome-cell-layer.md`](../re/46-strategy-chrome-cell-layer.md) | `ax = 0F01h`／`0801h` | 顏色／樣式的位元編碼未逐位對過 | 靜態 |
| [`re/47-main-screen-window-registry.md`](../re/47-main-screen-window-registry.md) | `0x80` | 繪製時 `and …, 7Fh` 清掉 / 未解 | 靜態 |
| [`re/47-main-screen-window-registry.md`](../re/47-main-screen-window-registry.md) | 熱區 5（x 464–496） | 登記了，`off_159D2` 對應 `nullsub_1`。是保留槽還是別處會改寫這一格，未讀 | 靜態 |
| [`re/47-main-screen-window-registry.md`](../re/47-main-screen-window-registry.md) | 選完君主之後的相機 | `sub_1D615(170, 98)` 只管 NEW GAME 對話框背後那張圖。主畫面開始時相機在哪、由誰寫，未讀——`word_1988E`／`word_19890` 的六個參考**全是讀**，寫入端走 `ds:988Eh` 這種形式，要用 `tools/ida_disp_users.py` 掃 | 靜態 |
| [`re/47-main-screen-window-registry.md`](../re/47-main-screen-window-registry.md) | 格子屬性 bit `0x80` | 擦除時被清掉，沒找到設它的地方 | 靜態 |
| [`re/47-main-screen-window-registry.md`](../re/47-main-screen-window-registry.md) | 系統視窗開著時時間停止 | 說明書明講，機器碼的實作位置未找（`sub_15FAA` 的等待迴圈是候選） | 靜態 |
| [`re/47-main-screen-window-registry.md`](../re/47-main-screen-window-registry.md) | `funcs_159C0` 的五筆內容 | 只確認是「擦除」對應表（`sub_1895D` 樣式 0），逐筆未 dump | 靜態 |
| [`re/47-main-screen-window-registry.md`](../re/47-main-screen-window-registry.md) | 開關圖示的圖形來源 | 五格 32×32 的圖從哪個圖庫來未讀 | 靜態 |
| [`re/47-main-screen-window-registry.md`](../re/47-main-screen-window-registry.md) | 原版執行期驗證 | **未做**。PC-98 oracle 上左鍵點這五格沒有反應（§6），原因未定 | 實測 |
| [`re/48-window-display-list.md`](../re/48-window-display-list.md) | `08` 的模式 byte | `03` 只畫字、`01` 連背景一起填，是**強推論**——兩個用例（系統選單的「 ＯＫ 」、注音聲母列）都只有這個讀法說得通，但 `sub_106F5` 沒逐行讀（`55` §3） | 靜態 |
| [`re/48-window-display-list.md`](../re/48-window-display-list.md) | `sub_1E9A7(bl=0, ax=1800h, cx=2020h)` | `sub_1030F` 登記的第二件事，未讀 | 靜態 |
| [`re/48-window-display-list.md`](../re/48-window-display-list.md) | `op 01` 的用法 | 它是直線（§2.2），但 handler 不展開座標而十個場景又沒用到它——**預期的呼叫方式無法驗證** | 靜態 |
| [`re/48-window-display-list.md`](../re/48-window-display-list.md) | `op 02` 與 `op 03` 的差別 | 兩支都畫矩形（`sub_1F020` 對 `cs:F1A3`），前者另有五個戰術區呼叫者。哪一支是實心、哪一支帶遮罩，沒有資料可分辨 | 靜態 |
| [`re/49-corps-formation-window.md`](../re/49-corps-formation-window.md) | `sub_1F9B0` 的 `ax = 1003h` | 貼圖的樣式參數；`sub_10C14` 用 `0801h`（`46` §3）。位元編碼未逐位對過 | 靜態 |
| [`re/49-corps-formation-window.md`](../re/49-corps-formation-window.md) | 段 3 `0x21A0` 那張空槽圖 | **已驗**：解得開、與綠組最後一張不同，內容是 384 個像素全 0——原版在空槽那一格貼一張全黑圖把前一張擦掉。「全黑」與「位移落在段尾之外」長得一樣，分辨的方法是看它後面還有沒有內容（`TestDOSVEmptySlotIcon`） | 靜態 |
| [`re/50-city-info-window.md`](../re/50-city-info-window.md) | `cs:word_1987C` | 據點圖的暫存段，由誰配置未讀 | 靜態 |
| [`re/50-city-info-window.md`](../re/50-city-info-window.md) | 熱區 | 這個視窗**沒有登記任何熱區**：`sub_17E1F` 的迴圈只等右鍵離開，是純顯示 | 靜態 |
| [`re/51-corps-info-window.md`](../re/51-corps-info-window.md) | `sub_17FDB` | 玩家軍團的指令輸入流程，未讀。`45` 解過「戰鬥指揮／委任／解體」那一段，兩者的接縫沒對過 | 靜態 |
| [`re/51-corps-info-window.md`](../re/51-corps-info-window.md) | `sub_14325` | 下完指令之後跑的一支，未讀 | 靜態 |
| [`re/51-corps-info-window.md`](../re/51-corps-info-window.md) | `or byte ptr [si], 2` | 位元 1 ＝「有指令」（`34`），這裡是它的其中一個寫入端 | 靜態 |
| [`re/52-slot-select-window.md`](../re/52-slot-select-window.md) | `ds:987Ch` | 四筆槽頭的暫存段，由誰配置未讀 | 靜態 |
| [`re/52-slot-select-window.md`](../re/52-slot-select-window.md) | 檔名 | `sub_18C20` 沒設 `dx`，靠 `sub_18B7C` 的 `push dx`／`pop dx` 從更上層傳進來 | 靜態 |
| [`re/52-slot-select-window.md`](../re/52-slot-select-window.md) | `sub_18C9F` | 關閉時擦除的那一支，未讀 | 靜態 |
| [`re/53-lord-select-window.md`](../re/53-lord-select-window.md) | 換勢力 | 這一支只顯示 `si` 指到的那一個勢力，**怎麼換下一個**在呼叫端 `sub_11AC3` 的迴圈裡，未讀 | 靜態 |
| [`re/53-lord-select-window.md`](../re/53-lord-select-window.md) | `ds:5222h` | 軍師名字的位置（推測），配置端未讀 | 靜態 |
| [`re/53-lord-select-window.md`](../re/53-lord-select-window.md) | `sub_18F6D`／`sub_18F7C` | 收尾與取消時擦除的兩支，未讀 | 靜態 |
| [`re/54-advisor-naming-window.md`](../re/54-advisor-naming-window.md) | 選字表 | 十個聲母各自對應哪些候選字、資料在哪、怎麼翻頁——全部未讀。`sub_18FC9`（呼叫端）是入口 | 靜態 |
| [`re/54-advisor-naming-window.md`](../re/54-advisor-naming-window.md) | 屬性低 byte | `01` 與 `03` 的差別未讀（§3） | 靜態 |
| [`re/54-advisor-naming-window.md`](../re/54-advisor-naming-window.md) | 「別　號」 | 軍師除了名字還有別號，寫進哪裡未讀 | 靜態 |
| [`re/54-advisor-naming-window.md`](../re/54-advisor-naming-window.md) | 「重來」「繼續」 | 三顆按鈕的 handler 未讀 | 靜態 |
| [`re/55-system-menu-window.md`](../re/55-system-menu-window.md) | 「資料儲存」與「遊戲結束」的 handler | `0x6084`／`0x60B4` 沒讀 | 靜態 |
| [`re/55-system-menu-window.md`](../re/55-system-menu-window.md) | `sub_15FAA` 的 `cmp bx, 0Ah` | 熱區碼 `0x2A` 不在這個視窗的 `0x20`–`0x25` 裡，哪來的沒查 | 靜態 |
| [`re/55-system-menu-window.md`](../re/55-system-menu-window.md) | `sub_106F5` 的屬性解碼 | §3 的低 byte 讀法是強推論，沒逐行驗 | 靜態 |
| [`re/55-system-menu-window.md`](../re/55-system-menu-window.md) | 設定表每筆的第 4 個 byte | 四筆都是 `00`，用途不明 | 靜態 |
| [`re/56-bgm-track-events.md`](../re/56-bgm-track-events.md) | 全音符 ＝ 192 tick | 從長度表的二分序列推的，**強證據不是 confirmed**。沒有樂譜可對 | 靜態 |
| [`re/56-bgm-track-events.md`](../re/56-bgm-track-events.md) | `+0x04` 那張表 | 見 `57` §8 | 靜態 |
| [`re/56-bgm-track-events.md`](../re/56-bgm-track-events.md) | PC-98 側 | 事件編碼共用，但音色與音源程式設計完全沒讀 | 靜態 |
| [`re/57-opl3-register-map.md`](../re/57-opl3-register-map.md) | 曲塊 `+0x04` 的表 | 固定 `0x30` B ＝ 3 × 16，與音效記錄同大小。但 parser 存進 `cs:099Ch` 之後，**整個驅動沒有任何一處讀它**（全庫掃立即值只有一筆寫入）。是舊版遺留還是由別處使用，未解 | 靜態 |
| [`re/57-opl3-register-map.md`](../re/57-opl3-register-map.md) | 曲塊 `+0x00`、`+0x06`–`+0x0F` | 12 B 未解 | 靜態 |
| [`re/57-opl3-register-map.md`](../re/57-opl3-register-map.md) | 音色記錄 `+0x16`–`+0x1F` | 驅動不讀，內容意義未解 | 靜態 |
| [`re/57-opl3-register-map.md`](../re/57-opl3-register-map.md) | `word [097Eh]` ＝ `0x0330` | MPU-401 的標準埠，但沒找到讀它的地方 | 靜態 |
| [`re/57-opl3-register-map.md`](../re/57-opl3-register-map.md) | 全域音量偏移 `[0996h]` | 誰設、範圍多少未解 | 靜態 |
| [`re/57-opl3-register-map.md`](../re/57-opl3-register-map.md) | 曲號 ↔ 場景 | `KI.EXE` 的呼叫端還沒對過（`23` §5） | 靜態 |
| [`re/57-opl3-register-map.md`](../re/57-opl3-register-map.md) | PC-98 側的音源程式設計 | 完全沒讀。YM2203 的暫存器路徑與音色版面都未解 | 靜態 |
| [`re/58-bgm-scene-mapping.md`](../re/58-bgm-scene-mapping.md) | 曲 1 | **DOS/V 的 `KI.EXE` 裡沒有任何呼叫端傳 1。** 掃過的範圍：`sub_10241` 的八個直接呼叫點（立即值全部列在 §3）、`cs:9309h` 那張表（只有 2–5）、`sub_19946` 的計算式（只到 7–10），以及全庫搜 `sub_10241` 的位址有沒有被當立即值取走（**沒… | 靜態 |
| [`re/58-bgm-scene-mapping.md`](../re/58-bgm-scene-mapping.md) | `AX=09F2h` | 換曲前送的服務號，TSR 那一側還沒讀（`17` §7） | 靜態 |
| [`re/58-bgm-scene-mapping.md`](../re/58-bgm-scene-mapping.md) | `AL` 的 6 vs 5 | `sub_10241` 對曲號 ≥ 2 把 `AL` 從 6 改成 5，語意未解 | 靜態 |
| [`re/58-bgm-scene-mapping.md`](../re/58-bgm-scene-mapping.md) | 音色聽感 | 這一份只解「哪一首」。**渲染出來像不像原版是另一回事**（`../playtest/26` §5） | 靜態 |
| [`re/59-game-over-exit-codes.md`](../re/59-game-over-exit-codes.md) | 結局的兩則訊息 | TALK `#75` 與 `#407` 的內容還沒對出來 | 靜態 |
| [`re/59-game-over-exit-codes.md`](../re/59-game-over-exit-codes.md) | `sub_14DF0` 的 CF | 「找不到替代據點」與「據點數 0」是不是同一件事，還沒逐行讀 | 靜態 |
| [`re/59-game-over-exit-codes.md`](../re/59-game-over-exit-codes.md) | 無主城 `0x18` | 值 24 落在 22 個勢力之外，但劇本裡有沒有無主城沒查過 | 靜態 |
| [`re/59-game-over-exit-codes.md`](../re/59-game-over-exit-codes.md) | `D7END.EXE` | 結局過場本身完全沒讀（`END_S*.DAT` 的用法未解） | 靜態 |
| [`re/60-tactical-sidebar.md`](../re/60-tactical-sidebar.md) | `▶▶` 列切換的語意 | 機制 confirmed（`byte_1A06A` 在 `0xEB`／`0x74` 間切、視點回 (128,128)），但 `loc_1A065` 那段自我修改碼還沒逐行讀，所以**擋掉的是什麼未解** | 靜態 |
| [`re/60-tactical-sidebar.md`](../re/60-tactical-sidebar.md) | 城兵臨時軍團的主將名 | §4.1：`0x4200` 照索引算式會指到武將表全零的那一筆。`sub_14F58`（`cx=0x1B`／`0x1C`）還沒讀 | 靜態 |
| [`re/60-tactical-sidebar.md`](../re/60-tactical-sidebar.md) | 段 1 `0x0000`／`0x0800`／`0x1000`／`0x1800`／`0x3500` 的圖形內容 | 貼點與尺寸 confirmed，**圖上畫了什麼**要另外解碼（`../formats/03` §5.3 的 UI 語意缺口） | 靜態 |
| [`re/60-tactical-sidebar.md`](../re/60-tactical-sidebar.md) | 熱區 `0x01`／`0x1F` | 兩張表裡都有 handler，但沒找到註冊它們的 `sub_1E3D7` 呼叫點 | 靜態 |
| [`re/60-tactical-sidebar.md`](../re/60-tactical-sidebar.md) | `0x1C21A`（退却）的 `sub_1A8F6` | 只知道它回 CF 與 `ah`，內容未讀 | 靜態 |
| [`re/60-tactical-sidebar.md`](../re/60-tactical-sidebar.md) | 側欄美術的調色盤 | 本份記的都是**調色盤索引**，不是 RGB。要比顏色得用 `GAMEPAL.BRG` 的當季 bank | 靜態 |
| [`re/61-timer-tick-source.md`](../re/61-timer-tick-source.md) | `ds:98A5h` 那個延時器實際在等什麼 | §5。呼叫端 `sub_11BE0` 沒讀 | 靜態 |
| [`re/61-timer-tick-source.md`](../re/61-timer-tick-source.md) | 音樂 tempo 分頻器 `cs:0B68h` 的算式 | `0x859` 那 20 條指令：`al = ((0FFh − ah) × 13) >> 3`，`ah` 從哪來沒讀 | 靜態 |
| [`re/61-timer-tick-source.md`](../re/61-timer-tick-source.md) | `cs:099Eh` 的 bit 1 | 「音樂啟用」是從用法推的，寫入端沒讀 | 靜態 |
| [`re/61-timer-tick-source.md`](../re/61-timer-tick-source.md) | 無音效驅動時的行為 | §3 推論「會卡死」，**沒有實測**——DOSBox 拿掉 `YNSOUND.COM` 跑一次就能驗 | 實測 |
| [`re/62-strategy-minimap.md`](../re/62-strategy-minimap.md) | 視野框的美術與確切尺寸 | `word_10D4C` offset 0，`sub_19752` 沒讀。20×12 是從畫面格數推的 | 靜態 |
| [`re/62-strategy-minimap.md`](../re/62-strategy-minimap.md) | `byte_198A7` 的初值 | 靜態影像裡是 `0`。**開新遊戲時有沒有被寫過沒查** | 靜態 |
| [`re/62-strategy-minimap.md`](../re/62-strategy-minimap.md) | 熱區 `0x16`（點地圖）做什麼 | 沒讀。合理猜測是把大地圖捲到該處，但**沒驗** | 靜態 |
| [`re/62-strategy-minimap.md`](../re/62-strategy-minimap.md) | `sub_15A3A` 裡 `si += 200h`／`bx += 0A00h` | 算完沒用到，是遺留的死碼 | 靜態 |
| [`re/62-strategy-minimap.md`](../re/62-strategy-minimap.md) | 圖例底圖在哪個資源 | `sub_1FA37` 的 `ds` 來自 `word_10D50`；`docs/re/47` 記成段 3 `0x09A0`，兩者沒對過 | 靜態 |
| [`re/63-ground-plane-map.md`](../re/63-ground-plane-map.md) | 段變數的配置迴圈 | `word_1D2F6`–`word_1D30E` 沒有直接寫入的 xref，§1 的相鄰關係是推論 | 靜態 |
| [`re/63-ground-plane-map.md`](../re/63-ground-plane-map.md) | `sub_1B186`／`sub_1B15D` | 爬升／下降時檢查上下一層的那兩支，沒讀 | 靜態 |
| [`re/63-ground-plane-map.md`](../re/63-ground-plane-map.md) | 命令 6 為什麼擋高平面橫移 | `[si+1Ah] == 6`，命令碼 6 是什麼沒對過 | 靜態 |
| [`re/63-ground-plane-map.md`](../re/63-ground-plane-map.md) | 打壞城壁之後地面層表不更新 | `sub_1B824` 只重算通行層（`sub_1BB6D`）與佔用表，**沒有重跑 `sub_1BC39`**。而那不影響結果：城壁的地面層本來就是拿打壞後的圖塊算的（§2 的 +0x10） | 靜態 |
| [`re/64-corps-arrival-state-machine.md`](../re/64-corps-arrival-state-machine.md) | `+0x00` 位元 1 | Stage 10／11 改目標時 `or byte [si], 2`；`sub_12662` 在 `0x126A0` 讀它、清掉並呼叫 `sub_147BB`，而 **`sub_147BB` 未讀**。疑似「路線要重算」 | 靜態 |
| [`re/64-corps-arrival-state-machine.md`](../re/64-corps-arrival-state-machine.md) | 已經解掉的**：非玩家的 Stage 0–3 四支 handler 與 `sub_13E11 | （未解小節內文） | 靜態 |
| [`re/65-ai-march-decision-chain.md`](../re/65-ai-march-decision-chain.md) | 非玩家的四個 entry 掛在 §6 未解。 | （散句） | 靜態 |
| [`re/65-ai-march-decision-chain.md`](../re/65-ai-march-decision-chain.md) | `+0x00` 位元 1（`or byte [si], 2`） | Stage 2 改目標、Stage 10／11 校正目標、`sub_1474A` 都會設；`sub_12662` 在 `0x126A0` 讀它並清掉，接著呼叫 `sub_147BB`。**`sub_147BB` 未讀**，位元 1 疑似「路線要重算」但沒有確認 | 靜態 |
| [`re/65-ai-march-decision-chain.md`](../re/65-ai-march-decision-chain.md) | `sub_1487B` 的挑格邏輯 | `sub_1474A`（AI 編成後的第一個目標）與 `sub_14DA4` 用它挑相鄰格；讀了外框但沒逐條解 | 靜態 |
| [`re/65-ai-march-decision-chain.md`](../re/65-ai-march-decision-chain.md) | `sub_128F4` 的 STC 分支 | 走到敵方據點時呼叫 `sub_1291A`（俘虜／脫離判定），之後 `di` 不可信。本文件的 `di` 推論只涵蓋一般路徑 | 靜態 |

## 2.4 驗收（54 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`playtest/10-event-message-modal.md`](../playtest/10-event-message-modal.md) | 直接證據為 `KI.EXE.asm` 的 IDA 線性位址 `0001237E`、`000134A6`、`000134B1`、 | （未解小節內文） | 靜態 |
| [`playtest/14-m7-review.md`](../playtest/14-m7-review.md) | #195 | `3` / 4 / 10 / 1 / `嗯，現在的話，{3}\n也不至於拒絕。好，准\n許停戰！\n` / `sub_13C99` 的狀態值 1 直接取 #193–#195；此格需保留 `{3}`，重用目前 #193 的既有繁中譯文，避免另造未驗證句子。 | 靜態 |
| [`playtest/17-expert-dosbox-remake.md`](../playtest/17-expert-dosbox-remake.md) | 松崗 DOS/V 原版 | **PASS（啟動至開場）** / 2026-08-12 證實空白確認／`0000`／`1234` 均越過密碼頁；完整自然長程驗證尚未執行 | 靜態 |
| [`playtest/21-dosboxx-bridge-sampling.md`](../playtest/21-dosboxx-bridge-sampling.md) | 據點換手之後遮罩會不會跟著變 | `sub_1890A` 的行為，靜態讀得出來，動態沒驗——要打下一座城才看得到 | 靜態 |
| [`playtest/21-dosboxx-bridge-sampling.md`](../playtest/21-dosboxx-bridge-sampling.md) | 松崗 DOS/V 側 | 這套 bridge 還沒在 DOS/V 上跑過。**密碼頁不構成阻礙**（四格留白按「確定」即可通過，`18`）——是還沒做 | 靜態 |
| [`playtest/21-dosboxx-bridge-sampling.md`](../playtest/21-dosboxx-bridge-sampling.md) | 上游授權 | `DOSBox-X-MCP-Debugger` 的原創碼**尚未選定授權條款**（README 明講是刻意留白）。本專案只在本機使用，未再散布 | 實測 |
| [`playtest/23-main-screen-geometry.md`](../playtest/23-main-screen-geometry.md) | 原版「四窗全開」的截圖 | **主畫面在模擬器上收不到任何點擊**，不只是開關——鑑別測試見下 | 實測 |
| [`playtest/23-main-screen-geometry.md`](../playtest/23-main-screen-geometry.md) | 松崗 DOS/V 側的主畫面 | 開新遊戲流程停在「確定」按鈕不回應。座標已照 PC-98 換算（Y ＋40）重試，仍不動；密碼頁不是障礙（`18`） | 靜態 |
| [`playtest/23-main-screen-geometry.md`](../playtest/23-main-screen-geometry.md) | 同版本同調色盤的對拍 | 上面兩項任一個通了就能做。**在那之前 `banner` 的 49.6% 不代表 remake 有錯** | 靜態 |
| [`playtest/23-main-screen-geometry.md`](../playtest/23-main-screen-geometry.md) | 四張 24×16 圖形的內容 | 顯示清單指到圖庫位移 `0x1200`／`0x12C0`／`0x1380`／`0x1440`，還沒畫出來看；remake 先用同位置同尺寸的色塊佔位 | 實測 |
| [`playtest/23-main-screen-geometry.md`](../playtest/23-main-screen-geometry.md) | 其他 10 個顯示清單場景 | 同一支直譯器有 11 個呼叫端，只讀了場景 0（`../re/48` §5） | 靜態 |
| [`playtest/24-window-toggles.md`](../playtest/24-window-toggles.md) | 各視窗內部的像素 | 只對過邊線位置，沒有對過內容 | 靜態 |
| [`playtest/24-window-toggles.md`](../playtest/24-window-toggles.md) | 原版執行期的開關行為 | **沒驗過**：模擬器上主畫面收不到點擊（`23` §4.1） | 靜態 |
| [`playtest/24-window-toggles.md`](../playtest/24-window-toggles.md) | 系統視窗的四個項目 | 存檔／畫面模式／音源／戰略速度，不在這一輪範圍 | 靜態 |
| [`playtest/25-audio-capture-feasibility.md`](../playtest/25-audio-capture-feasibility.md) | 這個實驗**沒有**證明下面這些事。 | （未解小節內文） | 靜態 |
| [`playtest/25-audio-capture-feasibility.md`](../playtest/25-audio-capture-feasibility.md) | **逐曲觸發** | 只錄到開場動畫（`D7OPEN.EXE` 自己會播）。⚠ 這一項**不再擋住任何事**——音檔改由 `tools/bgm2ogg.sh` 離線渲染，不需要在模擬器裡逐首觸發。要當對照組才需要它 | 靜態 |
| [`playtest/25-audio-capture-feasibility.md`](../playtest/25-audio-capture-feasibility.md) | **音效** | 戰術的三個 effect code 已知（`re/17` §3），但沒錄過 | 靜態 |
| [`playtest/25-audio-capture-feasibility.md`](../playtest/25-audio-capture-feasibility.md) | **音源正確性** | DOSBox 用 `sbtype=sb16`／`oplmode=auto` 模擬，與真實硬體的音色差異沒有對照組 | 實測 |
| [`playtest/26-bgm-render-vs-recording.md`](../playtest/26-bgm-render-vs-recording.md) | 音色的聽感 | 頻譜只驗了基頻。諧波結構（也就是「像不像那個音色」）沒有量化比對 | 靜態 |
| [`playtest/26-bgm-render-vs-recording.md`](../playtest/26-bgm-render-vs-recording.md) | 相關係數為什麼不是 0.9 | DOSBox 的 OPL 模擬與這顆的包絡實作不同，加上錄音有系統噪訊。**沒有排除「還有小錯」的可能** | 實測 |
| [`playtest/26-bgm-render-vs-recording.md`](../playtest/26-bgm-render-vs-recording.md) | 其他曲子 | 只有開場曲有錄音對照組。另外 13 首沒有 | 靜態 |
| [`playtest/27-original-video-frame-parity.md`](../playtest/27-original-video-frame-parity.md) | 逐像素 parity | 影片是再編碼的，做不到。要真的逐像素得回到模擬器，而主畫面的點擊閘還在 | 靜態 |
| [`playtest/27-original-video-frame-parity.md`](../playtest/27-original-video-frame-parity.md) | 色彩 | 只比了幾何。調色盤要另外用「同一格地形的色號」比，不能用影片的 RGB | 靜態 |
| [`playtest/27-original-video-frame-parity.md`](../playtest/27-original-video-frame-parity.md) | 戰術側欄的逐格對拍 | 組成已對齊（§7.3、`../re/60`），**同一場戰鬥的逐格比對還沒做** | 靜態 |
| [`playtest/27-original-video-frame-parity.md`](../playtest/27-original-video-frame-parity.md) | 一覽表視窗 | 影片裡有武將／據點／財政的實錄，**還沒量** | 靜態 |
| [`playtest/28-siege-breach-measurement.md`](../playtest/28-siege-breach-measurement.md) | 水平跨格的碰撞判定全在 `sub_1B1B1`（`0001B1B1`，143 B）： | （未解小節內文） | 靜態 |
| [`playtest/28-siege-breach-measurement.md`](../playtest/28-siege-breach-measurement.md) | 1. **擋路的是「單位」，不是地形高度。 | （未解小節內文） | 靜態 |
| [`playtest/28-siege-breach-measurement.md`](../playtest/28-siege-breach-measurement.md) | 2. **低平面的水平跨格沒有高度差上限**——目標格地面比自己高就升一層、 | （未解小節內文） | 靜態 |
| [`playtest/28-siege-breach-measurement.md`](../playtest/28-siege-breach-measurement.md) | 3. **高平面要高度完全相等**，而 `[si+1Eh]`（`PlaneHigh`）決定讀哪一張。 | （未解小節內文） | 靜態 |
| [`playtest/28-siege-breach-measurement.md`](../playtest/28-siege-breach-measurement.md) | `Field.gateX` 與登城點 | gateX 在三張抽樣圖上**正好就是可破壞城壁那一格**。`doScaleWall` 已經把它當目標，但爬不上去 | 靜態 |
| [`playtest/28-siege-breach-measurement.md`](../playtest/28-siege-breach-measurement.md) | 繞路點清單的演算法 | `0x1800 + 兵編號 × 128`，`loc_1BD46` 算（`../re/11` §5.15）。解出來之後 `FindPathForcing` 那條近似要換掉 | 靜態 |
| [`playtest/28-siege-breach-measurement.md`](../playtest/28-siege-breach-measurement.md) | 攻城計時器與突破時間的關係 | 大將體力 100 起、每 10 幀掉 1、50 觸發退卻 ⇒ **攻方只有約 500 幀**。爬牆接上之後要重量一次，確認這個預算合理 | 靜態 |
| [`playtest/29-strategy-minimap-markers.md`](../playtest/29-strategy-minimap-markers.md) | 視野框的美術 | 原版在 `word_10D4C`，尺寸沒從程式碼讀到 | 靜態 |
| [`playtest/29-strategy-minimap-markers.md`](../playtest/29-strategy-minimap-markers.md) | 點地圖區（熱區 `0x16`） | 原版做什麼沒讀 | 靜態 |
| [`playtest/30-ground-planes-implemented.md`](../playtest/30-ground-planes-implemented.md) | **攻方大多數不前進** | §3。目標選擇的問題，與地形無關 | 靜態 |
| [`playtest/30-ground-planes-implemented.md`](../playtest/30-ground-planes-implemented.md) | 一幀能有幾個兵撞牆 | 原版沒量過。前排寬度決定破牆速度，而破牆速度決定攻城打不打得下來 | 靜態 |
| [`playtest/30-ground-planes-implemented.md`](../playtest/30-ground-planes-implemented.md) | 打壞城壁之後地面層表不更新 | 原版就不更新，而且不影響結果——城壁的地面層本來就是拿打壞後的圖塊算的（`../re/63` §2） | 靜態 |
| [`playtest/30-ground-planes-implemented.md`](../playtest/30-ground-planes-implemented.md) | 高平面的橫向移動沒有實測 | 守方站到牆頂的情境還沒跑過 | 實測 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | 主畫面：視窗內部底紋 | 原版視窗內是**黑底 ＋ 深藍龍紋**（兩色、寬 32 的直條水平重複，`../formats/03` §5.5），remake 是純深藍 / 看得出來 / 靜態找法走完（69 檔、四段、8×8 拼塊都掃過）；下一步是動態追 VRAM 寫入 | 靜態 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | 主畫面：大地圖地形色調 | 原版偏黃綠、remake 偏綠 / 存疑 / 可能是影片的色彩取樣。要驗就比**同一格的色號**，不要比 RGB | 靜態 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | **戰鬥指揮／委任選單** | 版面是 remake 自訂的（`sub_193E9` 的矩形沒解）；**選項字串已改成原版的 TALK #76**，行軍指示的三選一也接上了（`../spec/39`） / 部分 / 缺 `sub_193E9` 的版面 | 靜態 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | 戰場：同一場的逐格對拍 | 沒做過 / 未對過 / 需要同狀態，難度同主畫面 | 靜態 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | 攻方大多數不前進 | 六個指令都能指揮部隊動作（`../spec/37` §4）；**AI 自己撞不進城是設計**——說明書第 11 章整章在講破城要換陣形 / 玩家要自己操作 | 靜態 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | `sub_13E11` 每「時」做什麼 | 未讀 / 行軍與 AI 的節拍 | 靜態 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | 防災值怎麼成長 | 欄位已知、規則未解 / 天災的長期經濟 | 靜態 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | 行軍費用 | 說明書 10.6 有專節，沒讀 / 經濟 | 靜態 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | 預備兵維持費單價 | `sub_15358` 尾段沒讀 / 經濟 | 靜態 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | 士氣值存在哪 | 戰術層 / 戰鬥判定 | 靜態 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | 地形色調差 | 分不出是影片取樣還是調色盤，要比色號 | 靜態 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | 戰鬥指揮／委任選單 | 影片裡沒有對照影格 | 靜態 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | 同一場戰鬥的逐格對拍 | 需要同狀態，還沒做 | 靜態 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | 勢力一覽的欄位 | 沒有逐欄比對 | 靜態 |
| [`playtest/32-talk-layout-fit.md`](../playtest/32-talk-layout-fit.md) | 肖像框的寬度 | remake 用 160 px，出處是 `sub_18810`／`sub_1895D` 的常數；**原版實錄影格上那個框的文字區看起來更寬**（f008 量到約 275 px）。要嘛常數讀錯、要嘛影格上那個是另一種框 | 靜態 |
| [`playtest/32-talk-layout-fit.md`](../playtest/32-talk-layout-fit.md) | 變數的實際長度分布 | 這一輪用固定三全形替身。人名多半是 2–3 全形、地名 2–3，但**軍團名與勢力名沒有逐一量過** | 靜態 |

## 2.5 外部資料（17 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`reference/02-jp-cht-diff.md`](../reference/02-jp-cht-diff.md) | #360–#1021 的逐句文意審查 | 第二批已列出，逐句審查未完成（§12） | 靜態 |
| [`reference/02-jp-cht-diff.md`](../reference/02-jp-cht-diff.md) | 校訂後的畫面抽樣與排版 parity | 未做。M7 因此未封口 | 靜態 |
| [`reference/02-jp-cht-diff.md`](../reference/02-jp-cht-diff.md) | `#223` 等訊息的欄位完整語意 | 只修已證實的標記編號，欄位語意仍未解（§9） | 靜態 |
| [`reference/03-baked-japanese.md`](../reference/03-baked-japanese.md) | 松崗版是移植不是重寫——**23 個檔在兩版之間 byte-for-byte 完全相同 | （未解小節內文） | 兩版對照 |
| [`reference/03-baked-japanese.md`](../reference/03-baked-japanese.md) | `END_S3`–`END_S11` | 九張全不同 / 有燒字，**松崗重繪過** | 靜態 |
| [`reference/03-baked-japanese.md`](../reference/03-baked-japanese.md) | `END_S1`／`S2`／`S12` | 相同 / 沒有燒字 | 靜態 |
| [`reference/03-baked-japanese.md`](../reference/03-baked-japanese.md) | `OPEN_S1`–`S6` | 六張全相同 / 開場沒有燒字，文字是疊繪的 | 靜態 |
| [`reference/03-baked-japanese.md`](../reference/03-baked-japanese.md) | **`ICONGRF.DAT`** | **相同** / **沒重繪 → 裡面的日文留著** | 靜態 |
| [`reference/03-baked-japanese.md`](../reference/03-baked-japanese.md) | `KAOGRF`／`KYOGRF`／`IVENTGRF` | 相同 / 純圖像，目前看過的都沒有文字 | 靜態 |
| [`reference/03-baked-japanese.md`](../reference/03-baked-japanese.md) | !標題橫幅 | （未解小節內文） | 靜態 |
| [`reference/04-first-survey.md`](../reference/04-first-survey.md) | `YNSOUND.COM` | 3,463 B / 音效驅動，**假說**：常駐 TSR | 靜態 |
| [`reference/04-first-survey.md`](../reference/04-first-survey.md) | `SHOW.O` | 57,148 B / 被 `INSTALL.EXE` 與 `LOGO.EXE` 引用。開頭 `3c df 00 00 11 af 01 00 50 00 80 07`。**未解** | 靜態 |
| [`reference/04-first-survey.md`](../reference/04-first-survey.md) | LOGO.EXE`／`YNFONT.EXE`／`INSTALL.EXE` 都有 | （未解小節內文） | 靜態 |
| [`reference/04-first-survey.md`](../reference/04-first-survey.md) | `KAOGRF.DAT` | 307,200 / 頭像（日文「顔」）。**假說**：307,200 ÷ 2,048 = **150 張 64×64 4bpp**，而武將 146 人 + 4 → 對得起來，但**這是算術巧合等級的證據，要驗** | 靜態 |
| [`reference/04-first-survey.md`](../reference/04-first-survey.md) | `KYOGRF.DAT` | 69,120 / 未解 | 靜態 |
| [`reference/04-first-survey.md`](../reference/04-first-survey.md) | FM 3 聲 ＋ SSG 3 聲，埠 `0x188`／`0x18A`。 DOS/V 側未解。 | （散句） | 靜態 |
| [`reference/05-eten-font-provenance.md`](../reference/05-eten-font-provenance.md) | `END_S13/S14/S15` 是中文版加的結局段 | S13／S14 是字型。**`END_S15` 仍未解** | 靜態 |

## 2.6 其他（91 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`promo/dosv-adlib-and-tactical-review.md`](../promo/dosv-adlib-and-tactical-review.md) | 原版雙 TALK 的 payload、肖像、文字 baseline 與出現時序尚未在同一戰況對拍。 | （未解小節內文） | 靜態 |
| [`promo/dosv-adlib-and-tactical-review.md`](../promo/dosv-adlib-and-tactical-review.md) | 原版右欄完整狀態資訊、旗標、軍名、命令 glyph 與裝飾仍未逐區等價。 | （未解小節內文） | 靜態 |
| [`promo/dosv-adlib-and-tactical-review.md`](../promo/dosv-adlib-and-tactical-review.md) | 原版底列按鈕 glyph、選取狀態與 remake 文字／簡化圖示仍有差異。 | （未解小節內文） | 靜態 |
| [`promo/dosv-adlib-and-tactical-review.md`](../promo/dosv-adlib-and-tactical-review.md) | 地形、部隊編成、鏡頭中心、動畫 frame 與戰況不同，不能用目前推廣片判定物件 | （未解小節內文） | 靜態 |
| [`promo/yt-remake-pixel-review.md`](../promo/yt-remake-pixel-review.md) | 中央 raw reserve glyph | 未解出原版圖形。remake 不冒充，改用自繪 | 靜態 |
| [`release/01-cross-build-gate.md`](../release/01-cross-build-gate.md) | 目標 OS 實跑 | **做不到**：這台是 Linux，沒有 Mac／Windows。檔頭驗過（PE32+／Mach-O），但視窗、輸入、音訊、字型載入都沒有在目標系統上跑過 | 實測 |
| [`release/01-cross-build-gate.md`](../release/01-cross-build-gate.md) | linux/arm64 的本體 | 要在 arm64 的 Linux 上建（Ebiten 的 cgo 沒有交叉工具鏈） | 靜態 |
| [`release/01-cross-build-gate.md`](../release/01-cross-build-gate.md) | Windows 的 smoke | 同第一項 | 靜態 |
| [`spec/00-index.md`](../spec/00-index.md) | **推論等級** | confirmed／強證據／假說／未知（`CLAUDE.md` §9）。假說也可以實作，但要標 | 靜態 |
| [`spec/00-index.md`](../spec/00-index.md) | 進言「請求君主出陣」（`sub_1699E`） | `11-ai-sortie.md` / 可實作，**尚未實作** | 靜態 |
| [`spec/00-index.md`](../spec/00-index.md) | 主畫面四個視窗的開關 | `13-main-window-toggles.md` / 已實作並留下截圖；原版執行期未驗 | 實測 |
| [`spec/00-index.md`](../spec/00-index.md) | 底列六格是選部隊 | `33-squad-selection.md` / 已實作並有單測；命令圖示的來源段未定案 | 靜態 |
| [`spec/00-index.md`](../spec/00-index.md) | 一覽表的欄位與版面 | `38-list-windows.md` / 四個家族照原版重做；捲軸未解 | 靜態 |
| [`spec/10-city-tick.md`](../spec/10-city-tick.md) | `sub_14194`／`sub_14269` | 內政與災害 marker 的細節在別的規格（`docs/mechanics/40`），本規格只保證呼叫順序 | 靜態 |
| [`spec/10-city-tick.md`](../spec/10-city-tick.md) | 據點換手之後 `+0x00` 低 4 位會不會跟著變 | `sub_1890A` 靜態讀過，動態沒驗——要打下一座城才看得到 | 靜態 |
| [`spec/10-city-tick.md`](../spec/10-city-tick.md) | 玩家據點求援的喇叭聲（`sub_10CDE`） | 呈現層未接 | 靜態 |
| [`spec/11-ai-sortie.md`](../spec/11-ai-sortie.md) | `資金高位 >= 0x80` 那一支 | `cmp bh, 80h / jnb` 會直接算「答應」，等於資金超過約 840 萬時門檻失效。**看起來像有號數的邊界處理**，未逐位對過 | 靜態 |
| [`spec/11-ai-sortie.md`](../spec/11-ai-sortie.md) | 君主出陣之後的行為 | 那支軍團跟一般軍團有沒有差別，未讀 | 靜態 |
| [`spec/12-strategy-chrome.md`](../spec/12-strategy-chrome.md) | 信賴度的呈現 | 原版是量條，remake 是數字。改成量條之前要先確定顏色與底圖，否則會畫出一條沒有背景的裸色塊 | 靜態 |
| [`spec/12-strategy-chrome.md`](../spec/12-strategy-chrome.md) | 勢力色標 | 原版怎麼畫未讀（`sub_15CE0` 是小地圖的四色點，不是這一列） | 靜態 |
| [`spec/12-strategy-chrome.md`](../spec/12-strategy-chrome.md) | opcode `06` 與其他 10 個場景 | 顯示清單解得開但只讀了場景 0（`docs/re/48` §5） | 靜態 |
| [`spec/12-strategy-chrome.md`](../spec/12-strategy-chrome.md) | 樣式碼的值域 | 只確定 `0`＝擦除、`0x0B`＝命令、`0x0Bh`／`0x10h`／`0x15h`／`0x1Fh` 各自出現在哪個視窗已知，完整值域未列 | 靜態 |
| [`spec/13-main-window-toggles.md`](../spec/13-main-window-toggles.md) | 對得上（`docs/playtest/24`）。 原版執行期的開關行為仍未驗。 | （散句） | 靜態 |
| [`spec/13-main-window-toggles.md`](../spec/13-main-window-toggles.md) | 熱區 5 | 原版登記了但不接任何常式，remake 照樣不做事 | 靜態 |
| [`spec/13-main-window-toggles.md`](../spec/13-main-window-toggles.md) | 六列的語意 | handler 讀到四支（`docs/re/55` §5）：畫面模式換調色盤組、音效走驅動、戰略速度只存值、戰術速度存值 ×16。**「資料儲存」與「遊戲結束」那兩支沒讀**，remake 照標籤字面接（開四槽視窗／走 ＹＥＳ／ＮＯ 確認） | 靜態 |
| [`spec/13-main-window-toggles.md`](../spec/13-main-window-toggles.md) | 「畫面模式」 | 兩個選項是「１６色」與「 液晶 」，切的是 `GAMEPAL.BRG` 的 bank 0–3 ↔ 4–7（`docs/re/02` §6.2）。**remake 只做第 0 組**，這一格固定顯示「１６色」——液晶那組是給 8 階調液晶的高飽和純色，現代螢幕沒有對照物 | 靜態 |
| [`spec/13-main-window-toggles.md`](../spec/13-main-window-toggles.md) | 「音效」 | 值由 `g.soundValue()` 填。原版五個選項是 ＯＦＦ／TYPE 1–4（音源型別），remake 只有開／關 | 靜態 |
| [`spec/13-main-window-toggles.md`](../spec/13-main-window-toggles.md) | 戰場內調速度 | 戰場獨佔輸入，所以 `updateBattle` 自己接一次 ＋／−（調戰術速度），調完浮一行 1.5 秒的提示。**原版戰場沒有速度指示**，常駐顯示會破壞版面 parity | 靜態 |
| [`spec/14-finance-window.md`](../spec/14-finance-window.md) | 收入的來源 | `cs:word_10D02` 由誰計算未讀（月結那一支是候選）。remake 的月結算得出 `res.Income` 但沒有留下來，所以這一格暫時顯示 0——**留白比填一個自己算的數字誠實** | 靜態 |
| [`spec/14-finance-window.md`](../spec/14-finance-window.md) | 徵兵數的上限 | 只有稅率的 100 是機器碼常數；三個兵種的上限沒看到 | 靜態 |
| [`spec/14-finance-window.md`](../spec/14-finance-window.md) | 數值輸入器 | `sub_17C6E` 已有 RE（`docs/re/13`），但 remake 的財政還沒接上去 | 靜態 |
| [`spec/20-save-format.md`](../spec/20-save-format.md) | 存檔區塊的 7 KB 未解區 | `+0x1EC0`–`+0x42C0`，靠 `raw` 原樣保存，但**內容仍不知道**（`docs/formats/08`） | 靜態 |
| [`spec/20-save-format.md`](../spec/20-save-format.md) | 原版 `SAVE.DAT` 的槽位語意 | 四個槽與 `SINARIO.DAT` 的四個劇本是不是同一個編號空間，未確認 | 靜態 |
| [`spec/21-corps-formation-reserves.md`](../spec/21-corps-formation-reserves.md) | 編成畫面的兵種切換 | remake 由呼叫端直接給 `kinds`，沒有原版那個「點一下 +1 → 全退回池 → 重跑分配」的迴圈（`sub_16C92`）。這是 UI 層的差異，不影響分配式 | 靜態 |
| [`spec/21-corps-formation-reserves.md`](../spec/21-corps-formation-reserves.md) | 退兵回池 | `sub_14717` 已讀（一點對一點、上限 65,500），remake 還沒有「解散軍團把兵退回去」的路徑 | 靜態 |
| [`spec/21-corps-formation-reserves.md`](../spec/21-corps-formation-reserves.md) | 池的上限 | `sub_155EC` 的 `0xFFDC` 只在退兵路徑上驗過；月結加兵是不是同一支未查 | 靜態 |
| [`spec/22-corps-formation-window.md`](../spec/22-corps-formation-window.md) | 頭像的邊框 | `sub_107D2` 只 blit 64×64 的圖塊，**框在哪裡畫的沒找到**——場景 5 的 op 清單裡沒有頭像那一格的框 | 靜態 |
| [`spec/22-corps-formation-window.md`](../spec/22-corps-formation-window.md) | 兵種標籤 | 畫面用場景 5 的「主將」，規則層的 `army.Position` 第一個是「大將」（原版 TALK #62 也這樣說）。兩處用語不同是原版就有的，不要統一 | 靜態 |
| [`spec/23-city-info-window.md`](../spec/23-city-info-window.md) | 進入方式 | 原版由地圖上點據點進來（`sub_11E46`），remake 走一覽表 | 靜態 |
| [`spec/23-city-info-window.md`](../spec/23-city-info-window.md) | `cs:word_1987C` | 原版每次開視窗都重讀一次檔；remake 的 `library` 是整檔載入，不需要這一層 | 靜態 |
| [`spec/24-corps-info-window.md`](../spec/24-corps-info-window.md) | 指令流程 | `sub_17FDB` 未讀（`docs/re/51` §5）。remake 的行軍指令走自己的流程（`M`） | 靜態 |
| [`spec/24-corps-info-window.md`](../spec/24-corps-info-window.md) | 進入方式 | 原版也可以在地圖上直接點軍團（`sub_11E46`），remake 只有一覽表 | 靜態 |
| [`spec/25-slot-select-window.md`](../spec/25-slot-select-window.md) | 空槽標記 | 原版用名稱欄第一個字 `0xD0A1`；remake 用「載得起來且玩家勢力有效」判定，兩者不等價 | 靜態 |
| [`spec/25-slot-select-window.md`](../spec/25-slot-select-window.md) | 新遊戲共用 | remake 的啟動殼層是自己的畫面，還沒有換成這個四槽視窗 | 靜態 |
| [`spec/26-yes-no-dialog.md`](../spec/26-yes-no-dialog.md) | 原版的使用者 | `sub_18DC8` 只有一個呼叫端 `sub_11AC3`（新遊戲流程），問題文字由那裡給，內容未讀 | 靜態 |
| [`spec/26-yes-no-dialog.md`](../spec/26-yes-no-dialog.md) | 背景保存 | `sub_19796`／`sub_197C3(cx=600Dh)` 是開關前後成對的一支，推測是保存／還原被蓋住的畫面，未讀 | 靜態 |
| [`spec/27-lord-select-window.md`](../spec/27-lord-select-window.md) | 「自定」 | 軍師命名（場景 9 的注音輸入）還沒做，這顆按鈕目前無效 | 靜態 |
| [`spec/27-lord-select-window.md`](../spec/27-lord-select-window.md) | 換勢力 | 原版的換法在 `sub_11AC3`，未讀 | 靜態 |
| [`spec/27-lord-select-window.md`](../spec/27-lord-select-window.md) | 頭像尺寸 | 軍師頭像的下緣照原版座標會略微超出那個 208×104 的底框；沒有 oracle 可比，先照機器碼畫 | 實測 |
| [`spec/28-scenario-json.md`](../spec/28-scenario-json.md) | 事件佇列 | 這一輪不進 JSON。編輯器要動它得先有 UI 語意 | 靜態 |
| [`spec/28-scenario-json.md`](../spec/28-scenario-json.md) | 未解區域 | `+0x1EC0` 那 7 KB 仍是黑盒，只能靠改寫保留 | 靜態 |
| [`spec/28-scenario-json.md`](../spec/28-scenario-json.md) | 編輯器 | 這一份只做資料層。UI 是另一份規格 | 靜態 |
| [`spec/29-audio.md`](../spec/29-audio.md) | 曲 1 | DOS/V 的 `KI.EXE` 裡沒有任何呼叫端（`re/58` §5）。PC-98 版還沒掃 | 靜態 |
| [`spec/29-audio.md`](../spec/29-audio.md) | 曲 6 的接法 | 原版是四支對話／事件常式，remake 用「事件訊息開著」與「進言對話開著」兩個狀態代替，**不是一對一** | 靜態 |
| [`spec/29-audio.md`](../spec/29-audio.md) | 換季的兩段時序 | 原版第 1 天停、第 2 天換曲，調色盤另外漸變 16 天。remake 只做了換曲那一半 | 靜態 |
| [`spec/29-audio.md`](../spec/29-audio.md) | 迴圈點怎麼呈現 | 原版靠控制事件 `C1`／`C3` 無限循環；ogg 是有限長度，要決定渲染幾輪或另存迴圈點 | 靜態 |
| [`spec/29-audio.md`](../spec/29-audio.md) | 全域音量偏移 | `cs:0996h` 誰設、範圍多少未解（`re/57` §8） | 靜態 |
| [`spec/29-audio.md`](../spec/29-audio.md) | PC-98 版 | 音源是 YM2203，暫存器路徑完全沒讀。要不要做是待裁定的問題 | 靜態 |
| [`spec/30-victory.md`](../spec/30-victory.md) | 結局畫面 | `D7END.EXE` 是另一支程式，`END_S*.DAT` 的過場沒有實作 | 靜態 |
| [`spec/30-victory.md`](../spec/30-victory.md) | 四個劇本的結局是否不同 | `END_S*` 檔的對應未解（`mechanics/80-victory.md` §116） | 靜態 |
| [`spec/30-victory.md`](../spec/30-victory.md) | 君主陣亡時軍師怎麼辦 | 未知（同上） | 靜態 |
| [`spec/30-victory.md`](../spec/30-victory.md) | 結局的訊息 | `sub_11CD0` 送 TALK `#0x4B` 與 `#0x197`，內容還沒對出來 | 靜態 |
| [`spec/31-tactical-sidebar.md`](../spec/31-tactical-sidebar.md) | 差異 | `▶▶` 列只畫美術，**不接行為**（原版切換的是 `loc_1A065` 的自我修改碼，語意未解）；兩面將旗的熱區原版是 `retn`，remake 同樣不接 | 靜態 |
| [`spec/31-tactical-sidebar.md`](../spec/31-tactical-sidebar.md) | `▶▶` 列的行為 | `byte_1A06A` 在 `0xEB`／`0x74` 間切，`loc_1A065` 未逐行讀（`../re/60` §12） | 靜態 |
| [`spec/31-tactical-sidebar.md`](../spec/31-tactical-sidebar.md) | 段 1 五塊美術的圖形語意 | 貼點與尺寸 confirmed，圖上畫什麼要另外解（`../formats/03` §5.3） | 靜態 |
| [`spec/31-tactical-sidebar.md`](../spec/31-tactical-sidebar.md) | 城兵臨時軍團的主將名 | `0x4200` 的索引算式指到武將表全零那一筆（`../re/60` §4.1） | 靜態 |
| [`spec/32-gate-strength-bar.md`](../spec/32-gate-strength-bar.md) | 右鍵提前收掉 | 原版是熱區 `0x1D` 的右鍵 handler（`../re/60` §10）；remake 沒有戰場區的右鍵熱區層 | 靜態 |
| [`spec/32-gate-strength-bar.md`](../spec/32-gate-strength-bar.md) | 「門強度」這三個字對城壁也照用 | 原版不分城壁與門，都用同一個標籤。照抄 | 靜態 |
| [`spec/32-gate-strength-bar.md`](../spec/32-gate-strength-bar.md) | ⚠ **攻方只在第 24–44 幀打城壁，之後就不再打** | 上面那次量測的副產品。城壁耐久 500，撞一次掉 1，照理應該持續撞；**這像是規則層的缺口**（攻方接觸城壁後停止攻擊），但本輪沒有查，也不當結論 | 靜態 |
| [`spec/33-squad-selection.md`](../spec/33-squad-selection.md) | 待機兵條的欄位語意 | `word_1D30A:+0x09 + 4k` 在 `../re/11` §3.9 記成「第 k 隊的待機兵數」；條的上限 76 遠小於一隊 100 兵，所以開局會頂在上限 | 靜態 |
| [`spec/34-speed-steps.md`](../spec/34-speed-steps.md) | 最高速在原版實機是多少 | 機器相依。DOSBox 固定 cycles 量得到「那台的上限」，量不到「原版的答案」 | 實測 |
| [`spec/34-speed-steps.md`](../spec/34-speed-steps.md) | 戰場幀是否等於 remake 的一次 `Step()` | 原版一幀做完整條戰場迴圈；remake 的 `Step()` 是規則層一步。**兩者對齊過但沒逐項比** | 靜態 |
| [`spec/34-speed-steps.md`](../spec/34-speed-steps.md) | 音效驅動不在時的行為 | `../re/61` §6 | 靜態 |
| [`spec/35-strategy-minimap.md`](../spec/35-strategy-minimap.md) | 視野框的美術 | 原版是 `word_10D4C` 的圖，尺寸沒從程式碼讀到（`../re/62` §5） | 靜態 |
| [`spec/35-strategy-minimap.md`](../spec/35-strategy-minimap.md) | 點地圖區（熱區 `0x16`） | 原版做什麼沒讀 | 靜態 |
| [`spec/36-ground-planes-and-climbing.md`](../spec/36-ground-planes-and-climbing.md) | 攻方大多數不前進 | `../playtest/30` §3。目標選擇的問題，與地形無關 | 靜態 |
| [`spec/36-ground-planes-and-climbing.md`](../spec/36-ground-planes-and-climbing.md) | 打壞城壁之後不重算地面表 | 原版不重算，而且**不影響結果**：城壁的地面層本來就是拿打壞後的圖塊算的。remake 為了合成戰場仍會重算，在真實資料上是恆等變換 | 靜態 |
| [`spec/36-ground-planes-and-climbing.md`](../spec/36-ground-planes-and-climbing.md) | 命令 6 為什麼擋高平面橫移 | 命令碼 6 是什麼沒對過 | 靜態 |
| [`spec/36-ground-planes-and-climbing.md`](../spec/36-ground-planes-and-climbing.md) | `sub_1B186`／`sub_1B15D` | 爬升／下降時檢查上下一層的那兩支沒讀，remake 用「目標平面有地面」代替 | 靜態 |
| [`spec/37-tactical-player-controls.md`](../spec/37-tactical-player-controls.md) | 選了陣形之後原版有沒有立刻重排 | 機器碼只寫偏移，**沒有看到立刻移動的呼叫**；remake 照抄（等命令） | 靜態 |
| [`spec/37-tactical-player-controls.md`](../spec/37-tactical-player-controls.md) | 陣形線在小地圖上的線寬與端點 | `sub_1C5AE` 沒逐行讀，remake 畫整條 1 px 的線 | 靜態 |
| [`spec/37-tactical-player-controls.md`](../spec/37-tactical-player-controls.md) | 敵方陣形線要不要顯示 | 原版只畫側 0 那條（`word_1D33C`） | 靜態 |
| [`spec/38-list-windows.md`](../spec/38-list-windows.md) | 俘虜身分 | remake 的 `Posted` 是 bool，存不下 `+0x17` 的 0–5；俘虜狀態目前推不出來 | 靜態 |
| [`spec/38-list-windows.md`](../spec/38-list-windows.md) | 「看」與「選」的內容差異 | 原版兩種取法的**列表內容**不同（`../re/26` §4.2），remake 只統一了欄位 | 靜態 |
| [`spec/38-list-windows.md`](../spec/38-list-windows.md) | 「委任」那一格的顏色 | 實錄影格上看起來是紅字，但影片是壓縮過的、也沒有機器碼證據。remake 先畫成一般色 | 靜態 |
| [`spec/39-march-order-menu.md`](../spec/39-march-order-menu.md) | `sub_193E9` 的選單版面 | 只知道 `cx = 0x4Ch`；矩形與列高沒解，remake 先用既有的對話框樣式並標成差異 | 靜態 |
| [`spec/40-ai-march-decision.md`](../spec/40-ai-march-decision.md) | `+0x00` 位元 1 | 原版改目標時會設，`sub_12662` 讀它並呼叫未讀的 `sub_147BB`。remake 直接重下一次行軍（`March`），行為等價但不是同一條路 | 靜態 |
| [`spec/40-ai-march-decision.md`](../spec/40-ai-march-decision.md) | `sub_1487B` | AI 編成後挑第一個目標用的相鄰格選擇，未逐條解；remake 沿用既有的 `nearestFactionCity` | 靜態 |
| [`spec/90-same-state-parity.md`](../spec/90-same-state-parity.md) | 各視窗**內部**的排版 | 分區的外框已由機器碼定死（§3），框內的頭像／文字列座標仍是影片估值（`docs/spec/12` §7） | 靜態 |
| [`spec/90-same-state-parity.md`](../spec/90-same-state-parity.md) | 原版的畫面輸出是 640×400 還是 640×480 | DOSBox-X 的視窗尺寸與 VGA 模式要確認，否則兩邊尺寸對不上 | 實測 |
| [`spec/90-same-state-parity.md`](../spec/90-same-state-parity.md) | 調色盤季節組 | 兩側都要鎖同一組，否則整片顏色不同（`docs/formats/02`） | 靜態 |

## 3. 這支工具的盲區

**目前是 0。** 每一份提到「未解／未定案／未定位／缺口／未驗」的文件，
要嘛抽得出至少一條，要嘛在檔尾寫了 `<!-- 缺口：無 -->` 明講自己沒有。
這一條由 `--strict` 把關，`check.sh` 帶著它跑——**新文件寫了「未解」卻沒有未解小節，提交會被擋下來**。

抽取只認四種結構（專門的未解小節、表格最後一欄標未解的列、
`**未解**：…`、收尾是「…未解」的句子）。**寫在段落中段、
或用別的詞說「這個還不知道」的缺口抽不到**——下列檔案提到未解
卻一列都沒抽出來，要嘛缺口寫成別的句式，要嘛那些字樣只是在講別的事：

（沒有）

只印抽得到的部分，會讓解析失敗長得像「那份文件沒有缺口」。
這一節就是為了讓那個差別看得見。
