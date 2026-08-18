# Verification matrix

| Subsystem | Unit/real-data test | Original oracle | Visual oracle | Player path | Result |
|---|---|---|---|---|---|
| startup/new game | 歷史 DOSBox／PC-98 playtest | `docs/playtest/01`、`02`、`06` | `docs/images/wlgame-release-smoke.png`（本輪） | PC-98 已可進本體 | Docker/Xvfb 啟動 smoke 通過；正常策略／戰術路徑仍在收斂 |
| movement/time | `internal/rules/clock`、`wlsim -years 1 -check` | `docs/re/06`、`15-realtime` | `wlgame` 歷史截圖 | 時間世界已跑通 | 一年／77760 tick smoke 通過 |
| map transitions | `internal/assets/world`／道路測試／`TestNormalScenarioMarchIntoGarrison` | `docs/re/04`、`05`、`08` | `wlview -world`、`docs/images/wlgame-normal-*.png`、`wlgame-ai-normal-encounter.png`、`wlgame-ai-battle-afterpatch.png`、`wlgame-ai-postbattle.png` | 編成／目的地／行軍／汝南城兵攻城／濮陽敵方遭遇→戰術畫面→戰後回戰略已正常操作 | 編成／行軍／城兵攻城／敵方 AI 遭遇、戰術接縫與 GUI 回戰略通過；狀態層正常攻城結算已由 `TestNormalScenarioTacticalBattleTerminates` 通過，戰後訊息與同狀態原版對拍未完 |
| economy/strategy | economy、diplomacy、persuasion tests | `docs/re/07`、`08` | 策略視窗歷史截圖 | 進言／說得已接畫面 | 大致完成 |
| battle/AI/rewards | tactical／combat/state 遭遇選擇與公式 tests、`TestNormalScenarioMarchIntoGarrison`、`TestNormalScenarioTacticalBattleTerminates`、`TestSpecialProjectileUsesCH20AndFallsVertically`、`TestSpecialProjectileCarriesPoseBitIntoSecondRawFrame`、`TestClimbingInfantryCanUseSpecialProjectile`、`TestRawPlaneHighAndTerrainFlag`、`TestSpecialProjectileUsesPlaneHighAndMaxAxisDistance`、`TestNormalProjectileVelocityMatchesRawSub1AD2D`、`TestLockOnPlanePenaltyMatchesRawBranches`、`TestProjectileRawDirectionGridAndHeightPower`、`TestProjectileChecksCurrentCellBeforeMoving`、`TestProjectileStopsAtSolidLayerAfterMoving`、`TestStrategicAIScenarioOneProducesEnemyWarPath` | `docs/re/09`、`11`、`07` §21 | `docs/images/wlgame-battle-choice.png`、`wlgame-save-replay-choice.png`、`wlgame-normal-garrison.png`、`wlgame-ai-normal-encounter.png`、`wlgame-ai-battle-afterpatch.png`、`wlgame-ai-battle-attack-afterpatch.png`、`wlgame-ai-battle-result.png`、`wlgame-ai-postbattle.png` | 城兵攻城、敵方 AI 宣戰／編成／行軍／遭遇選單→戰術畫面→攻擊命令→結果報告→回戰略已驗；正常真實攻城狀態層結算為守方勝、攻方 0／守方 100；CH=0x20 特殊效果、raw `PlaneHigh` 條件、普通箭初始速度公式、投射物逐幀狀態與特殊 raw `0x214/0x215` 已由獨立測試驗證；畫面有戰場內側別／特殊標記 | 城兵自動判定、遭遇決策兩分支、AI 狀態來源、核心戰術接縫、狀態層結算、結果資料流、GUI 回戰略、特殊效果、raw 平面條件與投射物逐幀順序已通過；完整原版投射物畫面與同狀態對拍未完 |
| events/quests | `TestEventQueueStorage`／`Timing`／`MonthlyCompaction`／`TestQueuedEventHandlers`／`TestQueuedDiplomacyReportHandlers`／`TestQueuedDiplomacyReportTalkNotices`／`TestPlayerDiplomacyProducers`／`TestQueuedDisasterAnimationHandlers`／`TestDisasterMarkerAppliesRawPersistentEffects`／`TestDisasterAnimationPhaseIsEightFrameAndReproducible`／`TestQueuedTalkNotices`／`TestQueuedEventReleaseGeneral`／`TestQueuedDiplomacyChoice`／`TestRawAmountEditorSemantics`／`TestQueuedFundingInitialTalkNotices`／`TestFundingRequestGenerators`／`TestQueuedFundingChoice`／`TestStrategicAIDiplomacyEventGenerators`／`TestEvent2To5FullTalkPageSampling`／`TestEvent9LongNaturalRoute`／`TestEvent9LongNotificationRoute` | `docs/re/07`、`08`、`docs/formats/01`、`docs/re/13`、`docs/playtest/12-event3-same-state-parity.md` | [`wlgame-event-modal.png`](docs/images/wlgame-event-modal.png)、[`pc98-oracle-event3-choice.png`](docs/images/pc98-oracle-event3-choice.png)、[`wlgame-event3-choice.png`](docs/images/wlgame-event3-choice.png)、[`wlgame-event3-amount.png`](docs/images/wlgame-event3-amount.png)、事件 2–5／事件 9 抽樣代表幀 | 短 fixture／長程 bounded fixture 通過；完整長程劇本依要求略過 | queue raw、節拍／月壓縮、事件 1–13 已接狀態接縫、事件 2／3／4／5 的 36 頁／18 組 TALK 回應、事件 9 27 小時 queue 與玩家通知條件、事件 10 raw message 邊界、事件 11／12 runtime 災害 marker 與物件 timer 已通過；原版完整自然路徑與逐像素 parity 仍是明確界線 |
| save/load | 四劇本 round-trip、`TestTrustStorage`、`TestPlayerStorage`、`TestEventQueueStorage`／`Timing`／`MonthlyCompaction`／`QueuedEventHandlers`／`TestQueuedDiplomacyReportHandlers`／`TestQueuedDiplomacyReportTalkNotices`／`TestPlayerDiplomacyProducers`／`TestQueuedDisasterAnimationHandlers`／`TestDisasterMarkerAppliesRawPersistentEffects`／`TestQueuedTalkNotices`／`TestQueuedEventReleaseGeneral`／`TestQueuedDiplomacyChoice`／`TestRawAmountEditorSemantics`／`TestQueuedFundingChoice`、`internal/savepath` 路徑測試 | `docs/formats/08` §1.2、`CONTEXT.md` §3、`docs/formats/01` | Xvfb `4→S→Return`、重新啟動載入 | writable overlay + pristine hash | overlay、Trust `+0x10`、Player `+0x0D,+0x0F`、事件佇列原始 u16、每十次節拍／月度壓縮與 1／2／3／4／5／6／7／8／9／10／11／12／13 handler、玩家外交／撥款／進言 producer 接縫與事件 6／7 主要 TALK、事件 11／12／13 `TalkNotice` 接通；事件 9 的 `ReleasedGenerals`／`TALK.DAT` 句型取用與事件 11／12 `sub_14269` 持久效果已接入；事件 2／3／4／5 的 raw 3×6 數值選取、TALK 五行分頁、event3 composite／消像短 smoke、DOS/V cursor／button glyph 資產解碼已接通；事件 6／7 次要反應／原版數值排版、事件 9 原版完整流程、事件 10 訊息、事件 11／12 物件動畫呈現、自然整張畫面 parity 與完整 save parity 未完成 |
| graphics/themes | palette／GRF／world parser tests | `docs/formats/02`、`03`、`05` | 原版素材截圖 | 四季／視窗／頭像已驗 | 段 1／龍紋未完 |
| localization/fonts | `tools/talkdat.py` round-trip；`talkdat_selftest.py`（60 筆套用與 runtime 產出一致）；`internal/ui/textdraw`／`cmd/wlgame` measured-wrap tests；`internal/assets/gfx` cursor／glyph tests；`tools/m7_review.py --check` | `docs/reference/02`、`docs/formats/01`、`docs/re/12`、`docs/re/13` | [`wlgame-event3-choice.png`](docs/images/wlgame-event3-choice.png)、[`wlgame-event3-amount.png`](docs/images/wlgame-event3-amount.png)、M7 六張代表幀與事件 2–5 四張代表幀 | M7 60 筆人工語意／硬換行／寬度／版面抽樣通過 | 60 筆定案校訂已逐筆審查並由 runtime gate 保護；事件 2–5 的 36 頁／18 組 TALK 已抽樣；原版 formatter 未定位分支與自然整張畫面逐像素 parity 仍是界線 |
| packaging | Docker 等價 `tools/release.sh` 跨平台矩陣、`dist/release-20260811/packages`、deny-list、封裝 Linux `wlgame -orig ... -font ... -shot` | ELF／PE／Mach-O 檔頭與解包檔案掃描 | packaged Linux Xvfb PNG（30 幀） | Windows/macOS GUI 目標平台實跑未完 | Linux amd64、Windows amd64、macOS Intel／Apple Silicon 候選包、每包 SHA-256、Linux smoke 與 deny-list 通過；Windows／macOS 原生 runtime 未驗 |

