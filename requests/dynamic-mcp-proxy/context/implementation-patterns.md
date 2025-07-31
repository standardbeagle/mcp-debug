# Implementation Patterns

## Code Patterns from Existing Codebase

### Tool Definition Pattern
From `main.go:65-71`:
```go
tool := mcp.NewTool("hello_world",
    mcp.WithDescription("Say hello to someone"),
    mcp.WithString("name",
        mcp.Required(),
        mcp.Description("Name of person to greet"),
    ),
)
```
**Application**: Use same pattern for defining proxy tools with prefixed names

### Handler Pattern
From `main.go:82-89`:
```go
func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    name, err := request.RequireString("name")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }
    return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
```
**Application**: Create proxy handlers that forward to remote servers

### Server Creation Pattern
From `main.go:58-62`:
```go
s := server.NewMCPServer(
    "Dynamic MCP Proxy Server",
    "1.0.0",
    server.WithToolCapabilities(true),
)
```

## New Implementation Patterns

### Remote Server Configuration
```go
type RemoteServer struct {
    Name      string
    Prefix    string
    Transport string
    Config    TransportConfig
    Client    MCPClient
    Tools     []RemoteTool
}

type TransportConfig struct {
    // For stdio
    Command string
    Args    []string
    
    // For HTTP
    URL  string
    Auth AuthConfig
}
```

### MCP Client Interface
```go
type MCPClient interface {
    Connect(ctx context.Context) error
    ListTools(ctx context.Context) ([]ToolInfo, error)
    CallTool(ctx context.Context, name string, args map[string]interface{}) (*mcp.CallToolResult, error)
    Close() error
}
```

### Proxy Handler Factory
```go
func createProxyHandler(client MCPClient, remoteTool string) server.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // Extract arguments
        args := make(map[string]interface{})
        // ... populate args from request
        
        // Forward to remote server
        result, err := client.CallTool(ctx, remoteTool, args)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Remote error: %v", err)), nil
        }
        
        return result, nil
    }
}
```

### Tool Discovery Pattern
```go
func discoverTools(ctx context.Context, server RemoteServer) ([]RemoteTool, error) {
    // Connect to remote server
    if err := server.Client.Connect(ctx); err != nil {
        return nil, fmt.Errorf("failed to connect to %s: %w", server.Name, err)
    }
    
    // List available tools
    tools, err := server.Client.ListTools(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to list tools from %s: %w", server.Name, err)
    }
    
    // Create prefixed tools
    var remoteTools []RemoteTool
    for _, tool := range tools {
        prefixedName := fmt.Sprintf("%s_%s", server.Prefix, tool.Name)
        remoteTools = append(remoteTools, RemoteTool{
            OriginalName: tool.Name,
            PrefixedName: prefixedName,
            Schema:       tool.Schema,
            Server:       server.Name,
        })
    }
    
    return remoteTools, nil
}
```

### Configuration Loading Pattern
```go
func loadConfiguration(path string) (*ProxyConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    
    var config ProxyConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    
    // Validate configuration
    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    
    // Expand environment variables
    config.ExpandEnvVars()
    
    return &config, nil
}
```

### JSON-RPC Client Pattern
```go
type StdioMCPClient struct {
    cmd    *exec.Cmd
    stdin  io.WriteCloser
    stdout io.ReadCloser
    reader *bufio.Reader
    reqID  int64
}

func (c *StdioMCPClient) sendRequest(method string, params interface{}) (*json.RawMessage, error) {
    reqID := atomic.AddInt64(&c.reqID, 1)
    
    request := JSONRPCRequest{
        JSONRPC: "2.0",
        Method:  method,
        Params:  params,
        ID:      reqID,
    }
    
    // Send request
    if err := json.NewEncoder(c.stdin).Encode(request); err != nil {
        return nil, err
    }
    
    // Read response
    // ... handle response parsing
}
```

### Error Handling Pattern
```go
type ProxyError struct {
    Server       string
    OriginalError error
    Context      string
}

func (e ProxyError) Error() string {
    return fmt.Sprintf("[%s] %s: %v", e.Server, e.Context, e.OriginalError)
}

func wrapRemoteError(server string, context string, err error) error {
    return ProxyError{
        Server:        server,
        OriginalError: err,
        Context:       context,
    }
}
```

### Health Check Pattern
```go
func (p *ProxyServer) startHealthChecks(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            for _, server := range p.remoteServers {
                go p.checkServerHealth(server)
            }
        }
    }
}
```

## Testing Patterns

### Mock MCP Server
```go
type MockMCPServer struct {
    tools map[string]ToolInfo
    responses map[string]interface{}
}

func (m *MockMCPServer) Start() (string, error) {
    // Start on random port
    // Return connection string
}
```

### Integration Test Pattern
```go
func TestProxyEndToEnd(t *testing.T) {
    // 1. Start mock MCP server
    // 2. Configure proxy to connect
    // 3. Start proxy server
    // 4. Connect as client
    // 5. List tools (verify prefix)
    // 6. Call tool (verify forwarding)
    // 7. Check response
}
```

## Performance Patterns

### Connection Pooling
```go
type ConnectionPool struct {
    servers map[string]MCPClient
    mu      sync.RWMutex
}

func (p *ConnectionPool) GetClient(server string) (MCPClient, error) {
    p.mu.RLock()
    client, exists := p.servers[server]
    p.mu.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("unknown server: %s", server)
    }
    
    return client, nil
}
```

### Concurrent Tool Discovery
```go
func discoverAllTools(ctx context.Context, servers []RemoteServer) (map[string][]RemoteTool, error) {
    var wg sync.WaitGroup
    toolsChan := make(chan serverTools, len(servers))
    errorsChan := make(chan error, len(servers))
    
    for _, server := range servers {
        wg.Add(1)
        go func(s RemoteServer) {
            defer wg.Done()
            
            tools, err := discoverTools(ctx, s)
            if err != nil {
                errorsChan <- err
                return
            }
            
            toolsChan <- serverTools{
                server: s.Name,
                tools:  tools,
            }
        }(server)
    }
    
    wg.Wait()
    close(toolsChan)
    close(errorsChan)
    
    // Collect results
    // ... handle errors and aggregate tools
}
```