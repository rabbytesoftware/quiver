//go:build darwin || linux

package process

import (
	"context"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sys/unix"

	"github.com/rabbytesoftware/quiver.core/internal/engine/wizard/internal/runtime/internal/models"
)

func TestUnixProcess_RunsInItsOwnSession(t *testing.T) {
	proc, pid, _ := startForkingShellProcess(t, context.Background(), 5*time.Second)

	childSID, err := unix.Getsid(pid)
	if err != nil {
		t.Fatalf("Getsid(%d) error = %v", pid, err)
	}
	if parentSID, _ := unix.Getsid(0); childSID == parentSID {
		t.Errorf("child session %d equals the daemon's: it still shares the controlling terminal", childSID)
	}
	if childSID != pid {
		t.Errorf("child session = %d, want its own pid %d", childSID, pid)
	}

	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		t.Fatalf("Getpgid(%d) error = %v", pid, err)
	}
	if pgid != pid {
		t.Errorf("pgid = %d, want %d: the child must lead its own group", pgid, pid)
	}
	_ = proc
}

func TestUnixProcess_Stop_OwnSession_TakesTheGrandchild(t *testing.T) {
	config := models.NewConfig([]string{"sleep 300 & echo $!; wait"})
	config.ShellWrap = true
	config.StopTimeout = 5 * time.Second
	proc, err := newProcess(context.Background(), config)
	if err != nil {
		t.Fatalf("newProcess() error = %v", err)
	}
	t.Cleanup(func() {
		_ = syscall.Kill(-proc.PID(), syscall.SIGKILL)
		_ = proc.Close()
	})

	var line string
	select {
	case line = <-proc.StreamOutput():
	case <-time.After(5 * time.Second):
		t.Fatal("shell never reported its background child's PID")
	}
	grandchild, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		t.Fatalf("grandchild PID %q: %v", line, err)
	}

	if err := proc.Stop(context.Background()); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for isAlive(grandchild) && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	if isAlive(grandchild) {
		t.Errorf("grandchild %d survived Stop of its session leader", grandchild)
	}
}
