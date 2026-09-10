# 191 — 城兵的將領是武將表第 127 筆，不是空的

**狀態：CONFORMED。** `sub_14F8A` 搭城兵的臨時軍團時寫
`mov byte ptr [bx+2], 7Fh`——那是**武將編號 127**，指向武將表最後一筆。
那筆記錄在劇本裡是武力 8、統率 8、適性 0，不是全 0。

remake 的 `combat.Garrison` 把 `Leader` 留成零值，於是
`leaderValue(0, 0)` 回 **0**，`factor` 跟著是 0，**守方戰力整個歸零**。

- 日期：2026-09-11
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_14F8A`（[`../re/09`](../re/09-combat.md) §7）；
  武將 127 的內容出自快照 `workplace/parity/orig-at-16h-v2.DAT`
  的 `0x42C0 + 127 × 32`
- 推論等級：confirmed（機器碼寫死 `7Fh`；比值由對拍分開兩種讀法）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §47.4
- remake 實作：`internal/rules/combat/combat.go` 的 `Garrison`、
  `internal/state/corps.go` 的 `fightGarrison`
- 相關：[`../re/09`](../re/09-combat.md) §7（臨時軍團的其餘欄位）

## 1. 武將 127 長什麼樣

```
00 ff a1 d0 a1 d0 a1 d0 a1 d0 a1 d0 a1 d0 00 00
00 08 08 08 00 00 00 00 00 ff 00 00 ff ff 03 00
```

| 欄位 | 值 |
|---|---|
| `+0x00` 旗標 | `00` ⇒ **不是活著的武將**（`>= 0x80` 才是）|
| `+0x02`–`+0x0D` 姓名／呼び名 | `A1D0` 重複 ＝ Big5 全形空白 |
| `+0x0E`–`+0x10` 三個適性 | 0 |
| `+0x11` 武力 | **8** |
| `+0x12` 統率 | **8** |
| `+0x13` 政治 | 8 |

⇒ 它是一筆**刻意填了 8 的佔位記錄**，專門給城兵用。
名字是空白、旗標 0，所以它不會出現在任何武將清單裡。

## 2. 為什麼零值會讓戰力歸零

```go
func leaderValue(l Leader, rng Rand) int {
	m, c := l.Martial, l.Command
	base := c
	if m >= c {
		if rng.Next()&3 != 0 { return m * 2 }   // 0 → 0
		base = m
	}
	return base - base>>2 + c                    // 0 → 0
}
```

`Power` 最後一步是 `base * factor >> 10`，`factor` 為 0 就整個是 0。
守方戰力 `pd = 0 + 8 = 8`，而攻方動輒兩百以上 ⇒
`ratio = hi × 8 ÷ lo` 直接撞上限 **100**。

⚠ **亂數還是有取**（`rng.Next()&3`），所以逐拍取數比對對這個 bug
**結構上是盲的**——它一路綠到 5/20。

## 3. 對拍上長什麼樣

同局面 196/5/17–18（remake 拍 6,387）軍團 82 攻據點 74 的城兵：

| | 比值 | 城損 ＝ `(0x3F − 比值) >> 2` |
|---|---:|---:|
| 原版（`WOLONG_DOSGOLEM_WATCH=151B3` 的 `al`）| **14** | 12 |
| remake（`rng_pace.go -battle-log`）| **100** | 54 |

城損差 42，而城損會同量扣掉上昇值、防災值與城兵三欄——
症狀就是「據點 74 三個欄位各差 42」。

⚠ 原版側的 `si`／`di` 是**執行期**的軍團記錄位址，基址 `0x2240`，
不是快照裡的 `0x22C0`（差兩筆）。照快照的基址換算會得到
「軍團 80 打軍團 125」，而正確答案是「軍團 82 打 `0x4200` 的臨時記錄」。
**表的基址在執行期與存檔佈局裡不一定相同。**

## 4. remake 的修法

- `World` 載入時多讀一筆武將 `127`（`generalBase + 127 × generalSize`）
  存成 `garrisonLeader`，**不改 `Generals` 的長度**（存檔寫回照舊只寫
  0–126，未解的 byte 一個都不碰）。
- `combat.Garrison` 多收一個 `Leader` 參數，由 `fightGarrison` 傳進去。
- 快照裡那一筆全 0 時（新局面、劇本沒填）退回武力 8／統率 8，
  並在註解裡標明那是觀測值不是常數。

## 5. 驗證

`tools/parity_ck.sh 196/5/18 196/5/20` 的據點表歸零。

<!-- 缺口：無 -->
