package execution

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages/execution/process"
	"github.com/rabbytesoftware/quiver/internal/packages/manifest"
	"github.com/rabbytesoftware/quiver/internal/packages/types"
)

// Engine handles execution of arrow methods
type Engine struct {
	logger            *logger.Logger
	tempDir           string
	httpClient        *http.Client
	netbridgeProcessor *NetbridgeProcessor
	processTracker    *process.ProcessTracker
}

// ExecutionOptions provides options for command execution
type ExecutionOptions struct {
	Timeout     time.Duration
	DryRun      bool
	MethodType  types.MethodType  // Type of method being executed
	ArrowName   string            // Name of the arrow being executed
}

// NewEngine creates a new execution engine
func NewEngine(logger *logger.Logger) *Engine {
	tempDir := filepath.Join(os.TempDir(), "quiver-exec")
	os.MkdirAll(tempDir, 0755)
	
	return &Engine{
		logger:             logger.WithService("execution-engine"),
		tempDir:            tempDir,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		netbridgeProcessor: NewNetbridgeProcessor(logger),
		processTracker:    process.NewProcessTracker(logger),
	}
}

// ExecuteMethod executes a specific method of an arrow
func (e *Engine) ExecuteMethod(arrow manifest.ArrowInterface, methodType types.MethodType, ctx *types.ExecutionContext) error {
	return e.ExecuteMethodWithOptions(arrow, methodType, ctx, &ExecutionOptions{
		Timeout:    5 * time.Minute,
		MethodType: methodType,
		ArrowName:  ctx.ArrowName,
	})
}

// ExecuteMethodWithOptions executes a specific method with custom options
func (e *Engine) ExecuteMethodWithOptions(arrow manifest.ArrowInterface, methodType types.MethodType, ctx *types.ExecutionContext, opts *ExecutionOptions) error {
	e.logger.Info("Starting execution of method %s for arrow %s", methodType, arrow.Name())
	
	// Process netbridge variables before method execution
	netbridgeResults, err := e.netbridgeProcessor.ProcessVariables(arrow, ctx)
	if err != nil {
		e.logger.Error("Failed to process netbridge variables: %v", err)
		return fmt.Errorf("netbridge variable processing failed: %v", err)
	}
	
	// Log netbridge results for method initialization status
	e.netbridgeProcessor.LogProcessingResults(netbridgeResults)
	
	// Validate pre-execution requirements
	if err := e.ValidateRequirements(arrow); err != nil {
		return fmt.Errorf("requirements validation failed: %v", err)
	}

	methods := arrow.GetMethods()
	methodMap := methods.GetMethod(string(methodType))
	
	if methodMap == nil {
		return fmt.Errorf("method %s not found", methodType)
	}

	// Get platform-specific commands with improved architecture detection
	osStr := runtime.GOOS
	archStr := runtime.GOARCH
	
	e.logger.Info("Detecting platform: OS=%s, Arch=%s", osStr, archStr)
	
	osMap, exists := methodMap[osStr]
	if !exists {
		return fmt.Errorf("method %s not supported on platform %s", methodType, osStr)
	}
	
	commands, exists := osMap[archStr]
	if !exists {
		// Try to find a compatible architecture
		supportedArchs := arrow.GetSupportedArchs(osStr)
		e.logger.Debug("Architecture %s not found, supported architectures: %v", archStr, supportedArchs)
		
		if len(supportedArchs) == 0 {
			return fmt.Errorf("method %s not supported on architecture %s for %s", methodType, archStr, osStr)
		}
		
		// Use the first available architecture as fallback
		commands, exists = osMap[supportedArchs[0]]
		if !exists {
			return fmt.Errorf("no compatible architecture found for method %s on %s", methodType, osStr)
		}
		
		e.logger.Warn("Using fallback architecture %s instead of %s", supportedArchs[0], archStr)
	}

	e.logger.Info("Found %d commands to execute for method %s", len(commands), methodType)

	// Set up environment
	env := e.prepareEnvironment(arrow, ctx)

	// Create execution context with timeout
	execCtx := context.Background()
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		execCtx, cancel = context.WithTimeout(execCtx, opts.Timeout)
		defer cancel()
	}

	// Execute commands
	for i, command := range commands {
		// Expand variables in command
		expandedCommand := e.expandVariables(command, ctx)
		
		if opts.DryRun {
			e.logger.Info("DRY RUN: Would execute command %d/%d: %s", i+1, len(commands), expandedCommand)
			continue
		}
		
		e.logger.Info("Executing command %d/%d: %s", i+1, len(commands), expandedCommand)
		
		if err := e.executeCommand(execCtx, expandedCommand, ctx.InstallPath, env, opts); err != nil {
			e.logger.Error("Command %d/%d failed: %s", i+1, len(commands), expandedCommand)
			return fmt.Errorf("command failed: %s - %v", expandedCommand, err)
		}
		
		e.logger.Info("Command %d/%d completed successfully", i+1, len(commands))
	}

	e.logger.Info("Successfully completed execution of method %s for arrow %s", methodType, arrow.Name())
	return nil
}

