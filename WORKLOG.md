# 工作歷程

較早工作紀錄保存在既有 [`WORKLIST.md`](WORKLIST.md) 與
[`RESEARCH-LOG.md`](RESEARCH-LOG.md)。目前狀態回查 [`CONTEXT.md`](CONTEXT.md)，
研究文件索引回查 [`docs/INDEX.md`](docs/INDEX.md)。

## 2026-10-08：比對式反編譯初步試點

- 使用者採用工具鏈辨識與三函式試點；Go 基線 `e8aeb99`，遊戲規則沒有修改。
- 核對既有 IDA archive 的 SHA-256，恢復 `ida-pro-9.4-idapython:py312-v1`，image ID 與私有 manifest 相符。沒有新建同功能映像。
- `tools/ida.sh` 新增 `probe` 入口。原始 EXE 唯讀、一次性 DB、非 root、限 CPU／記憶體／PID 與外層逾時；本次沒有使用舊 root 分支。
- 研究輸出使用既有 `workplace/` 下的 `matching-decompilation/`，沒有另建根層研究目錄。
- `tools/matching_decomp.sh` 重跑 IDA 與編譯試驗；三個函式的組語重組共 233 bytes 相同，沒有遮掉呼叫或重定位運算元。
- 每個函式的位元組突變、時鐘的錯誤呼叫目標與錯誤輸入雜湊均被拒絕。Borland 標記掃描有 `LOGO.EXE` 正對照。
- GCC 12.2 `-m16` 的十二組 C 試驗沒有匹配；原版編譯器、組譯器、連結器與語言維持未知，整檔重建與 Go 行為對照沒有完成聲明。
- 證據、輸入／輸出雜湊與重跑入口在 [`docs/re/89`](docs/re/89-matching-decompilation-pilot.md)。
- 新匯出自動合併 `docs/re/matching-semantic-index.json`，保留原始函式名、位址與運算元。未證實或未收錄的語意帶警示；索引本身綁定原始 EXE SHA-256。
- 新研究文件的三筆分流以 `evidence-only` 回填到既有台帳，原有 799 筆內容、分類與 Issue 對應保持不變。
- 完整重跑包含冷啟動 IDA、三函式組語匹配、十二組 C 試驗與正負對照。Shell 與 Python 語法檢查通過。
- 文件／資產門禁逐項執行。文件索引有兩張歷史圖片缺失，幽靈引用檢查有歷史研究檔缺失，教訓的 render／verify 有既有資料不同步與防線文字缺失。用 Docker 內的 `HEAD` 快照確認這四項同樣報錯，沒有以放寬門檻處理。
- `worklist` 與 `lessons` 自測的第一次執行被唯讀 `workplace/` 掛載擋住，分類為驗證環境問題。同一 Python image、同一命令改用容器 tmpfs 後通過。原始失敗與後續檢查 log 保存在 `workplace/matching-decompilation/checks/`。
- 分流驗證與嚴格研究索引生成通過，共 802 列。文件／資產檢查最終 19／23 通過，另外四項是上述已由 `HEAD` 重現的既有問題；沒有宣稱全專案門禁全綠。
- 研究輸出抽查均為 UID/GID `1000:1000`。本輪容器使用 `--rm` 並已退出；保留恢復的 IDA image 作目前研究工具鏈，其他專案容器與映像沒有清理。
- 尚未 commit 或 push。

## 2026-10-08：採用組語基準，完成完整 KI.EXE 重建

