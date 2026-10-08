#!/usr/bin/env python3
"""十三碼事件的其餘 handler、落下尾端與政治／災害依賴。"""
import sys
import traceback
from pathlib import Path
sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_pro
probe.TARGETS=tuple(0x10000+x for x in (
    0x320c,0x3220,0x3262,0x32a9,0x32e9,0x3327,0x3388,0x33ea,0x33fd,
    0x3485,0x34a6,0x34b1,0x351a,0x3526,0x35ab,0x35ed,0x3639,0x3669,
    0x3697,0x36c4,0x3712,0x3771,0x37d8,0x37f5,0x3138,0x4502,0x6a3d,
    0x23ff,0x2438,0x50d7))
try:
    probe.main()
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(),encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
