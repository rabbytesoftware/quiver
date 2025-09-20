package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rabbytesoftware/quiver/internal/core/watcher"
)

// RunUI starts the TUI application with the provided watcher service
func RunUI(w *watcher.Watcher) error {
	p := tea.NewProgram(
		NewModel(w),
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)
	
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}
	
	if m, ok := finalModel.(*Model); ok && m.cancel != nil {
		m.cancel()
	}
	
	return nil
}
