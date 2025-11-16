package chat

import (
	"context"
	"fmt"
	"net"
	"sort"
	"sync"

	"tcp-server/internal/games"
	"tcp-server/internal/games/blackjack"
	"tcp-server/internal/games/roulette"
)

const defaultRoomName = "lobby"

// Server is the main chat coordinator tying together rooms and games.
type Server struct {
	clients map[*Client]bool
	rooms   map[string]*Room
	mu      sync.RWMutex
	games   *games.Manager

	defaultRoom string
}

// NewServer wires the default rooms and registers the built-in games.
func NewServer() *Server {
	gameManager := games.NewManager()
	gameManager.Register("bl", func(bank games.Bank) games.Session {
		return blackjack.NewSession(bank)
	})
	gameManager.Register("roulette", func(bank games.Bank) games.Session {
		return roulette.NewSession(bank)
	})

	s := &Server{
		clients:     make(map[*Client]bool),
		rooms:       make(map[string]*Room),
		games:       gameManager,
		defaultRoom: defaultRoomName,
	}
	s.rooms[s.defaultRoom] = NewRoom(s.defaultRoom)
	return s
}

func (s *Server) HandleConnection(ctx context.Context, conn net.Conn) {
	client := NewClient(ctx, conn, s)

	s.mu.Lock()
	s.clients[client] = true
	s.mu.Unlock()

	lobby := s.GetOrCreateRoom(s.defaultRoom)
	lobby.Join(client)
	client.setRoom(lobby)

	client.Send("Welcome! Commands: /nick /join /rooms /list /msg /leave /who /topic /games /bl /roulette")
	s.Broadcast(s.defaultRoom, fmt.Sprintf("* %s joined", client.Nick()))

	client.Handle()

	s.mu.Lock()
	delete(s.clients, client)
	s.mu.Unlock()

	if room := client.Room(); room != nil {
		room.Leave(client)
		s.Broadcast(room.Name(), fmt.Sprintf("* %s left", client.Nick()))
		client.setRoom(nil)
		s.maybeDeleteRoom(room)
	}
}

// Broadcast sends a message to all members of the given room.
func (s *Server) Broadcast(roomName, msg string) {
	s.mu.RLock()
	room, exists := s.rooms[roomName]
	s.mu.RUnlock()

	if exists {
		room.Broadcast(msg)
	}
}

// GetOrCreateRoom retrieves a room by name or creates a new one.
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

func (s *Server) getRoom(name string) (*Room, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	room, ok := s.rooms[name]
	return room, ok
}

// ListRooms returns the room names sorted alphabetically.
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

// RoomSummaries returns metadata describing each room.
func (s *Server) RoomSummaries() []RoomInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	infos := make([]RoomInfo, 0, len(s.rooms))
	for _, room := range s.rooms {
		infos = append(infos, room.Info())
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})
	return infos
}

// GetRoomInfo returns the metadata for a single room if it exists.
func (s *Server) GetRoomInfo(name string) (RoomInfo, bool) {
	room, ok := s.getRoom(name)
	if !ok {
		return RoomInfo{}, false
	}
	return room.Info(), true
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

func (s *Server) GameFactory(name string) (games.Factory, bool) {
	return s.games.Factory(name)
}

func (s *Server) ListGames() []string {
	return s.games.List()
}

func (s *Server) maybeDeleteRoom(room *Room) {
	if room.Name() == s.defaultRoom {
		return
	}
	if !room.IsEmpty() {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, ok := s.rooms[room.Name()]; ok && existing.IsEmpty() {
		delete(s.rooms, room.Name())
	}
}

// RemoveRoomIfEmpty removes the room from the server if it's empty and not the lobby.
func (s *Server) RemoveRoomIfEmpty(room *Room) {
	s.maybeDeleteRoom(room)
}
