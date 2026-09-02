# 43 — 未解缺口總表（生成的）

**狀態：生成的清單，跑 `tools/py.sh tools/re_open_questions.py` 重出。
這一份不下結論，只把各文件的「未解」表集中到一處。**

- 日期：2026-09-02
- 產生工具：`tools/re_open_questions.py`
- 來源：`docs/` 底下所有文件的未解小節、表格裡標未解的列，與收尾是「…未解」的散句

既有的三張表回答別的問題（`CLAUDE.md` §10）：`docs/INDEX.md` 是**已解**的斷言、
[`21`](21-function-census.md) 是函式有沒有人寫過、[`24`](24-unread-function-catalogue.md) 是未讀函式在做什麼。
**這一份是唯一回答「還有什麼沒解」的。**

> 「擋住什麼」由來源目錄決定，「怎麼裁決」由關鍵字決定——兩欄都是機械算出來的，
> **不是逐條判斷過的優先序**。要排優先序請自己讀那一列指到的小節。

## 0. ⚠ 這個數字在量什麼

**491 列分布在 197 份文件，平均每份 2.5 列。**

⭐ **所以它比較接近「文件有多少份」，不是「原版還有多少沒解」。**
每寫一份新文件就帶進約三列自己的未解——而 `check.sh --strict` 還會
**要求**每份文件要嘛有未解小節、要嘛明講 `<!-- 缺口：無 -->`。
於是「解出新東西 → 寫一份文件 → 總數上升」是這個指標的常態，
不是退步。

⚠ 反過來也一樣：**數字變小不自動等於進度**。2026-08-21 的稽核
把它從 570 降到 431，而那 −139 沒有一列是靠解出新東西減掉的。

**要看進度請看別的東西**：`docs/spec/` 的 CONFORMED 份數、
`docs/playtest/` 的逐像素數字、`docs/re/21` 的覆蓋地圖。
這一份回答的是「還有什麼沒解」，**不是「還剩多少」**。

> ⭐ 另有 **6 列標成 `[DOS/BIOS]`**，**不計入下面的總數**——那是原版與 DOS／BIOS 之間的介面
> （`INT` 服務號、顯示卡暫存器、磁碟服務），而 remake 跑在 Go／Ebiten 上不跟它們講話。
> 清單在 §9。

## 1. 總量

| 擋住什麼 | 缺口數 | 靜態可解 | 要實測 | 兩版對照 |
|---|---:|---:|---:|---:|
| 規則正確性 | 15 | 12 | 3 | 0 |
| 資料保存 | 23 | 23 | 0 | 0 |
| 程式碼理解 | 160 | 154 | 6 | 0 |
| 驗收 | 84 | 73 | 11 | 0 |
| 外部資料 | 6 | 5 | 1 | 0 |
| 其他 | 203 | 190 | 13 | 0 |
| **合計** | **491** | 457 | 34 | 0 |

⚠ **這是列數，不是獨立問題數。** 索引檔的「現況」欄是別的文件的摘要，同一個缺口在那份文件自己的未解表裡還有一列——這類共 **0** 列（另有少數只是提到「未解」兩個字的圖例列）。

⭐ **狀態行自稱「已被…取代」的 9 份文件不計。**它們的「未解」是**當時**的未解，而每一批發行紀錄都會再列一次「Windows／macOS 實機」「沒有音效裝置」——一度讓 `docs/release/` 的 33 列其實只有 12 個獨立缺口。**還開著的缺口一定也在最新那一份**，所以跳過舊的不會漏。

| 來源目錄 | 列數 |
|---|---:|
| `docs/spec/` | 168 |
| `docs/re/` | 160 |
| `docs/playtest/` | 84 |
| `docs/formats/` | 23 |
| `docs/release/` | 18 |
| `docs/mechanics/` | 15 |
| `docs/mobile/` | 11 |
| `docs/promo/` | 6 |
| `docs/reference/` | 6 |

## 2.1 規則正確性（15 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`mechanics/15-realtime.md`](../mechanics/15-realtime.md) | `sub_10A65` 的內插演算法 | 只影響畫面 | 靜態 |
| [`mechanics/15-realtime.md`](../mechanics/15-realtime.md) | 最高速那一檔在原版實機上是多少 | 機器相依，要實測才有數字；只影響手感調校 | 實測 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | GUI 的戰後勝負／傷亡訊息 | 狀態層已結算，畫面上還沒顯示 | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | 原版投射物的圖形與動畫 | `BATTLE.SCH` 裡的圖形沒對出來，目前畫的是側別標記 | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | 少數戰術腳本與 `BATTLE` 資料欄位 | 十九個指令已全讀，剩的是資料側的殘留欄位 | 靜態 |
| [`mechanics/30-combat.md`](../mechanics/30-combat.md) | 重算路徑的時機 | 原版只在命令生效時算一次；remake 的兵每幀都可能被別人擋住，所以改成**每 30 幀可重算一次**（`replanInterval`，`internal/rules/tactical/soldier.go`）。這是 **remake 差異**，不是原版行為——原版被擋住之後怎麼處理沒有讀 | 靜態 |
| [`mechanics/40-economy.md`](../mechanics/40-economy.md) | 那兩個歸零的欄位很可能就是**本月累計的收入與支出 → 假說，待驗 | （散句） | 靜態 |
| [`mechanics/60-personnel.md`](../mechanics/60-personnel.md) | 武將 `+0` 旗標的 bit 0 | 只出現一次（劇本三的張衛），全庫沒有讀取端（`../re/77` §5）。要嘛是劇本作者的筆誤，要嘛是 PC-98 版才用的位元——後者要等 PC-98 的 `KI.EXE` 進 IDA 才驗得了 | 靜態 |
| [`mechanics/70-ai.md`](../mechanics/70-ai.md) | 正文散在各節的開放項目集中在這裡；每一條都寫下手點， | （未解小節內文） | 靜態 |
| [`mechanics/70-ai.md`](../mechanics/70-ai.md) | 委任 AI「統率力低就打得差」的**實作方式**（§3.2）——是決策頻率、判斷正確率，還是可用戰術的多寡 | 找戰術層讀武將統率力欄位的地方；`internal/rules/tactical` 目前不隨統率力變化 | 靜態 |
| [`mechanics/70-ai.md`](../mechanics/70-ai.md) | 說服請求裡 `SI+1Eh` 的完整結構語意、其他理由的索引公式、`cx ≥ 0x100` 分支（§1.3） | 三者都在 `sub_13830` 一帶；先補 `SI+1Eh` 的欄位表再談索引公式 | 靜態 |
| [`mechanics/70-ai.md`](../mechanics/70-ai.md) | 「拉玩家合攻一個對玩家友善的對象」的**意圖**（§1.2 末） | 機制已照抄，意圖無證據。要定案得靠實機長跑統計，不是反組譯 | 靜態 |
| [`mechanics/70-ai.md`](../mechanics/70-ai.md) | 事件 2／3 的完整接受流程與原版訊息、4／5 的逐位金額輸入、6／7 的 `TALK.DAT` 反應、9 的原版完整流程、11／12 的物件動畫 | 事件佇列的 runtime 已在，缺的是各事件的原版畫面流程；逐一對 `sub_14269` 一帶的分派做 oracle | 實測 |
| [`mechanics/70-ai.md`](../mechanics/70-ai.md) | 入佇列的隨機空格分佈 | 目前是 remake 的明示差異；要對齊得先解出原版排程時挑空格的那段 | 靜態 |
| [`mechanics/80-victory.md`](../mechanics/80-victory.md) | 四個劇本的結局是否不同 | **觸發條件四劇本共用**，差別只在初始勢力數；結局的十二幕也是一條路播完，沒有依劇本分支的證據（`../re/70` §3）。**但沒有實跑四個劇本對過** | 實測 |

## 2.2 資料保存（23 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`formats/01-talk-dat.md`](../formats/01-talk-dat.md) | `\1`／`\2`／`\3`／`\4` 的**文字語意** | 機制已解（§3：`sub_1075B` 換算索引、`sub_1084A` 逐字讀、`CS:[SI+08A4h]` 七項跳躍表分派），`\6`（排版控制，消耗一個 16 位元參數後 `DX` 左移 `0x30`）與 `\7`（走 `sub_1062F` 的數值繪製）兩支 handler 也定案了。**這四支的 ha… | 靜態 |
| [`formats/01-talk-dat.md`](../formats/01-talk-dat.md) | 只有格 0 有字的三組（`0x1A1`／`0x1AD`／`0x1AE`） | 呼叫端傳不傳 `ah`、為什麼不需要八格變體，未讀。`0x1AD`／`0x1AE` 像單挑台詞，但不在 `../re/75` 盤過的 `0x1B6`–`0x1CD` 範圍內 | 靜態 |
| [`formats/01-talk-dat.md`](../formats/01-talk-dat.md) | 訊息索引與遊戲事件的對應 | 逐則的顯示時機沒有全表。**已對出來的散在各規格**（進言 `../spec/44`、遷都 `../spec/64`、結局 `../spec/30`），未讀的部分見 `../re/24` | 靜態 |
| [`formats/03-grf-images.md`](../formats/03-grf-images.md) | `0x0480` | 24×16 × 3 / 兵種圖示的**橘色版**：馬／弓／步 / 尚未找到取用端 | 靜態 |
| [`formats/04-map-sch-container.md`](../formats/04-map-sch-container.md) | 狀態：容器格式的索引層 READY，壓縮演算法未解。 | （散句） | 靜態 |
| [`formats/05-mmap-worldmap.md`](../formats/05-mmap-worldmap.md) | `MMAP.MCH` 的 object **type 3** | 圖塊、`0xA000` metadata 與事件 12 的火災／暴動（type 1／2）查表已解（`../re/14`）。type 3 的事件語意、object timer 與逐 frame 的原版時序仍未知——remake 的 timer 是呈現層 substitute | 靜態 |
| [`formats/05-mmap-worldmap.md`](../formats/05-mmap-worldmap.md) | 比對的正是方向碼那個欄位）。 強證據，未定案。 | （散句） | 靜態 |
| [`formats/07-battle.md`](../formats/07-battle.md) | 每個戰場 4,096 B 的表頭與尾段各 64 byte | 內容未解（§2.1） | 靜態 |
| [`formats/07-battle.md`](../formats/07-battle.md) | 表頭 16 筆但資料只有 3 組 | 多出來的 13 筆是什麼未解（§3） | 靜態 |
| [`formats/07-battle.md`](../formats/07-battle.md) | 所以地形是 **64 寬 × 62 列**，而整塊 4,096 B 是 **64 × 64**—— 兩者內容未解。 | （散句） | 靜態 |
| [`formats/07-battle.md`](../formats/07-battle.md) | ——表頭筆數多於資料組數，原因未解。 | （散句） | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0x0008` | 51 / 未解的全域狀態（一起載入 `cs:0CF0h`） | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0x0080` | 2,112 / **勢力表：22 筆 × 64 B**（`docs/re/06` §5）＋ 其後未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0x1EC0` | 7,168 / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0x3AC0`…`+0x42C0` | — / 未解 | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0` | 1 / 旗標。四劇本 508 筆只出現 **8 種值**（0／0x80／0x81／0x90／0xA0／0xC0／0xD0／0xE0）。**bit 7 ＝ 在場**（381 筆）；**bit 6 ＝ 主公型**（43 筆，全是現任君主，被俘而舊主勢力已滅時清掉並把 `+30` 說話類型 `+3`，`sub_129C… | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+20`（`0x14`） | 1 / 值域 **0–7**、8 種，四劇本 508 筆都有值（0 佔 158、2 佔 98）。**與 `+0x1E` 的說話類型同值域但不同值**，交叉表也沒有結構 / 語意未解。**讀取端全庫掃過是零**（`../re/77` §5） | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+27`（`0x1B`） | 1 / 只有兩種值：0（481 筆）與 `0xFF`（27 筆）。劇本一的 `0xFF` 是笮融、張英、周倉、嚴輿、周昕、李別、刑道榮、劉賢、武安國九人 / 語意未解。**讀取端全庫掃過是零**（`../re/77` §5） | 靜態 |
| [`formats/08-sinario-save.md`](../formats/08-sinario-save.md) | `+0`／`+3` | 未解 | 靜態 |
| [`formats/09-cutscene-images.md`](../formats/09-cutscene-images.md) | `OPEN_S1` 的捲動節拍 | 版面三塊已解（`../re/76` §6），但每步捲幾列、26 格的動畫多久換一格，要先定 `int 61h` 的計時單位 | 靜態 |
| [`formats/09-cutscene-images.md`](../formats/09-cutscene-images.md) | `OPEN_S5` 的來源列距 | `sub_103EF` 每兩列消耗 720 B，推不出單純的「一列 240 B」；沒有畫面驗證（同 §10） | 靜態 |
| [`formats/09-cutscene-images.md`](../formats/09-cutscene-images.md) | `GAMEOVER.DAT` 誰播 | 不在 `D7END.EXE` 的十二幕裡。**推測是 `KI.EXE` 的敗北路徑**（`../re/59`），沒有找到取用端 | 靜態 |
| [`formats/10-end-s15-namechars.md`](../formats/10-end-s15-namechars.md) | 勢力 `+0x02 = 0x7F` 時，訊息裡的 `{4}` 從哪裡取名 | 推測從 `5222h`，`sub_1075B` 那條路沒回頭讀 | 靜態 |

