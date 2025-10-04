package common

import (
	"testing"
)

// TestCRUDInterface tests the CRUD interface definition
func TestCRUDInterface(t *testing.T) {
	// This test verifies that the CRUD interface is properly defined
	// We can't easily test the interface without a concrete implementation,
	// but we can verify the interface definition exists and is valid
	
	// Test that the interface can be referenced
	var _ CRUD[string] = (*mockCRUD[string])(nil)
	var _ CRUD[int] = (*mockCRUD[int])(nil)
	var _ CRUD[bool] = (*mockCRUD[bool])(nil)
}

// mockCRUD is a mock implementation for testing the CRUD interface
type mockCRUD[T any] struct {
	items []T
}

func newMockCRUD[T any]() *mockCRUD[T] {
	return &mockCRUD[T]{
		items: make([]T, 0),
	}
}

func (m *mockCRUD[T]) Get() []T {
	return m.items
}

func (m *mockCRUD[T]) GetById(id string) *T {
	// Mock implementation - just return nil
	return nil
}

func (m *mockCRUD[T]) Create(item *T) *T {
	if item != nil {
		m.items = append(m.items, *item)
	}
	return item
}

func (m *mockCRUD[T]) Update(item *T) *T {
	// Mock implementation - just return the item
	return item
}

func (m *mockCRUD[T]) DeleteById(id string) error {
	// Mock implementation - no error
	return nil
}

func TestCRUDInterfaceMethods(t *testing.T) {
	// Test with string type
	stringCRUD := newMockCRUD[string]()
	
	// Test Get method
	items := stringCRUD.Get()
	if items == nil {
		t.Error("Get() should return a slice, not nil")
	}
	
	// Test GetById method
	item := stringCRUD.GetById("test-id")
	if item != nil {
		t.Error("GetById() should return nil for mock implementation")
	}
	
	// Test Create method
	testItem := "test-item"
	created := stringCRUD.Create(&testItem)
	if created == nil {
		t.Error("Create() should return the created item")
	}
	if *created != testItem || created == nil {
		t.Error("Create() should return the same item")
	}
	
	// Test Update method
	updated := stringCRUD.Update(&testItem)
	if updated == nil {
		t.Error("Update() should return the updated item")
	}
	
	// Test DeleteById method
	err := stringCRUD.DeleteById("test-id")
	if err != nil {
		t.Errorf("DeleteById() should not return error: %v", err)
	}
}

func TestCRUDInterfaceWithInt(t *testing.T) {
	// Test with int type
	intCRUD := newMockCRUD[int]()
	
	// Test all methods with int type
	items := intCRUD.Get()
	if items == nil {
		t.Error("Get() should return a slice, not nil")
	}
	
	item := intCRUD.GetById("test-id")
	if item != nil {
		t.Error("GetById() should return nil for mock implementation")
	}
	
	testItem := 42
	created := intCRUD.Create(&testItem)
	if created == nil {
		t.Error("Create() should return the created item")
	}
	if *created != testItem || created == nil {
		t.Error("Create() should return the same item")
	}
	
	updated := intCRUD.Update(&testItem)
	if updated == nil {
		t.Error("Update() should return the updated item")
	}
	
	err := intCRUD.DeleteById("test-id")
	if err != nil {
		t.Errorf("DeleteById() should not return error: %v", err)
	}
}

func TestCRUDInterfaceWithBool(t *testing.T) {
	// Test with bool type
	boolCRUD := newMockCRUD[bool]()
	
	// Test all methods with bool type
	items := boolCRUD.Get()
	if items == nil {
		t.Error("Get() should return a slice, not nil")
	}
	
	item := boolCRUD.GetById("test-id")
	if item != nil {
		t.Error("GetById() should return nil for mock implementation")
	}
	
	testItem := true
	created := boolCRUD.Create(&testItem)
	if created == nil {
		t.Error("Create() should return the created item")
	}
	if *created != testItem {
		t.Error("Create() should return the same item")
	}
	
	updated := boolCRUD.Update(&testItem)
	if updated == nil {
		t.Error("Update() should return the updated item")
	}
	
	err := boolCRUD.DeleteById("test-id")
	if err != nil {
		t.Errorf("DeleteById() should not return error: %v", err)
	}
}

func TestCRUDInterfaceGenericBehavior(t *testing.T) {
	// Test that the interface works with different generic types
	testCases := []struct {
		name string
		crud CRUD[string]
	}{
		{"string CRUD", newMockCRUD[string]()},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test that all methods can be called
			items := tc.crud.Get()
			_ = items
			
			item := tc.crud.GetById("test")
			_ = item
			
			testItem := "test"
			created := tc.crud.Create(&testItem)
			_ = created
			
			updated := tc.crud.Update(&testItem)
			_ = updated
			
			err := tc.crud.DeleteById("test")
			if err != nil {
				t.Errorf("DeleteById() returned error: %v", err)
			}
		})
	}
}
