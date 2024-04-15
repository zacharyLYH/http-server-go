package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Listen on TCP port 4221 on all interfaces.
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	// Accept new connections.
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connection accepted")

	// Create a new reader from the connection
	reader := bufio.NewReader(conn)

	// Read data from the connection
	// You can specify your own buffer size or use bufio for simplicity
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err.Error())
			break
		}
		fmt.Print("Received: ", line)
		resp := "HTTP/1.1 200 OK\r\n\r\n"
		conn.Write([]byte(resp))
		if err != nil {
			fmt.Println("Error writing response to connection:", err.Error())
			break
		}
	}
}
