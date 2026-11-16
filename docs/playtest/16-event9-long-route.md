# 事件 9 長程通知流程

**狀態：27 小時 bounded queue、玩家／非玩家／在野通知條件已通過；2026-11-16補本人199h組請罪；
完整自然劇本依要求略過，未宣稱原版長程逐像素 parity。**

- 日期：2026-08-11

事件 9 的原始事件字低 byte 是 `0x09`，高 byte 是 General index。正常每時
queue consumer 取出後，state 執行 `sub_150D7` 對應的釋放流程；呈現層只有在
釋放後 General 的 Faction 等於玩家勢力時，才取 `TALK.DAT #37` 建立通知。
`199h` 是組編號，須先展開為430–437。舊文直接查原始409空槽而稱no-op，已由下方原版呼叫證據勘誤。

## 27 小時自然 queue fixture

可重跑入口：

```text
go test -p=1 -vet=off ./internal/state \
  -run 'TestEvent9LongNaturalRoute$' -count=1 -v
```

測試在正常 `World.Tick`、9 子刻／時與事件 queue `eventDelay=7` 起點下，預先
放入三筆 raw event 9。結果固定如下：

| 每時邊界 | raw 事件 | 釋放結果 | GUI 通知 |
|---:|---|---|---|
| 第 7 小時 | General #0 | 回到玩家勢力 | `#37` |
| 第 17 小時 | General #1 | 回到非玩家勢力 | 無 |
| 第 27 小時 | General #2 | 俘虜方已滅亡，回到在野 | 無 |

此測試同時確認 `Posted=false`、`Captor=0xFF`，以及存活俘虜方／已滅亡俘虜方
的 Faction 寫回規則；它是有界的長程自然路徑，不等同完整長時間遊戲測試。

## GUI 通知接縫

```text
DISPLAY=:99 go test -p=1 -vet=off ./cmd/wlgame \
  -run 'TestEvent9(LongNotificationRoute|ShortFixtureGate)$' -count=1 -v
```

`TestEvent9LongNotificationRoute` 以相同事件觀測值連續餵入 GUI，確認兩次玩家釋放
各依序追加 `#37` 及本人請罪，非玩家與在野結果不會插入空白訊息；既有短 fixture
則保留玩家／非玩家與 `#409` 邊界。

DOS/V modal 代表幀：[`event9-37.png`](../images/event9-37.png)。

本項完成的是已反組譯證實的事件 9 raw queue → state release → 玩家 TALK #37
通知流程與有界自然時序；未取得當時原版長程錄影，不宣稱同狀態逐像素證據。

## 2026-11-16：199h組編號與第二句勘誤

原版15107／15110明確先顯示37，再以武將肖像及raw說話型呼叫199h。
199h不是直接使用資料的409槽，而是依1075B展開430–437；此規則已在[re25](../re/25-message-variants-and-personnel.md)證實。
原版十個受控樣本及Go完整區塊／RNG全等。五種有文字變體3–7由武將本人說話，Go先前漏接，已修正。
raw0–2原版會開空框並等右鍵，Go維持既有空訊息省略，這項操作次數差異明列保留。
正常remake十條滑鼠51項／15圖、完整冷測／vet及野戰回歸通過，依[spec219本片](../spec/219-c-event-handlers.md)。
本頁舊27小時測試只證明Go隊列及玩家通知條件，不再被用來主張本人第二句是no-op。
