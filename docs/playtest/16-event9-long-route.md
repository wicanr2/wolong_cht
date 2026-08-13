# 事件 9 長程通知流程

**狀態：27 小時 bounded queue、玩家／非玩家／在野通知條件與 #409 no-op 已通過；
完整自然劇本依要求略過，未宣稱原版長程逐像素 parity。**

- 日期：2026-08-11

事件 9 的原始事件字低 byte 是 `0x09`，高 byte 是 General index。正常每時
queue consumer 取出後，state 執行 `sub_150D7` 對應的釋放流程；呈現層只有在
釋放後 General 的 Faction 等於玩家勢力時，才取 `TALK.DAT #37` 建立通知。
`#409` 是原版資料上的空槽，維持 no-op。

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

`TestEvent9LongNotificationRoute` 以相同事件觀測值連續餵入 GUI，確認玩家釋放
會依順序追加兩個 `#37` modal，非玩家與在野結果不會插入空白訊息；既有短 fixture
則保留玩家／非玩家與 `#409` 邊界。

DOS/V modal 代表幀：[`event9-37.png`](../images/event9-37.png)。

本項完成的是已反組譯證實的事件 9 raw queue → state release → 玩家 TALK #37
通知流程與有界自然時序；不把密碼保護下無法取得的原版長時間錄影，宣稱成
同狀態逐像素證據。
