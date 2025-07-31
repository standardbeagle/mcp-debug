#!/usr/bin/env python3
"""
Test script to verify tool invocation works after server updates
"""

import json
import subprocess
import sys
import time

def send_mcp_message(process, message):
    """Send a JSON-RPC message to the MCP server"""
    msg_json = json.dumps(message)
    process.stdin.write(msg_json + '\n')
    process.stdin.flush()

def read_mcp_response(process):
    """Read a JSON-RPC response from the MCP server"""
    try:
        line = process.stdout.readline()
        if line:
            return json.loads(line.strip())
        return None
    except json.JSONDecodeError as e:
        print(f"JSON decode error: {e}")
        return None

def test_updated_server_tools():
    """Test tool invocation with updated server (v2)"""
    print("Testing tool invocation with updated server...")
    
    # Create config pointing to v2 server
    config = '''servers:
  - name: "lifecycle-server"
    prefix: "test"
    transport: "stdio"
    command: "./test-servers/lifecycle-server-v2"
    timeout: "10s"

proxy:
  healthCheckInterval: "30s"
  connectionTimeout: "10s"
  maxRetries: 3'''
    
    with open('test-updated-config.yaml', 'w') as f:
        f.write(config)
    
    # Start proxy server
    cmd = ['./mcp-server', '--proxy', '--config', 'test-updated-config.yaml']
    process = subprocess.Popen(
        cmd,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        bufsize=0
    )
    
    try:
        time.sleep(2)
        
        # Initialize
        init_message = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {
                "protocolVersion": "2024-11-05",
                "capabilities": {"tools": {}},
                "clientInfo": {"name": "test-client", "version": "1.0.0"}
            }
        }
        
        send_mcp_message(process, init_message)
        init_response = read_mcp_response(process)
        print(f"Initialized: {init_response['result']['serverInfo']['name']}")
        
        # Send initialized notification
        send_mcp_message(process, {"jsonrpc": "2.0", "method": "notifications/initialized"})
        
        # List tools
        send_mcp_message(process, {"jsonrpc": "2.0", "id": 2, "method": "tools/list"})
        tools_response = read_mcp_response(process)
        
        tools = tools_response['result']['tools']
        print(f"Available tools: {[t['name'] for t in tools]}")
        
        # Test updated hello tool (should show v2 response)
        print("\nTesting updated hello tool...")
        send_mcp_message(process, {
            "jsonrpc": "2.0",
            "id": 3,
            "method": "tools/call",
            "params": {
                "name": "test_hello",
                "arguments": {"name": "Updated Server Test"}
            }
        })
        
        hello_response = read_mcp_response(process)
        hello_text = hello_response['result']['content'][0]['text']
        print(f"Hello response: {hello_text}")
        
        # Test new timestamp tool (only available in v2)
        print("\nTesting new timestamp tool...")
        send_mcp_message(process, {
            "jsonrpc": "2.0",
            "id": 4,
            "method": "tools/call",
            "params": {
                "name": "test_timestamp",
                "arguments": {}
            }
        })
        
        timestamp_response = read_mcp_response(process)
        timestamp_text = timestamp_response['result']['content'][0]['text']
        print(f"Timestamp response: {timestamp_text}")
        
        # Verify responses indicate v2
        if "v2" in hello_text and "Current time:" in timestamp_text:
            print("\n✅ SUCCESS: Tool invocation works with updated server!")
            print("✅ Updated tool functionality verified")
            print("✅ New tools from server update work correctly")
            return True
        else:
            print("\n❌ FAILED: Server update not reflected in tool responses")
            return False
            
    except Exception as e:
        print(f"Test error: {e}")
        return False
        
    finally:
        try:
            process.terminate()
            process.wait(timeout=5)
        except subprocess.TimeoutExpired:
            process.kill()
            process.wait()

if __name__ == "__main__":
    success = test_updated_server_tools()
    sys.exit(0 if success else 1)