package domain

import "errors"

type Color string

const (
	None   Color = "none"
	Black  Color = "black"
	Orange Color = "orange"
	Red    Color = "red"
	Blue   Color = "blue"
	White  Color = "white"
	Green  Color = "green"
	Yellow Color = "yellow"
)

var allColors = []Color{Red, Blue, White, Green, Yellow}

type Player struct {
	userId UserId // User aggregate id, extract name and other info using this reference

	color Color

	resources          []ResourceCard
	resourcesTypeCount map[ResourceCard]int64

	availableSettlements int64
	availableCities      int64
	availableRoads       int64

	victoryPoints int64

	longestRoad      int64
	longestRoadOwner bool
	largestArmyOwner bool

	devCardPlayed bool
	devCards      []string // todo
}

func NewPlayer(color Color, userId UserId) Player {
	player := Player{
		color:  color,
		userId: userId,

		availableRoads:       15,
		availableCities:      4,
		availableSettlements: 5,
		resourcesTypeCount:   make(map[ResourceCard]int64),
	}

	return player
}

func (player *Player) SetColor(color Color) {
	player.color = color
}

func (player Player) DevCardPlayed() bool {
	return player.devCardPlayed
}

func (player Player) LargestArmyOwner() bool {
	return player.largestArmyOwner
}

func (player Player) LongestRoadOwner() bool {
	return player.longestRoadOwner
}

func (player Player) LongestRoad() int64 {
	return player.longestRoad
}

func (player Player) VictoryPoints() int64 {
	return player.victoryPoints
}

func (player Player) AvailableRoads() int64 {
	return player.availableRoads
}

func (player Player) AvailableCities() int64 {
	return player.availableCities
}

func (player Player) AvailableSettlements() int64 {
	return player.availableSettlements
}

func (player Player) UserId() UserId {
	return player.userId
}

func (player Player) Color() Color {
	return player.color
}

func (player Player) Resources() []ResourceCard {
	return player.resources
}

func (player *Player) GainResources(resources []ResourceCard) {
	for _, resource := range resources {
		player.resourcesTypeCount[resource]++
	}

	player.resources = append(player.resources, resources...)
}

func (player Player) WithDisposedResources(resources []ResourceCard) Player {
	for _, resource := range resources {
		player.resourcesTypeCount[resource]-- // todo possible below zero case
	}

	player.resources = []ResourceCard{}
	for resource, count := range player.resourcesTypeCount {
		for i := 0; i < int(count); i++ {
			player.resources = append(player.resources, resource)
		}
	}

	return player
}

var (
	OutOfSettlementsErr   = errors.New("out of settlements")
	OutOfCitiesErr        = errors.New("out of cities")
	OutOfRoadsErr         = errors.New("out of roads")
	NotEnoughResourcesErr = errors.New("not enough resources")
)

func (player Player) CanBuildSettlement() error {
	if player.availableSettlements == 0 {
		return OutOfSettlementsErr
	}

	return nil
}

func (player Player) CanBuildCity() error {
	if player.availableCities == 0 {
		return OutOfCitiesErr
	}

	return nil
}

func (player Player) HasAvailableRoad() error {
	if player.availableRoads == 0 {
		return OutOfRoadsErr
	}

	return nil
}

func (player Player) CanBuy(buyable Buyable) error {
	resourcesTypeCount := player.resourcesTypeCount

	for _, resource := range buyable.Cost() {
		if resourcesTypeCount[resource] == 0 {
			return NotEnoughResourcesErr
		}

		resourcesTypeCount[resource]--
	}

	return nil
}

func (player Player) Buy(buyable Buyable) Player {
	return player.WithDisposedResources(buyable.Cost())
}

func (player Player) HasPlacedInitialBuildings() bool {
	return player.availableSettlements == 3
}

func (player Player) HasPlacedInitialBuildingsAndRoads() bool {
	return player.availableSettlements == 3 && player.availableRoads == 13
}
