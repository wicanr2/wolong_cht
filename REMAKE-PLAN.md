# 臥龍傳 remake 計畫

## Legal and integrity boundary

- Original files are user-supplied and not distributed.
- Pristine originals are read-only；所有存檔與修改資料使用 writable overlay。
- Clean rewrite language/runtime：Go + Ebiten；規則層以 tick 驅動，排版層不依賴 Ebiten。
- 玩家可見文字走語系資料；繁中為松崗母本，日文只作對照 oracle；英文／簡體本輪不做。

## Inventory

| Category | Files/builds | Known format | Evidence | Status |
|---|---|---|---|---|
| executable | `KI.EXE`、`YN*` 工具鏈、`D7*` | DOS/V／PC-98 執行檔 | `docs/re/01`、`CLAUDE.md` §3.1／§3.10 | 已盤點；`KI.EXE` 兩版已進 IDA |
| maps/events | `MMAP.*`、`BATTLE.*`、`PASS.*`、`MOUSE.*` | RLE／裸資料／松崗容器 | `docs/formats/04`–`07` | 主要格式 READY；MCH 語意及部分容器壓縮未解 |
| saves | `SINARIO.DAT`、`SAVE.DAT` | 4 × 22,208 B 區塊；改寫策略 | `docs/formats/08` | 結構與 round-trip 已有；部分欄位未知 |
| graphics/fonts | `*GRF.DAT`、`*.BRG`、`FONTGRF.DAT` | 4bpp planar、BRG 調色盤 | `docs/formats/02`、`03` | 主要圖庫與調色盤已解；ICONGRF 段 1 未解 |
| audio | `BGM.DAT`、`SOUND.DAT`、`YNSOUND.COM` | PC-98 FM + SSG confirmed；DOS/V 未解 | `CLAUDE.md` §3.9、說明書紀錄 | remake／DOS/V parity 待驗 |
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
| TALK.DAT text | `dosv` + `pc98` byte data | `tools/talkdat.py` + translations + `internal/ui/textdraw.WrapLines` | byte-for-byte round-trip、日中對照、實測字寬／硬斷行測試 | READY；60 筆校訂、ASCII／CJK 實測寬度、原始 hard line、結構尾空行與五行／16 px TALK 分頁已接；未定位 formatter 分支與自然整張畫面逐像素 parity 未完；DOS/V cursor／ICONGRF button glyph 已解碼 |
| world/time | PC-98／`SINARIO.DAT` | `internal/assets/world` + `internal/rules/clock` + `cmd/wlgame` | fixed-cycle DOSBox、歷史截圖 | 已實作；2026-08-10 Docker／Xvfb gate 已重跑 |
| strategy settlement | `KI.EXE` monthly/hourly paths | `internal/rules/economy`、diplomacy、persuasion、`internal/rules/strategyai` | deterministic tests + RE notes + PC-98 event3 fixture | 月結評估／事件 1 宣戰／事件 2 合作產生器、狀態、玩家三選一與 `sub_17C6E` 數值語意／事件 3 停戰產生器、狀態、玩家三選一與 `sub_17C6E` 數值語意／事件 4／5 官員撥款產生器、狀態、玩家三選一與 `sub_17C6E` 數值語意／事件 6／7 外交官回報狀態、主要 TALK #57／#58／#43–#45／#47–#49 與玩家進言 producer／事件 8 遷都與 `sub_14502` 軍團同步／事件 9 指定武將釋放狀態與可見 `#37` 通知（`#409` 是空槽 no-op）／事件 10 訊息邊界／事件 11／12 runtime 災害 marker、`sub_14269` 持久效果與 `TALK.DAT #70/#71/#72` 通知／事件 13 `#51` 通知／敵方編成狀態切片已接；事件 2／3 raw fixture→前置 TALK→三選一→3×6 滑鼠格位與消像、事件 4／5 共用數值選取與 TALK 五行分頁、DOS/V cursor／button glyph 資產已驗；自然整張畫面 parity、事件 6／7 次要反應、事件 10 訊息、事件 11／12 物件 timer parity、多軍團請求與完整行軍狀態機未完 |
| tactical battle | `KI.EXE` tactical paths + BATTLE data | `internal/rules/tactical`／battle renderer／遭遇決策 | parser invariants、公式測試、畫面／正常路徑 | 遭遇選單→戰鬥指揮→攻城戰場→攻擊命令→戰後結果報告已由正常 UI 路徑驗證；一般近戰／大將／投射物核心、退卻／繞路點修正與 `sub_1AD7F` CH=0x20 特殊效果已接；`sub_1B941` 的投射物先命中／再移動／後威力更新、raw 格索引與戰場標記已接；正常真實攻城狀態層結算、結果資料流與 GUI 回戰略已通過；`+0x1E` 的寫入端、鎖敵分支與 `sub_1AC55` raw 平面條件已接成 `PlaneHigh`；`sub_1AD2D` 初始速度公式與 `sub_1ECE0`／`sub_1EC82` RNG 已接並有單測；特殊投射物保留發射兵 `+0x02 bit 0` 並可取 raw `0x214/0x215`；原版完整投射物畫面／timer 與同狀態對拍未完 |
| save/load | `SINARIO.DAT`／`SAVE.DAT` 四槽 | `state.SaveInto` + `cmd/wlgame` system modal + `internal/savepath` | round-trip、overlay 差異、pristine hash、Trust `+0x10`／Player `+0x0D,+0x0F` round-trip、事件佇列 raw／節拍／月壓縮／1／2／3／4／5／6／7／8／9／13 handler 測試、`TestQueuedTalkNotices`／`TestQueuedDiplomacyReportTalkNotices`／`TestRawAmountEditorSemantics`、event3 fixture、Xvfb `4→S→Return` | 可玩 overlay、Trust、Player 雙欄、事件佇列原始 256 筆、每十次節拍、月度壓縮與事件 1／2／3／4／5／6／7／8／9（狀態、主要 TALK 句型取用）／13 handler、事件 11／12／13 的 `TalkNotice` 與 modal GUI 已接、玩家外交／撥款三選一與 raw 3×6 數值選取、event3 raw fixture→composite→消像已接；事件 6／7 次要反應／原版數值排版、事件 10、事件 11／12 物件動畫、原版外框／游標逐像素與完整原版 save parity 仍未完 |
| release | remake builds only | `tools/release.sh` + Docker 等價封裝流程 + deny-list | PE/ELF/Mach-O 檔頭、unpacked smoke、asset scan | Linux amd64、Windows amd64、macOS Intel／Apple Silicon 候選包已產出；Linux 封裝 Xvfb smoke 與 deny-list 通過；Windows／macOS GUI 目標 runtime 實跑仍未完 |

