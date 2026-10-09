# 117：C 君主出陣與自動編成

**狀態：CONFORMED。14 個 C 函式、501 條原始指令通過局部同狀態驗證。**

- 輸入：松崗 DOS/V `KI.EXE`，67,099 bytes。
- SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 日期：2026-10-09。
- IDA 資料庫 SHA-256：`0fa4e3cbe8e54c08aad40f9ece7aa4de0ea8bdc6e2f163fd2896427953c06fb3`。
- 工具與位址：IDA Pro 9.4 database linear，段基址 `0x10000`。
- 匯出：`tools/ida_verdict_probe.py`，本機證據在 `workplace/matching-decompilation/c-verdict/ida/`。
- 規格：[spec/237](../spec/237-c-ruler-sortie.md)。

## 已確認的控制流

| 原始定位 | 控制流 | 證據等級 |
|---|---|---|
| `sub_1699E` | 保留資金符號分支、16-bit 加總 carry、600 點門檻、TALK 基址 `0x18C`，接回原場景、編成、側欄與重畫 | 已證實，固定輸入原始指令 |
| `sub_16E8F`／`sub_16EC9` | 依原 `CS:6C4C` 表逐槽選兵種，每槽要求 50 點；失敗保留已寫槽與 CF | 已證實，固定輸入原始指令 |
| `sub_16F26`／`sub_16F86` | 建立軍團、複製首都座標、更新佔用圖與士氣 | 已證實，固定輸入原始指令 |
| `sub_1461D`／`sub_14717`／`sub_14698` | 先退回舊槽兵力，再分配預備兵；保留除法餘數加入與每槽 100 點上限 | 已證實，固定輸入原始指令 |
| `sub_16FD2` | 加總六槽，保留原 `cmp bx,12Ch` 與旗色計算，不把 BX 比較改解成總兵數比較 | 已證實，固定輸入原始指令；設計意圖未知 |
| `sub_15E80` 與四個 callee | 保留側欄開關、四個 bit、原版名稱／數字／兵力 raster | 已證實，固定輸入原始指令 |

既有編成證據見 [re/30](30-corps-formation-ui.md)、[re/49](49-corps-formation-window.md)，
完整場景依 [re/116](116-c-scene-resume-restoration.md) 的唯一 C 實作。

## 驗證結果

O0／O2 各 451 組完整 RAM、四 plane、暫存器／FLAGS、堆疊、呼叫入口與裝置 API 相同。
包含 64 次完整出陣 caller、160 組編成、96 組退兵／補兵、128 組側欄及 3 組 BX 邊界。
十個編譯錯版全被拒絕；固定 IDA 證據乾淨重生 C 一致。
來源與工具雜湊、實際入口次數、原始輸入與工具版本見
[版控收據](c-verdict-verification.json)，重跑 `bash tools/c_recovery_verdict.sh`。
原資金符號例外已回填 [AI 機制](../mechanics/70-ai.md)。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 原作者 C 與工具鏈 | 原始語言及 compiler 未確認，尚未證明 C 機器碼一致 |
| ~~遷都清單~~ | 原 caller、清單與排序已由 [re/118](118-c-city-list-restoration.md) 完成局部 C 驗證 |
| 正常玩家長程路徑 | 本輪局部主入口驗證不代證自然選單與通關 |
