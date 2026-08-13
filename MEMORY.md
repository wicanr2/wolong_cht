# 臥龍傳 remake — 代理持續記憶

> 最後更新：2026-08-12。這是快速恢復摘要；狀態以 `CONTEXT.md`、`docs/INDEX.md`、
> 目前程式與可重現測試為準。

## 2026-08-11 本輪恢復點

- 短 gate 入口是 `tools/parity_gate.sh`，必須在 `wolong-go` Docker 容器內執行；它
  覆蓋 `sub_1248A`、事件 2–5、事件 9、M7、投射物、事件 10 raw queue 與 TALK selftest。
- `sub_1248A` 只移動 raw object slots 16–31；`sub_124FF` 的 signed-byte／fixed-point
  行為、方向 byte、bounds／wrap 已有 state tests。
- 事件 10 原版 producer 是 **unknown**：IDA Pro 9.4 `.i64` 已完成 dispatcher／consumer／
  writer／direct-caller／data-ref 審查，沒有可證實的低碼 `0x0A` direct producer。這是
  負證據封口，不是允許猜測；詳見 `docs/re/15-event10-producer.md`、`RESEARCH-LOG.md`。
- `World.QueueEvent10` 是受控 raw fixture／劇本注入口，不得寫成原版自然 producer。
- `HANDOFF.md` 不存在且不得重新建立；進度只更新 `WORKLIST.md`。完整長程遊戲、密碼頁
  後自然畫面逐像素對拍、Windows／macOS 原生 GUI 與正式三平台包仍是明確邊界。

## 先記住

1. 先讀 `CONTEXT.md` §7.0，再看 `docs/INDEX.md` 斷言總表；不要因為「我不記得」就重推已解問題。
2. 反組譯固定先用 IDA Pro 9.4；保留原始函式名、全域名、線性位址、偏移與運算元，另加分級註記。
3. `dosv` 是松崗 DOS/V 繁中版，`pc98` 是 PC-98 日文版；不要把一版的結論直接外推到另一版。
4. 原版素材唯讀且不進發行；存檔採「改寫已解欄位」，未知 bytes 要 byte-for-byte 保留。
5. 所有分析／建置／測試／抓圖／GUI 都必須在 Docker；本機只做 `docker`、`git`、狀態檢查與檔案編輯。

## 目前基線

- 目標：Go／Ebiten 乾淨重寫《臥龍傳－三國制霸之計》，保留松崗繁中原文並以日文原版校訂。
- M0–M2 完成，M3–M6 大致完成；格式解碼、即時時鐘、月結、內政、外交、說服、政略 AI、
  戰術核心規則、原版道路行軍、存檔 round-trip 與主要畫面已有實作／歷史證據。
- 主要剩餘：M7 校訂後的畫面／排版驗收、投射物完整對拍、原版同狀態對拍，
  以及 M8 的目標平台建置／實機／發行驗收。
- 本輪已修正 `tools/go.sh` 的 Docker 資源／網路限制與 `tools/release.sh` 的主機 `go env`／`python3`
  違規；新增 `tools/py.sh` 作為固定 Docker Python 入口。
- M8 限定 smoke 已通過：`windows/amd64` 交叉產物檔頭為 PE、Linux 本機產物為 ELF、deny-list 通過，
  一年期 `wlsim -check` 無不變量違反，`wlgame` 在 Docker/Xvfb 由真實劇本啟動並產生畫面。
- `-shot` 的終止條件已改成「達到目標幀後下一次 Draw 成功寫 PNG，再由 Update 結束」；
  2026-08-09 直接啟動 `dist/linux-amd64/wlgame` 的 120 幀 Xvfb smoke 產出 74,256 B PNG，
  避免只跑 source build 或因幀邊界漏圖。
- 本輪把 `World.SaveInto` 接進 `wlgame` 系統視窗：指定 `-save-file` 後以 `S`／`L` 操作四槽 overlay；
  `internal/savepath` 做 fail-closed 原始路徑檢查，儲存用同目錄暫存檔加改名，讀取會重新掛道路圖／戰術來源。
  `4→S→Return` 的 Docker/Xvfb smoke 產出 88,832 B，前進後可見槽位差異，重新啟動讀到 196/4/8，原始雜湊未變。
  `Trust` 已由 IDA `seg000:10D00`／`byte_10D00` 與 `sub_18CAE`／`sub_18CFF` 定位為每區塊
  `+0x10` 的 u8，已接入讀寫與 0…255 鉗制測試；事件佇列 `+0x52C0` 的 256 筆原始
  `Code`／`Param`、每十次節拍、月度壓縮，以及事件 1 宣戰／2 合作產生器與狀態部分／3 停戰產生器與狀態部分／8 遷都（含 `sub_14502` 軍團同步）／9 指定武將釋放狀態／13 信賴度 −50
  handler 也已接入測試；Player 的 `+0x0D`／`+0x0F` 雙欄讀寫也已接入；事件 2／3 的玩家三選一接縫與 `sub_17C6E` 數值編輯語意已接，完整接受／PC-98 數值輸入畫面／原版訊息 UI、事件 4–7、10–12、事件 9 的通知／完整流程
  的效果與完整玩家 UI 仍未完，
  不可把這個切片宣稱為完整原版存檔 parity。
- 本輪另補上玩家遭遇的正常決策接縫：`World` 先掛 `EncounterChoice`，選「戰鬥指揮」才建
  `PendingBattle`，選「委任」走既有 `combat.Resolve`；選單期間時鐘不動。規則測試與真實素材
  Xvfb 選單截圖已通過；固定種子 17 的真實 `wlgame` 鍵盤路徑也已從編成／行軍走到
  「呂布 對 曹操／攻城／戰鬥指揮／委任」，證據為 `docs/images/wlgame-ai-normal-encounter.png`。
- 2026-08-09 已把戰術核心切片接回同一條正常玩家路徑：選「戰鬥指揮」後進入攻城戰場，
  再送出攻擊命令；`internal/rules/tactical` 依 IDA `sub_1B618`／`sub_1B6BC`／`sub_1B97E`
  實作一般近戰、大將命中與一般投射物命中，並修正 `sub_1B00D` 到達後才消費繞路點。
  `sub_1A83F` 支持隊長死亡後七名部下退卻；清除該隊待機兵是避免 remake 補兵狀態死鎖的強推論，
  已在 `docs/re/11-tactical-battle.md` 與 `RESEARCH-LOG.md` 分級記錄。證據見
  `docs/playtest/09-wlgame-normal-tactical-path.md` 及三張 `wlgame-ai-*-afterpatch.png`。
  另以 `TestNormalScenarioTacticalBattleTerminates` 完成正常真實攻城的狀態層結算
  （第 549 幀、守方勝、攻方 0／守方 100）；`sub_1AD7F` 的 `CH=0x20` 攻擊分支、
  GUI 戰後結果報告已由 `docs/images/wlgame-ai-battle-result.png` 留證（攻方 5770→0、
  守方 1000→1000、攻城損害 0）；原版同狀態對拍仍未完成；
  `docs/images/wlgame-ai-postbattle.png` 已證明正常 GUI 按 Enter 回到戰略地圖。
