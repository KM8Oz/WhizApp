package main

import (
	"context"
	"fmt"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	openai "github.com/sashabaranov/go-openai"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var str string = ""
var groups []string = []string{"all"}
var Settingsdb *gorm.DB

type Setting struct {
	gorm.Model
	Device string
	Apikey string
}

func createSettingPage() fyne.CanvasObject {
	// startStopToggle := widget.NewAccordion(widget.NewAccordionItem("Start", widget.NewLabel("Bot")), widget.NewAccordionItem("Stop", widget.NewLabel("Bot")))
	home, err := getUserHomeDir()
	if err != nil {
		panic("failed to get home dir")
	}
	_pathdb := path.Join(home, "_settings.db")
	Settingsdb, err = gorm.Open(sqlite.Open(_pathdb), &gorm.Config{})
	if err != nil {
		panic("failed to connect database groups.db")
	}
	Settingsdb.AutoMigrate(&Setting{})
	if GUIAPP.client != nil {
		joinedgroups, err = GUIAPP.client.GetJoinedGroups()
		if err == nil {
			for _, v := range joinedgroups {
				groups = append(groups, v.Name)
			}
		}
	} else {
		setclient(&active_account)
		if GUIAPP.client != nil {
			joinedgroups, err := GUIAPP.client.GetJoinedGroups()
			if err == nil {
				for _, v := range joinedgroups {
					groups = append(groups, v.Name)
				}
			}
		}
	}
	groupDropdown := widget.NewSelect(groups, func(selected string) {
		// Handle group selection
		targetedgroups = selected
	})
	groupDropdown.SetSelected(targetedgroups)
	botstatus := []string{"Start Bot", "Stop Bot"}
	boticons := []fyne.Resource{theme.ConfirmIcon(), theme.CancelIcon()}
	// botstatusicon := []fyne.Resource{ , "Stop Bot"}
	targetsCheckbox := widget.NewCheck("Targets Filter", func(b bool) {
		isfilter = b
	})
	targetsCheckbox.SetChecked(isfilter)
	apiKeyEntry := widget.NewEntry()
	contextinput := widget.NewEntry()
	var setting Setting
	err = Settingsdb.Where(&Setting{Device: active_account.Name}).First(&setting).Error
	apiKeyEntry.SetPlaceHolder("sk-.....")
	contextinput.SetPlaceHolder("Enter the context that you wan't your assistant to answer according to it (400 chars max)")
	if err == nil {
		apiKeyEntry.SetText(setting.Apikey)
		openaiAPIKey = setting.Apikey
	}
	contextinput.OnChanged = func(s string) {
		if len(s) > 400 {
			GUIAPP._app.SendNotification(fyne.NewNotification("Over Limit Context", fmt.Sprintf("over limit context (400 chars max) and you have %d", len(s))))
		}
		global_context = s
	}
	contextinput.SetText(global_context)
	form := widget.NewForm(
		widget.NewFormItem("Activated Group", groupDropdown),
		widget.NewFormItem("Targets Filter", targetsCheckbox),
		widget.NewFormItem("OpenAI API Key", apiKeyEntry),
		widget.NewFormItem("My Assistante context", contextinput),
	)
	str = botstatus[0]
	apiKeyEntry.OnChanged = func(s string) {
		_, _err := openai.NewClient(s).GetModel(context.Background(), openai.GPT3Dot5Turbo)
		if _err == nil {
			err = Settingsdb.Where(&Setting{Device: active_account.Name}).Attrs(&Setting{Apikey: s}).FirstOrCreate(&setting).Error
			if err == nil {
				openaiAPIKey = s
				GUIAPP._app.SendNotification(fyne.NewNotification("Device activated", fmt.Sprintf("Api Key added for device %s!", active_account.Name)))
			}
		}
	}
	var saveButton *widget.Button = nil
	saveButton = widget.NewButtonWithIcon(str, theme.ConfirmIcon(), func() {
		// Handle save button action
		if b, err := Status.Get(); b && err == nil {
			Status.Set(false)
			Listner(false, saveButton)
			if saveButton != nil {
				saveButton.Text = botstatus[0]
				saveButton.Icon = boticons[0]
				saveButton.Refresh()
			}
		} else {
			Status.Set(true)
			Listner(true, saveButton)
			if saveButton != nil {
				saveButton.Text = botstatus[1]
				saveButton.Icon = boticons[1]
				saveButton.Refresh()
			}
		}
	})

	content := container.NewBorder(nil, saveButton, nil, nil,
		form,
	)

	return content
}

func Listner(status bool, saveButton *widget.Button) {
	saveButton.Disable()
	// var i int = 0
	// for item, err := infoTargetsbng.GetValue(infoTargetsbng.Keys()[i]); err == nil; {
	// 	if contact, ok := item.(InfoTarget); ok {
	// 		if ok, err := contact.isactive.Get(); err == nil {
	// 			if ok {
	// 				fmt.Printf("active contact: %v\n", contact.info.FullName)
	// 			}
	// 		}
	// 	}
	// }
	if status {
		if GUIAPP.client != nil {
			gpt := openai.NewClient(openaiAPIKey)
			_, _err := gpt.GetModel(context.Background(), openai.GPT3Dot5Turbo)
			if _err != nil {
				GUIAPP._app.SendNotification(fyne.NewNotification("Can't start", fmt.Sprintf("openai Api Key needed for device %s!", active_account.Name)))
			}
			GUIAPP.client.RemoveEventHandlers()
			GUIAPP.client.AddEventHandler(GetEventHandler(GUIAPP.client, gpt))
			GUIAPP.client.Connect()
		}
	} else {
		GUIAPP.client.RemoveEventHandlers()
		GUIAPP.client.Disconnect()
	}
	saveButton.Enable()
}
