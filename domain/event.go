package domain

import (
	"github.com/rannoch/catan/grid"
	"time"
)

type eventCommon struct {
	Occurred time.Time
}

func (e eventCommon) withIncreasedVersion(game Game) Game {
	game.version++
	return game
}

type Event interface {
	Apply(game Game) Game
}

// todo before game started events

/// In-game events
type GameStartedEvent struct {
	eventCommon
	players []Player
	// game settings
}

func (event GameStartedEvent) Apply(game Game) Game {
	game = event.withIncreasedVersion(game)

	game = game.WithStatus(GameStatusStarted)

	game = game.AddPlayers(event.players...)

	return game
}

type BoardGeneratedEvent struct {
	eventCommon
	newBoard Board
}

func (event BoardGeneratedEvent) Apply(game Game) Game {
	game = event.withIncreasedVersion(game)

	game = game.WithBoard(event.newBoard)

	return game
}

type PlayersShuffledEvent struct {
	eventCommon
	playersInOrder []Color
}

func (event PlayersShuffledEvent) Apply(game Game) Game {
	game = event.withIncreasedVersion(game)

	game = game.WithTurnOrder(event.playersInOrder)
	return game
}

type InitialSetupPhaseStartedEvent struct {
	eventCommon
}

func (event InitialSetupPhaseStartedEvent) Apply(game Game) Game {
	game = event.withIncreasedVersion(game)

	return game.WithStatus(GameStatusInitialSetup)
}

type PlayPhaseStartedEvent struct {
	eventCommon
}

func (event PlayPhaseStartedEvent) Apply(game Game) Game {
	game = event.withIncreasedVersion(game)

	return game.WithStatus(GameStatusPlay)
}

type PlayerRolledDiceEvent struct {
	eventCommon
	roll roll
}

func (event PlayerRolledDiceEvent) Apply(game Game) Game {
	game = event.withIncreasedVersion(game)

	game.rollHistory = append(game.rollHistory, event.roll)
	return game
}

type PlayerPickedResourcesEvent struct {
	eventCommon
	playerColor     Color
	pickedResources []resource
}

func (event PlayerPickedResourcesEvent) Apply(game Game) Game {
	game = event.withIncreasedVersion(game)

	player, err := game.Player(event.playerColor)
	if err != nil {
		panic(err) // todo
	}

	player = player.WithGainedResources(event.pickedResources)

	game, err = game.WithUpdatedPlayer(player)
	if err != nil {
		panic(err) // todo
	}

	return game
}

type PlayerWasRobbedByRobberEvent struct {
	eventCommon
	robbedPlayerColor Color
	dumpedResources   []resource
}

func (event PlayerWasRobbedByRobberEvent) Apply(game Game) Game {
	game = event.withIncreasedVersion(game)

	panic("implement me")
}

type PlayerWasRobbedByPlayerEvent struct {
	eventCommon
	robbingPlayerColor Color
	robbedPlayerColor  Color
	dumpedResources    []resource
}

func (event PlayerWasRobbedByPlayerEvent) Apply(game Game) Game {
	game = event.withIncreasedVersion(game)

	panic("implement me")
}

type PlayerFinishedHisTurnEvent struct {
	eventCommon
	playerColor Color
}

func (event PlayerFinishedHisTurnEvent) Apply(game Game) Game {
	game = event.withIncreasedVersion(game)

	if game.InStatus(GameStatusPlay) {
		game.totalTurns++
	}

	game.currentTurn = None
	return game
}

type PlayerStartedHisTurnEvent struct {
	eventCommon
	playerColor Color
}

func (event PlayerStartedHisTurnEvent) Apply(game Game) Game {
	game = event.withIncreasedVersion(game)

	game.currentTurn = event.playerColor
	return game
}

type PlayerBuiltSettlementEvent struct {
	eventCommon
	playerColor       Color
	intersectionCoord grid.IntersectionCoord
	settlement        Settlement
}

func (event PlayerBuiltSettlementEvent) Apply(game Game) Game {
	game = event.withIncreasedVersion(game)

	player, err := game.Player(event.playerColor)
	if err != nil {
		panic(err)
	}

	player.victoryPoints += event.settlement.VictoryPoints()
	player.availableSettlements--

	game, err = game.WithUpdatedPlayer(player)
	if err != nil {
		panic(err)
	}

	game = game.WithBoard(game.Board().BuildSettlementOrCity(event.intersectionCoord, event.settlement))

	return game
}

type PlayerBuiltRoadEvent struct {
	eventCommon
	playerColor Color
	pathCoord   grid.PathCoord
	road        road
}

func (event PlayerBuiltRoadEvent) Apply(game Game) Game {
	game = event.withIncreasedVersion(game)

	player, err := game.Player(event.playerColor)
	if err != nil {
		panic(err)
	}

	player.availableRoads--

	game, err = game.WithUpdatedPlayer(player)
	if err != nil {
		panic(err)
	}

	game = game.WithBoard(game.Board().BuildRoad(event.pathCoord, event.road))

	return game
}
