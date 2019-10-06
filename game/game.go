package game

import (
	"deckofcards/deck"
	"fmt"
	"go-blackjack/hand"
	"log"
)

type Turn int

const (
	ShuffleStage Turn = iota
	BetStage
	DealCards
	PlayerTurn
	DealerTurn
	ScoreStage
	GameOver
)

type GameState struct {
	Rounds int
	Deck   deck.Deck
	Turn
	Player hand.Hand
	Dealer hand.Hand
	Bet    int
	Credit int
}

type Options struct {
	Rounds int
	Credit int
}

func (gs *GameState) CurrentPlayer() *hand.Hand {
	var ret *hand.Hand
	switch gs.Turn {
	case PlayerTurn:
		ret = &gs.Player
	case DealerTurn:
		ret = &gs.Dealer
	default:
		log.Panic("Not player turn")
	}
	return ret
}

func New(options Options) *GameState {
	return &GameState{
		Turn:   ShuffleStage,
		Credit: options.Credit,
		Rounds: options.Rounds,
	}
}

func Deal(gs *GameState) {
	if gs.Turn != DealCards {
		log.Panic("Game is not in DealCards")
	}

	gs.Player.Deck = *deck.NewFromCards(gs.Deck.DrawHand(2))
	gs.Dealer.Deck = *deck.NewFromCards(gs.Deck.DrawHand(2))

	if gs.Player.Score() == 21 {
		gs.Turn = DealerTurn
	} else {
		gs.Turn = PlayerTurn
	}
}

func Shuffle(gs *GameState) {
	if gs.Turn != ShuffleStage {
		log.Panic("Game is not in ShuffleStage")
	}
	if len(gs.Deck.GetCards()) < 10 {
		gs.Deck = *deck.New()
		gs.Deck.Shuffle()
	}
	gs.Turn = BetStage
}

func MoveHit(gs *GameState) {
	hand := gs.CurrentPlayer()
	hand.Add(gs.Deck.Draw())
	if hand.Score() >= 21 {
		MoveStand(gs)
	}
}

func MoveStand(gs *GameState) {
	gs.Turn++
}

func MoveDouble(gs *GameState) {
	if len(gs.Player.GetCards()) == 2 {
		gs.Bet = gs.Bet * 2
		MoveHit(gs)
		MoveStand(gs)
	}
}

func EndGame(gs *GameState) {
	if gs.Turn != ScoreStage {
		log.Panic("The game turn is not ScoreStage")
	}

	var res string
	switch {
	case gs.Player.Score() == 21:
		res = "You won with Blackjack"
		gs.Credit += gs.Bet + int(float64(gs.Bet)*1.5)
	case gs.Player.Score() > 21:
		res = "You busted"
		gs.Credit -= gs.Bet
	case gs.Dealer.Score() > 21:
		res = "Dealer busted"
		gs.Credit += gs.Bet
	case gs.Player.Score() == gs.Dealer.Score():
		res = "Draw"
	case gs.Player.Score() > gs.Dealer.Score():
		res = "You won"
		gs.Credit += gs.Bet
	case gs.Player.Score() < gs.Dealer.Score():
		res = "You lose"
		gs.Credit -= gs.Bet
	}

	fmt.Println("\nPlayer cards: " + fmt.Sprintf("%v", gs.Player.ShowFinal()))
	fmt.Println("Player score: " + fmt.Sprintf("%v", gs.Player.Score()))
	fmt.Println("Dealer cards: " + fmt.Sprintf("%v", gs.Dealer.ShowFinal()))
	fmt.Println("Dealer score: " + fmt.Sprintf("%v", gs.Dealer.Score()))
	fmt.Println(res)

	gs.Rounds--

	if gs.Rounds == 0 || gs.Credit <= 0 {
		gs.Turn = GameOver
	} else {
		gs.Turn = ShuffleStage
	}
}

func Bet(gs *GameState, amount int) {
	if gs.Turn != BetStage {
		log.Panic("The game is not currently in BetStage")
	}

	gs.Bet = amount
	gs.Turn = DealCards
}
