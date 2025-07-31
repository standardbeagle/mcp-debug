package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// ProxyConfig represents the main configuration for the proxy server
type ProxyConfig struct {
	Servers []ServerConfig `yaml:"servers"`
	Proxy   ProxySettings  `yaml:"proxy"`
}

// ServerConfig represents configuration for a remote MCP server
type ServerConfig struct {
	Name      string          `yaml:"name"`
	Prefix    string          `yaml:"prefix"`
	Transport string          `yaml:"transport"`
	Command   string          `yaml:"command,omitempty"`
	Args      []string        `yaml:"args,omitempty"`
	Env       map[string]string `yaml:"env,omitempty"`
	URL       string          `yaml:"url,omitempty"`
	Auth      *AuthConfig     `yaml:"auth,omitempty"`
	Timeout   string          `yaml:"timeout,omitempty"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Type     string `yaml:"type"`
	Token    string `yaml:"token,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// ProxySettings represents proxy-level settings
type ProxySettings struct {
	HealthCheckInterval string `yaml:"healthCheckInterval"`
	ConnectionTimeout   string `yaml:"connectionTimeout"`
	MaxRetries          int    `yaml:"maxRetries"`
}

// Validate validates the configuration
func (c *ProxyConfig) Validate() error {
	// Allow empty server lists for dynamic proxies
	if len(c.Servers) == 0 {
		return nil
	}
	
	// Check for unique server names and prefixes
	names := make(map[string]bool)
	prefixes := make(map[string]bool)
	
	for i, server := range c.Servers {
		// Validate server name
		if server.Name == "" {
			return fmt.Errorf("server %d: name is required", i)
		}
		if names[server.Name] {
			return fmt.Errorf("duplicate server name: %s", server.Name)
		}
		names[server.Name] = true
		
		// Validate prefix
		if server.Prefix == "" {
			return fmt.Errorf("server %s: prefix is required", server.Name)
		}
		if prefixes[server.Prefix] {
			return fmt.Errorf("duplicate server prefix: %s", server.Prefix)
		}
		prefixes[server.Prefix] = true
		
		// Validate transport
		if server.Transport != "stdio" && server.Transport != "http" {
			return fmt.Errorf("server %s: transport must be 'stdio' or 'http'", server.Name)
		}
		
		// Validate transport-specific fields
		if server.Transport == "stdio" {
			if server.Command == "" {
				return fmt.Errorf("server %s: command is required for stdio transport", server.Name)
			}
		} else if server.Transport == "http" {
			if server.URL == "" {
				return fmt.Errorf("server %s: url is required for http transport", server.Name)
			}
		}
		
		// Validate timeout format if specified
		if server.Timeout != "" {
			if _, err := time.ParseDuration(server.Timeout); err != nil {
				return fmt.Errorf("server %s: invalid timeout format: %w", server.Name, err)
			}
		}
	}
	
	// Validate proxy settings
	if c.Proxy.HealthCheckInterval != "" {
		if _, err := time.ParseDuration(c.Proxy.HealthCheckInterval); err != nil {
			return fmt.Errorf("invalid healthCheckInterval format: %w", err)
		}
	}
	
	if c.Proxy.ConnectionTimeout != "" {
		if _, err := time.ParseDuration(c.Proxy.ConnectionTimeout); err != nil {
			return fmt.Errorf("invalid connectionTimeout format: %w", err)
		}
	}
	
	return nil
}

// ExpandEnvVars expands environment variables in configuration values
func (c *ProxyConfig) ExpandEnvVars() {
	for i := range c.Servers {
		server := &c.Servers[i]
		
		// Expand command
		server.Command = expandEnvVar(server.Command)
		
		// Expand args
		for j := range server.Args {
			server.Args[j] = expandEnvVar(server.Args[j])
		}
		
		// Expand environment variables
		for key, value := range server.Env {
			server.Env[key] = expandEnvVar(value)
		}
		
		// Expand URL
		server.URL = expandEnvVar(server.URL)
		
		// Expand auth fields
		if server.Auth != nil {
			server.Auth.Token = expandEnvVar(server.Auth.Token)
			server.Auth.Username = expandEnvVar(server.Auth.Username)
			server.Auth.Password = expandEnvVar(server.Auth.Password)
		}
	}
}

// expandEnvVar expands environment variables in the format ${VAR}
func expandEnvVar(value string) string {
	if value == "" {
		return value
	}
	
	// Simple expansion of ${VAR} format
	if strings.Contains(value, "${") {
		return os.ExpandEnv(value)
	}
	
	return value
}

// GetServerTimeout returns the timeout duration for a server, with default
func (s *ServerConfig) GetServerTimeout() time.Duration {
	if s.Timeout == "" {
		return 30 * time.Second // default timeout
	}
	
	duration, err := time.ParseDuration(s.Timeout)
	if err != nil {
		return 30 * time.Second // fallback to default
	}
	
	return duration
}

// GetProxySettings returns proxy settings with defaults
func (c *ProxyConfig) GetProxySettings() ProxySettings {
	settings := c.Proxy
	
	// Apply defaults
	if settings.HealthCheckInterval == "" {
		settings.HealthCheckInterval = "30s"
	}
	if settings.ConnectionTimeout == "" {
		settings.ConnectionTimeout = "10s"
	}
	if settings.MaxRetries == 0 {
		settings.MaxRetries = 3
	}
	
	return settings
}