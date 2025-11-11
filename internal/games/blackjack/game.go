package blackjack

import "math/rand"

// GameState represents the current point in the round lifecycle.
type GameState int

const (
	StatePlayerTurn GameState = iota
	StateDealerTurn
	StateFinished
)

// Result models the final outcome of a round.
type Result int

const (
	ResultNone Result = iota
	ResultPlayerWin
	ResultDealerWin
	ResultPush
)

// Game models a single blackjack round between the player and the dealer.
type Game struct {
	Deck   Deck
	Player Hand
	Dealer Hand

	State  GameState
	Result Result
}

// NewGame returns a shuffled round ready for play. If rnd is nil the default
// global source is used for shuffling.
func NewGame(rnd *rand.Rand) *Game {
	deck := NewDeck()
	deck.Shuffle(rnd)

	g := &Game{
		Deck:   deck,
		Player: Hand{},
		Dealer: Hand{},
		State:  StatePlayerTurn,
		Result: ResultNone,
	}

	g.dealInitial()
	return g
}

func (g *Game) dealInitial() {
	g.Player.Add(g.Deck.Draw(2)...)
	g.Dealer.Add(g.Deck.Draw(2)...)

	if g.Player.IsBlackjack() || g.Dealer.IsBlackjack() {
		// finish round
		g.finishRound()
	}
}

// PlayerHit deals a single card to the player and updates the result if the
// player busts.
func (g *Game) PlayerHit() {
	if g.State != StatePlayerTurn || g.Result != ResultNone {
		return
	}

	g.Player.Add(g.Deck.Draw(1)...)

	if g.Player.IsBust() {
		g.Result = ResultDealerWin
		g.State = StateFinished
	}
}

// PlayerStand ends the player's turn and triggers the dealer's automatic play.
func (g *Game) PlayerStand() {
	if g.State != StatePlayerTurn || g.Result != ResultNone {
		return
	}

	g.State = StateDealerTurn
	g.dealerPlay()
}

func (g *Game) dealerPlay() {
	for g.Dealer.Value() < 17 {
		g.Dealer.Add(g.Deck.Draw(1)...)
	}
	g.finishRound()
}

func (g *Game) finishRound() {
	g.State = StateFinished

	playerVal := g.Player.Value()
	dealerVal := g.Dealer.Value()

	if g.Player.IsBust() {
		g.Result = ResultDealerWin
		return
	}
	if g.Dealer.IsBust() {
		g.Result = ResultPlayerWin
		return
	}

	switch {
	case playerVal > dealerVal:
		g.Result = ResultPlayerWin
	case dealerVal > playerVal:
		g.Result = ResultDealerWin
	default:
		g.Result = ResultPush

	}

}

// IsFinished reports whether the round has reached a terminal state.
func (g *Game) IsFinished() bool {
	return g.State == StateFinished
}
