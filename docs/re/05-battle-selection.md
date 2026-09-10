# 05 — 政略與戰術的接縫：戰場怎麼被選出來

**狀態：政略↔戰術的接縫已解；戰場記錄的部分欄位未解。**

- 日期：2026-08-07
- 輸入：`workplace/orig/dosv/KI.EXE`　SHA-256 `fffeba98…3868`（`.i64` 由它產生，當時未另記資料庫雜湊）
- 推論等級：**兩條路徑的存在與參數 confirmed**；地形對映表的內容也已解出，
  見 [`../mechanics/30-combat.md`](../mechanics/30-combat.md) §2

## 0. 為什麼追這一條

`BATTLE.MAP` 的載入器用 `ds:0D34h` 當戰場編號（`docs/formats/07` §2），
但誰寫這個變數沒查過。這條線是**大地圖（政略）與戰場（戰術）的接縫**
——remake 要能從政略進戰鬥，非解不可。

`byte_10D34`（＝ `ds:0D34h`）有兩個寫入點，對應**兩種完全不同的戰鬥**。

## 1. 攻城戰：據點編號就是戰場編號

`sub_14ADE`（由 `sub_12880` 呼叫）：

```asm
mov     cs:word_10D32, di
mov     ax, di
sub     ax, 840h            ; 減基址
shl     ax, 1
shl     ax, 1
shl     ax, 1               ; ×8
mov     cs:byte_10D34, ah   ; 取高位元組 ＝ (di − 0x840) / 32
mov     cs:byte_10D35, 0    ; ← 旗標 0：不旋轉
```

`ah` 是 `(di − 0x840) × 8` 的高位元組，等於 `(di − 0x840) / 32`。

| 結論 | 等級 |
|---|---|
| **據點記錄是 32 byte 一筆**，基址在段內偏移 `0x840` | confirmed |
| **據點編號直接當戰場編號** —— 每個據點有一個固定的戰場 | confirmed |
| 這條路徑 `byte_10D35 = 0` → **不觸發 180 度旋轉**（`docs/formats/07` §2.2） | confirmed |

## 2. 野戰：戰場由大地圖的地形即時決定

`sub_14A7B` → `sub_14B63`（由 `sub_12831` 呼叫）。

`sub_14B63` 先算出軍團在大地圖上的位置，再讀**周圍五格的地形值**：

```asm
mov     dx, [di+1Ch]        ; 軍團記錄 +0x1C
mov     bx, [di+1Ah]        ; 軍團記錄 +0x1A  ← 大地圖上的偏移
sub     dx, cs:word_19872
add     dx, cs:word_10D44
sub     dx, 18h
mov     ds, dx              ; ds 指向大地圖資料

mov     al, [bx+17Fh]       ; (X−1, Y)  左
call    sub_14C4C
mov     cl, al
mov     al, [bx+181h]       ; (X+1, Y)  右
mov     ch, al
mov     al, [bx]            ; (X, Y−1)  上
mov     dh, al
mov     al, [bx+300h]       ; (X, Y+1)  下
mov     dl, al
mov     al, [bx+180h]       ; ★ (X, Y)  **中心**——分流看這一格
mov     bl, al
```

**`0x180` ＝ 384 ＝ 地圖一列的格數**（`docs/formats/05` §2），
而 ⭐ **`sub dx, 18h` 讓 `ds` 指向 `Y − 1` 那一列**（佔用圖一列 24 段
＝ 0x180 byte）——所以整組取樣**往上位移一列**，
形狀是**以軍團所在格為中心的十字**：`[bx+180h]` 才是本格。

⚠ 少算那一行的位移會把「中心」讀成「正下方」，而**兩種讀法在同一格上
會給出不同的地形類型**：同局面拍 8,889（軍團 50 對 34，守方站在
(304,155)）中心是圖塊 `0xCA` ＝ 類型 8（要擲骰）、下方是 `0xC3`
＝ 類型 9（不擲骰），而原版擲了
（[`../spec/196`](../spec/196-field-battle-rolls-the-battlefield.md) §2）。

然後：

```asm
and     bl, bl
jz      short loc_14BD0     ; 正下方是 0 → 另一條分支
cmp     bl, 8
jnb     short loc_14BD5     ; ≥ 8 → 另一條分支
mov     cl, bl
add     cl, 0CEh            ; 戰場編號 ＝ 0xCE + 地形類型
xor     ch, ch
```

| 結論 | 等級 |
|---|---|
| **野戰的戰場不是查固定表，是依大地圖地形產生的** | confirmed |
| 取樣的是軍團所在格與其下方四格 | confirmed |
| 地形類型值域 1–7，戰場編號 ＝ `0xCE` ＋ 類型（`0xCF`–`0xD5`） | confirmed |
| `sub_14C4C` 是「地圖圖塊值 → 地形類型」的對映 | **confirmed**，14 筆範圍查表已攤開在
  [`../mechanics/30-combat.md`](../mechanics/30-combat.md) §2 |

> **這解釋了 214 個戰場的組成**：一部分是據點專屬的固定戰場，
> 另一部分（至少 `0xCF`–`0xD5` 這 7 個）是野戰用的地形樣板。
> 精確的分界還沒查。

## 3. 軍團記錄的欄位（部分）

從 `sub_14B63` 讀到的：

| 偏移 | 內容 | 等級 |
|---|---|---|
| `+0x01` | 與 `byte_10CFF`（當前勢力？）比對的欄位 → **所屬勢力** | 強證據 |
| `+0x08` | 勢力相符時取用的值 | 未解 |
| `+0x1A` | **大地圖上的位置**（格偏移，可直接當地圖陣列索引） | confirmed |
| `+0x1C` | 位置對應的段基底 | confirmed |

`+0x1A` 是格偏移而不是 (x, y) —— 直接就是 `y × 384 + x`。
**remake 的軍團座標要照這個表示法存**，否則存檔寫不回去。

## 4. 下一輪

1. ~~`sub_14C4C` 的對映~~ —— 已解，14 筆 × 3 byte 的範圍查表與各類型的
   驗證方法在 [`../mechanics/30-combat.md`](../mechanics/30-combat.md) §2。
2. `loc_14BD0`／`loc_14BD5` 兩條分支（正下方為 0、或類型 ≥ 8）。
3. 據點記錄 32 byte 的欄位表 —— 基址 `0x840`，配合 `SINARIO.DAT` 與
   `SAVE.DAT` 的 diff 應該能快速定位。
