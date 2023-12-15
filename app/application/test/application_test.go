package test

import (
	"fmt"
	"testing"

	test "github.com/andewx/microxt/app/application"
	"github.com/andewx/microxt/app/templates"
)

const CS_GREEN = "\033[32m"
const CS_WHITE = "\033[37m"

type Foo struct {
	Name    string
	Age     int
	Address string
}

func TestApplicationCreate(t *testing.T) {
	var app *test.Application
	var err error
	app, err = test.NewApplication()
	if err != nil {
		t.Errorf("Failed to create application with error %s", err.Error())
		return
	}

	m := app.NewSession()

	if app.GetSession(m.UID) == nil {
		t.Errorf("Failed to find user id session")
	} else {
		session := app.GetSession(m.UID)
		fmt.Printf("Session %s\n", session.UID)
	}

	//Find provision template
	prov := app.GetTemplate("REGISTER_VIEW")
	if prov == nil {
		t.Errorf("Failed to find provision template")
	}

	html := &templates.StringWriter{Str: ""}
	err = prov.Execute(m, html)

	if err != nil {
		t.Errorf("Failed to execute provision template with error %s", err.Error())
	}

	//Send Credentials over Bluetooth
	params := make(map[string]string)
	params["ssid"] = "test"
	params["password"] = "test"
	//app.Controller("UtilityController").Endpoint("Provision", params, m, app)

}
