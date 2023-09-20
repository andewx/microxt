package net

import (
	"fmt"
	"net"
)

func init() {
	// Define the broadcast address and port
	broadcastAddress := "192.168.1.255" // Replace with your network's broadcast address
	port := 9060                        // Replace with your chosen UDP port

	// Resolve the UDP address
	udpAddress := fmt.Sprintf("%s:%d", broadcastAddress, port)

	// Create a UDP connection
	conn, err := net.Dial("udp", udpAddress)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	// Compose and send the UDP broadcast message
	message := fmt.Sprintf("@handshake:[ip|byte|4]%b", localAddr.IP.To4())
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}

	fmt.Printf("Sent UDP broadcast message: %s\n", message)
}
