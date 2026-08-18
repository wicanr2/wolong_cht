# 臥龍傳 remake 工作清單

> ⚠ **待辦看 [`CONTEXT.md`](CONTEXT.md) §7，不是這裡。**
> 本檔是**按日期的完成紀錄**：每一節記下那一天封口了什麼、當時的邊界在哪。
> 每一輪都會更新的「現在該做什麼」在 `CONTEXT.md` §7.0
> （`CLAUDE.md` §10 也指向那裡）。
>
> 兩邊都寫「唯一來源」會讓接手的人拿到舊快照——本檔最後一次補紀錄是
> 2026-08-17，而 `CONTEXT.md` §7 每一輪都動。
>
> 自 2026-08-10 起不再建立或更新 `HANDOFF.md`；歷史交接內容已濃縮到本檔，
> 深層證據請回查 `RESEARCH-LOG.md`、`CONTEXT.md` 與 `docs/`。
>
> ⚠ **日期節裡的敘述是「當時的認知」**，其中有些後來被推翻
> （例如密碼頁一度被當成 oracle 的阻擋，2026-08-12 測出**不擋**）。
> 推翻紀錄集中在 `CONTEXT.md` §6。

## 目前目標

完成《臥龍傳－三國制霸之計》的證據導向 remake：核心規則、松崗 DOS/V
呈現、繁中 TALK、存檔與正常玩家路徑接到可重播的跨平台程式。

本輪以松崗繁中版為唯一畫面／行為 oracle（`workplace/orig/dosv` 是沿用的資料夾名稱）：

- 640×400 畫面外框與數值視窗以 DOS/V 對準。
- PC-98 與其他版本只保留歷史研究，不作本輪視覺、行為或 release 驗收基準。
- 依使用者要求不跑完整長程遊戲測試；採窄 fixture、單測、Docker/Xvfb 短 smoke。
- 三平台候選包可在短驗收 gate 通過後建立；原生 Windows／macOS runtime 仍是獨立 gate，不能由交叉編譯檔頭代替。

## 已完成且可回查

- 格式／時鐘／規則／狀態／存檔 round-trip 的主要垂直切片已接入；原始
  `workplace/orig/`、`workplace/ida/` 視為唯讀。
- `wlgame` 已接真實劇本、地圖、命令、軍團、行軍、戰術遭遇、四槽存檔 overlay、
  戰鬥指揮／委任與固定種子重播。
- 事件 1、2、3、4、5、8、9、11、12、13 的已證實狀態接縫已在 `internal/state`
  接入；事件 2／3／4／5 的玩家選單與數值編輯共用同一條 modal 狀態路徑。
- `sub_17C6E` 數值核心、DOS/V `CS:7D93h` 3×6 raw 格位、16×16 命中區與
  `AmountEdit` 已接入；DOS/V `sub_17D0D` 的 ICONGRF 第 3 段 96×64 內框也已
  接到 `(88,184)`，外圍保存區仍是 `(80,176)`／112×80。`KI.EXE` `seg002:031B`
  的兩層 16×16 硬體游標 mask 與 `ICONGRF` 下半部 3×6 靜態按鍵 glyph 已解碼、
  載入並接到面板；只保留資產缺失時的 fallback。
- TALK marker 展開、原始行硬斷行、五行／16 px 分頁、肖像與 IVENT 場景、
  pending 完成後的消像已接入；M7 校訂輸出由
  `translations/talk-dosv-corrected.json` 提供。
- `MMAP.MCH` 16×16 平面圖塊、`0xA000` metadata、`CS:985Ah` type 1／2／3
  八相位查表與火災／暴動圖像接到戰略地圖；runtime object record 已依
  `sub_123FF`／`sub_12459`／`sub_12533` 接入 state，包含 16-update timer、dirty
  與 render 後相位遞增。後半區 slots 16–31 的移動分支 `sub_1248A` 已依 raw fixed-point、
  方向 byte 與邊界 wrap 接入並有定向測試。
- Docker 內已有 Go test／vet、翻譯 selftest、文件索引與 Linux/Xvfb 截圖 smoke；
  完整長程遊戲測試依使用者要求略過。

### 2026-08-12 `dist-all` 統一交付封口

- [x] 新增可重播的 `tools/release_all.sh`／`tools/release_all_fs.py`：所有建置、封裝、
  雜湊與 AppImage 步驟都在 Docker 內執行，最終只保留 `dist-all/` 的交付檔；原始
  資料、完整 TALK 表與 `.work` 建置中間檔均被排除。
- [x] `dist-all/packages/` 已集中 Linux amd64、Windows amd64、macOS Intel＋Apple Silicon
  三個桌面完整包及 Linux amd64 AppImage；Linux arm64 無頭工具包另列為伴隨包，不假稱完整
  GUI 平台。
- [x] 四支推廣片已複製至 `dist-all/promo/`。其中
  `wolong-remake-dosv-live-comparison.mp4` 是 60 秒、1280×720 的「讓經典再現」
  DOS/V／remake 實機比較片；YouTube 代表幀短片是 24 秒、1280×400，兩者不可混稱為
  同狀態逐像素 parity。
- [x] Linux 解壓包與 `APPIMAGE_EXTRACT_AND_RUN=1` AppImage 都在 Docker/Xvfb 載入玩家
  自備松崗資料與倚天字型，從同包公開 `corrections.json` 成功啟動並得到同一張 640×400
  固定種子截圖（SHA-256 `45a68852335420dd7b22b4e240192dcd7a38fbbc62f72c8c59ec95acdc137b24`）。
