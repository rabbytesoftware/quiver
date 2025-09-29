package qtl

import (
	"context"

	fns "github.com/rabbytesoftware/quiver/internal/infrastructure/fetchnshare"
	translator "github.com/rabbytesoftware/quiver/internal/infrastructure/translator/models"
	"github.com/rabbytesoftware/quiver/internal/models/quiver"
)

type QuiverTranslationLayer struct {
	fns fns.FNSInterface
}

func NewQTL(
	fns fns.FNSInterface,
) translator.TranslatorLayerInterface[quiver.Quiver] {
	return &QuiverTranslationLayer{
		fns: fns,
	}
}

func (a *QuiverTranslationLayer) IsCompatible(
	ctx context.Context, 
	manifestPath string,
) (bool, error) {
	return false, nil
}

func (a *QuiverTranslationLayer) Translate(
	ctx context.Context, 
	manifestPath string,
) (*quiver.Quiver, error) {
	return nil, nil
}

func (a *QuiverTranslationLayer) GetManifestVersion(
	ctx context.Context, 
	manifestPath string,
) (string, error) {
	return "", nil
}

func (a *QuiverTranslationLayer) GetSupportedVersions(
	ctx context.Context,
) ([]string, error) {
	return nil, nil
}
