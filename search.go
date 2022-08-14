package main

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/hbollon/go-edlib"
	"log"
)

func searchCurrentData(search string, similarity float32, algorithm edlib.Algorithm) []map[string]string {
	// TODO: Change to use edlib
	//results := fuzzy.FindFrom(search, searchData)
	//searchResults := make([]map[string]string, 0)
	//for _, r := range results {
	//	searchResults = append(searchResults, dataSet[r.Index])
	//}
	//return searchResults
	results := fuzzySearch(search, searchData, similarity, algorithm)
	searchResults := make([]map[string]string, 0)
	for _, searchVal := range results {
		for i, resultVal := range searchData {
			if searchVal == resultVal {
				searchResults = append(searchResults, dataSet[i])
				break
			}
		}
	}
	return searchResults
}

func fuzzySearch(search string, data []string, similarity float32, algorithm edlib.Algorithm) []string {
	results, err := edlib.FuzzySearchSetThreshold(
		search,
		data,
		C.MaxEntries,
		similarity,
		algorithm,
	)
	if err != nil {
		log.Printf("error fuzzy searching: %v\n", err)
	}

	i := 0
	for range results {
		// Delete element from array if it is empty string and keep order
		if results[i] == "" {
			results = append(results[:i], results[i+1:]...)
		} else {
			i++
		}
	}

	return results
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
		//results := searchCurrentData(s)
		results := searchCurrentData(
			s,
			S.SearchSettings.Similarity,
			S.SearchSettings.GetSearchAlgorithm(),
		)

		var newData []interface{}
		for _, val := range results {
			jsonBytes, err := json.Marshal(val)
			if err != nil {
				log.Printf("error json.Marshal: %v\n", err)
			}
			jsonString := string(jsonBytes)
			// value appended to newData must be comparable
			newData = append(newData, jsonString)
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
				Height: content.MinSize().Height,
			})
		}
	}
}