- 已用真實 `SINARIO.DAT`、不帶 `-open-*` 走通 `A` 編成／`M` 目的地／`=` 行軍：
  開局資源只允許一槽步兵，確認後「曹操 向 虎牢關 行軍」，日期前進一天；證據見
  `docs/playtest/08-wlgame-normal-strategy-path.md`。同一路徑另選袁術的汝南，
  由 253 條原版道路抵達後自動攻打城兵；`TestNormalScenarioMarchIntoGarrison` 與
  `docs/images/wlgame-normal-garrison.png` 已驗證；敵方軍團遭遇正常選單另由固定種子
  17 的 `wlgame-ai-normal-encounter.png` 驗證。
- 使用者提供的 DOS/V `SAVE.DAT` 第 1 槽可複製成 overlay；不帶 `-open-*` 回放 600 幀
  後顯示原版風格的「張飛 對 許褚／攻城／戰鬥指揮／委任」選單，證據為
  `docs/images/wlgame-save-replay-choice.png`。但 `state.LoadScenario` 記錄起始時鐘為
  0 年 0 月 0 日，故這是隔離存檔與 UI 接縫證據，不是正常 AI 來源或時序 oracle；
  必須保留此 fail-closed 界線。
- 2026-08-09 已接上政略 AI 的可重播狀態切片：`City.Neighbours` 載入原始四槽，
  `sub_12C52` 候選排序、`sub_12DB8`／`sub_12DF3` 交友度漂移、`sub_12EFB` 三閘、
  `sub_13639` 雙向交戰值，以及 `CS:6C4C` 六槽編成表都已進入 `World`。真實劇本固定
  種子 17 跑六個月，`TestStrategicAIScenarioOneProducesEnemyWarPath` 驗到宣戰 5、
  編成 4、戰鬥 4 與逐 tick 不變量；事件 1／2（產生器、合作狀態、三選一接縫與數值編輯語意）／3（產生器、停戰狀態、三選一接縫與數值編輯語意）／8／9（指定武將釋放狀態）／13 已經過 queue handler，事件 2、事件 3 的完整接受／PC-98 數值輸入畫面／原版訊息 UI、事件 4–7、10–12、事件 9 的通知／完整流程、
  多軍團請求、完整原版行軍狀態機仍未完成。
- M7 目前是 `translations/corrections.json` 60 筆；`tools/talkdat.py correct`／`talkdat_selftest.py`
  已可安全套用 60 筆 `fix`。本輪已完成 #0–#1021 第一輪文意讀取；`#751` 採最小修正，在既有「我可不會這樣就輸的」中「我」後補 `{1}`；
  同一戰場台詞池的 `{1}` 位置提供強證據，但完整語氣與畫面排版仍未驗收。`#321` 已由 IDA 證實 `{6}` 處理器（handler）會消耗一個參數並
  左移 `DX`；`#192`–`#195` 則由 `sub_13BA9`／`sub_13C99` 證實是三變體直取槽位，
  以既有繁中句子完成重排。
- 明列未解的小項：`ICONGRF` 段 1、視窗內龍紋、`MMAP.MCH`／`BATTLE.MCH` 語意、DOS/V 音源、
  協力第 ④b 道閘、`sub_135AB` 的 `0x24` 門檻；詳見 `CONTEXT.md` §7.0。

## 可靠證據入口

- 即時模型：`docs/mechanics/15-realtime.md`、`docs/re/06-game-clock.md`。
- 月結／每時更新／亂數：`docs/re/07-monthly-settlement.md`、`docs/re/08-hourly-update.md`、
  `docs/re/10-rng.md`。
- 戰鬥：`docs/re/09-combat.md`、`docs/re/11-tactical-battle.md`、`docs/mechanics/30-combat.md`。
- 文件狀態與斷言：`docs/INDEX.md`；格式總覽：`docs/formats/`。
- PC-98 正常玩家 oracle：`docs/playtest/06-mouse-solved.md`、`docs/playtest/07-in-game-oracle.md`。

## 2026-08-09 最新增量：事件 4／5

- 月結後已接 `sub_15715`／`sub_1578F`：事件 4 用據點編號、事件 5 用勢力編號，Param 都是要求金額；事件 4 的 `Growth` 必須先還原成存檔 `+0x10` 的 `Growth+100`。
- 已接 `sub_132A9`／`sub_132E9` 的處理時重新驗證、`sub_139E8` 的 500 初始下限／30,000 自訂上限／`amount/128` 官員經費／玩家扣款／拒絕無副作用。
- `World.PendingFunding`／`ResolveFunding`、`EditFundingAmount` 與 `cmd/wlgame/funding.go` 是目前玩家 UI 接縫；`sub_17C6E` 的數值核心已由 `TestRawAmountEditorSemantics` 固定，原版 PC-98 掃描碼／數字視窗與 TALK.DAT 訊息仍未對拍。`TestFundingRequestGenerators`、`TestQueuedFundingChoice` 與全量 Go vet／test 已通過。

## 2026-08-09 最新增量：事件 6／7

- 事件 6／7 的處理端已接入 `internal/state/events.go`：事件 6 的來源是外交官回報停戰的對方勢力，玩家付款給對方；事件 7 的來源是協力方，Param 低 byte 是第三方，玩家付款給協力方後由協力方對第三方宣戰。
- 兩個 handler 都在處理時重新確認回報勢力仍有 `Faction.Diplomat`；缺少外交官時 fail-closed。`TestQueuedDiplomacyReportHandlers` 已固定付款方向、停戰／宣戰、交戰值與失效分支。
- 尚未完成：`sub_164F1`／`sub_16623` 的完整玩家進言選單與 queue 產生路徑、原版反應／TALK.DAT 訊息、事件 10–12 與完整長期外交 oracle。

## 2026-08-09 最新增量：玩家進言 producer 接回

- `cmd/wlgame/advise.go` 已把 `persuasion.FirstReaction` 接入玩家流程：敵對同意直接走 `World.ApplyPlayerHostility`；停戰／協力同意分別寫事件 6／7。
- 協力進言依原版先選協力方、再選侵攻目標；同一家仍可選，交給 `SameFaction` 分支拒絕。
- `World.queuePlayerEvent` 對應 `sub_1301C` 的 `BL=14h`：從 `eventCursor + 0x14*4` 掃完整 256 格；事件 6／7 的 Code／Param payload 與重複防護由 `TestPlayerDiplomacyProducers` 固定。
- 尚未完成的是 `TALK.DAT` 反應／外交官回報訊息的逐句內容、逐頁排版，以及事件 10–12／事件 9 通知與長期外交 oracle；不可把 producer 接縫誤記為完整文字 parity。

## 2026-08-09 最新增量：事件 10–12 災害 dispatch

- 月結產生端已將暴風雨排成事件 11、火災／暴動排成 `0x010C`／`0x020C`；事件 12 Param 是 runtime city record 位址，不是 city ID。
- 事件 10 保持訊息-only；事件 11／12 接入不序列化的 runtime `+0x15` marker、`sub_14269` 的持久效果與事件 12 的 6..13 格延遲清除。
- `TestQueuedDisasterAnimationHandlers` 與 `TestDisasterMarkerAppliesRawPersistentEffects` 通過；仍缺 TALK.DAT／物件動畫畫面與事件 10 訊息 parity。

## 下一輪優先順序

1. 保持已通過的 Docker 基線：文件索引、Go vet／test、deny-list、`git diff --check`，每次改動後重跑。
2. 以 `tools/talkdat.py` 重跑 `translations/corrections.json` 的 60 筆（全部套用；#751 為最小插入修正），
   並把剩餘時間放在校訂後的硬換行／行寬、formatter 與畫面抽樣驗收。
