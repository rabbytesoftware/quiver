package netbridge

import (
	"context"
	"testing"

	"github.com/rabbytesoftware/quiver/internal/models/port"
)

func TestNewNetbridge(t *testing.T) {
	nb := NewNetbridge()
	if nb == nil {
		t.Fatal("NewNetbridge() returned nil")
	}
}

func TestNetbridgeImpl_IsEnabled(t *testing.T) {
	nb := NewNetbridge()
	
	enabled := nb.IsEnabled()
	if !enabled {
		t.Error("IsEnabled() should return true")
	}
}

func TestNetbridgeImpl_IsAvailable(t *testing.T) {
	nb := NewNetbridge()
	
	available := nb.IsAvailable()
	if !available {
		t.Error("IsAvailable() should return true")
	}
}

func TestNetbridgeImpl_PublicIP(t *testing.T) {
	nb := NewNetbridge()
	ctx := context.Background()
	
	ip, err := nb.PublicIP(ctx)
	if err != nil {
		t.Errorf("PublicIP() returned error: %v", err)
	}
	if ip != "" {
		t.Error("PublicIP() should return empty string for unimplemented method")
	}
}

func TestNetbridgeImpl_LocalIP(t *testing.T) {
	nb := NewNetbridge()
	ctx := context.Background()
	
	ip, err := nb.LocalIP(ctx)
	if err != nil {
		t.Errorf("LocalIP() returned error: %v", err)
	}
	if ip != "" {
		t.Error("LocalIP() should return empty string for unimplemented method")
	}
}

func TestNetbridgeImpl_IsPortAvailable(t *testing.T) {
	nb := NewNetbridge()
	ctx := context.Background()
	
	available, err := nb.IsPortAvailable(ctx, 8080)
	if err != nil {
		t.Errorf("IsPortAvailable() returned error: %v", err)
	}
	if !available {
		t.Error("IsPortAvailable() should return true for unimplemented method")
	}
}

func TestNetbridgeImpl_ArePortsAvailable(t *testing.T) {
	nb := NewNetbridge()
	ctx := context.Background()
	
	ports := []int{8080, 8081, 8082}
	available, err := nb.ArePortsAvailable(ctx, ports)
	if err != nil {
		t.Errorf("ArePortsAvailable() returned error: %v", err)
	}
	if !available {
		t.Error("ArePortsAvailable() should return true for unimplemented method")
	}
}

func TestNetbridgeImpl_ForwardPort(t *testing.T) {
	nb := NewNetbridge()
	ctx := context.Background()
	
	rule, err := nb.ForwardPort(ctx, 8080)
	if err != nil {
		t.Errorf("ForwardPort() returned error: %v", err)
	}
	
	// Check that the returned rule has expected values
	if rule.StartPort != 8080 {
		t.Errorf("ForwardPort() returned wrong StartPort: got %d, want %d", rule.StartPort, 8080)
	}
	if rule.EndPort != 8080 {
		t.Errorf("ForwardPort() returned wrong EndPort: got %d, want %d", rule.EndPort, 8080)
	}
	if rule.Protocol != port.ProtocolTCP {
		t.Errorf("ForwardPort() returned wrong Protocol: got %v, want %v", rule.Protocol, port.ProtocolTCP)
	}
	if rule.ForwardingStatus != port.ForwardingStatusEnabled {
		t.Errorf("ForwardPort() returned wrong ForwardingStatus: got %v, want %v", rule.ForwardingStatus, port.ForwardingStatusEnabled)
	}
}

func TestNetbridgeImpl_ForwardPorts(t *testing.T) {
	nb := NewNetbridge()
	ctx := context.Background()
	
	ports := []int{8080, 8081, 8082}
	rules, err := nb.ForwardPorts(ctx, ports)
	if err != nil {
		t.Errorf("ForwardPorts() returned error: %v", err)
	}
	if rules == nil {
		t.Error("ForwardPorts() should return empty slice, not nil")
	}
	if len(rules) != 0 {
		t.Error("ForwardPorts() should return empty slice for unimplemented method")
	}
}

func TestNetbridgeImpl_ReversePort(t *testing.T) {
	nb := NewNetbridge()
	ctx := context.Background()
	
	rule, err := nb.ReversePort(ctx, 8080)
	if err != nil {
		t.Errorf("ReversePort() returned error: %v", err)
	}
	
	// Check that the returned rule has expected values
	if rule.StartPort != 8080 {
		t.Errorf("ReversePort() returned wrong StartPort: got %d, want %d", rule.StartPort, 8080)
	}
	if rule.EndPort != 8080 {
		t.Errorf("ReversePort() returned wrong EndPort: got %d, want %d", rule.EndPort, 8080)
	}
	if rule.Protocol != port.ProtocolTCP {
		t.Errorf("ReversePort() returned wrong Protocol: got %v, want %v", rule.Protocol, port.ProtocolTCP)
	}
	if rule.ForwardingStatus != port.ForwardingStatusEnabled {
		t.Errorf("ReversePort() returned wrong ForwardingStatus: got %v, want %v", rule.ForwardingStatus, port.ForwardingStatusEnabled)
	}
}

func TestNetbridgeImpl_ReversePorts(t *testing.T) {
	nb := NewNetbridge()
	ctx := context.Background()
	
	ports := []int{8080, 8081, 8082}
	rules, err := nb.ReversePorts(ctx, ports)
	if err != nil {
		t.Errorf("ReversePorts() returned error: %v", err)
	}
	if rules == nil {
		t.Error("ReversePorts() should return empty slice, not nil")
	}
	if len(rules) != 0 {
		t.Error("ReversePorts() should return empty slice for unimplemented method")
	}
}

func TestNetbridgeImpl_GetPortForwardingStatus(t *testing.T) {
	nb := NewNetbridge()
	ctx := context.Background()
	
	status, err := nb.GetPortForwardingStatus(ctx, 8080)
	if err != nil {
		t.Errorf("GetPortForwardingStatus() returned error: %v", err)
	}
	if status != port.ForwardingStatusEnabled {
		t.Errorf("GetPortForwardingStatus() returned wrong status: got %v, want %v", status, port.ForwardingStatusEnabled)
	}
}

func TestNetbridgeImpl_GetPortForwardingStatuses(t *testing.T) {
	nb := NewNetbridge()
	ctx := context.Background()
	
	ports := []int{8080, 8081, 8082}
	statuses, err := nb.GetPortForwardingStatuses(ctx, ports)
	if err != nil {
		t.Errorf("GetPortForwardingStatuses() returned error: %v", err)
	}
	if statuses == nil {
		t.Error("GetPortForwardingStatuses() should return empty slice, not nil")
	}
	if len(statuses) != 0 {
		t.Error("GetPortForwardingStatuses() should return empty slice for unimplemented method")
	}
}

func TestNetbridgeImpl_InterfaceCompliance(t *testing.T) {
	// Test that NetbridgeImpl implements NetbridgeInterface
	var _ NetbridgeInterface = &NetbridgeImpl{}
}

func TestNetbridgeImpl_MultipleInstances(t *testing.T) {
	nb1 := NewNetbridge()
	nb2 := NewNetbridge()
	
	// Both should be valid
	if nb1 == nil || nb2 == nil {
		t.Error("NewNetbridge() returned nil instance")
	}
	
	// Test that both instances work correctly
	if nb1.IsEnabled() != nb2.IsEnabled() {
		t.Error("Both instances should have same IsEnabled behavior")
	}
	
	if nb1.IsAvailable() != nb2.IsAvailable() {
		t.Error("Both instances should have same IsAvailable behavior")
	}
}
