package netbridge

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/rabbytesoftware/quiver/internal/config"
	natpmp "github.com/rabbytesoftware/quiver/internal/netbridge/natpmp"
	port "github.com/rabbytesoftware/quiver/internal/netbridge/port"
	upnp "github.com/rabbytesoftware/quiver/internal/netbridge/upnp"
)

// ForwardingMethod represents different port forwarding methods
type ForwardingMethod string

const (
	MethodUPnP   ForwardingMethod = "upnp"
	MethodManual ForwardingMethod = "manual"
	MethodNATPMP ForwardingMethod = "natpmp"
)

// PortForwardingResult represents the result of a port forwarding operation
type PortForwardingResult struct {
	Port     port.Port        `json:"port"`
	Method   ForwardingMethod `json:"method"`
	Success  bool             `json:"success"`
	Error    string           `json:"error,omitempty"`
}

type Netbridge struct {
	PublicIP   	net.IP
	Ports 		map[int32]port.Port
	config     *config.NetbridgeConfig

	upnpManager   *upnp.UPnPManager
	natpmpClient  *natpmp.NATPMPClient
}

func NewNetbridge(cfg *config.NetbridgeConfig) (*Netbridge, error) {
	var netbridge *Netbridge = &Netbridge{
		PublicIP:     net.IPv4(0, 0, 0, 0),
		Ports:        make(map[int32]port.Port),
		config:       cfg,
		upnpManager:  upnp.NewUPnPManager(),
		natpmpClient: nil, // Will be initialized if needed
	}

	// Try UPnP first
	publicIP, err := netbridge.upnpManager.GetPublicIP()
	if err != nil {
		// UPnP failed, try NAT-PMP
		natpmpClient, natErr := natpmp.NewNATPMPClient()
		if natErr == nil {
			netbridge.natpmpClient = natpmpClient
			publicIP, err = natpmpClient.GetExternalIP()
			if err != nil {
				natpmpClient.Close()
				return nil, fmt.Errorf("both UPnP and NAT-PMP failed: UPnP: %v, NAT-PMP: %v", err, natErr)
			}
		} else {
			return nil, fmt.Errorf("both UPnP and NAT-PMP failed: UPnP: %v, NAT-PMP: %v", err, natErr)
		}
	}

	netbridge.PublicIP = publicIP
	return netbridge, nil
}

// OpenPort opens a port using the best available method
func (n *Netbridge) OpenPort(portNum uint16, protocol string) (*PortForwardingResult, error) {
	// Normalize protocol
	protocol = strings.ToLower(protocol)
	
	// Handle combined protocols like "tcp/udp"
	protocols := []string{}
	if strings.Contains(protocol, "/") {
		protocols = strings.Split(protocol, "/")
	} else {
		protocols = []string{protocol}
	}

	results := []*PortForwardingResult{}
	
	for _, proto := range protocols {
		proto = strings.TrimSpace(proto)
		if proto != "tcp" && proto != "udp" {
			return nil, fmt.Errorf("unsupported protocol: %s (must be tcp, udp, or tcp/udp)", proto)
		}

		result := n.openSinglePort(portNum, proto)
		results = append(results, result)
		
		// Store successful port openings
		if result.Success {
			key := int32(portNum)
			if proto == "udp" {
				key = -int32(portNum) // Use negative for UDP to distinguish
			}
			n.Ports[key] = result.Port
		}
	}

	// Return the first result if single protocol, or combined result for multiple
	if len(results) == 1 {
		return results[0], nil
	}

	// For multiple protocols, return success if any succeeded
	combinedSuccess := false
	var combinedErrors []string
	for _, result := range results {
		if result.Success {
			combinedSuccess = true
		} else {
			combinedErrors = append(combinedErrors, fmt.Sprintf("%s: %s", result.Port.Protocol(), result.Error))
		}
	}

	combinedResult := &PortForwardingResult{
		Port:    port.NewPort(fmt.Sprintf("port-%d", portNum), portNum, "", protocol),
		Method:  MethodUPnP, // Primary method used
		Success: combinedSuccess,
	}

	if len(combinedErrors) > 0 {
		combinedResult.Error = strings.Join(combinedErrors, "; ")
	}

	return combinedResult, nil
}

// OpenPortAuto automatically finds and opens an available port
func (n *Netbridge) OpenPortAuto(protocol string) (*PortForwardingResult, error) {
	// Normalize protocol
	protocol = strings.ToLower(protocol)

	// Find an available port
	availablePort, err := n.findAvailablePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find available port: %v", err)
	}

	// Try to open the found port
	result, err := n.OpenPort(availablePort, protocol)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// findAvailablePort finds an available port on the local machine
func (n *Netbridge) findAvailablePort() (uint16, error) {
	// Use configured port range
	startPort := n.config.PortRangeStart
	endPort := n.config.PortRangeEnd

	for port := startPort; port <= endPort; port++ {
		// Skip ports that are already open in our netbridge
		key := int32(port)
		if _, exists := n.Ports[key]; exists {
			continue
		}
		if _, exists := n.Ports[-key]; exists { // Check UDP variant
			continue
		}

		// Test if port is available locally for TCP
		if n.isPortAvailable(port, "tcp") {
			return port, nil
		}
	}

	return 0, fmt.Errorf("no available ports found in range %d-%d", startPort, endPort)
}

