package main

import (
	"context"
	"fmt"
	"net"
	"sort"
	"sync"
)

const defaultRoom = "lobby"

type Server struct {
	clients map[*Client]bool
	rooms   map[string]*Room
	mu      sync.RWMutex
}

func NewServer() *Server {
	s := &Server{
		clients: make(map[*Client]bool),
		rooms:   make(map[string]*Room),
	}
	s.rooms[defaultRoom] = NewRoom(defaultRoom)
	return s
}

func (s *Server) HandleConnection(ctx context.Context, conn net.Conn) {
	client := NewClient(ctx, conn, s)

	s.mu.Lock()
	s.clients[client] = true
	s.mu.Unlock()

	lobby := s.GetOrCreateRoom(defaultRoom)
	lobby.Join(client)
	client.setRoom(lobby)

	client.Send("Welcome! Commands: /nick /join /rooms /list /msg /leave /who /blackjack")
	s.Broadcast(defaultRoom, fmt.Sprintf("* %s joined", client.Nick()))

	client.Handle()

	s.mu.Lock()
	delete(s.clients, client)
	s.mu.Unlock()

	if room := client.Room(); room != nil {
		room.Leave(client)
		s.Broadcast(room.Name(), fmt.Sprintf("* %s left", client.Nick()))
		client.setRoom(nil)
	}
}

func (s *Server) Broadcast(roomName, msg string) {
	s.mu.RLock()
	room, exists := s.rooms[roomName]
	s.mu.RUnlock()

	if exists {
		room.Broadcast(msg)
	}
}

func (s *Server) GetOrCreateRoom(name string) *Room {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, exists := s.rooms[name]
	if !exists {
		room = NewRoom(name)
		s.rooms[name] = room
	}
	return room
}

func (s *Server) ListRooms() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rooms := make([]string, 0, len(s.rooms))
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	sort.Strings(rooms)
	return rooms
}

func (s *Server) FindClientByNick(nick string) *Client {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for client := range s.clients {
		if client.Nick() == nick {
			return client
		}
	}
	return nil
}
