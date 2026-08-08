#!/usr/bin/env bash
# 提交前的單一入口。三件事一起跑，任何一件不過就回非 0。
#
#   tools/check.sh
#
# 為什麼要有這支：三種檢查分屬三個工具（go vet / go test / index.py），
# 分開記就會有一個被忘記。實際被忘記的是第三個——
# 文件的狀態行與內文矛盾了好幾輪都沒人發現
# （`docs/formats/07` 寫「像素格式未解」時，同一份文件已經解完了）。
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

echo "── go vet ──"
tools/go.sh vet ./...
echo "── go test ──"
tools/go.sh test ./...
echo "── 文件索引 ──"
python3 tools/index.py generate
echo "── 資產 deny-list ──"
python3 tools/denylist.py --selftest
python3 tools/denylist.py
