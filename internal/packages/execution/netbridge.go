package execution

import (
	"strconv"

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

// NetbridgeProcessor handles netbridge-related functionality for arrow execution
type NetbridgeProcessor struct {
	logger    *logger.Logger
	netbridge *netbridge.Netbridge
}

// NewNetbridgeProcessor creates a new netbridge processor
func NewNetbridgeProcessor(logger *logger.Logger) *NetbridgeProcessor {
	netbridgeInstance, err := netbridge.NewNetbridge()
	if err != nil {
		logger.Warn("Netbridge initialization failed: %v (port forwarding disabled)", err)
		netbridgeInstance = nil
	}

	return &NetbridgeProcessor{
		logger:    logger.WithService("netbridge-processor"),
		netbridge: netbridgeInstance,
	}
}

// ProcessVariables processes netbridge port assignments and adds them to execution context
func (np *NetbridgeProcessor) ProcessVariables(arrow manifest.ArrowInterface, ctx *types.ExecutionContext) ([]*NetbridgeResult, error) {
	specs := arrow.GetNetbridge()
	if len(specs) == 0 {
		np.logger.Debug("No netbridge variables for arrow %s", arrow.Name())
		return nil, nil
	}

	np.logger.Info("Processing %d netbridge variables for arrow %s", len(specs), arrow.Name())
	
	if ctx.Variables == nil {
		ctx.Variables = make(map[string]string)
	}

	results := make([]*NetbridgeResult, 0, len(specs))
	for _, spec := range specs {
		result := np.processVariable(spec, ctx)
		results = append(results, result)
		ctx.Variables[result.VariableName] = strconv.Itoa(int(result.Port))
		np.logger.Debug("Variable %s = port %d (success: %t)", result.VariableName, result.Port, result.Success)
	}

	return results, nil
}

// processVariable processes a single netbridge variable
func (np *NetbridgeProcessor) processVariable(spec manifest.Netbridge, ctx *types.ExecutionContext) *NetbridgeResult {
	varName := spec.GetName()
	protocol := spec.GetProtocol()
	userPortValue := ctx.Variables[varName]

	np.logger.Info("Processing netbridge variable %s (protocol: %s)", varName, protocol)

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

// GetResults returns netbridge processing results for API reporting
func (np *NetbridgeProcessor) GetResults(arrow manifest.ArrowInterface, ctx *types.ExecutionContext) ([]*NetbridgeResult, error) {
	return np.ProcessVariables(arrow, ctx)
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