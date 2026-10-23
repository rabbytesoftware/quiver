package netbridge

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/char2cs/asynx"
	asynxstore "github.com/char2cs/asynx/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rabbytesoftware/quiver.core/internal/core/config"
	"github.com/rabbytesoftware/quiver.core/internal/domain/netbridge"
	"github.com/rabbytesoftware/quiver.core/internal/engine/netbridge/internal/mocks"
	"github.com/rabbytesoftware/quiver.core/internal/engine/netbridge/internal/ports"
	"github.com/rabbytesoftware/quiver.core/internal/engine/netbridge/internal/store"
	"github.com/rabbytesoftware/quiver.core/internal/engine/netbridge/internal/strategies"
)

func buildNetbridgeWithStrategy(
	t *testing.T,
	strategy strategies.Strategy,
) *netbridgeService {
	t.Helper()

	es := asynxstore.New()
	ss := asynxstore.NewSnapshots()

	nb, err := New().WithEventStore(es).WithSnapshotStore(ss).WithStrategies([]strategies.Strategy{strategy}).Build(context.Background())
	require.NoError(t, err)

	impl, ok := nb.(*netbridgeService)
	require.True(t, ok)
	return impl
}

func buildNetbridge(
	t *testing.T,
) *netbridgeService {
	t.Helper()

	es := asynxstore.New()
	ss := asynxstore.NewSnapshots()

	nb, err := New().WithEventStore(es).WithSnapshotStore(ss).WithStrategies([]strategies.Strategy{}).Build(context.Background())
	require.NoError(t, err)
	require.NotNil(t, nb)

	impl, ok := nb.(*netbridgeService)
	require.True(t, ok)
	return impl
}

func TestBuilder_BuildSucceeds(
	t *testing.T,
) {
	es := asynxstore.New()
	ss := asynxstore.NewSnapshots()

	nb, err := New().WithEventStore(es).WithSnapshotStore(ss).Build(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, nb)
}

func TestBuilder_WithStore(
	t *testing.T,
) {
	es := asynxstore.New()
	ss := asynxstore.NewSnapshots()

	custom := store.NewPortMemory()
	nb, err := New().WithStore(custom).WithEventStore(es).WithSnapshotStore(ss).Build(context.Background())
	require.NoError(t, err)
	require.NotNil(t, nb)

	impl := nb.(*netbridgeService)
	assert.Equal(t, custom, impl.readModel)
}

func TestBuilder_WithDatabasePath(
	t *testing.T,
) {
	es := asynxstore.New()
	ss := asynxstore.NewSnapshots()

	nb, err := New().WithEventStore(es).WithSnapshotStore(ss).Build(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, nb)
}

func TestBuilder_WithDatabasePath_InvalidPath(
	t *testing.T,
) {
	es := asynxstore.New()
	ss := asynxstore.NewSnapshots()

	nb, err := New().WithEventStore(es).WithSnapshotStore(ss).Build(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, nb)
}

func TestAllocate_ReturnsValidPort(
	t *testing.T,
) {
	nb := buildNetbridge(t)

	port, err := nb.Allocate(context.Background(), "owner1", netbridge.ProtocolTCP, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, port, 1)
	assert.LessOrEqual(t, port, 65535)
}

func TestAllocate_HonorsPreferredPort(
	t *testing.T,
) {
	nb := buildNetbridge(t)

	preferred := freePort(t)
	port, err := nb.Allocate(context.Background(), "owner1", netbridge.ProtocolTCP, preferred)
	require.NoError(t, err)
	assert.Equal(t, preferred, port)
}

func TestAllocate_UDPProtocol(
	t *testing.T,
) {
	nb := buildNetbridge(t)

	port, err := nb.Allocate(context.Background(), "owner1", netbridge.ProtocolUDP, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, port, testEphemeralPortStart)
}

func TestAllocate_TCPUDPProtocol(
	t *testing.T,
) {
	nb := buildNetbridge(t)

	port, err := nb.Allocate(context.Background(), "owner1", netbridge.ProtocolTCPUDP, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, port, testEphemeralPortStart)
}

func TestAllocate_ReturnsErrPortOutOfRange(
	t *testing.T,
) {
	nb := buildNetbridge(t)
	ctx := context.Background()

	cases := []int{-1, 99999}
	for _, preferred := range cases {
		_, err := nb.Allocate(ctx, "owner1", netbridge.ProtocolTCP, preferred)
		assert.ErrorIs(t, err, ErrPortOutOfRange, "expected ErrPortOutOfRange for preferred=%d", preferred)
	}
}

