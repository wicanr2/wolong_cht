# 127 — 被俘的主公型武將當場變成臣下型

**狀態：CONFORMED。** `sub_129C3` 在寫完舊主之後，只要武將旗標的 bit 6
（主公型）成立就**清掉它並把說話類型 `+0x1E` 加 3**——0／1／2 正好搬到
3／4／5。remake 的 `Sovereign` 與 `TalkVariant` 兩個欄位本來就在，
**但先前沒有任何地方改它們**——存檔讀進來、原樣寫回去，中間一次都沒動，
所以被俘的君主之後講話還是用主公型的變體。已接上（§3）。

- 日期：2026-09-03
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_129C3`
  （`00012A12`–`00012A1D`），IDA `.i64` SHA-256 `65736f11…40572`（739 支函式）
- 推論等級：**confirmed（靜態）**
- 相關：[`../re/77`](../re/77-general-affinity-and-flags.md) §3（bit 6 的分布）、
  [`123`](123-captive-talk-messages.md)（被擒的訊息）、
  [`../re/25`](../re/25-message-variants-and-personnel.md) §4.2（說話類型 0–2 主公／3–7 臣下）

## 1. ⚠ 訂正：這一段**不以「舊主已滅」為條件**

`re/77` §3 與 [`../formats/08`](../formats/08-sinario-save.md) §「`+0`」
都把它寫成「被俘**而舊主勢力已滅時**清掉 bit 6」。逐條看控制流不是這樣：

```asm
000129F8  xchg al, [bx+1Ch]        ; 現主 ← 勝方，al ← 舊主
000129FB  mov  [bx+1Dh], al        ; +0x1D ← 舊主
000129FE  mov  ah, al / xor al, al / shr ax,1 / shr ax,1   ; ax ＝ 舊主 × 64
00012A06  mov  di, ax
00012A08  cmp  byte ptr [di], 80h  ; 舊主勢力還在？
00012A0B  jnb  short loc_12A12     ; ★ 還在 ⇒ 直接跳到 bit 6 這一段
00012A0D  test byte ptr [bx], 10h  ; （舊主已滅）bit 4 ＝ 不事二主
00012A10  jnz  short loc_12A57     ; ⇒ 自刎，整筆歸零
loc_12A12:
00012A12  test byte ptr [bx], 40h  ; ★ bit 6 ＝ 主公型
00012A15  jz   short loc_12A1E
00012A17  and  byte ptr [bx], 0BFh ; 清 bit 6
00012A1A  add  byte ptr [bx+1Eh], 3; 說話類型 +3
loc_12A1E:                          ; ← 訊息那一段（`123` §1.2）
```

`jnb loc_12A12`（舊主還在）與「舊主已滅但沒有 bit 4」**兩條路都落到
`loc_12A12`**。只有走自刎那一條才跳過去，而自刎的下一步是整筆歸零，
旗標怎樣已經無所謂。

**所以條件只有一個：被俘且沒自刎。** 舊主在不在只決定要不要檢查 bit 4。

## 2. 為什麼原版要這樣做

說話類型決定訊息的變體（`../re/25` §4.2）。一個當過君主的人被俘之後
不再是君主，**再用主公型的口氣講話就不對了**——而 43 名 bit 6 的武將
說話類型全部落在 0–2，`+3` 正好把整組平移到臣下型的 3–5。
這是一次不可逆的搬移：`sub_150D7` 釋放俘虜時**沒有把它加回去**。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 欄位 | `internal/state/state.go`：`General.Sovereign`（bit 6）、`General.TalkVariant`（`+0x1E`）——兩個都已經有，先前**只有存檔讀寫在用** |
| 被俘 | `internal/state/outcome.go` 的 `disperseFaction`：寫 `g.Captor, g.Faction` 之後 |
| 共用 | 抽一支 `demoteCapturedSovereign(g *General)`，兩個被俘入口都呼叫它（`CLAUDE.md` §7 第 6 條：一條規則只留一份實作）|

⚠ **原版有兩個入口**（`sub_129C3` 的呼叫者）：`sub_1291A`（戰場上武將被擒）
與 `sub_14FCE`（守城方潰散）。remake 目前只有 `disperseFaction` 這一條走到
「成為勝方的俘虜」，戰場那一條記在 [`123`](123-captive-talk-messages.md) §4。

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestCapturedSovereignBecomesRetainer`：bit 6 的武將被俘後 `Sovereign == false` 且 `TalkVariant` 加 3；沒有 bit 6 的不動 |
| 單元測試 | `TestDemoteCapturedSovereignIsOneWay`：呼叫兩次只加一次（先測 bit 再清）|
| 單元測試 | `TestReleasedCaptiveKeepsRetainerVariant`：釋放不還原（原版 `sub_150D7` 沒有加回去）|
| 突變測試 ✅ | 把 `+= 3` 改成 `+= 0`，前兩支當場紅——**確認測試真的在比**，不是綠在空氣上 |
| 存檔 round-trip | `+0x00` 的 bit 6 與 `+0x1E` 都在既有的 byte-for-byte 改寫路徑上 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| 被俘兩次會不會加兩次 | **不會，而且不需要另外防護**——`and [bx], 0BFh` 已經把 bit 6 清掉了，第二次 `test [bx], 40h` 不成立。所以 `+3` 至多發生一次，值域停在 3–5，不會溢出 `+0x1E` 的 0–7。remake 照抄同一個結構（先測 bit 再清）就自然有同樣的性質 |
| 劇本作者能不能給非君主 bit 6 | 四個劇本的 43 筆全是現任君主（`../re/77` §3），但那是**資料上的巧合還是規則**沒有讀出來。若有一筆說話類型 3–7 又帶 bit 6，`+3` 會把它推到 6–10 |