## 2.3 程式碼理解（160 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`re/01-first-recon.md`](../re/01-first-recon.md) | 是加了新過場、還是把原本的長段拆開，未解。 | （散句） | 靜態 |
| [`re/01-first-recon.md`](../re/01-first-recon.md) | `PASS.MAP`／`PASS.SCH` | dosv / **PC-98 沒有**。關隘資料，移植時新增或改名。未解 | 靜態 |
| [`re/02-palette-routine.md`](../re/02-palette-routine.md) | OPENPAL`（6 組）、`ENDPAL`（12 組）的分組對應哪些畫面。 | （未解小節內文） | 靜態 |
| [`re/02-palette-routine.md`](../re/02-palette-routine.md) | 設定表 `cs:0x5FF4` 每筆後三個 byte 是什麼（第 4 筆的第二個 word | （未解小節內文） | 靜態 |
| [`re/03-image-blitter.md`](../re/03-image-blitter.md) | ICONGRF` **段 1 `0x0000` 那一塊畫了什麼**。 | （未解小節內文） | 靜態 |
| [`re/03-image-blitter.md`](../re/03-image-blitter.md) | sub_1FAC2` 是另一支繪製常式（`shl al, 1` 後才 `mov cx, ax`），用途未解。 | （散句） | 靜態 |
| [`re/04-mmap-entry-points.md`](../re/04-mmap-entry-points.md) | MMAP.MCH` 的 object **type 3**：事件語意、原版 object timer 與畫面 | （未解小節內文） | 靜態 |
| [`re/05-battle-selection.md`](../re/05-battle-selection.md) | `+0x08` | 勢力相符時取用的值 / 未解 | 靜態 |
| [`re/06-game-clock.md`](../re/06-game-clock.md) | `cs:0CF0h` 那 59 byte 裡除了時鐘的其餘部分 | 對照存檔 diff | 靜態 |
| [`re/06-game-clock.md`](../re/06-game-clock.md) | `sub_10A65` 的內插演算法 | 直接讀 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | 10 | `sub_13496` / 訊息-only：建立武將／參數 formatter 游標；持久狀態尚未找到 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | `sub_15940` 的兩個分支 | 已派駐武將的每月行動，會發訊息 `0x41`／`0x42`。分支 2 有一行 `mov byte ptr [si+1Ch], 18h`（把所屬勢力寫成 24）**與「+1Ch 是勢力編號、只有 0–21」矛盾**，還沒解釋 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | `sub_14269`／`sub_13EFD` | 事件 11／12 寫入的據點 `+0x15` marker 在據點輪轉時先扣防災值；不足時再扣上昇值、生產力與城兵，已接入 `World.applyCityDisasterEffect`；物件動畫仍未完 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | 武將 `+1Ah` | 官員「要錢中」的旗標／金額，`sub_12FBF` 的事件會寫它 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | 武將 `+0Eh`／`+0Fh`／`+10h` 是不是兵種適性 | 需要與戰鬥程式對照 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | `sub_157FE` 觸發的事件內容 | `sub_12FBF(ax=0Dh, dx=196h)` | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | 軍團記錄剩下的欄位 | 段內 `2240h`，**64 B／筆、127 筆**（不是 32 B）。已具名的見 `08` §4 與 `34` | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | 段內 `2040h`／`2140h` 的兩張 16 × 16 B 表 | `sub_123FF` 會在 `2040h` 那張找空位配置；`2140h` 那張開局已有 16 筆，`+2`／`+4` 看起來是地圖座標 | 靜態 |
| [`re/07-monthly-settlement.md`](../re/07-monthly-settlement.md) | 123 個武將裡有 80 人只有一欄非零、32 人兩欄、9 人三欄—— → 強證據，待驗 | （散句） | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `sub_15456` 用 stride 32 掃軍團表 | 與 64 矛盾，疑似原版 bug | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `+0x0C` | 行軍中的暫存（`sub_12708` 寫） / 未解 | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `+0x04` | byte / `sub_1E81C` 的回傳 `ah` / 未解 | 靜態 |
| [`re/08-hourly-update.md`](../re/08-hourly-update.md) | `+0x08` | word / 另一張表的索引（`bx << 3`） / 未解 | 靜態 |
| [`re/09-combat.md`](../re/09-combat.md) | **戰術完整結算** | `TestNormalScenarioTacticalBattleTerminates` 已證實真實正常攻城的狀態層勝負／傷亡回寫，`wlgame-ai-postbattle.png` 證明正常 GUI 回戰略；GUI 戰後訊息、完整狀態對拍與少數分支仍未完 | 靜態 |
| [`re/09-combat.md`](../re/09-combat.md) | `sub_1AD7F` 攻擊分支 | `shootSpecial` 已接入 `CH=0x20` 的相鄰格／垂直效果；`+0x1E` 的初始化／上移／下移／交換來源與 `sub_1AC55` 的 raw 比較已確認並接成 `PlaneHigh`，普通箭原版 SCH 單幀圖形已接回，完整投射物動畫／同狀態對拍仍待確認 | 靜態 |
| [`re/09-combat.md`](../re/09-combat.md) | 原版／remake 同狀態對拍 | 需有效時序原版存檔或可重建的同狀態 oracle | 實測 |
| [`re/09-combat.md`](../re/09-combat.md) | 地形係數表的列 2、武將適性 `+0x10` | 兩者都要 `al = 2` 才取得到，戰略層沒有呼叫點 | 靜態 |
| [`re/09-combat.md`](../re/09-combat.md) | 武將旗標 `+0x00` 的 bit 4（自刎） | 值域已知有 7 種，只解出這一個位元 | 靜態 |
| [`re/09-combat.md`](../re/09-combat.md) | 據點 `+0x10`／`+0x11` 被攻城扣減 | 欄位語意已知（上昇值／防災值），但「被打過的城成長變慢」還沒在數值上驗過 | 靜態 |
| [`re/09-combat.md`](../re/09-combat.md) | `[si+3]` 的 0／1／≥2 是誰設的 | 決定哪一支軍團會進戰鬥畫面 | 靜態 |
| [`re/11-tactical-battle.md`](../re/11-tactical-battle.md) | `ds:0D2F8` | 4,096 / 未解（第二份戰場？） | 靜態 |
| [`re/11-tactical-battle.md`](../re/11-tactical-battle.md) | `ds:0D306` | 30,720 / 未解 | 靜態 |
| [`re/11-tactical-battle.md`](../re/11-tactical-battle.md) | 兵士記錄剩下的欄位 | `ds:0D30E`，32 B／筆 / 目前具名的有 `+0x00`／`+0x01`／`+0x02`／`+0x03` 體力／`+0x04` 大將／`+0x05` 面向／`+0x14` 陣形座標／`+0x16`・`+0x17` 繞路游標／`+0x19` 疲勞／`+0x1A`・`+0x1B` 命令／`+0x1E` Z… | 靜態 |
| [`re/11-tactical-battle.md`](../re/11-tactical-battle.md) | `ds:0D306` 那 30,720 B | 未解 / 大小是 `0x7800`，與任何已知的表都對不起來 | 靜態 |
| [`re/11-tactical-battle.md`](../re/11-tactical-battle.md) | `loc_1A065` 的自我修改碼 | `▶▶` 列切換的機制 confirmed（`byte_1A06A` 在 `0xEB`／`0x74` 間切），**擋掉的是什麼**沒逐行讀（`60` §） | 靜態 |
| [`re/12-diplomacy-dialogue.md`](../re/12-diplomacy-dialogue.md) | `AH` 的完整欄位名稱 | 語意由日中原文並列確認，欄位名本身未定（§3） | 靜態 |
| [`re/12-diplomacy-dialogue.md`](../re/12-diplomacy-dialogue.md) | #367–#372／#380–#385 的 AH／信賴度次要回覆 | 未解，不可當成完整的原版對話流程（§8） | 靜態 |
| [`re/12-diplomacy-dialogue.md`](../re/12-diplomacy-dialogue.md) | #73／#77 | 未定位，不得拿來補接事件 6／7（§9） | 靜態 |
| [`re/12-diplomacy-dialogue.md`](../re/12-diplomacy-dialogue.md) | 事件 6／7 次要 TALK 的 formatter 參數契約 | 缺參數且語意未知，維持 fail-closed（§10） | 靜態 |
| [`re/15-event10-producer.md`](../re/15-event10-producer.md) | 以下來源沒有證據，不能補成事實：未被 IDA 建成函式的 far code、以暫存器或指標 | （未解小節內文） | 靜態 |
| [`re/17-dosv-audio-tsr.md`](../re/17-dosv-audio-tsr.md) | `0x330` 的用途 | MPU-401 的標準埠，沒找到讀它的地方 | 靜態 |
| [`re/17-dosv-audio-tsr.md`](../re/17-dosv-audio-tsr.md) | 效果碼 ↔ 聽起來像什麼 | `SOUND.DAT` 的記錄結構已解（`57` §6），但哪一號對應哪個動作只有 §3 的三個 | 靜態 |
| [`re/19-outcome.md`](../re/19-outcome.md) | 勢力滅亡 selector | 未定位。remake 只顯示克制的 fallback 句，不冒充原版文字 | 靜態 |
| [`re/20-ida-re-coverage-audit.md`](../re/20-ida-re-coverage-audit.md) | 同狀態動態 oracle | 沒有可重放的存檔／輸入序列，所以「原版等價」目前無法驗。**這是還沒做，不是做不了**——DOS/V 的密碼頁空白確認就會過（`../playtest/18`），PC-98 側連除錯器都接好了（`../playtest/21`） | 實測 |
| [`re/20-ida-re-coverage-audit.md`](../re/20-ida-re-coverage-audit.md) | 逐幀執行順序 | 顯示串列與相機已重建，但整幀的呼叫順序沒有逐幀對過 | 靜態 |
| [`re/20-ida-re-coverage-audit.md`](../re/20-ida-re-coverage-audit.md) | `loc_1A065` 的 runtime bytes | 自我修改碼，靜態影像看不到每輪的實際內容（§2.2） | 靜態 |
| [`re/20-ida-re-coverage-audit.md`](../re/20-ida-re-coverage-audit.md) | 四層差分（terrain／display list／composited／HUD） | 沒有 machine-readable diff，目前只有 layout-only 比較 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | `+0x00` | 2 B / 未解 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | `+0x06`–`+0x0F` | 10 B / 未解 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | `+0x04` 那張表 | 大小與音效記錄相同（3 × 16 B），但驅動沒讀它。見 `57` §8 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | `+0x00`、`+0x06`–`+0x0F` | 12 B 未解 | 靜態 |
| [`re/23-bgm-resource-format.md`](../re/23-bgm-resource-format.md) | 曲號 ↔ 場景的對應 | `KI.EXE` 呼叫端傳哪個索引還沒對過 | 靜態 |
| [`re/27-list-row-fields.md`](../re/27-list-row-fields.md) | 開局選勢力的逐列 `sub_17BC0` | 未逐欄對照（欄位與勢力一覽重疊，但少了外交兩欄） | 靜態 |
| [`re/28-text-number-rendering.md`](../re/28-text-number-rendering.md) | 屬性的其餘位元 | bit 2 是陰影已證實；`0x9001`／`0x9000` 的 bit 0 差在哪未讀 | 靜態 |
| [`re/28-text-number-rendering.md`](../re/28-text-number-rendering.md) | `word_10D4C` 那一組 | 來源已解——`sub_100DF` 開機把 `ICONGRF` 段 3 切五塊，`word_10D54` 是 `+0x0840` 的 11 格 × 16 列數字字模（`../spec/52` §4）；緊接在後的 `+0x08F0` 另有一組 11 格，用途未解 | 靜態 |
| [`re/29-font-service-int15.md`](../re/29-font-service-int15.md) | `END_S10/S11` 與 `STR.EXE` 檔名不同步 | §6，要實跑裁決 | 實測 |
| [`re/29-font-service-int15.md`](../re/29-font-service-int15.md) | `END_S13.DAT` 前 408 格的來源 | 不是 `stdfont.15` 的任何一段，也不是 `usrfont.15m`（256 B） | 靜態 |
| [`re/29-font-service-int15.md`](../re/29-font-service-int15.md) | `END_S15.DAT`（5,242 B） | `KI.EXE` 的字串表引用它，但**不是字型**（大小不是 30／15 的倍數），也不是過場圖（沒有 `00 F4 01` 檔頭），Big5 解碼是亂碼 → 疑似壓縮 | 靜態 |
| [`re/30-corps-formation-ui.md`](../re/30-corps-formation-ui.md) | 軍團 `+0x00` 的位元 3／4／5 | 位元 1（有指令）、2（委任，`45`）已解；其餘仍未見成對的寫入端（`34` §4） | 靜態 |
| [`re/31-faction-picker-screen.md`](../re/31-faction-picker-screen.md) | 分派表已印出，但 `sub_15AD1 → sub_15AFC` 的進入路徑仍未定位。 | （散句） | 靜態 |
| [`re/31-faction-picker-screen.md`](../re/31-faction-picker-screen.md) | `cs:6056` 表的長度 | 前六筆是一組小 handler，後五筆疑似越過表尾（§1.2） | 靜態 |
| [`re/32-strategy-detail-panels.md`](../re/32-strategy-detail-panels.md) | 軍團 `+0x00` 的位元怎麼清 | 三處設定都找到了，清除點未找到 | 靜態 |
| [`re/33-shared-draw-helpers.md`](../re/33-shared-draw-helpers.md) | `cs:word_10D40` | 肖像圖庫所在的段，誰載入它未追 | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | 0 | `sub_12459`／`sub_126FF`（候選） / `sub_12533`（候選） / 未定 | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | 4 | **`sub_12B3C`（confirmed）** / `sub_12BA8`（候選） / **設定端定案**（2026-09-02）：`sub_12B3C` 開頭就是 `or byte ptr [si], 10h`，而同一支後面用 `[si+10h]`／`[si+12h]` 當地圖座標——**`si` 確定是… | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | 位元 0／4／5 的語意 | 有成對的設定與清除點，但那幾支函式的操作對象尚未逐支確認是軍團 | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | 位元 3 | 掃描裡沒出現。間接寫入抓不到，不能據此說它不存在 | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | `sub_12977` 的 `mov byte [si], 8` | 該函式同時碰武將表與軍團表，`si` 指哪一張未確認 | 靜態 |
| [`re/34-corps-status-bits.md`](../re/34-corps-status-bits.md) | `sub_12880` 的 `or [si], 20h` | 表歸屬指向據點表，語意要另外讀（`sub_144A9`／`sub_144D6` 已解：Stage 10／11 把目標校正成首都並設位元 1，見 `64` §2） | 靜態 |
| [`re/35-strategy-ui-module-map.md`](../re/35-strategy-ui-module-map.md) | `sub_18FC9` 叢 | — / 存檔畫面的槽位與按鈕對應未驗（§2.8） | 靜態 |
| [`re/40-garrison-relief-request.md`](../re/40-garrison-relief-request.md) | `+0x20` 與 `+0x14` 的關係 | §5 的張力，要實測 | 實測 |
| [`re/40-garrison-relief-request.md`](../re/40-garrison-relief-request.md) | 據點 `+0x00` 的 bit 4／5 | bit 6／7 已解（§2），中間兩位未見 | 靜態 |
| [`re/42-leaf-functions.md`](../re/42-leaf-functions.md) | `cs:byte_198A6` 位元 3 | **全庫沒有任何一處寫它**（2026-09-02 三種寫法全掃：`byte_198A6`、`cs:byte_198A6`、`ds:98A6h`，整個 `KI.EXE.asm` 只有 19 行提到這個位址）。位元 0 在 `or ds:98A6h, 1`／`and ds:98A6h, 0FEh` 成對、位元 1 … | 靜態 |
| [`re/42-leaf-functions.md`](../re/42-leaf-functions.md) | `sub_1E9A7` 的 8 bytes 參數表 | 寫入端已解（2026-09-02）：整支是 `bx = bl × 8 + 0EAF1h` 之後 `cs:[bx] = ax`／`cs:[bx+2] = dx`／`cs:[bx+4] = cx`——**一筆 8 B 但只寫前 6 B**，後兩個 byte 沒有任何寫入端。**表的內容（各筆代表什麼）仍未讀** | 靜態 |
| [`re/42-leaf-functions.md`](../re/42-leaf-functions.md) | `byte_1020E`／`byte_10CF9` | 音源相關的兩個旗標 | 靜態 |
| [`re/44-threat-and-reinforcement-ai.md`](../re/44-threat-and-reinforcement-ai.md) | 據點 `+0x00` 的 bit 4／5 | bit 6／7 是威脅旗標、低 4 位是敵方鄰居，中間兩位仍未見寫入端 | 靜態 |
| [`re/44-threat-and-reinforcement-ai.md`](../re/44-threat-and-reinforcement-ai.md) | `+0x20` 與 `+0x14` 的張力 | `sub_14575` 與 `sub_14155` 都只寫 `+0x20`，`40` §5 的張力還在 | 靜態 |
| [`re/46-strategy-chrome-cell-layer.md`](../re/46-strategy-chrome-cell-layer.md) | 樣式碼 | 只確定 `0` ＝ 擦除、`0x0B` ＝ 指令列、`0x0C`／`0x0F` 出現在別處；完整值域未列 | 靜態 |
| [`re/46-strategy-chrome-cell-layer.md`](../re/46-strategy-chrome-cell-layer.md) | `ax = 0F01h`／`0801h` | 顏色／樣式的位元編碼未逐位對過 | 靜態 |
| [`re/47-main-screen-window-registry.md`](../re/47-main-screen-window-registry.md) | `0x80` | 繪製時 `and …, 7Fh` 清掉 / 未解 | 靜態 |
| [`re/47-main-screen-window-registry.md`](../re/47-main-screen-window-registry.md) | 格子屬性 bit `0x80` | 擦除時被清掉，沒找到設它的地方 | 靜態 |
| [`re/47-main-screen-window-registry.md`](../re/47-main-screen-window-registry.md) | 系統視窗開著時時間停止 | 說明書明講，機器碼的實作位置未找（`sub_15FAA` 的等待迴圈是候選） | 靜態 |
| [`re/47-main-screen-window-registry.md`](../re/47-main-screen-window-registry.md) | 右鍵表 `funcs_159C0` 的真實表長 | 已 dump（`71` §2.1）：它與左鍵表 `off_159D2` **只差 9 個 word 且內容重疊**，前九筆沒有一筆是函式起點。是「表只有 9 筆」還是「兩張刻意重疊」，靜態分不出來 | 靜態 |
| [`re/48-window-display-list.md`](../re/48-window-display-list.md) | `08` 的模式 byte | `03` 只畫字、`01` 連背景一起填，是**強推論**——兩個用例（系統選單的「 ＯＫ 」、注音聲母列）都只有這個讀法說得通，但 `sub_106F5` 沒逐行讀（`55` §3） | 靜態 |
| [`re/48-window-display-list.md`](../re/48-window-display-list.md) | `sub_1E9A7(bl=0, ax=1800h, cx=2020h)` | **登記內容已解**（2026-09-02）：`sub_1030F` 在 `start` 裡做兩件事——① `sub_1E993(ax=0E16h, dx=cs, bx=cs, cx=word_10D50)` 把 far 指標 `cs:0E16`（顯示清單本身）寫進 `cs:dword_1EAE9`，另外兩個參數… | 靜態 |
| [`re/48-window-display-list.md`](../re/48-window-display-list.md) | `op 01` 的用法 | 它是直線（§2.2），但 handler 不展開座標而十個場景又沒用到它——**預期的呼叫方式無法驗證** | 靜態 |
| [`re/48-window-display-list.md`](../re/48-window-display-list.md) | `op 02` 與 `op 03` 的差別 | 兩支都畫矩形（`sub_1F020` 對 `cs:F1A3`），前者另有五個戰術區呼叫者。哪一支是實心、哪一支帶遮罩，沒有資料可分辨 | 靜態 |
| [`re/49-corps-formation-window.md`](../re/49-corps-formation-window.md) | `sub_1F9B0` 的 `ax = 1003h` | 貼圖的樣式參數；`sub_10C14` 用 `0801h`（`46` §3）。位元編碼未逐位對過 | 靜態 |
| [`re/51-corps-info-window.md`](../re/51-corps-info-window.md) | `or byte ptr [si], 2` | 位元 1 ＝「有指令」（`34`），這裡是它的其中一個寫入端 | 靜態 |
| [`re/52-slot-select-window.md`](../re/52-slot-select-window.md) | 檔名 | `sub_18C20` 沒設 `dx`，靠 `sub_18B7C` 的 `push dx`／`pop dx` 從更上層傳進來 | 靜態 |
| [`re/54-advisor-naming-window.md`](../re/54-advisor-naming-window.md) | 屬性低 byte | `01` 與 `03` 的差別未讀（§3） | 靜態 |
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
| [`re/59-game-over-exit-codes.md`](../re/59-game-over-exit-codes.md) | `sub_14DF0` 的 CF | 「找不到替代據點」與「據點數 0」是不是同一件事，還沒逐行讀 | 靜態 |
| [`re/59-game-over-exit-codes.md`](../re/59-game-over-exit-codes.md) | 無主城 `0x18` | 值 24 落在 22 個勢力之外，但劇本裡有沒有無主城沒查過 | 靜態 |
| [`re/60-tactical-sidebar.md`](../re/60-tactical-sidebar.md) | 城兵臨時軍團的主將名 | `sub_14F58` 已讀（2026-09-02）：`push [si+2] ＋ 0FF00h` 當 formatter 參數（`di = sp`）→ `sub_10CDE` 嗶一聲 → `al = 93h`（**預設肖像，不是訊息索引**）→ `sub_18810`，訊息索引在呼叫端的 `cx`（`0x1B`… | 靜態 |
| [`re/60-tactical-sidebar.md`](../re/60-tactical-sidebar.md) | 段 1 `0x0000`／`0x0800`／`0x1000`／`0x1800`／`0x3500` 的圖形內容 | 貼點與尺寸 confirmed，**圖上畫了什麼**要另外解碼（`../formats/03` §5.3 的 UI 語意缺口） | 靜態 |
| [`re/60-tactical-sidebar.md`](../re/60-tactical-sidebar.md) | 熱區 `0x01`／`0x1F` | 兩張表裡都有 handler，但沒找到註冊它們的 `sub_1E3D7` 呼叫點 | 靜態 |
| [`re/60-tactical-sidebar.md`](../re/60-tactical-sidebar.md) | 側欄美術的調色盤 | 本份記的都是**調色盤索引**，不是 RGB。要比顏色得用 `GAMEPAL.BRG` 的當季 bank | 靜態 |
| [`re/61-timer-tick-source.md`](../re/61-timer-tick-source.md) | 音樂 tempo 分頻器 `cs:0B68h` 的算式 | `0x859` 那 20 條指令：`al = ((0FFh − ah) × 13) >> 3`，`ah` 從哪來沒讀 | 靜態 |
| [`re/61-timer-tick-source.md`](../re/61-timer-tick-source.md) | `cs:099Eh` 的 bit 1 | 「音樂啟用」是從用法推的，寫入端沒讀 | 靜態 |
| [`re/61-timer-tick-source.md`](../re/61-timer-tick-source.md) | 無音效驅動時的行為 | §3 推論「會卡死」，**沒有實測**——DOSBox 拿掉 `YNSOUND.COM` 跑一次就能驗 | 實測 |
| [`re/62-strategy-minimap.md`](../re/62-strategy-minimap.md) | `byte_198A7` 的初值 | 靜態影像裡是 `0`。**開新遊戲時有沒有被寫過沒查** | 靜態 |
| [`re/62-strategy-minimap.md`](../re/62-strategy-minimap.md) | 圖例底圖在哪個資源 | `sub_1FA37` 的 `ds` 來自 `word_10D50`；`47` 記成段 3 `0x09A0`，兩者沒對過 | 靜態 |
| [`re/62-strategy-minimap.md`](../re/62-strategy-minimap.md) | ⭐ **熱區 `0x16` ＝ 點地圖捲鏡頭**：`sub_15AB6` 把螢幕座標減掉地圖區原點、 | （未解小節內文） | 靜態 |
| [`re/63-ground-plane-map.md`](../re/63-ground-plane-map.md) | 段變數的配置迴圈 | `word_1D2F6`–`word_1D30E` 沒有直接寫入的 xref，§1 的相鄰關係是推論 | 靜態 |
| [`re/63-ground-plane-map.md`](../re/63-ground-plane-map.md) | 命令 6 為什麼擋高平面橫移 | `[si+1Ah] == 6`，命令碼 6 是什麼沒對過 | 靜態 |
| [`re/65-ai-march-decision-chain.md`](../re/65-ai-march-decision-chain.md) | `loc_1491B` 的完整成本模型 | 只解出「穿過非己方據點 ＋0xA6 並設高位元」（§8.1）。廣度優先搜尋本身的佇列結構與其他成本項沒逐條讀 | 靜態 |
| [`re/65-ai-march-decision-chain.md`](../re/65-ai-march-decision-chain.md) | remake 的對應 | §8 已實作（`../spec/43`）；**§8.4 的「逐站前進」沒有移植**，remake 一次算完整條路 | 靜態 |
| [`re/65-ai-march-decision-chain.md`](../re/65-ai-march-decision-chain.md) | `sub_128F4` 的 STC 分支 | 走到敵方據點時呼叫 `sub_1291A`（俘虜／脫離判定），之後 `di` 不可信。本文件的 `di` 推論只涵蓋一般路徑 | 靜態 |
| [`re/66-message-box-geometry.md`](../re/66-message-box-geometry.md) | `sub_10AD9` 的 `cx = 40B0h` | 肖像繪製的尺寸／來源參數，沒逐位元解 | 靜態 |
| [`re/66-message-box-geometry.md`](../re/66-message-box-geometry.md) | `sub_189A4(al=1, dx=0, bx=2, cx=151Bh)` | `sub_13D09` 在貼完 `IVENTGRF` 之後畫的框，與 `sub_1895D` 是不是同一組單位沒驗 | 靜態 |
| [`re/66-message-box-geometry.md`](../re/66-message-box-geometry.md) | `IVENTGRF` 插圖本身的位置 | `sub_13D09` 的 `dx = 0E07h` 是餵給 `sub_1E38C`（讀檔）的參數，不是座標。插圖在畫面上的位置沒量 | 靜態 |
| [`re/67-city-emblem-on-strategy-map.md`](../re/67-city-emblem-on-strategy-map.md) | 推論等級：位置與顏色分類 **強證據**（四座 × 三種歸屬）；繪製常式 **未解 | （散句） | 靜態 |
| [`re/67-city-emblem-on-strategy-map.md`](../re/67-city-emblem-on-strategy-map.md) | 換圖塊的那一支機器碼 | **規則是從資料與畫面反推的，不是讀出來的。** 每一格都對得上，但「原版在哪裡做這件事」還沒定位——`sub_1D615` 只複製圖塊編號，`sub_11CC9` 的兩個 overlay 只畫軍團與災害 / 找誰在據點換手時改地圖或格子記錄（`sub_188CC` 一帶） | 靜態 |
| [`re/67-city-emblem-on-strategy-map.md`](../re/67-city-emblem-on-strategy-map.md) | 208 那一組 | 全圖 4 格，這一輪的兩張截圖裡沒有入鏡 / 找一局讓它入鏡，驗 206／207／208 | 實測 |
| [`re/67-city-emblem-on-strategy-map.md`](../re/67-city-emblem-on-strategy-map.md) | 「圖例選中的勢力」 | 縮小地圖有第四種顏色（`62` §2），大地圖有沒有對應的圖塊沒驗 / 開縮小地圖、切圖例第二格再截一張 | 靜態 |
| [`re/67-city-emblem-on-strategy-map.md`](../re/67-city-emblem-on-strategy-map.md) | 230 為什麼分位置 | 關隘上下換、大城左右不換。remake 照位置實作，但沒有機器碼解釋 / 同第一列 | 靜態 |
| [`re/68-t3-frontier-functions.md`](../re/68-t3-frontier-functions.md) | `sub_1304E(al=7)` 到底登記了什麼 | `sub_1676F` 只看它的進位旗標。那一支是 T1，但 `al` 的七種值各代表什麼還沒逐一對過 | 靜態 |
| [`re/68-t3-frontier-functions.md`](../re/68-t3-frontier-functions.md) | `sub_16D56` 的 `1,1,3,3,2,2` 對應哪三個兵種 | 位移確定是六個編成槽的兵種欄；值到兵種的對映要與 `kindFromByte` 對一次 | 靜態 |
| [`re/69-t2-cross-reference.md`](../re/69-t2-cross-reference.md) | `0x2040` 那張 16 筆表 | `sub_12438` 依 `(dx, bx)` 作廢一筆，`[si] ≥ 0x80` 是「這筆有效」 / 找誰寫 `[si+2]`／`[si+4]` | 靜態 |
| [`re/69-t2-cross-reference.md`](../re/69-t2-cross-reference.md) | `sub_1E6FF` 那張待繪表 | 欄位對應到什麼還沒查 / `byte_1E47F` 的其他使用端 | 靜態 |
| [`re/70-d7end-ending-player.md`](../re/70-d7end-ending-player.md) | `sub_1041E`（`ENDPAL.BRG`）怎麼套 | 只知道它載檔 / 與 `GAMEPAL.BRG` 同格式的話直接沿用（`../formats/02`） | 靜態 |
| [`re/70-d7end-ending-player.md`](../re/70-d7end-ending-player.md) | 淡入淡出的階數與色階 | 17 階（`cx` 0–0x10）已確定，每階怎麼算沒讀 / `sub_1035F`／`sub_103DC` | 靜態 |
| [`re/70-d7end-ending-player.md`](../re/70-d7end-ending-player.md) | `cs:0x780` 那張字幕描述子表 | §3.1 解出結構（幕序索引 → 筆數 ＋ 每筆三個 word），**表的內容沒 dump** / `ida_dump.py` 對 `D7END.EXE` 的 `0x780` 起 | 靜態 |
| [`re/70-d7end-ending-player.md`](../re/70-d7end-ending-player.md) | BGM 的**起訖時點** | `ENDBGM.DAT` 走 INT 61h、與 `KI.EXE` 同一條音源路徑（已解）。⚠ 2026-08-23 起 remake **整段結局都放 `endbgm-0`**（`musicTrack()`）——先前放的是 `overbgm-0`，那是 `D7OVER.EXE` 的遊戲結束曲。**剩下的缺口只有… | 靜態 |
| [`re/71-strategy-hotspot-dispatch.md`](../re/71-strategy-hotspot-dispatch.md) | 右鍵表的真實表長 | §2.1。兩種讀法都說得通，靜態影像分不出來，要動態取樣哪些碼會配右鍵 | 靜態 |
| [`re/71-strategy-hotspot-dispatch.md`](../re/71-strategy-hotspot-dispatch.md) | `funcs_159C0[0x00]`–`[0x08]` 那九筆 | 都不是函式起點。是別的資料還是真的 handler，沒查 | 靜態 |
| [`re/71-strategy-hotspot-dispatch.md`](../re/71-strategy-hotspot-dispatch.md) | `22` 的「`off_159D2` 的其餘槽位」、 | （未解小節內文） | 靜態 |
| [`re/72-world-map-display-list.md`](../re/72-world-map-display-list.md) | 那 110 張軍團圖的逐張外觀 | 算式定案、抽驗過勢力 0 靜止那一張（紅色軍旗），**22 × 5 沒有逐張看過** | 靜態 |
| [`re/73-new-game-faction-list.md`](../re/73-new-game-faction-list.md) | `sub_18DC8` 的 `si=98C8h` | 那一則字串沒取出來看 | 靜態 |
| [`re/73-new-game-faction-list.md`](../re/73-new-game-faction-list.md) | 欄位表的「型別」與「屬性」兩個 word | `0x76`／`0x73` 與 `0x0206`／`0x0204` 只由「名字欄 vs 數字欄」的對應推出語意，沒有讀 `sub_1820E` 裡消費它們的那一段 | 靜態 |
| [`re/73-new-game-faction-list.md`](../re/73-new-game-faction-list.md) | `sub_18607` | **算式已讀**（2026-09-02）：`ax = (word_181B4 − 30h) × dl ÷ byte_181BE`，**有餘數就進位**，結果存回 `cs:word_181A4` 再 `sub_18662`；前半段在 `word_181B0 + 20h + word_181A4` 不進位時用 `su… | 靜態 |
| [`re/74-battle-opening-duel.md`](../re/74-battle-opening-duel.md) | `word_1D311 += 6` | 疑似喊話框位置位移，未驗 | 靜態 |
| [`re/74-battle-opening-duel.md`](../re/74-battle-opening-duel.md) | 開場凍結與大將騎出的機器碼形式 | 實機定案（b0–b3），但 `sub_1A1C5` 內部的等待常式怎麼擋住實體更新未逐條讀——見 `spec/80` §3.1 | 靜態 |
| [`re/75-duel-talk-audit.md`](../re/75-duel-talk-audit.md) | 變體 0／2／3／5／6 的臨場抽驗 | 專屬句只在 1／4／7；預設句與它們共用選句機制，公式已 confirmed，抽驗優先度低 | 靜態 |
| [`re/76-d7open-opening-player.md`](../re/76-d7open-opening-player.md) | 第 5 幕來源的列距 | `sub_103EF` 每兩列消耗 720 B（240 × 2 ＋ `add si, 0F0h`），推不出單純的「一列 240 B」；要嘛來源有 stride，要嘛捲動是重疊取樣。沒有畫面驗證 | 靜態 |
| [`re/76-d7open-opening-player.md`](../re/76-d7open-opening-player.md) | `sub_10A70` | 第 5 幕與開場初始化各叫一次，內容未讀 | 靜態 |
| [`re/76-d7open-opening-player.md`](../re/76-d7open-opening-player.md) | 旁白的推進速度 | 逐字之間沒有顯式延遲，靠 `sub_10056` 主迴圈的 `sub_10777` 節拍；換算成秒要先定 `int 61h` 的計時單位 | 靜態 |
| [`re/76-d7open-opening-player.md`](../re/76-d7open-opening-player.md) | PC-98 版的差異 | 4,836 B，尚未進 IDA。旁白應為 Shift-JIS，位址不同 | 靜態 |
| [`re/77-general-affinity-and-flags.md`](../re/77-general-affinity-and-flags.md) | `+0x14` 的語意 | 值域 0–7（在場 381 筆：0 佔 126、2 佔 68），與 `+0x1E` 說話類型交叉表沒有結構。沒有讀取端（§5） | 靜態 |
| [`re/77-general-affinity-and-flags.md`](../re/77-general-affinity-and-flags.md) | `+0x1B` 的語意 | 只有 0（368）與 `0xFF`（13）。沒有讀取端 | 靜態 |
| [`re/77-general-affinity-and-flags.md`](../re/77-general-affinity-and-flags.md) | 旗標 bit 0 | 只出現一次（劇本三的張衛），沒有讀取端 | 靜態 |
| [`re/77-general-affinity-and-flags.md`](../re/77-general-affinity-and-flags.md) | `+0x19` 誰寫非 `0xFF` 值 | **只有劇本作者**：`KI.EXE` 裡只找得到清成 `0xFF` 的寫入端（`sub_15899`）。所以它是純劇本資料，執行期只會消耗不會補充 | 靜態 |

## 2.4 驗收（84 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`playtest/10-event-message-modal.md`](../playtest/10-event-message-modal.md) | 事件 10 producer 仍未定位。 | （未解小節內文） | 靜態 |
| [`playtest/10-event-message-modal.md`](../playtest/10-event-message-modal.md) | 事件 6 #72 的缺失 formatter payload 維持 fail-closed。 | （未解小節內文） | 靜態 |
| [`playtest/10-event-message-modal.md`](../playtest/10-event-message-modal.md) | 原版／remake 同狀態畫面對拍仍是剩餘驗收項。 | （未解小節內文） | 靜態 |
| [`playtest/17-expert-dosbox-remake.md`](../playtest/17-expert-dosbox-remake.md) | 松崗 DOS/V 原版 | **PASS（啟動至開場）** / 2026-08-12 證實空白確認／`0000`／`1234` 均越過密碼頁；完整自然長程驗證尚未執行 | 靜態 |
| [`playtest/21-dosboxx-bridge-sampling.md`](../playtest/21-dosboxx-bridge-sampling.md) | 據點換手之後遮罩會不會跟著變 | `sub_1890A` 的行為，靜態讀得出來，動態沒驗——要打下一座城才看得到 | 靜態 |
| [`playtest/21-dosboxx-bridge-sampling.md`](../playtest/21-dosboxx-bridge-sampling.md) | 松崗 DOS/V 側 | 這套 bridge 還沒在 DOS/V 上跑過。**密碼頁不構成阻礙**（四格留白按「確定」即可通過，`18`）——是還沒做 | 靜態 |
| [`playtest/21-dosboxx-bridge-sampling.md`](../playtest/21-dosboxx-bridge-sampling.md) | 上游授權 | `DOSBox-X-MCP-Debugger` 的原創碼**尚未選定授權條款**（README 明講是刻意留白）。本專案只在本機使用，未再散布 | 實測 |
| [`playtest/24-window-toggles.md`](../playtest/24-window-toggles.md) | 各視窗內部的像素 | 只對過邊線位置，沒有對過內容 | 靜態 |
| [`playtest/24-window-toggles.md`](../playtest/24-window-toggles.md) | 四個視窗**同時**開著的對拍 | 三個視窗開著已做完、三區逐像素 PASS（`38`），系統選單開著也做了（`39`）；**四個一起開的那一張還沒拍** | 靜態 |
| [`playtest/24-window-toggles.md`](../playtest/24-window-toggles.md) | 系統視窗的四個項目 | 存檔／畫面模式／音源／戰略速度，不在這一輪範圍 | 靜態 |
| [`playtest/25-audio-capture-feasibility.md`](../playtest/25-audio-capture-feasibility.md) | **逐曲觸發** | 只錄到開場動畫（`D7OPEN.EXE` 自己會播）。⚠ 這一項**不再擋住任何事**——音檔改由 `tools/bgm2ogg.sh` 離線渲染，不需要在模擬器裡逐首觸發。要當對照組才需要它 | 靜態 |
| [`playtest/25-audio-capture-feasibility.md`](../playtest/25-audio-capture-feasibility.md) | **音效** | 戰術的三個 effect code 已知（`re/17` §3），但沒錄過 | 靜態 |
| [`playtest/25-audio-capture-feasibility.md`](../playtest/25-audio-capture-feasibility.md) | **音源正確性** | DOSBox 用 `sbtype=sb16`／`oplmode=auto` 模擬，與真實硬體的音色差異沒有對照組 | 實測 |
| [`playtest/26-bgm-render-vs-recording.md`](../playtest/26-bgm-render-vs-recording.md) | 音色的聽感 | 頻譜只驗了基頻。諧波結構（也就是「像不像那個音色」）沒有量化比對 | 靜態 |
| [`playtest/26-bgm-render-vs-recording.md`](../playtest/26-bgm-render-vs-recording.md) | 相關係數為什麼不是 0.9 | DOSBox 的 OPL 模擬與這顆的包絡實作不同，加上錄音有系統噪訊。**沒有排除「還有小錯」的可能** | 實測 |
| [`playtest/26-bgm-render-vs-recording.md`](../playtest/26-bgm-render-vs-recording.md) | 其他曲子 | 只有開場曲有錄音對照組。另外 13 首沒有 | 靜態 |
| [`playtest/27-original-video-frame-parity.md`](../playtest/27-original-video-frame-parity.md) | 門強度條的 remake 截圖 | 規則與版面已解並實作（`../re/60` §11、`../spec/32`），但 remake 側沒截到——**這條只亮 20 幀** | 實測 |
| [`playtest/27-original-video-frame-parity.md`](../playtest/27-original-video-frame-parity.md) | 一覽表視窗 | 影片裡有武將／據點／財政的實錄，**還沒量** | 靜態 |
| [`playtest/29-strategy-minimap-markers.md`](../playtest/29-strategy-minimap-markers.md) | 22 勢力的選擇視窗 | 原版點圖例右半格會開一個兩欄的選單（`../re/62` §4.2）。**行為已解、版面未解**，remake 先用「點一下換下一個」代替 | 靜態 |
| [`playtest/29-strategy-minimap-markers.md`](../playtest/29-strategy-minimap-markers.md) | 視野框的美術 | 原版在 `word_10D4C`，尺寸沒從程式碼讀到 | 靜態 |
| [`playtest/29-strategy-minimap-markers.md`](../playtest/29-strategy-minimap-markers.md) | 點地圖區（熱區 `0x16`） | 原版做什麼沒讀 | 靜態 |
| [`playtest/30-ground-planes-implemented.md`](../playtest/30-ground-planes-implemented.md) | 一幀能有幾個兵撞牆 | 原版沒量過。前排寬度決定破牆速度，而破牆速度決定攻城打不打得下來 | 靜態 |
| [`playtest/30-ground-planes-implemented.md`](../playtest/30-ground-planes-implemented.md) | 高平面的橫向移動沒有實測 | 守方站到牆頂的情境還沒跑過 | 實測 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | 勢力一覽的欄位 | 沒有逐欄比對 | 靜態 |
| [`playtest/31-parity-inventory.md`](../playtest/31-parity-inventory.md) | 「委任」那一格的顏色 | 實錄影格上看起來是紅字，但影片是壓縮過的、也沒有機器碼證據。remake 先畫成一般色（`../spec/38`） | 靜態 |
| [`playtest/32-talk-layout-fit.md`](../playtest/32-talk-layout-fit.md) | 變數的實際長度分布 | 這一輪用固定三全形替身。人名多半是 2–3 全形、地名 2–3，但**軍團名與勢力名沒有逐一量過** | 靜態 |
| [`playtest/34-advise-scene-screens.md`](../playtest/34-advise-scene-screens.md) | 逐句節拍 | 原版每句要等玩家按鍵才往下走；remake 直接顯示最新一句 | 靜態 |
| [`playtest/34-advise-scene-screens.md`](../playtest/34-advise-scene-screens.md) | 插圖之外的畫面 | 原版這一頁底下是不是還留著大地圖沒驗過，remake 留著 | 靜態 |
| [`playtest/34-advise-scene-screens.md`](../playtest/34-advise-scene-screens.md) | 選單的反白樣式 | 原版怎麼畫游標列沒解，remake 用自己的反白條 ＋ `>` | 靜態 |
| [`playtest/35-advise-verdict-screens.md`](../playtest/35-advise-verdict-screens.md) | 遷都的畫面 | 沒有截圖。⚠ 目標用一覽表挑**與原版相同**——`sub_17400` 是據點一覽的呼叫端（`../re/26` §4），不是地圖選點 | 實測 |
| [`playtest/36-window-texture.md`](../playtest/36-window-texture.md) | 取用端 | `KI.EXE` 裡哪一段程式把這 128 byte 鋪上去的還沒找到（三條路都排除了）。**排法已經由實機畫面定案**，取用端只影響「還有沒有別的用法」 | 靜態 |
| [`playtest/36-window-texture.md`](../playtest/36-window-texture.md) | 米色視窗 | 一覽表那種米色底原版有沒有紋路沒量過（截圖裡那一片是純色） | 實測 |
| [`playtest/37-main-screen-parity.md`](../playtest/37-main-screen-parity.md) | 換圖塊的那一支機器碼 | 規則是從資料與畫面反推的，每一格都對得上，但**原版在哪裡做這件事**還沒定位（`../re/67` §5） / 找誰在據點換手時改地圖或格子記錄 | 靜態 |
| [`playtest/37-main-screen-parity.md`](../playtest/37-main-screen-parity.md) | 四個視窗**同時**開著時的對拍 | 還沒做。三個視窗開著（`38`）與系統選單開著（`39`）都做了 / 送點擊的方法已解：`click:x,y;press` 成對送，見 `38` §1 | 靜態 |
| [`playtest/37-main-screen-parity.md`](../playtest/37-main-screen-parity.md) | DOSBox 的滑鼠座標 | 視窗 640×480、遊戲 640×400 置中，而 INT 33 把**整個視窗**等比對映到遊戲畫面（送 y 要乘 1.2）。這是本機設定的性質，不是原版的 / 把 `int33 max y` 改成 400 再量一次 | 實測 |
| [`playtest/37-main-screen-parity.md`](../playtest/37-main-screen-parity.md) | 進到大地圖之後的滑鼠座標 | **又換一套**：`sub_120D6` 把 INT 33 的範圍改成 `0..0x17FF × 0..0x101F`（6143×4127 ＝ 整個世界的像素），螢幕座標 ＝ 原始座標 − 鏡頭原點。所以同一個視窗位置在選單裡與在地圖上指到完全不同的地方 / 用 `tools/cursor_probe.py` 在… | 靜態 |
| [`playtest/38-window-parity.md`](../playtest/38-window-parity.md) | 天候物件 | 原版跑了 10 天才截到，remake 停在第 1 天。**這是狀態差** / 要對就得讓兩邊同一天——用存檔定位（`../spec/90` §2） | 靜態 |
| [`playtest/39-system-window-parity.md`](../playtest/39-system-window-parity.md) | 「液晶」畫面模式 | 原版的畫面模式有兩個選項，對應 `GAMEPAL.BRG` 的 bank 0–3 與 4–7（`../re/55` §4）。remake 只做了 16 色那一組 / 載 bank 4–7 再對拍一次 | 靜態 |
| [`playtest/39-system-window-parity.md`](../playtest/39-system-window-parity.md) | 音效的 TYPE 2/3/4 | 原版有四種音源型別，remake 只有開／關 / 看 `sub_102D0` 那四型的差別 | 靜態 |
| [`playtest/39-system-window-parity.md`](../playtest/39-system-window-parity.md) | 日期對不上 | 原版跑到 4月9日才截到 / 要嘛用存檔定位，要嘛加一個「跑到指定日期」的驗收旗標 | 靜態 |
| [`playtest/40-tactical-parity.md`](../playtest/40-tactical-parity.md) | `sb-enemy` 的 44 px | 兩條都頂在上限，原版那一格已經打了 20 秒（§10） / 要對就得讓兩邊的**時刻**對齊，不是改算式 | 靜態 |
| [`playtest/40-tactical-parity.md`](../playtest/40-tactical-parity.md) | 右下那一小塊地形色（8 px，(361..364, 192..194)） | `field` 區剩下唯一沒歸類的一群 / 拿 §13.1 那一招再跑一次：把 192 個子圖塊逐張套進那個顯示格比一次。**假說**是同一類（另一個實體破損的時刻不同），但還沒驗 | 靜態 |
| [`playtest/40-tactical-parity.md`](../playtest/40-tactical-parity.md) | `sb-enemy`／`sb-self` 1.5% | 兩格將旗的內容 / — | 靜態 |
| [`playtest/40-tactical-parity.md`](../playtest/40-tactical-parity.md) | `sub_1DFBB` 的快路徑 | remake 一律走合成。兩條路在全畫面重繪下應該畫出同樣的像素（`../spec/58` §4），但沒有逐格驗過 / — | 靜態 |
| [`playtest/40-tactical-parity.md`](../playtest/40-tactical-parity.md) | unit 0 的第二趟 | 深度迴圈跑完後 `dl & 0x20` 成立時會對五個鄰格各跑一次 `ax = 0`；**觸發條件（旗標 bit 5）誰設還沒解** / 掃誰對顯示格的 `+0` 寫 `0x20` | 靜態 |
| [`playtest/41-m7-corrected-text-on-screen.md`](../playtest/41-m7-corrected-text-on-screen.md) | 原版側的同狀態對照 | 這一份只驗 remake 自己「有沒有溢出」。**原版同一則長什麼樣沒有並排比過**——要用 `-open-talk-index` 對應的原版操作序列，還沒做 | 靜態 |
| [`playtest/41-m7-corrected-text-on-screen.md`](../playtest/41-m7-corrected-text-on-screen.md) | 變數的實際長度分布 | 截圖用的是實際遊戲值（如「袁胤」兩字），而 `TestAllTalkLinesFitTheirBox` 用三全形替身。**軍團名與勢力名的長端沒有逐一量過**（`32`） | 實測 |
| [`playtest/42-window-parity.md`](../playtest/42-window-parity.md) | 進言五項選單的原版截圖 | §5 的輸入模型限制 / 需要能送「瞬時 click」的擷取動作（縮短按住時間），或改用鍵盤路徑（未驗證原版是否支援） | 實測 |
| [`playtest/43-field-battle-parity.md`](../playtest/43-field-battle-parity.md) | 遭遇訊息畫面的對拍 | 「遇上兵馬了」訊息與 remake 的遭遇戰選單版面沒有比 / 原版是 TALK 訊息框，remake 是自製選單——先讀原版遭遇後的選擇 UI 是什麼樣（影片 `parity-field13/enc.mp4` 15 秒附近有素材） | 靜態 |
| [`playtest/43-field-battle-parity.md`](../playtest/43-field-battle-parity.md) | 佔用圖快取欄的讀檔重建 | §3 是強證據不是 confirmed / 讀原版的讀檔常式（`sub_18CAE` 一帶）確認重建走哪個欄位 | 靜態 |
| [`playtest/45-ai-longrun-comparison.md`](../playtest/45-ai-longrun-comparison.md) | 原版側只有一次跑、三個資料點 | 存檔輪換的點擊序列要修（m1–m4 沒重存：疑似覆寫確認框的 YES 座標沒點中）；修好後可加拍幾輪拉出分佈 | 靜態 |
| [`playtest/45-ai-longrun-comparison.md`](../playtest/45-ai-longrun-comparison.md) | remake 孫策的擴張比原版猛（一個月滅劉繇 vs 原版半年拿兩城） | 幅度差在戰鬥節奏層（編成／行軍／攻城的耗時），不是決策層；要縮就得對拍單場攻城的時間線 | 靜態 |
| [`playtest/46-runtime-language-switch.md`](../playtest/46-runtime-language-switch.md) | Android 實機／模擬器沒實地切過 | 這裡的手機畫面是桌面 Xvfb 跑 `cmd/wlandroid` 拍的（同一份 `internal/ui/phone`）；實機驗收排在下一次 Android 打包 | 靜態 |
| [`playtest/47-latin-screen-survey.md`](../playtest/47-latin-screen-survey.md) | 戰術畫面的截圖不可重現 | 找出 `-shot-frames` 之下仍隨牆鐘走的那一段（動畫幀？音訊回呼？），改成照 tick 推進 | 實測 |
| [`playtest/47-latin-screen-survey.md`](../playtest/47-latin-screen-survey.md) | 半形語系的戰場標題沒有地名 | 見 `docs/spec/87` §9 | 靜態 |
| [`playtest/48-display-polish.md`](../playtest/48-display-polish.md) | 對拍沒有開闊地的 fixture | `playtest/40` 量的兩個局面都在城壁邊；要擋住這一類回歸得再加一個開闊地的 fixture | 靜態 |
| [`playtest/48-display-polish.md`](../playtest/48-display-polish.md) | 事件列本身是 remake 自創 | 原版怎麼提示月結（如果有）沒查過 | 靜態 |
| [`playtest/49-parity-retest-20260827.md`](../playtest/49-parity-retest-20260827.md) | 40` §8 的取樣點是**第 61 步**，當時 `field` 0.84%，再修兩處之後 0.17%。 | （未解小節內文） | 靜態 |
| [`playtest/49-parity-retest-20260827.md`](../playtest/49-parity-retest-20260827.md) | 修正前 | 45.55% / 45.18% / 45.83% | 靜態 |
| [`playtest/49-parity-retest-20260827.md`](../playtest/49-parity-retest-20260827.md) | 量表那一塊 | 0.13% / 0.09% / 0.06% | 靜態 |
| [`playtest/49-parity-retest-20260827.md`](../playtest/49-parity-retest-20260827.md) | `field` | 0.86% / 0.81% / 0.84% | 靜態 |
| [`playtest/49-parity-retest-20260827.md`](../playtest/49-parity-retest-20260827.md) | 61 | X 58..60（出發線） / X 3..5 | 靜態 |
| [`playtest/49-parity-retest-20260827.md`](../playtest/49-parity-retest-20260827.md) | 120 | X 50..57 / X 3..5 | 靜態 |
| [`playtest/49-parity-retest-20260827.md`](../playtest/49-parity-retest-20260827.md) | 300 | X 8..43（已越過城門） / X 3..5 | 靜態 |
| [`playtest/49-parity-retest-20260827.md`](../playtest/49-parity-retest-20260827.md) | 600 | X 5..60（散開） / X 3..5 | 靜態 |
| [`playtest/49-parity-retest-20260827.md`](../playtest/49-parity-retest-20260827.md) | 攻城對拍的取樣點 | 兩個缺口修好之後 `field` 0.86%，剩下的是局面不等價（§3.3） / 要回到 0 px 得讓原版與 remake 開出**同一場**（同一批武將、同一個 tick）。條件寫進 `../spec/91`，不要再寫死步數 | 靜態 |
| [`playtest/49-parity-retest-20260827.md`](../playtest/49-parity-retest-20260827.md) | 系統選單的兩列 remake 差異 | `map` 7.87%——主君編成與損害報告把視窗撐高 48 px（§2） / 這是**已標記的 remake 差異**，不是缺口。要對到 0 px 得有一個「照原版六列」的驗收模式 | 靜態 |
| [`playtest/49-parity-retest-20260827.md`](../playtest/49-parity-retest-20260827.md) | 原版側沒有重跑 | 這一輪用的是 08-16／17／24 的存檔擷取 / 要重跑得先建 `wolong-dosboxx`（`docker/dosboxx/Dockerfile`） | 靜態 |
| [`playtest/50-corps-on-map.md`](../playtest/50-corps-on-map.md) | 首都疊圖為什麼不畫 | 行為定案，機制未讀（§4） | 靜態 |
| [`playtest/50-corps-on-map.md`](../playtest/50-corps-on-map.md) | 別的疊圖組合 | 只驗過「首都 ＋ 軍團」。災害物件 ＋ 軍團、非首都據點 ＋ 軍團都沒有樣本 | 靜態 |
| [`playtest/50-corps-on-map.md`](../playtest/50-corps-on-map.md) | 「軍團在路上」對原版 | 要新的原版擷取：編成 → 行軍指示 → 等幾天。滑鼠腳本點得到行軍那一列之後才做得起來（`40` §1.2） | 靜態 |
| [`playtest/51-siege-deadlock.md`](../playtest/51-siege-deadlock.md) | 三個讀數仍對不上同一幀 | §2。兩軍會打了，但這條 fixture 的守方是**玩家側而玩家沒下令**，所以站著挨打——攻方只折損 1 點。⚠ **不是委任**：委任那條路走 `combat.Resolve` 抽象判定，根本不開戰場，而 e10 是戰場畫面 | 靜態 |
| [`playtest/51-siege-deadlock.md`](../playtest/51-siege-deadlock.md) | 原版擺兵寫的是哪一個表 | `spec/95` 是靠內部一致性定的，不是照抄。原版 `[si+0Ah]` 的來源沒讀出來 | 靜態 |
| [`playtest/51-siege-deadlock.md`](../playtest/51-siege-deadlock.md) | 攻方大將體力歸零 | `drainSiegeGeneral` 扣到 0，而戰鬥傷害最低留 1（`../re/11` §5.16）。攻城計時器要不要也留 1，沒讀出來 | 靜態 |
| [`playtest/52-siege-timeseries-parity.md`](../playtest/52-siege-timeseries-parity.md) | ⭐ **破牆速度** | remake 量得很乾淨（§5）：三個兵貼著打、每幀約 2 點、第 ~1,100 幀破。**原版那一側量不出來**——體力同時被攻城計時器與挨打扣，單張畫面分不開 / ⛔ **要新的原版擷取**：從存檔載入之後在開戰的**前兩百幀密集取樣**（現有的 `wait:6` 一格就跨掉五十幾幀）。這是這一項唯一的路，靜… | 靜態 |
| [`playtest/52-siege-timeseries-parity.md`](../playtest/52-siege-timeseries-parity.md) | ⛔ 試過但**不是**成因：不翻轉 `gateX` | `sub_1CAEB`（`0001CB0F`）寫進立即值的是**未翻轉**的索引第二欄，而 remake 對它套了 `RotateGateX`。實驗把翻轉拿掉（`GateX` 從 27 變回 36）跑同一條 fixture：**城壁照樣被磨穿**（1,660 → 0），門一點傷都沒有，結束時間只從 f1116 變… | 靜態 |
| [`playtest/52-siege-timeseries-parity.md`](../playtest/52-siege-timeseries-parity.md) | 索引第二欄到底是什麼 | `../re/11` §4.5 稱它「城門附近的 X」，§5.8i 定案為「命令 3 的目標 X」。但這張圖上它落在城壁本體、離最近的門有 3–7 格 / 掃 186 張攻城圖：第二欄那一欄的圖塊值是什麼、與門的距離分布 | 靜態 |
| [`playtest/52-siege-timeseries-parity.md`](../playtest/52-siege-timeseries-parity.md) | 守方為什麼在 t6 掉那麼快 | 原版守方大將體力兩段內 124 → 37 / 要先讓攻方以原版的方式進城，才比得了 | 靜態 |
| [`playtest/56-lubu-flow-parity.md`](../playtest/56-lubu-flow-parity.md) | 原版的行軍目的地一覽與三選一 | **「軍團」彈出選單的第二列「行軍指示」點不到。** 四輪都停在第一列「位置確認」：`tap:25,10,5`、`tap:25,9,5`、`click:25,10;press` 三種送法都一樣。這與 `docs/playtest/42` §5 記的是同一類限制——`playtest/54` 證明的是「**能把選單… | 靜態 |
| [`playtest/56-lubu-flow-parity.md`](../playtest/56-lubu-flow-parity.md) | 原版的攻城結算 | 卡在上一列。remake 側在 196/4/8「張遼 對 城兵　攻方勝　兵力 1000→960／910→60　據點損害 54　攻下 譙」 | 靜態 |
| [`playtest/56-lubu-flow-parity.md`](../playtest/56-lubu-flow-parity.md) | 捲軸滑塊的 1 px | 只有一個取樣點（`../spec/107` §7） | 靜態 |
| [`playtest/56-lubu-flow-parity.md`](../playtest/56-lubu-flow-parity.md) | 「軍團」彈出選單的逐像素對拍 | 位置與兩列文字對得上，**框本身沒有比過**——原版那張參考圖是另一個局面（4/20、讀了存檔、游標停在第一列） | 靜態 |
| [`playtest/56-lubu-flow-parity.md`](../playtest/56-lubu-flow-parity.md) | 攻城**戰場**（不是結算） | 兩邊都一樣：空城攻城是自動判定，不進戰術畫面（`internal/state/corps.go` 的 `fightGarrison`；原版 `sub_14ED7` 的 `cmp bx, 4200h`）。要看到戰場得等守方有軍團駐守，而那一刻兩邊不會同時發生 | 靜態 |
| [`playtest/57-appimage-20260902-reverify.md`](../playtest/57-appimage-20260902-reverify.md) | 沿用 `56` §6，沒有新增：軍團選單的同狀態原版參考、原版的行軍目的地一覽、 | （未解小節內文） | 靜態 |