- 使用者決定先以組語建立基準，後續還原 C 函式。新增 `tools/assembly_rebuild.sh` 與組譯工具，沿用先行試點的 Docker images。
- IDA 9.4 匯出全 segments，共 24,376 條指令、55,392 bytes；113 條指令有 loader 重定位差異，按原始 MZ 表還原段值。
- 所有 IDA 指令以助憶碼組譯回相同 bytes。任一指令不匹配即停止，沒有程式碼 bytes fallback。MZ 標頭、116 筆重定位與其餘 11,195 bytes 的資料宣告一起重建。
- 整份 67,099 bytes 的輸出與原始 KI.EXE SHA-256 相同。只掛載組語與 linker script 的獨立容器冷重建也相同。指令、立即數常數與 MZ 重定位三組負對照均拒絕。
- 來源對映完整覆蓋檔案，函式台帳保存 739 筆 IDA 定位與分級語意；沒有宣稱 C 還原、原作者來源切分或 Go 行為對拍完成。
- 全量匯出初次讀取中文 JSON 時缺少 UTF-8，修正後以同一入口、同一 image 從乾淨輸入重跑。分類 item 尾 byte 時改查 item head，並以完整指令 bytes 總數作正對照。
- 現況與收據：[`docs/re/90`](docs/re/90-assembly-reconstruction.md)。先行試點及其 DB／收據保留，沒有覆寫歷史 checkpoint。
- 新工具 Shell／Python 語法、803 列分流、嚴格研究索引、過期斷言與資產 deny-list 均通過。文件索引與幽靈引用仍回報先行試點已由 HEAD 複現的歷史缺檔，本輪新文件與工具沒有新增引用錯誤。
- 原始 EXE、完整重建與來源獨立重建的雜湊一致；來源／DB／收據與文件交叉核對通過。研究輸出與本輪修改檔案擁有權抽查為 UID/GID `1000:1000`。
- 本輪容器已退出並自動移除。只清理由本輪 IDA wrapper 留下的個別 `.ida-*` 暫存目錄，其他研究輸入、既有 root-owned 檔案與其他專案資源保持原狀。沒有 commit 或 push。

## 2026-10-08：第一個 C 函式 sub_1ECE0

- 先以固定雜湊的 IDA 指令建立 spec/201，審查取數、寫回順序、AX、FLAGS、PUSH／POP／RET 及堆疊重疊契約，READY 後寫研究用 C。
- `tools/c_recovery/rng.c` 提供可讀的 typed state 函式與 16-bit 呼叫介面適配層。正式 Go 規則沒有修改，也沒有在原版資料上寫入。
- dosgolem 載入原始 EXE，保留真實 CS 的 IDAIn 入口，每次執行原始 13 條指令。四張固定表各枚舉 65,536 種 c/s，另測堆疊邊界與重疊。
- GCC O0/O2 各 263,680 組原版/C 對照與 263,168 組 Go 狀態對照相同，每版 258 次完整 1 MB 記憶體核對。計數器增量、AH、旗標來源與 BX 取回方式四種突變皆被拒絕；BX 突變在 alias 案例才暴露。
- spec/201 更新 CONFORMED，範圍限受控局部常式與 IF/TF=0。C 機器碼匹配、播種、呼叫時機、GUI 與正常玩家路徑沒有新的完成聲明。
- 對照器的 C11 全域旗標使 cgo 的 POSIX prototype 不可見，加入 GNU feature macro 後同一工具鏈重跑通過。gofmt 的單檔掛載無法建立相鄰暫存檔，改用容器 tmp 產生內容後回寫。
- 全量枚舉第一次在原版計時器中斷送達時超出 20 條局部預算。根因是入口 IF=1；按已定的局部隔離契約固定 IF/TF=0，保留原始 13 條判準，沒有放寬預算或忽略差異。
- 收據逐項驗證 source hashes、O0/O2 摘要與負對照，dosgolem 被 import 的 oracle/internal 原始碼在執行前後雜湊相同。其既有 cmd/probe/main.go 修改沒有被本次改寫。
- 現況在 docs/re/91 與 c-recovery-status.json；組語基準台帳保留建立時的快照，避免改寫先行收據。
- C 完成索引的來源雜湊、807 列分流、嚴格研究索引、過期斷言與資產 deny-list 通過；README 的規格數已包含新增的研究契約。文件索引與幽靈引用仍有已複現的歷史缺檔，新 C 文件與工具沒有新增引用錯誤。
- 收尾核對原版 EXE 雜湊、研究輸出擁有權與容器清理；本輪輸出均由 UID/GID 1000:1000 擁有，沒有 commit 或 push。

