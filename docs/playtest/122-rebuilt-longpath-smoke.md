# 122：重建 EXE 從零開局長程實跑（matching slice A）

**狀態：通過。** 冷重建 `KI.EXE` 在 dosgolem 下從開機跑到 196 年 4 月 8 日，
8 張截圖、時鐘、累計指令數與原版側全同。

- 日期：2026-10-12。
- 規格：[`../spec/256`](../spec/256-matching-longpath-verification.md) slice A。
- 原版側：`tools/dosgolem.sh /tmp/lp-orig-full "<時間軸>"`（`GAMEDIR=orig/dosv`）。
- 重建側：`WOLONG_DOSGOLEM_GAMEDIR=longpath-rebuilt`（`workplace/longpath-rebuilt/`：
  複製 `orig/dosv` 69 檔，以冷重建 `KI.EXE` 蓋掉一支；`cmp` 與原版全同，
  不入版控）。
- 時間軸（兩側同一串）：
  `wait；shot:boot；click:320,171；wait；shot:scenario；click:315,151；wait；`
  `shot:faction；click:320,144；wait；click:320,152；wait；shot:confirm；`
  `click:337,252；wait；shot:naming；click:495,199；wait；shot:advisor；`
  `click:351,279；wait；rclick:320,350；wait；until:196/4/6；until:196/4/7；`
  `until:196/4/8；shot:later；clock`
  （開機 NEW GAME → 第一章呂布歸天 → 孫策／周瑜 → 命名確定 →
  開局 → 關進言對話框 → 跑到 4/8）。

## 比對結果

| 取樣點 | 遊戲時點 | 原版／重建截圖 SHA-256 |
|---|---|---|
| boot（NEW GAME 對話框） | 開機 | `4233ec31…0809d84` 全同 |
| scenario（四章選擇） | — | `b0a2e32a…0fe2e8b` 全同 |
| faction（勢力一覽） | — | `681eb292…a8a2fe` 全同 |
| confirm（君主軍師確認） | — | `ee0f6260…1827a` 全同 |
| naming（軍師命名） | — | `e16a5552…4942be` 全同 |
| advisor（命名後確認） | — | `b512bf26…0b9ef23` 全同 |
| map（開局大地圖＋進言） | 196/4/4 19時 | `9ab876ac…51a14a8` 全同 |
| later（同一地圖續跑） | 196/4/8 1時 | `c08fedf6…14c3925c` 全同 |

- 雙側最終累計指令數同為 130,849,387；`clock` 輸出同為 196 年 4 月 8 日 1 時。
- 完整雜湊在 `/tmp/lp-orig-full/`、`/tmp/lp-orig-L3/` 與重建側對應目錄
  （`cmp` 逐張比過，不只比雜湊）。

## 過程記錄（後人不用重踩）

- 開機直達 NEW GAME 對話框，不經密碼頁（dosgolem 路徑下如此；密碼頁結論見
  CLAUDE.md §4.0，不受影響）。
- 命名畫面的確定鈕是熱區 #32（`x 464..527，y 192..207`），不是目測的
  `(449,216)`（那是注音格，會翻頁）。
- 君主軍師確認面板上排的繼續／確定不是活熱區；活的是下排 #32/#33，
  開局按 `(351,279)`。
- 開局後的進言對話框擋住時間：`until` 之前先 `rclick:320,350` 關掉，
  否則 `until` 跑滿預算也到不了。
- `until` 單步預算 40M 道指令不夠走 4 天（4/5→4/9 失敗兩次）；
  改鏈式 `until:4/6；until:4/7；until:4/8` 即過。每步預算獨立。
