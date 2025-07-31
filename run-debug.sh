#!/bin/bash
# Simple wrapper to run the proxy with logging to file

LOG_FILE="${MCP_LOG_FILE:-/tmp/mcp-debug.log}"
RECORD_FILE="${MCP_RECORD_FILE:-}"

if [ -n "$RECORD_FILE" ]; then
  exec ./mcp-debug --proxy --config "$@" --log "$LOG_FILE" --record "$RECORD_FILE"
else
  exec ./mcp-debug --proxy --config "$@" --log "$LOG_FILE"
fi