// isPortAvailable checks if a port is available for binding
func (n *Netbridge) isPortAvailable(port uint16, protocol string) bool {
	// Try to bind to the port to see if it's available
	var network string
	if protocol == "tcp" {
		network = "tcp"
	} else {
		network = "udp"
	}

	addr := fmt.Sprintf(":%d", port)
	
	if network == "tcp" {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			return false
		}
		listener.Close()
		return true
	} else {
		conn, err := net.ListenPacket("udp", addr)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}
}

// openSinglePort opens a single port with a specific protocol
func (n *Netbridge) openSinglePort(portNum uint16, protocol string) *PortForwardingResult {
	portInstance := port.NewPortWithDescription(
		fmt.Sprintf("port-%d-%s", portNum, protocol),
		portNum,
		n.getLocalIP(),
		protocol,
		fmt.Sprintf("Port %d/%s opened via Quiver netbridge", portNum, protocol),
	)

	// Try UPnP first
	if err := n.upnpManager.ForwardPort(portInstance); err == nil {
		return &PortForwardingResult{
			Port:    portInstance,
			Method:  MethodUPnP,
			Success: true,
		}
	}

	// UPnP failed, try NAT-PMP if available
	if n.natpmpClient != nil {
		if err := n.natpmpClient.ForwardPort(portInstance); err == nil {
			return &PortForwardingResult{
				Port:    portInstance,
				Method:  MethodNATPMP,
				Success: true,
			}
		}
	}

	// Both methods failed
	return &PortForwardingResult{
		Port:    portInstance,
		Method:  MethodUPnP, // Primary method attempted
		Success: false,
		Error:   "Port forwarding failed: UPnP and NAT-PMP methods unsuccessful",
	}
}

// ClosePort closes a port forwarding
func (n *Netbridge) ClosePort(portNum uint16, protocol string) (*PortForwardingResult, error) {
	// Normalize protocol
	protocol = strings.ToLower(protocol)
	
	// Handle combined protocols like "tcp/udp"
	protocols := []string{}
	if strings.Contains(protocol, "/") {
		protocols = strings.Split(protocol, "/")
	} else {
		protocols = []string{protocol}
	}

	results := []*PortForwardingResult{}
	
	for _, proto := range protocols {
		proto = strings.TrimSpace(proto)
		if proto != "tcp" && proto != "udp" {
			return nil, fmt.Errorf("unsupported protocol: %s (must be tcp, udp, or tcp/udp)", proto)
		}

		result := n.closeSinglePort(portNum, proto)
		results = append(results, result)
		
		// Remove from stored ports
		if result.Success {
			key := int32(portNum)
			if proto == "udp" {
				key = -int32(portNum)
			}
			delete(n.Ports, key)
		}
	}

	// Return the first result if single protocol, or combined result for multiple
	if len(results) == 1 {
		return results[0], nil
	}

	// For multiple protocols, return success if any succeeded
	combinedSuccess := false
	var combinedErrors []string
	for _, result := range results {
		if result.Success {
			combinedSuccess = true
		} else {
			combinedErrors = append(combinedErrors, fmt.Sprintf("%s: %s", result.Port.Protocol(), result.Error))
		}
	}

	combinedResult := &PortForwardingResult{
		Port:    port.NewPort(fmt.Sprintf("port-%d", portNum), portNum, "", protocol),
		Method:  MethodUPnP,
		Success: combinedSuccess,
	}

	if len(combinedErrors) > 0 {
		combinedResult.Error = strings.Join(combinedErrors, "; ")
	}

	return combinedResult, nil
}

// closeSinglePort closes a single port with a specific protocol
func (n *Netbridge) closeSinglePort(portNum uint16, protocol string) *PortForwardingResult {
	portInstance := port.NewPortWithDescription(
		fmt.Sprintf("port-%d-%s", portNum, protocol),
		portNum,
		n.getLocalIP(),
		protocol,
		fmt.Sprintf("Port %d/%s closed via Quiver netbridge", portNum, protocol),
	)

	// Try UPnP first
	if err := n.upnpManager.ClosePort(portInstance); err == nil {
		return &PortForwardingResult{
			Port:    portInstance,
			Method:  MethodUPnP,
			Success: true,
		}
	}

	// UPnP failed, try NAT-PMP if available
	if n.natpmpClient != nil {
		if err := n.natpmpClient.ClosePort(portInstance); err == nil {
			return &PortForwardingResult{
				Port:    portInstance,
				Method:  MethodNATPMP,
				Success: true,
			}
		}
	}

	// Both methods failed
	return &PortForwardingResult{
		Port:    portInstance,
		Method:  MethodUPnP, // Primary method attempted
		Success: false,
		Error:   "Port closing failed: UPnP and NAT-PMP methods unsuccessful",
	}
}

