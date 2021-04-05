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

func (gameStateNew *GameStateNew) SetBoardGenerator(boardGenerator BoardGenerator, occurred time.Time) error {
	eventMessage := EventDescriptor{
		id:       gameStateNew.game.Id(),
		event:    BoardGeneratorSelectedEvent{BoardGenerator: boardGenerator},
		headers:  nil,
		version:  gameStateNew.game.Version(),
		occurred: occurred,
	}

	gameStateNew.game.Apply(eventMessage, true)
	return nil
}

func (gameStateNew *GameStateNew) SetPlayersShuffler(playersShuffler PlayersShuffler, occurred time.Time) error {
	eventMessage := EventDescriptor{
		id:       gameStateNew.game.Id(),
		event:    PlayersShufflerSelectedEvent{PlayersShuffler: playersShuffler},
		headers:  nil,
		version:  gameStateNew.game.Version(),
		occurred: occurred,
	}

	gameStateNew.game.Apply(eventMessage, true)
	return nil
}

func (gameStateNew *GameStateNew) SetDiceRoller(diceRoller DiceRoller, occurred time.Time) error {
	eventMessage := EventDescriptor{
		id:       gameStateNew.game.Id(),
		event:    DiceRollerSelected{DiceRoller: diceRoller},
		headers:  nil,
		version:  gameStateNew.game.Version(),
		occurred: occurred,
	}

	gameStateNew.game.Apply(eventMessage, true)
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
	game.currentState.EnterState(occurred)

	return nil
}

func (gameStateNew *GameStateNew) Apply(eventMessage EventMessage, _ bool) {
	game := gameStateNew.game

	switch event := eventMessage.Event().(type) {
	case PlayerJoinedTheGameEvent:
		game.addPlayer(event.Player)
	case PlayerLeftTheGameEvent:
		game.removePlayer(event.Player)
	case BoardGeneratorSelectedEvent:
		game.setBoardGenerator(event.BoardGenerator)
	case PlayersShufflerSelectedEvent:
		game.setPlayersShuffler(event.PlayersShuffler)
	case DiceRollerSelected:
		game.setDiceRoller(event.DiceRoller)
	case GameStartedEvent:
		game.setState(game.stateStarted)
	}
}
