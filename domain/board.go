package domain

import (
	"errors"
	"github.com/rannoch/catan/grid"
)

type Board interface {
	Intersection(intersectionCoord grid.IntersectionCoord) (Intersection, bool)

	Hex(hexCoord grid.HexCoord) (Hex, bool)

	Path(pathCoord grid.PathCoord) (Path, bool)

	HexesByNumberToken(numberToken int64) []Hex

	CanBuildSettlementOrCity(intersectionCoord grid.IntersectionCoord, building Building) error

	BuildSettlementOrCity(intersectionCoord grid.IntersectionCoord, building Building) Board

	CanBuildRoad(pathCoord grid.PathCoord, road road) error

	BuildRoad(pathCoord grid.PathCoord, road road) Board

	LongestRoad(playerColor Color) int64

	GetResourcesByRoll(roll int64) map[Color][]resource
}

var ErrBadIntersection = errors.New("bad intersection coord")
var ErrIntersectionAlreadyHasObject = errors.New("intersectionCoord already has object")

type BoardWithOffsetCoord struct {
	gridCalculator grid.HexagonGridWithOffsetCoordsCalculator

	hexes         map[grid.HexCoord]Hex
	intersections map[grid.IntersectionCoord]Intersection
	paths         map[grid.PathCoord]Path
}

func NewBoardWithOffsetCoord(
	hexes map[grid.HexCoord]Hex,
) BoardWithOffsetCoord {
	boardWithOffsetCoord := BoardWithOffsetCoord{
		hexes: hexes,
	}

	boardWithOffsetCoord.intersections = make(map[grid.IntersectionCoord]Intersection)
	boardWithOffsetCoord.paths = make(map[grid.PathCoord]Path)

	// calculate intersections and paths coords from hexes
	for hexCoord := range hexes {
		adjacentIntersectionCoords := boardWithOffsetCoord.gridCalculator.HexAdjacentIntersections(hexCoord)

		for _, intersectionCoord := range adjacentIntersectionCoords {
			boardWithOffsetCoord.intersections[intersectionCoord] = Intersection{}
		}

		adjacentPathCoords := boardWithOffsetCoord.gridCalculator.HexAdjacentPaths(hexCoord)

		for _, pathCoord := range adjacentPathCoords {
			boardWithOffsetCoord.paths[pathCoord] = Path{}
		}
	}

	return boardWithOffsetCoord
}

// interface compliance
var _ Board = BoardWithOffsetCoord{}

func (board BoardWithOffsetCoord) Intersection(intersectionCoord grid.IntersectionCoord) (Intersection, bool) {
	intersection, exists := board.intersections[intersectionCoord]
	return intersection, exists
}

func (board BoardWithOffsetCoord) Hex(hexCoord grid.HexCoord) (Hex, bool) {
	hex, exists := board.hexes[hexCoord]
	return hex, exists
}

func (board BoardWithOffsetCoord) Path(pathCoord grid.PathCoord) (Path, bool) {
	path, exists := board.paths[pathCoord]
	return path, exists
}

func (board BoardWithOffsetCoord) HexesByNumberToken(roll int64) []Hex {
	panic("implement me")
}

func (board BoardWithOffsetCoord) CanBuildSettlementOrCity(intersectionCoord grid.IntersectionCoord, building Building) error {
	intersection, exists := board.Intersection(intersectionCoord)
	if !exists {
		return ErrBadIntersection
	}

	// todo city check
	if intersection.Building != nil {
		return ErrIntersectionAlreadyHasObject
	}

	// todo add distance check
	return nil
}

func (board BoardWithOffsetCoord) BuildSettlementOrCity(intersectionCoord grid.IntersectionCoord, building Building) Board {
	intersection, exists := board.Intersection(intersectionCoord)
	if !exists {
		panic("todo") // todo
	}

	intersection.Building = building

	board.intersections[intersectionCoord] = intersection

	return board
}

func (board BoardWithOffsetCoord) CanBuildRoad(edgeCoord grid.PathCoord, road road) error {
	panic("implement me")
}

func (board BoardWithOffsetCoord) BuildRoad(edgeCoord grid.PathCoord, road road) Board {
	panic("implement me")
}

func (board BoardWithOffsetCoord) LongestRoad(playerColor Color) int64 {
	panic("implement me")
}

func (board BoardWithOffsetCoord) GetResourcesByRoll(roll int64) map[Color][]resource {
	panic("implement me")
}

// settlement, city, or knight in future
type Building interface {
	Color() Color
	VictoryPoints() int64

	ResourceCount() int64
}

type Settlement struct {
	color Color
}

func NewSettlement(color Color) Settlement {
	return Settlement{color: color}
}

var (
	_ Building = Settlement{}
	_ Buyable  = Settlement{}
)

func (s Settlement) Color() Color {
	return s.color
}

func (Settlement) VictoryPoints() int64 {
	return 1
}

func (Settlement) ResourceCount() int64 {
	return 1
}

func (s Settlement) Cost() []resource {
	return []resource{Wood, Brick, Sheep, Wheat}
}

type City struct {
	color Color
}

func NewCity(color Color) City {
	return City{color: color}
}

var (
	_ Building = City{}
	_ Buyable  = City{}
)

func (c City) Color() Color {
	return c.color
}

func (City) VictoryPoints() int64 {
	return 2
}

func (City) ResourceCount() int64 {
	return 2
}

func (c City) Cost() []resource {
	return []resource{Ore, Ore, Ore, Wheat, Wheat}
}

type road struct {
	color Color
}

var (
	_ Buyable = road{}
)

func (r road) Cost() []resource {
	return []resource{Wood, Brick}
}

var WaterHex = Hex{
	Type:     HexTypeWater,
	Resource: EmptyResource,
}

type Hex struct {
	NumberToken int64
	Type        hexType
	Resource    resource
}

type hexType string

const (
	HexTypeResource hexType = "resource"
	HexTypeDesert   hexType = "desert"
	HexTypeWater    hexType = "water"
	HexTypeEmpty    hexType = "empty"
)

type Intersection struct {
	Port     interface{} // todo port
	Building Building
}

type Path struct {
	Port interface{} // todo port
	Road road
}
