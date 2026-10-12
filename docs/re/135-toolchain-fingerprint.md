# 135：KI.EXE 工具鏈指紋（機器碼匹配前置判定）

**狀態：READY。2026-10-12。判定 `KI.EXE` 為手寫 16-bit 組語，非 C 編譯器產物；結論只限 `KI.EXE`。同日 IDA 9.4 從 signature 面覆核確認（§7）。**

- 日期：2026-10-12。
- 出處：spec [`255`](../spec/255-c-machine-code-feasibility.md)、試點 [`re/89`](89-matching-decompilation-pilot.md)（GCC `-m16` 0／12）、本輪 `tools/tc_fingerprint.py` 全綠輸出（見 §6）。
- 輸入：DOS/V `KI.EXE`，SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`；函式邊界 `function-ledger.json`（739 支，唯讀引用）；指令文本 `tools/c_recovery/KI.code.S`（25097 行）。

## 1. 判定

`KI.EXE` 是手寫 16-bit 組語（small 模型、register 呼叫慣例、DS／CS 全域變數）。
C 機器碼匹配在本檔上**不可行**：沒有編譯器＋選項能逐函式編出相同的機器碼，
因為輸入本來就不是編譯器產物。`c_machine_code_match` 旗標維持關閉。

| 判定三件套（spec 255 §決策準則） | 本筆記位置 |
|---|---|
| 指紋表 | §2–§5 |
| 證否線 | §5（位址逐項） |
| 試點編譯差距 | `re/89`：GCC `-m16` 12 支試點 0 匹配 |

## 2. Startup 形狀（非 C runtime）

`start`（IDA `0x10000`，67 bytes）前 7 條：

```
0x10000  mov  ah, 51h        ; 取 PSP 段址
0x10002  int  21h
0x10004  mov  ds, bx
0x10006  mov  al, ds:80h     ; 命令列尾長度
0x10009  mov  cs:byte_10D36, al
0x1000d  mov  ax, cs
0x1000f  mov  ds, ax         ; DS=CS，全域走 CS 段
```

遊戲開機流程（PSP 抓取、命令列尾存檔、DS 歸位），之後直入遊戲初始化。
無 heap／stack 設定、無 argv 建構、無 `_main` 轉跳——Borland `c0.asm`、
MSC `crt0`、Watcom `cstart` 都不是這個形狀。**confirmed**。

MZ 標頭旁證（`workplace/orig/dosv/KI.EXE` 前 0x200 bytes，只讀解析）：
116 個重定位項、標頭 32 paragraphs、**minpara=0／maxpara=1**（linker 預設多為
maxpara=`FFFFh`；只要載入映像、不留 heap，手工調過的形狀）。**強證據**（linker
可經選項調成這樣，故只列旁證，不列證否）。

## 3. 全庫慣例普查（檔案 bytes，非 IDA 語意）

`tools/tc_fingerprint.py` 輸出（assert 全過）：

- 739 支起點 bytes：`50 53`（push ax／bx）77 支、`1E 06`（push ds／es）45 支、
  `1E 50` 38 支；`0x55` 開頭僅 8 支。
- 8 支中 6 支是 `push bp；sub sp,N；mov bp,sp` 手寫 frame
  （`sub_17FDB／1820E／19479／194BF／1F4DF／1F593`）：
  暫存器傳參、`[bp+N]` 區域變數、`add sp,N` 收尾，
  **6 支體內零 `[bp+8]` 以上堆疊參數引用**（C frame 必有呼叫端推參）。
  另 2 支（`sub_13262／1812A`）只是 push bp 存暫存器。
- 單 chunk 737 支：尾 `C3`（retn）717 支；唯一 `CB`（retf）是驅動段
  `sub_20000`（far-call 目標，本來就該 retf）；多 chunk 2 支
  （`sub_1AF69／1DD22`，緊鄰 4-byte 跳過 stub，非 Watcom 式遠端 outline）。
- `enter`／`leave` 零條；`pusha／popa` 僅 2 組；32-bit GPR 引用零條
  （程式碼段內 98 個 `66h／67h` bytes 全是運算元或資料，非前綴）。
- `rep` 字串運算 27 條，散在 19 支函式（inline 字串搬移，非函式庫呼叫）。

## 4. 系統呼叫層（直呼 DOS，無 C 函式庫）

68 處 `int`（`tc_fingerprint.py` 嚴格計數；首輪 grep 漏數曾寫 19，已更正）：

| 向量 | 次數 | 用途 |
|---|---|---|
| `21h` | 34 | 直呼 DOS：`51h` 取 PSP、`09h` 印字串、`4Ch` 離開、`3Fh／40h／3Eh` 讀寫關檔、`48h／49h` 配放記憶體；AH 常由呼叫端經暫存器帶入，可見行之外無包裝層 |
| `61h` | 16 | 遊戲自訂常駐向量（tick／回呼，見既有 RE 筆記） |
| `33h` | 8 | 滑鼠 |
| `50h` | 4 | 自裝向量鉤：`sub_1E364／1E378／1E38C／1E3A6` 四支短截擊器，先 `call` 原處理、carry 才 `int 50h`（bytes 如 `1E 53 52 56 E8…5E 5A 5B 1F 73 06 B4 00 CD 50`）；編譯器產物不會自裝 `50h` 錯誤鉤 |
| `10h／15h` | 各 2 | 顯示／雜項（含字型服務呼叫端） |
| `1Ah` | 1 | `sub_1EC82` 讀時鐘 |
| `int 3` | 1 | `0x1CCC9`，**不屬任何函式 chunk**：bytes `…C3 CC EC…`（retn 後單一 `CCh`）。MS padding 是成串對齊，單一 `CC` 是手工斷點殘留或手動填充。**強證據** |

無執行期函式庫徵候：census 739 支全是 `sub_XXXX／start`，
無 `F_LXMUL`、`N_LXMUL`、字串／檔名／`printf`／`malloc` 家族；
MZ EXE 無 import 表，對外只走 INT。靜態連結的函式庫必帶自家 startup 與 helper
命名殘留，此處皆無。**confirmed**（斷言的是「全庫掃描零命中」，
不是邏輯上的「不存在」）。

## 5. 證否線（C 編譯器不會產生的結構，逐項定位）

| # | 結構 | 位址 | 等級 |
|---|---|---|---|
| 1 | 自修改碼 `mov cs:[…],…` | `0x153B2`（`mov cs:[di],ax`）；另 `re/28` 的 `loc_1A065`、`loc_10701` 兩處 | confirmed |
| 2 | byte 暫存器雜耍 `xchg ah,al`／`xchg bh,bl`／`xchg dh,cs:…` | `0x107F5`、`0x19115`、`0x192AE`、`0x193EE`；另 `xchg ah／bh,[si／bx／di+N]` 六處以上（`0x13408`、`0x14CF5`…） | confirmed |
| 3 | 中段入口／落空：`sub_1C300`（`0x1C300–0x1C30D`）直落 `sub_1C30D`（`0x1C30D–0x1C315`），入口即前一支的尾端 | `0x1C300／0x1C30D` | confirmed |
| 4 | retf 驅動段（far-call 目標，非函式） | `sub_20000`（`0x20000–0x2000E`，尾 `CB`） | confirmed |
| 5 | 自裝中斷向量鉤（`50h` 四截擊器，見 §4） | `sub_1E364／1E378／1E38C／1E3A6` | confirmed |
| 6 | 暫存器呼叫慣例貫穿全庫（AH／參數由呼叫端帶入，無堆疊參數） | 全庫；6 frame 函式體內零堆疊參數引用 | confirmed |
| 7 | CS 段全域變數（`mov cs:byte_…` 寫全域） | startup `0x10009` 起，全庫慣例 | confirmed |

以上任一項都足以排除「整支 C 編譯」；七項齊備，且分佈在 startup、
系統層、遊戲核心（`0x107F5`、`0x153B2`、`0x1C300` 皆在遊戲邏輯段）。

## 6. 工具收據

`tools/tc_fingerprint.py --repo .`（2026-10-12，docker `python:3.13-bookworm`）：

```
prologue top bytes: {(0x50,0x53): 77, (0x1E,0x06): 45, (0x1E,0x50): 38}
BP-frame funcs (6): [sub_17FDB, sub_1820E, sub_19479, sub_194BF, sub_1F4DF, sub_1F593]
push-bp-save funcs (2): [sub_13262, sub_1812A]
single-chunk: 737, tail C3: 717, tail CB: 1, multi: [sub_1AF69, sub_1DD22]
tail-CB owner: [sub_20000] (far-call driver segment, retf expected)
indexed instructions: 25097
int sites: {10:2, 15:2, 1A:1, 21:34, 3:1, 33:8, 50:4, 61:16} total=68
stack-arg refs in framed funcs: 0
rep string ops: 27 in 19 funcs
32-bit GPR refs: 0
toolchain fingerprint: PASS (hand-written 16-bit assembly shape)
```

## 7. IDA 9.4 確認（2026-10-12 追記）

以 `tools/ida.sh probe` 從唯讀原檔建一次性 DB（image
`ida-pro-9.4-idapython:py312-v1`，非 root），探針 `tools/ida_toolchain.py`，
輸出收據留 `/tmp`（`tc-probe3`、`tc-logo`），不入版控。

| 項目 | KI.EXE | LOGO.EXE（同隊已知 Borland C++ 1991，正對照） |
|---|---|---|
| probe（函式數／SHA） | 739／`fffeba98…3868` 相符 | 103／`2d36530b…aeef` |
| FLIRT：`FUNC_LIB` | **0** | 39（`_main`、`_malloc`、`_exit`、`_atexit`…） |
| 非 `sub_` 命名 | 僅 `start`＋7 個 `nullsub_N` | `start`＋33 個 runtime 名 |
| `cc_id`（IDA 自帶欄位） | 129 → API 渲染 `Visual C++` | 2 → `Borland C++` |
| imports | 0 | （同為 MZ，直接 INT） |

判讀（每條帶等級）：

- **FLIRT 零命中是 signature 面的答案**（confirmed）。同一支探針、同一顆 image
  在 LOGO 上看到 39 個函式庫函式，探針不是瞎的；KI 的 739 支裡 IDA 認不出任何
  runtime／library——Borland、MSC、Watcom 的 DOS signature 全數落空。
- **`cc_id=129／Visual C++` 是 loader 預設值，不是偵測結果**（強證據）。
  LOGO 的 `cc_id=2` 是 Borland signature 命中後由匹配寫入的；
  KI 沒有任何命中能寫入它，欄位停在 DOS MZ 預設。**不可讀成「KI 是 MSVC 編的」**——
  MSVC DOS 產物必帶自家 startup 與 `N_`／`F_` helper signature，此處零命中。
- **linker 也不同**（confirmed，只讀 MZ 標頭交叉比對）：
  三支 Borland 檔 `reloff=0x3E`（`0x1E` 起是 TLINK stub 區，LOGO 首 bytes
  `fb506a72…`）、`maxpara=65535`、`checksum=0`；
  KI `reloff=0x1E`（重定位緊接標頭）、`maxpara=1`、`minpara=0`、
  `checksum=0x764a`（非零，有人算過）。語言不同，linker 也不同。
- **Ghidra 交叉驗證未做**：本環境無 Ghidra（無 image、無 binary）。
  上面三條（FLIRT＋正對照、cc 欄位、linker）已各自獨立指向同一結論，
  Ghidra 若日後可跑，只做覆核，不預期翻案。

## 8. 未解與邊界

- 6 支 frame 函式的「手寫區域變數」是結構推定（暫存器傳參＋零堆疊參數）；
  若有人主張其中混了 C，需逐支體審查反駁，不影響全庫判定。
- `0x66／0x67` 98 bytes 只驗了「零 32-bit GPR 引用」，未逐 byte 解碼；
  不影響判定（即使有 386 指令，手寫仍成立）。
- 結論只限 DOS/V `KI.EXE`；PC-98 版與其他七支執行檔不在本判定範圍。
- 候選編譯器（Borland／MSC／Watcom）因判定為手寫而**不需取得**；
  若未來有人翻案，需先推翻 §5 證否線。
