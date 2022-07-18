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
	Name         string
	SourceFile   string
	FileType     string
	Fields       []FieldSettings
	FileSettings FileSettings
}

type FieldSettings struct {
	Name        string
	DisplayName string
	Search      bool
	Primary     bool
	Display     bool
	Qr          QRSettings
}

type FileSettings struct {
	Sheet            string
	NumberOfSkipRows int
}

type QRSettings struct {
	TemplateString string
}

var C Configuration
var Services []Service

func setupConfig() {
	// TODO: Change Config location depending on system
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

func (s *Service) GetPrimaryFields() []string {
	var fields []string
	for _, val := range s.Fields {
		if val.Primary {
			fields = append(fields, val.Name)
		}
	}
	if len(fields) > 0 {
		return fields
	}
	if len(s.GetDisplayFields()) > 0 {
		return []string{s.GetDisplayFields()[0].Name}
	}
	if len(s.GetSearchFields()) > 0 {
		return []string{s.GetSearchFields()[0]}
	}
	return nil
}

func (s *Service) GetDisplayFields() []FieldSettings {
	var fields []FieldSettings
	for _, val := range s.Fields {
		if val.Display {
			fields = append(fields, val)
		} else if val.DisplayName != "" {
			fields = append(fields, val)
		}
	}
	if len(fields) == 0 {
		for _, val := range s.Fields {
			fields = append(fields, val)
		}
	}
	return fields
}

func (s *Service) GetSearchFields() []string {
	var fields []string
	for _, val := range s.Fields {
		if val.Search {
			fields = append(fields, val.Name)
		}
	}
	if len(fields) == 0 {
		// TODO: Error out
	}
	return fields
}

func (f *FieldSettings) GetDisplayName() string {
	if f.DisplayName == "" {
		return f.Name
	}
	return f.DisplayName
}
