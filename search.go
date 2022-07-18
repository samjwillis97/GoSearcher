package main

import (
	"github.com/sahilm/fuzzy"
)

func searchCurrentData(search string) []map[string]string {
	results := fuzzy.FindFrom(search, searchData)
	searchResults := make([]map[string]string, 0)
	for _, r := range results {
		searchResults = append(searchResults, dataSet[r.Index])
	}
	return searchResults
}
