package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"File Operations Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Add list files tool
	listTool := mcp.NewTool("list_files",
		mcp.WithDescription("List files in a directory"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Directory path to list"),
		),
	)

	// Add read file tool
	readTool := mcp.NewTool("read_file",
		mcp.WithDescription("Read contents of a text file"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("File path to read"),
		),
	)

	// Add file info tool
	infoTool := mcp.NewTool("file_info",
		mcp.WithDescription("Get information about a file or directory"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("File or directory path"),
		),
	)

	// Add tool handlers
	s.AddTool(listTool, listFilesHandler)
	s.AddTool(readTool, readFileHandler)
	s.AddTool(infoTool, fileInfoHandler)

	// Start stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "File server error: %v\n", err)
		os.Exit(1)
	}
}

func listFilesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Clean the path
	path = filepath.Clean(path)

	// Read directory
	entries, err := os.ReadDir(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read directory: %v", err)), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Contents of %s:\n", path))
	
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		
		fileType := "file"
		if entry.IsDir() {
			fileType = "directory"
		}
		
		result.WriteString(fmt.Sprintf("- %s (%s, %d bytes)\n", 
			entry.Name(), fileType, info.Size()))
	}

	return mcp.NewToolResultText(result.String()), nil
}

func readFileHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Clean the path
	path = filepath.Clean(path)

	// Read file
	content, err := os.ReadFile(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read file: %v", err)), nil
	}

	// Limit output size for safety
	if len(content) > 10000 {
		return mcp.NewToolResultText(fmt.Sprintf("File content (first 10000 bytes):\n%s\n\n[File truncated - total size: %d bytes]", 
			string(content[:10000]), len(content))), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("File content:\n%s", string(content))), nil
}

func fileInfoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Clean the path
	path = filepath.Clean(path)

	// Get file info
	info, err := os.Stat(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get file info: %v", err)), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("File information for: %s\n", path))
	result.WriteString(fmt.Sprintf("Name: %s\n", info.Name()))
	result.WriteString(fmt.Sprintf("Size: %d bytes\n", info.Size()))
	result.WriteString(fmt.Sprintf("Mode: %s\n", info.Mode()))
	result.WriteString(fmt.Sprintf("Modified: %s\n", info.ModTime().Format("2006-01-02 15:04:05")))
	
	if info.IsDir() {
		result.WriteString("Type: Directory\n")
	} else {
		result.WriteString("Type: File\n")
	}

	return mcp.NewToolResultText(result.String()), nil
}