### 2026-08-09 outcome 勘誤

`TestDiplomacyAndFundingAmountOutcomeBounds` 已加入事件驗收：外交指定金額超過初始要求時
無外交收尾但依 `sub_13C3D` 扣 Trust 30；撥款指定 0 無副作用、超過初始要求仍照輸入金額
完成。這是 state／RNG 時序驗收，未把 PC-98 數值視窗或原版訊息排版標成完成。

### 2026-08-10 事件 10 producer 負證據

`sub_131AE`／`funcs_131E8` 的 0x0A handler 已定位到 `sub_13496`；策略 queue helper 的
已檢查 caller 未找到 0x0A producer。此為 raw RE 未完成項，不能把 `sub_13496` 的訊息邊界
誤寫成事件 10 完整串接，也不解除 packaging／影片 release gate。

### 2026-08-10 事件 6／7 次要 formatter 邊界

`sub_137D8`／`sub_13138` 的 `AH` 俘虜旗標與 `sub_13C3D` 的第二次 `CX+0x1D` 呼叫已由
IDA 線性位址 `00013138`／`000137D8`／`00013C3D` 交叉確認；第二次呼叫沒有重建第一次
`DI=SP` formatter 參數。直接 caller 的算術只證實 #72／#76；#73／#77 尚未定位，
仍是未知。這些 raw index 不能安全展開成目前通知欄位，因此事件 6／7 次要 TALK 仍列
未完成，不能解除正式包／影片 gate。

### 2026-08-10 災害 marker 視覺接線

新增 `TestDisasterMarkerReadOnlySnapshots`：驗證事件 11／12 runtime marker 與暴風雨範圍
的唯讀 snapshot；`wlgame` 已以向量替代層接到地圖，但原版物件動畫仍列未完成。事件／
外交 GUI 短測試與 `go build` 在 Docker／Xvfb 通過；不因此解除「三平台正式包與推廣影片
須等所有串接完成」的 release gate。

## Screenshot metadata

