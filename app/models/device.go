package models

import (
	"fmt"

	"github.com/andewx/microxt/common"
	"github.com/andewx/microxt/net"
)

type Device struct {
	UUID       []byte `json:"uuid"`
	Name       string `json:"name"`
	IP         []byte `json:"ip"`
	Port       []byte `json:"port"`
	LastUsed   int64  `json:"lastused"`
	BitRate    int64  `json:"bitrate"`
	Mode       int64  `json:"mode"`
	MaxVel     int64  `json:"maxvel"`
	MaxRange   int64  `json:"maxrange"`
	Network    string `json:"network"`
	Protocol   string `json:"protocol"`
	Bandwidth  int64  `json:"bandwidth"`
	Tx         int64  `json:"tx"`
	_ipv4      net.IP
	talker     *net.Talker
	connection *net.TCPConnection
}

func (p *Device) SetTalker(talk *net.Talker) {
	p.talker = talk
}

func (p *Device) GetTalker() *net.Talker {
	return p.talker
}

func (p *Device) SetConnection(conn *net.TCPConnection) {
	p.connection = conn
}

func (p *Device) GetPort() int {
	return int(p.Port[0]) + int(p.Port[1])*256
}

const KILL = -1

func (p *Device) IsReady() bool {
	if p.connection != nil && p.talker != nil {
		//send an ack to the device via the talker
		m := make(chan int)
		msg_id := common.RandUint32()
		go p.talker.Listen(m, msg_id, p.talker.RecieveAck)
		p.talker.SendAck(msg_id)
		select {
		case msg := <-m:
			if msg == net.GPR_ERROR {
				fmt.Printf("%sError%s, failed to parse message from device: %v\n", net.CS_RED, net.CS_WHITE, "Ack not recieved")
				m <- KILL
				return false
			} else if msg == net.GPR_RECIEVE {
				m <- KILL
				return true
			}
		}

	}

	return false
}
