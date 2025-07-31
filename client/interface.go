package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// MCPClient represents a client connection to an MCP server
type MCPClient interface {
	// Connect establishes connection to the MCP server
	Connect(ctx context.Context) error
	
	// Initialize performs MCP protocol handshake
	Initialize(ctx context.Context) (*InitializeResult, error)
	
	// ListTools discovers available tools from the server
	ListTools(ctx context.Context) ([]ToolInfo, error)
	
	// CallTool invokes a specific tool with arguments
	CallTool(ctx context.Context, name string, args map[string]interface{}) (*CallToolResult, error)
	
	// Close terminates the connection
	Close() error
	
	// ServerName returns the configured name of this server
	ServerName() string
	
	// IsConnected returns true if the client is currently connected
	IsConnected() bool
}

// InitializeResult represents the result of MCP initialize request
type InitializeResult struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ServerInfo      ServerInfo             `json:"serverInfo"`
}

// ServerInfo contains information about the MCP server
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ToolInfo represents information about a tool from the server
type ToolInfo struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

// CallToolResult represents the result of a tool invocation
type CallToolResult struct {
	Content []ContentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// ContentItem represents a piece of content in the tool result
type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// ClientError represents an error from the MCP client
type ClientError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Server  string `json:"server"`
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("[%s] %s (code: %d)", e.Server, e.Message, e.Code)
}

// NewClientError creates a new client error with server context
func NewClientError(server string, code int, message string) *ClientError {
	return &ClientError{
		Server:  server,
		Code:    code,
		Message: message,
	}
}