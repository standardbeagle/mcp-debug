# Task: Proxy Handlers Implementation
**Generated from Master Planning**: 2025-07-31
**Context Package**: `/requests/dynamic-mcp-proxy/context/`
**Next Phase**: [subtasks-execute.md](subtasks-execute.md)

## Task Sizing Assessment
**File Count**: 3-4 files - Within target
**Estimated Time**: 30 minutes - Target met
**Token Estimate**: ~90k tokens - Within target
**Complexity Level**: 3 (Complex) - Request/response transformation
**Parallelization Benefit**: MEDIUM - Some dependency on client
**Atomicity Assessment**: ✅ ATOMIC - Complete proxy functionality
**Boundary Analysis**: ✅ CLEAR - New proxy package

## Persona Assignment
**Persona**: Software Engineer
**Expertise Required**: Handler patterns, request forwarding, error handling
**Worktree**: `~/work/worktrees/dynamic-mcp-proxy/05-proxy/`

## Context Summary
**Risk Level**: HIGH - Core proxy functionality
**Integration Points**: Uses MCP client, integrates with mark3labs server
**Architecture Pattern**: Handler factory pattern with request forwarding
**Similar Reference**: main.go handler pattern, investigation results

### Codebase Context
**Files in Scope**:
```yaml
read_files: 
  - /main.go                          # Handler pattern reference
  - /investigations/test-tool-proxy.go # Proxy test results
  - /client/interface.go              # Client methods
create_files:
  - /proxy/handler.go      # Proxy handler factory
  - /proxy/transformer.go  # Request/response transformation
  - /proxy/errors.go      # Error wrapping and handling
  - /proxy/registry.go    # Tool-to-server mapping registry
```

### Task Scope Boundaries
**MODIFY Zone**:
```yaml
primary_files:
  - /proxy/*.go  # All proxy implementation
```

**REVIEW Zone**:
```yaml
check_integration:
  - /main.go                # Handler signature to match
  - /client/interface.go    # Client methods to use
  - /discovery/types.go     # RemoteTool structure
```

**IGNORE Zone**:
```yaml
ignore_completely:
  - /config/*       # Configuration system
  - /requests/*     # Planning documents
  - /server/*       # Server management (separate)
```

## Task Requirements
**Objective**: Create proxy handlers that forward tool calls to remote servers

**Success Criteria**:
- [ ] Create handler factory that generates proxy handlers
- [ ] Transform CallToolRequest to remote server format
- [ ] Forward requests using MCP client
- [ ] Transform responses back to CallToolResult
- [ ] Wrap errors with remote server context
- [ ] Maintain tool-to-server mapping

**Implementation Components**:

1. **Handler Factory** (`handler.go`):
```go
func CreateProxyHandler(client MCPClient, remoteTool RemoteTool) server.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // Extract arguments
        args := extractArguments(request)
        
        // Forward to remote server
        result, err := client.CallTool(ctx, remoteTool.OriginalName, args)
        if err != nil {
            return nil, wrapProxyError(remoteTool.ServerName, err)
        }
        
        // Transform result
        return transformResult(result), nil
    }
}
```

2. **Request/Response Transformer** (`transformer.go`):
   - Extract arguments from CallToolRequest
   - Convert to map[string]interface{} for client
   - Transform client result to CallToolResult
   - Handle different result types (text, error, etc.)

3. **Error Handling** (`errors.go`):
   - Wrap remote errors with context
   - Preserve error details
   - Add server identification
   - Handle timeout errors specially

4. **Tool Registry** (`registry.go`):
   - Map prefixed names to remote tools
   - Map tools to server clients
   - Efficient lookup during request handling

**Key Patterns**:
```go
// Registry for mapping tools to servers
type ToolRegistry struct {
    tools   map[string]RemoteTool
    clients map[string]MCPClient
}

// Argument extraction
func extractArguments(request mcp.CallToolRequest) map[string]interface{}

// Result transformation
func transformResult(clientResult *CallToolResult) *mcp.CallToolResult
```

**Validation Commands**:
```bash
go test ./proxy/...
# Test with mock client
go test -run TestProxyHandler ./proxy/
```

## Risk Mitigation
**High-Risk Mitigations**:
- Request format mismatch → Flexible argument extraction
- Response format issues → Handle multiple result types
- Connection failures → Clear error messages with server name

## Integration with Other Tasks
**Dependencies**: Task 02 (client), Task 04 (discovery types)
**Integration Points**: Used by task 07 (main integration)
**Shared Context**: Handler pattern from main.go

## Execution Notes
- Match handler signature from main.go exactly
- Test with various argument types
- Include request ID in error context
- Measure and log proxy latency