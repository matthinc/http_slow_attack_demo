package main

import (
	"fmt"
	"net"
	"bufio"
	"strings"
	"io/ioutil"
	"strconv"
)

func main() {
	sock, _ := net.Listen("tcp", "0.0.0.0:7000")
	defer sock.Close()

	for {
		connection, err := sock.Accept()
		
		if err != nil {
			fmt.Println(err)
		} else {
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

	fmt.Println("======================")
	fmt.Println(method + " " + resource + " " + http)
	fmt.Println("=== Headers: ")

	// Read additional data (headers, post body...)
	for {
		d, _, err := reader.ReadLine()
		if err != nil || len(d) == 0 {
			break
		}
		fmt.Println(string(d))
	}

	// Routing
	if strings.TrimSpace(method) == "GET" &&  strings.TrimSpace(resource) == "/" {
		fmt.Println("-> Send index")
		sendIndex(connection)
	} else if strings.TrimSpace(method) == "GET" &&  strings.TrimSpace(resource) == "/kitten.webp" {
		fmt.Println("-> Send kitten")
		sendKitten(connection)
	} else if strings.TrimSpace(method) == "GET" &&  strings.TrimSpace(resource) == "/ccwhp.jpg" {
		fmt.Println("-> Send ccwhp")
		sendCCWHP(connection)
	} else {
		fmt.Println("-> Send 404")
		send404(connection)
	}

	connection.Close()
}

func sendIndex(c net.Conn) {
	content := "<body><h1>Test site</h1><br />" +
		"<br /><img width=\"600\" src=\"ccwhp.jpg\"><br /><br />" +
	    "<a href=\"/kitten.webp\">another kitten</a></body>"
	
	c.Write([]byte("HTTP/1.1 200 OK\r\n"))
	c.Write([]byte("Content-Type: text/html\r\n"))
	c.Write([]byte("Content-Length: " + strconv.Itoa(len(content)) + "\r\n"))
	c.Write([]byte("\r\n"))
	c.Write([]byte(content))
}

func sendKitten(c net.Conn) {
	kitten, _ := ioutil.ReadFile("kitten.webp")
	
	c.Write([]byte("HTTP/1.1 200 OK\r\n"))
	c.Write([]byte("Content-Type: image/webp\r\n"))
	c.Write([]byte("Content-Length: " + strconv.Itoa(len(kitten)) + "\r\n"))
	c.Write([]byte("\r\n"))
	c.Write(kitten)
}

func sendCCWHP(c net.Conn) {
	kitten, _ := ioutil.ReadFile("cute_chick_with_hairy_pussy.jpg")
	
	c.Write([]byte("HTTP/1.1 200 OK\r\n"))
	c.Write([]byte("Content-Type: image/jpg\r\n"))
	c.Write([]byte("Content-Length: " + strconv.Itoa(len(kitten)) + "\r\n"))
	c.Write([]byte("\r\n"))
	c.Write(kitten)
}

func send404(c net.Conn) {
	c.Write([]byte("HTTP/1.1 404 Not found\r\n"))
	c.Write([]byte("Content-Type: text/html\r\n"))
	c.Write([]byte("Content-Length: 13\r\n"))
	c.Write([]byte("\r\n"))
	c.Write([]byte("404 Not Found"))
}
