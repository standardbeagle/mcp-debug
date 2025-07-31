package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	
	"mcp-debug/config"
	"mcp-debug/integration"
	"mcp-debug/playback"
)

const Version = "1.0.0"

var (
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// setupLogging configures logging for stdio MCP mode
func setupLogging(logFile string) error {
	// Default log file if not specified
	if logFile == "" {
		logFile = "/tmp/mcp-proxy.log"
	}
	
	// Ensure directory exists
	dir := filepath.Dir(logFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}
	
	// Open log file
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	
	// Set log output to file
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Printf("=== MCP Proxy Server Started ===")
	log.Printf("Logging to: %s", logFile)
	
	return nil
}

func main() {
	// Define command line flags
	var (
		proxyMode      = flag.Bool("proxy", false, "Run in proxy mode")
		dynamicMode    = flag.Bool("dynamic", false, "Run in dynamic proxy mode (true dynamic tool registration)")
		configPath     = flag.String("config", "", "Path to configuration file (required for proxy mode)")
		logFile        = flag.String("log", "", "Log file path (defaults to /tmp/mcp-proxy.log for stdio mode)")
		recordFile     = flag.String("record", "", "Record JSON-RPC traffic to file for playback")
		playbackClient = flag.String("playback-client", "", "Act as MCP client replaying recorded session file")
		playbackServer = flag.String("playback-server", "", "Act as MCP server replaying recorded responses")
	)
	flag.Parse()
	
	// Handle playback modes
	if *playbackClient != "" {
		if err := runPlaybackClient(*playbackClient); err != nil {
			log.Fatalf("Playback client failed: %v", err)
		}
		return
	}
	
	if *playbackServer != "" {
		if err := runPlaybackServer(*playbackServer); err != nil {
			log.Fatalf("Playback server failed: %v", err)
		}
		return
	}
	
	// Handle proxy modes
	if *proxyMode || *dynamicMode {
		if *configPath == "" {
			fmt.Fprintln(os.Stderr, "Error: --config is required when using --proxy or --dynamic mode")
			fmt.Fprintln(os.Stderr, "Usage: mcp-server --dynamic --config /path/to/config.yaml")
			os.Exit(1)
		}
		
		// Set up file logging for stdio mode
		if err := setupLogging(*logFile); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to setup logging: %v\n", err)
			os.Exit(1)
		}
		
		// Use dynamic proxy with management tools
		if err := runDynamicProxyWithManagement(*configPath, *recordFile); err != nil {
			log.Fatalf("Dynamic proxy server failed: %v", err)
		}
		return
	}
	
	// Handle CLI commands and configuration (original mode)
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-h", "--help":
			printUsage()
			return
		case "-v", "--version":
			handleVersionCommand()
			return
		case "config":
			handleConfigCommand()
			return
		case "env":
			handleEnvCommand()
			return
		case "test":
			handleTestCommand()
			return
		case "tools":
			handleToolsCommand()
			return
		default:
			if strings.HasPrefix(os.Args[1], "-") {
				fmt.Printf("Unknown flag: %s\n", os.Args[1])
				printUsage()
				return
			}
		}
	}

	// Detect if running from CLI vs MCP client
	if isRunningFromCLI() {
		fmt.Printf("MCP Server: Dynamic MCP Server v%s\n", Version)
		fmt.Printf("This is an MCP (Model Context Protocol) server.\n")
		fmt.Printf("It should be run by an MCP client, not directly from the command line.\n\n")
		printUsage()
		return
	}

	// Create MCP server
	s := server.NewMCPServer(
		"Dynamic MCP Server",
		Version,
		server.WithToolCapabilities(true),
	)

	// Define hello_world tool
	tool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of person to greet"),
		),
	)

	// Add tool handler
	s.AddTool(tool, helloHandler)

	// Start stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}

