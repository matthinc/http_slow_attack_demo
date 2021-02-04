package main

import (
	"fmt"
	"net"
	"bufio"
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

	fmt.Println(method + " " + resource + " " + http)

	// Read additional data (headers, post body...)
	for d, err := reader.ReadString('\r'); err == nil && d != ""; {
		fmt.Println("Received: " + d)
	}
}
