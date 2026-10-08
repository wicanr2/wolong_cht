# 103：C 熱區資料路徑與數值視窗接線

**狀態：十函式 O0/O2 各 261,367 組原版/C 相同，80 組 Go 數值相同，八個負對照拒絕；限原始 memory／ABI 與明示 pixel／VGA primitive。**

- 日期：2026-10-08
- DOS/V KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- SINARIO.DAT SHA-256：`21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c`
- IDA Pro 9.4 database linear；檔案偏移另列
- DB：`workplace/matching-decompilation/c-hotspot/ida/input.exe.i64`
- DB SHA-256：`e5b4874f382093fb14953bc6bb25fe8072338061f7ce72010621a37b8be71d5a`
- 匯出：[`ida_hotspot_probe.py`](../../tools/ida_hotspot_probe.py)
- 契約：[`spec/222`](../spec/222-c-hotspot-map.md)
- 既有定位：[`re/22`](22-strategy-command-tree.md)、[`re/47`](47-main-screen-window-registry.md)

## 1. 原始定位與分級語意

十函式共 491 bytes，沒有 MZ relocation。原始 names、chunks、operand、xref、file／IDA bytes
保留。以下 raw memory 與 caller 控制流為已證實，限本輪 map／fixed input／primitive contract。

| 原始函式 | IDA 起點 | 檔案起點 | bytes | 已證實的 raw 語意 |
|---|---|---|---:|---|
| `sub_1E3C0` | `0x1E3C0` | `0xE5C0` | 23 | 熱區 segment／base 初始化、CLD、4000 bytes 清除與 STOSW offset wrap |
| `sub_1E3D7` | `0x1E3D7` | `0xE5D7` | 68 | 熱區矩形 raw byte 登記、word 定址、零 byte loop 與重疊 CF |
| `sub_1E41B` | `0x1E41B` | `0xE61B` | 56 | 熱區矩形清除、原始 word 定址與保存 |
| `sub_1E453` | `0x1E453` | `0xE653` | 38 | 熱區 pixel query、保留 high Y byte、AL 及原始 ABI |
| `sub_1895D` | `0x1895D` | `0x8B5D` | 71 | 視窗 raw style、tile flags／border、Y+2 與熱區接線 |
| `sub_189DE` | `0x189DE` | `0x8BDE` | 18 | 16px wrapper 至8px熱區尺寸／座標的原始參數 |
| `sub_1D5D4` | `0x1D5D4` | `0xD7D4` | 65 | 原始 tile byte flags、8-byte record、40欄word offset與保存 |
| `sub_10C14` | `0x10C14` | `0xE14` | 76 | 原始外框幾何、水平／垂直 caller、word arith 與保存 |
| `sub_10C60` | `0x10C60` | `0xE60` | 23 | 原始水平 border 的 primitive 參數／LOOP／保存 |
| `sub_10C77` | `0x10C77` | `0xE77` | 53 | 原始垂直 border 的 primitive 參數／LOOP／保存 |

## 2. writer／reader 閉合與勘誤

`sub_1E3C0` 把 ES／DI 寫 CS:E479／E47B，REP STOSW 清 4000 bytes。
`sub_1E3D7`／`sub_1E41B` 寫該 segment／base 的矩形，`sub_1E453` 同址回讀 AL。
資料路徑為 80×50 熱區圖，每 cell 管8×8 pixels，沒有 glyph blit。

原始 writer `AND BL,F8h`、reader `AND DL,F8h` 都保留 high byte。完整 word 公式是
`offset = base + (x >> 3) + 10 × (y & FFF8h)`，在原始 u16 運算中 wrap。
[re/22](22-strategy-command-tree.md)、[re/47](47-main-screen-window-registry.md) 的先前 F8h
word 公式遺漏 high byte；原始 low-byte operand 留著，現行公式更正，全部640×400 query 含 Y≥256驗證。

[re/102](102-c-numeric-editor-restoration.md)、[spec/78](../spec/78-amount-input-editor.md) 與語意索引
把 0x1E3D7 導覽為 glyph 的前提不成立。舊 `glyph` 群組名／C ABI 收據保留，只是導航，
現行導覽改成熱區，新增證據比先前只看 caller 的推測更直接。
Spec/78 的假 glyph 未讀列由此閉合；實際 VGA plane／blit 與正常 UI 邊界沒有被宣稱完成。

DI=FFFFh 的 STOSW 依固定 dosgolem write16 模型，第二個 byte 在同 segment 0000h。
第一次 C 使用一般 linear write 跨到下一段，原版/C map／flags 的差異直接指出問題。
新 init 改兩個具 u16 offset 的 byte write，保留同一 DI=FFFFh 案例，不遮罩或排除邊界。

## 3. 真實 pixel 至財政鏈

原版與 C 保留真正 writer／reader／wrapper／border caller。每份 key 計畫先從原始十八格
raw table 找到 cell，再換成實際 pixel。輸入 leaf IDA 0x121E7 提供固定 CX／DX／cancel，
沒有覆寫 0x1E453 回傳鍵碼。原始熱區資料真的經 query 進數值主迴圈與財政 caller。
明示尾端取消保留，兩側輸入在執行前固定，不挑結果。

map 與 tile flags 各完整 65536-byte memory region 比較，包含正常4000-byte圖與rawword/base wrap。
Window wrapper 同時接回 tile byte flags、border 幾何與熱區登記。VGA blit 還是 primitive，
沒有把平面 RAM 模型當成真正的四 plane memory device 或像素匹配。

