#!/usr/bin/env bash
# 提交前的單一入口。三件事一起跑，任何一件不過就回非 0。
#
#   tools/check.sh          全部（Go ＋ 文件／資產）
#   tools/check.sh --docs   只跑文件／資產那半邊，跳過 go vet/test
#   tools/check.sh --go     只跑 Go 那半邊
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
#   worklist       未完成項的 verify 還成立嗎（做好了而條目沒改 → 當場開口）
#
# 第三類最貴也最晚才有檢查：**格式完全正確、連結都通、只有數字是舊的**。
# 2026-08-27 的稽核靠人工抓到六處，其中三處已經印在交付給使用者的檔案裡
# （`docs/release/09` §4）。
#
# ⚠ **自我測試要跟著跑。** 每一支的綠燈都可能是「沒問題」或「這一層根本
# 沒在比」，兩者輸出長得一樣；正對照是分辨它們的唯一方式。
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

WANT_GO=1 WANT_DOCS=1
case "${1:-}" in
    --docs) WANT_GO=0 ;;
    --go)   WANT_DOCS=0 ;;
    "")     ;;
    *) echo "用法：tools/check.sh [--docs|--go]" >&2; exit 2 ;;
esac

# 每一段都報耗時。**沒有這個就分不出「還在編譯」與「卡住了」**——
# 冷快取的 go vet 要編掉整棵依賴樹，看起來跟當掉一樣。
step() {
    local name=$1; shift
    local t0=$SECONDS
    echo "── $name ──"
    "$@"
    echo "   （$name：$((SECONDS - t0)) 秒）"
}

# vet/test 的範圍：**workplace/ 排除在外**。那是素材與實驗區，overlay 探針
# 宣告成 `package tactical` 卻放在那裡，單獨編譯必然 `undefined: Battle`；
# 不排除的話 `set -e` 會讓第一步就中止，後面五組檢查一次都跑不到。
#
# 直接列 pattern，不多開一次容器跑 `go list`（實測省 1.2 秒）。
# 下面那道迴圈擋「新增了頂層 Go 目錄卻忘了補進來」。
GO_PKGS="./cmd/... ./internal/... ./mobile/... ./tools/... ./translations/..."
for d in */; do
    d=${d%/}
    case "$d" in cmd|internal|mobile|tools|translations|workplace) continue ;; esac
    if compgen -G "$d/*.go" > /dev/null 2>&1 || compgen -G "$d/*/*.go" > /dev/null 2>&1; then
        echo "⚠ $d/ 底下有 Go 檔卻不在 GO_PKGS 裡，補進 tools/check.sh 再跑" >&2
        exit 1
    fi
done

if [[ $WANT_GO == 1 ]]; then
    step "go vet" tools/go.sh vet $GO_PKGS
    step "go test" tools/go.sh test $GO_PKGS
fi

if [[ $WANT_DOCS == 1 ]]; then
    step "文件索引" bash -c '
        tools/py.sh tools/index.py generate
        tools/py.sh tools/re_open_questions.py --strict > docs/re/43-open-questions.md'
    step "幽靈引用（指向不存在的東西）" tools/py.sh tools/phantom_scan.py
    step "過期斷言（指到的東西存在，但值不對）" bash -c '
        tools/py.sh tools/stale_scan.py --selftest
        tools/py.sh tools/stale_scan.py'
    step "未完成項（verify 還成立嗎）" bash -c '
        tools/py.sh tools/worklist.py --selftest
        tools/py.sh tools/worklist.py render --check
        tools/py.sh tools/worklist.py verify'
    step "對拍工具正對照" bash -c '
        tools/py.sh tools/parity_diff.py --selftest
        tools/py.sh tools/state_diff.py --selftest
        tools/py.sh tools/parity_save.py --selftest
        tools/py.sh tools/rng_state.py --selftest'
    step "發行目錄交換" tools/py.sh tools/release_all_fs.py --selftest
    step "資產 deny-list" bash -c '
        tools/py.sh tools/denylist.py --selftest
        tools/py.sh tools/denylist.py'
    step "TALK.DAT 校訂工具" tools/py.sh tools/talkdat_selftest.py
fi
