package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
)

// TODO: Fix this - surely a better way then a global
var a fyne.App
var w fyne.Window

// TODO: Add socket to lock file? or Add to Config..

type Listener struct{}
type Arg struct{}
type Reply struct{}

func (l *Listener) OpenService(arg Arg, reply *Reply) error {
	startServiceSearchInterface()
	return nil
}

func main() {
	// TODO: Log Better

	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	lockFilePath := path.Join(dir, "GoSearcher.lock")
	lockFile, err := createLockFile(lockFilePath)
	if err != nil {
		data, err := os.ReadFile(lockFilePath)
		log.Printf("Lock file read: %s", lockFilePath)
		if err != nil {
			log.Fatalf("error reading lock file: %v", err)
		}
		client, err := rpc.DialHTTP("tcp", string(data))
		if err != nil {
			log.Fatal("dialing:", err)
		}
		var reply struct{}
		err = client.Call("Listener.OpenService", Arg{}, &reply)
		if err != nil {
			log.Fatal("Listener.OpenService:", err)
		}
		return
	}

	log.Printf("Lock file created: %s", lockFile.Name())

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
				desk.SetSystemTrayMenu(setupTrayMenu())
			}
		}
	})

	setupConfig()
	readConfig()

	go startRPCServer(lockFile)

	viper.WatchConfig()
	a = app.New()
	if desk, ok := a.(desktop.App); ok {
		//desk.SetSystemTrayIcon()
		desk.SetSystemTrayMenu(setupTrayMenu())
	}

	go func() {
		killchan := make(chan os.Signal, 2)
		signal.Notify(killchan, os.Interrupt, syscall.SIGTERM)
		// wait for kill signal
		<-killchan
		log.Println("Kill sig!")
		// TODO
		//do clean up
		//now exit
		os.Exit(0)
	}()

	a.Run()
}

func setupTrayMenu() *fyne.Menu {
	var menus []*fyne.MenuItem

	for _, service := range Services {
		serviceToAssign := service

		menus = append(menus, fyne.NewMenuItem(serviceToAssign.Name, func() {
			startService(serviceToAssign)
		}))
	}

	return fyne.NewMenu("System Tray", menus...)
}

func createLockFile(filename string) (*os.File, error) {
	switch runtime.GOOS {
	case "windows":
		if _, err := os.Stat(filename); err == nil {
			err = os.Remove(filename)
			if err != nil {
				return nil, err
			}
		}
		return os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	case "linux":
	case "darwin":
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
	return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
}

func startRPCServer(lockFile *os.File) {
	log.Println("Starting RPC Server")
	listener := new(Listener)
	err := rpc.Register(listener)
	if err != nil {
		log.Fatalf("error registering: %v", err)
	}

	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	go func() {
		_, err = lockFile.WriteString(string(l.Addr().String()))
		if err != nil {
			log.Fatalf("file write error: %v", err)
		}
		defer func() {
			err := lockFile.Close()
			if err != nil {
				log.Fatalf("lock file close error: %v", err)
			}
		}()

		err := http.Serve(l, nil)
		if err != nil {
			log.Fatal("serve error:", err)
		}
	}()
}
