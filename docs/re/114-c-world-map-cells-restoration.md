# 114：原始大地圖顯示格與圖塊合成 C

**狀態：CONFORMED。10 支函式的局部原版／C 行為通過。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 固定 IDA DB SHA-256：`3e018b73e908af3f2da654dffd2fd050f063654442819fb153340941f736c40c`
- 工具／位址：IDA Pro 9.4 database linear；file offset 為 linear 減 `0x10000` 加 512
- 入口：[ida_scene_probe.py](../../tools/ida_scene_probe.py)
- 前置：[re/72](72-world-map-display-list.md)、[re/113](113-c-resource-cue-restoration.md)
- 契約：[spec/234](../spec/234-c-world-map-cells.md)

本輪由事件場景的 `sub_11D46` 退出重畫鏈定位顯示格閉包。原 `sub_13D09`／
`sub_13D45`／`sub_13D68` 已有 C；完整場景主入口仍依相依圖判讀。
`sub_1D46A` 設定 MDL／MCH、顯示記錄與來源地圖段；`sub_1D483` 初始化記錄；
`sub_1D615` 讀來源圖塊並設髒格；`sub_1D4C7` 裁切及推入最多四張覆蓋圖。
`sub_1D66A` 對 40×23 個記錄處理保護、快取、底圖、覆蓋圖與 `0xFF` 圖塊。
其五個繪製 helper 保留 128-byte 底圖、160-byte mask／color 與四個 VGA plane 的原始順序。

上述控制流以固定原始指令為證據，推論等級為已證實；局部 C 行為由下列收據確認。
IDA loader 的重定位 bytes 與 file bytes 分列，未將兩個位址空間混用。

## 實作與驗證

10 支原始函式共 710 bytes／328 指令，名稱、位址與運算元保留。
來源為 [mapcells.c](../../tools/c_recovery/mapcells.c)、[生成內容](../../tools/c_recovery/mapcells_generated.inc)。
重跑入口 [c_recovery_mapcells.sh](../../tools/c_recovery_mapcells.sh)，
[收據](c-mapcells-verification.json)記錄所有 source／工具／素材雜湊與實際編譯 flags。

| 矩陣 | 每個最佳化版本的組數 |
|---|---:|
| 初始化 | 12 |
| 真實來源地圖取樣 | 10 |
| 256 底圖與安全 DF 分支 | 768 |
| 256 遮罩覆蓋圖與 DF 分支 | 512 |
| 四平面 flush／位移邊界 | 16 |
| 256 旗標 × 0／1／4／5 個記錄 | 1,024 |
| 推入裁切／保護／容量／編號 | 256 |
| 完整 920 格重畫 | 8 |
| 初始化 → 地圖 → 推入 → 重畫 | 16 |

O0／O2 各 2,622 組完整 RAM／四 plane／ABI／FLAGS／I/O／入口 snapshot 相同。
另外 256 次直接 blit 與 512 次 mask／color 合成獨立比真實來源 bytes。
遮罩測試先放入真實 MDL 背景，避免零底圖遮住遮罩錯誤。12 個 O0 錯版全數拒絕，
原始 C 在乾淨容器重生相同，實際編譯 flags 綁定完整來源 digest。

原版最大指令步數與 400 萬上界同列於收據。16,384 個 trace 槽覆蓋最長的完整畫面，
沿用既有 VGA snapshot，沒有截短 trace 或原迴圈。

`sub_1D7E7` 的 DF=1 會由 `CS:0xD858` 反向寫入 `0xD7DA..0xD859`，
覆蓋自己的整支函式，原 guest 隨後改走損毀指令。原 renderer 先 `CLD`，
故局部 helper 契約為 DF=0。保留 [受控實驗](c-mapcells-df-precondition.json)；
這個破壞程式碼的輸入不被改寫成測試通過。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 完整場景退出與 producer | 小地圖、軍團／物件推入端及正常玩家流程仍需後續接線 |
| 原作者 C 工具鏈／機器碼 | 尚未確認 |

## 後續原始 C 場景退出

[re/116](116-c-scene-resume-restoration.md) 以真實原主入口接完整局部
world／小地圖／Mouse／TALK／resource，保留既有唯一 C world resume。
自然事件與長程玩家路徑仍按該收據界線判讀。
