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
var urlRequests []request.URLRequest

var conns = make(map[*websocket.Conn]bool)
var broadcast = make(chan []request.URLRequest)

func main() {
	// connect to url store
	if err := urlStore.Connect(); err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer urlStore.Close()

	// initaliaze urlRequests
	urlRequests = urlStore.RetrieveUrls()

	go func() {
		for {
			select {
			case urls := <-broadcast:
				for client := range conns {
					for _, url := range urls {
						err := client.WriteJSON(&url)
						if err != nil {
							log.Error("Failed to write to client", err)
							delete(conns, client)
							break
						}
					}
				}
			}
		}
	}()

	http.HandleFunc("/api/url", urlHandler)
	http.HandleFunc("/internal/worker", websocketHandler)
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
		// update in memory urls
		urlRequests = append(urlRequests, url)
		// notify clients of new url
		broadcast <- urls
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Failed to establish websocket connection", err)
		return
	}

	// TODO establish 2 way connect
	// gracefully shutdown client
	conns[c] = true

	for _, url := range urlRequests {
		err := c.WriteJSON(&url)
		if err != nil {
			log.Error("Failed to write to client", err)
		}
	}
}
