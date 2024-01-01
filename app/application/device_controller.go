package application

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/andewx/microxt/app/models"
	"github.com/andewx/microxt/common"
	"github.com/andewx/microxt/net"
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

func (p *DeviceController) GetProfile(app *Application) {
	device := app.ActiveDevice
	if device != nil {
		if device.GetTalker() != nil {
			msg_id := common.RandUint32()
			device.GetTalker().SendGetProfile(msg_id)
			m := make(chan int)
			go device.GetTalker().Listen(m, msg_id, device.GetTalker().RecieveDeviceProfile)
			select {
			case msg := <-m:
				if msg == net.GPR_ERROR {
					fmt.Printf("%sError%s, failed to parse message from device: %v\n", net.CS_RED, net.CS_WHITE, "Ack not recieved")
					m <- net.KILL
				} else if msg == net.GPR_RECIEVE {
					//check that our message_id matches.
					m <- net.KILL
				}
			}
		}
	}

	//Fill the device profile details
	device.MaxVel = int64(device.GetTalker().Remote.DeviceProfile.MaxVel)
	//...everything else
}

func (p *DeviceController) DeviceReady(app *Application) {
	device := app.ActiveDevice
	if device != nil {
		if device.IsReady() {
			fmt.Printf("Device is ready\n")
		} else {
			fmt.Printf("Device is not ready\n")
		}
	}
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
