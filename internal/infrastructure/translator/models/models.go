package models

import (
	"context"

	"github.com/rabbytesoftware/quiver/internal/models/arrow"
	"github.com/rabbytesoftware/quiver/internal/models/quiver"
)

// TranslatorInterface defines the Translator interface
type TranslatorInterface interface {
	GetArrowTranslator() TranslatorLayerInterface[arrow.Arrow]
	GetQuiverTranslator() TranslatorLayerInterface[quiver.Quiver]
}

// TranslatorInterface defines the Translator Layer interface
// This interface provides the core functionality for Arrow and Quiver translation operations
type TranslatorLayerInterface[t any] interface {
	// IsCompatible checks if an Arrow manifest can be translated
	IsCompatible(ctx context.Context, manifestPath string) (bool, error)

	// Translate performs the complete translation magic from manifest to Arrow model
	Translate(ctx context.Context, manifestPath string) (*t, error)

	// GetManifestVersion extracts the version from any Arrow manifest
	GetManifestVersion(ctx context.Context, manifestPath string) (string, error)

	// GetSupportedVersions returns all supported Arrow versions
	GetSupportedVersions(ctx context.Context) ([]string, error)
}
