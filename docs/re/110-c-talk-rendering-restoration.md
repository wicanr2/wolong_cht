# 110：C TALK、參數標記與肖像讀檔

**狀態：八個 C 函式／六個原始 handler 通過，O0/O2 各 1,277 全 RAM／VGA／ABI 相同，十錯版拒絕。**

- 日期：2026-10-09
- 輸入：DOS/V KI.EXE，67,099 bytes
- SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear，檔案偏移為線性位址減 `0x10000` 加 512
- DB：`workplace/matching-decompilation/c-talk/ida/input.exe.i64`
- DB SHA-256：`a163c3d4d90122c716b690e2af60b03d7bc60e6adaf79e9758650238d2e2f87e`
- 入口：[`ida_talk_probe.py`](../../tools/ida_talk_probe.py)
- 既有證據：[`re/79`](79-talk-marker-handlers.md)、[`re/33`](33-shared-draw-helpers.md)、[`re/109`](109-c-number-raster-restoration.md)
- 契約：[`spec/230`](../spec/230-c-talk-rendering.md)

還原 `sub_1075B` 的原始 ×8 索引、背景、肖像與字串呼叫，`sub_1084A` 的
逐字／換行／七項原始 near table，以及 `sub_107D2` 的四格 round-robin 快取。
字形和數字走既有真實 C，肖像走原始 DOS open／seek／read／close ABI。
參數 handler 保留原名稱與原位址，不推測原作者 C 函式邊界。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 原版 INT 50h 媒體錯誤介面 | 保留呼叫與重試分支；自然錯誤介面另有原始初始化依賴 |
| 自然玩家完整事件 UI | 此 renderer／caller 對照不代證完整輸入流程 |
| 原作者工具鏈與 C 機器碼 | 語意對照不證明 C 機器碼匹配 |

## 已證實的原始定位與資料流

| 原始函式 | IDA 線性位址 | bytes | 原始控制流與驗證邊界 |
|---|---|---:|---|
| `sub_106F9` | `0x106F9` | 4 | TALK 字串包裝 |
| `sub_1075B` | `0x1075B` | 119 | 八變體／TALK offset／肖像／背景／參數呼叫 |
| `sub_107D2` | `0x107D2` | 115 | 四格快取、miss 讀取與實際肖像 blit |
| `sub_1084A` | `0x1084A` | 90 | CF 碼寬、NUL 換行、live 七項 near table |
| `sub_18853` | `0x18853` | 48 | 狀態列與 FFFF 清除路徑 |
| `sub_189A4` | `0x189A4` | 58 | window map／背景 caller |
| `sub_1E38C` | `0x1E38C` | 26 | 讀檔包裝；INT 50h 分支保留，介面仍為證據界線 |
| `sub_1F4DF` | `0x1F4DF` | 73 | DOS open／seek／read／close，SS frame／CF／LAHF／SAHF |

原六個無函式邊界的 handler 在 `0x108B2`、`0x108DB`、`0x10904`、`0x10939`、
`0x1095B`、`0x1097E`，共 210 bytes；數字 `0x10984` 已由 re/109 還原。
八函式 533 bytes 加六 handler，共 743 bytes／376 原始指令。
索引為原始 CS:`0x8A4` 的七個 word，C 讀 live word，替換表後兩側也須同態。

TALK.DAT 34,182 bytes，SHA-256 `08a22e09791d0a6ec2968e87d8655e12c91b45e00fae460b28593b35ff85e384`。
KAOGRF.DAT 307,200 bytes，SHA-256 `b9c7745e3ed9b32f0c12003fe81f82756136f738427dc421ec96f6c4ed5c4ac8`。
四劇本 SINARIO.DAT 88,832 bytes，SHA-256 `21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c`。
三檔只作本機唯讀輸入，TALK offset、150 個 2,048-byte 肖像與四個 22,208-byte 區塊先驗形狀。

## 同狀態收據

