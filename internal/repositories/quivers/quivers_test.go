package quivers

import (
	"testing"

	"github.com/rabbytesoftware/quiver/internal/infrastructure"
	domain "github.com/rabbytesoftware/quiver/internal/models/quiver"
)

func TestNewQuiversRepository(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewQuiversRepository(infra)

	if repo == nil {
		t.Fatal("NewQuiversRepository() returned nil")
	}

	// Test interface compliance (repo is already QuiversInterface type)
	_ = repo

	// Test concrete type
	if _, ok := repo.(*QuiversRepository); !ok {
		t.Error("NewQuiversRepository() did not return *QuiversRepository")
	}
}

func TestNewQuiversRepository_WithNilInfrastructure(t *testing.T) {
	repo := NewQuiversRepository(nil)

	if repo == nil {
		t.Fatal("NewQuiversRepository() returned nil")
	}

	// Should still work with nil infrastructure
	if concreteRepo, ok := repo.(*QuiversRepository); ok {
		if concreteRepo.infrastructure != nil {
			t.Error("Expected infrastructure to be nil")
		}
	}
}

func TestQuiversRepository_Get(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewQuiversRepository(infra)

	result := repo.Get()

	// Should return empty slice (current implementation)
	if result == nil {
		t.Error("Get() should not return nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %d items", len(result))
	}
}

func TestQuiversRepository_GetById(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewQuiversRepository(infra)

	result := repo.GetById("test-id")

	// Should return nil (current implementation)
	if result != nil {
		t.Error("GetById() should return nil in current implementation")
	}
}

func TestQuiversRepository_GetById_EmptyId(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewQuiversRepository(infra)

	result := repo.GetById("")

	// Should return nil (current implementation)
	if result != nil {
		t.Error("GetById() with empty id should return nil")
	}
}

func TestQuiversRepository_Create(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewQuiversRepository(infra)

	quiver := &domain.Quiver{}
	result := repo.Create(quiver)

	// Should return nil (current implementation)
	if result != nil {
		t.Error("Create() should return nil in current implementation")
	}
}

func TestQuiversRepository_Create_WithNilQuiver(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewQuiversRepository(infra)

	result := repo.Create(nil)

	// Should return nil (current implementation)
	if result != nil {
		t.Error("Create() with nil quiver should return nil")
	}
}

func TestQuiversRepository_Update(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewQuiversRepository(infra)

	quiver := &domain.Quiver{}
	result := repo.Update(quiver)

	// Should return nil (current implementation)
	if result != nil {
		t.Error("Update() should return nil in current implementation")
	}
}

func TestQuiversRepository_Update_WithNilQuiver(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewQuiversRepository(infra)

	result := repo.Update(nil)

	// Should return nil (current implementation)
	if result != nil {
		t.Error("Update() with nil quiver should return nil")
	}
}

func TestQuiversRepository_DeleteById(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewQuiversRepository(infra)

	err := repo.DeleteById("test-id")

	// Should return nil (current implementation)
	if err != nil {
		t.Errorf("DeleteById() should return nil, got %v", err)
	}
}

func TestQuiversRepository_DeleteById_EmptyId(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewQuiversRepository(infra)

	err := repo.DeleteById("")

	// Should return nil (current implementation)
	if err != nil {
		t.Errorf("DeleteById() with empty id should return nil, got %v", err)
	}
}

func TestQuiversRepository_AllMethods(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewQuiversRepository(infra)

	// Test all methods in sequence to ensure they don't interfere
	quivers := repo.Get()
	if len(quivers) != 0 {
		t.Error("Expected empty quivers list")
	}

	quiver := repo.GetById("test")
	if quiver != nil {
		t.Error("Expected nil quiver")
	}

	created := repo.Create(&domain.Quiver{})
	if created != nil {
		t.Error("Expected nil created quiver")
	}

	updated := repo.Update(&domain.Quiver{})
	if updated != nil {
		t.Error("Expected nil updated quiver")
	}

	err := repo.DeleteById("test")
	if err != nil {
		t.Error("Expected nil error from delete")
	}
}
