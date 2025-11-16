package chat

import (
	"sort"
	"sync"
	"time"
)

// Room models a chat room/table collecting multiple clients.
type Room struct {
	name      string
	topic     string
	createdAt time.Time

	mu      sync.RWMutex
	clients map[*Client]bool
}

// RoomInfo provides read-only metadata about a room.
type RoomInfo struct {
	Name       string
	Topic      string
	CreatedAt  time.Time
	Population int
}

// NewRoom constructs a room with a default topic and metadata.
func NewRoom(name string) *Room {
	return &Room{
		name:      name,
		topic:     "General chat",
		createdAt: time.Now(),
		clients:   make(map[*Client]bool),
	}
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) Topic() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.topic
}

func (r *Room) SetTopic(topic string) {
	r.mu.Lock()
	r.topic = topic
	r.mu.Unlock()
}

func (r *Room) Info() RoomInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return RoomInfo{
		Name:       r.name,
		Topic:      r.topic,
		CreatedAt:  r.createdAt,
		Population: len(r.clients),
	}
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

func (r *Room) Population() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}

func (r *Room) IsEmpty() bool {
	return r.Population() == 0
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
