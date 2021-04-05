package domain

import (
	"errors"
	"time"

	"github.com/rannoch/catan/grid"
)

var CommandIsForbiddenErr = errors.New("command is forbidden")

type GameState interface {
	EnterState(occurred time.Time)

	SetBoardGenerator(boardGenerator BoardGenerator, occurred time.Time) error
	SetPlayersShuffler(playersShuffler PlayersShuffler, occurred time.Time) error
	SetDiceRoller(diceRoller DiceRoller, occurred time.Time) error

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

	PlaceRoad(playerColor Color, road Road, occurred time.Time) error

	PlaceRobber(playerColor Color, hexCoord grid.HexCoord) error

	RobPlayer(playerColor Color, targetColor Color) error

	BuyDevelopmentCard(playerColor Color) error

	PlayDevelopmentCard(playerColor Color, card DevelopmentCard) error

	TurnOrder() []Color
	EndTurn(playerColor Color, occurred time.Time) error
	CurrentTurn() Color
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

func (d GameStateDefault) SetBoardGenerator(BoardGenerator, time.Time) error {
	return CommandIsForbiddenErr
}

func (d GameStateDefault) SetPlayersShuffler(PlayersShuffler, time.Time) error {
	return CommandIsForbiddenErr
}

func (d GameStateDefault) SetDiceRoller(DiceRoller, time.Time) error {
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

func (GameStateDefault) PlaceRoad(Color, Road, time.Time) error {
	return CommandIsForbiddenErr
}

func (d GameStateDefault) PlaceRobber(Color, grid.HexCoord) error {
	return CommandIsForbiddenErr
}

func (d GameStateDefault) RobPlayer(Color, Color) error {
	return CommandIsForbiddenErr
}

func (GameStateDefault) BuyDevelopmentCard(Color) error {
	return CommandIsForbiddenErr
}

func (d GameStateDefault) PlayDevelopmentCard(playerColor Color, card DevelopmentCard) error {
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
