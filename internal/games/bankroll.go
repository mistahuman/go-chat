package games

import (
	"fmt"
	"sync"
)

// Bankroll is a threadsafe chip bank shared across multiple games.
type Bankroll struct {
	mu      sync.Mutex
	balance int
}

// NewBankroll returns a bank initialized with the provided amount.
func NewBankroll(initial int) *Bankroll {
	if initial < 0 {
		initial = 0
	}
	return &Bankroll{balance: initial}
}

// Balance returns the current amount of chips.
func (b *Bankroll) Balance() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.balance
}

// Credit adds chips to the bank.
func (b *Bankroll) Credit(amount int) {
	if amount <= 0 {
		return
	}
	b.mu.Lock()
	b.balance += amount
	b.mu.Unlock()
}

// Debit removes chips from the bank if available.
func (b *Bankroll) Debit(amount int) error {
	if amount <= 0 {
		return fmt.Errorf("invalid debit amount")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if amount > b.balance {
		return fmt.Errorf("Not enough chips. Bankroll: %d", b.balance)
	}
	b.balance -= amount
	return nil
}

var _ Bank = (*Bankroll)(nil)
