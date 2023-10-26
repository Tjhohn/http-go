package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}

func statusCodeToText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 404:
		return "NOT FOUND"
	default:
		return "NOT FOUND"
	}

}

func stringifyHttpResp(resp HTTPResponse) string {
	statusline := "HTTP/1.1 " + strconv.Itoa(resp.StatusCode) + "  " + statusCodeToText(resp.StatusCode) + "\r\n"
	var headers string
	for key, value := range resp.Headers {
		headers += key + ": " + value + "\r\n"
	}
	headers += "\r\n"

	return statusline + headers + resp.Body
}

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
	if strings.HasPrefix(request[1], "/") {
		if strings.HasPrefix(request[1], "/echo/") {
			val := request[1][6:len(request[1])]
			contentLength := strconv.Itoa(len(val))
			fmt.Println(val)
			response := HTTPResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"Content-Type":   "text/plain",
					"Content-Length": contentLength,
				},
				Body: val,
			}
			conn.Write([]byte(stringifyHttpResp(response)))
		} else {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		}

		fmt.Println("replied to valid request")
	} else {
		conn.Write([]byte("HTTP/1.1 404 NOT FOUND\r\n\r\n"))
		fmt.Println("replied to invalid request")
	}
}
