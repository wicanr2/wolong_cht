# 06 — `MMAP.MAP` 的 RLE 壓縮

**狀態：READY。** 解出來的地圖已用實機畫面與縮小地圖雙重驗證。

- 日期：2026-08-07
- 出處：`KI.EXE` 的 `sub_1F5E7`
- 工具：`tools/rle.py`、`internal/assets/rle`

## 1. 為什麼會漏掉這一支

`MMAP.MAP` 走的載入器**與其他資料檔不同**：

| 包裝 | 實作 | 誰用 |
|---|---|---|
| `sub_1E378` | `sub_1F4A2` | 一般讀檔（`TALK.DAT`、`SOUND.DAT`…） |
| `sub_1E38C` | `sub_1F4DF` | 帶位移讀檔（`KAOGRF` 那一族） |
| **`sub_1E364`** | **`sub_1F5E7`** | **只有 `MMAP.MAP`** |

三支包裝長得幾乎一樣（都是 push／call／pop ＋ 失敗重試），
**差別只在呼叫的內層函式**。掃過去很容易當成同一件事。

## 2. 演算法：用「連續兩個相同的 byte」當 run 的觸發

沒有逃脫字元。

```
逐 byte 複製；
一旦輸出的 byte 與前一個相同，下一個輸入 byte 是「再重複幾次」；
次數 0 表示那兩個相同的 byte 就只是字面值，回到逐 byte 模式。
```

所以一段 run 的總長是 **2 + count**。原始碼（`sub_1F5E7`）：

```asm
loc_1F602:
        call    sub_1F69E       ; 取一個輸入 byte → al（ah=FFh 表示 EOF）
        cmp     ah, 0FFh
        jz      short loc_1F638
        mov     [si], al        ; 輸出
loc_1F60C:
        mov     dl, al          ; dl ＝ 前一個
        call    sub_1F691       ; 輸出指標前進
        call    sub_1F69E
        cmp     ah, 0FFh
        jz      short loc_1F638
        mov     [si], al        ; 輸出
        cmp     al, dl
        jnz     short loc_1F60C ; 不同 → 繼續逐 byte
        ; 相同 → 進入 run
        call    sub_1F691
        call    sub_1F69E       ; 取重複次數
        cmp     ah, 0FFh
        jz      short loc_1F638
        cmp     al, 0
        jz      short loc_1F602 ; 次數 0 → 只是兩個字面值
loc_1F62E:
        mov     [si], dl
        call    sub_1F691
        dec     ax
        jnz     short loc_1F62E
        jmp     short loc_1F602
```

## 3. 尾巴會多幾個 byte，這是正常的

`MMAP.MAP` 80,716 B 解出來是 **98,308 B**，而地圖是 384 × 256 ＝ **98,304 格**。

**多的 4 個 byte 不是 bug。** 原版的解壓器**不知道目標長度**，
它一路解到檔尾；呼叫端只用前 98,304 B。
`internal/assets/world` 把尾巴留在 `Map.Extra`，不丟掉
——寫回時未解區域一個 byte 都不能動（`CLAUDE.md` §10）。

## 4. 驗收：四件事同時對才畫得出這張圖

![大地圖總覽](../images/mmap-overview.png)

（384 × 256，1 px ＝ 1 格，每格取該圖塊的平均色。）

黃河與長江的走向、渤海與黃海的海岸線、山脈與平原的分佈、
黃色的道路網——**全部連貫**。要畫出這張圖，下面四件事必須同時正確：

1. RLE 演算法（錯一個 byte 之後整條就亂）
2. 地圖尺寸 384 × 256（錯了會斜掉）
3. `MMAP.MDL` 的 256 塊圖塊解碼（`docs/formats/05`）
4. `.BRG` 調色盤的通道順序（`docs/formats/02`）

**大輪廓也與 `ICONGRF` 段 2 的縮小地圖（192 × 128）一致**
——河網走向與海岸位置對得上。

> ⚠ 但**縮小地圖不是世界地圖的降採樣**，是另外畫的一張圖。
> 逐格對應不成立（量化結果見 `docs/playtest/03` §2），
> 所以它驗證的是大輪廓，不是逐格的地形分類。

用到 **238 種圖塊**（共 256 種），最常見的一種佔 11,835 格（12%），
分佈合理，不是解壓爆掉的一片同值。

## 5. 兩版

`dosv` 與 `pc98` 的 `MMAP.MAP` **內容不同但解出來的長度相同**
——同一張地圖的兩次壓縮。測試 `TestRLESameShapeBothVersions` 釘住這一點。

## 6. 這支解壓器還有誰在用

目前只在 `MMAP.MAP` 上確認。`BATTLE.MAP`（877,056 B）、`MMAP.MCH`、
`BATTLE.MDL` 走哪一支載入器還沒查——**下一輪先查那個，不要預設同一支**
（`CLAUDE.md` §8 第 9 條）。
