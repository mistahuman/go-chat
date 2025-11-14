package blackjack

import (
	"fmt"
	"strings"
)

// Session exposes the blackjack game through the generic games.Session interface.
type Session struct {
	game *Game
}

// NewSession returns a new blackjack session ready to be started.
func NewSession() *Session {
	return &Session{}
}

// Name satisfies the games.Session interface.
func (s *Session) Name() string {
	return "blackjack"
}

// Handle interprets user actions and returns the resulting messages.
func (s *Session) Handle(action string, args []string) ([]string, bool, error) {
	switch strings.ToLower(action) {
	case "start":
		s.game = NewGame(nil)
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
		if err := s.ensureActive(); err != nil {
			return nil, false, err
		}
		messages := []string{
			formatBlackjackHand("Your hand", s.game.Player),
			fmt.Sprintf("Dealer shows %s and [?]", s.game.Dealer[0].String()),
		}
		return messages, false, nil
	case "quit":
		if s.game == nil {
			return []string{"No active blackjack game."}, true, nil
		}
		s.game = nil
		return []string{"Blackjack game discarded."}, true, nil
	default:
		return nil, false, fmt.Errorf("unknown blackjack action: %s", action)
	}
}

func (s *Session) ensureActive() error {
	if s.game == nil {
		return fmt.Errorf("start a blackjack game first with /blackjack start")
	}
	if s.game.IsFinished() {
		return fmt.Errorf("Round finished. Start a new game with /blackjack start")
	}
	return nil
}

func (s *Session) startMessages() []string {
	messages := []string{
		"Blackjack game started.",
		fmt.Sprintf("Dealer shows %s and [?]", s.game.Dealer[0].String()),
		formatBlackjackHand("Your hand", s.game.Player),
	}
	if s.game.IsFinished() {
		messages = append(messages, s.finishMessages()...)
	} else {
		messages = append(messages, "Use /blackjack hit or /blackjack stand.")
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
	return messages
}

func formatBlackjackHand(label string, hand Hand) string {
	return fmt.Sprintf("%s: %s", label, hand.String())
}
