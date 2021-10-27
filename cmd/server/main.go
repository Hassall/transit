package main

import (
	"encoding/json"
	"net/http"

	"github.com/Hassall/transit/pkg/hub"
	"github.com/Hassall/transit/pkg/request"
	"github.com/Hassall/transit/pkg/store"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{}
var urlStore store.URLStore = &store.DB{}
var urlRequests []request.URLRequest
var connHub = hub.NewHub()

func main() {
	go connHub.Run(func(msg []byte) {
		// TODO
	})

	// connect to url store
	if err := urlStore.Connect(); err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer urlStore.Close()

	// initaliaze urlRequests
	urlRequests = urlStore.RetrieveUrls()

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
		connHub.Broadcast <- url
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	if client := hub.ServeWs(connHub, w, r); client != nil {
		for _, url := range urlRequests {
			client.Send <- url
		}
	}
}
