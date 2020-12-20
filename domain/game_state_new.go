package domain

import (
	"errors"
	"time"
)

var (
	BoardGeneratorIsNotSelected  = errors.New("board generator is not selected")
	PlayersShufflerIsNotSelected = errors.New("players shuffler is not selected")
	NoPlayersError               = errors.New("cannot start the game without players")
)

type GameStateNew struct {
	GameStateDefault
	game *Game
}

func NewGameStateNew(game *Game) *GameStateNew {
	return &GameStateNew{game: game}
}

var _ GameState = (*GameStateNew)(nil)

func (gameStateNew GameStateNew) EnterState(occurred time.Time) {
	// todo if there are some
}

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
		id:       string(gameStateNew.game.Id()),
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
		id:       string(gameStateNew.game.Id()),
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
		return NoPlayersError
	}

	if game.PlayersShuffler() == nil {
		return PlayersShufflerIsNotSelected
	}

	if game.BoardGenerator() == nil {
		return BoardGeneratorIsNotSelected
	}

	gameStartedEventMessage := NewEventDescriptor(
		string(gameStateNew.game.Id()),
		GameStartedEvent{Occurred: occurred},
		nil,
		gameStateNew.game.Version(),
		occurred,
	)

	gameStateNew.game.Apply(gameStartedEventMessage, true)

	return nil
}

func (gameStateNew *GameStateNew) Apply(eventMessage EventMessage, _ bool) {
	switch event := eventMessage.Event().(type) {
	case PlayerJoinedTheGameEvent:
		if event.Player.color == None {
			// set first available color
			for _, color := range allColors {
				_, err := gameStateNew.game.Player(color)
				if err == PlayerNotExistsErr {
					event.Player.SetColor(color)
					break
				}
			}
		}

		gameStateNew.game.players[event.Player.color] = event.Player
	case PlayerLeftTheGameEvent:
		delete(gameStateNew.game.players, event.Player.color)
	case GameStartedEvent:
		gameStateNew.game.setState(NewGameStateStarted(gameStateNew.game), event.Occurred)
	}
}
