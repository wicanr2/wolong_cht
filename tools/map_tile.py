#!/usr/bin/env python3
"""印出 `MMAP.MAP` 某一塊的圖塊編號。

    tools/py.sh tools/map_tile.py workplace/orig/dosv/MMAP.MAP 200 110 8 6

用途是回答「畫面上那一格是幾號圖塊」。從截圖只看得到**畫出來的樣子**，
看不到編號；而要判斷「原版是不是換了圖塊」就得先知道檔案裡寫的是幾號。

RLE 與 `internal/assets/rle` 同一支（`KI.EXE` 的 `sub_1F5E7`）：
沒有逃脫字元，連續兩個相同的 byte 之後那一個 byte 是「再重複幾次」。
"""
import sys

WIDTH, HEIGHT = 384, 256


# ⭐ 解壓後**開頭**那 4 byte 是長度欄位（`00 80 01 00` ＝ 98,304 ＝ 384×256），
# 不是圖塊。從 offset 0 當圖塊讀不會報錯，只會讓整張地圖左移四格——
# 而那四格會以「據點中心在記錄座標 +4」「鏡頭比看到的那一欄小 4」的形式
# 散到各處，看起來像兩個獨立的原版怪癖（docs/formats/05 §2.1）。
HEADER = 4


def load_map(path):
    """解 MMAP.MAP 並跳過開頭的長度欄位，回傳地圖本體。"""
    raw = decode(open(path, "rb").read())
    want = WIDTH * HEIGHT
    got = int.from_bytes(raw[:HEADER], "little")
    if got != want:
        raise SystemExit("開頭那 4 byte 是 %d，不是 %d——排法與預期不同" % (got, want))
    return raw[HEADER:]


def decode(src):
    out = bytearray()
    i = 0
    n = len(src)
    while i < n:
        prev = src[i]
        out.append(prev)
        i += 1
        matched = False
        while i < n:
            cur = src[i]
            out.append(cur)
            i += 1
            if cur == prev:
                matched = True
                break
            prev = cur
        if not matched or i >= n:
            break
        count = src[i]
        i += 1
        out.extend(bytes([prev]) * count)
    return bytes(out)


def main():
    if len(sys.argv) not in (4, 6):
        raise SystemExit("用法: map_tile.py MMAP.MAP x y [寬 高]")
    path = sys.argv[1]
    x0, y0 = int(sys.argv[2]), int(sys.argv[3])
    w = int(sys.argv[4]) if len(sys.argv) == 6 else 8
    h = int(sys.argv[5]) if len(sys.argv) == 6 else 8
    raw = load_map(path)
    print("解出 %d B（跳過開頭 %d B 長度欄位；地圖本體 %d B，多餘 %d B）"
          % (len(raw) + HEADER, HEADER, WIDTH * HEIGHT,
             len(raw) - WIDTH * HEIGHT))
    print("    " + " ".join("%3d" % (x0 + dx) for dx in range(w)))
    for dy in range(h):
        y = y0 + dy
        row = [raw[y * WIDTH + x0 + dx] for dx in range(w)]
        print("%3d " % y + " ".join("%3d" % v for v in row))


if __name__ == "__main__":
    main()
