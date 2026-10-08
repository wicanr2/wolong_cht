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

## 2026-10-08：十三碼事件 handler 與政治／災害 C 閉包

- 前輪 `9f5a843` 已推送，分類為 progress。開工核對乾淨工作樹、現行狀態、逆向／IDA／文件職責路由與 GitHub Issue #27。
- 三十個原始函式共 1919 bytes，IDA 9.4 probe 保留名稱、linear address、operand、完整 chunks、xref 與 DB 身分。spec/219 READY 後實作 C 原始名稱入口。
- 十一 handler、event 8 落入 `sub_133FD` 的原始尾端及十八個政治／災害依賴均真實接線。沒有為落下尾端插入 CALL／RET；modal、金額、戰鬥準備、UI 仍是明示 fixture。
- O0/O2 各 43,476 組原版/C 相同，各版 340 次完整 1 MB 核對相同；三十函式都有實際入口次數，包含四劇本原始 dispatcher 與受控 queue／四種回應。
- 門檻、關係 SHR、raw 0x24、災害容量／timer、軍團分支、原始 count 比較與協力上界八個突變全拒絕。補齊 C 原始名入口後，以最終來源完整重生收據。
- 原版與 C 的狀態摘要分別讀各自記憶體，仍相同。正式 Go 本輪沒有改動，中立 byte 行為沒有升為玩家語意，Issue #27 保持 OPEN。
- 更正 spec/218 的每時矩陣為四劇本獨立入口，原收據保留。較早 spec/218、re/99、spec/64 的範圍 backlink 由新 verifier 核對。
- 現行 C 台帳加入三十函式，共八十四個；source／routine hash、位址、實際入口及 proven／re/100 出處保留。完整 Goal 保持 active，本輪依授權 commit、push。
- 收尾檢查的 DB 雜湊少了明確檔案路徑，stale_scan 誤配到上一個 Markdown 連結；補上 DB 路徑後通過。兩個自我測試需要 workplace 暫存案例，改掛有界 tmpfs 後以同一工具鏈重跑通過，首次環境失敗 log 保留。
- 收尾腳本用 basename 去比完整路徑，漏寫 docs/re/43；改為明確重生，並正查新兩份文件都已收錄。文件與工具檢查 21／25 通過、901 列分流與嚴格索引通過，四項既有缺檔／教訓失敗保留，本輪沒有新增引用錯誤。
- 八十四個 C 函式來源身分、三十筆分級出處、原版 EXE／DAT、實際入口與 backlink 均核對。輸出 UID/GID 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出移除。

## 2026-10-08：十七個外交／金額視窗 C 函式

- 前輪 `b5ad756` 已推送，分類為 progress。開工核對乾淨工作樹、現行狀態、逆向／IDA／視窗／文件職責路由與 GitHub Issue #22。
- 十七函式共 1280 bytes。第一次直接比 IDA 與 file chunks，在二十個 FAR segment word 失敗；確認全都來自 MZ loader 加 paragraph。保留舊 probe／DB，按原始 116 筆 MZ 表逐 word 記錄兩種 bytes 與 runtime paragraph，沒有遮罩差異。
- IDA 9.4 probe 審查通過、spec/220 READY 後實作 C 原始名入口。真正的 FAR／RETF、原始 SS frame、PUSHF／POPF、CF 輪詢／取消重試、RNG／trust、金額／扣款、code word 修改及恢復已接線。
- 受控裝置 FAR leaf 會 clobber FLAGS，讓 PUSHF／POPF 成為可拒絕的 gate。選單忙碌與數值取消分開，原始 caller 迴圈真實執行，底層繪圖／輸入／音效仍為明示 primitive fixture。
- O0/O2 各 53,322 組原版/C 相同，每版 417 次完整 1 MB 核對相同；十七函式皆有實際入口、252-byte callee 快照、控制欄位與中途 code word 證據，十個突變均拒絕。
- 收尾抽查查出訊息向量寫到 world bank，實際讀取端是 SS:BP。補 typed message、改寫真正的 SS frame 後，以完整矩陣與全部負對照重生，舊來源與收據保留在 c-modal 的 pre-message 檔案。
- 較早 spec/219、re/100 範圍 backlink 由 verifier 核對，另十一份同位址規格逐份確認不受影響。正式 Go 未改，原始未初始化 stack 槽不補值、raw 大數值不冒稱玩家輸入，Issue #22 保持 OPEN。
- 現行 C 台帳加入十七函式，共 101 個。完整 Goal 保持 active，本輪依授權完成後 commit、push；本輪不宣稱完整視窗／正常玩家／Go 三方或 C 機器碼完成。
- 收尾文件與工具檢查 21／25 通過，909 列分流、嚴格研究索引、過期斷言、校訂、資產 deny-list 與全部正對照通過。四項既有缺檔／教訓失敗保留，新文件沒有新增引用錯誤。
- 101 個 C 函式來源身分、十七筆分級出處、最終 driver／收據、原版 EXE／DAT、重定位與擁有權核對。輸出 UID/GID 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出移除。

