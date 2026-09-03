# 00 — RE 知識庫入口

**狀態：索引。44 份反組譯筆記按子系統分類，每份一句話說它回答什麼問題。**

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
| **還有什麼沒解** | [`43`](43-open-questions.md) 的缺口總表（生成的）|
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
| [`64`](64-corps-arrival-state-machine.md) | 軍團抵達時的狀態機：`+0x23` 的分派表與解體 |
| [`65`](65-ai-march-decision-chain.md) | **電腦勢力的行軍決策鏈**：Stage 0–3 與整條軍團生命週期 |
| [`59`](59-game-over-exit-codes.md) | 結局與敗北是靠 `KI.EXE` 的離開碼交出去的 |

## 3. 戰略層：玩家怎麼下指令

| 筆記 | 回答什麼 |
|---|---|
| [`22`](22-strategy-command-tree.md) | **指令列與熱區分派**——整個戰略層的入口 |
| [`25`](25-message-variants-and-personnel.md) | 訊息索引 ≥ `0x196` 的 ×8 變體；人事四個指令 |
| [`77`](77-general-affinity-and-flags.md) | ⭐ 武將記錄 `+0x19` ＝ 心向的勢力（在野出仕與俘虜歸降共用）；旗標 bit 5／bit 6 |
| [`30`](30-corps-formation-ui.md) | 軍團編成：兵員池、兵種循環、移動間隔公式 |
| [`31`](31-faction-picker-screen.md) | 勢力一覽（同時是領地圖）；主畫面的重繪分派表 |
| [`32`](32-strategy-detail-panels.md) | 據點詳情欄與軍團詳情欄的欄位對照 |
| [`26`](26-list-window-engine.md) | 一覽表視窗引擎與 23 個 callback |
| [`27`](27-list-row-fields.md) | 一覽表怎麼畫每一列；外交關係等級 |
| [`13`](13-pc98-numeric-window.md) | 數字輸入視窗的量測與 CJK 版面決策 |
| [`46`](46-strategy-chrome-cell-layer.md) | 主畫面的指令列沒有按鈕圖，外框取自 `ICONGRF` |
| [`47`](47-main-screen-window-registry.md) | 主畫面四個常駐視窗的開關、位元集與分派 |
| [`48`](48-window-display-list.md) | **視窗內容是一份顯示清單，不是一張圖**：記錄格式與十個場景 |
| [`49`](49-corps-formation-window.md) | 軍團編成視窗的版面與動態層 |
| [`50`](50-city-info-window.md) | 據點情報視窗；據點 `+0x16` 高 4 位的用途 |
| [`51`](51-corps-info-window.md) | 軍團情報視窗（顯示清單場景 4） |
| [`52`](52-slot-select-window.md) | 四槽選擇視窗：新遊戲、讀取、儲存共用 |
| [`53`](53-lord-select-window.md) | 君主選擇視窗（顯示清單場景 8） |
| [`54`](54-advisor-naming-window.md) | 軍師命名視窗（松崗版特有） |
| [`55`](55-system-menu-window.md) | 系統選單視窗；`op 08` 屬性的編碼 |
| [`71`](71-strategy-hotspot-dispatch.md) | 戰略層的兩張熱區分派表；點縮小地圖會發生什麼 |
| [`72`](72-world-map-display-list.md) | 大地圖的顯示表：地形一層 ＋ 最多四層疊圖 |
| [`73`](73-new-game-faction-list.md) | 新遊戲怎麼選君主：先一張清單，再一張卡 |

**指令列共八項**：進言／人事／財政／編成／軍團／據點／武將／勢力
（字串在 `cs:6181h`，見 [`31`](31-faction-picker-screen.md) §1）。

## 4. 戰術層：戰場

