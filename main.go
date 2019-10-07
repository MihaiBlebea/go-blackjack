package main

import (
	"fmt"
	"go-blackjack/blackjack"
	"go-blackjack/persistence"
	"time"
)

func main() {
	game := blackjack.New(blackjack.Options{
		Rounds: 1000,
		Credit: 5000,
		Logger: &persistence.Blackjack{},
	})

	// ai := blackjack.HumanAi{}
	ai := blackjack.ComputerAi{}

	start := time.Now()
	finalScore := game.Play(&ai)
	end := time.Now().Sub(start)

	fmt.Printf("Time: %v \n", end)
	fmt.Printf("Credit: %v \n", finalScore)
}