3. 以 PC-98 oracle 做正常玩家路徑與同狀態 remake 畫面／狀態驗收；敵方 AI 正常鍵盤路徑
   已在狀態層完成正常攻城結算且 GUI 可顯示戰後結果並回戰略地圖；`sub_1AD7F` CH=0x20 行為已有獨立測試，下一步補 `+0x1E` 來源／完整投射物對拍與原版／remake 同狀態對拍；不要用 debug 傳送、強制勝利或
   固定秒數假設跳過不定長度開場，改用 `until:<md5>`、`Ctrl+F10` 與閉迴路 `clickat`。
4. 依實際缺口補 M8：目標平台自己建 Ebiten、包裝、deny-list、可寫存檔與原始素材隔離，留下證據。
5. 小型 RE 項目只有在 `docs/INDEX.md` 確認尚未解且有明確問題時才做；每項同步寫 `docs/re/` 和
   `docs/mechanics/`，不把假說直接接進規則。

## 不要重犯的陷阱

- 固定 `wait:100` 不代表到達選單；開場長度不固定，必須以畫面雜湊判定。
- `0x80`–`0x83` 是道路方向碼，不是 autotiling 圖塊；`0xFF` 也可能是「無」而非布林旗標。
- 拓樸相同不代表沿路路徑相同；原版道路逐條比對，不能用 BFS 最短路冒充等價。
- 不要把解碼器用錯編碼（`cp950`／`cp932`）造成的字誤判成原版錯字。
- 不要把「解出機制但不懂意圖」當成「不能實作」；照抄已證實機制並標註原版怪行為。
- 不要在 `CONTEXT.md` 留著已不存在的工具或已解的閘；完成一項就更新狀態、索引與證據。

## 工作樹與容器

本輪基線開始時 HEAD 為 `e01c667c00cc120d4d819c591d33b701e271bd9e`、分支 `main`；目前工作樹
包含本輪新增交接文件、戰術規則／測試／playtest 證據與 `tools/{check,go,py,release}.sh`／`tools/release_fs.py` 修改，尚未 commit。原始資產與 IDA 研究資料依
`CLAUDE.md` 規則唯讀。每輪結束要重新檢查專案相關容器，
只清理由該輪建立且不再需要的容器；不要碰其他專案的 `cool_bose`、`portainer` 等資源。

## 2026-08-09 最新增量：事件 9 通知觀測

- 事件 9 的事件字高 byte 是武將索引；`World.dispatchQueuedEvent` 成功呼叫 `releaseGeneral` 後，寫入 `Event.ReleasedGenerals`，`wlgame` 顯示通用 `<武將名> 已釋放`。
- `TestQueuedEventReleaseGeneral` 已固定存活俘虜方與已滅原勢力方的狀態回寫，新增欄位只作 UI 觀測，不列入存檔。
- `TALK.DAT` 0x25 原句已由 `cmd/wlgame` 取用並代入武將名；原版 `sub_150D7` 的勢力通知對話框／排版與完整事件流程仍未完成，不要把句型取用記成完整訊息 parity。
- 證據：`KI.EXE.asm` SHA-256 `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`，IDA Pro 9.4 `seg000:13485`／`150D7`；runtime 在 `wolong-go:20260809` Docker。

## 2026-08-09 最新增量：戰術 raw `PlaneHigh`

- `sub_1B4EA`／`sub_1B0D3`／`sub_1B116`／`sub_1B732` 已確認兵記錄 `+0x1E` 的 0／上移 `0x10`／回地面 0／換位交換寫入端；`sub_1B240` 另確認 `+0x00 bit 1` 是堆疊至少 4 層的 `HighTerrain`，不是 `Climbing` 同義欄位。
- `sub_1A85B` 鎖敵懲罰的 raw 分支已接成 `targetPlanePenalty`；`sub_1AC55`／`sub_1ACA4` 的特殊投射物條件已接成 `PlaneHigh` 與 max-axis ≤2。`Climbing` 只保留給未經 `Place` 的舊測試資料。
- 新測試：`TestRawPlaneHighAndTerrainFlag`、`TestSpecialProjectileUsesPlaneHighAndMaxAxisDistance`、`TestLockOnPlanePenaltyMatchesRawBranches`；Docker 映像 `wolong-go:20260809` 內戰術套件通過。完整投射物／動畫與同狀態原版對拍仍是未完成邊界。
- 證據：DOS/V `KI.EXE.asm` SHA-256 `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`、`KI.EXE` SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`；`ida-pro-9.4-ver2:uidfix-v1`／IDA Pro 9.4，位址基準為 IDA 線性位址；匯出檔在 `workplace/ida/dosv/func-sub_*.txt`。

## 2026-08-09 最新增量：投射物 raw 更新與畫面接縫

- `sub_1B941` 的實體更新順序已確認為先 `sub_1B97E`、再 `sub_1BA2E`、後 `sub_1BAB7`。remake 已接方向、三維格索引、固定點高度、速度夾限／衰減、障礙清除與上升／下降威力變化。
- `TestProjectileRawDirectionGridAndHeightPower`、`TestProjectileChecksCurrentCellBeforeMoving`、`TestProjectileStopsAtSolidLayerAfterMoving` 通過；CH=0x20 的既有測試現在也驗證下降加成後的真實步兵傷害，而非固定 `0x20/4`。
- `Battle.Projectiles()` 與 `cmd/wlgame/battleview.go` 已接戰場內標記；使用 248×192 原生區域的空間標記，不冒充尚未解出的原版投射物圖形。普通箭 `sub_1AD2D` 初始速度公式與 `sub_1ECE0`／`sub_1EC82` RNG 已接入並有單測；原版同狀態時序／畫面對拍仍未完成。
- 證據：`KI.EXE.asm` SHA-256 `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`；IDA Pro 9.4 `seg000:1B941`／`1B97E`／`1BA2E`／`1BAB7`；原始 `.i64` 以唯讀掛載並在容器內副本匯出，runtime／測試使用 `wolong-go:20260809`。

## 2026-08-09 最新增量：事件 11／12 災害 marker 持久效果

- `sub_13EFD`（`00013EFD`）在 `00013F5A` 於 `sub_14194` 後呼叫 `sub_14269`（`00014269`）；事件 11／12 的城市 `+0x15` marker 會在據點輪轉時先扣 `+0x11` 防災值，防災不足才以差額扣 `+0x10` 上昇值存值、`+0x0E` 生產力與 `+0x13` 城兵。
- `World.applyCityDisasterEffect` 已接入，保留生產力 16 位元減法與其他欄位的原版下限；新增 `TestDisasterMarkerAppliesRawPersistentEffects` 驗證防災護盾與不足分支。物件動畫、事件 10 訊息與完整原版畫面／長期 oracle 仍未完成。
- 證據：DOS/V `KI.EXE.i64`／`KI.EXE.asm`／`KI.EXE`，`KI.EXE.i64` SHA-256 `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`，`KI.EXE.asm` SHA-256 `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`，`KI.EXE` SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`；IDA Pro 9.4／IDA 線性位址；runtime／測試使用 `wolong-go:20260809` Docker。

## 2026-08-09 最新增量：事件 TALK 通知與 modal 接縫

