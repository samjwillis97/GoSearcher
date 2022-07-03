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
	if len(S.PrimaryField) > 0 {
		if len(S.PrimaryField) > 1 {
			for _, val := range S.PrimaryField {
				windowName = windowName + data[val] + " "
			}
		} else {
			windowName = data[S.PrimaryField[0]]
		}
	} else if len(S.DisplayFields) > 0 {
		windowName = data[S.DisplayFields[0]]
	} else {
		windowName = data[S.SearchFields[0]]
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
	for _, key := range S.DisplayFields {
		// TODO: Add keybinding (Ctrl + 1) copies first field
		i++

		newWidget := widget.NewEntry()
		newWidget.Text = data[key]

		fieldDisplayName := key
		if val, ok := S.FieldSettings[key]; ok {
			fieldDisplayName = val.Display
		}

		copyCallback := func() {
			newWindow.Clipboard().SetContent(data[key])
			a.SendNotification(
				fyne.NewNotification(
					"Content Copied",
					fmt.Sprintf("%s copied to clipboard.", fieldDisplayName),
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
			Text:   fieldDisplayName,
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
