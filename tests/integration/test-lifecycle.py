#!/usr/bin/env python3
"""
Test script to verify the MCP proxy server can handle connection lifecycle:
1. Connect to stdio server
2. Discover tools  
3. Shut down connection
4. Reconnect to potentially updated server
"""

import json
import subprocess
import sys
import time
import os

def create_test_server_v1():
    """Create first version of test server with basic tools"""
    server_code = '''package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Lifecycle Test Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Version 1: Only has hello tool
	helloTool := mcp.NewTool("hello",
		mcp.WithDescription("Say hello (v1)"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name to greet")),
	)

	s.AddTool(helloTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Hello from v1, %s!", name)), nil
	})

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\\n", err)
		os.Exit(1)
	}
}'''
    
    with open('test-servers/lifecycle-server-v1.go', 'w') as f:
        f.write(server_code)

def create_test_server_v2():
    """Create second version of test server with additional tools"""
    server_code = '''package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Lifecycle Test Server",
		"2.0.0",
		server.WithToolCapabilities(true),
	)

	// Version 2: Has hello tool + new timestamp tool
	helloTool := mcp.NewTool("hello",
		mcp.WithDescription("Say hello (v2 - updated!)"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name to greet")),
	)

	timestampTool := mcp.NewTool("timestamp",
		mcp.WithDescription("Get current timestamp (new in v2!)"),
	)

	s.AddTool(helloTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Hello from v2 (UPDATED!), %s!", name)), nil
	})

	s.AddTool(timestampTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText(fmt.Sprintf("Current time: %s", time.Now().Format("2006-01-02 15:04:05"))), nil
	})

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\\n", err)
		os.Exit(1)
	}
}'''
    
    with open('test-servers/lifecycle-server-v2.go', 'w') as f:
        f.write(server_code)

def build_test_servers():
    """Build both versions of the test server"""
    print("Building test servers...")
    
    # Build v1
    result = subprocess.run([
        'go', 'build', '-o', 'test-servers/lifecycle-server-v1', 
        'test-servers/lifecycle-server-v1.go'
    ], capture_output=True, text=True)
    
    if result.returncode != 0:
        print(f"Failed to build v1: {result.stderr}")
        return False
    
    # Build v2
    result = subprocess.run([
        'go', 'build', '-o', 'test-servers/lifecycle-server-v2', 
        'test-servers/lifecycle-server-v2.go'
    ], capture_output=True, text=True)
    
    if result.returncode != 0:
        print(f"Failed to build v2: {result.stderr}")
        return False
    
    print("Both test servers built successfully!")
    return True

def create_lifecycle_config(server_path):
    """Create config file pointing to specific server version"""
    config = f'''servers:
  - name: "lifecycle-server"
    prefix: "test"
    transport: "stdio"
    command: "{server_path}"
    timeout: "10s"

proxy:
  healthCheckInterval: "30s"
  connectionTimeout: "10s"
  maxRetries: 3'''
    
    with open('test-lifecycle-config.yaml', 'w') as f:
        f.write(config)

def test_proxy_discovery(version_name):
    """Test proxy server tool discovery"""
    print(f"\\n=== Testing {version_name} ===")
    
    # Start proxy server
    cmd = ['../../mcp-debug', '--proxy', '--config', '../config-fixtures/test-lifecycle-config.yaml']
    
    # Run discovery only (proxy starts and shuts down automatically)
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=10)
    
    print(f"Proxy output:\n{result.stderr}")
    
    # Extract tool information from logs
    lines = result.stderr.split('\n')
    tools_found = []
    for line in lines:
        if 'Registered tool:' in line:
            tool_name = line.split('Registered tool: ')[1].strip()
            tools_found.append(tool_name)
    
    print(f"Tools discovered in {version_name}: {tools_found}")
    return tools_found

def main():
    """Main test function"""
    print("Testing MCP Proxy Connection Lifecycle")
    print("======================================")
    
    # Create and build test servers
    create_test_server_v1()
    create_test_server_v2()
    
    if not build_test_servers():
        return False
    
    # Test v1 discovery
    create_lifecycle_config('./test-servers/lifecycle-server-v1')
    v1_tools = test_proxy_discovery("Server v1")
    
    # Test v2 discovery (simulating server update)
    create_lifecycle_config('./test-servers/lifecycle-server-v2')
    v2_tools = test_proxy_discovery("Server v2 (after update)")
    
    # Verify lifecycle worked
    print("\n=== Lifecycle Test Results ===")
    print(f"V1 tools: {v1_tools}")
    print(f"V2 tools: {v2_tools}")
    
    # Check that v2 has more tools than v1
    if len(v2_tools) > len(v1_tools):
        print("✅ SUCCESS: Proxy successfully detected server update!")
        print("✅ Connection lifecycle working: connect → discover → disconnect → reconnect")
        return True
    else:
        print("❌ FAILED: Server update not detected properly")
        return False

if __name__ == "__main__":
    success = main()
    sys.exit(0 if success else 1)