- IDA 直接證據：`sub_1237E`（`0001237E`）對玩家城市使用 `CX=0x46`（TALK #70）並以 `DI` 傳城市記錄供 `\\2`；`sub_134B1`（`000134B1`）分派 #70／#71／#72；`sub_13507`（`00013507`）使用 `CX=0x33`（TALK #51）。輸入 `KI.EXE.asm` SHA-256 `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`、`KI.EXE` SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`；工具 `ida-pro-9.4-ver2:uidfix-v1`／IDA Pro 9.4，均為 IDA 線性位址。
- `Event.TalkNotices` 已保留原始 TALK index、city/general ID；`messages.go` 以 TALK.DAT 原始行邊界與 Big5／`{1}`／`{2}` 代換顯示 modal，modal 期間 `timeRuns()` 停止世界時間。事件 9 釋放武將也進入同一佇列（#37／`0x25`）。
- `TestQueuedTalkNotices` 與 Docker 目標測試通過；這只是通知資料／呈現接縫，不代表事件 10、完整 TALK 翻頁／肖像／動畫、災害物件動畫、原版畫面對拍或長程驗收已完成。
- `-open-message` 的 Docker/Xvfb 實機抽樣已產生 `docs/images/wlgame-event-modal.png`：TALK #70 顯示「許昌發生了暴風雨。」；30 幀後日期仍為 196/4/1。TALK marker 必須使用 ASCII `'1'`／`'2'`，不能用 state 數值 1／2。

## 2026-08-09 最新增量：事件 6／7 主要 TALK 回報

- `sub_13327`／`sub_13388` 先顯示 #57，再依 `sub_136C4`／`sub_13712` 的 response 產生 #43–#45／#47–#49；`DX` 是金額，`\\3` 是回報方君主名，`\\7` 是十進位金額。
- state `TalkNotice` 現保留 `Faction`／`Amount`；GUI 仍保留原始 TALK 行邊界，使用 `World.LordName` 與 ASCII 十進位 marker。`TestQueuedDiplomacyReportTalkNotices` 通過：事件 6 #57→#44／14,000，事件 7 #57→#48／15,000。
- 證據 `.i64` SHA-256 `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`；IDA Pro 9.4、IDA 線性位址 `00013327`、`00013388`、`000136C4`、`00013712`、`00013C3D`、`00010904`、`00010984`。次要 AH／`sub_13DC9`、原版數值欄寬／逐頁排版仍未完成。

## 2026-08-09 最新增量：M8 發行矩陣

- Docker `wolong-go:20260809` 已重建 `dist/`：Linux／Windows／Darwin 純 Go `wlsim`／`wlshot`，Windows `wlgame`／`wlview`，Linux 原生 `wlgame`／`wlview`。
- 十六進位檔頭確認 ELF／PE／Mach-O；`denylist.py dist` 通過。封裝 Linux `wlgame` 以 Xvfb `:101`、唯讀原版素材／字型、120 幀 `-shot` 成功輸出 74,269 bytes PNG。
- 尚未完成：Windows／macOS GUI 目標 runtime 實機、有效時序原版／remake 同狀態對拍與 M7 1,022 則文意全量審查。

## 2026-08-09 最新增量：投射物 raw 更新與畫面接縫

- `sub_1B941` 的實體更新順序已確認為先 `sub_1B97E`、再 `sub_1BA2E`、後 `sub_1BAB7`。remake 已接方向、三維格索引、固定點高度、速度夾限／衰減、障礙清除與上升／下降威力變化。
- `TestProjectileRawDirectionGridAndHeightPower`、`TestProjectileChecksCurrentCellBeforeMoving`、`TestProjectileStopsAtSolidLayerAfterMoving` 通過；CH=0x20 的既有測試現在也驗證下降加成後的真實步兵傷害，而非固定 `0x20/4`。
- `Battle.Projectiles()` 與 `cmd/wlgame/battleview.go` 已接戰場內標記；使用 248×192 原生區域的空間標記，不冒充尚未解出的原版投射物圖形。普通箭 `sub_1AD2D` 初始速度公式與 `sub_1ECE0`／`sub_1EC82` RNG 已接入並有單測；原版同狀態時序／畫面對拍仍未完成。
- 證據：`KI.EXE.asm` SHA-256 `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`；IDA Pro 9.4 `seg000:1B941`／`1B97E`／`1BA2E`／`1BAB7`；原始 `.i64` 以唯讀掛載並在容器內副本匯出，runtime／測試使用 `wolong-go:20260809`。
## 2026-08-09 最新增量：`sub_17C6E` 數值編輯核心

- IDA 線性位址 `00017C6E` 呼叫的操作表已確認：`00017DA5` 追加數字、`00017DC3` 追加 `00`、`00017DDD` 退位、`00017DEC` 還原初值、`00017DF1` 清零、`00017DEA` 結束並保留目前值；上限由呼叫端傳入，本專案事件 2／3／4／5 為 `0x7530`。
- `AmountEdit`、`EditDiplomacyOfferAmount`、`EditFundingAmount` 與 `TestRawAmountEditorSemantics` 已接入；跨平台 `wlgame` 指定金額列支援數字鍵、退格、Insert、Delete、Home。這是狀態／輸入語意接入，不是 PC-98 掃描碼或原版數字視窗 parity。
- 原始輸入為唯讀 DOS/V `KI.EXE.i64`／`KI.EXE.asm`／`KI.EXE`；工具 `ida-pro-9.4-ver2:uidfix-v1`／IDA Pro 9.4，位址均為 IDA 線性位址；Go 驗證在 `wolong-go:20260809` Docker 內完成。
## 2026-08-09 最新增量：事件 4／5 前置回報通知

- `sub_132A9`／`sub_132E9` 在 `sub_139E8` 前分別顯示 `CX=0x38`／TALK #56 與 `CX=0x39`／TALK #57；marker 目標已由原始 `DI`／堆疊資料確認。
- `dispatchQueuedEvent` case 4／5 已建立 `TalkNotice`，`wlgame` 會先顯示通知 modal，確認後才進入撥款三選一；`TestQueuedFundingInitialTalkNotices` 已通過。後續 `sub_139E8` 訊息池與 PC-98 數字視窗仍是未完成。

## 2026-08-09 最新增量：指定金額 outcome 與外交亂數時序

- `sub_13902` 的外交指定金額若高於初始 `DX`，回傳碼為 3，外層 `AL >= 2` 不執行外交收尾；`ResolveDiplomacy` 現在以 pending 的 `InitialAmount` 驗證，不能讓玩家輸入超額後仍宣戰／停戰。
- `sub_139E8` 的撥款指定金額則允許高於原始要求；只有 0 回傳無效碼 2，不改官員經費或玩家資金。這兩個事件不能共用同一個「超額拒絕」規則。
- `beginDiplomacy` 已保存初始 terms 結果，`ResolveDiplomacy` 不再重跑可能含平手亂數的代表政治計算，避免 modal 確認時多消耗一筆 RNG。測試為 `TestDiplomacyAndFundingAmountOutcomeBounds`，輸入／測試仍只在 `wolong-go:20260809` Docker 內執行。

## 2026-08-09 最新增量：事件 4／5 後續 TALK 分支

- IDA `CS:08A4` table 的可用 handler：`000108B2`=`\1` 武將、`000108DB`=`\2` 據點、`00010904`=`\3` 君主、`00010939`=`\4` 玩家軍師、`0001097E`=`\6` 不可見欄位控制、`00010984`=`\7` 十進位金額；位址基準是 IDA 線性位址。
- 事件 4 `base=0x116`、事件 5 `base=0x13F` 的 `sub_139E8` 結果／收尾 offset 已接入 `cmd/wlgame/funding.go`；`TestFundingTalkIndicesMatchRaw139E8Branches` 固定 12 個分支。
- `messages.go` 的 `\6` 必須代換成空字串，不可把它當作文字 marker；它原版只改 `DX`／消耗 formatter stack，真正像素欄位位置仍未移植。

