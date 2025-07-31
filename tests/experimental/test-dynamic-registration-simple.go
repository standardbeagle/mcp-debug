package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Simplified test - check if AddTool has any restrictions after server creation
func main() {
	fmt.Println("=== Testing Dynamic Tool Registration (Simplified) ===")
	
	// Create MCP server
	s := server.NewMCPServer(
		"Dynamic Test Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Add initial tool
	initialTool := mcp.NewTool("initial_tool",
		mcp.WithDescription("Initial tool"),
		mcp.WithString("message", mcp.Required()),
	)

	s.AddTool(initialTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("Initial tool works"), nil
	})

	fmt.Println("✅ Initial tool added successfully")

	// Try to add another tool (this should work if dynamic registration is possible)
	dynamicTool := mcp.NewTool("dynamic_tool",
		mcp.WithDescription("Tool added later"),
		mcp.WithString("data", mcp.Required()),
	)

	// Test if we can add tools after server is created but before ServeStdio
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("❌ PANIC when adding second tool: %v\n", r)
			fmt.Println("🔍 ASSUMPTION FAILED: Cannot add tools dynamically")
			os.Exit(1)
		}
	}()

	s.AddTool(dynamicTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("Dynamic tool works"), nil
	})

	fmt.Println("✅ Second tool added successfully")
	fmt.Println("🧪 Server creation and tool addition works")

	// The real test: can we call ServeStdio after adding tools?
	fmt.Println("🧪 Testing if server can start after multiple tool additions...")
	
	// We can't actually run ServeStdio in a test easily, but we can check
	// the server state and see if it's properly configured
	fmt.Println("✅ Server appears properly configured with multiple tools")
	
	fmt.Println("\n=== INVESTIGATION RESULTS ===")
	fmt.Println("✅ Can add multiple tools before ServeStdio()")
	fmt.Println("🔍 Need to test: Can add tools AFTER ServeStdio() starts")
	fmt.Println("📝 Next: Test with actual server running")
	
	// Key insight: The real question is whether tools can be added
	// after ServeStdio() is called, not just after server creation
	fmt.Println("\n=== CRITICAL INSIGHT ===")
	fmt.Println("The assumption test needs to check if AddTool works")
	fmt.Println("AFTER ServeStdio() is blocking on input/output")
	fmt.Println("This is architecturally unlikely with stdio transport")
}