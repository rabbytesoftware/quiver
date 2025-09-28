package system

import (
	"testing"

	"github.com/rabbytesoftware/quiver/internal/infrastructure"
	"github.com/rabbytesoftware/quiver/internal/repositories"
)

func TestNewSystemUsecase(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)

	usecase := NewSystemUsecase(repos)

	if usecase == nil {
		t.Fatal("NewSystemUsecase() returned nil")
	}

	if usecase.repositories != repos {
		t.Error("NewSystemUsecase() did not set repositories correctly")
	}
}

func TestNewSystemUsecase_WithNilRepositories(t *testing.T) {
	usecase := NewSystemUsecase(nil)

	if usecase == nil {
		t.Fatal("NewSystemUsecase() returned nil")
	}

	if usecase.repositories != nil {
		t.Error("NewSystemUsecase() should accept nil repositories")
	}
}

func TestSystemUsecase_Structure(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)

	usecase := NewSystemUsecase(repos)

	// Test that the struct has the expected fields
	if usecase.repositories == nil {
		t.Error("repositories field should not be nil when passed valid repositories")
	}
}

func TestSystemUsecase_MultipleInstances(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)

	usecase1 := NewSystemUsecase(repos)
	usecase2 := NewSystemUsecase(repos)

	// Should create different instances
	if usecase1 == usecase2 {
		t.Error("NewSystemUsecase() should create different instances")
	}

	// But both should reference the same repositories
	if usecase1.repositories != usecase2.repositories {
		t.Error("Both usecases should reference the same repositories instance")
	}
}

func TestSystemUsecase_Type(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repos := repositories.NewRepositories(infra)

	usecase := NewSystemUsecase(repos)

	// Verify the concrete type (usecase is already *SystemUsecase)
	_ = usecase
}