| Image | Build | Save/scenario | Map/coordinate | Seed/time | Theme | Match type |
|---|---|---|---|---|---|---|
| `docs/images/wlgame-*.png` | `cmd/wlgame`（部分於 2026-08-09 重拍） | real `SINARIO.DAT` | world map／戰術接縫 | fixed-cycle／clock state | original-derived | 視覺證據；未列出的歷史畫面仍待重拍 |
| `docs/images/wlgame-save-ui.png` | `cmd/wlgame`，2026-08-09 | real `SINARIO.DAT`、`-save-file` overlay | 系統視窗／第 1 槽 | `speed=0`、`4→S→Return` | original-derived chrome + remake save UI | Docker/Xvfb 操作證據 |
| `docs/images/wlgame-battle-choice.png` | `cmd/wlgame`，2026-08-09 | real DOS/V `SINARIO.DAT`、`-open-battle-choice` | 真實地圖／遭遇視窗 | scenario 0、player 0、`speed=0` | original-derived map/chrome + remake choice UI | Docker/Xvfb；驗收旗標，不冒充正常路徑 |
| `docs/images/wlgame-event-modal.png` | `cmd/wlgame`，2026-08-09 | real DOS/V `SINARIO.DAT`、`-open-message`（玩家首都、TALK #70） | 戰略地圖／事件通知 modal | scenario 0、player 0；196/4/1；第 30 幀 | original-derived map/chrome + remake TALK 行邊界 | Docker/Xvfb；驗收旗標；30 幀期間日期仍為 196/4/1，不冒充原版完整翻頁／肖像 |
| `docs/images/pc98-oracle-event3-choice.png` | PC-98 原版 `KI.EXE`，2026-08-10 | event3 PC-98 `SAVE.DAT` fixture；`0x0303`、來源 3、玩家 0 | 事件 3 前置通知後的原版三選一 | 196/4/1；PNG SHA-256 `56A1A16BC6D92F75A5DB3DC49F3C961609F1EB6B908106A39136D4E3E32FDB5C` | 原版 PC-98 | Docker/DOSBox-X；同狀態 oracle；含原版肖像／物件疊層 |
| `docs/images/wlgame-event3-choice.png` | `cmd/wlgame`，2026-08-10 | 同一 raw event3 fixture overlay；`-seed 1`、按一次 `Enter` | remake 外交通知後的三選一／數值面板 | 196/4/1；第 600 幀；PNG SHA-256 `C02183091189AF12B643FC0DB157DB02908C3A156C7434352F6D6A6B0172D1DA` | original-derived map/chrome + remake CHT modal | Docker/Xvfb；同狀態接線證據，不冒充像素 parity |
| `docs/images/wlgame-normal-formed.png` | `cmd/wlgame`，2026-08-09 | real DOS/V `SINARIO.DAT`、鍵盤 `A` 編成 | 真實地圖／編成完成事件 | scenario 0、player 0、196/4/1 | original-derived map/chrome + remake UI | Docker/Xvfb；正常資源限制 |
| `docs/images/wlgame-normal-destination.png` | `cmd/wlgame`，2026-08-09 | real DOS/V `SINARIO.DAT`、鍵盤 `M` | 目的地一覽／虎牢關 | scenario 0、player 0、196/4/1 | original-derived map/chrome + remake UI | Docker/Xvfb；正常兩段式選取 |
| `docs/images/wlgame-normal-march.png` | `cmd/wlgame`，2026-08-09 | real DOS/V `SINARIO.DAT`、確認虎牢關後行軍 | 真實地圖／行軍事件 | scenario 0、player 0、196/4/2 | original-derived map/chrome + remake UI | Docker/Xvfb；正常路徑切片 |
| `docs/images/wlgame-normal-garrison.png` | `cmd/wlgame`，2026-08-09 | real DOS/V `SINARIO.DAT`、無 `-open-*`，選汝南 | 真實地圖／城兵攻城事件 | scenario 0、player 0、曹操一槽步兵；汝南 owner=袁術、城兵=127 | original-derived map/chrome + remake UI | Docker/Xvfb；正常道路 79 格、城兵自動判定 |
| `docs/images/wlgame-ai-normal-encounter.png` | `cmd/wlgame`，2026-08-09 重拍 | real DOS/V `SINARIO.DAT`、`-seed 17`、正常鍵盤 `A`／`M`／方向鍵；無 `-open-*` | 濮陽／敵方軍團遭遇選單 | scenario 0、player 0；196/6/28；呂布對曹操、攻城 | original-derived map/chrome + remake choice UI | Docker/Xvfb；事件佇列接入後的正常 AI 月結／編成／道路行軍接點 |
| `docs/images/wlgame-ai-afterpatch.png` | `cmd/wlgame`，2026-08-09 重拍 | real DOS/V `SINARIO.DAT`、`-seed 17`、正常鍵盤；無 `-open-*` | 遭遇選單 | scenario 0、player 0；196/6/28；呂布對曹操、戰鬥指揮前 | original-derived map/chrome + remake encounter UI | Docker/Xvfb；正常玩家戰術入口前一幀 |
| `docs/images/wlgame-ai-battle-afterpatch.png` | `cmd/wlgame`，2026-08-09 重拍 | real DOS/V `SINARIO.DAT`、`-seed 17`、正常鍵盤＋戰鬥指揮 | 濮陽攻城戰場 | scenario 0、player 0；事件佇列版本；攻方 559 兵、守方 100 兵 | original-derived battle data + remake tactical UI | Docker/Xvfb；正常戰術畫面 |
| `docs/images/wlgame-ai-battle-attack-afterpatch.png` | `cmd/wlgame`，2026-08-09 重拍 | real DOS/V `SINARIO.DAT`、`-seed 17`、正常鍵盤＋戰鬥指揮＋攻擊命令 2；戰略／戰術速度 16 | 濮陽攻城戰場／命令列 | scenario 0、player 0；196/6/28；攻方 112 兵、守方 100 兵、城壁耐久 1830 | original-derived battle data + remake tactical UI | Docker/Xvfb；攻擊命令輸入接縫與戰鬥中間幀；未作為完整結算證據 |
| `docs/images/wlgame-ai-battle-result.png` | `cmd/wlgame`，2026-08-09 重拍 | real DOS/V `SINARIO.DAT`、`-seed 17`、正常鍵盤＋戰鬥指揮＋攻擊命令；無 `-open-*` | 濮陽攻城戰場／戰後結果報告 | scenario 0、player 0；事件佇列版本 | original-derived battle data + remake result UI | Docker/Xvfb；守方勝、攻方 5590→0、守方 1000→1000、損害 0；按 Enter 才回寫 |
| `docs/images/wlgame-ai-postbattle.png` | `cmd/wlgame`，2026-08-09 重拍 | real DOS/V `SINARIO.DAT`、`-seed 17`、正常鍵盤＋戰鬥完成後 Enter；無 `-open-*` | 戰略地圖／戰後回戰略接縫 | scenario 0、player 0；196/6/29；戰略地圖與行軍事件 | original-derived map/chrome + remake state/UI | Docker/Xvfb；結果報告另由 `wlgame-ai-battle-result.png` 證明 |
| `docs/images/wlgame-save-replay-choice.png` | `cmd/wlgame`，2026-08-09 | 使用者提供 DOS/V `SAVE.DAT` 第 1 槽的可寫複本；無 `-open-*` | 真實地圖／「張飛 對 許褚」遭遇視窗 | loader 起始 0 年 0 月 0 日；第 600 幀為 0 年 1 月 1 日 | original-derived map/chrome + remake choice UI | Docker/Xvfb；存檔回放證據，不作正常時序 oracle |
| `docs/images/pc98-oracle-scenarios.png` | PC-98 原版，2026-08-09 | `wolong-dosboxx:latest`、`machine=pc98`、`noopen`、固定 `cycles=20000` | 四章劇本選單 | 日期 196/4/1、208/9/1、212/6/1、225/5/1 | 原版 PC-98 畫面 | `until:<md5>` 到站；閉迴路相對滑鼠 |
| `docs/images/pc98-oracle-cao-highlight.png` | PC-98 原版，2026-08-09 | 同上；曹操列第一下點擊 | 君主／軍師／據點清單 | 曹操列整列反白，尚未進入確認 | 原版 PC-98 畫面 | 兩段式選取第一下 |
| `docs/images/pc98-oracle-in-game.png` | PC-98 原版，2026-08-09 | 同上；完整 `NEW GAME` → 劇本 1 → 曹操 → `決定` | 劇本 1 大地圖 | `196年 4月 1日` | 原版 PC-98 畫面 | 已真正進入遊戲本體 |
| `docs/images/pc98-oracle-city-panel.png` | PC-98 原版，2026-08-09 | 同一大地圖狀態再點據點 | 據點資訊框 | 城兵 1,140；生產力 15,857；上昇值 0；防災值 100 | 原版 PC-98 畫面 | 可作 remake 資訊窗對拍；尚非存檔 oracle |
| `docs/playtest/06` outputs | PC-98 oracle | NEW GAME path | menu／world | md5 until + controlled mouse | pc98 original | byte-for-byte repeat evidence |

