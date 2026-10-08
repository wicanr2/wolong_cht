#!/usr/bin/env python3
"""時鐘每時更新與事件分派核心，原始 IDA 定位／指令／xref 匯出。"""
import sys
import traceback
from pathlib import Path
sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_pro
probe.TARGETS=(0x13E11,0x13E65,0x13E8E,0x15673,0x131AE,0x13496,0x13507,0x13DC9)
try:
    probe.main()
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(),encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
