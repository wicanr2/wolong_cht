#!/usr/bin/env python3
"""月結 C 還原：沿用通用 IDA 探針，保留原始定位、指令與交叉參照。"""
import sys
import traceback
from pathlib import Path

sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_pro

probe.TARGETS = (0x15358, 0x15609, 0x1563B, 0x154FC, 0x155EC, 0x15828)
try:
    probe.main()
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
