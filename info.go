package main

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	widget "fyne.io/fyne/v2/widget"
	"github.com/skip2/go-qrcode"
	"strconv"
	"text/template"
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

	copyIcon, err := fyne.LoadResourceFromPath("./icons/content_copy.svg")
	if err != nil {
		panic(err)
	}

	qrIcon, err := fyne.LoadResourceFromPath("./icons/qr_icon.svg")
	if err != nil {
		panic(err)
	}

	i := 0
	for _, val := range S.GetDisplayFields() {
		i++

		newWidget := widget.NewEntry()
		newWidget.Text = data[val.Name]

		name := val.GetDisplayName()
		copyValue := data[val.Name]
		copyCallback := func() {
			newWindow.Clipboard().SetContent(copyValue)
			// TODO: Fix Notification
			a.SendNotification(
				fyne.NewNotification(
					"Content Copied",
					fmt.Sprintf("%s copied to clipboard.", name),
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

		var qrLayout fyne.CanvasObject

		if val.Qr.TemplateString != "" && data[val.Name] != "" {

			qrValue := data[val.Name]
			qrField := val
			qrCallback := func() {
				window := newQrWindow(qrValue, qrField)
				if window != nil {
					window.Show()
					newWindow.Close()
				}
			}

			widgetQr := widget.NewButtonWithIcon("", qrIcon, qrCallback)
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

func newQrWindow(value string, settings FieldSettings) fyne.Window {
	windowName := fmt.Sprintf("Generated %s QR Code", settings.GetDisplayName())
	newWindow := a.NewWindow(windowName)
	desiredWidth := float32(400)

	qrReader := generateQR(value, settings.Qr.TemplateString)

	image := canvas.NewImageFromReader(qrReader, "qr.png")

	newWindow.SetContent(image)
	newWindow.Resize(fyne.Size{
		Width:  desiredWidth,
		Height: desiredWidth,
	})
	newWindow.CenterOnScreen()

	return newWindow
}

func generateQR(value string, templateString string) *bytes.Reader {
	var b bytes.Buffer
	tmpl, err := template.New("").Parse(templateString)
	if err != nil {
		fmt.Printf("QR error parsing template: %v\n", err)
	}
	err = tmpl.Execute(&b, value)
	if err != nil {
		fmt.Printf("QR error executing template: %v\n", err)
	}

	value = b.String()

	qr, err := qrcode.Encode(value, qrcode.High, 1024)
	if err != nil {
		fmt.Printf("QR error encoding to bytes: %v\n", err)
	}

	return bytes.NewReader(qr)
}
