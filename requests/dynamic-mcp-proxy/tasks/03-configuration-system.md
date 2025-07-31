# Task: Configuration System Implementation
**Generated from Master Planning**: 2025-07-31
**Context Package**: `/requests/dynamic-mcp-proxy/context/`
**Next Phase**: [subtasks-execute.md](subtasks-execute.md)

## Task Sizing Assessment
**File Count**: 3-4 files - Within target
**Estimated Time**: 25 minutes - Target met
**Token Estimate**: ~80k tokens - Within target
**Complexity Level**: 2 (Moderate) - YAML parsing and validation
**Parallelization Benefit**: HIGH - Independent from client implementation
**Atomicity Assessment**: ✅ ATOMIC - Complete config system
**Boundary Analysis**: ✅ CLEAR - New config package

## Persona Assignment
**Persona**: Software Engineer
**Expertise Required**: YAML parsing, configuration validation, Go structs
**Worktree**: `~/work/worktrees/dynamic-mcp-proxy/03-config/`

## Context Summary
**Risk Level**: MEDIUM - Important but well-understood domain
**Integration Points**: Used by main server and connection manager
**Architecture Pattern**: YAML configuration with environment variable expansion
**Similar Reference**: Architecture decisions document

### Codebase Context
**Files in Scope**:
```yaml
read_files: 
  - /main.go  # Current config pattern (env vars)
create_files:
  - /config/types.go      # Configuration structs
  - /config/loader.go     # Loading and parsing logic
  - /config/validator.go  # Validation rules
  - /config/example.yaml  # Example configuration
```

### Task Scope Boundaries
**MODIFY Zone**:
```yaml
primary_files:
  - /config/*.go    # All configuration files
  - /config/*.yaml  # Example configs
```

**REVIEW Zone**:
```yaml
check_integration:
  - /go.mod  # Add gopkg.in/yaml.v3
```

**IGNORE Zone**:
```yaml
ignore_completely:
  - /client/*    # Client implementation
  - /main.go     # Server code
  - /requests/*  # Planning docs
```

## Task Requirements
**Objective**: Create configuration system for defining remote MCP servers

**Success Criteria**:
- [ ] Define configuration schema in Go structs
- [ ] Load YAML configuration files
- [ ] Support environment variable expansion
- [ ] Validate configuration (unique prefixes, valid transports)
- [ ] Create example configuration file

**Configuration Schema**:
```yaml
servers:
  - name: "math-server"
    prefix: "math"
    transport: "stdio"
    command: "/path/to/math-mcp-server"
    args: ["--flag", "value"]
    env:
      KEY: "value"
  
  - name: "api-server"
    prefix: "api"
    transport: "http"
    url: "http://localhost:8080"
    auth:
      type: "bearer"
      token: "${API_TOKEN}"
    timeout: "30s"
  
  - name: "data-server"
    prefix: "data"
    transport: "stdio"
    command: "${DATA_SERVER_PATH}/mcp-server"

proxy:
  healthCheckInterval: "30s"
  connectionTimeout: "10s"
  maxRetries: 3
```

**Implementation Requirements**:

1. **Type Definitions** (`types.go`):
   - ServerConfig struct with transport variants
   - ProxyConfig for global settings
   - AuthConfig for authentication

2. **Loader** (`loader.go`):
   - Read YAML file
   - Expand environment variables (${VAR} syntax)
   - Return parsed configuration

3. **Validator** (`validator.go`):
   - Unique prefix validation
   - Valid transport types
   - Required fields check
   - URL validation for HTTP transport

**Validation Commands**:
```bash
go test ./config/...
# Test with example config
go run main.go --config config/example.yaml
```

## Risk Mitigation
**Medium-Risk Mitigations**:
- Complex configs → Start with minimal fields
- Validation errors → Clear error messages
- Env var security → Don't log sensitive values

## Integration with Other Tasks
**Dependencies**: None - can start immediately
**Integration Points**: Used by tasks 06 (server manager) and 07 (integration)
**Shared Context**: Config schema used throughout

## Execution Notes
- Use gopkg.in/yaml.v3 for YAML parsing
- Include detailed validation error messages
- Support both file path and inline config
- Make example.yaml fully functional