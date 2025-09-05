package ui

import (
	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/tui"
)

type TUI struct {
	model *tui.Model
}

func NewTUI() *TUI {
	return &TUI{}
}

func (t *TUI) Run() error {
	return tui.RunTUI()
}

func (t *TUI) Stop() error {
	return nil
}
