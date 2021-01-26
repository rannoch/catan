package domain

import (
	"time"
)

type GameStateStarted struct {
	GameStateDefault
	game *Game
}

func NewGameStateStarted(game *Game) *GameStateStarted {
	return &GameStateStarted{game: game}
}

var _ GameState = (*GameStateStarted)(nil)

func (gameStateStarted *GameStateStarted) EnterState(occurred time.Time) {
	game := gameStateStarted.game

	if err := game.GenerateBoard(occurred); err != nil {
		panic(err)
	}
	if err := game.ShufflePlayers(occurred); err != nil {
		panic(err)
	}

	initialSetupPhaseStartedEventMessage := NewEventDescriptor(
		game.Id(),
		InitialSetupPhaseStartedEvent{},
		nil,
		game.Version(),
		occurred,
	)

	game.Apply(initialSetupPhaseStartedEventMessage, true)
	game.currentState.EnterState(occurred)
}

func (gameStateStarted GameStateStarted) Apply(eventMessage EventMessage, _ bool) {
	game := gameStateStarted.game

	switch event := eventMessage.Event().(type) {
	case BoardGeneratedEvent:
		game.setBoard(event.NewBoard)
	case PlayersShuffledEvent:
		game.setTurnOrder(event.PlayersInOrder)
	case InitialSetupPhaseStartedEvent:
		game.setState(game.stateInitialSetup)
	}
}

func (gameStateStarted GameStateStarted) GenerateBoard(occurred time.Time) error {
	if gameStateStarted.game.BoardGenerator() == nil {
		return BoardGeneratorIsNotSelectedErr
	}

	boardGeneratedEventMessage := NewEventDescriptor(
		gameStateStarted.game.Id(),
		BoardGeneratedEvent{
			NewBoard: gameStateStarted.game.BoardGenerator().GenerateBoard(),
		},
		nil,
		gameStateStarted.game.Version(),
		occurred,
	)

	gameStateStarted.game.Apply(boardGeneratedEventMessage, true)

	return nil
}

func (gameStateStarted GameStateStarted) ShufflePlayers(occurred time.Time) error {
	if gameStateStarted.game.PlayersShuffler() == nil {
		return PlayersShufflerIsNotSelectedErr
	}

	playersShuffledEventMessage := NewEventDescriptor(
		gameStateStarted.game.Id(),
		PlayersShuffledEvent{
			PlayersInOrder: gameStateStarted.game.playersShuffler.Shuffle(gameStateStarted.game.turnOrder),
		},
		nil,
		gameStateStarted.game.Version(),
		occurred,
	)

	gameStateStarted.game.Apply(playersShuffledEventMessage, true)

	return nil
}
