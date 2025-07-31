package main

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
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}