| 筆記 | 回答什麼 |
|---|---|
| [`05`](05-battle-selection.md) | 政略與戰術的接縫：戰場怎麼被選出來 |
| [`11`](11-tactical-battle.md) | 戰術模組結構與戰場資料模型 |
| [`80`](80-pathfind-request-queue.md) | ⭐ 尋路是**反應式**的：全域 128 格佇列、每幀只算兩個兵，走不動才排隊 |
| [`09`](09-combat.md) | 觸發、自動判定、傷亡、武將的下場 |
| [`19`](19-outcome.md) | 敗北 outcome 的接線 |
| [`20`](20-ida-re-coverage-audit.md) | 戰術管線被移植進 remake 的程度（**量的是移植度不是覆蓋度**）|
| [`74`](74-battle-opening-duel.md) | 開戰喊話是單挑狀態機的開頭：挑戰、拒戰、對嗆、決著 |
| [`75`](75-duel-talk-audit.md) | 單挑台詞的逐組逐變體抽驗：24 組 × 8 變體 |
| [`78`](78-soldier-power-from-command.md) | ⭐ 兵記錄 `+0x18`（戰力）來自**統率力**不是士氣；近戰的命中率與傷害都看它 |
| [`82`](82-display-slot-dead-flags.md) | ⭐ 顯示格旗標 **bit 5／bit 6 兩版都是死碼**：唯一的設定端與清除端沒有呼叫端，所以「unit 0 的第二趟」永遠不會執行 |
| [`83`](83-post-battle-troop-accounting.md) | ⭐ **打完的兵力是三項相加**（退場生還 ＋ 場上存活 ＋ 待機），士氣按兵力比例縮；順帶封閉「待機數沒有任何地方被清成 0」|

## 5. 呈現層：圖、字、聲

| 筆記 | 回答什麼 |
|---|---|
| [`02`](02-palette-routine.md) | `.BRG` 的通道順序與亮度縮放 |
| [`03`](03-image-blitter.md) | 圖庫載入器與 VGA 繪製 |
| [`28`](28-text-number-rendering.md) | **文字與數字怎麼畫**：EGA Set/Reset、屬性的 4-bit 配對 |
| [`29`](29-font-service-int15.md) | **原版怎麼顯示中文**：`INT 15h AH=50h` ＋ 倚天字型 |
| [`33`](33-shared-draw-helpers.md) | 共用繪圖層：字串包裝、肖像四格快取、小地圖上色 |
| [`18`](18-tactical-button-glyphs.md) | 戰術底列按鈕的 glyph 資產 |
| [`60`](60-tactical-sidebar.md) | **戰術側欄那一欄畫了什麼**：七格的美術／文字／計量條、小地圖上四種標記的顏色、30 個熱區碼的分派表，順帶解掉攻城的門強度條 |
| [`17`](17-dosv-audio-tsr.md) | DOS/V 音源 TSR 與戰術效果碼 |
| [`63`](63-ground-plane-map.md) | ⭐ **登城機制**：地面層表在 `word_1D2FC`（層 0–3 低平面、層 4–6 高平面），方向位元 ＋ 牆頂正規化，**上下城牆只在門那一格** |
| [`62`](62-strategy-minimap.md) | **主畫面縮小地圖**：192 個據點標記的四種顏色與 4×4 圖形、視野框、圖例兩格、點第二格開的勢力選擇視窗 |
| [`61`](61-timer-tick-source.md) | **計時中斷是誰發的**：音效驅動把 PIT 設成 4660.9 Hz、分頻 16 回呼遊戲，兩層速度設定的實際幀率與日長 |
| [`23`](23-bgm-resource-format.md) | `*BGM.DAT` 的音樂資源格式 |
| [`56`](56-bgm-track-events.md) | `*BGM.DAT` 的聲軌事件編碼 |
| [`57`](57-opl3-register-map.md) | DOS/V 的音源是 OPL3：六個聲軌各佔一組 4-operator |
| [`58`](58-bgm-scene-mapping.md) | 哪一首配哪個場景：`BGM.DAT` 的 11 首全部對出 |
| [`81`](81-sound-type-attenuation.md) | ⭐ 系統選單的 TYPE 1–4 是**四段衰減**不是四種音源：`AH=0Bh` 的參數加進載波的 Total Level |
| [`66`](66-message-box-geometry.md) | 訊息框的版面：一個框、一張肖像、四列字 |
| [`79`](79-talk-marker-handlers.md) | ⭐ `TALK.DAT` 七個 `\N` 標記的 handler：哪一張表、哪一個欄位、什麼顏色。**`\1` 取的是呼び名不是姓名** |
| [`67`](67-city-emblem-on-strategy-map.md) | 大地圖上的據點徽記：位置就在記錄座標 |
| [`70`](70-d7end-ending-player.md) | `D7END.EXE`：結局播放器與結局全文 |
| [`76`](76-d7open-opening-player.md) | ⭐ `D7OPEN.EXE`：開場六幕、開場旁白全文，與**資料檔的 4 byte 長度頭**（載入器 `LSEEK` 跳過它才解壓）|

