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
	userAgent := "" //stage 5
	if len(requestLines) > 2 {
		userAgent = strings.TrimSpace(strings.Split(requestLines[2], " ")[1])
		fmt.Println(len(userAgent))
	}
	path := (strings.Split(hasPath, " "))[1]
	omitEcho := strings.Split(path, "/echo/") //stage 4
	resp := ""
	if path == "/" || len(omitEcho) > 1 || path == "/user-agent" { //stage 3
		resp = "HTTP/1.1 200 OK\r\n"
	} else {
		resp = "HTTP/1.1 404 NOT FOUND\r\n"
	}
	resp += "Content-Type: text/plain\r\n"
	if path == "/user-agent" {
		resp += fmt.Sprintf("Content-Length: %d\r\n\r\n", len([]byte(userAgent)))
		resp += userAgent
	} else if len(omitEcho) > 1 {
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
