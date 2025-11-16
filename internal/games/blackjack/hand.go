package blackjack

import "fmt"

type Hand []Card

func (h *Hand) Add(cards ...Card) {
	*h = append(*h, cards...)
}

func (h Hand) Value() int {
	total := 0
	aces := 0

	for _, c := range h {
		switch c.Rank {
		case RankJack, RankQueen, RankKing:
			total += 10
		case RankAce:
			aces++
		default:
			total += int(c.Rank)
		}
	}

	total += aces

	for aces > 0 && total+10 <= 21 {
		total += 10
		aces--
	}

	return total
}

func (h Hand) IsBlackjack() bool {
	if len(h) != 2 {
		return false
	}
	return h.Value() == 21
}

func (h Hand) IsBust() bool {
	return h.Value() > 21
}

func (h Hand) String() string {
	if len(h) == 0 {
		return "[]"
	}
	s := ""
	for i, c := range h {
		if i > 0 {
			s += " "
		}
		s += c.String()
	}
	s += fmt.Sprintf(" (%d)", h.Value())
	return s
}
