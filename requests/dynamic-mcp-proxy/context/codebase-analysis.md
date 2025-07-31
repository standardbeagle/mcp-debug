# Codebase Context Documentation

## Existing Architecture Patterns

### Service Architecture
- **MCP Server Pattern**: Simple stdio-based server using mark3labs/mcp-go - Files: [`main.go`]
- **Tool Registration**: Static tool registration at startup via `s.AddTool()` - Files: [`main.go:70`]
- **Handler Pattern**: Context-based handlers with request/result types - Files: [`main.go:82-89`]
- **CLI Interface**: Dual-mode with CLI detection for testing - Files: [`main.go:37-49`]

### Data Layer
- **Configuration**: Environment variable based (MCP_DEBUG, MCP_CONFIG_PATH) - Files: [`main.go`]
- **State Management**: No persistent state, all in-memory
- **Tool Storage**: Tools stored in server instance via AddTool

### API Patterns
- **MCP Protocol**: Using mark3labs implementation
- **Transport**: stdio-based using `server.ServeStdio()`
- **Tool Interface**: `mcp.NewTool()` with schema definition
- **Result Types**: `mcp.CallToolResult` with Text/Error variants

### Authentication
- No authentication implemented in current server
- Would need to add for remote server connections

### Configuration
- Environment variables for debug and config path
- No structured config file implementation yet

## Similar Feature Implementations

### Tool Registration Pattern
- **Location**: `main.go:65-71`
- **Pattern**: Static tool creation and registration
- **Relevance**: Need to make this dynamic for proxy functionality

### Handler Pattern
- **Location**: `main.go:82-89`
- **Pattern**: Simple request/response handlers
- **Relevance**: Need to create proxy handler that forwards to remote servers

### Server Creation
- **Location**: `main.go:58-62`
- **Pattern**: Server instance with capabilities
- **Relevance**: May need custom server options for dynamic tools

## Dependency Analysis

### Core Dependencies
- **github.com/mark3labs/mcp-go@v0.36.0** - [GitHub](https://github.com/mark3labs/mcp-go) - Core MCP implementation
- **Go 1.24.2** - Latest Go version

### Dev Dependencies
- Testing frameworks: Standard Go testing
- Build tools: Go build

### External Services
- None currently - will need to connect to external MCP servers

## File Dependency Mapping

```yaml
high_change_areas:
  - /main.go: [Add dynamic tool registration, proxy configuration]
  - /NEW: config/: [Configuration for remote servers]
  - /NEW: proxy/: [Proxy client implementation]
  
medium_change_areas:
  - /go.mod: [May need additional dependencies]
  - /NEW: types/: [Configuration types]
  - /NEW: client/: [MCP client for connecting to remotes]

low_change_areas:
  - /.gitignore: [No changes needed]
  - /README.md: [Documentation updates]
```

## Current Implementation Limitations

### Static Tool Registration
The current implementation registers tools at startup:
```go
tool := mcp.NewTool("hello_world", ...)
s.AddTool(tool, helloHandler)
```
This happens before `ServeStdio()` is called. Need to investigate if tools can be added dynamically.

### No Client Implementation
Current code only implements server-side. Need to add MCP client capabilities to connect to remote servers.

### No Configuration System
Need to implement configuration loading for remote server definitions.

### Single Transport
Only stdio transport implemented. May need HTTP client support for remote servers.