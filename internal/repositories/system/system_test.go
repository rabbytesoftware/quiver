package system

import (
	"testing"

	"github.com/rabbytesoftware/quiver/internal/infrastructure"
)

func TestNewSystemRepository(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewSystemRepository(infra)

	if repo == nil {
		t.Fatal("NewSystemRepository() returned nil")
	}

	// Test interface compliance (repo is already SystemInterface type)
	_ = repo

	// Test concrete type
	if _, ok := repo.(*SystemRepository); !ok {
		t.Error("NewSystemRepository() did not return *SystemRepository")
	}
}

func TestSystemRepository_GetMetadata(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewSystemRepository(infra)

	metadata := repo.GetMetadata()

	// Current implementation returns nil
	if metadata != nil {
		t.Error("GetMetadata() expected nil, got non-nil result")
	}
}

func TestSystemRepository_UpdateQuiver(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewSystemRepository(infra)

	err := repo.UpdateQuiver()

	// Current implementation returns nil (no error)
	if err != nil {
		t.Errorf("UpdateQuiver() expected nil error, got %v", err)
	}
}

func TestSystemRepository_UninstallQuiver(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewSystemRepository(infra)

	err := repo.UninstallQuiver()

	// Current implementation returns nil (no error)
	if err != nil {
		t.Errorf("UninstallQuiver() expected nil error, got %v", err)
	}
}

func TestSystemRepository_GetLogs(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewSystemRepository(infra)

	logs := repo.GetLogs()

	// Current implementation returns empty string
	if logs != "" {
		t.Errorf("GetLogs() expected empty string, got %q", logs)
	}
}

func TestSystemRepository_RestartQuiver(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewSystemRepository(infra)

	err := repo.RestartQuiver()

	// Current implementation returns nil (no error)
	if err != nil {
		t.Errorf("RestartQuiver() expected nil error, got %v", err)
	}
}

func TestSystemRepository_Status(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewSystemRepository(infra)

	status := repo.Status()

	// Current implementation returns empty string
	if status != "" {
		t.Errorf("Status() expected empty string, got %q", status)
	}
}

func TestSystemRepository_StopQuiver(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewSystemRepository(infra)

	err := repo.StopQuiver()

	// Current implementation returns nil (no error)
	if err != nil {
		t.Errorf("StopQuiver() expected nil error, got %v", err)
	}
}

func TestSystemRepositoryStructure(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewSystemRepository(infra).(*SystemRepository)

	// Test that infrastructure is properly set
	if repo.infrastructure == nil {
		t.Error("SystemRepository.infrastructure is nil")
	}

	if repo.infrastructure != infra {
		t.Error("SystemRepository.infrastructure is not the same instance passed to constructor")
	}
}

func TestSystemRepositoryWithNilInfrastructure(t *testing.T) {
	// Test behavior when nil infrastructure is passed
	repo := NewSystemRepository(nil)

	if repo == nil {
		t.Fatal("NewSystemRepository(nil) returned nil")
	}

	// Test that methods still work with nil infrastructure
	metadata := repo.GetMetadata()
	if metadata != nil {
		t.Error("GetMetadata() with nil infrastructure should return nil")
	}

	logs := repo.GetLogs()
	if logs != "" {
		t.Error("GetLogs() with nil infrastructure should return empty string")
	}

	status := repo.Status()
	if status != "" {
		t.Error("Status() with nil infrastructure should return empty string")
	}

	// Test error-returning methods
	if err := repo.UpdateQuiver(); err != nil {
		t.Error("UpdateQuiver() with nil infrastructure should return nil error")
	}

	if err := repo.UninstallQuiver(); err != nil {
		t.Error("UninstallQuiver() with nil infrastructure should return nil error")
	}

	if err := repo.RestartQuiver(); err != nil {
		t.Error("RestartQuiver() with nil infrastructure should return nil error")
	}

	if err := repo.StopQuiver(); err != nil {
		t.Error("StopQuiver() with nil infrastructure should return nil error")
	}
}

func TestSystemRepositoryInterfaceMethods(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	var repo SystemInterface = NewSystemRepository(infra)

	// Test that all interface methods are callable
	metadata := repo.GetMetadata()
	if metadata != nil {
		t.Error("Interface GetMetadata() method should return nil")
	}

	logs := repo.GetLogs()
	if logs != "" {
		t.Error("Interface GetLogs() method should return empty string")
	}

	status := repo.Status()
	if status != "" {
		t.Error("Interface Status() method should return empty string")
	}

	if err := repo.UpdateQuiver(); err != nil {
		t.Error("Interface UpdateQuiver() method should return nil error")
	}

	if err := repo.UninstallQuiver(); err != nil {
		t.Error("Interface UninstallQuiver() method should return nil error")
	}

	if err := repo.RestartQuiver(); err != nil {
		t.Error("Interface RestartQuiver() method should return nil error")
	}

	if err := repo.StopQuiver(); err != nil {
		t.Error("Interface StopQuiver() method should return nil error")
	}
}
