package domain

import (
	"errors"

	"github.com/rannoch/catan/grid"
)

// Board board interface
type Board interface {
	grid.HexagonGridCalculator

	Intersections() []Intersection

	Paths() []Path

	Hexes() []Hex

	Intersection(intersectionCoord grid.IntersectionCoord) (Intersection, bool)

	Path(pathCoord grid.PathCoord) (Path, bool)

	Hex(hexCoord grid.HexCoord) (Hex, bool)

	UpdateIntersection(intersectionCoord grid.IntersectionCoord, intersection Intersection) error

	UpdatePath(pathCoord grid.PathCoord, path Path) error

	UpdateHex(hexCoord grid.HexCoord, hex Hex) error

	HexesByNumberToken(numberToken int64) []Hex
}

var (
	// BadIntersectionCoordErr is used when intersection is not is the board
	BadIntersectionCoordErr = errors.New("bad intersection coord")
	// IntersectionAlreadyHasObjectErr is used when intersection already has an object
	IntersectionAlreadyHasObjectErr = errors.New("intersectionCoord already has an object")
	// BadPathCoordErr is used when path in not the board
	BadPathCoordErr = errors.New("bad path coord")
	// BadPathCoordErr is used when path in not the board
	BadHexCoordErr = errors.New("bad hex coord")
)

type BoardWithOffsetCoord struct {
	grid.HexagonGridWithOffsetCoordsCalculator

	hexes         map[grid.HexCoord]Hex
	intersections map[grid.IntersectionCoord]Intersection
	paths         map[grid.PathCoord]Path
}

func NewBoardWithOffsetCoord(
	hexes map[grid.HexCoord]Hex,
) *BoardWithOffsetCoord {
	boardWithOffsetCoord := &BoardWithOffsetCoord{
		hexes: hexes,
	}

	boardWithOffsetCoord.intersections = make(map[grid.IntersectionCoord]Intersection)
	boardWithOffsetCoord.paths = make(map[grid.PathCoord]Path)

	// calculate intersections and paths coords from hexes
	for hexCoord := range hexes {
		adjacentIntersectionCoords := boardWithOffsetCoord.HexAdjacentIntersections(hexCoord)

		for _, intersectionCoord := range adjacentIntersectionCoords {
			boardWithOffsetCoord.intersections[intersectionCoord] = NewIntersection(intersectionCoord)
		}

		adjacentPathCoords := boardWithOffsetCoord.HexAdjacentPaths(hexCoord)

		for _, pathCoord := range adjacentPathCoords {
			boardWithOffsetCoord.paths[pathCoord] = Path{}
		}
	}

	return boardWithOffsetCoord
}

// interface compliance
var _ Board = (*BoardWithOffsetCoord)(nil)

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

func (board BoardWithOffsetCoord) Intersections() []Intersection {
	intersections := make([]Intersection, 0, len(board.intersections))

	for _, intersection := range board.intersections {
		intersections = append(intersections, intersection)
	}

	return intersections
}

func (board BoardWithOffsetCoord) Paths() []Path {
	paths := make([]Path, 0, len(board.intersections))

	for _, path := range board.paths {
		paths = append(paths, path)
	}

	return paths
}

func (board BoardWithOffsetCoord) Hexes() []Hex {
	hexes := make([]Hex, 0, len(board.hexes))

	for _, hex := range board.hexes {
		hexes = append(hexes, hex)
	}

	return hexes
}

func (board BoardWithOffsetCoord) UpdateIntersection(intersectionCoord grid.IntersectionCoord, intersection Intersection) error {
	_, exists := board.intersections[intersectionCoord]
	if !exists {
		return BadIntersectionCoordErr
	}

	board.intersections[intersectionCoord] = intersection
	return nil
}

func (board BoardWithOffsetCoord) UpdatePath(pathCoord grid.PathCoord, path Path) error {
	_, exists := board.paths[pathCoord]
	if !exists {
		return BadPathCoordErr
	}

	board.paths[pathCoord] = path
	return nil
}

func (board BoardWithOffsetCoord) UpdateHex(hexCoord grid.HexCoord, hex Hex) error {
	_, exists := board.hexes[hexCoord]
	if !exists {
		return BadHexCoordErr
	}

	board.hexes[hexCoord] = hex
	return nil
}

func (board BoardWithOffsetCoord) HexesByNumberToken(roll int64) []Hex {
	panic("implement me")
}

// settlement, city, or knight in future
type Building interface {
	IntersectionCoord() grid.IntersectionCoord

	Color() Color
	VictoryPoints() int64

	ResourceCount() int64
}

type Settlement struct {
	color             Color
	intersectionCoord grid.IntersectionCoord
}

func NewSettlement(color Color, intersectionCoord grid.IntersectionCoord) Settlement {
	return Settlement{color: color, intersectionCoord: intersectionCoord}
}

var (
	_ Building = Settlement{}
	_ Buyable  = Settlement{}
)

func (s Settlement) IntersectionCoord() grid.IntersectionCoord {
	return s.intersectionCoord
}

func (s Settlement) Color() Color {
	return s.color
}

func (Settlement) VictoryPoints() int64 {
	return 1
}

func (Settlement) ResourceCount() int64 {
	return 1
}

func (s Settlement) Cost() []ResourceCard {
	return []ResourceCard{ResourceCardWood, ResourceCardBrick, ResourceCardSheep, ResourceCardWheat}
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

func (c City) IntersectionCoord() grid.IntersectionCoord {
	panic("t") // todo add coord to city
}

func (c City) Color() Color {
	return c.color
}

func (City) VictoryPoints() int64 {
	return 2
}

func (City) ResourceCount() int64 {
	return 2
}

func (c City) Cost() []ResourceCard {
	return []ResourceCard{ResourceCardOre, ResourceCardOre, ResourceCardOre, ResourceCardWheat, ResourceCardWheat}
}

type Road struct {
	coord grid.PathCoord
	color Color
}

func NewRoad(coord grid.PathCoord, color Color) Road {
	return Road{coord: coord, color: color}
}

func (r Road) PathCoord() grid.PathCoord {
	return r.coord
}

var (
	_ Buyable = Road{}
)

func (r Road) Cost() []ResourceCard {
	return []ResourceCard{ResourceCardWood, ResourceCardBrick}
}

var WaterHex = Hex{
	Type:     HexTypeWater,
	Resource: EmptyResource,
}

type Hex struct {
	Coord       grid.HexCoord
	NumberToken NumberToken
	Type        hexType
	Resource    Resource
}

type hexType string

const (
	HexTypeResource hexType = "resource"
	HexTypeDesert   hexType = "desert"
	HexTypeWater    hexType = "water"
	HexTypeEmpty    hexType = "empty"
)

type Intersection struct {
	coord    grid.IntersectionCoord
	port     interface{} // todo port
	building Building
}

func NewIntersection(coord grid.IntersectionCoord) Intersection {
	return Intersection{coord: coord}
}

func (intersection Intersection) Building() Building {
	return intersection.building
}

func (intersection *Intersection) SetBuilding(building Building) {
	intersection.building = building
}

func (intersection Intersection) IsEmpty() bool {
	return intersection.building == nil
}

type Path struct {
	port interface{} // todo port
	road *Road
}

func (path Path) Road() *Road {
	return path.road
}

func (path *Path) SetRoad(road *Road) {
	path.road = road
}

func (path Path) IsEmpty() bool {
	return path.road == nil
}
