package main

import (
	"fmt"
	"go-blackjack/blackjack"
	"go-blackjack/persistence"
)

func main() {
	game := blackjack.New(blackjack.Options{
		Rounds: 1000,
		Credit: 5000,
		Logger: &persistence.Blackjack{},
	})

	// ai := blackjack.HumanAi{}
	ai := blackjack.ComputerAi{}
	finalScore := game.Play(&ai)

	fmt.Println(finalScore)
}
