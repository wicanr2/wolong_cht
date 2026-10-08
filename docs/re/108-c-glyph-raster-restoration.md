# 108：C 字型向量、glyph raster 與固定三字名稱

**狀態：四個原始 C 函式與兩個 code 入口，O0/O2 各 252 組真實字型／raster／全 RAM／plane／ABI 相同，九個錯版拒絕。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- IDA Pro 9.4 database linear；檔案偏移另列
- DB：`workplace/matching-decompilation/c-glyph/ida/input.exe.i64`
- DB SHA-256：`bd9456d704fb5f1d7ec6fe544675cd106c434d53bcad8afdbc11c7b1e4c6e9d7`
- 探針：[`ida_glyph_probe.py`](../../tools/ida_glyph_probe.py)
- 來源契約：[`spec/228`](../spec/228-c-glyph-raster.md)
- 既有證據：[`re/29`](29-font-service-int15.md)、[`re/107`](107-c-display-interpreter-restoration.md)

## 1. 原始定位與輸入

`sub_1F720` 62 bytes 設定兩個 INT 15h 字型向量；`sub_1F7A4` 212 bytes 處理真正
VGA glyph raster，`sub_106F5`／`sub_106FD` 各 4 bytes 是文字與三字 caller。
原始 glyph 區 `0x1F75E..0x1F7A1` 67 bytes 與名稱區 `0x10701..0x1075B` 90 bytes
不是原作者函式邊界證明，保留原名稱、資料分類、bytes、operands 與內部入口。

本機 END_S13.DAT 404,992 bytes，SHA-256 `0ab4c43920787bf87c91bba5730078e451cb6125b878f201d3e70d648412569b`。
END_S14.DAT 3,840 bytes，SHA-256 `1d0cf09d0a319a9e7039190688c6a905ba4370bd369fbfcfa43b2078480d6918`。
兩檔與原版美術一樣只作本機唯讀輸入，不進版控。

## 2. 真實 C 與平台界線

原版側完整跑原始字型向量、patched far call 與 glyph raster。
C 還原同一個 near／far ABI、32-byte SS 緩衝、signed／word 定址、mask、色彩／背景、
latch 與 16 列 source stride。C 不呼叫 guest CPU，不以原版畫面／結果回填。

兩側各持有獨立 dosgolem Machine／DOS／字型 cache／RAM／VGA。C 只將暫存器傳給
既有 DOS/V 服務鉤子，不執行 stub 指令。這是採用成熟平台服務契約，不還原 STR.EXE 的
逐字磁碟 open／read／close，也不宣稱實機時序。
初始化 API、來源檔名與 slot 0x410／0x414 契約見 dosgolem `internal/dos/font.go`。

`loc_10701` 會把輸入 AX 寫到原始 `0x71F`／`0x736`／`0x74D` operand。
native C 需讀取該 live operand，不能一直用解碼時的 9001h。glyph far operand 同樣由
初始化寫入後讀取。原有顯示清單原場景／字串／陰影要用真正字形重跑。

## 3. 驗證與未解範圍

O0/O2 比完整 RAM、四平面、latch／GC／seq／port、十四 register／FLAGS、SS buffer、
near／far／服務入口。覆蓋單／雙 byte、八種 X 對齊、右／下界、不畫分支、字色／背景／
模式 bit、固定三字來源、live operand、實際原場景。字形來源與大小先固定再執行。

| 項目 | 邊界 |
|---|---|
| STR.EXE 與封裝檔名考古 | 本輪使用已定案平台字型 API，不重新研究 TSR 安裝 |
| 真實硬體 wall-clock | 固定平台模型不代證 |
| 正常玩家完整 UI 操作 | glyph／場景矩陣不代證自然操作 |
| C 機器碼與完整 C | 本輪仍未完成 |

## 4. 真實字形與原始控制流收據

