# 112：C 選單 selector、捲動與文字 callback

**狀態：14 named／兩 raw 入口通過；O0/O2 各 138 全裝置／ABI 相同，三個 XOR 還原與十二錯版通過。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear；檔案偏移為線性位址減 `0x10000` 加 512
- DB：`workplace/matching-decompilation/c-choice/ida/input.exe.i64`
- DB SHA-256：`4d84b66126fbcb8377707155cf56f45c190b9321faa224b05ab15fb943a0e8d9`
- 入口：[`ida_choice_probe.py`](../../tools/ida_choice_probe.py)
- 前置：[`re/111`](111-c-input-mouse-restoration.md)、[`re/84`](84-popup-row-band-and-world-cursor.md)、[`spec/45`](../spec/45-advise-scene-layout.md)
- 契約：[`spec/232`](../spec/232-c-choice-selector.md)

原 `sub_1036F` 建 24-byte SS frame，保存 cursor、建立列 callback、畫反白並等待。
`sub_1054D`／`sub_1059B` 保留絕對列與可見列兩個 index，頁面端移動仍走原 callback。
`sub_19479`／`sub_194BF` 分別保護 hotspot map 與保存原畫面，保留 selector 的 carry 返回。
原 callback `0x1945A` 的 AH stride、SI table 與位置都從 live operand 讀，不用靜態解碼初值。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 自然玩家完整進言流程 | 固定 pointer／按鍵序列的局部 caller 不代證正常玩家入口 |
| 原版動態安裝 scroll helper | 原 `nullsub_5` 的已載入狀態另留證據界線 |
| 原作者工具鏈／C 機器碼 | 仍待確認 |

## 原始定位與既有 caller 接線

14 個原始函式為 `sub_1036F`、`sub_103C3`、`sub_103E6`、`sub_10414`、`sub_104B5`、
`sub_104FF`、`sub_1054D`、`sub_1059B`、`nullsub_5`、`sub_1061F`、`sub_10B46`、
`sub_10BAF`、`sub_19479`、`sub_194BF`，共 940 bytes。原始 `0x19409..0x1945A`
81 bytes 與 callback `0x1945A..0x19479` 31 bytes 保留 raw 邊界，總計 1,052 bytes／474 指令。

`0x19409` 量第一列、算 stride 並 patch AH／SI；`sub_193E9` 暫存原 DL／DH operand，
待 selector 返回後還原。callback 的原始解碼初值 `0Ch`／`1234h` 不代表執行時值，
C 讀 live `0x9465`／`0x9467`；row callback 用 `CS:0x36D` 的原 near pointer。

`sub_10414` 先看左鍵狀態，再看右鍵。Y=100..131 等待；小於 100 上移並回置 131，
大於 131 下移並回置 100。原始絕對列與可見列分開，頁面端操作仍走 `nullsub_5` 和
callback，不猜補像素搬移。re/84 舊函式位址是 file／IDA 基準混用，原兩回置 Y 值也寫反，
已依原指令與輸入回播勘誤。

既有 C `sub_13B7E`／`sub_193E9`、`sub_13C99`／`sub_13CDC` 重新接真正的 mouse、
字庫、背景與 patch 等待。第一版 pixel 保存段 `CS:0x987C` 未初始化，原版保存覆蓋低記憶體
服務區；由原 `sub_19796` 所讀欄位設定獨立段 7200 後同矩陣乾淨重跑，沒有加大 budget。

## 同狀態矩陣

兩側獨立 DOS／Mouse／Font／RAM／VGA；原版 guest、C native 運算，C 不跑 guest CPU。
每組比全 1 MB RAM／四平面／GC／seq／latch／IN／OUT／Mouse／十四 register／FLAGS、
near／far／API 與 SS 快照及內容區像素。press／位置在固定 API 邊界送入，不把固定輸入當自然玩家。

| 群組 | O0/O2 各組數 |
|---|---:|
| selector／map／pixel 保護／兩模式／三 cursor／accept-cancel | 48 |
| Y 帶與移動序列 | 8 |
| 絕對列與可見列／頁邊界 | 16 |
| 保存恢復／XOR／callback／標題 primitive | 27 |
| 三組 patched X/Y | 3 |
| 九個原始選單，accept／cancel | 18 |
| 四劇本 speaker／advisor | 8 |
| 進言 caller 2–5 列 | 4 |
| 三個兩步 XOR 還原 | 6 |
| 合計 | 138 |

Y 序列對應選列 `[0,0,1,0,1,3,0,0]`，全部 live、原選單與進言 choice 接線得到預定原列。
三個 XOR 區域含右下邊界都完整還原四 plane 背景。每版兩側各 2,153 次全形、6 次半形、缺字 0；
O0/O2 完整收據逐 byte 相同，既有與新層的主要輸入／繪圖入口有 raw 快照。

| 刻意改錯 | 群組 | 首次拒絕 |
|---|---|---:|
| stride 少一倍 | corpus | 1 |
| callback 不讀 live SI | helpers | 10 |
| Y=131 也觸發下移 | bands | 2 |
| 用可見上限當總列上限 | scroll | 9 |
| accept carry 反轉 | selector | 1 |
| cancel carry 反轉 | selector | 2 |
| XOR 不讀 latch | helpers | 13 |
| 丟舊 DX 恢復 | helpers | 1 |
| 不恢復 carry | selector | 14 |
| 丟絕對列上移 | scroll | 3 |
| 丟絕對列下移 | scroll | 10 |
| 不讀 live DL | live | 1 |

[公開摘要](c-choice-verification.json) 保存原 bytes／DB／來源／flags 與收據 hash。
[來源](../../tools/c_recovery/choice.c)、[header](../../tools/c_recovery/choice.h)、
[生成 C](../../tools/c_recovery/choice_generated.inc)、[產生器](../../tools/c_recovery_choice_generate.py)、
[driver](../../tools/c_recovery_choice.go) 進版控；[重跑入口](../../tools/c_recovery_choice.sh)／
[驗證器](../../tools/c_recovery_choice_verify.py) 驗全部來源、原入口與錯版、backlink 與 C 乾淨重生。
Go1.26.7／GCC12.2／IDA9.4，所有原輸入／平台來源唯讀、非 root 限資源 Docker。

live operand 區十條／19 bytes 已真正 GNU 助憶碼匹配，先前 359 條補充不覆寫。
24,745 指令／56,287 code bytes 與完整 67,099-byte EXE 相同，沒有 instruction-byte fallback。
原版資料／圖片／字庫只留本機。原 `nullsub_5` 的動態安裝與完整正常玩家流程仍有獨立證據界線。

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_choice_probe.py workplace/matching-decompilation/c-choice/ida KI.EXE
tools/c_recovery_choice.sh
```
