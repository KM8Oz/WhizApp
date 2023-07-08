package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func showtabs(win fyne.Window) {
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Accounts", resourceIcons8Whatsapp128Png, ShowAccountsWindow(win)),
		container.NewTabItemWithIcon("Targets", resourceIcons8Target72XxhdpiPng, groupstable()),
		container.NewTabItemWithIcon("Settings", resourceIcons8Settings128Png, createSettingPage()),
		// container.NewTabItemWithIcon("Logs", resourceIcons8About150Png, widget.NewLabel("Logs tab")),
	)

	tabs.SetTabLocation(container.TabLocationLeading)
	win.SetContent(tabs)
	win.Show()
}
