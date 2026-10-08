#!/usr/bin/env python3
"""核對矩形閉包、全 word 計量、正確資產 bank、來源與錯版收據。"""
import argparse
import hashlib
import json
import re
from pathlib import Path

EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
ICON='2154782c045b898aa5fafa74a4ff0c3745771ec85799b7882e1c4c009b3f1c1d'
NEW={0xf020:288,0xf140:59,0xf17b:39,0xf1a3:203,0xaaa:47,0xad9:109,0xcac:23,0xcc3:27,
     0xc61f:52,0xc6bf:55,0xc6ae:17,0xc6f6:86,0xc74c:41,0xc775:25,0xc78e:27,0xbcd:71}
GROUPS={'rectangle':3072,'edge':108,'mode':16,'gauge':131072,'bar':160,'span':192,
        'selection':184,'window-fill':18,'sidebar':75,'waiting':5,'actual-command':54,'actual-buttons':2}
BACKLINKS={
 'docs/re/48-window-display-list.md':('後續矩形證據見','106-c-rectangle-bars-restoration.md'),
 'docs/re/60-tactical-sidebar.md':('後續 C 計量證據見','106-c-rectangle-bars-restoration.md'),
 'docs/re/105-c-aligned-blit-restoration.md':('後續正確素材 bank 證據見','106-c-rectangle-bars-restoration.md'),
 'docs/spec/31-tactical-sidebar.md':('原始位移與 byte 長度','226-sidebar-bar-widths.md'),
}

def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()

