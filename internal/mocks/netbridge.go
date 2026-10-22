package mocks

import (
	"context"

	"github.com/rabbytesoftware/quiver.core/internal/domain/netbridge"
)

type Netbridge struct {
	AllocatePort  int
	AllocateErr   error
	DeallocateErr error
	ShutdownErr   error

	// AllocateCalls counts Allocate calls; DeallocatedOwners records the owner
	// key of every DeallocateByOwner call.
	AllocateCalls     int
	DeallocatedOwners []string
}

func (m *Netbridge) Allocate(_ context.Context, _ string, _ netbridge.Protocol, _ int) (int, error) {
	m.AllocateCalls++
	return m.AllocatePort, m.AllocateErr
}

func (m *Netbridge) DeallocateByOwner(_ context.Context, ownerKey string) error {
	m.DeallocatedOwners = append(m.DeallocatedOwners, ownerKey)
	return m.DeallocateErr
}

func (m *Netbridge) Shutdown(_ context.Context) error {
	return m.ShutdownErr
}
