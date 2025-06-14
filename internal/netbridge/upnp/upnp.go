package port

import (
	"fmt"
	"net"
	"sync"

	"github.com/huin/goupnp/dcps/internetgateway1"
	"github.com/huin/goupnp/dcps/internetgateway2"
	"github.com/rabbytesoftware/quiver/netbridge/port"
)

// UPnPClient defines a common interface for all UPnP client types
type UPnPClient interface {
	AddPortMapping(remoteHost string, externalPort uint16, protocol string, 
		internalPort uint16, internalClient string, enabled bool, 
		description string, leaseDuration uint32) error
	DeletePortMapping(remoteHost string, externalPort uint16, protocol string) error
	GetExternalIPAddress() (string, error)
}

// UPnPManager handles UPnP operations including port forwarding and IP detection
type UPnPManager struct {
	mu               sync.Mutex
	clients          []UPnPClient
	clientsDiscovered bool
	discoveryErr     error
}

// NewUPnPManager creates a new UPnP manager instance
func NewUPnPManager() *UPnPManager {
	return &UPnPManager{
		clientsDiscovered: false,
		clients:          []UPnPClient{},
	}
}

// discoverServices finds UPnP services on the network
func (m *UPnPManager) discoverServices() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// If already discovered, return the previous result
	if m.clientsDiscovered {
		return m.discoveryErr
	}

	// Start with empty client list
	m.clients = []UPnPClient{}

	// Helper function to collect errors without stopping discovery
	var discoveryErrors []error
	addError := func(err error) {
		if err != nil {
			discoveryErrors = append(discoveryErrors, err)
		}
	}

	// Discover and add IGD v1 IP clients
	ipClients1, _, err := internetgateway1.NewWANIPConnection1Clients()
	addError(err)
	for i := range ipClients1 {
		m.clients = append(m.clients, ipClients1[i])
	}

	// Discover and add IGD v1 PPP clients
	pppClients1, _, err := internetgateway1.NewWANPPPConnection1Clients()
	addError(err)
	for i := range pppClients1 {
		m.clients = append(m.clients, pppClients1[i])
	}

	// Discover and add IGD v2 IP clients (v1)
	ipClients2v1, _, err := internetgateway2.NewWANIPConnection1Clients()
	addError(err)
	for i := range ipClients2v1 {
		m.clients = append(m.clients, ipClients2v1[i])
	}

	// Discover and add IGD v2 IP clients (v2)
	ipClients2v2, _, err := internetgateway2.NewWANIPConnection2Clients()
	addError(err)
	for i := range ipClients2v2 {
		m.clients = append(m.clients, ipClients2v2[i])
	}

	// Discover and add IGD v2 PPP clients
	pppClients2, _, err := internetgateway2.NewWANPPPConnection1Clients()
	addError(err)
	for i := range pppClients2 {
		m.clients = append(m.clients, pppClients2[i])
	}

	m.clientsDiscovered = true
	
	// Check if any services were discovered
	if len(m.clients) == 0 {
		if len(discoveryErrors) > 0 {
			m.discoveryErr = fmt.Errorf("no UPnP services discovered, errors: %v", discoveryErrors)
		} else {
			m.discoveryErr = fmt.Errorf("no UPnP services discovered")
		}
		return m.discoveryErr
	}

	return nil
}

// ForwardPort forwards the specified port through UPnP
func (m *UPnPManager) ForwardPort(port port.Port) error {
	if err := m.discoverServices(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Port mapping parameters
	externalPort := port.Port()
	internalPort := port.Port()
	description := fmt.Sprintf("Watcher|%s:%d (%s)", port.Host, port.Port, port.Protocol)
	duration := uint32(86400) // 24 hours in seconds

	// Try all clients with a single loop
	for _, client := range m.clients {
		err := client.AddPortMapping(
			"", externalPort, port.Protocol(), internalPort,
			port.Host(), true, description, duration,
		)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("failed to add port mapping on all available UPnP services")
}

// ClosePort removes the port forwarding for the specified port
func (m *UPnPManager) ClosePort(port port.Port) error {
	if err := m.discoverServices(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	externalPort := uint16(port.Port())
	
	// Try all clients with a single loop
	for _, client := range m.clients {
		err := client.DeletePortMapping("", externalPort, port.Protocol())
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("failed to remove port mapping on all available UPnP services")
}

// GetPublicIP retrieves the router's public IP address
func (m *UPnPManager) GetPublicIP() (net.IP, error) {
	if err := m.discoverServices(); err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Try all clients with a single loop
	for _, client := range m.clients {
		ip, err := client.GetExternalIPAddress()
		if err == nil && ip != "" {
			return net.ParseIP(ip), nil
		}
	}

	return nil, fmt.Errorf("failed to get public IP address from any UPnP service")
}

// RefreshServices forces a rediscovery of UPnP services
func (m *UPnPManager) RefreshServices() error {
	m.mu.Lock()
	m.clientsDiscovered = false
	m.mu.Unlock()
	return m.discoverServices()
}