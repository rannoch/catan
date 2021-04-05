package catan_championship_premium_13_BUGGED_Semi_Final

import (
	"github.com/rannoch/catan/domain"
	"github.com/rannoch/catan/grid"
)

// Catan Championship Premium #13 - BUGGED Semi-Final
// game from https://www.youtube.com/watch?v=RoOFJk87_2E&ab_channel=RetiredHero
// todo
var rollHistory = []domain.Roll{
	domain.NewRoll(domain.D6Roll1, domain.D6Roll6),
}

type diceRoller struct {
	index int64
}

func (d diceRoller) Roll() domain.Roll {
	roll := rollHistory[d.index]
	d.index++
	return roll
}

var Events = []interface{}{
	domain.GameCreated{GameId: "Catan Championship Premium #13 - BUGGED Semi-Final"},

	domain.PlayerJoinedTheGameEvent{Player: domain.NewPlayer(domain.Black, "Kuuchi")},
	domain.PlayerJoinedTheGameEvent{Player: domain.NewPlayer(domain.Orange, "billyDESTROY")},
	domain.PlayerJoinedTheGameEvent{Player: domain.NewPlayer(domain.Blue, "buckass")},
	domain.PlayerJoinedTheGameEvent{Player: domain.NewPlayer(domain.Green, "HakMatata")},

	domain.BoardGeneratorSelectedEvent{},
	domain.PlayersShufflerSelectedEvent{},
	domain.DiceRollerSelected{DiceRoller: diceRoller{}},

	domain.GameStartedEvent{},

	domain.BoardGeneratedEvent{NewBoard: domain.NewBoardWithOffsetCoord(
		map[grid.HexCoord]domain.Hex{
			{R: 0, C: 0}: {Type: domain.HexTypeDesert, Resource: domain.EmptyResource},
			{R: 0, C: 1}: {NumberToken: 8, Type: domain.HexTypeResource, Resource: domain.Wheat},
			{R: 0, C: 2}: {NumberToken: 4, Type: domain.HexTypeResource, Resource: domain.Wheat},
			{R: 1, C: 0}: {NumberToken: 5, Type: domain.HexTypeResource, Resource: domain.Brick},
			{R: 1, C: 1}: {NumberToken: 10, Type: domain.HexTypeResource, Resource: domain.Brick},
			{R: 1, C: 2}: {NumberToken: 3, Type: domain.HexTypeResource, Resource: domain.Ore},
			{R: 1, C: 3}: {NumberToken: 11, Type: domain.HexTypeResource, Resource: domain.Wheat},
			{R: 2, C: 0}: {NumberToken: 2, Type: domain.HexTypeResource, Resource: domain.Wood},
			{R: 2, C: 1}: {NumberToken: 9, Type: domain.HexTypeResource, Resource: domain.Wood},
			{R: 2, C: 2}: {NumberToken: 11, Type: domain.HexTypeDesert, Resource: domain.Wood},
			{R: 2, C: 3}: {NumberToken: 6, Type: domain.HexTypeResource, Resource: domain.Brick},
			{R: 2, C: 4}: {NumberToken: 12, Type: domain.HexTypeResource, Resource: domain.Sheep},
			{R: 3, C: 1}: {NumberToken: 6, Type: domain.HexTypeResource, Resource: domain.Wood},
			{R: 3, C: 2}: {NumberToken: 4, Type: domain.HexTypeResource, Resource: domain.Ore},
			{R: 3, C: 3}: {NumberToken: 5, Type: domain.HexTypeResource, Resource: domain.Ore},
			{R: 3, C: 4}: {NumberToken: 9, Type: domain.HexTypeResource, Resource: domain.Wheat},
			{R: 4, C: 2}: {NumberToken: 3, Type: domain.HexTypeResource, Resource: domain.Sheep},
			{R: 4, C: 3}: {NumberToken: 8, Type: domain.HexTypeResource, Resource: domain.Sheep},
			{R: 4, C: 4}: {NumberToken: 10, Type: domain.HexTypeResource, Resource: domain.Sheep},
		},
	)},

	domain.PlayersShuffledEvent{
		PlayersInOrder: []domain.Color{
			domain.Black,
			domain.Orange,
			domain.Blue,
			domain.Green},
	},

	domain.InitialSetupPhaseStartedEvent{},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Black},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.Black,
		Settlement:  domain.NewSettlement(domain.Black, grid.IntersectionCoord{R: 3, C: 4, D: grid.L}),
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.Black,
		Road:        domain.NewRoad(grid.PathCoord{R: 3, C: 3, D: grid.E}, domain.Black),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Black},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Orange},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.Orange,
		Settlement:  domain.NewSettlement(domain.Orange, grid.IntersectionCoord{R: 3, C: 2, D: grid.R}),
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.Orange,
		Road:        domain.NewRoad(grid.PathCoord{R: 4, C: 3, D: grid.N}, domain.Orange),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Orange},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Blue},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.Blue,
		Settlement:  domain.NewSettlement(domain.Blue, grid.IntersectionCoord{R: 3, C: 2, D: grid.L}),
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.Blue,
		Road:        domain.NewRoad(grid.PathCoord{R: 3, C: 1, D: grid.E}, domain.Blue),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Blue},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Green},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.Green,
		Settlement:  domain.NewSettlement(domain.Green, grid.IntersectionCoord{R: 0, C: 1, D: grid.R}),
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.Green,
		Road:        domain.NewRoad(grid.PathCoord{R: 0, C: 1, D: grid.E}, domain.Green),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Green},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Green},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.Green,
		Settlement:  domain.NewSettlement(domain.Green, grid.IntersectionCoord{R: 1, C: 0, D: grid.R}),
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.Green,
		Road:        domain.NewRoad(grid.PathCoord{R: 2, C: 1, D: grid.W}, domain.Green),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Green},

	domain.PlayerPickedResourcesEvent{
		PlayerColor:     domain.Green,
		PickedResources: []domain.ResourceCard{domain.ResourceCardBrick, domain.ResourceCardBrick, domain.ResourceCardWood},
	},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Blue},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.Blue,
		Settlement:  domain.NewSettlement(domain.Blue, grid.IntersectionCoord{R: 1, C: 2, D: grid.R}),
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.Blue,
		Road:        domain.NewRoad(grid.PathCoord{R: 1, C: 2, D: grid.E}, domain.Blue),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Blue},

	domain.PlayerPickedResourcesEvent{
		PlayerColor:     domain.Blue,
		PickedResources: []domain.ResourceCard{domain.ResourceCardOre, domain.ResourceCardWheat, domain.ResourceCardBrick},
	},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Orange},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.Orange,
		Settlement:  domain.NewSettlement(domain.Orange, grid.IntersectionCoord{R: 1, C: 1, D: grid.R}),
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.Orange,
		Road:        domain.NewRoad(grid.PathCoord{R: 1, C: 1, D: grid.E}, domain.Orange),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Orange},

	domain.PlayerPickedResourcesEvent{
		PlayerColor:     domain.Orange,
		PickedResources: []domain.ResourceCard{domain.ResourceCardBrick, domain.ResourceCardOre, domain.ResourceCardWood},
	},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Black},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.Black,
		Settlement:  domain.NewSettlement(domain.Black, grid.IntersectionCoord{R: 3, C: 1, D: grid.L}),
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.Black,
		Road:        domain.NewRoad(grid.PathCoord{R: 3, C: 0, D: grid.E}, domain.Black),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Black},

	domain.PlayerPickedResourcesEvent{
		PlayerColor:     domain.Black,
		PickedResources: []domain.ResourceCard{domain.ResourceCardWood, domain.ResourceCardWood},
	},

	domain.PlayPhaseStartedEvent{},
	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Black},

	domain.PlayerRolledDiceEvent{Roll: domain.NewRoll(domain.D6Roll1, domain.D6Roll6)},
}
