# 104：C VGA blit 與畫面保存

**狀態：八函式 O0/O2 各 9,942 組原版/C 的完整 RAM／四 plane／latch／I/O 相同，八個負對照拒絕；限固定平台契約與局部 blit／保存鏈。**

- 日期：2026-10-08
- DOS/V KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- KAOGRF.DAT SHA-256：`b9c7745e3ed9b32f0c12003fe81f82756136f738427dc421ec96f6c4ed5c4ac8`
- IDA Pro 9.4 database linear；檔案偏移另列
- DB：`workplace/matching-decompilation/c-vga/ida/input.exe.i64`
- DB SHA-256：`7e11f676288faef8bf0e105d9a6ee8c338cf8a77fe8fe6064c867fc64abd7118`
- 原始匯出：[`ida_vga_probe.py`](../../tools/ida_vga_probe.py)
- 契約：[`spec/223`](../spec/223-c-vga-blit.md)
- 既有證據：[`re/03`](03-image-blitter.md)、[`re/103`](103-c-hotspot-restoration.md)

## 1. 原始定位與已證實控制流

八函式共 467 bytes，沒有 MZ relocation。原始名稱、chunks、operand、xref、file／IDA bytes
保留，語意只附分級註記。以下 raw blit 控制流為已證實，平台結果限固定成熟裝置模型。

| 原始函式 | IDA 起點 | 檔案起點 | bytes | 局部語意 |
|---|---|---|---:|---|
| `sub_19796` | `0x19796` | `0x9996` | 45 | 原始pixel-to-byte保存wrapper、80-byte row offset、ES／DI與register保存 |
| `sub_197C3` | `0x197C3` | `0x99C3` | 45 | 原始pixel-to-byte恢復wrapper、DS／SI與register保存 |
| `sub_1F9B0` | `0x1F9B0` | `0xFBB0` | 107 | 原始四plane單byte列blit、GC序列、dummy latch read與CLD |
| `sub_1FA1B` | `0x1FA1B` | `0xFC1B` | 28 | 原始單byte列MOVSB、dummy ES read、DF／word wrap與row pitch |
| `sub_1FA37` | `0x1FA37` | `0xFC37` | 107 | 原始四plane雙byte列blit、GC序列、來源連續性與CLD |
| `sub_1FAA2` | `0x1FAA2` | `0xFCA2` | 32 | 原始雙MOVSB列、dummy ES read、DF／word wrap與row pitch |
| `sub_1FAC2` | `0x1FAC2` | `0xFCC2` | 79 | 原始四plane保存、byte width倍乘、Read Map Select與保存 |
| `sub_1FB11` | `0x1FB11` | `0xFD11` | 24 | 原始REP MOVSB保存列、來源pitch80／目的連續、DF／word wrap |

## 2. 平台前提與獨立狀態