## 6. 外交與訊息

| 筆記 | 回答什麼 |
|---|---|
| [`12`](12-diplomacy-dialogue.md) | 停戰說服的訊息索引 `#190`–`#198` |
| [`15`](15-event10-producer.md) | 事件 10 producer 的深度逆向 |
| [`40`](40-garrison-relief-request.md) | 據點求援與援軍派遣 |
| [`44`](44-threat-and-reinforcement-ai.md) | **威脅偵測與 AI 出兵**：據點每 tick 掃一次、軍團數上限、派將只看武力 |
| [`45`](45-corps-command-mode.md) | **軍團的三種指令模式**：戰鬥指揮／委任／解體；求援只調得動委任中的軍團 |

## 7. 方法論與量測（讀之前先看這幾條）

| 筆記 | 回答什麼 |
|---|---|
| [`21`](21-function-census.md) | 全函式覆蓋普查；**§3.1 為什麼要排除目錄型文件** |
| [`24`](24-unread-function-catalogue.md) | 未讀函式的證據與下手順序 |
| [`39`](39-remaining-unread.md) | 未讀函式的逐支歸屬（生成的，可重跑）|
| [`43`](43-open-questions.md) | **缺口總表**：各文件「未解」表的集中版（生成的，`check.sh` 每次重出）|
| [`42`](42-leaf-functions.md) | 戰術以外的 47 支葉節點；`INT 61h`、`byte_198A6` 位元圖、第三處自我修改碼 |
| [`20`](20-ida-re-coverage-audit.md) | 手挑取樣的問題，與全量量測的差別 |
| [`68`](68-t3-frontier-functions.md) | T3 那九支逐支讀完；順帶解出顯示格表頭 `+1`／`+3` 的分工 |
| [`69`](69-t2-cross-reference.md) | T2 那 18 支逐支讀完；**新讀的與只是指路的分開標** |
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

## 8. 現況

**四個分級全部收斂到 T1，已解**——739 支函式每一支都有 `docs/re/` 層級的記錄。
最後兩批是 T3 九支（441 bytes，[`68`](68-t3-frontier-functions.md)）
與 T2 十八支（714 bytes，[`69`](69-t2-cross-reference.md)）。

> **收斂到 T1 不等於「全部讀懂了」。** T1 只保證**有人寫過**；
> 模組全圖的角色標籤是**強證據**不是 confirmed，各文件的「未解」表
> 才是真正的缺口清單。
>
> ⚠ T2 那一批**多半只是登記**：18 支裡有 5 支的機制早就逐行解過，
> 只是解釋寫在 `docs/spec/` 或 `docs/mechanics/`，分級卻因為
> 「`docs/re/` 沒提到符號名」算成未讀（[`69`](69-t2-cross-reference.md) §1）。

模組級全圖見
[`35`](35-strategy-ui-module-map.md)（戰略 UI）、
[`36`](36-tactical-module-map.md)（戰術戰鬥）、
[`37`](37-graphics-and-runtime-module-map.md)（圖庫與 C runtime）、
[`38`](38-strategy-core-module-map.md)（指令樹、到站處理、月結）。

**數字是 2026-08-18 的快照，要現況一律重跑**
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
