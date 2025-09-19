package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rabbytesoftware/quiver/internal/infrastructure/watcher"
)

// RunUI starts the TUI application with the provided watcher service
func RunUI(w *watcher.Watcher) error {
	// Create the model with the watcher
	model := NewModel(w)
	
	// Create the program with alt screen
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)
	
	// Run the program
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}
	
	// Clean up
	if m, ok := finalModel.(*Model); ok && m.cancel != nil {
		m.cancel()
	}
	
	return nil
}