// prepareEnvironment sets up the environment variables for command execution
func (e *Engine) prepareEnvironment(arrow manifest.ArrowInterface, ctx *types.ExecutionContext) []string {
	env := os.Environ()
	
	// Add install path (support both INSTALL_PATH and INSTALL_DIR for compatibility)
	env = append(env, fmt.Sprintf("INSTALL_PATH=%s", ctx.InstallPath))
	env = append(env, fmt.Sprintf("INSTALL_DIR=%s", ctx.InstallPath))
	e.logger.Debug("Environment: INSTALL_PATH=%s", ctx.InstallPath)
	e.logger.Debug("Environment: INSTALL_DIR=%s", ctx.InstallPath)
	
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
		
		// Log non-sensitive variables
		if variable.GetSensitive() {
			e.logger.Debug("Environment: %s=***REDACTED***", varName)
		} else {
			e.logger.Debug("Environment: %s=%s", varName, varValue)
		}
	}
	
	// Add any additional environment variables from context
	env = append(env, ctx.Environment...)
	
	return env
}

// expandVariables expands variables in a command string
func (e *Engine) expandVariables(command string, ctx *types.ExecutionContext) string {
	result := command
	
	// Replace INSTALL_PATH and INSTALL_DIR (both point to the same path for compatibility)
	result = strings.ReplaceAll(result, "${INSTALL_PATH}", ctx.InstallPath)
	result = strings.ReplaceAll(result, "${INSTALL_DIR}", ctx.InstallPath)
	
	// Replace user variables
	for name, value := range ctx.Variables {
		placeholder := fmt.Sprintf("${%s}", name)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	
	return result
}

// executeCommand executes a single command with special command interpretation
func (e *Engine) executeCommand(ctx context.Context, command, workDir string, env []string, opts *ExecutionOptions) error {
	// Parse and handle special commands
	if handled, err := e.handleSpecialCommand(ctx, command, workDir, env); handled {
		return err
	}
	
	// Fall back to shell execution for regular commands
	return e.executeShellCommand(ctx, command, workDir, env, opts)
}

// executeShellCommand executes regular shell commands with verbose logging
func (e *Engine) executeShellCommand(ctx context.Context, command, workDir string, env []string, opts *ExecutionOptions) error {
	// Split command for exec
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}

	e.logger.Info("Executing shell command: %s", command)
	e.logger.Debug("Working directory: %s", workDir)
	e.logger.Debug("Command parts: %v", parts)

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Env = env
	cmd.Dir = workDir
	
	// Create pipes for stdout and stderr to capture and log output in real-time
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}
	
	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}
	
	// Track process if this is an execute method (likely to be long-running)
	if opts.MethodType == types.MethodExecute && opts.ArrowName != "" {
		e.processTracker.TrackProcess(opts.ArrowName, cmd.Process, command)
		e.logger.Debug("Tracking process PID %d for arrow %s", cmd.Process.Pid, opts.ArrowName)
	}
	
	// Create channels to handle output
	done := make(chan bool, 2)
	
	// Handle stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			e.logger.Info("[STDOUT] %s", line)
		}
		done <- true
	}()
	
	// Handle stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			e.logger.Warn("[STDERR] %s", line)
		}
		done <- true
	}()
	
	// Wait for output handling to complete
	<-done
	<-done
	
	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		e.logger.Error("Command failed with error: %v", err)
		return fmt.Errorf("command execution failed: %v", err)
	}
	
	e.logger.Info("Shell command completed successfully")
	return nil
}

// GetNetbridgeResults returns netbridge processing results for API reporting
func (e *Engine) GetNetbridgeResults(arrow manifest.ArrowInterface, ctx *types.ExecutionContext) ([]*NetbridgeResult, error) {
	return e.netbridgeProcessor.GetResults(arrow, ctx)
}

// GetProcessTracker returns the process tracker for external access
func (e *Engine) GetProcessTracker() *process.ProcessTracker {
	return e.processTracker
}

// StopArrowProcesses stops all running processes for a given arrow
func (e *Engine) StopArrowProcesses(arrowName string, graceful bool, timeout time.Duration) error {
	return e.processTracker.StopProcesses(arrowName, graceful, timeout)
}

// GetArrowProcesses returns all running processes for a given arrow
func (e *Engine) GetArrowProcesses(arrowName string) []*process.ProcessInfo {
	return e.processTracker.GetProcesses(arrowName)
}

// HasRunningProcesses checks if an arrow has any running processes
func (e *Engine) HasRunningProcesses(arrowName string) bool {
	return e.processTracker.HasRunningProcesses(arrowName)
}

// Cleanup removes temporary files created during execution
func (e *Engine) Cleanup() error {
	if e.tempDir != "" {
		e.logger.Debug("Cleaning up temporary directory: %s", e.tempDir)
		return os.RemoveAll(e.tempDir)
	}
	return nil
} 