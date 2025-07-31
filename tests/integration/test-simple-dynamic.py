#!/usr/bin/env python3
"""
Simple test to verify dynamic MCP server starts and can handle initialization.
"""

import json
import subprocess
import sys
import time

def test_basic_connection():
    """Test basic connection to dynamic server"""
    print("Testing basic dynamic server connection...")
    
    # Create empty config (no servers initially)
    config = '''servers: []

proxy:
  healthCheckInterval: "30s"
  connectionTimeout: "10s"
  maxRetries: 3'''
    
    with open('../config-fixtures/test-empty-config-temp.yaml', 'w') as f:
        f.write(config)
    
    # Start the dynamic proxy server
    cmd = ['../../mcp-debug', '--dynamic', '--config', '../config-fixtures/test-empty-config-temp.yaml']
    process = subprocess.Popen(
        cmd,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        bufsize=0
    )
    
    try:
        time.sleep(1)
        
        # Test basic JSON-RPC initialization
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
        
        msg_json = json.dumps(init_message)
        process.stdin.write(msg_json + '\n')
        process.stdin.flush()
        
        # Read response
        response_line = process.stdout.readline()
        if not response_line:
            print("❌ No response from server")
            return False
            
        try:
            response = json.loads(response_line.strip())
            print(f"✅ Server response: {response}")
            
            if 'result' in response and 'serverInfo' in response['result']:
                server_name = response['result']['serverInfo']['name']
                print(f"✅ Connected to: {server_name}")
                
                # Send initialized notification
                notify_message = {
                    "jsonrpc": "2.0",
                    "method": "notifications/initialized"
                }
                
                notify_json = json.dumps(notify_message)
                process.stdin.write(notify_json + '\n')
                process.stdin.flush()
                
                # Request tools list
                tools_message = {
                    "jsonrpc": "2.0",
                    "id": 2,
                    "method": "tools/list"
                }
                
                tools_json = json.dumps(tools_message)
                process.stdin.write(tools_json + '\n')
                process.stdin.flush()
                
                # Read tools response
                tools_response_line = process.stdout.readline()
                if tools_response_line:
                    tools_response = json.loads(tools_response_line.strip())
                    tools = tools_response.get('result', {}).get('tools', [])
                    print(f"✅ Tools available: {len(tools)} tools")
                    for tool in tools:
                        print(f"  - {tool['name']}: {tool['description']}")
                    
                    print("✅ SUCCESS: Dynamic MCP server is working!")
                    return True
                else:
                    print("❌ No tools response")
                    return False
            else:
                print("❌ Invalid initialization response")
                return False
                
        except json.JSONDecodeError as e:
            print(f"❌ JSON decode error: {e}")
            print(f"Raw response: {response_line}")
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
    success = test_basic_connection()
    sys.exit(0 if success else 1)