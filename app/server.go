package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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

type HTTPRequest struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    string
}

func printRequest(request HTTPRequest) {
	fmt.Println(request.Method + " " + request.Path + " " + request.Version)
	for key, value := range request.Headers {
		fmt.Printf("%s: %s\n", key, value)
	}
	fmt.Println(request.Body)
}

func parseHTTPRequest(requestString string) (*HTTPRequest, error) {
	reader := bufio.NewReader(strings.NewReader(requestString))

	// gets request line
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("have to have this")
	}

	requestLineParts := strings.Fields(requestLine)

	method := requestLineParts[0]
	path := requestLineParts[1]
	version := requestLineParts[2]

	// Reads and parses headers
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break // End of headers
		}
		headerParts := strings.SplitN(line, ": ", 2)

		key := headerParts[0]
		value := strings.TrimSpace(headerParts[1])
		headers[key] = value
	}

	// gets body
	var body string
	bodyBytes, _ := io.ReadAll(reader)
	if len(bodyBytes) > 0 {
		body = string(bodyBytes)
	}

	return &HTTPRequest{
		Method:  method,
		Path:    path,
		Version: version,
		Headers: headers,
		Body:    body,
	}, nil
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

func handleConnection(conn net.Conn) {
	defer conn.Close() // unsure if this works like this but?
	fmt.Println("connection made")

	buffer := make([]byte, 1024)
	conn.Read(buffer)
	request, _ := parseHTTPRequest(string(buffer))
	printRequest(*request)
	if request.Path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if request.Path == "/user-agent" {
		var userAgent string = request.Headers["User-Agent"]
		contentLength := strconv.Itoa(len(userAgent))

		response := HTTPResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type":   "text/plain",
				"Content-Length": contentLength,
			},
			Body: userAgent,
		}
		conn.Write([]byte(stringifyHttpResp(response)))
	} else if strings.HasPrefix(request.Path, "/echo/") {
		val := request.Path[6:len(request.Path)]
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
	} else if strings.HasPrefix(request.Path, "/files/") {
		filename := request.Path[7:len(request.Path)]
		directory := flag.String("directory", "", "Specify the directory")
		flag.Parse() // unsure what does but see it in stack

		f, err := os.Open(*directory + filename)
		if err != nil {
			conn.Write([]byte("HTTP/1.1 404 NOT FOUND\r\n\r\n"))
		}
		defer f.Close()

		val, _ := io.ReadAll(f)
		contentLength := strconv.Itoa(len(val))
		response := HTTPResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type":   "application/octet-stream",
				"Content-Length": contentLength,
			},
			Body: string(val),
		}
		conn.Write([]byte(stringifyHttpResp(response)))
	} else {
		conn.Write([]byte("HTTP/1.1 404 NOT FOUND\r\n\r\n"))
		fmt.Println("replied to invalid request")
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}

}
