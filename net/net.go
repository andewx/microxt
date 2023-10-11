package net

import (
	"fmt"
	"net"
	"time"
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

// TCPConnection is a struct that contains the TCP connection, address, and message
type TCPConnection struct {
	Conn   *net.TCPConn
	Addr   *net.TCPAddr
	MyAddr *net.TCPAddr
	Talker []*Talker
	Status int
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
func NewTCPConnection() *TCPConnection {
	// Create New TCP Connection
	TCPConnection := new(TCPConnection)
	TCPConnection.Talker = make([]*Talker, 1)
	TCPConnection.Talker[0] = NewTalker()
	TCPConnection.Status = HANDSHAKE

	//UDP Connection For Broadcast to obtain device IP before we switch over to TCP communication

	// Get the local IP address
	if TCPConnection.MyAddr = IPV4Address(); TCPConnection.MyAddr == nil {
		fmt.Printf("%sError%s getting local IP address\n", CS_RED, CS_WHITE)
	}

	// Define the broadcast address and port
	broadcastAddress := "192.168.1.255" // Replace with your network's broadcast address
	port := 9060                        // Replace with your chosen TCP port

	// Resolve the TCP address
	BAddress := fmt.Sprintf("%s:%d", broadcastAddress, port)

	// Create a UDP connection we initiate request across the network
	var conn net.Conn
	var err error
	var waitTime int
	for conn, err = net.Dial("UDP", BAddress); err != nil; {
		time.Sleep(1 * time.Second)
		waitTime += 1
		if waitTime > TIMEOUT {
			fmt.Printf("%sError%s connecting with UDP broadcast address: %v\n", CS_RED, CS_WHITE, err)
			return nil
		}
	}

	// We only communicate with the local device
	localAddr := conn.LocalAddr().(*net.TCPAddr)
	TCPConnection.Addr = localAddr

	// Compose and send the UDP message to switch over to TCP
	message := "SWTCH"
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return nil
	}

	fmt.Printf("Sent UDP broadcast message: %s\n", message)

	// Delay and allow network to switch over to TCP
	time.Sleep(1 * time.Second)

	// Form TCP connection
	// Create a UDP connection we initiate request across the network
	var tcp_conn net.Conn
	var tcp_err error
	var tcp_waitTime int
	for tcp_conn, tcp_err = net.Dial("TCP", TCPConnection.Addr.String()); tcp_err != nil; {
		time.Sleep(1 * time.Second)
		tcp_waitTime += 1
		if tcp_waitTime > TIMEOUT {
			fmt.Printf("%sError%s connecting with TCP local address: %s\n", CS_RED, CS_WHITE, TCPConnection.Addr.String())
			return nil
		}
	}

	TCPConnection.Conn = tcp_conn.(*net.TCPConn)
	TCPConnection.Status = IDLE

	return TCPConnection
}

func (u *TCPConnection) Close() {
	u.Conn.Close()
}

func (u *TCPConnection) Listen(status chan int) {
	finished := false
	fmt.Printf("%sListening%s for TCP messages...\n", CS_YELLOW, CS_WHITE)
	for !finished {
		if msg := <-status; msg == EXIT {
			finished = true
		} else {
			// Read from the TCP connection
			buffer := make([]byte, 1024)
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

			}

		}

	}
}