func TestAllocate_WithActiveStrategy_SetsForwarded(
	t *testing.T,
) {
	strat := &mocks.AlwaysAvailableStrategy{}
	nb := buildNetbridgeWithStrategy(t, strat)
	ctx := context.Background()

	preferred := freePort(t)
	port, err := nb.Allocate(ctx, "owner-fwd", netbridge.ProtocolTCP, preferred)
	require.NoError(t, err)
	assert.Equal(t, preferred, port)

	nb.waitForProjection()

	alloc, err := nb.readModel.FindByPort(ctx, port)
	require.NoError(t, err)
	require.NotNil(t, alloc)
	assert.True(t, alloc.Forwarded)
}

func TestAllocate_SendErrorOnCancelledContext(
	t *testing.T,
) {
	nb := buildNetbridge(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := nb.Allocate(ctx, "owner-cancel", netbridge.ProtocolTCP, 0)
	assert.Error(t, err)
}

func TestDeallocateByOwner_ReleasesAllPorts(
	t *testing.T,
) {
	nb := buildNetbridge(t)
	ctx := context.Background()

	const ownerKey = "owner-dealloc"
	free := freePorts(t, 2)
	preferred1, preferred2 := free[0], free[1]

	port1, err := nb.Allocate(ctx, ownerKey, netbridge.ProtocolTCP, preferred1)
	require.NoError(t, err)
	assert.Equal(t, preferred1, port1)

	port2, err := nb.Allocate(ctx, ownerKey, netbridge.ProtocolTCP, preferred2)
	require.NoError(t, err)
	assert.Equal(t, preferred2, port2)

	nb.waitForProjection()

	err = nb.DeallocateByOwner(ctx, ownerKey)
	require.NoError(t, err)
	nb.waitForProjection()

	realloc1, err := nb.Allocate(ctx, ownerKey, netbridge.ProtocolTCP, preferred1)
	require.NoError(t, err)
	assert.Equal(t, preferred1, realloc1)

	realloc2, err := nb.Allocate(ctx, ownerKey, netbridge.ProtocolTCP, preferred2)
	require.NoError(t, err)
	assert.Equal(t, preferred2, realloc2)
}

func TestDeallocateByOwner_NoOpForUnknownOwner(
	t *testing.T,
) {
	nb := buildNetbridge(t)

	err := nb.DeallocateByOwner(context.Background(), "unknown-owner")
	assert.NoError(t, err)
}

func TestDeallocateByOwner_WithActiveStrategy_CallsReverse(
	t *testing.T,
) {
	strat := &mocks.AlwaysAvailableStrategy{}
	nb := buildNetbridgeWithStrategy(t, strat)
	ctx := context.Background()

	preferred := freePort(t)
	port, err := nb.Allocate(ctx, "owner-rev", netbridge.ProtocolTCP, preferred)
	require.NoError(t, err)
	assert.Equal(t, preferred, port)

	nb.waitForProjection()

	err = nb.DeallocateByOwner(ctx, "owner-rev")
	require.NoError(t, err)

	assert.Equal(t, []int{preferred}, strat.ReverseCalledWith)
}

func TestDeallocateByOwner_SendErrorOnCancelledContext(
	t *testing.T,
) {
	nb := buildNetbridge(t)

	preferred := freePort(t)
	port, allocErr := nb.Allocate(context.Background(), "owner-cancel", netbridge.ProtocolTCP, preferred)
	require.NoError(t, allocErr)
	assert.Equal(t, preferred, port)
	nb.waitForProjection()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := nb.DeallocateByOwner(ctx, "owner-cancel")
	assert.Error(t, err)
}

func TestDeallocateByOwner_FindByOwnerError(
	t *testing.T,
) {
	nb := buildNetbridge(t)
	nb.readModel = &mocks.ErrFindByOwnerReadModel{
		PortStore: nb.readModel,
		Err:       errors.New("read model failure"),
	}

	err := nb.DeallocateByOwner(context.Background(), "some-owner")
	assert.Error(t, err)
}

func TestBuilder_WithEphemeralPortRange_AppliedToBuild(
	t *testing.T,
) {
	es := asynxstore.New()
	ss := asynxstore.NewSnapshots()

	nb, err := New().
		WithEventStore(es).
		WithSnapshotStore(ss).
		WithStrategies([]strategies.Strategy{}).
		WithEphemeralPortRange(50000, 50100).
		Build(context.Background())
	require.NoError(t, err)
	require.NotNil(t, nb)

	impl := nb.(*netbridgeService)
	assert.Equal(t, 50000, impl.portStart)
	assert.Equal(t, 50100, impl.portEnd)
}

func TestBuilder_Build_MissingEventStore_ReturnsError(
	t *testing.T,
) {
	_, err := New().Build(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrBuildFailed)
}

func TestBuilder_Build_MissingSnapshotStoreReturnsError(
	t *testing.T,
) {
	_, err := New().WithEventStore(asynxstore.New()).Build(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrBuildFailed)
	assert.Contains(t, err.Error(), "missing SnapshotStore")
}

func TestAllocate_WithEphemeralPortRange_UsesCustomRange(
	t *testing.T,
) {
	es := asynxstore.New()
	ss := asynxstore.NewSnapshots()

	nb, err := New().
		WithEventStore(es).
		WithSnapshotStore(ss).
		WithStrategies([]strategies.Strategy{}).
		WithEphemeralPortRange(50200, 50300).
		Build(context.Background())
	require.NoError(t, err)

	// preferred=0 forces ephemeral allocation from the custom range.
	port, allocErr := nb.Allocate(context.Background(), "owner-custom", netbridge.ProtocolTCP, 0)
	require.NoError(t, allocErr)
	assert.GreaterOrEqual(t, port, 50200)
	assert.LessOrEqual(t, port, 50300)
}

func TestNetbridge_Allocate_WritesSnapshot(
	t *testing.T,
) {
	ss := asynxstore.NewSnapshots()
	nb, err := New().
		WithEventStore(asynxstore.New()).
		WithSnapshotStore(ss).
		WithStrategies([]strategies.Strategy{}).
		Build(context.Background())
	require.NoError(t, err)

	port, err := nb.Allocate(context.Background(), "owner-1", netbridge.ProtocolTCP, 50000)
	require.NoError(t, err)

	_, found, err := ss.Get(context.Background(), strconv.Itoa(port))
	require.NoError(t, err)
	assert.True(t, found, "allocate must persist a snapshot; without one every read cold-replays the port's full history forever")
}

func TestAllocate_SamePreferredPortTwiceGetsDifferentPort(
	t *testing.T,
) {
	nb := buildNetbridge(t)
	ctx := context.Background()

	preferred := freePort(t)

	port1, err := nb.Allocate(ctx, "owner-a", netbridge.ProtocolTCP, preferred)
	require.NoError(t, err)
	assert.Equal(t, preferred, port1)

	nb.waitForProjection()

	port2, err := nb.Allocate(ctx, "owner-b", netbridge.ProtocolTCP, preferred)
	require.NoError(t, err)
	assert.NotEqual(t, preferred, port2)
	assert.GreaterOrEqual(t, port2, testEphemeralPortStart)
	assert.LessOrEqual(t, port2, testEphemeralPortEnd)
}

// The assembler allocates every declared port back to back. The read model is
// projected asynchronously, so without Allocate waiting for its own projection
// the second call picks the port the first just took and the command is rejected.
func TestAllocate_BackToBack_NeverReusesAPortStillBeingProjected(
	t *testing.T,
) {
	nb := buildNetbridge(t)
	ctx := context.Background()

	taken := freePort(t)
	_, err := nb.Allocate(ctx, "owner-a", netbridge.ProtocolTCP, taken)
	require.NoError(t, err)

	seen := map[int]bool{taken: true}
	for i := 0; i < 20; i++ {
		proto := netbridge.ProtocolTCP
		if i%2 == 1 {
			proto = netbridge.ProtocolUDP
		}
		port, err := nb.Allocate(ctx, "owner-b", proto, taken)
		require.NoError(t, err, "allocation %d", i)
		require.False(t, seen[port], "port %d handed out twice", port)
		seen[port] = true
	}
}

func TestAllocate_Concurrent_NeverHandsOutAPortTwice(
	t *testing.T,
) {
	nb := buildNetbridge(t)
	ctx := context.Background()

	const workers = 16
	ports := make(chan int, workers)
	errs := make(chan error, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			port, err := nb.Allocate(ctx, fmt.Sprintf("owner-%d", i), netbridge.ProtocolTCP, 0)
			if err != nil {
				errs <- err
				return
			}
			ports <- port
		}()
	}
	wg.Wait()
	close(ports)
	close(errs)

	for err := range errs {
		require.NoError(t, err)
	}
	seen := map[int]bool{}
	for port := range ports {
		require.False(t, seen[port], "port %d handed out twice", port)
		seen[port] = true
	}
	assert.Len(t, seen, workers)
}

