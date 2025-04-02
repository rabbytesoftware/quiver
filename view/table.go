package view

import "github.com/pterm/pterm"

func Table(
	Title 	string,
	Header 	[]string,
	Data 	[][]string,
) {
	// Define the data for the table.
	// Each inner slice represents a row in the table.
	// The first row is considered as the header.
	tableData := pterm.TableData{
		Header,
	}

	for _, row := range Data {
		tableData = append(tableData, row)
	}

	// Create a table with the defined data.
	// The table has a header and the text in the cells is right-aligned.
	// The Render() method is used to print the table to the console.
	pterm.DefaultTable.WithHasHeader().WithLeftAlignment().WithData(tableData).Render()
}