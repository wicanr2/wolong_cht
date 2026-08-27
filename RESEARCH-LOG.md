# Reverse-engineering research log

## Target

- Game/build：NEO･GETEN《臥龍傳－三國制霸之計》；`dosv`／`pc98`
- Platform：DOS/V 與 PC-9801
- Input hashes：完整雜湊與輸入清單見 `docs/re/01-first-recon.md`；本台帳不複製省略雜湊
- Tools and versions：IDA Pro 9.4、既有 `wolong-dosboxx:latest`、專案 `tools/` 包裝器
- Original data location (read-only)：`workplace/orig/{dosv,pc98}`、`workplace/orig/pc98_fdi`
- Writable experiment location：`/tmp` 或明確測試輸出目錄

## Current question

本輪接手先不重新逆向已定案格式；問題是「如何以目前證據安全完成 M7／M8，並讓下一個
session 能在不重推舊結論的情況下繼續」。

## Prior claims

| Claim | Source | Evidence status | Conflict |
|---|---|---|---|
| M0–M2 完成、M3–M6 大致完成 | `CONTEXT.md` §2–§3 | 承接既有紀錄 | 本輪尚未重跑所有門禁 |
| `15-realtime.md` READY | `CONTEXT.md`、`docs/INDEX.md` | READY／confirmed | 無；不要再把即時制當 blocker |
| M7／M8 是主要剩餘工作 | `CONTEXT.md` §7.0 | 強證據（狀態摘要） | 需以本輪測試校準細節 |
| `#192`–`#195` 是槽位錯位 | `docs/re/12-diplomacy-dialogue.md`、`docs/reference/02-jp-cht-diff.md` | IDA 已證實三變體直取，16 筆校訂可重跑；#751 已採強證據最小修正 | 不應直接重譯；畫面行寬與全量文意仍未驗收 |

## Evidence

| Address/offset/artifact | Observation | Independent check |
|---|---|---|
| `docs/INDEX.md` | 本輪重生後列出 136 條斷言，並區分 confirmed／強證據／未解 | `tools/py.sh tools/index.py generate/check` 通過 |
| `docs/mechanics/15-realtime.md` | 五層時鐘與固定 remake 時間基準已有文件化證據 | `docs/re/06-game-clock.md`、歷史 `wlgame` 截圖 |
| `docs/playtest/06`、`07` | PC-98 可用畫面雜湊與閉迴路滑鼠進入遊戲本體 | 兩次相同操作 byte-for-byte 截圖紀錄 |
| `tools/release.sh windows/amd64` | 產出 `dist/windows-amd64` 與本機 `dist/linux-amd64`；PE／ELF 檔頭分別正確 | `tools/py.sh tools/denylist.py dist` 通過 |
| `tools/go.sh run ./cmd/wlsim ... -years 1 -check` | 一年 12 個月、77760 tick，無不變量違反 | 真實 `workplace/orig/dosv/SINARIO.DAT` 載入 |
| `tools/shot.sh docs/images/wlgame-release-smoke.png` | `wlgame` 在 Docker/Xvfb 啟動，道路圖 253 條 | app log 與展示 PNG；只作啟動 smoke，不冒充完整玩家路徑 |

## 2026-08-09 — PC-98 原版正常入口與遊戲本體畫面

| 斷言 | 證據 | 等級／邊界 |
|---|---|---|
| PC-98 原版可由新遊戲選單進入劇本 1 | `wolong-dosboxx:latest`、`machine=pc98`、`cycles=20000`、`noopen`；`until:<NEW GAME md5>`、`Ctrl+F10`、閉迴路相對滑鼠；最後 `(352,280)` | **confirmed**；截圖 `pc98-oracle-scenarios.png`、`pc98-oracle-in-game.png` |
| 兩段式選取確實存在 | `pc98-oracle-cao-highlight.png`：第一下只把曹操列反白；第二下才進君主確認 | **confirmed**；與 `internal/ui/listwin` 的兩段式狀態對上 |
| 遊戲本體的日期／據點資訊可作視覺 oracle | `pc98-oracle-in-game.png` 顯示 `196年 4月 1日`；`pc98-oracle-city-panel.png` 顯示城兵 1,140、生產力 15,857、上昇值 0、防災值 100 | **confirmed**；尚未證明真實月長或原版有效存檔寫回 |

原始 PC-98 檔案以唯讀掛載；截圖只寫入 `docs/images`／`workplace/shots`，容器以 UID/GID
`1000:1000` 執行。這筆證據把舊的「只到確認框」狀態修正為「可進入遊戲本體」，但不升格
為 remake 同狀態 parity 或完整 `SAVE.DAT` oracle。

## Conclusion

- Status：strong inference（接手範圍）；工具鏈與 smoke 已驗證，尚未宣稱本輪完成 remake
- Behavior：先重建可重現的驗證基線，再處理校訂與發行，不重新挖已 READY 的格式
- Remake mapping：以 `CONTEXT.md` worklist 為狀態入口，`docs/re/`／`docs/mechanics/` 為證據與機制雙帳
- Remaining uncertainty：M7 1,022 則逐句文意／畫面抽樣與 M8 目標平台實機／完整正常玩家路徑驗收

## Reproduction

```sh
# 以上命令須在 Docker 內執行；原版素材唯讀掛載，輸出使用 /tmp 或明確目錄。
tools/go.sh vet ./...
tools/go.sh test ./...
tools/py.sh tools/index.py generate
tools/py.sh tools/denylist.py --selftest
```

## Stale references updated

- [x] `AGENTS.md` 接手順序與硬規則
- [x] `MEMORY.md` 快速恢復摘要
- [x] `WORKLIST.md` 本輪交接邊界（`HANDOFF.md` 已刪除）
- [x] `CONTEXT.md`（補入 2026-08-09 Docker-only 發行 smoke 與 `tools/py.sh`）
- [x] `docs/formats/01`、`docs/re/07`、`docs/re/12`、`docs/reference/02`（本輪補入訊息格式化器、外交索引與校訂勘誤）
- [x] `docs/mechanics/70-ai.md`（本輪補入停戰說服三變體直取規則）

## 2026-08-09 — save/load vertical slice

| Claim | Evidence | Status |
|---|---|---|
| `wlgame` 可由系統視窗操作四槽存檔／讀取 | `cmd/wlgame/save.go`、`window.go`；Xvfb `4→S→Return` smoke | 已證實（remake 行為） |
| 原始素材不被就地改寫 | `internal/savepath` fail-closed path test；`SINARIO.DAT` SHA-256 `21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c` 在前後相同 | 已證實 |
| 寫入後仍是完整四槽格式 | overlay 大小 `88,832 B`；前進後第 1 槽與原始檔有差異，其餘格式由 `World.SaveInto`／round-trip 測試約束 | 已證實（已知欄位範圍） |
| 重新啟動可載入 overlay 並保留道路圖 | Xvfb app log：`讀取存檔 overlay`、`道路圖：253 條路`、日期 `196年4月8日` | 已證實 |
| `Trust` 與 `Player` 已完整持久化 | `Trust` 已定位為區塊 `+0x10`（`cs:0D00h`／IDA `byte_10D00`）並接入；`Player` 的 `word_10CFD`／`byte_10CFF` 對應區塊 `+0x0D`／`+0x0F`，並接入雙欄一致性 round-trip | `Trust`／Player 儲存位址 confirmed；新遊戲初始化、完整事件增減與 handler parity 未完成 |

Trust／事件佇列的實作仍不代表完整原版 save parity；Player 欄位則由下一筆靜態證據
另行定案，保留與這個 remake vertical slice 的邊界。

## 2026-08-09 — Player 存檔欄位定案

| 斷言 | 證據 | 等級／實作邊界 |
|---|---|---|
| 玩家勢力直接值是區塊 `+0x0F` | `sub_11AC3` 新遊戲選定勢力後把 `bh` 寫入 `byte_10CFF`；`cs:0CFFh−cs:0CF0h=0x0F`；其他 AI／軍團分支以 `byte_10CFF` 比較玩家勢力 | **confirmed**；`World.Player` 讀寫 |
| 玩家勢力表位址是區塊 `+0x0D` | 同一支先把勢力表位址 `bx=faction×0x40` 寫入 `word_10CFD`；`cs:0CFDh−cs:0CF0h=0x0D` | **confirmed**；寫回時與 `+0x0F` 同步 |
| 有效存檔會直接使用玩家欄位 | `sub_18B40` 以 `AH=1` 呼叫 `sub_18CAE` 載入後，`sub_11AC3` 從選單分支跳到 `loc_11B2C`；該處直接讀 `word_10CFD`，沒有重新選勢力 | **strong evidence → implementation**；兩欄不一致時 remake fail-closed |
| 出貨空檔仍是無玩家狀態 | DOS/V 與 PC-98 `SINARIO.DAT`／`SAVE.DAT` 四槽 raw：`+0x0D=0xFFFF`、`+0x0F=0xFF` | **confirmed**；新劇本／空槽 `World.Player=-1` |

逆向輸入：DOS/V `workplace/ida/dosv/KI.EXE.asm` 與對應 `KI.EXE.i64`；KI.EXE
SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`；
工具 IDA Pro 9.4、既有 `ida-pro-9.4-ver2:uidfix-v1`，位址採 IDA `seg000:xxxxx`。
remake 以 `TestPlayerStorage` 驗證 `+0x0D`／`+0x0F` 一致性與 round-trip；尚未取得
原版實際玩後的有效 `SAVE.DAT`，所以不把這筆靜態證據升格為完整存檔行為對拍。

## 2026-08-09 — M7 correction application（`#751` 裁定前）

| Claim | Evidence | Status |
|---|---|---|
| 已定案的 15 筆變數／槽位／內容修正可重跑套用 | `tools/talkdat.py correct`、`tools/talkdat_selftest.py`；輸入 `translations/extract/talk-dosv.json` 與 `translations/corrections.json` | 已證實 |
| `#192`–`#195` 的既有繁中句子可按原版三變體組重排 | `docs/re/12-diplomacy-dialogue.md`、`translations/corrections.json`；selftest 15 筆套用、`#751` 保留 | 已證實（資料校訂） |
| 校訂後資料仍可組回 TALK.DAT | corrected JSON → `build` 34208 B → `verify` round-trip 相同 | 已證實（校訂 overlay，不是原版 parity） |
| 校訂後行寬已與原版畫面完全相同 | `{N}` 目前按三格估算；#321、#192、#195 行數變更，#718、#193、#194 等有行寬警告 | 未知，需畫面抽樣 |

此切片仍不改動 `workplace/orig/`，也不把校訂 JSON 偷渡進 runtime；它先固定可審查的文本產物。

## 2026-08-09 — `#257/#258` 內容對調與 `#751` 保留裁定（歷史勘誤前）

| Claim | Evidence | Status |
|---|---|---|
| 繁中 `#257` 與 `#258` 的內容互換 | `translations/extract/talk-{dosv,pc98}.json`：日文 #257 對應繁中 #258，日文 #258 對應繁中 #257；兩則皆無變數 | 已證實；逐字內容對照 |
| `#257/#258` 可安全重用既有繁中句子 | `translations/corrections.json` 的 `content-swap` fix；`tools/talkdat_selftest.py` | 已證實（資料校訂）；不宣稱重譯 |
| `#751` 少了日文的 `{1}` | 日文 `{6}なんの！これしきのことで参る{1}では無いわ！` 對繁中 `{6}沒什麼！我可不會這樣就輸的！` | 已證實（標記缺失） |
| `#751` 的繁中插入位置與完整語氣已定案 | 目前沒有同句既有譯文或畫面／正常路徑證據 | 未知；`fix:null`，工具保留人工裁定 |

本輪 `corrections.json` 共 16 筆，其中 15 筆可重跑套用、`#751` 保留；selftest
同時驗證校訂後 round-trip 與現況不符時拒絕。行數變更與行寬警告仍需 CJK 畫面抽樣。

## 2026-08-09 — TALK.DAT formatter 證據與 #321 校訂

| Claim | Evidence | Status |
|---|---|---|
| `sub_1075B` 以 `CX` 取 TALK.DAT 偏移並把 `DI` 傳給文字格式化器（formatter） | IDA Pro 9.4 `func-sub_1075B-current.txt`：`0001075B` 的偏移表讀取、`000107C9` 呼叫 `sub_106F9` | 已證實；IDA 線性位址 |
| `sub_1084A` 遇 `\N` 從 `CS:[SI+08A4h]` 七項表分派 | `func-sub_1084A-current.txt`：`0001085D`–`0001089A`；原始 `KI.EXE` 檔案偏移 `0x0AA4` 的七個 little-endian handler 位址 | 已證實；檔案偏移與線性位址分開記錄 |
| `\6` 消耗一個 16 位元參數並使 `DX -= 0x30`，不繪出參數 | `KI.EXE` handler 線性位址 `0001097E`–`00010983`：`inc di`×2、`sub dx,30h`、`retn` | 已證實 |
| `sub_10CDE` 不是訊息參數準備 | `func-sub_10CDE.txt`：`00010CDE` → `sub_1EB11`；`func-sub_1EB11.txt`：`IN/OUT 61h`。後續參數資料流是 `DI=SP` → `sub_18810` → `sub_1075B` → `sub_1084A` | 已證實；推翻 `docs/re/07` 舊註解 |
| #321 應補回開頭 `{6}` 並修正「以」 | 日文 #321 有 `{6}{6}{7}` 的出現次數；格式化器已證實第一個 `{6}` 會消耗一個參數；`translations/corrections.json` fix 為 `{6}外交資金已用完，大約需要金額{6}{7}左右。` | 已證實（文本校訂）；尚未宣稱原版畫面排版 parity |

本輪逆向輸入為 `workplace/ida/dosv/KI.EXE.i64`／同目錄 `KI.EXE`，SHA-256
`FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`；工具為
`ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）。所有 `000xxxxx` 是 IDA 線性位址，
`0x0AA4` 是原始 EXE 檔案偏移，不能混用。

## 2026-08-09 — 停戰說服 `#190`–`#198` 槽位索引

| Claim | Evidence | Status |
|---|---|---|
| `sub_13830` 以 `CX=0096h` 作停戰說服訊息基底；要求說明後加 `10h` | `func-sub_13830.txt`：`00013830`–`00013889`；`00013885` 的 `add word ptr [bp+0],10h` | 已證實；IDA 線性位址 |
| `sub_13B5A` 把理由／結果值交給 `sub_13BA9` | `func-sub_13B5A.txt`：`00013B5A`–`00013B70` | 已證實 |
| `sub_13BA9` 對 `AL=2` 產生 `#190/#193/#196` 三個組的起點 | `func-sub_13BA9-current.txt`：`00013BFA`–`00013C19`；`bl = AL×9 + AH×3 + 6`，再加訊息基底 | 已證實；`AH` 的完整欄位語意未知 |
| `sub_13C99` 依 `SI+1Eh` 將變體索引直接加到 `CX` | `func-sub_13C99-current.txt`：`00013CA6`–`00013CB4` | 已證實；值達 3 的分支另減 3，未宣稱欄位只限 `0..2` |
| `#190–#192`、`#193–#195`、`#196–#198` 是三個連續三變體組，不是隨機抽句 | 上述兩支索引計算＋`TALK.DAT` 日中槽位內容對照 | 已證實（取用規則）；組內中文語氣只作既有譯文重排 |
| `#192`–`#195` 可用既有繁中句子重排 | `translations/corrections.json` #192–#195；`tools/talkdat_selftest.py` 15 筆套用、`#751` 保留、round-trip、mismatch guard | 已證實（資料校訂）；尚未宣稱畫面行寬 parity |

