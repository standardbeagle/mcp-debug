# Dynamic MCP Server

A Model Context Protocol (MCP) server implementation in Go with stdio transport and comprehensive CLI tooling.

## Features

- MCP server with stdio transport for Claude Desktop integration
- CLI detection with helpful usage information
- Built-in tool testing and configuration management
- Hello World example tool implementation
- Environment variable management
- Dual-mode operation (MCP server and CLI tool)

## Installation

```bash
# Clone or navigate to this directory
cd dynamic-mcp

# Install dependencies
go mod download

# Build the server
go build -o mcp-server

# Make it executable (Unix/Linux/macOS)
chmod +x mcp-server
```

## Usage

### CLI Mode

When run directly from the command line, the server detects it's not being called by an MCP client and provides helpful CLI tools:

```bash
# Show help
./mcp-server --help

# Show version
./mcp-server --version

# Test the hello_world tool
./mcp-server test hello_world name="Alice"

# List all available tools
./mcp-server tools list

# Describe a specific tool
./mcp-server tools describe hello_world

# Run tools with CLI interface
./mcp-server tools run hello_world --name "Bob"
```

### Configuration Management

```bash
# Initialize configuration
./mcp-server config init

# Show configuration file path
./mcp-server config path

# Set configuration values
./mcp-server config set api_key "your-api-key"
./mcp-server config set database_url "postgres://localhost/mydb"

# Show current configuration
./mcp-server config show
```

### Environment Variables

```bash
# Generate .env template
./mcp-server env template > .env

# List environment variables
./mcp-server env list

# Check required environment variables
./mcp-server env check

# Validate environment
./mcp-server env validate
```

### MCP Server Mode

To use with Claude Desktop or other MCP clients:

1. Build the server:
   ```bash
   go build -o mcp-server
   ```

2. Add to Claude Desktop configuration (`claude_desktop_config.json`):
   ```json
   {
     "mcpServers": {
       "dynamic-mcp": {
         "command": "/full/path/to/mcp-server"
       }
     }
   }
   ```

3. Restart Claude Desktop

4. The server's tools will be available in your Claude conversations

## Available Tools

### hello_world
Say hello to someone.

**Parameters:**
- `name` (string, required): Name of person to greet

**Example:**
```
Tool: hello_world
Arguments: {"name": "World"}
Result: "Hello, World!"
```

## Environment Variables

- `MCP_DEBUG`: Set to `1` to enable debug logging
- `MCP_CONFIG_PATH`: Path to configuration file (default: `./config.json`)

## Development

### Adding New Tools

To add a new tool to the server:

1. Define the tool schema:
   ```go
   tool := mcp.NewTool("tool_name",
       mcp.WithDescription("Tool description"),
       mcp.WithString("param_name",
           mcp.Required(),
           mcp.Description("Parameter description"),
       ),
   )
   ```

2. Create a handler function:
   ```go
   func toolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
       // Implementation
   }
   ```

3. Register the tool:
   ```go
   s.AddTool(tool, toolHandler)
   ```

4. Update CLI commands to include the new tool in listings and test commands

### Testing

Test your tools directly using the CLI:

```bash
# Test with simple syntax
./mcp-server test tool_name param="value"

# Test with CLI syntax
./mcp-server tools run tool_name --param "value"
```

### Building with Version Information

```bash
go build -ldflags "-X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.GitCommit=$(git rev-parse HEAD)" -o mcp-server
```

## MCP Inspector Testing

You can test the server using MCP Inspector:

```bash
# Install MCP Inspector
npm install -g @modelcontextprotocol/inspector

# Test the server
mcp-inspector ./mcp-server
```

## Security Notes

- Store sensitive configuration in environment variables
- Use secure file permissions (600) for config files
- Validate all inputs to prevent injection attacks
- Never commit secrets to version control

## Troubleshooting

### Server won't start in Claude Desktop
- Check that the path in `claude_desktop_config.json` is absolute
- Ensure the binary is executable: `chmod +x mcp-server`
- Check logs for any error messages

### CLI detection not working
- The server detects CLI usage by checking if stdin is a TTY
- If running through a script, it may think it's being called by an MCP client
- Use explicit CLI commands to force CLI mode

## Resources

- [Model Context Protocol Documentation](https://modelcontextprotocol.io/)
- [MCP Go Library (mark3labs)](https://github.com/mark3labs/mcp-go)
- [Claude Desktop MCP Integration](https://docs.anthropic.com/en/docs/claude-code/mcp)

## License

This project is provided as-is for educational and development purposes.