## 2026-10-08：數值輸入器與財政 caller 的真實 C

- 前輪 `7bdafd3` 已推送，分類為 progress。開工核對乾淨工作樹、現行狀態、逆向／IDA／spec／文件職責路由及 GitHub Issue #22。
- 十七函式共 598 bytes，IDA 9.4 probe 保留原始名稱、完整 chunks、operand、xref、file／IDA bytes、三筆重定位、CS:7D01 間接 table 與 CS:7D93 十八格。spec/221 READY 後實作 C。
- 真實六鍵、主迴圈、LAHF／SAHF、word MUL／carry／cap、DIV model、CLD／LODSB、裝置保護／popup wrapper 與財政四 caller 接線。圖形／文字／輸入／popup leaf 為明示 primitive。
- 新 driver 少一個 scenario 外層括號，補齊後編譯器又指明 goto 跨 plans 宣告；回查路由，將宣告移到所有跳躍前，用同一工具鏈重跑。FAR fixture 跳距多一 byte，誤入未覆寫的 code；按指令邊界修正，沒有把它當成 CPU 或產品缺陷。
- 完整 O0/O2 各原版/C 788,931 組、Go 有效數值 787,839 組相同，每版 6164 次完整 1 MB 核對相同，十七函式都有實際入口。Go 使用既有正式 API，原版 CF 與 Go action-valid 不混稱。
- 第一版跳表錯版忽略完成鍵，讀到第 513 筆非宣告事件才返回；加入明示尾端取消並完整重生，最終錯版由第二筆事件取消，verifier 正查 index=2／cancel=1。舊來源／收據留在 c-numeric 的 pre-terminal-cancel 檔案，正常矩陣與判準不變。
- 十個 carry／乘百／除十／最大／FLAGS／glyph step／table／SAHF／取消寫回／徵兵單位突變全拒絕；六個原始 handler 都有真正 table 主路徑向量。正式 Go 原始碼未改，圖形／自然玩家／完整 Go transaction 與 C 機器碼仍未驗證。
- spec/220、re/101 與 spec/78 的同版範圍 backlink 由 verifier 核對。現行 C 台帳累計 118 個，Go 欄只記每函式直接 scalar 比較；完整 Goal 保持 active，本輪依授權完成後 commit、push。
- 收尾文件與工具檢查 21／25 通過，917 列分流、嚴格研究索引、過期斷言、校訂、資產 deny-list 與全部正對照通過。四項既有缺檔／教訓失敗保留，本輪沒有新增引用錯誤。
- 118 個 C 函式來源身分、十七筆分級出處、最終原版/C/Go 收據、原版 EXE／DAT、重定位／終止輸入與擁有權核對。輸出 UID/GID 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出移除。

## 2026-10-08：熱區 map／query 與數值 pixel 接線

