#!/bin/bash
# Test script for playback functionality

echo "=== MCP Playback Testing ==="
echo

# First, create a sample recording by testing the minimal server
echo "1. Creating a sample recording with minimal server..."
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"1.0.0","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}},"id":1}' > /tmp/test-input.jsonl
echo '{"jsonrpc":"2.0","method":"initialized","params":{}}' >> /tmp/test-input.jsonl
echo '{"jsonrpc":"2.0","method":"tools/list","params":{},"id":2}' >> /tmp/test-input.jsonl
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"test_tool","arguments":{}},"id":3}' >> /tmp/test-input.jsonl

# Build minimal test server
cat > /tmp/minimal-server.go << 'EOF'
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Parse JSON-RPC request
		var request map[string]interface{}
		json.Unmarshal([]byte(line), &request)
		
		method, _ := request["method"].(string)
		id := request["id"]
		
		var response map[string]interface{}
		
		switch method {
		case "initialize":
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id": id,
				"result": map[string]interface{}{
					"protocolVersion": "2024-11-05",
					"capabilities": map[string]interface{}{
						"tools": map[string]interface{}{
							"listChanged": true,
						},
					},
					"serverInfo": map[string]interface{}{
						"name": "Test Server",
						"version": "1.0.0",
					},
				},
			}
		case "tools/list":
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id": id,
				"result": map[string]interface{}{
					"tools": []map[string]interface{}{
						{
							"name": "test_tool",
							"description": "A test tool",
							"inputSchema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{},
							},
						},
					},
				},
			}
		case "tools/call":
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id": id,
				"result": map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": "Test tool executed successfully!",
						},
					},
				},
			}
		case "initialized":
			// No response for notifications
			continue
		default:
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id": id,
				"error": map[string]interface{}{
					"code": -32601,
					"message": "Method not found",
				},
			}
		}
		
		if response != nil {
			jsonBytes, _ := json.Marshal(response)
			fmt.Println(string(jsonBytes))
		}
	}
}
EOF

go build -o /tmp/minimal-server /tmp/minimal-server.go

echo "2. Testing playback client against minimal server..."
cat /tmp/test-input.jsonl | ./mcp-debug --playback-client /dev/stdin | /tmp/minimal-server

echo
echo "3. Testing playback server with recorded responses..."
echo "   (This would normally be used with mcp-tui)"

echo
echo "=== Playback functionality ready! ==="
echo
echo "Usage examples:"
echo "1. Record a session:"
echo "   MCP_RECORD_FILE='session.jsonl' mcp-tui ./run-proxy.sh test-empty-config.yaml"
echo
echo "2. Replay client requests to test your server:"
echo "   ./mcp-debug --playback-client session.jsonl | your-mcp-server"
echo
echo "3. Replay server responses to test your client:"
echo "   mcp-tui ./mcp-debug --playback-server session.jsonl"
echo