package strategies

import (
	"context"
	"net"
	"time"

	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"

	"github.com/rabbytesoftware/quiver.core/internal/domain/netbridge"
)

// natpmpClient defines the interface for NAT-PMP client operations.
type natpmpClient interface {
	GetExternalAddress() (*natpmp.GetExternalAddressResult, error)
	AddPortMapping(
		protocol string,
		internalPort int,
		externalPort int,
		lifetime int,
	) (*natpmp.AddPortMappingResult, error)
}

type natpmpStrategy struct {
	gateway   func() (net.IP, error)
	newClient func(net.IP) natpmpClient
}

// NewNATPMP returns a NAT-PMP port-forwarding strategy.
func NewNATPMP() Strategy {
	return &natpmpStrategy{
		gateway:   gateway.DiscoverGateway,
		newClient: natpmpNewClient,
	}
}

// natpmpTimeout bounds every exchange with the gateway. go-nat-pmp's default
// retries for about two minutes, which a gateway that does not speak NAT-PMP
// (most of them) turns into a two-minute stall of whoever asked.
const natpmpTimeout = 3 * time.Second

func natpmpNewClient(gw net.IP) natpmpClient {
	return natpmp.NewClientWithTimeout(gw, natpmpTimeout)
}

func (s *natpmpStrategy) Name() string {
	return "NAT-PMP"
}

func (s *natpmpStrategy) Available(
	ctx context.Context,
) bool {
	gw, err := s.gateway()
	if err != nil {
		return false
	}

	client := s.newClient(gw)
	_, err = client.GetExternalAddress()
	return err == nil
}

func (s *natpmpStrategy) Forward(
	ctx context.Context,
	port int,
	protocol netbridge.Protocol,
) error {
	gw, err := s.gateway()
	if err != nil {
		return err
	}

	client := s.newClient(gw)
	protos := natpmpProtocols(protocol)

	for _, proto := range protos {
		_, err := client.AddPortMapping(proto, port, port, 3600)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *natpmpStrategy) Reverse(
	ctx context.Context,
	port int,
	protocol netbridge.Protocol,
) error {
	gw, err := s.gateway()
	if err != nil {
		return err
	}

	client := s.newClient(gw)
	protos := natpmpProtocols(protocol)

	for _, proto := range protos {
		_, err := client.AddPortMapping(proto, port, port, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func natpmpProtocols(p netbridge.Protocol) []string {
	switch p {
	case netbridge.ProtocolTCP:
		return []string{"tcp"}
	case netbridge.ProtocolUDP:
		return []string{"udp"}
	case netbridge.ProtocolTCPUDP:
		return []string{"tcp", "udp"}
	default:
		return []string{}
	}
}
