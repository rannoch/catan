package domain

import (
	"time"
)

type gameStatePlayerIsRollingDice struct {
	game *Game

	GameStateDefault
}

func NewGameStatePlayerIsRollingDice(game *Game) *gameStatePlayerIsRollingDice {
	return &gameStatePlayerIsRollingDice{game: game}
}

func (g *gameStatePlayerIsRollingDice) RollDice(playerColor Color, occurred time.Time) error {
	// todo
	//
	return nil
}

func (g *gameStatePlayerIsRollingDice) Apply(eventMessage EventMessage, isNew bool) {
	panic("implement me")
}
