# MCP Debug

The essential debugging and development tool for [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) servers.

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![MCP Spec](https://img.shields.io/badge/MCP-2025--06--18-green.svg)](https://modelcontextprotocol.io/specification/2025-06-18)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**MCP Debug** enables rapid development and testing of MCP servers with hot-swapping, session recording, and automated playback testing. Built for developers building the MCP ecosystem.

## ✨ Features

### 🔄 **Hot-Swap Development**
- **Replace server binaries** without disconnecting MCP clients
- **Add/remove servers** dynamically during development
- **Tool name preservation** - same interface, new implementation
- **Graceful disconnect/reconnect** workflow for binary replacement

### 📹 **Session Recording & Playback**
- **Record JSON-RPC traffic** for debugging and documentation
- **Playback client mode** - replay requests to test servers
- **Playback server mode** - replay responses to test clients  
- **Regression testing** with recorded sessions

### 🛠️ **Development Proxy**
- **Multi-server aggregation** with tool prefixing
- **Real-time connection monitoring** with automatic failure detection
- **Management API** for server lifecycle control
- **Comprehensive logging** with configurable output

## 🚀 Quick Start

```bash
# Install
go install github.com/your-org/mcp-debug@latest

# Or build from source
git clone https://github.com/your-org/mcp-debug
cd mcp-debug
go build -o mcp-debug .

# Start debugging proxy
mcp-debug --proxy --config empty-config.yaml

# Connect with mcp-tui
mcp-tui ./mcp-debug --proxy --config empty-config.yaml
```

## 🎯 Core Workflow

### Development Cycle

```bash
# 1. Start MCP Debug proxy
mcp-tui ./mcp-debug --proxy --config config.yaml

# 2. Add your server dynamically
server_add:
  name: myserver
  command: ./my-mcp-server-v1

# 3. Test tools: myserver_read_file, myserver_process, etc.

# 4. Make changes, rebuild
# Edit code, fix bugs...
go build -o my-mcp-server-v2

# 5. Hot-swap the server
server_disconnect: {name: myserver}
server_reconnect: {name: myserver, command: ./my-mcp-server-v2}

# 6. Same tools work immediately with new implementation! 🎉
```

### Testing Workflow

```bash
# Record a working session
MCP_RECORD_FILE="working.jsonl" mcp-tui ./mcp-debug --proxy --config config.yaml

# Test changes automatically
./mcp-debug --playback-client working.jsonl | ./my-new-server > results.txt

# Compare with expected output
diff expected-results.txt results.txt
```

## 📖 Usage Modes

### 1. Proxy Mode (Primary)

Dynamic proxy with hot-swapping capabilities:

```bash
# Basic proxy
./mcp-debug --proxy --config config.yaml

# With recording
./mcp-debug --proxy --config config.yaml --record session.jsonl

# With custom logging
./mcp-debug --proxy --config config.yaml --log /tmp/debug.log
```

**Management Tools Available:**
- **`server_add`** - Add server: `{name: "fs", command: "npx -y @mcp/filesystem /path"}`
- **`server_remove`** - Remove server completely
- **`server_disconnect`** - Disconnect (tools return errors, enables binary swap)  
- **`server_reconnect`** - Reconnect with new command (after disconnect)
- **`server_list`** - Show all servers and connection status

### 2. Recording Mode

Capture JSON-RPC traffic:

```bash
# Environment variable method
MCP_RECORD_FILE="debug-session.jsonl" mcp-tui ./run-proxy.sh config.yaml

# Command line flag method  
./mcp-debug --proxy --config config.yaml --record session.jsonl
```

**Recording Format:**
```json
{"timestamp":"2025-07-31T12:00:00Z","direction":"request","message_type":"tool_call","tool_name":"fs_read_file","server_name":"fs","message":{...}}
{"timestamp":"2025-07-31T12:00:01Z","direction":"response","message_type":"tool_call","tool_name":"fs_read_file","server_name":"fs","message":{...}}
```

### 3. Playback Client Mode

Replay recorded client requests to test servers:

```bash
# Test server with recorded requests
./mcp-debug --playback-client session.jsonl | ./your-mcp-server

# Regression testing
./mcp-debug --playback-client baseline.jsonl | ./new-server > new-results.txt
./mcp-debug --playback-client baseline.jsonl | ./old-server > old-results.txt
diff new-results.txt old-results.txt
```

### 4. Playback Server Mode

Replay recorded server responses to test clients:

```bash
# Test mcp-tui with recorded responses
mcp-tui ./mcp-debug --playback-server session.jsonl

# Test your custom client
your-mcp-client ./mcp-debug --playback-server session.jsonl
```

## ⚙️ Configuration

### Basic Configuration

```yaml
# config.yaml
servers:
  - name: "filesystem"
    prefix: "fs"
    transport: "stdio"
    command: "npx"
    args: ["-y", "@modelcontextprotocol/filesystem", "/home/user"]
    timeout: "30s"
    
  - name: "database"  
    prefix: "db"
    transport: "stdio"
    command: "./db-mcp-server"
    args: ["--conn", "postgres://localhost/db"]

proxy:
  healthCheckInterval: "30s"
  connectionTimeout: "10s" 
  maxRetries: 3
```

### Environment Variables

```bash
# Logging
MCP_LOG_FILE="/tmp/mcp-debug.log"     # Log location
MCP_DEBUG=1                           # Enable debug logging

# Recording
MCP_RECORD_FILE="session.jsonl"       # Auto-record sessions

# Configuration  
MCP_CONFIG_PATH="./config.yaml"       # Default config
```

### Empty Configuration

For dynamic-only usage:

```yaml
# empty-config.yaml
servers: []
proxy:
  healthCheckInterval: "30s"
  connectionTimeout: "10s"
  maxRetries: 3
```

## 🏗️ Use Cases

### MCP Server Development

**Problem**: Constantly restarting MCP clients during development is slow and breaks flow.

**Solution**: Hot-swap servers without client disconnection.

```bash
# Start development session
mcp-tui ./mcp-debug --proxy --config empty-config.yaml

# Add initial server
server_add: {name: api, command: ./api-server-v1}

# Develop, test, iterate...
# server_disconnect: {name: api}
# <rebuild binary>
# server_reconnect: {name: api, command: ./api-server-v2}

# Same tools (api_get_user, api_create_post) work immediately!
```

### Multi-Server Integration Testing

**Problem**: Testing interactions between multiple MCP servers is complex.

**Solution**: Aggregate multiple servers with tool prefixing.

```bash
server_add: {name: auth, command: ./auth-server}
server_add: {name: db, command: ./db-server}  
server_add: {name: cache, command: ./cache-server}

# Now test workflows across servers:
# auth_login -> db_get_user -> cache_set_session
```

### Regression Testing

**Problem**: Manual testing is time-consuming and error-prone.

**Solution**: Record working sessions, replay for testing.

```bash
# 1. Record working session
MCP_RECORD_FILE="regression.jsonl" mcp-tui ./working-server

# 2. Automate testing
./test-regression.sh:
  ./mcp-debug --playback-client regression.jsonl | ./new-server > results.txt
  if ! diff expected.txt results.txt; then
    echo "Regression detected!"
    exit 1
  fi
```

### Client Development

**Problem**: Need consistent server responses for client testing.

**Solution**: Record server responses, replay for client testing.

```bash
# Record server responses  
MCP_RECORD_FILE="responses.jsonl" mcp-tui ./stable-server

# Test client against recorded responses
mcp-tui ./mcp-debug --playback-server responses.jsonl
```

## 🔧 Development

### Project Structure

```
mcp-debug/
├── main.go              # CLI entry point and mode routing
├── config/              # Configuration loading and types
├── client/              # MCP client implementations
├── integration/         # Proxy server and dynamic wrapper
├── discovery/           # Tool discovery and registration
├── proxy/               # Request forwarding handlers
├── playback/            # Recording and playback system
└── docs/                # Detailed documentation
```

### Building

```bash
# Development build
go build -o mcp-debug .

# Production build with version info
go build -ldflags "-X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.GitCommit=$(git rev-parse HEAD)" -o mcp-debug .
```

### Testing

```bash
# Unit tests
go test ./...

# Integration testing
./test-playback.sh

# Manual testing with mcp-tui  
mcp-tui ./mcp-debug --proxy --config test-empty-config.yaml
```

## 📚 Examples

### Example 1: API Server Development

```bash
# Start debugging session with recording
MCP_RECORD_FILE="api-dev.jsonl" mcp-tui ./mcp-debug --proxy --config empty-config.yaml

# Add initial API server
server_add:
  name: api
  command: go run ./api-server

# Test endpoints: api_get_users, api_create_user, api_update_user

# Fix bug in create_user, rebuild
# go build -o api-server-fixed ./api-server

# Hot-swap without losing session
server_disconnect: {name: api}
server_reconnect: {name: api, command: ./api-server-fixed}

# Test create_user immediately - bug fixed!
```

### Example 2: Multi-Service Architecture

```bash
# Add all services
server_add: {name: auth, command: ./auth-service}
server_add: {name: users, command: ./user-service}  
server_add: {name: posts, command: ./post-service}
server_add: {name: notifications, command: ./notification-service}

# Test full user flow:
# 1. auth_login
# 2. users_get_profile  
# 3. posts_create
# 4. notifications_send

# Update just the notification service
server_disconnect: {name: notifications}
server_reconnect: {name: notifications, command: ./notification-service-v2}

# All other services stay connected!
```

### Example 3: Continuous Integration

```bash
#!/bin/bash
# ci-test.sh

echo "Recording baseline..."
timeout 60 mcp-tui ./baseline-server --record baseline.jsonl < test-sequence.txt

echo "Testing PR changes..."  
./mcp-debug --playback-client baseline.jsonl | ./pr-server > pr-results.txt

echo "Comparing results..."
if diff baseline-results.txt pr-results.txt; then
  echo "✅ All tests passed"
else
  echo "❌ Regression detected"
  exit 1
fi
```

## 🐛 Troubleshooting

### Common Issues

**Proxy won't start:**
```bash
# Check config syntax
./mcp-debug --proxy --config config.yaml --log debug.log

# Review logs
tail -f /tmp/mcp-debug.log
```

**Tools not appearing:**
```bash
# Check server status
server_list

# Verify server command manually
./your-mcp-server
```

**Recording/playback issues:**
```bash
# Validate recording file
jq . session.jsonl

# Check message format
head -5 session.jsonl
```

### Debug Mode

```bash
# Enable verbose logging
MCP_DEBUG=1 ./mcp-debug --proxy --config config.yaml --log debug.log

# Monitor in real-time
tail -f debug.log
```

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Contribution Guide

```bash
# Fork and clone
git clone https://github.com/your-fork/mcp-debug
cd mcp-debug

# Create feature branch
git checkout -b feature/awesome-feature

# Make changes, test
go test ./...
./test-playback.sh

# Submit PR
git push origin feature/awesome-feature
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **[MCP Specification](https://modelcontextprotocol.io/)** - The foundation protocol
- **[mark3labs/mcp-go](https://github.com/mark3labs/mcp-go)** - Excellent Go MCP implementation
- **[mcp-tui](https://github.com/mcp-tools/mcp-tui)** - Perfect development companion tool
- **MCP Community** - For building the ecosystem

---

**🚀 Built for MCP developers, by MCP developers.**

*Make your MCP development workflow 10x faster with hot-swapping, recording, and automated testing.*