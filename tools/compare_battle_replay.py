#!/usr/bin/env python3
"""比較 dosgolem 與 remake 的戰況收據；完成結算不等於逐拍相同。"""
import hashlib
import json
from pathlib import Path
import sys

FIELDS = ('X', 'Y', 'Stamina', 'Order', 'NewOrder', 'Alive', 'Kind', 'Z', 'Power', 'Fatigue', 'PlaneHigh')


def units(data):
    rows = {(u['Side'], u['Squad'], u['Slot']): dict(u) for u in data['units']}
    expected = {(s, q, k) for s in range(2) for q in range(6) for k in range(8)}
    if len(data['units']) != 96 or set(rows) != expected:
        raise ValueError('必須恰好包含兩側各 48 個唯一單位槽')
    raw = bytes.fromhex(data['unit_bytes']) if 'unit_bytes' in data else None
    if raw is not None and len(raw) != 0xC00:
        raise ValueError('原版單位區應為 0xC00 byte')
    for (side, squad, slot), u in rows.items():
        if raw is not None:
            off = side*0x600+squad*0x100+slot*32
            b = raw[off:off+32]
            u.update(Alive=bool(b[0]&0x80), Kind=b[4], Z=b[10], Power=b[24], Fatigue=b[25], PlaneHigh=b[30])
        elif 'state' in u:
            s = u['state']
            u.update(Alive=s['Alive'], Kind=s['Kind'], Z=s['Z'], Power=s['Power'], Fatigue=s['Stamina'], PlaneHigh=s['PlaneHigh'])
    return rows


def compare(a, b):
    aa, bb = units(a), units(b)
    diffs = []
    count = 0
    for key in sorted(aa):
        fields = FIELDS if aa[key]['Alive'] or bb[key]['Alive'] else ('Alive',)
        for field in fields:
            count += 1
            if aa[key][field] != bb[key][field]:
                diffs.append({'unit': key, 'field': field, 'original': aa[key][field], 'remake': bb[key][field]})
    return {'original_tick': a['tick'], 'remake_tick': b['tick'],
            'same_tick': a['tick'] == b['tick'], 'compared_fields': count,
            'difference_fields': len(diffs), 'differences': diffs,
            'rng_prefix_matches': a['rng'][:4] == b['rng_prefix']}


def selftest():
    a = {'tick': 0, 'rng': '1234', 'rng_prefix': '1234', 'units': [
        dict(Side=s, Squad=q, Slot=k, **{f: (True if f=='Alive' else 0) for f in FIELDS})
        for s in range(2) for q in range(6) for k in range(8)]}
    b = json.loads(json.dumps(a))
    assert compare(a, b)['difference_fields'] == 0
    b['units'][9]['Y'] = 1
    assert compare(a, b)['difference_fields'] == 1
    b['units'].pop()
    try:
        compare(a, b)
    except ValueError:
        return
    raise AssertionError('缺槽沒有被拒絕')


def main():
    if sys.argv[1:] == ['--selftest']:
        selftest()
        print('戰況比較器正負對照通過')
        return
    original, remake, out = map(Path, sys.argv[1:])
    read = lambda p: json.loads(p.read_text())
    command = read(remake/'command.json')
    rows = []
    for p in sorted(original.glob('tick-*.json')):
        q = remake / p.name
        if p.name == f"tick-{command['actual_tick']:04d}.json":
            q = remake / 'command-state.json'
        if q.exists():
            result = compare(read(p), read(q))
            result['snapshot'] = p.name
            rows.append(result)
    if not rows:
        raise ValueError('沒有共同節拍，不能宣稱逐拍比較')
    init = compare(read(original/'initialized.json'), read(remake/'initialized.json'))
    result = {
        'classification': '同一初始化狀態；逐拍差異與完成結算分開判定，非正常 GUI 通過收據',
        'original': str(original), 'remake': str(remake),
        'hp_label': 'Stamina 沿用 dosgolem 匯出名稱，實際比較原記錄 +03 體力，非 +19 餘力',
        'initialized': init, 'timeline': rows,
        'first_divergence': next((r['original_tick'] for r in rows if r['difference_fields'] or not r['rng_prefix_matches']), None),
        'completed_both': read(original/'result.json')['completed'] and read(remake/'result.json')['completed'],
        'ending_ticks': {'original': read(original/'before-settlement.json')['tick'], 'remake': read(remake/'before-settlement.json')['tick']},
        'command': command,
        'input_sha256': {str(p): hashlib.sha256(p.read_bytes()).hexdigest() for p in (original/'entry-SAVE.DAT',original/'spawn-rng.bin')},
    }
    result['trajectory_pass'] = result['first_divergence'] is None and result['ending_ticks']['original'] == result['ending_ticks']['remake']
    out.write_text(json.dumps(result, ensure_ascii=False, indent=2))
    print(json.dumps({k: result[k] for k in ('first_divergence','completed_both','ending_ticks','trajectory_pass')},ensure_ascii=False))


if __name__ == '__main__':
    main()
