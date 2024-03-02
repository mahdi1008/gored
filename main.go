package main

import (
	"fmt"
	"net"
)

func main() {
	// Listen on TCP port 6380 on all interfaces.
	listener, err := net.Listen("tcp", ":6380")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Listening on :6380")

	for {
		// Wait for a connection.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	defer conn.Close()

	fmt.Println("handling response")
	for {
		buffer := make([]byte, 1024)
		length, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}
		fmt.Println(string(buffer[:length]))

		// This is where you can customize the response.
		// For example, sending "resp" as a byte array.

		// command := string(buffer[:length])

		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Error writing:", err.Error())
		}
		// Close the connection when you're done with it.

		if err != nil {
			fmt.Println("Error sending response: ", err.Error())
		}
		// time.Sleep(time.Second * 1)
	}
}
