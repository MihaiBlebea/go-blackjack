package blackjack

import (
	"fmt"
	"log"
	"strconv"
)

type Ai interface {
	Bet(gs *GameState)
	Play(gs *GameState)
	Result(gs *GameState)
}

type HumanAi struct{}

func (ai *HumanAi) Bet(gs *GameState) {
	var input string
	fmt.Println("\nCredit left: " + strconv.FormatInt(int64(gs.Credit), 10))
	fmt.Println("How much do you want to bet?")
	fmt.Scanf("%s\n", &input)
	amount, err := strconv.Atoi(input)
	if err != nil {
		log.Panic(err)
	}
	Bet(gs, amount)
}

func (ai *HumanAi) Play(gs *GameState) {
	var input string
	fmt.Println("\nPlayer: " + gs.Player.ShowFinal())
	fmt.Println("Score: " + strconv.FormatInt(int64(gs.Player.Score()), 10))
	fmt.Println("Dealer: " + gs.Dealer.ShowPartial())
	fmt.Println("Score: " + strconv.FormatInt(int64(gs.Dealer.Score()), 10))
	fmt.Println("What will you do? (h)it, (s)tand or (d)ouble")
	fmt.Scanf("%s\n", &input)

	switch input {
	case "h":
		MoveHit(gs)
	case "s":
		MoveStand(gs)
	case "d":
		MoveDouble(gs)
	}
}

func (ai *HumanAi) Result(gs *GameState) {
	var result string
	switch {
	case gs.Player.Score() == 21:
		result = "You won with Blackjack"
	case gs.Player.Score() > 21:
		result = "You busted"
	case gs.Dealer.Score() > 21:
		result = "Dealer busted"
	case gs.Player.Score() == gs.Dealer.Score():
		result = "Draw"
	case gs.Player.Score() > gs.Dealer.Score():
		result = "You won"
	case gs.Player.Score() < gs.Dealer.Score():
		result = "You lost"
	}

	fmt.Println("\nPlayer cards: " + fmt.Sprintf("%v", gs.Player.ShowFinal()))
	fmt.Println("Player score: " + fmt.Sprintf("%v", gs.Player.Score()))
	fmt.Println("Dealer cards: " + fmt.Sprintf("%v", gs.Dealer.ShowFinal()))
	fmt.Println("Dealer score: " + fmt.Sprintf("%v", gs.Dealer.Score()))
	fmt.Println(result)
}

type ComputerAi struct {
	seen  int
	score int
}

func (ai *ComputerAi) Bet(gs *GameState) {
	if ai.seen >= 52 {
		ai.seen = 0
		ai.score = 0
	}

	switch {
	case ai.score > 4:
		Bet(gs, 1000)
	case ai.score < 0:
		Bet(gs, 10)
	default:
		Bet(gs, 100)
	}

	// fmt.Println(ai.seen)
}

func (ai *ComputerAi) Play(gs *GameState) {
	// if gs.Player.Score() < 16 || (gs.Player.Score() == 17 && gs.Player.MinScore() != 17) {
	// 	MoveHit(gs)
	// } else {
	// 	MoveStand(gs)
	// }
	if gs.Player.Score() < 16 {
		MoveDouble(gs)
	} else {
		MoveStand(gs)
	}
}

func (ai *ComputerAi) Result(gs *GameState) {
	for _, card := range gs.Dealer.GetCards() {
		switch {
		case card.Rank < 6:
			ai.score++
		case card.Rank <= 9:
		case card.Rank > 9:
			ai.score--
		}
		ai.seen++
	}

	for _, card1 := range gs.Player.GetCards() {
		switch {
		case card1.Rank < 6:
			ai.score++
		case card1.Rank <= 9:
		case card1.Rank > 9:
			ai.score--
		}
		ai.seen++
	}
}