- [x] 公開校訂覆蓋在真實松崗 `TALK.DAT` 上逐則與既有 1,022 則產生表位元組一致；完整表
  僅作測試 oracle，不進發行包。`denylist.py dist-all` 通過，根目錄 `SHA256SUMS.txt` 已驗證。
- [x] Android APK 收在 `dist-all/experimental/android/`，並清楚標為觸控 shell 原型，不計入
  三平台完整遊戲發行。
- [ ] Windows／macOS 尚缺目標作業系統的原生 GUI／輸入／音訊／字型短 smoke；交叉 ABI
  檔頭已驗證，不能取代該 gate。

## 本輪高優先工作

### P0：DOS/V 畫面外框與數值選取（內框已完成）

1. 已由 DOS/V `sub_17D0D` 接入 `ICONGRF` 第 3 段相對 `0x14A0` 的 96×64
   平面內框，目的座標為 `(88,184)`；`sub_19796` 保存區與 `(88,200)` 格位仍分離。
2. 實際數值選取已使用 `CS:7D93h` raw byte 與同一套 `AmountEdit`；DOS/V
   `KI.EXE` cursor mask（file `0x1051B`）與 `ICONGRF` 3×6 靜態 glyph 已由原始
   bytes 解碼並接線，`internal/assets/gfx` 單測固定像素統計與每格 glyph。
