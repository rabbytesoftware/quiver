package netbridge

import (
	"fmt"
	"net"
	"strings"

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

	upnpManager   *upnp.UPnPManager
	natpmpClient  *natpmp.NATPMPClient
}

func NewNetbridge() (*Netbridge, error) {
	var netbridge *Netbridge = &Netbridge{
		PublicIP:     net.IPv4(0, 0, 0, 0),
		Ports:        make(map[int32]port.Port),
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
	// Start from a reasonable range for user applications
	startPort := uint16(8000)
	endPort := uint16(9000)

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
