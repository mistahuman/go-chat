package blackjack

import "math/rand"

type GameState int

const (
	StatePlayerTurn GameState = iota
	StateDealerTurn
	StateFinished
)

type Result int

const (
	ResultNone Result = iota
	ResultPlayerWin
	ResultDealerWin
	ResultPush
)

type Game struct {
	Deck   Deck
	Player Hand
	Dealer Hand

	State  GameState
	Result Result
}

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

func (g *Game) IsFinished() bool {
	return g.State == StateFinished
}
