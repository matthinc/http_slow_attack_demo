package main

import (
	"fmt"
	"net"
	"bufio"
	"strings"
	"io/ioutil"
	"strconv"
	"mime"
	"path/filepath"
)

var (
	numConnections = 0
)

func main() {
	sock, _ := net.Listen("tcp", "0.0.0.0:7000")
	defer sock.Close()

	for {
		connection, err := sock.Accept()
		
		if err != nil {
			fmt.Println(err)
		} else {
			numConnections++
			fmt.Println("Connections: " + strconv.Itoa(numConnections))
			go requestWorker(connection)
		}
	}
}

func requestWorker(connection net.Conn) {
	reader := bufio.NewReader(connection)

	// Read first line GET / HTTP/1.1
	method, _ := reader.ReadString(' ')
	resource, _ := reader.ReadString(' ')
	http, _ := reader.ReadString('\r')

	// Read newline before headers
	reader.ReadLine()

	fmt.Println(method + " " + resource + " " + http)

	// Read additional data (headers, post body...)
	for {
		d, _, err := reader.ReadLine()
		if err != nil || len(d) == 0 {
			break
		}
	}

	method = strings.TrimSpace(method)
	resource = strings.TrimSpace(resource)

	if resource == "/" {
		resource = "index.html"
	}

	if method == "GET" {
		sendFile(connection, resource)
	}
	
	connection.Close()
	numConnections--
}

func sendFile(c net.Conn, filename string) {
	file, err := ioutil.ReadFile("./www/" + filename)

	if err == nil {
		mimetype := mime.TypeByExtension(filepath.Ext(filename))
		c.Write([]byte("HTTP/1.1 200 OK\r\n"))
		c.Write([]byte("Content-Type: " + mimetype + "\r\n"))
		c.Write([]byte("Content-Length: " + strconv.Itoa(len(file)) + "\r\n"))
		c.Write([]byte("\r\n"))
		c.Write(file)
	} else {
		c.Write([]byte("HTTP/1.1 404 Not found\r\n"))
		c.Write([]byte("Content-Type: text/html\r\n"))
		c.Write([]byte("Content-Length: 13\r\n"))
		c.Write([]byte("\r\n"))
		c.Write([]byte("404 Not Found"))
	}
}

