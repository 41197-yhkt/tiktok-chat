// Copyright 2017 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file may have been modified by CloudWeGo authors. All CloudWeGo
// Modifications are Copyright 2022 CloudWeGo Authors.

package main

import "github.com/cloudwego/hertz/pkg/common/hlog"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[int64]*Client

	// Inbound messages from the clients.
	broadcast chan C2SMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan C2SMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int64]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.uid] = client
			hlog.Info("uid=", client.uid, " login success")
		case client := <-h.unregister:
			if _, ok := h.clients[client.uid]; ok {
				delete(h.clients, client.uid)
				close(client.send)
			}
		case msg := <-h.broadcast:
			h.chooseClient(msg)
		}
	}
}

func (h Hub) chooseClient(msg C2SMessage) {
	if client, ok := h.clients[msg.To_user_id]; ok {
		select {
		case client.send <- msg:
		default:
			close(client.send)
			delete(h.clients, client.uid)
		}
	} else {
		// 返回用户不存在信息
		sendClient := h.clients[msg.User_id]
		sendClient.send <- C2SMessage{
			User_id:     0,
			To_user_id:  msg.User_id,
			Msg_content: "recv user not login",
		}
	}
}
