package blackjack

import (
	"deckofcards/deck"
	"go-blackjack/persistence"
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

type GameResult int

const (
	Draw GameResult = iota
	Blackjack
	DealerWon
	PlayerWon
	DealerBust
	PlayerBust
)

type GameState struct {
	Rounds int
	Deck   deck.Deck
	Turn
	Player Hand
	Dealer Hand
	Bet    int
	Credit int
	Logger persistence.Logger
}

type Options struct {
	Rounds int
	Credit int
	Logger persistence.Logger
}

func (gs *GameState) CurrentPlayer() *Hand {
	var ret *Hand
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
		Logger: options.Logger,
	}
}

func (gs *GameState) Play(ai Ai) int {
	for gs.Turn != GameOver {
		Shuffle(gs)

		// Betting stage
		ai.Bet(gs)

		Deal(gs)

		// Player turn
		for gs.Turn == PlayerTurn {
			ai.Play(gs)
		}

		// Dealer turn
		for gs.Turn == DealerTurn {
			if gs.Dealer.Score() < 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
				MoveHit(gs)
			} else {
				MoveStand(gs)
			}
		}

		EndGame(gs)
		ai.Result(gs)
	}

	return gs.Credit
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

	var result GameResult
	switch {
	case gs.Player.Score() == 21:
		result = Blackjack
		gs.Credit += gs.Bet + int(float64(gs.Bet)*1.5)
	case gs.Player.Score() > 21:
		result = PlayerBust
		gs.Credit -= gs.Bet
	case gs.Dealer.Score() > 21:
		result = DealerBust
		gs.Credit += gs.Bet
	case gs.Player.Score() == gs.Dealer.Score():
		result = Draw
	case gs.Player.Score() > gs.Dealer.Score():
		result = PlayerWon
		gs.Credit += gs.Bet
	case gs.Player.Score() < gs.Dealer.Score():
		result = DealerWon
		gs.Credit -= gs.Bet
	}

	// Log game
	gs.Logger.LogDealerHand(gs.Dealer.GetCards())
	gs.Logger.LogPlayerHand(gs.Player.GetCards())
	gs.Logger.LogDealerScore(gs.Dealer.Score())
	gs.Logger.LogPlayerScore(gs.Player.Score())
	gs.Logger.LogBet(gs.Bet)
	gs.Logger.LogResult(int(result))
	gs.Logger.LogCredit(gs.Credit)
	gs.Logger.Persist()

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
