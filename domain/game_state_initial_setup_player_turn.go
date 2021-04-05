package domain

import (
	"time"
)

type GameStateInitialSetupPlayerTurn struct {
	playerColor Color

	GameState

	game *Game
}

func NewGameStateInitialSetupPlayerTurn(game *Game, playerColor Color, parentState GameState) *GameStateInitialSetupPlayerTurn {
	g := &GameStateInitialSetupPlayerTurn{playerColor: playerColor}
	g.game = game
	g.GameState = parentState

	return g
}

func (g *GameStateInitialSetupPlayerTurn) EnterState(occurred time.Time) {
	game := g.game

	switch event := game.LastEvent().(type) {
	case PlayerPlacedInitialRoadEvent:
		game.Apply(
			NewEventDescriptor(
				game.Id(),
				PlayerFinishedInitialSetupTurn{},
				nil,
				game.version,
				occurred,
			),
			true,
		)

		game.ChangeState(g.GameState, occurred)
	case PlayerPlacedInitialSettlementEvent:
		gameStateInitialSetupPlayerIsPlacingRoad := NewGameStateInitialSetupPlayerIsPlacingRoad(event.Settlement.intersectionCoord, g, g.game)

		game.ChangeState(gameStateInitialSetupPlayerIsPlacingRoad, occurred)
	default:
		gameStateInitialSetupPlayerIsPlacingSettlement := NewGameStateInitialSetupPlayerIsPlacingSettlement(g, g.game)

		game.ChangeState(gameStateInitialSetupPlayerIsPlacingSettlement, occurred)
	}

	return
}
