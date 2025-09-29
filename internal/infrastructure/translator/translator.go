package translator

import (
	fns "github.com/rabbytesoftware/quiver/internal/infrastructure/fetchnshare"
	translator "github.com/rabbytesoftware/quiver/internal/infrastructure/translator/models"
	"github.com/rabbytesoftware/quiver/internal/models/arrow"
	"github.com/rabbytesoftware/quiver/internal/models/quiver"

	atl "github.com/rabbytesoftware/quiver/internal/infrastructure/translator/atl"
	qtl "github.com/rabbytesoftware/quiver/internal/infrastructure/translator/qtl"
)

// TranslatorImplementation acts as DI container for the Translator module

type TranslatorImplementation struct {
	fns    fns.FNSInterface

	atl 	 translator.TranslatorLayerInterface[arrow.Arrow]
	qtl 	 translator.TranslatorLayerInterface[quiver.Quiver]
}

func NewTranslator(fns fns.FNSInterface) translator.TranslatorInterface {
	return TranslatorImplementation{
		fns:    fns,
		atl:    atl.NewATL(fns),
		qtl:    qtl.NewQTL(fns),
	}
}

// GetArrowTranslator returns the Arrow translation layer
func (t TranslatorImplementation) GetArrowTranslator() translator.TranslatorLayerInterface[arrow.Arrow] {
	return t.atl
}

// GetQuiverTranslator returns the Quiver translation layer
func (t TranslatorImplementation) GetQuiverTranslator() translator.TranslatorLayerInterface[quiver.Quiver] {
	return t.qtl
}
