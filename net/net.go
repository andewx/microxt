package net

import (
	"fmt"
	"net"

	"github.com/andewx/microxt/common"
)

const APPLICATION_PORT = 9060
const DEFAULT_BAUD = 9600
const CS_RED = "\033[31m"
const CS_GREEN = "\033[32m"
const CS_YELLOW = "\033[33m"
const CS_BLUE = "\033[34m"
const CS_WHITE = "\033[37m"

const EXIT = 1
const IDLE = 0
const HANDSHAKE = 2

const (
	NO_DEVICE = iota
	DEVICE_CONNECTED
	DEVICE_SCANNING
	DEVICE_ON
	DEVICE_DISCONNECTED
)

// TCPConnection is a struct that contains the TCP connection, address, and message
type TCPConnection struct {
	Conn   *net.TCPConn
	Addr   *net.TCPAddr
	MyAddr *net.TCPAddr
	Talker *Talker
	Status int
}

type IP struct {
	IP net.IP
}

func (ip *IP) String() string {
	return ip.IP.String()
}

func (ip *IP) To4() []byte {
	return ip.IP.To4()
}

func (ip *IP) To16() []byte {
	return ip.IP.To16()
}

func Int32ToIP(in uint32) net.IP {
	bytes := common.GetBytes32(int(in), common.BIG_ENDIAN)
	return net.IPv4(bytes[0], bytes[1], bytes[2], bytes[3])
}

func ByteToIP(bytes []byte) net.IP {
	return net.IPv4(bytes[0], bytes[1], bytes[2], bytes[3])
}

func IPV4Address() *net.TCPAddr {

	var localAddr *net.TCPAddr

	// Get a list of network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("%sError%s getting network interfaces: %v\n", CS_RED, CS_WHITE, err)
		return nil
	}

	// Loop through the network interfaces
	for _, iface := range interfaces {
		// Skip interfaces that are down or not connected
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Get a list of addresses for this interface
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Printf("%sError%s getting addresses for interface %v: %v\n", CS_RED, CS_WHITE, iface.Name, err)
			continue
		}

		// Loop through the addresses
		for _, addr := range addrs {
			// Check if this is an IP address
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					// This is an IPv4 address
					fmt.Printf("%sLocal IP address:%s %v\n", CS_GREEN, CS_WHITE, ipnet.IP.String())
					localAddr = &net.TCPAddr{IP: ipnet.IP, Port: APPLICATION_PORT}
				}
			}
		}
	}
	return localAddr
}

// TCPBroadcast is a function that sends a TCP broadcast message, we send a specialized message for connection
func NewTCPConnection(addressport string) (*TCPConnection, error) {
	// Create New TCP Connection
	TCPConnection := new(TCPConnection)
	TCPConnection.Talker = NewTalker()
	TCPConnection.Status = HANDSHAKE

	//UDP Connection For Broadcast to obtain device IP before we switch over to TCP communication

	// Get the local IP address
	if TCPConnection.MyAddr = IPV4Address(); TCPConnection.MyAddr == nil {
		fmt.Printf("%sError%s getting local IP address\n", CS_RED, CS_WHITE)
	}

	// Create a TCP connection we initiate request across the network
	var conn net.Conn
	var err error
	var waitTime int
	for conn, err = net.Dial("tcp", addressport); err != nil; {
		//	time.Sleep(1 * time.Second)
		waitTime += 1
		if waitTime > TIMEOUT {
			return nil, fmt.Errorf("%sError%s connecting with TCP address: %v\n", CS_RED, CS_WHITE, err)
		}
	}

	TCPConnection.Conn = conn.(*net.TCPConn)
	TCPConnection.Status = IDLE

	return TCPConnection, nil
}

func (u *TCPConnection) Close() {
	u.Conn.Close()
}

func (u *TCPConnection) Listen(status chan int) {
	finished := false
	buffer := make([]byte, 1024*4)
	u.Talker.LocalStatus = IDLE
	u.Talker.DeviceStatus = IDLE
	u.Talker.ConversationID = "0"
	u.Talker.Inbox = Messages{}
	fmt.Printf("%sListening%s for TCP messages...\n", CS_YELLOW, CS_WHITE)
	for !finished {
		// Read from the TCP connection
		_, err := u.Conn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Timed out waiting for data")
			} else {
				fmt.Printf("Error reading from TCP connection: %v\n", err)
			}
		} else {
			// Print the TCP message
			fmt.Printf("%sReceived TCP broadcast: %s\n", CS_GREEN, buffer[:len(buffer)-1])
			err = u.Talker.Receive(buffer)

			if err != nil {
				fmt.Printf("%sError%s, failed to parse message from device: %v\n", CS_RED, CS_WHITE, err)
			} else {
				fmt.Printf("%sMessage%s, parsed message from device:\n", CS_GREEN, CS_WHITE)
			}
			//Interrogate and route commands etc , handle the protobuf message here
		}
	}
}

func (u *TCPConnection) Send(bytes []byte) (int, error) {
	return u.Conn.Write(bytes)
}
