package quivers

import (
	"testing"

	"github.com/rabbytesoftware/quiver/internal/infrastructure"
	"github.com/rabbytesoftware/quiver/internal/repositories"
)

func TestNewQuiversUsecase(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)

	usecase := NewQuiversUsecase(repos)

	if usecase == nil {
		t.Fatal("NewQuiversUsecase() returned nil")
	}

	if usecase.repositories != repos {
		t.Error("NewQuiversUsecase() did not set repositories correctly")
	}
}

func TestNewQuiversUsecase_WithNilRepositories(t *testing.T) {
	usecase := NewQuiversUsecase(nil)

	if usecase == nil {
		t.Fatal("NewQuiversUsecase() returned nil")
	}

	if usecase.repositories != nil {
		t.Error("NewQuiversUsecase() should accept nil repositories")
	}
}

func TestQuiversUsecase_Structure(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)

	usecase := NewQuiversUsecase(repos)

	// Test that the struct has the expected fields
	if usecase.repositories == nil {
		t.Error("repositories field should not be nil when passed valid repositories")
	}
}

func TestQuiversUsecase_MultipleInstances(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)

	usecase1 := NewQuiversUsecase(repos)
	usecase2 := NewQuiversUsecase(repos)

	// Should create different instances
	if usecase1 == usecase2 {
		t.Error("NewQuiversUsecase() should create different instances")
	}

	// But both should reference the same repositories
	if usecase1.repositories != usecase2.repositories {
		t.Error("Both usecases should reference the same repositories instance")
	}
}

func TestQuiversUsecase_Type(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)

	usecase := NewQuiversUsecase(repos)

	// Verify the concrete type (usecase is already *QuiversUsecase)
	_ = usecase
}
