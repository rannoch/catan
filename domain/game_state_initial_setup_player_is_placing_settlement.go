package domain

import (
	"time"
)

type GameStateInitialSetupPlayerIsPlacingSettlement struct {
	GameState

	game *Game
}

func NewGameStateInitialSetupPlayerIsPlacingSettlement(gameState GameState, game *Game) *GameStateInitialSetupPlayerIsPlacingSettlement {
	return &GameStateInitialSetupPlayerIsPlacingSettlement{GameState: gameState, game: game}
}

func (g *GameStateInitialSetupPlayerIsPlacingSettlement) EnterState(time.Time) {}

func (g *GameStateInitialSetupPlayerIsPlacingSettlement) PlaceSettlement(playerColor Color, settlement Settlement, occurred time.Time) error {
	game := g.game

	if game.CurrentTurn() != playerColor {
		return WrongTurnErr
	}

	if err := g.canBuildSettlement(settlement); err != nil {
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

func (g *GameStateInitialSetupPlayerIsPlacingSettlement) canBuildSettlement(settlement Settlement) error {
	game := g.game

	intersection, exists := game.Board().Intersection(settlement.IntersectionCoord())
	if !exists {
		return BadIntersectionCoordErr
	}

	if !intersection.IsEmpty() {
		return IntersectionAlreadyHasObjectErr
	}

	// distance check
	adjacentIntersectionsCoords := game.Board().IntersectionAdjacentIntersections(settlement.intersectionCoord)
	for _, adjacentIntersectionCoord := range adjacentIntersectionsCoords {
		adjacentIntersection, exists := game.Board().Intersection(adjacentIntersectionCoord)
		if !exists {
			continue
		}

		if !adjacentIntersection.IsEmpty() {
			return CommandIsForbiddenErr
		}
	}

	return nil
}