// ListOpenPorts returns all currently open ports
func (n *Netbridge) ListOpenPorts() []port.Port {
	ports := make([]port.Port, 0, len(n.Ports))
	for _, p := range n.Ports {
		ports = append(ports, p)
	}
	return ports
}

// GetPublicIP returns the public IP address
func (n *Netbridge) GetPublicIP() net.IP {
	return n.PublicIP
}

// RefreshPublicIP refreshes the public IP address
func (n *Netbridge) RefreshPublicIP() error {
	publicIP, err := n.upnpManager.GetPublicIP()
	if err != nil {
		return err
	}
	n.PublicIP = publicIP
	return nil
}

// getLocalIP gets the local IP address (helper method)
func (n *Netbridge) getLocalIP() string {
	// This is a simple implementation - in a real scenario you might want to detect
	// the actual local IP address more accurately
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// AssignPortVariable handles port assignment for arrow variables
// It checks for user-specified ports, validates them, and falls back to auto-assignment
func (n *Netbridge) AssignPortVariable(userPortValue string, protocol string) (*PortForwardingResult, error) {
	// Normalize protocol
	protocol = strings.ToLower(protocol)
	
	// Check if user provided a specific port value
	if userPortValue != "" {
		if portNum, err := strconv.ParseUint(userPortValue, 10, 16); err == nil && portNum > 0 && portNum <= 65535 {
			// User specified a valid port, try to open it
			result, err := n.OpenPort(uint16(portNum), protocol)
			if err != nil {
				// Failed to open user-specified port, but still assign the port number
				return &PortForwardingResult{
					Port:    port.NewPort(fmt.Sprintf("port-%d", portNum), uint16(portNum), "", protocol),
					Method:  MethodManual,
					Success: false,
					Error:   fmt.Sprintf("Failed to open user-specified port: %v", err),
				}, nil
			}
			return result, nil
		}
		// Invalid port value provided, fall back to auto-assignment
	}
	
	// Auto-assign port
	result, err := n.OpenPortAuto(protocol)
	if err != nil {
		// Auto-assignment failed, find an available port without opening
		if availablePort, findErr := n.findAvailablePort(); findErr == nil {
			return &PortForwardingResult{
				Port:    port.NewPort(fmt.Sprintf("port-%d", availablePort), availablePort, "", protocol),
				Method:  MethodManual,
				Success: false,
				Error:   fmt.Sprintf("Auto-assignment failed: %v", err),
			}, nil
		}
		// Even finding an available port failed, use default fallback
		return &PortForwardingResult{
			Port:    port.NewPort("port-8080", 8080, "", protocol),
			Method:  MethodManual,
			Success: false,
			Error:   fmt.Sprintf("Auto-assignment failed and no available port found: %v", err),
		}, nil
	}
	
	return result, nil
}

// AssignPortVariableWithFallback handles port assignment with a fallback port when netbridge is unavailable
func (n *Netbridge) AssignPortVariableWithFallback(userPortValue string, protocol string) *PortForwardingResult {
	if n == nil {
		// Netbridge not available, find an available port without opening
		if userPortValue != "" {
			if portNum, err := strconv.ParseUint(userPortValue, 10, 16); err == nil && portNum > 0 && portNum <= 65535 {
				return &PortForwardingResult{
					Port:    port.NewPort(fmt.Sprintf("port-%d", portNum), uint16(portNum), "", protocol),
					Method:  MethodManual,
					Success: false,
					Error:   "Netbridge not available",
				}
			}
		}
		
		// Find available port or use default
		if availablePort := findAvailablePortStatic(nil); availablePort > 0 {
			return &PortForwardingResult{
				Port:    port.NewPort(fmt.Sprintf("port-%d", availablePort), availablePort, "", protocol),
				Method:  MethodManual,
				Success: false,
				Error:   "Netbridge not available",
			}
		}
		
		return &PortForwardingResult{
			Port:    port.NewPort("port-65534", 65534, "", protocol),
			Method:  MethodManual,
			Success: false,
			Error:   "Netbridge not available",
		}
	}
	
	result, _ := n.AssignPortVariable(userPortValue, protocol)
	return result
}

// findAvailablePortStatic finds an available port without requiring a Netbridge instance
func findAvailablePortStatic(cfg *config.NetbridgeConfig) uint16 {
	// Use default range if config is nil
	startPort := uint16(65000)
	endPort := uint16(65534)
	
	if cfg != nil {
		startPort = cfg.PortRangeStart
		endPort = cfg.PortRangeEnd
	}
	
	for port := startPort; port <= endPort; port++ {
		if isPortAvailableStatic(port) {
			return port
		}
	}
	
	return 0
}

// isPortAvailableStatic checks if a port is available without requiring a Netbridge instance
func isPortAvailableStatic(port uint16) bool {
	addr := fmt.Sprintf(":%d", port)
	
	// Test TCP
	if listener, err := net.Listen("tcp", addr); err == nil {
		listener.Close()
		return true
	}
	
	return false
}