### 2026-08-09 事件 4／5 後續 TALK

`TestFundingTalkIndicesMatchRaw139E8Branches` 已固定 `sub_139E8` 的事件 4／5 base+offset：全額、指定等額／低額／零／超額與拒絕共 12 個分支；`cmd/wlgame` 的 `messages.go` 已接 `\4`／`\6`／`\7` 的已證實 formatter 行為。這只完成訊息索引與文字 modal 接縫，未將 PC-98 數值視窗欄寬、游標、逐頁動畫或跨平台 GUI 實跑列為完成。

### 2026-08-09 事件 2／3 前置外交 TALK

`TestQueuedDiplomacyChoiceTalkNotices` 固定事件 3 的 `CX=0x168`／#360 與事件 2 的
`CX=0x175`／#373；`TestDiplomacyTalkExpansionUsesOriginalRequestMarkers` 用真實 DOS/V
`TALK.DAT`／劇本驗證 `{3}` 君主名代換。這只完成進入三選一前的前置報告，不把後續回應、
PC-98 數值輸入視窗或逐頁排版列為完成。

### 2026-08-09 事件 2／3 主要 TALK

`TestDiplomacyTalkIndicesMatchRaw13902Branches` 固定建言 #364–#366／#377–#379 與主要
結果 #43–#45／#47–#49，並驗證外交超額輸入映射 response 2；真實 TALK 展開測試驗證
低額成功的 `\\7` 金額與 `\\3` 君主名；超額 response 2 與 Trust−30 另有 state 測試。
AH／信賴度次要訊息、PC-98 數值視窗與長程路徑
仍進行中；完整遊戲測試不作為本輪阻塞。

### 2026-08-09 事件 9 通知條件

`TestReleasedGeneralTalkOnlyTargetsPlayerFaction` 驗證 raw `sub_150D7` 的 #37 owner
條件；非玩家勢力／在野釋放不產生玩家通知。#409 空槽與原版肖像／逐頁呈現仍未驗。

補查確認 PC-98／DOS/V `TALK.DAT` 的 #409（`CX=0x199`）均只有資料上的空行、沒有文字／marker；
`TestReleasedGeneralRawFollowup409IsEmptyNoOp` 驗證 remake 不排入空白 modal。事件 9 的
可見文字接縫因此封閉；原版空呼叫時序與完整長程畫面仍不宣稱 parity。

