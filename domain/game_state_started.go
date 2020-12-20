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

func (gameStateStarted GameStateStarted) EnterState(occurred time.Time) {
	boardGeneratedEventMessage := NewEventDescriptor(
		string(gameStateStarted.game.Id()),
		BoardGeneratedEvent{
			newBoard: gameStateStarted.game.BoardGenerator().GenerateBoard(),
		},
		nil,
		gameStateStarted.game.Version(),
		occurred,
	)

	gameStateStarted.game.Apply(boardGeneratedEventMessage, true)

	playersShuffledEventMessage := NewEventDescriptor(
		string(gameStateStarted.game.Id()),
		PlayersShuffledEvent{
			playersInOrder: gameStateStarted.game.playersShuffler.Shuffle(gameStateStarted.game.Players()),
		},
		nil,
		gameStateStarted.game.Version(),
		occurred,
	)

	gameStateStarted.game.Apply(playersShuffledEventMessage, true)

	initialSetupPhaseStartedEventMessage := NewEventDescriptor(
		string(gameStateStarted.game.Id()),
		InitialSetupPhaseStartedEvent{},
		nil,
		gameStateStarted.game.Version(),
		occurred,
	)

	gameStateStarted.game.Apply(initialSetupPhaseStartedEventMessage, true)
}

func (gameStateStarted GameStateStarted) Apply(eventMessage EventMessage, _ bool) {
	switch event := eventMessage.Event().(type) {
	case BoardGeneratedEvent:
		gameStateStarted.game.setBoard(event.newBoard)
	case PlayersShuffledEvent:
		gameStateStarted.game.setTurnOrder(event.playersInOrder)
	case InitialSetupPhaseStartedEvent:
		gameStateStarted.game.setState(NewGameStateInitialSetup(gameStateStarted.game), event.occurred)
	}
}
