package ui

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/pterm/pterm"
	"github.com/rabbytesoftware/quiver/internal/metadata"
)

// getTerminalWidth returns the terminal width, with a fallback
func getTerminalWidth() int {
	// Try to get terminal width from pterm
	width := pterm.GetTerminalWidth()
	if width <= 0 {
		width = 80 // Fallback to common terminal width
	}
	return width
}

// wrapText wraps text to fit within the specified width
func wrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}
	
	var result []string
	words := strings.Fields(text)
	var line string
	
	for _, word := range words {
		// If adding this word would exceed the width, start a new line
		if len(line)+len(word)+1 > width {
			if line != "" {
				result = append(result, line)
				line = word
			} else {
				// Single word is longer than width, just add it
				result = append(result, word)
			}
		} else {
			if line == "" {
				line = word
			} else {
				line += " " + word
			}
		}
	}
	
	if line != "" {
		result = append(result, line)
	}
	
	return result
}

// prepareInfoLines prepares and wraps all info lines to fit the available width
func prepareInfoLines(username, hostname, version, description, author, license string, maintainers []metadata.Maintainer, maxWidth int) []string {
	var infoLines []string
	
	// Add basic info
	infoLines = append(infoLines, fmt.Sprintf("%s@%s", username, hostname))
	infoLines = append(infoLines, "--------------------------")	
	
	// Wrap long lines
	softwareLines := wrapText(fmt.Sprintf("Quiver v%s", version), maxWidth)
	infoLines = append(infoLines, softwareLines...)
	
	descriptionLines := wrapText(fmt.Sprintf("%s", description), maxWidth)
	for i, line := range descriptionLines {
		if i == 0 {
			infoLines = append(infoLines, line)
		} else {
			infoLines = append(infoLines, "             "+line)
		}
	}

	infoLines = append(infoLines, fmt.Sprintf("OS: %s %s", runtime.GOOS, runtime.GOARCH))
	infoLines = append(infoLines, fmt.Sprintf("Kernel: %s", runtime.Version()))
	
	authorLines := wrapText(fmt.Sprintf("Copyright by: %s", author), maxWidth)
	infoLines = append(infoLines, authorLines...)
	
	licenseLines := wrapText(fmt.Sprintf("Published under: %s", license), maxWidth)
	infoLines = append(infoLines, licenseLines...)
	
	// Add maintainers
	if len(maintainers) > 0 {
		infoLines = append(infoLines, "Maintainers:")
		for _, maintainer := range maintainers {
			maintainerStr := fmt.Sprintf("  • %s", maintainer.Name)
			if maintainer.Email != "" {
				maintainerStr += fmt.Sprintf(" <%s>", maintainer.Email)
			}
			if maintainer.URL != "" {
				maintainerStr += fmt.Sprintf(" - %s", maintainer.URL)
			}
			
			// Wrap maintainer lines
			maintainerLines := wrapText(maintainerStr, maxWidth)
			for i, line := range maintainerLines {
				if i == 0 {
					infoLines = append(infoLines, line)
				} else {
					infoLines = append(infoLines, "    "+line) // Indent continuation for maintainer
				}
			}
		}
	}
	
	return infoLines
}

