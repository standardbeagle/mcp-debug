package main

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
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}