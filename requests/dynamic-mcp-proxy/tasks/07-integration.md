# Task: Main Integration - Wire Everything Together
**Generated from Master Planning**: 2025-07-31
**Context Package**: `/requests/dynamic-mcp-proxy/context/`
**Next Phase**: [subtasks-execute.md](subtasks-execute.md)

## Task Sizing Assessment
**File Count**: 2-3 files (modify main.go, create cmd) - Within target
**Estimated Time**: 40 minutes - Acceptable for integration complexity
**Token Estimate**: ~110k tokens - Within target
**Complexity Level**: 3 (Complex) - Integrating all components
**Parallelization Benefit**: NONE - Requires all other tasks complete
**Atomicity Assessment**: ✅ ATOMIC - Complete integration
**Boundary Analysis**: ✅ CLEAR - Modify main.go, create integration code

## Persona Assignment
**Persona**: Senior Software Engineer (Integration Specialist)
**Expertise Required**: System integration, Go application architecture
**Worktree**: `~/work/worktrees/dynamic-mcp-proxy/07-integration/`

## Context Summary
**Risk Level**: HIGH - Brings all components together
**Integration Points**: All previous tasks
**Architecture Pattern**: Pre-discovery proxy with static registration
**Similar Reference**: Current main.go structure

### Codebase Context
**Files in Scope**:
```yaml
read_files: 
  - /main.go                 # Current server structure
  - All interfaces from previous tasks
modify_files:
  - /main.go                 # Add proxy functionality
create_files:
  - /cmd/proxy/main.go      # Alternative entry point
  - /integration/setup.go    # Integration helpers
```

### Task Scope Boundaries
**MODIFY Zone**:
```yaml
primary_files:
  - /main.go               # Extend with proxy features
  - /cmd/proxy/main.go     # New proxy-specific entry
  - /integration/setup.go  # Setup helpers
```

**REVIEW Zone**:
```yaml
check_integration:
  - /config/*.go      # Configuration loading
  - /client/*.go      # Client creation
  - /discovery/*.go   # Tool discovery
  - /proxy/*.go       # Handler creation
  - /manager/*.go     # Server management
```

**IGNORE Zone**:
```yaml
ignore_completely:
  - /investigations/*  # Test apps
  - /requests/*       # Planning docs
```

## Task Requirements
**Objective**: Integrate all components into working proxy server

**Success Criteria**:
- [ ] Load configuration from file/environment
- [ ] Initialize server manager with all remote servers
- [ ] Discover tools from all configured servers
- [ ] Register prefixed tools with MCP server
- [ ] Create proxy handlers for each tool
- [ ] Start MCP server with all tools
- [ ] Handle graceful shutdown

**Integration Flow**:

1. **Main Function Enhancement**:
```go
func main() {
    // Existing CLI handling...
    
    // Add proxy mode
    if proxyMode {
        if err := runProxyServer(); err != nil {
            log.Fatal(err)
        }
        return
    }
    
    // Original server mode...
}
```

2. **Proxy Server Setup** (`runProxyServer`):
   - Load configuration
   - Create server manager
   - Initialize all connections
   - Discover all tools
   - Create MCP server
   - Register all tools with handlers
   - Start server

3. **Integration Helpers** (`integration/setup.go`):
   - Configuration validation
   - Tool registration helper
   - Error aggregation
   - Shutdown coordination

**Key Integration Points**:
```go
func runProxyServer() error {
    // 1. Load configuration
    config, err := config.LoadConfig(*configPath)
    
    // 2. Create manager
    manager := manager.NewServerManager(config)
    
    // 3. Initialize connections
    ctx := context.Background()
    if err := manager.Initialize(ctx); err != nil {
        return fmt.Errorf("failed to initialize: %w", err)
    }
    
    // 4. Discover tools
    results, err := manager.DiscoverAllTools(ctx)
    
    // 5. Create MCP server
    s := server.NewMCPServer("Dynamic MCP Proxy", "1.0.0")
    
    // 6. Register tools
    registry := proxy.NewToolRegistry()
    for _, result := range results {
        for _, tool := range result.Tools {
            // Create prefixed tool definition
            mcpTool := createMCPTool(tool)
            
            // Get client for this server
            client, _ := manager.GetClient(result.Server)
            
            // Create and register handler
            handler := proxy.CreateProxyHandler(client, tool)
            s.AddTool(mcpTool, handler)
            
            registry.Register(tool, client)
        }
    }
    
    // 7. Start server
    return server.ServeStdio(s)
}
```

**Shutdown Handling**:
- Graceful shutdown on signals
- Close all client connections
- Clean up resources

**Validation Commands**:
```bash
# Build and test
go build -o mcp-proxy
./mcp-proxy --config example-config.yaml

# Test with MCP inspector
mcp-inspector ./mcp-proxy --config test-config.yaml
```

## Risk Mitigation
**High-Risk Mitigations**:
- Integration failures → Clear error messages at each step
- Partial tool discovery → Continue with available tools
- Configuration errors → Validate before starting

## Integration with Other Tasks
**Dependencies**: ALL previous tasks (01-06)
**Integration Points**: Uses all components
**Shared Context**: Architecture from planning phase

## Execution Notes
- Test each integration step separately
- Add detailed logging for debugging
- Handle partial failures gracefully
- Ensure backward compatibility with original server mode