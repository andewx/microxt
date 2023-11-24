package net

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

const (
	UUID                = "22e0e945-cc6e-4a9a-955f-f8f12350cb14"
	SSID_CHARACTERISTIC = "22e0e946-cc6e-4a9a-955f-f8f12350cb14"
	PASS_CHARACTERISTIC = "22e0e947-cc6e-4a9a-955f-f8f12350cb14"
	BLE_ON              = 1
	BLE_OFF             = 0
	BLE_SCANNING        = 2
	BLE_CONNECTING      = 3
	BLE_CONNECTED       = 4
	BLE_DISCONNECTED    = 5
	BLE_ERROR           = 6
	BLE_FOUND           = 7
	MAX_BLE_BYTES       = 96
)

var adapter = bluetooth.DefaultAdapter

type BluetoothConnection struct {
	local           *bluetooth.Adapter
	remote          *bluetooth.Device
	services        map[string]*bluetooth.DeviceService
	characteristics map[string]*BLECharacteristic
	available       []bluetooth.ScanResult
	Status          chan int
}

type BLECharacteristic struct {
	characteristic *bluetooth.DeviceCharacteristic
	data           []byte
}

func NewBluetoothConnection() (*BluetoothConnection, error) {
	var err error
	var conn = &BluetoothConnection{}
	conn.Status = make(chan int)
	must("enable BLE stack", adapter.Enable())
	conn.services = make(map[string]*bluetooth.DeviceService)
	conn.characteristics = make(map[string]*BLECharacteristic)
	conn.local = adapter
	return conn, err
}

// Available returns a list of available devices
func (conn *BluetoothConnection) Scan() error {
	var err error
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
	conn.Status <- BLE_SCANNING
	err = adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.AdvertisementPayload.HasServiceUUID(b_uuid) {
			conn.Status <- BLE_FOUND
			conn.available = append(conn.available, result)
			conn.local.StopScan()
			params := bluetooth.ConnectionParams{
				ConnectionTimeout: 200, //2000ms
			}
			conn.Connect(0, params)
		}
	})
	return err
}

func (conn *BluetoothConnection) Close() {
	conn.remote.Disconnect()
}

func (conn *BluetoothConnection) Connect(device int, params bluetooth.ConnectionParams) error {
	var err error
	if device < 0 || device >= len(conn.available) {
		return err
	}
	conn.local.StopScan()
	selected := conn.available[device]
	conn.Status <- BLE_CONNECTING
	conn.remote, err = adapter.Connect(selected.Address, params)
	if err != nil {
		fmt.Printf("Failed to connect to device: %v\n", err)
		return err
	}

	// Discover services and characteristics
	services, err := conn.remote.DiscoverServices(nil)
	if err != nil {
		fmt.Printf("Failed to discover services: %v\n", err)
		return err
	}

	// Iterate over services
	for _, service := range services {

		// Add by uuid to our connection
		conn.services[service.UUID().String()] = &service

		// Discover characteristics for each service
		characteristics, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			fmt.Printf("Failed to discover characteristics: %v\n", err)
			return err
		}

		// Iterate over characteristics
		for _, characteristic := range characteristics {

			// Add by uuid to our connection
			newCharacteristic := &BLECharacteristic{characteristic: &characteristic}
			newCharacteristic.data = make([]byte, MAX_BLE_BYTES)
			conn.characteristics[characteristic.UUID().String()] = newCharacteristic

			// Read data from the characteristic
			size, err := characteristic.Read(newCharacteristic.data)
			if err != nil || size > MAX_BLE_BYTES {
				fmt.Printf("Failed to read data from characteristic: %v\n", err)
				return err
			}
			fmt.Printf("Characteristic value: %s\n", string(newCharacteristic.data))
		}
	}

	conn.Status <- BLE_CONNECTED

	return err
}

func (conn *BluetoothConnection) Write(uuid string, data []byte) error {
	var err error
	characteristic := conn.characteristics[uuid]
	size, err := characteristic.characteristic.WriteWithoutResponse(data)
	if err != nil || size > MAX_BLE_BYTES {
		fmt.Printf("Failed to write data to characteristic: %v\n", err)
		return err
	}
	return err
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
