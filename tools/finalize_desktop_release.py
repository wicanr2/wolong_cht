#!/usr/bin/env python3
"""整理桌面交付清單與影片來源；不執行任何遠端發布。"""
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import sys

version = sys.argv[1]
assert re.fullmatch(r'v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}', version)
root = Path('/src')
dest = root / 'dist-all' / version
assert dest.stat().st_uid == os.getuid()
public, promo, verification = (dest / name for name in ('release', 'promo', 'verification'))


def sha(p):
    with p.open('rb') as f:
        return hashlib.file_digest(f, 'sha256').hexdigest()


def write(p, data):
    p.write_text(json.dumps(data, ensure_ascii=False, indent=2) + '\n')


inputs = json.loads((dest / 'build-inputs.json').read_text())
assert all(sha(root / p) == h for p, h in inputs['source_files'].items())
shutil.copy2(root / 'LICENSE', public / 'LICENSE')
shutil.copy2(root / 'docs/release/15-desktop-closeout-20260908.md', verification / 'release-report.md')
shutil.copy2(root / 'docs/promo/desktop-closeout-20260908.md', promo / 'README.md')
source_paths = ['dist/promo/wolong-remake-trailer.mp4', 'dist/promo/wolong-remake-dosv-realmachine.mp4',
                'workplace/promo-live/original-audio/original-adlib.wav',
                'workplace/promo-live/original-audio/capture-metadata.txt',
                'workplace/desktop-polish-20260908/v18-closeout/defaults.png',
                'workplace/desktop-polish-20260908/v18-closeout/preferences-after-restart.png']
write(promo / 'sources.json', {'version': version, 'public_upload': False,
      'audio': 'DOSBox-X 原版混音擷取；僅作配樂，不作 dosgolem 對拍證據',
      'sources': {p: sha(root / p) for p in source_paths},
      'timeline_seconds': [[0, 4, '標題'], [4, 16, '歷史主預告 8–20'],
                           [16, 22, '本版預設設定'], [22, 28, '本版重啟設定'],
                           [28, 42, '歷史主預告 24–38'], [42, 50, '歷史對照片 44–52'],
                           [50, 56, '結尾與限制']]})
shutil.copy2(root / 'workplace/promo-live/original-audio/capture-metadata.txt', promo / 'audio-capture-metadata.txt')
info = json.loads((promo / 'ffprobe.json').read_text())
assert float(info['format']['duration']) == 56
v = next(s for s in info['streams'] if s['codec_type'] == 'video')
a = next(s for s in info['streams'] if s['codec_type'] == 'audio')
assert (v['width'], v['height'], v['r_frame_rate']) == (1280, 720, '30/1')
assert abs(float(v['duration']) - float(a['duration'])) < .05
log = (promo / 'media-check.log').read_text()
assert 'black_start:' not in log and 'silence_start:' not in log
mean = float(re.search(r'mean_volume: ([-\d.]+) dB', log)[1])
peak = float(re.search(r'max_volume: ([-\d.]+) dB', log)[1])
assert -35 < mean < -10 and -6 < peak < 0
write(promo / 'validation.json', {'duration_seconds': 56, 'mean_db': mean, 'peak_db': peak,
      'black_segments': 0, 'silence_segments_over_2s': 0,
      'freeze_review': '字卡與設定圖最長 6 秒；歷史場景短暫靜止，接觸表已檢視',
      'subtitle_clipping': '抽樣未見', 'human_listening': '未聲稱人耳驗收'})
notes = f'''# {version}

桌面操作修正版：滑鼠持續外推才捲圖、停手即停並保留方向鍵；快速切焦跳圖已修。
移除額外確認頁，改善選章、行軍返回與存讀檔流程。戰後結果頁預設關閉，
可選 3／5／10／15／30 秒自動返回，偏好設定可跨重啟保存。

AI 比較方向寫反造成的誤退卻已修；dosgolem 驗證呂布委任獲勝，戰術不下令重播也恢復呂布勝。
完整逐拍戰況、單挑位置與戰後退路仍有已知差異，不宣稱完全原版一致。

附件為 Linux AppImage、Windows amd64、macOS amd64／arm64 引擎包，
不含原版資料、音樂與倚天字型。請自備完整合法松崗 DOS/V 資料，依包內 README 的
`-orig` 命令啟動；繁中使用原版內建字型，不必另備倚天字型。
Windows／macOS 原生操作仍待人工驗收；Android 不在本批。

Linux AppImage 正常流程、公開包自備資料啟動、三平台封包／架構／雜湊均已查核。
自評 84／100（主觀證據評分，不是完成率）。採 RRSAL-1.0：非商業免費、
實況及影片分潤明示允許、貢獻回授、商業另洽。第三方元件依各自授權。
'''
(public / 'RELEASE-NOTES.md').write_text(notes)
(public / 'README.md').write_text(f'# {version}\n\n本目錄僅含可公開的桌面引擎包與發行文件。\n安裝方式見包內 README；素材須自行合法提供。\n本機 full/、promo/ 與 verification/ 不在上傳清單。\n')
(dest / 'README.md').write_text(f'# 臥龍傳 {version}\n\n- full/：含遊戲素材的本機完整版，不可公開散布。\n- release/：不含原版素材的公開引擎包。\n- promo/：56 秒桌面推廣片及來源／驗收，留在本機。\n- verification/：實跑與封包證據。\n\nWindows／macOS 原生操作待人工驗收；Android 排除。\nGitHub 發布狀態見專案 CONTEXT.md，不以檔案存在冒稱已發布。\n')
scripts = verification / 'scripts'
scripts.mkdir(exist_ok=True)
for name in ('package_public_desktop.py', 'verify_public_desktop.py', 'promo_desktop_render.sh',
             'finalize_desktop_release.py', 'release_desktop.sh', 'release_desktop_fs.py'):
    shutil.copy2(root / 'tools' / name, scripts / name)
for folder in (public, promo):
    (folder / 'SHA256SUMS.txt').write_text(''.join(f'{sha(p)}  {p.relative_to(folder)}\n'
        for p in sorted(folder.rglob('*')) if p.is_file() and p.name != 'SHA256SUMS.txt'))
os.environ['WOLONG_RELEASE_VERSION'] = version
sys.path.insert(0, str(root / 'tools'))
import release_desktop_fs as release
release.DEST = dest
release.manifest()
manifest = json.loads((dest / 'manifest.json').read_text())
manifest['public_release'] = json.loads((public / 'manifest.json').read_text())
manifest['promo'] = {'path': f'promo/wolong-remake-promo-{version}.mp4',
                     'sha256': sha(promo / f'wolong-remake-promo-{version}.mp4'),
                     'public_upload': False, 'duration_seconds': 56}
write(dest / 'manifest.json', manifest)
(dest / 'SHA256SUMS.txt').write_text(''.join(f'{sha(p)}  {p.relative_to(dest)}\n'
    for p in sorted(dest.rglob('*')) if p.is_file() and p.name != 'SHA256SUMS.txt'))
count = 0
for line in (dest / 'SHA256SUMS.txt').read_text().splitlines():
    expected, name = line.split('  ', 1)
    assert sha(dest / name) == expected
    count += 1
assert all(p.stat().st_uid == os.getuid() for p in dest.rglob('*'))
print(json.dumps({'version': version, 'verified_files': count, 'source_files': len(inputs['source_files']),
                  'published': False}, ensure_ascii=False))