## 2026-08-09 最新增量：事件 2／3 前置外交請求 TALK

- 原版 `sub_138C7`（IDA `000138C7`）使用 TALK #360（`CX=0x168`）顯示停戰請求；
  `sub_138E6`（`000138E6`）使用 TALK #373（`CX=0x175`）顯示協力請求；`{3}` 是請求方
  君主名。不要把事件字高直接當成 `{3}` 文字。
- `beginDiplomacy` 已在 pending 三選一前 append `TalkNotice`；state 與真實 TALK 展開測試
  分別是 `TestQueuedDiplomacyChoiceTalkNotices`、
  `TestDiplomacyTalkExpansionUsesOriginalRequestMarkers`。
- 只完成前置請求通知；後續結果 TALK、PC-98 數值視窗與逐頁排版仍未完成。

## 2026-08-09 最新增量：事件 2／3 主要 TALK 分支

- `sub_13902` base+4/+5/+6 是玩家三選一後的建言句：事件 3 #364–#366、事件 2
  #377–#379；主要結果由 `sub_13C3D` 形成 #43–#45／#47–#49。
- `diplomacyTalkResponse` 把正值且不超過初值映射為 response 1／`\\7`，零映射 response 0，
  超額與拒絕映射 response 2；這只決定 TALK 與原版分支，不取代 state 的 outcome 測試。
- 主要接線測試：`TestDiplomacyTalkIndicesMatchRaw13902Branches`、
  `TestDiplomacyTalkExpansionUsesOriginalRequestMarkers`。AH／信賴度次要回覆仍未知。

## 2026-08-09 最新增量：事件 9 玩家通知條件

- raw `sub_150D7` 只有釋放後所屬勢力等於 `byte_10CFF` 玩家勢力時才顯示 TALK #37；
  不可對所有 `ReleasedGenerals` 無條件排入通知。
- `enqueueEventMessages` 已加 owner 條件；測試為 `TestReleasedGeneralTalkOnlyTargetsPlayerFaction`。
- 後續 `CX=0x199`／TALK #409 是空槽／未知 formatter 路徑，保持未接，不猜測文案。

## 2026-08-10 最新增量：災害 marker 唯讀呈現

- state 新增 `DisasterMarkerAt`／`StormAreaSnapshot`，供 `wlgame` 讀取 runtime marker；
  不序列化、不暴露可寫陣列。`TestDisasterMarkerReadOnlySnapshots` 通過。
- 主畫面已接火災／暴動／暴風雨的向量替代標記與暴風雨範圍輪廓；這只證明狀態到畫面的
  接線，不能升格成原版災害物件圖形／動畫完成。GUI 短測試與建置在 Docker／Xvfb 通過。

## 2026-08-10 最新增量：事件 10 暫不猜接

- `sub_131AE` dispatch table 確認低位碼 `0x0A` → `sub_13496`；handler 只呼叫訊息入口
  `sub_18810`，尚無 producer／TALK 槽位證據。
- 已檢查 `sub_12FBF`／`sub_12FB1`／`sub_1301C` 的策略 caller，沒有找到 0x0A 寫入；
  記為負證據／未知，不把其他 UI／戰術的 `0Ah` 常數當事件 10。保留為接班下一個 raw RE
  項目，不能因此宣稱全量串接完成。

### IDA xref 補查記憶

- `ida-pro-9.4-ver2:uidfix-v1`／`tools/ida_xref.idc` 在 `.i64` 副本確認：
  `sub_12FBF` 7 callers、`sub_12FB1` 5 callers、`sub_1301C` 4 callers，直接 caller
  沒有新增事件碼 `0x0A`。`sub_12286` 是 `0x010C／0x020C`，`sub_122DB` 是 `0x0B`，
  `sub_12D3A` 是事件 8。
- 保留界線：這只封閉已查的直接 xref；函式指標／間接寫入仍 unknown。不要建立事件 10
  producer、TALK index 或泛用訊息。

## 2026-08-10 最新增量：進言／說得 Trust 與外交超額分支

- 原版 `sub_13830`（IDA 線性位址 `00013830`）：`AL=1` 直接成功 `+20`，`AL=0／3`
  失敗 `−20`，`AL=2` 多理由完成 `+10`；`sub_13BA9` 錯選理由 `−20`，撤回／重複不變。
- `sub_13C1E`（`00013C1E`）讀 `byte_10D00`，`E0／90／20` 分段要求 1／2／3／4 個
  成立理由。`persuasion.Situation.Trust`、`Begin`、`Session.Offer`、
  `ReactionTrustDelta` 與 `cmd/wlgame.adjustTrust` 已接入；`cmd/wlgame`／state／rules
  測試在 `wolong-go:20260809` Docker 通過。
- 事件 2／3 超額外交：`sub_13902` `AL=3` 不收尾；`sub_13C3D` 再以 `AL=1Eh` 呼叫
  `sub_13DC9`，玩家 Trust `−30`。`World.ResolveDiplomacy` 與
  `TestDiplomacyAndFundingAmountOutcomeBounds` 已接入／驗證，狀態仍不收尾。
- 文件證據：`RESEARCH-LOG.md` 2026-08-10 兩節、`docs/mechanics/70-ai.md` §1.5.1、
  `docs/re/12-diplomacy-dialogue.md` §7。正式三平台包與推廣影片仍被事件 6／7 次要
  TALK、事件 10、災害物件動畫、目標平台 GUI 與同狀態對拍 gate 阻擋；不要提前建立。

## 2026-08-10 最新增量：PC-98 數字視窗量測與驗證 gate

- 讀取並套用 `research-pc98-golden-box-ui` 技能的版面原則：保留 640×400 邏輯畫布、
  16×16 CJK 格、內容／狀態／指令分區；不複製 Golden Box 的框線、肖像或文字。
- raw `sub_13902`／`sub_139E8` 在 `sub_17C6E` 前都設 `DX=0058h`／`BX=00B8h`；
  `sub_17D5F` 從 `(88,200)` 起畫 3×6、16×16 格；`sub_17D0D`→`sub_1FA37` 的
  `AX=4006h`／`sub_1FAA2` 已證實外框資源 blit 為 96×64。外框內容、每格語意、游標與
  TALK 逐頁契約仍未知，已記於 [`docs/re/13-pc98-numeric-window.md`](docs/re/13-pc98-numeric-window.md)。
- Docker/Xvfb 只做 PC-98 原版啟動短抽樣，確認 640×400 原生畫布；沒有進行長程遊戲或
  把啟動畫面冒充數字視窗 parity。`go vet ./...`、`go test -p=1 -vet=off ./... -count=1`
  在 `wolong-go:20260809` 通過。
- 目前 release gate 仍被事件 6／7 次要 TALK、事件 10、原版災害／投射物動畫、PC-98
  數字視窗與完整訊息排版、Windows／macOS GUI runtime、原版／remake 同狀態對拍阻擋；
  不建立正式包或推廣影片。

## 2026-08-10 最新增量：事件 6／7 次要 TALK 的負證據收斂

