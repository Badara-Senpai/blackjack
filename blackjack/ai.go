package blackjack

import (
	"fmt"
	"github.com/Badara-Senpai/go-deck/deck"
)

type AI interface {
	Bet(shuffled bool) int
	Play(hand []deck.Card, dealer deck.Card) Move
	Results(hand [][]deck.Card, dealer []deck.Card)
}

func HumanAI() AI {
	return humanAI{}
}

type humanAI struct{}

type dealerAI struct{}

func (ai humanAI) Bet(shuffled bool) int {
	if shuffled {
		fmt.Println("The deck was just shuffled!")
	}

	fmt.Println("What would you like to bet?")

	var bet int
	fmt.Scanf("%d\n", &bet)

	return bet
}

func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	for {
		var input string

		fmt.Println("Player: ", hand, "\t*Score: ", Score(hand...))
		fmt.Println("Dealer: ", dealer)

		fmt.Println("What will you do? (h)Hit, (s)Stand, (d)Double")
		fmt.Scanf("%s\n", &input)

		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		case "d":
			return MoveDouble
		default:
			fmt.Println("Invalid Option: ", input)
		}
	}
}

func (ai humanAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	fmt.Println("****** FINAL HANDS ******")
	fmt.Println("Player: ", hand, "\nScore: ", Score(hand[0]...))
	fmt.Println("Dealer: ", dealer, "\nScore: ", Score(dealer...))
}

func (ai dealerAI) Bet(shuffled bool) int {
	return 1
}

func (ai dealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	dScore := Score(hand...)

	if dScore <= 16 || (dScore == 17 && Soft(hand...)) {
		return MoveHit
	}

	return MoveStand
}

func (ai dealerAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	// nothing for now
}
