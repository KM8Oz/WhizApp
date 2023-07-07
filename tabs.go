package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func showtabs(win fyne.Window) {
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Accounts", resourceIcons8Whatsapp128Png, ShowAccountsWindow(win)),
		container.NewTabItemWithIcon("Targets", resourceIcons8Target72XxhdpiPng, groupstable()),
		container.NewTabItemWithIcon("Logs", resourceIcons8Settings128Png, widget.NewLabel("Settings tab")),
		container.NewTabItemWithIcon("About", resourceIcons8About150Png, widget.NewLabel("About tab")),
	)

	tabs.SetTabLocation(container.TabLocationLeading)
	win.SetContent(tabs)
	win.Show()
}
