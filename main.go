package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/joho/godotenv"
)

var GUIAPP = struct {
	_app fyne.App
	_win fyne.Window
}{
	_app: nil,
	_win: nil,
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load envirement file:", err)
	}
	app := app.New()
	GUIAPP._app = app
	app.SetIcon(resourceIcon256Png)
	win := app.NewWindow("Whatsapp business manager")
	win.Resize(fyne.NewSize(800, 500))
	GUIAPP._win = win
	showtabs(win)
	app.Run()
}
