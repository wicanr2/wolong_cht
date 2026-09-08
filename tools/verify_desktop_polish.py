#!/usr/bin/env python3
"""在隔離 Xvfb 中以正常輸入驗證桌面偏好與戰後結果；由同名 shell 啟動。"""
import ctypes
from collections import Counter
import json
import os
from pathlib import Path
import subprocess as sp
import time

OUT = Path('/out')
TRACE = []


def record(*event):
    TRACE.append([time.monotonic(), list(event)])
    (OUT / 'input-trace.json').write_text(json.dumps(TRACE, ensure_ascii=False, indent=2))


def run(*args):
    return sp.check_output(args, text=True, timeout=12).strip()


class Player:
    def __init__(self):
        self.app = None
        self.log = None
        self.x11 = ctypes.CDLL('libX11.so.6')
        self.xt = ctypes.CDLL('libXtst.so.6')
        self.x11.XOpenDisplay.restype = ctypes.c_void_p
        self.display = self.x11.XOpenDisplay(None)
        self.xt.XTestFakeRelativeMotionEvent.argtypes = [ctypes.c_void_p, ctypes.c_int, ctypes.c_int, ctypes.c_ulong]
        self.x11.XFlush.argtypes = [ctypes.c_void_p]

    def start(self, name):
        self.log = (OUT / f'{name}.log').open('w')
        self.app = sp.Popen(['/tmp/squashfs-root/AppRun', '-save-file', '/out/SAVE.DAT'], stdout=self.log, stderr=sp.STDOUT)
        record('start', name)
        time.sleep(5)
        self.wid = run('xdotool', 'search', '--onlyvisible', '--name', '臥龍傳').splitlines()[0]
        run('xdotool', 'windowfocus', self.wid)
        geo = dict(line.split('=') for line in run('xdotool', 'getwindowgeometry', '--shell', self.wid).splitlines())
        self.x, self.y, self.w, self.h = [int(geo[k]) for k in ('X', 'Y', 'WIDTH', 'HEIGHT')]
        self.vx = self.vy = self.clicks = 0
        self.click(320, 196)
        self.click(300, 152)

    def stop(self):
        if self.app is not None:
            record('terminate-process')
            self.app.terminate()
            self.app.wait(timeout=12)
            self.app = None
            self.log.close()

    def click(self, a, b, button=1, hold=.35):
        record('click', a, b, button, hold)
        if self.clicks < 2:
            run('xdotool', 'mousemove', str(self.x+a*self.w//640), str(self.y+b*self.h//400))
            time.sleep(.3)
        else:
            dx, dy = a-self.vx, b-self.vy
            while dx or dy:
                sx, sy = max(-16, min(16, dx)), max(-16, min(16, dy))
                self.xt.XTestFakeRelativeMotionEvent(self.display, sx*self.w//640, sy*self.h//400, 0)
                self.x11.XFlush(self.display)
                time.sleep(.08)
                dx -= sx
                dy -= sy
        self.vx, self.vy = a, b
        self.clicks += 1
        if button == 0:
            time.sleep(.5)
            return
        run('xdotool', 'mousedown', str(button))
        time.sleep(hold)
        run('xdotool', 'mouseup', str(button))
        time.sleep(.5)

    def key(self, name):
        record('key', name)
        run('xdotool', 'keydown', name)
        time.sleep(.5)
        run('xdotool', 'keyup', name)
        time.sleep(.5)

    def park_pointer_in_menu(self):
        # 系統視窗阻擋地圖捲動時，把捕捉游標送到已知邊界；避免 Xvfb 相對事件
        # 與繪圖不同步造成箭頭殘差。此操作只用於開著的系統視窗。
        record('park-pointer-in-system-menu')
        for _ in range(48):
            self.xt.XTestFakeRelativeMotionEvent(self.display, -16*self.w//640, -16*self.h//400, 0)
            self.x11.XFlush(self.display)
            time.sleep(.12)
        self.vx = self.vy = 0
        time.sleep(.5)

    def shot(self, name):
        path = OUT / f'{name}.png'
        run('import', '-window', 'root', '/tmp/screen.png')
        run('convert', '/tmp/screen.png', '-crop', f'{self.w}x{self.h}+{self.x}+{self.y}', '+repage', '-filter', 'point', '-resize', '640x400!', f'PNG24:{path}')
        record('shot', name)
        return path


def difference(a, b, crop):
    for path, target in [(a, '/tmp/a.png'), (b, '/tmp/b.png')]:
        run('convert', str(path), '-crop', crop, '+repage', target)
    p = sp.run(['compare', '-metric', 'AE', '/tmp/a.png', '/tmp/b.png', 'null:'], capture_output=True, text=True, timeout=12)
    if p.returncode not in (0, 1):
        raise RuntimeError(p.stderr)
    return float(p.stderr.strip())


def shape_difference(a, b, crop):
    # 四季色盤不同時，以一對一顏色對應驗固定圖案，不能把純色背景誤認為視窗。
    images = [sp.check_output(['convert', str(p), '-crop', crop, '+repage', '-depth', '8', 'RGB:-'], timeout=12) for p in (a, b)]
    colors = [[raw[i:i+3] for i in range(0, len(raw), 3)] for raw in images]
    pairs = Counter(zip(*colors))
    used_a, used_b, matched = set(), set(), 0
    for (ca, cb), count in pairs.most_common():
        if ca not in used_a and cb not in used_b:
            used_a.add(ca)
            used_b.add(cb)
            matched += count
    return len(colors[0]) - matched


RESULT_REFERENCE = Path('/src/workplace/battle-return-20260908/app/post-7.png')
BATTLE_REFERENCE = Path('/src/workplace/player-march-20260908/battle-ready-original/select.png')


def result_visible(path):
    return shape_difference(path, RESULT_REFERENCE, '16x8+72+104') == 0


def wait_battle(player, prefix):
    for i in range(48):
        path = player.shot(f'{prefix}-encounter-{i:02d}')
        if shape_difference(path, BATTLE_REFERENCE, '128x32+496+248') < 100:
            break
        player.key('Return')
        time.sleep(3)
    else:
        raise AssertionError('未進入正常遭遇')
    for i in range(60):
        path = player.shot(f'{prefix}-opening-{i:02d}')
        if shape_difference(path, BATTLE_REFERENCE, '8x8+0+0') == 0:
            return
        time.sleep(1)
    raise AssertionError('開場等待未結束')


def form_and_wait(player, prefix):
    player.click(352, 16)
    player.click(192, 47)
    player.click(150, 112)
    player.click(150, 112)
    player.shot(f'{prefix}-formation')
    player.click(324, 279)
    player.click(320, 200)
    player.click(320, 200, 3)
    player.click(320, 200, 3)
    wait_battle(player, prefix)


def retreat(player, prefix):
    player.click(510, 367, hold=2)
    player.click(510, 367, hold=2)
    previous = time.monotonic()
    for i in range(70):
        path = player.shot(f'{prefix}-after-{i:02d}')
        now = time.monotonic()
        if result_visible(path):
            return previous, now
        previous = now
        time.sleep(.2)
    raise AssertionError('退卻後未看見結果頁')


def main():
    player = Player()
    try:
        player.start('first-start')
        player.click(448, 16)
        player.shot('defaults')
        for y in (184, 208, 232, 256, 304, 328, 352):
            # 實際速度列中心是 232、256；原版列間距 24。
            player.click(376, y)
        player.park_pointer_in_menu()
        before = player.shot('preferences-before-restart')
        pref_path = Path(os.environ['XDG_CONFIG_HOME']) / 'wolong-remake/preferences.json'
        pref = json.loads(pref_path.read_text())
        assert pref['battle_result_seconds'] == 3 and pref['lord_corps'] is False
        assert pref['strategy_speed'] == 3 and pref['tactical_speed'] == 3
        assert pref['video_lcd'] is True and pref['sound'] == 2 and pref['damage_report'] is True
        player.stop()
        raw = pref_path.read_bytes()
        player.start('restart')
        player.click(448, 16)
        player.click(376, 352, button=3)  # 右鍵只關系統，先重新開啟以確認沒有改值。
        player.click(448, 16)
        player.park_pointer_in_menu()
        after = player.shot('preferences-after-restart')
        metric = difference(before, after, '208x264+208+112')
        assert metric == 0, f'重啟後設定標籤差異 {metric}'
        assert pref_path.read_bytes() == raw, '啟動或右鍵取消不應寫回偏好'
        (OUT / 'preferences-result.json').write_text(json.dumps({'restart': True, 'panel_difference_pixels': metric}))
        if os.environ.get('WOLONG_VERIFY_SCOPE') == 'preferences':
            return
        # 回到正常彩色／無損害報告，速度加至最高；結果頁設 30 秒驗提前關閉。
        player.click(376, 184)
        player.click(376, 328)
        for y in (232, 256):
            for _ in range(2):
                player.click(376, y)
        for _ in range(4):
            player.click(376, 352)
        player.shot('result-30-seconds')
        assert json.loads(pref_path.read_text())['battle_result_seconds'] == 30
        player.click(300, 200, 3)
        form_and_wait(player, 'first')
        _, first_seen = retreat(player, 'first')
        player.click(510, 367, hold=2)
        for i in range(5):
            closed = player.shot(f'early-dismissed-{i}')
            if not result_visible(closed):
                break
            time.sleep(.5)
        early_elapsed = time.monotonic() - first_seen
        assert not result_visible(closed) and early_elapsed < 30
        (OUT / 'early-result.json').write_text(json.dumps({'dismissed': True, 'seconds_upper_bound': early_elapsed}))
        # 高速策略下退卻軍團會再次自然遭遇，直接驗第二場，不在兩場之間注入或讀檔。
        wait_battle(player, 'second')
        first_lower, first_upper = retreat(player, 'second')
        last_visible = first_upper
        observations = []
        for i in range(100):
            path = player.shot(f'second-countdown-{i:02d}')
            now = time.monotonic()
            visible = result_visible(path)
            observations.append([now, visible])
            if not visible:
                break
            last_visible = now
            time.sleep(.2)
        else:
            raise AssertionError('第二場結果頁未自動關閉')
        bounds = [last_visible-first_upper, now-first_lower]
        # 精確 30 秒閾值由固定時間單元測試驗證。GUI 包含 Update／呈現延遲，
        # 這裡只驗重設後保留約 30 秒且有界返回，不能拿截圖時間當毫秒級時鐘。
        assert bounds[0] >= 25 and bounds[1] <= 35, bounds
        (OUT / 'result.json').write_text(json.dumps({
            'preferences_restart': True, 'settings_label_difference_pixels': metric,
            'early_dismiss_seconds_upper_bound': early_elapsed,
            'second_battle': '同一程序連續自然遭遇，兩場之間不讀檔',
            'second_countdown_bounds_seconds': bounds, 'second_observations': observations,
            'native_windows_macos': '人工驗收',
        }, ensure_ascii=False, indent=2))
    finally:
        player.stop()


if __name__ == "__main__":
    main()
