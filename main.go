package main

import (
	"fmt"

	"github.com/andewx/microxt/app"
)

func main() {

	electronApplication, err := app.New()

	if err != nil {
		fmt.Printf("Failed to launch application %s", err)
		return

	}

	err = electronApplication.Init() /*Application Blocks*/

	if err != nil {
		fmt.Printf("Failed to launch window %s", err)
		return
	}
}
