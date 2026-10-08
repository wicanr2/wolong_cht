#!/usr/bin/env python3
"""月結世界更新：保留原始函式、operand、bytes、chunk 與 xref。"""
import sys
import traceback
from pathlib import Path
sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_pro
probe.TARGETS = (0x15695,0x155A6,0x12FBF,0x157FE,0x15715,0x1578F,
                 0x130CB,0x13119,0x122DB,0x12286,0x1237E)
try:
    probe.main()
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
