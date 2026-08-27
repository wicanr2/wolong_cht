# 臥龍傳 remake 計畫

> ⚠ **待辦看 [`CONTEXT.md`](CONTEXT.md) §7，不是這裡。**
> 本檔是**架構與垂直切片的定位**：法務邊界、素材盤點、資料流、
> 刻意的 remake 差異，以及各切片的驗收方式。
> 下半部按日期的節是**當時的封口紀錄**，不是現況。

## Legal and integrity boundary

- Original files are user-supplied and not distributed.
- Pristine originals are read-only；所有存檔與修改資料使用 writable overlay。
- Clean rewrite language/runtime：Go + Ebiten；規則層以 tick 驅動，排版層不依賴 Ebiten。
- 玩家可見文字走語系資料；繁中為松崗母本。**四個語系都已接通**（繁中／簡體／日文／英文，
  `docs/spec/84`）：日文直接讀 PC-98 原版而不是翻回去，簡體是 OpenCC 機轉，英文是逐則英譯。

## Inventory

| Category | Files/builds | Known format | Evidence | Status |
|---|---|---|---|---|
| executable | `KI.EXE`、`YN*` 工具鏈、`D7*` | DOS/V／PC-98 執行檔 | `docs/re/01`、`CLAUDE.md` §3.1／§3.10 | 已盤點；`KI.EXE` 兩版已進 IDA |
| maps/events | `MMAP.*`、`BATTLE.*`、`PASS.*`、`MOUSE.*` | RLE／裸資料／松崗容器 | `docs/formats/04`–`07` | 主要格式 READY；MCH 語意及部分容器壓縮未解 |
| saves | `SINARIO.DAT`、`SAVE.DAT` | 4 × 22,208 B 區塊；改寫策略 | `docs/formats/08` | 結構與 round-trip 已有；部分欄位未知 |
| graphics/fonts | `*GRF.DAT`、`*.BRG`、`FONTGRF.DAT` | 4bpp planar、BRG 調色盤 | `docs/formats/02`、`03` | 主要圖庫與調色盤已解；ICONGRF 段 1 未解 |
| audio | `BGM.DAT`、`SOUND.DAT`、`YNSOUND.COM` | PC-98 是 YM2203；**DOS/V 是 OPL3（YMF262）**，靠 `0x104`／`0x105` 兩個 OPL2 沒有的暫存器定案 | `docs/re/23`、`docs/re/57`、`docs/re/58` | 已解並實作：純 Go 的 OPL3 渲染成 ogg，場景對應出自機器碼（`docs/spec/29`）|
| manuals/walkthrough | `Garyouden_Manual.pdf`、追補頁 | 純掃描影像 | `docs/reference/01` | 有實質機制頁已判讀；少數附錄未讀 |

## Architecture

```text
platform/UI → internal/rules + internal/state → internal/assets / parsed data
            → reusable clock / renderer / storage / RNG
```

Android 先沿用同一條資料流，只在平台殼增加安全區、視口轉換、觸控手勢與生命週期；手機規劃與驗收邊界見 [`docs/mobile/android-plan.md`](docs/mobile/android-plan.md)。

## Vertical slices

