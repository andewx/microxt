package models

import "net"

type Device struct {
	UUID     []byte `json:"uuid"`
	Name     string `json:"name"`
	IP       []byte `json:"ip"`
	Port     []byte `json:"port"`
	LastUsed int64  `json:"lastused"`
	_ipv4    net.IP
}
