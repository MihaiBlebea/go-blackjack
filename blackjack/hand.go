package blackjack

import (
	"deckofcards/card"
	"deckofcards/deck"
	"strings"
)

type Hand struct {
	deck.Deck
}

func (h *Hand) MinScore() int {
	var score int
	for _, crd := range h.GetCards() {
		score += min(crd.GetRankInt(), 10)
	}
	return score
}

func (h *Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}
	for _, crd := range h.GetCards() {
		if crd.GetRank() == card.Ace {
			minScore += 10
		}
	}
	return minScore
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (h *Hand) ShowFinal() string {
	var playerCards []string
	for _, crd := range h.GetCards() {
		playerCards = append(playerCards, crd.Display()["full"])
	}
	return strings.Join(playerCards, ", ")
}

func (h *Hand) ShowPartial() string {
	var playerCards []string
	for index, crd := range h.GetCards() {
		if index == len(h.GetCards())-1 {
			playerCards = append(playerCards, "** HIDDEN **")
		} else {
			playerCards = append(playerCards, crd.Display()["full"])
		}
	}
	return strings.Join(playerCards, ", ")
}