輸入為 `workplace/ida/dosv/KI.EXE.i64`／同目錄 `KI.EXE`，SHA-256
`FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`；工具為
`ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）。`000xxxxx` 是 IDA 線性位址，
不是 EXE 檔案偏移。完整資料流見 `docs/re/12-diplomacy-dialogue.md`。

## 2026-08-09 — packaged screenshot gate

| Claim | Evidence | Status |
|---|---|---|
| 發行目錄的 Linux `wlgame` 可啟動並截圖 | `tools/release.sh windows/amd64` 後直接執行 `dist/linux-amd64/wlgame -shot -shot-frames 120` | 已證實（Docker/Xvfb） |
| 截圖旗標不會在 Draw 前提前終止 | `cmd/wlgame/main.go` 的 `shotDone`；log 顯示第 120 幀 PNG、大小 74,256 B | 已證實 |
| Windows／Linux／mac 目標平台都已實機執行 | 目前只有 Docker Linux packaged smoke；Windows 為 PE 交叉建置，mac 未在本環境實跑 | 未完成 |
| 本輪遭遇選單改動後，Linux packaged `wlgame` 仍可啟動並截圖 | `dist/linux-amd64/wlgame -open-battle-choice -shot -shot-frames 60`；Xvfb PNG 62,610 B、UID/GID `1000:1000` | 已證實（Docker packaged smoke） |

## 2026-08-09 — encounter decision vertical slice

| Claim | Evidence | Status |
|---|---|---|
| 玩家勢力捲入軍團遭遇時先進入「戰鬥指揮／委任」狀態 | 原版分派條件 `docs/re/09-combat.md` §2；`internal/state/corps.go` 建立 `EncounterChoice` | 已證實（原版條件／remake 接點） |
| 遭遇選單掛起時戰略時鐘不推進 | `internal/state/state_test.go` 的 `TestPlayerBattleGoesTactical`、`TestPlayerBattleCanBeDelegated` | 已證實（remake 行為） |
| 「戰鬥指揮」與「委任」共用一致的戰後出口 | `ChooseBattleCommand` → `PendingBattle`；`ChooseBattleDelegate` → `combat.Resolve`／`CorpsEvent`；完整 Go tests | 已證實（remake 行為） |
| 真實劇本素材下選單版面可顯示 | Docker/Xvfb `-open-battle-choice`；`docs/images/wlgame-battle-choice.png`；道路圖 253 條 | 已證實（畫面 smoke；驗收旗標，不是正常路徑證據） |

## 2026-08-09 — `wlgame` 正常編成／行軍切片

| Claim | Evidence | Status |
|---|---|---|
| 不使用 `-open-*` 也能以真實預備兵資源完成編成 | Docker/Xvfb 鍵盤序列 `A` → 武將兩段式 `Enter`×2 → `3` → `Down+Space`×5 → `Enter`；`wlgame-normal-formed.png` | 已證實（正常玩家輸入） |
| 編成後能以 `M` 選軍團與目的地 | `wlgame-normal-destination.png`；目的地一覽按距離排序，第一項為虎牢關 | 已證實（正常玩家輸入） |
| 確認目的地後世界時鐘與行軍事件推進 | `wlgame-normal-march.png`；`=` 將速度由 0 調為 1，日期 196/4/1 → 196/4/2，事件為「曹操 向 虎牢關 行軍」 | 已證實（正常玩家輸入） |
| 正常行軍能抵達敵方城兵攻城 | 無 `-open-*` 另選袁術據點「汝南」；`wlgame-normal-garrison.png` 顯示「曹操 對 城兵　攻下 汝南」；`TestNormalScenarioMarchIntoGarrison` 通過 | 已證實（正常玩家輸入／規則回歸） |
| 同一條正常路徑已抵達敵方軍團遭遇戰 | 劇本開局沒有敵方軍團；城兵攻城依原版自動判定，不產生 `EncounterChoice` | 未完成 |

此切片沒有替未解的戰術傷害數值、委任 AI 品質或完整玩家路徑宣稱 parity；它補上
真實道路到城兵自動攻城的閉環，也保留敵方軍團與戰鬥選單仍缺正常 AI 來源的界線。

## 2026-08-09 — 使用者存檔遭遇回放核查

| Claim | Evidence | Status |
|---|---|---|
| 不使用 `-open-*` 也能從含軍團的存檔回放到遭遇選單 | Docker/Xvfb 以 `workplace/orig/dosv/SAVE.DAT` 第 1 槽複本作 `-save-file`；`wlgame-save-replay-choice.png` 顯示「張飛 對 許褚／攻城」及「戰鬥指揮／委任」 | 已證實（存檔回放／畫面接縫） |
| 回放狀態具備結構化軍團資料，且原始檔沒有被寫入 | Docker 內讀取 overlay；道路圖 253 條，畫面進入遭遇選單；工作樹原始 `SAVE.DAT` 唯讀掛載 | 已證實（隔離回放） |
| 該存檔可作正常開局或正常 AI 時序 oracle | `state.LoadScenario` 記錄起始為 0 年 0 月 0 日，截圖時為 0 年 1 月 1 日；與真實劇本 196/4/1 開局不一致 | 未證實（fail-closed） |

這筆證據把「無除錯旗標的遭遇 UI」與「正常開局的戰略時序」分開記錄；不能把
使用者提供的零時鐘存檔回放升格為敵方 AI 已完成或原版存檔 parity。

## 2026-08-09 — 政略 AI 目標、宣戰與敵方編成垂直切片

| 斷言 | 證據 | 狀態 |
|---|---|---|
| 據點的四個鄰接槽可作為政略 AI 候選來源 | DOS/V `KI.EXE` 的 `sub_12C52`／`sub_12CDF` 讀 `+0x1C..+0x1F`；`state.LoadScenario` 保留原始四 bytes；`TestScenarioOneStrategyNeighbourOrder` 得到曹操 `[13,12,11,21,2]` | 已證實／已接入 |
| 非玩家與玩家勢力的月度交友度漂移不對稱 | `sub_12DB8`：和平 −2 下限 20、交戰 +1 上限 50 且玩家例外；`sub_12DF3`：玩家第一鄰居 −1、交戰再 −7，其他勢力對玩家 −1 | 已證實／已接入 |
| `sub_12EFB` 的宣戰條件是三道嚴格閘 | 資金 high word 必須嚴格大於 `min(據點數×16+64, 0x61A)`；交友度不得大於 `0x80+20+好戰+floor(好戰/2)`；己方國力至少目標的四分之三；`internal/rules/strategyai` 單元測試 | 已證實／已接入 |
| 國力計算不可用一般有號資金右移取代 | `sub_13091` 先以 24 位原始資金的高兩 bytes 作 unsigned word 與 `0x13` 比較；負資金測試固定住這個分支 | 已證實／已接入 |
| 敵方編成六槽兵種順序來自 `CS:6C4C` | DOS/V `KI.EXE` 檔案段位址對應 bytes：`01 03 02 / 01 03 02 / 03 01 02 / 03 01 02 / 02 03 01 / 02 03 01`；1＝騎馬、2＝弓、3＝步兵；每槽預備兵門檻 `0x32` | 已證實／已接入 |
| 編成武將與分兵規則可接到 remake 軍團 | `sub_145C1` 選未出陣最高武力武將；`sub_14698` 同兵種剩餘兵力平均分配、單槽最多 100；`formAICorps` 與真實劇本測試 | 已證實／已接入 |
| 事件 8 不是侵攻起點 | `sub_12D3A` 發事件 8；`sub_133EA` 無 `retn` 並落入 `sub_16A3D`，後者掃自身據點選新首都；侵攻候選在 `sub_12D58`／`sub_12EFB` | 已證實；修正舊文件 |
| 固定亂數下真實劇本能產生敵方戰略戰鬥 | `TestStrategicAIScenarioOneProducesEnemyWarPath`：玩家 0、種子 17、六個 settled months；宣戰 5、編成 4、戰鬥 4；每 tick `CheckInvariants` 通過 | 已證實（remake 狀態層） |

本輪逆向輸入為 DOS/V `workplace/ida/dosv/KI.EXE.i64` 與同目錄 `KI.EXE`，
SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`；
工具為 `ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）。所有 `000xxxxx` 是 IDA
線性位址；`CS:6C4C` 是原始程式段位址；其對應 EXE 檔案偏移 `0x6E4C` 只作
獨立檔案定位，三種位址基準不可混用。

remake 端的明示差異：`World.runStrategicAI` 尚未持久化或重現原版事件佇列的
延遲與月度壓縮，而是在月結邊界直接套用宣戰；敵方目前至多每勢力一支軍團，
目的地採 MMAP 道路圖上最近的敵方據點。`sub_14575` 多軍團請求、
`sub_142AB`／`sub_14300`／`sub_14325` 的完整行軍狀態機，以及協力／停戰 AI
決策仍列為未完成；因此本紀錄是可重播的垂直切片，不是完整 AI parity。

## 2026-08-09 — 正常 `wlgame` 敵方 AI 遭遇路徑

| 斷言 | 證據 | 狀態 |
|---|---|---|
| 固定亂數可在正常 `wlgame` 輸入中重播 AI 軌跡 | `cmd/wlgame -seed 17`；預設仍為 `rng.Now()`，固定種子只在驗收時取代播種來源 | 已證實（驗收介面） |
| 真實劇本可由正常鍵盤完成玩家軍團編成與目的地選擇 | Docker/Xvfb：`A`、選曹操、`3`、`Down+Space×5`、`Enter`；`M`、軍團兩段選取、距離排序第 22 列、兩段確認 | 已證實（正常玩家輸入） |
| 敵方 AI 軍團能沿現有道路自然抵達並使玩家進入遭遇選單 | `wlgame-ai-normal-encounter.png`；196/6/19；畫面為「呂布 對 曹操／攻城／戰鬥指揮／委任」；狀態列為「呂布 軍團向 濮陽 行軍」 | 已證實（正常玩家路徑／Docker/Xvfb） |
| 這張畫面不是 `-open-*` 或座標傳送捷徑 | 啟動沒有 `-open-battle`、`-open-battle-choice`、`-open-form`、`-open-corps`；只使用一般鍵盤輸入與固定亂數 | 已證實 |

回放使用 DOS/V 真實 `SINARIO.DAT`、`workplace/eten/`，原始素材唯讀掛載；輸出 PNG
由 UID/GID `1000:1000` 寫入。`-seed` 是 remake 的可重播驗收入口，不是原版存檔欄位，
也不宣稱原版亂數時鐘播種 parity。畫面證據與操作序列見
`docs/playtest/08-wlgame-normal-strategy-path.md`。

仍未由此證明的項目：原版有效時鐘存檔與 remake 的逐畫面對拍、戰鬥指揮後完整戰術
戰鬥／戰後出口、事件佇列延遲、多軍團請求，以及原版完整行軍狀態機。

## 2026-08-09 — 戰術近戰公式與繞路狀態修正

| 斷言 | 證據 | 狀態 |
|---|---|---|
| 一般近戰不是固定 6 點 | DOS/V `KI.EXE` `seg000:1B618`：`rand & 0x7F` 加攻擊者 `+0x18`，`0x46` 為命中門檻；傷害使用攻擊者戰力 | 已證實／已接入 `meleeHit` |
| 有利／不利的近戰效果分別作用在傷害／命中值 | `seg000:1ADC8` 產生 `byte_1D31E`；`seg000:1B618` 有利 `+0x40`、不利 `−0x32` | 已證實／已接入 |
| 突擊命中有 `+0xC8` 飽和加成 | `seg000:1B618`：`cmp byte ptr [si+1Ah],2` → `add ah,0C8h`，溢位設 `0FFh` | 已證實／已接入 |
| 大將命中旗標即使未命中也會留下 | `seg000:1B6BC` 尾端無條件 `or byte ptr [si+2],8`；remake `hitGeneral` 已在所有分支保留 `HitGeneral` | 已證實／已接入 |
| 隊長離場後不能補進無隊長的待機兵 | `seg000:1A83F` 對該隊七名隊員排命令 5；remake 同時清該隊 `Reserve`，避免 §5.9 永遠補不完 | 強推論／已接入；待原版待機摘要欄位的獨立對拍 |
| 繞路點不是每幀消費 | `seg000:1B00D` 只在 `sub_1AF69` 抵達目前 X/Y/Z 後呼叫；`Waypoints.Current/Advance` 與回歸測試固定此時序 | 已證實／已接入 |
| 普通飛道具初始威力為 `0x1C` | `seg000:1AD2D` 設 `CH=1Ch`，`sub_1B8AA` 寫入飛道具 `+0x04`；`sub_1AD7F` 的 `CH=20h` 仍分開列為未完成分支 | 已證實／普通分支已接入 |

本輪逆向輸入：`workplace/ida/dosv/KI.EXE.i64`、`workplace/ida/dosv/KI.EXE.asm`、
同目錄 `KI.EXE`；`KI.EXE` SHA-256 為
`FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
工具為 `ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）；本段 `seg000:xxxxx` 全是
IDA 線性位址，沒有與 EXE 檔案偏移混列。

## 2026-08-09 — 正常玩家戰術畫面延伸

| 斷言 | 證據 | 狀態 |
|---|---|---|
| 正常 AI 遭遇可由「戰鬥指揮」進入攻城戰術畫面 | 真實 `SINARIO.DAT`、固定種子 17、無 `-open-*`；從正常遭遇選單按 `Enter`；`wlgame-ai-battle-afterpatch.png` | 已證實（正常玩家輸入／Docker/Xvfb） |
| 戰術畫面顯示原版素材城壁、雙方兵力、城壁耐久與 1–6 命令列 | `wlgame-ai-battle-afterpatch.png`；`wlgame -shot` 第 1500 幀 | 已證實（呈現層切片） |
| 同一流程可送出 2 號攻擊命令並維持戰場更新 | 正常輸入序列最後送 `2`；`wlgame-ai-battle-attack-afterpatch.png`；`cmd/wlgame/battle.go` 的 1–6 → `tactical.Command(0–5)` 對應 | 已證實（remake 接縫）；尚未宣稱戰後出口 parity |

輸出 PNG 均由 UID/GID `1000:1000` 寫入，原始素材唯讀掛載。完整操作與未完成界線見
`docs/playtest/09-wlgame-normal-tactical-path.md`；本段不把中途戰術截圖升格為一場
完整戰鬥已結算，也不宣稱原版同狀態逐像素相同。

## 2026-08-09 — 正常真實攻城結算與退卻／高度移動修正

| 斷言 | 證據 | 狀態 |
|---|---|---|
| 正常玩家路徑的真實攻城可以在狀態層完成並回寫戰略層 | `internal/state/state_test.go` 的 `TestNormalScenarioTacticalBattleTerminates`；真實 `SINARIO.DAT`、固定種子 17、濮陽戰場 56、真實 `BATTLE.MAP`／`BATTLE.MDL`／`BATTLE.DAT`／腳本／編成；第 549 幀守方勝，攻方 0、守方 100，`ResolvePending` 後 pending 清除 | 已證實（remake 狀態層；不是原版 parity） |
| 退卻出口不是當前邊緣格的站立高度 | `sub_1AAED`：0x000 側寫 X=1、0x600 側寫 X=0x3E；Y 夾在 0x10..0x2F；目標 Z 寫 0；`TestRetreatUsesOriginalExitTarget` | 已證實／已接入 |
| 隊長消失時七名隊員改排退卻 | `sub_1A83F`：`mov cx,7`，逐筆把 Next command 寫成 5；`TestLeaderLossRetreatsSquadAndDropsReserve` | 已證實；清除 remake 待機槽是避免無隊長補兵死鎖的強推論 |
| 非爬牆兵可在水平跨格時同步一層高度 | `sub_1B1B1` 的水平移動分支會將 `[si+0A]` 調整 ±1；`sub_1AF69` 的兵種門檻只包住純 Z 軸嘗試；`TestHorizontalStepAdjustsOneLevelForNonClimber` | 已證實／已接入 |
| 大將／騎馬仍不能做純 Z 軸爬牆 | `sub_1AF69`：`cmp byte ptr [si+4],12h / jbe` 跳過 Z 移動；`TestCavalryCannotClimb` 保留四層城牆負例 | 已證實／已接入 |

本輪逆向輸入為 DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm` 與同目錄
`KI.EXE`；SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
工具為 `ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）；本段 `seg000:xxxxx` 是 IDA
線性位址，`.asm` 只作搜尋／可重生證據，未取代 `.i64` 的函式邊界與交叉參照。

這次修正只把正常攻城的狀態層結算、GUI 回戰略接縫與已證實的移動／退卻規則接起來；當時 GUI
戰後勝負／傷亡訊息、`sub_1AD7F` 的 `CH=0x20` 分支、完整原版投射物分支與原版／remake 同狀態
對拍仍不得宣稱完成。後續的結果報告證據見下一節。

## 2026-08-09 — 正常 GUI 戰術完成後回戰略層

| 斷言 | 證據 | 狀態 |
|---|---|---|
| 正常 GUI 戰術戰鬥完成後可按 Enter 回到戰略地圖 | Docker／Xvfb；真實 `SINARIO.DAT`、固定種子 17、無 `-open-*`；沿正常編成／行軍／遭遇／戰鬥指揮／攻擊命令路徑，戰場完成後以正常 `Enter`，輸出 `docs/images/wlgame-ai-postbattle.png` | 已證實（remake GUI 接縫；不單獨證明戰果數字） |
| GUI 戰後畫面與狀態層勝負證據可分開回查 | GUI 圖顯示已回戰略地圖；`TestNormalScenarioTacticalBattleTerminates` 顯示第 549 幀守方勝、攻方 0／守方 100，並由 `ResolvePending` 清除 pending | 已證實（兩種證據互補） |

本次 GUI 輸出由 UID/GID `1000:1000` 寫入，原始素材唯讀掛載；`wlgame` 只用一般鍵盤，
沒有 `-open-*`、座標傳送或強制勝利。畫面上的後續行軍事件不是戰後勝負的獨立文字 oracle，
所以不能用這張 PNG 取代狀態層傷亡測試，也不能宣稱原版／remake 同狀態 parity。

## 2026-08-09 — `sub_1AD7F` CH=0x20 特殊投射物分支

| 斷言 | 證據 | 狀態 |
|---|---|---|
| `sub_1AD7F` 的威力不是普通箭的 0x1C，而是 0x20 | IDA `seg000:1AD7F`：`mov ch,20h`；`sub_1B8AA` 寫入 projectile `+0x04`；`TestSpecialProjectileUsesCH20AndFallsVertically` | 已證實／已接入 |
| CH=0x20 會把方向加上 0x80，停止 X/Y 逐格移動 | `seg000:1AD7F`：`or bl,80h`；`seg000:1BA2E`：`cmp al,80h / jnb loc_1BA6A`，跳過 X/Y 位移；remake `projectile.special` | 已證實／已接入 |
| 特殊效果建立在攻擊者朝向相鄰格，垂直速度初值為 -0x100 | `seg000:1B8AA`：高位方向分支調整 X/Y；`seg000:1AD7F`：`mov ax,0FF00h`；`sub_1BA2E` 先命中檢查再更新高度；測試驗證高處弓兵三幀後命中 | 已證實／已接入 |
| 弓兵高處近距離會走特殊分支 | `seg000:1ABFF`：攻擊者 Z 高於目標且 `bx ≤ 1` 時呼叫 `sub_1AD7F`；`TestSpecialProjectileUsesCH20AndFallsVertically` | 已證實／已接入 |
| 步兵高層近距離會走特殊分支 | `seg000:1AC55`：`[si+1E] > [di+1E]` 且 `bx ≤ 2` 時呼叫 `sub_1AD7F`；remake 以 `Climbing` 適配，`TestClimbingInfantryCanUseSpecialProjectile` | 強推論／已接入；`+0x1E` 的完整資料來源仍待獨立對拍 |

本輪逆向輸入仍為 DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm` 與同目錄
`KI.EXE`；SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
工具為 `ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）；上述 `seg000:xxxxx` 是 IDA
線性位址。`sub_1AD2D` 普通箭路徑的完整原版運動／動畫與 `sub_1AD7F` 的原版同狀態
畫面仍未完成對拍，故本條不升格為完整投射物 parity。

## 2026-08-09 — 正常 GUI 戰後結果報告

| 斷言 | 證據 | 狀態 |
|---|---|---|
| 正常戰術戰鬥在 GUI 上先顯示可讀的勝負／兵力／攻城損害，再等待確認 | `docs/images/wlgame-ai-battle-result.png`；真實 `SINARIO.DAT`、固定種子 17、無 `-open-*`；第 549 幀，畫面顯示守方勝、攻方 5770→0、守方 1000→1000、損害 0 | 已證實（remake GUI） |
| 結果畫面與狀態層使用同一份戰果 | `CorpsEvent.BattleBefore`／`BattleAfter`／`BattleCityDamage`；`TestNormalScenarioTacticalBattleTerminates` 斷言事件欄位等於 `World` 回寫後的軍團狀態；戰略層 0／100 點 × 10 與畫面 0／1000 人一致 | 已證實（remake 狀態／GUI 接縫） |
| 結果畫面按 Enter 才回戰略層 | `cmd/wlgame/battle.go` 的 `updateBattle`；同一路徑 `docs/images/wlgame-ai-postbattle.png` | 已證實（正常玩家輸入／Docker-Xvfb） |

輸入與環境：`wolong-go:20260809`、Docker/Xvfb、DOS/V `SINARIO.DAT` 與
`BATTLE.*` 唯讀掛載、`-seed 17`、正常 `A`／`M`／方向鍵／速度／遭遇／戰術命令輸入；
輸出由目前 UID/GID `1000:1000` 寫入。這條證據只完成 remake 自身的結果資料流與 GUI 留證，
不代表原版同狀態畫面、完整投射物動畫或跨平台實機 parity 已完成。

## 2026-08-09 — 信賴度原始存檔欄位定案

| 斷言 | 證據 | 等級／實作邊界 |
|---|---|---|
| 信賴度不是勢力 `+0x1D` | `sub_16F26` 將勢力 `+0x1D` 複製到軍團 `+0x06` 士氣；四個劇本均為 200 | **confirmed**：`+0x1D` 是士氣基準 |
| 信賴度的原始 runtime 位址是 `cs:0D00h` | `byte_10D00` 的靜態讀寫：`sub_13D91` 增加、`sub_13DC9` 減少並在 0／`0xFF` 飽和；`sub_13C1E`／`sub_15F27` 用它繪製 UI | **confirmed**；IDA 線性位址 `seg000:10D00` |
| 信賴度會持久化到每個存檔區塊 `+0x10` | `sub_18CAE`／`sub_18CFF` 對 `cs:0CF0h` 起始的 0x3B byte 做載入／存檔；`0x10D00−0x0CF0=0x10` | **confirmed**；不是勢力表欄位 |
| 新遊戲初始值與所有事件增減已完全解出 | `sub_18B12` 仍有讀取勢力 `+0x2B` 的初始化路徑；`sub_13830` 可證實至少 `+20`／`−20` 結果分支 | **未完成**：`+0x2B` 的時序與其餘分支仍需原版動態 oracle |

逆向輸入：DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm` 與同目錄 `KI.EXE`；
SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
工具為 `ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）；本文 `seg000:xxxxx` 是 IDA 線性位址，
未與檔案偏移混列。`internal/state` 現以 `World.Trust` 讀寫 `+0x10`，並以
`TestTrustStorage` 驗證 0…255 鉗制與 round-trip；這不代表完整原版存檔／新遊戲 parity。

## 2026-08-09 — 事件佇列原始保存切片

| 斷言 | 證據 | 等級／實作邊界 |
|---|---|---|
| 事件佇列是區塊尾端的 1,024 byte | `sub_18CAE`／`sub_18CFF` 的第三段搬移、`docs/formats/08` §0；`+0x52C0`–`+0x56BF` | **confirmed**：256 筆 × 4 B |
| 每筆必須保留完整事件字與參數 | `sub_12FBF`／`sub_1301C` 的入佇列證據；`sub_131AE` 以事件字低 byte dispatch，高 byte 在災害／勢力路徑承載變體 | **confirmed**：`QueuedEvent.Code`／`Param` 均為 u16 |
| remake 可無損保存並重現已解出的 queue 邊界 | `World.LoadScenario`／`Bytes` 讀寫 `+0x52C0`；`TestEventQueueStorage` 驗證第 0、63、64、255 格與 `0x010C`／`0x020C` 高 byte；`TestEventQueueTiming`／`TestEventQueueMonthlyCompaction` 驗證 `byte_131AD=7`、每十次一筆與 `sub_12BD9` 的 64／192 壓縮 | **已接入**：raw／節拍／月壓縮通過，未宣稱完整 dispatch、玩家 UI 或事件增減 parity |

逆向輸入仍為 DOS/V `KI.EXE`／`SINARIO.DAT`（前述 SHA-256 與 IDA Pro 9.4 工具鏈）；
本切片不改寫原始資料與反組譯名稱，只把已定案的檔案區段接到可測試的存檔模型。

## 2026-08-09 — 事件佇列三個狀態 handler 接入（`sub_14502` 補接前）

| 斷言 | 證據 | 等級／實作邊界 |
|---|---|---|
| `sub_12FB1` 會把發起勢力放在事件字高 byte | `KI.EXE.asm` `sub_12FB1`：`si × 4` 後 `mov ah,ch`，再呼叫 `sub_12FBF`；`sub_131AE` 只以低 byte 查 dispatch 表 | **confirmed**；`queueEvent` 寫入 `Code = source<<8 | code` |
| 事件 1 的來源／目標資料流 | `seg000:1320C` 先以 `AH` 驗證來源，再以 `DL` 驗證目標，最後落入 `sub_13526`；`seg000:13526` 寫 `+0x19`、呼叫 `sub_135AB` 與 `sub_13639` | **confirmed**；`applyQueuedDeclaration` 接入宣戰、回頭宣戰與雙向交戰值 |
| 事件 8 是延遲處理的遷都 | `seg000:133EA` 以事件字高 byte 驗證勢力後呼叫 `sub_16A3D`，無 `retn` 落入 `sub_133FD` 寫回首都 | **confirmed**；本段記錄建立時 `dispatchQueuedEvent` 已接入遷都，`sub_14502` 軍團同步尚未補接（見後續勘誤） |
| 事件 13 固定扣 50 點信賴度 | `seg000:13507` 將 `AL` 設為 `0x32` 後呼叫 `sub_13DC9`；`sub_13DC9` 對 `byte_10D00` 做減法並歸零 | **confirmed**；`dispatchQueuedEvent` 以 `clampU8(Trust−50)` 接入 |
| remake 不再在月結直接套用宣戰 | `World.runStrategicAI` 以 `queueEvent` 寫入 code 1；`World.hourly` 先呼叫 `dispatchQueuedEvent`；`TestQueuedEventHandlers`、`TestStrategicAIScenarioOneProducesEnemyWarPath` 通過 | **已接入**；事件 2–7、9–12 的 UI／外交／災害 handler 仍未完成 |

本輪輸入為 DOS/V `workplace/ida/dosv/KI.EXE.i64`、同目錄 `KI.EXE.asm` 與 `KI.EXE`；
`KI.EXE` SHA-256 為 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
工具為 `ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）；本文 `seg000:xxxxx` 是 IDA
線性位址，未與 EXE 檔案偏移混用。`sub_135AB` 對中立目標的原版越界讀仍採
fail-closed remake 邊界，沒有把該差異寫成原版規則。

## 2026-08-09 — `#751` 變數缺失的最小校訂

| 斷言 | 證據 | 等級／實作邊界 |
|---|---|---|
| `#751` 的日文原文比繁中多一個 `{1}` | Docker 內以 `tools/talkdat.py export` 讀取 `workplace/orig/pc98/TALK.DAT`（SHA-256 `537e563269e414da79381ff48184a98e062ca454eb5d12e16a5fcbd52b79cf6f`）與 `workplace/orig/dosv/TALK.DAT`（SHA-256 `08a22e09791d0a6ec2968e87d8655e12c91b45e00fae460b28593b35ff85e384`）；索引 751 為 `{6}なんの！これしきのことで参る{1}では無いわ！`／`{6}沒什麼！我可不會這樣就輸的！` | **confirmed**：變數缺失 |
| `#751` 屬於武將戰場台詞池，而 `{1}` 是該語境的武將名插入值 | 同一批 export 的 `#681/#682/#683/#685/#741` 均保留 `{1}`；其中文句分別使用「我{1}」「我名叫{1}」「就由我{1}」等既有插入模式；`KI.EXE` formatter 證據見 `docs/formats/01-talk-dat.md` | **strong inference**：語境與插入值有交叉對照；不把 `{1}` 的所有全域語意升格為已解 |
| `#751` 的可重跑 fix | `translations/corrections.json`：`{6}沒什麼！我{1}可不會這樣就輸的！`；只在既有「我」後補 `{1}`，不另改中文語氣；`tools/talkdat_selftest.py` 的 APPLIED 集合已由 15 筆改為 16 筆 | **已定案（最小校訂）**；不是日文逐字重譯，也未宣稱畫面行寬 parity |