O0/O2 各 252 組全 1 MB RAM、四個 65536-byte plane、GC／seq／latch／port、
十四 register／FLAGS、原始 near／far 入口與 SS buffer 相同。每組皆比較 640×400 內容區像素。
兩側各實際讀取 740 個全形與 150 個半形字模，缺字計數皆 0；原型用的 glyph no-op 已移除。
服務 stub 標為真正 `0080:0410`／`0080:0414`，不包裝成 KI 線性位址。

| 群組 | 每種最佳化組數 |
|---|---:|
| 原始 INT 15h 向量初始化 | 2 |
| 單／雙 byte／顏色／裁切 | 140 |
| 八種 X 對齊 | 64 |
| 原始字串／固定三字／live operand | 16 |
| 十個原場景 | 30 |
| 合計 | **252** |

原字形入口 `0x1F75E` 呼叫 996 次，raster `sub_1F7A4` 實際進入 890 次，其餘
106 次為右／下裁切不畫分支。字庫 bytes、大小與 hash 在執行前固定，沒有用輸出的圖形回填 C。
O0/O2 JSON 逐 byte 相同，SHA-256：`6ac33b1d26e622f54f6f3e5fe0af1dd9c835840dc50945db049d8f14bfebcf05`。

| 刻意改錯 | 首次拒絕 |
|---|---:|
| X 上界改 623 | glyph 第 3 組 |
| Y 上界改 383 | glyph 第 3 組 |
| 字模 word 不交換 byte | glyph 第 21 組 |
| 半形不回退 source byte | glyph 第 21 組 |
| 刪 dummy latch read | glyph 第 1 組 |
| 列距少一 byte | glyph 第 1 組 |
| 不保留背景略過 bit | glyph 第 6 組 |
| 三字不讀 patched live operand | text-name 第 10 組 |
| 第二字只前進 8 px | text-name 第 9 組 |

[平台 adapter](../../tools/c_recovery_glyph_platform.go) 把 C register 交給成熟 DOS/V
服務鉤子，再拿回服務 register；從未呼叫 `CPU.Step`。兩側獨立 machine／DOS／
cache／RAM／VGA，只有固定平台程式碼相同。初始化回入口、32-byte 緩衝和補第 16 列
保持該平台契約，磁碟 I/O 次數和原始 TSR 檔名差異不由此收據裁決。

C [glyph.c](../../tools/c_recovery/glyph.c)、[header](../../tools/c_recovery/glyph.h)、
[fixture](../../tools/c_recovery/glyph_fixture.h)、[生成 C](../../tools/c_recovery/glyph_generated.inc)、
[產生器](../../tools/c_recovery_glyph_generate.py) 與 [driver](../../tools/c_recovery_glyph.go) 進版控。
[重跑入口](../../tools/c_recovery_glyph.sh) 與 [驗證器](../../tools/c_recovery_glyph_verify.py) 核對
實際 binary flags、原始函式、平台 CS:IP、字庫與錯版；[公開摘要](c-glyph-verification.json) 保留雜湊。

Compiled source manifest：`workplace/matching-decompilation/c-glyph/results/c-source.sha256`，
SHA-256 `e121f9359837bcf6494fea4cae7314108ec1820e40264ee8946f0c9c3c2356f1`。
Go 1.26.7／GCC12.2.0／IDA9.4／非 root 限資源 Docker，原版、字庫與來源唯讀。
生成來源從固定 IDA 在乾淨容器重生逐 byte 相同。

三字範圍第一次截到 `mov bx,bp` 後，補原始四 pop與ret，完整輸出 208 原始指令為具體 C。
Go adapter 未用 import 刪除後同矩陣重跑；服務平台位址誤寫成 sub_10410 的導航名已修正
為實際 CS:IP。字形控制與資料名稱未改寫原始 DB。
新增 21 組語指令／59 bytes 與原 bytes 相同，其中包括兩個 patched far call；
保留原先 317 條／746 bytes 補充，合併為 24,714 條／56,197 bytes，整檔仍匹配。

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_glyph_probe.py workplace/matching-decompilation/c-glyph/ida KI.EXE
tools/c_recovery_glyph.sh
```
