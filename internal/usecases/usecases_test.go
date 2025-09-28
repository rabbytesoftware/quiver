package usecases

import (
	"testing"

	"github.com/rabbytesoftware/quiver/internal/infrastructure"
	"github.com/rabbytesoftware/quiver/internal/repositories"
)

func TestNewUsecases(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)
	usecases := NewUsecases(repos)

	if usecases == nil {
		t.Fatal("NewUsecases() returned nil")
	}

	// Test that all usecases are initialized
	if usecases.Arrows == nil {
		t.Error("Arrows usecase is not initialized")
	}

	if usecases.Quivers == nil {
		t.Error("Quivers usecase is not initialized")
	}

	if usecases.System == nil {
		t.Error("System usecase is not initialized")
	}
}

func TestUsecasesStructure(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)
	usecases := NewUsecases(repos)

	// Test that Usecases struct has the expected fields
	if usecases.Arrows == nil {
		t.Error("Arrows usecase is nil")
	}

	if usecases.Quivers == nil {
		t.Error("Quivers usecase is nil")
	}

	if usecases.System == nil {
		t.Error("System usecase is nil")
	}
}

func TestUsecasesDependencyInjection(t *testing.T) {
	// Test that the same repositories instance is passed to all usecases
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)
	usecases := NewUsecases(repos)

	// We can't directly test that the same repositories is used
	// but we can verify that all usecases are created successfully
	if usecases.Arrows == nil {
		t.Error("Arrows usecase was not created with repositories")
	}

	if usecases.Quivers == nil {
		t.Error("Quivers usecase was not created with repositories")
	}

	if usecases.System == nil {
		t.Error("System usecase was not created with repositories")
	}
}

func TestMultipleUsecasesInstances(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)
	usecases1 := NewUsecases(repos)
	usecases2 := NewUsecases(repos)

	// Different Usecases instances
	if usecases1 == usecases2 {
		t.Error("NewUsecases() should create new instances each time")
	}

	// But they should have different usecase instances too
	if usecases1.Arrows == usecases2.Arrows {
		t.Error("Different Usecases instances should have different Arrows usecases")
	}

	if usecases1.Quivers == usecases2.Quivers {
		t.Error("Different Usecases instances should have different Quivers usecases")
	}

	if usecases1.System == usecases2.System {
		t.Error("Different Usecases instances should have different System usecases")
	}
}

func TestUsecasesWithNilRepositories(t *testing.T) {
	// Test behavior when nil repositories is passed
	defer func() {
		if r := recover(); r != nil {
			// If it panics, that's acceptable behavior for nil repositories
			t.Logf("NewUsecases panicked with nil repositories: %v", r)
		}
	}()

	usecases := NewUsecases(nil)

	// If it doesn't panic, verify the usecases are still created
	if usecases != nil {
		if usecases.Arrows == nil {
			t.Error("Arrows usecase is nil with nil repositories")
		}
		if usecases.Quivers == nil {
			t.Error("Quivers usecase is nil with nil repositories")
		}
		if usecases.System == nil {
			t.Error("System usecase is nil with nil repositories")
		}
	}
}

func TestUsecasesFieldAccess(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)
	usecases := NewUsecases(repos)

	// Test that we can access all public fields
	arrows := usecases.Arrows
	if arrows == nil {
		t.Error("Cannot access Arrows field")
	}

	quivers := usecases.Quivers
	if quivers == nil {
		t.Error("Cannot access Quivers field")
	}

	system := usecases.System
	if system == nil {
		t.Error("Cannot access System field")
	}

	// Test that accessing the same field returns the same instance
	arrows2 := usecases.Arrows
	if arrows != arrows2 {
		t.Error("Arrows field should return the same instance")
	}
}
