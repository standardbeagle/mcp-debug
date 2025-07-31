# External Context Sources

## Primary Documentation

### MCP Protocol Specification
- **MCP Specification**: [https://modelcontextprotocol.io/specification/2025-06-18](https://modelcontextprotocol.io/specification/2025-06-18)
  - JSON-RPC 2.0 based protocol
  - Tool security requires explicit user consent
  - Client-server architecture with Hosts, Clients, and Servers
  - Tools expose functions for AI model execution

### mark3labs/mcp-go Library
- **GitHub Repository**: [https://github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go)
  - Go implementation of MCP protocol
  - Supports dynamic tool registration via AddTool
  - Tools must be added before ServeStdio() is called
  - Provides type-safe parameter definitions

- **Package Documentation**: [https://pkg.go.dev/github.com/mark3labs/mcp-go](https://pkg.go.dev/github.com/mark3labs/mcp-go)
  - Detailed API documentation
  - Server package for MCP server implementation
  - MCP package for protocol types and tool definitions

### Library Documentation
- **mcp-go Getting Started**: [https://mcp-go.dev/getting-started/](https://mcp-go.dev/getting-started/)
  - Basic usage patterns
  - Server creation and tool registration
  - Transport setup (stdio, HTTP)

## Industry Standards & Best Practices

### Security Standards
- **MCP Security Model**: Tools require explicit user consent before execution
  - Never execute tools without user authorization
  - Clear tool descriptions for user understanding
  - Secure handling of credentials for remote servers

### API Design
- **JSON-RPC 2.0**: [https://www.jsonrpc.org/specification](https://www.jsonrpc.org/specification)
  - Request/response protocol used by MCP
  - Error handling standards
  - Batch request support

### Testing Standards
- **Go Testing**: Standard Go testing practices
  - Unit tests for proxy handlers
  - Integration tests with mock MCP servers
  - Performance benchmarks for proxy overhead

## Reference Implementations

### Grafana MCP Implementation
- **Repository**: [https://github.com/grafana/mcp-grafana](https://github.com/grafana/mcp-grafana)
  - Shows pattern for tool organization
  - Demonstrates tool registration helpers
  - Example of complex tool implementations

### Example MCP Servers
- **MCP Hub**: [https://mcphub.tools/](https://mcphub.tools/)
  - Collection of MCP server implementations
  - Various patterns for tool organization
  - Examples of different transport types

## Standards Applied

### Coding Standards
- **Go Code Style**: Standard Go formatting (gofmt)
  - Error handling patterns
  - Context propagation
  - Interface design for extensibility

### MCP Protocol Standards
- **Tool Naming**: Use descriptive names with underscores
  - Pattern: `^[a-zA-Z0-9_-]{1,128}$`
  - Prefix pattern for proxy: `<prefix>_<original_tool_name>`

### Configuration Standards
- **Configuration Format**: JSON/YAML for server definitions
  - Support for multiple transport types
  - Authentication credentials handling
  - Connection pooling settings

## Key Implementation Insights

### Dynamic Tool Registration Limitations
From the research, mark3labs/mcp-go requires tools to be registered **before** calling `ServeStdio()`. This means:
1. Cannot add tools while server is running
2. Need to discover remote tools before starting our server
3. May need to restart server to add new remote connections

### Potential Solutions
1. **Pre-discovery**: Connect to all remote servers at startup, discover tools, then start
2. **Wrapper Server**: Create custom server that can restart with new tools
3. **Protocol Extension**: Investigate if MCP protocol supports dynamic tool updates

### Transport Considerations
- **stdio**: Process-based communication, need to spawn child processes
- **HTTP**: REST-based, easier for remote connections
- **WebSocket**: May be needed for persistent connections

### Authentication Patterns
- **API Keys**: Pass as headers or environment variables
- **OAuth**: For services requiring OAuth flows
- **mTLS**: For enterprise security requirements