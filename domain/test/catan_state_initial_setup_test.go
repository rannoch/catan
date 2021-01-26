package test_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rannoch/catan/domain"
	"github.com/rannoch/catan/grid"
)

//
//                         _ _
//                       /     \
//                  _ _ / (0,4) \ _ _
//                /     \       /     \
//           _ _ / (0,3) \ _ _ / (1,4) \
//         /     \       /     \       /
//        / (0,2) \ _ _ / (1,3) \ _ _ /
//        \       /     \       /     \
//         \ _ _ / (1,2) \ _ _ / (2,3) \
//         /     \       /     \       /
//        / (1,1) \ _ _ / (2,2) \ _ _ /
//        \       /     \       /
//         \ _ _ /       \ _ _ /
//
//
//
//

type testBoardGenerator struct{}

func (testBoardGenerator) GenerateBoard() domain.Board {
	return domain.NewBoardWithOffsetCoord(
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
	)
}

type simplePlayersShuffler struct{}

func (simplePlayersShuffler) Shuffle(playerColors []domain.Color) []domain.Color {
	return playerColors
}

type gameCommand struct {
	playerColor    domain.Color
	coord          interface{}
	buildingOrRoad interface{}
	occurred       time.Time
}

var _ = Describe("Catan state initial setup", func() {
	var (
		game    *domain.Game
		err     error
		players = []domain.Player{
			domain.NewPlayer(domain.Blue, "baska"),
			domain.NewPlayer(domain.White, "bot"),
			domain.NewPlayer(domain.Red, "masha"),
			domain.NewPlayer(domain.Yellow, "vasya"),
		}
		startGameCommandOccurred = time.Now()
	)

	BeforeEach(func() {
		game = domain.NewGame("test_id", time.Now())

		// add players
		for _, player := range players {
			Expect(game.AddPlayer(player, time.Now())).To(BeNil())
		}
		// set board generator
		Expect(game.SetBoardGenerator(testBoardGenerator{})).To(BeNil())
		// set players shuffler
		Expect(game.SetPlayersShuffler(simplePlayersShuffler{})).To(BeNil())

		Expect(game.StartGame(time.Now())).To(BeNil())
	})

	It("aggregate id should match", func() {
		Expect(game.Id()).To(Equal("test_id"))
	})

	Specify("turn order is right", func() {
		Expect(game.TurnOrder()).To(Equal([]domain.Color{
			domain.Blue,
			domain.White,
			domain.Red,
			domain.Yellow,
			domain.Yellow,
			domain.Red,
			domain.White,
			domain.Blue,
		}))
	})

	Specify("turns order", func() {
		Expect(game.CurrentTurn()).To(Equal(domain.Blue))
		Expect(game.NextTurnColor()).To(Equal(domain.White))
		Expect(game.TotalTurns()).To(Equal(int64(0)))
	})

	It("should have correct version", func() {
		Expect(game.Version()).To(Equal(int64(10)))
	})

	Specify("no error", func() {
		Expect(err).To(BeNil())
	})

	It("game should enter set-up phase", func() {
		Expect(game.InState(&domain.GameStateInitialSetup{})).To(BeTrue())
	})

	Specify("player is added to the game", func() {
		_, err := game.Player(domain.Red)
		Expect(err).To(BeNil())
	})

	Specify("players available buildings are correct", func() {
		for _, player := range game.Players() {
			Expect(player.AvailableSettlements()).To(Equal(int64(5)))
			Expect(player.AvailableCities()).To(Equal(int64(4)))
			Expect(player.AvailableRoads()).To(Equal(int64(15)))
		}
		Expect(len(game.Players())).To(Equal(len(players)))
	})

	It("intersection (3,1,R) should exist and for example (3213,1,R) not ", func() {
		intersection, exists := game.Board().Intersection(grid.IntersectionCoord{R: 3, C: 1, D: grid.R})
		Expect(exists).To(Equal(true))
		Expect(intersection.Building()).To(BeNil())

		intersection, exists = game.Board().Intersection(grid.IntersectionCoord{R: 3213, C: 1, D: grid.R})
		Expect(exists).To(Equal(false))
		Expect(intersection.Building()).To(BeNil())
	})

	When("first player builds an initial settlement", func() {
		BeforeEach(func() {
			Expect(game.PlaceSettlement(game.CurrentTurn(), domain.NewSettlement(game.CurrentTurn(), grid.IntersectionCoord{R: 3, C: 1, D: grid.R}), startGameCommandOccurred)).To(BeNil())
		})

		It("board should have the settlement with right color", func() {
			intersection, exists := game.Board().Intersection(grid.IntersectionCoord{R: 3, C: 1, D: grid.R})
			Expect(exists).To(Equal(true))
			Expect(intersection.Building()).To(Equal(domain.NewSettlement(domain.Blue, grid.IntersectionCoord{R: 3, C: 1, D: grid.R})))
		})

		When("first player builds a legal road", func() {
			BeforeEach(func() {
				err = game.PlaceRoad(
					game.CurrentTurn(),
					grid.PathCoord{R: 3, C: 1, D: grid.E},
					domain.NewRoad(game.CurrentTurn()),
					time.Now(),
				)
			})

			It("board should have the road with right color ", func() {
				Expect(err).To(BeNil())
			})
			It("turn passes to next player", func() {
				Expect(game.CurrentTurn()).To(Equal(domain.White))
			})
		})

		When("first player tries to build an illegal road", func() {
			It("should receive an error", func() {
				Expect(game.PlaceRoad(
					game.CurrentTurn(),
					grid.PathCoord{R: 1, C: 2, D: grid.E},
					domain.NewRoad(game.CurrentTurn()),
					time.Now(),
				)).To(Equal(domain.CommandIsForbiddenErr))
			})
		})

		When("first player tries to end turn before placing a road", func() {
			It("should receive an error", func() {
				Expect(game.EndTurn(game.CurrentTurn(), time.Now())).
					To(Equal(domain.CommandIsForbiddenErr))
			})
		})

		When("first player tries to build second settlement", func() {
			It("should receive an error", func() {
				Expect(game.PlaceSettlement(game.CurrentTurn(), domain.NewSettlement(game.CurrentTurn(), grid.IntersectionCoord{R: 3, C: 3, D: grid.R}), startGameCommandOccurred)).To(Equal(domain.CommandIsForbiddenErr))
			})
		})
	})

	When("current player tries to buy a settlement", func() {
		It("should receive an error", func() {
			Expect(game.BuySettlement(game.CurrentTurn(), time.Now())).To(Equal(domain.CommandIsForbiddenErr))
		})
	})

	When("current player tries to buy a road", func() {
		It("should receive an error", func() {
			Expect(game.BuyRoad(game.CurrentTurn(), time.Now())).To(Equal(domain.CommandIsForbiddenErr))
		})
	})

	When("current player tries to buy a city", func() {
		It("should receive an error", func() {
			Expect(game.BuyCity(game.CurrentTurn(), time.Now())).To(Equal(domain.CommandIsForbiddenErr))
		})
	})

	When("current player tries to buy a development card", func() {
		It("should receive an error", func() {
			Expect(game.BuyDevelopmentCard(game.CurrentTurn(), domain.DevelopmentCard{})).To(Equal(domain.CommandIsForbiddenErr))
		})
	})

	When("current player tries to roll a dice", func() {
		It("should receive an error", func() {
			Expect(game.RollDice(game.CurrentTurn(), time.Now())).To(Equal(domain.CommandIsForbiddenErr))
		})
	})

	When("first player tries to place road before building", func() {
		It("should receive an error", func() {
			Expect(game.PlaceRoad(
				game.CurrentTurn(),
				grid.PathCoord{R: 3, C: 1, D: grid.W},
				domain.NewRoad(game.CurrentTurn()),
				time.Now(),
			)).To(Equal(domain.CommandIsForbiddenErr))
		})
	})

	When("first player tries to end turn before placing a building and a road", func() {
		It("should receive an error", func() {
			Expect(game.EndTurn(game.CurrentTurn(), time.Now())).
				To(Equal(domain.CommandIsForbiddenErr))
		})
	})

	When("first player tries to place city", func() {
		It("should receive an error", func() {
		})
	})

	When("not first player tries to build settlement", func() {
		It("should receive an error", func() {
			err := game.PlaceSettlement(game.NextTurnColor(), domain.NewSettlement(game.NextTurnColor(), grid.IntersectionCoord{R: 3, C: 3, D: grid.R}), startGameCommandOccurred)
			Expect(err).To(Equal(domain.WrongTurnErr))
		})
	})

	Context("last player is to place second road", func() {
		var gameCommands []gameCommand

		BeforeEach(func() {
			gameCommands = []gameCommand{
				{
					playerColor:    domain.Blue,
					buildingOrRoad: domain.NewSettlement(domain.Blue, grid.IntersectionCoord{R: 3, C: 3, D: grid.R}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Blue,
					coord:          grid.PathCoord{R: 3, C: 3, D: grid.E},
					buildingOrRoad: domain.NewRoad(domain.Blue),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.White,
					buildingOrRoad: domain.NewSettlement(domain.White, grid.IntersectionCoord{R: 2, C: 3, D: grid.R}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.White,
					coord:          grid.PathCoord{R: 2, C: 3, D: grid.E},
					buildingOrRoad: domain.NewRoad(domain.White),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Red,
					buildingOrRoad: domain.NewSettlement(domain.Red, grid.IntersectionCoord{R: 0, C: 0, D: grid.R}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Red,
					coord:          grid.PathCoord{R: 1, C: 1, D: grid.N},
					buildingOrRoad: domain.NewRoad(domain.Red),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Yellow,
					buildingOrRoad: domain.NewSettlement(domain.Yellow, grid.IntersectionCoord{R: 1, C: 3, D: grid.L}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Yellow,
					coord:          grid.PathCoord{R: 1, C: 2, D: grid.N},
					buildingOrRoad: domain.NewRoad(domain.Yellow),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Yellow,
					buildingOrRoad: domain.NewSettlement(domain.Yellow, grid.IntersectionCoord{R: 3, C: 2, D: grid.R}),
					occurred:       startGameCommandOccurred,
				},
			}
		})

		JustBeforeEach(func() {
			for _, gameCommand := range gameCommands {
				switch gameCommand.buildingOrRoad.(type) {
				case domain.Settlement:
					err := game.PlaceSettlement(
						gameCommand.playerColor,
						gameCommand.buildingOrRoad.(domain.Settlement),
						gameCommand.occurred,
					)
					Expect(err).To(BeNil())
				case domain.Road:
					err := game.PlaceRoad(
						gameCommand.playerColor,
						gameCommand.coord.(grid.PathCoord),
						gameCommand.buildingOrRoad.(domain.Road),
						gameCommand.occurred,
					)
					Expect(err).To(BeNil())
				}
			}
		})

		When("he tries to build a road connected to first settlement", func() {
			It("should receive an error", func() {
				Expect(game.PlaceRoad(
					game.CurrentTurn(),
					grid.PathCoord{R: 1, C: 3, D: grid.W},
					domain.NewRoad(game.CurrentTurn()),
					time.Now(),
				)).To(Equal(domain.CommandIsForbiddenErr))
			})
		})
		When("he tries to build a road connected to second settlement", func() {
			It("should be ok", func() {
				Expect(game.PlaceRoad(
					game.CurrentTurn(),
					grid.PathCoord{R: 4, C: 3, D: grid.N},
					domain.NewRoad(game.CurrentTurn()),
					time.Now(),
				)).NotTo(HaveOccurred())
			})
		})
	})

	Context("last player is to place second settlement", func() {
		var gameCommands []gameCommand

		BeforeEach(func() {
			gameCommands = []gameCommand{
				{
					playerColor:    domain.Blue,
					buildingOrRoad: domain.NewSettlement(domain.Blue, grid.IntersectionCoord{R: 3, C: 3, D: grid.R}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Blue,
					coord:          grid.PathCoord{R: 3, C: 3, D: grid.E},
					buildingOrRoad: domain.NewRoad(domain.Blue),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.White,
					buildingOrRoad: domain.NewSettlement(domain.White, grid.IntersectionCoord{R: 2, C: 3, D: grid.R}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.White,
					coord:          grid.PathCoord{R: 2, C: 3, D: grid.E},
					buildingOrRoad: domain.NewRoad(domain.White),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Red,
					buildingOrRoad: domain.NewSettlement(domain.Red, grid.IntersectionCoord{R: 0, C: 0, D: grid.R}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Red,
					coord:          grid.PathCoord{R: 1, C: 1, D: grid.N},
					buildingOrRoad: domain.NewRoad(domain.Red),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Yellow,
					buildingOrRoad: domain.NewSettlement(domain.Yellow, grid.IntersectionCoord{R: 1, C: 3, D: grid.L}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Yellow,
					coord:          grid.PathCoord{R: 1, C: 2, D: grid.N},
					buildingOrRoad: domain.NewRoad(domain.Yellow),
					occurred:       startGameCommandOccurred,
				},
			}
		})

		JustBeforeEach(func() {
			for _, gameCommand := range gameCommands {
				switch gameCommand.buildingOrRoad.(type) {
				case domain.Settlement:
					err := game.PlaceSettlement(
						gameCommand.playerColor,
						gameCommand.buildingOrRoad.(domain.Settlement),
						gameCommand.occurred,
					)
					Expect(err).NotTo(HaveOccurred())
				case domain.Road:
					err := game.PlaceRoad(
						gameCommand.playerColor,
						gameCommand.coord.(grid.PathCoord),
						gameCommand.buildingOrRoad.(domain.Road),
						gameCommand.occurred,
					)
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})

		When("he tries to build a settlement in a legal spot", func() {
			It("should be ok", func() {
				Expect(game.PlaceSettlement(
					game.CurrentTurn(),
					domain.NewSettlement(game.CurrentTurn(), grid.IntersectionCoord{R: 3, C: 2, D: grid.R}),
					time.Now(),
				)).NotTo(HaveOccurred())
			})
		})

		When("he tries to build a settlement that is too close to another settlement", func() {
			It("should receive an error", func() {
				Expect(game.PlaceSettlement(
					game.CurrentTurn(),
					domain.NewSettlement(game.CurrentTurn(), grid.IntersectionCoord{R: 1, C: 1, D: grid.L}),
					time.Now(),
				)).To(Equal(domain.CommandIsForbiddenErr))
			})
		})
	})

	When("all players placed initial buildings", func() {
		var gameCommands []gameCommand

		BeforeEach(func() {
			gameCommands = []gameCommand{
				{
					playerColor:    domain.Blue,
					buildingOrRoad: domain.NewSettlement(domain.Blue, grid.IntersectionCoord{R: 3, C: 3, D: grid.R}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Blue,
					coord:          grid.PathCoord{R: 3, C: 3, D: grid.E},
					buildingOrRoad: domain.NewRoad(domain.Blue),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.White,
					buildingOrRoad: domain.NewSettlement(domain.White, grid.IntersectionCoord{R: 2, C: 3, D: grid.R}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.White,
					coord:          grid.PathCoord{R: 2, C: 3, D: grid.E},
					buildingOrRoad: domain.NewRoad(domain.White),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Red,
					buildingOrRoad: domain.NewSettlement(domain.Red, grid.IntersectionCoord{R: 0, C: 0, D: grid.R}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Red,
					coord:          grid.PathCoord{R: 1, C: 1, D: grid.N},
					buildingOrRoad: domain.NewRoad(domain.Red),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Yellow,
					buildingOrRoad: domain.NewSettlement(domain.Yellow, grid.IntersectionCoord{R: 1, C: 3, D: grid.L}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Yellow,
					coord:          grid.PathCoord{R: 1, C: 2, D: grid.N},
					buildingOrRoad: domain.NewRoad(domain.Yellow),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Yellow,
					buildingOrRoad: domain.NewSettlement(domain.Yellow, grid.IntersectionCoord{R: 3, C: 2, D: grid.R}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Yellow,
					coord:          grid.PathCoord{R: 4, C: 3, D: grid.N},
					buildingOrRoad: domain.NewRoad(domain.Yellow),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Red,
					buildingOrRoad: domain.NewSettlement(domain.Red, grid.IntersectionCoord{R: 2, C: 0, D: grid.R}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Red,
					coord:          grid.PathCoord{R: 3, C: 1, D: grid.N},
					buildingOrRoad: domain.NewRoad(domain.Red),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.White,
					buildingOrRoad: domain.NewSettlement(domain.White, grid.IntersectionCoord{R: 1, C: 0, D: grid.R}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.White,
					coord:          grid.PathCoord{R: 2, C: 1, D: grid.W},
					buildingOrRoad: domain.NewRoad(domain.White),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Blue,
					buildingOrRoad: domain.NewSettlement(domain.Blue, grid.IntersectionCoord{R: 3, C: 1, D: grid.R}),
					occurred:       startGameCommandOccurred,
				},
				{
					playerColor:    domain.Blue,
					coord:          grid.PathCoord{R: 4, C: 2, D: grid.N},
					buildingOrRoad: domain.NewRoad(domain.Blue),
					occurred:       startGameCommandOccurred,
				},
			}
		})

		JustBeforeEach(func() {
			for _, gameCommand := range gameCommands {
				switch gameCommand.buildingOrRoad.(type) {
				case domain.Settlement:
					Expect(game.PlaceSettlement(gameCommand.playerColor, gameCommand.buildingOrRoad.(domain.Settlement), gameCommand.occurred)).To(BeNil())
				case domain.Road:
					Expect(game.PlaceRoad(
						gameCommand.playerColor,
						gameCommand.coord.(grid.PathCoord),
						gameCommand.buildingOrRoad.(domain.Road),
						gameCommand.occurred,
					)).To(BeNil())
				}
			}
		})

		It("should enter play state", func() {
			Expect(game.InState(domain.NewGameStatePlay(game))).To(BeTrue())
		})

		It("players should receive initial resources", func() {
			expectedInitialResources := map[domain.Color][]domain.ResourceCard{
				domain.Blue:   {domain.ResourceCardWood, domain.ResourceCardOre, domain.ResourceCardBrick},
				domain.White:  {domain.ResourceCardWheat, domain.ResourceCardBrick, domain.ResourceCardWood},
				domain.Red:    {domain.ResourceCardWheat, domain.ResourceCardWood, domain.ResourceCardWood},
				domain.Yellow: {domain.ResourceCardOre, domain.ResourceCardWheat, domain.ResourceCardWheat},
			}

			for color, expectedResources := range expectedInitialResources {
				player, err := game.Player(color)
				Expect(err).NotTo(HaveOccurred())
				Expect(player.Resources()).To(Equal(expectedResources))
			}
		})
	})
})
