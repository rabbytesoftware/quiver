package execution

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages/manifest"
	"github.com/rabbytesoftware/quiver/internal/packages/types"
)

// Engine handles execution of arrow methods
type Engine struct {
	logger *logger.Logger
}

// NewEngine creates a new execution engine
func NewEngine(logger *logger.Logger) *Engine {
	return &Engine{
		logger: logger.WithService("execution-engine"),
	}
}

// ExecuteMethod executes a specific method of an arrow
func (e *Engine) ExecuteMethod(arrow manifest.ArrowInterface, methodType types.MethodType, ctx *types.ExecutionContext) error {
	methods := arrow.GetMethods()
	methodMap := methods.GetMethod(string(methodType))
	
	if methodMap == nil {
		return fmt.Errorf("method %s not found", methodType)
	}

	// Get platform-specific commands
	platform := runtime.GOOS
	commands, exists := methodMap[platform]
	if !exists {
		return fmt.Errorf("method %s not supported on platform %s", methodType, platform)
	}

	// Set up environment
	env := e.prepareEnvironment(arrow, ctx)

	// Execute commands
	for _, command := range commands {
		// Expand variables in command
		expandedCommand := e.expandVariables(command, ctx)
		
		e.logger.Debug("Executing command: %s", expandedCommand)
		
		if err := e.executeCommand(expandedCommand, ctx.InstallPath, env); err != nil {
			return fmt.Errorf("command failed: %s - %v", expandedCommand, err)
		}
	}

	return nil
}

// prepareEnvironment sets up the environment variables for command execution
func (e *Engine) prepareEnvironment(arrow manifest.ArrowInterface, ctx *types.ExecutionContext) []string {
	env := os.Environ()
	
	// Add install path
	env = append(env, fmt.Sprintf("INSTALL_PATH=%s", ctx.InstallPath))
	
	// Add arrow variables
	arrowVars := arrow.GetVariables()
	for _, variable := range arrowVars {
		varName := variable.GetName()
		var varValue string
		
		// Use user-provided value if available, otherwise use default
		if userValue, exists := ctx.Variables[varName]; exists {
			varValue = userValue
		} else if variable.GetDefault() != nil {
			varValue = fmt.Sprintf("%v", variable.GetDefault())
		}
		
		env = append(env, fmt.Sprintf("%s=%s", varName, varValue))
	}
	
	// Add any additional environment variables from context
	env = append(env, ctx.Environment...)
	
	return env
}

// expandVariables expands variables in a command string
func (e *Engine) expandVariables(command string, ctx *types.ExecutionContext) string {
	result := command
	
	// Replace INSTALL_PATH
	result = strings.ReplaceAll(result, "${INSTALL_PATH}", ctx.InstallPath)
	
	// Replace user variables
	for name, value := range ctx.Variables {
		placeholder := fmt.Sprintf("${%s}", name)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	
	return result
}

// executeCommand executes a single command
func (e *Engine) executeCommand(command, workDir string, env []string) error {
	// Split command for exec
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Env = env
	cmd.Dir = workDir
	
	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		e.logger.Error("Command failed: %s\nOutput: %s\nError: %v", command, string(output), err)
		return err
	}
	
	e.logger.Debug("Command output: %s", string(output))
	return nil
}

// ValidateEnvironment checks if the environment is suitable for execution
func (e *Engine) ValidateEnvironment(arrow manifest.ArrowInterface, installPath string) error {
	// Check if install path exists
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		return fmt.Errorf("install path does not exist: %s", installPath)
	}

	// Check platform support
	platform := runtime.GOOS
	methods := arrow.GetMethods()
	
	// Check if any method supports this platform
	supportedMethods := []types.MethodType{
		types.MethodInstall,
		types.MethodExecute,
		types.MethodUninstall,
	}
	
	hasSupport := false
	for _, methodType := range supportedMethods {
		if methodMap := methods.GetMethod(string(methodType)); methodMap != nil {
			if _, exists := methodMap[platform]; exists {
				hasSupport = true
				break
			}
		}
	}
	
	if !hasSupport {
		return fmt.Errorf("arrow does not support platform: %s", platform)
	}
	
	return nil
} 