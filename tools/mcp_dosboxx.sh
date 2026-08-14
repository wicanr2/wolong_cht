#!/usr/bin/env bash
# 啟動 DOSBox-X MCP Debugger 的 stdio server（給 Claude Code 的 .mcp.json 用）。
#
# 它不是給人手動跑的——stdin/stdout 是 MCP 協定本身。要確認環境好不好，
# 跑 `tools/dosboxx_bridge.sh health` 而不是這一支。
#
# ⚠ **`--network host` 是必要的**：MCP server 只是轉發層，真正的狀態在
# DOSBox-X 裡 `debug_ai.cpp` 開的 TCP 127.0.0.1:9876。沒有 host 網路，
# 它會連不上而且**看起來像 DOSBox 沒開**。
set -euo pipefail
MCP_REPO="${WOLONG_MCP_REPO:-$HOME/cht/DOSBox-X-MCP-Debugger}"
IMAGE="${WOLONG_MCP_IMAGE:-wolong-mcp-dosboxx}"
exec docker run --rm -i --log-opt max-size=10m --log-opt max-file=3 \
    --network host --memory 512m --pids-limit 64 \
    -v "$MCP_REPO:/app:ro" -w /app \
    "$IMAGE" python -m ai.server_phase5c "$@"
