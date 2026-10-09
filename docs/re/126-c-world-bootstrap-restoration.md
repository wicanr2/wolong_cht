# 126：城市標記、軍團佔用與主迴圈初始化的 C

**狀態：CONFORMED。五函式及有限前綴，O0／O2各5,008組、5,416次全狀態與十錯版通過。**

- 日期：2026-10-09
- 範圍：松崗 DOS/V KI.EXE；SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 工具：IDA Pro 9.4；位址均為 IDA database linear，段基址 10000。
- Database SHA-256：`3792dd686a7aba92be5adac03a996bfbf54caac3e5f64e36ec5fe65ef7ea9750`。
- Probe SHA-256：`a35b6c350b0e83a36843351f25cf563722dced75ed1324a8f6f47bba592f8d91`。
- 探針：[ida_bootstrap_probe.py](../../tools/ida_bootstrap_probe.py)；前輪 frame：[re/125](125-c-nonlocal-exit-restoration.md)；規格：[spec/246](../spec/246-c-world-bootstrap.md)。

## 已證實原始契約

| 入口 | 契約 |
|---|---|
| `sub_189F0`，IDA `0x189F0`，46 bytes | 無城市存在檢查，順序處理 192 城；再處理軍團槽 0–126，共 127 槽。保存 DS／ES／AX／BX／CX／SI。 |
| `sub_18A1E`，179 bytes | CS:D44 對應 decoded map；城市中心為 Y×384+X，原 tile 減 CB 後以 byte 除 3、乘 3、加 CB，再加玩家 0／他方 1／中立 24 的色碼 2。 |
| `sub_18AD1`，25 bytes | 只改 DE..F1；折回 DE..E7，再依 AH 加 0 或 10。其他 raw byte 不寫。 |
| `sub_18AEA`，40 bytes | flags≥80 才重寫軍團 +1A=X、+1C=word(occupancy base+Y×24)，然後增加該 segment:[X] 的 byte；不清佔用表，重跑會累加／回繞。 |
| `sub_1533D`，27 bytes | 入口 DS:D52 提供 world；共用 `sub_155A6` 重算武將評分，再由既有 `sub_11E17` 畫日期。退出 DS=CS，不能假設恢復入口 DS。 |

五函式共 140 指令／317 bytes。四鄰格依城市 +16 的低 nibble：
0 使用中心 ±302／±2FE，3 使用 ±180／±1，其他使用 ±181／±17F。
邊緣運算保留原 segment 與 16-bit offset，不自行加座標保護。

Occupancy 使用 384×256 byte 的執行期緩衝區。原範例基址為 B000，
這是 dosgolem 記憶體實址 B0000..C7FFF，與上述 IDA 程式位址分開。
軍團槽 127 不由 189F0 處理；單獨呼叫 helper 的有效記錄仍依其 flags 判斷。

## 有限前綴與完整主迴圈的界線

`sub_11BE0` 的 `[0x11BE0,0x11BF6)` 為七指令／22 bytes：設 DS／ES、
呼叫兩個初始化 wrapper，再保存 SS／SP。這段沒有 RET。
其 SHA-256 為 `0c5138467de03c907f023374d58d25b2f72aec8cb487581777966c1f39d2d07e`。
它另列原始程式碼片段，不計入完整 C 函式；一般函式分派不把它當成已完成的 11BE0。

- 原版從真正 CALL10067 開始，到 11BF6 停止；C 由同一份 raw RAM／獨立 VGA 狀態執行相同前綴與完整 helpers。
- C 前置階段不使用原版執行後的 Snapshot 作初始輸入。
- 原始資產先核對，城市標記、occupancy、descriptor、評分與日期由兩側各自產生。
- 同格多軍團、FF 計數、poison descriptor、active 7F／80、槽 126／127 與重複初始化是必要分支。
- 原生初始化後再進入既有信賴度／進言 caller，核對正常返回及真正恢復 outer frame 的退出鏈。

## 驗證結果

- O0／O2 各 5,008 組：city 964、fringe 2,048、army 1,052、bulk 60、score 464、prefix 12、chain 408；每版 5,416 次全狀態比對。
- 原版先經獨立 byte 模型，再比完整 RAM／四 plane／DAC／暫存器／FLAGS／SS frame／IN／OUT／API；十個實編譯錯版均由狀態比較拒絕。
- 軍團有效邊界、同格累加／FF回繞、poison descriptor、126／127、重複 bulk 與合法 alternate DS 均已驗。D44／D850 不同、occupancy=C000 的反例能分辨真正位址讀取端。
- C 從兩側各自的 raw 初態計算所有初始化資料，沒有使用原版已執行 prefix 的 Snapshot；另驗 native prefix→既有 C 退出鏈。
- [收據](c-bootstrap-verification.json)、[指令覆蓋](c-bootstrap-code.json)、[Go 冷測](c-bootstrap-go-verification.json)保存來源與範圍；重跑 `bash tools/c_recovery_bootstrap.sh`。
- 147 指令／339 bytes 獨立重組精確且無新增組語；五個完整函式與 22-byte 無 RET 前綴分開計數。整檔 67,099-byte EXE 仍逐 byte 相同。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 11BF6 之後的完整主迴圈與自然玩家路徑 | 七指令切片不代證 |
| 非法座標、原硬體時間與原 C 機器碼 | 本輪不推廣其驗收聲明 |
