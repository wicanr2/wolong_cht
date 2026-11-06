# 121 — 碼頭與棧道的野戰用錯了戰場：`SelectWater` 從來沒有被接上

**狀態：CONFORMED。** 地形類型 8／9 的戰場編號算成 `TerrainBase + 類型`
＝ **214／215**，而戰場只有 0–213——呼叫端看到超出範圍就退回**合成戰場**。
`SelectWater`（正確的那一支）寫好也有單測，**只有它自己的測試在呼叫它**。

- 日期：2026-09-03
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_14C1A`
  （`0001C1A` 段內 `0x4C1A`）——[`../re/05`](../re/05-battle-selection.md)、
  [`../mechanics/30`](../mechanics/30-combat.md) §2
- 推論等級：**confirmed（原始 bytes 交叉解碼）**
- 相關：[`115`](115-soldier-power.md)（同一條鏈的另一個接線錯誤）

## 1. 原版做什麼

```asm
sub_14C1A:                      ; bl ＝ 中心格的地形類型（≥ 8 才進來）
        cmp bl, 9 / jz .kind9
        call sub_1ECE0          ; ⭐ 亂數
        and al, 3
        mov cl, al / add cl, 0D1h   ; ⇒ 209–212（0xD1–0xD4）
        xor ch, ch / retn
.kind9: ...另一條（遠指標讀取，見spec273勘誤）...  ⇒ 213
```

| 地形類型 | 大地圖圖塊 | 戰場編號 |
|---|---|---|
| 8（碼頭）| `0xCA` | **209–212 隨機挑一張** |
| 9（橋／棧道）| `0xC0`–`0xC3` | **213** |

⭐ **戰場編號決定的不只是地圖**：`byte_1D34B` 由它分類，而類別決定
側欄標題（「陸上」／「海上」）、腳本段、以及主將的哪一個適性欄
（[`115`](115-soldier-power.md)）。209–213 全部 ≥ `0xD1` ⇒ **海上作戰**。

## 2. remake 錯在哪

```go
case b >= 8:
    // 隨機那一支要亂數，交給 SelectWater。
    return TerrainBase + b, false      // ⛔ 0xCE + 8 ＝ 214、+9 ＝ 215
```

註解寫著「交給 SelectWater」，但**沒有任何呼叫端接**——
`grep -rn "SelectWater"` 只有 `battlefield_test.go` 兩筆。
於是 `BuildField` 的 `n >= battle.NumFields` 那道閘成立，
碼頭與棧道的野戰開出來的是 `SyntheticField`：
**平坦的合成地形，不是原版那五張**。

⚠ 適性與標題**碰巧是對的**（214／215 也 ≥ `0xD1` ⇒ 海上），
所以症狀只出現在地形上——而地形對不對沒有任何測試在看。

## 3. remake 怎麼改

| 項目 | 位置 |
|---|---|
| 帶亂數的選擇 | ✅ `battlefield.SelectWith(dir, n, roll)`；`Select` 變成 `roll = 0` 的薄殼 |
| 只有類型 8 抽 | ✅ `battlefield.NeedsWaterRoll`——**多抽一次會讓亂數流錯位**，而戰術與戰略共用同一條 |
| 抽亂數 | 2026-11-04由規則層在接戰時抽一次；舊Provider抽骰接線由spec273取代 |
| 同一場只算一次 | Pending保存編號及翻轉；舊Provider按據點永久快取由spec273移除，連續遭遇重新選場 |

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 ✅ | `TestSelectAlwaysInRange`：**十種地形類型 × 五個方向 × 四個亂數值，算出來的編號都要落在 0–213**。這一支原本不存在，所以 214／215 沒有人擋 |
| 單元測試 ✅ | `TestNeedsWaterRollOnlyForKind8`：只有類型 8 抽亂數 |
| 既有 ✅ | `SelectWater` 自己的單測（209+roll&3、類型 9 固定 213）本來就在，**只是沒有人呼叫它** |
| 對原版 ✅ | 2026-11-04正常孫策開局CA野戰212，96等待點、選場後完整RNG及海戰適性通過；見playtest129。完整分支矩陣仍待驗 |

## 5. ⭐ 教訓

**寫好、測好、沒有人呼叫**——這是本專案第三次踩同一個形狀：
`state.TacticalSetup.Tile` 從頭到尾沒有人設（[`115`](115-soldier-power.md) §3.1）、
`SelectWater` 只有測試在呼叫、`replanInterval` 取代了原版根本不同的機制
（[`120`](120-pathfind-request-queue.md)）。

三次的共同點是**單元測試驗的是「這支函式算得對不對」，
不是「有沒有人用它」**。便宜的解法是加一條**值域斷言**：
`TestSelectAlwaysInRange` 不看語意，只問「算出來的東西合不合法」——
它會在任何一條分支算出 214 的當下就紅。

## 6. 未解

| 項目 | 現況 |
|---|---|
| 類型 9 的另一條分支 | 213編號已證實，remake維持不翻轉近似。`14C33`的LDS切換DS，`14C36`讀攻方佔用圖byte，`14C41`仍讀該DS的`[si+1]`；完整語意未知，不能解釋為鄰格或攻方勢力 |
| 其餘海戰分支與結算 | playtest129只補類型8的戰場212，未代證其餘編號、類型9或完整自動／戰術結算 |

2026-11-04 勘誤來源：[spec273](273-encounter-battlefield-context.md) 的固定IDA匯出與原始bytes。
舊「鄰格CA」說法未追蹤DS，不能作規則依據。局部terrain矩陣60案通過只證明原版函式與既有固定C，
Go接線另以非平原／水域測試及正常道路野戰回歸驗證，不冒稱正常海戰完成。

同日後續：[playtest129](../playtest/129-normal-water-field-battle.md) 已補正常CA野戰。
原版275,647,559指令自然進14A7B，選場前raw RNG僅取一次即到11B5A狀態及212編號，
Go在Spawn前完整RNG一致；錯座標版本選198並被同一比較器拒絕。上述局部terrain矩陣與新正常場景分開記錄。
