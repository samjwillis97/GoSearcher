package main

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/skip2/go-qrcode"
	"text/template"
)

func createQRWindow(value string, template string) fyne.Window {
	windowName := fmt.Sprintf("Generated QR Code")
	newWindow := a.NewWindow(windowName)
	desiredWidth := float32(400)

	valueMap := map[string]interface{}{
		"value": value,
	}
	qrReader := generateQRCode(valueMap, template)

	image := canvas.NewImageFromReader(qrReader, "qr.png")

	newWindow.SetContent(image)
	newWindow.Resize(fyne.Size{
		Width:  desiredWidth,
		Height: desiredWidth,
	})
	newWindow.CenterOnScreen()

	return newWindow
}

func (q *QRGeneratorSettings) createPromptWindow() fyne.Window {
	windowName := "QR Code Generator"
	mainWindow := a.NewWindow(windowName)
	desiredWidth := float32(400)

	valueMap := map[string]interface{}{}

	form := widget.NewForm()
	for _, val := range q.Inputs {
		key := val
		valueMap[val] = ""

		textEntry := widget.NewEntry()
		textEntry.OnChanged = func(text string) {
			valueMap[key] = text
		}

		item := widget.NewFormItem(key, textEntry)

		form.AppendItem(item)
	}

	button := widget.NewButton("Generate", func() {
		qrReader := generateQRCode(valueMap, q.TemplateString)

		newWindow := a.NewWindow("Generated QR Code")
		width := float32(400)
		image := canvas.NewImageFromReader(qrReader, "qr.png")

		newWindow.SetContent(image)
		newWindow.Resize(fyne.Size{
			Width:  width,
			Height: width,
		})
		newWindow.CenterOnScreen()

		newWindow.Show()
		mainWindow.Close()
	})

	content := container.NewVBox(form, button)

	mainWindow.SetContent(content)
	mainWindow.Resize(fyne.Size{
		Width:  desiredWidth,
		Height: 0,
	})
	mainWindow.CenterOnScreen()

	return mainWindow
}

func generateQRCode(value map[string]interface{}, templateString string) *bytes.Reader {
	var b bytes.Buffer
	tmpl, err := template.New("").Parse(templateString)
	if err != nil {
		fmt.Printf("QR error parsing template: %v\n", err)
	}
	err = tmpl.Execute(&b, value)
	if err != nil {
		fmt.Printf("QR error executing template: %v\n", err)
	}

	outputString := b.String()

	qr, err := qrcode.Encode(outputString, qrcode.High, 1024)
	if err != nil {
		fmt.Printf("QR error encoding to bytes: %v\n", err)
	}

	return bytes.NewReader(qr)
}
