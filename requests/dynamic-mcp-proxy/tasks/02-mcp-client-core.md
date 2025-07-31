# Task: MCP Client Core Implementation
**Generated from Master Planning**: 2025-07-31
**Context Package**: `/requests/dynamic-mcp-proxy/context/`
**Next Phase**: [subtasks-execute.md](subtasks-execute.md)

## Task Sizing Assessment
**File Count**: 4-5 files - Within target
**Estimated Time**: 45 minutes - Acceptable for complex atomic unit
**Token Estimate**: ~120k tokens - Within target
**Complexity Level**: 3 (Complex) - Protocol implementation
**Parallelization Benefit**: HIGH - Independent from other tasks
**Atomicity Assessment**: ✅ ATOMIC - Complete client implementation
**Boundary Analysis**: ✅ CLEAR - New package, no conflicts

## Persona Assignment
**Persona**: Software Engineer (Protocol Specialist)
**Expertise Required**: JSON-RPC, Go interfaces, stdio/HTTP clients
**Worktree**: `~/work/worktrees/dynamic-mcp-proxy/02-client-core/`

## Context Summary
**Risk Level**: HIGH - Core dependency for proxy functionality
**Integration Points**: Used by tool discovery and proxy handlers
**Architecture Pattern**: Interface-based client with transport abstraction
**Similar Reference**: Investigation results from task 01

### Codebase Context
**Files in Scope**:
```yaml
read_files: 
  - /investigations/test-mcp-client.go  # Working patterns from investigation
create_files:
  - /client/interface.go      # MCPClient interface
  - /client/stdio_client.go   # Stdio transport implementation
  - /client/http_client.go    # HTTP transport implementation
  - /client/jsonrpc.go       # JSON-RPC protocol helpers
  - /client/types.go         # Request/response types
```

### Task Scope Boundaries
**MODIFY Zone**:
```yaml
primary_files:
  - /client/*.go  # All client implementation files
```

**REVIEW Zone**:
```yaml
check_integration:
  - /investigations/test-mcp-client.go  # Reference implementation
  - /go.mod  # May need new dependencies
```

**IGNORE Zone**:
```yaml
ignore_completely:
  - /main.go     # Server implementation
  - /requests/*  # Planning documents
```

## Task Requirements
**Objective**: Implement MCP client library for connecting to remote servers

**Success Criteria**:
- [ ] Define MCPClient interface with Connect, ListTools, CallTool, Close
- [ ] Implement stdio transport client
- [ ] Implement HTTP transport client (basic)
- [ ] Handle JSON-RPC request/response cycle
- [ ] Proper error handling and timeouts

**Implementation Details**:

1. **Interface Definition** (`interface.go`):
```go
type MCPClient interface {
    Connect(ctx context.Context) error
    Initialize(ctx context.Context) (*InitializeResult, error)
    ListTools(ctx context.Context) ([]ToolInfo, error)
    CallTool(ctx context.Context, name string, args map[string]interface{}) (*CallToolResult, error)
    Close() error
}
```

2. **Stdio Client** (`stdio_client.go`):
   - Process management with exec.Cmd
   - Bidirectional communication
   - Line-based JSON-RPC parsing

3. **HTTP Client** (`http_client.go`):
   - Standard HTTP POST for JSON-RPC
   - Connection pooling
   - Authentication support

4. **JSON-RPC Helpers** (`jsonrpc.go`):
   - Request ID generation
   - Response correlation
   - Error parsing

**Validation Commands**:
```bash
go test ./client/...
go build ./client/...
```

## Risk Mitigation
**High-Risk Mitigations**:
- Complex protocol → Start with minimal methods
- Transport differences → Abstract common logic
- Timeout handling → Use context everywhere

## Integration with Other Tasks
**Dependencies**: Investigation results from task 01
**Integration Points**: Used by tasks 04 (discovery) and 05 (proxy)
**Shared Context**: Protocol findings from investigations

## Execution Notes
- Use investigation findings for protocol details
- Implement timeouts on all operations
- Include debug logging for troubleshooting
- Test with real MCP server early