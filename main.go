package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

// TODO: Fix this - surely a better way then a global

// TODO: Setup Window Size

var a fyne.App
var w fyne.Window

func main() {

	viper.OnConfigChange(func(e fsnotify.Event) {
		readConfig()
	})
	readConfig()
	viper.WatchConfig()

	readCSV(viper.GetString("sourceFile"))

	a = app.New()
	if desk, ok := a.(desktop.App); ok {
		desk.SetSystemTrayMenu(setupMenu())
	}

	w = a.NewWindow("Client Search")
	setupWindow(w)

	a.Run()
}

func setupMenu() *fyne.Menu {
	var menus []*fyne.MenuItem

	menus = append(menus, fyne.NewMenuItem("Search", searchInterface))
	menus = append(menus, fyne.NewMenuItem("Refresh Source", refreshSource))

	return fyne.NewMenu("Menu :)", menus...)
}

func setupWindow(w fyne.Window) {
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter Text...")

	initData := make([]interface{}, 3)
	data := binding.BindUntypedList(&initData)

	list := widget.NewListWithData(
		data,
		func() fyne.CanvasObject {
			listItem := widget.NewButton("", func() {
			})
			listItem.Alignment = widget.ButtonAlignLeading
			return listItem
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			item, _ := i.(binding.Untyped).Get()
			switch item := item.(type) {
			case clientData:
				callBackFn := func() {
					// TODO: New Window here !
					window := createInfoWindow(item)
					if window != nil {
						window.Show()
					}
					if viper.GetBool("clearOnHide") {
						input.SetText("")
						_ = data.Set([]interface{}{})
					}
					input.SetPlaceHolder("Enter Text...")
					w.Hide()
					w.Resize(fyne.Size{
						Width:  500,
						Height: 0,
					})
				}
				o.(*widget.Button).SetText(item.BusinessName)
				o.(*widget.Button).OnTapped = callBackFn
				// FIXME: Try make enter click the button - space currently does
			}
		},
	)

	content := container.New(layout.NewBorderLayout(
		input,
		nil,
		nil,
		nil,
	),
		input,
		list,
	)

	w.SetContent(content)
	w.Resize(fyne.Size{
		Width:  500,
		Height: 0,
	})
	w.CenterOnScreen()

	w.SetCloseIntercept(func() {
		if viper.GetBool("clearOnHide") {
			input.SetText("")
			_ = data.Set([]interface{}{})
		}
		input.SetPlaceHolder("Enter Text...")
		w.Hide()
		w.Resize(fyne.Size{
			Width:  500,
			Height: 0,
		})
	})

	w.Canvas().Focus(input) // FIXME
	w.Canvas().Unfocus()

	input.OnChanged = func(s string) {
		results := searchStruct(s)

		var newData []interface{}
		for _, val := range results {
			newData = append(newData, val)
		}

		if len(newData) > 0 {
			_ = data.Set(newData)

			maxShown := float32(viper.GetInt("maxEntriesShown") - 1)
			baseListHeight := list.MinSize().Height
			newListHeight := maxShown * baseListHeight

			if len(newData) < int(maxShown) {
				newListHeight = float32(len(newData)-1) * baseListHeight
			}

			// Shows Input with 4 List items
			w.Resize(fyne.Size{
				Width:  500,
				Height: content.MinSize().Height + newListHeight,
			})
		}
	}
}

func searchInterface() {
	w.Show()
	w.CenterOnScreen()
	w.RequestFocus()
}

func refreshSource() {
	readCSV(viper.Get("sourceFile").(string))
	a.SendNotification(
		fyne.NewNotification(
			"Source Refreshed",
			"Source successfully refreshed",
		),
	)
}

func readConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			a.SendNotification(
				fyne.NewNotification(
					"Config File Not Found",
					"Reading of config file failed.",
				),
			)
		}
		panic(fmt.Errorf("error in readConfig viper: %v", err))
	}

	return
}

func onExit() {
	log.Println("Exiting")
	log.Println("Clean me up")
	// clean up here
}