- 前輪`99d914b`已推送，分類為progress。核對乾淨工作樹、現況、逆向／IDA／GUI／文件職責路由與Issue #22。
- 十函式491 bytes，IDA9.4保留原始names／chunks／operand／xref／file bytes，無重定位。Writer／reader與舊re/22、re/47證據確認0x1E3D7為熱區，不是glyph，spec/222 READY後實作C。
- 原始map init／register／clear／query、window wrapper、tile flags與border caller真實接線。輸入改事先固定CX／DX pixels／cancel，從十八個raw key table定位，不直接覆寫query回傳AL。
- 生成器找舊模板位置失敗，核對後修正；新C的hex literal尾E與加號連為preprocessing number，補空白再跑，均為工具問題。DI=FFFFh抽樣找出word第二byte應同段wrap，C init改byte write，不移除邊界。
- Init修正後結果仍舊，回查build發現外部C include未觸發cgo cache。加入source manifest SHA-256至CGO flags，另匯出每個binary的build info並驗證，強制來源變更真正重編。
- O0/O2各261,367組原版/C、80組正式Go數值相同，各版2042次完整1 MB相同。全640×400 query、raw word／base wrap、原始0尺寸、重疊CF、flags／border、四劇本及財政／數值pixel鏈均比完整ABI。
- 八個init／high-byte／collision／pitch／tile flags／border／Y offset突變全拒絕，query丟high Y在第163,841組被拒絕，不以隨機pattern相同掩蓋。正式Go原始碼沒改，VGA primitive和正常玩家仍未驗證。
- 勘誤現行glyph導覽與FFF8h word公式，保留舊群組名、operand與歷史收據，取消spec/78的一列假glyph未讀前提；其餘分流保留，backlink由verifier機器核對。
- C台帳累計128函式，完整Goal保持active；本輪依授權完成後commit、push，不把map／scalar結果稱完整UI或C機器碼匹配。
- 收尾文件與工具21／25通過、924列分流與嚴格索引通過；四項既有缺檔／教訓失敗保留。最終研究改寫後四列source fingerprint按實際來源重審，檢查容器cwd補為/repo後同一矩陣重跑，沒有放寬分類。
- 128個C來源、十筆分級證據、最終ABI／map／pixel／compiled source digest與原版EXE／DAT身分核對。輸出UID/GID 1000:1000，既有root-owned路徑仍27筆，本輪容器已全部退出移除。

## 2026-10-08：八函式VGA四plane與保存／恢復

- 前輪`a8605b5`已推送，分類為progress。核對乾淨工作樹、現行狀態、逆向／IDA／GUI／平台規格／文件職責路由與Issue #22。
- 八個原始chunk467 bytes，IDA9.4無relocation；平台probe確認MZ loader與兩個獨立plane／latch狀態。標準VGA引用IBM規格，遊戲RE只追原始port參數、source連續性、dummy read與loop。
- C算法保留原始register／FLAGS／DF／SS、MOVSB／REP、single/double byte列、GC序列與保存major layout；bus callback到第二台獨立成熟VGA。不讀原版結果回填C、不用平面RAM假證明VRAM。
- 新Go工具goto跨configs宣告，移到所有跳躍前以同一矩陣重跑；smoke數量8截斷第三次roundtrip，改為9保留完整三步，沒有改保存／恢復算法。
- O0/O2各9,942例全部RAM／四plane／GC／seq／latch／index／port序列與原始register相同，310次640×400 indexed pixels相同，18條save/draw/restore全plane回到保存前。
- 四個原始2KB頭像×兩destination，原版／C本機debug-palette PNG逐byte相同；來源檔count／大小／SHA核對。自動審查拒絕PNG base64外送檢視，理由可能披露未授權美術；未執行或繞過，改以本機技術收據完成驗證。
- 原版port log的Step保留定位，C不宣稱相同CPU instruction count／wall-clock。所有CGO build flags含source manifest digest，由實際binary build info核對，研究檔用matching_vga tag隔離後完整重生。
- 八個dummy latch／列寬／pitch／OR／plane／read map／wrapper X／DF突變全拒絕，原版資料與衍生圖沒有進Git。正常玩家／整體UI／完整Go圖形與C機器碼仍未驗證，Issue #22保持OPEN。
- re/103、spec/222與re/03加同版後續範圍backlink。C台帳累計136，完整Goal保持active，本輪依持續授權完成後commit、push。

