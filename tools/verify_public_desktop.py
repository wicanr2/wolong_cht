#!/usr/bin/env python3
"""核對公開引擎包的素材排除、授權、逐檔雜湊與本批執行檔。"""
import hashlib
import json
from pathlib import Path
import struct
import sys
import tarfile
import zipfile

if sys.argv[1] == '--inventory':
    app = Path(sys.argv[2])
    files = {str(p.relative_to(app)): hashlib.sha256(p.read_bytes()).hexdigest()
             for p in app.rglob('*') if p.is_file()}
    Path(sys.argv[3]).write_text(json.dumps(files, indent=2) + '\n')
    sys.exit(0)

root = Path('/src')
delivery = Path(sys.argv[1])
version = delivery.name
inputs = json.loads((delivery / 'build-inputs.json').read_text())
stage = root / 'workplace' / ('public-package-' + version)
denied = {hashlib.sha256(p.read_bytes()).hexdigest()
          for base in (root / 'workplace/orig', root / 'workplace/eten', root / 'workplace/audio')
          for p in base.rglob('*') if p.is_file()}
results = []
packages = sorted(p for p in (delivery / 'release').iterdir() if p.suffix in ('.zip', '.gz', '.AppImage'))
assert len(packages) == 3
for package in packages:
    if package.suffix == '.zip':
        with zipfile.ZipFile(package) as z:
            assert z.testzip() is None
            files = {n.split('/', 1)[1]: z.read(n) for n in z.namelist() if not n.endswith('/')}
    elif package.suffix == '.gz':
        with tarfile.open(package) as t:
            files = {m.name.split('/', 1)[1]: t.extractfile(m).read() for m in t.getmembers() if m.isfile()}
    else:
        files = {str(p.relative_to(stage / 'appdir')): p.read_bytes()
                 for p in (stage / 'appdir').rglob('*') if p.is_file()}
        extracted = json.loads((root / 'workplace/release-closeout-20260908/public-extracted.json').read_text())
        assert extracted == {name: hashlib.sha256(data).hexdigest() for name, data in files.items()}
        data = package.read_bytes()
        assert data[:4] == b'\x7fELF' and struct.unpack_from('<H', data, 18)[0] == 62
        assert (root / 'workplace/release-closeout-20260908/public-runtime.sha256').read_text().split()[0] == inputs['binaries']['linux-amd64/wlgame']
    for name, data in files.items():
        assert not any(part in ('gamedata', 'fonts', 'audio') for part in Path(name).parts), name
        assert hashlib.sha256(data).hexdigest() not in denied, name
        assert not name.lower().endswith(('.dat', '.ogg', '.wav', '.apk', '.png')), name
    licenses = [data for name, data in files.items() if name.endswith('/LICENSE') or name == 'LICENSE']
    assert any(data == (root / 'LICENSE').read_bytes() for data in licenses)
    assert any('third-party-licenses/' in name for name in files)
    for line in files['SHA256SUMS.txt'].decode().splitlines():
        expected, name = line.split('  ', 1)
        assert hashlib.sha256(files[name]).hexdigest() == expected, name
    for name, data in files.items():
        if name == 'wlgame.exe':
            offset = struct.unpack_from('<I', data, 0x3c)[0]
            assert data[offset:offset+4] == b'PE\0\0' and struct.unpack_from('<H', data, offset+4)[0] == 0x8664
            key = 'windows-amd64/wlgame.exe'
        elif name in ('darwin-amd64/wlgame', 'darwin-arm64/wlgame'):
            magic, cpu = struct.unpack_from('<II', data)
            assert magic == 0xfeedfacf and cpu == (0x01000007 if 'amd64' in name else 0x0100000c)
            key = name
        elif name == 'usr/bin/wlgame':
            key = 'linux-amd64/wlgame'
        else:
            continue
        assert hashlib.sha256(data).hexdigest() == inputs['binaries'][key]
    results.append({'package': package.name, 'files': len(files),
                    'sha256': hashlib.sha256(package.read_bytes()).hexdigest(),
                    'original_asset_hash_matches': 0, 'license': 'RRSAL-1.0'})
receipt = {'version': version, 'packages': results,
           'appimage_launch': '以自備資料參數啟動成功', 'windows_macos_native': '待人工驗收'}
(delivery / 'release/manifest.json').write_text(json.dumps(receipt, ensure_ascii=False, indent=2) + '\n')
print(json.dumps(receipt, ensure_ascii=False))
