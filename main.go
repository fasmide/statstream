package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

func init() {
	connections = make([]net.Conn, 0, 5)
	stats = Stats{}
}

var (
	connections []net.Conn
	connLock    sync.Mutex
)

func main() {
	go acceptConnections()

	for {
		time.Sleep(time.Second)

		findStats()
		payload, _ := json.Marshal(stats)

		connLock.Lock()
		for i, conn := range connections {
			_, err := conn.Write(payload)
			if err != nil {
				log.Printf("Could not write to client %s: %s", conn.RemoteAddr(), err.Error())
				conn.Close()
				connections = append(connections[:i], connections[i+1:]...)
			}
		}
		connLock.Unlock()
	}
}

type Stats struct {
	Device    string
	RxBytes   uint64
	TxBytes   uint64
	RxPackets uint64
	TxPackets uint64
}

var lastNetStat SingleNetStats
var currentNetStat SingleNetStats
var stats Stats

func findStats() error {

	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()

		data := strings.Fields(strings.Trim(text, " "))
		if len(data) != 17 {
			continue
		}

		if strings.HasSuffix(data[0], ":") {

			currentNetStat.ReadArray(data)
		}
		stats.Device = data[0]
	}

	stats.RxBytes = currentNetStat.RxBytes - lastNetStat.RxBytes
	stats.TxBytes = currentNetStat.TxBytes - lastNetStat.TxBytes
	stats.RxPackets = currentNetStat.RxPackets - lastNetStat.RxPackets
	stats.TxPackets = currentNetStat.TxPackets - lastNetStat.TxPackets

	lastNetStat = currentNetStat

	return nil
}

func acceptConnections() {
	l, err := net.Listen("tcp", ":1337")

	if err != nil {
		log.Fatalln("Error listening:", err.Error())
	}
	// Close the listener when the application closes.
	defer l.Close()
	log.Println("Listening on :1337")
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
	log.Printf("New conn from %s", conn.RemoteAddr().String())

	connLock.Lock()
	connections = append(connections, conn)
	connLock.Unlock()
}