## 2026-10-08：持續 Goal 與每輪提交

- 使用者設定持續目標「完成臥龍傳 matching decompilation」，並明確授權每輪完成後 commit、push。
- 前一輪分類為 progress：整檔組語來源、第一個 C 函式、原版/C/Go 局部收據已形成可回查交付物。
- 本輪先核對 main 與 origin/main 沒有差異，再提交這批已驗證研究來源與文件；本機原版、DB、整檔組語、重建 EXE、含原版資料的研究產物保持 Git 忽略。
- 完整 Goal 維持 active。現有成果只證明組語基準及 sub_1ECE0 的局部 C 語意，不當成整個 C 還原與 remake 比較完成。
- 基準與第一個 C 函式已 commit 為 `1d21147`，push 成功，origin/main 指向同一提交。

## 2026-10-08：第二個 C 函式 sub_1EC82

- 先核對完整 IDA inventory 的 94 bytes 與 RTC 呼叫流程，建立 spec/202。原版讀 BIOS 回覆的 AL，不能默默假設為零；C 介面保留 RTC_AL。
- dosgolem 的現有 RTC 成功回覆不在支援範圍，以明示的 RAM fixture ISR 提供 CX、DX、AL，保留原始 INT／IRET 與原始 94 bytes。來源與固定回覆隨收據記錄。
- 全部 86,400 個合法時分秒、12 個非零 AL 案例與 8 個堆疊邊界，O0/O2 各 86,420 組完整原版/C 狀態相同；AL=0 的 86,408 組 Go 播種狀態相同。每版 85 次完整記憶體核對相同。
- 每案原版固定 3,369 條指令，沒有挑選時間或放寬停點。索引步長、小時倍率、XOR 來源三個 C 突變全部被拒絕。
- XOR 的 AF 本輪按 dosgolem 模型為零，明示硬體未定義；合法 BCD、IF/TF=0、分離堆疊與固定 RTC 是證據邊界，未外推實機時序或正常玩家流程。
- spec/202 更新 CONFORMED，C 狀態索引加入第二個函式；原版取數的來源與先前收據保持原樣。Goal 未達整個 C 還原完成。
- 完整收據、來源雜湊與重跑入口在 docs/re/92；所有 workload 仍在 Docker 內，原始輸入唯讀。
- 新工具語法、815 列分流、嚴格研究索引、過期斷言及資產 deny-list 通過。文件索引與幽靈引用的既有歷史缺檔保留，本輪新文件沒有新增引用錯誤。
- 本輪完成後依授權 commit、push；Goal 保持 active，兩個 C 函式不代表全 executable 的 C 還原與 remake 比較已完成。

## 2026-10-08：第三個 C 函式 sub_11D8E 與 Go 年份訂正

- 先核對固定 IDA inventory 與原始 137 bytes，建立 spec/203 的時鐘契約。原版年份比較是 1000，999 會再加到 1000；五組原版 probe 揭露舊 Go 的固定 999 上限差異。
- spec/204 審查 READY 後，修正 Go 換年，新增連續 999→1000→999 與五種入口年份的冷測。機制現況更新在 docs/mechanics/15，舊結論的訂正追加在 CONTEXT.md §6。
- 原版/C O0/O2 各 292,297 組相同：226,665 組合法日期、65,536 個年份入口、96 組 callee／等待案例。Go 排除 MOV callee 後 292,249 組相同，各版 286 次完整 1 MB 記憶體抽查相同。
- 四個 callee 的入口快照比較全部 14 個暫存器／段／FLAGS、日期與堆疊；等待比較輪詢次數、FLAGS 與清零欄位。原始受測 137 bytes 不變，RET／MOV 與 ready／count fixture 是明示的輸入控制。
- 日進位門檻、年份比較、呼叫順序與等待比較四種突變全部被拒絕。O0/O2 收據相同，來源、EXE、原始函式 bytes 與 dosgolem import 來源前後身分核對通過。
- 最初年份 probe 放在 /output，遭 Go internal import 邊界拒絕；移到工作樹的研究區並補正明確輸出掛載後，以同一原版與條件重跑。這是驗證環境問題。
- spec/203、spec/204 更新 CONFORMED，研究入口為 docs/re/93；現行 C 索引已有三個函式。證據限局部時鐘，完整 callee、硬體時序、1000 年 UI 與長期玩家流程未驗證，沒有 C 機器碼匹配完成聲明。
- 所有 workload 在非 root Docker 內執行，原版與外部 dosgolem 來源唯讀；本輪依既有授權 commit、push，完整 Goal 保持 active。
- 文件與工具檢查 21／25 通過。825 列分流、嚴格研究索引、過期斷言、校訂、資產 deny-list、新工具語法及所有正對照通過。四項失敗沿用已在 HEAD 重現的歷史缺檔、教訓資料不同步與既有防線缺失；新文件與工具沒有新增引用錯誤。
- 三個 C 狀態索引的來源身分全部通過。時鐘產物 UID/GID 為 1000:1000，原版 EXE 雜湊相同，專案既有 root-owned 路徑仍為 27 筆，沒有擴大修復範圍。

