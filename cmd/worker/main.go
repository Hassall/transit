package main

import (
	"log"
	"time"

	"github.com/Hassall/transit/pkg/request"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://server:8080", nil)

	if err != nil {
		log.Fatal("dial:", err)
	}

	defer c.Close()

	go func() {
		for {
			msg := request.URLRequest{}

			err := c.ReadJSON(&msg)
			if err != nil {
				log.Println("read: ", err)
			}
			log.Printf("recv: %s", msg)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}
