package blackjack

import "testing"

func TestSessionBlackjackPayout(t *testing.T) {
	s := NewSession()
	s.bankroll = 100
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
	if s.bankroll != 150 {
		t.Fatalf("expected bankroll 150, got %d", s.bankroll)
	}
}

func TestSessionPushRefund(t *testing.T) {
	s := NewSession()
	s.bankroll = 50
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
	if s.bankroll != 60 {
		t.Fatalf("expected bankroll 60, got %d", s.bankroll)
	}
}