### 2026-08-09 事件指定金額 outcome 勘誤

`sub_13902`（外交）與 `sub_139E8`（撥款）的指定金額邊界已分開接入：外交超過初始要求
不收尾但依 `sub_13C3D` 扣信賴度 30；撥款 0 無效但超過初始要求仍完成。
`ResolveDiplomacy` 保存初始 terms，避免確認時重抽平手政治 RNG；
`TestDiplomacyAndFundingAmountOutcomeBounds` 固定此差異。PC-98 數值視窗與原版訊息排版
仍列為未完成。

### 2026-08-09 事件 4／5 TALK 結果／收尾接入

IDA `CS:08A4` handler 表已證實 `\4` 是玩家軍師、`\6` 是不可見欄位控制、`\7` 是十進位金額；`cmd/wlgame/funding.go` 依 `sub_139E8` 的 base+offset 接回事件 4／5 的結果與收尾訊息。`TestFundingTalkIndicesMatchRaw139E8Branches` 通過 12 個分支；原版數值視窗像素排版、逐頁動畫與長程正常劇本仍不能標成完成。

### 2026-08-10 進言／說得信賴度切片

`sub_13830`／`sub_13C1E`／`sub_13B5A`／`sub_13BA9` 的結果碼、`+20`／`+10`／`−20`、
四段信賴度理由門檻已接到 `internal/rules/persuasion` 與 `cmd/wlgame`：
`Situation.Trust` 在說服開始時取 `World.Trust`，`ReactionTrustDelta` 處理直接反應，
`Session.Offer` 處理多理由與錯選。固定邊界測試與 Docker 內受影響套件驗收通過。
這只完成一條規則／GUI 垂直切片，不改變事件 2／3 其他外交增減、事件 6／7 次要 TALK、
事件 10、物件動畫、目標平台 runtime 與同狀態對拍的未完成狀態；正式包與推廣影片 gate
仍保持關閉。

## Worklist

