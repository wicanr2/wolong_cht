#!/usr/bin/env python3
"""由本批桌面建置的原始執行檔建立不含遊戲資料的公開包。須在 Docker 執行。"""
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import sys
import tarfile
import zipfile

version = sys.argv[1]
assert re.fullmatch(r'v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}', version)
root = Path('/src')
delivery = root / 'dist-all' / version
stage = root / 'workplace' / ('public-package-' + version)
assert not stage.exists()
assert delivery.stat().st_uid == os.getuid()
stage.mkdir()
out = delivery / 'release'
raw = root / 'workplace' / ('desktop-package-' + version) / '.work/raw'
inputs = json.loads((delivery / 'build-inputs.json').read_text())


def digest(p):
    with p.open('rb') as f:
        return hashlib.file_digest(f, 'sha256').hexdigest()


def common(p):
    shutil.copy2(root / 'LICENSE', p / 'LICENSE')
    (p / 'README.md').write_text(f'''# 臥龍傳 {version} 桌面引擎包

本包採 RRSAL-1.0，不含原版遊戲資料、音樂或倚天字型。
請自備完整合法松崗 DOS/V 版資料；`-orig` 指向包含 SINARIO.DAT、TALK.DAT、
END_S13.DAT、END_S14.DAT 等檔案的資料夾。繁中直接使用原版內建字型，不必另備倚天字型。
不要使用 PC-98 資料替代。只有原版內建字型缺失、且您有合法替代字型時，才加 `-font`。

Linux：`./wolong-remake-engine-linux-amd64-{version}.AppImage -orig /完整路徑/dosv`
Windows PowerShell：`./wlgame.exe -orig C:/Games/dosv`
macOS：依架構執行 `./darwin-arm64/wlgame` 或 `./darwin-amd64/wlgame`，加上相同 `-orig` 參數。
AppImage 若無 FUSE，可先 `--appimage-extract`，再以 `squashfs-root/AppRun` 加上述參數啟動。
音樂為可選項，`-audio` 指向自行合法轉出的 ogg 資料夾；未提供時靜音。

存檔與偏好使用作業系統的使用者設定目錄內 wolong-remake/；可用 `-save-file` 指定可寫存檔。
滑鼠持續外推捲圖，停手即停，方向鍵亦可捲圖。戰後結果頁預設關閉，可於系統選單設定秒數。
Linux AppImage 已抽驗正常流程；Windows／macOS 封包與架構已驗，原生操作待人工驗收。
已修 AI 比較方向造成的誤退卻；完整逐拍戰況與戰後退路仍有差異。Android 不在本批。

來源：https://github.com/wicanr2/wolong_cht
''')
    licenses = p / 'third-party-licenses'
    licenses.mkdir()
    for mod in Path('/gomod').glob('**/LICENSE*'):
        if mod.is_file() and 'cache/download' not in str(mod):
            rel = str(mod.relative_to('/gomod')).replace('/', '__')
            shutil.copy2(mod, licenses / rel)


def hashes(p):
    (p / 'SHA256SUMS.txt').write_text(''.join(f'{digest(f)}  {f.relative_to(p)}\n'
        for f in sorted(p.rglob('*')) if f.is_file() and f.name != 'SHA256SUMS.txt'))


for label, platforms in [('windows-amd64', ['windows-amd64']),
                         ('macos-universal', ['darwin-amd64', 'darwin-arm64'])]:
    name = f'wolong-remake-engine-{label}-{version}'
    p = stage / name
    p.mkdir()
    for platform in platforms:
        target = p if label.startswith('windows') else p / platform
        shutil.copytree(raw / platform, target, dirs_exist_ok=True)
        for binary in (raw / platform).iterdir():
            assert digest(binary) == inputs['binaries'][str(binary.relative_to(raw))]
        (target / 'translations').mkdir()
        shutil.copy2(root / 'translations/corrections.json', target / 'translations/corrections.json')
    common(p)
    hashes(p)
    if label.startswith('windows'):
        with zipfile.ZipFile(out / (name + '.zip'), 'x', zipfile.ZIP_DEFLATED) as z:
            for f in sorted(p.rglob('*')):
                if f.is_file():
                    z.write(f, str(Path(name) / f.relative_to(p)))
    else:
        with tarfile.open(out / (name + '.tar.gz'), 'x:gz') as t:
            t.add(p, arcname=name)

# 公開 AppImage 使用文字圖示，避免把原版遊戲截圖當作可散布圖示。
app = stage / 'appdir'
shutil.copytree(root / 'workplace' / ('desktop-package-' + version) / '.work/appdir', app, symlinks=True)
for name in ('gamedata', 'fonts', 'audio'):
    p = app / 'usr/share/wolong-remake' / name
    assert p.is_dir()
    shutil.rmtree(p)
for name in ('wolong-remake.png', '.DirIcon', 'SHA256SUMS.txt'):
    (app / name).unlink()
(app / 'wolong-remake.svg').write_text('<svg xmlns="http://www.w3.org/2000/svg" width="256" height="256"><rect width="256" height="256" rx="32" fill="#122638"/><text x="128" y="160" text-anchor="middle" font-family="sans-serif" font-size="92" fill="#ead9aa">WL</text></svg>')
(app / '.DirIcon').symlink_to('wolong-remake.svg')
doc = app / 'usr/share/doc/wolong-remake'
shutil.rmtree(doc)
doc.mkdir()
common(doc)
hashes(app)
print(stage)