## 2.5 外部資料（6 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`reference/02-jp-cht-diff.md`](../reference/02-jp-cht-diff.md) | 逐句對照的**日文原版側畫面** | 兩版的同一則並排截圖沒有做過——這一輪的畫面抽樣只驗 remake 自己（`../playtest/41` §6） | 實測 |
| [`reference/02-jp-cht-diff.md`](../reference/02-jp-cht-diff.md) | `#223` 等訊息的欄位完整語意 | 只修已證實的標記編號，欄位語意仍未解（§9） | 靜態 |
| [`reference/03-baked-japanese.md`](../reference/03-baked-japanese.md) | 橫幅上寫的是**「臥竜伝」**——日文漢字，不是「臥龍傳」。 | （未解小節內文） | 靜態 |
| [`reference/04-first-survey.md`](../reference/04-first-survey.md) | `SHOW.O` | 57,148 B / 被 `INSTALL.EXE` 與 `LOGO.EXE` 引用。開頭 `3c df 00 00 11 af 01 00 50 00 80 07`。**未解** | 靜態 |
| [`reference/04-first-survey.md`](../reference/04-first-survey.md) | 不要憑「同一份專案應該用同一個編譯器」外推——**`KI.EXE` 的編譯器未解。 | （散句） | 靜態 |
| [`reference/05-eten-font-provenance.md`](../reference/05-eten-font-provenance.md) | `END_S13/S14/S15` 是中文版加的結局段 | S13／S14 是字型。**`END_S15` 仍未解** | 靜態 |

