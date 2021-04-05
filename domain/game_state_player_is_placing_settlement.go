package domain

import (
	"time"
)

type GameStatePlayerIsPlacingSettlement struct {
	game *Game



	GameStateDefault
}

func NewGameStatePlayerIsToPlaceSettlement(game *Game) *GameStatePlayerIsPlacingSettlement {
	return &GameStatePlayerIsPlacingSettlement{game: game}
}

func (gameStatePlayerIsToPlaceSettlement *GameStatePlayerIsPlacingSettlement) PlaceSettlement(playerColor Color, settlement Settlement, occurred time.Time) error {
	game := gameStatePlayerIsToPlaceSettlement.game

	if game.CurrentTurn() != playerColor {
		return WrongTurnErr
	}

	if err := gameStatePlayerIsToPlaceSettlement.canBuildSettlement(settlement); err != nil {
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

func (gameStatePlayerIsToPlaceSettlement *GameStatePlayerIsPlacingSettlement) canBuildSettlement(settlement Settlement) error {
	game := gameStatePlayerIsToPlaceSettlement.game

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

func (gameStatePlayerIsToPlaceSettlement *GameStatePlayerIsPlacingSettlement) Apply(eventMessage EventMessage, _ bool) {
	game := gameStatePlayerIsToPlaceSettlement.game

	switch event := eventMessage.Event().(type) {
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
	}
}
