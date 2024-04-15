package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Listen on TCP port 4221 on all interfaces.
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Connection accepted")

	reader := bufio.NewReader(conn)

	var requestLines []string
	for {
		line, err := reader.ReadString('\n')
		if strings.TrimRight(line, "\r\n") == "" {
			break
		}
		if err != nil {
			fmt.Println("Error reading from connection:", err.Error())
			break
		}
		requestLines = append(requestLines, line)
	}
	hasPath := requestLines[0]
	path := (strings.Split(hasPath, " "))[1]
	resp := ""
	if path == "/" {
		resp = "HTTP/1.1 200 OK\r\n\r\n"
	} else {
		resp = "HTTP/1.1 404 NOT FOUND\r\n\r\n"
	}
	_, err := conn.Write([]byte(resp))
	if err != nil {
		fmt.Println("Error writing response to connection:", err.Error())
	}
}
