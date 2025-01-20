package main

import (
	"fmt"
	"github.com/Badara-Senpai/go-deck/deck"
	"strings"
)

type Hand []deck.Card

func (h Hand) String() string {
	cardsInHand := make([]string, len(h))

	for i := range h {
		cardsInHand[i] = h[i].String()
	}

	return strings.Join(cardsInHand, ", ")
}

func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

func (h Hand) Score() int {
	minScore := h.MinScore()

	if minScore <= 11 {
		for _, card := range h {
			if card.Rank == deck.Ace {
				// Ace is currently worth 1, and the player has a score of 11 or less
				// making the Ace worth 11 would make the player bust
				// 11 - 1 = 10
				return minScore + 10
			}
		}
	}

	return minScore
}

func (h Hand) MinScore() int {
	score := 0
	for _, card := range h {
		score += min(int(card.Rank), 10)
	}

	return score
}

func main() {
	cards := deck.New(deck.Deck(3), deck.Shuffle)

	var card deck.Card
	var player, dealer Hand

	for i := 0; i < 2; i++ {
		for _, hand := range []*Hand{&player, &dealer} {
			card, cards = drawCards(cards)

			*hand = append(*hand, card)
		}
	}

	var input string
	for input != "s" {
		fmt.Println("Player: ", player, "\t*Score: ", player.Score())
		fmt.Println("Dealer: ", dealer.DealerString())

		fmt.Println("What will you do? (h)Hit, (s)Stand")
		fmt.Scanf("%s\n", &input)

		switch input {
		case "h":
			card, cards = drawCards(cards)
			player = append(player, card)
		}
	}

	// Lets add some logic for the dealer
	// If dealer score <= 16, we hit
	// If dealer has soft 17, then we hit
	for dealer.Score() <= 16 || (dealer.Score() == 17 && dealer.MinScore() != 17) {
		card, cards = drawCards(cards)
		dealer = append(dealer, card)
	}

	pScore, dScore := player.Score(), dealer.Score()
	fmt.Println("****** FINAL HANDS ******")
	fmt.Println("Player: ", player, "\nScore: ", pScore)
	fmt.Println("Dealer: ", dealer, "\nScore: ", dScore)

	switch {
	case pScore > 21:
		fmt.Println("You busted!")
	case dScore > 21:
		fmt.Println("Dealer busted")
	case pScore > dScore:
		fmt.Println("You win!")
	case dScore > pScore:
		fmt.Println("You lose!")
	case pScore == dScore:
		fmt.Println("Draw")
	}
}

func drawCards(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}
