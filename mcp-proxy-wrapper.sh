#!/bin/bash
# Wrapper script to redirect logs away from stdout/stderr that mcp-tui sees

# Create log directory
LOG_DIR="/tmp/mcp-proxy-logs"
mkdir -p "$LOG_DIR"

# Generate unique log file name with timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
LOG_FILE="$LOG_DIR/mcp-proxy-$TIMESTAMP.log"

# Run the server with logs redirected
exec ./mcp-server --proxy --config "$@" 2>"$LOG_FILE"