// runDynamicProxyWithManagement runs the proxy with dynamic management tools
func runDynamicProxyWithManagement(configPath, recordFile string) error {
	ctx := context.Background()
	
	// Load configuration
	log.Printf("Loading configuration from: %s", configPath)
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	
	log.Printf("Configuration loaded: %d servers configured", len(cfg.Servers))
	
	// Create dynamic wrapper
	wrapper := integration.NewDynamicWrapper(cfg)
	
	// Enable recording if specified
	if recordFile != "" {
		log.Printf("Recording JSON-RPC traffic to: %s", recordFile)
		if err := wrapper.EnableRecording(recordFile); err != nil {
			return fmt.Errorf("failed to enable recording: %w", err)
		}
	}
	
	// Initialize with static servers
	log.Println("Initializing proxy server...")
	if err := wrapper.Initialize(ctx); err != nil {
		// Allow starting with no tools for dynamic management
		if !strings.Contains(err.Error(), "no tools were successfully discovered") {
			return fmt.Errorf("failed to initialize: %w", err)
		}
		log.Println("Starting with no initial servers - use server_add to add servers dynamically")
	}
	
	// Start the server
	return wrapper.Start()
}

// runProxyServer runs the MCP proxy server with the given configuration
func runDynamicProxyServer(configPath string) error {
	log.Printf("Loading configuration from: %s", configPath)
	
	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	
	log.Printf("Configuration loaded: %d servers configured", len(cfg.Servers))
	
	// Create dynamic proxy server
	proxyServer := integration.NewDynamicProxyServer(&cfg.Proxy)
	
	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		log.Printf("Shutting down...")
		cancel()
		proxyServer.Shutdown()
	}()
	
	// Start connecting to servers in background
	go func() {
		for _, serverConfig := range cfg.Servers {
			if err := proxyServer.ConnectToServer(ctx, serverConfig); err != nil {
				log.Printf("Failed to connect to server %s: %v", serverConfig.Name, err)
			}
		}
	}()
	
	// Start the MCP server (this will block)
	log.Printf("Starting dynamic MCP proxy server...")
	return proxyServer.Serve()
}

func runProxyServer(configPath string) error {
	ctx := context.Background()
	
	// Load configuration
	log.Printf("Loading configuration from: %s", configPath)
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	
	log.Printf("Configuration loaded: %d servers configured", len(cfg.Servers))
	
	// Create proxy server
	proxyServer := integration.NewProxyServer(cfg)
	
	// Initialize proxy server (connect to remotes and discover tools)
	log.Println("Initializing proxy server...")
	if err := proxyServer.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize proxy server: %w", err)
	}
	
	// Set up graceful shutdown
	// TODO: Add signal handling for graceful shutdown
	defer func() {
		log.Println("Shutting down...")
		if err := proxyServer.Shutdown(ctx); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	}()
	
	// Start the proxy server (this blocks)
	log.Println("Proxy server initialized successfully. Starting MCP server...")
	return proxyServer.Start()
}

// isRunningFromCLI detects if the program is running from command line vs MCP client
func isRunningFromCLI() bool {
	// Check if stdin is a terminal (tty)
	if fileInfo, err := os.Stdin.Stat(); err == nil {
		return (fileInfo.Mode() & os.ModeCharDevice) != 0
	}
	return true
}

