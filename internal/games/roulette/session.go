package roulette

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"tcp-server/internal/games"
)

type randomizer interface {
	Intn(n int) int
}

// Session models a simple roulette table bound to a single client.
type Session struct {
	bank       games.Bank
	minBet     int
	maxBet     int
	currentBet int

	selectionKind  string
	selectionColor string
	selectionValue int

	rng randomizer
}

const (
	kindNone   = ""
	kindColor  = "color"
	kindNumber = "number"
)

// NewSession returns a roulette session with a seeded RNG.
func NewSession(bank games.Bank) *Session {
	if bank == nil {
		panic("roulette session requires bank")
	}
	return newSessionWithRand(bank, rand.New(rand.NewSource(time.Now().UnixNano())))
}

func newSessionWithRand(bank games.Bank, r randomizer) *Session {
	return &Session{
		bank:   bank,
		minBet: 5,
		maxBet: 200,
		rng:    r,
	}
}

// Name satisfies the games.Session interface.
func (s *Session) Name() string {
	return "roulette"
}

// Handle interprets roulette commands.
func (s *Session) Handle(action string, args []string) ([]string, bool, error) {
	switch strings.ToLower(action) {
	case "start":
		s.selectionKind = kindNone
		return []string{
			"Roulette table ready.",
			"Place a bet with /roulette wager <amount> <red|black|green|0-36>.",
			s.statusMessage(),
		}, false, nil
	case "wager":
		if len(args) < 2 {
			return nil, false, fmt.Errorf("usage: /roulette wager <amount> <choice>")
		}
		amount, err := strconv.Atoi(args[0])
		if err != nil || amount <= 0 {
			return nil, false, fmt.Errorf("invalid bet amount")
		}
		if amount < s.minBet || amount > s.maxBet {
			return nil, false, fmt.Errorf("bet must be between %d and %d", s.minBet, s.maxBet)
		}
		if s.bank.Balance() < amount {
			return nil, false, fmt.Errorf("Not enough chips. Bankroll: %d", s.bank.Balance())
		}
		choice := strings.ToLower(strings.Join(args[1:], ""))
		if err := s.setSelection(choice); err != nil {
			return nil, false, err
		}
		s.currentBet = amount
		return []string{fmt.Sprintf("Betting %d chips on %s.", s.currentBet, s.describeSelection())}, false, nil
	case "spin":
		if s.selectionKind == kindNone || s.currentBet == 0 {
			return nil, false, fmt.Errorf("place a bet first with /roulette wager <amount> <choice>")
		}
		if err := s.bank.Debit(s.currentBet); err != nil {
			return nil, false, err
		}
		outcome := s.rng.Intn(37)
		color := colorFor(outcome)
		won, multiplier := s.evaluateBet(outcome, color)
		messages := []string{fmt.Sprintf("Wheel landed on %d (%s).", outcome, color)}
		if won {
			winnings := s.currentBet * multiplier
			s.bank.Credit(s.currentBet + winnings)
			messages = append(messages, fmt.Sprintf("You won %d chips!", winnings))
		} else {
			messages = append(messages, fmt.Sprintf("You lost %d chips.", s.currentBet))
		}
		messages = append(messages, s.statusMessage())
		s.selectionKind = kindNone
		s.currentBet = 0
		return messages, false, nil
	case "status":
		return []string{s.statusMessage()}, false, nil
	case "reset":
		s.currentBet = 0
		s.selectionKind = kindNone
		return []string{"Roulette session reset.", s.statusMessage()}, false, nil
	case "quit":
		return []string{"Roulette session closed."}, true, nil
	default:
		return nil, false, fmt.Errorf("unknown roulette action: %s", action)
	}
}

func (s *Session) setSelection(choice string) error {
	switch choice {
	case "red", "black", "green":
		s.selectionKind = kindColor
		s.selectionColor = choice
		return nil
	default:
		value, err := strconv.Atoi(choice)
		if err != nil || value < 0 || value > 36 {
			return fmt.Errorf("invalid bet choice: %s", choice)
		}
		s.selectionKind = kindNumber
		s.selectionValue = value
		return nil
	}
}

func (s *Session) describeSelection() string {
	switch s.selectionKind {
	case kindColor:
		return fmt.Sprintf("%s", s.selectionColor)
	case kindNumber:
		return fmt.Sprintf("number %d", s.selectionValue)
	default:
		return "no selection"
	}
}

func (s *Session) evaluateBet(outcome int, color string) (bool, int) {
	switch s.selectionKind {
	case kindColor:
		if s.selectionColor == color {
			if color == "green" {
				return true, 35
			}
			return true, 1
		}
	case kindNumber:
		if s.selectionValue == outcome {
			return true, 35
		}
	}
	return false, 0
}

func (s *Session) statusMessage() string {
	bet := "No active bet."
	if s.selectionKind != kindNone {
		bet = fmt.Sprintf("Pending bet: %d on %s.", s.currentBet, s.describeSelection())
	}
	return fmt.Sprintf("Bankroll: %d chips. %s", s.bank.Balance(), bet)
}

func colorFor(value int) string {
	if value == 0 {
		return "green"
	}
	if redNumbers[value] {
		return "red"
	}
	return "black"
}

var redNumbers = map[int]bool{
	1:  true,
	3:  true,
	5:  true,
	7:  true,
	9:  true,
	12: true,
	14: true,
	16: true,
	18: true,
	19: true,
	21: true,
	23: true,
	25: true,
	27: true,
	30: true,
	32: true,
	34: true,
	36: true,
}
