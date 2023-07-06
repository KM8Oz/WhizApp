package main

import (
	"fmt"
	"log"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Account represents an individual account.
type Account struct {
	Name         string
	DatabasePath string
}

var (
	store AccountStore
	_db   *gorm.DB
)

type AccountStore struct {
	Accounts []Account
}

// SaveAccounts saves the accounts to the database.
func SaveAccount(account *Account) {
	var _account Account
	result := _db.First(&_account, "Name = ?", _account.Name)
	if result.RowsAffected > 0 {
		result.Update("Name", account.Name)
		result.Update("DatabasePath", account.DatabasePath)
	} else {
		_db.Create(&account)
	}
	log.Println("Accounts saved successfully.")
}

// LoadAccounts loads the accounts from the database.
func LoadAccounts() ([]Account, error) {
	var accounts []Account
	result := _db.Order("Name DESC").Find(&accounts)
	log.Default().Printf("result:%v", result.RowsAffected)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Println("Accounts loaded successfully.%s", len(accounts))
	return accounts, nil
}

// ShowAccountsWindow displays the main window with the list of accounts.
func ShowAccountsWindow(window fyne.Window) fyne.CanvasObject {
	GUIAPP._win.SetTitle("Whatsapp business manager: accounts")
	db, err := gorm.Open(sqlite.Open("accounts.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	_db = db
	db.AutoMigrate(&Account{})
	accounts, err := LoadAccounts()
	if err != nil {
		fmt.Printf("failed to query accounts: %w", err)
	}
	store.Accounts = accounts
	addButton := widget.NewButtonWithIcon("Add Account", theme.ContentAddIcon(), func() {
		ShowAddAccountWindow(window)
	})
	addButton.Resize(fyne.NewSize(100, 50))
	var vbox fyne.CanvasObject
	// vbox.Resize(fyne.NewSize(600, 500))
	if len(accounts) == 0 {
		vbox = container.NewVBox(
			widget.NewLabel("Account List"),
			widget.NewLabel("No accounts found."),
			addButton,
		)
	} else {
		accountList := widget.NewList(
			func() int {
				return len(store.Accounts)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("")
			},
			func(i widget.ListItemID, obj fyne.CanvasObject) {
				obj.(*widget.Label).SetText(store.Accounts[i].Name)
			},
		)

		accountList.OnSelected = func(id widget.ListItemID) {
			account := &store.Accounts[id]
			log.Println("Selected Account:", account.Name)
		}
		// accountList.Resize(fyne.NewSize(accountList.Size().Width, 300))

		vbox = container.New(layout.NewPaddedLayout(),
			container.NewBorder(widget.NewLabel("Account List"), addButton, nil, nil, container.NewAdaptiveGrid(1,
				accountList,
			)),
		)
	}
	return vbox
}

// ShowAddAccountWindow displays the window to add a new account.
func ShowAddAccountWindow(window fyne.Window) {
	var wg sync.WaitGroup
	accountNameEntry := widget.NewEntry()
	accountNameEntry.SetPlaceHolder("Account Name")

	saveButton := widget.NewButton("Save", func() {
		name := accountNameEntry.Text

		if name == "" {
			log.Println("Account name are required.")
			return
		}
		go tryconnect(&name, &wg)
	})

	cancelButton := widget.NewButton("Cancel", func() {
		showtabs(window)
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Account Name", Widget: accountNameEntry},
		},
		// OnCancel: cancelButton.OnTapped,
		// OnSubmit: saveButton.OnTapped,
	}

	form.Append("Connect", widget.NewFormItem("", saveButton).Widget)
	form.Append("Cancel", widget.NewFormItem("", cancelButton).Widget)

	window.SetContent(form)
	window.Show()
}
