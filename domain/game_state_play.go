package domain

import (
	"time"

	"github.com/rannoch/catan/grid"
)

type GameStatePlay struct {
	GameStateDefault
	game *Game
}

func NewGameStatePlay(game *Game) *GameStatePlay {
	return &GameStatePlay{game: game}
}

var _ GameState = (*GameStatePlay)(nil)

func (gameStatePlay *GameStatePlay) StartGame(time.Time) error {
	return GameAlreadyStartedErr
}

func (gameStatePlay *GameStatePlay) EnterState(occurred time.Time) {
	game := gameStatePlay.game

	playerStartedHisTurnEventMessage := NewEventDescriptor(
		game.Id(),
		PlayerStartedHisTurnEvent{
			PlayerColor: game.TurnOrder()[0],
		},
		nil,
		game.Version(),
		occurred,
	)

	game.Apply(playerStartedHisTurnEventMessage, true)
}

func (gameStatePlay *GameStatePlay) PlaceSettlement(playerColor Color, settlement Settlement, occurred time.Time) error {
	game := gameStatePlay.game

	player, err := game.Player(playerColor)
	if err != nil {
		return PlayerNotExistsErr
	}

	if game.CurrentTurn() != playerColor {
		return WrongTurnErr
	}

	if err := player.CanBuy(settlement); err != nil {
		return err
	}

	if err := player.CanBuildSettlement(); err != nil {
		return err
	}

	playerBuiltSettlementEventMessage := NewEventDescriptor(
		game.Id(),
		PlayerPlacedSettlementEvent{
			PlayerColor: playerColor,
			Settlement:  settlement,
		},
		nil,
		game.Version(),
		occurred,
	)

	game.Apply(playerBuiltSettlementEventMessage, true)

	return nil
}

func (gameStatePlay *GameStatePlay) PlaceRoad(playerColor Color, road Road, occurred time.Time) error {
	game := gameStatePlay.game

	if game.CurrentTurn() != playerColor {
		return WrongTurnErr
	}

	player, err := game.Player(playerColor)
	if err != nil {
		return err
	}

	// check if the player has available road
	if err := player.HasAvailableRoad(); err != nil {
		return err
	}

	// check if the player has enough resources to buy a road
	if err := player.CanBuy(road); err != nil {
		return err
	}

	// check if the board allowing to build
	if err := gameStatePlay.canBuildRoad(road.PathCoord(), road); err != nil {
		return err
	}

	gameStatePlay.Apply(
		NewEventDescriptor(game.Id(), PlayerPlacedRoadEvent{
			PlayerColor: playerColor,
			Road:        road,
		}, nil, game.version, occurred,
		),
		true,
	)

	return nil
}

func (gameStatePlay *GameStatePlay) canBuildRoad(pathCoord grid.PathCoord, road Road) error {
	game := gameStatePlay.game

	path, exists := game.Board().Path(pathCoord)
	if !exists {
		return BadPathCoordErr
	}

	if !path.IsEmpty() {
		return BadPathCoordErr
	}

	// check if road is adjacent to existing and doesn't cross the building
	canBuildRoad := false

	adjacentIntersections := game.Board().PathAdjacentIntersections(pathCoord)
	for _, adjacentIntersectionCoord := range adjacentIntersections {
		intersection, exists := game.Board().Intersection(adjacentIntersectionCoord)
		if !exists {
			continue
		}

		if intersection.building == nil {
			continue
		}

		if intersection.building.Color() == road.color {
			canBuildRoad = true
			break
		}
	}

	if canBuildRoad {
		return nil
	}

	adjacentPaths := game.Board().PathAdjacentPaths(pathCoord)
	for _, adjacentPathCoord := range adjacentPaths {
		adjacentPath, exists := game.Board().Path(adjacentPathCoord)
		if !exists {
			continue
		}

		if !adjacentPath.IsEmpty() {
			continue
		}

		jointIntersectionCoord, found := game.Board().PathsJointIntersection(pathCoord, adjacentPathCoord)
		if !found {
			continue
		}

		intersection, exists := game.Board().Intersection(jointIntersectionCoord)
		if exists {
			continue
		}

		if !intersection.IsEmpty() && intersection.building.Color() != road.color {
			continue
		}

		canBuildRoad = true
		break
	}

	if canBuildRoad {
		return nil
	}

	return CommandIsForbiddenErr
}

func (gameStatePlay *GameStatePlay) BuyDevelopmentCard(playerColor Color, card DevelopmentCard) error {
	panic("implement me")
}

func (gameStatePlay *GameStatePlay) EndTurn(playerColor Color, occurred time.Time) error {
	panic("implement me")
}

func (gameStatePlay *GameStatePlay) CurrentTurn() Color {
	panic("implement me")
}

func (gameStatePlay *GameStatePlay) TurnOrder() []Color {
	return gameStatePlay.game.turnOrder
}

func (gameStatePlay *GameStatePlay) Apply(eventMessage EventMessage, isNew bool) {
	game := gameStatePlay.game

	switch event := eventMessage.Event().(type) {
	case PlayerStartedHisTurnEvent:
		game.setCurrentTurn(event.PlayerColor)
	case PlayerFinishedHisTurnEvent:
		game.incrementTotalTurns()
		game.setCurrentTurn(None)

		// todo invoke next player start his turn
	case PlayerPlacedSettlementEvent:
		player, err := game.Player(event.PlayerColor)
		if err != nil {
			panic(err)
		}

		player.victoryPoints += event.Settlement.VictoryPoints()
		player.availableSettlements--

		err = game.updatePlayer(player)
		if err != nil {
			panic(err)
		}

		err = game.placeSettlement(event.Settlement)
		if err != nil {
			panic(err)
		}
	case PlayerPlacedRoadEvent:
		player, err := game.Player(event.PlayerColor)
		if err != nil {
			panic(err)
		}

		player.availableRoads--

		err = game.updatePlayer(player)
		if err != nil {
			panic(err)
		}

		//game.Board().BuildRoad(event.PathCoord, event.Road)
	case PlayerWasRobbedByRobberEvent:
		// todo
	case PlayerWasRobbedByPlayerEvent:
	case PlayerPickedResourcesEvent:
		player, err := game.Player(event.PlayerColor)
		if err != nil {
			panic(err) // todo
		}

		player.GainResources(event.PickedResources)
		err = game.updatePlayer(player)
		if err != nil {
			panic(err)
		}
	case PlayerRolledDiceEvent:
		game.rollHistory = append(game.rollHistory, event.Roll)
	}
}
