# Task: Tool Discovery Implementation
**Generated from Master Planning**: 2025-07-31
**Context Package**: `/requests/dynamic-mcp-proxy/context/`
**Next Phase**: [subtasks-execute.md](subtasks-execute.md)

## Task Sizing Assessment
**File Count**: 3 files - Within target
**Estimated Time**: 25 minutes - Target met
**Token Estimate**: ~70k tokens - Within target
**Complexity Level**: 2 (Moderate) - Protocol-specific implementation
**Parallelization Benefit**: MEDIUM - Depends on client but separate concern
**Atomicity Assessment**: ✅ ATOMIC - Complete discovery feature
**Boundary Analysis**: ✅ CLEAR - Extends client, new discovery package

## Persona Assignment
**Persona**: Software Engineer
**Expertise Required**: MCP protocol, concurrent programming
**Worktree**: `~/work/worktrees/dynamic-mcp-proxy/04-discovery/`

## Context Summary
**Risk Level**: MEDIUM - Depends on protocol understanding
**Integration Points**: Uses MCP client, provides tools to proxy
**Architecture Pattern**: Concurrent discovery with timeout handling
**Similar Reference**: Investigation results from test-tool-discovery.go

### Codebase Context
**Files in Scope**:
```yaml
read_files: 
  - /investigations/test-tool-discovery.go  # Discovery method findings
  - /client/interface.go                    # Client interface to use
create_files:
  - /discovery/discoverer.go    # Main discovery logic
  - /discovery/types.go        # Discovery-specific types
  - /discovery/concurrent.go   # Concurrent discovery for multiple servers
```

### Task Scope Boundaries
**MODIFY Zone**:
```yaml
primary_files:
  - /discovery/*.go  # All discovery implementation
```

**REVIEW Zone**:
```yaml
check_integration:
  - /client/interface.go  # Understand client methods
  - /config/types.go     # Server configuration types
```

**IGNORE Zone**:
```yaml
ignore_completely:
  - /main.go        # Server implementation
  - /proxy/*        # Proxy handlers (separate task)
  - /requests/*     # Planning documents
```

## Task Requirements
**Objective**: Implement tool discovery from remote MCP servers

**Success Criteria**:
- [ ] Discover tools from single server using client
- [ ] Handle discovery failures gracefully
- [ ] Support concurrent discovery from multiple servers
- [ ] Create prefixed tool definitions
- [ ] Return structured discovery results

**Implementation Components**:

1. **Discovery Types** (`types.go`):
```go
type DiscoveryResult struct {
    Server     string
    Tools      []RemoteTool
    Error      error
    Duration   time.Duration
}

type RemoteTool struct {
    OriginalName string
    PrefixedName string
    Description  string
    Schema       json.RawMessage
    ServerName   string
    ServerPrefix string
}
```

2. **Single Server Discovery** (`discoverer.go`):
   - Connect to server using client
   - Call Initialize to get capabilities
   - Extract tool information
   - Apply prefix to tool names

3. **Concurrent Discovery** (`concurrent.go`):
   - Discover from multiple servers in parallel
   - Timeout handling per server
   - Aggregate results
   - Report partial failures

**Key Functions**:
```go
// Discover tools from a single server
func DiscoverServer(ctx context.Context, client MCPClient, config ServerConfig) (*DiscoveryResult, error)

// Discover from multiple servers concurrently
func DiscoverAll(ctx context.Context, configs []ServerConfig) ([]*DiscoveryResult, error)

// Create prefixed tool from remote tool info
func CreatePrefixedTool(prefix string, remoteTool ToolInfo) RemoteTool
```

**Validation Commands**:
```bash
go test ./discovery/...
# Integration test with real server
go run discovery/cmd/test/main.go
```

## Risk Mitigation
**Medium-Risk Mitigations**:
- Protocol variations → Try multiple discovery methods
- Timeout issues → Configurable timeouts per server
- Partial failures → Continue with working servers

## Integration with Other Tasks
**Dependencies**: Task 02 (MCP client)
**Integration Points**: Results used by task 05 (proxy handlers) and 07 (integration)
**Shared Context**: Discovery patterns from investigations

## Execution Notes
- Use investigation findings for protocol details
- Handle servers that don't support tool listing
- Include server name in all errors
- Test with multiple server types