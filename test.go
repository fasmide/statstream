package main

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

func main() {

	// func (d *Dialer) Dial(urlStr string, requestHeader http.Header) (*Conn, *http.Response, error)
	d := websocket.Dialer{}

	conn, httpRequest, err := d.Dial("ws://10.42.2.236:81/", nil)
	if err != nil {
		log.Printf("Connect error: %s", err)
	}

	log.Printf("http: %+v", httpRequest)

	payload := []byte("/0")
	log.Printf("Sending: %s", payload)
	conn.WriteMessage(websocket.TextMessage, payload)

	payload = []byte("%255")
	log.Printf("Sending: %s", payload)
	conn.WriteMessage(websocket.TextMessage, payload)

	payload = []byte("!FF0000")
	log.Printf("Sending: %s", payload)
	conn.WriteMessage(websocket.TextMessage, payload)

	// payload = []byte("!00FFFFFF")
	// log.Printf("Sending: %s", payload)
	// conn.WriteMessage(websocket.TextMessage, payload)

	var wg sync.WaitGroup
	go func() {
		for {

			msgType, response, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading from leds: %s", err)
				wg.Done()
				break
			}

			log.Printf("Response: %d: %s", msgType, response)
		}
	}()
	wg.Add(1)
	wg.Wait()
	conn.Close()
}
