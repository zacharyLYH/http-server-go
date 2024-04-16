package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	extraArg := ""
	for i := 1; i < len(os.Args)-1; i++ {
		if os.Args[i] == "--directory" {
			extraArg = os.Args[i+1] // Get the next argument after "--directory"
		}
	}
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

		go handleConnection(conn, extraArg) //stage 6. should handle concurrent connections easily
	}
}

func handleConnection(conn net.Conn, extraArg string) {
	defer conn.Close()

	fmt.Println("Connection accepted")

	if extraArg != "" {
		fmt.Println("Passed in argument: ", extraArg)
	}

	reader := bufio.NewReader(conn)

	var requestLines []string
	bodyLength := 0
	for {
		line, err := reader.ReadString('\n')
		trimmedLine := strings.TrimRight(line, "\r\n")
		if trimmedLine == "" {
			break
		}
		fmt.Println("Header: ", line)
		if err != nil {
			fmt.Println("Error reading from connection:", err.Error())
			break
		}
		requestLines = append(requestLines, line)
		if strings.HasPrefix(trimmedLine, "Content-Length:") {
			parts := strings.SplitN(trimmedLine, ":", 2)
			if len(parts) == 2 {
				fmt.Sscanf(parts[1], "%d", &bodyLength)
			}
		}
	}
	hasPath := requestLines[0]
	userAgent := "" //stage 5
	if len(requestLines) > 2 {
		fmt.Println("Request liens length greater than 2: ", requestLines[2])
		userAgent = strings.TrimSpace(strings.Split(requestLines[2], " ")[1])
		fmt.Println(len(userAgent))
	}
	method := (strings.Split(hasPath, " "))[0]
	path := (strings.Split(hasPath, " "))[1]
	omitEcho := strings.Split(path, "/echo/") //stage 4
	var linesRead string
	filesPrefix := strings.Split(path, "/files/") //stage 7
	resp := ""
	if path == "/" || len(omitEcho) > 1 || path == "/user-agent" { //stage 3
		resp = "HTTP/1.1 200 OK\r\n"
		resp += "Content-Type: text/plain\r\n"
	} else if len(filesPrefix) > 1 {
		fmt.Println("File name being requested: ", filesPrefix[1])
		fileName := filesPrefix[1]
		if method == "GET" {
			file, err := os.Open(extraArg + fileName)
			if err != nil {
				resp = "HTTP/1.1 404 NOT FOUND\r\n"
			} else {
				defer file.Close()
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()
					linesRead += line
				}
				resp = "HTTP/1.1 200 OK\r\n"
			}
			resp += "Content-Type: application/octet-stream\r\n"
		} else if method == "POST" {
			fmt.Println("Writing to file: ", extraArg+fileName)
			file, err := os.Create(extraArg + fileName)
			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}
			defer file.Close()
			if bodyLength > 1 {
				data := make([]byte, bodyLength)
				_, err := reader.Read(data)
				if err != nil {
					fmt.Println("Error reading body from connection:", err.Error())
					return
				}
				fmt.Println("Data to write: ", data)
				_, err = file.Write(data)
				if err != nil {
					fmt.Println("Error writing to file:", err)
					return
				}
				resp = "HTTP/1.1 201 CREATED\r\n"
			}
		}
	} else {
		resp = "HTTP/1.1 404 NOT FOUND\r\n"
		resp += "Content-Type: text/plain\r\n"
	}
	//Set content-length and return body
	if path == "/user-agent" {
		resp += fmt.Sprintf("Content-Length: %d\r\n\r\n", len([]byte(userAgent)))
		resp += userAgent
	} else if len(omitEcho) > 1 {
		fmt.Println(omitEcho[1])
		resp += fmt.Sprintf("Content-Length: %d\r\n\r\n", len([]byte(omitEcho[1])))
		resp += omitEcho[1]
	} else if len(filesPrefix) > 1 {
		resp += fmt.Sprintf("Content-Length: %d\r\n\r\n", len([]byte(linesRead)))
		resp += linesRead
	} else {
		resp += "Content-Length: 0\r\n\r\n"
	}
	_, err := conn.Write([]byte(resp))
	if err != nil {
		fmt.Println("Error writing response to connection:", err.Error())
	}
}
