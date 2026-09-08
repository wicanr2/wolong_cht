#!/usr/bin/env python3
"""彙整已通過的桌面打磨收據；之後仍需封包查核與 manifest／SHA-256 重生。"""
import hashlib
import json
from pathlib import Path
import shutil
import sys

delivery, evidence = map(Path, sys.argv[1:3])
verification = delivery / 'verification'
assert verification.is_dir() and not list(verification.iterdir()), '只能彙整到尚未寫入收據的候選包'
inputs = json.loads((delivery / 'build-inputs.json').read_text())
version = json.loads((delivery / 'manifest.json').read_text())['version']


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


cases = {}
for source, name in [('final-polish', 'preferences-and-battles'), ('final-scroll-2', 'scroll'), ('final-march', 'march')]:
    case = evidence / source
    result = json.loads((case / 'result.json').read_text())
    assert (case / 'runtime.sha256').read_text().split()[0] == inputs['binaries']['linux-amd64/wlgame']
    cases[name] = result
    shutil.copytree(case, verification / name)
assert cases['preferences-and-battles']['preferences_restart']
assert cases['scroll']['passed'] and all(c['difference_pixels'] == 0 for c in cases['scroll']['comparisons'])
assert cases['march']['stop_difference_pixels'] == cases['march']['target_menu_difference_pixels'] == 0
for name in ('original-scroll', 'original-scroll-y', 'original-system'):
    shutil.copytree(evidence / name, verification / name)
for name in ('current-formation-diff.json', 'load-diff-corrected.json', 'formation-diff-corrected.json'):
    shutil.copy2(evidence / name, verification / name)
shutil.copy2(evidence / 'final-polish/runtime.sha256', verification / 'runtime.sha256')
shutil.copy2(evidence / 'final-polish/first-start.log', verification / 'app.log')
assert version in (verification / 'app.log').read_text()

source_root = Path('/src')
bad = [name for name, expected in inputs['source_files'].items() if sha(source_root / name) != expected]
assert not bad, bad
(verification / 'source-check.json').write_text(json.dumps({'matched': len(inputs['source_files']), 'mismatches': bad}, indent=2))
identity = {'dosgolem_commit': 'fd746405cf0ef9e535121e215b935b0bf59cd1de', 'files': {}}
for name in ('workplace/dosgolem-main-20260908/dosgolem-shot', 'workplace/orig/dosv/KI.EXE',
             'workplace/dosgolem/root-noclouds/SAVE.DAT', 'workplace/dosgolem/root-saveb/SAVE.DAT'):
    identity['files'][name] = sha(source_root / name)
(verification / 'oracle-identity.json').write_text(json.dumps(identity, indent=2))
scripts = verification / 'scripts'
scripts.mkdir()
for name in ('verify_desktop_polish.py', 'verify_desktop_polish.sh', 'verify_desktop_scroll.py',
             'verify_desktop_march.py', 'verify_visual_residuals.py', 'verify_scroll_steps.py',
             'finalize_desktop_polish.py', 'verify_desktop_packages.py', 'release_desktop.sh', 'release_desktop_fs.py'):
    shutil.copy2(source_root / 'tools' / name, scripts / name)
result = {'version': version, 'actual_appimage': True, 'cases': cases,
          'windows_macos': '使用者指定人工檢驗，尚未回報原生操作結果', 'android': 'excluded',
          'limits': ['未證實完整同狀態戰況或全劇本通關',
                     '低更新率下半秒快速切焦曾跳圖；兩秒停留驗收不涵蓋快速切焦',
                     '戰術按住兩秒供低更新率採樣，不代表實體短按延遲']}
(verification / 'normal-input-check.json').write_text(json.dumps(result, ensure_ascii=False, indent=2))
print(json.dumps({'version': version, 'source_files': len(inputs['source_files']), 'cases': list(cases)}, ensure_ascii=False))
