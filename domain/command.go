package domain

import (
	"time"

	"github.com/rannoch/catan/grid"
)

type PlayersShuffler interface {
	Shuffle(playerColors []Color) []Color
}

type RandomPlayersShuffler struct{}

func NewRandomPlayersShuffler() RandomPlayersShuffler {
	return RandomPlayersShuffler{}
}

// todo true random
func (RandomPlayersShuffler) Shuffle(playerColors []Color) []Color {
	return playerColors
}

type BoardGenerator interface {
	GenerateBoard() Board
}

type RandomBoardGenerator struct{}

func NewRandomBoardGenerator() RandomBoardGenerator {
	return RandomBoardGenerator{}
}

func (RandomBoardGenerator) GenerateBoard() Board {
	return NewBoardWithOffsetCoord(
		map[grid.HexCoord]Hex{},
	) // todo
}

type StartGameCommand struct {
	occurred        time.Time
	players         []Player
	playersShuffler PlayersShuffler
	boardGenerator  BoardGenerator
	// todo more settings
}

type BuildSettlementCommand struct {
	occurred          time.Time
	playerColor       Color
	intersectionCoord grid.IntersectionCoord
	settlement        Settlement
}

type BuildRoadCommand struct {
	occurred    time.Time
	playerColor Color
	pathCoord   grid.PathCoord
	road        Road
}