- `sub_13C3D`（IDA `00013C3D`）第一次 `sub_18810` 以 `DI=SP` 建 formatter stack；
  第二次在恢復 `DI` 後只做 `CX+0x1D`、`AL=0x93`，沒有重新建立 stack。事件 6／3
  直接落 #72，事件 7／2 直接落 #76。
- `sub_137D8`／`sub_13138` 的 `AH` 是 General `+0x1D`／`+0x1C` 雙向俘虜關係旗標，
  不是 Trust；不能用它填 `TalkNotice.City`。#72 的 `\\2` 會讀 `SS:[DI]` word，
  #76 是選單文字。
- 既有「#72/#73、#76/#77」範圍已勘誤為「#72/#76 直接證實；#73/#77 unknown」。
  PC-98 暫時 oracle 未越過啟動選擇畫面，沒有有效事件畫面證據。不要新增次要通知
  或猜測 state 欄位；正式包／影片 gate 維持關閉。

## 2026-08-10 最新增量：事件 9 #409 空槽

- `sub_150D7` 的第二個 formatter 呼叫是 `CX=0x199`／TALK #409；PC-98／DOS/V 原始
  `TALK.DAT` 均只有資料上的空行、沒有文字／marker。玩家可見結果只有 #37，remake 不建立空白 modal。
- `cmd/wlgame/messages_test.go` 的 `TestReleasedGeneralRawFollowup409IsEmptyNoOp`
  固定原始空槽與 UI no-op；不要把 #409 當成未知文案或再造泛用訊息。

## 2026-08-10 最新增量：M7 校訂表已接入 wlgame

- 版控檔 `translations/talk-dosv-corrected.json` 是 `talkdat.py correct` 的 1,022 則繁中
  產出，套用 60 筆已定案 `fix`；`talkdat_selftest.py` 會 byte-by-byte 比對重建結果，防止
  `extract`／`corrections` 與 runtime 產出失配。
- `internal/assets/text.LoadJSON` 只服務呈現層，保留 JSON 的行邊界、尾端空行與 `{N}` marker，
  以 cp950 轉回 raw Part；`Parse(TALK.DAT)` 的原始 byte round-trip 與原版唯讀規則不變。
- `cmd/wlgame` 預設載入校訂表，可用 `-talk-json ""` 退回 raw TALK；新 JSON 單測、完整
  Docker/Xvfb Go gate、TALK 兩版 verify、deny-list selftest、index gate 皆通過。
- 仍不可宣稱完整 remake：校訂後行寬／畫面抽樣、formatter／翻頁、事件 6／7 次要 TALK、
  事件 10、原版災害／投射物動畫、PC-98 數字視窗、Windows／macOS GUI 與同狀態對拍仍未完；
  因此不建立正式三平台包或推廣影片。

## 2026-08-10 最新增量：M7 modal 行寬 guard

- `tools/talkdat_selftest.py` 新增 22 格保守行寬檢查，校訂表目前最大 13 格；Docker selftest
  通過。這只保證 remake `drawMessage` 不會因校訂文字溢出 384 px modal，不可當成原版硬換行、
  formatter 或逐頁畫面證據。

## 2026-08-10 最新增量：數值面板／訊息分頁已接入

- `cmd/wlgame/amountpanel.go`：沿用 raw `sub_17C6E` caller 的已證實幾何 `(80,176)`、
  `(88,200)`、3×6／16×16；內容明確是 remake 的目前值／初值／上限呈現，未知的
  `CS:7D93h` 格位資料不可當原版語意。
- `diplomacy.go`／`funding.go` 已使用共用面板；`messages.go`／`main.go` 以 TALK 硬斷行
  分頁，每頁 5 行（對齊原版 `CX=510h`），動態寬度但不重切行。相關幾何／分頁測試與 Docker 無網路 Xvfb
  `go test ./...`、`go vet ./...`、`go build ./cmd/wlgame` 已通過。
- 這完成的是 remake UI 安全接縫，不解除事件 6／7 次要 formatter、事件 10 producer、
  原版災害／投射物動畫、Windows／macOS GUI runtime、原版同狀態對拍等 release gate；
  不要因此建立正式三平台包或推廣影片。

## 2026-08-10 最新增量：第 3／4／5 項定向接手

- 原版 PC-98 event3 fixture 已在一次性 DOSBox-X 容器中穩定進入前置外交通知與三選一；
  `sub_13B7E`（`00013B7E`）→ `sub_193E9`（`000193E9`）的通用輸入保存鏈已由 IDA
  Pro 9.4／IDA 線性位址補證。PC-98 `KI.EXE` SHA-256：
  `061917F9F3F5C03E29397A9C636D546052128A99B8C8CE31DED0E84CF2A481E8`。
- 原版／remake 同一 raw event3 fixture 截圖已保存：
  `docs/images/pc98-oracle-event3-choice.png` SHA-256
  `56A1A16BC6D92F75A5DB3DC49F3C961609F1EB6B908106A39136D4E3E32FDB5C`；
  `docs/images/wlgame-event3-choice.png` SHA-256
  `CA40B865B44A6EA13ED5B4F2C0B6AB913A0BC895EF48D7A19E1825501E535151`。
  這證明同狀態可進入原版式 composite；功能性肖像／場景／選項與高亮已接。DOS/V
  `KI.EXE` 硬體游標 mask 與 `ICONGRF` 3×6 靜態 glyph 已在後續增量解碼並接線；
  尚未完成的是自然 DOS/V／remake 整張截圖逐像素對拍。
- `internal/ui/textdraw.WrapLines` 現在在 marker 展開後，以 ASCII 8 px／非 ASCII 16 px
  和 22 full-width cell modal 寬度換行；保留 TALK hard line／空列，關閉標點不落列首，
  每頁 5 行。`internal/ui/textdraw` 與 `cmd/wlgame` 的排版／Xvfb 測試通過。
- `internal/rules/tactical.Soldier.PoseStep` → `ProjectileView.SpecialFrame` 保留原版
  `+0x02 bit 0`，`cmd/wlgame` 可選特殊投射物 raw `0x214`／`0x215`；事件 12 MCH 八相位
  矩陣仍使用原版 type 1／2，並新增固定 8-frame fallback phase 回歸測試。
- 本輪未跑完整劇情，符合使用者要求；尚未建立三平台正式包或推廣影片。仍開啟的 gate：
  事件 6／7 次要 formatter、事件 10、MCH／PC-98 物件 timer 與 UI parity、Windows／macOS
  GUI runtime、完整 packaged smoke。

## 2026-08-10 最新記憶：原版游標行為／數值格／TALK 分頁與消像

- 原版 `sub_17C6E`／`sub_17D5F` 證據已落成可執行表：`59 5A 5B 5D 5E 5E / 56 57 58 52 5F 5F /
  53 54 55 5C 60 60`。remake 的 `(88,200)` 起點、3×6／16×16、raw action 與 `0x60`
  完成格均由 `amountPanelButtonAtPoint` 命中；滑鼠與鍵盤走同一 `AmountEdit`。
- `messagePageRows=5` 對齊原版 `sub_18810`／`sub_1895D` 的 `CX=510h`。TALK marker 先展開，
  再以 ASCII 8 px／非 ASCII 16 px 實測換行；尾端結構空行移除，中間空行保留，Enter／Space
  逐五行翻頁。
