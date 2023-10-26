package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("connection made")

	buffer := make([]byte, 1024)
	conn.Read(buffer)
	request := strings.Split(string(buffer), " ")
	if request[1] == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		fmt.Println("replied to valid request")
	} else {
		conn.Write([]byte("HTTP/1.1 404 NOT FOUND\r\n\r\n"))
		fmt.Println("replied to invalid request")
	}
}
