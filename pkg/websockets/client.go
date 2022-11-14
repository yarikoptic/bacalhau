// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://github.com/gorilla/websocket/blob/master/examples/chat/client.go

// Hub maintains the set of active clients and broadcasts messages to the
// clients.

// We are going to just reuse the functionality, but read from the event stream and rebroadcast it to the websocket clients.

package websockets

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case event, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Warn().Err(err).Msg("error while setting write deadline")
			}
			if !ok {
				// The hub closed the channel.
				err = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Warn().Err(err).Msg("error while writing close message")
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Warn().Err(err).Msg("error while getting next writer")
				return
			}
			n, err := w.Write(event)
			log.Trace().Int("bytes", n).Msg("wrote message")
			if err != nil {
				log.Warn().Err(err).Msg("error while writing message")
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// AddClientListener handles websocket requests from the peer.
func AddClientListener(h *Hub, w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg(fmt.Sprintf("w's type is %T\n", w))

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn().Err(err).Msg("error while upgrading websocket connection")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("error: not a websocket connection"))
		if err != nil {
			log.Warn().Err(err).Msg("error while writing error response")
		}
		return
	}
	client := &Client{hub: h, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
}
