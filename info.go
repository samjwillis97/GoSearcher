package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	widget "fyne.io/fyne/v2/widget"
	"strconv"
)

func createInfoWindow(data map[string]string, S Service) fyne.Window {

	var windowName string
	var primaryFields = S.GetPrimaryFields()
	if len(primaryFields) > 0 {
		if len(primaryFields) > 1 {
			for _, val := range primaryFields {
				windowName = windowName + data[val] + " "
			}
		} else {
			windowName = data[primaryFields[0]]
		}
	}

	newWindow := a.NewWindow(windowName)
	desiredWidth := float32(500)

	form := widget.NewForm()
	shortCuts := map[fyne.Shortcut]func(shortcut fyne.Shortcut){}

	// TODO: do this better
	copyIcon, err := fyne.LoadResourceFromPath("./content_copy.svg")
	if err != nil {
		panic(err)
	}

	i := 0
	//for key, val := range dataMap {
	for _, val := range S.GetDisplayFields() {
		// TODO: Add keybinding (Ctrl + 1) copies first field
		i++

		newWidget := widget.NewEntry()
		newWidget.Text = data[val]

		copyCallback := func() {
			newWindow.Clipboard().SetContent(data[val])
			a.SendNotification(
				fyne.NewNotification(
					"Content Copied",
					fmt.Sprintf("%s copied to clipboard.", val),
				),
			)
			newWindow.Close()
		}

		buttonLabel := ""

		if i < 10 {
			keyString := strconv.Itoa(i)
			copyKey := &desktop.CustomShortcut{KeyName: fyne.KeyName(keyString), Modifier: fyne.KeyModifierSuper}
			shortCuts[copyKey] = func(shortcut fyne.Shortcut) {
				copyCallback()
			}
			buttonLabel = " - " + keyString
		}

		widgetCopy := widget.NewButtonWithIcon(buttonLabel, copyIcon, copyCallback)

		widgetContent := container.New(layout.NewBorderLayout(nil, nil, nil, widgetCopy), newWidget, widgetCopy)
		newItem := widget.FormItem{
			Text:   val,
			Widget: widgetContent,
		}
		form.AppendItem(&newItem)
	}

	newWindow.SetContent(form)
	for key, val := range shortCuts {
		newWindow.Canvas().AddShortcut(key, val)
	}
	newWindow.Resize(fyne.Size{
		Width:  desiredWidth,
		Height: 0,
	})
	newWindow.CenterOnScreen()

	return newWindow
}
