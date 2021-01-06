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
	domain.PlayerBuiltSettlementEvent{
		PlayerColor:       domain.Blue,
		IntersectionCoord: grid.IntersectionCoord{R: 3, C: 3, D: grid.R},
		Settlement:        domain.NewSettlement(domain.Blue),
	},
	domain.PlayerBuiltRoadEvent{
		PlayerColor: domain.Blue,
		PathCoord:   grid.PathCoord{R: 3, C: 3, D: grid.E},
		Road:        domain.NewRoad(domain.Blue),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Blue},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.White},
	domain.PlayerBuiltSettlementEvent{
		PlayerColor:       domain.White,
		IntersectionCoord: grid.IntersectionCoord{R: 2, C: 3, D: grid.R},
		Settlement:        domain.NewSettlement(domain.White),
	},
	domain.PlayerBuiltRoadEvent{
		PlayerColor: domain.White,
		PathCoord:   grid.PathCoord{R: 2, C: 3, D: grid.E},
		Road:        domain.NewRoad(domain.White),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.White},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Red},
	domain.PlayerBuiltSettlementEvent{
		PlayerColor:       domain.Red,
		IntersectionCoord: grid.IntersectionCoord{R: 0, C: 0, D: grid.R},
		Settlement:        domain.NewSettlement(domain.Red),
	},
	domain.PlayerBuiltRoadEvent{
		PlayerColor: domain.Red,
		PathCoord:   grid.PathCoord{R: 1, C: 1, D: grid.N},
		Road:        domain.NewRoad(domain.Red),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Red},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Yellow},
	domain.PlayerBuiltSettlementEvent{
		PlayerColor:       domain.Yellow,
		IntersectionCoord: grid.IntersectionCoord{R: 1, C: 3, D: grid.L},
		Settlement:        domain.NewSettlement(domain.Yellow),
	},
	domain.PlayerBuiltRoadEvent{
		PlayerColor: domain.Yellow,
		PathCoord:   grid.PathCoord{R: 1, C: 2, D: grid.N},
		Road:        domain.NewRoad(domain.Yellow),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Yellow},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Yellow},
	domain.PlayerBuiltSettlementEvent{
		PlayerColor:       domain.Yellow,
		IntersectionCoord: grid.IntersectionCoord{R: 3, C: 2, D: grid.R},
		Settlement:        domain.NewSettlement(domain.Yellow),
	},
	domain.PlayerPickedResourcesEvent{
		PlayerColor:     domain.Yellow,
		PickedResources: []domain.ResourceCard{domain.ResourceCardOre, domain.ResourceCardWheat, domain.ResourceCardWheat},
	},
	domain.PlayerBuiltRoadEvent{
		PlayerColor: domain.Yellow,
		PathCoord:   grid.PathCoord{R: 4, C: 3, D: grid.N},
		Road:        domain.NewRoad(domain.Yellow),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Yellow},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Red},
	domain.PlayerBuiltSettlementEvent{
		PlayerColor:       domain.Red,
		IntersectionCoord: grid.IntersectionCoord{R: 2, C: 0, D: grid.R},
		Settlement:        domain.NewSettlement(domain.Red),
	},
	domain.PlayerPickedResourcesEvent{
		PlayerColor:     domain.Red,
		PickedResources: []domain.ResourceCard{domain.ResourceCardWheat, domain.ResourceCardWood, domain.ResourceCardWood},
	},
	domain.PlayerBuiltRoadEvent{
		PlayerColor: domain.Red,
		PathCoord:   grid.PathCoord{R: 3, C: 1, D: grid.N},
		Road:        domain.NewRoad(domain.Red),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Red},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.White},
	domain.PlayerBuiltSettlementEvent{
		PlayerColor:       domain.White,
		IntersectionCoord: grid.IntersectionCoord{R: 1, C: 0, D: grid.R},
		Settlement:        domain.NewSettlement(domain.White),
	},
	domain.PlayerPickedResourcesEvent{
		PlayerColor:     domain.White,
		PickedResources: []domain.ResourceCard{domain.ResourceCardWheat, domain.ResourceCardBrick, domain.ResourceCardWood},
	},
	domain.PlayerBuiltRoadEvent{
		PlayerColor: domain.White,
		PathCoord:   grid.PathCoord{R: 2, C: 1, D: grid.W},
		Road:        domain.NewRoad(domain.White),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.White},

	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Blue},
	domain.PlayerBuiltSettlementEvent{
		PlayerColor:       domain.Blue,
		IntersectionCoord: grid.IntersectionCoord{R: 3, C: 1, D: grid.R},
		Settlement:        domain.NewSettlement(domain.Blue),
	},
	domain.PlayerPickedResourcesEvent{
		PlayerColor:     domain.Blue,
		PickedResources: []domain.ResourceCard{domain.ResourceCardWood, domain.ResourceCardOre, domain.ResourceCardBrick},
	},
	domain.PlayerBuiltRoadEvent{
		PlayerColor: domain.Blue,
		PathCoord:   grid.PathCoord{R: 4, C: 2, D: grid.N},
		Road:        domain.NewRoad(domain.Blue),
	},
	domain.PlayerFinishedHisTurnEvent{PlayerColor: domain.Blue},

	domain.PlayPhaseStartedEvent{},
	domain.PlayerStartedHisTurnEvent{PlayerColor: domain.Blue},
}
