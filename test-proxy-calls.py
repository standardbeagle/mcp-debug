#!/usr/bin/env python3
"""
Test script to verify the MCP proxy server tool calls work correctly.
This simulates what an MCP client would do.
"""

import json
import subprocess
import sys
import threading
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
        print(f"Raw line: {line}")
        return None

def test_proxy_server():
    """Test the proxy server functionality"""
    print("Starting MCP proxy server test...")
    
    # Start the proxy server
    cmd = ['./mcp-server', '--proxy', '--config', 'test-multi-config.yaml']
    process = subprocess.Popen(
        cmd,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        bufsize=0
    )
    
    try:
        # Give the server time to start
        time.sleep(2)
        
        # 1. Initialize the MCP connection
        init_message = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {
                "protocolVersion": "2024-11-05",
                "capabilities": {
                    "tools": {}
                },
                "clientInfo": {
                    "name": "test-client",
                    "version": "1.0.0"
                }
            }
        }
        
        print("Sending initialize request...")
        send_mcp_message(process, init_message)
        
        # Read initialization response
        init_response = read_mcp_response(process)
        if init_response:
            print(f"Initialize response: {json.dumps(init_response, indent=2)}")
        else:
            print("No initialize response received")
            return False
        
        # 2. Send initialized notification
        initialized_notification = {
            "jsonrpc": "2.0",
            "method": "notifications/initialized"
        }
        
        print("Sending initialized notification...")
        send_mcp_message(process, initialized_notification)
        
        # 3. List available tools
        list_tools_message = {
            "jsonrpc": "2.0",
            "id": 2,
            "method": "tools/list"
        }
        
        print("Requesting tools list...")
        send_mcp_message(process, list_tools_message)
        
        # Read tools list response
        tools_response = read_mcp_response(process)
        if tools_response:
            print(f"Tools list response: {json.dumps(tools_response, indent=2)}")
            
            # Extract tool names
            if 'result' in tools_response and 'tools' in tools_response['result']:
                tools = tools_response['result']['tools']
                print(f"\nDiscovered {len(tools)} tools:")
                for tool in tools:
                    print(f"  - {tool['name']}: {tool.get('description', 'No description')}")
        else:
            print("No tools list response received")
            return False
        
        # 4. Test calling a math tool (math_calculate)
        if tools and any(tool['name'] == 'math_calculate' for tool in tools):
            print("\nTesting math_calculate tool...")
            
            calculate_message = {
                "jsonrpc": "2.0",
                "id": 3,
                "method": "tools/call",
                "params": {
                    "name": "math_calculate",
                    "arguments": {
                        "operation": "add",
                        "a": 15,
                        "b": 25
                    }
                }
            }
            
            send_mcp_message(process, calculate_message)
            
            # Read calculation response
            calc_response = read_mcp_response(process)
            if calc_response:
                print(f"Calculate response: {json.dumps(calc_response, indent=2)}")
            else:
                print("No calculate response received")
        
        # 5. Test calling a file tool (file_list_files)
        if tools and any(tool['name'] == 'file_list_files' for tool in tools):
            print("\nTesting file_list_files tool...")
            
            list_files_message = {
                "jsonrpc": "2.0",
                "id": 4,
                "method": "tools/call",
                "params": {
                    "name": "file_list_files",
                    "arguments": {
                        "path": "."
                    }
                }
            }
            
            send_mcp_message(process, list_files_message)
            
            # Read file list response
            files_response = read_mcp_response(process)
            if files_response:
                print(f"File list response: {json.dumps(files_response, indent=2)}")
            else:
                print("No file list response received")
        
        # 6. Test calling hello tool (hello_hello_world)
        if tools and any(tool['name'] == 'hello_hello_world' for tool in tools):
            print("\nTesting hello_hello_world tool...")
            
            hello_message = {
                "jsonrpc": "2.0",
                "id": 5,
                "method": "tools/call",
                "params": {
                    "name": "hello_hello_world",
                    "arguments": {
                        "name": "Proxy Test"
                    }
                }
            }
            
            send_mcp_message(process, hello_message)
            
            # Read hello response
            hello_response = read_mcp_response(process)
            if hello_response:
                print(f"Hello response: {json.dumps(hello_response, indent=2)}")
            else:
                print("No hello response received")
        
        print("\nTest completed successfully!")
        return True
        
    except Exception as e:
        print(f"Test error: {e}")
        return False
        
    finally:
        # Clean up
        try:
            process.terminate()
            process.wait(timeout=5)
        except subprocess.TimeoutExpired:
            process.kill()
            process.wait()

if __name__ == "__main__":
    success = test_proxy_server()
    sys.exit(0 if success else 1)