## 2026-10-08：還原程式碼加入 GitHub 紀錄

- 使用者要求把驗證通過的 matching decompilation 程式碼放進現有 GitHub repo。C 來源已逐輪版控，本輪另加入完整組語指令來源與 linker script。
- 沿用 re/90 與既有研究目錄，新增來源索引、公開驗證摘要與 Docker 重跑入口。所有指令保留原始 IDA 定位，沒有加入資料宣告，原版資料區只從本機自備 EXE 匯入。
- 版控來源冷組譯 24,376 條／55,392 bytes 全部匹配；匯入 11,707 個非指令 bytes 後，完整 67,099-byte EXE 相同。指令、linker 常數與私有輸入身分三個負對照均拒絕。
- README、CONTEXT、re/90、re/104與研究索引同輪連結來源、台帳與驗證紀錄；原版 EXE、DB、圖庫、RAM、plane 與 PNG 留在 ignored 本機研究目錄。
- 只掛載版控指令來源、索引與自備 EXE 的第二個容器同樣冷重建通過，不依賴本機 DB 或含資料的組語。公開摘要與實際收據逐 byte 相同。
- 文件與工具檢查 21／25 通過，932 列分流與嚴格索引通過。四項既有失敗仍是歷史圖片／研究產物缺檔、教訓文件不同步與一條教訓文字缺失，沒有新增來源引用問題。
- Docker stdin 檢查首次未開 interactive 而沒有執行，補 `-i` 後同一檢查矩陣有完整 probe。一般 tools 套件本來沒有非研究 Go 檔，改用 `go list -e` 核對新兩檔確實被排除；格式檢查通過。
- 最終核對全部 136 個 C 來源、八筆 VGA 分級、組語自動合併的語意／原始運算元、公開與本機收據。原版與來源雜湊不變，輸出 UID/GID 1000:1000，既有 root-owned 路徑仍 27 筆，本輪容器全部退出移除。

## 2026-10-09：五個位元對齊與戰術按鈕 C 函式

- 前輪 `47e16c9` 已推送，分類為 progress。核對乾淨 main、現況、逆向／IDA／平台／文件職責路由與 Issue #22，沒有重開已還原的外框函式。
- IDA9.4 十八函式原始 chunks／xref／operand 與三筆既有 MZ 重定位核對，新五函式共 429 bytes。原始位置／hit 表另核對，契約 READY 後才實作 C。
- 新 driver include 首次用了不存在的 world.c，核對現有名稱後改 world_update.c，以同矩陣重跑。C 首次用連續實體 word 讀取，SI=FFFFh 反例只有 plane 不同；回查 CPU.read16 後修正同段回繞，保留反例與原邊界。
- FAR fixture 首次漏 CS 段覆寫，原版與 C 的裝置 restore 分支分歧；補 fixture 契約後同矩陣重生，沒有改遊戲 caller。原始失敗收據與編譯來源摘要留在本機。
- O0/O2 各 4,954 組全部 RAM／四 plane／GC／seq／latch／port／十四暫存器／FLAGS／入口快照相同，683 次內容區像素、36 次 hit 消費通過。六按鈕來源連續、單一重畫與外框直接接真實 C／VGA。
- 十一個首／尾遮罩、dummy read、來源步幅、word 交換、平面重設、JL／OF、六圖重設、位置表、DF 與外框 pitch 錯版全拒絕。ROW 遮罩不再與忽略資料的 mode 1 關聯，不放寬判準。
- 五個 C 來源與分級證據加入台帳，累計 141 個。re/18、re/103、re/104 與 spec/223 回填新範圍入口，原始圖片／收據保留。正式 Go 規則、UI 與存檔未改，正常玩家與 C 機器碼仍未驗證，完整 Goal 保持 active。
- 更新版控組語的五函式分級註記後，24,376 條／55,392 bytes 與完整 67,099-byte EXE 冷重建仍相同，三個負對照拒絕；來源與公開收據身分核對。
- 收尾文件／工具 21／25 通過，940 列分流、嚴格索引、來源引用與研究 tag 隔離通過。四項既有缺檔／教訓問題保留，沒有新增引用錯誤。
- 全部 141 個 C 來源、五筆分級、原版／資產與編譯摘要核對。輸出 UID/GID 1000:1000，既有 root-owned 路徑仍 27 筆，本輪容器已全部退出移除；依授權提交並推送來源與證據。

