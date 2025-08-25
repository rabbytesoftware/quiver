package natpmp

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/rabbytesoftware/quiver/internal/netbridge/port"
)

const (
	// NAT-PMP protocol constants
	natPMPPort        = 5351
	natPMPVersion     = 0
	opCodeExternalIP  = 0
	opCodeMapTCP      = 1
	opCodeMapUDP      = 2
	
	// Response codes
	resultSuccess             = 0
	resultUnsupportedVersion  = 1
	resultNotAuthorized       = 2
	resultNetworkFailure      = 3
	resultOutOfResources      = 4
	resultUnsupportedOpcode   = 5
)

// NATPMPClient represents a NAT-PMP client
type NATPMPClient struct {
	gateway net.IP
	conn    *net.UDPConn
}

// NewNATPMPClient creates a new NAT-PMP client
func NewNATPMPClient() (*NATPMPClient, error) {
	gateway, err := findDefaultGateway()
	if err != nil {
		return nil, fmt.Errorf("failed to find default gateway: %v", err)
	}

	// Create UDP connection
	laddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve local UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP connection: %v", err)
	}

	client := &NATPMPClient{
		gateway: gateway,
		conn:    conn,
	}

	// Test if NAT-PMP is available by getting external IP
	_, err = client.GetExternalIP()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("NAT-PMP not available: %v", err)
	}

	return client, nil
}

// Close closes the NAT-PMP client connection
func (c *NATPMPClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetExternalIP gets the external IP address via NAT-PMP
func (c *NATPMPClient) GetExternalIP() (net.IP, error) {
	// Build request packet
	request := make([]byte, 2)
	request[0] = natPMPVersion
	request[1] = opCodeExternalIP

	// Send request
	gatewayAddr := &net.UDPAddr{IP: c.gateway, Port: natPMPPort}
	_, err := c.conn.WriteToUDP(request, gatewayAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to send NAT-PMP request: %v", err)
	}

	// Read response with timeout
	c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	response := make([]byte, 12)
	n, _, err := c.conn.ReadFromUDP(response)
	if err != nil {
		return nil, fmt.Errorf("failed to read NAT-PMP response: %v", err)
	}

	if n < 12 {
		return nil, fmt.Errorf("invalid NAT-PMP response length: %d", n)
	}

	// Parse response
	version := response[0]
	opcode := response[1]
	resultCode := binary.BigEndian.Uint16(response[2:4])

	if version != natPMPVersion {
		return nil, fmt.Errorf("invalid NAT-PMP version: %d", version)
	}

	if opcode != opCodeExternalIP+128 { // Response opcode is request + 128
		return nil, fmt.Errorf("invalid NAT-PMP opcode: %d", opcode)
	}

	if resultCode != resultSuccess {
		return nil, fmt.Errorf("NAT-PMP error: %d", resultCode)
	}

	// Extract external IP (bytes 8-11)
	externalIP := net.IP(response[8:12])
	return externalIP, nil
}

// AddPortMapping adds a port mapping via NAT-PMP
func (c *NATPMPClient) AddPortMapping(protocol string, internalPort, externalPort uint16, lifetime uint32) error {
	var opcode byte
	if protocol == "tcp" {
		opcode = opCodeMapTCP
	} else if protocol == "udp" {
		opcode = opCodeMapUDP
	} else {
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}

	// Build request packet (12 bytes)
	request := make([]byte, 12)
	request[0] = natPMPVersion
	request[1] = opcode
	// bytes 2-3 are reserved (0)
	binary.BigEndian.PutUint16(request[4:6], internalPort)
	binary.BigEndian.PutUint16(request[6:8], externalPort)
	binary.BigEndian.PutUint32(request[8:12], lifetime)

	// Send request
	gatewayAddr := &net.UDPAddr{IP: c.gateway, Port: natPMPPort}
	_, err := c.conn.WriteToUDP(request, gatewayAddr)
	if err != nil {
		return fmt.Errorf("failed to send NAT-PMP port mapping request: %v", err)
	}

	// Read response with timeout
	c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	response := make([]byte, 16)
	n, _, err := c.conn.ReadFromUDP(response)
	if err != nil {
		return fmt.Errorf("failed to read NAT-PMP port mapping response: %v", err)
	}

	if n < 16 {
		return fmt.Errorf("invalid NAT-PMP port mapping response length: %d", n)
	}

	// Parse response
	version := response[0]
	responseOpcode := response[1]
	resultCode := binary.BigEndian.Uint16(response[2:4])

	if version != natPMPVersion {
		return fmt.Errorf("invalid NAT-PMP version: %d", version)
	}

	if responseOpcode != opcode+128 { // Response opcode is request + 128
		return fmt.Errorf("invalid NAT-PMP response opcode: %d", responseOpcode)
	}

	if resultCode != resultSuccess {
		return fmt.Errorf("NAT-PMP port mapping error: %d", resultCode)
	}

	return nil
}

// RemovePortMapping removes a port mapping via NAT-PMP
func (c *NATPMPClient) RemovePortMapping(protocol string, externalPort uint16) error {
	// To remove a mapping, set lifetime to 0
	return c.AddPortMapping(protocol, 0, externalPort, 0)
}

// ForwardPort forwards a port using NAT-PMP
func (c *NATPMPClient) ForwardPort(p port.Port) error {
	// Use 24-hour lifetime (86400 seconds)
	return c.AddPortMapping(p.Protocol(), p.Port(), p.Port(), 86400)
}

// ClosePort removes port forwarding using NAT-PMP
func (c *NATPMPClient) ClosePort(p port.Port) error {
	return c.RemovePortMapping(p.Protocol(), p.Port())
}

// findDefaultGateway finds the default gateway IP address
func findDefaultGateway() (net.IP, error) {
	// Try to find gateway by creating a connection to a remote address
	// This will route through the default gateway
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, fmt.Errorf("failed to create connection to find gateway: %v", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	localIP := localAddr.IP

	// Common gateway patterns
	if localIP.To4() != nil {
		ip := localIP.To4()
		// Try .1 as gateway (most common)
		gateway := net.IPv4(ip[0], ip[1], ip[2], 1)
		
		// Test if this gateway responds to NAT-PMP
		if testNATPMPGateway(gateway) {
			return gateway, nil
		}

		// Try .254 as gateway (some routers use this)
		gateway = net.IPv4(ip[0], ip[1], ip[2], 254)
		if testNATPMPGateway(gateway) {
			return gateway, nil
		}
	}

	return nil, fmt.Errorf("no responsive NAT-PMP gateway found")
}

// testNATPMPGateway tests if a given IP responds to NAT-PMP requests
func testNATPMPGateway(gateway net.IP) bool {
	// Create a temporary UDP connection
	laddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return false
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return false
	}
	defer conn.Close()

	// Build external IP request
	request := make([]byte, 2)
	request[0] = natPMPVersion
	request[1] = opCodeExternalIP

	// Send request
	gatewayAddr := &net.UDPAddr{IP: gateway, Port: natPMPPort}
	_, err = conn.WriteToUDP(request, gatewayAddr)
	if err != nil {
		return false
	}

	// Try to read response with short timeout
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	response := make([]byte, 12)
	n, _, err := conn.ReadFromUDP(response)
	if err != nil {
		return false
	}

	// Basic validation of response
	if n >= 4 && response[0] == natPMPVersion && response[1] == (opCodeExternalIP+128) {
		resultCode := binary.BigEndian.Uint16(response[2:4])
		return resultCode == resultSuccess
	}

	return false
} 