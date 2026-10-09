# 237：C 君主出陣與自動編成

**狀態：CONFORMED。O0／O2 各 451 組完整裝置案例與十個錯版通過。**

- 日期：2026-10-09

證據：[re/117](../re/117-c-ruler-sortie-restoration.md)。
輸入 `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。

## 契約

保留 14 個原始函式的名稱、線性位址與運算元。C 實作原指令的固定控制流，
沿用既有原版場景、TALK、mouse、VGA、數字與資源實作。
資金符號、carry、600／50／100 門檻、原始候選表、失敗的部分寫入、
LAHF／SAHF、佔用圖、側欄 bit 與 `cmp bx,12Ch` 均須保持原意。

## 驗證閘門

- O0／O2 比較完整 RAM、四 plane、14 暫存器與 FLAGS、SS frame、呼叫入口與裝置 API。
- 四個原始劇本；出陣接受／拒絕、資金符號、兵力邊界與加總 carry。
- 自動編成六槽成功／中途失敗、三類預備兵、補兵除法餘數及上限。
- 側欄開關與四個 bit；真實名稱、數字與圖庫，不替換成 UI stub。
- 編譯時綁定來源雜湊；錯誤實作必須被拒絕；固定 IDA 證據可乾淨重生 C。

## 實作與收據

原始 C 在 `tools/c_recovery/verdict.c` 與 `verdict_generated.inc`，
沿用唯一場景／資源／圖形與 `sub_155EC` 實作。
[版控收據](../re/c-verdict-verification.json)綁定來源、IDA DB、編譯 flags 與獨立原版 oracle。
本輪未修改 Go 規則。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 正常玩家長程路徑與 C 機器碼 | 本輪不據局部通過宣稱完成 |
| ~~遷都清單~~ | [spec/238](238-c-city-list.md) 已完成原 caller 與清單局部 C 驗證 |
