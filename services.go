package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"runtime"
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
			return listItem
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			item, _ := i.(binding.Untyped).Get()
			switch item := item.(type) {
			case string:
				serviceName := item

				var foundService Service
				for _, val := range Services {
					if val.Name == serviceName {
						foundService = val
					}
				}

				callBackFn := func() {
					found := foundService
					startService(found)
					w.Close()
				}

				o.(*widget.Button).SetText(item)
				o.(*widget.Button).OnTapped = callBackFn
			}
		},
	)

	var allServices []string
	for _, val := range getServicesInBindingOrder() {
		initData = append(initData, val.Name)
		allServices = append(allServices, val.Name)
	}

	shortCuts := map[fyne.Shortcut]func(shortcut fyne.Shortcut){}
	assignedBindings := make(map[string]struct{})
	modifier := fyne.KeyModifierControl
	if runtime.GOOS == "darwin" {
		modifier = fyne.KeyModifierSuper
	}
	for _, val := range getServicesInBindingOrder() {
		boundKey := val.GetServiceKeybinding()
		if _, ok := assignedBindings[boundKey]; ok {
			continue
		}

		assignedBindings[boundKey] = struct{}{}
		keyBinding := &desktop.CustomShortcut{
			KeyName:  fyne.KeyName(boundKey),
			Modifier: modifier,
		}

		callback := func() {
			found := val
			startService(found)
			w.Close()
		}

		shortCuts[keyBinding] = func(shortcut fyne.Shortcut) {
			callback()
		}
	}

	content := container.New(layout.NewBorderLayout(
		input,
		nil,
		nil,
		nil,
	),
		input,
		list,
	)

	// Set Initial List
	newListHeight := float32(0)
	if len(initData) > 0 {
		_ = data.Set(initData)

		maxShown := float32(C.MaxEntries)
		baseListHeight := list.MinSize().Height
		newListHeight = maxShown * baseListHeight

		if len(initData) < int(maxShown) {
			newListHeight = float32(len(initData)-1)*baseListHeight + 10
		}
	}

	input.OnChanged = func(text string) {
		results := allServices
		if text != "" {
			results = fuzzySearch(text, results, C.Similarity, C.GetSearchAlgorithm())
		}

		var newData []interface{}
		for _, val := range results {
			newData = append(newData, val)
		}

		_ = data.Set(newData)
		if len(newData) > 0 {
			maxShown := float32(C.MaxEntries)
			baseListHeight := list.MinSize().Height
			newListHeight := maxShown * baseListHeight

			if len(newData) < int(maxShown) {
				newListHeight = float32(len(newData)-1)*baseListHeight + 10
			}

			w.Resize(fyne.Size{
				Width:  500,
				Height: content.MinSize().Height + newListHeight,
			})
		} else {
			w.Resize(fyne.Size{
				Width:  500,
				Height: input.MinSize().Height,
			})
		}
	}

	w.SetContent(content)
	for key, val := range shortCuts {
		w.Canvas().AddShortcut(key, val)
	}
	w.Resize(fyne.Size{
		Width:  500,
		Height: content.MinSize().Height + newListHeight,
	})
	w.CenterOnScreen()

	w.Canvas().Focus(input)

}
