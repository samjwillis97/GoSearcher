package main

import (
	"encoding/csv"
	"fyne.io/fyne/v2"
	"log"
	"os"
	"strings"
)

type SearchDataSet []string

func (s SearchDataSet) String(i int) string {
	return s[i]
}

func (s SearchDataSet) Len() int {
	return len(s)
}

var dataSet []map[string]string
var searchData = SearchDataSet{}

func (s *Service) loadData() {
	file, err := os.Open(s.SourceFile)
	if err != nil {
		log.Printf("error in loadData: %v\n", err)
		if a != nil {
			a.SendNotification(
				fyne.NewNotification(
					"Error Loading Service",
					"Could not open Source File.",
				),
			)
		}
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("error in readCSV closing file: %v\n", err)
		}
	}(file)

	switch strings.ToLower(s.FileType) {
	case "csv":
		s.loadFromCSV(file)
	}

}

func (s *Service) loadFromCSV(f *os.File) {
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("error in readCSV reader: %v\n", err)
	}

	clearData()

	searchCols := make(map[int]string)
	displayCols := make(map[int]string)

	// Need LUT type of thing
	for i, row := range data {
		searchValue := ""
		var dataSetEntry = map[string]string{}
		for j, col := range row {
			if i == 0 { // Header Row
				for _, val := range s.GetSearchFields() {
					if col == val {
						searchCols[j] = val
					}
				}
				for _, val := range s.GetDisplayFields() {
					if col == val.Name {
						displayCols[j] = val.Name
					}
				}
			} else {
				if _, ok := searchCols[j]; ok {
					searchValue += col + " "
				}
				if val, ok := displayCols[j]; ok {
					dataSetEntry[val] = col
				}
			}
		}
		if searchValue != "" {
			// THESE MUST STAY TOGETHER
			searchData = append(searchData, searchValue)
			dataSet = append(dataSet, dataSetEntry)
		}
	}
}

func clearData() {
	dataSet = []map[string]string{}
	searchData = SearchDataSet{}
}