本輪使用 Docker-only Python 標準函式庫與 `demonwinter-go:latest` 讀取／比較資料，原始檔案唯讀掛載；
沒有改寫 `workplace/orig/`。`#751` 的文字修正已納入 `talkdat.py correct` 的 round-trip 與 mismatch
guard；M7 剩餘仍是 1,022 則逐句文意審查、CJK 畫面抽樣與目標平台驗收。

## 2026-08-09 — 事件 2／3 外交佇列資料流追查（事件 3 接入前）

| 斷言 | 證據 | 等級／實作邊界 |
|---|---|---|
| 事件 2 是非玩家勢力向玩家提出合作請求，事件字來源是玩家、參數低 byte 是請求方、high byte 是該勢力的入侵目標 | `KI.EXE.asm` `sub_12E33`：以目前 AI 勢力的 `+0x19` 入侵目標篩選，找到玩家後 `xchg si,di`、`al=2`、`bl=0xFF` 呼叫 `sub_12FB1`；`sub_12FB1` 將事件字高 byte 取自交換後的玩家勢力，參數保留 `DL`／`DH` | **confirmed**：佇列 payload 與觸發前置條件；remake 尚未產生此事件 |
| 事件 2 的接受分支會先走合作判定，成功時顯示訊息索引 `0x2F`，再由玩家對請求方宣戰 | `seg000:13220`：驗證事件來源／參數後呼叫 `sub_13712`；玩家路徑呼叫 `sub_138E6`／`sub_13C3D` 顯示 `0x2F`，再呼叫 `sub_135ED` 與 `sub_13526`（`DL` 為請求方） | **confirmed**：控制流與呼叫關係；`sub_13712` 的完整接受條件仍未接入，不能以現有 `CanRequestCooperation` 近似替代 |
| `sub_13771` 的外交／政治比較不是單純交友度門檻 | `seg000:13771`：`+0x2A` 有外交官時取其 `+0x13` 政治，否則由 `sub_137F5` 取對方勢力最高政治且 `+0x17=0` 的武將；再與本方君主／未出陣最高政治武將比較，平手消費 `sub_1ECE0` 的低位元 | **confirmed**：欄位與選人資料流；`sub_137D8`／`sub_13138` 的俘虜旗標結果與 UI 回應仍未接入 |
| `sub_13712`／`sub_136C4` 內含不同於玩家外交規則的算術 gate | `seg000:13712`：比較兩次 `sub_130CB` 交友度讀取、`+0x28×2+0x28` 好戰門檻，再以 `0x5A−value` 取半；`seg000:136C4`：若對方等於 `+0x19` 改反應碼，並以 `value−(+0x28+2)`、`0x1E−value` 取半；兩者最後乘 `0x3E8` 呼叫 `sub_137D8` | **confirmed**：算術與欄位來源；事件方向、`sub_137D8` 的完整接受／俘虜語意及回應碼仍需固定狀態 oracle |
| 事件 3 是勢力在無法同時承受多個標記鄰國時提出停戰，事件字高 byte 是發起勢力、參數低 byte 是對方、high byte 為 `0xFF` | `KI.EXE.asm` `sub_12E89`：以 `sub_13091` 計算勢力力量，逐一扣除標記鄰國力量，餘額不足時對剩餘標記鄰國以 `al=3`、`bl=0xFF` 呼叫 `sub_12FB1` | **confirmed**：觸發計算與 payload；remake 尚未產生此事件 |
| 事件 3 通過判定後會清除互指入侵目標、同步相關城市的記錄擁有者，並以雙方原始關係值較小者建立和平值 | `seg000:13262`：呼叫 `sub_136C4`；玩家路徑顯示 `0x2B`；後續 `sub_145F8` 清除互指 `+0x19`，`sub_14236` 同步城市 `OwnerRecorded`，`sub_13669` 以雙向關係最小值寫回和平值 | **confirmed**：成功後的狀態副作用；`sub_136C4` 的接受條件、玩家 UI 選項與延遲仍未解出 |
| 當時不把事件 2／3 接入 runtime 是有意的證據邊界，而非遺漏宣稱 | `sub_13712`／`sub_136C4` 仍會讀取雙向友好度、君主／外交官／武將政治、出場狀態與俘虜相關資料；既有 `internal/rules/diplomacy` 只涵蓋玩家勸說的近似規則 | **當時已記錄未完成**：事件 2／3 暫為 no-op／不生成；後續只接入事件 3 的已證實狀態副作用，事件 2 與事件 3 完整接受／UI 仍未完成 |

逆向輸入為 DOS/V `workplace/ida/dosv/KI.EXE.i64`、同目錄 `KI.EXE.asm` 與 `KI.EXE`；
`KI.EXE` SHA-256 為 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
工具為 `ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）；本文 `seg000:xxxxx` 是 IDA 線性位址，
未與檔案偏移混列。這次追查沒有修改反組譯資料庫、原始素材或 runtime 行為。

## 2026-08-09 — 事件 3 停戰狀態部分接入

| 斷言 | 證據 | 等級／實作邊界 |
|---|---|---|
| 事件 3 的代表政治值可由現有欄位重現 | `seg000:13771`：`Faction +0x2A` 有外交官時讀其 `General +0x13`；否則用 `sub_137F5` 掃對方勢力最高政治且 `+0x17=0`；本方君主未出陣時改取本方同條件最高政治者；比較結果為 `政治×2` 或 `(16−政治)×2`，平手取 `sub_1ECE0` bit 0 | **confirmed／已接入**：`World.diplomacyRepresentative`；代表不存在時 fail-closed |
| `sub_136C4` 的停戰金額與拒絕反應碼可接到狀態層 | `seg000:136C4`：讀 `Friendship[target][source]` 的 7-bit 值，扣 `Faction[target] +0x28 + 2` 後以 `0x1E−value` 加入政治結果，除二、乘 `0x3E8`；對方 `+0x19` 正在指向提出方時將 `AL` 變成 2，金額為零時為 0 | **confirmed／已接入**：`World.applyQueuedCeasefire`；不宣稱 `sub_13138` 的 AH 旗標語意 |
| 成功停戰的資金與俘虜副作用 | `seg000:135ED`：`AL=1` 時對 `SI` 勢力加金額、對 `DI` 勢力減金額；掃 General `+0x1C/+0x1D` 的雙向配對後呼叫 `sub_150D7` 釋放；`sub_150D7` 把原 Captor 還活著的武將歸回原勢力，否則設為在野 | **confirmed／已接入**：資金使用 `economy.ClampFunds`；UI／dirty bit 未新增 |
| 成功停戰的三個收尾效果 | `seg000:145F8` 清除彼此互指的 `+0x19`；`seg000:14236` 只在 `Owner`／`OwnerRecorded` 都屬兩方時同步記錄欄位；`seg000:13669` 取雙向交友度 7-bit 較小值並設和平位元 | **confirmed／已接入**：`TestQueuedEventHandlers` 驗證資金、侵攻目標、俘虜、城市記錄與雙向和平值 |
| 事件 3 的完整 runtime parity 已完成 | 目前只接入已證實的 `sub_13262` 狀態副作用；尚未接事件 3 的產生器、玩家訊息／接受 UI，也未把 `sub_13138` 的 `cmp ax,ax` 直接升格為俘虜旗標語意 | **未完成**：事件 2 仍不生成／不套用；事件 3 的完整接受／UI／長期劇本路徑仍待固定狀態 oracle |

本輪仍使用 DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm` 與 `KI.EXE`；
`KI.EXE` SHA-256 為 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
工具為 `ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）；`seg000:xxxxx` 是 IDA 線性位址，未與
EXE 檔案偏移混列。狀態層測試在 `wolong-go:20260809` Docker 映像內完成；沒有修改原始執行檔、
IDA 資料庫或原始素材。

## 2026-08-09 — 事件 2 合作狀態部分接入

| 斷言 | 證據 | 等級／實作邊界 |
|---|---|---|
| 事件 2 的 payload 可直接轉成三個勢力 | `seg000:13220`：事件字高 byte 經 `sub_1351A` 成為合作方；參數低 byte 成為侵攻方 `BX`；參數高 byte 成為被侵攻方 `DI`；最後以 `DL=參數低 byte` 呼叫 `sub_13526` | **confirmed／已接入**：`World.applyQueuedCooperation`；仍不生成事件 2 |
| `sub_13712` 的合作 gate 與金額 | `seg000:13712`：先呼叫 `sub_13771`；比較合作方→被侵攻方與合作方→侵攻方的原始交友度，並比較被侵攻方值與 `0x28 + 2×Faction +0x28`；金額為 `(0x5A−value)` 鉗 0..0x3C、除二、乘 `0x3E8`；結果 2 不套用，金額 0 為結果 0 | **confirmed／已接入**：`World.applyQueuedCooperation`；`sub_137D8` 的 AH 旗標仍不命名 |
| 合作成立的資金／俘虜副作用與宣戰方向 | `seg000:135ED` 在 `AL=1` 時對 `SI=合作方` 加金額、對 `DI=被侵攻方` 減金額，並釋放兩方俘虜；`sub_13220` 後續以 `DL=侵攻方` 呼叫 `sub_13526` | **confirmed／已接入**：合作方與侵攻方重用 `applyQueuedDeclaration`；玩家的 `+0x19` 保留不直接寫入，雙向交友度進入交戰值 |
| 事件 2 的完整 runtime parity 已完成 | `sub_13220` 對玩家仍會顯示訊息 `0x2F`；事件產生來自 `sub_12E33`；接受 UI／延遲／`sub_137D8` 回應旗標與長期多勢力路徑尚未接入 | **未完成**：`TestQueuedEventHandlers` 只驗證固定狀態的金額、宣戰、交戰值與俘虜副作用 |

本輪仍使用 DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm` 與 `KI.EXE`；
`KI.EXE` SHA-256 為 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
工具為 `ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）；`seg000:xxxxx` 是 IDA 線性位址，未與
EXE 檔案偏移混列。狀態層測試在 `wolong-go:20260809` Docker 映像內完成；沒有修改原始執行檔、
IDA 資料庫或原始素材。

## 2026-08-09 — 事件 8 遷都後軍團同步補接

| 斷言 | 證據 | 等級／實作邊界 |
|---|---|---|
| 事件 8 與首都失守的兩條遷都路徑都會在首都改變後呼叫 `sub_14502` | `KI.EXE.asm` `seg000:133FD` 與 `seg000:14DF0` 都在 `xchg ah,[...+3]`、確認新舊不同後呼叫 `sub_14502` | **confirmed**：`World.relocateCapital` 統一接入同步；事件 8 與據點易主共用 |
| `sub_14502` 的第一段只同步同勢力、存在旗標 ≥ `0x80`、`Home=舊首都` 的軍團 | `seg000:14502`：掃 `ds:2240h` 起 127 筆 × 64 B；比較 `+0x01`、`+0x00`、`+0x20`，命中後寫 `+0x20=AL` | **confirmed／已接入**：`syncCorpsAfterCapitalChange` |
| 第二段雖方向反直覺但可由暫存器定案 | `seg000:14502` 先以 `AL` 形成 `CX=新首都×8`，以 `AH` 形成 `DX=舊首都×8`；若 `+0x14==CX` 就寫 `+0x14=DX`，沒有寫 `+0x16/+0x18` | **confirmed／已接入**：只改 `TargetNode`，保留 `TargetX/TargetY`，不把反向寫入改成推測的「新首都」 |
| 非匹配軍團不會因遷都被重寫 | `seg000:14502` 的勢力／存在／`Home` 三重比較；`TestQueuedEventHandlers` 放入同勢力兩種目標與其他勢力軍團，驗證匹配與不匹配分支 | **已驗證**：事件 8 首都、Home、目標節點與 X/Y 保留值均通過 |

本輪只接入可直接回查的欄位寫入；原版 `or byte ptr [si],2` 的 dirty bit 沒有在 `Corps`
模型新增未證實用途，因為目前存檔／runtime 沒有等價的獨立旗標。輸入仍為上述 DOS/V
`KI.EXE` 與 IDA Pro 9.4 Docker 工具鏈，沒有修改原始執行檔或素材。

## 2026-08-09 — 事件 2／3 產生器接入

| 斷言 | 證據 | 等級／實作邊界 |
|---|---|---|
| `sub_12C52` 的高位元是交戰標記，而不是鄰接資料本身 | `seg000:12CA5`–`12CBF`：以 `xlat` 讀交友度原始值；低於 `0x80` 時對候選勢力記錄位址 `or di,8000h`，再進行選擇排序。`sub_12CDF` 只在 `seg000:12D12` 寫入乾淨的勢力記錄位址 | **confirmed**；`World.strategyCandidates` 的排序快照保留此標記語意 |
| 事件 2 的四道產生閘與 payload 已接入 | `seg000:12E33`：侵攻目標不得為 `0xFF` 或玩家；玩家須在 `sub_12C52` 清單；`Friendship[AI][玩家]` 原始值須 `≥0x80`，`Friendship[目標][玩家]` 須 `≥0xA3`；`al=2`、交換 `si/di` 後呼叫 `sub_12FB1`。runtime 為 `World.queueCooperationProposal`，Code 高 byte＝玩家、Param＝目標<<8|侵攻方 | **confirmed／已接入**：`TestStrategicAIDiplomacyEventGenerators` |
| 事件 3 的國力累減與排入範圍已接入 | `seg000:12E89`：玩家勢力跳過；`sub_13091` 取得自身國力後，只累減已標記交戰鄰居；第一次降至 `≤0` 時，從該筆往後掃描，排除第一筆清單項目，對每個標記鄰居以 `al=3`、`DH=0xFF`、`DL=對象` 呼叫 `sub_12FB1`。runtime 為 `World.queueCeasefireProposals` | **confirmed／已接入**：測試驗證「第二筆排入、第一筆不排入」邊界 |
| 月結呼叫順序已接到 live state | `seg000:12D58`：交友度漂移後依序呼叫 `sub_12E33`、`sub_12E89`、`sub_12EFB`；`World.runStrategicAI` 在每個勢力的漂移後產生事件 2／3，再處理既有宣戰佇列 | **已接入**；事件 2／3 的完整接受回應、玩家訊息／UI、`sub_137D8`／`sub_13138` 旗標語意與長期 oracle 仍未完成 |

逆向輸入為 DOS/V `workplace/ida/dosv/KI.EXE.i64`、同目錄 `KI.EXE.asm` 與 `KI.EXE`；
`KI.EXE` SHA-256 為 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
工具為 `ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）；上述 `seg000:xxxxx` 全是 IDA 線性位址，
沒有與檔案偏移混列。runtime 與測試在 `wolong-go:20260809` Docker 映像內完成；沒有修改
IDA 資料庫、原始素材或反組譯輸出。

## 2026-08-09 — 事件 9 俘虜釋放狀態接入與外交代表欄位校訂

| 斷言 | 證據 | 等級／實作邊界 |
|---|---|---|
| 事件 9 的事件字高 byte 是指定武將索引，不是勢力來源 | `seg000:13485`：以完整 `AX=Code` 右移三次後加 `0x4240`，形成 General 記錄位址，再呼叫 `sub_150D7`；低 byte 9 只用於 dispatch | **confirmed／已接入**：`World.dispatchQueuedEvent` 的 case 9 呼叫 `World.releaseGeneral(source)`；`TestQueuedEventReleaseGeneral` 驗證兩條回寫分支 |
| `sub_150D7` 的核心狀態寫入 | `seg000:150D7`：清 General `+0x17`（出陣）與 `+0x1D`（Captor）；原 Captor 勢力仍存在時把 `+0x1C` 回寫為該勢力，否則寫 `0xFF` 在野 | **confirmed／已接入**；原版訊息／dirty bit／完整事件互動仍未接入 |
| `sub_13771` 的已出陣君主需要平行軍團記錄存在 | `seg000:13771`：先讀 `ds:2240h` 平行 Corps 表的 `+0x00` 存在旗標，再依 General `+0x17` 決定使用君主或 fallback；remake 對 posted lord 採同一 fail-closed gate | **confirmed／已接入**：`World.diplomacyRepresentative`；`TestQueuedEventHandlers` 的 event 2／3 合成狀態提供平行 Corps 記錄 |

逆向輸入為 DOS/V `workplace/ida/dosv/KI.EXE.i64`、同目錄 `KI.EXE.asm` 與 `KI.EXE`；
`KI.EXE` SHA-256 為 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
工具為 `ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）；`seg000:xxxxx` 是 IDA 線性位址，未與
EXE 檔案偏移混列。狀態層實作／測試在 `wolong-go:20260809` Docker 映像內完成；沒有修改
IDA 資料庫、原始素材或反組譯輸出。

## 2026-08-09 — 事件 2／3 玩家外交三選一接縫

| 斷言 | 證據 | 等級／實作邊界 |
|---|---|---|
| 玩家涉入事件 2／3 時，事件處理先停在外交選擇，不立即套用狀態 | `seg000:13220`／`seg000:13262`：玩家分支分別進 `sub_138E6`／`sub_138C7`，兩者再進 `sub_13902` 三選項；runtime 為 `World.beginDiplomacy`，`Event.Diplomacy` 與 `World.PendingDiplomacy` 保存暫存狀態 | **confirmed 控制流／已接入**；`World.Tick` 與同小時後續財政處理在 pending 時早退 |
| 三列選項的文字與 response 方向 | DOS/V `TALK.DAT`（SHA-256 `08a22e09791d0a6ec2968e87d8655e12c91b45e00fae460b28593b35ff85e384`）#363／#376：`無條件同意／提供資金／拒絕`；`sub_13902` 的 response 0／1／reject 對應 `World.ResolveDiplomacy` 的 `DiplomacyAcceptFree`／`DiplomacyOfferFunds`／`DiplomacyReject` | **confirmed 選項／部分接入**；提供資金暫用已解出的預設金額，numeric input、Trust advice、完整訊息池仍未接 |
| GUI 正常模態接縫 | `cmd/wlgame/diplomacy.go`：事件 2／3 的三列選擇、Enter／1–3／方向鍵、ESC 拒絕；`TestQueuedDiplomacyChoice` 驗證兩種接受、pending 清除與時鐘不動 | **已接入／未宣稱完整原版 UI**；尚缺原版 TALK.DAT 訊息順序、建議文字與金額輸入對拍 |

逆向輸入為 DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm`、`KI.EXE` 與
`TALK.DAT`；`KI.EXE` SHA-256 為 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
工具為 `ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）；`seg000:xxxxx` 是 IDA 線性位址，未與
EXE 檔案偏移混列。runtime／測試在 `wolong-go:20260809` Docker 映像內完成；GUI 程式碼只在
工作樹寫入，沒有修改 IDA 資料庫、原始素材或反組譯輸出。

## 2026-08-09 — 事件 4／5 官員撥款產生器與狀態接入

| 斷言 | 證據 | 等級／實作邊界 |
|---|---|---|
| 事件 4 的高 byte 是據點編號，Param 是要求金額 | `KI.EXE.asm` `seg000:15715`：`di` 從據點表 `0x0840` 起每次加 `0x20`；扣回基址後左移三次形成 `city×0x100`，`mov al,4` 後呼叫 `sub_12FBF`；`DX` 由三個 gap 算式得到 | **confirmed／已接入**：`World.queueFundingRequests`；`TestFundingRequestGenerators` |
| 事件 4 的上昇值要用存檔 byte `+0x10` | `seg000:15715` 直接比較 `[di+10h]` 與 `0xB4`；`docs/formats/08` 已證實該欄位是「上昇值＋100」，runtime 的 `City.Growth` 已扣除 100 | **confirmed**；接入時使用 `Growth+100`，避免偏移錯誤 |
| 事件 5 的高 byte 是勢力編號，Param 是外交官要求金額 | `seg000:1578F`：`di` 從勢力表 `0` 起每次加 `0x40`；左移兩次形成 `faction×0x100`，`mov al,5` 後呼叫 `sub_12FBF` | **confirmed／已接入**：`World.queueFundingRequests`；`TestFundingRequestGenerators` |
| 事件 5 的要價取兩向友好度原始 byte 較小者 | `seg000:1578F`：第一次 `sub_130CB` 暫存 `DL`，交換 `si/di` 再讀第二向，`cmp al,dl` 後保留較小原始 byte；`0x80` 以上基準 100，否則 125，低七位相減後乘 200 | **confirmed／已接入**：`diplomacy.Demand` 改為 raw min；外交旗標混合案例有回歸測試 |
| 事件 4／5 處理時會重新確認官員仍被派駐 | `seg000:132A9` 先讀據點 `+0x19`，`seg000:132E9` 先讀勢力 `+0x2A`，`0xFF` 直接返回；兩支再把 General 指標交給 `sub_139E8` | **confirmed／已接入**：`World.beginFunding`／`fundingOfficer` fail-closed |
| 撥款的已證實狀態效果 | `seg000:139E8`：非零初始金額低於 `0x1F4` 拉到 500；`sub_17C6E` 上限 `0x7530`；response 2 跳過；否則 `shl ax,1` 後存 General `+0x1A` 高位，再以 `sub_1563B` 從玩家資金扣款 | **confirmed／已接入**：`PendingFunding`／`SetFundingAmount`／`ResolveFunding`；`TestQueuedFundingChoice` |
| 玩家 GUI 的狀態接縫 | `cmd/wlgame/funding.go`：三列選項、Enter／1–3／方向鍵、ESC、調整金額；`World.Tick`／`hourly` 在 pending 時凍結 | **已接入／未宣稱原版畫面 parity**：TALK.DAT 訊息順序與逐位數字輸入仍待對拍 |

逆向輸入為 DOS/V `workplace/ida/dosv/KI.EXE.i64`／`KI.EXE.asm`／`KI.EXE`；
`KI.EXE` SHA-256 為 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
工具為 `ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）；本文 `seg000:xxxxx` 是 IDA 線性位址，
未與檔案偏移混列。runtime／測試在 `wolong-go:20260809` Docker 映像內完成；沒有修改
IDA 資料庫、原始素材或反組譯輸出。

## 2026-08-09 — 事件 6／7 外交官回報狀態接入