- 事件 2／3 使用 IVENTGRF page 0、玩家君主肖像、General `+0x1E` 的 TALK variant 與
  三選一；事件 4／5 使用 page 1、官員肖像 fallback、prompt／選項與數值器。pending 清除
  後 `Draw` 的地圖路徑覆蓋舊場景，後續 TALK 才出現，完成功能性消像 parity。
- 新證據圖為 `docs/images/wlgame-event3-choice.png`（`CA40B865...535151`）與
  `docs/images/wlgame-event3-amount.png`（`27A5474E...04609D`）。完整長程遊戲測試略過；
  自然 DOS/V／remake 整張截圖逐像素、事件 6／7 次要 formatter、事件 10、MCH timer、
  目標平台 GUI 仍是開啟 gate，不能打包或做推廣影片。

## 2026-08-10 本輪 DOS/V 接手增量

- DOS/V 數值框資源不是 8×8 chrome tile：`sub_17D0D`（IDA `00017D0D`）的
  `word_10D50:0600h` 依 `sub_100DF` 段配置換算到 `ICONGRF` 第 3 段 `0x14A0`，
  以 `AX=4006h` 解為 96×64，目的 `(88,184)`；背景保存區仍 `(80,176)`／112×80，
  首個 raw 操作格 `(88,200)`。程式入口是 `gfx.DOSVAmountPanel`／`Library.DOSVAmountPanel`。
- 事件 6／7 次要 TALK：`sub_13C3D` 條件是 `AH != 0` 且 response 非 2／3；事件 6
  追加 #72 (`0x48`)，事件 7 追加 #76 (`0x4C`)。#76 `NoPortrait=true`；#72 的 `\\2`
  因第二次呼叫沒有重建 `DI=SP` stack，未知 payload 必須 fail-closed，不能把 `AH`
  猜成城市。測試已固定 sentinel `-1`，避免 Go 零值誤把未知城市變成城市 0。
- 事件 10 `sub_13496` consumer 保留 raw `Param` TALK index，事件高 byte 僅在有效
  General 範圍映射 `\\1`；producer 尚無證據。
- 災害 runtime 不再使用 global 8-frame clock：每筆 typed object 保存 active/type/city/
  timer/interval/phase/dirty，`AdvanceDisasterObjects` 對應 `sub_12459`，
  `RenderDisasterObjects` 對應 `sub_12533` 的「先畫舊 phase、再遞增」；state 單測驗證
  16-update cadence 與清除。只剩最後 slot 的 `sub_1248A` 移動分支。
- `HANDOFF.md` 已刪除，交接與待辦唯一寫在 `WORKLIST.md`；不要重新建立 HANDOFF。

## 2026-08-10 最新記憶：DOS/V 硬體游標／按鍵 glyph 已解碼

- `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
  IDA Pro 9.4 線性位址 `seg002:031B`（由 `sub_201E4` 兩次呼叫 `sub_2020C`）對應
  file `0x1051B`，兩層各 32 bytes；`0x0F` 是白框、`0x0A` 是紅填。
- `DecodeDOSVCursor` 按 `AH`／`AL` 寫入順序反轉每列 source byte、MSB-first 展開；
  16×16 結果計數固定白 39／紅 56／透明 161，palette-index SHA-256 為
  `385c2f1949d3d1e331399316305db7d7f2fd489a0626de1c6b2b8375aadfc6fe`。
- `ICONGRF` 第 3 段相對 `0x14A0` 的 96×64 資源下半部含 3×6 個 16×16 靜態 button
  glyph；`TestDOSVAmountPanelContainsStaticButtonGlyphs` 逐格驗證。`amountFrame` 存在
  時直接繪製原始資源，不再覆蓋 vector／CJK label；`Library.DOSVCursor` 提供 cursor，
  缺資源才 fallback。
- 全量 Go 測試、`go vet`、`docs/INDEX` generate/check、30 幀 Xvfb smoke 通過；完整長程
  遊戲測試略過。自然 DOS/V／remake 整張畫面对拍、事件 6／7 次要 formatter、事件 10、
  物件 timer／目標平台 GUI 仍是 gate；不建立三平台正式包或影片。

## 2026-08-11 接手記憶：事件 raw 接線與自然 GUI gate

- 事件 6／7 次要 TALK 已完成 raw formatter 邊界：事件 6 #72 只在回報勢力 record
  pointer 落於 DOS/V 原版 0x400-byte stack 範圍時保存 `RawFormatterWord=0` 且
  `RawFormatterWordValid=true`；事件 7 #76 沒有可知 formatter word。呈現層經
  `World.ResolveTalkFormatter2` 取 bytes，invalid 或越界必須整則 fail-closed；不可
  以 Go 零值推測城市 0。
- 事件 10 已增加 `World.QueueEvent10(general, talkIndex)`，沿完整 256 格 queue 寫入
  `(general<<8)|0x0A` 與原始 TALK index。這是 remake 的受控 fixture／劇本入口；IDA
  直接 caller 與現有 DOS/V queue 沒有找到原版自然 producer，不能把它寫成原版時序。
- 自然 DOS/V 在 Docker／DOSBox 只能到「密碼輸入：第 09 頁」，既有原版頁面與限制記在
  `docs/playtest/13-dosv-natural-and-target-gui.md`；remake 的固定種子 30 幀基準是
  `docs/images/wlgame-dosv-natural-remake.png`，不能冒充逐像素 parity。
- Linux／Xvfb GUI smoke、Windows amd64 `PE32+`、macOS amd64／arm64 `Mach-O` 交叉建置
  已驗證；無原生 Windows／macOS runtime，故正式三平台包與推廣片仍關閉。`HANDOFF.md`
  已捨棄，所有接手資訊只寫 `WORKLIST.md`、`CONTEXT.md`、`MEMORY.md` 與研究／驗證文件。

## 2026-08-11 YouTube oracle 與自然 HUD

- 使用者指定 `https://www.youtube.com/watch?v=af6xqcicXoI` 作原版自然畫面參考；metadata
  為 567 秒、478×360、30 fps。代表幀已放在 `docs/images/yt-wolong-natural-*.png`，
  80 秒 640×400 去黑邊參考幀 hash 為
  `c0217b8722bd44a22a112a2981b626126d5ee53d3e9f00498c6cbd018e08d6`。
- 影片確認舊的「自然策略畫面滿版地圖、沒有常駐側欄」假設不成立：自然 HUD 是
  32 px banner＋32 px command strip＋左側 27×21 地圖＋右側 208 px sidebar。新增
  `cmd/wlgame/strategyhud.go`，使用原版 minimap／portrait／chrome 與 state 數值。
- 最新 remake 自然截圖 `docs/images/wlgame-dosv-natural-remake.png` hash 為
  `961e583915d2e0e7b65cd51f637ec214530b68040ba0da5770add4b35cb46e30`。影片視覺／
  幾何對拍通過；影片壓縮、縮放且與 remake 日期／鏡頭不同，不宣稱嚴格同狀態逐像素。

## 2026-08-11 記憶：事件 10 是 idle clock 的下游通知

- DOS/V 已由 IDA `.i64` 證實無輸入鏈：`sub_11BE0` → `sub_11F7F`（滑鼠座標不變時
  設 `byte_198A3` bit 7）→ `sub_11CD0` → `sub_13EFD`／`sub_125A3`／`sub_12459`／
  `sub_11D8E`。所以玩家停手時日期、已下達路徑的軍團與 MCH runtime 物件仍會走。
