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
	"time"
	"bytes"
	"sync"
)

type headerField struct {
	key string
	value string
}

type connectionData struct {
	conn *net.Conn
	addr string
	tst int64
}

var (
	connections = make([]connectionData, 0)
	connectionManagementMutex sync.Mutex
)

const (
	maxNumberConnections = 30
	paramU = 20
	paramL = 10
	paramT = 5
)

func main() {
	sock, _ := net.Listen("tcp", "0.0.0.0:7000")
	defer sock.Close()

	go dosProtector()

	for {
		// Wait for new free connections
		for len(connections) >= maxNumberConnections {}

		// Accept new connection
		connection, err := sock.Accept()
		
		if err != nil {
			fmt.Println(err)
		} else {
			connectionData := registerConnection(connection)
			
			// Handle request
			go requestWorker(connectionData)
		}
	}
}

func dosProtector() {
	dosState := false
	
	for {
		if !dosState && len(connections) > paramU {
			dosState = true
			fmt.Println("Slow DOS detected!")
		}
		if dosState && len(connections) < paramL {
			dosState = false
			fmt.Println("Slow DOS ended.")
		}
		if dosState {
			connectionsByIP := make(map[string]int)
			
			// Find number of connections by source IP
			for i := 0; i < len(connections); i++ {
				connection := connections[i]
				
				val, ok := connectionsByIP[connection.addr]
				if ok {
					connectionsByIP[connection.addr] = val + 1
				} else {
					connectionsByIP[connection.addr] = 1
				}
			}

			// Find IP with most connections
			maxNum := 0
			maxIP := ""
			
			for ip, num := range connectionsByIP {
				if num > maxNum {
					maxNum = num
					maxIP = ip
				}
			}

			// Kill all connections with this IP where duration > paramT
			for i := 0; i < len(connections); i++ {
				connection := connections[i]

				if connection.addr == maxIP {
					connectionDuration := time.Now().Unix() - connection.tst

					if connectionDuration > paramT {
						closeConnection(connection)
					}
				}
			}
		}
	}
}

func requestWorker(connection connectionData) {
	defer closeConnection(connection)
	
	reader := bufio.NewReader(*connection.conn)

	// Read first line GET / HTTP/1.1
	method, _ := reader.ReadString(' ')
	resource, _ := reader.ReadString(' ')
	reader.ReadString('\r')

	// Read newline before headers
	reader.ReadLine()

	// Read additional data (headers, post body...)

	contentLength := 0
	contentType := "text/plain"
	
	for {
		d, _, err := reader.ReadLine()
		if err != nil || len(d) == 0 {
			break
		}
		
		headerLine := parseHeaderLine(string(d))

		if headerLine.key == "content-type" {
			contentType = headerLine.value;
		} else if headerLine.key == "content-length" {
			contentLength, _ = strconv.Atoi(headerLine.value)
		}
	}

	method = strings.TrimSpace(method)
	resource = strings.TrimSpace(resource)

	if resource == "/" {
		resource = "index.html"
	}

	if resource == "/connections" && method == "GET" {
		sendConnectionList(*connection.conn)
	}

	if method == "GET" {
		sendFile(*connection.conn, resource)
	}

	if method == "POST" && resource == "/echo" {
		data := receivePostData(reader, contentLength)
		sendEchoReply(*connection.conn, contentType, data)
	}
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

func receivePostData(reader *bufio.Reader, length int) []byte {
	buff := make([]byte, length)

	for i := 0; i < length;  i++ {
		b, err := reader.ReadByte()
		
		if err == nil {
			buff[i] = b
		}
	}

	return buff
}

func sendEchoReply(c net.Conn, contentType string, data []byte) {
	c.Write([]byte("HTTP/1.1 200 OK\r\n"))
	c.Write([]byte("Content-Type: " + contentType + "\r\n"))
	c.Write([]byte("Content-Length: " + strconv.Itoa(len(data)) + "\r\n"))
	c.Write([]byte("\r\n"))
	c.Write(data)
}

func sendConnectionList(c net.Conn) {
	var buffer bytes.Buffer

	buffer.WriteString("<style>td, th { padding: 5px}</style>")
	buffer.WriteString("<h1>Current connections<h1>")
	buffer.WriteString("<table border=\"1\">")
	buffer.WriteString("<tr><th>Connection index</th><th>Source IP</th><th>Duration</th></tr>")

	for i := 0; i < len(connections); i++ {
		buffer.WriteString("<tr>")
		
		buffer.WriteString("<td>")
		buffer.WriteString(strconv.Itoa(i))
		buffer.WriteString("</td>")

		buffer.WriteString("<td>")
		buffer.WriteString(connections[i].addr)
		buffer.WriteString("</td>")

		buffer.WriteString("<td>")
		buffer.WriteString(strconv.Itoa(int(time.Now().Unix() - connections[i].tst)))
		buffer.WriteString("</td>")

		buffer.WriteString("</tr>")
	}

	buffer.WriteString("</table>")
	
	data := buffer.Bytes()
	
	c.Write([]byte("HTTP/1.1 200 OK\r\n"))
	c.Write([]byte("Content-Type: text/html\r\n"))
	c.Write([]byte("Content-Length: " + strconv.Itoa(len(data)) + "\r\n"))
	c.Write([]byte("\r\n"))
	c.Write(data)
}

func parseHeaderLine(header string) *headerField {
	split := strings.Split(header, ": ")
	
	return &headerField{
		key: strings.ToLower(split[0]),
		value: split[1],
	}
}

func closeConnection(connection connectionData) {
	connectionManagementMutex.Lock()
	(*connection.conn).Close()

	for i := 0; i < len(connections); i++ {
		if connection == connections[i] {
			connections = append(connections[:i], connections[i+1:]...)
			break
		}
	}

	fmt.Println("Connections: " + strconv.Itoa(len(connections)))
	connectionManagementMutex.Unlock()
} 

func registerConnection(connection net.Conn) connectionData {
	connectionManagementMutex.Lock()

	cData := connectionData{
		conn : &connection,
		addr : strings.Split(connection.RemoteAddr().String(), ":")[0],
		tst  : time.Now().Unix(),
	}
	
	// Add connection to list
	connections = append(connections, cData)
	
	fmt.Println("Connections: " + strconv.Itoa(len(connections)))
	connectionManagementMutex.Unlock()
	
	return cData
}
