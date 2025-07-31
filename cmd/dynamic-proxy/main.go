package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	
	"mcp-debug/client"
	"mcp-debug/config"
)

// DynamicProxy is our main proxy server
type DynamicProxy struct {
	server        *mcp.Server
	clients       map[string]client.MCPClient
	tools         map[string]string // tool name -> server name
	mu            sync.RWMutex
}

func NewDynamicProxy() *DynamicProxy {
	return &DynamicProxy{
		clients: make(map[string]client.MCPClient),
		tools:   make(map[string]string),
	}
}

func (p *DynamicProxy) Initialize() error {
	// Create stdio transport
	transport := stdio.NewStdioServerTransport()
	
	// Create MCP server
	p.server = mcp.NewServer(
		transport,
		mcp.WithName("Dynamic MCP Proxy"),
		mcp.WithVersion("1.0.0"),
	)
	
	// Register management tools
	p.registerManagementTools()
	
	return nil
}

func (p *DynamicProxy) registerManagementTools() error {
	// server_add tool
	type AddServerArgs struct {
		Name    string `json:"name" jsonschema:"required,description=Name/prefix for the server"`
		Command string `json:"command" jsonschema:"required,description=Command to run (e.g. 'npx -y @modelcontextprotocol/filesystem /path')"`
	}
	
	err := p.server.RegisterTool("server_add", "Add a new MCP server to the proxy", 
		func(args AddServerArgs) (*mcp.ToolResponse, error) {
			return p.handleServerAdd(args.Name, args.Command)
		})
	if err != nil {
		return fmt.Errorf("failed to register server_add: %w", err)
	}
	
	// server_remove tool
	type RemoveServerArgs struct {
		Name string `json:"name" jsonschema:"required,description=Name of the server to remove"`
	}
	
	err = p.server.RegisterTool("server_remove", "Remove an MCP server from the proxy",
		func(args RemoveServerArgs) (*mcp.ToolResponse, error) {
			return p.handleServerRemove(args.Name)
		})
	if err != nil {
		return fmt.Errorf("failed to register server_remove: %w", err)
	}
	
	// server_list tool
	err = p.server.RegisterTool("server_list", "List all connected MCP servers",
		func() (*mcp.ToolResponse, error) {
			return p.handleServerList()
		})
	if err != nil {
		return fmt.Errorf("failed to register server_list: %w", err)
	}
	
	return nil
}

func (p *DynamicProxy) handleServerAdd(name, command string) (*mcp.ToolResponse, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// Check if already exists
	if _, exists := p.clients[name]; exists {
		return nil, fmt.Errorf("server '%s' already exists", name)
	}
	
	// Parse command
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid command")
	}
	
	// Create client
	client := client.NewStdioClient(name, parts[0], parts[1:])
	
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	
	if _, err := client.Initialize(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to initialize: %w", err)
	}
	
	// List tools
	tools, err := client.ListTools(ctx)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}
	
	// Store client
	p.clients[name] = client
	
	// Register each tool
	registeredCount := 0
	for _, tool := range tools {
		prefixedName := fmt.Sprintf("%s_%s", name, tool.Name)
		
		// Create a closure that captures the tool name
		originalName := tool.Name
		serverName := name
		
		err := p.server.RegisterTool(prefixedName, 
			fmt.Sprintf("[%s] %s", name, tool.Description),
			func(args map[string]interface{}) (*mcp.ToolResponse, error) {
				return p.callRemoteTool(serverName, originalName, args)
			})
		
		if err != nil {
			log.Printf("Failed to register tool %s: %v", prefixedName, err)
			continue
		}
		
		p.tools[prefixedName] = name
		registeredCount++
		log.Printf("Registered tool: %s", prefixedName)
	}
	
	result := fmt.Sprintf("Added server '%s' with %d tools", name, registeredCount)
	return mcp.NewToolResponse(mcp.NewTextContent(result)), nil
}

func (p *DynamicProxy) handleServerRemove(name string) (*mcp.ToolResponse, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	client, exists := p.clients[name]
	if !exists {
		return nil, fmt.Errorf("server '%s' not found", name)
	}
	
	// Remove tools
	removedCount := 0
	for toolName, serverName := range p.tools {
		if serverName == name {
			if err := p.server.DeregisterTool(toolName); err != nil {
				log.Printf("Failed to deregister tool %s: %v", toolName, err)
			} else {
				removedCount++
			}
			delete(p.tools, toolName)
		}
	}
	
	// Close client
	client.Close()
	delete(p.clients, name)
	
	result := fmt.Sprintf("Removed server '%s' and %d tools", name, removedCount)
	return mcp.NewToolResponse(mcp.NewTextContent(result)), nil
}

func (p *DynamicProxy) handleServerList() (*mcp.ToolResponse, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	var result strings.Builder
	result.WriteString("Connected MCP Servers:\n")
	result.WriteString("=====================\n\n")
	
	if len(p.clients) == 0 {
		result.WriteString("No servers connected.\n")
	} else {
		for name, client := range p.clients {
			// Count tools
			toolCount := 0
			for _, serverName := range p.tools {
				if serverName == name {
					toolCount++
				}
			}
			
			status := "connected"
			if !client.IsConnected() {
				status = "disconnected"
			}
			
			result.WriteString(fmt.Sprintf("- %s [%s] - %d tools\n", name, status, toolCount))
		}
	}
	
	result.WriteString(fmt.Sprintf("\nTotal servers: %d\n", len(p.clients)))
	return mcp.NewToolResponse(mcp.NewTextContent(result.String())), nil
}

func (p *DynamicProxy) callRemoteTool(serverName, toolName string, args map[string]interface{}) (*mcp.ToolResponse, error) {
	p.mu.RLock()
	client, exists := p.clients[serverName]
	p.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("server '%s' not found", serverName)
	}
	
	// Call remote tool
	result, err := client.CallTool(context.Background(), toolName, args)
	if err != nil {
		return nil, fmt.Errorf("remote call failed: %w", err)
	}
	
	// Convert result
	if result.IsError {
		if len(result.Content) > 0 {
			return nil, fmt.Errorf("remote error: %s", result.Content[0].Text)
		}
		return nil, fmt.Errorf("remote error: unknown")
	}
	
	// Convert content items
	var contents []*mcp.Content
	for _, item := range result.Content {
		if item.Type == "text" {
			contents = append(contents, mcp.NewTextContent(item.Text))
		}
	}
	
	if len(contents) == 0 {
		contents = append(contents, mcp.NewTextContent("No content"))
	}
	
	return mcp.NewToolResponse(contents...), nil
}

func (p *DynamicProxy) Serve() error {
	return p.server.Serve()
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.Parse()
	
	// Set up logging
	log.SetOutput(os.Stderr)
	
	proxy := NewDynamicProxy()
	
	if err := proxy.Initialize(); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	
	// If config provided, load initial servers
	if configPath != "" {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Printf("Warning: Failed to load config: %v", err)
		} else {
			// Add servers from config
			for _, server := range cfg.Servers {
				cmd := server.Command
				if len(server.Args) > 0 {
					cmd = fmt.Sprintf("%s %s", cmd, strings.Join(server.Args, " "))
				}
				
				log.Printf("Adding server from config: %s", server.Name)
				if _, err := proxy.handleServerAdd(server.Name, cmd); err != nil {
					log.Printf("Failed to add server %s: %v", server.Name, err)
				}
			}
		}
	}
	
	fmt.Fprintf(os.Stderr, "Starting Dynamic MCP Proxy Server...\n")
	if err := proxy.Serve(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}