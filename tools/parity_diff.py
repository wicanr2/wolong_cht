#!/usr/bin/env python3
"""逐區比對原版與 remake 的畫面（規格：`docs/spec/90-same-state-parity.md`）。

    tools/py.sh tools/parity_diff.py 原版.png remake.png [--out 差分.png]
    tools/py.sh tools/parity_diff.py 原版.png remake.png --regions tactical
    tools/py.sh tools/parity_diff.py --selftest

分區座標出自 `docs/re/46`，是**從機器碼算出來的**不是量的——指令列的
432×32 來自 `sub_1614A` 的 `cx=0x021B`，而它又與 `cs:6181h` 那串字的
排版寬度互相印證。

## 為什麼要分區

整張比只會得到一個數字，而那個數字永遠不是 0，於是沒有人知道該修哪裡。
分區之後每一區各自有結論：右欄 FAIL 而地圖 PASS，就知道問題在版面常數
不在圖庫解碼。

## 判定

    PASS   不同像素 ＝ 0
    NEAR   ≤ 該區 0.5%
    FAIL   其餘

**這是保存專案，PASS 才是目標**；NEAR 只是把「差在少數像素」與
「整片不對」分開，不是及格。

## ⚠ 這支工具自己也要有正對照

`--selftest` 拿同一張圖比自己（每區必須 0），再把圖平移 1 px 比
（每區必須不是 0）。少了後者，一支「永遠回 0」的壞工具會讓所有對拍
看起來全部通過——而那正是最想避免的失敗模式。
"""
import argparse
import struct
import sys
import zlib

# 分區。座標 (x, y, w, h)，**每一個都出自機器碼**（docs/spec/12 §1）：
# 橫幅是 sub_18755 的熱區 6，另外三個是各自的 sub_1895D 呼叫。
# 右欄三段相加 32 + 160 + 208 = 400，剛好鋪滿畫面高度。
REGIONS = [
    ("banner", 0, 0, 640, 32),
    ("command", 0, 32, 432, 32),
    ("map", 0, 64, 432, 336),
    ("minimap", 432, 32, 208, 160),
    ("faction", 432, 192, 208, 208),
]

# 戰場的分區（docs/spec/91 §2）。側欄七段相加 48+32+128+40+32+96+24 = 400。
# ⭐ 側欄切成七段是為了把「圖庫解錯」與「狀態沒對上」分開：
# sb-formation／sb-command／sb-arrow 是純美術，不隨戰況變。
TACTICAL_REGIONS = [
    ("field", 0, 0, 480, 368),
    ("bottom", 0, 368, 480, 32),
    ("sb-title", 480, 0, 160, 48),
    ("sb-enemy", 480, 48, 160, 32),
    ("sb-minimap", 480, 80, 160, 128),
    ("sb-self", 480, 208, 160, 40),
    ("sb-formation", 480, 248, 160, 32),
    ("sb-command", 480, 280, 160, 96),
    ("sb-arrow", 480, 376, 160, 24),
]

REGION_SETS = {"strategy": REGIONS, "tactical": TACTICAL_REGIONS}

NEAR_RATIO = 0.005


