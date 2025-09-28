package repositories

import (
	"testing"

	"github.com/rabbytesoftware/quiver/internal/infrastructure"
)

func TestNewRepositories(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := NewRepositories(infra)

	if repos == nil {
		t.Fatal("NewRepositories() returned nil")
	}

	// Test that all repositories are initialized
	if repos.arrows == nil {
		t.Error("arrows repository is not initialized")
	}

	if repos.system == nil {
		t.Error("system repository is not initialized")
	}

	if repos.quivers == nil {
		t.Error("quivers repository is not initialized")
	}
}

func TestRepositories_GetArrows(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := NewRepositories(infra)
	arrowsRepo := repos.GetArrows()

	if arrowsRepo == nil {
		t.Error("GetArrows() returned nil")
	}

	// Test that it returns the same instance
	arrowsRepo2 := repos.GetArrows()
	if arrowsRepo != arrowsRepo2 {
		t.Error("GetArrows() should return the same instance")
	}

	// Test interface compliance (arrowsRepo is already ArrowsInterface type)
	_ = arrowsRepo
}

func TestRepositories_GetSystem(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := NewRepositories(infra)
	systemRepo := repos.GetSystem()

	if systemRepo == nil {
		t.Error("GetSystem() returned nil")
	}

	// Test that it returns the same instance
	systemRepo2 := repos.GetSystem()
	if systemRepo != systemRepo2 {
		t.Error("GetSystem() should return the same instance")
	}

	// Test interface compliance (systemRepo is already SystemInterface type)
	_ = systemRepo
}

func TestRepositories_GetQuivers(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := NewRepositories(infra)
	quiversRepo := repos.GetQuivers()

	if quiversRepo == nil {
		t.Error("GetQuivers() returned nil")
	}

	// Test that it returns the same instance
	quiversRepo2 := repos.GetQuivers()
	if quiversRepo != quiversRepo2 {
		t.Error("GetQuivers() should return the same instance")
	}

	// Test interface compliance (quiversRepo is already QuiversInterface type)
	_ = quiversRepo
}

func TestRepositoriesStructure(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := NewRepositories(infra)

	// Test that Repositories struct has the expected fields
	if repos.arrows == nil {
		t.Error("Repositories.arrows field is nil")
	}

	if repos.system == nil {
		t.Error("Repositories.system field is nil")
	}

	if repos.quivers == nil {
		t.Error("Repositories.quivers field is nil")
	}
}

func TestRepositoriesDependencyInjection(t *testing.T) {
	// Test that the same infrastructure instance is passed to all repositories
	infra := infrastructure.NewInfrastructure()
	repos := NewRepositories(infra)

	// We can't directly test that the same infrastructure is used
	// but we can verify that all repositories are created successfully
	if repos.GetArrows() == nil {
		t.Error("arrows repository was not created with infrastructure")
	}

	if repos.GetSystem() == nil {
		t.Error("system repository was not created with infrastructure")
	}

	if repos.GetQuivers() == nil {
		t.Error("quivers repository was not created with infrastructure")
	}
}

func TestMultipleRepositoriesInstances(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos1 := NewRepositories(infra)
	repos2 := NewRepositories(infra)

	// Different Repositories instances
	if repos1 == repos2 {
		t.Error("NewRepositories() should create new instances each time")
	}

	// But they should have different repository instances too
	if repos1.GetArrows() == repos2.GetArrows() {
		t.Error("Different Repositories instances should have different arrows repositories")
	}

	if repos1.GetSystem() == repos2.GetSystem() {
		t.Error("Different Repositories instances should have different system repositories")
	}

	if repos1.GetQuivers() == repos2.GetQuivers() {
		t.Error("Different Repositories instances should have different quivers repositories")
	}
}

func TestRepositoriesWithNilInfrastructure(t *testing.T) {
	// Test behavior when nil infrastructure is passed
	defer func() {
		if r := recover(); r != nil {
			// If it panics, that's acceptable behavior for nil infrastructure
			t.Logf("NewRepositories panicked with nil infrastructure: %v", r)
		}
	}()

	repos := NewRepositories(nil)

	// If it doesn't panic, verify the repositories are still created
	if repos != nil {
		if repos.GetArrows() == nil {
			t.Error("arrows repository is nil with nil infrastructure")
		}
		if repos.GetSystem() == nil {
			t.Error("system repository is nil with nil infrastructure")
		}
		if repos.GetQuivers() == nil {
			t.Error("quivers repository is nil with nil infrastructure")
		}
	}
}
