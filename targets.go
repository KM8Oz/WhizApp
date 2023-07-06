package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	GroupID      uint
	Name         string
	Fullname     string
	Businessname string
	Jid          string
}

type Group struct {
	gorm.Model
	Name    string
	Members []Member
}

var (
	groupsdb *gorm.DB
)

//lint: Ignore func table is unused (U1000)go-staticcheck
func groupstable() fyne.CanvasObject {
	var err error
	groupsdb, err = gorm.Open(sqlite.Open("groups.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database groups.db")
	}
	groupsdb.AutoMigrate(&Group{})
	GUIAPP._win.SetTitle("Whatsapp business manager: targets")
	// rectangle := canvas.NewRectangle(color.RGBA{152, 160, 197, 1})
	// text1 := container.NewPadded(rectangle, container.NewCenter(canvas.NewText("1", color.Black)))
	// text2 := container.NewPadded(rectangle, container.NewCenter(canvas.NewText("2", color.Black)))
	// text3 := container.NewPadded(rectangle, container.NewCenter(canvas.NewText("3", color.Black)))

	// grid := container.NewBorder(widget.NewLabel("Targets List"), nil, nil, nil, container.New(layout.NewGridWrapLayout(fyne.NewSize(100, 100)),
	// 	text1, text2, text3))
	return (createRectGrid())
}
func createRectGrid() fyne.CanvasObject {
	grid := container.NewGridWrap(fyne.Size{Width: 100, Height: 100})
	scroll := container.NewVScroll(container.NewPadded(grid))
	space := layout.NewSpacer()
	screen := container.NewBorder(widget.NewLabel("Targets List"), createBottomButtons(), space, space, scroll)
	// Create rectangles with click animations
	for i := 0; i < 10; i++ {
		rect := widget.NewButton("", func() {})
		rect.Importance = widget.MediumImportance
		rect.ExtendBaseWidget(rect)
		grid.Add(container.NewPadded(rect))
	}
	return screen
}
func createBottomButtons() fyne.CanvasObject {
	importCSVButton := widget.NewButton("Import from CSV", func() {
		// TODO: Implement CSV import functionality
		addgroup()
	})
	importWhatsAppButton := widget.NewButton("Import from WhatsApp", func() {
		// TODO: Implement WhatsApp import functionality
	})

	buttons := container.NewBorder(nil, nil, nil, container.NewHBox(importCSVButton, importWhatsAppButton))
	return buttons
}
func addgroup() {
	newapp := GUIAPP._app.NewWindow("import")
	newapp.Resize(newapp.Content().MinSize())
	newapp.SetIcon(theme.FileApplicationIcon())
	popup := container.NewMax(container.NewCenter(widget.NewCard(
		"Import from file CSV",
		"click import to finish",
		widget.NewButtonWithIcon("Confirm", theme.FileIcon(), func() {
			newapp.Resize(fyne.NewSize(400, 300))
			dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err == nil && reader != nil {
					// loadCSVFile(reader, table)
					loadCSVFile(reader, newapp)
					newapp.Hide()
				}
			}, newapp)

		}),
	)))
	newapp.SetContent(popup)
	newapp.Show()
}
