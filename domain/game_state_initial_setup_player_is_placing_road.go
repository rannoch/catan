package domain

import (
	"time"

	"github.com/rannoch/catan/grid"
)

type GameStateInitialSetupPlayerIsPlacingRoad struct {
	playerIsToPlaceRoadAdjacentToBuilding grid.IntersectionCoord

	GameState

	game *Game
}

func NewGameStateInitialSetupPlayerIsPlacingRoad(
	playerIsToPlaceRoadAdjacentToBuilding grid.IntersectionCoord,
	gameState GameState,
	game *Game,
) *GameStateInitialSetupPlayerIsPlacingRoad {
	return &GameStateInitialSetupPlayerIsPlacingRoad{
		playerIsToPlaceRoadAdjacentToBuilding: playerIsToPlaceRoadAdjacentToBuilding,
		GameState:                             gameState,
		game:                                  game,
	}
}

func (GameStateInitialSetupPlayerIsPlacingRoad) EnterState(time.Time) {}

func (g *GameStateInitialSetupPlayerIsPlacingRoad) PlaceRoad(playerColor Color, road Road, occurred time.Time) error {
	game := g.game

	if game.CurrentTurn() != playerColor {
		return WrongTurnErr
	}

	_, err := game.Player(playerColor)
	if err != nil {
		return err
	}

	// check if the board allows to build
	if err := g.canBuildRoad(road); err != nil {
		return err
	}

	game.Apply(
		NewEventDescriptor(
			game.Id(),
			PlayerPlacedInitialRoadEvent{
				PlayerColor: playerColor,
				Road:        road,
			},
			nil,
			game.version,
			occurred,
		),
		true,
	)

	// back to turn
	game.setState(g.GameState)
	g.GameState.EnterState(occurred)

	return nil
}

func (g *GameStateInitialSetupPlayerIsPlacingRoad) canBuildRoad(road Road) error {
	game := g.game

	path, exists := game.Board().Path(road.PathCoord())
	if !exists {
		return BadPathCoordErr
	}

	if !path.IsEmpty() {
		return BadPathCoordErr
	}

	if !g.isRoadAdjacentToLastSettlement(road.PathCoord()) {
		return CommandIsForbiddenErr
	}

	// check if road is adjacent to existing and doesn't cross the building
	canBuildRoad := false

	adjacentIntersections := game.Board().PathAdjacentIntersections(road.PathCoord())
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

	adjacentPaths := game.Board().PathAdjacentPaths(road.PathCoord())
	for _, adjacentPathCoord := range adjacentPaths {
		adjacentPath, exists := game.Board().Path(adjacentPathCoord)
		if !exists {
			continue
		}

		if !adjacentPath.IsEmpty() {
			continue
		}

		jointIntersectionCoord, found := game.Board().PathsJointIntersection(road.PathCoord(), adjacentPathCoord)
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

func (g *GameStateInitialSetupPlayerIsPlacingRoad) isRoadAdjacentToLastSettlement(pathCoord grid.PathCoord) bool {
	game := g.game

	intersectionAdjacentPaths := game.board.IntersectionAdjacentPaths(g.playerIsToPlaceRoadAdjacentToBuilding)

	for _, adjacentPathCoord := range intersectionAdjacentPaths {
		if adjacentPathCoord == pathCoord {
			return true
		}
	}

	return false
}