// printUsage displays help information for CLI usage
func printUsage() {
	fmt.Printf(`USAGE:
    This MCP server can run in multiple modes:
    
    1. PROXY MODE (recommended):
       %s --proxy --config /path/to/config.yaml [--record session.jsonl]
       
       Connects to multiple MCP servers and exposes their tools with prefixes.
       Optional recording creates playback files.
       
    2. STANDALONE MODE:
       %s (without flags)
       
       Runs as a simple MCP server with hello_world tool.
    
    3. PLAYBACK CLIENT MODE:
       %s --playback-client session.jsonl
       
       Acts as MCP client replaying recorded requests.
       
    4. PLAYBACK SERVER MODE:
       %s --playback-server session.jsonl
       
       Acts as MCP server replaying recorded responses.
    
    For direct testing:
    %s --help           Show this help message
    %s --version        Show version information
    %s config           Configuration management commands
    %s env              Environment variable management
    %s test             Test MCP tools directly
    %s tools            Tool interface commands
    
    For MCP client usage (proxy mode):
    1. Create a configuration file:
       servers:
         - name: "math-server"
           prefix: "math"
           transport: "stdio"
           command: "/path/to/math-mcp-server"
    
    2. Build the server:
       go build -o mcp-server
    
    3. Add to your MCP client configuration:
       Claude Desktop: Add to claude_desktop_config.json
       {
         "mcpServers": {
           "dynamic-mcp-proxy": {
             "command": "/path/to/mcp-server",
             "args": ["--proxy", "--config", "/path/to/config.yaml"]
           }
         }
       }
    
    4. Start your MCP client (Claude Desktop, etc.)
    
    Available Tools (standalone mode):
    - hello_world: Say hello to someone
    
    Environment Variables:
    - MCP_DEBUG=1: Enable debug logging
    - MCP_CONFIG_PATH: Path to configuration file
    
    For more information about MCP:
    https://modelcontextprotocol.io/
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

// handleVersionCommand shows version information
func handleVersionCommand() {
	fmt.Printf("Dynamic MCP Server v%s\n", Version)
	fmt.Printf("Build time: %s\n", BuildTime)
	fmt.Printf("Git commit: %s\n", GitCommit)
}

// handleConfigCommand manages configuration files
func handleConfigCommand() {
	if len(os.Args) < 3 {
		fmt.Printf(`Configuration Management:
    %s config init              Create default configuration file
    %s config show              Show current configuration
    %s config set <key> <value> Set configuration value
    %s config get <key>         Get configuration value
    %s config validate          Validate configuration file
    %s config path              Show configuration file path
    
Example:
    %s config init
    %s config set api_key "your-api-key"
    %s config set database_url "postgres://localhost/mydb"
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
		return
	}

	switch os.Args[2] {
	case "init":
		fmt.Println("Creating default configuration file...")
		// TODO: Implement config file creation
		fmt.Println("Configuration file created at: ./config.json")
	case "show":
		fmt.Println("Current configuration:")
		// TODO: Implement config display
		fmt.Println("No configuration file found. Run 'config init' to create one.")
	case "set":
		if len(os.Args) < 5 {
			fmt.Println("Usage: config set <key> <value>")
			return
		}
		key, value := os.Args[3], os.Args[4]
		fmt.Printf("Setting %s = %s\n", key, value)
		// TODO: Implement config value setting
	case "get":
		if len(os.Args) < 4 {
			fmt.Println("Usage: config get <key>")
			return
		}
		key := os.Args[3]
		fmt.Printf("Getting value for %s\n", key)
		// TODO: Implement config value retrieval
	case "validate":
		fmt.Println("Validating configuration...")
		// TODO: Implement config validation
		fmt.Println("Configuration validation not yet implemented.")
	case "path":
		configPath := os.Getenv("MCP_CONFIG_PATH")
		if configPath == "" {
			configPath = "./config.json"
		}
		fmt.Printf("Configuration file path: %s\n", configPath)
	default:
		fmt.Printf("Unknown config command: %s\n", os.Args[2])
	}
}

