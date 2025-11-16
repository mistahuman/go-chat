package blackjack

import (
	"fmt"
	"strconv"
	"strings"

	"tcp-server/internal/games"
)

// Session exposes the blackjack game through the generic games.Session interface.
type Session struct {
	game *Game

	bank       games.Bank
	currentBet int
	wagered    int
	minBet     int
	maxBet     int
	settled    bool
}

// NewSession returns a new blackjack session ready to be started.
func NewSession(bank games.Bank) *Session {
	if bank == nil {
		panic("blackjack session requires bank")
	}
	return &Session{
		bank:       bank,
		currentBet: 10,
		minBet:     5,
		maxBet:     200,
	}
}

// Name satisfies the games.Session interface.
func (s *Session) Name() string {
	return "bl"
}

// Handle interprets user actions and returns the resulting messages.
func (s *Session) Handle(action string, args []string) ([]string, bool, error) {
	switch strings.ToLower(action) {
	case "start":
		if err := s.startRound(); err != nil {
			return nil, false, err
		}
		return s.startMessages(), s.game.IsFinished(), nil
	case "hit":
		if err := s.ensureActive(); err != nil {
			return nil, false, err
		}
		s.game.PlayerHit()
		messages := []string{"You drew a card.", formatBlackjackHand("Your hand", s.game.Player)}
		if s.game.IsFinished() {
			messages = append(messages, s.finishMessages()...)
			return messages, true, nil
		}
		messages = append(messages, fmt.Sprintf("Dealer shows %s and [?]", s.game.Dealer[0].String()))
		return messages, false, nil
	case "stand":
		if err := s.ensureActive(); err != nil {
			return nil, false, err
		}
		s.game.PlayerStand()
		return s.finishMessages(), true, nil
	case "status":
		if s.game == nil {
			return []string{s.bankrollStatus()}, false, nil
		}
		messages := []string{
			formatBlackjackHand("Your hand", s.game.Player),
			fmt.Sprintf("Dealer shows %s and [?]", s.game.Dealer[0].String()),
			s.bankrollStatus(),
		}
		return messages, false, nil
	case "quit":
		if s.game == nil {
			return []string{"No active blackjack game."}, true, nil
		}
		s.game = nil
		s.wagered = 0
		s.settled = true
		return []string{"Blackjack game discarded.", s.bankrollStatus()}, true, nil
	case "bet":
		if len(args) == 0 {
			return nil, false, fmt.Errorf("usage: /bl bet <amount>")
		}
		amount, err := strconv.Atoi(args[0])
		if err != nil || amount <= 0 {
			return nil, false, fmt.Errorf("invalid bet amount")
		}
		if amount < s.minBet || amount > s.maxBet {
			return nil, false, fmt.Errorf("bet must be between %d and %d", s.minBet, s.maxBet)
		}
		s.currentBet = amount
		return []string{fmt.Sprintf("Bet set to %d chips", s.currentBet)}, false, nil
	case "bankroll":
		return []string{s.bankrollStatus()}, false, nil
	case "reset":
		s.currentBet = 10
		s.wagered = 0
		s.settled = true
		s.game = nil
		return []string{"Blackjack session reset.", s.bankrollStatus()}, true, nil
	default:
		return nil, false, fmt.Errorf("unknown blackjack action: %s", action)
	}
}

func (s *Session) startRound() error {
	if s.game != nil && !s.game.IsFinished() {
		return fmt.Errorf("finish the current round before starting a new one")
	}
	if s.currentBet < s.minBet || s.currentBet > s.maxBet {
		return fmt.Errorf("bet must be between %d and %d", s.minBet, s.maxBet)
	}
	if err := s.bank.Debit(s.currentBet); err != nil {
		return err
	}
	s.wagered = s.currentBet
	s.game = NewGame(nil)
	s.settled = false
	return nil
}

func (s *Session) ensureActive() error {
	if s.game == nil {
		return fmt.Errorf("start a blackjack game first with /bl start")
	}
	if s.game.IsFinished() {
		return fmt.Errorf("Round finished. Start a new game with /bl start")
	}
	return nil
}

func (s *Session) startMessages() []string {
	messages := []string{
		"Blackjack game started.",
		fmt.Sprintf("Dealer shows %s and [?]", s.game.Dealer[0].String()),
		formatBlackjackHand("Your hand", s.game.Player),
		s.bankrollStatus(),
	}
	if s.game.IsFinished() {
		messages = append(messages, s.finishMessages()...)
	} else {
		messages = append(messages, "Use /bl hit, /bl stand or /bl status.")
	}
	return messages
}

func (s *Session) finishMessages() []string {
	messages := []string{
		"Round finished.",
		formatBlackjackHand("Dealer", s.game.Dealer),
		formatBlackjackHand("Player", s.game.Player),
	}

	switch s.game.Result {
	case ResultPlayerWin:
		messages = append(messages, "Result: you win!")
	case ResultDealerWin:
		messages = append(messages, "Result: dealer wins.")
	case ResultPush:
		messages = append(messages, "Result: push.")
	default:
		messages = append(messages, "Result: no winner.")
	}
	if payout := s.settleRound(); payout != "" {
		messages = append(messages, payout)
	}
	messages = append(messages, s.bankrollStatus())
	messages = append(messages, "Use /bl start to play again.")
	return messages
}

func formatBlackjackHand(label string, hand Hand) string {
	return fmt.Sprintf("%s: %s", label, hand.String())
}

func (s *Session) settleRound() string {
	if s.game == nil || s.settled || !s.game.IsFinished() {
		return ""
	}
	payout := 0
	switch s.game.Result {
	case ResultPlayerWin:
		payout = s.wagered * 2
		if s.game.Player.IsBlackjack() && !s.game.Dealer.IsBlackjack() {
			payout = s.wagered + (s.wagered*3)/2
		}
	case ResultPush:
		payout = s.wagered
	}
	s.bank.Credit(payout)
	s.wagered = 0
	s.settled = true
	if payout == 0 {
		return "Bankroll unchanged."
	}
	return fmt.Sprintf("Bankroll credited with %d chips.", payout)
}

func (s *Session) bankrollStatus() string {
	return fmt.Sprintf("Bankroll: %d chips. Current bet: %d", s.bank.Balance(), s.currentBet)
}
