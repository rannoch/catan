package domain

import (
	"errors"

	"github.com/rannoch/catan/grid"
)

// Board board interface
type Board interface {
	grid.HexagonGridCalculator

	Intersection(intersectionCoord grid.IntersectionCoord) (Intersection, bool)

	Hex(hexCoord grid.HexCoord) (Hex, bool)

	Path(pathCoord grid.PathCoord) (Path, bool)

	HexesByNumberToken(numberToken int64) []Hex

	PlaceSettlementOrCity(building Building)

	BuildRoad(pathCoord grid.PathCoord, road Road)

	LongestRoad(playerColor Color) int64

	GetResourcesByRoll(roll int64) map[Color][]ResourceCard

	IntersectionInitialResources(intersectionCoord grid.IntersectionCoord) []ResourceCard
}

var (
	// BadIntersectionCoordErr is used when intersection is not is the board
	BadIntersectionCoordErr = errors.New("bad intersection coord")
	// IntersectionAlreadyHasObjectErr is used when intersection already has an object
	IntersectionAlreadyHasObjectErr = errors.New("intersectionCoord already has an object")
	// BadPathCoordErr is used when path in not the board
	BadPathCoordErr = errors.New("bad path coord")
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
			boardWithOffsetCoord.intersections[intersectionCoord] = Intersection{}
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

func (board BoardWithOffsetCoord) HexesByNumberToken(roll int64) []Hex {
	panic("implement me")
}

func (board *BoardWithOffsetCoord) PlaceSettlementOrCity(building Building) {
	intersection, exists := board.Intersection(building.IntersectionCoord())
	if !exists {
		panic("todo") // todo
	}

	intersection.building = building

	board.intersections[building.IntersectionCoord()] = intersection
}

func (board *BoardWithOffsetCoord) BuildRoad(pathCoord grid.PathCoord, road Road) {
	path, exists := board.Path(pathCoord)
	if !exists {
		panic("todo") // todo
	}

	path.road = &road

	board.paths[pathCoord] = path
}

func (board BoardWithOffsetCoord) LongestRoad(playerColor Color) int64 {
	panic("implement me")
}

func (board BoardWithOffsetCoord) GetResourcesByRoll(roll int64) map[Color][]ResourceCard {
	panic("implement me")
}

func (board BoardWithOffsetCoord) IntersectionInitialResources(intersectionCoord grid.IntersectionCoord) []ResourceCard {
	hexCoords := board.IntersectionAdjacentHexes(intersectionCoord)

	var resources []ResourceCard

	for _, hexCoord := range hexCoords {
		hex, exists := board.Hex(hexCoord)
		if !exists {
			continue
		}

		resources = append(resources, hex.Resource.GetResourceCard(1)...)
	}

	return resources
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
	color Color
}

func NewRoad(color Color) Road {
	return Road{color: color}
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
	NumberToken int64
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
	port     interface{} // todo port
	building Building
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
