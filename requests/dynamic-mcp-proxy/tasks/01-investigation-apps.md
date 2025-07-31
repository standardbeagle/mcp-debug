# Task: Investigation Apps - Test Critical Assumptions
**Generated from Master Planning**: 2025-07-31
**Context Package**: `/requests/dynamic-mcp-proxy/context/`
**Next Phase**: [subtasks-execute.md](subtasks-execute.md)

## Task Sizing Assessment
**File Count**: 5 investigation apps - Within target
**Estimated Time**: 30 minutes - Target met
**Token Estimate**: ~100k tokens - Within target
**Complexity Level**: 2 (Moderate) - Multiple test scenarios
**Parallelization Benefit**: LOW - Sequential testing preferred
**Atomicity Assessment**: ✅ ATOMIC - Complete investigation phase
**Boundary Analysis**: ✅ CLEAR - Create new test files only

## Persona Assignment
**Persona**: Research Analyst
**Expertise Required**: Go development, MCP protocol understanding, JSON-RPC
**Worktree**: `~/work/worktrees/dynamic-mcp-proxy/01-investigation/`

## Context Summary
**Risk Level**: HIGH - Testing fundamental assumptions
**Integration Points**: Results inform all subsequent tasks
**Architecture Pattern**: Investigation apps to validate approach
**Similar Reference**: None - new investigation

### Codebase Context
**Files in Scope**:
```yaml
read_files: [/main.go, /go.mod]
create_files: 
  - /investigations/test-dynamic-registration.go
  - /investigations/test-mcp-client.go
  - /investigations/test-tool-discovery.go
  - /investigations/test-concurrent-servers.go
  - /investigations/test-tool-proxy.go
```

### Task Scope Boundaries
**MODIFY Zone**: 
```yaml
primary_files:
  - /investigations/*.go  # All investigation apps
```

**REVIEW Zone**:
```yaml
check_integration:
  - /main.go  # Understand current server pattern
  - /go.mod   # Check dependencies
```

**IGNORE Zone**:
```yaml
ignore_completely:
  - /requests/*  # Planning documents
  - /*.md       # Documentation files
```

## Task Requirements
**Objective**: Create and run 5 investigation apps to test critical assumptions

**Success Criteria**:
- [ ] Test if AddTool works after ServeStdio() starts
- [ ] Test connecting to MCP server as client
- [ ] Test tool discovery protocol methods
- [ ] Test concurrent server connections
- [ ] Test tool invocation proxying

**Investigation Apps**:

1. **test-dynamic-registration.go**
   - Start MCP server in goroutine
   - Attempt AddTool after start
   - Document result: WORKS/FAILS

2. **test-mcp-client.go**
   - Implement basic JSON-RPC client
   - Connect to stdio MCP server
   - Send requests and parse responses

3. **test-tool-discovery.go**
   - Try various methods: "tools/list", "initialize"
   - Document working discovery method

4. **test-concurrent-servers.go**
   - Start multiple servers
   - Connect to all simultaneously
   - Measure resource usage

5. **test-tool-proxy.go**
   - Forward tool request to remote
   - Transform request/response
   - Measure latency

**Validation Commands**:
```bash
cd investigations
go run test-dynamic-registration.go
go run test-mcp-client.go
go run test-tool-discovery.go
go run test-concurrent-servers.go
go run test-tool-proxy.go
```

## Risk Mitigation
**High-Risk Mitigations**:
- If dynamic registration fails → Document pre-discovery requirement
- If no standard discovery → Try multiple protocol methods
- If client implementation complex → Start with minimal subset

## Execution Notes
- Start with simplest tests first
- Document all findings in comments
- Include timing measurements
- Save working code patterns for reuse