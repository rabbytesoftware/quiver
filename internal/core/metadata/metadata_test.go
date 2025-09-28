package metadata

import (
	"testing"
)

func TestGet(t *testing.T) {
	metadata := Get()

	if metadata == nil {
		t.Fatal("Get() returned nil")
	}

	// Test that it returns the same instance on subsequent calls (singleton pattern)
	metadata2 := Get()
	if metadata != metadata2 {
		t.Error("Get() should return the same instance (singleton)")
	}
}

func TestGetVersion(t *testing.T) {
	version := GetVersion()

	if version == "" {
		t.Error("GetVersion() returned empty string")
	}
}

func TestGetVersionCodename(t *testing.T) {
	codename := GetVersionCodename()

	if codename == "" {
		t.Error("GetVersionCodename() returned empty string")
	}
}

func TestGetName(t *testing.T) {
	name := GetName()

	if name == "" {
		t.Error("GetName() returned empty string")
	}
}

func TestGetDescription(t *testing.T) {
	description := GetDescription()

	if description == "" {
		t.Error("GetDescription() returned empty string")
	}
}

func TestGetAuthor(t *testing.T) {
	author := GetAuthor()

	if author == "" {
		t.Error("GetAuthor() returned empty string")
	}
}

func TestGetURL(t *testing.T) {
	url := GetURL()

	if url == "" {
		t.Error("GetURL() returned empty string")
	}
}

func TestGetLicense(t *testing.T) {
	license := GetLicense()

	if license == "" {
		t.Error("GetLicense() returned empty string")
	}
}

func TestGetCopyright(t *testing.T) {
	copyright := GetCopyright()

	if copyright == "" {
		t.Error("GetCopyright() returned empty string")
	}
}

func TestGetMaintainers(t *testing.T) {
	maintainers := GetMaintainers()

	// Maintainers can be empty, so we just test it doesn't panic
	_ = maintainers
}

func TestGetVariables(t *testing.T) {
	variables := GetVariables()

	// Variables can be empty, so we just test it doesn't panic
	_ = variables
}

func TestGetDefaultConfigPath(t *testing.T) {
	path := GetDefaultConfigPath()

	if path == "" {
		t.Error("GetDefaultConfigPath() returned empty string")
	}
}

func TestMetadataStructure(t *testing.T) {
	metadata := Get()

	// Test that metadata has all expected fields
	if metadata.Version.Number == "" {
		t.Error("Metadata.Version.Number is empty")
	}

	if metadata.Version.Codename == "" {
		t.Error("Metadata.Version.Codename is empty")
	}

	if metadata.Metadata.Name == "" {
		t.Error("Metadata.Metadata.Name is empty")
	}

	if metadata.Metadata.Description == "" {
		t.Error("Metadata.Metadata.Description is empty")
	}

	if metadata.Metadata.Author == "" {
		t.Error("Metadata.Metadata.Author is empty")
	}

	if metadata.Metadata.URL == "" {
		t.Error("Metadata.Metadata.URL is empty")
	}

	if metadata.Metadata.License == "" {
		t.Error("Metadata.Metadata.License is empty")
	}

	if metadata.Metadata.Copyright == "" {
		t.Error("Metadata.Metadata.Copyright is empty")
	}
}

func TestDefaultMetadata(t *testing.T) {
	// Test that defaultMetadata function works
	defaultMeta := defaultMetadata()

	if defaultMeta == nil {
		t.Fatal("defaultMetadata() returned nil")
	}

	// Test that default metadata has reasonable values
	if defaultMeta.Metadata.Name == "" {
		t.Error("Default metadata Name is empty")
	}

	if defaultMeta.Version.Number == "" {
		t.Error("Default metadata Version is empty")
	}
}

func TestMetadataConsistency(t *testing.T) {
	// Test that all getter functions return consistent values with the metadata struct
	metadata := Get()

	if GetVersion() != metadata.Version.Number {
		t.Error("GetVersion() inconsistent with metadata.Version.Number")
	}

	if GetVersionCodename() != metadata.Version.Codename {
		t.Error("GetVersionCodename() inconsistent with metadata.Version.Codename")
	}

	if GetName() != metadata.Metadata.Name {
		t.Error("GetName() inconsistent with metadata.Metadata.Name")
	}

	if GetDescription() != metadata.Metadata.Description {
		t.Error("GetDescription() inconsistent with metadata.Metadata.Description")
	}

	if GetAuthor() != metadata.Metadata.Author {
		t.Error("GetAuthor() inconsistent with metadata.Metadata.Author")
	}

	if GetURL() != metadata.Metadata.URL {
		t.Error("GetURL() inconsistent with metadata.Metadata.URL")
	}

	if GetLicense() != metadata.Metadata.License {
		t.Error("GetLicense() inconsistent with metadata.Metadata.License")
	}

	if GetCopyright() != metadata.Metadata.Copyright {
		t.Error("GetCopyright() inconsistent with metadata.Metadata.Copyright")
	}
}