| Slice | Original oracle | Parser/rule/UI/save path | Verification | Status |
|---|---|---|---|---|
| TALK.DAT text | `dosv` + `pc98` byte data | `tools/talkdat.py` + translations + `internal/ui/textdraw.WrapLines` | byte-for-byte round-trip、日中對照、實測字寬／硬斷行測試 | READY；60 筆校訂、ASCII／CJK 實測寬度、原始 hard line、結構尾空行與五行／16 px TALK 分頁已接；未定位 formatter 分支未完；DOS/V cursor／ICONGRF button glyph 已解碼 |
| world/time | PC-98／`SINARIO.DAT` | `internal/assets/world` + `internal/rules/clock` + `cmd/wlgame` | fixed-cycle DOSBox、歷史截圖 | 已實作；2026-08-10 Docker／Xvfb gate 已重跑 |
| strategy settlement | `KI.EXE` monthly/hourly paths | `internal/rules/economy`、diplomacy、persuasion、`internal/rules/strategyai` | deterministic tests + RE notes + PC-98 event3 fixture | 月結評估／事件 1 宣戰／事件 2 合作產生器、狀態、玩家三選一與 `sub_17C6E` 數值語意／事件 3 停戰產生器、狀態、玩家三選一與 `sub_17C6E` 數值語意／事件 4／5 官員撥款產生器、狀態、玩家三選一與 `sub_17C6E` 數值語意／事件 6／7 外交官回報狀態、主要 TALK #57／#58／#43–#45／#47–#49 與玩家進言 producer／事件 8 遷都與 `sub_14502` 軍團同步／事件 9 指定武將釋放狀態與可見 `#37` 通知（`#409` 是空槽 no-op）／事件 10 訊息邊界／事件 11／12 runtime 災害 marker、`sub_14269` 持久效果與 `TALK.DAT #70/#71/#72` 通知／事件 13 `#51` 通知／敵方編成狀態切片已接；事件 2／3 raw fixture→前置 TALK→三選一→3×6 滑鼠格位與消像、事件 4／5 共用數值選取與 TALK 五行分頁、DOS/V cursor／button glyph 資產已驗；事件 6／7 次要反應、事件 10 訊息、事件 11／12 物件 timer parity、多軍團請求與完整行軍狀態機未完 |
| tactical battle | `KI.EXE` tactical paths + BATTLE data | `internal/rules/tactical`／battle renderer／遭遇決策 | parser invariants、公式測試、畫面／正常路徑 | 遭遇選單→戰鬥指揮→攻城戰場→攻擊命令→戰後結果報告已由正常 UI 路徑驗證；一般近戰／大將／投射物核心、退卻／繞路點修正與 `sub_1AD7F` CH=0x20 特殊效果已接；`sub_1B941` 的投射物先命中／再移動／後威力更新、raw 格索引與戰場標記已接；正常真實攻城狀態層結算、結果資料流與 GUI 回戰略已通過；`+0x1E` 的寫入端、鎖敵分支與 `sub_1AC55` raw 平面條件已接成 `PlaneHigh`；`sub_1AD2D` 初始速度公式與 `sub_1ECE0`／`sub_1EC82` RNG 已接並有單測；特殊投射物保留發射兵 `+0x02 bit 0` 並可取 raw `0x214/0x215`；⭐ **同狀態逐區對拍：九區裡六區逐像素相同**（`docs/playtest/40`，2026-08-18；⚠ 該組重跑未重現，`docs/playtest/49`）；原版完整投射物畫面／timer 未完 |
| save/load | `SINARIO.DAT`／`SAVE.DAT` 四槽 | `state.SaveInto` + `cmd/wlgame` system modal + `internal/savepath` | round-trip、overlay 差異、pristine hash、Trust `+0x10`／Player `+0x0D,+0x0F` round-trip、事件佇列 raw／節拍／月壓縮／1／2／3／4／5／6／7／8／9／13 handler 測試、`TestQueuedTalkNotices`／`TestQueuedDiplomacyReportTalkNotices`／`TestRawAmountEditorSemantics`、event3 fixture、Xvfb `4→S→Return` | 可玩 overlay、Trust、Player 雙欄、事件佇列原始 256 筆、每十次節拍、月度壓縮與事件 1／2／3／4／5／6／7／8／9（狀態、主要 TALK 句型取用）／13 handler、事件 11／12／13 的 `TalkNotice` 與 modal GUI 已接、玩家外交／撥款三選一與 raw 3×6 數值選取、event3 raw fixture→composite→消像已接；事件 6／7 次要反應／原版數值排版、事件 10、事件 11／12 物件動畫與完整原版 save parity 仍未完 |
| release | remake builds only | `tools/release.sh` + Docker 等價封裝流程 + deny-list | PE/ELF/Mach-O 檔頭、unpacked smoke、asset scan | Linux amd64、Windows amd64、macOS Intel／Apple Silicon 候選包已產出；Linux 封裝 Xvfb smoke 與 deny-list 通過；Windows／macOS GUI 目標 runtime 實跑仍未完 |

## Worklist

