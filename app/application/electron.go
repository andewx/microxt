package application

import (
	"fmt"
	"log"
	"os"

	"github.com/andewx/microxt/app/models"
	"github.com/andewx/microxt/common"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

const EVENT = 2
const EXIT = 1
const START = 12

type ElectronApp struct {
	Electron *astilectron.Astilectron
	Window   *astilectron.Window
	App      Application
	Finished bool
	Port     int
}

func NewElectron(application Application) (*ElectronApp, error) {
	var myApp = new(ElectronApp)
	var err error
	myApp.Port = 8000
	myApp.Electron, err = astilectron.New(log.New(os.Stderr, "", 0), astilectron.Options{
		AppName:            "AirXT",
		AppIconDefaultPath: common.ProjectRelativePath("microxt/app/icons/icon.png"),  // If path is relative, it must be relative to the data directory
		AppIconDarwinPath:  common.ProjectRelativePath("microxt/app/icons/icon.icns"), // Same here
		BaseDirectoryPath:  common.ProjectRelativePath("microxt/app"),
		SkipSetup:          false,
		TCPPort:            &myApp.Port,
	})

	myApp.App = application

	if err != nil {
		fmt.Printf("Error Initializing \n%s", err)
		return myApp, err
	}

	return myApp, err
}

// Launches electron completely as its own process
func LaunchElectron(status chan int, app Application) error {
	el, err1 := NewElectron(app)
	if err1 != nil {
		fmt.Printf("Error Initializing \n%s", err1)
		status <- EXIT
		return err1
	}

	defer el.Electron.Close()

	el.Electron.HandleSignals()

	var sErr = el.Electron.Start()

	if sErr != nil {
		fmt.Printf("Error Initializing \n%s", sErr)
		status <- EXIT
		return sErr
	}

	var w, err = el.Electron.NewWindow(common.ProjectRelativePath("microxt/app/site/index.html"), &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(1024),
		Width:  astikit.IntPtr(1280),
	})

	if err != nil {
		fmt.Printf("Error Initializing \n%s", err)
		status <- EXIT
		return err
	}

	el.Window = w

	err = el.Window.Create()
	if err != nil {
		fmt.Printf("Error Initializing \n%s", err)
		status <- EXIT
		return err
	}

	status <- START

	el.App.SetWindow(el.Window)

	el.Window.OpenDevTools()

	el.Window.OnMessage(func(m *astilectron.EventMessage) interface{} {
		var s string

		m.Unmarshal(&s)
		fmt.Printf("Application recieved message:\n%s\n", s)
		request := models.NewJsonRequest(s)
		var session *Session
		if request != nil {

			if request.SessionKey == "0" {
				session = el.App.NewSession()
			} else {
				session = el.App.GetSession(request.SessionKey)
				if session == nil {
					fmt.Printf("Session not found for key %s\n", request.SessionKey)
					//Print all sessions
					for k := range el.App.GetKeys() {
						fmt.Printf("Session %s\n", k)
					}
					return nil
				}
			}
			route := el.App.GetRoute(request.RouteKey)
			fmt.Printf("Calling route %s\n", request.RouteKey)
			if route != nil {
				fmt.Printf("Route found\n")
				go route.Handler(request.Paramaters, session, el.App)
			} else {
				fmt.Printf("Route not found\n")
			}
		}
		return nil
	})

	// Add a listener on Astilectron
	el.Electron.On(astilectron.EventNameAppCrash, func(e astilectron.Event) (deleteListener bool) {
		fmt.Println("App has crashed")
		el.App.Close()
		el.Finished = true
		return
	})

	// Add a listener on the window
	el.Window.On(astilectron.EventNameWindowEventResize, func(e astilectron.Event) (deleteListener bool) {
		return
	})

	// Add a listener on the window
	el.Window.On(astilectron.EventNameWindowEventClosed, func(e astilectron.Event) (deleteListener bool) {
		el.Finished = true
		el.App.Close()
		return
	})

	// Add a listener on the window
	el.Window.On(astilectron.EventNameWindowCmdDestroy, func(e astilectron.Event) (deleteListener bool) {
		el.Finished = true
		el.App.Close()
		return
	})

	// Add a listener on the window
	el.Window.On(astilectron.EventNameWindowCmdClose, func(e astilectron.Event) (deleteListener bool) {
		el.Finished = true
		el.App.Close()
		return
	})
	el.Electron.Wait()
	status <- EXIT
	return nil
}
