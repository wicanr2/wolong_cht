# 229：C 原始數字 raster

**狀態：CONFORMED。O0/O2 各 4,881 組完整 RAM／VGA／ABI 相同，八個錯版拒絕。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear
- 證據：[`re/109`](../re/109-c-number-raster-restoration.md)

## 原始契約

還原 `sub_1062F`、`sub_1069A`、`sub_106DE` 的原始 register／FLAGS、
stack、segment、DF、顏色與 VGA latch／port 副作用。
字模取自本機唯讀 ICONGRF.DAT 第三段加 `0x840`，CS:`0xD54` 為實際 bank。
日期及 SS 參數 code 入口呼叫真實數字 C，輸出不得由原版畫面回填。

第一個 DIV 的 quotient 必須可放入 AX，數值絕對值至多 655,359。
BL 為原始 byte 寬度，保留 0 的回繞與負數寬度 1 的保護分支。
直接 digit／blank 入口另外驗證 DF 與 CX 高 byte，主入口則自行 CLD。

## 驗證閘門

O0/O2 比較全 1 MB RAM、四個 plane、GC／seq／latch／port、ABI 與呼叫入口快照。
覆蓋正負值、十進位邊界、寬度／截斷／背景、256 種顏色、word 定址回繞、
直接 primitive 與兩個原始 caller。錯版必須被拒絕，生成 C 與實際編譯來源必須相符。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 原作者編譯器與 C 機器碼 | 尚未確認 |
| 正常玩家 UI／訊息參數 producer | 此局部驗證不代證 |
| DIV 例外 | 只驗證可執行數值範圍，不宣稱原版接受任意 signed 32-bit |
