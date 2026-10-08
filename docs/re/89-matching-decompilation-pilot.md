# 89：比對式反編譯試點：三個函式可以重組回原始機器碼嗎

**狀態：局部組語重組與 C 編譯試驗完成。結論限於三個函式。**

- 日期：2026-10-08
- 範圍：松崗 DOS/V `KI.EXE`；`LOGO.EXE` 只作工具鏈掃描正對照
- 專案來源：`e8aeb99`，本試點沒有修改 Go 遊戲規則
- 輸入 `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 正對照 `LOGO.EXE` SHA-256：`2d36530b998888da0f35cef843c38daebb957905c58a5102ca5bbe9bf1baaeef`
- 一次性 `.i64` SHA-256：`5c1ae91a175fc28a8f6276b1df15326dc705227cd2df8c50e7743b5c5acd4a17`
- 工具：IDA Pro 9.4、Python 3.12、GNU assembler／ld 2.40、GCC 12.2.0
- 位址空間：以下函式位址使用 IDA database linear address；段內位移與檔案偏移分欄記錄

本檔保存先行試點的收據。2026-10-08 後續已完成整檔組語重建，現況見
[`90`](90-assembly-reconstruction.md)；以下表格仍限於先行試點。

## 1. 問題與本次驗收

使用者希望另開比對式反編譯研究線，重建可以編譯回原始二進位的程式碼，再審查 Go remake。
本次先檢驗三個代表性函式，沒有把整檔重建或高階原始碼還原列為已完成。

| 項目 | 本次結果 | 推論等級與範圍 |
|---|---|---|
| 三函式組語重組 | 3／3 相同，共 233 bytes，沒有遮掉任何位元組 | 已證實，局部重組 |
| 高階 C 匹配 | 三個 C 探針各試四組 GCC 最佳化選項，共 0／12 匹配 | 已證實，限所列探針與選項 |
| 原版工具鏈與語言 | 沒有足以定案的證據 | 未知 |
| 整份 `KI.EXE` | 尚未重建 | 未驗證 |
| Go remake | 已列出對照入口，本次沒有執行行為對拍 | 未驗證 |

重組來源使用實際指令，不使用 `.byte`、`.word` 或 `.incbin` 注入原始程式碼。
輸出裡因 `.org` 產生的補零不計入已還原範圍。原始 EXE、分析 DB 與重建二進位都留在
被 Git 忽略的 `workplace/matching-decompilation/`。

## 2. 三個函式的逐位元組收據

| 原始名稱 | IDA 線性起點 | 段內起點 | 檔案起點 | 長度 | 原版與重組 SHA-256 |
|---|---|---|---|---:|---|
| `sub_11D8E` | `0x11D8E` | `0x1D8E` | `0x1F8E` | 137 | `4e2b6ecf618756da0f51e5dfb29c34b502fc70ffbf2d7050c5f5f51cb1a311dd` |
| `sub_131AE` | `0x131AE` | `0x31AE` | `0x33AE` | 68 | `c9b8ef2aa1fbeebb4bc0be77b70542babe82723f6d4d8305adfe0565af930cec` |
| `sub_1ECE0` | `0x1ECE0` | `0xECE0` | `0xEEE0` | 28 | `f56e5a1d08956bc9c68476dd5245b85a7e0169d52757b5e476c32b28e1911a9d` |

三組範圍先由 IDA 的函式邊界取得，再逐指令核對原始 EXE 的檔案偏移與 bytes。
重組後以 GNU ld 解決呼叫位址，再比最終機器碼；沒有遮掉重定位運算元。

匯出自動合併受版控的 [`matching-semantic-index.json`](matching-semantic-index.json)。
索引綁定原始 EXE 雜湊，只收既有已分級證據。原始函式名、運算元與位址保留，
附加語意另列等級與來源；未收錄或未證實的語意顯示警示。

GNU assembler 的 `{load}` 選擇暫存器指令的編碼方向。這會影響 bytes，
所以相同助憶碼與相同暫存器還不足以當作匹配證據。

| 驗證 | 結果 |
|---|---|
| 三組原始 bytes 與 IDA 匯出相同 | 通過 |
| 三組指令完整覆蓋各自函式範圍 | 通過 |
| 重組並連結後逐位元組相同 | 通過 |
| 連結後無待處理重定位 | 通過 |
| 每個函式首 byte 改一個 bit | 三組皆拒絕 |
| 把 `sub_13E11` 的呼叫目標由段內 `0x3E11` 改成 `0x3E12` | 拒絕 |
| 把輸入 SHA-256 換成全零 | 拒絕 |

重組只涵蓋函式本體。外部函式、月份表、亂數置換表與事件跳表沒有因此獲得新驗證。

## 3. 工具鏈辨識與正對照

| 觀察 | `KI.EXE` | `LOGO.EXE` | 證據等級 |
|---|---|---|---|
| 可列印的編譯器標記 | 本次字串掃描沒有命中 | 檔案偏移 `0x7324` 有 Borland C++ 標記 | 已證實，限此掃描條件 |
| IDA 導航用 compiler ID | `129` | `2` | 導航資訊，不能單獨辨識工具鏈 |
| IDA 函式數 | 739 | 103 | 導航資訊，不是完成率 |
| 函式起點為 `55 8B EC` 或 `55 89 E5` | 0 | 22 | 已證實，限兩種 BP frame 前綴 |

`LOGO.EXE` 的標記證明它連入 Borland C++ 家族的執行期。標記的 copyright 年份
不能指定精確編譯器版本，也不能外推到 `KI.EXE`。

`KI.EXE` 的三個試點含暫存器介面、DS／ES 切換與 CS 間接呼叫。
手寫組語或混合語言仍是假說。沒有傳統 BP frame 也可能來自編譯器最佳化，
不能據此宣稱整份遊戲由組語撰寫。既有工具鏈限制見
[`reference/04` §2.2](../reference/04-first-survey.md)。

## 4. C 試驗的意義

每個函式以 `-m16` 配合 `-O0`、`-O1`、`-O2`、`-Os` 編譯。其他固定選項包含
`-ffreestanding`、`-fno-pic`、`-fno-pie`、`-fno-stack-protector` 與關閉 unwind tables。
原始大小為 137、68、28 bytes；C 產物大小與逐位元組差異保存在 `report.json`。

GCC 12.2 的 `-m16` 沿用 `-m32` 的資料與呼叫模型，只改成適合 16-bit 模式的
`.code16gcc` 輸出。它沒有自動提供原版的 16-bit 暫存器、near call 與段暫存器介面。
來源：[GCC 12.2 x86 options](https://gcc.gnu.org/onlinedocs/gcc-12.2.0/gcc/x86-Options.html)，
查閱日期 2026-10-08。

這些 C 探針只是編碼試驗。尤其事件分派器使用平面記憶體模型，沒有重建 DS／ES
與 handler 的暫存器介面，不能當成可執行或語意已驗證的 DOS 還原碼。
0／12 的結果不能排除其他編譯器、選項或來源寫法。

比對式反編譯的方法定義見 [decomp.me FAQ](https://www.decomp.me/faq)；
同類 16-bit PC-98 工作可參考 [ReC98](https://github.com/nmlgc/ReC98#building)。
兩份皆於 2026-10-08 查閱，只提供方法背景，不作臥龍傳工具鏈證據。

## 5. remake 的對照入口

| 已匹配函式 | 既有原版語意證據 | Go 對照入口 | 本次驗證範圍 |
|---|---|---|---|
| `sub_1ECE0` | [`10-rng.md`](10-rng.md) | `internal/rules/rng/rng.go` 的 `Next` | 函式本體 bytes |
| `sub_11D8E` | [`06-game-clock.md`](06-game-clock.md)、[`15-realtime.md`](../mechanics/15-realtime.md) | `internal/rules/clock/clock.go` 的 `Advance`、`internal/state/state.go` 的月結／每時更新 | 函式本體 bytes，含呼叫順序 |
| `sub_131AE` | [`15-event10-producer.md`](15-event10-producer.md) | `internal/state/events.go` 的 `takeNextQueuedEvent` 與 `dispatchQueuedEvent` | 函式本體 bytes，不含各 handler |

組語匹配讓這三組原始控制流可以從可重組來源回查。它沒有證明 Go 的整合行為，
也沒有替代 dosgolem 的同狀態與正常玩家路徑驗收。核心規則維持既有證據等級。

## 6. 未解與證據限制

| 項目 | 目前邊界 |
|---|---|
| 原版編譯器、組譯器與連結器 | 未知，沒有版本綁定的 codegen 或工具輸出證據 |
| 原始語言與 translation-unit 邊界 | 未知，不能由三個函式或零 BP frame 外推 |
| 高階 C matching | 本次 12 組探針皆不匹配，沒有達成 |
| 整檔匹配 | 尚未處理 MZ 標頭、116 筆重定位、資料區、其他函式與外部模組 |
| 原作者原稿 | 命名、註解與來源切分無法由本次試點取回 |
| Go 行為對照 | 本次只有入口對映，沒有新的行為收據 |

這些是獨立研究線的證據限制，不擴張 remake 的既有發行閘門。
現行專案工作仍以 [GitHub Issues](https://github.com/wicanr2/wolong_cht/issues) 為入口。

## 7. 重跑入口與產物

在專案根目錄執行：

```sh
tools/matching_decomp.sh
```

這支包裝器只在主機編排 Docker；分析與編譯都在限資源、無網路、非 root 容器裡執行。
原始 EXE 唯讀掛載，寫入只限研究輸出。IDA 由 [`tools/ida.sh`](../../tools/ida.sh) 的
`probe` 模式建立一次性 DB，再由 [`ida_matching_probe.py`](../../tools/ida_matching_probe.py)
匯出。重組、C 試驗與負對照由 [`matching_decomp.py`](../../tools/matching_decomp.py) 執行。

| 工具鏈 | 本次身分與入口 |
|---|---|
| IDA image | `ida-pro-9.4-idapython:py312-v1`；ID `sha256:4ac62de83339c215bab10e455cee3a22d9c6efed9fd0d8ed0f068327b83a06ab` |
| IDA 恢復 | `/home/anr2/ida_94_official/backups/docker-images/MANIFEST.md`；archive SHA-256 `d21d5072359ea427d182f1795aee158e95843085d0c460e1ec0aae2eb48bced0` |
| 組譯／C 編譯 image | 沿用既有 `pto2-remake-build:latest`；固定 ID `sha256:474f41ef91c354dd4754b08ef9302e965271e32417d6fba1772aecca0a5f9e2e` |
| 替換工具鏈 | 可設定 `WOLONG_IDA_PY_IMAGE`／`WOLONG_MATCH_BUILD_IMAGE`；新結果記錄實際 image ID，不能沿用本次收據身分 |

包裝器優先選已安裝的 `locked-v1` IDA image；本次它不存在，所以使用核對過 archive
與 image ID 的歷史 `py312-v1`。image 內的授權資料沒有進入專案輸出。
本次只是沿用已安裝的建置 image，不承諾從頭重建它的 OCI digest。

| 本機產物 | 用途 |
|---|---|
| `workplace/matching-decompilation/dosv/ida-probe.json` | 原始函式名、邊界、指令、xref、檔案偏移、輸入雜湊與 IDA 版本 |
| `workplace/matching-decompilation/dosv/input.exe.i64` | 本次從原始 EXE 新建的資料庫，不取代歷史 DB |
| `workplace/matching-decompilation/borland-control/ida-probe.json` | 正對照的標記與前綴掃描 |
| `workplace/matching-decompilation/build/reconstructed.S` | 三個函式的指令來源 |
| `workplace/matching-decompilation/build/report.json` | 匹配範圍、C 差異、輸入／DB／image 雜湊與負對照 |

本次 `report.json` SHA-256：`c4bcae3ed6b75458b1008ec354088fe225930fce3fa4f4e02f71950f3a633f35`。
DB 的容器路徑與建立時間可能讓重跑後的 DB 雜湊變動；每次都用該次收據綁定身分。
所有輸出抽查為 UID/GID `1000:1000`；本次容器均已退出並自動移除。
工具語法與專案文件門禁的驗證結果另記在 [`WORKLOG.md`](../../WORKLOG.md)。
