# 事件 2–5 TALK 完整分支抽樣

**狀態：36 個 raw TALK 頁面、18 組雙頁回應的分支、marker、硬換行、字寬與五列版面
抽樣通過；不宣稱完整自然劇本長程或原版逐像素 parity。**

- 日期：2026-08-11

本輪以 DOS/V `TALK.DAT` 與 `translations/talk-dosv-corrected.json` 為輸入，
把事件 2／3 的外交 prompt、三選一、成功／指定金額／拒絕／超額回應，以及事件
4／5 的撥款 prompt、三選一與六種結果全部走過 runtime modal。

## 可重跑驗收

```text
Xvfb :99 -screen 0 640x400x24 -nolisten tcp
DISPLAY=:99 go test -p=1 -vet=off ./cmd/wlgame \
  -run 'TestEvent2To5FullTalkPageSampling$' -count=1 -v
```

實際使用的是 `demonwinter-go:latest` Docker 容器、DOS/V 原始素材唯讀掛載、
倚天字型與 640×400 Xvfb。測試結果為 36 個 raw TALK 頁面、18 組雙頁回應全部
通過；每頁實際為 1 頁，原始硬行為 2–3 行，沒有超過五行或 352 px modal 內容寬度。

## 分支索引

| 事件 | 前置 prompt／選單 | 結果分支與 raw TALK 索引 | 結果 |
|---|---|---|---|
| 3 停戰 | prompt `#360–#362`；選單 `#363` | 無條件 `#364 → #43`；指定金額 `#365 → #44`；拒絕 `#366 → #45`；超額 `#365 → #45` | PASS |
| 2 協力 | prompt `#373–#375`；選單 `#376` | 無條件 `#377 → #47`；指定金額 `#378 → #48`；拒絕 `#379 → #49`；超額 `#378 → #49` | PASS |
| 4 內政官撥款 | prompt `#278`；選單 `#283` | 全額 `#284 → #288`；等額 `#284 → #293`；低額 `#285 → #293`；零額 `#286 → #293`；超額 `#287 → #293`；拒絕 `#286 → #298` | PASS |
| 5 外交官撥款 | prompt `#319`；選單 `#324` | 全額 `#325 → #329`；等額 `#325 → #334`；低額 `#326 → #334`；零額 `#327 → #334`；超額 `#328 → #334`；拒絕 `#327 → #339` | PASS |

事件 2／3 的超額路徑只驗證原版已證實的 response=2 TALK 選取；狀態層的信賴度
扣除與合作／停戰收尾由既有 state tests 另外驗證。事件 4／5 的「超額」不套用
外交事件的拒絕規則，維持 `sub_139E8` 的獨立分支。

## DOS/V modal 代表幀

- 事件 2 合作結果：[`event2-47.png`](../images/event2-47.png)
- 事件 3 停戰金額結果：[`event3-44.png`](../images/event3-44.png)
- 事件 4 內政官全額結果：[`event4-284.png`](../images/event4-284.png)
- 事件 5 外交官拒絕結果：[`event5-339.png`](../images/event5-339.png)
- 事件 3 的 composite 選擇／數值視窗：[`wlgame-event3-choice.png`](../images/wlgame-event3-choice.png)、
  [`wlgame-event3-amount.png`](../images/wlgame-event3-amount.png)

代表幀確認文字框位置、DOS/V 外框、繁中字寬、硬換行、數值 marker 與
`Enter／Space` 分頁提示。此項完成的是所有已接入分支的 TALK／排版抽樣，
不宣稱原版自然長程的逐像素對拍。
（⚠ **密碼頁不擋**——四格留白按「確定」就進開場，
[`18`](18-dosv-password-verification.md)；長程對拍是還沒做。）
