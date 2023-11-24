package main

import (
	"os"

	"github.com/andewx/microxt/app/application"
	app "github.com/andewx/microxt/app/application"
)

const EXIT = 1

func main() {

	//Create a channel to monitor the application status
	os.Setenv("GODEBUG", "cgocheck=0")
	main_exit := false
	backend, _ := application.NewAirXTApplication()
	status := backend.GetChannel()
	{
		go app.LaunchElectron(status, backend)
	}

	for !main_exit {
		msg := <-status
		if msg == EXIT {
			main_exit = true
		}
	}
}
