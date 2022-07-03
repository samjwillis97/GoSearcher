package main

import (
	"encoding/json"
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

	setupConfig()
	readConfig()

	viper.WatchConfig()

	a = app.New()
	if desk, ok := a.(desktop.App); ok {
		desk.SetSystemTrayMenu(setupMenu())
	}

	a.Run()
}

func setupMenu() *fyne.Menu {
	var menus []*fyne.MenuItem

	for _, service := range Services {
		serviceToAssign := service
		menus = append(menus, fyne.NewMenuItem(serviceToAssign.Name, func() {
			searchInterface(serviceToAssign)
		}))
	}

	return fyne.NewMenu("Menu :)", menus...)
}

func setupWindow(w fyne.Window, S Service) {
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter Text...")

	initData := make([]interface{}, 0)
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
			case string:
				var value map[string]string
				byteValue := []byte(item)
				err := json.Unmarshal(byteValue, &value)
				if err != nil {
					log.Printf("error unmarshalling json: %v\n", err)
				}

				if len(S.PrimaryField) > 0 {
					var text string
					if len(S.PrimaryField) > 1 {
						for _, val := range S.PrimaryField {
							text = text + value[val] + " "
						}
					} else {
						text = value[S.PrimaryField[0]]
					}
					o.(*widget.Button).SetText(text)
				} else if len(S.DisplayFields) > 0 {
					o.(*widget.Button).SetText(value[S.DisplayFields[0]])
				} else {
					o.(*widget.Button).SetText(value[S.SearchFields[0]])
				}

				callBackFn := func() {
					// TODO: New Window here !
					window := createInfoWindow(value, S)
					if window != nil {
						window.Show()
					}
					if C.ClearOnHide {
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
				o.(*widget.Button).OnTapped = callBackFn
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
		w.Close()

		// Clear out memory
		dataSet = []map[string]string{}
		searchData = []string{}
	})

	w.Canvas().Focus(input) // FIXME
	w.Canvas().Unfocus()

	input.OnChanged = func(s string) {
		results := searchCurrentData(s)

		var newData []interface{}
		for _, val := range results {
			jsonBytes, err := json.Marshal(val)
			if err != nil {
				log.Printf("error json.Marshal: %v\n", err)
			}
			jsonString := string(jsonBytes)
			// value appended to newData must be comparable
			// could use indices instead lol
			newData = append(newData, jsonString)
		}

		if len(newData) > 0 {
			_ = data.Set(newData)

			maxShown := float32(C.MaxEntries - 1)
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

func searchInterface(s Service) {
	// Implementation loads data into memory
	log.Println(s)

	s.loadData()

	w = a.NewWindow(s.Name)
	setupWindow(w, s)
	w.Show()
	w.CenterOnScreen()
	w.RequestFocus()
}
