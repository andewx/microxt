package application

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/andewx/microxt/app/templates"
	"github.com/andewx/microxt/common"
	"github.com/andewx/microxt/net"
	"github.com/asticode/go-astilectron"
)

type MainController struct {
	*ControllerBase
}

func NewMainController() *MainController {
	c := &MainController{NewControllerBase()}
	c.handles["ReadDevices"] = c.ReadDevices
	c.handles["SaveDevices"] = c.SaveDevices
	c.handles["SelectMostRecentDevice"] = c.SelectMostRecentDevice
	c.handles["ConnectActiveDevice"] = c.ConnectActiveDevice
	c.handles["OpenProvisioningPage"] = c.OpenProvisioningPage
	c.handles["OpenIDEPage"] = c.OpenIDEPage
	c.handles["SetActiveDevice"] = c.SetActiveDevice
	return c
}

// Implements Functional Components of Application State Controller for Application Requirements
func (m *MainController) ReadDevices(params map[string]string, session *Session, app *Application) error {
	//Reads devices.json file and returns the contents as a string
	filename := common.ProjectRelativePath("microxt/devices.json")
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to read devices.json file %s\n", err.Error())
	}

	//Unmarshal JSON file into the app.Devices
	err = json.Unmarshal(file, &app.Devices)
	if err != nil {
		fmt.Printf("Failed to unmarshal devices.json file %s\n", err.Error())
	}

	fmt.Printf("Read devices.json file successfully\n")
	return nil
}

func (m *MainController) SaveDevices(params map[string]string, session *Session, app *Application) error {
	//Reads devices.json file and returns the contents as a string
	filename := common.ProjectRelativePath("microxt/devices.json")
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("Failed to read devices.json file %s\n", err.Error())
	}

	//Unmarshal JSON file into the app.Devices
	bytes, err_bytes := json.Marshal(app.Devices)
	if err_bytes != nil {
		fmt.Printf("Failed to marshal devices.json file %s\n", err.Error())
	}

	_, err = file.Write(bytes)

	if err != nil {
		fmt.Printf("Failed to write devices.json file %s\n", err.Error())

	}

	fmt.Printf("Read devices.json file successfully\n")
	return nil
}

func (m *MainController) SelectMostRecentDevice(params map[string]string, session *Session, app *Application) error {
	var maxLastUsed int64
	for _, device := range app.Devices {
		if device.LastUsed > maxLastUsed {
			maxLastUsed = device.LastUsed
			app.ActiveDevice = device
		}
	}
	return nil
}

func (m *MainController) ConnectActiveDevice(params map[string]string, session *Session, app *Application) error {
	var err error
	device := app.ActiveDevice
	if device == nil {
		return fmt.Errorf("Device has nil handle, application can't connect\n")
	}
	ip := net.ByteToIP(device.IP)
	port := common.Int16(device.Port, common.LITTLE_ENDIAN)
	app.Connection, err = net.NewTCPConnection(ip.String() + ":" + fmt.Sprintf("%d", port))
	if err != nil {
		app.ActiveDevice.LastUsed = time.Now().Unix()
	}
	return err
}

func (m *MainController) SetActiveDevice(params map[string]string, session *Session, app *Application) error {
	var err error
	device_uuid := params["uuid"]
	device := app.Devices[device_uuid]
	if device == nil {
		return fmt.Errorf("Device has nil handle, application can't connect\n")
	}
	app.ActiveDevice = device
	return err
}

func (m *MainController) OpenProvisioningPage(params map[string]string, session *Session, app *Application) error {
	req := NewRequest("@endpoint", session)
	html := &templates.StringWriter{Str: ""}
	tmpl := app.GetTemplate("Provision")

	if tmpl != nil {
		tmpl.Execute(session.state, html)
	} else {
		fmt.Printf("Failed to find Provision\n")
	}

	req.Extensions["name"] = "@dom"
	req.Extensions["selectors"] = []*DomElement{{HTML: html.Str, Selector: "#root"}}

	//Send session request information to electron
	app.GetElectron().SendMessage(req.JSON(), func(m *astilectron.EventMessage) {})
	return nil

}

func (m *MainController) OpenIDEPage(params map[string]string, session *Session, app *Application) error {
	req := NewRequest("@endpoint", session)
	html := &templates.StringWriter{Str: ""}
	tmpl := app.GetTemplate("Ide")

	if tmpl != nil {
		tmpl.Execute(session.state, html)
	} else {
		fmt.Printf("Failed to find Main template\n")
	}

	req.Extensions["name"] = "@dom"
	req.Extensions["selectors"] = []*DomElement{{HTML: html.Str, Selector: "#root"}}

	//Send session request information to electron
	app.GetElectron().SendMessage(req.JSON(), func(m *astilectron.EventMessage) {})
	return nil
}
