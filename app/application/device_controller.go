package application

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/andewx/microxt/app/models"
	"github.com/andewx/microxt/common"
	"github.com/asticode/go-astilectron"
)

type DeviceController struct {
	*ControllerBase
}

func NewDeviceController() *DeviceController {
	c := &DeviceController{NewControllerBase()}
	return c
}

func (p *DeviceController) AddDevice(uuid string, name string, app *Application) *models.Device {
	mDevice := app.Devices[uuid]
	if mDevice == nil {
		//Generate a valid IP/Port Configuration for the device and we need to configure the device over bluetooth
		//IP Address Class C: 192.168.0.0 - 192.168.255.255
		//Port Range: 49152 - 65535
		netIP, port := app.GenUniqueLocation()
		port_bytes := common.GetBytes32(port, common.LITTLE_ENDIAN)
		device := &models.Device{IP: netIP.To4(), Port: port_bytes, Name: name, UUID: []byte(uuid), LastUsed: time.Now().Unix()}
		app.Devices[uuid] = device

		//Set the active device
		app.ActiveDevice = device
		return device
	} else {
		app.ActiveDevice = mDevice
	}

	//Write all devices to the file devices.json
	p.SaveDevices(app)

	return nil
}

func (p *DeviceController) SaveDevices(app *Application) error {
	bytes, err := json.Marshal(app.Devices)
	fn := common.ProjectRelativePath("microxt/devices.json")
	if err == nil {
		file, err_ := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0644)
		defer file.Close()
		if err_ == nil {
			file.Write(bytes)
		} else {
			fmt.Printf("Failed to open devices.json file %s\n", err_.Error())
		}
	} else {
		fmt.Printf("Failed to marshal user %s\n", err.Error())
	}
	return err
}

func (p *DeviceController) GetDevices(params map[string]string, session *Session, app *Application) error {
	bytes, err := json.Marshal(app.Devices)
	if err != nil {
		return err
	}
	req := NewRequest("@dom", session)
	req.Extensions["devices"] = string(bytes)
	app.GetElectron().SendMessage(req.JSON(), func(m *astilectron.EventMessage) {})
	return nil
}