def verify(repo,out):
 raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(raw).hexdigest()==EXPECTED
 icons=(repo/'workplace/orig/dosv/ICONGRF.DAT').read_bytes();assert len(icons)==47776 and hashlib.sha256(icons).hexdigest()==ICON
 assert len(icons[0x9700:0x9dc0])==9*192 and len(icons[0x9dc0:0x9f40])==0x180
 probe=json.loads((out/'ida/ida-probe.json').read_text());assert probe['tool_version']=='9.4' and probe['input_sha256']==probe['ida_input_sha256']==EXPECTED
 assert probe['function_count']==739 and len(probe['targets'])==36 and probe['target_relocations_applied']==4
 routines={};new={}
 for t in probe['targets']:
  b=b''.join(bytes.fromhex(c['file_bytes']) for c in t['chunks'])
  assert b==b''.join(raw[c['start']-0x10000+512:c['end']-0x10000+512] for c in t['chunks'])
  routines[t['name']]=hashlib.sha256(b).hexdigest();assert routines[t['name']]==t['file_sha256']
  off=t['ida_linear']-0x10000
  if off in NEW:assert len(b)==NEW[off];new[t['name']]=routines[t['name']]
 assert len(new)==16 and sum(NEW.values())==1169
 receipts={};full={};total=sum(GROUPS.values())
 for level in ['O0','O2']:
  f=out/'results'/f'{level}.json';r=json.loads(f.read_text())
  assert r['schema']=='wolong-c-rect-parity-v1' and r['input_sha256']==EXPECTED and r['icon_sha256']==ICON
  assert all(r['routine_sha256'][n]==v for n,v in routines.items() if n!='sub_100DF')
  assert r['passed'] and r['mismatch'] is None and r['cases']==total and r['groups']==GROUPS
  assert r['full_ram_plane_audits']==total and r['go_bar_cases']==131072 and r['indexed_content_audits']==2441
  assert r['content_top']==40 and r['content_height']==400 and r['trace_entry_size']==90
  assert all(r['entries_seen'].get(n,0)>0 for n in new) and r['entries_seen']['fixture-far']>0
  assert r['asset_banks']=='ICONGRF segment1 at 3000; whole segment3 at 2200; command 2200, border 226c'
  assert r['original_state_sha256']==r['c_state_sha256'] and r['c_machine_code_match'] is False
  receipts[level]=sha(f);full[level]=r
 assert receipts['O0']==receipts['O2']
 controls={}
 for n,g in [(1,'rectangle'),(2,'rectangle'),(3,'rectangle'),(4,'rectangle'),(5,'rectangle'),(6,'rectangle'),(7,'gauge'),(8,'gauge'),(9,'waiting'),(10,'bar'),(11,'rectangle'),(12,'selection')]:
  f=out/'results'/f'mutant-{n}.json';r=json.loads(f.read_text());assert not r['passed'] and r['input_sha256']==EXPECTED
  assert 0<r['cases']<=GROUPS[g] and r['groups']=={g:r['cases']};m=r['mismatch'];assert m['group']==g and m['case']==r['cases']-1
  assert any(m['original'+s]!=m['c'+s] for s in ['', '_bar','_trace','_ports','_device','_planes','_ram'])
  controls[str(n)]={'group':g,'first_rejected':r['cases'],'receipt_sha256':sha(f)}
 source={}
 for line in (out/'results/c-source.sha256').read_text().splitlines():
  digest,name=line.split(None,1);relative=name.strip().removeprefix('/repo/');assert sha(repo/relative)==digest;source[relative]=digest
 for p in ['tools/c_recovery/rect.c','tools/c_recovery/rect.h','tools/c_recovery/rect_fixture.h','tools/c_recovery_rect.go','cmd/wlgame/battlelayout.go']:assert p in source
 compiled=sha(out/'results/c-source.sha256');assert (out/'results/compiled-source-digest.txt').read_text().strip()==compiled
 for level in ['O0','O2']:assert '-DKI_RECT_SOURCE_DIGEST=0x'+compiled[:16] in (out/'results'/f'buildinfo-{level}.txt').read_text()
 helper=(out/'results/go-helper-source.txt').read_text();production=(repo/'cmd/wlgame/battlelayout.go').read_text()
 match=re.search(r'^func battleSideBarLengths\([^\n]*\n.*?^}\n',production,re.M|re.S);assert match and helper==match[0]
 constant=(out/'results/go-helper-constant.txt').read_text().strip();line=next(x for x in production.splitlines() if re.match(r'\s*battleSideBarMaxLen\s*=',x));assert constant=='const '+line
 before=(out/'results/golem-source-before.sha256').read_bytes();assert before==(out/'results/golem-source-after.sha256').read_bytes()
 for path,(marker,link) in BACKLINKS.items():s=(repo/path).read_text();assert marker in s and link in s
 assert 'ok  ' in (out/'go-ui-tests.log').read_text()
 assert 'ok  ' in (out/'go-ui-full-tests.log').read_text()
 supplement=json.loads((repo/'docs/re/rectangle-handler-code.json').read_text());assert supplement['input_sha256']==EXPECTED and len(supplement['instructions'])==8
 assert sum(r['file_end']-r['file_start'] for r in supplement['instructions'])==20
 for r in supplement['instructions']:assert raw[r['file_start']:r['file_end']]==bytes.fromhex(r['bytes']) and not r['ida_code_classified']
 result={'schema':'wolong-c-rect-verification-v1','status':'semantic-conformed','input_sha256':EXPECTED,'icon_sha256':ICON,
         'new_routine_sha256':new,'dependency_routine_sha256':routines,'ida_database_sha256':sha(out/'ida/input.exe.i64'),
         'source_sha256':source,'compiled_source_manifest_sha256':compiled,'go_helper_source_sha256':sha(out/'results/go-helper-source.txt'),
         'go_helper_constant_sha256':sha(out/'results/go-helper-constant.txt'),'cases_per_optimization':total,'groups':GROUPS,
         'go_bar_cases_per_optimization':131072,'full_ram_plane_audits_per_optimization':total,'indexed_content_audits_per_optimization':2441,
         'receipt_sha256':receipts,'negative_controls_rejected':12,'mutants':controls,'entries_seen':full['O2']['entries_seen'],
         'oracle_revision':(out/'results/golem-revision.txt').read_text().strip(),'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
         'ui_cold_test_log_sha256':sha(out/'go-ui-tests.log'),'ui_full_cold_test_log_sha256':sha(out/'go-ui-full-tests.log'),'scope_backlinks':BACKLINKS,'handler_code_supplement_sha256':sha(repo/'docs/re/rectangle-handler-code.json'),
         'platform_contract':full['O2']['platform'],'asset_banks':full['O2']['asset_banks'],'c_machine_code_match':False}
 (out/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
 print(f'16 new C routines; {total} complete RAM/plane cases; 131072 original/C/Go bars; 12 mutations; source identity: PASS')

if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True)
 args=p.parse_args();verify(args.repo,args.output)
