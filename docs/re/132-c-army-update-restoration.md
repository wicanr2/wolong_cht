# 132：軍團輪轉與行軍更新的 C

**狀態：原始來源、直接控制流與窄驗證契約已核對；spec252 READY，原生C尚未驗收。**

- 日期：2026-10-10。
- 輸入：松崗DOS/V `KI.EXE`，SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 原始證據：`workplace/matching-decompilation/c-tick/census/ida-probe.json`，SHA-256 `937f8738225ae0f9e5efc133245316dd075d845ffded97a22f5f211b4804f2a7`。
- 對應IDA database：同目錄`input.exe.i64`，SHA-256 `4f01452cc2df29a9306aa4b444d1d4125e3f27f7730fcfdadbefd9f9da3a4667`。
- 工具：IDA Pro 9.4；位址為IDA database linear，段基址10000。檔案偏移另按`linear−10000+512`換算。
- 規格入口：[spec/252](../spec/252-c-army-update.md)；既有機制：[軍事](../mechanics/20-military.md)、[spec/168](../spec/168-corps-and-hour-cursors.md)。

[準備工具](../../tools/c_recovery_army_prepare.py)從固定舊probe抽取17函式，保存[本機來源](../../workplace/matching-decompilation/c-army/ida-source.json)的原始名稱、chunks、運算元、推論等級及交叉參照。
這份來源是既有IDA證據的範圍抽取，沒有執行新的IDA分析。
[指令重組收據](c-army-code.json)記錄434指令／1068bytes獨立匹配，以及現行`KI.code.S`對應切片的第二次組譯匹配；新增指令為0，全域組語與manifest不變。

## 原始範圍與閉包

17支原始函式為`sub_125A3`、`sub_12600`、`sub_1264A`、`sub_12662`、`sub_126FF`、`sub_12708`、`sub_127A2`、`sub_127F6`、`sub_12804`、`sub_12808`、`sub_12831`、`sub_12880`、`sub_128F4`、`sub_12A7E`、`sub_142AB`、`sub_14300`及`sub_1562B`。

已證實，434條指令／1068指令bytes均與固定KI原始檔逐chunk一致；本範圍沒有MZ重定位或間接CALL／JMP。`0x142E2`的XLAT是資料讀取，不當作間接呼叫。十個外部callee皆已在665函式／62原始區段的CONFORMED台帳內，包含[交戰與戰術](131-c-engagement-tactical-restoration.md)。

已證實，`0x126FF→0x12708`與`0x12804→0x12808`是原始落下邊，不補RET。原`0x12620/0x1262E`呼叫`sub_1562B`；其`0x15636`呼叫已還原的`sub_1563B`。

## 已證實的輸入與邊界

| 原始定位 | 窄結論與證據 |
|---|---|
| `0x125AD/0x125B2/0x125E9/0x125EE/0x125F6` | 讀CS:D18、處理16格、每格加40、批後比1FC0再寫回。由初值0出發的批首為0、400、800、C00、1000、1400、1800、1C00。 |
| KI檔案offset `0xF18`；SINARIO/SAVE槽內`+0x28` | 四劇本目前輸入均為0。這是輸入身分證據，不宣稱窮盡所有間接writer。 |
| `0x125B6`及尾批 | 原碼會讀到DS:4200的第128格；`sub_14F8A`建立此臨時城兵，`0x14FC8`清其flags。原始C不能套用Go只建模127個常駐軍團的跳過規則。 |
| `0x12613–0x12619` | 軍費是`(men>>1)+(men>>2)`，保留兩次無符號右移的各自截斷。 |
| `0x1263D–0x12646` | 士氣先以byte加10，再與勢力基準比較，保留回繞與原高byte。 |
| `0x12697–0x1269E/0x126F5–0x126FC` | 由原LDS取得佔用格指標，先減、處理移動，再加；不能用猜定的固定bank取代。 |
| `0x12831–0x1287F` | 依127筆原始順序找第一個同座標active軍團，再判owner；不改成跳過自己或最高戰力。 |
| `0x12705/0x12807` | 原CBW保留有號byte方向轉換；跨節點的ES索引及DS恢復按原碼處理。 |

CS:D18從存檔槽`+0x28`載入的間接鏈已有[spec/168](../spec/168-corps-and-hour-cursors.md)及既有時鐘研究。此次重用的probe不包含完整`sub_18CAE`，不由這份證據宣稱完整讀檔C已還原。

## 驗證與實作界線

原生工作由[容器入口](../../tools/c_recovery_army.sh)及[有界執行腳本](../../tools/c_recovery_army_container.sh)執行；[驗證工具](../../tools/c_recovery_army_verify.py)核對固定來源、乾淨重生、實際編譯旗標與同狀態收據。工具存在不代表原生驗收已通過。
容器工作結束後仍須另行執行`python tools/c_recovery_army_verify.py --repo /repo --output /output`，使用相同唯讀repo及可寫c-army掛載；payload退出0只代表該模式執行完成。

[C來源](../../tools/c_recovery/army.c)、[介面](../../tools/c_recovery/army.h)與[產生檔](../../tools/c_recovery/army_generated.inc)由[固定控制流產生器](../../tools/c_recovery_army_generate.py)重生。[十二個錯版定義](../../tools/c_recovery_army_mutants.tsv)對應批量、timer、軍費、士氣、方向、佔用格、碰撞、停戰、潰散及到達條件。17函式434指令的乾淨重生已通過；此項只證明來源一致，尚未證明同狀態行為。

[原生比較程式](../../tools/c_recovery_army.go)、[raw初態與獨立模型](../../tools/c_recovery_army_data.go)及[C呼叫適配](../../tools/c_recovery/army_fixture.h)為本輪研究來源。軍費／士氣窄矩陣32例已通過原版first與全狀態比較；完整工廠仍在核對，O0／O2與十二個錯版收據未齊前，維持READY，不提升C台帳。

後續C須保留原始名稱、chunks、位址與運算元，並與既有真C callee接線。兩側各自由唯讀原版與明示raw狀態初始化，固定亂數條件在執行前設定，不複製原版執行後Snapshot，不以原版CPU指令替代C。

驗證應覆蓋批首及末批、timer回繞、軍費與士氣字寬、道路節點與佔用格、野戰／攻城、暫存器／FLAGS／SS堆疊與原始落下邊。原Go規則與正常玩家驗收範圍不因這個C切片提升。

## 未解範圍

| 項目 | 狀態 |
|---|---|
| 原生C與完整同狀態矩陣 | C原型已產生；完整動態矩陣尚未完成。 |
| 完整主排程、正常玩家長程及原C機器碼 | 本輪局部軍團切片不代證。 |