原版跑 guest 指令，C 跑具體 native 運算／原位址 goto。兩側各有獨立 DOS／字庫 cache／
檔案 handle／RAM／VGA。C 僅把暫存器交給成熟平台 API，不執行原版 CPU。
比較全 1 MB RAM、四 plane、GC／seq／latch／port、十四 register／FLAGS、
新函式／handler／字形／數字入口與 32-byte stack、DOS API 前後 register，以及內容區像素。
IF／TF 固定為 0，未定義旗標沿用 dosgolem 模型，不宣稱實機時間。

| 群組 | O0／O2 各組數 |
|---|---:|
| 四劇本／七標記／直接與 FF ID | 56 |
| 半形／Big5 trail 5C／換行／替代 live table | 20 |
| 全部 TALK 槽與原始八變體 | 1,022 |
| 全部肖像 | 150 |
| 冷／熱快取替換序列 | 13 |
| 狀態列及清除 | 7 |
| 0／1／2048／4096-byte 讀取與開檔錯誤 | 9 |
| 合計 | 1,277 |

每個 TALK 槽都在原始 `sub_106F9` 入口核對實際 SI 與原 offset，1,022 槽各一次。
每種最佳化兩側各讀取 15,287 個全形／178 個半形字模，缺字 0。
初期把像素 X=80 當成文字格座標，造成多數 corpus glyph 被裁切；兩側雖然相同，
字庫工作量僅 209／76，不能作完整字形證據。改回原 caller 的格座標 0，再完整重跑。
驗證器以 944 個非空槽建立字庫工作量下限，並拒絕舊裁切工作量形狀，防止只靠 digest 宣稱完整。

快取序列 `0,1,2,3,0,4,1,5,2,6,3,7,7` 的 API 次數是
`4,4,4,4,0,4,0,4,0,4,0,4,0`。每個 miss 依序 open／seek／read／close，
hit 完全不讀檔。每版 open 1,195 次，seek／read／close 各 1,194 次；
差一是低階入口的受控不存在檔案，原 CF／error／stack 返回一致。
INT 50h 媒體重試 UI 沒有借這個成功路徑升級完成聲明。

| 刻意改錯 | 群組內首次拒絕 |
|---|---:|
| 變體 stride 少一個 shift | 415 |
| 半形不回退 byte | 1 |
| 換行只前進 8 px | 1 |
| 標記只前進 32 px | 13 |
| 呼び名 +8 改成姓名 +2 | 1 |
| 快取 cursor 只模 2 | 2 |
| 肖像位移少一個 shift | 2 |
| 肖像只讀 1,024 bytes | 1 |
| 控制標記不抵銷 cursor | 13 |
| 變體不加 AH | 408 |

[公開摘要](c-talk-verification.json) 記錄來源、輸入與收據雜湊。
來源：[talk.c](../../tools/c_recovery/talk.c)、[header](../../tools/c_recovery/talk.h)、
[生成 C](../../tools/c_recovery/talk_generated.inc)、[產生器](../../tools/c_recovery_talk_generate.py)、
[driver](../../tools/c_recovery_talk.go)。[重跑入口](../../tools/c_recovery_talk.sh)／
[驗證器](../../tools/c_recovery_talk_verify.py) 檢查來源 digest、編譯 flags、字形工作量、錯版與 backlink。
Go 1.26.7／GCC12.2.0／IDA9.4；dosgolem revision `a9714ebdab2ad6b529f81225472680f2b11f2842`，
輸入與外部平台來源唯讀、非 root 限資源 Docker。生成 C 從固定 IDA 在乾淨容器重生相同。

原探針只以 basename 驗證語意來源，上一輪 docs/spec/229 因而無法解析。
`ida.sh probe` 現在唯讀掛載完整 docs，`ida_matching_probe.evidence_path` 依完整路徑查證。
同名錯目錄、絕對與越界路徑有拒絕對照；不退回忽略出處的匯出。

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_talk_probe.py workplace/matching-decompilation/c-talk/ida KI.EXE
tools/c_recovery_talk.sh
```
