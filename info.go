package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	widget "fyne.io/fyne/v2/widget"
	"strconv"
	"strings"
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

	//copyIcon, err := fyne.LoadResourceFromPath("./icons/content_copy.svg")
	//if err != nil {
	//	panic(err)
	//}
	//
	//qrIcon, err := fyne.LoadResourceFromPath("./icons/qr_icon.svg")
	//if err != nil {
	//	panic(err)
	//}

	// Create Shortcuts for Copying
	var copyButtonLabels []string
	assignedBindings := make(map[string]struct{})
	modifier := S.GetSearchModifierKey()
	i := 0
	for _, val := range S.GetDisplayFields() {
		i++
		if i < 10 {
			boundKey := val.GetKeyBinding()
			if boundKey == "" {
				boundKey = strconv.Itoa(i)
			}
			if _, ok := assignedBindings[boundKey]; ok {
				boundKey = strconv.Itoa(i)
			}
			assignedBindings[boundKey] = struct{}{}

			keyBinding := &desktop.CustomShortcut{
				KeyName:  fyne.KeyName(boundKey),
				Modifier: modifier,
			}

			callback := createCopyCallback(
				newWindow,
				val.GetDisplayName(),
				data[val.Name],
			)
			shortCuts[keyBinding] = func(shortcut fyne.Shortcut) {
				callback()
			}
			copyButtonLabels = append(copyButtonLabels, " - "+strings.ToLower(boundKey))
		}
	}

	i = 0
	for _, val := range S.GetDisplayFields() {
		i++

		newWidget := widget.NewEntry()
		newWidget.Text = data[val.Name]

		buttonLabel := ""
		if i < 10 {
			buttonLabel = copyButtonLabels[i-1]
		}

		widgetCopy := widget.NewButton(
			buttonLabel,
			createCopyCallback(
				newWindow,
				val.GetDisplayName(),
				data[val.Name],
			),
		)
		//widgetCopy := widget.NewButtonWithIcon(buttonLabel, copyIcon, copyCallback)

		var qrLayout fyne.CanvasObject

		if val.Qr.TemplateString != "" && data[val.Name] != "" {
			widgetQr := widget.NewButton(
				"QR",
				val.createQRCodeCallback(
					newWindow,
					data[val.Name],
				),
			)
			//widgetQr := widget.NewButtonWithIcon("", qrIcon, qrCallback)
			qrLayout = container.New(
				layout.NewBorderLayout(nil, nil, nil, widgetQr),
				widgetQr,
				newWidget,
			)
		}

		var widgetContent fyne.CanvasObject

		if qrLayout != nil {
			widgetContent = container.New(
				layout.NewBorderLayout(nil, nil, nil, widgetCopy),
				qrLayout,
				widgetCopy,
			)
		} else {
			widgetContent = container.New(
				layout.NewBorderLayout(nil, nil, nil, widgetCopy),
				newWidget,
				widgetCopy,
			)
		}

		newItem := widget.FormItem{
			Text:   val.GetDisplayName(),
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

func createCopyCallback(window fyne.Window, name string, value string) func() {
	return func() {
		window.Clipboard().SetContent(value)
		if a != nil {
			a.SendNotification(
				fyne.NewNotification(
					"Content Copied",
					fmt.Sprintf("%s copied to clipboard.", name),
				),
			)
		}
		window.Close()
	}
}

func (f *SearchFieldSettings) createQRCodeCallback(window fyne.Window, value string) func() {
	template := f.Qr.TemplateString // Do this to ensure correct assignment
	return func() {
		newWindow := createQRWindow(value, template)
		if newWindow != nil {
			newWindow.Show()
			window.Close()
		}
	}
}
