package main

import (
	"encoding/json"
	"net/http"

	"github.com/Hassall/transit/pkg/request"
	"github.com/Hassall/transit/pkg/store"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{}
var urlStore store.URLStore = &store.DB{}

func main() {
	// connect to url store
	if err := urlStore.Connect(); err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer urlStore.Close()

	http.HandleFunc("/api/url", urlHandler)
	http.ListenAndServe(":8080", nil)
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		var url request.URLRequest
		// TODO handle malformed requests
		if err := decoder.Decode(&url); err != nil {
			log.Error("Failed to decode JSON", err)
		}
		urls := []request.URLRequest{url}
		// TODO perform batch updates via goroutine
		urlStore.StoreUrls(urls)
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
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
