package net

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

// We manage a BLE Central Device
type BluetoothConnection struct {
	Conn *bluetooth.Device
	Addr *bluetooth.Address
}

func NewBluetoothConnection() (*BluetoothConnection, error) {
	var conn = new(BluetoothConnection)
	must("Enable BLE Adapter", adapter.Enable())
	fmt.Println("Scanning for devices...")
	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		fmt.Printf("Found %s\n", device.Address.String())
		if device.LocalName() == "MicroXT" {
			fmt.Println("Found the device we are looking for!")
			nConn, err := adapter.Connect(device.Address, bluetooth.ConnectionParams{})
			if err != nil {
				panic(err)
			}
			conn.Conn = nConn
			conn.Addr = &device.Address
		}
	})
	return conn, err
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
