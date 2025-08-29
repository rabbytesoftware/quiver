package ui

import (
	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/tui"
)

type UI struct {
	
}

func NewUI() *UI {
	return &UI{}
}

func (ui *UI) Run() error {
	return tui.RunTUI()
}