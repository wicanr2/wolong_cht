#!/usr/bin/env python3
"""正常軍團清單、外推選點、返回父清單與存檔回歸。"""
import hashlib
import json
from pathlib import Path
import time

from verify_desktop_polish import Player, OUT, difference

player = Player()
try:
    player.start('march')
    player.click(352, 16)
    player.click(240, 47)
    player.click(240, 96)
    player.shot('march-list')
    player.click(200, 112)
    player.click(200, 112)
    player.shot('march-pick')
    player.click(768, 112, 0)
    player.vx = 639
    right = player.shot('march-right')
    time.sleep(2)
    stopped = player.shot('march-stopped')
    metric = difference(right, stopped, '350x240+32+64')
    assert metric == 0, metric
    player.click(-256, 112, 0)
    player.vx = 0
    player.shot('march-left')
    player.click(152, 120)
    target = player.shot('march-target')
    reference = Path('/src/dist-all/v.1.0.11-20260908/verification/march/march-target.png')
    target_metric = difference(target, reference, '112x48+144+112')
    assert target_metric == 0, target_metric
    player.click(152, 120, 3)
    player.shot('march-parent-list')
    player.click(152, 120, 3)
    player.click(152, 120, 3)
    player.click(448, 16)
    player.click(376, 160)
    player.click(300, 200)
    player.shot('saved-after-march')
    save = OUT / 'SAVE-slot2.wlsave'
    assert save.is_file(), '正常第二槽存檔未完成'
    (OUT / 'result.json').write_text(json.dumps({
        'stop_difference_pixels': metric, 'target_menu_difference_pixels': target_metric,
        'target_reference': str(reference), 'save_sha256': hashlib.sha256(save.read_bytes()).hexdigest(),
        'scope': '正常行軍選陝、取消回父清單與存第二槽；不宣稱重新完成長程行軍對拍',
    }, ensure_ascii=False, indent=2))
finally:
    player.stop()
