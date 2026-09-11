#!/usr/bin/env bash
# 畫面逐區對拍閘：把 `tools/parity_screens.json` 的每一組重跑一次。
#
#   tools/parity_screens.sh              # 全部
#   tools/parity_screens.sh main field   # 只跑這兩組
#
# ⭐ **原版圖每次都重新裁**，不用 `workplace/parity/` 底下既有的裁切檔——
# 存下來的 PNG 不會告訴你它是哪一版裁的（docs/lessons.json 的
# `baseline-from-stale-artifact`，犯過 3 次）。
#
# 規則層的對應物是 `tools/parity_ck.sh`。兩支都跑才算「規則與畫面都對過」。
set -uo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

OUT=${WOLONG_SCREENS_DIR:-workplace/parity/screens}
mkdir -p "$OUT"

want=("$@")
in_want() {
    [[ ${#want[@]} -eq 0 ]] && return 0
    local n
    for n in "${want[@]}"; do [[ "$n" == "$1" ]] && return 0; done
    return 1
}

fail=0
ran=0
gaps=0
while IFS=$'\t' read -r name orig crop regions args rects note; do
    in_want "$name" || continue
    ran=$((ran + 1))
    echo "── $name ── $note"
    if [[ ! -f "$orig" ]]; then
        echo "  ✗ 原版圖不在：$orig"; fail=1; continue
    fi
    ref="$OUT/orig-$name.png"
    if [[ "$crop" == "1" ]]; then
        tools/py.sh tools/parity_crop.py "$orig" "$ref" > /dev/null || {
            echo "  ✗ 裁切失敗：$orig"; fail=1; continue; }
    else
        cp "$orig" "$ref"
    fi
    # shellcheck disable=SC2086
    if ! tools/parity_shot.sh "$OUT/rm-$name.png" $args > "$OUT/$name.log" 2>&1; then
        echo "  ✗ remake 截圖失敗，見 $OUT/$name.log"; tail -3 "$OUT/$name.log" | sed 's/^/    /'
        fail=1; continue
    fi
    tools/py.sh tools/parity_diff.py "$ref" "$OUT/rm-$name.png" \
        --regions "$regions" --out "$OUT/diff-$name.png" > "$OUT/$name.md" 2>&1
    # 自訂矩形：文件當初比的是視窗本體。分區那一份照印，兩邊都要看得到。
    [[ "$rects" == "-" ]] && rects=""
    for spec in $rects; do
        tools/py.sh tools/parity_diff.py "$ref" "$OUT/rm-$name.png" --rect "${spec#*:}" 2>&1 |
            sed "0,/rect/s//${spec%%:*}/" >> "$OUT/$name.md"
    done
    tools/py.sh tools/parity_screens.py judge "$name" "$OUT/$name.md"
    case $? in 0) ;; 3) gaps=$((gaps + 1)) ;; *) fail=1 ;; esac
done < <(tools/py.sh tools/parity_screens.py list)

echo
printf '跑了 %d 組：%d 組符合預期、%d 組帶已知缺口%s。逐區數字在 %s/*.md\n' \
    "$ran" "$((ran - gaps - fail))" "$gaps" \
    "$([[ $fail -ne 0 ]] && echo '、有組超出上限')" "$OUT"
exit $fail
