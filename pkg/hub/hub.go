package hub

import (
	"github.com/Hassall/transit/pkg/request"
	log "github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Broadcast message to workers.
	Broadcast chan request.URLRequest

	// Receive message from workers.
	Receive chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

// NewHub returns a new hub
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan request.URLRequest),
		Receive:    make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run runs rub
func (h *Hub) Run(receiveHandler func([]byte)) {
	for {
		select {
		case client := <-h.Register:
			log.Debug("Registering client.")
			h.Clients[client] = true
		case client := <-h.Unregister:
			log.Debug("Unregistering client.")
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			log.Debug("Broadcasting message.")
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		case message := <-h.Receive:
			receiveHandler(message)
		}
	}
}
