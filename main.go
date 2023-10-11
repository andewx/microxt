package main

import (
	"fmt"

	"github.com/andewx/microxt/app"
	"github.com/andewx/microxt/net"
)

const EXIT = 1

func main() {

	//Create a channel to monitor the application status
	app_status := make(chan int)
	main_exit := false

	//Establish a UDP Connection
	tcpConnection := net.NewTCPConnection()
	if tcpConnection != nil {
		defer tcpConnection.Close()
	} else {
		fmt.Printf("Failed to establish UDP Connection")
		return
	}

	{
		go tcpConnection.Listen(app_status)
	}

	{
		app := launchElectronApplication()
		if app != nil {
			go app.RunEventLoop(app_status)
		}
	}

	for !main_exit {
		msg := <-app_status
		if msg == EXIT {
			main_exit = true
		}
	}
}

func launchElectronApplication() *app.ElectronApp {
	//Launch Our Electron Application
	electronApplication, err := app.New()
	if err != nil {
		fmt.Printf("Failed to launch application %s", err)
		return nil
	}
	err = electronApplication.Init() /*Application Blocks*/
	if err != nil {
		fmt.Printf("Failed to launch window %s", err)
		return nil
	}
	return electronApplication
}