## 2026-10-08：月結呼叫鏈、六個 C 函式與赤字捨位訂正

- 前一輪時鐘成果 `228f30d` 已推到 origin/main，分類為 progress。開工確認工作樹乾淨，再按逆向、IDA 及文件職責路由載入當前入口。
- 新 IDA 9.4 probe 匯出月結主流程、資金加減、距離除數、預備兵上限與赤字函式，共 306 原始 bytes。DB、原始 bytes、operand、xref 與原始名稱保留在本機 c-economy/ida。
- spec/205 READY 後實作六個 C 函式。月結的 C 資金／赤字／RNG 實際相接，原版側保留相同真實指令；未還原的據點結算與尾端 callee 用明示 MOV／RET fixture。
- 小量 probe 通過 192 組。固定 RNG 的原版赤字 probe 揭露 Go 捨位錯誤：原版先讀負資金高位再 NEG，Go 曾先取絕對值再右移；合法欠款非 256 倍數時每兵種少扣 16。
- 依 spec/206 READY 修正 Go，將舊 -16000 的 992 預期訂正為原版 1008，新增六個邊界與非零固定 RNG 測試。原始證據、機制現況及 CONTEXT.md 訂正台帳同步。
- O0/O2 各 1,775,568 組原版/C 與 1,326,175 組 Go 對照相同，各版 867 次完整 1 MB 記憶體抽查相同。高低位資金、carry、距離、赤字與月結順序六個突變全部被拒絕。
- 經濟及 internal/state 冷測通過。狀態層最初缺少 x/text v0.40.0，既有 image 的 GOPROXY=off 使首次開網路仍無法載入；明示官方 proxy、按 go.mod／go.sum 下載後，以同一 image、同一測試範圍重跑通過，分類為環境問題。
- C 初次編譯的 hex literal 與 + 貼連成 preprocessing token，修正空白；Go 對照器的 Bytes 長度由 uint32 改成正式 API 的 int。沒有放寬判準、挑樣本或忽略差異。
- 六個新函式加入現行 C 台帳，共九個；語意索引附 proven 等級及 re/94 出處。組語與先前 C 收據保留，沒有完整 C matching 或完整玩家月結的完成聲明。
- 本輪依授權完成後 commit、push，完整 Goal 保持 active。
- 收尾文件與工具檢查 21／25 通過，833 列分流、嚴格研究索引、新工具語法及所有正對照通過。四項既有失敗仍為歷史圖片／研究輸入缺檔與教訓資料不同步／防線缺失；新文件沒有新增引用錯誤。
- DB 雜湊與探針路徑原先寫在同一行，過期斷言工具誤綁為腳本雜湊；改為分列工具、DB 路徑與 DB 雜湊後，原工具重跑通過。沒有放寬掃描器。
- 經濟及狀態層 Go vet 通過，九個 C 函式的現行來源雜湊核對通過。研究產物 UID/GID 為 1000:1000，原版 EXE 雜湊相同，既有 root-owned 路徑仍為 27 筆。