## 2026-10-09：十六矩形／選取／計量 C 函式與 Go 校訂

- 前輪 `6139f9e` 已推送，分類為 progress。核對現況、乾淨 main、逆向／IDA／文件／平台路由與 Issue #22；直接沿用 renderer／VGA adapter。
- 新 16 函式 1,169 bytes、36 原始定位、四筆 MZ 核對；spec/225–226 READY 後才實作 C 與正式 Go 修正。兩次右移與取 byte 原樣保留；體力 3=1、1024 兵力回繞與 9999 體力=74 成為相鄰反例。
- 舊素材矩陣給單一命令／外框第 1 段資料；依 sub_100DF 配置與圖庫實際 47,776 bytes，補第 3 段九命令圖與 0x6C0 框 bank。舊 raw receipts 保留，新的實際素材 54／18 組通過。
- 本機原 go.sh 預設 SDK 缺失，沿用現成 hr-go-ebiten image；僅在容器下載 go.mod／go.sum 鎖版依賴，Xvfb 有界 trap。正式 sidebar／selection／layout 冷測通過。
- O0/O2 各 134,958 全 RAM／四 plane／I/O／ABI 相同，131,072 原版／C／Go 長度與 2,441 次內容區相同；十二錯版全拒絕。正式 Go 函式與上限常數逐字擷取，摘要進實際 binary flags。
- IDA 資料項 mnemonic 空白造成窄 decoder 首次無法停止，改用 decoded insn 的 canonical mnemonic；原資料行與分類保留。Opcode 02／03 八條／20 bytes 合併到完整版控組語，新 24,384 條／55,412 bytes 冷組譯、整檔 67,099 bytes 相同，三個負對照拒絕。
- C 台帳累計 157，完整 Goal 保持 active。新的來源與原始位址／分級／收據進 GitHub；正式計量只依原版校訂，規則／存檔格式未改，自然戰術與 C 機器碼仍未完成。
- 收尾文件／工具 21／25 通過、950 列分流與嚴格索引通過。原始雜湊曾被掃描器誤綁到同段工具連結，分開 DB／manifest 檔案段落後過期斷言檢查通過；四項既有缺檔／教訓問題保留。
- 正式 wlgame 全套 `go test -count=1` 在 SDK／Xvfb 通過。最終 157 個 C 來源、16 分級、8 追加組語、公開／本機收據、compiled source 與擁有權核對；輸出 UID/GID 1000:1000，既有 root-owned 路徑仍 27 筆，容器全部退出移除，依授權提交並推送。

## 2026-10-09：顯示清單九 C 函式與十七原始入口

- 前輪 `c7dc8df` 已推送，分類為 progress。核對現況、乾淨 main、逆向／IDA／文件與 Issue #22，沿用獨立成熟 VGA，還原原始登記／DS table 與九個 opcode。
- 臨時 DB 解碼先保存原名稱／資料行／bytes，遇 ret 後 word 資料與切在前綴中間的邊界後回查寫入／分支；保留 1EE62 與 1F45B 資料，完整 code 不以資料硬解。
- 固定 IDA 832 指令產生 native C 具體運算與 label，不讀 opcode／不呼叫 guest CPU。明示 MUL／DIV operand 修正，原始像素 EFEF 與舊矩形 helper 給兩側對称擷取；有限字串 fixture DS 改讀宣告 bank。
- 十個原始場景與九 opcode、替代 DS／交換 table、線段／底紋／雙色／文字迴圈 O0/O2 各 382 全 RAM／四 plane／I/O／ABI 相同。X／Y 相同使錯源突變不能拒絕，改不對称輸入後同矩陣重生，十個錯版全拒絕。
- 新增九個原始 named C 函式與 17 原始 code 入口，函式界線導覽與無函式 data 項分開記錄，glyph raster 還是明示 fixture。正式 Go 未改，完整 Goal 保持 active。
- 317 條／746 bytes 未分類 code 經實際組譯匹配，來源同輪進 GitHub。合併為 24,693 條／56,138 bytes，完整 67,099-byte EXE 與 SHA 相同，三個錯版拒絕，資料槽維持原始 bytes。
- 生成 C 在乾淨容器從同一 IDA 證據逐 byte 重生相同。完整 C 來源摘要、166 named 函式／17 code 入口、317 補充指令與公開／本機收據核對；只保留 glyph raster 的明示邊界。
- 收尾文件／工具 21／25 通過、958 列分流與嚴格索引通過，四項既有缺檔／教訓問題保留。舊 verifier 的補充範圍改為原兩個 handler 子集，修正縮排後全部 AST／語法通過，未改其原矩陣。
- 研究 driver 以 matching_display 隔離，輸出 UID/GID 1000:1000，既有 root-owned 路徑仍 27 筆，容器已全部退出移除；依授權提交並推送原始來源與收據。

