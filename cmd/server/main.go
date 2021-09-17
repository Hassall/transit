package main

import (
	"log"
	"net/http"

	"github.com/Hassall/transit/pkg/request"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("Hello world"))
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		msg := request.URLRequest{URL: "google.com"}

		log.Printf("recv: %s", message)
		err = c.WriteJSON(&msg)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
