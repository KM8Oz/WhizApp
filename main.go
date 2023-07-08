package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"github.com/tmc/langchaingo/schema"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

var GUIAPP = struct {
	_app   fyne.App
	_win   fyne.Window
	client *whatsmeow.Client
}{
	_app:   nil,
	_win:   nil,
	client: nil,
}
var list_contacts map[types.JID]types.ContactInfo = map[types.JID]types.ContactInfo{}
var infoTargetsbng binding.UntypedMap
var active_account Account
var targetedgroups string = "all"
var joinedgroups []*types.GroupInfo = []*types.GroupInfo{}
var isfilter bool = false
var Status binding.Bool = binding.NewBool()
var openaiAPIKey string
var globaldocs map[string][]schema.Document = map[string][]schema.Document{}
var global_context string = "you are a helpful personal assistant"

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Failed to load envirement file:", err)
	// }
	app := app.New()
	GUIAPP._app = app
	app.SetIcon(resourceIcon256Png)
	win := app.NewWindow("Whatsapp business manager")
	win.Resize(fyne.NewSize(800, 500))
	GUIAPP._win = win
	showtabs(win)
	app.Run()
}
