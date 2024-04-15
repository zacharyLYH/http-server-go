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
	} else {
		fmt.Println("Port 4221 bound")
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
	if path == "/" { //stage 3
		resp = "HTTP/1.1 200 OK\r\n"
	} else {
		resp = "HTTP/1.1 400 NOT FOUND\r\n"
	}
	resp += "Content-Type: text/plain\r\n"
	omitEcho := strings.Split(path, "/echo/") //stage 4
	if len(omitEcho) > 1 {
		fmt.Println(omitEcho[1])
		resp += fmt.Sprintf("Content-Length: %d\r\n\r\n", len([]byte(omitEcho[1])))
		resp += omitEcho[1]
	} else {
		resp += "Content-Length: 0\r\n\r\n"
	}
	_, err := conn.Write([]byte(resp))
	if err != nil {
		fmt.Println("Error writing response to connection:", err.Error())
	}
}
