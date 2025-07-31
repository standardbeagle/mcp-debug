package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Test if tools can be added after ServeStdio() starts
func main() {
	fmt.Println("=== Testing Dynamic Tool Registration ===")
	
	// Create MCP server
	s := server.NewMCPServer(
		"Dynamic Test Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Add initial tool
	initialTool := mcp.NewTool("initial_tool",
		mcp.WithDescription("Initial tool added before server start"),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("Message to return"),
		),
	)

	s.AddTool(initialTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		message, _ := request.RequireString("message")
		return mcp.NewToolResultText(fmt.Sprintf("Initial: %s", message)), nil
	})

	fmt.Println("✅ Initial tool added before server start")

	// Start server in a goroutine to test if we can add tools after
	done := make(chan error, 1)
	go func() {
		fmt.Println("🚀 Starting MCP server...")
		// This will block, so we run it in goroutine
		err := server.ServeStdio(s)
		done <- err
	}()

	// Wait a moment for server to initialize
	time.Sleep(2 * time.Second)
	fmt.Println("⏰ Server should be running now...")

	// Try to add a new tool after server started
	fmt.Println("🧪 Attempting to add tool after server start...")
	
	dynamicTool := mcp.NewTool("dynamic_tool",
		mcp.WithDescription("Tool added after server start"),
		mcp.WithString("data",
			mcp.Required(),
			mcp.Description("Data to process"),
		),
	)

	// This is the critical test - can we add tools after ServeStdio()?
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("❌ PANIC when adding dynamic tool: %v\n", r)
			fmt.Println("🔍 ASSUMPTION FAILED: Cannot add tools after server starts")
			os.Exit(1)
		}
	}()

	// Attempt the dynamic addition
	s.AddTool(dynamicTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		data, _ := request.RequireString("data")
		return mcp.NewToolResultText(fmt.Sprintf("Dynamic: %s", data)), nil
	})

	fmt.Println("✅ Dynamic tool addition succeeded!")
	fmt.Println("🎉 ASSUMPTION VALIDATED: Can add tools after server starts")

	// Wait a bit more to see if server is still stable
	time.Sleep(1 * time.Second)

	select {
	case err := <-done:
		if err != nil {
			fmt.Printf("❌ Server error: %v\n", err)
		} else {
			fmt.Println("✅ Server completed successfully")
		}
	case <-time.After(1 * time.Second):
		fmt.Println("✅ Server still running stable after dynamic tool addition")
	}

	fmt.Println("\n=== RESULTS ===")
	fmt.Println("Dynamic Tool Registration: SUCCESS")
	fmt.Println("Plan Impact: Original approach is viable")
}