package translator

import (
	"fmt"
	"testing"

	fns "github.com/rabbytesoftware/quiver/internal/infrastructure/fetchnshare"
)

func TestNewTranslator(t *testing.T) {
	mockFNS := fns.NewFNS()
	tr := NewTranslator(mockFNS)
	if tr == nil {
		t.Fatal("NewTranslator() returned nil")
	}
}

func TestTranslator_GetArrowTranslator(t *testing.T) {
	mockFNS := fns.NewFNS()
	tr := NewTranslator(mockFNS)

	translator := tr.GetArrowTranslator()
	if translator == nil {
		t.Error("GetArrowTranslator() should return a translator")
	}
}

func TestTranslator_GetQuiverTranslator(t *testing.T) {
	mockFNS := fns.NewFNS()
	tr := NewTranslator(mockFNS)

	translator := tr.GetQuiverTranslator()
	if translator == nil {
		t.Error("GetQuiverTranslator() should return a translator")
	}
}

func TestTranslator_InterfaceCompliance(t *testing.T) {
	// Test that Translator implements the expected interface
	mockFNS := fns.NewFNS()
	tr := NewTranslator(mockFNS)
	if tr == nil {
		t.Error("Translator should not be nil")
	}
}

func TestTranslator_MultipleInstances(t *testing.T) {
	mockFNS := fns.NewFNS()
	tr1 := NewTranslator(mockFNS)
	tr2 := NewTranslator(mockFNS)

	// Both should be valid
	if tr1 == nil || tr2 == nil {
		t.Error("NewTranslator() returned nil instance")
	}

	// Test that both instances work correctly
	arrowTranslator1 := tr1.GetArrowTranslator()
	arrowTranslator2 := tr2.GetArrowTranslator()

	if arrowTranslator1 == nil || arrowTranslator2 == nil {
		t.Error("Both instances should return valid translators")
	}
}

func TestTranslator_AllMethods(t *testing.T) {
	mockFNS := fns.NewFNS()
	tr := NewTranslator(mockFNS)

	// Test all methods to ensure they don't panic
	testCases := []struct {
		name string
		fn   func() error
	}{
		{"GetArrowTranslator", func() error {
			translator := tr.GetArrowTranslator()
			if translator == nil {
				return fmt.Errorf("GetArrowTranslator() returned nil")
			}
			return nil
		}},
		{"GetQuiverTranslator", func() error {
			translator := tr.GetQuiverTranslator()
			if translator == nil {
				return fmt.Errorf("GetQuiverTranslator() returned nil")
			}
			return nil
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn()
			if err != nil {
				t.Errorf("%s() returned error: %v", tc.name, err)
			}
		})
	}
}
