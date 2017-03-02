package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"sync"
	"time"
)

var (
	connections []net.Conn
	connLock    sync.Mutex
)

var listen *string = flag.String("listen", "192.168.1.1:1337", "Specify interface which statstream will listen")
var statInterface *string = flag.String("interface", "enp1s0", "Specify interface which should have its stats collected")

func init() {
	flag.Parse()
	connections = make([]net.Conn, 0, 5)
}

func main() {

	stats := NewStats(*statInterface, 5)

	go acceptConnections()

	for {
		time.Sleep(time.Millisecond * 200)

		stats.findNetStats()
		stats.findFlowStats()
		payload, _ := json.Marshal(stats)
		payload = append(payload, []byte("\n")...)
		connLock.Lock()
		for i, conn := range connections {
			_, err := conn.Write(payload)
			if err != nil {
				log.Printf("Client disconnected: %s: %s", conn.RemoteAddr(), err.Error())
				conn.Close()
				connections = append(connections[:i], connections[i+1:]...)
				if len(connections) == 0 {
					log.Println("No clients left")
				}
			}
		}
		connLock.Unlock()
	}
}

func acceptConnections() {
	l, err := net.Listen("tcp", *listen)

	if err != nil {
		log.Fatalln("Error listening:", err.Error())
	}
	// Close the listener when the application closes.
	defer l.Close()
	log.Println("Listening on", *listen)
	for {

		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {

			log.Fatalln("Error accepting: ", err.Error())

		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	log.Printf("Client connected: %s", conn.RemoteAddr().String())

	connLock.Lock()
	connections = append(connections, conn)
	connLock.Unlock()
}