| 斷言 | 證據 | 等級／實作邊界 |
|---|---|---|
| 事件 6 的來源是回報停戰的對方勢力 | `seg000:13327`：事件字高經 `sub_1351A` 成為 `SI`，讀 `SI+0x2A` 的外交官；以 `DI=word_10CFD` 呼叫 `sub_136C4`，後續 `sub_135ED` 由 `SI` 收款、`DI` 付款 | **confirmed／已接入**：`dispatchQueuedEvent` case 6 → `applyQueuedCeasefireReport` → `applyQueuedCeasefire(Player, other)` |
| 事件 7 的 payload 是協力方與第三方 | `seg000:13388`：事件字高成為協力方；`DL` 經 `sub_1351A` 成為第三方；`sub_13712` 的 `SI=協力方`、`BX=第三方`、`DI=玩家`；完成後 `sub_13526` 以協力方對第三方設目標 | **confirmed／已接入**：`dispatchQueuedEvent` case 7 → `applyQueuedCooperationReport` → `finishQueuedCooperation(ally, invader, Player, …)` |
| 事件 6／7 都要求回報勢力仍有外交官 | `seg000:13327`／`13388`：讀 `Faction +0x2A`，`0xFF` 直接返回；不是依事件入列時的舊指標 | **confirmed／已接入**：處理時重新驗證，缺少外交官 fail-closed |
| 事件 6／7 的核心狀態效果可重用既有收尾 | `sub_135ED`、`sub_145F8`、`sub_14236`、`sub_13669` 與 `sub_13526` 的共用呼叫鏈；固定測試驗證付款方向、停戰、第三方宣戰與交戰值 | **confirmed／已接入**：`TestQueuedDiplomacyReportHandlers` |
| 事件 6／7 的完整玩家進言 parity 已完成 | `sub_164F1`／`sub_16623` 仍包含原版選單、反應訊息與 `sub_1300E`／`sub_1301C` 的產生路徑；目前 `cmd/wlgame/advise.go` 尚未把同一條流程接到這兩個 payload | **未完成**：只完成 queue 處理端，不宣稱原版訊息／玩家產生路徑 |

## 2026-08-09 — 勘誤：玩家進言 producer 接回

上一節最後一列是當時的狀態快照；本節保留該歷史結論並記錄後續變更，不覆寫
原證據。

| 資料流 | 證據／對應 | 最新狀態 |
|---|---|---|
| 敵對進言同意 | `seg000:16405` 在 `sub_16475` 後直接呼叫 `sub_13526`；`World.ApplyPlayerHostility`；`TestPlayerDiplomacyProducers` | **confirmed／已接入**：玩家不走事件 1 延遲 producer，並在套用前重驗證和平與存活 |
| 停戰進言同意 | `seg000:164F1` 的 `AL=6`、`DX=0`、`BL=14h` → `sub_1300E`／`sub_1301C`；`World.QueuePlayerCeasefire` | **confirmed／已接入**：事件 Code 高 byte 為回報方，從 `eventCursor+0x14*4` 搜尋完整 256 格，重複事件拒絕 |
| 協力進言同意 | `seg000:16623` 的 Code 高 byte＝協力方、Param 高低 byte 同為第三方、`BL=14h`；`World.QueuePlayerCooperation` | **confirmed／已接入**：兩段選取與 payload 已接入，重複事件／存活／外交官／交戰條件 fail-closed |
| 玩家第一反應 | `sub_16475`／`sub_16577`／`sub_166D9` 與 `persuasion.FirstReaction` | **confirmed／已接入**：UI 分流直接同意、拒絕、已交戰／同一家與說服入口；TALK.DAT 原文／逐頁排版仍未完成 |

新增測試與程式均以 Docker `wolong-go:20260809`、Go `gofmt`／`go test` 執行；
輸入反組譯證據仍為 `KI.EXE.asm` SHA-256
`FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`，IDA Pro 9.4
線性位址（`seg000:16405`、`164F1`、`16623`、`1301C`）。

## 2026-08-09 — 事件 10–12 dispatch 與災害 runtime marker

| 資料流 | 證據／對應 | 最新狀態 |
|---|---|---|
| 事件 10 | `seg000:13496`：`AL=AH`、`AH=FF`、`CX=DX` 後只呼叫 `sub_18810` | **confirmed／訊息-only 邊界**：state 只取出事件，不假設持久寫入 |
| 事件 11 | `seg000:134A6` → `sub_1237E`：`rand&0x0F+0x18`，按暴風雨中心距離寫城市 `+0x15` | **confirmed／已接入**：`World.applyQueuedStormMarker` 寫 runtime marker／level；`World.Tick` 月結排入事件 11 |
| 事件 12 | `seg000:134B1`：高 byte 0 清除；高 byte 1／2 以 `sub_123FF` 建火災／暴動物件，兩次亂數寫強度／6..13 延遲，再用 `sub_1301C` 排清除事件 | **confirmed／部分已接入**：`World.applyQueuedDisasterMarker` 保留高 byte／runtime city pointer、marker、延遲清除與完整 256 格搜尋 |
| 持久傷害 | 目前 handler 片段只直接寫 `City+0x15`／物件表；未找到可證實的生產力／上昇值／城兵寫入端 | **未知／刻意未接入**：不得把 `Event.Disaster` 回報 map 當成傷害 parity |

輸入仍為同一份 `KI.EXE.asm` SHA-256
`FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`，IDA Pro 9.4
線性位址（`seg000:13496`、`134A6`、`134B1`、`1237E`、`123FF`、`12438`、`1301C`）。
新增 `TestQueuedDisasterAnimationHandlers` 以 Docker `wolong-go:20260809` 執行。

逆向輸入為 DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm` 與 `KI.EXE`；
`KI.EXE` SHA-256 為 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
工具為 `ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）；本文 `seg000:xxxxx` 是 IDA
線性位址，未與檔案偏移混列。狀態層測試在 `wolong-go:20260809` Docker 映像內完成；
沒有修改 IDA 資料庫、原始素材或反組譯輸出。

## 2026-08-09 — 事件 9 釋放結果的通用通知觀測

| 斷言 | 證據 | 最新狀態 |
|---|---|---|
| 事件 9 的高 byte 指定武將，處理成功後才應報告 | `seg000:13485` 以事件字高 byte 形成 General 記錄位址，再呼叫 `sub_150D7`；`World.dispatchQueuedEvent` case 9 現在只在 `releaseGeneral` 成功時追加 `Event.ReleasedGenerals` | **confirmed／已接入**：`TestQueuedEventReleaseGeneral` 驗證兩條狀態分支與成功回報 |
| runtime 需要把釋放結果交給玩家 | `cmd/wlgame/main.go` 消費 `Event.ReleasedGenerals`，以武將名顯示通用「已釋放」訊息 | **已接入／觀測接縫**：不冒充原版訊息排版 |
| 原版通知內容與完整流程 | `seg000:150D7` 的勢力通知及 `TALK.DAT` 原句「被敵軍所擒的{1}大人回來了」；句型已由 `cmd/wlgame` 的 `talkMessage(0x25)` 取用並代入武將名，原版訊息選擇／逐頁顯示尚未接回 | **部分已接入／未完成**：仍需原版／remake 同狀態畫面與長期事件 oracle |

逆向輸入為 `KI.EXE.asm`，SHA-256 為
`FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`；工具為
`ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4），本文位址均為 IDA 線性位址，未與檔案偏移混列。
新增欄位、`TALK.DAT` 句型取用、GUI 與測試在 `wolong-go:20260809` Docker 內完成；沒有修改 IDA 資料庫、原始素材或反組譯輸出。

## 2026-08-09 — 戰術 raw `PlaneHigh` 寫入端與分支條件勘誤

### 輸入與定位契約

- 輸入：DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm`、`KI.EXE`。
- 雜湊：`KI.EXE.asm` SHA-256 `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`；`KI.EXE` SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
- 工具：`ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）；本文所有 `seg000:xxxxx` 均為 IDA 線性位址，不與檔案偏移混用。函式匯出為 `workplace/ida/dosv/func-sub_1A85B.txt`、`func-sub_1AB7C.txt`、`func-sub_1AB9C.txt`、`func-sub_1ABB2.txt`、`func-sub_1AC55.txt`、`func-sub_1ACA4.txt`、`func-sub_1AD7F.txt`、`func-sub_1B0D3.txt`、`func-sub_1B116.txt`、`func-sub_1B240.txt`、`func-sub_1B4EA.txt`、`func-sub_1B732.txt`、`func-sub_1B8AA.txt`、`func-sub_1BA2E.txt`。
- 原始 `.i64`、原始資產與反組譯輸入保持唯讀；本輪只新增非破壞性函式匯出與文件紀錄。

### 證據與實作邊界

| 斷言 | 原始證據 | 等級／接入狀態 |
|---|---|---|
| 兵 `+0x1E` 的生命週期是地面 0／上移 0x10／回地面 0，換位時隨兵交換 | `sub_1B4EA` 寫 0；`sub_1B0D3` 寫 `0x10`；`sub_1B116` 清 0；`sub_1B732` 交換 `[si+1Eh]`；`sub_1B240` 將它複製至 `+0x1F` 再組入格索引 | **confirmed／已接入**：`Soldier.PlaneHigh`，`Place`／移動／增援／交換同步；`Climbing` 只作舊夾具相容 |
| `+0x00 bit 1` 是高地面旗標 | `sub_1B240` 由目前圖塊值判斷堆疊高度，`>=4` 時 `or byte ptr [si],2` | **confirmed／已接入**：`Field.HighTerrain`／`Soldier.HighTerrain` |
| 鎖敵 64 格懲罰的分支 | `sub_1A85B` 讀目前／候選 `+0x1E`；目前較高時看候選 `+0x00 bit 1`；目前不高於候選時比較目前兵種 `+0x04 <= 0x12` 與候選平面是否非 0，成立才 `add al,40h` | **confirmed／已接入**：`targetPlanePenalty`；取代先前 `Z > 0 && !CanClimb` 的近似 |
| CH=0x20 步兵特殊投射物的進入條件 | `sub_1AB7C`／`1AB9C`／`1ABB2` 依兵種分派；`sub_1AC55` 直接比較 raw `+0x1E`；`sub_1ACA4` 回傳 `max(abs(dx),abs(dy))`；`sub_1AD7F` 建立 `CH=0x20` | **confirmed／部分已接入**：`specialAttackAvailable` 已用 raw `PlaneHigh`、步兵、max-axis ≤2；完整 `sub_1BA2E` 投射物動畫／畫面對拍仍未完成 |

### 驗證

- Docker 映像：`wolong-go:20260809`，目前使用者 UID/GID，網路關閉、有界 CPU／記憶體／PID。
- 新增並通過：`TestRawPlaneHighAndTerrainFlag`、`TestSpecialProjectileUsesPlaneHighAndMaxAxisDistance`、`TestLockOnPlanePenaltyMatchesRawBranches`；原有 `TestClimbingInfantryCanUseSpecialProjectile` 也在移除 `Z>0` 偽 fallback 後通過。
- 本輪仍不得把「raw 欄位與分支已確認」寫成「完整戰術 parity」；未完成項為投射物逐幀運動／命中動畫、GUI 戰後訊息、有效時序的原版／remake 同狀態對拍與 M7／M8 全範圍。

## 2026-08-09 — 投射物 `sub_1B941` 更新順序與 raw 狀態接入

### 輸入／工具／位址基準

- 輸入：DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm`、`KI.EXE`；`KI.EXE.asm` SHA-256 `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`，`KI.EXE` SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
- 工具：`ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4），`seg000:xxxxx` 均為 IDA 線性位址；輸入 `.i64` 唯讀掛載後在容器內複製副本匯出，沒有修改原始資料庫。Go／UI 驗證使用 `wolong-go:20260809`。
- 直接證據：`sub_1B941` `seg000:1B941–1B97E`、`sub_1B97E` `1B97E–1BA2E`、`sub_1BA2E` `1BA2E–1BAB7`、`sub_1BAB7` `1BAB7–1BB10`；可回查 `KI.EXE.asm` 與隔離匯出，未把檔案偏移混進位址欄。

### 證據表

| 斷言 | 證據 | 等級／接入狀態 |
|---|---|---|
| 每幀先查目前格、再移動、後畫／更新 | `sub_1B941` 依序呼叫 `sub_1B97E`、`sub_1BA2E`、`sub_1BAB7` | **confirmed／已接入**：`Battle.stepProjectiles` |
| raw 方向與三維格索引 | `sub_1BA2E` 讀 `+0x05`／`+0x10`；西東改 X ±1、北南改 Y ±1／`0x40`；`+0x10 = Z×0x1000 + Y×0x40 + X` | **confirmed／已接入**：`direction`／`gridIndex`／`previousGridIndex` |
| 固定點高度與速度 | `sub_1BA2E` 以 `+0x14` 夾在 `±0x100` 加到 `+0x0A`；負值清除；成功後 `−0x14`；`+0x0B` 的 Z 上限 5 | **confirmed／已接入**：`heightFP`／`velocityFP`，普通箭初始速度公式與 `sub_1ECE0` RNG 已接入；圖形／動畫仍未完成 |
| 命中與障礙生命週期 | `sub_1B97E` 目前格空則繼續、`>=0x80` 清除、敵我編號區間避免誤傷；`sub_1BA2E` 移動後再查障礙 | **confirmed／部分接入**：以 `soldierAt`／`Field.solid` 取代原始佔用圖；完整原版 occupancy map 尚未做同狀態對拍 |
| 高度威力 | `sub_1BAB7` 新 Z > 舊 Z 時 `power − power/4`；下降時 `power + power/4 + 1` | **confirmed／已接入**：測試固定上升、下降與 CH=0x20 步兵傷害 |
| 戰場視覺觀測 | 原版繪製呼叫在 `sub_1BAB7`，投射物圖形資料尚未完整解出；PC-98 UI 參考只提供 640×400／戰場優先層次，不是原版圖形 oracle | **部分已接入／明示替代**：`Battle.Projectiles()`＋248×192 戰場內側別／特殊標記；不宣稱像素 parity |

### 驗證與剩餘邊界

- 新增通過：`TestProjectileRawDirectionGridAndHeightPower`、`TestProjectileChecksCurrentCellBeforeMoving`、`TestProjectileStopsAtSolidLayerAfterMoving`；全量 `go vet ./...`、`go test ./... -count=1`、`go build -o /tmp/wlgame-projectile-check ./cmd/wlgame` 通過。
- PC-98 UI 技能只用於區域層次與原生解析度決策：戰場保留空間推理、命令／狀態區不被標記侵入；不複製 Golden Box 的邊框、圖像或文字。普通箭 `sub_1AD2D` 初始速度公式與 `sub_1ECE0`／`sub_1EC82` RNG 已由 raw 與單測接入；完整 `BATTLE.SCH` 投射物圖形／動畫、原版／remake 同狀態對拍仍未完成。

## 2026-08-09 — M8 跨平台發行與封裝 Linux smoke

| 斷言 | 證據 | 等級／邊界 |
|---|---|---|
| 純 Go 交叉產物不是同一平台假成功 | Docker 內分別以 `GOOS/GOARCH` 建 `linux/amd64`、`linux/arm64`、`windows/amd64`、`darwin/amd64`、`darwin/arm64`；檔頭分別核對 ELF `7f454c46`、PE `4d5a`、Mach-O `cffaedfe` | **confirmed／已接入發行目錄**：純 Go `wlsim`／`wlshot`；Windows 另有 `wlgame`／`wlview` |
| Linux Ebiten 本體可由封裝目錄啟動 | `dist/linux-amd64/wlgame` 由 `wolong-go:20260809`、Xvfb `:101`、唯讀 `workplace/orig/dosv`／`workplace/eten` 執行 `-shot ... -shot-frames 120`，輸出 74,269 bytes PNG | **confirmed／已通過 packaged smoke**；Xvfb 背景程序以 trap 結束 |
| 發行包不帶原版資產 | `python3 tools/denylist.py dist`：掃描 18 個檔、比對 120 個原版檔，無命中 | **confirmed／已接入 gate** |
| 目標平台完整 runtime parity | 本輪只在 Linux Docker 執行封裝 smoke；Windows／macOS GUI 尚未在目標平台實機啟動 | **未知／未完成**：不可由檔頭或交叉編譯宣稱完成 |

輸入／工具：專案工作樹與既有 `wolong-go:20260809`，所有建置／Python／Xvfb 均在 Docker；沒有修改原始素材。M8 產物／Linux smoke 已更新 `WORKLIST.md`、`REMAKE-PLAN.md`、`VERIFICATION-MATRIX.md`，而 `HANDOFF.md` 已刪除；M7、Windows／macOS runtime 與原版／remake 同狀態對拍仍保留在 Remaining。

## 2026-08-09 — 事件 11／12／13 TALK 通知接縫

### 證據輸入與定位契約

- 輸入：DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm`、`KI.EXE`、唯讀 `workplace/orig/dosv/TALK.DAT`。
- 雜湊：`KI.EXE.asm` SHA-256 `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`；`KI.EXE` SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
- 工具與位址：`ida-pro-9.4-ver2:uidfix-v1`／IDA Pro 9.4；下列 `000xxxxx` 均為 IDA 線性位址，不是檔案偏移。原始 `.i64` 唯讀掛載後只在容器內臨時副本匯出，沒有修改資料庫或原始素材。實作／測試使用 `wolong-go:20260809`。

### 直接證據與接入

| 斷言 | 直接證據 | 等級／狀態 |
|---|---|---|
| 玩家城市暴風雨通知使用 TALK #70 | `sub_1237E`（`0001237E`）呼叫 `sub_18810`，`CX=0x46`；呼叫前以 `DI` 指向城市記錄，供 formatter 的 `\\2` 取值 | **confirmed／已接入**：`TalkNotice{Index:0x46, City:...}` |
| 事件 12 的火災／暴動通知分別使用 #71／#72 | `sub_134B1`（`000134B1`）高來源分支依事件來源形成 `0x46`、`0x47`、`0x48`；`TALK.DAT` #70／#71／#72 分別為暴風雨／大火／暴動 | **confirmed／已接入**：事件 11／12 玩家城市通知 |
| 事件 13 的玩家通知使用 TALK #51 | `sub_13507`（`00013507`）以 `AL=0x93`、`CX=0x33` 呼叫 `sub_18810`；`TALK.DAT` #51 為主公盛怒訊息 | **confirmed／已接入**：事件 13 `TalkNotice{Index:0x33}` |
| 事件 9 釋放武將可沿用同一呈現佇列 | `sub_15940` 的 #65／#66 來源與現有事件 9 狀態；TALK #37（索引 `0x25`）為釋放武將回來訊息 | **strong inference／已接入**：只共用 remake modal 接縫，未宣稱原版完整流程 parity |
| 事件 10 已有可安全接入的 producer | 本輪檢查 `sub_12FBF` 呼叫清單與 formatter 路徑，未找到可確認的事件 10 producer；`sub_13496` 為 message-only formatter 邊界，尚不能反推事件來源 | **unknown／保留未完成** |

### Remake 行為與驗證

- `internal/state.Event.TalkNotices` 儲存原始 TALK index、city/general ID，不把翻譯文字或 GUI 狀態污染規則層。
- `cmd/wlgame/messages.go` 保留 TALK.DAT 行邊界、Big5 解碼與 `{1}`／`{2}` 代換；modal 固定在 640×400 內容座標的命令／狀態區外，以 Enter／Space 關閉並讓 `timeRuns()` 停止世界時間。這個空間層次參考 `research-pc98-golden-box-ui` 的原生解析度／戰場與命令區分界原則，不複製其圖像或文案。
- `TestQueuedTalkNotices` 驗證事件 11／12／13 產生的 TALK index 與玩家城市／一般通知參數；Docker `go test ./internal/state ./cmd/wlgame -count=1` 通過。
- 實作勘誤：第一次 `-open-message` 截圖沒有 modal，追查後確認 `TALK.DAT` parser 保存的是 ASCII marker `0x31`／`0x32`（`'1'`／`'2'`），而 UI 暫時用了數值鍵 `1`／`2`；這使既有 fail-closed 邏輯正確拒絕整則訊息。已改用 ASCII marker，重新建置後 `docs/images/wlgame-event-modal.png` 顯示「許昌發生了暴風雨。」；30 幀後日期仍為 196/4/1。這是實作修正，不改變原版證據或推論等級。
- 明確邊界：這是通知資料與呈現接縫，不是原版肖像、逐字翻頁、完整事件 6／7 TALK 觸發、事件 9 完整流程、事件 10、災害持續傷害／物件動畫或原版／remake 同狀態畫面對拍。

## 2026-08-09 — 事件 6／7 主要 TALK 回報接入

### 證據輸入與位址

