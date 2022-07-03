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

func createInfoWindow(data clientData) fyne.Window {
	newWindow := a.NewWindow(data.BusinessName)
	desiredWidth := float32(500)

	dataMap, dataFields, _ := ToMap(data, "struct")
	form := widget.NewForm()
	shortCuts := map[fyne.Shortcut]func(shortcut fyne.Shortcut){}

	// TODO: do this better
	copyIcon, err := fyne.LoadResourceFromPath("./content_copy.svg")
	if err != nil {
		panic(err)
	}

	i := 0
	//for key, val := range dataMap {
	for _, key := range dataFields {
		stringVal := ""
		switch val := dataMap[key].(type) {
		case string:
			stringVal = val
		case int:
			stringVal = strconv.Itoa(val)
		}
		if stringVal != "" {
			// TODO: Add keybinding (Ctrl + 1) copies first field
			i++

			fieldName := key
			newWidget := widget.NewEntry()
			newWidget.Text = stringVal

			copyCallback := func() {
				newWindow.Clipboard().SetContent(stringVal)
				a.SendNotification(
					fyne.NewNotification(
						"Content Copied",
						fmt.Sprintf("%s copied to clipboard.", snakeToDisplay(fieldName)),
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
				Text:   snakeToDisplay(key),
				Widget: widgetContent,
			}
			form.AppendItem(&newItem)
		}
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