## 2.6 其他（203 條）

| 出處 | 缺口 | 現況 | 裁決 |
|---|---|---|---|
| [`docs/mobile/android-plan.md`](../mobile/android-plan.md) | 實機驗收 | ⛔ 沒有裝置。里程碑 H 保持未完成 | 靜態 |
| [`docs/mobile/android-plan.md`](../mobile/android-plan.md) | SAF 匯入的複製流程 | 入口做完了，但「選資料夾 → 複製 69 檔」沒有自動驗過：要驅動系統的檔案選擇器。smoke 走的是 `adb` ＋ `run-as`，那是驗收路徑不是玩家路徑 | 靜態 |
| [`docs/mobile/android-plan.md`](../mobile/android-plan.md) | 高 DPI 下的點陣字 | 定為 2 倍（`../spec/100`）；實機上要不要另一檔未驗 | 靜態 |
| [`docs/mobile/android-plan.md`](../mobile/android-plan.md) | release signing | keystore 怎麼保管還沒決定；目前出的是 debug 簽章 | 靜態 |
| [`docs/mobile/android-plan.md`](../mobile/android-plan.md) | 16 KB 對齊只驗到建置那一層 | `readelf` 確認 LOAD 段是 `0x4000`，但**沒有 16 KB page size 的裝置或 AVD 實際載過**。這一條與「實機驗收」是同一個缺口 | 靜態 |
| [`docs/mobile/android-ux.md`](../mobile/android-ux.md) | 點陣字在高 DPI 上要放大幾倍 | 定為 2（§6）；平板或小手機要不要另一檔，得有實機回饋 | 靜態 |
| [`docs/mobile/android-ux.md`](../mobile/android-ux.md) | 小卡要放哪些欄位 | 目前放名稱／歸屬／生產力／防災／城兵五項。原版一覽表的欄位全表在 `docs/spec/38`，還沒逐項比對過取捨 | 靜態 |
| [`docs/mobile/android-ux.md`](../mobile/android-ux.md) | 縮放的下限 | 整張大地圖 384×256 格全塞進手機會小到看不見，最小縮放級距還沒定 | 靜態 |
| [`docs/mobile/android-ux.md`](../mobile/android-ux.md) | 戰場的縮放 | 目前固定 1×（原版的 480×368 剛好塞進主區）。放大之後看得到的格子會變少，**那會改變決策**，所以沒有做 | 靜態 |
| [`docs/mobile/android-ux.md`](../mobile/android-ux.md) | 戰場縮圖的點選 | 原版點縮圖可以移動鏡頭，手機版目前只顯示 | 靜態 |
| [`docs/mobile/android-ux.md`](../mobile/android-ux.md) | 事件訊息的停留時間 | 六秒是估的。原版沒有這個機制（它要按鍵才消），所以沒有可抄的數字 | 靜態 |
| [`promo/android.md`](../promo/android.md) | 實機錄影 | ⛔ 沒有裝置。片中畫面出自桌面的同一份 `internal/ui/phone`，與 APK 是同一份程式碼，但**不是實機錄影** | 靜態 |
| [`promo/android.md`](../promo/android.md) | 模擬器錄影 | 模擬器在這台機器上只有個位數 fps，錄出來會頓 | 靜態 |
| [`promo/combined.md`](../promo/combined.md) | 三段的響度沒有逐段量 | 配樂那一軌過了 `loudnorm=I=-18`，但三段畫面本身沒有音效，所以整片只有一條音軌——目前沒有段間響度落差的來源。**換成保留各段原音的版本時要重量** | 靜態 |
| [`promo/dosv-realmachine.md`](../promo/dosv-realmachine.md) | 原版戰術戰場的實機擷取 | 四次未觸發（§4）。要嘛接受原版 RNG 的變異多跑幾次，要嘛從存檔直接進戰場——後者要先解「怎麼從存檔載入到開戰的那一刻」 | 靜態 |
| [`promo/dosv-realmachine.md`](../promo/dosv-realmachine.md) | 原版 AdLib 的同場錄音 | `ctrl+F6` 的 WAV 擷取這次沒生效，配樂沿用 2026-08-12 那次的實錄 | 靜態 |
| [`promo/dosv-realmachine.md`](../promo/dosv-realmachine.md) | 兩側時鐘速度可比 | remake 用最高速檔才看得到動靜；要真的可比，得先量原版預設檔的每日實時秒數 | 靜態 |
| [`release/01-cross-build-gate.md`](../release/01-cross-build-gate.md) | 目標 OS 實跑 | **做不到**：這台是 Linux，沒有 Mac／Windows。檔頭驗過（PE32+／Mach-O），但視窗、輸入、音訊、字型載入都沒有在目標系統上跑過 | 實測 |
| [`release/01-cross-build-gate.md`](../release/01-cross-build-gate.md) | linux/arm64 的本體 | 要在 arm64 的 Linux 上建（Ebiten 的 cgo 沒有交叉工具鏈） | 靜態 |
| [`release/01-cross-build-gate.md`](../release/01-cross-build-gate.md) | Windows 的 smoke | 同第一項 | 靜態 |
| [`release/10-full-20260828.md`](../release/10-full-20260828.md) | Windows／macOS 原生 GUI 實機驗收 | ⛔ 沒有那兩個平台的機器。交叉建置的產物只驗了檔頭 | 靜態 |
| [`release/10-full-20260828.md`](../release/10-full-20260828.md) | Android 實機驗收 ＋ release signing | ⛔ 沒有裝置；目前是 debug 簽章 | 靜態 |
| [`release/10-full-20260828.md`](../release/10-full-20260828.md) | Android 模擬器 smoke | 這一批跳過（§3）。要跑得先重建 13.7 GB 的模擬器映像 | 靜態 |
| [`release/10-full-20260828.md`](../release/10-full-20260828.md) | 沒有音效裝置時遊戲會掛 | 完整版會自己找到音檔並開音訊，而 Ebiten 沒有可查詢的音訊 API（`../spec/75` §5） | 靜態 |
| [`release/11-full-20260830.md`](../release/11-full-20260830.md) | Windows／macOS 原生 GUI 實機驗收 | 沒有 Mac／Windows 主機。交叉建置的檔頭與 ABI 已驗，但不以檔頭取代實測 | 實測 |
| [`release/11-full-20260830.md`](../release/11-full-20260830.md) | Android 實機／模擬器 smoke | 沒有裝置；模擬器映像不在機器上，重建要重下 13.7 GB | 靜態 |
| [`release/11-full-20260830.md`](../release/11-full-20260830.md) | 合成片的段間響度 | 目前整片只有一條配樂軌，沒有段間落差的來源；換成保留各段原音的版本時要重量（`../promo/combined.md` §6） | 靜態 |
| [`release/12-full-20260902.md`](../release/12-full-20260902.md) | Windows／macOS 原生 GUI | 沒有實機或 VM，只驗了檔頭 | 靜態 |
| [`release/12-full-20260902.md`](../release/12-full-20260902.md) | Android 實機與 release signing | 沒有裝置；金鑰保管方式未定 | 靜態 |
| [`release/12-full-20260902.md`](../release/12-full-20260902.md) | Android 模擬器 smoke | `wolong-android-emulator` 映像不在這台機器上 | 靜態 |
| [`release/12-full-20260902.md`](../release/12-full-20260902.md) | 可散布批次 | 仍停在 `20260830`，這一輪的六項修正還沒進公開發行 | 靜態 |
| [`release/README-RELEASE.md`](../release/README-RELEASE.md) | Windows／macOS 原生 GUI | 交叉建置的產物只驗了檔頭，沒有在目標作業系統跑過。M8 唯一的閘 | 靜態 |
| [`release/README-RELEASE.md`](../release/README-RELEASE.md) | Android 實機驗收 | 只有 Docker 模擬器；觸控手感、真實 GPU、高 DPI 上的點陣字可讀性都驗不到 | 靜態 |
| [`release/README-RELEASE.md`](../release/README-RELEASE.md) | Android 正式簽章 | 出的是 debug 簽章，keystore 怎麼保管還沒決定 | 靜態 |
| [`release/README-RELEASE.md`](../release/README-RELEASE.md) | 16 KB page size 裝置 | `.so` 的 LOAD 段已是 `0x4000`，但沒有那種裝置或 AVD 實際載過 | 靜態 |
| [`spec/10-city-tick.md`](../spec/10-city-tick.md) | 據點換手之後 `+0x00` 低 4 位會不會跟著變 | `sub_1890A` 靜態讀過，動態沒驗——要打下一座城才看得到 | 靜態 |
| [`spec/10-city-tick.md`](../spec/10-city-tick.md) | 玩家據點求援的喇叭聲（`sub_10CDE`） | 呈現層未接 | 靜態 |
| [`spec/100-phone-text-scale.md`](../spec/100-phone-text-scale.md) | 倍率不能在遊戲內調 | 固定 2。平板或小手機可能要 1 或 3，得先有實機回饋 | 靜態 |
| [`spec/101-phone-glyph-scale2x.md`](../spec/101-phone-glyph-scale2x.md) | 高 DPI 上 Ebiten 把 960×540 再縮到螢幕（非整數倍）的殘餘鋸齒 | 那一層是 Ebiten 的螢幕縮放（預設線性濾波），本規格只處理字模那一層。要不要改成整數倍畫布得有實機回饋 | 靜態 |
| [`spec/102-battle-fast-forward.md`](../spec/102-battle-fast-forward.md) | 快轉時戰場區的底紋 | 原版是**藍底綠線的菱形格**（`playtest/53` 的裁切），不是龍紋。已收斂到「讀哪裡」：`sub_1DC9D` 的 `es` ＝ `word_1E15E`，**唯一的寫入端**是 `sub_1D958`（`sub_19946` 用 `cx = word_1D2F6` 傳進去），而 `word_1D2F… | 靜態 |
| [`spec/103-phone-diplomacy-amount-keypad.md`](../spec/103-phone-diplomacy-amount-keypad.md) | 撥款請求的「指定金額」 | 原版 `sub_17C6E` 的另一個呼叫端；手機版仍只給「照要求撥款／拒絕」，同一套鍵盤可以接，先不做 | 靜態 |
| [`spec/104-advisor-naming-window.md`](../spec/104-advisor-naming-window.md) | 六格與本頁字的顏色 | 六格 15、游標底線 15／1、字 9：`sub_19223`／`sub_1928A` 的屬性值直讀，但沒逐像素對過 | 靜態 |
| [`spec/105-encounter-goes-straight-to-battle.md`](../spec/105-encounter-goes-straight-to-battle.md) | 遭遇當天的日期差一天 | **量到剩 2 個子刻**（§6）。原本記的「時鐘推進速率」不是成因——**接觸在第幾個子刻與速度檔無關**，節流只改牆鐘秒數不改 tick 數。缺的是原版接觸 tick 的一手數字 | 靜態 |
| [`spec/106-message-box-reporter-portrait.md`](../spec/106-message-box-reporter-portrait.md) | `0x94`／`0x95` 兩頁保留肖像的用途 | `0x94` 是一張紅臉武將、`0x95` 是空白。沒找到傳這兩個值的呼叫點 | 靜態 |
| [`spec/106-message-box-reporter-portrait.md`](../spec/106-message-box-reporter-portrait.md) | #217 的兩個 `{3}` | 機制已通（`SeqFactions`），但**發那一則的呼叫端還沒讀**，所以第二個勢力是誰未定 | 靜態 |
| [`spec/107-launcher-ui-colours.md`](../spec/107-launcher-ui-colours.md) | 殼層其餘幾頁（ＹＥＳ／ＮＯ、劇本、四槽讀檔）的配色 | 這一輪只對過勢力清單與君主卡。其餘幾頁同樣走 `paletteInk`，修完應該一起好，但**沒有逐像素比過** | 靜態 |
| [`spec/107-launcher-ui-colours.md`](../spec/107-launcher-ui-colours.md) | 捲軸滑塊的位置差 1 px | 量到：22 筆、`top` ＝ 4 時，原版的綠面在 y 161–216，remake 在 160–215（高度都是 56）。`⌊128×4/22⌋ ＝ 23` 給 159，原版對應的是 24。**只有這一個取樣點**，分不出是無條件進位、四捨五入還是槽的起點差 1；`38` §1.6 的實機量測只釘住高度沒釘位… | 靜態 |
| [`spec/107-launcher-ui-colours.md`](../spec/107-launcher-ui-colours.md) | 反白條 | remake 的鍵盤游標，原版沒有。要不要照戰略層一覽表那樣「碰過才畫」（`g.listTouched`）沒有定案 | 靜態 |
| [`spec/11-ai-sortie.md`](../spec/11-ai-sortie.md) | `資金高位 >= 0x80` 那一支 | `cmp bh, 80h / jnb` 會直接算「答應」，等於資金超過約 840 萬時門檻失效。**看起來像有號數的邊界處理**，未逐位對過 | 靜態 |
| [`spec/11-ai-sortie.md`](../spec/11-ai-sortie.md) | 君主出陣之後的行為 | 那支軍團跟一般軍團有沒有差別，未讀 | 靜態 |
| [`spec/110-corps-command-menu.md`](../spec/110-corps-command-menu.md) | 選單本身的逐像素對拍 | 原版選單拍到了（`parity-tap5/menu.png`），**但沒有做逐像素比對**——remake 的彈出選單走 `drawLegacyChoiceBox`，位置與底色未量 | 靜態 |
| [`spec/110-corps-command-menu.md`](../spec/110-corps-command-menu.md) | 「據點」指令同樣是兩項選單 | `../re/22` §3.4：TALK #82「首都確認／據點一覽」。**這一輪沒動**，remake 的「據點」目前直接開一覽 ＋ 情報卡 | 靜態 |
| [`spec/112-cursor-idle-resume-delay.md`](../spec/112-cursor-idle-resume-delay.md) | `sub_20000(ax=2)` 為什麼在恢復時被呼叫 | 那是滑鼠層跳表的第 2 號（`sub_2009A`，回傳快取座標並清旗標），在這個時點叫它的理由沒讀 | 靜態 |
| [`spec/112-cursor-idle-resume-delay.md`](../spec/112-cursor-idle-resume-delay.md) | 訊息框那一路的 `8` 與游標的 `0Ch` 為什麼不同 | 兩個立即值都是 confirmed，但差別的用意沒有證據；remake 兩邊都用同一個延遲 | 靜態 |
| [`spec/113-rle-length-header.md`](../spec/113-rle-length-header.md) | `BATTLE.MAP`／`MMAP.MCH`／`BATTLE.MDL` 走哪一支載入器 | 沒查（`../formats/06` §6）。它們的前 4 byte 不是長度，所以**至少不是這一族** | 靜態 |
| [`spec/114-general-affinity.md`](../spec/114-general-affinity.md) | 隨機投靠那一條 | `sub_15899` 的 `+0x19 == 0xFF` 分支還沒接（§3 的第二個 ⚠）。要先把勢力的武將數與據點數接進規則層 | 靜態 |
| [`spec/114-general-affinity.md`](../spec/114-general-affinity.md) | 出仕的畫面通知 | 原版對象是玩家勢力時跳訊息 `0x29`；這一版只改狀態，不排事件 | 靜態 |
| [`spec/114-general-affinity.md`](../spec/114-general-affinity.md) | 旗標 bit 5 之外的退場條件 | `sub_15899` 只在「心向的勢力已滅」時看 bit 5；bit 5 沒設的武將會留在原地等下一輪，這一點沒有實機驗證 | 靜態 |
| [`spec/12-strategy-chrome.md`](../spec/12-strategy-chrome.md) | 樣式碼的值域 | 只確定 `0`＝擦除、`0x0B`＝命令、`0x0Bh`／`0x10h`／`0x15h`／`0x1Fh` 各自出現在哪個視窗已知，完整值域未列 | 靜態 |
| [`spec/13-main-window-toggles.md`](../spec/13-main-window-toggles.md) | 對得上（`docs/playtest/24`）。 原版執行期的開關行為仍未驗。 | （散句） | 靜態 |
| [`spec/13-main-window-toggles.md`](../spec/13-main-window-toggles.md) | 熱區 5 | 原版登記了但不接任何常式，remake 照樣不做事 | 靜態 |
| [`spec/13-main-window-toggles.md`](../spec/13-main-window-toggles.md) | **remake 沒有邊緣捲動** | **機制早就解了**（`../re/47` §6）：`sub_120D6` 進大地圖時把 INT 33 的範圍換成**整個世界**（水平 0–`17FFh` ＝ 384 格 × 16、垂直 0–`101Fh`），`sub_11F7F` 再把原始座標減掉鏡頭原點、夾在 0–639／0–399，**夾掉的量同時加回鏡… | 靜態 |
| [`spec/13-main-window-toggles.md`](../spec/13-main-window-toggles.md) | 六列的語意 | handler 讀到四支（`docs/re/55` §5）：畫面模式換調色盤組、音效走驅動、戰略速度只存值、戰術速度存值 ×16。**「資料儲存」與「遊戲結束」那兩支沒讀**，remake 照標籤字面接（開四槽視窗／走 ＹＥＳ／ＮＯ 確認） | 靜態 |
| [`spec/13-main-window-toggles.md`](../spec/13-main-window-toggles.md) | 「畫面模式」 | 兩個選項是「１６色」與「 液晶 」，切的是 `GAMEPAL.BRG` 的 bank 0–3 ↔ 4–7（`docs/re/02` §6.2）。**remake 只做第 0 組**，這一格固定顯示「１６色」——液晶那組是給 8 階調液晶的高飽和純色，現代螢幕沒有對照物 | 靜態 |
| [`spec/13-main-window-toggles.md`](../spec/13-main-window-toggles.md) | 「音效」 | 值由 `g.soundValue()` 填。原版五個選項是 ＯＦＦ／TYPE 1–4（音源型別），remake 只有開／關 | 靜態 |
| [`spec/13-main-window-toggles.md`](../spec/13-main-window-toggles.md) | 戰場內調速度 | 戰場獨佔輸入，所以 `updateBattle` 自己接一次 ＋／−（調戰術速度），調完浮一行 1.5 秒的提示。**原版戰場沒有速度指示**，常駐顯示會破壞版面 parity | 靜態 |
| [`spec/20-save-format.md`](../spec/20-save-format.md) | 存檔區塊的 7 KB 未解區 | `+0x1EC0`–`+0x42C0`，靠 `raw` 原樣保存，但**內容仍不知道**（`docs/formats/08`） | 靜態 |
| [`spec/20-save-format.md`](../spec/20-save-format.md) | 原版 `SAVE.DAT` 的槽位語意 | 四個槽與 `SINARIO.DAT` 的四個劇本是不是同一個編號空間，未確認 | 靜態 |
| [`spec/21-corps-formation-reserves.md`](../spec/21-corps-formation-reserves.md) | 編成畫面的兵種切換 | remake 由呼叫端直接給 `kinds`，沒有原版那個「點一下 +1 → 全退回池 → 重跑分配」的迴圈（`sub_16C92`）。這是 UI 層的差異，不影響分配式 | 靜態 |
| [`spec/21-corps-formation-reserves.md`](../spec/21-corps-formation-reserves.md) | 池的上限 | `sub_155EC` 的 `0xFFDC` 只在退兵路徑上驗過；月結加兵是不是同一支未查。**remake 兩條路徑現在都夾**（`economy.ClampReserve`），但那是照著同一個常數做的，不是證明原版共用同一支 | 靜態 |
| [`spec/22-corps-formation-window.md`](../spec/22-corps-formation-window.md) | 頭像的邊框 | `sub_107D2` 只 blit 64×64 的圖塊，**框在哪裡畫的沒找到**——場景 5 的 op 清單裡沒有頭像那一格的框 | 靜態 |
| [`spec/22-corps-formation-window.md`](../spec/22-corps-formation-window.md) | 兵種標籤 | 畫面用場景 5 的「主將」，規則層的 `army.Position` 第一個是「大將」（原版 TALK #62 也這樣說）。兩處用語不同是原版就有的，不要統一 | 靜態 |
| [`spec/24-corps-info-window.md`](../spec/24-corps-info-window.md) | 指令流程的**入口**與 remake 不同 | `sub_17FDB` 已解（`../re/45` §1：選據點 → 選「戰鬥指揮／委任／解體」→ 寫 `+0x00` 位元 2、`+0x0B`、`+0x20`）。⚠ 先前這一列寫「remake 沒有那三個選項」——**那是錯的**：`cmd/wlgame/marchmode.go` 三個選項都在（`docs/s… | 靜態 |
| [`spec/25-slot-select-window.md`](../spec/25-slot-select-window.md) | 空槽標記 | 原版用名稱欄第一個字 `0xD0A1`；remake 用「載得起來且玩家勢力有效」判定，兩者不等價 | 靜態 |
| [`spec/25-slot-select-window.md`](../spec/25-slot-select-window.md) | 新遊戲共用 | remake 的啟動殼層是自己的畫面，還沒有換成這個四槽視窗 | 靜態 |
| [`spec/26-yes-no-dialog.md`](../spec/26-yes-no-dialog.md) | 原版的使用者 | `sub_18DC8` 只有一個呼叫端 `sub_11AC3`（新遊戲流程），問題文字由那裡給，內容未讀 | 靜態 |
| [`spec/26-yes-no-dialog.md`](../spec/26-yes-no-dialog.md) | `cx = 600Dh` 的尺寸編碼 | `sub_19796`／`sub_197C3` 是**保存／還原被蓋住的畫面**，`dx`／`bx` 是像素座標、換算成 VRAM 位址（`45` §2 逐行解過）。這個呼叫端的 `cx` 高低位元組怎麼對到寬高沒逐位對過 | 靜態 |
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
| [`spec/30-victory.md`](../spec/30-victory.md) | 四個劇本的結局是否不同 | 十二幕依序播，**沒有依劇本分支的證據**（`D7END.EXE` 的 `start` 只有一條路，`../re/70` §3）；四劇本是否真的共用同一段沒有實跑對過 | 實測 |
| [`spec/30-victory.md`](../spec/30-victory.md) | 君主陣亡時軍師怎麼辦 | 未知（同上） | 靜態 |
| [`spec/31-tactical-sidebar.md`](../spec/31-tactical-sidebar.md) | 段 1 五塊美術的圖形語意 | 貼點與尺寸 confirmed，圖上畫什麼要另外解（`../formats/03` §5.3） | 靜態 |
| [`spec/31-tactical-sidebar.md`](../spec/31-tactical-sidebar.md) | 城兵臨時軍團的主將名 | `0x4200` 的索引算式指到武將表全零那一筆（`../re/60` §4.1） | 靜態 |
| [`spec/33-squad-selection.md`](../spec/33-squad-selection.md) | 待機兵條的欄位語意 | `word_1D30A:+0x09 + 4k` 在 `../re/11` §3.9 記成「第 k 隊的待機兵數」；條的上限 76 遠小於一隊 100 兵，所以開局會頂在上限 | 靜態 |
| [`spec/34-speed-steps.md`](../spec/34-speed-steps.md) | 最高速在原版實機是多少 | 機器相依。DOSBox 固定 cycles 量得到「那台的上限」，量不到「原版的答案」 | 實測 |
| [`spec/34-speed-steps.md`](../spec/34-speed-steps.md) | 戰場幀是否等於 remake 的一次 `Step()` | 原版一幀做完整條戰場迴圈；remake 的 `Step()` 是規則層一步。**兩者對齊過但沒逐項比** | 靜態 |
| [`spec/34-speed-steps.md`](../spec/34-speed-steps.md) | 音效驅動不在時的行為 | `../re/61` §6 | 靜態 |
| [`spec/36-ground-planes-and-climbing.md`](../spec/36-ground-planes-and-climbing.md) | 命令 6 為什麼擋高平面橫移 | 命令碼 6 是什麼沒對過 | 靜態 |
| [`spec/36-ground-planes-and-climbing.md`](../spec/36-ground-planes-and-climbing.md) | 擋路時的換位 | 原版被兵擋住還會試 `loc_1B533`（§1.4），remake 的 `tryClimb` 直接失敗。水平移動那條路的換位已經有了（`swapWith`） | 靜態 |
| [`spec/37-tactical-player-controls.md`](../spec/37-tactical-player-controls.md) | 選了陣形之後原版有沒有立刻重排 | 機器碼只寫偏移，**沒有看到立刻移動的呼叫**；remake 照抄（等命令） | 靜態 |
| [`spec/37-tactical-player-controls.md`](../spec/37-tactical-player-controls.md) | 陣形線在小地圖上的線寬與端點 | `sub_1C5AE` 沒逐行讀，remake 畫整條 1 px 的線 | 靜態 |
| [`spec/38-list-windows.md`](../spec/38-list-windows.md) | 俘虜身分 | remake 的 `Posted` 是 bool，存不下 `+0x17` 的 0–5；俘虜狀態目前推不出來 | 靜態 |
| [`spec/38-list-windows.md`](../spec/38-list-windows.md) | 「看」與「選」的內容差異 | 原版兩種取法的**列表內容**不同（`../re/26` §4.2），remake 只統一了欄位 | 靜態 |
| [`spec/38-list-windows.md`](../spec/38-list-windows.md) | 「委任」那一格的顏色 | 實錄影格上看起來是紅字，但影片是壓縮過的、也沒有機器碼證據。remake 先畫成一般色 | 靜態 |
| [`spec/41-message-box-geometry.md`](../spec/41-message-box-geometry.md) | 君主那一側的回話 | 原版事件場景會同時出現兩個框（`docs/re/66` §5.1 的影格就是），remake 只畫一個 | 靜態 |
| [`spec/41-message-box-geometry.md`](../spec/41-message-box-geometry.md) | 框的底紋 | 龍紋的點陣找到了（`../formats/03` §5.5），但 96 列的排法還沒解，remake 仍用純色 | 靜態 |
| [`spec/42-event-scene-speakers.md`](../spec/42-event-scene-speakers.md) | 撥款事件（4／5） | 同上，還沒對過哪一則進下框 | 靜態 |
| [`spec/43-rout-on-blocked-return.md`](../spec/43-rout-on-blocked-return.md) | `loc_1491B` 的其他成本項 | 只解出「非己方據點 ＋0xA6」，廣度優先搜尋本身沒逐條讀 | 靜態 |
| [`spec/44-advise-original-text.md`](../spec/44-advise-original-text.md) | 逐句節拍 | 原版每句要等玩家按鍵才往下走，remake 直接顯示最新一句（`45` §3.1） | 靜態 |
| [`spec/45-advise-scene-layout.md`](../spec/45-advise-scene-layout.md) | 選單的反白樣式 | 原版怎麼畫游標列沒解，remake 用自己的反白條 | 靜態 |
| [`spec/46-post-battle-retreat.md`](../spec/46-post-battle-retreat.md) | `loc_1491B` 的方向回傳 | `±4` 決定讀哪一個鄰接槽，remake 用 `Route` 的第 2 個節點取代，沒有逐條對過兩者選的是不是同一站 | 靜態 |
| [`spec/46-post-battle-retreat.md`](../spec/46-post-battle-retreat.md) | 野外那一格的鄰接槽 | remake 的 `Node` 在行軍途中停在上一個據點，所以走的是「從上一個據點找路」，不是原版的「從野外那一格的鄰接槽挑」 | 靜態 |
| [`spec/47-city-fall-corps-redirect.md`](../spec/47-city-fall-corps-redirect.md) | `[si+1Ah]` | 據點記錄記下舊主，remake 的 `OwnerRecorded` 是同一格但語意沒逐位元對過 | 靜態 |
| [`spec/48-governor-returns-on-city-fall.md`](../spec/48-governor-returns-on-city-fall.md) | 武將 `+0x1E` 的值域 | 534/535/536 與 542 是空的，所以實際用到的變體大概只有 3–6。哪些武將拿到哪個值沒統計過 | 靜態 |
| [`spec/48-governor-returns-on-city-fall.md`](../spec/48-governor-returns-on-city-fall.md) | `sub_10CE7` 的變數表 | 這裡推出 `{1}` ＝ 武將、`{2}` ＝ 據點（照 push 的順序與譯文），沒有逐個 handler 讀 | 靜態 |
| [`spec/49-advise-relocate-and-sortie.md`](../spec/49-advise-relocate-and-sortie.md) | `sub_16E8F` 編成前的其餘檢查 | 只確認「君主還沒帶軍團」這一條。⚠ 它呼叫的 `sub_16EC9` 本身已解（六槽 × 三候選兵種表、每槽門檻 `0x32`、試算在堆疊副本上做，見 `../re/30` §7.3） | 靜態 |
| [`spec/49-advise-relocate-and-sortie.md`](../spec/49-advise-relocate-and-sortie.md) | 進言的指令列 | 五項在原版指令樹裡的排法（`docs/re/22`）沒有逐格對過，remake 用自己的小視窗 | 靜態 |
| [`spec/52-main-screen-camera-and-banner-date.md`](../spec/52-main-screen-camera-and-banner-date.md) | `+0x0000` | 0x6C0 / — / 未解 | 靜態 |
| [`spec/52-main-screen-camera-and-banner-date.md`](../spec/52-main-screen-camera-and-banner-date.md) | `+0x08F0` | 0x0B0 / `word_10D4C` / 另一組 11 格（未解，可能是別的字重） | 靜態 |
| [`spec/52-main-screen-camera-and-banner-date.md`](../spec/52-main-screen-camera-and-banner-date.md) | `word_10D4C` 那一組 | 與數字字模同樣是 11 格 × 16 列，緊接在後面，用途未解 / 找誰把 `ds` 設成 `cs:word_10D4C`（`KI.EXE.asm` 只有一處） | 靜態 |
| [`spec/54-ui-colours-from-palette.md`](../spec/54-ui-colours-from-palette.md) | 季節換色 | 五個索引在四季調色盤裡的值只有色 14 會變（`../formats/02` §4），而這五個都不是 14，所以目前用第 0 組。**若之後有視窗在別的調色盤組下畫，要改成跟著組走** | 靜態 |
| [`spec/55-minimap-view-box.md`](../spec/55-minimap-view-box.md) | 剩下的 11 byte | `+0x8F0` 那一塊有 176 byte，框只用 165 | 靜態 |
| [`spec/56-battlefield-rotation.md`](../spec/56-battlefield-rotation.md) | 表頭與尾段那各 64 byte | 轉的時候原版**不動它們**（迴圈只掃 `0x40`–`0xFBF`）。內容仍未解 | 靜態 |
| [`spec/56-battlefield-rotation.md`](../spec/56-battlefield-rotation.md) | 鏡頭差一個等角格 | 翻轉之後戰場區還差 (−16, −8)（`../playtest/40` §4.1）。小地圖沒有位移，所以不是翻轉中心的問題 | 靜態 |
| [`spec/57-tactical-projection.md`](../spec/57-tactical-projection.md) | 物件與地形差一列會不會看得出來 | 奇數鏡頭時 anchor 那一半的物件比自己腳下的地形低一格。**原版就是這樣算的**，但沒有找到能單獨驗證這一點的畫面 | 靜態 |
| [`spec/58-display-slot-depth-range.md`](../spec/58-display-slot-depth-range.md) | 旗標 bit 5（`0x20`）／bit 6（`0x40`）誰設 | `sub_1DD22` 只設 bit 7。bit 6 是快路徑那道 `dl & 0x50` 的一半，bit 5 決定要不要跑「unit 0 的第二趟」 | 靜態 |
| [`spec/58-display-slot-depth-range.md`](../spec/58-display-slot-depth-range.md) | unit 0 的第二趟 | 深度迴圈跑完後，`dl & 0x20` 成立時對五個鄰格各跑一次 `ax = 0`。**remake 沒做**，而觸發條件還沒解 | 靜態 |
| [`spec/59-battle-opening-orders.md`](../spec/59-battle-opening-orders.md) | 玩家側的開場常令 | 畫面上看起來是「站在陣形上」，但原版是哪一個命令碼（`Form`／`Holding`／`Guard`）沒有直接證據。⚠ **三個在開場那一幀畫出來一模一樣**——守陣沒有敵人靠近時也是站在陣形位置上，所以截圖分不出來 | 實測 |
| [`spec/59-battle-opening-orders.md`](../spec/59-battle-opening-orders.md) | ⤷ 試過的一條線索 | 原版 `probe-march/e10.png` 那一刻**攻方已折損 115 點兵力**（兵力條反推，`../playtest/51` §2），而 remake 的守方站著不還手。把玩家側的開場常令改成守陣試跑（2026-08-27）：攻方只折損 **30 點**，離 115 還差得遠，而且守方仍然全滅。**所… | 靜態 |
| [`spec/59-battle-opening-orders.md`](../spec/59-battle-opening-orders.md) | 腳本節奏與原版的 tick 對應 | 第 40 步對上那一張截圖，但「原版的 40 個 tick 是多久」還沒對過（`34`） | 實測 |
| [`spec/60-battle-talk-duration.md`](../spec/60-battle-talk-duration.md) | 開戰 pair 的側別對應 | `0x1BA` → 上格、`0x1BB` → 下格是**強推論**（照影格位置接的）；`sub_1A3C3` 怎麼決定側別沒讀（§3.5） | 靜態 |
| [`spec/60-battle-talk-duration.md`](../spec/60-battle-talk-duration.md) | `byte_1D349` 的三個值 | `sub_1A69F` 拿它當「這句要不要顯示」的閘（`al & 6` 那一段還沒逐位讀）。0／1／2 三種值由 `sub_1A6FA` 切換 | 靜態 |
| [`spec/60-battle-talk-duration.md`](../spec/60-battle-talk-duration.md) | 玩家按鍵能不能提早關掉 | remake 可以按鍵推進；原版是否有這條路沒讀 | 靜態 |
| [`spec/61-soldier-initial-hp-from-morale.md`](../spec/61-soldier-initial-hp-from-morale.md) | `+0x18`（戰力）的算式 | `sub_19B6D` 由士氣、`ch` 與 `cs:[bx-63F1h]` 的每兵種係數算出來，還沒逐項拆；remake 目前直接用士氣 | 靜態 |
| [`spec/61-soldier-initial-hp-from-morale.md`](../spec/61-soldier-initial-hp-from-morale.md) | 大將的體力怎麼掉 | 原版那一格在 20 秒內從 200 掉到 140；remake 第 61 步還是滿的。兩邊的時刻本來就不同，**掉法還沒對過** | 靜態 |
| [`spec/63-hit-stun.md`](../spec/63-hit-stun.md) | `+0x13` ← 8 | `sub_1B618` 寫、`sub_1B6BC` 不寫。那個欄位誰讀還沒查 | 靜態 |
| [`spec/63-hit-stun.md`](../spec/63-hit-stun.md) | 倒地動畫（§1.2） | 4 幀之後 `sub_1B4B8` 收掉，remake 直接把 `Alive` 設成 false | 靜態 |
| [`spec/65-retreated-soldiers-survive.md`](../spec/65-retreated-soldiers-survive.md) | `sub_19F2C` 打完那一次數的是什麼 | `../re/11` §3.9 記成「打完時數」；是不是把還站在場上的補進存活數，還沒逐行讀 | 靜態 |
| [`spec/65-retreated-soldiers-survive.md`](../spec/65-retreated-soldiers-survive.md) | 隊長離場時清掉待機數 | remake 的 `squadLeaderGone` 這樣做（`docs/re/11` §5.9），原版是否也把那些待機兵算掉還沒對 | 靜態 |
| [`spec/66-broken-walls-repaint.md`](../spec/66-broken-walls-repaint.md) | 縮小地圖要不要跟著換 | 側欄的縮圖也是從同一個緩衝區來的，但重畫時機還沒讀。這一版只換戰場本身 | 靜態 |
| [`spec/66-broken-walls-repaint.md`](../spec/66-broken-walls-repaint.md) | 對拍那 88 px | 兩邊的時刻不同（§1.2），要對就得讓門在同一個 tick 破 | 靜態 |
| [`spec/67-ending-playback.md`](../spec/67-ending-playback.md) | 淡入淡出的色階算式 | 17 階已確定，每階怎麼算色值沒讀（`sub_1035F`／`sub_103DC`）；remake 先用疊黑 | 靜態 |
| [`spec/67-ending-playback.md`](../spec/67-ending-playback.md) | `END_S12` 右半 | 用 640 版面畫出來右邊是雜訊，可能還有第二塊（`formats/09` §6） | 實測 |
| [`spec/67-ending-playback.md`](../spec/67-ending-playback.md) | 第一幕的捲動 | §3 標成 remake 差異；要做就得先對 `sub_10094` 那一段的逐列位移 | 靜態 |
| [`spec/67-ending-playback.md`](../spec/67-ending-playback.md) | 音樂的**起訖時點** | ⚠ 2026-08-23 起整段結局都放 `endbgm-0`（`cmd/wlgame` 的 `musicTrack()`，排在 `world == nil` 之前——`-open-ending` 那條 fixture 沒有世界）。**先前放的是 `overbgm-0`**，那是另一支執行檔的遊戲結束曲。剩下的缺… | 靜態 |
| [`spec/68-death-animation.md`](../spec/68-death-animation.md) | 大將陣亡 | 大將不會死（`sub_1B618` 的 `IsGeneral` 那一條），所以 `+0` 那一組實際只有騎馬會用到；大將的倒地圖是不是死碼還沒查 | 靜態 |
| [`spec/69-world-fingerprint.md`](../spec/69-world-fingerprint.md) | 跨平台實測 | Android 端還沒有東西可以跑（里程碑 A 本身） | 實測 |
| [`spec/69-world-fingerprint.md`](../spec/69-world-fingerprint.md) | 戰術戰鬥要不要進指紋 | 目前不進。要驗戰場的決定性得另外做一個，`tactical.Battle` 的欄位更多 | 靜態 |
| [`spec/70-phone-chrome.md`](../spec/70-phone-chrome.md) | 外框在高 DPI 上的觀感 | 8×8 的點陣框在 960×540 上是原尺寸；字已放大 2 倍（`100`），框要不要跟著放大**沒量過**，要實機回饋 | 靜態 |
| [`spec/70-phone-chrome.md`](../spec/70-phone-chrome.md) | 龍紋的對齊 | 手機版的面板不是 640×400 的格子，龍紋仍釘在螢幕上，與原版的相位不同。**視覺上看得出來的差異只有相位，不是圖案** | 靜態 |
| [`spec/72-bundled-game-data.md`](../spec/72-bundled-game-data.md) | Windows／macOS 上「解開就能跑」 | ⛔ 沒有那兩個平台的機器。包內版面驗過（`gamedata/`、`fonts/` 位置正確），但 `resolveDataDir` 在那兩個 OS 上沒實跑過 | 實測 |
| [`spec/72-bundled-game-data.md`](../spec/72-bundled-game-data.md) | APK 內嵌後的實機驗收 | ⛔ 沒有裝置。模擬器驗到了解包與指紋，驗不到真實儲存空間與 DPI | 靜態 |
| [`spec/72-bundled-game-data.md`](../spec/72-bundled-game-data.md) | 25 MB 的 APK 在低容量裝置上 | 解包後 app 私有目錄再佔約 4.8 MB，總共約 30 MB。**沒有量過安裝失敗的門檻** | 靜態 |
| [`spec/73-right-click-cancel.md`](../spec/73-right-click-cancel.md) | 原版右鍵是否也關常駐視窗 | 沒量過。常駐視窗不走模態等待常式，推測不關，但**沒有實機證據** | 靜態 |
| [`spec/74-corps-on-world-map.md`](../spec/74-corps-on-world-map.md) | 那 110 張圖在 MCH 裡的實際外觀 | 算式定案，但**沒有逐張看過** 22 勢力 × 5 方向長什麼樣。已看過的：勢力 0 的靜止與行進（`../playtest/50`） | 靜態 |
| [`spec/74-corps-on-world-map.md`](../spec/74-corps-on-world-map.md) | 每格 4 層上限 | 刻意沒做（§4） | 靜態 |
| [`spec/74-corps-on-world-map.md`](../spec/74-corps-on-world-map.md) | 首都疊圖為什麼不畫 | **機制讀出來了**（2026-09-02）：`sub_1D4C7` 是「往一格推一個疊圖」，**一格最多 4 個**——`bl = [si+1]`（目前疊圖數）`cmp bl, 4 / jnb` 就整個跳過，而且 `[si]` 的 **bit 4 設起來的格子完全不接受疊圖**。顯示清單一格 8 B：`+0` … | 靜態 |
| [`spec/74-corps-on-world-map.md`](../spec/74-corps-on-world-map.md) | 別的疊圖組合 | 只驗過「首都 ＋ 軍團」。**災害物件 ＋ 軍團**、**非首都據點 ＋ 軍團**都沒有樣本 | 靜態 |
| [`spec/75-bundled-audio.md`](../spec/75-bundled-audio.md) | 音檔大小 | ogg 全套 19 MB，桌面包從 11.7 MB 漲到 29 MB | 靜態 |
| [`spec/75-bundled-audio.md`](../spec/75-bundled-audio.md) | 沒有音效裝置的**真實玩家** | ⛔ 仍然會掛。驗收模式擋住的是截圖路徑，一般啟動沒有擋——Ebiten 沒有可查詢的音訊 API，目前沒有乾淨的偵測法 | 實測 |
| [`spec/75-bundled-audio.md`](../spec/75-bundled-audio.md) | 音效與場景的對應完整度 | 見 `29`，本規格不重複 | 靜態 |
| [`spec/76-lord-not-in-formation.md`](../spec/76-lord-not-in-formation.md) | 君主被編成之後原版會怎樣 | 沒試過。若原版其實允許、只是清單排序讓人以為不行，§2 要推翻（但開關本身照樣成立） | 靜態 |
| [`spec/76-lord-not-in-formation.md`](../spec/76-lord-not-in-formation.md) | 開關要不要進存檔 | **不進**。與旁邊的速度設定一樣是 session 設定，讀檔不會帶回來 | 靜態 |
| [`spec/77-rout-talk-messages.md`](../spec/77-rout-talk-messages.md) | 戰鬥脫身的 `#1F`／`#20` | 同一支 `sub_12977` 也服務戰鬥脫身；remake 那條路現在畫的是自己的句子，沒接原文 | 靜態 |
| [`spec/78-amount-input-editor.md`](../spec/78-amount-input-editor.md) | `sub_17D5F` 讀 `CS:7D93` 之外還做什麼 | 每格的 raw byte 表已解，但那一支怎麼把 glyph 貼上去沒逐行讀 | 靜態 |
| [`spec/78-amount-input-editor.md`](../spec/78-amount-input-editor.md) | 稅率上限 100 的意義 | 是「100%」還是別的刻度沒有第二個證據；remake 照抄 100 | 靜態 |
| [`spec/79-new-game-faction-list.md`](../spec/79-new-game-faction-list.md) | 欄位表的「屬性」與「型別」兩個 word | `0x0206`／`0x0204` 與 `0x76`／`0x73` 只由「名字欄 vs 數字欄」推語意，消費它們的那一段沒讀（`../re/73` §6） | 靜態 |
| [`spec/79-new-game-faction-list.md`](../spec/79-new-game-faction-list.md) | 捲軸的滑塊樣式 | 同 `38` §4，原版那一支沒讀 | 靜態 |
| [`spec/79-new-game-faction-list.md`](../spec/79-new-game-faction-list.md) | 標題列的底色與字色 | 屬性 `0x9000`／`0x9001` 沒有換算成調色盤索引；remake 沿用一覽表既有的用色 | 靜態 |
| [`spec/79-new-game-faction-list.md`](../spec/79-new-game-faction-list.md) | 無頭點擊 | §3.1：建置 image 沒有視窗管理員，滑鼠按鍵送不進 Ebiten。加一個 WM 就能把所有點擊路徑納入自動驗收 | 靜態 |
| [`spec/79-new-game-faction-list.md`](../spec/79-new-game-faction-list.md) | 橫幅在不在 | remake 的啟動殼層一直有畫橫幅（`ICONGRF` 段 0）。`sub_11A6E` 沒有明顯的橫幅呼叫，**原版那 32 px 是什麼沒驗過**——地圖只佔 y 32–400 | 靜態 |
| [`spec/79-new-game-faction-list.md`](../spec/79-new-game-faction-list.md) | 調色盤組 | `sub_10241(al=0)` 取的是第 0 組；remake 照抄成季 0。**那一組是不是「春」沒有另外驗** | 靜態 |
| [`spec/80-duel-opening.md`](../spec/80-duel-opening.md) | 變體 0／2／3／5／6 的臨場抽驗 | 專屬句只在變體 1／4／7（`../re/75` §1.1），預設句共用同一選句機制；優先度低 | 靜態 |
| [`spec/81-disaster-quantities.md`](../spec/81-disaster-quantities.md) | 事件 11 在風暴期間被重排的節奏 | 產生端一次排 32–60 延遲；期間事件 11 重擲強度的觸發頻率沒逐格讀（推定隨佇列輪到） | 靜態 |
| [`spec/82-defender-selection.md`](../spec/82-defender-selection.md) | 名單上限 127 的邊界 | 原版緩衝區 0x100 bytes ÷ 2；remake 不設上限（軍團總數 127 本來就到不了）——記為等效差異 | 靜態 |
| [`spec/84-multilanguage.md`](../spec/84-multilanguage.md) | UI 詞表只涵蓋靜態 literal 與 `%d` 樣板 | `%s` 樣板（內容是任意文字）反查不回去，這一類仍顯示繁中 | 靜態 |
| [`spec/84-multilanguage.md`](../spec/84-multilanguage.md) | 版面沒有為非全形語系重算 | 訊息框 10 全形字／列是原版幾何；英文靠折行與翻頁吃下去，固定寬的欄位只做了「不可超過原文寬度」的約束 | 靜態 |
| [`spec/84-multilanguage.md`](../spec/84-multilanguage.md) | UI 詞表與名表的第二人覆核 | **做了**（2026-08-28，374 條 UI 詞、343 個名）：改 4 條 UI 詞（`敵正侵攻他國` 縮成 `Enemy invades others` 對齊 TALK #102、`原版格式`→`Original format`、`提出外交要求` 補 diplomatic、`成立`→`agreed` … | 靜態 |
| [`spec/84-multilanguage.md`](../spec/84-multilanguage.md) | 兩個專有名詞的英文寫法未經查證 | `松崗` 目前寫 `Sung Kang`、遊戲名寫 `Wolong`。原版包裝上有沒有官方英文名沒查過——**拿到一手資料再定案**，不要照現在這個沿用 | 靜態 |
| [`spec/84-multilanguage.md`](../spec/84-multilanguage.md) | **16 px 字型在清單裡列距為 0** | 版面是照倚天的 16×15 排的（列高 16 ⇒ 1 px 間隙），JIS／GB 的字是 16×16 ⇒ 貼在一起。**裁掉首列或末列會傷到 89% 的字**（量過：兩份字型的第 0 列與第 15 列各有 82–91% 帶墨），所以不裁。PC-98 原版的文字模式本來就是 16 px 列距貼齊，日文這樣看是對的；… | 靜態 |
| [`spec/85-latin-list-layout.md`](../spec/85-latin-list-layout.md) | 軍團／勢力兩家仍有 8–9 字的欄 | 要再寬就得動視窗本身的寬度（384 px 是原版幾何），那會影響對拍基準——**要做就整個視窗一起重排，不要只加寬一欄** | 靜態 |
| [`spec/85-latin-list-layout.md`](../spec/85-latin-list-layout.md) | 其他視窗（HUD、據點情報卡、編成）還沒為半形語系重排 | 那些是固定座標的欄位，逐一量過再改；目前靠「英文不可超過原文寬度」的約束擋住溢出 | 靜態 |
| [`spec/86-runtime-language-switch.md`](../spec/86-runtime-language-switch.md) | Android 實機／模擬器還沒實地切過語言 | 面板與 `syncLanguage` 都有單測，但手機版的畫面沒拍過；下一次 Android 打包驗收時補 | 靜態 |
| [`spec/86-runtime-language-switch.md`](../spec/86-runtime-language-switch.md) | 語言不進存檔 | 原版存檔格式沒有這一欄，塞進去會破壞 round-trip。要記住偏好得另存 remake 自己的設定檔 | 靜態 |
| [`spec/86-runtime-language-switch.md`](../spec/86-runtime-language-switch.md) | F9 是 remake 自創的鍵 | 原版沒有這個鍵；`docs/spec/13` 的按鍵表要同步記一筆 | 靜態 |
| [`spec/87-latin-screen-layout.md`](../spec/87-latin-screen-layout.md) | 半形語系的戰場標題沒有地名 | 要顯示就得動側欄幾何（128 px 是原版的），那會影響對拍基準——要做就整格重排 | 靜態 |
| [`spec/87-latin-screen-layout.md`](../spec/87-latin-screen-layout.md) | 原版美術上的中文（§2） | 只能重畫美術，本專案不做 | 靜態 |
| [`spec/87-latin-screen-layout.md`](../spec/87-latin-screen-layout.md) | 將旗名字 8 字仍有 10 個撞號 | 地板是 8（同音不同字）；再多兩個是 `XIAHOU-D`／`XIAHOU-Y` 這類，要解得換呈現方式不是換寬度 | 靜態 |
| [`spec/88-display-polish-parity.md`](../spec/88-display-polish-parity.md) | 事件列本身是 remake 自創 | 原版怎麼提示月結（如果有）沒查過；目前只是讓它不擋畫面 | 靜態 |
| [`spec/88-display-polish-parity.md`](../spec/88-display-polish-parity.md) | `playtest/40` 沒有涵蓋開闊地的兵 | 那一份量的兩個局面都在城壁邊。要擋住這一類回歸，得再加一個**開闊地**的對拍 fixture | 靜態 |
| [`spec/89-siege-damage-report-toggle.md`](../spec/89-siege-damage-report-toggle.md) | 設定不進存檔 | 與語言同一條（`86` §7）：原版存檔沒有這一欄，要記住偏好得另存 remake 自己的設定檔 | 靜態 |
| [`spec/89-siege-damage-report-toggle.md`](../spec/89-siege-damage-report-toggle.md) | 結果畫面本身原版沒有 | 原版打完直接回戰略畫面。要不要整頁拿掉是另一個裁定，本規格只讓多出來的那一行可關 | 靜態 |
| [`spec/90-same-state-parity.md`](../spec/90-same-state-parity.md) | 各視窗**內部**的排版 | 分區的外框已由機器碼定死（§3），框內的頭像／文字列座標仍是影片估值（`docs/spec/12` §7） | 靜態 |
| [`spec/90-same-state-parity.md`](../spec/90-same-state-parity.md) | 送點擊的座標 | DOSBox-X 的**視窗**是 640×480，遊戲的 640×400 在 y 偏移 40（`tools/parity_crop.py` 量的），而 INT 33 把整個視窗等比對映到遊戲畫面——**送 y 要乘 1.2，不是減 40**。這是本機設定的性質，把 `int33 max y` 改成 400 應該… | 實測 |
| [`spec/90-same-state-parity.md`](../spec/90-same-state-parity.md) | 主畫面的四窗狀態 | 開局四個視窗全關（`sub_11A6E` 結尾 `mov cs:byte_198A6, 0`）。要開得先移游標再按同一點（`docs/re/47` §3.1），單純 `click` 會被當成移動吃掉 | 靜態 |
| [`spec/90-same-state-parity.md`](../spec/90-same-state-parity.md) | 調色盤季節組 | 兩側都要鎖同一組，否則整片顏色不同（`docs/formats/02`） | 靜態 |
| [`spec/91-tactical-parity.md`](../spec/91-tactical-parity.md) | 動畫幀序 | 原版的兵有 `PoseStep`，截圖時機差一幀就整批不同。這是 `field` 剩下那 299 px 的來源之一（`../playtest/40` §13） | 實測 |
| [`spec/91-tactical-parity.md`](../spec/91-tactical-parity.md) | 野戰打起來之後的戰場 | **開戰第一拍已經對過**（九區裡七區 0 px、`field` 0.05%，`../playtest/43`）。沒對過的是**開戰之後**的畫面——野戰的地形是從大地圖即時長出來的，兩邊要走到同一個 tick 比攻城更難湊 | 靜態 |
| [`spec/92-android-music.md`](../spec/92-android-music.md) | 手機沒有實機聽過 | 沒有裝置；模擬器是 `-no-audio` 起的 | 靜態 |
| [`spec/92-android-music.md`](../spec/92-android-music.md) | 桌面「沒有音效裝置就掛」那一條 | `../release/08` §5 的老問題，手機端要確認 Ebiten 在 Android 上不會踩同一個 | 靜態 |
| [`spec/94-retreat-path-not-cleared-every-frame.md`](../spec/94-retreat-path-not-cleared-every-frame.md) | 守方一兵未損 | 這條 fixture 從頭到尾守方 48 人全活、兵 600 不變（`../playtest/51` §3）。守方站在 Z=1，攻方在 Z=0，而 `doAttack` 的碰撞要求 `s.Z == e.Z`——攻方走到腳下卻上不去。原版靠命令 3（城壁）登城，這條 fixture 的腳本沒有下過那道命令 | 靜態 |
| [`spec/94-retreat-path-not-cleared-every-frame.md`](../spec/94-retreat-path-not-cleared-every-frame.md) | 攻方大將體力歸零卻還在 | `drainSiegeGeneral` 會扣到 0，而 `docs/re/11` §5.16 說戰鬥傷害最低留 1。攻城計時器要不要也留 1，沒有讀出來 | 靜態 |
| [`spec/95-spawn-height-uses-ground-plane.md`](../spec/95-spawn-height-uses-ground-plane.md) | 原版擺兵寫進 `[si+0Ah]` 的是哪一個表 | 沒讀出來。remake 這一版是靠內部一致性定的（§2），不是照抄 | 靜態 |
| [`spec/98-climb-into-a-gate-hits-it.md`](../spec/98-climb-into-a-gate-hits-it.md) | `loc_1B533` 的完整分流 | 這裡只用到「撞到結構物」那一條。它同時也是敵我碰撞的入口（`../re/11` §5.16），另外兩條在水平移動那邊已經接了 | 靜態 |
| [`spec/99-about-page-license.md`](../spec/99-about-page-license.md) | 條款全文沒有在手機上顯示 | 只顯示摘要與全文出處。全文 104 行，要另做可捲動的文字頁；摘要已滿足「收到的人知道自己被授權了什麼」 | 靜態 |

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

