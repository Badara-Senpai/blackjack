package main

import (
	"fmt"
	"github.com/Badara-Senpai/blackjack/blackjack"
)

func main() {
	opts := blackjack.Options{
		Decks:           3,
		Hands:           1,
		BlackjackPayout: 1.5,
	}
	game := blackjack.New(opts)
	winnings := game.Play(blackjack.HumanAI())

	fmt.Println(winnings)

}

//package main
//
//import (
//	"fmt"
//	"github.com/Badara-Senpai/go-deck/deck"
//	"strings"
//)
//
//type State int8
//
//const (
//	StatePlayerTurn State = iota
//	StateDealerTurn
//	StateHandOver
//)
//
//type Hand []deck.Card
//
//type GameState struct {
//	Deck   []deck.Card
//	State  State
//	Player Hand
//	Dealer Hand
//}
//
//func (h Hand) String() string {
//	cardsInHand := make([]string, len(h))
//
//	for i := range h {
//		cardsInHand[i] = h[i].String()
//	}
//
//	return strings.Join(cardsInHand, ", ")
//}
//
//func (h Hand) DealerString() string {
//	return h[0].String() + ", **HIDDEN**"
//}
//
//func (h Hand) Score() int {
//	minScore := h.MinScore()
//
//	if minScore <= 11 {
//		for _, card := range h {
//			if card.Rank == deck.Ace {
//				// Ace is currently worth 1, and the player has a score of 11 or less
//				// making the Ace worth 11 would make the player bust
//				// 11 - 1 = 10
//				return minScore + 10
//			}
//		}
//	}
//
//	return minScore
//}
//
//func (h Hand) MinScore() int {
//	score := 0
//	for _, card := range h {
//		score += min(int(card.Rank), 10)
//	}
//
//	return score
//}
//
//func Shuffle(gs GameState) GameState {
//	newGS := clone(gs)
//	newGS.Deck = deck.New(deck.Deck(3), deck.Shuffle)
//
//	return newGS
//}
//
//func Deal(gs GameState) GameState {
//	newGS := clone(gs)
//
//	newGS.Player = make(Hand, 0, 5)
//	newGS.Dealer = make(Hand, 0, 5)
//
//	var card deck.Card
//	for i := 0; i < 2; i++ {
//		card, newGS.Deck = drawCards(newGS.Deck)
//		newGS.Player = append(newGS.Player, card)
//
//		card, newGS.Deck = drawCards(newGS.Deck)
//		newGS.Dealer = append(newGS.Dealer, card)
//	}
//
//	newGS.State = StatePlayerTurn
//	return newGS
//}
//
//func Hit(gs GameState) GameState {
//	newGS := clone(gs)
//	hand := newGS.CurrentPlayer()
//
//	var card deck.Card
//	card, newGS.Deck = drawCards(newGS.Deck)
//
//	*hand = append(*hand, card)
//	if hand.Score() > 21 {
//		return Stand(newGS)
//	}
//
//	return newGS
//}
//
//func Stand(gs GameState) GameState {
//	newGS := clone(gs)
//	newGS.State++
//
//	return newGS
//}
//
//func EndHand(gs GameState) GameState {
//	newGS := clone(gs)
//
//	pScore, dScore := newGS.Player.Score(), newGS.Dealer.Score()
//	fmt.Println("****** FINAL HANDS ******")
//	fmt.Println("Player: ", newGS.Player, "\nScore: ", pScore)
//	fmt.Println("Dealer: ", newGS.Dealer, "\nScore: ", dScore)
//
//	switch {
//	case pScore > 21:
//		fmt.Println("You busted!")
//	case dScore > 21:
//		fmt.Println("Dealer busted")
//	case pScore > dScore:
//		fmt.Println("You win!")
//	case dScore > pScore:
//		fmt.Println("You lose!")
//	case pScore == dScore:
//		fmt.Println("Draw")
//	}
//
//	fmt.Println()
//	newGS.Player = nil
//	newGS.Dealer = nil
//
//	return newGS
//}
//
//func main() {
//	var gs GameState
//	gs = Shuffle(gs)
//
//	for i := 0; i < 10; i++ {
//		gs = Deal(gs)
//
//		var input string
//		for input != "s" {
//			fmt.Println("Player: ", gs.Player, "\t*Score: ", gs.Player.Score())
//			fmt.Println("Dealer: ", gs.Dealer.DealerString())
//
//			fmt.Println("What will you do? (h)Hit, (s)Stand")
//			fmt.Scanf("%s\n", &input)
//
//			switch input {
//			case "h":
//				gs = Hit(gs)
//			case "s":
//				gs = Stand(gs)
//			default:
//				fmt.Println("Invalid Option: ", input)
//			}
//		}
//
//		// Let's add some logic for the dealer
//		// If dealer score <= 16, we hit
//		// If dealer has soft 17, then we hit
//		for gs.State == StateDealerTurn {
//			if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
//				gs = Hit(gs)
//			} else {
//				gs = Stand(gs)
//			}
//		}
//
//		gs = EndHand(gs)
//	}
//}
//
//func drawCards(cards []deck.Card) (deck.Card, []deck.Card) {
//	return cards[0], cards[1:]
//}
//
//func (gs *GameState) CurrentPlayer() *Hand {
//	switch gs.State {
//	case StatePlayerTurn:
//		return &gs.Player
//	case StateDealerTurn:
//		return &gs.Dealer
//	default:
//		panic("It isn't currently ay player's turn")
//	}
//}
//
//func clone(gs GameState) GameState {
//	newGs := GameState{
//		Deck:   make([]deck.Card, len(gs.Deck)),
//		State:  gs.State,
//		Player: make(Hand, len(gs.Player)),
//		Dealer: make(Hand, len(gs.Dealer)),
//	}
//
//	copy(newGs.Deck, gs.Deck)
//	copy(newGs.Player, gs.Player)
//	copy(newGs.Dealer, gs.Dealer)
//
//	return newGs
//}
