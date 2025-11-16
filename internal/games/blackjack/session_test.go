package blackjack

import (
	"testing"

	"tcp-server/internal/games"
)

func TestSessionBlackjackPayout(t *testing.T) {
	bank := games.NewBankroll(100)
	s := NewSession(bank)
	s.wagered = 20
	s.game = &Game{
		Player: Hand{
			{Rank: RankAce, Suit: SuitClubs},
			{Rank: RankKing, Suit: SuitSpades},
		},
		Dealer: Hand{
			{Rank: RankNine, Suit: SuitDiamonds},
			{Rank: RankSeven, Suit: SuitHearts},
		},
		State:  StateFinished,
		Result: ResultPlayerWin,
	}

	msg := s.settleRound()
	if msg == "" {
		t.Fatalf("expected payout message")
	}
	if bank.Balance() != 150 {
		t.Fatalf("expected bankroll 150, got %d", bank.Balance())
	}
}

func TestSessionPushRefund(t *testing.T) {
	bank := games.NewBankroll(50)
	s := NewSession(bank)
	s.wagered = 10
	s.game = &Game{
		Player: Hand{
			{Rank: RankTen, Suit: SuitClubs},
			{Rank: RankQueen, Suit: SuitHearts},
		},
		Dealer: Hand{
			{Rank: RankTen, Suit: SuitSpades},
			{Rank: RankQueen, Suit: SuitDiamonds},
		},
		State:  StateFinished,
		Result: ResultPush,
	}

	_ = s.settleRound()
	if bank.Balance() != 60 {
		t.Fatalf("expected bankroll 60, got %d", bank.Balance())
	}
}
