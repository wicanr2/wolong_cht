# 167 — 一個遊戲日是 **23** 小時 ＝ 207 子刻，時從 1 走到 23

**狀態：READY**

- 日期：2026-09-09
- 出處：`KI.EXE`（松崗 DOS/V）`sub_11D8E`（IDA 線性 00011D8E）
- 實測：dosgolem，`workplace/dosgolem/root-liubei13`，
  [`../playtest/119`](../playtest/119-rng-pace-comparison.md) §17
- 相關：[`../re/06`](../re/06-game-clock.md)、
  [`../mechanics/15`](../mechanics/15-realtime.md)、`internal/rules/clock`

## 1. 結論

原版的「時」只走到 **23**，而且**跨日之後從 1 開始**——`0 時不存在`。
所以一個遊戲日是 **23 × 9 ＝ 207 個子刻**，不是 216。

remake 的 `clock.HoursPerDay = 24`（時 1–24、一天 216 子刻），
每天比原版多 9 個子刻，**誤差 4.3%**，而且是累積的。

## 2. 組語

```asm
sub_11D8E:
        cmp     byte ptr ds:0CF2h, 8       ; 子刻 < 8？
        jb      short loc_11DF4            ;   → 只 inc 子刻
        cmp     byte ptr ds:0CF3h, 17h     ; ⭐ 時 < 23？
        jb      short loc_11DE0            ;   → 進位到「時」
        mov     ax, ds:0CF0h               ; al ＝ 日、ah ＝ 該月天數
        cmp     al, ah
        jb      short loc_11DD7
        …換月／換年…
loc_11DD7:
        inc     byte ptr ds:0CF0h          ; 日 +1
        mov     byte ptr ds:0CF3h, 0       ; 時 ← 0
loc_11DE0:                                 ; ⭐ **落下來**
        inc     byte ptr ds:0CF3h          ;   時 ← 1
        mov     byte ptr ds:0CF2h, 0
```

兩件事一起決定了「23」：

1. `cmp …, 17h` ＋ `jb` 是**嚴格小於 23**，所以時只被 `inc` 到 23 為止；
   時 ＝ 23 那一拍走的是換日那條路。
2. 換日把時歸零之後**落到 `loc_11DE0` 又 `inc` 一次**，所以新的一天從
   **1 時**開始。`0 時`在整個遊戲裡從來不會出現。

⇒ 時的值域是 **1–23**（23 階），一天 ＝ 23 × 9 ＝ **207** 子刻。

## 3. 實測

同一份存檔、同一個起點（196/4/16 16 時 子刻 3），
攔 `sub_11CD0`／`sub_13EFD`／`sub_14194`／`sub_11D8E` 四支：

```
快照之後各支被呼叫次數: {'11CD0': 200, '13EFD': 200, '14194': 200, '11D8E': 200}
遊戲時鐘：196年4月16日 16時  →  196年4月17日 15時
```

四支都是 200 次（**每個子刻剛好一次，沒有漏攔**），而時鐘走了 **209** 個子刻。
多出來的 9 個正是那一天的換日：`23 時 → 1 時` 跳過了 0 時。

拿 23 小時的模型往前推 1,900 個子刻，落點是 **196/4/25 20 時 子刻 4**——
與原版記憶體讀回來的完全相同；24 小時的模型會差 9 小時。

| 對拍拍 | 原版 | remake（24 小時模型）| 差 |
|---:|---|---|---|
| 200 | 4/17 15 時 | 4/17 14 時 | 1 小時 |
| 700 | 4/20 2 時 | 4/19 22 時 | 4 小時 |
| 1200 | 4/22 11 時 | 4/22 5 時 | 6 小時 |
| 1900 | 4/25 20 時 | 4/25 11 時 | 9 小時 |

⭐ **差距永遠等於這段期間的換日次數。**

## 4. remake 要改什麼

`internal/rules/clock/clock.go`：

```go
// HoursPerDay 是「日」進位前的時數。原版 `cmp byte ptr ds:0CF3h, 17h` ＋ `jb`
// ⇒ 時只 inc 到 23；時 ＝ 23 那一拍換日，而換日把時歸零之後**落下來又 inc**，
// 所以新的一天從 1 時開始——**0 時不存在**，值域是 1–23。
HoursPerDay = 23
TicksPerDay = HoursPerDay * SubticksPerHour // 207
```

`Advance` 的 `case c.Hour < HoursPerDay` 不用改，換常數就對。

連帶要改的敘述（都寫著 216／24）：
`internal/rules/clock` 的套件註解、`internal/rules/governor` 的套件註解、
`internal/state/state.go` 的據點巡迴註解、
[`../re/06`](../re/06-game-clock.md)、[`../mechanics/15`](../mechanics/15-realtime.md)。

## 5. 影響面

「一天幾個 tick」是所有日／月節奏的分母：

- **據點巡迴**：192 個據點、207 拍一天 ⇒ 一天掃 1.08 圈（原本以為 1.125 圈）。
- **月結與收稅**的到達時刻整體提前。
- 季節漸變掛在「時 ＝ 1」，本來就只在換日之後成立——這一條不受影響。

## 6. 怎麼驗

1. `internal/rules/clock` 的測試加「23 時的下一拍是隔天 1 時」與
   「一天 207 拍」兩個案例。
2. 節拍對拍：從同一份快照跑 1,900 拍，remake 的時鐘要落在 4/25 20 時 子刻 4。
