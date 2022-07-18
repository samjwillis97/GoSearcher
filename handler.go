package main

import (
	"encoding/csv"
	"fyne.io/fyne/v2"
	"github.com/xuri/excelize/v2"
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
	switch strings.ToLower(s.FileType) {
	case "csv":
		s.loadFromCSV(s.SourceFile)
	case "xlsx":
		s.loadFromXLSX(s.SourceFile)
	}

}

// FIXME: Consider refactoring these shitloads of repeated code
func (s *Service) loadFromCSV(path string) {
	f, err := os.Open(path)
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

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("error in readCSV closing file: %v\n", err)
		}
	}(f)

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("error in readCSV reader: %v\n", err)
	}

	clearData()

	searchCols := make(map[int]string)
	displayCols := make(map[int]string)

	for i, row := range data {
		searchValue := ""
		var dataSetEntry = map[string]string{}
		if i >= s.FileSettings.NumberOfSkipRows {
			for j, col := range row {
				if i == s.FileSettings.NumberOfSkipRows { // Header Row
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
		}
		if searchValue != "" {
			// THESE MUST STAY TOGETHER
			searchData = append(searchData, searchValue)
			dataSet = append(dataSet, dataSetEntry)
		}
	}
}

func (s *Service) loadFromXLSX(path string) {
	f, err := excelize.OpenFile(path)
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
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			log.Fatalf("error in readXLSX closing file: %v\n", err)
		}
	}()

	clearData()

	rows, err := f.GetRows(s.FileSettings.Sheet)
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

	searchCols := make(map[int]string)
	displayCols := make(map[int]string)

	for i, row := range rows {
		searchValue := ""
		var dataSetEntry = map[string]string{}
		if i >= s.FileSettings.NumberOfSkipRows {
			for j, col := range row {
				if i == s.FileSettings.NumberOfSkipRows { // Header Row
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
