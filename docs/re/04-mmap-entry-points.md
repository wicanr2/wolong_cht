# 04 — 大地圖：入口點與記憶體佈局（進行中）

**狀態：只找到入口，格式還沒解。** 這份是為了不讓下一輪重找一次。

- 日期：2026-08-07
- 輸入：`workplace/ida/dosv/KI.EXE.i64`　SHA-256 `fffeba98…3868`
- 推論等級：**位址與載入關係 confirmed，資料語意全部未解**

## 1. 三個檔的載入處

| 檔案 | 大小 | 檔名偏移 | 載入處 | 目的 |
|---|---:|---|---|---|
| `MMAP.MAP` | 80,716 | `0DAFh` | `sub_1006B`（開機） | 段存在 `word_10D44` |
| `MMAP.MDL` | 32,768 | `0DB8h` | `sub_187AF` | 段 `word_19876` |
| `MMAP.MCH` | 43,058 | `0DC1h` | `sub_187AF` | 段 `word_19876 + 800h`（緊接在 `MDL` 之後） |

`MMAP.MDL` 剛好 `0x8000` ＝ 32,768 B，而 `MCH` 被放在 `+0x800` 段
（＝ +32,768 B）——**兩個檔在記憶體裡是連續的一整塊**。

> `MMAP.MDL` ÷ 128 ＝ **256**。128 B 正好是一塊 16×16 的 4bpp 圖
> （`sub_1CA3B` 用 `ax=1001h` 畫 16×16 時就是 128 B／塊）。
> **假說：`MDL` 是 256 塊 16×16 的地形圖塊。** 還沒驗——
> 要找到讀這個緩衝區的繪製呼叫，看它的 `ax` 參數。

## 2. `sub_187CC`：記憶體佈局

開機時一次算好六個緩衝區的段位址：

```asm
mov  ax, ds:0D46h        ; 基底
add  ax, 1800h  → ds:9872h      ; 98,304 B
add  ax, 1800h  → ds:9876h      ; 98,304 B  ← MMAP.MDL + MCH 放這裡
add  ax, 1200h  → ds:987Ah      ; 73,728 B
add  ax, 410h   → ds:9878h      ; 16,640 B
add  ax, 1CCh   → ds:0D42h      ;  7,360 B
add  ax, 0FAh   → ds:9874h      ;  4,000 B
```

（單位是段＝16 byte。）

**`word_19872` 的 98,304 B 是渲染用的畫布**：`sub_119CA` 會把它整塊清零
（`rep stosw` 0x8000 words ＋ 0x4000 words ＝ 0x18000 B ＝ 98,304）。

`word_19876` 這一塊**會被重複使用**——`sub_187AF` 放 `MMAP.MDL`＋`MCH`，
但 `sub_13D09` 也把 `IVENTGRF.DAT` 載到同一個段。
**寫筆記時要標明是在哪個時間點**，不要把兩件事混成一件。

## 3. `sub_119CA`：大地圖初始化的主流程

```asm
call sub_187CC                    ; 算緩衝區位址
mov  ax, word_19876 / bx, word_19878 / cx, word_10D44
call sub_1D46A                    ; ← 未解，吃 MMAP.MAP 的段
mov  es, cs:word_10D42 / xor di,di
call sub_1E3C0                    ; ← 未解
mov  ax, word_19874 / bx, word_10D44 / cx, word_19876 / dx, word_19872
call sub_1E48A                    ; ← 未解，四個緩衝區都吃 → 疑似展開／解碼
call sub_187AF                    ; 載入 MMAP.MDL + MCH
call sub_18755                    ; 載入 ICONGRF 段 1
mov ax,7 / cx,0 / dx,17FFh : call sub_20000
mov ax,8 / cx,0 / dx,101Fh : call sub_20000
mov ax,2                   : call sub_20000
…清空 word_19872 的 98,304 B…
…依 byte_198A6 的位元逐一呼叫跳表 off_159D2…
```

**下一輪的三個入口，依投報排**：

1. **`sub_1E48A`** —— 同時吃 `MMAP.MAP`、`MDL`、`19872`、`19874` 四個緩衝區，
   最可能是「把地圖資料展開成畫布」的那一支。解開它就解開整條鏈。
2. **`sub_1D46A`** —— 吃 `MMAP.MAP`，在 `sub_1E48A` 之前跑，疑似索引建立。
3. **`off_159D2` 跳表** —— `byte_198A6` 的每個位元對應一項，
   `bx` 從 2 走到 8（四項）。這是圖層開關（地形／據點／部隊／格線？）的候選。

## 4. 還沒解的

- `MMAP.MAP`（80,716 B）的結構。
- `MMAP.MDL` 是不是 256 塊 16×16 圖塊（§1 的假說）。
- `MMAP.MCH` 是什麼（`MCH` 這個副檔名在 `MOUSE.MCH` 也出現過，588 B）。
- `sub_20000` 的三次呼叫在設定什麼（參數像是範圍：`17FFh`、`101Fh`）。
- `word_19872` 的 98,304 B 畫布怎麼對應到螢幕。
