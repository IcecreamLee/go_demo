// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "log"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[int]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	broadcastID chan map[string]interface{}

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	ClientIdMax int
}

func newHub() *Hub {
	return &Hub{
		broadcast:   make(chan []byte),
		broadcastID: make(chan map[string]interface{}),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[int]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			log.Printf("client: %+v\n", client)
			h.clients[client.ID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.send)
			}
		case message := <-h.broadcast:
			for i, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, i)
				}
			}
		case mapInfo := <-h.broadcastID:
			id := mapInfo["id"].(int)
			message := mapInfo["msg"].([]byte)
			h.clients[id].send <- message
		}
	}
}
