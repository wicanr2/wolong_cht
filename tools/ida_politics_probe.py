#!/usr/bin/env python3
"""月結政治與俘虜依賴閉包的原始 IDA 證據匯出。"""
import sys
import traceback
from pathlib import Path
sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_pro
probe.TARGETS = (0x1585F,0x15899,0x15940,0x12AD2,0x15990,0x1301C,
                 0x12BD9,0x12C52,0x12CDF,0x12D3A,0x12FB1,0x12D58,
                 0x12DB8,0x12DF3,0x130F0,0x1310A,0x13091,0x12E33,
                 0x12E89,0x12EFB,0x12F71)
try:
    probe.main()
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
