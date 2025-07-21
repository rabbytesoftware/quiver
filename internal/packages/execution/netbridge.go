package execution

import (
	"strconv"

	"github.com/rabbytesoftware/quiver/internal/config"
	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/netbridge"
	"github.com/rabbytesoftware/quiver/internal/packages/manifest"
	"github.com/rabbytesoftware/quiver/internal/packages/types"
)

// NetbridgeResult represents the result of netbridge variable processing
type NetbridgeResult struct {
	VariableName string                          `json:"variable_name"`
	Port         uint16                          `json:"port"`
	Protocol     string                          `json:"protocol"`
	Success      bool                            `json:"success"`
	Error        string                          `json:"error,omitempty"`
	Result       *netbridge.PortForwardingResult `json:"result,omitempty"`
}

// NetbridgeIdentification represents identification of netbridge variables without port assignment
type NetbridgeIdentification struct {
	VariableName string `json:"variable_name"`
	Protocol     string `json:"protocol"`
	UserSpecified bool  `json:"user_specified"`
	UserValue    string `json:"user_value,omitempty"`
}

// NetbridgeProcessor handles netbridge-related functionality for arrow execution
type NetbridgeProcessor struct {
	logger    *logger.Logger
	netbridge *netbridge.Netbridge
}

// NewNetbridgeProcessor creates a new netbridge processor
func NewNetbridgeProcessor(logger *logger.Logger) *NetbridgeProcessor {
	return NewNetbridgeProcessorWithConfig(logger, nil)
}

// NewNetbridgeProcessorWithConfig creates a new netbridge processor with config
func NewNetbridgeProcessorWithConfig(logger *logger.Logger, cfg *config.NetbridgeConfig) *NetbridgeProcessor {
	// Use default config if none provided
	if cfg == nil {
		cfg = &config.NetbridgeConfig{
			PortRangeStart: 65000,
			PortRangeEnd:   65534,
		}
	}

	netbridgeInstance, err := netbridge.NewNetbridge(cfg)
	if err != nil {
		logger.Warn("Netbridge initialization failed: %v (port forwarding disabled)", err)
		netbridgeInstance = nil
	}

	return &NetbridgeProcessor{
		logger:    logger.WithService("netbridge-processor"),
		netbridge: netbridgeInstance,
	}
}

// IdentifyVariables identifies netbridge variables for API reporting without assigning ports
// This is used during initialization to show what netbridge variables exist
func (np *NetbridgeProcessor) IdentifyVariables(arrow manifest.ArrowInterface, ctx *types.ExecutionContext) ([]*NetbridgeIdentification, error) {
	specs := arrow.GetNetbridge()
	if len(specs) == 0 {
		np.logger.Debug("No netbridge variables for arrow %s", arrow.Name())
		return nil, nil
	}

	np.logger.Info("Identifying %d netbridge variables for arrow %s (no ports assigned)", len(specs), arrow.Name())
	
	results := make([]*NetbridgeIdentification, 0, len(specs))
	for _, spec := range specs {
		varName := spec.GetName()
		protocol := spec.GetProtocol()
		userValue := ""
		userSpecified := false
		
		if ctx.Variables != nil {
			if value, exists := ctx.Variables[varName]; exists && value != "" {
				userValue = value
				userSpecified = true
			}
		}
		
		identification := &NetbridgeIdentification{
			VariableName:  varName,
			Protocol:      protocol,
			UserSpecified: userSpecified,
			UserValue:     userValue,
		}
		
		results = append(results, identification)
		np.logger.Debug("Identified netbridge variable %s (protocol: %s, user_specified: %t)", varName, protocol, userSpecified)
	}

	return results, nil
}

// ProcessVariablesRuntime processes netbridge port assignments at runtime and adds them to execution context
// This should ONLY be called during actual execution of the execute method
func (np *NetbridgeProcessor) ProcessVariablesRuntime(arrow manifest.ArrowInterface, ctx *types.ExecutionContext) ([]*NetbridgeResult, error) {
	specs := arrow.GetNetbridge()
	if len(specs) == 0 {
		np.logger.Debug("No netbridge variables for runtime processing in arrow %s", arrow.Name())
		return nil, nil
	}

	np.logger.Info("Runtime processing of %d netbridge variables for arrow %s", len(specs), arrow.Name())
	
	if ctx.Variables == nil {
		ctx.Variables = make(map[string]string)
	}

	results := make([]*NetbridgeResult, 0, len(specs))
	for _, spec := range specs {
		result := np.processVariableRuntime(spec, ctx)
		results = append(results, result)
		ctx.Variables[result.VariableName] = strconv.Itoa(int(result.Port))
		np.logger.Debug("Runtime assignment: Variable %s = port %d (success: %t)", result.VariableName, result.Port, result.Success)
	}

	return results, nil
}

// processVariableRuntime processes a single netbridge variable at runtime
func (np *NetbridgeProcessor) processVariableRuntime(spec manifest.Netbridge, ctx *types.ExecutionContext) *NetbridgeResult {
	varName := spec.GetName()
	protocol := spec.GetProtocol()
	userPortValue := ctx.Variables[varName]

	np.logger.Info("Runtime processing netbridge variable %s (protocol: %s)", varName, protocol)

	// Use netbridge package to handle all port logic
	portResult := np.netbridge.AssignPortVariableWithFallback(userPortValue, protocol)

	return &NetbridgeResult{
		VariableName: varName,
		Port:         portResult.Port.Port(),
		Protocol:     protocol,
		Success:      portResult.Success,
		Error:        portResult.Error,
		Result:       portResult,
	}
}

// ProcessVariables is the legacy method that should only be used for backwards compatibility
// DEPRECATED: Use IdentifyVariables for initialization and ProcessVariablesRuntime for execution
func (np *NetbridgeProcessor) ProcessVariables(arrow manifest.ArrowInterface, ctx *types.ExecutionContext) ([]*NetbridgeResult, error) {
	np.logger.Warn("DEPRECATED: ProcessVariables called - this should be replaced with IdentifyVariables or ProcessVariablesRuntime")
	return np.ProcessVariablesRuntime(arrow, ctx)
}

// GetResults returns netbridge processing results for API reporting
// DEPRECATED: Use IdentifyVariables for initialization reporting
func (np *NetbridgeProcessor) GetResults(arrow manifest.ArrowInterface, ctx *types.ExecutionContext) ([]*NetbridgeResult, error) {
	np.logger.Warn("DEPRECATED: GetResults called - this should be replaced with IdentifyVariables for initialization")
	return np.ProcessVariablesRuntime(arrow, ctx)
}

// IsNetbridgeAvailable returns whether netbridge functionality is available
func (np *NetbridgeProcessor) IsNetbridgeAvailable() bool {
	return np.netbridge != nil
}

// LogProcessingResults logs the results of netbridge processing
func (np *NetbridgeProcessor) LogProcessingResults(results []*NetbridgeResult) {
	if len(results) == 0 {
		return
	}

	np.logger.Info("Netbridge processing completed for %d variables", len(results))
	for _, result := range results {
		if result.Success {
			np.logger.Info("✓ %s: port %d/%s opened successfully", result.VariableName, result.Port, result.Protocol)
		} else {
			np.logger.Warn("⚠ %s: port %d/%s failed - %s", result.VariableName, result.Port, result.Protocol, result.Error)
			np.logger.Info("Port %d assigned to %s for manual configuration", result.Port, result.VariableName)
		}
	}
} 