3. DOS/V 執行檔仍有複製保護，但使用者提供的 YouTube 錄製
   [`af6xqcicXoI`](https://www.youtube.com/watch?v=af6xqcicXoI) 已提供自然遊戲 oracle；
   原版 DOSBox 密碼頁只保留作啟動邊界，不再阻擋自然畫面結構對照。

4. `cmd/wlgame` 有真實 96×64 資源時不再疊畫 vector 矩形／CJK 按鈕標籤；解出的
   16×16 白框／紅填 cursor 只在數值面板選取狀態繪製。自然策略 HUD 已依影片／說明書
   對齊 32 px 命令列、左側 27×21 地圖與右側 208 px minimap／情報欄。

### P0：事件 6／7 次要 TALK（raw formatter 已接線）

1. 已追完 `sub_13327`、`sub_13388` → `sub_13C3D`：事件 6 是 `CX=0x48`／#72，
   事件 7 是 `CX=0x4C`／#76；條件是 `AH != 0` 且 response 不是 2／3。
2. `TalkNotice` 已保留 `Secondary`／`NoPortrait`；事件 6／7 state 單測固定雙向
   俘虜旗標、索引與順序。#76 以無肖像文字 modal 顯示。
3. 第二次呼叫沒有重建 `DI=SP`，#72 的 `\\2` 參數不能安全映射成城市；呈現層缺
   formatter 時整則 fail-closed，禁止把 `AH` 猜成城市、勢力或信賴度。最新 raw handler
   證據顯示第二次會讀 `SS:[DI]`，而此 transient stack word 不能從 World state 重建；
   因此事件 6 現在固定 `RawFormatterWord=-1`、`RawFormatterWordValid=false`，不再把
   Go 零值 `0` 偽裝成原版 payload。只有未來動態原版 trace 擷取到該 word 時才能標為有效。

### P0：事件 10（raw producer／consumer 已接）

1. `sub_131AE` dispatch 的低碼 `0x0A` 已確定進 `sub_13496`；該 handler 將
   `AL=事件字高 byte`、`AH=0xFF`、`CX=Param` 後呼叫 `sub_18810`。state consumer
   已保留 raw `Param` 為 TALK index，並在高 byte 是有效 General 時提供 `\\1`。
2. 直接檢查 `sub_12FBF`、`sub_12FB1`、`sub_1301C` 的 caller，尚未找到可證實的
   `0x0A` producer；這是負證據，不是全 binary 排除。
3. `World.QueueEvent10` 保留已證實的受控 raw producer；另新增
   `World.produceApproximateEvent10`，在月結後依玩家俘虜狀態近似排入 TALK `0x41／0x42`。
   近似規則預設開啟但可由 `SetApproximateEvent10(false)` 關閉，不冒充原版自然時序。
4. consumer 與兩種 producer 都沿 raw queue／formatter marker 接入；未知原版 TALK index
   不以泛用文案替代，近似路徑只使用已查到的俘虜文字並標為 substitute。
5. [x] `cmd/wlgame` 的 `idleClockGate` 已對應松崗繁中 `sub_11F7F` 的穩定游標條件：
   首次座標觀測、游標位移、按鈕或命令 frame 都停住世界；游標穩定且無輸入時才按
   據點／軍團／物件／時鐘順序跑 `World.TickMap`。`TestIdleClockGateRequiresStablePointerAndNoCommand`
   與 `TestIdleClockDispatchesQueuedEvent10OnHourlyCadence` 已納入 parity gate。

### P0：物件動畫時序（timer／render 已完成）

1. 已依 DOS/V IDA 線性位址實作 runtime object record：
   - `sub_123FF` `000123FF` 建立物件時寫 `[si+0C]=1`、`[si+0D]=0x10`、
     `[si+0F]=1`，並把 type 寫入 `[si+0E]`。
   - `sub_12459` `00012459` 每次主畫面更新對 active object 遞減 `[si+0C]`；
     到零時補回 `[si+0D]` 並設 dirty bit 0。
   - `sub_12533` `00012533` 只在 dirty bit 命中時把 `[si+0F]` 加一並 `& 7`，
     再依 `CS:985Ah` 查表取 MCH frame。
   - `sub_12438` `00012438` 依同一城市座標清除 active object。
2. 原本的 `disasterFrameTicks=8` fallback 已移除，改成 typed、可測試的 per-object
   `Phase`／`Timer`／`Delay` 狀態；`World.DisasterMarker` 仍不序列化。
3. `TestDisasterObjectAnimationTiming` 已驗證「建立初值 → 首次 timer cadence →
   16 次 update → dirty → render 舊 phase → phase +1 → 清除」；`wlgame` 只在可見
   map-loop Update 推進，模態視窗不穿透。
4. MCH type 1／2 的圖像接線與固定時序已完成；type 3 查表保留但事件語意未知，
   暴風雨仍依 `sub_134A6`／`sub_1237E` 的已證實範圍處理。後半區 slots 16–31 的物件移動
   `sub_1248A` 分支已依 raw fixed-point word、方向 byte、邊界與 wrap 規則接入，
   由 `TestSub124FFMatchesRawSignedByteContract`、兩個 `MovingDisasterSub1248A*`
   測試固定。

### P1：DOS/V 自然畫面对拍與目標平台 GUI

1. 使用者影片 `af6xqcicXoI`（567 秒、478×360、30 fps）已取 20／80／160／240／
   320／400／480／550 秒代表幀；80 秒幀另去黑邊／縮放成 640×400。比較確認橫幅 32 px、
   命令列 32 px、左側 27×21 地圖與右側 208 px minimap／情報欄。
2. `cmd/wlgame/strategyhud.go` 已接入這個自然 HUD；remake 固定 `seed=17`、30 幀
   輸出 [`wlgame-dosv-natural-remake.png`](docs/images/wlgame-dosv-natural-remake.png)，
   並與 [`yt-wolong-natural-80s-640x400.png`](docs/images/yt-wolong-natural-80s-640x400.png)
   做結構／色彩／欄位位置對照。
3. 影片是有損縮放來源，且 80 秒為 196 年 4 月 5 日、remake smoke 為 196 年 4 月 1 日；
   因此「影片視覺 oracle 對拍」已通過，但嚴格同狀態逐像素 diff 不宣稱。
4. Linux／Xvfb GUI smoke 通過；Windows amd64 交叉建置產生 `PE32+ x86-64`，macOS
   amd64／arm64 產生 `Mach-O`。目前 Docker 沒有 Windows／macOS 原生桌面執行環境，
   這兩項仍是 GUI 編譯／格式 gate，不寫成原生 runtime parity。

## 2026-08-10 最新增量：DOS/V 硬體游標／button glyph

- `KI.EXE` `sub_201E4`（IDA `seg002:01E4`）設定 `SI=031Bh`，兩次呼叫
  `sub_2020C`；`seg002:031B` 的 64 bytes 對應 file `0x1051B`，第一層
  `AX=0F00h` 白色、第二層 `AX=0A00h` 紅色，解成兩層 16×16 mask。
- `internal/assets/gfx.DecodeDOSVCursor` 已固定原始每列 `AH`／`AL` 反轉與 MSB-first
  展開；單測固定白 39、紅 56、透明 161 與完整 16 列形狀。
- `ICONGRF` 第 3 段相對 `0x14A0` 的 96×64 資源下半部是 3×6 個 16×16 靜態
  button glyph；單測逐格確認非背景像素。`amountFrame` 存在時直接畫資源，不再覆蓋
  vector 矩形／CJK button label；`Library.DOSVCursor` 提供原版 cursor overlay，
  缺資源才 fallback。
- Docker 全量 `go test -p=1 -vet=off ./... -count=1`、`go vet ./...`、文件
  `index.py generate/check` 與 30 幀 Xvfb 短 smoke 已通過；完整長程遊戲仍依要求略過。

## 2026-08-11 本輪封口

本輪以短 fixture、可重跑單測與 IDA 靜態證據封口，不以完整長程遊戲測試作為必要條件：

- [x] `sub_1248A`：只對 raw slot 16–31 執行移動；兩次 `sub_124FF` 的有號 byte
  漂移、方向 byte、邊界與 `-0x10..0x190`／`-0x10..0x110` wrap 已實作並測試。
- [x] 事件 2／3／4／5：外交與資助的成功、拒絕、金額邊界及 TALK 展開／硬換行／5 行
  分頁由 `TestEvent2To5TalkBranchParityGate` 表格化驗收。
- [x] 事件 9：玩家勢力／非玩家勢力／在野與空白 #409 的短 fixture 由
  `TestEvent9ShortFixtureGate` 封口；不宣稱原版長程時序。
- [x] M7：`corrections.json` 的 60 筆校訂逐筆展開，逐筆檢查寬度與最多 5 行分頁，並由
  `TestM7CorrectedTalkLayoutGate` 加 `tools/talkdat_selftest.py` 驗收。
- [x] 投射物：一般水平／垂直與特殊 `0x214`／`0x215` frame、發射姿態位元及運動規則由
  `TestProjectileParityGate` 與 tactical projectile tests 驗收。
- [x] 事件 10：已完成 IDA Pro 9.4 `.i64` 的 queue dispatcher、consumer、writer、caller
  與 data-ref 深度追查。原版自然 producer 仍是未知；這是有界的負證據結論，
  `World.QueueEvent10` 保留為受控 raw producer，已不再把未知來源列為無限追查阻塞。

### 2026-08-11 事件 10 與無輸入自動 clock 勘誤

- [x] 已證實 DOS/V 無輸入路徑：`sub_11BE0` → `sub_11F7F`（座標不變設
  `byte_198A3` bit 7）→ `sub_11CD0` → `sub_13EFD`／`sub_125A3`／`sub_12459`／
  `sub_11D8E`；因此日期、已下達目的地的軍團、據點 runtime 與 MCH 物件會自動前進。
- [x] 已證實事件 10 只是這條路徑的下游 queue consumer：每小時進入 `sub_13E11`，
  再由 `sub_131AE` 依 `byte_131AD` 節流；初始化為 7，取到一筆後重設 0x0A，
  不是每小時同步取，也不是 clock／行軍 driver。
- [x] 新增 `TestIdleClockDispatchesQueuedEvent10OnHourlyCadence`，固定 7 個每時邊界
  驗證預先注入的 `Code=0x030A`／`Param=0x42` 只在節拍邊界產生 TALK；完整證據見
  [`docs/re/16-idle-clock-event10.md`](docs/re/16-idle-clock-event10.md)。
- [x] 原版自然 `0x0A` producer 仍未知，維持限時封口；不要把 `World.QueueEvent10`
  接成 clock。正常 UI map-loop 已用 `World.TickMap` 對齊原版據點／軍團／物件／時鐘；
  同一畫面的額外規則 tick 使用不含物件的 `World.Tick`，不污染 MCH 動畫 cadence。

### 2026-08-12 追加 IDA 再審：未知項的已證實邊界

- [x] 事件 10：逐一展開所有目前可見 direct queue producer；已證實集合僅寫出
  `0x01`–`0x09`、`0x0B`–`0x0D` 的已列舉路徑，沒有 `0x0A`。這強化有限負證據，
  仍不排除間接／外部／受保護流程，原版 natural producer 保持 **unknown**。
- [x] 事件 6／7：`\\2` handler `000108DB` 已證實第二次 formatter 讀 `SS:[DI]`；
  舊有「raw word 0」已撤回，#72 以無 payload fail-closed，#76 保持無肖像次要文字。
- [x] MCH：`sub_123FF` 唯一 direct caller 是 event 12 handler；`sub_12286` 只直接排
  `0x010C`／`0x020C`。type 3 仍未知；`sub_1248A` 是 slots 16–31，不是最後一筆。
- [x] 戰術：投射物已接原版 normal／special branch、`0x214／0x215` source frame 與
  raw `[si+13]` 的 8／6 發射冷卻；完整同狀態逐像素生命週期仍非已證實。
- [x] 音源：新增 `docs/re/17-dosv-audio-tsr.md`，記錄 `YNSOUND.COM` 的 INT 61h command
  table、遊戲效果碼與 address/data 硬體寫入；精確音色格式與硬體型號仍未知。

### 2026-08-11 事件 10 近似自然 producer

- [x] 月結後新增 `produceApproximateEvent10`：只選玩家勢力目前收容的活武將，每月最多
  一筆，保留 `General.Timer` 倒數閘與固定 RNG 邊界；逃走／歸降分別排入 raw TALK
  `0x41`／`0x42`。
- [x] raw payload 維持已證實的 `(general<<8)|0x0A`／`Param`，下一個每時 queue 節拍
  才由 `sub_13496` 對應的 remake consumer 產生 `TalkNotice`；不把 producer 接成 clock。
- [x] `SetApproximateEvent10(false)` 可關閉替代規則，`TestApproximateEvent10*` 固定
  payload、逃走／歸降狀態、queue 滿／倒數邊界與 idle clock consumer 接縫。
- [x] 原版自然 producer 仍維持 **unknown**；本項完成的是可遊玩的 substitute，不提升
  原版 parity 證據等級。完整說明見 [`docs/re/15-event10-producer.md`](docs/re/15-event10-producer.md)。

### 2026-08-11 M7 人工文字與畫面抽樣封口

- [x] `tools/m7_review.py --check` 逐筆確認 60 筆校訂、TALK marker、校訂表產出一致性、
  原始硬行與保守字寬；`TestM7CorrectedTalkLayoutGate` 再以 runtime `textdraw` 實測像素寬度。
- [x] 逐筆閱讀 `translations/corrections.json` 的語意備註，確認人名／勢力／據點／金額
  marker 沒有被文字修正誤換；#321／#322／#718 群組／#751 等槽位差異保留原始證據。
- [x] 以 `#321`、`#258`、`#663`、`#718`、`#889`、`#967` 代表幀回看硬換行、五列分頁、
  寬度、標點、尾端空行與戰場命令；完整表與畫面連結見
  [`docs/playtest/14-m7-review.md`](docs/playtest/14-m7-review.md)。
- [x] 本項完成的是 60 筆已定案校訂的人工審查，不把它擴大宣稱成 1,022 則全部重譯或
  同狀態逐像素 parity。

### 2026-08-11 事件 2–5 完整 TALK 抽樣封口

- [x] `TestEvent2To5FullTalkPageSampling` 以真實校訂 TALK、原始索引、marker 展開、硬換行、
  runtime 字寬與五列上限逐頁檢查事件 2／3／4／5 的 36 個 raw TALK 頁面、18 組雙頁回應。
- [x] 已覆蓋外交的自由／有資金／拒絕／超額，以及撥款的全額／等額／低額／零額／超額／
  拒絕分支；各分支的 raw index 與下一頁 index 均由 fixture 固定，不用泛用文案代替。
- [x] 四張代表幀與完整索引、限制條件見
  [`docs/playtest/15-event2-5-talk-sampling.md`](docs/playtest/15-event2-5-talk-sampling.md)。
- [x] 本項封口的是完整分支／逐頁抽樣與版面 contract；不把受控 fixture 宣稱成完整自然劇本
  長程或原版同狀態逐像素 parity。

### 2026-08-11 事件 9 長程通知流程封口

- [x] `TestEvent9LongNaturalRoute` 以正常 `World.Tick` 跑 27 小時、每小時 9 個 subtick，
  驗證 queue delay=7 的三筆事件在第 7／17／27 小時依序取出，並確認據點／軍團／物件／時鐘
  的 idle clock 順序沒有被事件 9 取用污染。
- [x] `TestEvent9LongNotificationRoute` 逐段驗證玩家勢力釋放武將才產生 #37 modal；非玩家勢力、
  在野武將與空白 #409 均不產生錯誤通知，玩家後續釋放仍可再次通知。
- [x] 完整 fixture、raw index、通知條件與代表幀見
  [`docs/playtest/16-event9-long-route.md`](docs/playtest/16-event9-long-route.md)。
- [x] 本項封口的是可重跑的長程 queue／通知流程；完整原版自然劇本仍依使用者要求不跑，
  當時也沒有做同狀態逐像素對拍。

### 2026-08-11 推廣影片產出

- [x] 已用 remake 實機流程代表幀串接自然策略、事件 2–5 TALK、戰術、投射物、戰果、事件 9、
  M7 與存檔畫面，產出 60 秒、1280×720、H.264/AAC 影片。
- [x] 配樂由 `tools/promo_score.py` 以原創合成音色產生；未使用原版 `SOUND.DAT` 或原版 BGM，
  權利與重現命令見 [`docs/promo/README.md`](docs/promo/README.md)。
- [x] 影片輸出：`dist/promo/wolong-remake-trailer.mp4`；正式三平台可執行包與目標平台原生
  GUI smoke 仍是獨立 release gate。

### 2026-08-12 DOS/V／remake 實機動態推廣片

- [x] 依使用者指定錄製／剪輯 60 秒、1280×720、H.264/AAC 的
  [`wolong-remake-dosv-live-comparison.mp4`](dist-all/promo/wolong-remake-dosv-live-comparison.mp4)。
  原版側是使用者指定松崗 DOS/V 遊玩錄影與受控 DOSBox-X 新遊戲畫面；remake 側是正常
  鍵盤路徑的策略、編成、目的地與行軍實機擷取。
- [x] 影片明示同類畫面比較「非同狀態逐像素判定」；戰術 remake 段標為獨立 fixture，
  不被挪作自然路徑驗收證據。原版音訊已排除，唯一音軌為本專案原創合成配樂。
- [x] 完整來源、秒數、雜湊、抽樣驗收與離線重播規則見
  [`docs/promo/dosv-live-comparison.md`](docs/promo/dosv-live-comparison.md)；暫存原版錄影在驗收後
  刪除，不隨平台遊戲包或工作區長期保存。

### 2026-08-11 DOSBox／remake 可玩性專家驗證

- [x] DOS/V 原版以固定 `cycles=20000` 在 DOSBox 啟動至密碼保護頁；不繞過密碼，
  原版自然玩法保持阻擋界線。PC-98 DOSBox-X `NEW GAME` 與既有自然 oracle 截圖可用。
- [x] 目前 remake 建置以 `-seed 17`、無 `-open-*` 旗標走過編成、關閉通知、行軍、
  玩家不下命令時的日期流逝與 196/6/28 遭遇選單；逐階段證據見
  [`docs/playtest/17-expert-dosbox-remake.md`](docs/playtest/17-expert-dosbox-remake.md)。
- [x] 正常系統視窗儲存／讀取第 1 槽通過，overlay 寫入 88,832 bytes；原始 DOS/V 檔案
  唯讀掛載，沒有被覆寫。
- [x] 目前建置的戰術 GUI／`2` 號攻擊完成 debug smoke；正常無旗標戰術接續沿用
  [`docs/playtest/09-wlgame-normal-tactical-path.md`](docs/playtest/09-wlgame-normal-tactical-path.md)，
  不把 debug fixture 宣稱成原版 parity。
- PC-98 DOSBox-X 的 bus-mouse／焦點輸入重播仍不穩定，但依使用者決策，PC-98 不再是
  本輪 DOS/V remake 的畫面或 release gate；保留作未來研究，不列為收尾阻塞。

### 2026-08-11 DOS/V 自然策略骨架對齊

- [x] 以 YouTube 80 秒 DOS/V 640×400 參考幀重對常駐骨架：32 px banner、32 px 命令列、
  左側 432×336 地圖與右側 208 px 欄位維持同一座標契約。
- [x] 修正右欄 minimap／情報框的共用分隔邊：原版 16 px 紅／藍勢力色標列覆蓋共用
  8 px 邊，不再產生 remake 的雙倍分隔；情報區改為君主／首都／軍師、信賴度、黑底
  資金／預備兵結構。最新畫面為
  [`docs/images/wlgame-dosv-natural-remake-skeleton.png`](docs/images/wlgame-dosv-natural-remake-skeleton.png)。
- [x] 依新骨架重錄 60 秒推廣片與 24 秒 YouTube 對照片；raw `AE` 由 `255003` 降為
  `249178`，`RMSE` 由 `0.338208` 降為 `0.329145`。差異解讀仍不升格為同狀態逐像素 parity。

### 本輪 gate 入口

在 `wolong-go` Docker 容器內執行 [`tools/parity_gate.sh`](tools/parity_gate.sh)。
事件 10 的非破壞性 IDA 匯出器是 [`tools/ida_event10_producer.idc`](tools/ida_event10_producer.idc)，
完整證據與推論等級見 `RESEARCH-LOG.md` 的 2026-08-11 章節。

### 2026-08-11 三平台候選封裝與 Android 規劃（歷史起點）

- [x] Linux amd64、Windows amd64、macOS Intel／Apple Silicon 的候選產物已在 Docker 建置；Linux 原生 `wlgame`／`wlview`、Windows PE32+、macOS `Mach-O` 均已核對檔頭。
- [x] `dist/release-20260811/packages/` 已產出 Linux amd64、Windows amd64 與 macOS universal 三個主 `.tar.gz`，另附 Linux arm64 邏輯工具包；解包清單、每包 `SHA256SUMS.txt`、外層雜湊與 deny-list 均通過。
- [x] 以封裝 Linux `wlgame` 執行 640×400、固定 seed、30 幀 Xvfb smoke；輸出 hash 為 `45a68852335420dd7b22b4e240192dcd7a38fbbc62f72c8c59ec95acdc137b24`。
- [x] Android 規劃已建立：[`docs/mobile/android-plan.md`](docs/mobile/android-plan.md)。第一版鎖定橫向、保留 640×400 邏輯畫布、用安全區／觸控抽屜／手勢轉接改善手機操作。
- 當時尚未開始 Android 觸控層、固定 SDK／NDK／`gomobile` Docker 工具鏈與 APK smoke；
  這是歷史起點，現況以下方「Android 原型與 AppImage 封口」為準，仍不宣稱 Android 完整支援。

### 2026-08-11 Android 原型、AppImage 與「經典再現」推廣片封口

- [x] Linux amd64 AppImage 已產出：`wolong-remake-linux-amd64-20260811.AppImage`（已被 8/12 的 [`dist-all/packages`](dist-all/packages) 取代）。
  AppDir 根目錄含 `.desktop`／`AppRun`，deny-list 通過，Docker／Xvfb `APPIMAGE_EXTRACT_AND_RUN=1`
  啟動與 640×400 固定 seed 截圖 smoke 通過；不含原版資料與字型。
- [x] 新增「經典再現」原版／remake 比較片：
  [`wolong-remake-classic-revival.mp4`](dist-all/promo/wolong-remake-classic-revival.mp4)，60 秒、
  1280×720、H.264/AAC。原版側使用使用者 YouTube 的代表幀，remake 側使用固定 `seed=17`
  實機代表幀；影片與 [`docs/promo/classic-revival.md`](docs/promo/classic-revival.md) 明示
  `core=normal`、`cputype=486`、`cycles=20000` 的 DOSBox 重播原則，以及不宣稱同狀態逐像素 parity。
- [x] Android 觸控 shell debug APK 已建置並在 API 35 `google_apis;x86_64` 模擬器安裝／啟動；
  1080×1920 實體畫布旋轉成 1920×1080 橫向畫面，`CONTINUE` 顯示 TALK 頁、`MENU` 開啟命令抽屜。
  三張證據圖見 [`docs/mobile/android-plan.md`](docs/mobile/android-plan.md)；完整自然時鐘、事件、
  存檔／讀檔、實機與 release signing 仍未完成。
- [x] Android 驗證容器已停止並確認沒有留下 `wolong` 專案容器；本輪沒有建立 `HANDOFF.md`，
  後續工作仍只追加在本檔。

### 2026-08-11 事件 10 remake 封口

- [x] 使用者所指的「不下命令、不移動滑鼠時，部隊依已下達指令自動跑、日期時間流逝」
  已在正常 `wlgame` 主迴圈完成，不再僅是 state fixture。游標每次移動都會暫停該 frame，
  穩定下一 frame 才恢復；modal 仍依既有規則暫停。
- [x] 月結俘虜的 `0x41`／`0x42` substitute、raw `0x0A` queue 與每時 TALK consumer
  仍完整串接。原版自然 `0x0A` writer 仍是 **unknown**，但不再阻擋這項 remake 功能；
  不把 substitute 寫成松崗原版已證實的 producer。
- [x] 松崗繁中唯一 oracle 的 30 frame Xvfb smoke hash：
  `45a68852335420dd7b22b4e240192dcd7a38fbbc62f72c8c59ec95acdc137b24`。

## 仍未宣稱完成的邊界

- YouTube 原版遊玩影片與 remake 推廣片已完成並排比較，並以 640×400 自然畫面保存
  raw pixel diff、差異圖與 24 秒研究對照片；這已封閉可見像素差異／畫面骨架的驗收，
  詳見 [`docs/promo/yt-remake-pixel-review.md`](docs/promo/yt-remake-pixel-review.md)。
  當時還沒做「同日期／同輸入／同狀態」的逐像素 parity。
- Windows／macOS 目前有交叉建置／檔頭 gate，沒有目標作業系統原生 GUI runtime；三平台候選
  封裝已產出，但原生 GUI smoke 尚未完成。本輪推廣影片已產出，不取代目標平台驗收。
- Android 目前只有已在模擬器驗證的觸控 shell debug APK；完整核心、自然 clock、事件流程、存檔／
  讀檔、實機／平板與 release signing 仍未完成，不能標成 Android release。
- 完整長程遊戲測試依使用者指示略過；這不影響上述短 fixture 的封口。

## 證據入口

| 主題 | 入口 |
|---|---|
| 總體狀態與歷史勘誤 | [`CONTEXT.md`](CONTEXT.md)、[`RESEARCH-LOG.md`](RESEARCH-LOG.md) |
| 事件 6／7 對話索引 | [`docs/re/12-diplomacy-dialogue.md`](docs/re/12-diplomacy-dialogue.md) |
| DOS/V 數值視窗 | [`docs/re/13-pc98-numeric-window.md`](docs/re/13-pc98-numeric-window.md)（檔名為歷史名稱，本輪內容須以 DOS/V 為準） |
| MCH 物件格式 | [`docs/re/14-mmap-mch-objects.md`](docs/re/14-mmap-mch-objects.md) |
| 事件原版 fixture | [`docs/playtest/11-event6-original-fixture.md`](docs/playtest/11-event6-original-fixture.md) |
| 同狀態截圖規則 | [`VERIFICATION-MATRIX.md`](VERIFICATION-MATRIX.md)、[`docs/playtest/12-event3-same-state-parity.md`](docs/playtest/12-event3-same-state-parity.md) |
| M7／事件 2–5／事件 9 抽樣 | [`docs/playtest/14-m7-review.md`](docs/playtest/14-m7-review.md)、[`docs/playtest/15-event2-5-talk-sampling.md`](docs/playtest/15-event2-5-talk-sampling.md)、[`docs/playtest/16-event9-long-route.md`](docs/playtest/16-event9-long-route.md) |
| 推廣影片與配樂 | [`docs/promo/README.md`](docs/promo/README.md)、[`docs/promo/dosv-live-comparison.md`](docs/promo/dosv-live-comparison.md)、`dist/promo/wolong-remake-trailer.mp4`、`dist/promo/wolong-remake-classic-revival.mp4`、`dist/promo/wolong-remake-dosv-live-comparison.mp4` |
| IDA DOS/V 證據 | `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm`、`func-sub_*.txt` |

DOS/V 證據輸入：`KI.EXE` SHA-256
`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`，IDA `.i64`
SHA-256 `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`，工具為
IDA Pro 9.4，位址基準為線性位址。

## 驗證與清理

所有 Go、Python、GUI、DOSBox、截圖與建置都在 Docker 內執行；主機只做 Docker 控制、
Git 狀態檢查與 `apply_patch` 編輯。每次工作結束至少執行：

```text
docker run --rm --network none ... go test -p=1 -vet=off ./... -count=1
docker run --rm --network none ... go vet ./...
docker run --rm --network none ... tools/talkdat_selftest.py
docker run --rm --network none ... tools/index.py generate
docker run --rm --network none ... tools/index.py check
git diff --check
git status --short
docker ps -a
```

一次性 container 必須 `--rm`、有界 CPU／記憶體／PID、以目前 UID/GID 寫入；原始素材
唯讀掛載。明確刪除本輪 `/tmp` fixture、截圖與 binary，不做廣域 `/tmp` 清理；不得停止
其他專案既有容器。交接時記錄沒有留下本輪持續執行 container。

## 不要重做或誤宣稱

- 不要把 PC-98 screenshot、游標或外框當 DOS/V release oracle。
- 不要把通過 remake 單測寫成原版 parity；必須附 DOS/V 位址、raw bytes、fixture 或
  同狀態畫面證據。
- 不要猜事件 6／7 次要 formatter、事件 10 producer 或未知 marker；#72 缺參數時
  必須維持 fail-closed。
- 不要把 `HANDOFF.md` 重新建立；後續接手、進度與剩餘工作只更新本檔。
- 不要把三平台候選包的交叉建置／檔頭驗證寫成 Windows／macOS 原生 GUI parity；正式發布仍
  需要目標平台短 smoke 與輸入、音訊、字型載入驗收。

## 2026-08-12 DOS/V 密碼頁驗證勘誤

- [x] 松崗 DOS/V 密碼頁的有效輸入橋接已由 `wolong-dosboxx:latest` 的
  `mouse_emulation=integration`／INT 33 640×480 建立；不是舊 `dosbox-run` 腳本的
  DOSBox 命令列假陽性。
- [x] 新的唯讀原版副本中，空白確認、`0000`、`1234` 都進入原版開場；空白與 `1234`
  的 10 秒後畫面雜湊相同。這表示「任意數字會過」可重現，且數字非必要。
- [x] 密碼頁不納入 remake，亦不再作為 DOS/V 原版自然流程／畫面採樣的阻擋理由。
  `PASS.*`／`YNFONT.EXE` 的實際比較語意、真實硬體行為、完整長程流程與同狀態逐像素
  parity 仍各自保持未驗證；詳見 `docs/playtest/18-dosv-password-verification.md`。

## 2026-08-12 一般玩家啟動／新遊戲殼層切片

- [x] 新增 `cmd/wlgame/launcher.go` 純狀態啟動殼層：`NEW GAME` 確認、四劇本選擇、由實際劇本資料篩出的合法玩家勢力／君主、確認與返回。
- [x] `-save-file` 存在時提供四槽 `LOAD DATA`；以既有四槽 overlay 解析判定可用槽，空槽 fail-closed，不讀取或修改 `SINARIO.DAT`。
- [x] 一般流程在確認後才呼叫 `startWorld`，再掛道路、AI、戰術資料、季節外框、硬體游標與相機；既有驗收入口使用明確 direct-start 白名單，不以任意旗標跳過 launcher。
- [x] 鍵盤與滑鼠共用 launcher selection state；首次靜止滑鼠位置不覆蓋鍵盤選取，實際滑鼠移動／點擊才改變選取。
- [x] `cmd/wlgame/launcher_test.go` 覆蓋新局成功、取消／返回、非法玩家、空槽拒讀與成功讀槽狀態轉移。
- [x] 正式視窗標題改為「臥龍傳－三國制霸之計」；本切片不處理戰術或自然 HUD polish，不宣稱逐像素 parity。
- [x] Docker/Xvfb 實際擷取修正後三張 640×400 launcher 畫格：`/tmp/wolong-launcher-title.png`（`60def98e0cf54726ad62794b92906017f200863241442db68aa3f536eb3b5150`）、`/tmp/wolong-launcher-scenario.png`（`9fc49e523b75e45b5a939177272a4b199fe1a36d4dd7fb67fc506b70d70eeb04`）、`/tmp/wolong-launcher-player.png`（`0b3a7f4709b538a6c9729a3bb745e3b49c523d999c8e7bfc596e1178cf5f38eb`）；確認 title／scenario 未選文字可讀，player 8 列、標題、反白與提示均未越過 panel 安全區。
- [ ] 四劇本與可用存檔槽的逐一代表畫面仍未錄製；不阻擋本小切片交付，完整長程測試依使用者要求不跑。

## 2026-08-12 Mentor polish 收斂

- [ ] DOS/V 戰術畫面已有可玩骨架與鍵鼠 dispatcher，但推廣片複驗證實戰場 viewport、上下 TALK 區、右側縮圖／狀態／命令 glyph 與底列配置尚未對齊原版；撤銷「完整骨架」宣稱。
- [x] 戰術縮圖由 `BATTLE.MAP` 與 `BATTLE.MDL` attribute 動態產生 128×128 圖，不再使用高度圖替代。
- [x] 依 `sub_1075B` 已證實公式接入開戰兩筆 `TALK.DAT`；每場只初始化一次、未知 marker fail-closed、戰場時間不因對話停止。實機證據：`docs/images/wlgame-tactical-opening-talk.png`。
- [x] 已證實兩種敗北 latch：信賴度歸零優先顯示 TALK #414；最後據點失守使用克制 fallback。研究備註不進 GUI。
- [x] 存讀檔四槽與一覽表新增滑鼠／觸控列、分頁、確認、取消；遵守原版兩段式選取且 modal 不穿透背景。
- [ ] 指揮／事件／一覽畫面的 DOS/V 幾何 parity 未完成；現有推廣片「指令與事件」左右不是同類畫面，不能作為還原證據，必須重做同類畫面對拍。
- [x] 推廣片異類「事件 vs 目的地」鏡頭已撤換為「原版系統設定 vs remake 系統設定」；戰術段已重錄為 `240:80` 主要幾何。影片視訊與音訊均為 60 秒。
- [x] 系統設定中央五列、事件左下 TALK、一覽第一層主要外框已依松崗錄影正規化座標修正。
- [ ] 戰術右欄原版命令 glyph／內框、一覽左側捲軸與選取後的前層武將詳細窗仍未完成；不得以本輪主要幾何修正宣稱完整 parity。
- [x] Docker＋Xvfb `go test -p=1 -vet=off ./... -count=1` 與 `go vet ./...` 通過；文件索引 65 份通過。
- [x] 三平台包與「經典再現」推廣片已以本輪畫面重建；`dist-all/` deny-list 與全部 SHA-256 驗證通過。實機對照主片的新版戰術段含 TALK、縮圖、六指令與動態戰場。
## 2026-08-12 攻城／兩軍遭遇共用戰術骨架

- [x] `combat.Siege` 與 `combat.Field` 的實際繪製、TALK slot、右欄、底部六命令及鍵鼠命中區均無模式分支；唯一幾何來源為 `dosvBattleLayoutFor`。
- [x] 新增 `TestFieldAndSiegeShareExactDOSVBattleChrome`，防止後續 polish 分裂成兩套畫面或命中區。
- [x] 修正 `demoBattle` 預跑 900 tick 導致野戰在 GUI 出現前結束的假差異；改為 120 tick。
- [x] Docker／Xvfb 六格實機對照完成，見 [`docs/playtest/22-field-siege-shared-layout.md`](docs/playtest/22-field-siege-shared-layout.md)；全專案 `go test -p=1 ./...` 與 `go vet ./...` 通過。
- [x] 純城兵據點攻擊保留原版自動判定，不擅自導入手動戰術畫面。
- [x] 原版指令 glyph／右欄複合面板已依 `sub_1C7F4`、`sub_1C863`、`sub_1F888`、`sub_1C6BF` 接入；舊的 2×3 文字格與綠色選取框已移除。
- [x] 三平台完整包與 AppImage 已重建至 `dist-all/`；Linux tar／AppImage Xvfb smoke 皆維持基準 SHA-256 `45a68852…b5150`，deny-list 掃描 19 個交付檔並比對 120 個原版檔後通過。
- [x] DOS/V／remake 推廣片的 40–45 秒戰術段已替換為原版 glyph／右欄面板實機畫面；繁中 overlay 以 Noto CJK 代表幀驗收，成品 60 秒、1800 幀，SHA-256 `feddd663…3961e8`。
- [ ] 兩模式仍共用同一組未完成項：戰術動畫時序與同狀態逐像素 parity。
# 2026-08-12 推廣片 AdLib／戰術骨架勘誤

- [x] 移除推廣片對 `tools/promo_score.py` 的依賴，改用使用者松崗 DOS 原版錄影中
  的實際遊戲 AdLib 音軌；來源與權利邊界見 `docs/promo/dosv-adlib-and-tactical-review.md`。
- [x] 原版錄影先正規化／裁成 640×400，再與 remake 共用最近鄰縮放鏈。
- [x] 戰術離屏 buffer 從 496×384 收斂到 480×368，消除右、下各 16 px 的遮切。
- [ ] 完成同戰況戰術 capture pair：同攻城節點、攻守方、編成、命令、鏡頭與 frame。
- [ ] 依 pair 修完雙 TALK、右欄完整狀態、底列／側欄 glyph 與選取時序。
