# DOSBox-X MCP Debugger 的 Python 端（stdio MCP server）
#
# 它自己不模擬任何東西：把 MCP 呼叫轉成 newline-delimited JSON，
# 送到 DOSBox-X 裡 `debug_ai.cpp` 開的 TCP 127.0.0.1:9876。
# 所以這顆容器要 `--network host` 才看得到那個 port。
FROM python:3.12-slim
RUN pip install --no-cache-dir mcp==2.0.0
WORKDIR /app
