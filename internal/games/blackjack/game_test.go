package blackjack

import "testing"

func TestNewGame_DealsTwoCardsEach(t *testing.T) {
	game := NewGame(nil)

	if len(game.Player) != 2 {
		t.Fatalf("expected player to have 2 cards, got %d", len(game.Player))
	}
	if len(game.Dealer) != 2 {
		t.Fatalf("expected dealer to have 2 cards, got %d", len(game.Dealer))
	}

	if len(game.Deck) != 52-4 {
		t.Fatalf("expected deck to have %d cards left, got %d", 52-4, len(game.Deck))
	}
}

func TestPlayerHit_Busts(t *testing.T) {
	game := &Game{
		Deck:   Deck{{Rank: RankTen, Suit: SuitHearts}},
		Player: Hand{},
		Dealer: Hand{},
		State:  StatePlayerTurn,
		Result: ResultNone,
	}

	game.Player.Add(
		Card{Rank: RankKing, Suit: SuitClubs},   // 10
		Card{Rank: RankQueen, Suit: SuitSpades}, // 10 -> 20 totale
	)

	game.PlayerHit()

	if !game.IsFinished() {
		t.Fatalf("expected game to be finished after player busts")
	}
	if game.Result != ResultDealerWin {
		t.Fatalf("expected dealer to win after player bust, got %v", game.Result)
	}
}

func TestPlayerStand_DealerPlaysAndFinishes(t *testing.T) {
	// Dealer ha 10 (5+5) e pesca un 7 → 17
	game := &Game{
		Deck:   Deck{{Rank: RankSeven, Suit: SuitHearts}},
		Player: Hand{},
		Dealer: Hand{},
		State:  StatePlayerTurn,
		Result: ResultNone,
	}

	// Player 18 (vincerà contro 17)
	game.Player.Add(
		Card{Rank: RankTen, Suit: SuitClubs},
		Card{Rank: RankEight, Suit: SuitDiamonds},
	)

	// Dealer 10
	game.Dealer.Add(
		Card{Rank: RankFive, Suit: SuitClubs},
		Card{Rank: RankFive, Suit: SuitDiamonds},
	)

	game.PlayerStand()

	if !game.IsFinished() {
		t.Fatalf("expected game to be finished after dealer plays")
	}

	dealerVal := game.Dealer.Value()
	if dealerVal < 17 {
		t.Fatalf("expected dealer to have at least 17, got %d", dealerVal)
	}

	if game.Result != ResultPlayerWin {
		t.Fatalf("expected player to win (18 vs dealer %d), got %v", dealerVal, game.Result)
	}
}
