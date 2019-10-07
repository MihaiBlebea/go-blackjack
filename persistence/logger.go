package persistence

import (
	"deckofcards/card"
	"encoding/json"
	"time"
)

type Logger interface {
	LogPlayerHand(cards []card.Card)
	LogDealerHand(cards []card.Card)
	LogPlayerScore(score int)
	LogDealerScore(score int)
	LogBet(bet int)
	LogCredit(credit int)
	LogResult(result int)
	Persist()
}

type Blackjack struct {
	id          int
	dealerHand  string
	playerHand  string
	dealerScore int
	playerScore int
	bet         int
	credit      int
	result      int
	created     time.Time
	updated     time.Time
}

func (b *Blackjack) StringifyPlayerCards(cardIds []int) {
	res, err := json.Marshal(cardIds)
	checkError(err)
	b.playerHand = string(res)
}

func (b *Blackjack) StringifyDealerCards(cardIds []int) {
	res, err := json.Marshal(cardIds)
	checkError(err)
	b.dealerHand = string(res)
}

func (b *Blackjack) LogPlayerHand(cards []card.Card) {
	cardsInt := cardsAsInt(cards)
	b.StringifyPlayerCards(cardsInt)
}

func (b *Blackjack) LogDealerHand(cards []card.Card) {
	cardsInt := cardsAsInt(cards)
	b.StringifyDealerCards(cardsInt)
}

func (b *Blackjack) LogPlayerScore(score int) {
	b.playerScore = score
}

func (b *Blackjack) LogDealerScore(score int) {
	b.dealerScore = score
}

func (b *Blackjack) LogBet(bet int) {
	b.bet = bet
}

func (b *Blackjack) LogCredit(credit int) {
	b.credit = credit
}

func (b *Blackjack) LogResult(result int) {
	b.result = result
}

func (b *Blackjack) Persist() {
	insert(*b)
}

func cardsAsInt(cards []card.Card) []int {
	var ret []int
	for _, crd := range cards {
		ret = append(ret, crd.GetRankInt())
	}
	return ret
}