### 2026-08-10 進言／說得信賴度切片

`sub_13830`、`sub_13C1E`、`sub_13B5A`、`sub_13BA9` 已以唯讀 DOS/V 反組譯證據固定
第一反應 `+20／−20`、多理由完成 `+10`、錯選 `−20`，以及 `byte_10D00` 的
`E0／90／20` 四段理由數。`internal/rules/persuasion` 的邊界測試與
`cmd/wlgame` 的直接反應／Trust 接線在 Docker 通過；這是局部垂直切片，不等於外交事件
2／3 的所有增減、事件 6／7 次要 TALK 或完整玩家長程 parity。

### 2026-08-10 M7 校訂文字 runtime 接線

- `translations/talk-dosv-corrected.json` 由 `talkdat.py correct` 產出，`talkdat_selftest.py`
  固定 60 筆修正、校訂後 round-trip 與 checked-in runtime 產出一致。
- `internal/assets/text.LoadJSON` 的 marker／尾端空行單測、完整 Docker/Xvfb Go gate、兩版
  raw TALK verify、deny-list selftest 與文件索引均通過；`wlgame` 短 smoke 已載入校訂表。
- 行寬警告、formatter／翻頁、原版同狀態畫面與跨平台 GUI 仍未驗收，因此本節不改變 packaging
  或推廣影片的未完成狀態。

### 2026-08-10 第 3／4／5 項定向驗收

- `sub_17C6E` 的數值狀態核心已由既有 `AmountEdit` 單測固定；`docs/re/13` 補上 PC-98
  event3 fixture、`sub_13B7E`／`sub_193E9` 輸入保存鏈與原始函式位址。數值面板仍只
  宣稱已量得的 `(80,176)`、`112×80`、`(88,200)` 3×6／16×16 幾何。
- `textdraw.WrapLines` 將 marker 展開後的 TALK 硬斷行按 ASCII 8 px／CJK 16 px 重新
  排版；`cmd/wlgame` 的 22 full-width cell、空列、標點與分頁測試已在 Docker＋Xvfb 通過。
- 特殊投射物保留發射兵 `+0x02 bit 0`，畫面層現在可選 raw `0x214`／`0x215`；事件 12
  MCH type 1／2 八相位已有資產測試，另以固定 8-frame fallback clock 加上回歸測試。
- 原版／remake event3 同狀態截圖已保存，但結果明確顯示原版肖像／物件疊層、PC-98 游標
  與數值外框資源尚未對齊；因此 original/remake screenshot review、目標平台封裝與影片 gate
  仍保持未通過。

### 2026-08-10 MCH／release gate 複核

- `wolong-go:20260809`、無網路 Docker＋有界 Xvfb：`go test -p=1 -vet=off ./... -count=1`、
  `go vet ./...`、Linux `cmd/wlgame` build 與 `tools/index.py generate/check` 通過；完整
  長程遊戲測試依使用者指示略過。
- 依 `tools/release.sh` 的實際矩陣，`wlsim`／`wlshot` 已在 Linux amd64/arm64、Windows
  amd64、Darwin amd64/arm64 建置，Windows `wlgame` 候選也建置成功；Linux／macOS
  `wlgame` 仍需目標平台原生建置。沒有把交叉編譯當成 GUI runtime 證據。
- `MMAP.MCH` type 1／2 的原版圖形已接入事件 12 火災／暴動；目前仍缺原版 timer／frame
  parity、事件 6／7 次要 formatter、事件 10、PC-98 數字／排版 parity、Windows／macOS
  GUI 實跑、原版／remake 同狀態對拍與正式封裝影片，因此 release gates 不全通過。

## Release gates

- [x] full test and static analysis in Docker（長程完整遊玩依使用者指示略過）
- [x] `tools/index.py generate/check` and documentation links
- [x] translation／encoding checks and 60 corrections review（`tools/m7_review.py`、TALK selftest、
  runtime 寬度 gate 與 6 張 DOS/V modal 代表幀通過；文字層的逐像素原版對拍尚未做）
- [x] normal player path without debug hooks（策略／戰術短路徑已驗）
- [x] writable save and pristine-original check
- [x] video-oracle／remake screenshot review with metadata
- [x] ⭐ 同狀態逐區對拍：拿原版存檔開同一個局面比像素。主畫面五區逐像素相同、
      戰場九區裡六區逐像素相同（`docs/playtest/37`、`40`；方法與判定在 `docs/spec/90`／`91`）
- [ ] packaged build smoke on target platforms
- [ ] deny-list scan and final dirty-tree review

### 2026-08-10 DOS/V 接線勘誤

本輪新增／取代下列過時描述，完整長程遊戲測試仍依使用者要求略過：

- `TestQueuedEvent10TalkNotice` 已固定事件 10 的 raw `Param`／高 byte General formatter
  邊界；`0x0A` producer 仍沒有已證實來源。
- `TestQueuedDiplomacySecondaryTalkConditions` 已固定事件 6 #72／事件 7 #76 的條件、
  順序與 `NoPortrait`；#72 因第二次呼叫缺 `DI=SP` 的 `\\2` payload，呈現層
  fail-closed，不把未知資料轉成城市。
- `TestDisasterObjectAnimationTiming` 已取代舊的 global 8-frame fallback 測試，固定
  typed object 的 timer=1、interval=16、dirty render 舊 phase、八相位遞增與清除。
- DOS/V 數值框的 `ICONGRF` 第 3 段 `0x14A0`／96×64 內框、下半部 3×6 靜態 button glyph
  與 `KI.EXE` `seg002:031B` 兩層 16×16 硬體 cursor 已接入；`cursor_test.go` 固定白 39、
  紅 56、透明 161，`gfx_test.go` 逐格固定 glyph 存在。最後 slot 移動分支、原版自然畫面
  對拍與目標平台 GUI 仍是未完成 gate。

