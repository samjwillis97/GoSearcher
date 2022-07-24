package main

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/skip2/go-qrcode"
	"text/template"
)

func createQRWindow(value string, template string) fyne.Window {
	windowName := fmt.Sprintf("Generated QR Code")
	newWindow := a.NewWindow(windowName)
	desiredWidth := float32(400)

	qrReader := generateQRCode(value, template)

	image := canvas.NewImageFromReader(qrReader, "qr.png")

	newWindow.SetContent(image)
	newWindow.Resize(fyne.Size{
		Width:  desiredWidth,
		Height: desiredWidth,
	})
	newWindow.CenterOnScreen()

	return newWindow
}

func generateQRCode(value string, templateString string) *bytes.Reader {
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
