package websockets

import (
	"context"

	"github.com/filecoin-project/bacalhau/pkg/model"
	"github.com/rs/zerolog/log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case event := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- event:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Event subscribe function
func (h *Hub) JobEventBroadcaster(ctx context.Context, je model.JobEvent) error {
	go func() {
		for {
			b, err := model.JSONMarshalWithMax(model.NewEventEnvelope(je))
			if err != nil {
				log.Error().Msgf("Error marshaling event: %v", err)
				return
			}
			h.broadcast <- b
		}
	}()
	// TODO: #1121 Is there value in returning an error here?
	return nil
}
