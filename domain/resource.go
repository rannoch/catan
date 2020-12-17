package domain

//todo resource type
type resource string

const EmptyResource resource = "emptyResource"
const Ore resource = "ore"
const Wheat resource = "wheat"
const Sheep resource = "sheep"
const Brick resource = "brick"
const Wood resource = "wood"

type Buyable interface {
	Cost() []resource
}