- 輸入：DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm`、`KI.EXE`、唯讀 `workplace/orig/dosv/TALK.DAT`。
- 雜湊：`KI.EXE.i64` SHA-256 `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`；`KI.EXE.asm` SHA-256 `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`；`KI.EXE` SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
- 工具：`ida-pro-9.4-ver2:uidfix-v1`／IDA Pro 9.4；本節 `000xxxxx` 均為 IDA 線性位址。`.i64` 唯讀掛載後複製到容器 `/tmp`，只讀查詢 `CS:08A4` table，臨時副本與 log 已清理。TALK 文字由唯讀 DOS/V `TALK.DAT` 以 `tools/talkdat.py`／cp950 解析；runtime／測試使用 `wolong-go:20260809`。

### 直接證據

| 斷言 | 直接證據 | 等級／狀態 |
|---|---|---|
| 事件 6 先報告外交官／回報勢力 | `sub_13327`（`00013327`）在檢查 `SI+0x2A` 後把外交官 ID 以高 byte `FF`、勢力記錄指標 `SI` 依序放入 formatter stack，`CX=0x39` 呼叫 `sub_18810` | **confirmed／已接入**：`TalkNotice{Index:0x39, Faction:other, General:diplomat}` |
| 事件 7 共用同一則 #57 | `sub_13388`（`00013388`）交換後以協力方 `SI` 檢查 `+0x2A`，同樣 `CX=0x39` 呼叫 `sub_18810` | **confirmed／已接入**：`TalkNotice{Index:0x39, Faction:ally, General:diplomat}` |
| 事件 6 的結果索引與金額 | `sub_136C4`（`000136C4`）成功路徑把計算後 `DX` 乘 `0x3E8` 的值保留，`AL=CL` 回傳 0／1／2；`sub_13327` 以 `CX=0x2B` 呼叫 `sub_13C3D`，其 `CX += AL`，所以是 #43／#44／#45 | **confirmed／主要接入**：`Amount` 僅在 response 1 的 #44 保存；#45 不套用狀態 |
| 事件 7 的結果索引與金額 | `sub_13712`（`00013712`）同樣以 `DX` 保留乘 `0x3E8` 金額與 `AL=CL` response；`sub_13388` 以 `CX=0x2F` 呼叫 `sub_13C3D`，所以是 #47／#48／#49 | **confirmed／主要接入**：`Amount` 僅在 response 1 的 #48 保存；#49 不套用狀態 |
| `\\3` 的來源 | `sub_13C3D`（`00013C3D`）將相關勢力記錄指標放入 formatter stack；`CS:08A4` table 的 marker `\\3` handler（`00010904`）讀該記錄的 `+1` 勢力 ID，再定位君主／勢力名稱。`TALK.DAT` #43–#49 與 #57 的文字 marker 交叉吻合 | **confirmed／已接入**：`TalkNotice.Faction` → `World.LordName` |
| `\\7` 的來源 | `CS:08A4` table 的 marker `\\7` handler（`00010984`）讀下一個 stack word，經 `sub_1062F`（`0001062F`）逐位除 10 的十進位繪製；`sub_13C3D` stack 中該 word 是 `DX` | **confirmed／語意接入**：`TalkNotice.Amount` → ASCII 十進位；原版欄寬／glyph 位置仍未 parity |

### TALK 索引對照

`TALK.DAT` 解析結果：#43（`0x2B`）「與{3}停戰交涉的結果，無條件地達成了。」、#44
（`0x2C`）含 `{3}`／`{7}` 的金額結果、#45（`0x2D`）停戰破裂；#47（`0x2F`）、#48
（`0x30`）、#49（`0x31`）是對應合作結果；#57（`0x39`）是外交官／勢力報告；#58
（`0x3A`）是「敵方的君主已不在了。」。原始行邊界與 Big5 bytes 沒有寫入 state，
由 `messages.go` 在 modal 端展開；ASCII marker 修正與 #70 smoke 的勘誤見上一節。

### Remake 接入與邊界

- `TalkNotice` 新增 `Faction`／`Amount`，保留 -1 未使用契約；事件 6／7 dispatch 以單次
  `ceasefireTerms`／`cooperationTerms` 取得 response／amount，再同一份值產生通知與套用
  狀態，避免 tie RNG 被重抽。
- `TestQueuedDiplomacyReportTalkNotices` 驗證事件 6 的 #57→#44（14,000）與事件 7 的
  #57→#48（15,000）；原有 `TestQueuedDiplomacyReportHandlers` 仍驗證付款／停戰／宣戰。
- 仍未完成：`sub_13C3D` 在 AH／AL 特殊狀態下追加的次要 `CX+0x1D`／`sub_13DC9` 流程、
  事件 2／3 完整玩家回應、原版逐頁訊息／數值欄寬與長期正常劇本 oracle。主要 TALK index
與 marker 已接入，不等於完整外交訊息 parity。

## 2026-08-09 — 勘誤：事件 11／12 災害 marker 的持久效果

### 證據輸入與定位契約

- 輸入：唯讀 DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm`、`KI.EXE`；`.i64` 只在容器內複製副本查詢，沒有修改原始資料庫或原始素材。
- 雜湊：`KI.EXE.i64` SHA-256 `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`；`KI.EXE.asm` SHA-256 `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`；`KI.EXE` SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
- 工具：`ida-pro-9.4-ver2:uidfix-v1`／IDA Pro 9.4；本節所有 `000xxxxx` 均為 IDA 線性位址，不是檔案偏移。runtime／測試使用 `wolong-go:20260809`。

### 勘誤前後與直接證據

先前的事件 10–12 dispatch 段只確認 handler 直接寫 `+0x15`，因此把持久傷害列成未知；本輪沿主迴圈的**取址端**補查後，舊結論不再成立，保留該歷史紀錄但以本節為後續狀態：

| 斷言 | 直接證據 | 等級／接入狀態 |
|---|---|---|
| marker 會進入持久效果 | `sub_13EFD`（`00013EFD`）於 `00013F5A` 在 `sub_14194` 後無條件呼叫 `sub_14269`（`00014269`） | **confirmed／已接入**：`tickCity` → `applyCityDisasterEffect` |
| 防災護盾與不足分支 | `sub_14269` `00014269`–`000142A9`：`[si+851h]` 先減 `+0x15`；不足時寫 0，差額扣 `[si+850h]` | **confirmed／已接入** |
| 生產力損失 | `sub_14269` 的 `mul byte ptr [si+84Fh]`、右移兩次、`sub [si+84Eh], ax` | **confirmed／已接入**：`(deficit × (Production >> 8)) >> 2`，保留 u16 減法 |
| 城兵損失 | `sub_14269` `shr al, 1` 後 `sub [si+853h], al`，不足寫 0 | **confirmed／已接入** |
| marker 清除時機 | `sub_134B1` 高 byte 0 寫 `[si+15h]=0`；事件 12 高 byte 1／2 另排 6..13 格清除事件 | **confirmed／既有接入** |

### Remake 與驗證

- `internal/state/state.go` 在 governor 更新後呼叫 `applyCityDisasterEffect`；它只以已證實的 `disasterMarkerLevels` 對應 runtime `+0x15`，不把 marker 陣列序列化。
- 新增 `TestDisasterMarkerAppliesRawPersistentEffects`：固定 marker 7、防災 3，驗證差額 4 對上昇值、生產力與城兵的原版結果；另驗證防災值足夠時只扣護盾，且 marker 不由 `sub_14269` 自動清除。
- 針對性 state 測試與後續完整 Go／文件／發行驗證均在 Docker 執行；事件 10 formatter／producer、物件動畫、完整 TALK 翻頁、原版／remake 同狀態對拍仍是未完成邊界。
## 2026-08-09 — `sub_17C6E` 數值編輯核心接入

### 輸入與工具識別

- 輸入：唯讀 DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm`、`KI.EXE`；本段引用的 `.i64` SHA-256 為 `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`，`.asm` SHA-256 為 `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`，binary SHA-256 為 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
- 工具：`ida-pro-9.4-ver2:uidfix-v1`／IDA Pro 9.4；下列 `000xxxxx` 均為 IDA 線性位址，不是檔案偏移。runtime／測試使用 `wolong-go:20260809` Docker；原始輸入唯讀，未修改 `.i64`。

### 已證實資料流

| 原始定位 | 直接行為 | 推論等級 |
|---|---|---|
| `sub_17C6E` `00017C6E` | 呼叫端把目前值、初值畫面、上限傳入；事件 2／3 與事件 4／5 的上限皆為 `0x7530` | 已證實 |
| `sub_17DA5` `00017DA5` | `SI = SI × 10 + AL`；16 位元溢位後再與 `[BP+0]` 上限比較，結果不超過上限 | 已證實 |
| `sub_17DC3` `00017DC3` | `SI = SI × 100`，溢位／超過上限時鉗到上限 | 已證實 |
| `sub_17DDD` `00017DDD` | `SI = SI ÷ 10` | 已證實 |
| `sub_17DEC` `00017DEC` | `SI = [BP+0]`，還原呼叫前的初值 | 已證實 |
| `sub_17DF1` `00017DF1` | `SI = 0`，清零後繼續輸入迴圈 | 已證實 |
| `sub_17DEA` `00017DEA` | `STC`；由 `sub_17C6E` 把它轉成離開操作迴圈，保留目前 `SI` | 已證實 |

### Remake 接入與邊界

- 新增 `internal/state.AmountEdit`，以 `editAmount` 集中實作追加數字、追加 `00`、退位、**設成上限**、清零與結束輸入；`World.EditDiplomacyOfferAmount`／`World.EditFundingAmount` 分別接入事件 2／3、4／5。（2026-08-24 訂正：那一鍵是「最大」不是「還原初值」，`sub_17C6E` 收到的是上限——`docs/spec/78`。）
- `cmd/wlgame` 指定金額列現在接受跨平台數字鍵、退格、Insert、Delete、Home；原有方向鍵仍保留為 remake 輔助控制。`TestRawAmountEditorSemantics` 驗證追加、`00`、退位、還原、清零、非法數字 fail-closed 與 `30,000` 上限。
- 本段沒有把 PC-98 掃描碼、原版數字視窗字型／欄寬／游標或 TALK.DAT 順序外推成已完成；事件 2／3、4／5 的原版畫面與完整訊息流程仍是未完成邊界。此前把「numeric input 尚未接入」寫在現況中的說法，現在由本段修正為「核心語意與跨平台輸入已接，原版呈現仍未知」。
## 2026-08-09 — 事件 4／5 前置 TALK #56／#57 接入

### 直接證據

- `sub_132A9`（IDA 線性位址 `000132A9`）在呼叫 `sub_139E8` 前，以 `CX=0x38` 呼叫 `sub_18810`；`DI`／堆疊參數指向同一筆城市／內政官資料。DOS/V `TALK.DAT` 的 #56 是「{2}內政官的{1}大人，前來報告現狀。」。
- `sub_132E9`（`000132E9`）在呼叫 `sub_139E8` 前，以 `CX=0x39` 呼叫 `sub_18810`；`TALK.DAT` 的 #57 是「駐{3}勢力的外交官{1}大人前來報告。」。
- `sub_139E8` 後續仍會依初始要求、玩家選項與 `sub_17C6E` 結果進入多個訊息／資金分支；本段只接前置報告，沒有把後續索引猜成完成。

### Remake 接入

- `dispatchQueuedEvent` case 4／5 在 `beginFunding` 成功後分別 append `TalkNotice{Index:0x38, City, General}` 與 `TalkNotice{Index:0x39, Faction, General}`；派駐驗證失敗時不產生通知，維持 fail-closed。
- `cmd/wlgame/main.go` 讓通知 modal 優先於同一 tick 的撥款／外交 pending，按 Enter／Space 關閉後才進入三選一撥款視窗。`TestQueuedFundingInitialTalkNotices` 固定兩個索引與 marker 目標；Go／GUI 編譯於 `wolong-go:20260809` Docker 通過。
- 邊界：事件 4／5 的 `sub_139E8` 後續 TALK 池、原版三列訊息順序、PC-98 數字視窗、欄寬／游標與完整玩家流程仍未完成。

## 2026-08-09 — `sub_13902`／`sub_139E8` 指定金額 outcome 勘誤

### 證據輸入與定位契約

- 輸入：唯讀 DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm`、`KI.EXE`。
- 雜湊：`.i64` SHA-256 `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`；`.asm` SHA-256 `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`；binary SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
- 工具：`ida-pro-9.4-ver2:uidfix-v1`／IDA Pro 9.4；以下 `000xxxxx` 均為 IDA 線性位址，不是檔案偏移。Go 測試使用 `wolong-go:20260809` Docker。

### 直接證據與接入

| 原始定位 | 直接行為 | 推論等級／狀態 |
|---|---|---|
| `sub_13902` `00013902` | 指定外交金額回到 `DX` 後，比較 `[BP+0Ah]` 初始要求；超過時把回傳碼寫成 3，事件外層 `cmp AL,2` 不呼叫 `sub_135ED` | **已證實／已接入**：`ResolveDiplomacy` 超額 fail-closed |
| `sub_13902` `00013902` | 輸入 0 時把回傳碼改成 0，等同無條件收尾；輸入 1..初始要求仍保留付費／指定分支 | **已證實／狀態效果接入**：0 不扣款，正值沿用輸入 |
| `sub_139E8` `000139E8` | 指定撥款金額為 0 時選項碼為 2；高於原始要求時為 3，但尾段仍寫入輸入金額並扣款 | **已證實／已接入**：撥款 0 無效，超額可完成 |
| `sub_138C7`／`sub_138E6` + `sub_13771` | 初始外交 terms 在進入 `sub_13902` 前已計算；確認流程沒有第二次呼叫 terms | **已證實／已接入**：pending 保存 `InitialAmount`，避免平手 RNG 重抽 |

### 驗證與邊界

- `TestDiplomacyAndFundingAmountOutcomeBounds` 固定外交超額無外交狀態收尾但扣 Trust 30、撥款 0 無副作用、撥款超額照輸入金額完成；`TestQueuedDiplomacyChoice`、`TestQueuedFundingChoice`、`TestRawAmountEditorSemantics` 仍通過。
- 這次只修正 state outcome 與 RNG 時序；PC-98 數值視窗掃描碼、欄寬／游標、原版 TALK 訊息順序與逐頁排版仍是未完成 parity 邊界。

## 2026-08-09 — `CS:08A4` TALK formatter 表與事件 4／5 後續訊息

### 證據輸入與定位契約

- 輸入：唯讀 DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm`、`KI.EXE`；`.i64` 只在 `ida-pro-9.4-ver2:uidfix-v1` 容器內複製副本分析，沒有修改原始資料庫或原始素材。
- 雜湊：`.i64` SHA-256 `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`；`.asm` SHA-256 `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`；binary SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
- 工具：IDA Pro 9.4；下列 `000xxxxx` 均為 IDA 線性位址，不是檔案偏移。TALK 資料以 `tools/talkdat.py` 在 `wolong-go:20260809` Docker 內讀取。

### formatter handler 的直接資料流

`sub_1084A`（由 `sub_1075B` 呼叫）遇到反斜線後，把 ASCII marker 減去 `0x31`，以 `CS:[SI+08A4h]` 的 word table 取 handler；因此 marker 的 state 數值不能直接代替 ASCII 字元。IDA 讀出的已用表項如下：

| TALK marker | table word（CS offset） | IDA 線性 handler | 直接行為／語意 | 等級 |
|---|---:|---:|---|---|
| `\1` | `08B2h` | `000108B2` | 讀 formatter stack 的武將索引，定位 `0x4240` 武將表並繪姓名 | 已證實 |
| `\2` | `08DBh` | `000108DB` | 讀 formatter stack 的據點索引，定位 `0x0840` 據點表並繪姓名 | 已證實 |
| `\3` | `0904h` | `00010904` | 讀勢力記錄的 `+1` 君主索引，再繪君主姓名 | 已證實 |
| `\4` | `0939h` | `00010939` | 取 `word_10CFD` 玩家勢力記錄的 `+2` 軍師索引，再繪軍師姓名 | 已證實 |
| `\6` | `097Eh` | `0001097E` | `DI += 2`、`DX -= 30h`，不呼叫字串／數字繪製器 | 已證實；排版控制 |
| `\7` | `0984h` | `00010984` | 讀 formatter stack 的 word，交給十進位數值繪製器 | 已證實 |

`\5`、`\8`、`\9` 的表項雖可由 IDA 讀出，這次事件 4／5 接縫沒有需要替它們建立新的 state 語意；未使用項維持未知，不批次命名。

### `sub_139E8` 的 TALK base／offset

事件 4 呼叫 base `CX=0116h`（TALK #278），事件 5 呼叫 base `CX=013Fh`（TALK #319）。依 `sub_139E8` 的 `base+5` menu、`base+6+[code]` 結果與 `base+10+5×choice` 收尾控制流，及 TALK.DAT 原始行交叉確認：

| 玩家路徑 | 事件 4（內政） | 事件 5（外交） | state 效果 |
|---|---:|---:|---|
| 前置要求 | #278 | #319 | 尚未寫入 |
| menu | #283 | #324 | 尚未寫入 |
| 全額 | #284 → #288 | #325 → #329 | 以原始要求撥款 |
| 指定＝初值 | #284 → #293 | #325 → #334 | 撥款 |
| 指定＜初值 | #285 → #293 | #326 → #334 | 以輸入金額撥款 |
| 指定＝0 | #286 → #293 | #327 → #334 | code 2；無副作用 |
| 指定＞初值 | #287 → #293 | #328 → #334 | code 3；仍以輸入金額撥款 |
| 拒絕 | #286 → #298 | #327 → #339 | 無副作用 |

### Remake 接入與驗證

- `cmd/wlgame/messages.go` 將已證實的 `\4` 映射到玩家勢力 `Advisor`，`\6` 映射為不可見排版控制；原有 `\1`／`\2`／`\3`／`\7` 仍走結構化 state 目標與 Big5／十進位展開。
- `cmd/wlgame/funding.go` 新增 `fundingTalkIndices` 與 `enqueueFundingTalk`，在玩家確認後依上述原始分支排入結果與收尾 modal；指定 0／拒絕雖然 state 回傳 false，仍顯示各自原版 TALK，不落入泛用錯誤句。
- `TestFundingTalkIndicesMatchRaw139E8Branches` 固定 12 個事件 4／5 分支。`cmd/wlgame` 測試在有界 Docker `Xvfb` 內通過；完整 Go／文件／發行閘門仍需在本輪收尾重跑。
- 呈現邊界：文字 modal 保留 TALK.DAT 硬斷行，但沒有冒充原版 `DX` 欄位位置、數值字型、游標、逐頁動畫或 `\6` 的像素排版；事件 4／5 初始零要求的未使用流程也不由本節外推。

## 2026-08-09 — 事件 2／3 玩家外交請求前置 TALK

### 證據輸入與定位契約

- 輸入：唯讀 DOS/V `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm`、`KI.EXE`，以及
  `workplace/orig/dosv/TALK.DAT`。
- 雜湊：`.i64` SHA-256 `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`；
  `.asm` SHA-256 `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`；
  binary SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
- 工具：`ida-pro-9.4-ver2:uidfix-v1`／IDA Pro 9.4；本節 `000xxxxx` 均為 IDA 線性位址，
  不是檔案偏移。TALK 文字由 `tools/talkdat.py` 在 `wolong-go:20260809` Docker 內以
  Big5 解析；原始輸入唯讀。

### 直接證據

| 玩家事件 | 原始定位 | 直接行為 | 等級 |
|---|---|---|---|
| 事件 3 停戰請求 | `sub_138C7` `000138C7` | 呼叫 `sub_13902` 前將 `CX=0168h`，對應 TALK #360 | **confirmed** |
| 事件 2 協力請求 | `sub_138E6` `000138E6` | 呼叫 `sub_13902` 前將 `CX=0175h`，對應 TALK #373 | **confirmed** |
| `{3}` 參數 | `sub_13C3D` `00013C3D`、formatter `00010904` | 使用請求方勢力記錄定位君主姓名；不是把事件字高轉成數字 | **confirmed** |

原始 `TALK.DAT` 的 #360 是「{3}前來請求停戰……意下如何？」；#373 是
「{3}前來請求協助……意下如何？」。這一段只確認進入玩家三選一前的前置報告，
不把 `sub_13902` 後續的回應／理由變體或數值欄位外推成已知。

### Remake 接入與驗證

- `internal/state.beginDiplomacy` 在 pending 三選一成立時，分別建立
  `TalkNotice{Index:0x168,Faction:Source}` 或
  `TalkNotice{Index:0x175,Faction:Invader}`；通知先由 `cmd/wlgame/messages.go`
  展開，下一個 input 才進入外交選項，世界時間不穿透 modal。
- `TestQueuedDiplomacyChoiceTalkNotices` 驗證事件 2／3 的 index 與請求方；
  `TestDiplomacyTalkExpansionUsesOriginalRequestMarkers` 以真實 `TALK.DAT`／劇本驗證
  `{3}` 代換、原始中文行與 fail-closed marker 行為。兩者均在有界 Docker（GUI 測試含
  Xvfb）通過。
- 明確邊界：事件 2／3 的接受／拒絕後 TALK 池、PC-98 數值輸入視窗、原版欄寬／游標／
  逐頁動畫與長期正常劇本仍未完成；本節只封閉前置請求報告。

## 2026-08-09 — 事件 2／3 `sub_13902` 選項／主要結果 TALK

### 直接證據

`sub_138C7`／`sub_138E6` 都呼叫 `sub_13902`；其控制流在選項後先把
`base+4+[choice]` 傳給 `sub_13CDC`，再由事件 handler 呼叫 `sub_13C3D` 顯示主要結果：

| 路徑 | 建言／選項（base+4/+5/+6） | 主要結果 base | response 0／1／2 |
|---|---|---:|---|
| 事件 3 停戰（base #360） | #364／#365／#366 | #43（`0x2B`） | #43／#44／#45 |
| 事件 2 協力（base #373） | #377／#378／#379 | #47（`0x2F`） | #47／#48／#49 |

`sub_13902` 的指定外交金額超過初值時回傳 `AL=3`；`sub_13C3D` 對這個特殊碼仍以
response 2 的主要破裂句呈現，事件外層也以 `AL>=2` 不套用外交狀態。指定金額為 0
回到 response 0；正值且不超過初值才帶 `\\7` 金額進入 response 1。這些判斷與 state
收尾的 `TestDiplomacyAndFundingAmountOutcomeBounds` 相互獨立，避免把結果文字誤當成
資金效果證據。

### Remake 接入與驗證

- `cmd/wlgame/diplomacy.go` 新增 `diplomacyTalkChoiceIndex`、
  `diplomacyTalkResponse` 與 `enqueueDiplomacyTalk`；選項建言與主要結果按順序進入既有
  TALK modal，`\6` 保留為空的不可見控制，`\3`／`\7` 分別代換君主／十進位金額。
- `TestDiplomacyTalkIndicesMatchRaw13902Branches` 固定 6 個事件 2／3 選項索引與
  超額 response 2；`TestDiplomacyTalkExpansionUsesOriginalRequestMarkers` 以真實
  DOS/V TALK／劇本驗證低額成功的 #378→#48 與 marker 展開。Docker／Xvfb 目標測試通過。
- 邊界：`sub_13C3D` 的 AH／信賴度次要訊息（#367–#372／#380–#385 的選擇條件）仍未知，
  未接入；PC-98 數字欄位、逐頁／肖像動畫、跨平台 GUI 仍不能標成完成。

## 2026-08-09 — 事件 9 釋放武將的玩家通知條件勘誤

### 直接證據

- `sub_13485`（IDA `00013485`）把事件字高轉成 `0x4240` General 記錄指標後呼叫
  `sub_150D7`（`000150D7`）。
- `sub_150D7` 先清除 General `+0x1D` 俘虜方，依原俘虜方是否存活寫回 `+0x1C`；
  接著以 `cmp AL, byte_10CFF` 比較釋放後所屬勢力與玩家勢力，只有相等時才以
  `CX=0x25` 呼叫 `sub_18810`，即 TALK #37。之後的 `CX=0x199` 是另一個 formatter
  參數流程，TALK #409 在 DOS/V 資料中為空槽，不能把它猜成第二句。
- 輸入為唯讀 DOS/V `KI.EXE.asm`／`TALK.DAT`；binary／asm 雜湊與 IDA Pro 9.4 位址契約
  沿本日志前節，沒有修改原始輸入。

### Remake 勘誤與驗證

- `cmd/wlgame/messages.go` 現在只有在 `World.Generals[id].Faction == World.Player`
  時排入 #37；釋放回其他勢力或在野不再錯誤顯示玩家通知。
- `TestReleasedGeneralTalkOnlyTargetsPlayerFaction` 以真實劇本／TALK 驗證玩家勢力與
  非玩家勢力兩條分支；Docker／Xvfb 目標測試通過。
- state 的 `releaseGeneral`／`ReleasedGenerals` 仍保留已證實持久欄位，事件 9 的未知空槽／
  formatter 參數、原版肖像／逐頁畫面與完整長程 oracle 不外推。

## 2026-08-10 — 事件 11／12 runtime marker 的地圖呈現接縫

### 變更範圍與證據等級

本節是呈現層接線，不新增反組譯語意。事件 11／12 的 marker、持久效果與 TALK #70／#71／#72
證據沿用本日志前段：`KI.EXE.i64`／`KI.EXE.asm`／`KI.EXE` 為唯讀輸入，IDA Pro 9.4 線性
位址與雜湊契約不變。新增程式碼的證據等級是 **integration／substitute**：已證實 state
狀態被讀取，但原版 `sub_123FF`／`sub_12438` 的物件圖形、動畫幀與畫面座標尚未解出。

### Remake 接線

- `World.DisasterMarkerAt(cityID)` 回傳 runtime marker 的 `economy.Disaster` 與強度副本；
  `World.StormAreaSnapshot()` 回傳暴風雨範圍副本。兩者對無效索引／空值 fail-closed，且不
  暴露可寫內部陣列、不進存檔。
- `cmd/wlgame.drawDisasterOverlay` 以城市地圖格座標畫火災／暴動／暴風雨低干擾向量標記，
  並以輪廓畫 11×11 `StormArea`；繪製順序為地圖後、橫幅與浮動視窗前。
- `TestDisasterMarkerReadOnlySnapshots` 固定 marker、強度、範圍副本隔離；`cmd/wlgame`
  事件／外交短測試與 Docker build 通過。這些測試驗證接線，不宣稱原版動畫 parity。

### 邊界

事件 11／12 的原版物件資源／幀序列、事件 10 producer／formatter、完整 TALK 翻頁與
跨平台 GUI 實機仍未完成；本切片不改變「三平台正式包與推廣影片須等所有串接完成」的決策。

## 2026-08-10 — 事件 10 producer／formatter 負證據收斂

### 直接 dispatch 證據

- DOS/V `KI.EXE.asm` 的 `sub_131AE`（IDA 線性位址 `000131AE`）先取事件佇列低位元組、
  減一後索引 `funcs_131E8`；表中第 10 個 1-based 槽（低位事件碼 `0x0A`）指向
  `sub_13496`（`00013496`）。因此「事件 10 有 handler」是 **confirmed**。
- `sub_13496` 只做 `AL ← AH`、`AH ← 0xFF`、`CX ← DX`，把堆疊上的兩 byte 與 `AL=0x93`
  交給 `sub_18810`，隨後返回；本函式沒有可證實的存檔／世界狀態寫入。這裡只能確認
  formatter／訊息邊界，不能從 `0x93` 猜出 TALK 槽或完整句子。

### producer 搜尋結果與等級

- 在同一份唯讀 `.asm` 中，直接呼叫 `sub_12FBF` 的策略 producer 只看到事件碼
  `0x010C`／`0x020C`（事件 12 火災／暴動）、`0x0B`（事件 11 暴風雨）、`0x04`／`0x05`
  （撥款）、`0x0D`（事件 13）；`sub_12FB1` 的動態呼叫前置只寫入 1／2／3／8，
  `sub_1301C` 的事件寫入只看到 6／7／9／12。已檢查的 caller 沒有把低位碼 0x0A
  送入佇列，這是 **negative evidence／未完成**，不是「全 binary 不可能產生」。
- `mov al,0Ah`／`mov ah,0Ah` 的其餘命中落在數值顯示、戰術繪圖／輸入與硬體繪圖路徑，
  未形成「事件碼 → queue helper」證據；沒有把它們誤接成事件 10 producer。

### IDA `.i64` 交叉參照補查

在 `ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）內，將唯讀
`workplace/ida/dosv/KI.EXE.i64` 複製到一次性 Docker 暫存目錄，使用受版控的
`tools/ida_xref.idc`；輸出保留原始函式名與 IDA 線性位址，沒有修改來源資料庫：

