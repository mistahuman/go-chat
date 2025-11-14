package blackjack

import "fmt"

type Suit string
type Rank int

const (
	SuitHearts   Suit = "♥"
	SuitDiamonds Suit = "♦"
	SuitClubs    Suit = "♣"
	SuitSpades   Suit = "♠"
)

const (
	RankAce   Rank = 1
	RankTwo   Rank = 2
	RankThree Rank = 3
	RankFour  Rank = 4
	RankFive  Rank = 5
	RankSix   Rank = 6
	RankSeven Rank = 7
	RankEight Rank = 8
	RankNine  Rank = 9
	RankTen   Rank = 10
	RankJack  Rank = 11
	RankQueen Rank = 12
	RankKing  Rank = 13
)

type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string {
	var r string
	switch c.Rank {
	case RankAce:
		r = "A"
	case RankJack:
		r = "J"
	case RankQueen:
		r = "Q"
	case RankKing:
		r = "K"
	default:
		r = fmt.Sprintf("%d", int(c.Rank))
	}
	return r + string(c.Suit)
}
