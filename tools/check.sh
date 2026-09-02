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
#
# ⭐ **文件的錯有三類，三支工具各擋一類**：
#
#   index.py       狀態行與內文矛盾（說未解、內文全是 confirmed）
#   phantom_scan   指向不存在的東西（檔案、Go 識別字、IDA 符號、測試名）
#   stale_scan     指到的東西存在，但值不對（雜湊、映像標籤、旗標、覆蓋率、未解列數、規格份數）
#
# 第三類最貴也最晚才有檢查：**格式完全正確、連結都通、只有數字是舊的**。
# 2026-08-27 的稽核靠人工抓到六處，其中三處已經印在交付給使用者的檔案裡
# （`docs/release/09` §4）。
#
# ⚠ **自我測試要跟著跑。** 每一支的綠燈都可能是「沒問題」或「這一層根本
# 沒在比」，兩者輸出長得一樣；正對照是分辨它們的唯一方式。
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

echo "── go vet ──"
tools/go.sh vet ./...
echo "── go test ──"
tools/go.sh test ./...
echo "── 文件索引 ──"
tools/py.sh tools/index.py generate
tools/py.sh tools/re_open_questions.py --strict > docs/re/43-open-questions.md
echo "── 幽靈引用（指向不存在的東西）──"
tools/py.sh tools/phantom_scan.py
echo "── 過期斷言（指到的東西存在，但值不對）──"
tools/py.sh tools/stale_scan.py --selftest
tools/py.sh tools/stale_scan.py
echo "── 對拍工具正對照 ──"
tools/py.sh tools/parity_diff.py --selftest
echo "── 發行目錄交換 ──"
tools/py.sh tools/release_all_fs.py --selftest
echo "── 資產 deny-list ──"
tools/py.sh tools/denylist.py --selftest
tools/py.sh tools/denylist.py
echo "── TALK.DAT 校訂工具 ──"
tools/py.sh tools/talkdat_selftest.py
