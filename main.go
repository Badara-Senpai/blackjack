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
		fmt.Println("Player: ", player)
		fmt.Println("Dealer: ", dealer.DealerString())

		fmt.Println("What will you do? (h)Hit, (s)Stand")
		fmt.Scanf("%s\n", &input)

		switch input {
		case "h":
			card, cards = drawCards(cards)
			player = append(player, card)
		}
	}

	fmt.Println("****** FINAL HANDS ******")
	fmt.Println("Player: ", player)
	fmt.Println("Dealer: ", dealer)
}

func drawCards(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}
