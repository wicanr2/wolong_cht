#!/usr/bin/env python3
"""據點結算：沿用通用 IDA 探針，保留原始定位與交叉參照。"""
import sys
import traceback
from pathlib import Path
sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_pro
probe.TARGETS = (0x153C6, 0x15456, 0x1548F, 0x15538, 0x15547)
try:
    probe.main()
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
