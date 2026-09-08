#!/usr/bin/env python3
"""彙整快速切焦與完整戰況驗證；不將完成重播誤標成原版一致。"""
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import sys
import tarfile

root = Path('/src')
delivery = Path(sys.argv[1])
verification = delivery / 'verification'
assert verification.is_dir() and not list(verification.iterdir()), '收據目錄必須為空'
assert delivery.stat().st_uid == verification.stat().st_uid == os.getuid()
inputs = json.loads((delivery / 'build-inputs.json').read_text())
version = json.loads((delivery / 'manifest.json').read_text())['version']
evidence = root / 'workplace/desktop-polish-20260908'
battles = root / 'workplace/same-battle-20260908'


def sha(path):
    with path.open('rb') as f:
        return hashlib.file_digest(f, 'sha256').hexdigest()


def write(name, value):
    (verification / name).write_text(json.dumps(value, ensure_ascii=False, indent=2) + '\n')


bad = [name for name, expected in inputs['source_files'].items() if sha(root / name) != expected]
assert not bad, bad
write('source-check.json', {'matched': len(inputs['source_files']), 'mismatches': bad})
cases = {}
for source, name in [('v16-focus-half-clean-2', 'focus-half'),
                     ('v16-focus-frame-gap-clean', 'focus-frame-gap'),
                     ('v16-normal-battles', 'normal-play')]:
    case = evidence / source
    result = json.loads((case / 'result.json').read_text())
    assert (case / 'runtime.sha256').read_text().split()[0] == inputs['binaries']['linux-amd64/wlgame']
    if name.startswith('focus-'):
        assert result['passed'] and all(c['difference_pixels'] == 0 for c in result['comparisons'])
    else:
        assert result['preferences_restart'] and result['settings_label_difference_pixels'] == 0
        low, high = result['second_countdown_bounds_seconds']
        assert 25 <= low <= high <= 35
    cases[name] = result
    shutil.copytree(case, verification / name)
for name in ('original-scroll', 'original-scroll-y'):
    shutil.copytree(evidence / name, verification / name)
shutil.copy2(evidence / 'v16-normal-battles/runtime.sha256', verification / 'runtime.sha256')
shutil.copy2(evidence / 'v16-normal-battles/first-start.log', verification / 'app.log')
assert version in (verification / 'app.log').read_text()

# 完整逐拍原記錄與重播壓成單一私有證據包，並逐成員驗證內容未變。
names = ['original-siege-continuous', 'remake-siege-same-input',
         'original-field-combat', 'remake-field-same-input', 'appimage-v16',
         'siege-same-input-comparison.json', 'field-same-input-comparison.json']
members = {}
archive = verification / 'battle-traces.tar.gz'
with tarfile.open(archive, 'w:gz') as t:
    for name in names:
        p = battles / name
        for f in sorted(p.rglob('*')) if p.is_dir() else [p]:
            if f.is_file():
                rel = str(f.relative_to(battles))
                members[rel] = sha(f)
                t.add(f, arcname=rel, recursive=False)
with tarfile.open(archive) as t:
    assert len(t.getmembers()) == len(members)
    for m in t.getmembers():
        assert hashlib.sha256(t.extractfile(m).read()).hexdigest() == members[m.name]
write('battle-trace-index.json', {'archive_sha256': sha(archive), 'members': members})

comparisons = {}
initial = {}
for name, orig in [('siege', 'original-siege-continuous'), ('field', 'original-field-combat')]:
    data = json.loads((battles / f'{name}-same-input-comparison.json').read_text())
    assert data['completed_both'] and data['initialized']['difference_fields'] == 0
    assert data['initialized']['compared_fields'] == 1056
    assert data['initialized']['rng_prefix_matches']
    assert not data['trajectory_pass'] and data['first_divergence'] == 1
    comparisons[name] = {k: v for k, v in data.items() if k != 'timeline'}
    shutil.copy2(battles / f'{name}-same-input-comparison.json', verification / f'{name}-comparison.json')
    expected = {(1-u['Side'], u['Squad'], u['Slot']): (u['X'], u['Y'], u['Stamina'])
                for u in json.loads((battles / orig / 'initialized.json').read_text())['units']}
    observed = {}
    log = (battles / 'appimage-v16' / f'{name}.log').read_text()
    assert version in log
    for m in re.finditer(r'兵 (\d)/(\d)/(\d) \(\s*(-?\d+),\s*(-?\d+)\) z-?\d+ 體(\d+)', log):
        side, squad, slot, x, y, hp = map(int, m.groups())
        key = (side, squad, slot)
        assert key not in observed
        observed[key] = (x, y, hp)
    assert len(expected) == len(observed) == 96 and observed == expected
    initial[name] = {'slots': 96, 'coordinate_stamina_differences': 0,
                     'image_same_frame': '未證實；首次 Draw 圖僅供呈現參考'}
assert (battles / 'appimage-v16/runtime.sha256').read_text().split()[0] == inputs['binaries']['linux-amd64/wlgame']
write('appimage-initial-check.json', initial)
write('battle-summary.json', comparisons)
write('oracle-identity.json', {
    'dosgolem_commit': 'fd746405cf0ef9e535121e215b935b0bf59cd1de',
    'address_space': 'IDA DOS/V linear',
    'files': {str(p.relative_to(root)): sha(p) for p in [
        root / 'workplace/orig/dosv/KI.EXE',
        root / 'workplace/dosgolem/root-noclouds/SAVE.DAT',
        root / 'workplace/dosgolem/root-saveb/SAVE.DAT',
        root / 'workplace/parity/SAVE-FIELD.DAT',
        battles / 'oracle-replay-combat']}})
scripts = verification / 'scripts'
scripts.mkdir()
for name in ('verify_desktop_polish.py', 'verify_desktop_polish.sh', 'verify_desktop_scroll.py',
             'dosgolem_battle_replay.go', 'battle_state_replay.go', 'compare_battle_replay.py',
             'finalize_focus_battle.py', 'verify_desktop_packages.py',
             'release_desktop.sh', 'release_desktop_fs.py'):
    shutil.copy2(root / 'tools' / name, scripts / name)
shutil.copy2(battles / 'check-complete.log', verification / 'check-complete.log')
shutil.copy2(battles / 'check-final.log', verification / 'check-final.log')
shutil.copy2(root / 'docs/playtest/114-focus-and-same-battle.md', verification / 'report.md')
write('normal-input-check.json', {
    'version': version, 'actual_appimage': True, 'cases': cases,
    'battle_trajectory_pass': False,
    'battle_verdict': '同一初始已映射欄位及亂數、同拍指令下，第 1 拍分歧，野戰勝負相反',
    'windows_macos': '使用者指定人工檢驗，原生操作尚未回報',
    'android': 'excluded', 'full_campaign': '使用者排除',
    'limits': ['完整戰鬥為共用規則重播，正常 GUI 另抽驗兩場自然遭遇',
               '未以 Linux 容器宣稱實機短按延遲、音訊或所有焦點時序一致']})
assert all(p.stat().st_uid == os.getuid() for p in verification.rglob('*'))
print(json.dumps({'version': version, 'source_files': len(inputs['source_files']),
                  'archive_members': len(members), 'battle_trajectory_pass': False}, ensure_ascii=False))
