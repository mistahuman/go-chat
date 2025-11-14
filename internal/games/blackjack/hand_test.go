package blackjack

import "testing"

func TestHandValue_NoAces(t *testing.T) {
	h := Hand{}
	h.Add(
		Card{Rank: RankTen, Suit: SuitHearts},
		Card{Rank: RankNine, Suit: SuitClubs},
	)

	if got := h.Value(); got != 19 {
		t.Fatalf("expected 19, got %d", got)
	}
}

func TestHandValue_SingleAce(t *testing.T) {
	h := Hand{}
	h.Add(
		Card{Rank: RankAce, Suit: SuitHearts},
		Card{Rank: RankNine, Suit: SuitClubs},
	)

	if got := h.Value(); got != 20 {
		t.Fatalf("expected 20, got %d", got)
	}
}

func TestHandValue_MultipleAces(t *testing.T) {
	h := Hand{}
	h.Add(
		Card{Rank: RankAce, Suit: SuitHearts},
		Card{Rank: RankAce, Suit: SuitClubs},
		Card{Rank: RankNine, Suit: SuitSpades},
	)

	if got := h.Value(); got != 21 {
		t.Fatalf("expected 21, got %d", got)
	}
}

func TestIsBlackjack(t *testing.T) {
	h := Hand{}
	h.Add(
		Card{Rank: RankAce, Suit: SuitHearts},
		Card{Rank: RankKing, Suit: SuitClubs},
	)

	if !h.IsBlackjack() {
		t.Fatalf("expected blackjack")
	}
}

func TestIsBust(t *testing.T) {
	h := Hand{}
	h.Add(
		Card{Rank: RankTen, Suit: SuitHearts},
		Card{Rank: RankTen, Suit: SuitClubs},
		Card{Rank: RankThree, Suit: SuitSpades},
	)

	if !h.IsBust() {
		t.Fatalf("expected bust")
	}
}
