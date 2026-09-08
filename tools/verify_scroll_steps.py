#!/usr/bin/env python3
"""量測 dosgolem 逐像素外推收據的實際畫面位移；不修改輸入圖片。"""
import hashlib
import json
from pathlib import Path
import sys

from verify_visual_residuals import pixels

root = Path(sys.argv[1])
axis = sys.argv[2]
assert axis in ('x', 'y')
base = pixels(root / 'edge.png')
results = []
for n in range(1, 17):
    path = root / f'plus-{n}.png'
    raw = pixels(path)
    scores = []
    for shift in range(18):
        total = 0
        for y in range(64, 184):
            for x in range(32, 332):
                a = (y*640+x)*3
                b = ((y+(shift if axis == 'y' else 0))*640+x+(shift if axis == 'x' else 0))*3
                total += raw[a:a+3] != base[b:b+3]
        scores.append(total)
    best = min(range(18), key=lambda i: scores[i])
    results.append({'origin_delta': n, 'rendered_delta': best, 'difference_pixels': scores[best],
                    'sha256': hashlib.sha256(path.read_bytes()).hexdigest()})
(root / 'rendered-deltas.json').write_text(json.dumps({'axis': axis, 'rect': [32, 64, 300, 120], 'results': results}, indent=2))
assert all(r['rendered_delta'] == r['origin_delta']//16*16 and r['difference_pixels'] == 0 for r in results)
print(f'{axis}：16 組原版收據皆以整格顯示')