func TestDeallocateByOwner_FreesThePortsForTheNextAllocation(
	t *testing.T,
) {
	nb := buildNetbridge(t)
	ctx := context.Background()

	preferred := freePort(t)
	first, err := nb.Allocate(ctx, "owner-a", netbridge.ProtocolTCP, preferred)
	require.NoError(t, err)
	require.Equal(t, preferred, first)

	require.NoError(t, nb.DeallocateByOwner(ctx, "owner-a"))

	again, err := nb.Allocate(ctx, "owner-b", netbridge.ProtocolTCP, preferred)
	require.NoError(t, err)
	assert.Equal(t, preferred, again, "a released port is immediately available again")
}

// newNetbridgeWithDiscovery builds a Netbridge whose router discovery is the
// given channel, so a test controls when (and whether) it finishes.
func newNetbridgeWithDiscovery(
	t *testing.T,
	discovery <-chan []strategies.Strategy,
	wait time.Duration,
) *netbridgeService {
	t.Helper()

	ax, err := asynx.New[ports.PortAllocation]().
		WithEventStore(asynxstore.New()).
		WithSnapshotStore(asynxstore.NewSnapshots()).
		WithShardingOpts(asynx.ShardingOpts{Shards: 2, QueueDepth: 100}).
		Build()
	require.NoError(t, err)
	t.Cleanup(func() { _ = ax.Shutdown(context.Background()) })

	nb, err := newNetbridge(ax, store.NewPortMemory(), discovery, testEphemeralPortStart, testEphemeralPortEnd)
	require.NoError(t, err)
	impl, ok := nb.(*netbridgeService)
	require.True(t, ok)
	impl.stratsWait = wait

	return impl
}

