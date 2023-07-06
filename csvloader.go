package main

import (
	"encoding/csv"
	"io"
	"os"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var Selected struct {
	cells []string
}

func NewCell(label fyne.CanvasObject, t widget.TableCellID, table *widget.Table) fyne.CanvasObject {
	c := container.NewVBox(label, widget.NewCheck("nil", func(b bool) {
		table.Select(t)
	}))
	return c
}
func csvtablescreen() {
	var table *widget.Table
	w := GUIAPP._app.NewWindow("CSV Viewer")
	w.Resize(fyne.NewSize(400, 300))

	// Create table widget to display CSV content
	table = widget.NewTable(
		func() (int, int) {
			return 0, 0 // Initial dimensions, will be updated when loading the CSV
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("No data loaded") // Initial placeholder label
		},
		func(t widget.TableCellID, content fyne.CanvasObject) {
			// Implement this function to update the cell content when needed
			NewCell(content, t, table)
		},
	)

	// Create button widget to open file picker
	openButton := widget.NewButtonWithIcon("Open CSV", theme.FolderOpenIcon(), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				// loadCSVFile(reader, table)
			}
		}, w)
	})

	// Create layout to hold buttons
	buttons := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), openButton)

	// Create layout to hold table and buttons
	content := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), table, buttons)

	w.SetContent(content)
	w.ShowAndRun()
}

func loadCSVFile(reader fyne.URIReadCloser, win fyne.Window) {
	defer reader.Close()
	csvFile, err := os.Open(reader.URI().String())
	if err != nil {
		dialog.ShowError(err, win)
		return
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	// Read the CSV file and populate the table
	data := [][]string{}
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			dialog.ShowError(err, win)
			return
		}

		data = append(data, record)
	}
	println(data)
	// Update table dimensions and content
	// table.SetRowCount(len(data))
	// table.SetColumnCount(len(data[0]))
	// table.CreateCell = func() fyne.CanvasObject {
	// 	return widget.NewLabel("")
	// }
	// table.UpdateCell = func(t widget.TableCellID, content fyne.CanvasObject) {
	// 	if t.Row < len(data) && t.Col < len(data[0]) {
	// 		content.(*widget.Label).SetText(data[t.Row][t.Col])
	// 	}
	// }

	// table.Refresh()
}
