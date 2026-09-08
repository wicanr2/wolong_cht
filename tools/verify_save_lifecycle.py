#!/usr/bin/env python3
"""於 Docker 內驗證正常操作產生的存檔；不改寫任何遊戲資料。"""
import hashlib
import json
from pathlib import Path
import sys

seed, original, app, receipt = map(Path, sys.argv[1:5])
base = seed.read_bytes()
assert len(base) == 22208 * 4
results = {}
for name, path in [("dosgolem", original), ("appimage", app / "SAVE.DAT")]:
    data = path.read_bytes()
    assert len(data) == len(base), f"{name} 存檔長度不符"
    differences = []
    for slot in range(4):
        start, end = slot * 22208, (slot + 1) * 22208
        count = sum(a != b for a, b in zip(base[start:end], data[start:end]))
        differences.append(count)
        assert (count > 0) == (slot == 1), f"{name} 槽 {slot + 1} 寫入範圍不符"
    title = data[22208 + 64:22208 + 128].split(b"\0")[0]
    assert title == base[64:128].split(b"\0")[0], f"{name} 第二槽標題未保存第一槽勢力"
    results[name] = {"sha256": hashlib.sha256(data).hexdigest(),
                     "changed_bytes_by_slot": differences,
                     "saved_title": title.decode("big5")}

before = (app / "save/save.sha256").read_text().split()[0]
after = (app / "reload/save.sha256").read_text().split()[0]
assert before == after == results["appimage"]["sha256"], "重啟讀檔改動了 SAVE.DAT"
native = app / "SAVE-slot2.wlsave"
assert native.is_file() and native.stat().st_size > 0, "缺少第二槽原生存檔"
results["native_save_sha256"] = hashlib.sha256(native.read_bytes()).hexdigest()
results["restart_save_unchanged"] = True
results["limits"] = "檔案檢查不替代畫面與輸入收據；未對齊隱藏亂數與時鐘，不要求兩端存檔完全相同。"
receipt.write_text(json.dumps(results, ensure_ascii=False, indent=2) + "\n")
print(json.dumps(results, ensure_ascii=False, indent=2))
