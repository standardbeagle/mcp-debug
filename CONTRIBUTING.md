# Contributing to MCP Debug

Thank you for your interest in contributing to MCP Debug! We welcome contributions from the community.

## 🚀 Quick Start

```bash
# Fork and clone
git clone https://github.com/your-username/mcp-debug
cd mcp-debug

# Install dependencies
go mod download

# Build and test
go build -o mcp-debug .
go test ./...

# Test integration
./test-playback.sh
mcp-tui ./mcp-debug --proxy --config test-empty-config.yaml
```

## 📋 How to Contribute

### Reporting Issues

- Search existing issues first
- Use the issue templates when available
- Include steps to reproduce
- Provide system information (OS, Go version, etc.)

### Submitting Changes

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes**
4. **Add tests** for new functionality
5. **Run tests**: `go test ./...`
6. **Test manually** with mcp-tui
7. **Commit with descriptive messages**
8. **Push to your fork**: `git push origin feature/amazing-feature`
9. **Create a Pull Request**

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Add comments for exported functions
- Keep functions focused and small
- Use meaningful variable names

### Testing Requirements

- Add unit tests for new functionality
- Test with `mcp-tui` for integration testing
- Update `test-playback.sh` if needed
- Ensure all existing tests pass

## 🏗️ Development Areas

We welcome contributions in these areas:

### High Priority
- Additional transport types (HTTP, WebSocket)
- Enhanced recording formats
- Performance optimizations
- Better error handling and recovery

### Medium Priority
- Additional management tools
- Configuration validation
- Documentation improvements
- Example servers and workflows

### Ideas Welcome
- Web UI for server management
- Plugin system for custom tools
- Integration with CI/CD systems
- Advanced playback features

## 🔧 Development Setup

### Prerequisites
- Go 1.24+
- mcp-tui for testing
- Git

### Project Structure
```
mcp-debug/
├── main.go              # CLI entry point
├── config/              # Configuration system
├── client/              # MCP client implementations  
├── integration/         # Proxy and dynamic wrapper
├── discovery/           # Tool discovery
├── proxy/               # Request handlers
├── playback/            # Recording and playback
└── test-servers/        # Test MCP servers
```

### Adding New Features

1. **Understand the architecture** - Read existing code
2. **Plan the change** - Discuss in issues if significant
3. **Implement incrementally** - Small, focused changes
4. **Test thoroughly** - Unit tests + manual testing
5. **Document** - Update README if needed

### Common Tasks

**Adding a new management tool:**
1. Add tool definition in `integration/dynamic_wrapper.go`
2. Implement handler function
3. Add tests
4. Update documentation

**Adding new transport:**
1. Implement client interface in `client/`
2. Update proxy server to handle new transport
3. Add configuration options
4. Test with real servers

**Improving playback:**
1. Update recording format in `playback/`
2. Enhance parser for new features
3. Update client/server playback modes
4. Test with various scenarios

## 🧪 Testing Guidelines

### Manual Testing Checklist
- [ ] Basic proxy functionality
- [ ] Server add/remove/disconnect/reconnect
- [ ] Recording session to file
- [ ] Playback client with recorded session
- [ ] Playback server with mcp-tui
- [ ] Error handling and recovery
- [ ] Hot-swap workflow

### Automated Testing
```bash
# Run all tests
go test ./...

# Test specific package
go test ./playback

# Integration test
./test-playback.sh
```

## 📝 Commit Message Guidelines

Use conventional commits:

```
feat: add WebSocket transport support
fix: handle connection errors gracefully  
docs: update README with new examples
test: add playback integration tests
refactor: simplify proxy handler logic
```

## 🤝 Community Guidelines

- Be respectful and inclusive
- Help others learn and contribute
- Focus on constructive feedback
- Celebrate contributions of all sizes

## 🎯 Release Process

1. Features merged to `main`
2. Version tagged following semver
3. Release notes generated
4. Binaries built and published

## 📞 Getting Help

- **Issues**: [GitHub Issues](https://github.com/standardbeagle/mcp-debug/issues)
- **Discussions**: [GitHub Discussions](https://github.com/standardbeagle/mcp-debug/discussions)
- **MCP Community**: [MCP Discord](https://discord.gg/mcp)

Thank you for contributing to MCP Debug! 🚀