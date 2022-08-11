package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
)

func startService(s Service) {
	switch s.GetServiceType() {
	case "search":
		startSearchInterface(s)
	case "qr-code-generator":
		startQRGeneratorInterface(s)
	default:
		log.Printf("Unknown Service Type: %s", s.GetServiceType())
	}
}

func startSearchInterface(s Service) {
	s.loadSearchData()
	w = a.NewWindow(s.Name)
	initSearchWindow(w, s)
	w.Show()
	w.CenterOnScreen()
	w.RequestFocus()
}

func startQRGeneratorInterface(s Service) {
	w = s.QRSettings.createPromptWindow()
	w.Show()
	w.CenterOnScreen()
	w.RequestFocus()
}

func startServiceSearchInterface() {
	w = a.NewWindow("Services")
	initServiceSearchWindow(w)
	w.Show()
	w.CenterOnScreen()
	w.RequestFocus()
}

// TODO: Clean this up? Code duplication
func initServiceSearchWindow(w fyne.Window) {
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
			log.Println("CreateItem")
			return listItem
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			item, _ := i.(binding.Untyped).Get()
			switch item := item.(type) {
			case string:
				o.(*widget.Button).SetText(item)

				callBackFn := func() {
					serviceName := item
					var foundService Service
					for _, val := range Services {
						if val.Name == serviceName {
							foundService = val
						}
					}
					startService(foundService)
					w.Close()
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

	var allServices []string
	for _, val := range Services {
		initData = append(initData, val.Name)
		allServices = append(allServices, val.Name)
	}
	if len(initData) > 0 {
		_ = data.Set(initData)

		maxShown := float32(C.MaxEntries - 1)
		baseListHeight := list.MinSize().Height
		newListHeight := maxShown * baseListHeight

		if len(initData) < int(maxShown) {
			newListHeight = float32(len(initData)-1) * baseListHeight
		}

		// Shows Input with 4 List items
		w.Resize(fyne.Size{
			Width:  500,
			Height: content.MinSize().Height + newListHeight,
		})
	}

	input.OnChanged = func(text string) {
		services := allServices
		results := fuzzySearch(text, services)

		var newData []interface{}
		for _, val := range results {
			newData = append(newData, val)
		}

		if len(newData) > 0 {
			_ = data.Set(newData)

			maxShown := float32(C.MaxEntries - 1)
			baseListHeight := list.MinSize().Height
			newListHeight := maxShown * baseListHeight

			if len(newData) < int(maxShown) {
				newListHeight = float32(len(newData)-1) * baseListHeight
			}

			w.Resize(fyne.Size{
				Width:  500,
				Height: content.MinSize().Height + newListHeight,
			})
		}
	}
}
