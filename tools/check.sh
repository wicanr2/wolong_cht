#!/usr/bin/env bash
# 提交前的單一入口。三件事一起跑，任何一件不過就回非 0。
#
#   tools/check.sh
#
# 為什麼要有這支：檢查分屬多個工具，分開記就會有一個被忘記。
# 實際被忘記過的是文件那一組——狀態行與內文矛盾了好幾輪都沒人發現
# （`docs/formats/07` 寫「像素格式未解」時，同一份文件已經解完了），
# 以及規則指向不存在的工具與目錄（`tools/addr.py`、`internal/game/`），
# 那種待辦永遠不會完成卻一直佔著缺口欄。
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

echo "── go vet ──"
tools/go.sh vet ./...
echo "── go test ──"
tools/go.sh test ./...
echo "── 文件索引 ──"
tools/py.sh tools/index.py generate
tools/py.sh tools/re_open_questions.py > docs/re/43-open-questions.md
echo "── 幽靈引用 ──"
tools/py.sh tools/phantom_scan.py
echo "── 資產 deny-list ──"
tools/py.sh tools/denylist.py --selftest
tools/py.sh tools/denylist.py
echo "── TALK.DAT 校訂工具 ──"
tools/py.sh tools/talkdat_selftest.py
