package atl

import (
	"context"

	fns "github.com/rabbytesoftware/quiver/internal/infrastructure/fetchnshare"
	translator "github.com/rabbytesoftware/quiver/internal/infrastructure/translator/models"
	"github.com/rabbytesoftware/quiver/internal/models/arrow"
)

type ArrowTranslationLayer struct {
	fns fns.FNSInterface
}

func NewATL(
	fns fns.FNSInterface,
) translator.TranslatorLayerInterface[arrow.Arrow] {
	return &ArrowTranslationLayer{
		fns: fns,
	}
}

func (a *ArrowTranslationLayer) IsCompatible(
	ctx context.Context, 
	manifestPath string,
) (bool, error) {
	return false, nil
}

func (a *ArrowTranslationLayer) Translate(
	ctx context.Context, 
	manifestPath string,
) (*arrow.Arrow, error) {
	return nil, nil
}

func (a *ArrowTranslationLayer) GetManifestVersion(
	ctx context.Context, 
	manifestPath string,
) (string, error) {
	return "", nil
}

func (a *ArrowTranslationLayer) GetSupportedVersions(
	ctx context.Context,
) ([]string, error) {
	return nil, nil
}
