package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type clientData struct {
	OldCode          string `struct:"old_code"`
	Code             string `struct:"code"`
	BusinessName     string `struct:"business_name"`
	Tier             int    `struct:"tier"`
	PrintList        int    `struct:"print_list"`
	TradingName      string `struct:"trading_name"`
	ContactFirstName string `struct:"contact_first_name"`
	ContactLastName  string `struct:"contact_last_name"`
	ContactEmail     string `struct:"contact_email"`
	ContactPhone     string `struct:"contact_phone"`
	ContactMobile    string `struct:"contact_mobile"`
	OfficeAddress    string `struct:"office_address"`
}

type clientDataSet []clientData

func (c clientDataSet) String(i int) string {
	meh := fmt.Sprintf(
		"%s %s %s %s",
		c[i].OldCode,
		c[i].Code,
		c[i].BusinessName,
		c[i].TradingName,
	)
	return meh
}

func (c clientDataSet) Len() int {
	return len(c)
}

var globalData = clientDataSet{}
var globalHeaders []string

func readCSV(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error in readCSV file open: %v\n", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("error in readCSV closing file: %v\n", err)
		}
	}(file)

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("error in readCSV reader: %v\n", err)
	}

	csvDataToStruct(data)
}

func csvDataToStruct(data [][]string) {
	var records []clientData
	for i, row := range data {
		var record clientData
		for j, col := range row {
			if i == 0 {
				header := strings.Trim(col, " ")
				header = strings.Replace(header, string('\uFEFF'), "", -1)
				header = strings.ToLower(header)
				globalHeaders = append(globalHeaders, header)
			} else {
				switch globalHeaders[j] {
				case "code":
					record.OldCode = col
				case "new code":
					record.Code = col
				case "business name":
					record.BusinessName = col
				case "tier level":
					record.Tier, _ = strconv.Atoi(col)
					// TODO: Handle Error
				case "print list":
					record.PrintList, _ = strconv.Atoi(col)
					// TODO: Handle Error
				case "trading as":
					record.TradingName = col
				case "name":
					record.ContactFirstName = col
				case "last name":
					record.ContactLastName = col
				case "email":
					record.ContactEmail = col
				case "phone number":
					record.ContactPhone = col
				case "mobile number":
					record.ContactMobile = col
				case "office address / location":
					record.OfficeAddress = col
				}
			}
		}
		records = append(records, record)
	}
	globalData = records
}
