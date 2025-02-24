package blackjack

import (
	"fmt"
	"github.com/Badara-Senpai/go-deck/deck"
)

type state int8

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

type Game struct {
	// unexported fields
	deck     []deck.Card
	state    state
	player   []deck.Card
	dealer   []deck.Card
	balance  int
	dealerAI AI
}

func New() Game {
	return Game{
		state:    statePlayerTurn,
		dealerAI: dealerAI{},
		balance:  0,
	}
}

func (g *Game) currentPlayer() *[]deck.Card {
	switch g.state {
	case statePlayerTurn:
		return &g.player
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("It isn't currently ay player's turn")
	}
}

type Move func(*Game)

func deal(g *Game) {
	g.player = make([]deck.Card, 0, 5)
	g.dealer = make([]deck.Card, 0, 5)

	var card deck.Card
	for i := 0; i < 2; i++ {
		card, g.deck = drawCards(g.deck)
		g.player = append(g.player, card)

		card, g.deck = drawCards(g.deck)
		g.dealer = append(g.dealer, card)
	}

	g.state = statePlayerTurn
}

func (g *Game) Play(ai AI) int {
	g.deck = deck.New(deck.Deck(3), deck.Shuffle)

	for i := 0; i < 2; i++ {
		deal(g)

		for g.state == statePlayerTurn {
			hand := make([]deck.Card, len(g.player))
			copy(hand, g.player)

			move := ai.Play(hand, g.dealer[0])
			move(g)
		}

		// Let's add some logic for the dealer
		// If dealer score <= 16, we hit
		// If dealer has soft 17, then we hit
		for g.state == stateDealerTurn {
			hand := make([]deck.Card, len(g.dealer))
			copy(hand, g.dealer)

			move := g.dealerAI.Play(hand, g.dealer[0])
			move(g)
		}

		endHand(g, ai)
	}

	return g.balance
}

func MoveHit(g *Game) {
	hand := g.currentPlayer()

	var card deck.Card
	card, g.deck = drawCards(g.deck)

	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		MoveStand(g)
	}
}

func MoveStand(g *Game) {
	g.state++
}

func drawCards(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

func Score(hand ...deck.Card) int {
	minScore := minScore(hand...)

	if minScore <= 11 {
		for _, card := range hand {
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

func Soft(hand ...deck.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)

	return minScore != score
}

func minScore(hand ...deck.Card) int {
	score := 0
	for _, card := range hand {
		score += min(int(card.Rank), 10)
	}

	return score
}

func endHand(g *Game, ai AI) {
	pScore, dScore := Score(g.player...), Score(g.dealer...)

	fmt.Println("****** FINAL HANDS ******")
	fmt.Println("Player: ", g.player, "\nScore: ", pScore)
	fmt.Println("Dealer: ", g.dealer, "\nScore: ", dScore)

	switch {
	case pScore > 21:
		fmt.Println("You busted!")
		g.balance--
	case dScore > 21:
		fmt.Println("Dealer busted")
		g.balance++
	case pScore > dScore:
		fmt.Println("You win!")
		g.balance++
	case dScore > pScore:
		fmt.Println("You lose!")
		g.balance--
	case pScore == dScore:
		fmt.Println("Draw")
	}

	fmt.Println()

	ai.Results([][]deck.Card{g.player}, g.dealer)
	g.player = nil
	g.dealer = nil
}
