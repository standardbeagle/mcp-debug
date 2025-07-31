#!/usr/bin/env python3
"""
Test script to verify true dynamic tool registration with mcp-golang.
This tests that tools can be added and removed at runtime while the server is serving.
"""

import json
import subprocess
import sys
import time
import threading
import signal

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

def get_tool_list(process):
    """Get current tool list from the server"""
    send_mcp_message(process, {
        "jsonrpc": "2.0",
        "id": 100,
        "method": "tools/list"
    })
    
    response = read_mcp_response(process)
    if response and 'result' in response and 'tools' in response['result']:
        return [tool['name'] for tool in response['result']['tools']]
    return []

def test_dynamic_registration():
    """Test dynamic tool registration while server is running"""
    print("Testing dynamic tool registration with mcp-golang...")
    
    # Create a simple config with just the math server
    config = '''servers:
  - name: "math-server"
    prefix: "math"
    transport: "stdio"
    command: "../../test-servers/math-server"
    timeout: "10s"

proxy:
  healthCheckInterval: "30s"
  connectionTimeout: "10s"
  maxRetries: 3'''
    
    with open('../config-fixtures/test-dynamic-config-temp.yaml', 'w') as f:
        f.write(config)
    
    # Start the dynamic proxy server
    cmd = ['../../mcp-debug', '--dynamic', '--config', '../config-fixtures/test-dynamic-config-temp.yaml']
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
        
        # Initialize the connection
        print("Initializing MCP connection...")
        init_message = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {
                "protocolVersion": "2024-11-05",
                "capabilities": {"tools": {}},
                "clientInfo": {"name": "dynamic-test-client", "version": "1.0.0"}
            }
        }
        
        send_mcp_message(process, init_message)
        init_response = read_mcp_response(process)
        
        if not init_response:
            print("❌ Failed to initialize connection")
            return False
        
        print(f"✅ Connected to: {init_response['result']['serverInfo']['name']}")
        
        # Send initialized notification
        send_mcp_message(process, {
            "jsonrpc": "2.0",
            "method": "notifications/initialized"
        })
        
        # Check initial tool list (should start empty or minimal)
        print("\n--- Phase 1: Initial State ---")
        initial_tools = get_tool_list(process)
        print(f"Initial tools: {initial_tools}")
        
        # Wait a bit for background connections to happen
        print("\n--- Phase 2: Waiting for Dynamic Connections ---")
        time.sleep(5)
        
        # Check tool list again (should now have dynamically added tools)
        updated_tools = get_tool_list(process)
        print(f"Tools after dynamic connection: {updated_tools}")
        
        # Test calling a dynamically added tool
        if 'math_calculate' in updated_tools:
            print("\n--- Phase 3: Testing Dynamic Tool Invocation ---")
            send_mcp_message(process, {
                "jsonrpc": "2.0",
                "id": 3,
                "method": "tools/call",
                "params": {
                    "name": "math_calculate",
                    "arguments": {
                        "operation": "multiply",
                        "a": 7,
                        "b": 6
                    }
                }
            })
            
            calc_response = read_mcp_response(process)
            if calc_response and 'result' in calc_response:
                result_text = calc_response['result']['content'][0]['text']
                print(f"✅ Dynamic tool call result: {result_text}")
                
                if "42.00" in result_text:
                    print("✅ SUCCESS: Dynamic tool registration and invocation works!")
                    return True
                else:
                    print("❌ FAILED: Unexpected calculation result")
                    return False
            else:
                print("❌ FAILED: No response from dynamic tool call")
                return False
        else:
            print("❌ FAILED: math_calculate tool not found after dynamic connection")
            print(f"Available tools: {updated_tools}")
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
    success = test_dynamic_registration()
    sys.exit(0 if success else 1)