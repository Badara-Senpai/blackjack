package blackjack

import (
	"fmt"
	"github.com/Badara-Senpai/go-deck/deck"
)

type state int8
type Options struct {
	Decks           int
	Hands           int
	BlackjackPayout float64
}

const (
	stateBet state = iota
	statePlayerTurn
	stateDealerTurn
	stateHandOver
)

type Game struct {
	// unexported fields
	nDecks          int
	nHands          int
	blackjackPayout float64

	deck  []deck.Card
	state state

	player    []deck.Card
	playerBet int
	balance   int

	dealer   []deck.Card
	dealerAI AI
}

func New(opts Options) Game {
	g := Game{
		state:    statePlayerTurn,
		dealerAI: dealerAI{},
		balance:  0,
	}

	if opts.Decks == 0 {
		opts.Decks = 3
	}

	if opts.Hands == 0 {
		opts.Hands = 100
	}

	if opts.BlackjackPayout == 0.0 {
		opts.BlackjackPayout = 1.5
	}

	g.nDecks = opts.Decks
	g.nHands = opts.Hands
	g.blackjackPayout = opts.BlackjackPayout

	return g
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

func bet(g *Game, ai AI, shuffled bool) {
	bet := ai.Bet(shuffled)
	g.playerBet = bet
}

func (g *Game) Play(ai AI) int {
	g.deck = nil

	minLength := 52 * g.nDecks / 3

	for i := 0; i < g.nHands; i++ {
		shuffled := false

		if len(g.deck) < minLength {
			g.deck = deck.New(deck.Deck(g.nDecks), deck.Shuffle)
			shuffled = true
		}

		bet(g, ai, shuffled)

		deal(g)

		if BlackJack(g.dealer...) {
			endHand(g, ai)
			continue
		}

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

// BlackJack returns true if a hand is two cards and have a score of 21
func BlackJack(hand ...deck.Card) bool {
	return len(hand) == 2 && Score(hand...) == 21
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
	pBlackJack, dBlackJack := BlackJack(g.player...), BlackJack(g.dealer...)
	winnings := g.playerBet

	switch {
	case pBlackJack && dBlackJack:
		fmt.Println("BlackJack Draw")
		winnings = 0
	case dBlackJack:
		fmt.Println("*BlackJack* for Dealer, You Lose!")
		winnings = -winnings
	case pBlackJack:
		fmt.Println("*BlackJack* You win!")
		winnings = int(float64(winnings) * g.blackjackPayout)
	case pScore > 21:
		fmt.Println("You busted!")
		winnings = -winnings
	case dScore > 21:
		fmt.Println("Dealer busted")
	case pScore > dScore:
		fmt.Println("You win!")
	case dScore > pScore:
		fmt.Println("You lose!")
		winnings = -winnings
	case pScore == dScore:
		fmt.Println("Draw")
		winnings = 0
	}
	g.balance += winnings

	fmt.Println()

	ai.Results([][]deck.Card{g.player}, g.dealer)
	g.player = nil
	g.dealer = nil
}