## 2026-10-09：原始 C glyph raster 與三字名稱

- 前輪 `4d31520`已推送，分類progress。核對現況、逆向／IDA／文件／字型平台路由與Issue #22，目標為移除 glyph raster no-op。
- 原始C四named函式282bytes與兩入口157bytes由IDA9.4核對，三字caller最初遺漏尾端epilogue後補完整邊界。spec/228 READY後產生C，所有原位址、far與self operand保留。
- 兩側各DOS／Machine／Font cache；C adapter只呼叫原平台IntHook，沒有CPU.Step。原始F720取得向量、32-byteSSbuffer、字庫讀取與真正VGAraster皆執行。未用import移除後同輸入乾淨重跑。
- O0/O2各252組全RAM／四plane／latch／port／ABI相同，每組取內容區像素。兩側740全形／150半形、缺字0，八種對齊、上界、透明／背景、三字live字色與十場景覆蓋，九錯版全部拒絕。
- 服務stub早先導覽名錯寫成sub_10410，改真正0080:0410／0414再重生所有最新收據，KI位址與平台位址不混用。固定生成C在乾淨容器逐byte再生一致。
- 組語補充保留先前317指令，再追加21／59bytes全部實際組譯相同；總24714／56197與完整67099-byteEXE匹配，資料槽、字庫與原版不進Git。
- C台帳170函式，19個原始code入口；glyph fixture已取消，正式Go未改，自然UI、原TSR硬體時序與C機器碼仍未完成。完整Goal保持active。
- 最終公開摘要在所有錯版程序結束後重新綁定最新收據，四C來源／兩入口與真實字庫／compiled flags 逐項核對，生成C在乾淨容器重生相同。
- 文件／工具21／25通過、965列分流與嚴格索引通過，四項既有缺檔／教訓問題保留。舊display verifier保留原317指令子集，新增glyph指令不破壞歷史檢查條件。
- 170個C來源／19code入口、24714組語與整檔、原版／字庫、source摘要及擁有權核對；輸出UID/GID1000:1000，既有root-owned路徑仍27筆，本輪容器已全部退出移除，依授權提交與推送來源證據。

## 2026-10-09：數字 C 還原與 GitHub 紀錄