- `sub_12FBF` 的 7 個直接 caller 是 `sub_12286`（`000122AF`／`000122CC`）、
  `sub_122DB`（`00012346`）、`sub_12FB1`（`00012FBA`）、`sub_15715`
  （`0001577E`）、`sub_1578F`（`000157ED`）、`sub_157FE`（`00015824`）。
- `sub_12FB1` 的 5 個直接 caller 是 `sub_12D3A`（`00012D51`）、`sub_12E33`
  （`00012E81`）、`sub_12E89`（`00012EE7`）、`sub_12EFB`（`00012F68`）、
  `sub_12F71`（`00012FA8`）。
- `sub_1301C` 的 4 個直接 caller 是 `sub_1300E`（`00013017`）、`sub_134B1`
  （`00013503`）、`sub_15940`（`00015981`）、`sub_16623`（`000166A9`）。

對 caller 函式的 IDA 函式匯出再核對：`sub_12286` 只建立 `0x010C／0x020C` 災害事件，
`sub_122DB` 只建立 `0x0B` 暴風雨事件，`sub_12D3A` 只以 `AL=8` 建立事件 8；其餘
caller 已由事件 1／2／3／4／5／6／7／9／12／13 的既有研究對上。這把「已檢查的
直接 xref 範圍」從平面 grep 升格為 IDA database 證據，但不排除仍未定位的函式指標或
間接寫入；事件 10 仍是 unknown／未完成。

### Remake 邊界

目前仍不建立事件 10 state producer、TALK index 或泛用訊息；`World` 的事件 10 handler
維持無狀態邊界。下一次若要收斂，必須以 IDA `.i64` 的完整交叉參照／資料流追出外部或
間接 queue writer，再以原始 `TALK.DAT` 槽位與呼叫參數交叉確認。此負證據不解除正式包與
推廣影片的 release gate。

## 2026-08-10 — 事件 6／7 次要 formatter 邊界再確認

### 原始控制流

- `sub_13327`（IDA 線性位址 `00013327`）與 `sub_13388`（`00013388`）在主要結果前分別
  呼叫 `sub_136C4`／`sub_13712`，再以 `CX=0x2B`／`0x2F` 呼叫 `sub_13C3D`；response
  0／1／2 對應主要 TALK #43–#45／#47–#49，這部分維持已接入的 **confirmed** 結論。
- `sub_137D8`（`000137D8`）將 `AH` 的 bit 0／1 由兩次 `sub_13138` 結果組成；
  `sub_13138`（`00013138`）掃描 `0x4240` General 表的 `+0x1D`／`+0x1C`，所以可
  **confirmed** 判定為兩方向「持有對方舊主的俘虜」旗標，不把它命名成信賴度。
- `sub_13C3D`（`00013C3D`）第一次呼叫 `sub_18810` 前建立 `DI=SP` 的 formatter
  參數堆疊；恢復 `DI` 後，若 `AH` 非零且 response 不是 2，第二次只做
  `CX += 0x1D`、`AL=0x93`、再呼叫 `sub_18810`，沒有重新建立同一個 `DI=SP` 參數
  堆疊。這是直接反組譯證據。

### 為何本輪不接 UI

以 raw index 算術，事件 6 的次要呼叫會落在 #72／#73，事件 7 會落在 #76／#77；原始
`TALK.DAT` 這些槽位分別含城市／勢力 marker 或選單文字，與第二次呼叫缺少參數堆疊的
狀態不一致。這是 **strong inference／未完成邊界**：可以確認程式確實有一次次要呼叫，
但不能把它當成目前 `TalkNotice{City,Faction,General,Amount}` 可安全展開的普通通知，
也不能用事件 2／3 的 #367–#372／#380–#385 反推其語意。下一步需用原版可重播狀態或
IDA `.i64` 的 formatter 資料流與畫面結果共同確認；本輪不新增猜測文字或錯誤 marker。

## 2026-08-10 — 進言／說得信賴度結果碼與理由級別定案

### 直接反組譯證據

- `sub_13830`（IDA 線性位址 `00013830`）先呼叫 `sub_13C1E`，把全域
  `byte_10D00` 轉成 `AH=1/2/3/4`，再把外交常式回傳的 `AL` 保存到 `[bp+2]`。
  `AL=1` 走 `sub_13D91` 並加入 `0x14`；`AL=0` 或 `3` 走 `sub_13DC9` 並扣除
  `0x14`；`AL=2` 進入 `sub_13B5A`，完成後先將 `0x14` 右移一位再呼叫
  `sub_13D91`，所以是 `+0x0A`。`AL>3` 只走訊息邊界，不改信賴度。
- `sub_13B5A`（`00013B5A`）把理由選擇交給 `sub_13BA9`（`00013BA9`）。
  `[bp+3]` 是 `sub_13C1E` 的信賴度級別；成立理由會遞減它，歸零才回報成功。
  選到不成立理由時，`sub_13BA9` 直接以 `AL=0x14` 呼叫 `sub_13DC9`；撤回碼
  `AL=4` 與重複理由不改信賴度。
- `sub_13C1E` 的比較常數是 `0xE0`、`0x90`、`0x20`，故原始分段為
  `信賴度≥E0／≥90／≥20／<20 → 1／2／3／4` 個成立理由。
- `sub_13D91`／`sub_13DC9` 對 `byte_10D00` 分別做加／減，並在 `0xFF`／`0x00`
  飽和。輸入為唯讀 DOS/V `KI.EXE`／`KI.EXE.asm`；SHA-256 與 IDA Pro 9.4
  位址契約沿本日志前段，完整局部輸出保留於
  `workplace/ida/dosv/func-sub_13830.txt`、`func-sub_13B5A.txt`、
  `func-sub_13BA9-current.txt`。

### Remake 接線與測試

- `persuasion.Situation.Trust` 保存說服開始時的 `World.Trust`；`Begin` 依四段
  原始級別建立 Session。`TrustOnReasonSuccess=10`、`TrustOnImmediateSuccess=20`、
  `TrustOnFailure=-20`，`ReactionTrustDelta` 對應第一反應 `AL=0/1/2/3/4`。
- `cmd/wlgame` 的 `situation` 傳入 Trust；直接拒絕、直接同意與「已交戰」分支也
  透過 `adjustTrust` 接回 byte 飽和，理由選擇沿用同一 Session 變更。
- 新增四段邊界、理由成功／錯選增減與第一反應碼測試；Docker 內
  `go vet ./cmd/wlgame ./internal/state ./internal/rules/persuasion` 及三套件
  `go test -p=1 -vet=off ... -count=1` 通過。

本節只封閉玩家進言／說得的信賴度切片；事件 2／3 外交回報的其他增減、事件 6／7
次要 TALK、事件 10 producer、事件 11／12 原版物件動畫、Windows／macOS 目標 GUI
與原版／remake 同狀態對拍仍未完成，因此不解除三平台正式包與推廣影片 release gate。

## 2026-08-10 — 事件 2／3 超額外交的信賴度副作用

### 直接證據

- `sub_13902`（IDA 線性位址 `00013902`）在玩家指定金額大於 `[bp+0Ah]` 的原始要求
  時把回傳 `AL` 設為 `3`；事件 2／3 外層以 `cmp AL,2 / jnb` 跳過外交狀態收尾。
- `sub_13C3D`（`00013C3D`）先把 `AL=3` 暫時改成 `2` 顯示主要破裂句，恢復後命中
  `cmp al,3 / jz`，以 `AL=1Eh`、`CX=1A5h` 呼叫 `sub_13DC9`。因此這個玩家分支
  的確是信賴度 `−30`，不是「完全沒有副作用」；`response=2` 的一般拒絕則不走這個
  `−30` 呼叫。

### Remake 勘誤與驗證

- `World.ResolveDiplomacy(DiplomacyOfferFunds)` 在超額 fail-closed 前，以
  `clampU8(Trust−0x1E)` 接回這個已證實副作用；資金、交友度、侵攻目標與 pending
  外交仍不收尾。
- `TestDiplomacyAndFundingAmountOutcomeBounds` 現在同時固定外交狀態無收尾與 Trust
  扣 30；Docker 內 `cmd/wlgame`／`internal/state`／`internal/rules/persuasion`
  vet、unit tests 通過。事件 2／3 AH／次要 TALK、PC-98 數字視窗與其餘結果仍列未完成。

## 2026-08-10 — 普通箭初始速度與 `sub_1ECE0` 接線勘誤

### 直接證據

- `sub_1ACA4`（IDA 線性位址 `0001ACA4`）取 X／Y 兩軸絕對差，回傳較大軸差於 `BX`，
  並以 `AL=目標高度−射手高度` 回傳有號高度差。
- `sub_1AD2D`（`0001AD2D`）執行 `shr bx,1`、`add bx,ax`、`call sub_1ECE0`、
  `and ax,3`、`add ax,bx`，再以 `0x14` 乘出普通箭初始速度；這是已證實公式，
  不是只依外觀補出的近似。
- `sub_1ECE0`（`0001ECE0`）以 `byte_1ECFD` 查 256-byte table，加上
  `byte_1ECFC`，再以 `0x89` 推進；`sub_1EC82`（`0001EC82`）負責遞增表、時鐘播種
  與 256 次交換。`internal/rules/rng` 已按此實作。

證據輸入是唯讀 DOS/V `workplace/ida/dosv/KI.EXE.asm`（SHA-256
`FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`）與
`KI.EXE`（SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`）；
工具為 `ida-pro-9.4-ver2:uidfix-v1`／IDA Pro 9.4，所有位址均為 IDA 線性位址。

### Remake 接線與勘誤

- `internal/rules/tactical.normalProjectileVelocity` 現在明示重現上述公式，使用
  `Battle.rng` 的 `Next() & 3`；新增 `TestNormalProjectileVelocityMatchesRawSub1AD2D`。
- `wolong-go:20260809` Docker 內 `go vet ./internal/rules/tactical` 與
  `go test -p=1 -vet=off ./internal/rules/tactical -count=1` 通過。
- 舊交接文字把這條寫成「普通箭 RNG 未接入」；新證據足以否定的是「公式／RNG 演算法
  尚未接入」這一層，不否定仍未完成的原版／remake 同狀態時序、`BATTLE.SCH` 圖形、
  完整動畫與畫面對拍。因此 release gate 仍維持關閉。

## 2026-08-10 — M7 文意審查第一批（#0–#359）

### 輸入與工具

- 日文輸入：`workplace/orig/pc98/TALK.DAT`，SHA-256
  `537e563269e414da79381ff48184a98e062ca454eb5d12e16a5fcbd52b79cf6f`；抽出檔
  `translations/extract/talk-pc98.json`，SHA-256
  `5821772a97255f6b5f2d84679e5e6b7b289460f8301aa0142ba884f1a24bb19a`。
- 繁中輸入：`workplace/orig/dosv/TALK.DAT`，SHA-256
  `08a22e09791d0a6ec2968e87d8655e12c91b45e00fae460b28593b35ff85e384`；抽出檔
  `translations/extract/talk-dosv.json`，SHA-256
  `f98cfea4f80333630a99817988ae2a9b2eb76e8e8edbd144e12d46aa0fc667ed`。
- 解析／校訂工具：`tools/talkdat.py`、`tools/talkdat_selftest.py`；研究工作樹基準
  `e01c667c00cc120d4d819c591d33b701e271bd9e`；命令在 `wolong-go:20260809`、
  Python 3 標準函式庫、`--network none` 容器內執行。

### 逐句裁定

第一輪實際讀取 #0–#359 的日中句對；只有可由原句直接證明的項目才寫入
`translations/corrections.json`：

| 類型 | 編號 | 證據／處理 |
|---|---|---|
| 明確文字／標點／名詞錯誤 | #16、#31、#32、#34、#66、#114、#135、#138、#139、#149、#213、#266、#268、#314、#347、#351 | 修正日文殘留、錯字、殘字、標點、`備蓄`／`軍情` 等可直接對照項 |
| 付款／句法錯誤 | #44、#48 | 日文 `金{7}` 與中文代入結果共同證明缺少付款語意，並移除重複助詞 |
| 語意方向錯誤 | #211、#259、#352、#359 | 依日文否定／肯定、主語與外交結果直接修正，不改 marker 集合 |

既有 16 筆變數／槽位／內容校訂保留；本批新增 22 筆，合計 38 筆 `fix`。驗證結果：

```text
tools/talkdat.py correct        套用 38 筆
tools/talkdat_selftest.py       38 筆、round-trip、mismatch guard 通過
TALK.DAT dosv／pc98 verify      34,182／45,718 bytes，1022 則，byte-for-byte 相同
```

`correct` 顯示的行數／行寬變更只是待原生 640×400 畫面抽樣的警告，不升格成排版
parity。#360–#1021 尚未逐句審查，因此 M7 與正式三平台包／推廣影片 release gate
仍未解除。

## 2026-08-10 — M7 文意審查第二批（#360–#1021）

### 輸入與研究邊界

- 本批在 Docker 內讀取 `translations/extract/talk-pc98.json` 與
  `translations/extract/talk-dosv.json` 的 #360–#1021；兩份抽取檔與原始 `TALK.DAT`
  雜湊沿用上一節，未修改 `workplace/orig/` 唯讀素材。
- 工具為 `wolong-go:20260809` 內的 Python 3 標準函式庫、`tools/talkdat.py` 與
  `tools/talkdat_selftest.py`；容器使用 `--network none`、UID 1000:1000。
- 裁定原則是只升格日文原句與繁中現況可以直接證明的錯誤；未有 formatter／畫面證據的
  語氣、行寬與重譯選擇保留，不宣稱已完成原版排版 parity。

### 逐句裁定

| 類型 | 編號 | 證據／處理 |
|---|---|---|
| 明確文字／標點錯誤 | #405、#426、#440、#459、#470、#518、#581、#617、#660、#797、#886、#906、#967 | 修正 `再／在`、`遵／尊`、職稱、方向詞、日文詞形、殘留標點與重複句段 |
| 明確語意錯誤 | #428、#524、#649、#663、#770、#840、#843、#844、#889 | 依原文主語、秘密遷都、條件句、控制城壁、戰鬥吶喊與暫時退卻重接繁中語意 |

本批新增 22 筆；連同既有 38 筆，`translations/corrections.json` 現有 60 筆可重跑
`fix`。Docker 驗證如下：

```text
tools/talkdat.py correct        套用 60 筆
tools/talkdat_selftest.py       60 筆、round-trip、mismatch guard 通過
TALK.DAT dosv／pc98 verify      34,182／45,718 bytes，1022 則，byte-for-byte 相同
```

因此 #0–#1021 的第一輪逐句讀取已完成，並留下 60 筆有直接證據的校訂；下一個 gate 是
把校訂輸出接到 remake 的可重現載入路徑，逐句抽樣 640×400／16×16 CJK 硬換行、行寬、
formatter marker 與翻頁，再做原版／remake 同狀態畫面／狀態對拍。M7 的畫面 parity、
其餘事件接線、目標平台 runtime 與正式三平台包／推廣影片 gate 仍維持關閉。

## 2026-08-10 — PC-98 數值視窗外框 blit 尺寸收斂

- 輸入：唯讀 DOS/V `workplace/ida/dosv/KI.EXE.asm`，SHA-256
  `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`；binary `KI.EXE`
  SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
  工具為 `ida-pro-9.4-ver2:uidfix-v1`／IDA Pro 9.4，位址基準為 IDA 線性位址。
- `sub_17D0D`（`00017D0D`）在 `sub_17C6E` 的數值視窗流程中使用 `AL=51h`，設定
  `DS=CS:word_10D50`、`SI=0600h`、`AX=4006h` 後呼叫 `sub_1FA37`。
- `sub_1FA37` 把 `AX` 傳給 `sub_1FAA2`；`sub_1FAA2` 以 `DH=40h` 迴圈 64 列，
  每列 `DL=06h` 次、每次複製兩個 byte。以 EGA planar 每 byte 8 pixels 計，外框資源
  的拷貝矩形是 96×64 像素；這是尺寸／資源偏移的已證實結果，不是外框內容或游標語意。
- `docs/re/13-pc98-numeric-window.md` 已更新。仍未知：資源圖像實際四角／邊界內容、
  `CS:7D93h` 18 個 cell byte 的語意、游標動畫與完整 TALK 逐頁契約；不以這條尺寸證據
  提前解除事件 2／3／4／5 數值視窗 release gate。

## 2026-08-10 — 事件 6／7 次要 TALK：formatter 缺參數負證據

### 輸入與位址契約

- 唯讀輸入：`workplace/ida/dosv/KI.EXE.asm` SHA-256
  `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`、
  `workplace/orig/dosv/KI.EXE` SHA-256
  `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
- 工具：`ida-pro-9.4-ver2:uidfix-v1`／IDA Pro 9.4；所有反組譯位址是 IDA 線性位址。
  原始 `.i64`、`KI.EXE`、`TALK.DAT` 均只讀，沒有改寫或覆蓋。

### 已證實結果

- `sub_13327`（事件 6）與 `sub_13388`（事件 7）分別以 `CX=0x2B`／`0x2F` 呼叫
  `sub_13C3D`（`00013C3D`）。第一個結果呼叫建立 `DI=SP` formatter stack；恢復
  `DI/CX/AX` 後，條件 `AH!=0 && AL!=2` 才以 `CX+0x1D`、`AL=0x93` 第二次呼叫。
- 因此事件 6 的直接次要 index 是 `0x48`／#72，事件 7 是 `0x4C`／#76；事件 2／3
  的同一 helper caller 也分別得到 #76／#72。`AL=3` 則先走主要 response 2，再由
  `sub_13DC9(AL=0x1E,CX=0x1A5)` 扣 Trust 30，不要和次要呼叫混為一談。
- `sub_137D8`（`000137D8`）將 `sub_13138`（`00013138`）掃描兩方向的 General
  `+0x1D`／`+0x1C` 關係寫入 `AH` bit；這是俘虜關係 raw flag，非 Trust。`\\2` handler
  （`000108DB` 附近）從 `SS:[DI]` 消耗 word，而第二次呼叫沒有提供同樣的 stack word。

### 判定與勘誤

- PC-98 原文 #72 含 `{2}` 城市 marker；#76 是選單文字。現有 `TalkNotice` 沒有可靠
  的第二次 formatter payload，故不新增 state／GUI 映射。#73／#77 不是四個直接 caller
  的 `CX+0x1D` 結果；它們仍是 unknown。這收窄了先前交接寫得過寬的 #72/#73、#76/#77
  範圍，並保留原始歷史紀錄，不升格猜測。
