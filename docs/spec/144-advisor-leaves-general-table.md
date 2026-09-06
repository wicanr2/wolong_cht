# 144 — 選了軍師，那個人就從武將表消失

**狀態：CONFORMED。** 新遊戲定案時原版把**選中的軍師整筆從武將表移除**
（記錄 `+0x00` 寫 0，存在旗標一起沒了），並把勢力的武將數 `+0x18` 減一。
所以他不出現在任何清單裡——編成、人事任命、武將一覽都一樣。
remake 沒有這一步，改用「候選過濾時比對 `Advisor` 編號」擋，
而那個補丁只貼在編成，人事任命漏了。

- 日期：2026-09-06
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_11AC3`（`00011AC3`）
  的 `loc_11AF8`（`00011AF8`）
- 推論等級：**confirmed**（逐行讀完 ＋ dosgolem 讀原版記憶體）
- 驗證：[`../playtest/85`](../playtest/85-general-duty-field.md) §5
- 相關：[`76`](76-lord-not-in-formation.md) §2（原本掛在這裡的缺口）、
  [`../re/73`](../re/73-new-game-faction-list.md) §1、
  [`143`](143-general-duty-field.md)（同一張清單的另一個過濾條件）

## 1. 原版做什麼

```asm
loc_11AF8:                       ; 君主卡按下「確定」之後
    mov ds, cs:word_10D52
    mov ah, [bx+2]               ; bx ＝ 勢力記錄，+0x02 ＝ 軍師編號
    cmp ah, 7Fh
    jz  short loc_11B08
    dec byte ptr [bx+18h]        ; ① 勢力武將數 −1（自定軍師時跳過）
loc_11B08:
    xor al, al / shr ax,1 ×3     ; ax = 軍師編號 × 32
    mov si, ax
    mov byte ptr [si+4240h], 0   ; ② ★ 武將記錄 +0x00 ＝ 0（整筆抹掉）
    mov cs:word_10CFD, bx        ; ③ 玩家勢力（記錄位址）
    shl bx,1 ×2 / mov cs:byte_10CFF, bh   ;    玩家勢力編號
```

⭐ **②寫的是 `+0x00` 不是 `+0x17`。** 存在旗標（bit 7）在那個 byte 裡，
所以清成 0 等於**這一筆武將不存在了**——不是「他有職務」，是「他不在表上」。

**這就是軍師為什麼不出現在任何清單裡。** 建清單的 callback（線性 `0x176A0`
等）第一個條件都是 `cmp byte ptr [si], 80h / jb 跳過`，一次擋掉。
不需要另外的身分欄位，也不需要每張清單各寫一次排除。

⚠ 自定軍師（`+0x02 ＝ 0x7F`）時 ①跳過而②照寫，`si` ＝ `0x7F × 32` ＝ `0xFE0`，
落在武將表尾端之後一個 byte。**那是原版的越界寫，不要照抄。**

## 2. 現場證據

`root-saveb` 的存檔（曹操，軍師欄 62）在 dosgolem 上讀出來：

```
勢力  0 君 16 師 62 都 82 …… 將 11 ……
武將 62  2754:4A00 = [00] 79 AFFB D17B A140 …    ← +0x00 ＝ 0，不在表上
武將 63  2754:4A20 = [80] 45 AFFB A7F1 A140 …    ← 荀彧，正常
```

- `generals:0` 只印得出 11 人，與勢力的 `將 11` 一致——**武將數已經減過**。
- 同一局面開「編成」與「人事 → 內政官任命」，兩張清單**逐列相同**
  （夏侯淵、許褚、荀彧、徐晃、曹洪、曹純、曹仁、程昱、典韋），
  少的只有君主（曹操）與軍團長（夏侯惇）。**荀彧在列**。

## 3. 這推翻了什麼

[`76`](76-lord-not-in-formation.md) §2 記著「軍師也不在候選裡……
是執行期有人把軍師的身分寫進 `+0x17`；**寫入者未讀**」。
那個假說不成立——`+0x17` 沒有人為軍師寫過任何值
（[`143`](143-general-duty-field.md) §2 的寫入端全表），
擋住他的是 `+0x00`。缺口關閉。

[`../re/73`](../re/73-new-game-faction-list.md) §1 把 ② 記成
「把選中的武將 `+0x17` 職務清 0」，位移抄錯一個欄位。

## 4. remake 實作

| 檔案 | 改了什麼 |
|---|---|
| `internal/state/advisor.go` | 新增 `TakeAdvisor(faction, who int)`：`Generals--`、`Generals[who].Alive = false`，並寫 `Factions[faction].Advisor = who` |
| `cmd/wlgame/launcher.go` | `launcherStartNewGame` 定案時呼叫它（自定軍師走既有的 `SetCustomAdvisor`，不動武將表）|
| `cmd/wlgame/corps.go` | `formCandidates` 拿掉 `i != advisor` 的補丁——真正的機制已經涵蓋 |

⚠ **`LoadScenario` 不做這件事。** 劇本檔裡軍師還在表上（他要能被選），
移除是「選定之後」才發生的，所以下手點在定案那一步，不在載入。

## 5. 驗證

| 怎麼驗 | 結果 |
|---|---|
| `TestTakeAdvisorRemovesHimFromTheGeneralTable`：`Generals` 減一、`Alive == false`、寫回的 `r[0x00]` 除了未解的 bit 0 之外全 0 | 通過 |
| `TestAdvisorInNeitherCandidateList`：`formCandidates` 與 `freeGenerals` **兩張**都不含軍師 | 通過 |
| 勢力表對拍：原版 `root-saveb`（曹操局）勢力 0 將 11／勢力 1 將 12；remake `-player 0` 得 11／12、`-player 1` 得 12／11 | **逐格相同**，而且**減一只發生在玩家自己的勢力** |

## 6. 未解

| 項目 | 現況 |
|---|---|
| 自定軍師時原版那個越界寫 | 位置在武將表尾端後一個 byte，寫進去的是什麼欄位沒查；remake 不照抄 |
| 軍師退場時（如果有）會不會放回表上 | 沒找到反向的寫入端 |