def read_png(path):
    """回傳 (width, height, pixels)，pixels 是每點 (r, g, b) 的 list of rows。

    支援 8-bit 的灰階／truecolor／truecolor+alpha／**調色盤**，
    調色盤另外支援 **1/2/4 bit**。
    調色盤那一項是必要的：兩邊的擷取管線對 16 色畫面都會輸出 color type 3，
    少了它每一次對拍都要先手動轉檔，而「忘了轉」與「真的不一樣」看起來一樣。

    ⚠ **sub-byte 深度不是理論情況**：原版主畫面只用 16 色，ImageMagick
    因此輸出 `depth=4`，同一條管線的別張圖卻是 `depth=8`——
    差別只在那張圖用了幾個顏色。少了這一段，對拍會在「剛好夠簡單的畫面」
    上失敗，而那正是最該拿來比的畫面。

    遇到別的格式**明確報錯**，不要猜著解——猜錯會產生看似合理的差分圖。
    """
    data = open(path, "rb").read()
    if data[:8] != b"\x89PNG\r\n\x1a\n":
        raise ValueError("%s 不是 PNG" % path)
    pos, idat, w, h, depth, color, plte = 8, [], 0, 0, 0, 0, b""
    while pos < len(data):
        (length,) = struct.unpack(">I", data[pos:pos + 4])
        typ = data[pos + 4:pos + 8]
        body = data[pos + 8:pos + 8 + length]
        if typ == b"IHDR":
            w, h, depth, color = struct.unpack(">IIBB", body[:10])
        elif typ == b"PLTE":
            plte = body
        elif typ == b"IDAT":
            idat.append(body)
        elif typ == b"IEND":
            break
        pos += 12 + length
    ok = (depth == 8 and color in (0, 2, 3, 6)) or (depth in (1, 2, 4) and color == 3)
    if not ok:
        raise ValueError("%s 是 depth=%d color=%d，這支工具只認 8-bit 灰階／RGB／RGBA／"
                         "調色盤，或 1/2/4-bit 調色盤" % (path, depth, color))
    if color == 3 and not plte:
        raise ValueError("%s 是調色盤 PNG 卻沒有 PLTE" % path)
    channels = {0: 1, 2: 3, 3: 1, 6: 4}[color]
    raw = zlib.decompress(b"".join(idat))
    # sub-byte 深度的解濾波以 byte 為單位（濾波器不認識像素邊界），
    # 展開成一像素一 byte 要等解完濾波之後。
    stride = (w * depth + 7) // 8 if depth < 8 else w * channels
    # 濾波器的「左邊那一點」是 bpp（byte 為單位，不足 1 byte 算 1）。
    bpp = max(1, channels * depth // 8)
    out, prev, p = [], bytearray(stride), 0
    for _ in range(h):
        f = raw[p]
        line = bytearray(raw[p + 1:p + 1 + stride])
        p += 1 + stride
        for i in range(stride):
            a = line[i - bpp] if i >= bpp else 0
            b = prev[i]
            c = prev[i - bpp] if i >= bpp else 0
            if f == 1:
                line[i] = (line[i] + a) & 0xFF
            elif f == 2:
                line[i] = (line[i] + b) & 0xFF
            elif f == 3:
                line[i] = (line[i] + (a + b) // 2) & 0xFF
            elif f == 4:
                pa, pb, pc = abs(b - c), abs(a - c), abs(a + b - 2 * c)
                pr = a if (pa <= pb and pa <= pc) else (b if pb <= pc else c)
                line[i] = (line[i] + pr) & 0xFF
        prev = line
        if depth < 8:
            idx, mask = [], (1 << depth) - 1
            for byte in line:
                for shift in range(8 - depth, -1, -depth):
                    idx.append((byte >> shift) & mask)
            idx = idx[:w]
            out.append([tuple(plte[3 * v:3 * v + 3]) for v in idx])
        elif color == 3:
            out.append([tuple(plte[3 * v:3 * v + 3]) for v in line])
        elif channels == 1:
            out.append([(v, v, v) for v in line])
        else:
            out.append([tuple(line[i:i + 3]) for i in range(0, stride, channels)])
    return w, h, out


def write_png(path, w, h, rows):
    def chunk(typ, body):
        return (struct.pack(">I", len(body)) + typ + body
                + struct.pack(">I", zlib.crc32(typ + body) & 0xFFFFFFFF))
    raw = b"".join(b"\x00" + bytes(v for px in row for v in px) for row in rows)
    with open(path, "wb") as fh:
        fh.write(b"\x89PNG\r\n\x1a\n")
        fh.write(chunk(b"IHDR", struct.pack(">IIBBBBB", w, h, 8, 2, 0, 0, 0)))
        fh.write(chunk(b"IDAT", zlib.compress(raw, 9)))
        fh.write(chunk(b"IEND", b""))


def compare(a, b, w, h, regions=REGIONS):
    """回傳每一區的統計，以及一張差分列（不同的畫紅色）。"""
    diff = [[(0, 0, 0)] * w for _ in range(h)]
    stats = []
    for name, rx, ry, rw, rh in regions:
        n = worst = 0
        total = 0
        for y in range(ry, min(ry + rh, h)):
            for x in range(rx, min(rx + rw, w)):
                total += 1
                pa, pb = a[y][x], b[y][x]
                if pa == pb:
                    diff[y][x] = (pa[0] // 4, pa[1] // 4, pa[2] // 4)
                    continue
                n += 1
                d = max(abs(pa[i] - pb[i]) for i in range(3))
                worst = max(worst, d)
                diff[y][x] = (255, 0, 0)
        ratio = n / total if total else 0.0
        verdict = "PASS" if n == 0 else ("NEAR" if ratio <= NEAR_RATIO else "FAIL")
        stats.append((name, n, total, ratio, worst, verdict))
    return stats, diff


def selftest():
    w, h = 640, 400
    base = [[((x * 7 + y * 3) % 256, (x * 5) % 256, (y * 11) % 256)
             for x in range(w)] for y in range(h)]
    shifted = [row[1:] + [row[-1]] for row in base]
    for set_name, regions in REGION_SETS.items():
        stats, _ = compare(base, base, w, h, regions)
        for name, n, _, _, _, verdict in stats:
            if n != 0 or verdict != "PASS":
                print("✗ %s：同一張圖比自己，%s 竟然有 %d 個不同像素" % (set_name, name, n))
                return 1
        stats, _ = compare(base, shifted, w, h, regions)
        for name, n, _, _, _, _ in stats:
            if n == 0:
                print("✗ %s：平移 1 px 之後 %s 仍回 0——這支工具沒有真的在比"
                      % (set_name, name))
                return 1
    print("✓ 正對照通過（%s）：同圖每區 0、平移 1 px 每區非 0"
          % "／".join(REGION_SETS))
    return 0


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("original", nargs="?")
    ap.add_argument("remake", nargs="?")
    ap.add_argument("--out", help="差分圖輸出路徑")
    ap.add_argument("--selftest", action="store_true")
    ap.add_argument("--rect", help="改比一個自訂矩形 `x,y,w,h`（單一視窗用）")
    ap.add_argument("--regions", choices=sorted(REGION_SETS), default="strategy",
                    help="分區組：strategy ＝ 主畫面五區、tactical ＝ 戰場九區")
    ns = ap.parse_args()
    if ns.selftest:
        return selftest()
    if not ns.original or not ns.remake:
        ap.error("要兩張圖，或 --selftest")

    aw, ah, a = read_png(ns.original)
    bw, bh, b = read_png(ns.remake)
    if (aw, ah) != (bw, bh):
        # 縮放後硬比會把「尺寸不對」這個最重要的線索洗掉。
        print("✗ 尺寸不同：原版 %dx%d、remake %dx%d——先把兩邊調成同一個模式"
              % (aw, ah, bw, bh))
        return 2

    regions = REGION_SETS[ns.regions]
    if ns.rect:
        try:
            rx, ry, rw, rh = (int(v) for v in ns.rect.split(","))
        except ValueError:
            ap.error("--rect 要四個整數 x,y,w,h")
        regions = [("rect", rx, ry, rw, rh)]
    stats, diff = compare(a, b, aw, ah, regions)
    print("| 區 | 不同像素 | 佔比 | 最大色差 | 判定 |")
    print("|---|---:|---:|---:|---|")
    worst = "PASS"
    for name, n, total, ratio, d, verdict in stats:
        print("| `%s` | %d / %d | %.2f%% | %d | %s |" % (name, n, total, ratio * 100, d, verdict))
        if verdict == "FAIL" or (verdict == "NEAR" and worst == "PASS"):
            worst = verdict
    if ns.out:
        write_png(ns.out, aw, ah, diff)
        print("\n差分圖：%s（不同的畫紅色，相同的壓暗）" % ns.out)
    return 0 if worst == "PASS" else 1


if __name__ == "__main__":
    sys.exit(main())