| ID | Deliverable | Evidence needed | Acceptance gate | Status |
|---|---|---|---|---|
| M7-A | 60 筆 `translations/corrections.json` | `TALK.DAT` 兩版索引、IDA 槽位索引證據、逐句內容對照 | `talkdat.py correct`／`verify` + selftest + `WrapLines`／Xvfb 抽樣 | 60 筆已可重跑套用並接入 `wlgame`；#0–#1021 第一輪已讀；remake 實測行寬／hard line／五行分頁／尾空行已測，未定位 formatter 與逐句畫面抽樣未完 |
| M7-B | 1,022 則文意層校訂 | 日文／繁中逐句對照 | 變數、名詞、漏譯、刪節與決策可回查 | 第一輪逐句讀取完成；60 筆 runtime 產出已鎖定，硬換行／行寬、formatter 與畫面 parity 進行中 |
| M8-A | 目標平台建置 | `tools/release.sh` 等價 Docker 矩陣、PE/ELF/Mach-O 檢查、packaged Linux `-shot` | 交叉編譯目標正確、產物非同一平台假成功、發行目錄可啟動 | Linux／Windows／Darwin 純 Go 產物、Linux 原生本體與封裝 smoke 通過；Windows／macOS GUI runtime 未實機驗證 |
| M8-B | 正常玩家路徑與畫面 | PC-98 固定狀態 oracle、event3 raw fixture、有效時鐘的原版／AI 存檔 | 無 debug hook、同狀態截圖／狀態對拍 | 編成／行軍／城兵攻城／敵方 AI 遭遇選單／戰術畫面、攻擊命令、結果報告與 GUI 回戰略接縫已完成；事件 3 raw fixture→前置 TALK→三選一→3×6 實際點擊→消像短路徑已完成；DOS/V 96×64 內框、3×6 button glyph、`KI.EXE` 16×16 hardware cursor 已解碼接線；仍缺自然整張畫面对拍、其他事件物件與跨平台實機 |
| M8-C | 發行隔離 | deny-list、可寫 save overlay | 原始資產零命中、解包 smoke | deny-list／overlay smoke 已通過；完整目標平台矩陣未完成 |
| M9-A | Android 手機版規劃 | `docs/mobile/android-plan.md`、固定 Android Docker 工具鏈、`arm64-v8a` debug APK | 橫向安全區、觸控 hitbox、TALK／數值二段確認、pause/resume 不重複 tick | 規格已建立；尚未實作觸控層、Android 工具鏈或 APK |
| RE-1 | ICONGRF 段 1／龍紋／MCH 等 | IDA／原版畫面或檔案不變量 | `docs/re/` + `docs/mechanics/` 雙寫 | `MMAP.MCH` type 1／2 已完成格式與接線；ICONGRF 段 1／龍紋仍待排程 |

## Intentional differences

| Area | Original | Remake | Reason | Player impact |
|---|---|---|---|---|
| realtime pacing | 速度上限跟機器效能 | 固定且可重現的 remake 時間基準 | 現代硬體避免不可玩 | 文件需明標；規則不變 |
| presentation | 640×400 原版畫面／平台繪製 | 高解析、縮放、現代輸入 | 可跨平台與可讀性 | 不得改核心規則 |
| text templates | 原版片段與 ASCII 插入碼 | 具名參數語系模板 | 可維護與可校訂 | 原版格式仍保存在研究文件 |
| font path | 平台字型驅動／原版環境 | `internal/assets/cjk`，玩家指定字庫 | 不散布倚天字型 | 需在發行文件說明 |

## 2026-08-09 事件 2／3 前置外交通知切片

`sub_138C7`／`sub_138E6` 的 TALK 基底已接到 pending 外交選單：停戰為 #360，協力為
#373，兩者以 `{3}` 展開請求方君主。`TestQueuedDiplomacyChoiceTalkNotices` 與
`TestDiplomacyTalkExpansionUsesOriginalRequestMarkers` 已通過；事件 2／3 後續訊息池、
PC-98 數值視窗與完整玩家長程路徑仍列為未完成，不提高 M8／M7 狀態。

### 事件 2／3 主要 TALK

`sub_13902` 的建言 base+4/+5/+6 與 `sub_13C3D` 的主要結果 #43–#45／#47–#49 已接入
`enqueueDiplomacyTalk`；`TestDiplomacyTalkIndicesMatchRaw13902Branches` 與真實 TALK
展開測試通過。AH／信賴度次要回覆與 PC-98 數值視窗仍未完成；長時間完整遊戲測試不列為
本輪阻塞，但發行前仍需短 smoke 與封裝檢查。

### 事件 9 通知條件與空槽

`sub_150D7` 的 #37 只對釋放後歸屬玩家勢力的武將顯示；`enqueueEventMessages` 與
`TestReleasedGeneralTalkOnlyTargetsPlayerFaction` 已接入。後續 `CX=0x199`／#409 在
PC-98／DOS/V `TALK.DAT` 只有資料空行，`TestReleasedGeneralRawFollowup409IsEmptyNoOp`
固定不產生空白 modal；原版空呼叫時序與長程 oracle不冒充已驗。

### 災害 marker 視覺接線

`World.DisasterMarkerAt`／`StormAreaSnapshot` 與 `wlgame.drawDisasterOverlay` 已把事件
11／12 的 runtime 狀態接到地圖；事件 12 火災／暴動現在使用由 `MMAP.MCH` 解出的原版
type 1／2 圖塊矩陣，缺檔才回退向量 marker。固定 phase clock、暴風雨範圍輪廓與同狀態
畫面對拍仍是明示替代／P1 未完成項。`TestDisasterMarkerReadOnlySnapshots`、MCH 資產測試、
GUI 短測試與 Docker build 通過。