## 2026-10-08：五個據點結算 C 函式與原始整數寬度

- 前輪 `4e50f80` 已推送，分類為 progress。核對乾淨工作樹、現行 C 台帳與逆向／IDA／文件職責路由，接續據點收入與募兵。
- 新 IDA 9.4 probe 保留五個原始函式的 420 bytes、operand、chunk、xref 及 DB 身分。spec/207 READY 後實作 C，沿用未改動的經濟算術／堆疊介面。
- 據點結算接入既有 C 月結，替換前輪 MOV 收入 fixture。距離、收入、募兵、玩家／AI、資金、預備兵、赤字與 RNG 在 C 實際相接；九個世界更新尾端及重畫仍是明示 RET fixture。
- 固定原版反例查出玩家 low-word carry 與募兵 16-bit 累計兩個 Go 差異。依 spec/208 READY 修正；經濟、internal/state 冷測及 Go vet 通過。
- 密集募兵測試最初用概略比率推期望，忽略分次移位的捨位。直接讀原版後改為 51584／5952／7808，四組稅率／累計反例的原版/C/Go 相同。舊差異 probe 保留，沒有挑樣本。
- 首次單獨執行探針缺少 /orig 掛載，補上原版唯讀掛載後，用同一 binary、同一案例重跑；分類為驗證環境問題。
- O2 的 1,779,300 組原版/C、1,385,882 組 Go 與 435 次完整記憶體抽查相同。整批 600 秒上限在 O0 結束前送達；已確認工作終止，保留 O2，只補跑 O0 與負對照，沒有減少案例或變更判準。
- 既有六函式收據已用新 Go 來源完整重生，1,775,568 組原版/C、1,326,175 組 Go 與六個負對照仍通過；先前收據保留在 c-economy/pre-settlement-proof。
- O0 補跑完成，與已取得的 O2 收據逐位元組相同；五組負對照皆被拒絕，來源、DB、原版 EXE／SINARIO.DAT 與收據核對通過。
- 現行 C 台帳已加入五個函式，共十四個；語意索引保留原始名稱、位址與 proven／re/95 出處，沒有改寫先前 DB 或證據。
- 文件與工具檢查 21／25 通過，842 列分流、嚴格研究索引、過期斷言、新工具語法、校訂、資產 deny-list 與全部正對照通過。四項既有缺檔／教訓檢查失敗仍保留，新文件沒有新增引用錯誤。
- 本輪按既有授權 commit、push；完整 matching decompilation Goal 保持 active，十四個 C 函式與經濟接線不代表完整世界月結、玩家垂直鏈或 C 機器碼匹配完成。
- 原版 EXE／SINARIO.DAT 身分、十四個 C 函式來源雜湊與分級語意均核對。輸出 UID/GID 為 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出並移除。

## 2026-10-08：十一個月結世界更新 C 函式與真實 writer

- 前輪 `259df4f` 已推送，分類為 progress。開工核對工作樹、現行台帳與逆向／IDA／文件職責路由，選擇七個月結尾端及其 event writer、關係取址與暴風雨標記依賴。
- IDA 9.4 probe 保留十一個原始函式，共 1068 bytes，原始名稱、operand、chunk、xref 及 DB 身分不改寫。spec/209 READY 後實作 C，沿用原有經濟、據點及 RNG 介面。
- 原版與 C 均執行真實事件 writer 與 RNG，並比較 queue code 的低 byte 判空、指定／隨機起點、滿槽、向後搜尋及 CF。月底的官員、外交、災害與赤字通知不再用事件清單或 RET 替代。
- O0/O2 各 142,864 組原版/C 相同，各版 1117 次完整 1 MB 記憶體抽查相同；四劇本接線與七個刻意錯誤版本均通過原定拒絕閘門。
- 固定原版生產力 probe 揭露 Go 普通整數加法及乘積符號的寬度差異：65000／上昇值100／稅率30 的原版為 12114。依 spec/210 READY 修正後，三個邊界冷測與 560 個原版/C/Go 向量相同，最終 RNG 相同。
- 經濟、internal/state 冷測與 Go vet 通過。新矩陣最初缺少八個群組的外層括號，依語法定位補回後在同一工具鏈重跑，沒有改案例範圍或判準。
- 十一函式加入現行 C 台帳，共二十五個；每筆保留原始位址、C source hash、routine hash、proven 及 re/96 出處。
- 先前經濟及據點收據已保留在各工作根的 pre-world-proof，現行收據按新 Go 檔案身分重生。原始 EXE／DAT 一直唯讀。
- 剩餘 sub_1585F、sub_12BD9 與 UI／音效仍為明示 fixture，完整玩家月結及 C 機器碼尚未驗證。Goal 保持 active，本輪依授權完成後 commit、push。
- 收尾文件與工具檢查 21／25 通過，852 列分流、嚴格研究索引、過期斷言、新工具語法、校訂、資產 deny-list 與所有正對照通過。四項既有缺檔／教訓檢查失敗保留，新文件沒有新增引用錯誤。
- 二十五個 C 函式來源身分、分級語意、原版 EXE／DAT 雜湊與輸出擁有權均核對。輸出 UID/GID 為 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出並移除。