- Docker 內暫存存檔的 PC-98 事件 oracle 沒有越過 `NEW GAME`／讀檔選擇畫面，沒有取得
  有效事件畫面；標記為 inconclusive，不作行為證據。原始素材未修改。
- 結論：事件 6／7 次要 TALK 是「raw call 已證實、formatter 契約未知」；正式包／影片
  gate 不變。

## 2026-08-10 — 事件 9 `CX=199h` 空槽收斂

- `sub_13485`（`00013485`）將事件字高轉成 General 指標後呼叫 `sub_150D7`
  （`000150D7`）。若釋放後所屬勢力是玩家，先以 `CX=25h` 顯示 #37，再以
  `CX=199h` 呼叫第二次 formatter。
- 在 Docker 內以 `tools/talkdat.py dump` 讀取唯讀 PC-98／DOS/V `TALK.DAT`，#409
  均只有資料上的空行，沒有任何文字 byte 或 marker。這不是「尚未找到文案」；
  原始表明確把該槽留空。
- `cmd/wlgame/messages.go` 現在對全空行訊息 fail-closed，不建立空白 modal；新增
  `TestReleasedGeneralRawFollowup409IsEmptyNoOp`，並與既有玩家勢力 owner 條件測試在
  Docker／Xvfb 通過。這封閉事件 9 的可見文字接縫，不宣稱原版空呼叫的堆疊時序 parity。

## 2026-08-10 — M7 校訂表 runtime 接線

### 可重現輸入與產出

- `translations/extract/talk-dosv.json` SHA-256：
  `f98cfea4f80333630a99817988ae2a9b2eb76e8e8edbd144e12d46aa0fc667ed`。
- `translations/corrections.json` SHA-256：
  `f97ad2191adcb93fda2f5819dc647ef6fbe445eb42d18307d85c067d1c8bd93e`。
- 產出 `translations/talk-dosv-corrected.json` SHA-256：
  `a8c2c4335fff1b08fde1f7bd48d7a13a969cdee03716a7c98995cdc4821cba7e`。
- 產生器 `tools/talkdat.py` SHA-256：
  `5413123acb3db73ae44a530513f8f114413e052612295fc9cbd1a26bbd23cfb2`；
  `tools/talkdat_selftest.py` SHA-256：
  `6fe7e86a0b8f380cdc25655780dd1e6e3a8e11727c25eb65cdb7cf1de95a18d3`。

### 接線與判定

- `internal/assets/text.LoadJSON` 只把 UTF-8 的 1,022 則呈現行重建成 `text.Part`，保留
  `{N}` marker 與尾端空行，使用 `golang.org/x/text` cp950 encoder；raw `TALK.DAT`
  的 `Parse`／byte-for-byte round-trip 不走此路徑，也不會被覆寫。
- `internal/assets/library.LoadWithOptions` 提供可選 `TalkJSON`；`cmd/wlgame` 預設載入
  版控校訂表，`-talk-json ""` 才回到 raw。`talkdat_selftest.py` 先以目前來源重建，再
  逐 byte 比對 checked-in 產出，並驗證校訂後 round-trip。
- Docker 工具映像 `wolong-go:20260809`、Go `go1.25.12 linux/amd64`：文字 selftest、兩版
  `TALK.DAT verify`、JSON 單測、完整 `go vet`／`go test`、有界 Xvfb 短 smoke、deny-list
  selftest 與文件索引均通過。校訂工具仍報出需畫面抽樣的行寬警告；本項不升格為 formatter、
  PC-98 數字視窗或逐頁畫面 parity。

## 2026-08-10 — M7 remake modal 行寬 guard

### 可重現檢查

- `tools/talkdat_selftest.py`（本輪 SHA-256
  `1f8fff1c8b5065f4999949c7f1d6a49d84d37827417b9b497d807cc98ace929d`）在校訂後 JSON 上以
  22 個全形格（352 px）檢查每一行；這是
  `cmd/wlgame.drawMessage` 384 px modal、左右 12 px 內容內距的保守上限。`{N}` formatter
  token 按三個全形格計，故未依賴某個具體武將／勢力名稱。
- 執行環境：`wolong-go:20260809`、Python 3、`--network none`、UID/GID 1000:1000；目前
  1,022 則校訂表最大保守寬度 13 格，selftest 通過，產出仍與 checked-in artifact 逐 byte 相同。

### 判定

- **已封閉**：remake modal 的保守像素行寬安全性。
- **仍未知／未完成**：原版硬換行與欄寬、formatter 缺參數、游標、逐頁動畫及原版／remake
  同狀態畫面。此 guard 不解除事件 6／7 次要 TALK、事件 10、PC-98 數字視窗或三平台
  正式包／推廣影片 release gate。

## 2026-08-10 — 原版事件 6 fixture oracle 勘誤

### 輸入、fixture 與工具

- PC-98 原始 `KI.EXE` SHA-256：
  `061917f9f3f5c03e29397a9c636d546052128a99b8c8ce31ded0e84cf2a481e8`；PC-98 原始
  `SAVE.DAT` SHA-256：
  `18aa181327d0a6f1410ebdd47a1a4281ef50a45d443acd52712724318aa1f62c`。
- 使用者提供的 DOS/V `SAVE.DAT` SHA-256：
  `59e27270ee8192f63b08e012bf31a0b8da1477ce2c643fd487e6f8181b7650d2`；fixture
  `SAVE.DAT` SHA-256：
  `c695eec5ce61d5eda2eb5927d0ee33dd023d639c16ca47e01cdb45eb0fe24245`。
- 原始 PC-98 `TALK.DAT` SHA-256：
  `537e563269e414da79381ff48184a98e062ca454eb5d12e16a5fcbd52b79cf6f`。
- 執行環境：`wolong-dosboxx:latest`、DOSBox-X `2025.02.01` commit `32b2c24`，
  `machine=pc98`、`cycles=20000`、`--network none`、UID/GID `1000:1000`；PC-98
  原始目錄唯讀掛載，fixture 與截圖只寫 `/tmp`。

fixture 是一次性測試資料，不是原版自然遊戲存檔：

1. 以 DOS/V 第 1 槽作為狀態基礎，將區塊 `+0x0000` 設為已由劇本資料證實的
   `196/4/1` 時鐘。
2. 將玩家勢力指標 `+0x0D` 設為 `0x0040`；將勢力 3 的 `+0x2A` 外交官設為同一
   存檔中所屬勢力 3 的有效武將 17。
3. 在 `+0x52C0` 寫入事件字 `0x0306`、Param `0`。事件 6 的高 byte＝回報勢力 3；
   `sub_13327`／`sub_136C4` 的狀態計算不使用此 Param。上述欄位與位址以 IDA
   線性位址對照 `KI.EXE`，不是對原始素材就地改寫。

### 原版操作與觀測

- 用既有的 PC-98 紅色游標閉迴路，完成 `NEW GAME` → `NO` → 第 1 個 `LOAD DATA`
  槽；讀檔畫面顯示 `196年 4月 1日`。
- 讀檔後原版在地圖上顯示事件 6 結果視窗：外交官肖像與日文訊息的語意為
  「停戰交涉的結果，金 14000 成立」。這是 `sub_13327` 的主要結果畫面證據，
  不是 remake 轉譯或泛用訊息。
- 第一次原地按下後回到地圖；第二次按下因游標仍位於讀檔槽對應位置而開啟據點
  情報。這個 fixture 沒有觀察到可辨識的次要 TALK；它只能證明此輸入狀態的主結果，
  不能證明所有 `sub_13C3D` 次要分支不存在或已完成。

證據截圖：[`pc98-oracle-event6-result.png`](docs/images/pc98-oracle-event6-result.png)，
SHA-256 `fd7cbfa4cf2d4e0773a33181c4c4ea7944020552f5facc847750f69f670fa72b`。
先前「暫時 oracle 未越過啟動選擇畫面」的紀錄保留為歷史嘗試；本節是後續勘誤，將
事件 6 的**主要結果**提升為原版畫面證據，但事件 6／7 次要 formatter、原版數值視窗、
事件 10、原版物件動畫、同狀態 parity 與正式三平台包／影片 gate 仍未封口。

## 2026-08-10 — `MMAP.MCH` 火災／暴動物件圖形證據與接線

### 輸入、工具與位址基準

- 唯讀輸入為 DOS/V `KI.EXE`（SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`）與
  `MMAP.MCH`（SHA-256
  `b10a5b64bbffa672c1fb5cb37703ac4c14b18bf1166cc47c4e802c19aae9f8f7`）；PC-98 與
  DOS/V 的 MCH 雜湊相同。IDA 資料庫為 IDA Pro 9.4，`KI.EXE.i64` SHA-256 為
  `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`；以下位址均為
  IDA 線性位址、載入基準 `0x10000`，未修改 `.i64` 或原始素材。

### 可證實資料流

- `sub_187AF`（`000187AF`）載入 `MMAP.MDL`，並把 `MMAP.MCH` 放在其後的 `+0x800`
  paragraphs；`sub_187CC`（`000187CC`）設定工作段，故 `word_1987A` 指向 MCH 檔案
  `0xA000`。`sub_1D804`（`0001D804`）顯示 MCH 圖塊是 16×16 遮罩加四個 16 色平面，
  每塊 160 bytes，前 `0xA000` bytes 共 256 塊。
- 事件 12 的 `sub_134B1` 將高 byte 1／2 交給 `sub_123FF`（`000123FF`）建立火災／
  暴動物件；`sub_12533`（`00012533`）以 `type*8+frame` 查 `CS:[bx-67A6h]`，IDA
  資料庫的線性位置 `0001985A` 實際內容為：type 1
  `18 19 1A 1B 1C 18 19 1A`、type 2 `20 21 22 23 20 21 22 23`、type 3
  `28 29 2A 2B 28 29 2A 2B`。`0001D51F`（該區無函式邊界，由 `sub_12533`／`sub_12B3C` 呼叫）再把 MCH 矩陣中的 tile ID
  置入戰略地圖 40×23 cell buffer。
- MCH `0xA000` metadata 的 4-byte entry 指向 `0xA100` 後的矩陣資料；原始資料確認
  type 1 pattern 為 16×9（144 tile IDs），type 2／3 為 5×5（25 tile IDs）。type 1
  frame 0 的矩陣第 5 格為 `0xD0`，並非以寬高猜測或向量圖替代。

### Remake 接線與邊界

- `internal/assets/world/mmapmch.go` 新增長度、metadata、平面圖塊遮罩與 pattern 解碼；
  `internal/assets/library.Library.MCH` 讀取 `MMAP.MCH`。`cmd/wlgame` 將事件 12 火災／
  暴動 marker 對回 type 1／2，依 season palette 合成圖像並快取，缺檔或不合法資料
  fail-closed 回退原有低干擾向量 marker。
- `internal/assets/world/world_test.go` 固定原始大小、pattern 維度、查表 index、透明／
  不透明色值與 `0xD0` fixture。`drawDisasterOverlay` 的大型圖像裁切改按實際 image
  bounds 判斷，避免城市在畫面邊緣時因仍使用單一 16×16 marker bounds 而錯誤跳過。
- `sub_12459`／`sub_12533` 的 `[si+0x0C]` timer 與 `[si+0x0F]` frame 狀態未保存在目前
  remake `DisasterMarker`；因此目前 `(frame/8+1)&7` 只是呈現層固定相位時鐘，屬明示
  替代，不是原版動畫時序證據。事件 11 `sub_134A6`→`sub_1237E` 沒有呼叫物件建立
  函式，故暴風雨仍保留已證實的範圍輪廓，不硬接 type 3。

### 驗證判定

- `wolong-go:20260809`、`--network none`、UID/GID `1000:1000`、有界 Docker／Xvfb 下，
  `go test -p=1 -vet=off ./internal/assets/world ./internal/assets/library` 與
  `go test -p=1 -vet=off ./cmd/wlgame` 通過；完整長程遊玩依使用者指示略過。
- 本項只封閉 MCH 資產與事件 12 火災／暴動圖形來源；事件 6／7 次要 formatter、事件
  10 producer、原版動畫 timer parity、PC-98 數字／排版 parity、Windows／macOS GUI
  runtime、同狀態原版對拍與正式三平台包／推廣影片仍是開啟 gate。

## 2026-08-10 DOS/V parity 增量：數值內框／次要 TALK／物件 timer

### 數值內框資源

- 輸入：DOS/V `KI.EXE` SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`、
  `ICONGRF.DAT`；IDA Pro 9.4，位址基準為線性位址。
- `sub_17D0D`（`00017D0D`）設定 `DS:SI=word_10D50:0600h`、`AX=4006h` 呼叫
  `sub_1FA37`；`sub_1FA37`／`sub_1FAA2` 證明來源矩形是 96×64 平面圖。
- `sub_100DF` 的段指標鏈將該來源換算到 `ICONGRF.DAT` 第 3 段相對 `0x14A0`；
  `internal/assets/gfx.DOSVAmountPanel` 與 `Spec.DecodeAt` 以 byte offset 解碼，
  不把非連續段內資源誤切成 frame index。`cmd/wlgame` 目的地為 `(88,184)`，保存區
  `(80,176)`／112×80、raw 格位 `(88,200)` 不變。

### 事件 6／7 次要 TALK

- `sub_13C3D`（`00013C3D`）第二次呼叫的條件與 `CX+0x1D` 已直接接入 state：事件 6
  `0x48`／#72，事件 7 `0x4C`／#76；`TalkNotice.Secondary` 保存 raw 性質，#76
  `NoPortrait` 保存 direct text／menu 呈現政策。
- 第一次呼叫才把 `DI=SP` 建立 formatter stack；第二次沿用恢復後 DI，#72 的 `\\2`
  不可安全映射城市。呈現層缺 marker 整則 fail-closed；不把 `AH` 解成城市、勢力或
  信賴度。state／message 單測固定 sentinel `-1`，避免 Go 零值污染 raw 邊界。

### 事件 10

- `sub_13496`（`00013496`）的 consumer 已保存 `Param` 為 raw TALK index；事件字高
  byte 對應原版 `FFxx` formatter word，只有 `0..126` 才映射 `General`／`\\1`。
- 尚無 `0x0A` producer 的直接證據；目前結論是可消費、無已知產生路徑，不虛構劇本。

### 物件動畫

- `sub_123FF`／`sub_12459`／`sub_12533` 的 runtime record 已落到
  `internal/state/disaster_objects.go`。建立 timer=1、interval=16、phase=1；每次
  可見 map-loop update 遞減，dirty render 先畫舊 phase 再遞增 `&7`，`sub_12438` 同城
  清除同步。`TestDisasterObjectAnimationTiming` 固定此時序；不寫入存檔。
- `sub_12459` 最後一筆 object 的 `sub_1248A` 移動分支仍未接，type 1／2 靜態城市物件
  的 timer／frame parity 已完成。完整長程遊玩與 DOS/V 自然畫面 oracle 依使用者要求略過。

## 2026-08-10 DOS/V 硬體游標／按鍵 glyph 解碼

### 證據與位址基準

- 唯讀 DOS/V `KI.EXE` SHA-256：
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 反組譯輸入：`workplace/ida/dosv/KI.EXE.asm`；IDA Pro 9.4／
  `ida-pro-9.4-ver2:uidfix-v1`；以下 `seg002`／`sub_201E4`／`sub_2020C` 都是
  IDA 線性位址，不把 file offset 當成同一個位址空間。
- `seg002` 的 image offset 是 `0x10000`，MZ header 是 `0x200`；因此
  `seg002:031B` 對應原始檔 `0x1051B`。`sub_201E4`（`seg002:01E4`）設定
  `SI=031Bh`，連續呼叫 `sub_2020C` 兩次；每次消費 32 bytes。

### 游標 mask

`sub_2020C` 每列讀一個 word，先把 `AH` 寫到 `[bx]`、再把 `AL` 寫到 `[bx+1]`。
因此每列的兩個 source byte 在畫面順序中必須反轉，再以 MSB-first 展開。
`sub_201E4` 第一次先設 `AX=0F00h`（EGA set/reset 0x0F，白色外框），第二次設
`AX=0A00h`（EGA set/reset 0x0A，紅色填色）。兩個 32-byte mask 合成 16×16 游標；
解碼計數是白 39、紅 56、透明 161，palette-index buffer SHA-256：
`385c2f1949d3d1e331399316305db7d7f2fd489a0626de1c6b2b8375aadfc6fe`。

`internal/assets/gfx/cursor.go` 保留 `seg002:031B`、file `0x1051B` 與原始顏色／
透明常數；`DecodeDOSVCursor` 拒絕不符合 MZ header／image layout 的檔案，避免靜默
讀錯 offset。`cursor_test.go` 固定 16 列形狀與上述三種像素計數。

### 按鍵 glyph 資源

`sub_17D0D`（IDA `00017D0D`）的 `DS:SI=word_10D50:0600h`、`AX=4006h` 經
`sub_1FA37`／`sub_1FAA2` 對應 `ICONGRF.DAT` 第 3 段相對 `0x14A0` 的 96×64
平面資源。其下半部從 `(88,200)` 開始正好是 3×6 個 16×16 格；各格的內容像素
均大量不同於格背景。`TestDOSVAmountPanelContainsStaticButtonGlyphs` 對六欄三列
逐格要求至少 128 個非背景像素，證明可見 glyph 是資源內靜態圖像，不是
`CS:7D93h` 的 18-byte hit-test table。

### Remake 接線與勘誤

- `Library` 讀取 `KI.EXE` 後快取 cursor；`cmd/wlgame` 由 `Library.DOSVCursor` 建立
  Ebiten image。數值面板有 `amountFrame` 時直接畫原始 96×64，舊 vector 矩形／CJK
  button label 只在資源缺失時 fallback，不覆蓋真實 glyph。
- `go test -p=1 -vet=off ./internal/assets/gfx ./internal/assets/library ./cmd/wlgame`
  已在無網路 Docker／有界 Xvfb 通過；沒有在主機啟動遊戲或圖形工作負載。
- 此節修正先前「硬體游標與按鍵 bitmap 尚未取得足夠資源證據」的現況文字。DOS/V
  資產解碼與接線已完成；自然 DOS/V／remake 整張畫面逐像素對拍、PC-98 視覺基準、
  其他事件 formatter／物件動畫與目標平台 GUI 仍各自是未完成 gate。

## 2026-08-11 事件 6／7 raw formatter、事件 10 producer 與 GUI 邊界

### 事件 6／7 次要 formatter

