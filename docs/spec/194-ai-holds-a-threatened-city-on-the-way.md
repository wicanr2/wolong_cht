# 194 — AI 行軍中經過受威脅的據點就留守

**狀態：CONFORMED。** `sub_12662` 有**第二條到站路徑**：軍團站在據點上
（`+0x0E < 800h`）而且**不是玩家勢力**時先問 `sub_14300`，回 CF=1 就跳到
`loc_1266A`——設朝向 4、`sub_128F4`、`sub_14325` 分派，**不再往前走**。

規則本身早就解出來了（[`../re/86`](../re/86-march-turnback-at-peace.md) §2.1，
confirmed），**remake 一直沒接**——`docs/spec/170` 的「出處」列了
`sub_14300` 卻沒有實作欄，這正是 `CLAUDE.md` §8 說的
「**『remake 實作』欄留白就是缺口**」。

- 日期：2026-09-11
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_14300`（IDA 線性 00014300–00014324）、
  `sub_12662`（00012662，`0001268D` 的呼叫點）
- 推論等級：confirmed（[`../re/86`](../re/86-march-turnback-at-peace.md) §2.1
  逐條讀出；原版的分派軌跡與意圖欄位對上）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §48.2
- remake 實作：`internal/state/corps.go` 的 `tickCorps`、`holdThreatenedCity`
- 相關：[`132`](132-march-turnback-at-peace.md) §「折返之後會不會再出發」、
  [`170`](170-ai-target-takes-two-ticks.md)、[`179`](179-recalc-on-battle-and-arrival-test.md)

## 1. 三個條件

```asm
sub_14300:                        ; bx ＝ +0x0E ＝ 據點 × 8
        shl bx,1 / shl bx,1       ; bx ＝ 據點 × 32
        cmp byte ptr [bx+858h], 1 / ja .no   ; ① 這一格的軍團數 > 1 ⇒ 不留守
        cmp byte ptr [si+8], 4    / jz .no   ; ② 已經是靜止（朝向 4）⇒ 不留守
        cmp byte ptr [bx+840h], 80h / jb .no ; ③ 據點沒受威脅（bit 7）⇒ 不留守
        shl bx,1 / shl bx,1 / shl bx,1
        mov [si+20h], bh          ; ★ **意圖 ← 腳下這座據點**
        stc / retn
.no:    clc / retn
```

呼叫端只在**非玩家**時問（`cmp al, cs:byte_10CFF / jz` 跳過玩家），
而且只在 `+0x0E < 800h`（站在據點上，不是走在路段上）時。

## 2. 對拍上長什麼樣

同局面拍 7,091：軍團 88（勢力 0，AI）在往據點 56 的路上打下據點 74
（原主勢力 13）。那座城剛換手、只有它一支、而且受威脅——三個條件都成立。

| | 之後 |
|---|---|
| 原版 | **停下來守**：`+0x20` ← 74、朝向 4、當拍分派。`WOLONG_DOSGOLEM_WATCH=1439D,143AF,1440F,14483` 的軌跡是 `6421:1439D(意圖56) 7125:14483(意圖74) 7725 7749 7773 7797`——7,125 起意圖就是 74 |
| remake | 照舊往 56 走，7,747 才到；到了才等士氣、挑目標，於是拿到的待辦與原版不同（勢力 0 的 `+0x17` 失守據點被 remake 取走、原版留著）|

到 5/25（拍 7,935）差 22 個 byte／17 座據點。

⚠ **這個差異不取亂數**，所以逐拍取數比對對它結構上是盲的——
7,935 拍裡只有兩拍不一致，而且都是下游效應。

## 3. remake 的修法

`tickCorps` 的「輪到移動」分支，在 `atTargetNode` 之後、走一步之前：

```go
} else if w.holdThreatenedCity(i) {
	// sub_14300 回 STC ⇒ 走 loc_1266A：朝向 4 ＋ 分派，不動佔用圖
	c.Heading = headingIdle
	w.arriveCorps(i, rng)
} else {
	…原本的移動…
}
```

⚠ **不要動佔用圖**：原版那一條在 `sub_14325` 之後直接 `retn`，
不跑結尾的 `lds di,[si+1Ah] / inc byte ptr [di]`。

⚠ **「站在據點上」的判準是座標，不是 `Node`。** 原版問 `+0x0E < 800h`，
而 remake 的 `Node` 在行軍中留著出發那一站——只看它會把走在半路上的
軍團讀成「站在城裡」。第一版就是這樣寫的，症狀是
`TestStageEncounterArmsDuelBeforeStepping` 掛掉（fixture 把敵方軍團放到
野外座標而 `Node` 沒動，於是它「留守」而不走過來，戰鬥開不出來）。
判準要是 `LinkAddr == 0 && onCity(i)`。

## 4. 驗證

單獨接這一條時 5/25 從 22 個 byte 變成 26——**它自己不夠**。
配上 [`195`](195-ai-stage1-asks-about-the-intent.md)（Stage 1 問的是意圖）
之後，`tools/parity_ck.sh 196/5/25`（拍 7,935）**六張表全 0**。

⭐ **對照實驗**：留著 `195` 而把這一條關掉，5/25 是 **180 個 byte／118 座**
——兩條都必要（`CLAUDE.md` §7 第 18 條：結果對不代表理由對）。

<!-- 缺口：無 -->
