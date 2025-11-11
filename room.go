package main

import (
	"sort"
	"sync"
)

type Room struct {
	name    string
	clients map[*Client]bool
	mu      sync.RWMutex
}

func NewRoom(name string) *Room {
	return &Room{
		name:    name,
		clients: make(map[*Client]bool),
	}
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) Join(client *Client) {
	r.mu.Lock()
	r.clients[client] = true
	r.mu.Unlock()
}

func (r *Room) Leave(client *Client) {
	r.mu.Lock()
	delete(r.clients, client)
	r.mu.Unlock()
}

func (r *Room) Broadcast(msg string) {
	r.mu.RLock()
	clients := make([]*Client, 0, len(r.clients))
	for client := range r.clients {
		clients = append(clients, client)
	}
	r.mu.RUnlock()

	for _, client := range clients {
		client.Send(msg)
	}
}

func (r *Room) ListUsers() []string {
	r.mu.RLock()
	clients := make([]*Client, 0, len(r.clients))
	for client := range r.clients {
		clients = append(clients, client)
	}
	r.mu.RUnlock()

	users := make([]string, 0, len(clients))
	for _, client := range clients {
		users = append(users, client.Nick())
	}
	sort.Strings(users)
	return users
}
