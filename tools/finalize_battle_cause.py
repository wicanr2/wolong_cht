#!/usr/bin/env python3
"""保存呂布野戰根因收據；勝方一致不代表逐拍或完整結算一致。"""
import hashlib
import json
import os
from pathlib import Path
import shutil
import sys
import tarfile

root = Path('/src')
delivery = Path(sys.argv[1])
gui = Path(sys.argv[2])
out = delivery / 'verification'
evidence = root / 'workplace/delegated-battle-20260908'
assert out.is_dir() and not list(out.iterdir()), '收據目錄必須為空'
assert delivery.stat().st_uid == out.stat().st_uid == os.getuid()
inputs = json.loads((delivery / 'build-inputs.json').read_text())
version = json.loads((delivery / 'manifest.json').read_text())['version']


def sha(path):
    with path.open('rb') as f:
        return hashlib.file_digest(f, 'sha256').hexdigest()


def write(name, data):
    (out / name).write_text(json.dumps(data, ensure_ascii=False, indent=2) + '\n')


bad = [p for p, expected in inputs['source_files'].items() if sha(root / p) != expected]
assert not bad, bad
result = json.loads((gui / 'result.json').read_text())
assert result['preferences_restart'] and result['settings_label_difference_pixels'] == 0
low, high = result['second_countdown_bounds_seconds']
assert 25 <= low <= high <= 35
assert (gui / 'runtime.sha256').read_text().split()[0] == inputs['binaries']['linux-amd64/wlgame']
assert version in (gui / 'first-start.log').read_text()
fixed = {}
for name in ('remake-fixed-no-orders', 'remake-fixed-same-input', 'prototype-branch-only'):
    value = json.loads((evidence / name / 'result.json').read_text())
    assert value['completed'] and value['battle']['AttackerWins']
    assert value['battle']['Frames'] == 1037
    fixed[name] = value
auto = json.loads((evidence / 'remake-auto.json').read_text())
assert not auto['result']['DefenderWins'] and auto['delegated'] == [True, True]
assert [(c['Men'], c['Morale']) for c in auto['after']] == [(570, 190), (484, 80)]
assert auto['rng_prefix'] == [212, 5]
original_auto = json.loads((evidence / 'original-auto/after-auto.json').read_text())
corps = {c['Index']: c for c in original_auto['corps']}
for index, remake in zip((35, 39), auto['after']):
    original = corps[index]
    assert (original['Troops'], original['Morale']) == (remake['Men'], remake['Morale'])
    assert [(u['Men'], u['Kind'] - 1) for u in original['Slots']] == [
        (u['Men'], u['Kind']) for u in remake['Units']]
assert json.loads((evidence / 'original-auto/result.json').read_text())['ax'] == 512
assert json.loads((evidence / 'original-no-orders/result.json').read_text())['completed']

shutil.copytree(gui, out / 'normal-play')
shutil.copy2(gui / 'runtime.sha256', out / 'runtime.sha256')
shutil.copy2(gui / 'first-start.log', out / 'app.log')
names = ['original-auto', 'original-no-orders', 'original-start-trace',
         'remake-no-orders', 'remake-fixed-no-orders', 'remake-fixed-same-input',
         'prototype-branch-only', 'remake-auto.json', 'delegation-fixture.json',
         'ida-script/battle-cause.json', 'targeted-tests.log', 'fixture-retest.log']
members = {}
archive = out / 'battle-cause-traces.tar.gz'
with tarfile.open(archive, 'w:gz') as t:
    for name in names:
        p = evidence / name
        for f in sorted(p.rglob('*')) if p.is_dir() else [p]:
            if f.is_file():
                rel = str(f.relative_to(evidence))
                members[rel] = sha(f)
                t.add(f, arcname=rel, recursive=False)
with tarfile.open(archive) as t:
    assert len(t.getmembers()) == len(members)
    for member in t.getmembers():
        assert hashlib.sha256(t.extractfile(member).read()).hexdigest() == members[member.name]
write('battle-trace-index.json', {'archive_sha256': sha(archive), 'members': members})
write('source-check.json', {'matched': len(inputs['source_files']), 'mismatches': bad})
identities = {
    'workplace/orig/dosv/KI.EXE': 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868',
    'workplace/ida/dosv/KI.EXE.i64': '6deb8e9c15804ef566c692771249fbb2fdf13a21ea5b899c5cec167e87e5dd62',
}
assert all(sha(root / name) == expected for name, expected in identities.items())
write('oracle-identity.json', {'dosgolem_commit': 'fd746405cf0ef9e535121e215b935b0bf59cd1de',
                             'tool': 'IDA Pro 9.4', 'address_space': 'IDA DOS/V linear',
                             'files': identities})
write('battle-summary.json', {'fixed': fixed, 'delegation_core': auto,
      'full_battle_parity': False, 'full_world_settlement_parity': False})
scripts = out / 'scripts'
scripts.mkdir()
for name in ('dosgolem_auto_battle.go', 'dosgolem_battle_replay.go',
             'auto_battle_state_replay.go', 'battle_state_replay.go', 'ida_battle_cause.py',
             'verify_desktop_polish.py', 'verify_desktop_polish.sh',
             'verify_desktop_packages.py', 'finalize_battle_cause.py',
             'release_desktop.sh', 'release_desktop_fs.py'):
    shutil.copy2(root / 'tools' / name, scripts / name)
shutil.copy2(evidence / 'check-final.log', out / 'check-final.log')
shutil.copy2(root / 'docs/playtest/115-delegated-battle-cause.md', out / 'report.md')
shutil.copy2(root / 'docs/spec/157-battle-script-comparisons.md', out / 'spec.md')
write('normal-input-check.json', {
    'version': version, 'actual_appimage': True, 'normal_play': result,
    'battle_verdict': '呂布勝；誤退卻直接原因已修，逐拍及完整世界結算仍不同',
    'battle_trajectory_pass': False,
    'windows_macos': '原生操作待人工檢驗', 'android': 'excluded',
    'full_campaign': '使用者排除',
    'limits': ['完整戰鬥採共用規則重播；正常 AppImage 另驗兩場自然遭遇',
               '先前 v.1.0.16 切焦收據保留，不當成本版重新實跑收據']})
assert all(p.stat().st_uid == os.getuid() for p in out.rglob('*'))
print(json.dumps({'version': version, 'source_files': len(inputs['source_files']),
                  'archive_members': len(members)}, ensure_ascii=False))
