package models

import serial "github.com/andewx/microxt/app/serialize"

type Model interface {
	serial.Serialize
	Create() interface{}
	ToString() string
	UpdateKey(string, interface{}) error
	GetKey(string) interface{}
	GetKeys() []string
	GetValues() []interface{}
	GetKeysValues() map[string]interface{}
	GetModelName() string
	GetModelID() string
	DeleteKey(string) error
}
