#!/usr/bin/env python3
"""從唯讀 MMAP.MAP 重生驗證輸入，不從原版執行結果回填。"""
import argparse
import hashlib
import importlib.util
import os
from pathlib import Path

def prepare(repo,out):
    assert out.is_dir() and out.stat().st_uid==os.getuid()
    source=repo/'workplace/orig/dosv/MMAP.MAP';raw=source.read_bytes()
    assert len(raw)==80716 and hashlib.sha256(raw).hexdigest()=='51b6fcaa390c80bd8a358dcacbf7d8dbb6dfeb0e048d8c3bdd86329df3401bcf'
    spec=importlib.util.spec_from_file_location('rle',repo/'tools/rle.py')
    rle=importlib.util.module_from_spec(spec);spec.loader.exec_module(rle)
    data=rle.decode_file(raw)
    assert len(data)==98304 and hashlib.sha256(data).hexdigest()=='740708c27a89db0a7f82865be623b5732099ec97aabf399e7dbde1450c87c861'
    folder=out/'fixtures';folder.mkdir(exist_ok=True);assert folder.stat().st_uid==os.getuid()
    target=folder/'MMAP.raw';assert not target.exists() or target.stat().st_uid==os.getuid()
    target.write_bytes(data)
    print('MMAP fixture: 98,304 bytes; input and decoded SHA-256 verified')

if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True)
    a=p.parse_args();prepare(a.repo,a.output)