## 9. 平台層與工具鏈產物（不計入總數）

兩類東西掛 `[DOS/BIOS]`，都**不算 remake 的缺口**：

1. **原版與 DOS／BIOS 之間的介面**：`INT` 服務號、顯示卡暫存器、磁碟服務。
   知道 `INT 61h` 的 `ah=4` 是什麼，不會改變任何一行 Go。
2. **編譯器 runtime、連結進來的驅動與程式庫**：C runtime 啟動碼、算術／字串
   輔助常式、被靜態連結進執行檔的驅動模組（本作的 segment `0x2000` 整段是
   `INT 33h` 滑鼠包裝，見 [`24`](24-unread-function-catalogue.md) §3）。
   **這些是 toolchain 產物，不是玩家看得到的遊戲邏輯。**

使用者裁定：第 1 類 2026-08-23、第 2 類 2026-09-02。

⚠ **排除的是「怎麼跟平台講話」，不是「原版選了什麼參數」。**
滑鼠驅動把游標範圍設成 640×400 是**遊戲行為**，仍然要讀；
`INT 33h` 本身的呼叫慣例不用。判準與流程見
`~/.claude/knowledge-base/retro/compiler-runtime-helper-fingerprints.md`。

⚠ **分開數不是不數。** 這些仍然是原版的未解之處，只是**不擋 remake**；
哪天要寫「原版怎麼跟 DOS 講話」的文件，這一節就是清單。

