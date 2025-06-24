package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rabbytesoftware/quiver/internal/config"
	"github.com/rabbytesoftware/quiver/internal/logger"
)

// CLI represents the command line interface
type CLI struct {
	config    *config.Config
	logger    *logger.Logger
	baseURL   string
	httpClient *http.Client
}

// Command represents a CLI command with its HTTP mapping
type Command struct {
	Name        string
	Description string
	Method      string
	Endpoint    string
	ParamTypes  []ParamType
	Example     string
}

// ParamType represents parameter type and validation
type ParamType struct {
	Name        string
	Type        string // "string", "int", "bool"
	Required    bool
	Position    int // 0-based position in args
	URLParam    bool // true if it's a URL parameter (like {id})
}

// New creates a new CLI instance
func New(cfg *config.Config, logger *logger.Logger) *CLI {
	baseURL := fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.Port)
	
	return &CLI{
		config:  cfg,
		logger:  logger.WithService("cli"),
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Execute executes a CLI command by parsing arguments and sending HTTP request
func (c *CLI) Execute(args []string) error {
	if len(args) == 0 {
		return c.showHelp()
	}

	commandName := args[0]
	commandArgs := args[1:]

	// Special case for help
	if commandName == "help" || commandName == "-h" || commandName == "--help" {
		return c.showHelp()
	}

	// Look up command in registry
	command, exists := CommandRegistry[commandName]
	if !exists {
		return fmt.Errorf("unknown command: %s\nUse './quiver help' to see available commands", commandName)
	}

	// Validate arguments
	if err := c.validateArgs(command, commandArgs); err != nil {
		return fmt.Errorf("invalid arguments for command '%s': %v\nExample: %s", commandName, err, command.Example)
	}

	// Build and execute HTTP request
	return c.executeHTTPRequest(command, commandArgs)
}

// validateArgs validates command arguments against parameter types
func (c *CLI) validateArgs(command Command, args []string) error {
	requiredParams := 0
	for _, param := range command.ParamTypes {
		if param.Required {
			requiredParams++
		}
	}

	if len(args) < requiredParams {
		return fmt.Errorf("expected at least %d arguments, got %d", requiredParams, len(args))
	}

	// Validate each parameter type
	for _, param := range command.ParamTypes {
		if param.Position >= len(args) {
			if param.Required {
				return fmt.Errorf("missing required parameter: %s", param.Name)
			}
			continue
		}

		value := args[param.Position]
		if err := c.validateParamType(param, value); err != nil {
			return fmt.Errorf("invalid value for parameter %s: %v", param.Name, err)
		}
	}

	return nil
}

// validateParamType validates a single parameter value
func (c *CLI) validateParamType(param ParamType, value string) error {
	switch param.Type {
	case "string":
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("string parameter cannot be empty")
		}
	case "int":
		if _, err := strconv.Atoi(value); err != nil {
			return fmt.Errorf("expected integer, got %s", value)
		}
	case "bool":
		if _, err := strconv.ParseBool(value); err != nil {
			return fmt.Errorf("expected boolean (true/false), got %s", value)
		}
	}
	return nil
}

// executeHTTPRequest builds and executes the HTTP request
func (c *CLI) executeHTTPRequest(command Command, args []string) error {
	// Build URL by replacing URL parameters
	url := c.baseURL + command.Endpoint
	var queryParams []string
	
	for _, param := range command.ParamTypes {
		if param.URLParam && param.Position < len(args) {
			// Handle URL parameters
			placeholder := "{" + param.Name + "}"
			if strings.Contains(url, "{id}") && param.Name == "package_id" {
				url = strings.Replace(url, "{id}", args[param.Position], 1)
			} else if strings.Contains(url, "{name}") && param.Name == "name" {
				url = strings.Replace(url, "{name}", args[param.Position], 1)
			} else {
				url = strings.Replace(url, placeholder, args[param.Position], 1)
			}
		} else if !param.URLParam && param.Position < len(args) {
			// Handle query parameters
			if param.Name == "query" {
				queryParams = append(queryParams, fmt.Sprintf("q=%s", args[param.Position]))
			} else {
				queryParams = append(queryParams, fmt.Sprintf("%s=%s", param.Name, args[param.Position]))
			}
		}
	}
	
	// Add query parameters to URL
	if len(queryParams) > 0 {
		url += "?" + strings.Join(queryParams, "&")
	}

	// Create request body for POST/PUT/DELETE requests
	var body io.Reader
	if command.Method == "POST" || command.Method == "PUT" || command.Method == "DELETE" {
		// For repository management commands, create appropriate JSON body
		if strings.Contains(command.Endpoint, "/repositories") && len(args) > 0 {
			bodyData := map[string]string{
				"repository": args[0],
			}
			jsonData, err := json.Marshal(bodyData)
			if err != nil {
				return fmt.Errorf("failed to marshal request body: %v", err)
			}
			body = bytes.NewReader(jsonData)
		}
	}
	
	req, err := http.NewRequest(command.Method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Execute request
	c.logger.Debug("Executing %s %s", command.Method, url)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Quiver server at %s\nMake sure the Quiver server is running with: ./quiver\nError: %v", c.baseURL, err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// Handle response
	return c.handleResponse(resp.StatusCode, respBody)
}

// handleResponse processes and displays the HTTP response
func (c *CLI) handleResponse(statusCode int, body []byte) error {
	// Pretty print JSON response
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		// If it's not JSON, just print as is
		fmt.Print(string(body))
	} else {
		fmt.Print(prettyJSON.String())
	}

	// Handle non-success status codes
	if statusCode >= 400 {
		return fmt.Errorf("\nHTTP Error %d", statusCode)
	}

	return nil
}

// showHelp displays help information
func (c *CLI) showHelp() error {
	fmt.Println("Quiver CLI - Game Server Management Platform")
	fmt.Println("\nUsage:")
	fmt.Println("  ./quiver                    Start the REST API server")
	fmt.Println("  ./quiver [command] [args]   Execute command against running server")
	fmt.Println("\nNote: Commands require a running Quiver server. Start the server first with:")
	fmt.Println("  ./quiver")
	fmt.Println("\nAvailable commands:")

	for _, command := range CommandRegistry {
		fmt.Printf("  %-15s %s\n", command.Name, command.Description)
		fmt.Printf("  %18s Example: %s\n", "", command.Example)
		fmt.Println()
	}

	fmt.Println("  help            Show this help message")
	fmt.Println("\nFor more information, visit: https://github.com/rabbytesoftware/quiver")

	return nil
} 