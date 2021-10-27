package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Hassall/transit/pkg/request"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var timedRequests = make(map[string][]request.RequestStatistic)
var requests = make(map[string]bool)
var mu sync.Mutex

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/internal/worker", nil)

	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	go func() {
		ticker := time.NewTicker(5000 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				mu.Lock()
				for url := range requests {
					timedHTTPRequest := request.TimedHTTPRequest(url)
					timedRequests[url] = append(timedRequests[url], timedHTTPRequest)
					fmt.Println(timedHTTPRequest)
				}
				mu.Unlock()
			}
		}
	}()

	for {
		msg := request.URLRequest{}

		err := c.ReadJSON(&msg)
		if err != nil {
			log.Println("read: ", err)
		}
		log.Printf("recv: %s", msg)
		mu.Lock()
		requests[msg.URL] = true
		mu.Unlock()

		err = c.WriteMessage(websocket.TextMessage, []byte("Received Msg."))
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
