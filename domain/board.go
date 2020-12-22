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

	CanBuildRoad(pathCoord grid.PathCoord, road Road) error

	BuildRoad(pathCoord grid.PathCoord, road Road)

	LongestRoad(playerColor Color) int64

	GetResourcesByRoll(roll int64) map[Color][]resource
}

var (
	ErrBadIntersectionCoord         = errors.New("bad intersection coord")
	ErrIntersectionAlreadyHasObject = errors.New("intersectionCoord already has object")
	ErrBadPathCoord                 = errors.New("bad path coord")
)

type BoardWithOffsetCoord struct {
	gridCalculator grid.HexagonGridWithOffsetCoordsCalculator

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

func (board BoardWithOffsetCoord) CanBuildSettlementOrCity(intersectionCoord grid.IntersectionCoord, building Building) error {
	intersection, exists := board.Intersection(intersectionCoord)
	if !exists {
		return ErrBadIntersectionCoord
	}

	// todo city check
	if !intersection.IsEmpty() {
		return ErrIntersectionAlreadyHasObject
	}

	// todo add distance check
	return nil
}

func (board *BoardWithOffsetCoord) BuildSettlementOrCity(intersectionCoord grid.IntersectionCoord, building Building) Board {
	intersection, exists := board.Intersection(intersectionCoord)
	if !exists {
		panic("todo") // todo
	}

	intersection.building = building

	board.intersections[intersectionCoord] = intersection

	return board
}

func (board BoardWithOffsetCoord) CanBuildRoad(pathCoord grid.PathCoord, road Road) error {
	path, exists := board.Path(pathCoord)
	if !exists {
		return ErrBadPathCoord
	}

	if !path.IsEmpty() {
		return ErrBadPathCoord
	}

	// check if road is adjacent to existing and doesn't cross the building
	canBuildRoad := false

	adjacentIntersections := board.gridCalculator.PathAdjacentIntersections(pathCoord)
	for _, adjacentIntersectionCoord := range adjacentIntersections {
		intersection, exists := board.Intersection(adjacentIntersectionCoord)
		if !exists {
			continue
		}

		if intersection.building == nil {
			continue
		}

		if intersection.building.Color() == road.color {
			canBuildRoad = true
			break
		}
	}

	if canBuildRoad {
		return nil
	}

	adjacentPaths := board.gridCalculator.PathAdjacentPaths(pathCoord)
	for _, adjacentPathCoord := range adjacentPaths {
		adjacentPath, exists := board.Path(adjacentPathCoord)
		if !exists {
			continue
		}

		if !adjacentPath.IsEmpty() {
			continue
		}

		jointIntersectionCoord, found := board.gridCalculator.PathsJointIntersection(pathCoord, adjacentPathCoord)
		if !found {
			continue
		}

		intersection, exists := board.Intersection(jointIntersectionCoord)
		if exists {
			continue
		}

		if !intersection.IsEmpty() && intersection.building.Color() != road.color {
			continue
		}

		canBuildRoad = true
		break
	}

	if canBuildRoad {
		return nil
	}

	return CommandIsForbiddenErr
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

type Road struct {
	color Color
}

func NewRoad(color Color) Road {
	return Road{color: color}
}

var (
	_ Buyable = Road{}
)

func (r Road) Cost() []resource {
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
