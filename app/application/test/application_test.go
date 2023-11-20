package test

import (
	"fmt"
	"testing"

	test "github.com/andewx/microxt/app/application"
)

const CS_GREEN = "\033[32m"
const CS_WHITE = "\033[37m"

type Foo struct {
	Name    string
	Age     int
	Address string
}

func TestApplicationCreate(t *testing.T) {
	var app test.Application
	var err error
	app, err = test.NewAirXTApplication()
	if err != nil {
		if err.Error() != "NO_DEVICE" {
			t.Errorf("Unexpected error %s", err)
			panic(err)
		}
	}

	m := app.NewSession()

	if app.GetSession(m.UID) == nil {
		t.Errorf("Failed to find user id session")
	} else {
		session := app.GetSession(m.UID)
		fmt.Printf("Session %s\n", session.UID)
	}

	//Find provision template
	prov := app.GetTemplate("Provision")
	if prov == nil {
		t.Errorf("Failed to find provision template")
	}

}