## 2026-10-08：完整 C 月結規則與政治／俘虜依賴

- 前輪 `5402deb` 已推送，分類為 progress。中斷前只讀取依賴，恢復後確認工作樹乾淨、沒有存活容器，再沿逆向／IDA／文件職責契約續作。
- 二十一個原始函式的 1508 bytes、operand、chunk、xref 與固定 IDA DB 匯出，形成政治／俘虜依賴閉包。spec/211 READY 後實作 C。
- 九個月結尾端全部以真實 C 接線，含政治、俘虜、queue 月度初始化、scratch 邊境去重排序與 writer。僅通知、音效、重畫仍是明示 RET fixture。
- O0/O2 各 11,632 組原版/C 相同，各版 91 次完整記憶體核對相同，含四個原版劇本 × 三個玩家選擇。新政治／俘虜原版證據沒有外推完整 Go 月結 parity。
- 八個 writer、relation war bit、queue 搬移、任官、力量、timer、排序及協力門檻突變均被拒絕。協力門檻曾通過原先 80 個政治樣本，查明沒有滿足深層前提；新增八個相鄰門檻後，正常版本通過、突變拒絕。
- 收據驗證器舊計數 11,624 在新增案例後拒絕新收據，按實際新增群組與 source identity 更新為 11,632，沒有改來源身分或降低比較判準。
- 現行 C 台帳加入二十一函式，共四十六個，每筆保留原始位址、routine／source hash 與 proven／re/97 出處。正式 Go 原始碼本輪沒有變更。
- 本輪依授權完成後 commit、push；完整 matching decompilation Goal 保持 active，月結規則接線不代表整檔 C 或正常玩家 UI 完成。
- 收尾文件與工具檢查 21／25 通過，860 列分流、嚴格研究索引、過期斷言、新工具語法、校訂、資產 deny-list 與全部正對照通過。四項既有缺檔／教訓檢查失敗保留，新文件沒有新增引用錯誤。
- 四十六個 C 函式 source identity 與分級出處核對，原版 EXE／DAT 雜湊相同。輸出 UID/GID 為 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出移除。

## 2026-10-08：完整月結原版／C／Go 36 向量比較

