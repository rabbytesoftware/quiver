package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// InputPrompt handles the command input
type InputPrompt struct {
	value  string
	cursor int
	focus  bool
	width  int
}

// NewInputPrompt creates a new input prompt
func NewInputPrompt() *InputPrompt {
	return &InputPrompt{
		focus: true,
	}
}

// ProcessInput processes the current input value
func (ip *InputPrompt) ProcessInput(input string) error {
	ip.value = input
	ip.cursor = len(input)
	return nil
}

// GetPrompt returns the formatted prompt string
func (ip *InputPrompt) GetPrompt() string {
	prompt := "quiver> "
	
	// Add cursor
	value := ip.value
	if ip.focus {
		if ip.cursor >= len(value) {
			value += "█"
		} else {
			value = value[:ip.cursor] + "█" + value[ip.cursor+1:]
		}
	}

	return prompt + value
}

// SetFocus sets the focus state
func (ip *InputPrompt) SetFocus(focused bool) {
	ip.focus = focused
}

// IsFocused returns whether the prompt is focused
func (ip *InputPrompt) IsFocused() bool {
	return ip.focus
}

// View renders the input prompt
func (ip *InputPrompt) View(width int) string {
	ip.width = width
	promptText := ip.GetPrompt()
	
	return ip.getPromptStyle().
		Width(width - 4).
		Render(promptText)
}

// HandleKey handles key input
func (ip *InputPrompt) HandleKey(key string) {
	switch key {
	case "backspace":
		if ip.cursor > 0 {
			ip.value = ip.value[:ip.cursor-1] + ip.value[ip.cursor:]
			ip.cursor--
		}
	case "left":
		if ip.cursor > 0 {
			ip.cursor--
		}
	case "right":
		if ip.cursor < len(ip.value) {
			ip.cursor++
		}
	default:
		// Add character to input
		if len(key) == 1 {
			ip.value = ip.value[:ip.cursor] + key + ip.value[ip.cursor:]
			ip.cursor++
		}
	}
}

// GetValue returns the current input value
func (ip *InputPrompt) GetValue() string {
	return ip.value
}

// Clear clears the input
func (ip *InputPrompt) Clear() {
	ip.value = ""
	ip.cursor = 0
}

// getPromptStyle returns the styling for the prompt
func (ip *InputPrompt) getPromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#333333")).
		Padding(0, 1).
		Margin(0, 1)
}
