package main

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/sahilm/fuzzy"
	"log"
)

func searchCurrentData(search string) []map[string]string {
	// TODO: Change this to the one Jeremy recommended
	results := fuzzy.FindFrom(search, searchData)
	searchResults := make([]map[string]string, 0)
	for _, r := range results {
		searchResults = append(searchResults, dataSet[r.Index])
	}
	return searchResults
}

func fuzzySearch(search string, data []string) []string {
	// TODO: Change this to the one Jeremy recommended
	results := fuzzy.Find(search, data)
	var searchResults []string
	for _, r := range results {
		searchResults = append(searchResults, r.Str)
	}
	return searchResults
}

// TODO: Clean this up - break up into smaller functions?
func initSearchWindow(w fyne.Window, S Service) {
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

				var primaryField = S.GetPrimaryFields()
				if len(primaryField) > 0 {
					var text string
					if len(primaryField) > 1 {
						for _, val := range primaryField {
							text = text + value[val] + " "
						}
					} else {
						text = value[primaryField[0]]
					}
					o.(*widget.Button).SetText(text)
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

	w.Canvas().Focus(input)

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
