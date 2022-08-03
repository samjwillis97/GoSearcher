package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/spf13/viper"
	"log"
	"os"
	"runtime"
	"strings"
)

type Configuration struct {
	MaxEntries  int
	ClearOnHide bool
}

type Service struct {
	Name           string
	Type           string
	SearchSettings SearchSettings
	Fields         []SearchFieldSettings
	FileSettings   FileSettings
	QRSettings     QRGeneratorSettings
}

type SearchSettings struct {
	Modifier string
}

type SearchFieldSettings struct {
	Name        string
	DisplayName string
	Search      bool
	Primary     bool
	Display     bool
	KeyBinding  string
	Qr          SearchQRSettings
}

type FileSettings struct {
	Source           string
	Type             string
	Sheet            string
	NumberOfSkipRows int
}

type SearchQRSettings struct {
	TemplateString string
}

type QRGeneratorSettings struct {
	Inputs         []string
	TemplateString string
}

var C Configuration
var Services []Service

func setupConfig() {
	configDir, _ := os.UserConfigDir()
	viper.SetConfigName("config")
	viper.AddConfigPath(configDir + string(os.PathSeparator) + "GoSearcher")
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

	return
}

func (s *Service) GetServiceType() string {
	return strings.ToLower(s.Type)
}

func (s *Service) GetSourceFilePath() string {
	// TODO: Throw error
	return s.FileSettings.Source
}

func (s *Service) GetSourceFileType() string {
	// TODO: default to SourceFile
	return s.FileSettings.Type
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

func (s *Service) GetDisplayFields() []SearchFieldSettings {
	var fields []SearchFieldSettings
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

func (s *Service) GetSearchModifierKey() fyne.KeyModifier {
	switch strings.ToLower(s.SearchSettings.Modifier) {
	case "ctrl":
	case "control":
		return fyne.KeyModifierControl
	case "shift":
		return fyne.KeyModifierShift
	case "alt":
	case "option":
		return fyne.KeyModifierAlt
	case "windows":
	case "super":
	case "meta":
	case "cmd":
	case "command":
		return fyne.KeyModifierSuper
	}

	switch runtime.GOOS {
	case "windows":
		return fyne.KeyModifierControl
	case "darwin":
	case "linux":
		return fyne.KeyModifierAlt
	}

	return fyne.KeyModifierControl
}

func (f *SearchFieldSettings) GetDisplayName() string {
	if f.DisplayName == "" {
		return f.Name
	}
	return f.DisplayName
}

func (f *SearchFieldSettings) GetKeyBinding() string {
	return strings.ToUpper(f.KeyBinding)
}
