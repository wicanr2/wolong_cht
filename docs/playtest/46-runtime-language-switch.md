# 46 — 執行期切換語言：桌面與手機的實跑

**狀態：confirmed。** 四個語言在**桌面啟動殼層**與**手機系統面板**都切得動，
畫面上各自用自己的文字顯示。F9 熱鍵切到的畫面與 `-lang` 啟動的畫面逐像素相同。

- 日期：2026-08-26
- 對象：remake（`cmd/wlgame`、`cmd/wlandroid`），不涉及原版
- 規格：[`../spec/86`](../spec/86-runtime-language-switch.md)

## 1. 桌面：啟動殼層的 `LANGUAGE`

![桌面語言頁](../images/launcher-language-page.png)

```
WOLONG_SHOT_CMD=wlgame tools/shot.sh out.png "KEYS=Down,Return" \
    -orig workplace/orig/dosv -shot-frames 90
```

四列各自用自己的寫法寫，目前語言前面一個 `●`。

## 2. 手機：系統面板的第三頁

![手機語言頁](../images/phone-language-page.png)

```
WOLONG_SHEET=3 WOLONG_TAB=2 tools/phone_shot.sh out.png 60
WOLONG_LANG=en tools/phone_shot.sh out-en.png 60
```

![手機英文主畫面](../images/phone-english-main.png)

手機的 HUD 與指令列一起換（`Funds`／`Res.`／`Prop`／`List`／`Corp`／`Sys.`），
日期欄變成 `Y196 M4 D1`。

## 3. 兩個字型陷阱

| 症狀 | 成因 | 處置 |
|---|---|---|
| 「简体中文」的「简」是方框 | 「简」不在 Big5，而選單當時掛的是倚天 Big5 | `fontChain` 每個語系都把另外兩套字集接在後面墊底 |
| 選中記號是方框 | `▶`／`✓` 都不在倚天 Big5 | 換成 `●` |

⚠ 手機端還有第三個：**母本繁中原本不走 `SetLanguage`**，於是 `LangPack()` 是 nil、
字型停在單一倚天 Big5，墊底鏈根本沒掛上。**「預設語言不需要載語系」這個省略
會讓修好的字型鏈在預設狀態下失效**——四個語言一律走同一條載入路徑才擋得住。

## 4. 未解

| 缺口 | 下手點 |
|---|---|
| Android 實機／模擬器沒實地切過 | 這裡的手機畫面是桌面 Xvfb 跑 `cmd/wlandroid` 拍的（同一份 `internal/ui/phone`）；實機驗收排在下一次 Android 打包 |
