package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

// TODO: Fix this - surely a better way then a global

// TODO: Setup Window Size

var a fyne.App
var w fyne.Window

func main() {
	// TODO: Log Better

	// open a file
	f, err := os.OpenFile(os.TempDir()+string(os.PathSeparator)+"GoSearcher.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// don't forget to close it
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	// Output to stderr instead of stdout, could also be a file.
	//log.SetOutput(f)
	log.Println("GoSearcher Initd")

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config Modified - Reloading")
		readConfig()
		if a != nil {
			if desk, ok := a.(desktop.App); ok {
				//desk.SetSystemTrayIcon()
				desk.SetSystemTrayMenu(setupTrayMenu())
			}
		}
	})

	setupConfig()
	readConfig()

	viper.WatchConfig()
	a = app.New()
	if desk, ok := a.(desktop.App); ok {
		//desk.SetSystemTrayIcon()
		desk.SetSystemTrayMenu(setupTrayMenu())
	}
	a.Run()
}

func setupTrayMenu() *fyne.Menu {
	var menus []*fyne.MenuItem

	for _, service := range Services {
		serviceToAssign := service

		menus = append(menus, fyne.NewMenuItem(serviceToAssign.Name, func() {
			switch serviceToAssign.GetServiceType() {
			case "search":
				createSearchInterface(serviceToAssign)
			default:
				log.Printf("Unknown Service Type: %s", serviceToAssign.GetServiceType())
			}
		}))
	}

	return fyne.NewMenu("System Tray", menus...)
}

func createSearchInterface(s Service) {
	s.loadSearchData()
	w = a.NewWindow(s.Name)
	initSearchWindow(w, s)
	w.Show()
	w.CenterOnScreen()
	w.RequestFocus()
}
