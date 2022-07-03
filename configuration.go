package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/spf13/viper"
	"log"
)

type Configuration struct {
	MaxEntries  int
	ClearOnHide bool
}

type Service struct {
	Name          string
	SourceFile    string
	FileType      string
	PrimaryField  []string
	SearchFields  []string
	DisplayFields []string
	FieldSettings map[string]FieldSettings
}

type FieldSettings struct {
	Display string
}

var C Configuration
var Services []Service

func setupConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
}

func readConfig() {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if a != nil {
				a.SendNotification(
					fyne.NewNotification(
						"Error Config File Not Found",
						"Reading of config file failed.",
					),
				)
			}
			return
		}
		panic(fmt.Errorf("error in readConfig viper: %v", err))
		return
	}

	err := viper.UnmarshalKey("configuration", &C)
	if err != nil {
		if a != nil {
			a.SendNotification(
				fyne.NewNotification(
					"Error Reading Config File",
					"Reading of config file failed.",
				),
			)
		}
		log.Fatalf("unable to decode into struct, %v", err)
	}

	err = viper.UnmarshalKey("services", &Services)
	if err != nil {
		if a != nil {
			a.SendNotification(
				fyne.NewNotification(
					"Error Reading Config File",
					"Reading of config file failed.",
				),
			)
		}
		log.Fatalf("unable to decode into struct, %v", err)
	}

	if w != nil {
		w.Close()
		w = a.NewWindow("Client Search")
	}

	return
}
