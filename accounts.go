package main

import (
	"fmt"
	"log"
	"path"
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
	Isselected   bool
	DatabasePath string
}

var (
	store AccountStore
	_db   *gorm.DB
)

type AccountStore struct {
	Accounts []Account
}

func getAccountByName(accounts []Account, name string) *Account {
	for i := 0; i < len(accounts); i++ {
		if accounts[i].Name == name {
			return &accounts[i]
		}
	}
	return nil
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
	home, err := getUserHomeDir()
	if err != nil {
		panic("failed to get home dir")
	}
	_pathdb := path.Join(home, "accounts.db")
	db, err := gorm.Open(sqlite.Open(_pathdb), &gorm.Config{})
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
	var fistload bool = true
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
				toolbar1 := widget.NewToolbar(
					widget.NewToolbarSeparator(),
					widget.NewToolbarAction(theme.DeleteIcon(), func() {
						// _db.Model(&Account{}).Where("Name = ?", account.Name).Update("isselected", true)
					}),
				)
				return container.NewBorder(nil, nil, nil, toolbar1, widget.NewLabel(""))
			},
			nil,
		)
		accountList.UpdateItem = func(i widget.ListItemID, obj fyne.CanvasObject) {
			account := store.Accounts[i]
			if obj.(*fyne.Container).Objects[0].(*widget.Label).Text == "" {
				obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(account.Name)
				if account.Isselected {
					accountList.Select(i)
					fistload = false
				}
			}
		}
		accountList.OnSelected = func(id widget.ListItemID) {
			if !fistload {
				account := &store.Accounts[id]
				_db.Model(&Account{}).Where("Name = ?", account.Name).Update("isselected", true)
				_db.Model(&Account{}).Where("Name != ?", account.Name).Update("isselected", false)
				err = setclient(account)
				println(err)
			}
		}
		var isnothing_selected bool = true
		for _, v := range store.Accounts {
			if v.Isselected {
				isnothing_selected = false
			}
		}
		if isnothing_selected {
			_db.Model(&Account{}).Where("Name = ?", &store.Accounts[0].Name).Update("isselected", true)
			err = setclient(&store.Accounts[0])
			println(err)
		}
		vbox = container.New(layout.NewPaddedLayout(),
			container.NewBorder(widget.NewLabel("Account List"), addButton, nil, nil, container.NewAdaptiveGrid(1,
				accountList,
			)),
		)
	}
	return vbox
}
func setclient(account *Account) error {
	client, err := getclient(account.DatabasePath)
	if err == nil {
		GUIAPP.client = client
		return nil
	}
	return err
}

// ShowAddAccountWindow displays the window to add a new account.
func ShowAddAccountWindow(window fyne.Window) {
	var wg sync.WaitGroup
	accountNameEntry := widget.NewEntry()
	accountNameEntry.SetPlaceHolder("Account Name")
	win := GUIAPP._app.NewWindow("New account form")
	// pattern := "^(?=.*[a-zA-Z0-9])[a-zA-Z0-9]+$"
	// accountNameEntry.Validator = validation.NewRegexp(pattern, "Invalid input [not empty string or number]")
	// accountNameEntry.SetValidationError(fmt.Errorf("Account name are required"))
	saveButton := widget.NewButton("Save", func() {
		name := accountNameEntry.Text
		if name == "" {
			log.Println("Account name are required.")
			accountNameEntry.Validate()
			return
		}
		go tryconnect(&name, &wg)
	})

	cancelButton := widget.NewButton("Cancel", func() {
		showtabs(window)
		win.Close()
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
	win.Resize(fyne.NewSize(400, 150))
	win.SetContent(form)
	win.Show()
}
