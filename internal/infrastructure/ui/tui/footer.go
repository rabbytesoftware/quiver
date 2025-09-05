package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ResourceMetrics holds system resource information
type ResourceMetrics struct {
	CPUPercent     float64
	MemoryPercent  float64
	MemoryUsed     uint64
	MemoryTotal    uint64
	LastUpdate     time.Time
}

// FooterBar handles the resource metrics display
type FooterBar struct {
	metrics ResourceMetrics
	width   int
}

// NewFooterBar creates a new footer bar
func NewFooterBar() *FooterBar {
	return &FooterBar{}
}

// UpdateMetrics updates the metrics displayed in the footer
func (fb *FooterBar) UpdateMetrics(metrics ResourceMetrics) {
	fb.metrics = metrics
}

// View renders the footer
func (fb *FooterBar) View(width int) string {
	fb.width = width
	
	cpu := fb.getCpuStyle().Render(fmt.Sprintf("CPU: %.1f%%", fb.metrics.CPUPercent))
	memory := fb.getMemoryStyle().Render(fmt.Sprintf("RAM: %.1f%% (%.1fGB/%.1fGB)", 
		fb.metrics.MemoryPercent, 
		float64(fb.metrics.MemoryUsed)/1024/1024/1024,
		float64(fb.metrics.MemoryTotal)/1024/1024/1024))
	
	footerContent := fmt.Sprintf("%s  %s", cpu, memory)
	
	return fb.getFooterStyle().
		Width(width).
		Render(footerContent)
}

// getFooterStyle returns the main footer styling
func (fb *FooterBar) getFooterStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#0F4C75")).
		Padding(0, 1).
		Bold(true)
}

// getCpuStyle returns the CPU metrics styling
func (fb *FooterBar) getCpuStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Bold(true)
}

// getMemoryStyle returns the memory metrics styling
func (fb *FooterBar) getMemoryStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4ECDC4")).
		Bold(true)
}

// getNetworkStyle returns the network metrics styling
func (fb *FooterBar) getNetworkStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#45B7D1")).
		Bold(true)
}
