package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"path"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	Name         string
	Fullname     string
	Businessname string
	Jid          string
	IsActive     bool
	AccountID    uint
}

var WhatsAppGreen = color.RGBA{
	R: 37,
	G: 211,
	B: 102,
	A: 255,
}

type Targets struct {
	gorm.Model
	AccountID *uint `gorm:"primaryKey;autoIncrement:true"`
	Account   string
	Members   []Member `gorm:"foreignkey:AccountID"`
}

var (
	Targetsdb *gorm.DB
)

func getclient(path string) (*whatsmeow.Client, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, err
	}
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:"+path+"?_foreign_keys=on", dbLog)
	if err != nil {
		return nil, err
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, err
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.EnableAutoReconnect = true
	return client, nil
}

//lint: Ignore func table is unused (U1000)go-staticcheck
func groupstable() fyne.CanvasObject {
	var err error
	home, err := getUserHomeDir()
	if err != nil {
		panic("failed to get home dir")
	}
	_pathdb := path.Join(home, "target.db")
	Targetsdb, err = gorm.Open(sqlite.Open(_pathdb), &gorm.Config{})
	if err != nil {
		panic("failed to connect database groups.db")
	}
	Targetsdb.AutoMigrate(&Targets{}, &Member{})
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
	var wg sync.WaitGroup
	grid := container.NewGridWrap(fyne.Size{Width: 120, Height: 100})
	scroll := container.NewVScroll(container.NewPadded(grid))
	space := layout.NewSpacer()
	search := binding.NewString()
	// Create rectangles with click animations
	accounts, err := LoadAccounts()
	if err != nil {
		fmt.Printf("failed to query accounts: %w", err)
	}
	for _, a := range accounts {
		if a.Isselected {
			active_account = a
		}
	}
	screen := container.NewBorder(widget.NewLabel("Targets List"), createBottomButtons(grid, active_account, search), space, space, scroll)
	if active_account.Name != "" {
		if GUIAPP.client == nil {
			_ = setclient(&active_account)
			fillgrid(grid, active_account, &wg)
		}
		fillgrid(grid, active_account, &wg)
	}

	return screen
}

type InfoTarget struct {
	jid      types.JID
	info     types.ContactInfo
	isactive binding.Bool
}

