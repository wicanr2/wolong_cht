# 239：C 四類清單與原始回呼

**狀態：CONFORMED。O0／O2 各 804 完整裝置與十三錯版通過。**

- 日期：2026-10-09
- 證據：[re/119](../re/119-c-list-families-restoration.md)。
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。

## 契約

保留軍團、武將、勢力與開局清單的七個 caller、builder 與 renderer。
原始 126／128／22 筆邊界、SS 記錄位移、X/Y patch、DS:0CFF、資格與哨兵均須一致。
武將入口沿用排序狀態，軍團士氣讀 byte，勢力數字保留原三位呼叫。
外交 cache 依 LOOP 目標每筆清 DL，不猜補敵情或新增排序安全分支。
共用清單、排序、特殊返回、滑鼠、字型與 VGA 使用既有唯一 C 實作。

## 驗證閘門

- O0／O2 比完整 RAM、四 plane、暫存器／FLAGS、SS frame、呼叫與裝置軌跡。
- 四劇本與多勢力；builder 獨立核對、各欄排序、真實表頭點擊、選取與取消。
- 同座標敵我軍團、士氣／兵力換色、外交 cache 邊界、武將原屬／現屬與職務優先序。
- 空列與空清單的合法狀態；原下溢不當成安全修正。
- 完整真實字型與素材雜湊、零缺字、編譯來源綁定、錯版拒絕、乾淨 C 重生。
- 全部原指令比對組語覆蓋，新增未知區以固定 binutils 匹配。
- 既有 Go 包裝器使用原始素材唯讀、專用鎖版快取與有界 Xvfb，正式檢查冷跑。

## 實作與收據

`tools/c_recovery/catalog.c` 與 `catalog_generated.inc` 接回唯一共用清單引擎。
獨立資料參照在 `tools/c_recovery_catalog_data.go`。
[版控收據](../re/c-catalog-verification.json)記錄全 RAM／plane／FLAGS／SS frame／裝置、
168 builder／240 排序／16 cache、156 表頭與 28 成功選取；五條舊證據回鏈自動檢查。
正式 Go 引擎未更動；[冷測收據](../re/c-catalog-go-verification.json)記錄 vet／test、鎖版依賴與素材唯讀掛載。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 非法排序狀態與自然玩家長程 | 不由合法局部案例代證 |
| 原作者 C 機器碼 | 尚未匹配 |
