package domain

import (
	"errors"
	"time"

	"github.com/rannoch/catan/grid"
)

var CommandIsForbiddenErr = errors.New("command is forbidden")

type GameState interface {
	EnterState(occurred time.Time)

	SetBoardGenerator(boardGenerator BoardGenerator) error
	SetPlayersShuffler(playersShuffler PlayersShuffler) error

	GenerateBoard(occurred time.Time) error
	ShufflePlayers(occurred time.Time) error

	AddPlayer(player Player, occurred time.Time) error
	RemovePlayer(player Player, occurred time.Time) error

	StartGame(occurred time.Time) error

	RollDice(playerColor Color, occurred time.Time) error

	BuyRoad(playerColor Color, occurred time.Time) error

	BuySettlement(playerColor Color, occurred time.Time) error

	BuyCity(playerColor Color, occurred time.Time) error

	PlaceSettlement(playerColor Color, settlement Settlement, occurred time.Time) error

	PlaceRoad(
		playerColor Color,
		pathCoord grid.PathCoord,
		road Road,
		occurred time.Time,
	) error

	BuyDevelopmentCard(playerColor Color, card DevelopmentCard) error

	TurnOrder() []Color
	EndTurn(playerColor Color, occurred time.Time) error
	CurrentTurn() Color
	Apply(eventMessage EventMessage, isNew bool)
}

type GameStateDefault struct{}

var _ GameState = GameStateDefault{}

func (GameStateDefault) EnterState(_ time.Time) {}

func (GameStateDefault) AddPlayer(Player, time.Time) error {
	return CommandIsForbiddenErr
}

func (GameStateDefault) RemovePlayer(Player, time.Time) error {
	return CommandIsForbiddenErr
}

func (d GameStateDefault) SetBoardGenerator(BoardGenerator) error {
	return CommandIsForbiddenErr
}

func (d GameStateDefault) SetPlayersShuffler(PlayersShuffler) error {
	return CommandIsForbiddenErr
}

func (d GameStateDefault) GenerateBoard(time.Time) error {
	return CommandIsForbiddenErr
}

func (d GameStateDefault) ShufflePlayers(time.Time) error {
	return CommandIsForbiddenErr
}

func (GameStateDefault) StartGame(time.Time) error {
	return CommandIsForbiddenErr
}

func (GameStateDefault) RollDice(Color, time.Time) error {
	return CommandIsForbiddenErr
}

func (GameStateDefault) BuyRoad(Color, time.Time) error {
	return CommandIsForbiddenErr
}

func (GameStateDefault) BuySettlement(Color, time.Time) error {
	return CommandIsForbiddenErr
}

func (GameStateDefault) BuyCity(Color, time.Time) error {
	return CommandIsForbiddenErr
}

func (GameStateDefault) PlaceSettlement(Color, Settlement, time.Time) error {
	return CommandIsForbiddenErr
}

func (GameStateDefault) PlaceRoad(Color, grid.PathCoord, Road, time.Time) error {
	return CommandIsForbiddenErr
}

func (GameStateDefault) BuyDevelopmentCard(Color, DevelopmentCard) error {
	return CommandIsForbiddenErr
}

func (GameStateDefault) TurnOrder() []Color {
	return nil
}

func (GameStateDefault) EndTurn(Color, time.Time) error {
	return CommandIsForbiddenErr
}

func (GameStateDefault) CurrentTurn() Color {
	return None
}

func (GameStateDefault) Apply(EventMessage, bool) {

}