### 事件 10 raw 邊界

已確認 `sub_131AE` dispatch table 的 `0x0A` handler 是 `sub_13496`；IDA `.i64` 直接
xref 也已覆核 `sub_12FBF`／`sub_12FB1`／`sub_1301C` 的 caller 集合與事件碼，沒有找到
0x0A producer。其他 `0Ah` 常數沒有 queue 資料流證據。維持事件 10 未知／負證據，不以
猜測補 state 或 TALK；若繼續只追函式指標／間接 writer。

### 事件 6／7 次要 formatter 邊界

`sub_137D8`／`sub_13138` 的 `AH` 是雙向俘虜配對旗標；`sub_13C3D` 的次要呼叫則在恢復
`DI` 後不重建第一次 TALK 所用的 stack formatter 參數。直接 caller 的 `CX+0x1D` 可證實
#72／#76；#73／#77 尚未定位，原始槽位的 marker／選單內容不能可靠映射到目前
`TalkNotice`。列為
**strong inference／未完成**，不得用猜測文案解除事件 6／7 或 release gate。

### 2026-08-10 普通箭初始速度 RNG 勘誤

`sub_1ACA4`／`sub_1AD2D` 的兩軸最大距離、目標／射手高度差、`sub_1ECE0 & 3` 與
`0x14` 乘法已接入 `normalProjectileVelocity`；`internal/rules/rng` 已重現
`sub_1ECE0`／`sub_1EC82`，新增公式測試並在 Docker 通過。原版投射物圖形／完整動畫、
同狀態時序與畫面對拍、目標平台 GUI runtime 仍未完成，release gate 不變。

### 2026-08-10 第 3／4／5 項定向接手結果

1. 原版 PC-98 event3 fixture 已真正進入前置通知與三選一；`sub_13B7E`／`sub_193E9`
   的選單輸入保存鏈、`sub_17C6E` 的數值編輯核心與 3×6／16×16／96×64 幾何證據已
   回填 `docs/re/13-pc98-numeric-window.md`。
2. `textdraw.WrapLines` 接到 TALK modal：先展開 marker，再以 ASCII 8 px／CJK 16 px、
   22 full-width cell 進行單一 hard line 內換行，空列／標點／分頁均有 Docker 測試。
3. 特殊投射物新增 `PoseStep` → `SpecialFrame`，可取原版 raw `0x214`／`0x215`；事件
   12 MCH type 1／2 仍使用八相位矩陣，固定 phase fallback 也有回歸測試。
4. 同一 raw save fixture 的原版／remake 截圖已保存於 `docs/images/`；本輪已補上原版式
   肖像／IVENTGRF composite、3×6 真實格位點擊、TALK 五行分頁與 pending 結束後消像。
   自然 DOS/V／remake 整張畫面與 PC-98 視覺 oracle 的逐像素 parity 仍未接，因此不解除三平台包與
   影片 release gate。完整劇情測試依使用者要求不列入本輪阻塞。

### 2026-08-10 游標／數值選取／TALK 封口補記

- `CS:7D93h` 的 18 bytes 已固定為 3×6 raw 格位表；`amountPanelButtonAtPoint` 以
  `(88,200)` 起點與 16×16 格直接命中，滑鼠點擊、游標高亮、鍵盤 fallback 都呼叫同一
  `AmountEdit`。`0x60` 的完成格會先保留目前值，再離開數值器。
- `messages.go` 已對齊原版五行／16 px TALK 頁面；marker 展開後按 ASCII 8 px／非 ASCII
  16 px 測量換行，保留 hard line／中間空行，只移除結構性的最後空行。事件 2／3／4／5
  的場景、肖像、prompt、選項與數值器由同一 Draw 層組合。
- `PendingDiplomacy`／`PendingFunding` 清除後，下一個 Draw 先重畫地圖，才顯示後續 TALK，
  因此 IVENTGRF／肖像不殘留；這是功能性消像 parity。DOS/V 硬體游標與數值外框資產
  已由原始 bytes 解碼接線；自然整張畫面的逐像素 parity 仍是明示替代，不能以短 smoke
  截圖宣稱完成。
- 畫面證據：`docs/images/wlgame-event3-choice.png`（SHA-256
  `CA40B865B44A6EA13ED5B4F2C0B6AB913A0BC895EF48D7A19E1825501E535151`）與
  `docs/images/wlgame-event3-amount.png`（SHA-256
  `27A5474EBA79C92C23B24A79938CA4E1D376B9FA52C0956AE3D3359C0404609D`）。
- 完整長程遊戲測試依使用者指示略過；事件 6／7 未定位次要 formatter、事件 10、MCH
  timer／投射物逐像素動畫、Windows／macOS GUI runtime 仍是開啟 gate，所以本輪不建立
  三平台正式包或推廣影片。
