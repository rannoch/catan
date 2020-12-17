package domain

import (
	"github.com/rannoch/catan/grid"
	"time"
)

type Command interface {
	Process(game Game) ([]Event, error)
}

type PlayersShuffler interface {
	Shuffle(players []Player) []Color
}

type RandomPlayersShuffler struct{}

func NewRandomPlayersShuffler() RandomPlayersShuffler {
	return RandomPlayersShuffler{}
}

// todo true random
func (RandomPlayersShuffler) Shuffle(players []Player) []Color {
	var shuffledColors []Color
	for _, player := range players {
		shuffledColors = append(shuffledColors, player.color)
	}

	return shuffledColors
}

type BoardGenerator interface {
	GenerateBoard() Board
}

type RandomBoardGenerator struct{}

func NewRandomBoardGenerator() RandomBoardGenerator {
	return RandomBoardGenerator{}
}

func (RandomBoardGenerator) GenerateBoard() Board {
	return NewBoardWithOffsetCoord(
		map[grid.HexCoord]Hex{},
	) // todo
}

type StartGameCommand struct {
	occurred        time.Time
	players         []Player
	playersShuffler PlayersShuffler
	boardGenerator  BoardGenerator
	// todo more settings
}

func NewStartGameCommand(
	occurred time.Time,
	players []Player,
	playersShuffler PlayersShuffler,
	boardGenerator BoardGenerator,
) *StartGameCommand {
	return &StartGameCommand{
		occurred:        occurred,
		players:         players,
		playersShuffler: playersShuffler,
		boardGenerator:  boardGenerator,
	}
}

func (command StartGameCommand) Process(game Game) ([]Event, error) {
	// check conditions to start the game
	if game.IsStarted() {
		return nil, GameAlreadyStartedErr
	}

	gameStartedEvent := GameStartedEvent{
		eventCommon: eventCommon{
			Occurred: command.occurred,
		},
		players: command.players,
	}

	boardGeneratedEvent := BoardGeneratedEvent{
		eventCommon: eventCommon{
			Occurred: command.occurred,
		},
		newBoard: command.boardGenerator.GenerateBoard(),
	}

	playersShuffledEvent := PlayersShuffledEvent{
		eventCommon: eventCommon{
			Occurred: command.occurred,
		},
		playersInOrder: command.playersShuffler.Shuffle(command.players),
	}

	phaseOneStartedEvent := InitialSetupPhaseStartedEvent{eventCommon{command.occurred}}

	playerStartedHisTurnEvent := PlayerStartedHisTurnEvent{
		eventCommon: eventCommon{
			Occurred: command.occurred,
		},
		playerColor: playersShuffledEvent.playersInOrder[0],
	}

	return []Event{
		gameStartedEvent,
		boardGeneratedEvent,
		playersShuffledEvent,
		phaseOneStartedEvent,
		playerStartedHisTurnEvent,
	}, nil
}

type BuildSettlementCommand struct {
	occurred          time.Time
	playerColor       Color
	intersectionCoord grid.IntersectionCoord
	settlement        Settlement
}

func NewBuildSettlementCommand(
	occurred time.Time,
	playerColor Color,
	intersectionCoord grid.IntersectionCoord,
	settlement Settlement,
) BuildSettlementCommand {
	return BuildSettlementCommand{
		occurred:          occurred,
		playerColor:       playerColor,
		intersectionCoord: intersectionCoord,
		settlement:        settlement,
	}
}

func (command BuildSettlementCommand) Process(game Game) ([]Event, error) {
	var events []Event

	player, err := game.Player(command.playerColor)
	if err != nil {
		return nil, err
	}

	if err := player.CanBuildSettlement(); err != nil {
		return nil, err
	}

	if err := game.Board().CanBuildSettlementOrCity(command.intersectionCoord, command.settlement); err != nil {
		return nil, err
	}

	playerBuiltSettlementEvent := PlayerBuiltSettlementEvent{
		eventCommon: eventCommon{
			Occurred: command.occurred,
		},
		playerColor:       command.playerColor,
		intersectionCoord: command.intersectionCoord,
		settlement:        command.settlement,
	}

	events = append(events, playerBuiltSettlementEvent)

	if game.InStatus(GameStatusInitialSetup) {
		game.availableCommands = []Command{BuildRoadCommand{}}
	}

	return events, nil
}

type BuildRoadCommand struct {
	occurred    time.Time
	playerColor Color
	pathCoord   grid.PathCoord
	road        road
}

func (command BuildRoadCommand) Process(game Game) ([]Event, error) {
	var events []Event

	player, err := game.Player(command.playerColor)
	if err != nil {
		return nil, err
	}

	// check if the player has available road
	if err := player.HasAvailableRoad(); err != nil {
		return nil, err
	}

	if game.InStatus(GameStatusPlay) {
		if err := player.CanBuy(command.road); err != nil {
			return nil, err
		}
	}

	// check if the board allowing to build
	if err := game.Board().CanBuildRoad(command.pathCoord, command.road); err != nil {
		return nil, err
	}

	playerBuiltRoadEvent := PlayerBuiltRoadEvent{
		eventCommon: eventCommon{
			Occurred: command.occurred,
		},
		playerColor: command.playerColor,
		pathCoord:   command.pathCoord,
		road:        command.road,
	}

	events = append(events, playerBuiltRoadEvent)

	if game.InStatus(GameStatusInitialSetup) {
		playerFinishedHisTurnEvent := PlayerFinishedHisTurnEvent{
			eventCommon: eventCommon{
				Occurred: command.occurred,
			},
			playerColor: command.playerColor,
		}

		events = append(events, playerFinishedHisTurnEvent)
	}

	return events, nil
}