func fillgrid(grid *fyne.Container, account Account, wg *sync.WaitGroup) {
	if GUIAPP.client != nil {
		var err error
		if len(list_contacts) == 0 || account.Name != active_account.Name {
			err = GUIAPP.client.Connect()
		} else {
			err = nil
		}
		if err == nil {
			if len(list_contacts) == 0 || account.Name != active_account.Name {
				list_contacts, err = GUIAPP.client.Store.Contacts.GetAllContacts()
			}
			if err == nil {
				infoTargets := make(map[string]interface{}, len(list_contacts))
				for jid, v := range list_contacts {
					infoTargets[jid.User] = InfoTarget{
						info:     v,
						isactive: binding.NewBool(),
					}
				}
				infoTargetsbng = binding.NewUntypedMap()
				infoTargetsbng.Set(infoTargets)
				for jid, v := range list_contacts {
					var _profileimage *canvas.Image = nil
					// profileimage, _err := GUIAPP.client.GetProfilePictureInfo(jid, &whatsmeow.GetProfilePictureParams{
					// 	IsCommunity: true,
					// })
					// fmt.Printf("profileimage.URL:%s \n", profileimage.DirectPath)
					_profileimage = canvas.NewImageFromResource(resourceIcons8Target96XxxhdpiPng)
					_profileimage.FillMode = canvas.ImageFillOriginal
					// if _err != nil {

					// } else {
					// 	_bytes, err := GUIAPP.client.DownloadAny(&proto.Message{
					// 		ImageMessage: &proto.ImageMessage{
					// 			Url:        &profileimage.URL,
					// 			Mimetype:   &profileimage.Type,
					// 			DirectPath: &profileimage.DirectPath,
					// 		},
					// 	})
					// 	if err != nil {
					// 		_profileimage := canvas.NewImageFromReader(bytes.NewReader(_bytes), "qrcode")
					// 		_profileimage.FillMode = canvas.ImageFillOriginal
					// 	} else {
					// 		_profileimage = canvas.NewImageFromResource(resourceIcons8Target96XxxhdpiPng)
					// 		_profileimage.FillMode = canvas.ImageFillOriginal
					// 	}
					// }
					title := widget.NewLabel(v.FullName + " " + v.FirstName + " " + v.BusinessName + " " + v.PushName)
					title.TextStyle = fyne.TextStyle{
						Monospace: true,
					}
					title.Alignment = fyne.TextAlignCenter
					title.Wrapping = fyne.TextTruncate
					title.ExtendBaseWidget(&widget.Button{})
					// item.isactive.AddListener(binding.NewDataListener(func() {
					// 	if b, _ := item.isactive.Get(); b {
					// 		item.isactive.Set(false)
					// 	} else {
					// 		item.isactive.Set(true)
					// 	}
					// }))
					item, err := infoTargetsbng.GetValue(jid.User)
					boolean := item.(InfoTarget).isactive
					boolean.AddListener(binding.NewDataListener(func() {
						myitem := item.(InfoTarget)
						if v, err := boolean.Get(); v && err == nil {
							var COUNT int64
							var _targets Targets
							res := Targetsdb.Model(&Targets{}).Assign(&Targets{
								Account: active_account.Name,
							}).FirstOrCreate(_targets, "Account = ?", active_account.Name)
							res.Count(&COUNT)
							if res.Error == nil && COUNT > 0 {
								eer := res.Association("Members").Append(&Member{
									Name:         myitem.info.FirstName,
									Fullname:     myitem.info.FullName,
									Businessname: myitem.info.BusinessName,
									Jid:          myitem.jid.User,
								})
								if eer != nil {
									fmt.Println(eer)
								}
							} else {
								Targetsdb.Create(&Targets{
									Account: active_account.Name,
									Members: []Member{
										{
											Name:         myitem.info.FirstName,
											Fullname:     myitem.info.FullName,
											Businessname: myitem.info.BusinessName,
											Jid:          myitem.jid.User,
										},
									},
								})
							}
						} else {
							var _targets Targets
							res := Targetsdb.Model(&Targets{}).First(_targets, "Account = ?", active_account.Name)
							if res.Error == nil {
								eer := res.Association("Members").Delete(&Member{
									Jid: myitem.jid.User,
								})
								if eer != nil {
									fmt.Println(eer)
								}
							}
						}
					}))
					if err != nil {
						fmt.Printf("item: %v", err)
					}
					card := container.NewPadded(
						title,
						widget.NewCheckWithData(jid.User, boolean),
					)
					grid.Add(card)
				}
			}
		}
	}
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// <-c
}
func createBottomButtons(grid *fyne.Container, active_account Account, str binding.String) fyne.CanvasObject {
	var wg sync.WaitGroup
	importWhatsAppButton := widget.NewButtonWithIcon("Reload contacts", theme.ViewRefreshIcon(), func() {
		// TODO: Implement WhatsApp import functionality
		if GUIAPP.client == nil {
			_ = setclient(&active_account)
			fillgrid(grid, active_account, &wg)
		}
		grid.Refresh()
		fmt.Printf("GUIAPP.client: %v, grid count: %v\n", GUIAPP.client != nil, len(grid.Objects))
	})
	importWhatsAppButton.ExtendBaseWidget(&widget.Icon{})
	intry := widget.NewEntryWithData(str)
	intry.Resize(intry.MinSize().AddWidthHeight(200, 0))
	intry.SetPlaceHolder("phone ..")
	content := container.NewHBox(intry, widget.NewButton("Search", func() {
		log.Println("Content was:", intry.Text)
	}))
	content.Resize(content.MinSize().AddWidthHeight(100, 0))
	buttons := container.NewBorder(nil, nil, content, container.NewHBox(importWhatsAppButton), nil)
	return buttons
}

// func addgroup() {
// 	newapp := GUIAPP._app.NewWindow("import")
// 	newapp.Resize(newapp.Content().MinSize())
// 	newapp.SetIcon(theme.FileApplicationIcon())
// 	popup := container.NewMax(container.NewCenter(widget.NewCard(
// 		"Import from file CSV",
// 		"click import to finish",
// 		widget.NewButtonWithIcon("Confirm", theme.FileIcon(), func() {
// 			newapp.Resize(fyne.NewSize(400, 300))
// 			dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
// 				if err == nil && reader != nil {
// 					// loadCSVFile(reader, table)
// 					loadCSVFile(reader, newapp)
// 					newapp.Hide()
// 				}
// 			}, newapp)

// 		}),
// 	)))
// 	newapp.SetContent(popup)
// 	newapp.Show()
// }
