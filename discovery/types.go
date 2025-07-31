package discovery

import (
	"encoding/json"
	"time"
)

// DiscoveryResult represents the result of discovering tools from a server
type DiscoveryResult struct {
	ServerName   string        `json:"serverName"`
	ServerPrefix string        `json:"serverPrefix"`
	Tools        []RemoteTool  `json:"tools"`
	Error        error         `json:"error,omitempty"`
	Duration     time.Duration `json:"duration"`
}

// RemoteTool represents a tool discovered from a remote server
type RemoteTool struct {
	OriginalName string          `json:"originalName"`
	PrefixedName string          `json:"prefixedName"`
	Description  string          `json:"description"`
	InputSchema  json.RawMessage `json:"inputSchema"`
	ServerName   string          `json:"serverName"`
	ServerPrefix string          `json:"serverPrefix"`
}

// IsSuccessful returns true if the discovery was successful
func (r *DiscoveryResult) IsSuccessful() bool {
	return r.Error == nil
}

// ToolCount returns the number of tools discovered
func (r *DiscoveryResult) ToolCount() int {
	return len(r.Tools)
}

// CreatePrefixedTool creates a RemoteTool with proper prefixing
func CreatePrefixedTool(serverName, serverPrefix string, originalTool ToolInfo) RemoteTool {
	prefixedName := serverPrefix + "_" + originalTool.Name
	
	return RemoteTool{
		OriginalName: originalTool.Name,
		PrefixedName: prefixedName,
		Description:  originalTool.Description,
		InputSchema:  originalTool.InputSchema,
		ServerName:   serverName,
		ServerPrefix: serverPrefix,
	}
}

// ToolInfo represents tool information from the MCP client
type ToolInfo struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}