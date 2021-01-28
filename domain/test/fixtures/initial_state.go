package fixtures

import (
	"github.com/rannoch/catan/domain"
	"github.com/rannoch/catan/grid"
)

var JustAfterInitialStateEvents = []interface{}{
	domain.GameCreated{GameId: "test_id"},
	domain.PlayerJoinedTheGameEvent{Player: domain.NewPlayer(domain.Blue, "baska")},
	domain.PlayerJoinedTheGameEvent{Player: domain.NewPlayer(domain.White, "bot")},
	domain.PlayerJoinedTheGameEvent{Player: domain.NewPlayer(domain.Red, "masha")},
	domain.PlayerJoinedTheGameEvent{Player: domain.NewPlayer(domain.Yellow, "vasya")},
	domain.GameStartedEvent{},
	domain.BoardGeneratedEvent{NewBoard: domain.NewBoardWithOffsetCoord(
		map[grid.HexCoord]domain.Hex{
			{R: 0, C: 0}: {NumberToken: 10, Type: domain.HexTypeResource, Resource: domain.Ore},
			{R: 0, C: 1}: {NumberToken: 2, Type: domain.HexTypeResource, Resource: domain.Sheep},
			{R: 0, C: 2}: {NumberToken: 9, Type: domain.HexTypeResource, Resource: domain.Wood},
			{R: 1, C: 0}: {NumberToken: 12, Type: domain.HexTypeResource, Resource: domain.Wheat},
			{R: 1, C: 1}: {NumberToken: 6, Type: domain.HexTypeResource, Resource: domain.Brick},
			{R: 1, C: 2}: {NumberToken: 4, Type: domain.HexTypeResource, Resource: domain.Sheep},
			{R: 1, C: 3}: {NumberToken: 10, Type: domain.HexTypeResource, Resource: domain.Brick},
			{R: 2, C: 0}: {NumberToken: 9, Type: domain.HexTypeResource, Resource: domain.Wheat},
			{R: 2, C: 1}: {NumberToken: 11, Type: domain.HexTypeResource, Resource: domain.Wood},
			{R: 2, C: 2}: {NumberToken: 0, Type: domain.HexTypeDesert, Resource: domain.EmptyResource},
			{R: 2, C: 3}: {NumberToken: 3, Type: domain.HexTypeResource, Resource: domain.Wood},
			{R: 2, C: 4}: {NumberToken: 8, Type: domain.HexTypeResource, Resource: domain.Ore},
			{R: 3, C: 1}: {NumberToken: 8, Type: domain.HexTypeResource, Resource: domain.Wood},
			{R: 3, C: 2}: {NumberToken: 3, Type: domain.HexTypeResource, Resource: domain.Ore},
			{R: 3, C: 3}: {NumberToken: 4, Type: domain.HexTypeResource, Resource: domain.Wheat},
			{R: 3, C: 4}: {NumberToken: 5, Type: domain.HexTypeResource, Resource: domain.Sheep},
			{R: 4, C: 2}: {NumberToken: 5, Type: domain.HexTypeResource, Resource: domain.Brick},
			{R: 4, C: 3}: {NumberToken: 6, Type: domain.HexTypeResource, Resource: domain.Wheat},
			{R: 4, C: 4}: {NumberToken: 11, Type: domain.HexTypeResource, Resource: domain.Sheep},
		},
	)},
	domain.PlayersShuffledEvent{
		PlayersInOrder: []domain.Color{
			domain.Blue,
			domain.White,
			domain.Red,
			domain.Yellow},
	},

	domain.InitialSetupPhaseStartedEvent{},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Blue},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.Blue,
		Settlement:  domain.NewSettlement(domain.Blue, grid.IntersectionCoord{R: 3, C: 3, D: grid.R}),
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.Blue,
		Road:        domain.NewRoad(grid.PathCoord{R: 3, C: 3, D: grid.E}, domain.Blue),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Blue},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.White},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.White,
		Settlement:  domain.NewSettlement(domain.White, grid.IntersectionCoord{R: 2, C: 3, D: grid.R}),
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.White,
		Road:        domain.NewRoad(grid.PathCoord{R: 2, C: 3, D: grid.E}, domain.White),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.White},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Red},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.Red,
		Settlement:  domain.NewSettlement(domain.Red, grid.IntersectionCoord{R: 0, C: 0, D: grid.R}),
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.Red,
		Road:        domain.NewRoad(grid.PathCoord{R: 1, C: 1, D: grid.N}, domain.Red),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Red},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Yellow},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.Yellow,
		Settlement:  domain.NewSettlement(domain.Yellow, grid.IntersectionCoord{R: 1, C: 3, D: grid.L}),
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.Yellow,
		Road:        domain.NewRoad(grid.PathCoord{R: 1, C: 2, D: grid.N}, domain.Yellow),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Yellow},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Yellow},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.Yellow,
		Settlement:  domain.NewSettlement(domain.Yellow, grid.IntersectionCoord{R: 3, C: 2, D: grid.R}),
	},
	domain.PlayerPickedResourcesEvent{
		PlayerColor:     domain.Yellow,
		PickedResources: []domain.ResourceCard{domain.ResourceCardOre, domain.ResourceCardWheat, domain.ResourceCardWheat},
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.Yellow,
		Road:        domain.NewRoad(grid.PathCoord{R: 4, C: 3, D: grid.N}, domain.Yellow),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Yellow},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Red},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.Red,
		Settlement:  domain.NewSettlement(domain.Red, grid.IntersectionCoord{R: 2, C: 0, D: grid.R}),
	},
	domain.PlayerPickedResourcesEvent{
		PlayerColor:     domain.Red,
		PickedResources: []domain.ResourceCard{domain.ResourceCardWheat, domain.ResourceCardWood, domain.ResourceCardWood},
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.Red,
		Road:        domain.NewRoad(grid.PathCoord{R: 3, C: 1, D: grid.N}, domain.Red),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Red},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.White},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.White,
		Settlement:  domain.NewSettlement(domain.White, grid.IntersectionCoord{R: 1, C: 0, D: grid.R}),
	},
	domain.PlayerPickedResourcesEvent{
		PlayerColor:     domain.White,
		PickedResources: []domain.ResourceCard{domain.ResourceCardWheat, domain.ResourceCardBrick, domain.ResourceCardWood},
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.White,
		Road:        domain.NewRoad(grid.PathCoord{R: 2, C: 1, D: grid.W}, domain.White),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.White},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Blue},
	domain.PlayerPlacedSettlementEvent{
		PlayerColor: domain.Blue,
		Settlement:  domain.NewSettlement(domain.Blue, grid.IntersectionCoord{R: 3, C: 1, D: grid.R}),
	},
	domain.PlayerPickedResourcesEvent{
		PlayerColor:     domain.Blue,
		PickedResources: []domain.ResourceCard{domain.ResourceCardWood, domain.ResourceCardOre, domain.ResourceCardBrick},
	},
	domain.PlayerPlacedRoadEvent{
		PlayerColor: domain.Blue,
		Road:        domain.NewRoad(grid.PathCoord{R: 4, C: 2, D: grid.N}, domain.Blue),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Blue},

	domain.PlayPhaseStartedEvent{},
	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Blue},
}
