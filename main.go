package main

import (
	"fmt"

	"github.com/andewx/microxt/app"
	"github.com/andewx/microxt/app/application"
	"github.com/andewx/microxt/net"
)

const EXIT = 1

func main() {

	//Create a channel to monitor the application status
	app_status := make(chan int)
	main_exit := false
	backend, _ := application.NewAirXTApplication()
	//Establish a TCP Connection
	tcpConnection := net.NewTCPConnection("192.168.0.10:9060")
	if tcpConnection != nil {
		defer tcpConnection.Close()
	} else {
		fmt.Printf("Failed to establish TCP Connection")
		return
	}

	//Send SSID/Password as raw byte messages to the connection in successio
	var err error

	_, err = tcpConnection.Send([]byte("ObamaPhone"))

	if err != nil {
		fmt.Printf("%sError$%s, failed to send SSID/Password to device", net.CS_RED, net.CS_WHITE)
		return
	}

	_, err = tcpConnection.Send([]byte("Ijoyflqt9"))

	if err != nil {
		fmt.Printf("%sError$%s, failed to send SSID/Password to device", net.CS_RED, net.CS_WHITE)
		return
	}

	{
		app := launchElectronApplication(backend)
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

func launchElectronApplication(backend application.Application) *app.ElectronApp {
	//Launch Our Electron Application
	electronApplication, err := app.New(backend)
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
