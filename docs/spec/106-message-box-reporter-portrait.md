# 106 — 訊息框那張臉是固定的通報者，不是說話者

**狀態：CONFORMED（2026-08-29 機器碼全量稽核 ＋ 實機對照，已實作與單測）。**

- 日期：2026-08-29
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_18810` 與它的 **60 個呼叫點**
  （全量掃 `KI.EXE.asm`）、`sub_1075B`（`al` → `sub_10AD9` 畫肖像、`cx` → TALK 索引）。
  實機 [`../playtest/55`](../playtest/55-encounter-menu-parity.md) 的 `e4.png`。
- 推論等級：**confirmed**

## 1. 原版做什麼

`sub_18810`（訊息框）把呼叫端的 `al` 與 `cx` 原封不動傳給 `sub_1075B`：

```asm
00018810  push bx / push dx / push ax / push cx
00018814  dx=0Ah bx=8 al=1Eh cx=510h → sub_1895D   ; 畫框
00018822  pop cx / pop ax                          ; ★ 還原呼叫端的 al 與 cx
00018827  dx=0Ah bx=0Ah → sub_1075B                ; al ＝ 肖像頁、cx ＝ TALK 索引
```

**60 個呼叫點的 `al`**：

| `al` | 幾處 | 是什麼 |
|---|---:|---|
| `93h` | **40** | 固定的**通報者**肖像 |
| `[bx+4241h]`／`[bx+1]`／`[di+1]` | 20 | **說話者**的肖像（武將記錄 `+0x01`）|

用說話者肖像的那 20 處，`cx` 幾乎都是 ≥`0x196` 的變體組
（`0x197` 財政赤字、`0x198` 財政困境、`0x19A`、`0x1A4`…），例外是
`sub_13327`／`sub_13388` 的 **#58**「敵方的君主已不在了。」。

⭐ **`0x93` 是 KAOGRF 的第 147 頁**——一個戴帽的老者。KAOGRF 共 **150 頁**
（0x00–0x95），而命名視窗可選的自訂軍師肖像只到 `0x92`
（[`104`](104-advisor-naming-window.md) §1 的 `sub_1912D`），所以
**`0x93`–`0x95` 三頁是保留頁**，`0x93` 就是通報者那張臉。

## 2. remake 先前做錯了什麼

`noticePortraitPage` 只要 `notice.General` 有值就用那個人的肖像——
所以每一則通報都畫**被通報的人**，而不是通報者。實機截圖裡
「夏侯惇大人的兵馬，遇上呂布的兵馬了！！」旁邊是那位老者，
remake 畫的是夏侯惇。

## 3. 重複標記：一個標記出現兩次是兩個值

同一批機器碼還解掉另一件事：原版的 formatter 是一個**共用的堆疊游標**
（`sub_14EB9` 的 `mov di, sp`），**每個標記消耗下一個參數**。所以
#29「{1}大人的兵馬，遇上{1}的兵馬了！！」的兩個 `{1}` 是兩個不同的武將。

全庫 1,022 則裡有 9 則出現重複標記，其中只有 **#29（兩個 `{1}`）** 與
**#217（兩個 `{3}`，「和{3}大人合作，去打倒{3}如何？」）** 真的要不同值；
其餘 7 則重複的是 `{6}`（排版控制，代入空字串）。

## 4. remake 實作

| 項目 | 位置 |
|---|---|
| 依序取值 | `internal/assets/text`：`Table.LinesSeq(index, vars, seq)`；`Lines` 轉呼叫它 |
| 序列來源 | `state.TalkNotice.SeqGenerals`／`SeqFactions` ＋ `World.TalkNoticeSeq` |
| 肖像 | `cmd/wlgame/messages.go`：`reporterPortraitPage = 0x93`；≥`0x196` 的變體組、#58 與 `SpeakerPortrait` 才用說話者 |
| 差異 | 無 |

## 5. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestLinesSeqGivesEachRepeatedMarkerItsOwnValue`（`internal/assets/text`）、`TestNoticePortraitUsesTheReporterFaceExceptForSpeakers`（`cmd/wlgame`）|
| 對原版 | 遭遇訊息與 `playtest/55` 的 `e4.png` **文字一字不差、肖像同一張**（[`105`](105-encounter-goes-straight-to-battle.md) §4）|

## 6. 未解

| 項目 | 現況 |
|---|---|
| `0x94`／`0x95` 兩頁保留肖像的用途 | `0x94` 是一張紅臉武將、`0x95` 是空白。沒找到傳這兩個值的呼叫點 |
| #217 的兩個 `{3}` | 機制已通（`SeqFactions`），但**發那一則的呼叫端還沒讀**，所以第二個勢力是誰未定 |
| ~~推廣片的靜態圖~~ | ⚠ `docs/images/promo-talk.png` 這類訊息框截圖是**改動前**拍的（肖像會變）。重拍會動到別份文件記的 fixture 雜湊，下次重錄推廣片時一起處理 |