- 前輪 `85f06de` 已推送，核對乾淨工作樹、復古逆向／IDA／文件與 backlink 路由，讀取 Issue #22 現況。沿用每輪 commit／push 授權，完整 Goal 保持 active。
- spec/229 READY 後還原 `sub_1062F`、`sub_1069A`、`sub_106DE`、`sub_11E17` 與原始無名稱 `0x10984` 入口。固定 IDA 143 指令生成 native C，字庫取實際 ICONGRF 第三段加 `0x840`。
- O0/O2 各 4,881 組完整 RAM／四 plane／GC／seq／latch／port／ABI 與內容區像素相同。八個錯版全部拒絕；產生器在乾淨容器重生相同，source digest 綁定實際編譯 flags。
- 新收據限數字 raster、日期與 SS 參數入口，較早 editor 收據仍保留原明示邊界。第一個 DIV 限制、原作者工具鏈與自然玩家 UI 不猜補，正式 Go 本輪沒有改動。
- 分級語意、四函式／一 code 入口與公開驗證摘要進版控，台帳目前 174 函式／20 code 入口。組語分級註解重生後，24,714 指令與完整原 EXE 相同，三個組語錯版拒絕。
- 文件與工具 21／25 通過，971 列分流與嚴格問題索引通過。新增規格份數已回填 README；四項既有缺圖／研究產物／教訓同步問題維持基線。Go package 檢查補上容器 HOME 與 cache 後同命令通過，研究 driver 正常被預設 build 排除。
- 公開與本機收據、174 函式／20 入口來源、整檔組語、原版字庫、source digest 與 UID/GID 逐項核對。既有 root-owned 路徑仍 27 筆，沒有新 root 輸出；本輪容器已退出移除，依授權提交並推送來源及驗證紀錄。

## 2026-10-09：TALK、參數與原始肖像載入

- 前輪 `06fff77` 已推送，分類 progress，核對乾淨工作樹、逆向／IDA／backlink 路由與 Issue #22。依既有 commit／push 授權繼續完整 Goal，不把本輪 renderer 完成當成全專案完成。
- 八個 named C 函式與六個原始 handler 共 743 bytes／376 指令，spec/230 READY 後接到真實字庫、數字、背景與肖像。保留原 near table、word offset、CF、LAHF／SAHF、原始四格替換與 DOS open／seek／read／close。
- IDA 探針原先只按 basename 找出處，上一輪 spec/229 使匯出失敗。改完整 docs 唯讀掛載與完整路徑 resolver，同名錯目錄與越界拒絕通過，再從原 EXE 建一次性 DB。
- driver 少了入口計數後補回，status 樣本的 AX personality 與 speaker 分離，避免越界索引。初期 corpus 格／像素單位混用讓 glyph 被裁切，字庫計數暴露假完整；改格座標 0，停止舊工作並以同 wrapper 全部重跑。
- O0/O2 各 1,277 組全 RAM／四 plane／ABI／port／DOS API／像素相同，1,022 槽各一次、150 肖像、15,287 全形／178 半形、缺字 0。十個錯版與裁切工作量護欄拒絕通過，生成 C 在乾淨容器重生相同。
- 正式 Go 保持既有實作，INT50 媒體 UI、自然事件／producer 與 C 機器碼沒有擴大完成聲明。分級語意與公開來源台帳目前 182 函式／26 code 入口；整檔 24,714 組語與 67,099-byte 原 EXE 仍匹配。
- 較早三份訊息／快取文件回填可驗 backlink；直接位址 consumer 的「五支」改為三支，數字格式移出已解缺口，其餘 producer／屬性限制保留。976 列分流與嚴格問題索引通過，文件與工具 21／25 通過，四項既有缺檔／教訓問題維持基線。
- 最新公開收據在驗完六個唯讀素材與 resolver 拒絕對照後重生，182 函式／26 入口的來源 binding、生成 C、編譯 flags、所有八個新入口與完整 RAM／plane 收據相符。研究 driver 預設 build 排除，所有輸出 UID/GID 1000:1000，既有 root-owned 路徑仍 27 筆。
- 舊裁切工作已停止且刪除，最後全矩陣與每批工具容器皆已退出移除；依授權把還原來源、原始定位及驗證紀錄提交並推送 GitHub，完整 Goal 保持 active。

## 2026-10-09：訊息等待與完整 mouse segment

