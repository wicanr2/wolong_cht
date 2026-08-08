# 03 — 圖庫載入器與 VGA 繪製常式

**狀態：READY。圖庫載入器與四平面繪製常式都讀完。**

- 日期：2026-08-07
- 輸入：`workplace/ida/dosv/KI.EXE.i64`　SHA-256 `fffeba98…3868`
- 推論等級：**每筆大小、定位算式、寬高、平面順序全部 confirmed**

## 0. 先講一個工具上的坑

想查「誰用了 `KAOGRF.DAT` 這個字串」，`tools/ida_xref.idc` 回**零筆**。
不是腳本壞了——**IDA 根本沒建立那條 xref**，因為檔名的位址是以立即值出現的：

```asm
mov     dx, 0D79h        ; ← 'KAOGRF.DAT' 在 seg000 的偏移
```

IDA 把它當成一個普通的數字常數，不是位址參考。

> **這正是 `CLAUDE.md` §5.1 說的「手上拿到的是位址，不是符號」。**
> xref 圖只涵蓋 IDA 認得出來的參考；**xref 回零筆不等於沒人用**。
> 這裡的解法是反過來：先從 `.asm` 拿到符號的位址，
> 再 grep 那個立即值。**兩個方向都要會。**

檔名表在 seg000 的偏移：

| 檔名 | 偏移 | 檔名 | 偏移 |
|---|---|---|---|
| `GAMEPAL.BRG` | `0D58h` | `TALK.DAT` | `0D64h` |
| `ICONGRF.DAT` | `0D6Dh` | `KAOGRF.DAT` | `0D79h` |

（檔名是 NUL 結尾、緊密排列，不是固定長度的槽位。）

## 1. `sub_107D2`：頭像載入器 ＋ 4 格快取

```asm
; 進入時 al = 頭像編號
        xor     bx, bx
loc_107E0:
        cmp     al, [bx+846h]       ; 快取表 4 格
        jz      short loc_10828     ; 命中 → 跳過載入
        inc     bx
        cmp     bl, 4
        jb      short loc_107E0

        mov     bl, byte_10845      ; 未命中：round-robin 挑一格
        push    bx
        mov     [bx+846h], al       ; 記錄這格現在放誰
        xchg    bh, bl
        mov     si, bx
        inc     bh
        and     bh, 3               ; ← 4 格輪替
        mov     byte_10845, bh

        mov     dx, 0D79h           ; 'KAOGRF.DAT'
        shl     si, 1
        shl     si, 1
        shl     si, 1               ; si = 快取格 × 8（段位址單位）
        mov     bx, word_10D40      ; 目的段
        mov     ah, al              ; ── 檔案位移計算 ──
        xor     cx, cx
        mov     al, cl              ; ax = 編號 << 8
        shl     ax, 1
        rcl     cx, 1
        shl     ax, 1
        rcl     cx, 1
        shl     ax, 1
        rcl     cx, 1               ; cx:ax = 編號 << 11 = 編號 × 2048
        mov     di, 800h            ; 讀 2,048 byte
        call    sub_1E38C           ; 帶 32-bit 位移的讀檔
```

**定案兩件事：**

- `KAOGRF.DAT` 每筆 **2,048 byte**（`di = 800h`）
- 定位是 **編號 × 2048**（三次 `shl`／`rcl` 的 32 位元左移）

307,200 ÷ 2,048 = **150 筆**，餘 0。

## 2. `sub_1FA37`：四平面繪製

進入時 `ax` 是尺寸參數；頭像的呼叫端是 `mov ax, 4004h`。

常式對 VGA Graphics Controller 做設定，然後**呼叫 `sub_1FAA2` 四次**，
每次改一次 Enable Set/Reset（GC index 1），讓不同的平面吃 CPU 資料：

| 第幾次 | GC index 1 的值 | 二進位 | 寫入的平面 | 寫入函式 |
|---|---|---|---|---|
| 1 | `0Eh` | 1110 | **plane 0** | 直接寫（Set/Reset ＝ 0，其餘平面清零） |
| 2 | `0Dh` | 1101 | **plane 1** | Data Rotate ＝ `10h` ＝ **OR** |
| 3 | `0Bh` | 1011 | **plane 2** | OR |
| 4 | `07h` | 0111 | **plane 3** | OR |

`sub_1FAA2` **不重設來源指標 `si`** ——四次呼叫連續吃掉四段資料。
**所以檔案佈局是 plane-major**：plane0 整張、plane1 整張、plane2、plane3。

目的地 `es = 0A0C8h` → 線性位址 `0A0C80h`，即 VGA 段 `A000` 的偏移 `0C80h`。

## 3. `sub_1FAA2`：複製迴圈，寬高就在這裡

```asm
        push    cx
        mov     dx, cx              ; dh = ch = 40h、dl = cl = 04h
        mov     bp, bx              ; bx = 目的 VRAM 偏移
        xor     ch, ch
loc_1FAA9:                          ; ── 每一列 ──
        mov     dl, cl              ; 每列跑 4 次
        mov     di, bp
loc_1FAAD:
        mov     al, es:[di]         ; 讀 VRAM → 載入 VGA latch
        movsb                       ; es:[di] ← ds:[si]
        mov     al, es:[di]
        movsb
        dec     dl
        jnz     short loc_1FAAD     ; 4 次 × 2 byte = 每列 8 byte
        add     bp, 50h             ; 下一列：+80 byte
        dec     dh
        jnz     short loc_1FAA9     ; 64 列
```

參數 `ax = 4004h` 拆開來就是尺寸：

| 欄位 | 值 | 意義 |
|---|---|---|
| `ah` | `40h` ＝ 64 | **列數（高）** |
| `al` | `04h` | 每列的迴圈次數，每次 2 byte → **每列 8 byte ＝ 64 pixel（寬）** |

列跨距 `50h` ＝ 80 byte → **螢幕寬 640 pixel**。

驗算：8 byte／列 × 64 列 × 4 平面 = **2,048 byte**，
與載入器的 `di = 800h` 完全吻合。**兩條獨立的證據對上。**

`mov al, es:[di]` 是 VGA 的 latch 載入慣用法——先讀一次讓四個平面的資料進 latch，
寫入時未選中的平面才會維持原值。**這行不是死碼，拿掉會壞。**

## 4. 驗收

`tools/grf.py` 照這份規格解 `KAOGRF.DAT`：150 張 64×64，**餘 0 byte**，
配 `GAMEPAL.BRG` 第 0 組渲染出可辨識的武將頭像
（`docs/images/kaogrf-sheet.png`）。膚色、盔甲、頭盔的顏色都正確
——這同時反向驗證了 `docs/formats/02` 的調色盤通道順序。

## 5. 還沒解的

- `ICONGRF`／`KYOGRF`／`IVENTGRF` 的尺寸。**照 §1 的路徑各走一次**：
  grep 檔名偏移的立即值（`0D6Dh` 等）→ 載入器 → 繪製呼叫端的 `ax`。
- `sub_1E38C`（帶位移的讀檔）與 `sub_1F4A2`（開檔／讀檔／關檔）的完整介面。
- `sub_1FAC2` 是另一支繪製常式（`shl al, 1` 後才 `mov cx, ax`），用途未解。