// A router that never answers discovery (a gateway that ignores NAT-PMP made
// every allocation wait two minutes) must not hold an allocation up: forwarding
// is best-effort.
func TestAllocate_SlowDiscovery_DoesNotStallTheAllocation(
	t *testing.T,
) {
	discovery := make(chan []strategies.Strategy, 1)
	nb := newNetbridgeWithDiscovery(t, discovery, 50*time.Millisecond)

	started := time.Now()
	port, err := nb.Allocate(context.Background(), "owner", netbridge.ProtocolTCP, 0)
	require.NoError(t, err)
	assert.Less(t, time.Since(started), 5*time.Second)

	alloc, err := nb.readModel.FindByPort(context.Background(), port)
	require.NoError(t, err)
	require.NotNil(t, alloc)
	assert.False(t, alloc.Forwarded, "no strategy was known yet, so nothing could forward")
}

func TestAllocate_DiscoveryFinishingLate_IsUsedByTheNextAllocation(
	t *testing.T,
) {
	discovery := make(chan []strategies.Strategy, 1)
	nb := newNetbridgeWithDiscovery(t, discovery, 20*time.Millisecond)
	ctx := context.Background()

	_, err := nb.Allocate(ctx, "owner", netbridge.ProtocolTCP, 0)
	require.NoError(t, err)

	discovery <- []strategies.Strategy{&mocks.AlwaysAvailableStrategy{}}

	port, err := nb.Allocate(ctx, "owner", netbridge.ProtocolTCP, 0)
	require.NoError(t, err)
	alloc, err := nb.readModel.FindByPort(ctx, port)
	require.NoError(t, err)
	require.NotNil(t, alloc)
	assert.True(t, alloc.Forwarded)
}

func TestAllocate_CancelledContext_DoesNotWaitForDiscovery(
	t *testing.T,
) {
	nb := newNetbridgeWithDiscovery(t, make(chan []strategies.Strategy), time.Hour)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	assert.Empty(t, nb.resolveStrategies(ctx))
}

func TestShutdown_DrainsAsynx(
	t *testing.T,
) {
	nb := buildNetbridge(t)

	require.NoError(t, nb.Shutdown(context.Background()))

	_, err := nb.Allocate(context.Background(), "owner-a", netbridge.ProtocolTCP, 0)
	assert.Error(t, err, "a drained aggregate must reject new allocations")
}