- `sub_11D8E` 每小時進 `sub_13E11`，再呼叫 `sub_131AE`。事件 queue dispatcher
  由 `byte_131AD` 節流：`sub_12BD9` 初值 7，取到 record 後重設 0x0A；事件 10
  只有 queue low byte `0x0A` 時才由 `sub_13496` 轉成 TALK。不要把這個 `0x0A`
  和節流 byte 的 `0x0A` 混在一起。
- remake 對應為 `game.timeRuns` → `World.Tick`；`World.QueueEvent10` 只可作受控
  raw fixture／劇本注入，不能驅動時鐘。驗收為
  `TestIdleClockDispatchesQueuedEvent10OnHourlyCadence`；詳細證據在
  `docs/re/16-idle-clock-event10.md`。
- 原版 idle path 順序是據點／軍團／物件／時鐘；remake 正常 UI map-loop 已改用
  `World.TickMap` 對齊，額外 `g.speed` tick 使用不含物件的 `World.Tick`，維持 MCH
  動畫 cadence。原版自然 event10 producer 仍 unknown，限時封口不再無限追查。

## 2026-08-11 最新記憶：四項依序封口

- M7 60 筆校訂已完成逐筆人工語意／marker／硬換行／寬度／代表畫面抽樣；入口為
  `docs/playtest/14-m7-review.md`，可重跑 `tools/m7_review.py --check`。
- 事件 2–5 已由 `TestEvent2To5FullTalkPageSampling` 封口 36 個 raw TALK 頁面／18 組
  雙頁回應，入口為 `docs/playtest/15-event2-5-talk-sampling.md`。
- 事件 9 已由 27 小時 `World.Tick` queue fixture 與 GUI 通知測試封口；玩家勢力才有 #37，
  非玩家／在野與 #409 空白不產生錯誤通知，入口為 `docs/playtest/16-event9-long-route.md`。
- 推廣片已生成於 `dist/promo/wolong-remake-trailer.mp4`，60 秒 1280×720 H.264/AAC；
  配樂是 `tools/promo_score.py` 新作，詳見 `docs/promo/README.md`。正式三平台包與
  Windows／macOS 原生 GUI smoke 尚未宣稱完成。

## 2026-08-11 最新記憶：事件 10 改用可關閉近似 producer

- 原版 `0x0A` producer 仍 unknown；不要把近似規則寫成原版逆向結論。
- `LoadScenario` 預設開啟 `World.produceApproximateEvent10`。月結後從玩家勢力的活俘虜
  每月最多選一名，依 `General.Timer` 與 `rand&0xFF` 的 `0x20／0x40` 邊界，近似排入
  raw `(general<<8)|0x0A`、`Param=0x41`（逃走）或 `0x42`（歸降）。
- `SetApproximateEvent10(false)` 可關閉 substitute；事件仍由下一次每時 consumer 顯示。
  驗收為 `TestApproximateEvent10ProducerUsesKnownRawContract`、
  `TestApproximateEvent10ProducerIsBoundedAndDisableable`、
  `TestApproximateEvent10ReentersIdleClockConsumer`。

## 2026-08-11 記憶：DOSBox／remake 可玩性驗證

- DOS/V DOSBox 原版可啟動，但本輪停在複製保護第 15 頁；PC-98 DOSBox-X 的既有 oracle
  提供 `NEW GAME`、劇本選單與戰略地圖畫面。密碼保護仍是 DOS/V 自然對拍阻擋。
- 目前 remake 建置以固定 seed 17、無 `-open-*` 走過編成、行軍、時鐘自然流逝與
  196/6/28 遭遇；實際系統視窗第 1 槽 save／load overlay 88,832 bytes 通過。
- 戰術 GUI／2 號命令 current-build debug smoke 通過；正常無旗標戰術與原版同狀態 parity
  仍依 `docs/playtest/09-wlgame-normal-tactical-path.md` 與原有矩陣，不把 smoke 升格。
- 測試報告與截圖 hash 在 `docs/playtest/17-expert-dosbox-remake.md`；PC-98 容器目前
  缺少 window manager，完整 bus-mouse 重播是下一個最小輸入橋接工作。

## 2026-08-11 記憶：YouTube／推廣片像素差異

- 使用者要求以 YouTube 遊玩影片與 remake 推廣片比較像素差異；已生成
  `dist/promo/wolong-remake-yt-comparison.mp4`、`docs/promo/yt-remake-natural-side-by-side.png`
  與 `docs/promo/yt-remake-natural-difference.png`。
- 640×400 自然幀 raw metric 為 `AE=255003/256000 (99.61%)`、`RMSE=0.338208`。
  這是不同日期／鏡頭／狀態的可見差異量測，不是同狀態 renderer 錯誤率；完整判讀見
  `docs/promo/yt-remake-pixel-review.md`。
- 目前決策：YouTube／推廣片視覺 oracle 已足夠，不再以 DOS/V 密碼保護後的嚴格逐像素
  diff 阻擋 remake；正式三平台包與 Windows／macOS 原生 GUI 仍是獨立 release gate。

## 2026-08-11 記憶：DOS/V 自然 HUD 骨架已對齊

- 右欄不是兩個獨立框直接背靠背：minimap 色標列與下方情報框共用一條 8 px 分隔邊；
  `strategyFactionY=176`、君主頭像內容從 y=184 開始。
- 常駐資訊區已改成原版順序：君主／首都／軍師 → 信賴度 → 紅色分隔線 → 黑底資金／
  預備兵與右對齊數值。reserve 中央 raw glyph 仍是未知，只保留幾何 fallback。
- 重新產生 `dist/promo/wolong-remake-trailer.mp4`、`dist/promo/wolong-remake-yt-comparison.mp4`
  與差異圖；最新基準幀為 `docs/images/wlgame-dosv-natural-remake-skeleton.png`。

## 2026-08-12 勘誤記憶：反組譯再審

- 事件 6 #72 的 `RawFormatterWord=0` 已被撤回：第二次 formatter 的 `\\2` 會讀未保存的
  `SS:[DI]` transient word，現行 remake 必須 `-1/false` fail-closed；詳見
  `docs/re/12-diplomacy-dialogue.md`。
- 所有目前可見 event 10 direct queue producer 都未寫 `0x0A`；這是有限負證據，原版
  natural producer 仍 **unknown**，不能把 substitute 當逆向結論。
- `sub_1248A` 是 slots 16–31，type 3 仍未知；投射物已補 raw 8／6 冷卻與
  `0x214／0x215` special frame。音源 TSR 的已證實／未知邊界見
  `docs/re/17-dosv-audio-tsr.md`。

## 2026-08-12 勘誤記憶：DOS/V 密碼頁不再阻擋

- 在 `wolong-dosboxx:latest`（DOSBox-X 2025.02.01）以
  `mouse_emulation=integration`、INT 33 範圍 640×480 重新輸入後，空白確認、`0000`、
  `1234` 都由松崗 DOS/V 密碼頁進入原版開場。
- `1234` 與空白控制組的 10 秒後開場截圖 SHA-256 同為
  `10cd0e199e7bd944a4664f3a3e4debd94a3986df1c36ca7eac4ffaf208ebbd34`；四格沒有回顯數字，
  故結論是「確認流程未阻擋」，不是「已解出密碼比較語意」。
- 密碼頁不重製到 remake，也不再列為 DOS/V oracle／自然畫面採樣阻擋。完整長程驗證與
  同狀態逐像素 parity 仍未做，不能混為已完成；證據入口是
  `docs/playtest/18-dosv-password-verification.md`。