| 出處 | 缺口 | 現況 |
|---|---|---|
| [`re/17-dosv-audio-tsr.md`](../re/17-dosv-audio-tsr.md) | `INT 61h` 的四個服務號 `[DOS/BIOS]` | `ah=4`／`7`／`8` 與 `ax=09F2h`／`0C01h`，對應什麼動作要看 `YNSOUND.COM`（[`42`](42-leaf-functions.md) §7）。⚠ **這是原版與音效 TSR 的介面，不擋 remake**——音訊走純 Go 的 OPL3 渲染（[`../spec/29`](../spec/29-audio.md)），不經過 DOS |
| [`re/28-text-number-rendering.md`](../re/28-text-number-rendering.md) | `sub_1F7A4` `[DOS/BIOS]` | 把 32 B 字模緩衝畫上 VRAM 的實際迴圈，未逐行讀。⚠ remake 要的是**畫什麼**（字模版面，已解），不是**怎麼寫 VRAM**——Ebiten 不碰 VGA 平面（同 [`29`](29-font-service-int15.md) §9） |
| [`re/29-font-service-int15.md`](../re/29-font-service-int15.md) | `sub_1F7A4` `[DOS/BIOS]` | 把 32 B 緩衝畫上 VRAM 的實際迴圈，未逐行讀。⚠ remake 要的是**畫什麼**（字模版面，已解），不是**怎麼寫 VRAM** |
| [`re/29-font-service-int15.md`](../re/29-font-service-int15.md) | `YNFONT.EXE` 怎麼顯示中文 `[DOS/BIOS]` | 它不走 INT 15h（0 次），密碼輸入畫面的中文是它自己畫的。⚠ 那是一支 DOS TSR，remake 沒有對應物；密碼頁本身也不擋任何事（`CLAUDE.md` §4.0） |
| [`re/37-graphics-and-runtime-module-map.md`](../re/37-graphics-and-runtime-module-map.md) | `sub_1F7A4` `[DOS/BIOS]` | 212 / 字型 blitter，逐行未解。**同一支函式在 [`29`](29-font-service-int15.md) §9 也列著，那裡是正本**；未解的是「怎麼寫 VRAM」，而 remake 不碰 VGA 平面，所以不擋 remake |
| [`re/42-leaf-functions.md`](../re/42-leaf-functions.md) | `INT 61h` 的四個服務號（`ah=4`／`7`／`8`、`ax=09F2h`／`0C01h`）`[DOS/BIOS]` | 對應什麼音效動作要看 `YNSOUND.COM`（[`17`](17-dosv-audio-tsr.md)）。⚠ 原版與音效 TSR 的介面，**不擋 remake** |