- 前輪 `e2fe24f` 已推送，分類為 progress。核對現況、逆向／IDA／spec／文件職責路由，將已還原 C 呼叫鏈用來直接審查 remake。
- 從正式 tick 純抽取 153 行原有月結規則至 private monthlyRules，正常 tick 的前置據點／軍團／時鐘與最後每時更新次序保留。MatchingMonthly 僅在 matching build tag 可見，沒有 production shortcut。
- 同四個原版劇本、玩家 0／7／21、raw c=0／77／255 與 s=uint8(c×37)，原版/C 每向量全暫存器及 1 MB 相同。Go 在同一 rules 邊界比較全部 22,208 bytes 與 258-byte RNG，所有載入基準零差異，沒有 mask 或剔除 queue。
- 最初 0／36 完全相同，general +0x1F 共 3429 bytes 差異。按 READY spec/213 接回 typed score、月結 byte 計算與保存，9／36 通過。
- 原始交換式排序修正依 READY spec/214，12／36 通過；中立 FF18 producer 依 spec/215，20／36 通過；storm globals 與玩家宣戰 gate 依 spec/216、217，最終 36／36 完整區塊與最終 RNG 相同。
- 各 before-score／after-score／after-order／after-neutral checkpoint 的 audit、vectors 與 source hashes 保留。本機原版區塊與 RNG binary 不進 Git。
- Storm globals 使用 typed 更新欄位，Bytes 套到原版位置，RawBlock 原始錨點保持不變。新增原版格式 round-trip、停用／終止槽、byte wrap、同值尾端、FF 邊界、中立非 producer 與玩家規則冷測。
- strategyai／state 正常與 matching tag 冷測、Go vet 通過；驗證器 globals／勢力／據點／武將／queue／RNG 六種單 byte 突變全部拒絕。
- 首次排序實作以 byte 暫存 Raw() 的 int，型別檢查拒絕，改用正式 API 型別後同一向量重跑。沒有放寬原始比較或挑 seed。
- spec/212–217 更新 CONFORMED，範圍限這 36 固定月結 rules 向量，不能當作任意後期世界、正常 GUI／音效／長程、完整 C 或 C 機器碼完成。Goal 保持 active，本輪依授權 commit、push。
- 收尾文件與工具檢查 21／25 通過，885 列分流、嚴格研究索引、過期斷言、校訂、資產 deny-list 與全部正對照通過。四項既有缺檔／教訓檢查失敗保留，新文件與研究入口沒有新增引用錯誤。
- 36 向量 source identity、原版 EXE／DAT 雜湊、46 個未修改 C 函式來源與產物擁有權核對。輸出 UID/GID 為 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出移除。

## 2026-10-08：八個每時更新與 dispatcher C 核心

- 前輪 `8251289` 已推送，分類為 progress。開工核對乾淨工作樹與現況，命中逆向／IDA／文件職責路由，沿時鐘的每時下游續作。
- 原始八函式共 445 bytes，IDA 9.4 probe 保留原始名稱、位址、operand、chunk、xref 與固定 DB。spec/218 READY 後以 C 接線財政、維持費、外交與原始 event table。
- event 10／13 和 trust helper 為真實 C，其他 11 event handler 與 UI／音效／退出／重畫仍為明示 RET fixture，沒有額外 Alive gate 或更換表格。
- O0/O2 各 21,871 組原版/C 相同，各版 171 次全 1 MB 記憶體核對相同，含全部 cadence byte、13 碼及空事件、四劇本每時入口、外交與支出／兵池邊界。
- 初次 trace 的暫存器與資料已一致，根因是 event 10 位址重複登錄、sub_1310A 沒登錄。按位址去重並補正式依賴後，用同樣矩陣、原版與 source 重跑，沒有放寬 trace 判準。
- cadence 突變被原始 comparator 拒絕，但失敗 JSON 沒保存真正不同的 globals／cadence 欄位；補齊 schema 並重生收據，validator 才能機器核對差異。
- 支出上界、carry、cadence、cursor stride、trust borrow、次序與外交 byte 經費七個突變全拒絕。正式 Go 本輪沒有修改，完整 Go 每時對拍與其餘 handler 尚待後續。
- 現行 C 台帳加入八函式，共五十四個；每筆 source／routine hash、原始位址與 proven／re/99 出處保留。Goal 保持 active，本輪依授權 commit、push。
- 收尾文件與工具檢查 21／25 通過，893 列分流、嚴格研究索引、過期斷言、校訂、資產 deny-list 與全部正對照通過。四項既有缺檔／教訓檢查失敗保留，新文件沒有新增引用錯誤。
- 五十四個 C 函式來源身分、分級語意、原版 EXE／DAT 雜湊與產物擁有權均核對。輸出 UID/GID 為 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出移除。
