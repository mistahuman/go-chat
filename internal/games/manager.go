package games

import (
	"sort"
	"strings"
	"sync"
)

// Session represents an interactive game session bound to a single client.
type Session interface {
	// Name returns the canonical command name used to trigger the game.
	Name() string
	// Handle processes the provided action and optional arguments, returning
	// the messages to display to the client and whether the session has
	// finished as a result of the action.
	Handle(action string, args []string) (messages []string, finished bool, err error)
}

// Bank exposes the chip operations shared across games.
type Bank interface {
	Balance() int
	Credit(amount int)
	Debit(amount int) error
}

// Factory builds new Session instances bound to a shared bank.
type Factory func(bank Bank) Session

// Manager holds the set of available game factories and provides lookup helpers.
type Manager struct {
	mu        sync.RWMutex
	factories map[string]Factory
}

// NewManager returns an initialized Manager.
func NewManager() *Manager {
	return &Manager{
		factories: make(map[string]Factory),
	}
}

// Register associates a command name with a factory.
func (m *Manager) Register(name string, factory Factory) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.factories[strings.ToLower(name)] = factory
}

// Factory returns the factory bound to the provided command name if present.
func (m *Manager) Factory(name string) (Factory, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, ok := m.factories[strings.ToLower(name)]
	return f, ok
}

// List returns the registered command names sorted alphabetically.
func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.factories))
	for name := range m.factories {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
