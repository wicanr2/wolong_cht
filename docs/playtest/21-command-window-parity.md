# 21 — 松崗 DOS/V 指揮／事件／一覽畫面 parity 重開

**狀態：歷史量測紀錄（2026-08-12，影片幀對幾何）。** 當時列為未完成的
**一覽詳細層與捲軸都已收掉**：四個家族的欄位、標題、視窗矩形、捲軸的四個熱區
與「委任」格都照原版（[`../spec/38`](../spec/38-list-windows.md)，CONFORMED，
`TestListScrollbarMatchesRawHotzones`）。
⭐ **這一份的量法已被取代**——它拿的是壓縮影片的代表幀、縮到 320×200 比幾何；
現在有同狀態逐像素對拍（[`37`](37-main-screen-parity.md)–[`40`](40-tactical-parity.md)）。

- 日期：2026-08-12
- 原版來源：使用者松崗 DOS/V 錄影代表幀 `20s`、`240s`、`550s`
- 正規化：影片遊戲本體裁切後最近鄰縮放為 320×200。

| 畫面 | 原版 logical rect | remake rect | 結果 |
|---|---:|---:|---|
| 一覽後層清單 | 約 `49,41..273,161` | `98,82,448,...` | 第一層主要外框已對齊 |
| 一覽前層詳細窗 | 約 `97,62..271,193` | 尚無同等呈現 | 未完成 |
| 一覽左側捲軸 | 約 `67..84` | 尚無 | 未完成 |
| 左下事件 TALK | 約 `0,161..128,200` | `0,320,256,80` | 已對齊，右側 HUD 保留 |
| 中央系統面板 | 約 `104,57..208,154` | `208,114,208,194` | 已對齊主要外框與五列值格 |

同類並排證據：

- [`parity/message-layout-side-by-side.png`](parity/message-layout-side-by-side.png)
- [`parity/system-layout-side-by-side.png`](parity/system-layout-side-by-side.png)
- [`parity/list-layout-side-by-side.png`](parity/list-layout-side-by-side.png)

先前推廣片把原版事件／一覽畫面與 remake 目的地畫面並排，屬異類比較，不能支持 parity；該鏡頭必須撤換。
