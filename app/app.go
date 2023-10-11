package app

import (
	"fmt"
	"log"
	"os"

	"github.com/andewx/microxt/common"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

const EVENT = 2
const EXIT = 1

type ElectronApp struct {
	App       *astilectron.Astilectron
	Endpoints map[string]*RemoteHandler
	Window    *astilectron.Window
	Finished  bool
	Port      int
}

type RemoteHandler struct {
	F func(string) (string, error)
}

func New() (*ElectronApp, error) {
	var myApp = new(ElectronApp)
	var err error
	myApp.Port = 8000
	myApp.App, err = astilectron.New(log.New(os.Stderr, "", 0), astilectron.Options{
		AppName:            "MicroXT Radar Application",
		AppIconDefaultPath: common.ProjectRelativePath("microxt/app/icons/microxt.png"),  // If path is relative, it must be relative to the data directory
		AppIconDarwinPath:  common.ProjectRelativePath("microxt/app/icons/microxt.icns"), // Same here
		BaseDirectoryPath:  common.ProjectRelativePath("microxt/app"),
		TCPPort:            &myApp.Port,
	})

	if err != nil {
		fmt.Printf("Error Initializing Application \n%s", err)
		return myApp, err
	}

	myApp.Endpoints = make(map[string]*RemoteHandler)
	myApp.Finished = false
	return myApp, err
}

func (app *ElectronApp) Init() error {
	var sErr = app.App.Start() //Blocking

	if sErr != nil {
		return sErr
	}

	var w, err = app.App.NewWindow(common.ProjectRelativePath("microxt/app/site/index.html"), &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(1024),
		Width:  astikit.IntPtr(1280),
	})

	if err != nil {
		return err
	}

	app.Window = w
	err = app.Window.Create()
	if err != nil {
		return err
	}

	app.Window.OnMessage(func(m *astilectron.EventMessage) interface{} {
		var s string
		m.Unmarshal(&s)
		//	var apiCall = app.Endpoints[s]
		//Call API Function
		return nil
	})

	// Add a listener on Astilectron
	app.App.On(astilectron.EventNameAppCrash, func(e astilectron.Event) (deleteListener bool) {
		fmt.Println("App has crashed")
		app.Finished = true
		return
	})

	// Add a listener on the window
	app.Window.On(astilectron.EventNameWindowEventResize, func(e astilectron.Event) (deleteListener bool) {
		return
	})

	// Add a listener on the window
	app.Window.On(astilectron.EventNameWindowEventClosed, func(e astilectron.Event) (deleteListener bool) {
		app.Finished = true
		return
	})

	// Add a listener on the window
	app.Window.On(astilectron.EventNameWindowCmdDestroy, func(e astilectron.Event) (deleteListener bool) {
		app.Finished = true
		return
	})

	// Add a listener on the window
	app.Window.On(astilectron.EventNameWindowCmdClose, func(e astilectron.Event) (deleteListener bool) {
		app.Finished = true
		return
	})

	return err
}

func (app *ElectronApp) RunEventLoop(status chan int) {
	defer app.App.Close()
	for !app.Finished {
		app.App.Wait()
		status <- EVENT
	}
	status <- EXIT
}

func (app *ElectronApp) PassContext() (*astilectron.Astilectron, *astilectron.Window) {
	return app.App, app.Window
}