RET primitive 沿用 [re/102](102-c-numeric-editor-restoration.md) 的裝置／文字／其他世界邊界，
另有 IDA 0x1F9B0。0x10C14、0x1895D、0x1E3D7、0x1E453 都真正執行，實際 VGA
0x1F9B0／0x1FA37 等不在本輪完成範圍。Popup／FAR 裝置仍使用明示契約，見舊收據。

## 4. 完整矩陣

| 群組 | O0/O2 各原版/C | Go scalar |
|---|---:|---:|
| query-pixels | 256,000 | 0 |
| query-word | 240 | 0 |
| init | 5 | 0 |
| rectangle | 432 | 0 |
| register-cells | 4,000 | 0 |
| wrapper | 450 | 0 |
| numeric-pixels | 80 | 80 |
| finance-pixels | 64 | 0 |
| scenario-hotspot | 96 | 0 |
| 合計 | **261,367** | **80** |

- query-pixels：全部640×400 pixel，map pattern 特別區分 Y=0與Y=256的列。
- query-word：五個base × 六X × 八Y，保留word wrap與高byte。
- init：五base，含奇數／FFFFh起點。
- rectangle：兩writer × 六width × 六height × 兩base × 三code，含原始0→256 byte loop。
- register-cells：全部4000 cells，各註冊一格。
- wrapper：六caller × 五width × 五height × 三style。Border dimensions≥2，原始零尺寸另由map驗。
- numeric-pixels：五cap × 十六真實pixel輸入計畫，同時比正式Go數值核心。
- finance-pixels：四caller × 十六計畫，真實map／query與CSglobals。
- scenario-hotspot：四劇本 × 三玩家 × 兩modal × 四計畫。

每版2042次完整1 MB抽查；每案例全十四register／FLAGS、world、globals、queue、RNG、
map／flags、原始stack與控制欄位。每callee入口292 bytes，十函式皆有實際 `entries_seen`。
Go 80組只比 `state.EditAmountValue` 的數值，不代表完整財政／UI；正式Go沒有改動。

O0/O2 JSON 相同，SHA-256：`8c8020bd8aeed24acfb0e8d2bace25bff9553ee9b7cb653b9a26b83b8c9f7fa0`。
兩側分讀自己的 memory，狀態摘要相同：`ffd29827711bd7b1454a3c7f6281c85c61bfa89756c2b5767149a77033e0a7fe`。

| 刻意改錯的 C | 首次拒絕 |
|---|---:|
| init 少清一個 word | init 第 1 組 |
| writer 丟 high Y byte | rectangle 第 1 組 |
| 重疊不回 CF | rectangle 第 1 組 |
| map row pitch 改 79 | rectangle 第 1 組 |
| query 丟 high Y byte | query-pixels 第 163841 組 |
| tile flags bit 4 未設 | wrapper 第 2 組 |
| border pitch 改 639 | wrapper 第 2 組 |
| window Y 只加一格 | wrapper 第 1 組 |

## 5. 實際編譯來源與重跑

外部 C include 修改未觸發原有 cgo cache 重編，曾使 init修正後仍跑舊碼。
[`c_recovery_hotspot.sh`](../../tools/c_recovery_hotspot.sh) 將全C／header／fixture／driver清單
的 SHA-256 放入實際 CGO flags，另存 binary `go version -m` build info。
Verifier 機器核對 digest與兩個optimization的實際flags，不能只有新的 source清單與舊binary。

Compiled source manifest SHA-256：`37037190a861969d5eae4d56dfa6592fa127d871a7e4d86d547f1085858c5802`。

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_hotspot_probe.py workplace/matching-decompilation/c-hotspot/ida KI.EXE
tools/c_recovery_hotspot.sh
```

Go1.26.7／GCC12.2.0／IDA9.4，非root、無網路、限資源Docker，原版／專案／dosgolem／
既有modulecache唯讀。dosgolem revision `a9714ebdab2ad6b529f81225472680f2b11f2842`，import來源前後雜湊相同。
RNG表固定12:34:56、raw c/s執行前同寫兩側，沒有重擲；IF/TF=0，SS分離，未定義FLAGS沿用固定模型。

| 來源 | SHA-256 |
|---|---|
| hotspot.c | `805a2963363f151d01478d1eefd956ac482e2df22189d15638c7f79ea94dd5c1` |
| hotspot.h | `21f1dcb13ecbc301237487a12d6db80de4c944d816892a989a9a48711f496b20` |
| hotspot_fixture.h | `ed5a4f5b1b4ecdd5982689110d7b82c2372ead6abd83551a9a18357d06862411` |
| Go對拍工具 | `03438c27ed659a1006631d578087936d5b3333e3123a306654cd59259c54125e` |

本機收據在 `workplace/matching-decompilation/c-hotspot/`，核對入口
[`c_recovery_hotspot_verify.py`](../../tools/c_recovery_hotspot_verify.py)，現行C台帳
[`c-recovery-status.json`](c-recovery-status.json) 共128個函式。原始private資料不進Git。

## 6. 未解範圍

後續 C VGA 證據見 [re/104](104-c-vga-blit-restoration.md)。本規格原始熱區收據的blit仍是RET；
新收據對0x1F9B0／0x1FA37及保存wrapper使用獨立VGA裝置，整條正常玩家與其他primitive範圍保留。

| 項目 | 邊界 |
|---|---|
| VGA plane／blit／畫面保存 | primitive 不代證真實 memory device 或像素 |
| 真實滑鼠polling與正常玩家視窗 | 固定pixel只量控制流／map，Issue #22保持OPEN |
| 完整Go視窗／財政 | 局部map／scalar不外推完整玩家路徑 |
| C機器碼與原作者工具鏈 | 未驗證，整檔binary match仍是組語基準 |