### 2026-08-10 DOS/V 游標／按鍵 glyph 驗收

- `KI.EXE`：MZ header `0x200` + `seg002:031B` → file `0x1051B`；
  `sub_201E4`／`sub_2020C` 的兩層 32-byte mask 已以原始 byte 順序解碼。
- `ICONGRF.DAT`：`sub_17D0D`／`AX=4006h` → 第 3 段 `0x14A0`／96×64；下半部
  3×6、16×16 靜態 glyph 已直接繪製，不再被 vector／CJK fallback 覆蓋。
- Docker／無網路／有界 Xvfb 下，`./internal/assets/gfx`、`./internal/assets/library`、
  `./cmd/wlgame` 針對性測試通過；未跑完整長程遊戲，符合本輪範圍。

### 2026-08-11 事件 raw producer、DOS/V 自然邊界與目標 GUI

- 事件 6／7：`TalkNotice.RawFormatterWord`、valid sentinel、DOS/V `SS:[DI]` resolver
  與呈現層 fail-closed 已由 `TestSecondaryTalkUsesOriginalRawFormatterWord` 及 state
  fixture 固定；事件 6 只在已知 0x400-byte stack 範圍保留 raw word，事件 7 無 marker。
- 事件 10：`World.QueueEvent10` 只建立已證實的 `(general<<8)|0x0A`／raw `Param`，
  使用完整 256 格受控注入口；`TestEvent10ProducerWritesRawTalkPayload` 固定 payload、
  前 64 格滿時的搜尋與輸入邊界。IDA 直接 caller 仍未找到原版自然 producer，保留負證據。
- 自然 DOS/V：當時 Docker／DOSBox／Xvfb 只跑到「密碼輸入：第 09 頁」就停下，
  所以那一輪沒有做原版／remake 的自然逐像素 gate。
  （**密碼頁本身不擋**——空白確認就進開場，`docs/playtest/18`。）remake 短 smoke 截圖是
  `docs/images/wlgame-dosv-natural-remake.png`，SHA-256 為
  `8420d97955be60af16da403544b47e84b3f44363ef75f867930e022d1bc2f916`。
- 目標 GUI：Linux／Xvfb smoke 通過；Windows amd64 為 `PE32+ x86-64`，macOS amd64／
  arm64 為 `Mach-O`。沒有 Windows／macOS 原生桌面 runtime，故只關閉建置 gate，不
  宣稱平台 parity 或建立正式包。

### 2026-08-11 YouTube 自然畫面 oracle 對拍

