# 39 — 剩餘未讀函式的逐支歸屬

**狀態：清單。**0 支**——全部讀完了。
每一支都已有叢集歸屬與呼叫證據（模組全圖 [`35`](35-strategy-ui-module-map.md)–[`38`](38-strategy-core-module-map.md)），這一份是那份清單的生成器輸出，現在是空的。**

- 日期：2026-08-14
- 範圍：只驗松崗 DOS/V
- 原始 `KI.EXE` SHA-256：
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具：`tools/ida_module.py`（IDAPython）＋ `tools/re_coverage.py`
- 位址空間：IDA DOS/V linear address，segment base `0x10000`

這一份是**生出來的不是手打的**（`tools/re_coverage.py` ＋
`tools/ida_module.py` 的 JSON），所以可以重跑。

> ⚠ **它會把自己算進覆蓋率。** 逐支列出未讀函式的文件正是
> [`21`](21-function-census.md) §3.1 講的那種——重跑覆蓋率時要一併排除，
> 否則 T4 會歸零而實際上一支都沒多讀。
> `tools/re_coverage.py` 的 `CATALOGUES` 已經加了這一份。

「證據」欄的內容來自它呼叫的共用常式與 I/O 埠；「父」是它的直接呼叫者。
**父已被記錄的葉節點，讀起來最便宜**——上下文已知。


## 清單是空的

 目前量到 **T4 ＝ 0**：739 支函式每一支都有
`docs/re/` 層級的記錄。收尾的路徑是模組級掃描補歸屬
（[`35`](35-strategy-ui-module-map.md)–[`38`](38-strategy-core-module-map.md)）
再逐行讀葉節點（[`42`](42-leaf-functions.md)）。

> **空清單不代表沒有缺口。** 真正的缺口在各文件的「未解」表裡——
> 那些是「讀過但沒解出語意」的項目，不會出現在覆蓋率上。

重跑方式見下。清單再度變長，就表示 census 重新分析出了新函式。

## 怎麼重跑

```sh
tools/ida.sh script dosv tools/ida_module.py KI.EXE.i64
tools/py.sh tools/re_coverage.py workplace/ida/dosv/census/census.tsv
```
