package net

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

const (
	UUID                = "22e0e945-cc6e-4a9a-955f-f8f12350cb14"
	SSID_CHARACTERISTIC = "22e0e946-cc6e-4a9a-955f-f8f12350cb14"
	PASS_CHARACTERISTIC = "22e0e947-cc6e-4a9a-955f-f8f12350cb14"
	IP_CHARACTERISTIC   = "22e0e948-cc6e-4a9a-955f-f8f12350cb14"
	PORT_CHARACTERISTIC = "22e0e949-cc6e-4a9a-955f-f8f12350cb14"
	UUID_CHARACTERISTIC = "22e0e950-cc6e-4a9a-955f-f8f12350cb14"
	DONE_CHARACTERISTIC = "22e0e951-cc6e-4a9a-955f-f8f12350cb14"

	BLE_ON           = 1
	BLE_OFF          = 0
	BLE_SCANNING     = 2
	BLE_CONNECTING   = 3
	BLE_CONNECTED    = 4
	BLE_DISCONNECTED = 5
	BLE_ERROR        = 6
	BLE_FOUND        = 7
	BLE_SUCCESS      = 8
	MAX_BLE_BYTES    = 96
)

var adapter = bluetooth.DefaultAdapter

type BluetoothConnection struct {
	Names           map[string]string
	local           *bluetooth.Adapter
	remote          *bluetooth.Device
	services        map[string]*bluetooth.DeviceService
	characteristics map[string]*BLECharacteristic
	available       []bluetooth.ScanResult
	Status          chan int
}

type BLECharacteristic struct {
	attribute *bluetooth.DeviceCharacteristic
	data      []byte
}

func NewBluetoothConnection() (*BluetoothConnection, error) {
	var err error
	var conn = &BluetoothConnection{}
	conn.Status = make(chan int)
	conn.services = make(map[string]*bluetooth.DeviceService)
	conn.characteristics = make(map[string]*BLECharacteristic)
	conn.Names = make(map[string]string)
	err = adapter.Enable()
	conn.local = adapter
	return conn, err
}

func (conn *BluetoothConnection) Service(uuid string) *bluetooth.DeviceService {
	return conn.services[uuid]
}

// Available returns a list of available devices
func (conn *BluetoothConnection) Scan() error {
	var err error

	conn.Status <- BLE_ON

	conn.Status <- BLE_SCANNING
	err = adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		conn.available = append(conn.available, result)
	})
	return err
}

func (conn *BluetoothConnection) ScanUUID(uuid string) error {
	var err error
	b_uuid, err := bluetooth.ParseUUID(uuid)
	must("parse UUID", err)

	ch := make(chan bluetooth.ScanResult, 1)

	conn.Status <- BLE_ON
	conn.Status <- BLE_SCANNING
	fmt.Printf("Scanning...\n")
	go adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.AdvertisementPayload.HasServiceUUID(b_uuid) {
			if result.LocalName() == "DieselX Bluetooth KLD7" {
				conn.Names[UUID] = result.LocalName()
				fmt.Printf("Found device %s\n", result.LocalName())
				conn.available = append(conn.available, result)
				ch <- result
			}
		}
	})

	fmt.Printf("Connecting...\n")
	select {
	case result := <-ch:
		fmt.Printf("Connecting to %s...\n", result.Address.String())
		err = conn.Connect(0, bluetooth.ConnectionParams{})
		if err != nil {
			fmt.Printf("Error connecting to %s device: %s\n", result.LocalName(), err)
			conn.Status <- BLE_ERROR
			return err
		}
	}
	return err
}

func (conn *BluetoothConnection) Close() {
	conn.local.StopScan()
	if conn.remote != nil {
		conn.remote.Disconnect()
	}
}

func (conn *BluetoothConnection) Connect(device int, params bluetooth.ConnectionParams) error {
	var err error
	if device < 0 || device >= len(conn.available) {
		return err
	}
	selected := conn.available[device]

	fmt.Printf("Connecting to device %s with address %s...\n", selected.LocalName(), selected.Address.String())
	conn.Status <- BLE_CONNECTING
	conn.remote, err = adapter.Connect(selected.Address, params)
	if err != nil {
		fmt.Printf("Failed to connect to device: %v\n", err)
		return err
	}

	fmt.Print("Connected to device discovering services...\n")

	// Discover services and characteristics
	services, err := conn.remote.DiscoverServices(nil)
	if err != nil {
		fmt.Printf("Failed to discover services: %v\n", err)
		return err
	}

	if len(services) == 0 {
		fmt.Printf("No services found\n")
		return err
	}

	// Iterate over services
	for _, service := range services {

		// Add by uuid to our connection
		conn.services[service.UUID().String()] = &service
		fmt.Printf("Service: %s\n", service.UUID().String())

		// Discover characteristics for each service
		characteristics, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			fmt.Printf("Failed to discover characteristics: %v\n", err)
			return err
		}

		if len(characteristics) == 0 {
			fmt.Printf("No characteristics found for Service UUID %s \n", service.UUID().String())
		} else {
			fmt.Printf("Found %d characteristics\n", len(characteristics))
		}

		//Rewrite loop as not a range loop
		for i := 0; i < len(characteristics); i++ {
			characteristic := characteristics[i]
			uuid := characteristic.UUID().String()

			bytes := 64
			if uuid == PORT_CHARACTERISTIC {
				bytes = 2
			} else if uuid == UUID_CHARACTERISTIC {
				bytes = 8
			} else if uuid == IP_CHARACTERISTIC {
				bytes = 4
			} else if uuid == DONE_CHARACTERISTIC {
				bytes = 1
			}
			newCharacteristic := &BLECharacteristic{attribute: &characteristic}
			newCharacteristic.data = make([]byte, bytes)
			characteristic.Read(newCharacteristic.data)
			conn.characteristics[uuid] = newCharacteristic
			fmt.Printf("Characteristic: %s\n", uuid)
			fmt.Println(string(newCharacteristic.data))

		}
	}

	conn.Status <- BLE_CONNECTED

	return err
}

func (conn *BluetoothConnection) Write(uuid string, data []byte) error {
	var err error
	characteristic := conn.characteristics[uuid]

	if data == nil || len(data) == 0 {
		return fmt.Errorf("Failed to write data to characteristic, data is nil or empty\n")
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in bluetooth write %s\n", r)
			err = fmt.Errorf("Recovered from panic in bluetooth write %s\n", r)
		}
	}()

	if characteristic != nil {
		fmt.Printf("Writing to characteristic %s\n%s\n", uuid, string(data))
		size, err := characteristic.attribute.WriteWithoutResponse(data)
		if err != nil || size == 0 {
			fmt.Printf("Failed to write data to characteristic: %v\n", err)
			return err
		}
	} else {
		err = fmt.Errorf("Failed to find characteristic %s\n", uuid)
	}

	return err
}

func (conn *BluetoothConnection) Read(uuid string) ([]byte, error) {
	var err error
	characteristic := conn.characteristics[uuid]
	if characteristic != nil {
		_, err = characteristic.attribute.Read(characteristic.data)
		return characteristic.data, err
	} else {
		err = fmt.Errorf("Failed to find characteristic %s\n", uuid)
	}
	return []byte{}, err
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