- 輸入：DOS/V `KI.EXE` SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`；IDA Pro 9.4，
  `ida-pro-9.4-ver2:uidfix-v1`，IDA 線性位址。
- `seg000:000108DB` 的 formatter 讀取 `SS:[DI]` word；高 byte 為 `FF` 時轉成
  `0x0840 + city×0x20`，否則直接當 DS offset，兩條路徑都再前進 2 bytes。remake
  以 `World.ResolveTalkFormatter2` 保存這個 raw 邊界，越界或不存在資料時回傳
  invalid，不猜城市。
- `sub_13327`／`sub_13388` → `sub_13C3D` 已證實事件 6／7 次要索引為 #72／#76，
  條件為 `AH != 0` 且 response 非 2／3。事件 6 回報勢力的 record pointer 在已知
  0x400-byte original stack 內時保存 raw word 0；事件 7 沒有可知 formatter word。
  `TalkNotice.RawFormatterWordValid` 防止 Go 零值誤被當成有效 word 0。
- 驗證：`TestSecondaryTalkUsesCapturedRawFormatterWord`（當時叫 `…UsesOriginalRawFormatterWord`）以 DOS/V raw TALK table
  解出 word 0 的 bytes，確認 explicit raw word 會顯示；valid flag 關閉時同一通知
  fail-closed。state fixture 另固定雙向俘虜旗標、索引、順序與 `NoPortrait`。

### 事件 10 producer

- IDA 交叉參照：`word_10D56` 的 writer 是 `sub_100DF`；`sub_12FBF` 的 callers 為
  `sub_12286`、`sub_122DB`、`sub_12FB1`、`sub_15715`、`sub_1578F`、`sub_157FE`；
  `sub_12FB1` 的 callers 為 `sub_12D3A`、`sub_12E33`、`sub_12E89`、`sub_12EFB`、
  `sub_12F71`。已檢查的事件 producer 常數為 `0x010C`、`0x020C`、`0x0B`、`0x06`、
  `0x07`、`0x09`、`0x0C`，沒有直接證實的低碼 `0x0A`；這是負證據，不是全 binary
  排除。DOS/V 四個 `SINARIO.DAT` event queue 及 SAVE block 0 的現有事件也沒有
  低碼 `0x0A`。
- `sub_13496` consumer 的已證實資料流仍是事件字高 General、`CX=Param` TALK index。
  `World.QueueEvent10` 因此只做受控 raw producer：驗證 `general`／`talkIndex`，以完整
  256 格 queue 寫入 `(general<<8)|0x0A` 與原始 `Param`。它是 remake fixture／劇本
  注入口，不是原版自然時序的宣稱。
- 驗證：`TestEvent10ProducerWritesRawTalkPayload` 覆蓋 payload、前 64 格滿後仍寫入
  第 65 個空槽，以及 General／TALK index 邊界拒絕；`TestQueuedEvent10TalkNotice`
  覆蓋 consumer 的有效／無效 General formatter 邊界。

### DOS/V 自然畫面與目標平台 GUI

- Docker／DOSBox／Xvfb 固定 `machine=vgaonly`、`cputype=486`、`cycles=fixed 20000`
  執行 DOS/V `START`，自然啟動停在「密碼輸入：第 09 頁」，那一輪就沒有再往下走，
  也沒有把原版密碼頁與 remake 策略畫面宣稱為同狀態 parity。
  （**後來測出密碼頁不擋**：四格留白按確定就進開場，`docs/playtest/18`。）
- remake 使用 DOS/V `workplace/orig/dosv`、倚天字型、`scenario=0`、`player=0`、
  `seed=17`、`speed=0`，30 幀輸出 `docs/images/wlgame-dosv-natural-remake.png`；
  SHA-256 `8420d97955be60af16da403544b47e84b3f44363ef75f867930e022d1bc2f916`。
- Linux／Xvfb GUI smoke 通過；Windows amd64 交叉產物為 `PE32+ x86-64`，macOS
  amd64／arm64 交叉產物為 `Mach-O`。使用的 macOS SDK 是 `MacOSX15.5.sdk`，編譯器
  為 `x86_64-apple-darwin24.5-clang` 與 `aarch64-apple-darwin24.5-clang`。沒有
  Windows／macOS 原生桌面 runtime，故只記錄 build gate，不升格為 GUI parity。暫存
  產物 SHA-256：Windows amd64 `ea7bbb9747cb37bc47ae716ac1795f9cd15cce5b110eb8f78d3ed602489b8816`；
  macOS amd64 `a8894b8245cac1895c8341b18082efb9a8e76cbc9fa4dc795f19484fa76de840`；
  macOS arm64 `50f8c9378b41a595df99e8495a06b78aeda23b05251bb4ae3744557f2366b8e2`。

## 2026-08-11 YouTube 自然畫面 oracle 勘誤與 HUD 對齊

- 使用者指定來源：[臥龍傳 呂布開局滅曹操](https://www.youtube.com/watch?v=af6xqcicXoI)。
  Docker／暫存 `yt-dlp` 讀得 metadata：標題相同、長度 567 秒、30 fps、478×360、
  上傳日期 2018-07-01；使用 format `18` 約 28.64 MiB，只作容器暫存。
- 擷取的代表幀為 20／80／160／240／320／400／480／550 秒，檔案分別保存在
  `docs/images/yt-wolong-natural-*.png`。80 秒原始幀 SHA-256
  `d33fff8d664e24321274310287dce38b82c82cfb62f3d0427e70dfd5bd301e08`；去除上下黑邊
  `y=29..327`、縮放至 640×400 的參考幀 SHA-256
  `c0217b8722bd44a22a112a2981b626126d5ee53d3e9f00498c6cbd018e08d6`。
- 影片策略畫面與說明書主畫面共同支持以下幾何：banner 32 px、command strip 32 px、
  左側大地圖 432×336（27×21 個 16 px 格）、右側 208 px；右側上方 192×128 minimap，
  下方為君主／軍師、信賴度、資金與三種預備兵。這推翻舊的「DOS/V 自然畫面只有滿版
  地圖、沒有常駐 HUD」工作假設；舊假設保留在歷史勘誤，不覆寫證據。
- 新增 `cmd/wlgame/strategyhud.go`：以原版 `ICONGRF` minimap、KAOGRF portrait、
  chrome 與同一份 state 數值來源畫出常駐 HUD；`main.go` 的自然地圖改為 27×21，
  災害 marker／storm clip 也改用 `strategyMapY`／`strategyMapW`。浮動命令／列表／
  事件視窗仍保留，避免模態驗收路徑被自然 HUD 取代。
- 重跑 Linux／Xvfb 30 幀 smoke，`docs/images/wlgame-dosv-natural-remake.png` SHA-256
  為 `961e583915d2e0e7b65cd51f637ec214530b68040ba0da5770add4b35cb46e30`；影片視覺／
  幾何對拍 gate 通過。影片 80 秒為 196 年 4 月 5 日，remake 為 196 年 4 月 1 日，
  來源又經有損縮放，故不把此結果升格成嚴格同狀態逐像素 diff。

## 2026-08-11 `sub_1248A` typed slice 與短 parity gate

- `internal/state/disaster_objects.go` 將 MCH moving slot 的 raw `[si+08]`／`[si+0A]`
  fixed-point word 保存為 `xDrift`／`yDrift`，只讓 slot index `>= 16` 進入移動分支，
  對應 `sub_12459` 呼叫 `sub_1248A` 的最後半區。沒有 RNG 時只維持 timer／dirty，
  不猜原版移動。
- `sub_124FF` 的低 byte 以 raw signed-byte 規則正規化：`random&7`、高 byte 變動
  `[-15,+15]`、低 byte 跨 `±15` 時帶出 whole-word `±1`，保留 `uint16` raw word。
  `advanceMovingDisasterObject` 再套用方向 byte、storm bounds 與 `-0x10..0x190`／
  `-0x10..0x110` wrap。
- `TestSub124FFMatchesRawSignedByteContract`、
  `TestMovingDisasterSub1248AUsesOnlyLastHalfOfRawSlots`、
  `TestMovingDisasterSub1248ARawWrapAndDirectionByte` 在 Docker 內通過。
- `cmd/wlgame/parity_gate_test.go` 新增四個可重跑 fixture：事件 2–5 的 TALK 分支與
  行寬／分頁、M7 60 筆校訂的 marker／寬度／5 行、一般／特殊投射物 raw frame，以及
  事件 9 玩家勢力／非玩家／在野／#409 空行。`tools/parity_gate.sh` 將它們和 raw
  state／tactical tests、`tools/talkdat_selftest.py` 組成短 gate。

## 2026-08-11 事件 10 producer 深度 IDA 審查

- 輸入為 DOS/V `KI.EXE` SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`；IDA `.i64`
  SHA-256 `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`；工具為
  IDA Pro 9.4、`ida-pro-9.4-ver2:uidfix-v1`，位址均為 DOS/V 線性位址。實際匯出由
  非破壞性 `tools/ida_event10_producer.idc` 完成；IDAPython runtime 在該 image 缺少
  可載入的 `libpython3.14.so.1.0`，故不以失效的 Python wrapper 取代 IDC／`.i64`。
- 已證實：`sub_131AE`（`000131AE`）每次取 queue 的 4-byte `Code`／`Param`，
  以 low byte `Code−1` 索引 `CS:funcs_131E8`；事件 10 entry 是 `sub_13496`
  （`00013496`）。handler 將 high byte 轉 `FFxx` raw formatter word、把 `DX` 原樣
  放入 `CX`，再呼叫 `sub_18810`，本身沒有 queue write。
- 已證實：`sub_12FBF`（`00012FBF`）與 `sub_1301C`（`0001301C`）是主要 queue
  writer，分別由 IDA code xref 列出的 callers 供給 `AX=Code`／`DX=Param`；
  `word_10D20` 是 queue offset、`word_10D56` 是 queue segment。`byte_131AD` 的
  `0x0A` 是 dispatcher cadence reset，不是 event code。
- 已證實的 caller 常數包含事件 4 `0x04`、事件 5 `0x05`、事件 7 `0x07`、事件 9
  `0x09`、事件 12 `0x0C`、事件 13 `0x0D` 與 `0x010C`／`0x020C` 路徑；在已列
  direct caller 中沒有低 byte `0x0A`。`sub_130CB` 是狀態讀取，不是 queue writer。
- 強推論：保存的 `.i64` function boundary／direct xref producer graph 內沒有原版
  自然 event 10 producer。這是負證據，不能排除未建函式的 far code、register／pointer
  間接寫入、外部 loader 或密碼保護後資料。
- 未知：原版自然觸發時序與 TALK index；DOS/V dynamic oracle 仍停在密碼頁，沒有以
  猜測補出 producer。結論依使用者指定限時封口為「深度靜態逆向完成、producer unknown」；
  `World.QueueEvent10` 僅作受控 raw fixture／劇本注入口，並由
  `TestEvent10ProducerWritesRawTalkPayload`／`TestQueuedEvent10TalkNotice` 固定。
- 完整 raw instruction、caller／data-ref 表與證據分級見
  [`docs/re/15-event10-producer.md`](docs/re/15-event10-producer.md)。

## 2026-08-11 — 無輸入自動 clock 與事件 10 節拍定位

| 斷言 | 證據 | 等級／邊界 |
|---|---|---|
| 玩家不下命令、滑鼠座標不變時，原版仍會推進世界 | DOS/V `.i64`：`sub_11BE0`（`00011BE0`）呼叫 `sub_11F7F`（`00011F7F`）；座標相同時 `00011FC8` 設 `byte_198A3` bit 7，主迴圈在 `00011C5B` 呼叫 `sub_11CD0` | **已證實**；完整鏈見 `docs/re/16-idle-clock-event10.md` |
| 無輸入路徑會推進據點、軍團、MCH 物件與時鐘 | `sub_11CD0`（`00011CD0`）的 `00011D13`／`00011D16`／`00011D19`／`00011D1C` 依序呼叫 `sub_13EFD`、`sub_125A3`、`sub_12459`、`sub_11D8E` | **已證實**；原版呼叫順序保留 |
| 日期的 driver 是時鐘，不是事件 10 | `sub_11D8E`（`00011D8E`）在每小時進位時呼叫 `sub_13E11`（`00011DEC`），並在換日／換月時進入日期／月結分支 | **已證實** |
| 事件 queue 每小時只「被呼叫」，不一定每小時取一筆 | `sub_131AE`（`000131AE`）先 `dec byte_131AD`；`sub_12BD9`（`00012BD9`）初始化 `7`，取 record 後重設 `0x0A` | **已證實**；第一次每時邊界第 7 次取，之後每 10 次取，且受 queue 尾端條件限制 |
| 已注入的事件 10 會在自動 clock 期間變成 TALK | `World.Tick` → `hourly` → `dispatchQueuedEvent`；`TestIdleClockDispatchesQueuedEvent10OnHourlyCadence` 以 7 個每時邊界驗證 `Code=0x030A`／`Param=0x42` 只產生一筆 `TalkNotice` | **已證實（remake state fixture）**；不是原版 dynamic trace |
| 事件 10 是自動 clock 的 driver | `sub_13496` 無 `sub_11D8E`／`sub_125A3`／`sub_12662` 呼叫；已知 writer graph 仍無直接低碼 `0x0A` producer | **否定；負證據**，仍不能排除未建函式／間接寫入／密碼保護後資料 |

這筆證據修正「事件 10 是每小時同步通知」的簡化說法：正確模型是

```text
滑鼠不動／玩家停手
  → sub_11CD0 idle path
  → sub_11D8E 推進子刻、時、日、月
  → 每小時 sub_13E11
  → 受 byte_131AD 節流的 sub_131AE
  → 若 queue record 的 low byte 是 0x0A，才進 sub_13496 TALK
```

remake 的功能層對應為 `game.timeRuns` → `World.TickMap`，不是
`World.QueueEvent10` → clock。`World.TickMap` 已將正常 map-loop 對齊為據點／軍團／
物件／時鐘；同一畫面的額外 `g.speed` 規則 tick 使用不含物件的 `World.Tick`，保留
物件每個可見 map-loop 一次的 cadence。完整 raw 指令、IDA hash、工具版本與位址基準
見 [`docs/re/16-idle-clock-event10.md`](docs/re/16-idle-clock-event10.md)。

## 2026-08-11 M7／事件 2–5／事件 9／推廣片驗收紀錄

- M7 以 `tools/m7_review.py --check`、`tools/talkdat_selftest.py`、
  `TestM7CorrectedTalkLayoutGate` 驗證；60 筆校訂均保留原始 marker 與硬行結構，人工抽樣
  使用 `docs/playtest/14-m7-review.md` 的六張 DOS/V modal 代表幀。自動警告的行數／寬度
  變更已納入人工抽樣，不把保守上限警告誤判成失敗。
- 事件 2–5 的 `TestEvent2To5FullTalkPageSampling` 以 DOS/V `TALK.DAT` 與校訂產出驗證
  36 個 raw TALK 頁面／18 組雙頁回應；每頁實際不超過五列與 modal 寬度，分支索引與四張
  代表幀見 `docs/playtest/15-event2-5-talk-sampling.md`。
- 事件 9 的 `internal/state/TestEvent9LongNaturalRoute` 以 9 subtick／時、delay=7、27
  小時 bounded fixture 驗證第 7／17／27 小時取件；`cmd/wlgame/TestEvent9LongNotificationRoute`
  驗證玩家 #37、非玩家／在野抑制與 #409 no-op。證據見 `docs/playtest/16-event9-long-route.md`。
- `tools/parity_gate.sh`、Docker 內全量 `go test -p=1 -vet=off ./... -count=1`（Xvfb）與
  `go vet ./...` 通過；`tools/index.py generate/check` 通過。未跑完整長程遊戲，符合使用者
  指示。
- 推廣片輸出 `dist/promo/wolong-remake-trailer.mp4` SHA-256：
  `fa4ebe6147a7050a56fc390a822e61da7f65fa07661a66d06d047ecdcceb5f40`；ffprobe 確認
  H.264 1280×720、AAC 44.1 kHz 立體聲、60.000 秒。配樂為 `tools/promo_score.py` 原創
  合成音，未使用原版 `BGM.DAT`／`SOUND.DAT`。

## 2026-08-11 事件 10 remake substitute

### 問題與證據邊界

原版 `sub_13496` 只證實 `Code.high → FFxx formatter`、`Param → TALK index`，
且在 `sub_11D8E` 的每時 queue consumer 下游；保存的 IDA Pro 9.4 `.i64` writer graph
仍沒有可證實的低碼 `0x0A` producer。這個 unknown 結論保留，不用近似實作回填原版證據。

已查到可借用的鄰近文字／狀態證據是 DOS/V `sub_1585F` → `sub_15940` 的月結俘虜路徑：
`General +0x18` 是倒數閘，`sub_15940` 使用 `rand < 0x20`／`0x40` 分支，並以
TALK `0x41`（逃走）／`0x42`（歸降）回報。這只足以支撐 substitute 的候選規則，不能
證明它原本透過事件 10 queue。

### Remake mapping

`internal/state/event10_approx.go` 在月結既有 queue producer 後執行，限玩家勢力、活著、
目前 `Faction == Player` 且 `Captor != noFaction` 的武將；每月最多一名。狀態先以
副本計算，raw queue 寫入成功後才提交：

| RNG／狀態 | substitute 行為 | raw payload |
|---|---|---|
| `rand&0xFF < 0x20` | 清 `Faction`／`Captor`／`Posted`，轉在野 | `(general<<8)|0x0A`, `0x41` |
| `0x20 <= rand&0xFF < 0x40` | 保留玩家勢力，清 `Captor`／`Posted` | `(general<<8)|0x0A`, `0x42` |
| 其他 | 不變、不排隊 | — |

`SetApproximateEvent10(false)` 是 raw fixture 的關閉開關；正常 `LoadScenario` 預設開啟。
事件寫入後仍由 `World.Tick`／`World.TickMap` 的每時 dispatcher 消費，因此沒有改變
「據點／軍團／物件／時鐘」更新順序，也沒有把 queue producer 接成 clock。

### 驗收

Docker／Go 定向測試：`TestApproximateEvent10ProducerUsesKnownRawContract`、
`TestApproximateEvent10ProducerIsBoundedAndDisableable`、
`TestApproximateEvent10ReentersIdleClockConsumer`，另與既有
`TestIdleClockDispatchesQueuedEvent10OnHourlyCadence`、`TestQueuedEvent10TalkNotice`
一併通過。此節的 remake producer 證據等級是 **substitute／強推論**，不是原版 parity。

## 2026-08-11 DOSBox／remake 可玩性專家驗證

### 問題

確認目前 remake 是否能以正常玩家輸入完成策略開局、編成、行軍、idle clock、遭遇、
存檔／讀檔與戰術入口，並與 DOSBox 原版的可觀測畫面分開記錄。完整長程遊戲依使用者
要求不執行。

### 工具與輸入

- DOS/V：`dosbox-run:latest`，DOSBox 0.74-3，`machine=vgaonly`、`cycles=fixed 20000`。
- PC-98：`wolong-dosboxx:latest`，DOSBox-X，`machine=pc98`、`cputype=486`、
  `cycles=20000`。
- remake：目前工作樹建置的 Linux binary，`-seed 17`，`dosbox-run:latest` 承載 Xvfb／
  xdotool／ImageMagick 截圖；原始資料唯讀掛載，存檔用 `/tmp` overlay。

### 觀察與證據

1. DOS/V 原版成功啟動至密碼保護第 15 頁；`original-dosv-password.png` SHA-256
   `a9a972edbb4c896a914a84acfe65b4e55a8f93a6ea532b6feef7036d527ed5bf`。這是已證實的
   啟動／阻擋，不是密碼頁後玩法證據。
2. PC-98 原版本輪成功到 `NEW GAME`；既有受控 oracle 的 `pc98-oracle-scenarios.png`
   與 `pc98-oracle-in-game.png` 證實劇本選單／自然戰略地圖。當前 headless image 沒有
   window manager，bus-mouse／焦點輸入重播不穩，故不把本輪輸入失敗解釋成原版不可玩。
3. remake 正常無 debug 旗標路徑：`A；Enter；Enter；3；Down + Space × 5；Enter；
   Enter；M；Enter；Enter；Down × 22；Enter；Enter；= × 64`。穩定重播到 196/6/28
   的「呂布 對 曹操／攻城／戰鬥指揮／委任」；代表幀與 hash 詳見
   `docs/playtest/17-expert-dosbox-remake.md`。
4. remake 系統視窗 `4 → S → 1 → Enter → L → 1 → Enter` 產生 88,832 bytes overlay，
   顯示儲存及讀取成功；原始檔案沒有寫入。
5. current-build `-open-siege` 戰術 GUI／`2` 號命令 smoke 成功；這是 debug fixture，
   既有無旗標正常戰術證據與本輪 DOSBox 原版輸入橋接問題分開保存。

### 判定

remake 的正常策略可玩性、idle clock、遭遇接縫與保存回讀通過；DOS/V 原版受密碼頁
阻擋，PC-98 原版仍是規則 oracle。畫面比較只判定邏輯尺寸／HUD／資料結構的可比性：
原版 640×400，remake 1280×800（2×）；鏡頭、日期或模式不同時不稱為逐像素 parity。
剩餘最小工作是為 PC-98 DOSBox-X 補隔離 window manager／輸入橋接，再重播一次原版
完整新局至戰術，而不是重新逆向已知規則。

## 2026-08-11 YouTube／推廣片像素差異驗收

### 輸入與工具

- YouTube 原版輸入是既有保存的 20／80／160／240／320／400／480／550 秒代表幀；
  80 秒幀已去黑邊／還原為 640×400，來源 metadata 為 478×360、30 fps、567 秒。
- remake 輸入為推廣片採用的 `docs/images/wlgame-dosv-natural-remake.png`，640×400。
- 研究對照片由 `tools/promo_yt_compare.sh` 在 `u5cht/dev` Docker 內產生；ImageMagick
  `compare` 只對保存幀執行，原始 YouTube 影片未加入儲存庫或發行包。

### 結果

全畫面 `AE=255003/256000 (99.61%)`、`RMSE=0.338208`；橫幅、命令列、地圖與右欄
的統計與 SHA-256 由 `docs/promo/yt-remake-pixel-review.md` 保存。新增 24 秒、1280×400
H.264 對照片，並保留自然畫面並排／差異 PNG。

### 證據等級與決策

原始 raw metric 是 **已證實的影像差異**；它不是同狀態原版 parity，因兩張畫面日期、
輸入、部隊／地圖狀態不同，且 YouTube 經有損縮放。依使用者決策，YouTube／推廣片比較
作為 DOS/V 視覺 oracle，**影片這條路**不宣稱嚴格同狀態逐像素 diff。
（後來另走一條：拿原版存檔開同一個局面、640×400 原尺寸逐區比，
見 `docs/spec/90`／`91`。密碼頁也測出**不擋**，`docs/playtest/18`。）

## 2026-08-11 DOS/V 自然策略骨架調整

### 證據與變更

- YouTube 80 秒 640×400 幀可重現的常駐分區仍為 banner 32 px、命令列 32 px、左側
  地圖 432×336、右側 208 px。這些幾何維持不變。
- remake 原先把 minimap 與下方情報各自畫成完整 `chrome.Window`，導致共享分隔邊多出
  8 px。現在先畫下方框，再畫上方 minimap／16 px 紅藍色標列，使色標列覆蓋共同邊界，
  君主頭像從原版對應的 y=184 開始。
- 情報框新增原版可見的三列身份區、信賴度、紅色分隔線與黑底資源區。reserve 中央
  raw glyph 仍 **unknown**，只畫幾何 fallback，證據等級沒有升格。

### 驗收

`TestDOSVNaturalStrategySkeleton` 固定 640×400、32／32、432×336、208、共同分隔邊
與下方框貼齊畫布。Docker／Xvfb 重拍基準幀
`docs/images/wlgame-dosv-natural-remake-skeleton.png`；YouTube 對照 raw metric 改善為
`AE=249178/256000 (97.34%)`、`RMSE=0.329145`。這是畫面骨架驗收，不是不同日期／狀態
的同狀態逐像素 parity。

## 2026-08-12 松崗 DOS/V 密碼頁動態勘誤

### 問題與舊結論

舊紀錄把「密碼輸入：第 NN 頁」列為 DOS/V dynamic oracle 的阻擋。重查後發現舊
`tools/dosbox.sh` 的 autoexec 只掛載 C: 而未呼叫 `START.BAT`，所以當時輸入 `1234`
的畫面其實是 DOSBox 命令列，不能作為密碼比較失敗的證據。

### 受控重播與證據

使用既有 `wolong-dosboxx:latest`（DOSBox-X 2025.02.01）而非新工具鏈；原始
`workplace/orig/dosv` 唯讀掛載，容器內複製到 `/tmp/game`。設定：
`machine=vgaonly`、`core=normal`、`cputype=486`、`cycles=fixed 20000`、
`mouse_emulation=integration`、`int33 max x=640`、`int33 max y=480`。

密碼頁穩定後重新取得 X11 client geometry，逐格點選／以 XTEST 送鍵後按畫面上的
「確定」。三組各自在全新容器與全新遊戲副本執行：空白、`0000`、`1234` 都進入原版
開場敘事。空白與 `1234` 的 10 秒後 640×480 截圖 SHA-256 同為
`10cd0e199e7bd944a4664f3a3e4debd94a3986df1c36ca7eac4ffaf208ebbd34`；`0000` 取樣在
同段文字漸現的不同時間點，SHA-256 為
`cba3762abb82dd14919acaa3da22849760a2621b04ba488d09b21df44b861055`。

### 判定與邊界

**已證實：** 在這個原版／DOSBox-X 輸入路徑中，密碼頁的確認流程不阻擋遊戲開場。
使用者所述「任意數字會過」成立，並且空白確認同樣會過。**未知：** 數字沒有在四格
回顯，故未證實它們是否被 `PASS.*`／`YNFONT.EXE` 儲存或比較；也未把 DOSBox-X 行為
外推成真實硬體的密碼比較結論。沒有修改任何二進位檔，且密碼頁不會實作到 remake。
可重查摘要、來源雜湊與座標契約見
[`docs/playtest/18-dosv-password-verification.md`](docs/playtest/18-dosv-password-verification.md)。
