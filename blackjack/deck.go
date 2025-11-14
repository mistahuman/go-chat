package blackjack

import (
	"math/rand"
	"time"
)

// Deck deck of cards
type Deck []Card

// NewDeck create a standard deck of unordered cards
func NewDeck() Deck {
	deck := make([]Card, 0, 52)
	suits := []Suit{SuitHearts, SuitDiamonds, SuitClubs, SuitSpades}

	for _, s := range suits {
		for r := RankAce; r <= RankKing; r++ {
			deck = append(deck, Card{Rank: r, Suit: s})
		}
	}

	return deck
}

func (d Deck) Shuffle(rnd *rand.Rand) {
	if len(d) == 0 {
		return
	}

	if rnd == nil {
		rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	rnd.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})

}

func (d *Deck) Draw(n int) []Card {
	if n <= 0 || len(*d) == 0 {
		return nil
	}

	if n > len(*d) {
		n = len(*d)
	}

	cards := (*d)[:n]
	*d = (*d)[n:]

	return cards

}