- 前輪 `52c2b09` 已推送，分類 progress，核對原始現況與逆向／IDA 路由後接訊息顯示／等待／擦除。19 named 函式與兩 raw 入口共 957 bytes／418 指令，spec/231 READY 後還原。
- mouse 遠呼叫與 16 項原 table、cursor 保存／恢復／兩層繪圖、CLI／STI／FLAGS 及 callback 為真實 C。INT33 用獨立平台 API，C 不跑原 guest，timer 以明示 counter checkpoint 輸入，不外推硬體 cadence。
- 原 `EB 0F` 跳過左鍵但問右鍵及 counter，NOP 模式加入左鍵，兩者均有十個 counter 間隔 timeout。spec/45、re/42、re/25 與 mechanics/15 回填分級勘誤，正式 Go 沒有猜改。
- Driver 原先錯認 CPU 有 InHook；查實 Bus 接口後用只記錄 IN 的委派 bus，其他記憶體／OUT 契約不變。生成 C 的 NULL 尾端與 raw callback／patch 保留原邊界和實際 CS。
- O0/O2 各 210 全 RAM／四 plane／IN／OUT／mouse 狀態／ABI 相同；三個背景還原、十九 named／兩 raw 入口與十二錯版通過，生成 C 乾淨重生。C 台帳 201 函式／28 raw 入口。
- 組語新增21指令／71bytes與三個far重定位，比對 file bytes 與 IDA loader bytes。匯出器舊補充數量約束修正後重生，24,735指令／56,268codebytes與整檔相同，三組語錯版拒絕；舊glyph verifier只驗原338指令子集，保留歷史矩陣。
- 收尾讀取 Issue #22 現況，仍為 OPEN，未以局部 C 換成正常玩家完成聲明。982 列分流、嚴格問題索引、生成 C／source flags 與所有來源 binding 通過；文件與工具 21／25，四項既有缺檔／教訓問題維持基線。
- 六個唯讀資料／字庫身分重驗，201 函式／28 入口與公開／本機收據相同。舊 glyph verifier 在隔離暫存輸出通過，不覆寫其歷史收據；本輪所有輸出 UID/GID 1000:1000，既有 root-owned 路徑仍 27 筆。
- 全矩陣與每批 Docker 容器皆已退出移除；依授權提交並推送來源、原始定位與驗證紀錄，完整 Goal 保持 active。

## 2026-10-09：selector 與原進言 GUI C 接線

- 前輪 `95b951a` 已推送，分類 progress，核對原始狀態、逆向／IDA 路由和 Issue #22 後接 popup 選列、捲動與 cursor 保存還原；spec/232 READY 後還原 14 named 函式和兩 raw 入口，474 指令 native C。
- 原 `0x9441`／`0x9443`／`0x9465`／`0x9467` live operand與 `0x36D` row callback保留。map／pixel兩分支與原 `sub_193E9`／`sub_13B7E`、speaker／advisor caller接入真正字形／mouse／等待C；原 nullsub_5 照原image，沒有猜補。
- Prototype未初始化CS:987C畫面保存段，原版保存覆蓋服務區導致超budget；查原sub_19796後設定7200獨立段，用同wrapper和輸入重新跑。沒有調高steps或以fixture代替原分支。
- O0/O2各138完整RAM／四plane／IN／OUT／Mouse／ABI／像素相同，三個XOR雙切還原、2153全形／6半形、缺字0；十二錯版全拒絕，生成C在乾淨容器重生相同。原版正常玩家與scroll helper動態改寫仍有獨立界線，正式Go不改。
- 原re/84錯置file／IDA位址與Y回置值已按完整原branch和固定回播勘誤，保留原斷言形成背景。新增10組語／19bytes，先前369以前的補充原記錄保留，24,745指令／56,287codebytes與完整EXE匹配。
- 全錯版終止後獨立驗證最新完整來源與 flags；987 列分流與嚴格索引通過，re/84 已閉合缺口加明示無缺口標記。文件／工具21／25通過，四項既有缺檔／教訓問題維持基線；gofmt與預設build排除通過。
- 舊glyph與input verifier在隔離暫存輸出通過，原始歷史收據不覆寫；215函式／30raw入口binding、六個原始唯讀素材、完整組語／C收據與擁有權核對。所有輸出UID/GID1000:1000，既有root-owned路徑仍27筆。
- 最後全矩陣與每批工具容器皆已退出移除，依授權把原始定位、C來源與驗證紀錄提交並推送GitHub；完整Goal保持active，未把局部consumer取代正常玩家驗收。
