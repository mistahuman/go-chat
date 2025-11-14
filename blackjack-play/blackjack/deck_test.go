package blackjack

import (
	"testing"
)

func TestNewDeck_Has52UniqueCards(t *testing.T) {
	deck := NewDeck()
	t.Log(deck)
	if len(deck) != 52 {
		t.Fatalf("expected d52 cards, got %d", len(deck))
	}

	seen := make(map[Card]bool)
	for _, c := range deck {
		if seen[c] {
			t.Fatalf("duplicate card found: %v", c)
		}
		seen[c] = true
	}
}

func TestDraw_RemovesCardsFromDeck(t *testing.T) {
	deck := NewDeck()
	deck.Shuffle(nil)

	t.Log(deck)

	n := 2

	draw := deck.Draw(2)
	t.Log(draw)

	if len(draw) != n {
		t.Fatalf("expected to draw %d cards, got %d", n, len(draw))
	}

	if len(deck) != 52-n {
		t.Fatalf("expected deck len %d after draw, got %d", 52-n, len(draw))
	}
}
