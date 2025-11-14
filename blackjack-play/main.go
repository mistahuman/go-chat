package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"blackjack-play/blackjack"
)

// func main() {
// 	// deck := blackjack.NewDeck()
// 	// deck.Shuffle(nil)

// 	// var hand blackjack.Hand
// 	// hand.Add(deck.Draw(2)...)

// 	// fmt.Println("Player hand:", hand.String())
// 	// fmt.Println("Value:", hand.Value())

// 	// if hand.IsBlackjack() {
// 	// 	fmt.Println("BLACKJACK!")
// 	// }

// 	// hand.Add(deck.Draw(1)...)
// 	// fmt.Println("After hit:", hand.String())

// 	// if hand.IsBust() {
// 	// 	fmt.Println("BUST!")
// 	// }

// 	game := blackjack.NewGame(nil)

// 	fmt.Println("Player:", game.Player.String())
// 	fmt.Println("Dealer:", game.Dealer.String())
// 	fmt.Println("State:", game.State, "Result:", game.Result)

// 	game.PlayerHit()
// 	fmt.Println("After hit, player:", game.Player.String(), "value:", game.Player.Value())

// 	if game.IsFinished() {
// 		fmt.Println("Game finished, result:", game.Result)
// 		return
// 	}

// 	game.PlayerStand()
// 	fmt.Println("Final player:", game.Player.String(), "value:", game.Player.Value())
// 	fmt.Println("Final dealer:", game.Dealer.String(), "value:", game.Dealer.Value())
// 	fmt.Println("Result:", game.Result)
// }

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("=== New game ===")
		game := blackjack.NewGame(nil)

		playRound(reader, game)

		fmt.Print("\nPlay again? (y/n): ")
		ans, _ := reader.ReadString('\n')
		ans = strings.TrimSpace(strings.ToLower(ans))
		if ans != "y" && ans != "yes" {
			break
		}
	}
}

func playRound(reader *bufio.Reader, game *blackjack.Game) {
	fmt.Println("Dealer shows:", game.Dealer[0].String(), "and [?]")
	fmt.Println("Your hand:", game.Player.String())

	if game.IsFinished() {
		printResult(game)
		return
	}

	for !game.IsFinished() {
		fmt.Print("\n(hit = h, stand = s, quit = q) > ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "h", "hit":
			game.PlayerHit()
			fmt.Println("You drew. Your hand:", game.Player.String())

		case "s", "stand":
			game.PlayerStand()

		case "q", "quit":
			fmt.Println("Quitting game.")
			return

		default:
			fmt.Println("Unknown command, use h/s/q.")
		}
	}

	fmt.Println()
	printResult(game)
}

func printResult(game *blackjack.Game) {
	fmt.Println("=== Final hands ===")
	fmt.Println("Dealer:", game.Dealer.String())
	fmt.Println("Player:", game.Player.String())

	switch game.Result {
	case blackjack.ResultPlayerWin:
		fmt.Println(">>> You win! 🎉")
	case blackjack.ResultDealerWin:
		fmt.Println(">>> Dealer wins. 😭")
	case blackjack.ResultPush:
		fmt.Println(">>> Push (tie). 🤝")
	default:
		fmt.Println(">>> No result? Something went wrong.")
	}
}
