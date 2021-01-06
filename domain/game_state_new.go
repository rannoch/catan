package domain

import (
	"errors"
	"time"
)

var (
	BoardGeneratorIsNotSelectedErr  = errors.New("board generator is not selected")
	PlayersShufflerIsNotSelectedErr = errors.New("players shuffler is not selected")
	NoPlayersErr                    = errors.New("cannot start the game without players")
)

type GameStateNew struct {
	GameStateDefault
	game *Game
}

func NewGameStateNew(game *Game) *GameStateNew {
	return &GameStateNew{game: game}
}

var _ GameState = (*GameStateNew)(nil)

func (gameStateNew *GameStateNew) SetBoardGenerator(boardGenerator BoardGenerator) error {
	gameStateNew.game.setBoardGenerator(boardGenerator)
	return nil
}

func (gameStateNew *GameStateNew) SetPlayersShuffler(playersShuffler PlayersShuffler) error {
	gameStateNew.game.setPlayersShuffler(playersShuffler)
	return nil
}

func (gameStateNew GameStateNew) AddPlayer(player Player, occurred time.Time) error {
	// todo game is full condition

	eventMessage := EventDescriptor{
		id:       gameStateNew.game.Id(),
		event:    PlayerJoinedTheGameEvent{Player: player},
		headers:  nil,
		version:  gameStateNew.game.Version(),
		occurred: occurred,
	}

	gameStateNew.game.Apply(eventMessage, true)
	return nil
}

func (gameStateNew GameStateNew) RemovePlayer(player Player, occurred time.Time) error {
	// todo
	eventMessage := EventDescriptor{
		id:       gameStateNew.game.Id(),
		event:    PlayerLeftTheGameEvent{Player: player},
		headers:  nil,
		version:  gameStateNew.game.Version(),
		occurred: occurred,
	}

	gameStateNew.game.Apply(eventMessage, true)
	return nil
}

func (gameStateNew *GameStateNew) StartGame(occurred time.Time) error {
	game := gameStateNew.game

	if len(game.Players()) == 0 {
		return NoPlayersErr
	}

	if game.PlayersShuffler() == nil {
		return PlayersShufflerIsNotSelectedErr
	}

	if game.BoardGenerator() == nil {
		return BoardGeneratorIsNotSelectedErr
	}

	gameStartedEventMessage := NewEventDescriptor(
		game.Id(),
		GameStartedEvent{},
		nil,
		game.Version(),
		occurred,
	)

	game.Apply(gameStartedEventMessage, true)
	game.state.EnterState(occurred)

	return nil
}

func (gameStateNew *GameStateNew) Apply(eventMessage EventMessage, _ bool) {
	game := gameStateNew.game

	switch event := eventMessage.Event().(type) {
	case PlayerJoinedTheGameEvent:
		game.addPlayer(event.Player)
	case PlayerLeftTheGameEvent:
		game.removePlayer(event.Player)
	case GameStartedEvent:
		game.setState(NewGameStateStarted(game))
	}
}
