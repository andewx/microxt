package main

import (
	"fmt"
	"os"

	"github.com/andewx/microxt/app/application"
	app "github.com/andewx/microxt/app/application"
)

const EXIT = 1
const READY = 11

func main() {

	//Create a channel to monitor the application status
	os.Setenv("GODEBUG", "cgocheck=0")
	main_exit := false

	if backend, err := application.NewApplication(); err != nil {
		fmt.Printf("Failed to create application with error %s", err.Error())
		return
	} else {
		status := backend.GetChannel()
		{
			go app.LaunchElectron(status, backend)
		}

		for !main_exit {
			msg := <-status
			if msg == EXIT {
				main_exit = true
			}
			if msg == READY {
				backend.Init()
			}
		}
	}

}
