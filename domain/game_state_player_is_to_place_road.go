package domain

import (
	"time"

	"github.com/rannoch/catan/grid"
)

type GameStatePlayerIsToPlaceRoad struct {
	game *Game

	GameStateDefault
}

func NewGameStatePlayerIsToPlaceRoad(game *Game) *GameStatePlayerIsToPlaceRoad {
	return &GameStatePlayerIsToPlaceRoad{game: game}
}

func (gameStatePlayerIsToPlaceRoad *GameStatePlayerIsToPlaceRoad) PlaceRoad(playerColor Color, pathCoord grid.PathCoord, road Road, occurred time.Time) error {
	game := gameStatePlayerIsToPlaceRoad.game

	if game.CurrentTurn() != playerColor {
		return WrongTurnErr
	}

	_, err := game.Player(playerColor)
	if err != nil {
		return err
	}

	// check if the board allows to build
	if err := gameStatePlayerIsToPlaceRoad.canBuildRoad(pathCoord, road); err != nil {
		return err
	}

	game.Apply(
		NewEventDescriptor(
			game.Id(),
			PlayerPlacedRoadEvent{
				PlayerColor: playerColor,
				PathCoord:   pathCoord,
				Road:        road,
			},
			nil,
			game.version,
			occurred,
		),
		true,
	)

	return nil
}

func (gameStatePlayerIsToPlaceRoad *GameStatePlayerIsToPlaceRoad) canBuildRoad(pathCoord grid.PathCoord, road Road) error {
	game := gameStatePlayerIsToPlaceRoad.game

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

func (gameStatePlayerIsToPlaceRoad *GameStatePlayerIsToPlaceRoad) Apply(eventMessage EventMessage, _ bool) {
	game := gameStatePlayerIsToPlaceRoad.game

	switch event := eventMessage.Event().(type) {
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

		game.Board().BuildRoad(event.PathCoord, event.Road)
	}
}
