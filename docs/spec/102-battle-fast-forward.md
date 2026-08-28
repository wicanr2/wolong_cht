# 102 — 戰場的 `▶▶`：快轉＝跳過戰場重畫

**狀態：CONFORMED（2026-08-29 實作、單測、實機對照過行為；底紋見 §5）。**

- 日期：2026-08-29
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `loc_1A065`（自我修改碼，
  `tools/ida_range.py` 逐 byte 手解）、handler `0x1C234`／`0x1C260`
  （[`../re/60`](../re/60-tactical-sidebar.md) §9）、`sub_1DC9D`（鏡頭走訪）。
  實機：[`../playtest/53`](../playtest/53-battle-fast-forward.md)。
- 推論等級：**confirmed**（機器碼 ＋ 實機兩條獨立證據）

## 1. 原版做什麼

`loc_1A065` 是每一幀都跑的一段（`sub_1A04B` 呼叫）：

```
0001A065  cmp  byte ptr ds:0D348h, 0     ; 鏡頭 dirty flag（byte_1D348）
0001A06A  74 10 / EB 10                  ; ★ jz／jmp short → 1A07C
0001A06C  mov  dx, [0D328h]              ; 鏡頭 X（word_1D328）
0001A070  mov  bx, [0D32Ah]              ; 鏡頭 Y（word_1D32A）
0001A074  call sub_1DC9D                 ; 從鏡頭重算整個顯示格（重畫戰場）
0001A077  mov  byte ptr ds:0D348h, 0     ; 清 dirty
0001A07C  call sub_1A12A …               ; 之後是兵的移動、小地圖等
```

IDA 把 `1A06A` 那一個 byte 當資料，後面整段就解錯了（`adc [bp+di+2816h], cl`
那幾行是假的）；拿原始 bytes 手解才對得起來。

⭐ **`byte_1D348` 的讀取端就在這裡**——[`57`](57-tactical-projection.md) §6 掛著
「三處寫、零處讀」，零是因為讀取端在自我修改碼裡，xref 掃不到。

`▶▶`（熱區 `0x0F`）的 handler 是一個切換：

| | 做什麼 |
|---|---|
| 開 | opcode 改成 `EB`（**永遠跳過**重算）、`sub_1DC9D(0x80, 0x80)` 把視點設到 (128,128)、在 (497,377)-(622,390) 描一圈色 **12** |
| 關 | opcode 還原 `74`、`byte_1D348 = 1`（下一幀從鏡頭重算）、同一圈描色 **10** |

(128,128) 在 64×64 的戰場外：`sub_1DC9D` 從 `Y×64＋X ＝ 0x2080` 起走，讀到的是
戰場資料之後的記憶體，畫出來是一片均勻的底紋（§5）。**效果是戰場不再重畫**，
兵照樣動、小地圖與側欄照樣更新——也就是**快轉**：省掉每幀最貴的那一步。

## 2. remake

| 項目 | 位置 |
|---|---|
| 切換 | `cmd/wlgame/battle.go`：`toggleBattleFastForward`（點 `SideFooter` 熱區）|
| 開著時 | 戰場區不呼叫 `drawBattleIso`，鋪視窗底的龍紋（`chrome.FillMenu`）；規則層照跑 |
| 按下框 | 按過之後才描：色 12（開）／色 10（關），126×14 於 (497,377) |
| 驗收旗標 | `-battle-ff`（配 `-open-battle`／`-open-siege`）|
| 差異 | 底紋（§5）；原版的視點重設不需要（關掉時鏡頭沒動，重算回原處）|

## 3. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestBattleFastForwardToggle`（`cmd/wlgame`）：熱區在 `SideFooter`、開→關→開 |
| 實機 | `playtest/53`：按一下戰場整片變底紋、側欄與小地圖照動；再按一下戰場回來 |
| 截圖 | `tools/parity_shot.sh … -open-battle -battle-ff`：戰場區鋪龍紋、`▶▶` 框描色 12 |

## 4. 為什麼先前解不開

`re/60` §12 記「`loc_1A065` 未逐行讀」——IDA 對自我修改碼的那一格
（`db 74h`）之後全部解錯，看起來像一段垃圾。`CLAUDE.md` §7 第 28 條的判準
（有沒有 `mov cs:[…]` 寫進程式碼段）早就指出它是自我修改碼，缺的只是
拿原始 bytes 交叉解碼那一步。

## 5. 未解

| 項目 | 現況 |
|---|---|
| 快轉時戰場區的底紋 | 原版是**藍底綠線的菱形格**（`playtest/53` 的裁切），不是龍紋——那是 `sub_1DC9D` 從 `es:0x2080` 讀到的「戰場之後的記憶體」畫成圖塊的樣子。要對上得知道那一段是什麼（候選：均勻填值的表）；可用 `tools/dosboxx_bridge.sh` 在戰場中讀 `word_1E15E:0x2080` 起 64 B 定案。remake 先鋪龍紋，**標為差異** |
