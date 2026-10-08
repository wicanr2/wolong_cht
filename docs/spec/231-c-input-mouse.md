# 231：C 原始訊息等待與滑鼠游標

**狀態：CONFORMED。O0/O2 各 210 全裝置／ABI 相同，三個背景還原與十二錯版通過。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear
- 證據：[`re/111`](../re/111-c-input-mouse-restoration.md)

## 原始契約

保留 main／mouse 兩個 segment 的 near／far ABI、CLI／STI、PUSHF／POPF、live 16 項
mouse table、press／release 次數、游標兩層 mask、保存／恢復、GC／seq／latch／IN／OUT。
mouse reset 及 range／handler 登記用成熟 dosgolem INT 33h API，C 不跑 guest CPU。

等待讀取 live CS:`0x2239`，只接受原 `EB 0F` 與 caller patch 的 `90 90`。
原模式跳過左鍵，兩模式均可右鍵結束；counter 每到 `0xFF` 重設並將 SI 從 10 減一。
輸入在 API／counter 比較的可重播邊界送入，兩側獨立平台服務與狀態，不挑結果重送。

## 驗證閘門

O0/O2 比全 1 MB RAM／四 plane／GC／seq／latch／讀寫 port、十四 register／FLAGS、
完整 near／far 與 API 快照。覆蓋 16 項服務、cursor show／hide／move／callback、八種 X 對齊、
右／下邊界、保存／恢復序列、左右鍵與無按鍵 timeout、patch 還原、訊息繪圖等待擦除。
callback 明示 far 輸入，不使用 C-side CPU.Step；原版 timer wall-clock 不由控制輸入代證。

## 未解範圍

| 項目 | 限制 |
|---|---|
| counter 的真實初始化／timer cadence | 只驗原 consumer，不推算硬體時間 |
| 正常玩家 callback／訊息流程 | 局部完整 device 對照仍有 direct-entry 界線 |
| C 機器碼與原作者工具鏈 | 尚未確認 |
