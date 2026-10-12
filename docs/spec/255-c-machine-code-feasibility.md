# 255：C 機器碼匹配可行性判定

**狀態：CONFORMED。判定完成：`KI.EXE` 為手寫 16-bit 組語（[`re/135`](../re/135-toolchain-fingerprint.md)），C 機器碼匹配不可行，`c_machine_code_match` 維持關閉。同日 IDA 9.4 以 FLIRT＋正對照覆核確認（`re/135` §7）；Ghidra 本環境無，未跑。其餘七支同日判完（[`re/136`](../re/136-other-binaries-toolchain.md)：3 Borland＋4 手寫 ASM），本 spec 的「只限 KI.EXE」範圍維持，不擴張。**

- 日期：2026-10-12。
- 出處：試點 [`re/89`](../re/89-matching-decompilation-pilot.md)（工具鏈未知、GCC `-m16` 0／12）、[`reference/04`](../reference/04-first-survey.md) §2.2（三支 Borland、一支不確定，不外推）、本輪 prologue 普查（見下）。
- 輸入：DOS/V `KI.EXE`，SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。

## 初步觀察（假說，非結論）

- 739 支函式起點 bytes 普查（檔案 bytes，`linear-0x10000+header`）：`50 53`（push ax／bx）77 支、`1E 06`（push ds／es）45 支、`1E 50` 38 支；`0x55` 開頭僅 8 支，其中 6 支是 `push bp；sub sp,N；mov bp,sp` 手寫 frame（`17FDB／1820E／19479／194BF／1F4DF／1F593`，暫存器傳參＋`[bp+N]` 區域變數＋`add sp,N` 收尾），2 支只是 push bp 存暫存器。
- 單 chunk 函式尾 byte：717／737 支 `C3`（retn）；唯一 `CB`（retf）是驅動段 `sub_20000`（far-call 目標，本來就該 retf）；多 chunk 函式 2 支（皆為緊鄰 4-byte 跳過 stub，非遠端 tail）。
- 已知手寫徵候：暫存器呼叫慣例（全庫）、DS／CS 全域、自我修改碼 2 處（[`re/28`](../re/28-text-number-rendering.md)：`loc_1A065`、`loc_10701` 的 `mov cs:[…]`）、函式中段入口（`sub_1C30D`／`sub_1C300`）、retf tick 回呼、中斷驅動常駐。
- 反例待查：6 支 frame 函式是否為 C 混編；有無執行期函式庫（32 位元乘除、字串、檔名；目前所見皆 inline 或直呼 INT21）。

## 判定程序

1. **指紋表**：startup 形狀（`start` 0x10000：遊戲開機流程，非 C runtime 形）、6 frame 函式體、MZ 標頭／stub／重定位風格、全庫 `int`／`in／out`／`rep` 分佈、疑似函式庫的候選（高 fan-in 共用常式）。
2. **編譯器家族測試**（若有家族證據才做）：取候選編譯器（Borland C++ 3.x／Turbo C／MSC 6-8／Watcom 9.x 之一），寫最小探針（frame、呼叫、除法、字串），比 codegen 習慣；習慣不合即停。
3. **試點匹配**（若家族成立才做）：挑一支小函式，寫 C 迭代到逐位元組相同；記錄編譯器版本、選項、來源寫法。
4. **證否線**（與 2 並行）：列出「任何 C 編譯器都不會產生」的結構（中段入口、自修改、`xchg dh,bl` 類 byte 雜耍、retf 回呼、CS 跳表計算），逐項定位；若遊戲核心充滿這類結構，判手寫。

## 決策準則

- 判「編譯器產物」：需 startup ABI ＋ codegen 習慣 ＋ 至少一支試點匹配三者齊備；版本需版本字串或 release-specific bytes，否則只定家族。
- 判「手寫組語」：需指紋表 ＋ 證否線 ＋ 試點編譯差距三者齊備；結論只限 `KI.EXE`，不碰其他七支。
- 無論哪邊，`c_machine_code_match` 的旗標語意不變：只有逐函式機器碼收據能翻轉它。

## 實作與驗收

- 產出：指紋掃描工具（`tools/tc_fingerprint.py`，讀 census＋原檔 bytes，不跑 IDA，全綠）、判定筆記 [`re/135`](../re/135-toolchain-fingerprint.md)、本 spec 升 CONFORMED（判定完成即 CONFORMED，不是匹配完成）。
- 驗收：每個斷言帶位址＋bytes＋推論等級；試點編譯（若有）可重跑；不把「像手寫」寫成已證實。
- 原版素材唯讀，不入 Git；候選編譯器若需取得，走使用者授權。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 其他七支執行檔的工具鏈 | 本輪只判 `KI.EXE`。 |
| 原作者語言（若手寫） | 組語來源無法取回命名與切分；判定止於「非編譯器產物」。 |
| 完整匹配 | 判定通過才開下一輪；本 spec 不承諾匹配。 |
