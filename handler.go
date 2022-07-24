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

// TODO: Work out how to not use globals :)
var dataSet []map[string]string
var searchData = SearchDataSet{}

func (s *Service) loadSearchData() {

	var data [][]string

	switch strings.ToLower(s.GetSourceFileType()) {
	case "csv":
		data = s.loadRowsFromCSV()
	case "xlsx":
		data = s.loadRowsFromXLSX()
	}

	clearGlobalData()
	searchData, dataSet = s.parseTabulatedData(data)
}

func (s *Service) loadRowsFromCSV() [][]string {
	f, err := os.Open(s.GetSourceFilePath())
	if err != nil {
		log.Printf("error in loadSearchData: %v\n", err)
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
	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("error in readCSV reader: %v\n", err)
	}

	return rows
}

func (s *Service) loadRowsFromXLSX() [][]string {
	f, err := excelize.OpenFile(s.GetSourceFilePath())
	if err != nil {
		log.Printf("error in loadSearchData: %v\n", err)
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

	rows, err := f.GetRows(s.FileSettings.Sheet)
	if err != nil {
		log.Printf("error in loadSearchData: %v\n", err)
		if a != nil {
			a.SendNotification(
				fyne.NewNotification(
					"Error Loading Service",
					"Could not open Source File.",
				),
			)
		}
	}

	return rows
}

func (s *Service) parseTabulatedData(table [][]string) ([]string, []map[string]string) {
	searchCols := make(map[int]string)
	displayCols := make(map[int]string)

	var searchStrings []string
	var data []map[string]string

	for i, row := range table {
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
			searchStrings = append(searchStrings, searchValue)
			data = append(data, dataSetEntry)
		}
	}
	return searchStrings, data
}

func clearGlobalData() {
	dataSet = []map[string]string{}
	searchData = SearchDataSet{}
}