// ShowWelcome displays the welcome message in neofetch style
func ShowWelcome() {
	// Clear screen first
	ClearScreen()
	
	// Quiver ASCII Art (provided by user)
	quiverArt := []string{
		"         ;;;;;;;;;;;;;;;;;",
		"         ;;;;;;;;;;;;;;;;;",
		"         ;;;;;;;;;;;;;;;;;", 
		"         ;;;;;;;;;;;;;;;;;",
		"                          ",
		";;;;;;;;          ;;;;;;;;",
		";;;;;;;;          ;;;;;;;;",
		";;;;;;;;          ;;;;;;;;",
		";;;;;;;;          ;;;;;;;;",
		";;;;;;;;          ;;;;;;;;",
		";;;;;;;;          ;;;;;;;;",
		";;;;;;;;          ;;;;;;;;",
		";;;;;;;;          ;;;;;;;;",
		";;;;;;;;          ;;;;;;;;",
		";;;;;;;;;;;;;;;;;         ",
		";;;;;;;;;;;;;;;;;         ",
		";;;;;;;;;;;;;;;;;         ",
		";;;;;;;;;;;;;;;;;         ",
		"                          ",
		"         ;;;;;;;;;;;;;;;;;",
		"         ;;;;;;;;;;;;;;;;;",
		"         ;;;;;;;;;;;;;;;;;",
		"         ;;;;;;;;;;;;;;;;;",
	}

	// Get system information
	hostname, _ := os.Hostname()
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME") // Windows fallback
	}
	
	// Get metadata
	version := metadata.GetVersion()
	description := metadata.GetDescription()
	author := metadata.GetAuthor()
	license := metadata.GetLicense()
	maintainers := metadata.GetMaintainers()
	
	// Calculate available width for text dynamically
	const asciiWidth = 57
	const spacing = 4 // A bit more spacing for safety
	terminalWidth := getTerminalWidth()
	maxTextWidth := terminalWidth - asciiWidth - spacing
	
	// Ensure reasonable bounds
	if maxTextWidth < 30 {
		maxTextWidth = 30 // Minimum reasonable width
	} else if maxTextWidth > 60 {
		maxTextWidth = 60 // Maximum to maintain readability
	}
	
	// Prepare info lines with wrapping
	infoLines := prepareInfoLines(username, hostname, version, description, author, license, maintainers, maxTextWidth)
	
	// Style configuration
	artStyle := pterm.NewStyle(pterm.FgCyan, pterm.Bold)
	infoStyle := pterm.NewStyle(pterm.FgWhite, pterm.Bold)
	labelStyle := pterm.NewStyle(pterm.FgCyan, pterm.Bold)
	
	// Display ASCII art and info side by side
	maxLines := len(quiverArt)
	if len(infoLines) > maxLines {
		maxLines = len(infoLines)
	}
	
	fmt.Println()
	for i := 0; i < maxLines; i++ {
		// Print ASCII art line
		artLine := ""
		if i < len(quiverArt) {
			artLine = quiverArt[i]
		} else {
			artLine = strings.Repeat(" ", asciiWidth)
		}
		artStyle.Print(artLine)
		
		// Print info line
		if i < len(infoLines) {
			infoLine := infoLines[i]
			if i == 1 { // Username@hostname line - make it bold
				labelStyle.Print("  ")
				labelStyle.Println(infoLine)
			} else if strings.Contains(infoLine, ":") && !strings.HasPrefix(infoLine, "─") && !strings.HasPrefix(infoLine, " ") {
				// Lines with labels - colorize the label part (but not indented continuation lines)
				parts := strings.SplitN(infoLine, ":", 2)
				if len(parts) == 2 {
					labelStyle.Print("  ")
					labelStyle.Print(parts[0] + ":")
					infoStyle.Println(parts[1])
				} else {
					infoStyle.Print("  ")
					infoStyle.Println(infoLine)
				}
			} else {
				// Regular lines (including continuation lines)
				infoStyle.Print("  ")
				infoStyle.Println(infoLine)
			}
		} else {
			fmt.Println()
		}
	}
	
	fmt.Println()
	
	// Print status
	pterm.DefaultSection.WithLevel(2).Println("Initializing Quiver...")
}

// ShowTable displays a formatted table
func ShowTable(title string, headers []string, data [][]string) {
	// Create table data with headers
	tableData := make([][]string, 0, len(data)+1)
	tableData = append(tableData, headers)
	tableData = append(tableData, data...)

	// Render table
	pterm.DefaultSection.Println(title)
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

// ShowPackageStatus displays package loading status
func ShowPackageStatus(packages map[string]interface{}) {
	if len(packages) == 0 {
		pterm.Warning.Println("No packages loaded")
		return
	}

	pterm.DefaultSection.Println("Loaded Packages")
	
	for name, pkg := range packages {
		pterm.Success.Printf("✓ %s\n", name)
		// Add more package details here if needed
		_ = pkg // Avoid unused variable warning
	}
}

// ShowProgress displays a progress indicator
func ShowProgress(message string) *pterm.SpinnerPrinter {
	spinner, _ := pterm.DefaultSpinner.Start(message)
	return spinner
}

// ShowSuccess displays a success message
func ShowSuccess(message string) {
	pterm.Success.Println(message)
}

// ShowError displays an error message
func ShowError(message string) {
	pterm.Error.Println(message)
}

// ShowWarning displays a warning message
func ShowWarning(message string) {
	pterm.Warning.Println(message)
}

// ShowInfo displays an info message
func ShowInfo(message string) {
	pterm.Info.Println(message)
}

// Confirm prompts user for confirmation
func Confirm(message string) bool {
	result, _ := pterm.DefaultInteractiveConfirm.WithDefaultValue(false).WithDefaultText(message).Show()
	return result
}

// GetInput prompts user for input
func GetInput(prompt string) string {
	result, _ := pterm.DefaultInteractiveTextInput.WithDefaultText(prompt).Show()
	return result
}

// ShowServerInfo displays server information
func ShowServerInfo(host string, port int) {
	info := fmt.Sprintf(`
Server Configuration:
• Host: %s
• Port: %d
• API Endpoint: http://%s:%d/api/v1
• Health Check: http://%s:%d/health
`, host, port, host, port, host, port)

	pterm.DefaultBox.WithTitle("Server Info").WithTitleTopCenter().Println(info)
}

// ShowShutdown displays shutdown message
func ShowShutdown() {
	pterm.DefaultCenter.Println(
		pterm.DefaultHeader.WithFullWidth().
			WithBackgroundStyle(pterm.NewStyle(pterm.BgRed)).
			WithTextStyle(pterm.NewStyle(pterm.FgWhite)).
			Sprint("Quiver Server Shutdown"),
	)
}

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// ExitWithError exits with error message
func ExitWithError(message string, code int) {
	ShowError(message)
	os.Exit(code)
} 