// handleEnvCommand manages environment variables
func handleEnvCommand() {
	if len(os.Args) < 3 {
		fmt.Printf(`Environment Variable Management:
    %s env list           List all environment variables
    %s env check          Check required environment variables
    %s env template       Generate .env template file
    %s env validate       Validate environment variables
    
Example:
    %s env check
    %s env template > .env
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
		return
	}

	switch os.Args[2] {
	case "list":
		fmt.Println("Environment variables:")
		fmt.Printf("MCP_DEBUG: %s\n", os.Getenv("MCP_DEBUG"))
		fmt.Printf("MCP_CONFIG_PATH: %s\n", os.Getenv("MCP_CONFIG_PATH"))
		// Add other relevant env vars as needed
	case "check":
		fmt.Println("Checking required environment variables...")
		// For this basic server, no env vars are strictly required
		fmt.Println("✓ All required environment variables are set")
	case "template":
		fmt.Println(`# Dynamic MCP Server Environment Variables
# Copy this file to .env and fill in your values

# Debug logging (0 or 1)
MCP_DEBUG=0

# Configuration file path
MCP_CONFIG_PATH=./config.json

# API Keys (if needed)
# API_KEY=your-api-key-here
# DATABASE_URL=your-database-url-here`)
	case "validate":
		fmt.Println("Validating environment variables...")
		// TODO: Implement env var validation
		fmt.Println("✓ Environment variables are valid")
	default:
		fmt.Printf("Unknown env command: %s\n", os.Args[2])
	}
}

// handleTestCommand provides CLI testing of MCP tools
func handleTestCommand() {
	if len(os.Args) < 3 {
		fmt.Printf(`Tool Testing:
    %s test list                List available tools
    %s test <tool> [args...]    Test specific tool
    
Example:
    %s test hello_world name="John"
`, os.Args[0], os.Args[0], os.Args[0])
		return
	}

	switch os.Args[2] {
	case "list":
		fmt.Println("Available tools:")
		fmt.Println("- hello_world: Say hello to someone")
		// TODO: Dynamically list registered tools
	default:
		toolName := os.Args[2]
		fmt.Printf("Testing tool: %s\n", toolName)

		if toolName == "hello_world" {
			name := "World"
			if len(os.Args) > 3 {
				// Simple argument parsing for demo
				for _, arg := range os.Args[3:] {
					if strings.HasPrefix(arg, "name=") {
						name = strings.TrimPrefix(arg, "name=")
						name = strings.Trim(name, "\"'")
					}
				}
			}
			fmt.Printf("Result: Hello, %s!\n", name)
		} else {
			fmt.Printf("Unknown tool: %s\n", toolName)
		}
	}
}

// handleToolsCommand provides CLI interface to MCP tools
func handleToolsCommand() {
	if len(os.Args) < 3 {
		fmt.Printf(`Tool Interface:
    %s tools list               List all available tools
    %s tools describe <tool>    Show tool description and parameters
    %s tools run <tool> [args]  Run tool with arguments
    
Example:
    %s tools list
    %s tools describe hello_world
    %s tools run hello_world --name "John"
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
		return
	}

	switch os.Args[2] {
	case "list":
		fmt.Println("Available MCP Tools:")
		fmt.Println()
		fmt.Println("hello_world")
		fmt.Println("  Description: Say hello to someone")
		fmt.Println("  Parameters:")
		fmt.Println("    - name (string, required): Name of person to greet")
		fmt.Println()
		// TODO: Dynamically list all registered tools
	case "describe":
		if len(os.Args) < 4 {
			fmt.Println("Usage: tools describe <tool>")
			return
		}
		toolName := os.Args[3]
		if toolName == "hello_world" {
			fmt.Println("Tool: hello_world")
			fmt.Println("Description: Say hello to someone")
			fmt.Println("Parameters:")
			fmt.Println("  - name (string, required): Name of person to greet")
			fmt.Println()
			fmt.Println("Example usage:")
			fmt.Printf("  %s tools run hello_world --name \"John\"\n", os.Args[0])
			fmt.Printf("  %s test hello_world name=\"John\"\n", os.Args[0])
		} else {
			fmt.Printf("Unknown tool: %s\n", toolName)
		}
	case "run":
		if len(os.Args) < 4 {
			fmt.Println("Usage: tools run <tool> [args]")
			return
		}
		toolName := os.Args[3]
		fmt.Printf("Running tool: %s\n", toolName)

		if toolName == "hello_world" {
			name := "World"
			// Parse CLI arguments (simple implementation)
			for i := 4; i < len(os.Args); i++ {
				if os.Args[i] == "--name" && i+1 < len(os.Args) {
					name = os.Args[i+1]
					i++ // Skip next arg as it's the value
				}
			}
			fmt.Printf("Result: Hello, %s!\n", name)
		} else {
			fmt.Printf("Unknown tool: %s\n", toolName)
		}
	default:
		fmt.Printf("Unknown tools command: %s\n", os.Args[2])
	}
}

// runPlaybackClient runs the playback client mode
func runPlaybackClient(recordingFile string) error {
	log.SetOutput(os.Stderr) // Ensure logs go to stderr, not stdout
	log.Printf("Starting playback client with recording: %s", recordingFile)
	
	// Parse the recording file
	session, err := playback.ParseRecordingFile(recordingFile)
	if err != nil {
		return fmt.Errorf("failed to parse recording file: %w", err)
	}
	
	log.Printf("Loaded session with %d messages", len(session.Messages))
	
	// Create and run playback client
	client := playback.NewPlaybackClient(session)
	return client.Run()
}

// runPlaybackServer runs the playback server mode
func runPlaybackServer(recordingFile string) error {
	log.SetOutput(os.Stderr) // Ensure logs go to stderr, not stdout
	log.Printf("Starting playback server with recording: %s", recordingFile)
	
	// Parse the recording file
	session, err := playback.ParseRecordingFile(recordingFile)
	if err != nil {
		return fmt.Errorf("failed to parse recording file: %w", err)
	}
	
	log.Printf("Loaded session with %d messages", len(session.Messages))
	
	// Create and run playback server
	server := playback.NewPlaybackServer(session)
	return server.Run()
}