package ui

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/rabbytesoftware/quiver/internal/metadata"
)

// ShowWelcome displays the welcome message
func ShowWelcome() {
	// Create a big text for the title
	bigText, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("QUIVER", pterm.NewStyle(pterm.FgCyan)),
	).Srender()

	// Print the big text
	pterm.Println(bigText)

	// Print description
	description := metadata.GetDescription()
	version := metadata.GetVersion()
	author := metadata.GetAuthor()
	license := metadata.GetLicense()
	maintainers := metadata.GetMaintainers()
	
	pterm.DefaultBox.WithTitle("About").WithTitleTopCenter().Println(description)
	pterm.DefaultBox.WithTitle("Version").WithTitleTopCenter().Println(version + " by " + author + " - " + license)

	for _, maintainer := range maintainers {
		pterm.DefaultBox.WithTitle("Maintainer").WithTitleTopCenter().Println(maintainer.Name + " - " + maintainer.Email + " - " + maintainer.URL)
	}

	// Print status
	pterm.DefaultSection.Println("Initializing Quiver...")
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