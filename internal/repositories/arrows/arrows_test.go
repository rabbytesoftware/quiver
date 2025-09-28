package arrows

import (
	"testing"

	"github.com/rabbytesoftware/quiver/internal/infrastructure"
	domain "github.com/rabbytesoftware/quiver/internal/models/arrow"
)

func TestNewArrowsRepository(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewArrowsRepository(infra)

	if repo == nil {
		t.Fatal("NewArrowsRepository() returned nil")
	}

	// Test interface compliance (repo is already ArrowsInterface type)
	_ = repo

	// Test concrete type
	if _, ok := repo.(*ArrowsRepository); !ok {
		t.Error("NewArrowsRepository() did not return *ArrowsRepository")
	}
}

func TestArrowsRepository_Get(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewArrowsRepository(infra)

	arrows := repo.Get()

	// Test that it returns a slice (even if empty)
	if arrows == nil {
		t.Error("Get() returned nil instead of empty slice")
	}

	// Test that it returns a slice
	_ = arrows // arrows is already []domain.Arrow type

	// Current implementation returns empty slice
	if len(arrows) != 0 {
		t.Errorf("Expected empty slice, got %d arrows", len(arrows))
	}
}

func TestArrowsRepository_GetById(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewArrowsRepository(infra)

	// Test with various IDs
	testCases := []string{"", "test-id", "123", "non-existent"}

	for _, id := range testCases {
		arrow := repo.GetById(id)

		// Current implementation returns nil
		if arrow != nil {
			t.Errorf("GetById(%q) expected nil, got %v", id, arrow)
		}
	}
}

func TestArrowsRepository_Create(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewArrowsRepository(infra)

	// Test with nil arrow
	result := repo.Create(nil)
	if result != nil {
		t.Error("Create(nil) expected nil, got non-nil result")
	}

	// Test with valid arrow
	arrow := &domain.Arrow{}
	result = repo.Create(arrow)

	// Current implementation returns nil
	if result != nil {
		t.Error("Create() expected nil, got non-nil result")
	}
}

func TestArrowsRepository_Update(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewArrowsRepository(infra)

	// Test with nil arrow
	result := repo.Update(nil)
	if result != nil {
		t.Error("Update(nil) expected nil, got non-nil result")
	}

	// Test with valid arrow
	arrow := &domain.Arrow{}
	result = repo.Update(arrow)

	// Current implementation returns nil
	if result != nil {
		t.Error("Update() expected nil, got non-nil result")
	}
}

func TestArrowsRepository_DeleteById(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewArrowsRepository(infra)

	// Test with various IDs
	testCases := []string{"", "test-id", "123", "non-existent"}

	for _, id := range testCases {
		err := repo.DeleteById(id)

		// Current implementation returns nil (no error)
		if err != nil {
			t.Errorf("DeleteById(%q) expected nil error, got %v", id, err)
		}
	}
}

func TestArrowsRepositoryStructure(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	repo := NewArrowsRepository(infra).(*ArrowsRepository)

	// Test that infrastructure is properly set
	if repo.infrastructure == nil {
		t.Error("ArrowsRepository.infrastructure is nil")
	}

	if repo.infrastructure != infra {
		t.Error("ArrowsRepository.infrastructure is not the same instance passed to constructor")
	}
}

func TestArrowsRepositoryWithNilInfrastructure(t *testing.T) {
	// Test behavior when nil infrastructure is passed
	repo := NewArrowsRepository(nil)

	if repo == nil {
		t.Fatal("NewArrowsRepository(nil) returned nil")
	}

	// Test that methods still work with nil infrastructure
	arrows := repo.Get()
	if arrows == nil {
		t.Error("Get() with nil infrastructure returned nil")
	}

	arrow := repo.GetById("test")
	if arrow != nil {
		t.Error("GetById() with nil infrastructure should return nil")
	}

	err := repo.DeleteById("test")
	if err != nil {
		t.Error("DeleteById() with nil infrastructure should return nil error")
	}
}

func TestArrowsRepositoryInterfaceMethods(t *testing.T) {
	infra := infrastructure.NewInfrastructure()
	var repo ArrowsInterface = NewArrowsRepository(infra)

	// Test that all interface methods are callable
	arrows := repo.Get()
	if arrows == nil {
		t.Error("Interface Get() method returned nil")
	}

	arrow := repo.GetById("test")
	if arrow != nil {
		t.Error("Interface GetById() method should return nil")
	}

	result := repo.Create(&domain.Arrow{})
	if result != nil {
		t.Error("Interface Create() method should return nil")
	}

	result = repo.Update(&domain.Arrow{})
	if result != nil {
		t.Error("Interface Update() method should return nil")
	}

	err := repo.DeleteById("test")
	if err != nil {
		t.Error("Interface DeleteById() method should return nil error")
	}
}
