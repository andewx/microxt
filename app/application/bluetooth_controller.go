package application

import (
	"fmt"
	"time"

	"github.com/andewx/microxt/net"
	"github.com/asticode/go-astilectron"
)

type BluetoothController struct {
	*ControllerBase
	Bluetooth *net.BluetoothConnection
	hasStack  bool
}

func NewBluetoothController() *BluetoothController {
	bluetooth_conenction, err := net.NewBluetoothConnection()
	c := &BluetoothController{ControllerBase: NewControllerBase(), Bluetooth: bluetooth_conenction}
	c.handles["BluetoothProvisionWifi"] = c.BluetoothProvisionWifi
	c.handles["BluetoothAddDevice"] = c.BluetoothAddDevice
	c.handles["BluetoothDisconnect"] = c.BluetoothDisconnect
	c.handles["BluetoothEnable"] = c.BluetoothEnable

	if err != nil {
		fmt.Printf("Failed to initiate bluetooth stack on your device %s\n", err.Error())
		c.hasStack = false
	} else {
		c.hasStack = true
	}

	return c
}

func (c *BluetoothController) BluetoothWrite(characteristics map[string][]byte, session *Session) error {
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in bluetooth write %s\n", r)
			err = fmt.Errorf("Recovered from panic in bluetooth write %s\n", r)
		}
	}()

	//Lets rediscover the characteristics from our service
	service := c.Bluetooth.Service(net.UUID)
	if service == nil {
		err = fmt.Errorf("Failed to find service %s\n", net.UUID)
		return err
	}

	//Write from Cache
	for uuid, data := range characteristics {
		err = c.Bluetooth.Write(uuid, data) //Call raises NSInternalConsistencyException
		if err != nil {
			fmt.Printf("Failed to write %s to bluetooth device %s\n", uuid, err.Error())
		}
		time.Sleep(1 * time.Second)
	}

	return err
}

// BluetoothProvisionWifi - Writes the params["ssid"], and params["pass"] strings to a device matching the AirXT Device UUID
func (c *BluetoothController) BluetoothProvisionWifi(params map[string]string, session *Session, app *Application) error {
	connected := false
	ssid := params["ssid"]
	password := params["pass"]
	writeCharacteristics := make(map[string][]byte)
	doneCharacteristic := make(map[string][]byte)

	doneCharacteristic[net.DONE_CHARACTERISTIC] = []byte{1}

	writeCharacteristics[net.SSID_CHARACTERISTIC] = []byte(ssid)
	writeCharacteristics[net.PASS_CHARACTERISTIC] = []byte(password)

	readCharacteristics := make(map[string][]byte)
	readCharacteristics[net.UUID_CHARACTERISTIC] = []byte{}
	readCharacteristics[net.IP_CHARACTERISTIC] = []byte{}
	readCharacteristics[net.PORT_CHARACTERISTIC] = []byte{}
	readCharacteristics[net.SSID_CHARACTERISTIC] = []byte{}
	readCharacteristics[net.PASS_CHARACTERISTIC] = []byte{}

	//Scan and connect to the device
	go c.Bluetooth.ScanUUID(net.UUID)

	//Send to the electron application that we are scanning for a device
	req := NewRequest("@endpoint", session)
	req.Extensions["name"] = "@bluetoothOn"
	app.GetElectron().SendMessage(req.JSON(), func(m *astilectron.EventMessage) {})
	req.Extensions["name"] = "@bluetoothScanning"
	app.GetElectron().SendMessage(req.JSON(), func(m *astilectron.EventMessage) {})

	for !connected {
		if msg := <-c.Bluetooth.Status; msg == net.BLE_CONNECTED {
			req.Extensions["name"] = "@bluetoothConnected"
			app.GetElectron().SendMessage(req.JSON(), func(m *astilectron.EventMessage) {})
			connected = true
			if c.BluetoothRead(readCharacteristics, session) == nil {
				uuid := string(readCharacteristics[net.UUID_CHARACTERISTIC])
				device := app.Devices[uuid]
				if device == nil {
					device := app.Controller("DeviceController").(*DeviceController).AddDevice(uuid, c.Bluetooth.Names[net.UUID], app) //Bluetooth Names is only attached to the GATT UUID
					if device != nil {
						fmt.Printf("Writing IP and Port: %d , %d\n", device.IP, device.Port)
						writeCharacteristics[net.IP_CHARACTERISTIC] = []byte(device.IP)
						writeCharacteristics[net.PORT_CHARACTERISTIC] = []byte(device.Port)
					}
				}
			} else {
				fmt.Printf("Failed to read UUID from device\n")
			}

			if c.BluetoothWrite(writeCharacteristics, session) != nil {
				fmt.Printf("Failed to write characteristics to the device\n")
				return nil
			}

			if c.BluetoothWrite(doneCharacteristic, session) != nil {
				fmt.Printf("Failed to write done characteristic to the device\n")
				return nil
			}
		}
	}

	c.Bluetooth.Close()

	//Set the application to the main page and call
	time.Sleep(5 * time.Second)
	app.ActiveView = MAIN_VIEW
	app.Controller("UtilityController").Endpoint("Scaffold", nil, session, app)

	return nil
}

// Adds a device by provisioning a Unique IP/Port application entry location for the device to use over TCP, to
// do this we need to connect to the device over bluetooth and obtain the unique device ID
func (c *BluetoothController) BluetoothAddDevice(params map[string]string, session *Session, app *Application) error {

	//Read Characteristics
	readCharacteristics := make(map[string][]byte)
	readCharacteristics[net.UUID_CHARACTERISTIC] = []byte{}

	if c.BluetoothRead(readCharacteristics, session) != nil {
		name := c.Bluetooth.Names[net.UUID]
		app.Controller("DeviceController").(*DeviceController).AddDevice(string(readCharacteristics[net.UUID_CHARACTERISTIC]), name, app)
	} else {
		req := NewRequest("@error", session)
		req.Extensions["error"] = "Failed to add device"
		app.GetElectron().SendMessage(req.JSON(), func(m *astilectron.EventMessage) {})
	}

	return nil
}

func (c *BluetoothController) BluetoothRead(characteristics map[string][]byte, session *Session) error {

	//Block until data is receieved
	for uuid := range characteristics {
		data, err := c.Bluetooth.Read(uuid)
		if err != nil {
			err = fmt.Errorf("Failed to read  %s from bluetooth device %s\n", uuid, err.Error())
			fmt.Printf(err.Error())
		} else {
			characteristics[uuid] = data
		}
	}

	return nil
}

func (c *BluetoothController) BluetoothEnable(params map[string]string, session *Session, app *Application) error {
	req := NewRequest("@bluetooth", session)
	req.Extensions["name"] = "@bluetoothEnable"
	app.GetElectron().SendMessage(req.JSON(), func(m *astilectron.EventMessage) {})
	return nil
}

func (c *BluetoothController) BluetoothDisconnect(params map[string]string, session *Session, app *Application) error {
	req := NewRequest("@bluetooth", session)
	req.Extensions["name"] = "@bluetoothDisconnect"
	app.GetElectron().SendMessage(req.JSON(), func(m *astilectron.EventMessage) {})
	return nil
}
