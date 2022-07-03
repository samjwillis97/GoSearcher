package main

import (
	"github.com/sahilm/fuzzy"
)

func searchStruct(search string) []clientData {
	results := fuzzy.FindFrom(search, globalData)
	var searchResults []clientData
	for _, r := range results {
		searchResults = append(searchResults, globalData[r.Index])
	}
	return searchResults
}