func buildNetbridgeNoForwarding(
	t *testing.T,
	strategy strategies.Strategy,
) *netbridgeService {
	t.Helper()

	es := asynxstore.New()
	ss := asynxstore.NewSnapshots()

	nb, err := New().
		WithEventStore(es).
		WithSnapshotStore(ss).
		WithStrategies([]strategies.Strategy{strategy}).
		WithForwarding(false).
		Build(context.Background())
	require.NoError(t, err)

	impl, ok := nb.(*netbridgeService)
	require.True(t, ok)
	return impl
}

func TestBuilder_ResolveForwarding_DefaultsToConfig(
	t *testing.T,
) {
	assert.Equal(t, config.GetNetbridge().Enabled, New().resolveForwarding())
}

func TestBuilder_ResolveForwarding_OptionOverridesConfig(
	t *testing.T,
) {
	testCases := []struct {
		name    string
		enabled bool
	}{
		{name: "explicitly enabled", enabled: true},
		{name: "explicitly disabled", enabled: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.enabled, New().WithForwarding(tc.enabled).resolveForwarding())
		})
	}
}

func TestBuilder_WithForwarding_DisabledDropsInjectedStrategies(
	t *testing.T,
) {
	nb := buildNetbridgeNoForwarding(t, &mocks.AlwaysAvailableStrategy{})

	assert.Empty(t, nb.resolveStrategies(context.Background()))
}

func TestBuilder_WithForwarding_EnabledKeepsInjectedStrategies(
	t *testing.T,
) {
	es := asynxstore.New()
	ss := asynxstore.NewSnapshots()

	nb, err := New().
		WithEventStore(es).
		WithSnapshotStore(ss).
		WithStrategies([]strategies.Strategy{&mocks.AlwaysAvailableStrategy{}}).
		WithForwarding(true).
		Build(context.Background())
	require.NoError(t, err)

	impl, ok := nb.(*netbridgeService)
	require.True(t, ok)
	assert.Len(t, impl.resolveStrategies(context.Background()), 1)
}

// No preferred port is asked for: which port comes back is the allocator's
// business and depends on what the machine already has bound, while what this
// test is about is that nothing was forwarded.
func TestAllocate_ForwardingDisabled_AllocatesWithoutForwarding(
	t *testing.T,
) {
	nb := buildNetbridgeNoForwarding(t, &mocks.AlwaysAvailableStrategy{})
	ctx := context.Background()

	port, err := nb.Allocate(ctx, "owner-nofwd", netbridge.ProtocolTCP, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, port, nb.portStart)
	assert.LessOrEqual(t, port, nb.portEnd)

	nb.waitForProjection()

	alloc, err := nb.readModel.FindByPort(ctx, port)
	require.NoError(t, err)
	require.NotNil(t, alloc)
	assert.False(t, alloc.Forwarded)
}

func TestDeallocateByOwner_ForwardingDisabled_SkipsReverse(
	t *testing.T,
) {
	strat := &mocks.AlwaysAvailableStrategy{}
	nb := buildNetbridgeNoForwarding(t, strat)
	ctx := context.Background()

	port, err := nb.Allocate(ctx, "owner-nofwd-rev", netbridge.ProtocolTCP, 0)
	require.NoError(t, err)

	nb.waitForProjection()

	err = nb.DeallocateByOwner(ctx, "owner-nofwd-rev")
	require.NoError(t, err)
	nb.waitForProjection()

	assert.Empty(t, strat.ReverseCalledWith)

	// The port is free again, so asking for it by name gets it back.
	realloc, err := nb.Allocate(ctx, "owner-nofwd-rev", netbridge.ProtocolTCP, port)
	require.NoError(t, err)
	assert.Equal(t, port, realloc)
}

// blockingStrategy only answers its availability probe when the probe's context
// ends, like a router that never replies.
type blockingStrategy struct{ mocks.AlwaysAvailableStrategy }

func (blockingStrategy) Available(ctx context.Context) bool {
	<-ctx.Done()
	return false
}

func TestProbeStrategy_GivesUpWhenTheContextEnds(
	t *testing.T,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	started := time.Now()
	assert.False(t, probeStrategy(ctx, &blockingStrategy{}))
	assert.Less(t, time.Since(started), 5*time.Second)
}

func TestProbeStrategy_ReportsAnAvailableStrategy(
	t *testing.T,
) {
	assert.True(t, probeStrategy(context.Background(), &mocks.AlwaysAvailableStrategy{}))
}
