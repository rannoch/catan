package domain

import "github.com/rannoch/catan/grid"

type GameStatePlayerIsPlacingRobber struct {
	game *Game

	GameStateDefault
}

func (g *GameStatePlayerIsPlacingRobber) PlaceRobber(Color, grid.HexCoord) error {
	return CommandIsForbiddenErr
}

func (g *GameStatePlayerIsPlacingRobber) Apply(eventMessage EventMessage, isNew bool) {
	panic("implement me")
}
