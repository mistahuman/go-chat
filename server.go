package main

import (
    "context"
    "fmt"
    "net"
    "sync"
)

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
    s.rooms["lobby"] = NewRoom("lobby")
    return s
}

func (s *Server) HandleConnection(ctx context.Context, conn net.Conn) {
    client := NewClient(conn, s)
    
    s.mu.Lock()
    s.clients[client] = true
    s.mu.Unlock()
    
    s.rooms["lobby"].Join(client)
    
    client.Send("Welcome! Commands: /nick /join /rooms /list /msg /leave /who")
    s.Broadcast("lobby", fmt.Sprintf("* %s joined", client.Nick()))
    
    client.Handle(ctx)
    
    s.mu.Lock()
    delete(s.clients, client)
    s.mu.Unlock()
    
    if client.room != nil {
        client.room.Leave(client)
        s.Broadcast(client.room.name, fmt.Sprintf("* %s left", client.Nick()))
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