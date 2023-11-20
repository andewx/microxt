package net

import (
	"fmt"
	"log"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

const (
	UUID                = "22e0e945-cc6e-4a9a-955f-f8f12350cb14"
	SSID_CHARACTERISTIC = "22e0e946-cc6e-4a9a-955f-f8f12350cb14"
	PASS_CHARACTERISTIC = "22e0e947-cc6e-4a9a-955f-f8f12350cb14"
)

type BluetoothConnection struct {
	local           gatt.Device
	remote          gatt.Peripheral
	services        map[string]*gatt.Service
	characteristics map[string]*gatt.Characteristic
	ssid_write      gatt.WriteHandlerFunc
	pass_write      gatt.WriteHandlerFunc
	Status          chan int
}

func NewBluetoothConnection() (*BluetoothConnection, error) {
	var err error
	var conn = new(BluetoothConnection)
	conn.Status = make(chan int)
	conn.local, err = gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s", err)
	}

	conn.Status <- DEVICE_ON

	conn.services = make(map[string]*gatt.Service)
	conn.characteristics = make(map[string]*gatt.Characteristic)
	conn.local.Handle(gatt.PeripheralDiscovered(func(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
		conn.onPeripheralDiscovered(p, a, rssi)
	}))
	conn.local.Handle(gatt.PeripheralConnected(func(p gatt.Peripheral, err error) {
		conn.onPeripheralConnected(p, err)
	}))
	conn.local.Handle(gatt.PeripheralDisconnected(func(p gatt.Peripheral, err error) {
		conn.onPeripheralDisconnected(p, err)
	}))
	conn.local.Init(func(device gatt.Device, state gatt.State) {
		conn.onStateChanged(device, state)
	})

	return conn, err
}

func (c *BluetoothConnection) Write(uuid string, data []byte) error {
	characteristic := c.characteristics[uuid]
	if characteristic == nil {
		return c.remote.WriteCharacteristic(c.characteristics[uuid], data, false)
	} else {
		return fmt.Errorf("Chracteristic %s not found", uuid)
	}
}

func (c *BluetoothConnection) Close() {
	c.local.CancelConnection(c.remote)
}

func (c *BluetoothConnection) onStateChanged(device gatt.Device, state gatt.State) {
	switch state {
	case gatt.StatePoweredOn:
		fmt.Println("Bluetooth powered on. Scanning for peripherals...")
		device.Scan([]gatt.UUID{gatt.MustParseUUID(UUID)}, false)
	default:
		device.StopScanning()
	}
}

func (c *BluetoothConnection) onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	c.Status <- DEVICE_SCANNING
	if a.LocalName == "DieselX Bluetooth KLD7" {
		fmt.Printf("Peripheral discovered: %s\n", a.LocalName)
		c.remote = p
		p.Device().StopScanning()
		p.Device().Connect(p)
	}
}

func (c *BluetoothConnection) onPeripheralConnected(p gatt.Peripheral, err error) {
	fmt.Printf("Connected to peripheral: %s\n", p.ID())
	// Discover services and characteristics
	services := []gatt.UUID{gatt.MustParseUUID("22e0e945-cc6e-4a9a-955f-f8f12350cb14")}
	discovered, err := p.DiscoverServices(services)
	if err != nil {
		fmt.Printf("Failed to discover services: %s\n", err)
		return
	}
	for _, s := range discovered {
		fmt.Printf("Service discovered: %s\n", s.UUID().String())
		c.services[s.UUID().String()] = s
		cs, cs_err := p.DiscoverCharacteristics(nil, s)
		//Error handle characteristics
		if cs_err != nil {
			err = cs_err
			return
		}
		//Read from the characteristics into the service handler
		for _, ch := range cs {
			fmt.Printf("Characteristic discovered: %s\n", ch.UUID().String())
			c.characteristics[ch.UUID().String()] = ch
			if (ch.Properties() & gatt.CharRead) != 0 {
				fmt.Printf("Characteristic is readable\n")
				b, err := p.ReadCharacteristic(ch)
				if err != nil {
					fmt.Printf("Failed to read characteristic: %s\n", err)
					return
				}
				fmt.Printf("Characteristic value: %s\n", string(b))
			}
		}
	}
	c.Status <- DEVICE_CONNECTED
}

func (c *BluetoothConnection) GetCharacteristicBytes(uuid string) ([]byte, error) {
	characteristic := c.characteristics[uuid]
	if characteristic == nil {
		return nil, fmt.Errorf("Characteristic %s not found", uuid)
	} else {
		return c.remote.ReadCharacteristic(characteristic)
	}
}

func (c *BluetoothConnection) onPeripheralDisconnected(p gatt.Peripheral, err error) {
	c.Status <- DEVICE_DISCONNECTED
	fmt.Printf("Disconnected from peripheral: %s\n", p.ID())
}