- 使用者提供影片：[臥龍傳 呂布開局滅曹操](https://www.youtube.com/watch?v=af6xqcicXoI)。
  metadata 為 567 秒、478×360、30 fps；Docker 內擷取 20／80／160／240／320／400／
  480／550 秒代表幀，未保存原始影片。
- 80 秒幀 `yt-wolong-natural-80s.png` SHA-256：
  `d33fff8d664e24321274310287dce38b82c82cfb62f3d0427e70dfd5bd301e08`；去黑邊／縮放
  至 640×400 的 `yt-wolong-natural-80s-640x400.png` SHA-256：
  `c0217b8722bd44a22a112a2981b626126d5ee53d3e9f00498c6cbd018e08d6`。
- 影片與 remake 對照已確認 32 px banner、32 px command strip、左側 432×336（27×21）
  map、右側 208 px sidebar、192×128 minimap 與自勢力數值欄；`strategyhud.go` 已
  接入自然 HUD，remake 截圖 `wlgame-dosv-natural-remake.png` SHA-256 為
  `961e583915d2e0e7b65cd51f637ec214530b68040ba0da5770add4b35cb46e30`。
- 影片視覺／幾何對拍 gate 通過；因影片為有損縮放且與 remake 的日期／鏡頭不同，
  **影片這條路**不宣稱嚴格同狀態逐像素 diff。
  （同狀態逐像素是另一條路走出來的：拿原版存檔開同一個局面、640×400 原尺寸逐區比，
  見 `docs/spec/90`／`91` 與 `docs/playtest/37`／`40`。）
  Windows／macOS 原生 GUI gate 與正式包 gate 不受此項取代。
- 依使用者要求，另產出 `dist/promo/wolong-remake-yt-comparison.mp4`、自然畫面並排圖與
  差異圖；骨架調整前的 640×400 自然幀基線為 `AE=255003/256000 (99.61%)`、
  `RMSE=0.338208`。這封閉的是 YouTube／推廣片可見像素差異與畫面骨架比較，不把
  不同日期／鏡頭／狀態升格為原版同狀態 parity。完整量測見
  [`docs/promo/yt-remake-pixel-review.md`](docs/promo/yt-remake-pixel-review.md)。
- 2026-08-11 骨架調整後以 `wlgame-dosv-natural-remake-skeleton.png` 重拍：右欄的
  minimap 色標列、共用 8 px 分隔邊、君主／首都／軍師區、信賴度與黑底資源區已對回
  原版參考幀；最新 raw metric 為 `AE=249178/256000 (97.34%)`、`RMSE=0.329145`。
  新推廣片與對照片已重新產出，未解出的 reserve raw glyph 不冒充已解資產。

### 2026-08-11 短 parity gate 與事件 10 深度逆向

本節是目前最高優先的驗收狀態；較早章節保留為歷史驗證紀錄，不覆寫其當時的未完成
描述。

| 範圍 | 自動化入口 | 結果 | 證據界線 |
|---|---|---|---|
| `sub_1248A` | `internal/state` 的 `TestSub124FFMatchesRawSignedByteContract`、兩個 `MovingDisasterSub1248A*` | **PASS** | raw slot 半區、fixed-point drift、方向 byte、bounds／wrap；不把未知 type 3 語意補成已知 |
| 事件 2／3／4／5 | `cmd/wlgame/TestEvent2To5TalkBranchParityGate` + `TestEvent2To5FullTalkPageSampling` | **PASS** | 36 raw TALK 頁面／18 組雙頁回應；成功／拒絕／金額分支、marker、硬換行、寬度與五列分頁；不宣稱完整自然長程 |
| 事件 9 | `cmd/wlgame/TestEvent9ShortFixtureGate` + `TestEvent9LongNotificationRoute` + `internal/state/TestEvent9LongNaturalRoute` | **PASS** | 27 小時 bounded queue 在第 7／17／27 小時取出；玩家通知、非玩家／在野抑制與 #409 空白 no-op；完整自然劇本依要求略過 |
| M7 | `cmd/wlgame/TestM7CorrectedTalkLayoutGate` + `tools/talkdat_selftest.py` | **PASS** | 60 筆 checked-in correction、marker 展開、寬度與硬分頁；人工文意審查不由單測冒充 |
| 投射物 | `cmd/wlgame/TestProjectileParityGate` + `internal/rules/tactical` projectile tests | **PASS** | 一般水平／垂直、特殊 raw `0x214`／`0x215`、姿態位元與運動規則；同狀態原版像素對拍另列界線 |
| 事件 10 producer | `tools/ida_event10_producer.idc` → `docs/re/15-event10-producer.md`；`TestApproximateEvent10Producer*` | **PASS（原版 unknown＋remake substitute）** | 原版 dispatcher／consumer／writer／caller／data refs 已證實；remake 近似 producer 只對玩家俘虜產生 TALK `0x41／0x42`，可關閉，不冒充原版自然 producer |

#### 2026-08-11 M7 人工抽樣

`translations/corrections.json` 的 60 筆修正已逐筆閱讀，並以
[`docs/playtest/14-m7-review.md`](docs/playtest/14-m7-review.md) 留下語意、marker、硬行與
代表幀紀錄。`m7_review.py --check` 與 `TestM7CorrectedTalkLayoutGate` 只負責可重跑的
結構／寬度保護；人工審查結論仍不冒充原版同狀態 pixel parity。

#### 事件 10 raw 證據索引

- 輸入 `KI.EXE` SHA-256：
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- IDA `.i64` SHA-256：
  `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`。
- IDA Pro 9.4、DOS/V 線性位址、segment base `0x10000`；原始名稱與 operand 保留。
- 已證實 dispatch：`sub_131AE` `000131AE` → `CS:funcs_131E8` → `sub_13496`
  `00013496`；queue record 是 `Code`／`Param`、stride 4。
- 已檢查 writer caller：`sub_12FBF`／`sub_12FB1`／`sub_1301C` graph 的已列常數有
  `0x04`、`0x05`、`0x07`、`0x09`、`0x0C`、`0x0D`、`0x010C`、`0x020C`，沒有
  直接可證實的低碼 `0x0A`。這是負證據，不是對間接／外部寫入的全 binary 排除。

#### 2026-08-11 事件 10 近似自然 producer

`LoadScenario` 預設開啟 `World.produceApproximateEvent10`：月結後從玩家勢力的活俘虜
中至多選一名，沿已知 `General.Timer` 倒數與固定 RNG 邊界，排入
`(general<<8)|0x0A`／`Param=0x41`（逃走）或 `0x42`（歸降）。下一個每時 dispatcher
才消費，故仍保留「據點／軍團／物件／時鐘」與 queue 節拍順序。這是 **substitute／強推論**，
不是原版 `0x0A` producer 證據；`SetApproximateEvent10(false)` 可回到純 raw fixture。

自動驗收：`TestApproximateEvent10ProducerUsesKnownRawContract`、
`TestApproximateEvent10ProducerIsBoundedAndDisableable`、
`TestApproximateEvent10ReentersIdleClockConsumer`。

#### 2026-08-11 DOSBox／remake 可玩性專家驗證

| 範圍 | 執行方式／證據 | 結果 | 邊界 |
|---|---|---|---|
| DOS/V 原版啟動 | `wolong-dosboxx:latest`、`START.BAT`、`mouse_emulation=integration`、固定 `cycles=20000`；playtest 18 | **PASS（密碼頁後開場）** | 空白確認、`0000`、`1234` 均進入開場；⭐ 同狀態逐區對拍已做（`docs/playtest/37`／`40`），完整自然長程尚未執行 |
| PC-98 原版 oracle | `wolong-dosboxx:latest`、`machine=pc98`；`pc98-oracle-scenarios.png`、`pc98-oracle-in-game.png` | **PASS（oracle）** | 目前 headless image 缺 window manager，bus-mouse 完整重播仍待輸入橋接 |
| remake 正常策略 | 目前工作樹建置、`-seed 17`、無 `-open-*`；編成／行軍／速度／196/6/28 遭遇截圖 | **PASS** | 有一輪因等待／速度窗口越過遭遇；不把短 replay 當完整長程測試 |
| remake 存檔／讀檔 | 系統視窗 `4 → S → 1 → Enter → L → 1 → Enter`；overlay 88,832 bytes | **PASS** | 原始資料唯讀，驗證的是 remake overlay contract |
| remake 戰術 GUI | `-open-siege` current-build debug smoke；既有無旗標正常 tactical path | **PASS（smoke）** | ⭐ 同狀態逐區對拍另外做過：九區裡六區逐像素相同、戰場區 0.17%（`docs/playtest/40`）|

詳細報告與截圖 hash：[`docs/playtest/17-expert-dosbox-remake.md`](docs/playtest/17-expert-dosbox-remake.md)。

#### 可重跑命令

在 `wolong-go` Docker 容器內、`DISPLAY` 未設定時，執行：

```text
tools/parity_gate.sh
```

完整長程遊戲、DOS/V 的原版自然畫面，以及 Windows／macOS 原生 GUI 仍依
使用者決策不作為本輪小型 gate 的必要條件。

### 2026-08-11 三平台候選包與 Android 規劃（歷史起點）

- `dist/release-20260811/packages/` 已建立 Linux amd64、Windows amd64、macOS universal 三個主 `.tar.gz`，另有 Linux arm64 邏輯工具包；解包內容各自有 `README-RELEASE.md`、`translations/corrections.json` 與 `SHA256SUMS.txt`，不含原版資產。
- Linux 封裝執行 640×400、`seed=17`、30 幀 Xvfb smoke；輸出與最新 DOS/V 骨架相同，hash `45a68852335420dd7b22b4e240192dcd7a38fbbc62f72c8c59ec95acdc137b24`。
- Android 規劃起始紀錄：[`docs/mobile/android-plan.md`](docs/mobile/android-plan.md) 定義橫向 640×400 邏輯畫布、安全區、觸控抽屜、手勢與 pause/resume 不重複推進 clock；現況以後續 Android 原型 smoke 紀錄為準。

### 2026-08-11 idle clock 更新順序修正

- `World.TickMap` 現在對應原版 `sub_11CD0` 的完整順序：據點 `sub_13EFD` → 軍團
  `sub_125A3` → MCH 物件 `sub_12459` → 時鐘 `sub_11D8E`。
- `wlgame` 每個可見 map-loop 的第一個規則 tick 使用 `TickMap`；同一畫面的額外
  `g.speed` tick 使用不含物件的 `Tick`，因此物件動畫仍維持每個 map-loop 一次。
- `TestTickRunsCorpsBeforeClockAdvance` 驗證時鐘進位前的軍費／士氣讀取舊小時；
  `TestMovingDisasterSub1248AUsesOnlyLastHalfOfRawSlots` 改用
  `AdvanceMapObjects` 驗證 map-loop 亂數在物件更新前綁定。
- Docker 定向 state／GUI 測試通過；事件 10 的自然 `0x0A` producer 仍是限時封口的
  unknown，不影響本次更新順序修正。

### 2026-08-11 AppImage、經典再現影片與 Android 模擬器 smoke

| 範圍 | 執行方式／證據 | 結果 | 邊界 |
|---|---|---|---|
| Linux amd64 AppImage | `u5cht/appimage:latest` 建置；AppDir root `.desktop`／`AppRun`；deny-list；`APPIMAGE_EXTRACT_AND_RUN=1` + Xvfb | **PASS** | 仍需 Linux 主機的 X11／GL／ALSA；不含原版資料與字型 |
| 「經典再現」比較片 | `tools/promo_classic_revival.sh`；YouTube 代表幀＋remake 固定 `seed=17`；`ffprobe` | **PASS** | 60 秒、1280×720、H.264/AAC；是視覺展示，不是同狀態 pixel parity |
| Android shell APK | 固定 Docker 工具鏈；Android 35 `google_apis;x86_64`、KVM；1080×1920 physical／1920×1080 app | **PASS（原型）** | 只驗證安裝、啟動、橫向畫面與 dock 觸控；不是完整 Android 遊戲 |
| Android `CONTINUE` | 有界長按 `(500,980)`；`android-wolong-touch-after-continue.png` | **PASS** | 顯示 `TALK page 1/3`；尚未接完整 TALK 核心 |
| Android `MENU` | 有界長按 `(960,980)`；`android-wolong-touch-after-menu.png` | **PASS** | 顯示 `COMMAND DRAWER OPEN`；尚未驗返回鍵／完整手勢 |

Android 產物與限制詳見 [`docs/mobile/android-plan.md`](docs/mobile/android-plan.md)。模擬器容器已停止；
沒有把原版資產放入 APK，也沒有把這次 shell smoke 寫成 Android release gate。

### 2026-08-11 松崗繁中事件 10／idle gate

| 範圍 | 執行方式／證據 | 結果 | 邊界 |
|---|---|---|---|
| 穩定游標 UI 閘門 | `cmd/wlgame/TestIdleClockGateRequiresStablePointerAndNoCommand` | **PASS** | 首次座標、移動、按鈕／命令均停住；下一個穩定無輸入 frame 才允許世界更新 |
| 已注入事件 10 的每時節拍 | `internal/state/TestIdleClockDispatchesQueuedEvent10OnHourlyCadence` | **PASS** | 第 7 個每時邊界才得到 `Code=0x030A`／`Param=0x42` 的 TALK；不是原版 dynamic trace |
| substitute 月結 producer | `TestApproximateEvent10*` | **PASS** | `0x41`／`0x42`、raw queue 與狀態原子性固定；仍是 substitute／強推論 |
| 松崗繁中自然畫面 smoke | Docker/Xvfb，`seed=17`、速度 1、30 frame | **PASS** | hash `45a68852335420dd7b22b4e240192dcd7a38fbbc62f72c8c59ec95acdc137b24`；不是受密碼保護原版的同狀態 trace |

本輪唯一驗證 oracle 是松崗繁中版；`workplace/orig/dosv` 為既有資料夾名稱。PC-98 與其他
版本不再是本輪事件 10 的行為或 release gate。原版自然低碼 `0x0A` writer 維持 **unknown**，
不因上述 remake 功能通過而提升為已證實。
