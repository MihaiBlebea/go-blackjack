package main

import (
	"fmt"
	"go-blackjack/blackjack"
)

func main() {
	game := blackjack.New(blackjack.Options{
		Rounds: 100,
		Credit: 5000,
	})

	// ai := blackjack.HumanAi{}
	ai := blackjack.ComputerAi{}
	finalScore := game.Play(&ai)

	fmt.Println(finalScore)
}