| ID | Deliverable | Evidence needed | Acceptance gate | Status |
|---|---|---|---|---|
| M7-A | 60 筆 `translations/corrections.json` | `TALK.DAT` 兩版索引、IDA 槽位索引證據、逐句內容對照 | `talkdat.py correct`／`verify` + selftest + `WrapLines`／Xvfb 抽樣 | 60 筆已可重跑套用並接入 `wlgame`；**#0–#1021 兩批逐句讀完**（`docs/reference/02` §11／§12）；remake 實測行寬／hard line／五行分頁／尾空行已測；**校訂後的畫面抽樣已做**（18 則，`docs/playtest/41`）。未定位 formatter 仍未完 |
| M7-B | 1,022 則文意層校訂 | 日文／繁中逐句對照 | 變數、名詞、漏譯、刪節與決策可回查 | **兩批逐句讀取完成**；60 筆 runtime 產出已鎖定；**排版 parity 全量量過**（`docs/playtest/32`，單行超寬 0 行）、**畫面抽樣已做**（`docs/playtest/41`）。缺的是**兩版並排的畫面對照** |
| M8-A | 目標平台建置 | `tools/release.sh` 等價 Docker 矩陣、PE/ELF/Mach-O 檢查、packaged Linux `-shot` | 交叉編譯目標正確、產物非同一平台假成功、發行目錄可啟動 | Linux／Windows／Darwin 純 Go 產物、Linux 原生本體與封裝 smoke 通過；Windows／macOS GUI runtime 未實機驗證 |
| M8-B | 正常玩家路徑與畫面 | PC-98 固定狀態 oracle、event3 raw fixture、有效時鐘的原版／AI 存檔 | 無 debug hook、同狀態截圖／狀態對拍 | 編成／行軍／城兵攻城／敵方 AI 遭遇選單／戰術畫面、攻擊命令、結果報告與 GUI 回戰略接縫已完成；事件 3 raw fixture→前置 TALK→三選一→3×6 實際點擊→消像短路徑已完成；DOS/V 96×64 內框、3×6 button glyph、`KI.EXE` 16×16 hardware cursor 已解碼接線；⭐ **同狀態逐區對拍已完成**（主畫面五區逐像素相同、戰場九區裡六區；`docs/playtest/37`／`40`）；仍缺其他事件物件與跨平台實機 |
| M8-C | 發行隔離 | deny-list、可寫 save overlay | 原始資產零命中、解包 smoke | deny-list／overlay smoke 已通過；完整目標平台矩陣未完成 |
| M9-A | Android 手機版 | `docs/mobile/android-plan.md`、`docs/mobile/android-ux.md`、固定 Android Docker 工具鏈、`arm64-v8a` debug APK | 橫向安全區、觸控 hitbox、TALK／數值二段確認、pause/resume 不重複 tick | **核心已接入**：手機端共用 `internal/rules`／`internal/state`，模擬器與桌面在 frame 1／60／120 的指紋逐字相同；四個入口、戰場、存讀檔、四語系切換與音樂都在。剩下的兩件不是程式：**實機驗收**（沒有裝置）與 **release signing**（keystore 保管未決）|
| RE-1 | ICONGRF 段 1／龍紋／MCH 等 | IDA／原版畫面或檔案不變量 | `docs/re/` + `docs/mechanics/` 雙寫 | `MMAP.MCH` type 1／2 已完成；⭐ **全函式靜態分析收斂到 T1**（739/739 有 `docs/re/` 筆記）；ICONGRF 段 1 的 UI 語意／龍紋仍待排程 |

## Intentional differences

| Area | Original | Remake | Reason | Player impact |
|---|---|---|---|---|
| realtime pacing | 速度上限跟機器效能 | 固定且可重現的 remake 時間基準 | 現代硬體避免不可玩 | 文件需明標；規則不變 |
| presentation | 640×400 原版畫面／平台繪製 | 高解析、縮放、現代輸入 | 可跨平台與可讀性 | 不得改核心規則 |
| text templates | 原版片段與 ASCII 插入碼 | 具名參數語系模板 | 可維護與可校訂 | 原版格式仍保存在研究文件 |
| font path | 平台字型驅動／原版環境 | `internal/assets/cjk`，玩家指定字庫 | 不散布倚天字型 | 需在發行文件說明 |

## 按日期的紀錄在哪裡

本檔只寫**定位**：法務邊界、素材盤點、資料流、垂直切片與刻意的 remake 差異。
按日期的封口紀錄在 [`WORKLIST.md`](WORKLIST.md)，逆向的證據台帳在
[`RESEARCH-LOG.md`](RESEARCH-LOG.md)。

> 本檔曾經在下半部掛著 2026-08-09／08-10 的封口小節，內容與 `RESEARCH-LOG.md`
> 同期條目重疊，而句子是現在式的——「release gate 仍被…阻擋」「本輪不建立
> 三平台正式包或推廣影片」在四個平台出貨、兩支推廣片上線之後就變成假斷言。
> ⭐ **現在式的紀錄會在事情做完的那一刻變成錯的**，所以它只能有一份，
> 而且要放在按日期讀的地方。
