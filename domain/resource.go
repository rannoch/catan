package domain

// ResourceCard
// is used by players in building, buying development cards, trading etc
type ResourceCard struct {
	resource Resource
}

// Resource represents resource types
type Resource string

const (
	// EmptyResource is used for desert hex, water etc
	EmptyResource Resource = "empty"

	// Ore
	Ore Resource = "ore"
	// Wheat
	Wheat Resource = "wheat"
	// Sheep
	Sheep Resource = "sheep"
	// Brick
	Brick Resource = "brick"
	// Wood
	Wood Resource = "wood"
)

// GetResourceCard todo move to hex method?
func (resource Resource) GetResourceCard(count int64) []ResourceCard {
	if resource == EmptyResource {
		return nil
	}

	var resourceCards []ResourceCard

	for i := 0; i < int(count); i++ {
		resourceCards = append(resourceCards, ResourceCard{resource: resource})
	}

	return resourceCards
}

var (
	ResourceCardOre   = ResourceCard{resource: Ore}
	ResourceCardWheat = ResourceCard{resource: Wheat}
	ResourceCardSheep = ResourceCard{resource: Sheep}
	ResourceCardBrick = ResourceCard{resource: Brick}
	ResourceCardWood  = ResourceCard{resource: Wood}
)

// Buyable staff can be bought by ResourceCard
type Buyable interface {
	Cost() []ResourceCard
}