VGA標準契約依 [IBM VGA/XGA Technical Reference，1992，2-82至2-85](https://bitsavers.trailing-edge.com/pdf/ibm/pc/cards/IBM_VGA_XGA_Technical_Reference_Manual_May92.pdf#page=99)：
讀取載入四個 latch、Read Map Select 決定資料，write mode／set-reset／logical function 決定實際plane。
這是平台前提，不把標準硬體重新逆向成遊戲 gate。

原版側使用固定 dosgolem `machine.New`／`LoadEXE` 與真實 guest CPU。C 依原始 register、
stack、byte／word／DF 自行計算 loop／I/O，沒有呼叫原版 guest 函式或複製它的結果。
C 使用第二台獨立 machine／VGA；兩側 plane、GC、seq、latch、RAM、input 都獨立。
普通 C RAM 由 C malloc 持有，裝置 bus 經 [callback](../../tools/c_recovery_vga_bus.go) 處理真正 memory map，
沒有把A0000當平面RAM。共享的只有成熟平台實作，不共享遊戲CPU算法或 mutable 狀態。

平台 probe 已驗MZ loader、獨立plane與latch。私有 build module 位於容器 `/tmp/vga-build`，
import 固定dosgolem internal API，專用 `matching_vga` tag 不混入 production 套件。
沒有修改外部dosgolem專案，所有輸入與來源唯讀。內部API版本由revision與逐檔雜湊固定。

## 3. 原始演算法接線

GC setup 的 write mode0、set/reset0、bit maskFF、enable set/reset 0E／0D／0B／07、
第二plane起OR 10h與四次row helper保持。來源SI連續跨plane，不重設；dummy ES read在
每個MOVSB前載入latch。單byte與雙byte列分開，row pitch80以原始word運算前進。
保存會依Read Map Select0／1／2／3讀四plane，目的buffer連續；寬度只倍乘AL，不讓carry
默默進AH。原始CLD與leaf的DF、SI／DI word wrap、全部register保存保持。

原版I/O port/value序列與C序列逐筆相同。原版Step計數保留在平台log，C沒有相同CPU
指令數，驗證不宣稱instruction-time或實機wall-clock parity。GC／seq／latch與目前index另逐byte比較。

## 4. 完整矩陣與實際資產

| 群組 | O0/O2各原版/C |
|---|---:|
| draw | 2,800 |
| leaf | 5,880 |
| save | 200 |
| wrapper | 1,000 |
| real-assets | 8 |
| roundtrip | 54 |
| 合計 | **9,942** |

- draw：兩blit × 七尺寸 × 五map mask × 五destination × 四source × 兩DF。
- leaf：三helper × 七尺寸 × 七write-mode／logic條件 × 五destination × 四source × 兩DF。
- save：五尺寸 × 五source-rectangle offset × 四buffer offset × 兩DF。
- wrapper：兩保存／恢復 × 五尺寸 × 五X × 五Y × 四read map。
- real-assets：原始150筆／每筆2048 bytes圖庫中的0、5、49、149四筆 × 兩destination。
- roundtrip：18組三步save／draw覆寫／restore，每一步都比完整RAM／plane／ABI，最後全plane與保存前相同。

每版全部9,942例都比十四register／FLAGS、完整1 MB RAM、四個65536-byteplane、
GC／seq／latch／index、port序列與90-bytecallee快照。另310次640×400 indexed pixels抽查相同。
每函式都有 `entries_seen`。原始plane pattern與每次input／register在執行前固定，沒有mask差異或挑結果。

O0/O2 JSON逐byte相同，SHA-256：`67413f4a383ce39fbe1c297d9fb6743d018cc3d773630d7c0485481f4b4fbbe3`。
原版/C從各自狀態計算，摘要同為 `63044c32589fe246d8f4e6d78780f89f4cb9d2641bc192be855b74abc8765a8d`。

八組實際頭像blit各寫兩張PNG，本機原版／C PNG逐byte相同；16個PNG雜湊記在verification.json。
它們使用明示debug灰階palette，不是原版palette或正常玩家畫面。原版美術與衍生圖僅留本機ignored目錄，
不進Git／Release。圖片外送檢視被自動審查拒絕後，改以容器內planes／indexed pixels／PNG bytes核對，沒有繞過拒絕。

| 刻意改錯的 C | 首次拒絕 |
|---|---:|
| 不讀dummy destination latch | leaf 第1組 |
| 雙byte列只複製一byte | draw 第1401組 |
| row pitch改79 | leaf 第1組 |
| 第二plane起不做OR | draw 第1組 |
| plane2的ESR改plane1 | draw 第1組 |
| 保存一直讀plane0 | save 第1組 |
| wrapper X少一次SHR | wrapper 第21組 |
| 不清DF | draw 第2組 |

## 5. 編譯來源與重跑

全C／header／fixture／driver manifest SHA-256進實際CGOflags，O0/O2 binary的build info
由verifier核對。Compiled source manifest digest：`03632a6dd02c01cb857608b40167b58c0e1031f9fc466e4ce1259805fde557cf`。

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_vga_probe.py workplace/matching-decompilation/c-vga/ida KI.EXE
tools/c_recovery_vga.sh
```

Go1.26.7／GCC12.2.0／IDA9.4，dosgolem revision `a9714ebdab2ad6b529f81225472680f2b11f2842`，
實際import tree雜湊前後相同。非root、無網路、限資源Docker，只明確輸出可寫，
IF/TF=0、獨立SS／來源／保存buffer。未定義FLAGS依固定CPU模型，不外推實機硬體時序。

| 來源 | SHA-256 |
|---|---|
| vga.c | `9c9c75270422a8950d9d09c1321d5093423d2eada6d5bde3e417c9d711ba9b5b` |
| vga.h | `40f9e9ee054830df5c95e5bbe0185f82f3b8d7d8e9e780510da29794b2ab7699` |
| vga_fixture.h | `97ab311f79b371c26b90d450466173e9ac90e70c55f6862a204bef40db569096` |
| Go driver | `3e0b75bce9fde3175dc056672134b0a256c9a58d601a927f65b867e2dc37bfbd` |
| Go bus callback | `b842e1561a96da289bac5ed83f4756c2d6e77883681227651034aad3625d0c86` |

本機完整RAM／plane與debug PNG在 `workplace/matching-decompilation/c-vga/`，驗證入口
[`c_recovery_vga_verify.py`](../../tools/c_recovery_vga_verify.py)，C現況分級台帳
[`c-recovery-status.json`](c-recovery-status.json) 累計136函式。
版控驗證摘要在 [`c-vga-verification.json`](c-vga-verification.json)，保留來源雜湊、
工具版本、各函式入口、原版／C結果雜湊與八個負對照。原版美術沒有納入摘要。

## 6. 未解範圍

| 項目 | 邊界 |
|---|---|
| 完整normal-player畫面／視窗合成 | 局部primitive與debug PNG不代證正常玩家路徑，Issue #22保持OPEN |
| 實機VGA逐週期／wall-clock | 固定成熟平台模型的局部對拍，不稱同硬體時間 |
| 完整Go圖形三方 | 未驗證 |
| C機器碼與原作者工具鏈 | 未驗證，整檔binary match仍是組語基準 |
