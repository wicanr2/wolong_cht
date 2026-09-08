#!/usr/bin/env python3
"""以既有原圖計算差異座標，不平移或遮掉差異來宣稱全圖通過。"""
import hashlib
import json
from pathlib import Path
import subprocess
import sys


def pixels(path):
    return subprocess.check_output(['convert', str(path), '-depth', '8', 'RGB:-'])


def compare(a, b, rect):
    pa, pb = pixels(a), pixels(b)
    assert len(pa) == len(pb) == 640 * 400 * 3
    x0, y0, width, height = rect
    changed = []
    for y in range(y0, y0+height):
        for x in range(x0, x0+width):
            p = (y*640+x)*3
            if pa[p:p+3] != pb[p:p+3]:
                changed.append((x, y))
    return {
        'a': str(a), 'b': str(b), 'rect': rect, 'difference_pixels': len(changed),
        'bbox': ([min(p[0] for p in changed), min(p[1] for p in changed),
                  max(p[0] for p in changed)+1, max(p[1] for p in changed)+1] if changed else None),
        'input_sha256': [hashlib.sha256(p.read_bytes()).hexdigest() for p in (a, b)],
        'coordinates': changed,
    }


if __name__ == '__main__':
    a, b, out = map(Path, sys.argv[1:4])
    rect = list(map(int, sys.argv[4:8])) if len(sys.argv) > 4 else [0, 0, 640, 400]
    result = compare(a, b, rect)
    out.write_text(json.dumps(result, ensure_ascii=False, indent=2))
    print(json.dumps({k: v for k, v in result.items() if k != 'coordinates'}, ensure_ascii=False))
