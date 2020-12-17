package domain

import "errors"

type Color string

const (
	None   Color = "none"
	Red    Color = "red"
	Blue   Color = "blue"
	White  Color = "white"
	Green  Color = "green"
	Yellow Color = "yellow"
)

type Player struct {
	userId UserId // User aggregate id, extract name and other info using this reference

	color Color

	resources          []resource
	resourcesTypeCount map[resource]int64

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

func NewPlayer(color Color, userId UserId) Player {
	player := Player{
		color:  color,
		userId: userId,

		availableRoads:       15,
		availableCities:      4,
		availableSettlements: 5,
	}

	return player
}

func (player Player) UserId() UserId {
	return player.userId
}

func (player Player) Color() Color {
	return player.color
}

func (player Player) Resources() []resource {
	return player.resources
}

func (player Player) WithGainedResources(resources []resource) Player {
	for _, resource := range resources {
		player.resourcesTypeCount[resource]++
	}

	player.resources = append(player.resources, resources...)

	return player
}

func (player Player) WithDisposedResources(resources []resource) Player {
	for _, resource := range resources {
		player.resourcesTypeCount[resource]-- // todo possible below zero case
	}

	player.resources = []resource{}
	for resource, count := range player.resourcesTypeCount {
		for i := 0; i < int(count); i++ {
			player.resources = append(player.resources, resource)
		}
	}

	return player
}

var (
	ErrOutOfSettlements   = errors.New("out of settlements")
	ErrOutOfRoads         = errors.New("out of roads")
	ErrNotEnoughResources = errors.New("not enough resources")
)

func (player Player) CanBuildSettlement() error {
	if player.availableSettlements == 0 {
		return ErrOutOfSettlements
	}

	return nil
}

func (player Player) HasAvailableRoad() error {
	if player.availableRoads == 0 {
		return ErrOutOfRoads
	}

	return nil
}

func (player Player) CanBuy(buyable Buyable) error {
	resourcesTypeCount := player.resourcesTypeCount

	for _, resource := range buyable.Cost() {
		if resourcesTypeCount[resource] == 0 {
			return ErrNotEnoughResources
		}

		resourcesTypeCount[resource]--
	}

	return nil
}

func (player Player) Buy(buyable Buyable) Player {
	return player.WithDisposedResources(buyable.Cost())
}
