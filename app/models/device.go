package models

type Device struct {
	DeviceID     string
	DeviceName   string
	DeviceStatus int
	AddressPort  string
	Driver       Driver
}

type Driver struct {
	DriverID     string
	DriverName   string
	DriverStatus int
}
