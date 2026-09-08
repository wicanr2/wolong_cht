#!/usr/bin/env python3
"""正常外推與 dosgolem 原版逐像素對照，不注入相機或遊戲狀態。"""
import json
import ctypes
from pathlib import Path
import time
import subprocess
import os
import signal

from verify_desktop_polish import Player, OUT, difference, run, record

player = Player()
results = []
original = Path('/src/workplace/desktop-polish-20260908')
focus_wait = float(os.environ.get('WOLONG_FOCUS_WAIT', '2'))
focus_repeats = int(os.environ.get('WOLONG_FOCUS_REPEATS', '1'))
focus_freeze = os.environ.get('WOLONG_FOCUS_FREEZE') == '1'


def prepare(name):
    player.start(name)
    player.click(448, 16)
    # 正常設定最低戰略速度，避免多張截圖期間自然事件遮住地圖取樣區。
    pref = Path(os.environ['XDG_CONFIG_HOME']) / 'wolong-remake/preferences.json'
    speed = json.loads(pref.read_text())['strategy_speed'] if pref.exists() else 2
    while speed != 4:
        player.click(376, 232)
        speed = json.loads(pref.read_text())['strategy_speed']
    player.park_pointer_in_menu()
    player.click(320, 224, 0)
    # 逐像素量測須先由畫出的箭頭回讀起點，不能把 XTest 指令量當成實際位置。
    mask = ('WW............', 'W#WW..........', '.W##WW........', '.W####WW......',
            '..W#####WW....', '..W#######WW..', '...W########WW', '...W#######W..',
            '....W#####W...', '....W######W..', '.....W##W###W.', '.....W#W.W###W',
            '......W...W##W', '......W....WW.')
    for attempt in range(3):
        shot = player.shot(f'{name}-pointer-{attempt}')
        raw = subprocess.check_output(['convert', str(shot), '-depth', '8', 'RGB:-'])
        def color(x, y):
            i = (y*640+x)*3
            return raw[i:i+3]
        found = []
        for y in range(208, 245):
            for x in range(300, 341):
                colors = {'W': color(x, y), '#': color(x+1, y+1)}
                if colors['W'] == colors['#']:
                    continue
                if all(ch == '.' or color(x+dx, y+dy) == colors[ch]
                       for dy, row in enumerate(mask) for dx, ch in enumerate(row)):
                    found.append((x, y))
        assert len(found) == 1, found
        record('observed-pointer', *found[0])
        player.vx, player.vy = found[0]
        if found[0] == (320, 224):
            break
        player.click(320, 224, 0)
    else:
        raise AssertionError('游標起點未能對齊')
    player.click(320, 224, 3)


def compare(path, reference, crop, label):
    pixels = difference(path, reference, crop)
    results.append({'label': label, 'difference_pixels': pixels, 'rect': crop})
    (OUT / 'comparisons.json').write_text(json.dumps(results, ensure_ascii=False, indent=2))
    assert pixels == 0, f'{label} 差 {pixels} 像素'


try:
    for axis in ('x', 'y'):
        prepare(f'start-{axis}')
        initial = player.shot(f'{axis}-initial')
        if axis == 'x':
            compare(initial, original/'original-scroll/loaded.png', '16x16+320+224', '一般地圖格框')
        base = original / ('original-scroll' if axis == 'x' else 'original-scroll-y')
        for n in range(1, 17):
            player.click(640 if axis == 'x' else 320, 400 if axis == 'y' else 224, 0)
            player.vx, player.vy = min(639, player.vx), min(399, player.vy)
            image = player.shot(f'{axis}-plus-{n:02d}')
            compare(image, base/f'plus-{n}.png', '560x240+32+64', f'{axis} 外推 {n} 像素')
        time.sleep(2)
        stopped = player.shot(f'{axis}-stopped')
        compare(stopped, image, '560x240+32+64', f'{axis} 停手')
        # 回復焦點不可把視窗外的移動算成外推。
        player.x11.XDefaultRootWindow.argtypes = [ctypes.c_void_p]
        player.x11.XDefaultRootWindow.restype = ctypes.c_ulong
        player.x11.XSetInputFocus.argtypes = [ctypes.c_void_p, ctypes.c_ulong, ctypes.c_int, ctypes.c_ulong]
        for repeat in range(focus_repeats):
            record('focus-cycle', axis, repeat, focus_wait)
            if focus_freeze:
                os.kill(player.app.pid, signal.SIGSTOP)
            try:
                player.x11.XSetInputFocus(player.display, player.x11.XDefaultRootWindow(player.display), 1, 0)
                player.x11.XFlush(player.display)
                time.sleep(focus_wait)
                run('xdotool', 'mousemove', '100', '100')
                time.sleep(focus_wait)
                run('xdotool', 'windowfocus', player.wid)
            finally:
                if focus_freeze:
                    os.kill(player.app.pid, signal.SIGCONT)
            time.sleep(focus_wait)
            focused = player.shot(f'{axis}-refocused-{repeat}')
            # 回焦可能改變游標位置。兩個分離的固定地標區驗相機，完整圖仍保留。
            for rect in ('200x80+32+64', '200x64+352+256'):
                compare(focused, stopped, rect, f'{axis} 焦點回復 {repeat}')
            time.sleep(2)
            settled = player.shot(f'{axis}-refocused-settled-{repeat}')
            for rect in ('200x80+32+64', '200x64+352+256'):
                compare(settled, stopped, rect, f'{axis} 回焦後停手 {repeat}')
        player.stop()
    (OUT / 'result.json').write_text(json.dumps({'passed': True, 'comparisons': results}, ensure_ascii=False, indent=2))
finally:
    player.stop()
