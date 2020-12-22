package domain_test

import (
	"sort"
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
			grid.HexCoord{R: 0, C: 0}: {NumberToken: 10, Type: domain.HexTypeResource, Resource: domain.Ore},
			grid.HexCoord{R: 0, C: 1}: {NumberToken: 2, Type: domain.HexTypeResource, Resource: domain.Sheep},
			grid.HexCoord{R: 0, C: 2}: {NumberToken: 9, Type: domain.HexTypeResource, Resource: domain.Wood},
			grid.HexCoord{R: 1, C: 0}: {NumberToken: 12, Type: domain.HexTypeResource, Resource: domain.Wheat},
			grid.HexCoord{R: 1, C: 1}: {NumberToken: 6, Type: domain.HexTypeResource, Resource: domain.Brick},
			grid.HexCoord{R: 1, C: 2}: {NumberToken: 4, Type: domain.HexTypeResource, Resource: domain.Sheep},
			grid.HexCoord{R: 1, C: 3}: {NumberToken: 10, Type: domain.HexTypeResource, Resource: domain.Brick},
			grid.HexCoord{R: 2, C: 0}: {NumberToken: 9, Type: domain.HexTypeResource, Resource: domain.Wheat},
			grid.HexCoord{R: 2, C: 1}: {NumberToken: 11, Type: domain.HexTypeResource, Resource: domain.Wood},
			grid.HexCoord{R: 2, C: 2}: {NumberToken: 0, Type: domain.HexTypeDesert, Resource: domain.EmptyResource},
			grid.HexCoord{R: 2, C: 3}: {NumberToken: 3, Type: domain.HexTypeResource, Resource: domain.Wood},
			grid.HexCoord{R: 2, C: 4}: {NumberToken: 8, Type: domain.HexTypeResource, Resource: domain.Ore},
			grid.HexCoord{R: 3, C: 1}: {NumberToken: 8, Type: domain.HexTypeResource, Resource: domain.Wood},
			grid.HexCoord{R: 3, C: 2}: {NumberToken: 3, Type: domain.HexTypeResource, Resource: domain.Ore},
			grid.HexCoord{R: 3, C: 3}: {NumberToken: 4, Type: domain.HexTypeResource, Resource: domain.Wheat},
			grid.HexCoord{R: 3, C: 4}: {NumberToken: 5, Type: domain.HexTypeResource, Resource: domain.Sheep},
			grid.HexCoord{R: 4, C: 2}: {NumberToken: 5, Type: domain.HexTypeResource, Resource: domain.Brick},
			grid.HexCoord{R: 4, C: 3}: {NumberToken: 6, Type: domain.HexTypeResource, Resource: domain.Wheat},
			grid.HexCoord{R: 4, C: 4}: {NumberToken: 11, Type: domain.HexTypeResource, Resource: domain.Sheep},
		},
	)
}

type alphaBetPlayersShuffler struct{}

func (alphaBetPlayersShuffler) Shuffle(players []domain.Player) []domain.Color {
	var colors []domain.Color

	for _, player := range players {
		colors = append(colors, player.Color())
	}

	sort.Slice(colors, func(i, j int) bool {
		return colors[i] < colors[j]
	})

	return colors
}

var _ = Describe("Catan initial state", func() {
	var (
		game    *domain.Game
		err     error
		players = []domain.Player{
			domain.NewPlayer(domain.Blue, "baska"),
			domain.NewPlayer(domain.Red, "masha"),
			domain.NewPlayer(domain.Yellow, "vasya"),
			domain.NewPlayer(domain.Green, "bot"),
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
		Expect(game.SetPlayersShuffler(alphaBetPlayersShuffler{})).To(BeNil())

		Expect(game.StartGame(time.Now())).To(BeNil())
	})

	It("game should enter set-up phase", func() {
		Expect(err).To(BeNil())
		Expect(game.Id()).To(Equal(domain.GameId("test_id")))
		Expect(game.InState(&domain.GameStateInitialSetup{})).To(BeTrue())
		Expect(game.TurnOrder()).To(Equal([]domain.Color{
			domain.Blue,
			domain.Green,
			domain.Red,
			domain.Yellow,
			domain.Yellow,
			domain.Red,
			domain.Green,
			domain.Blue,
		}))
		Expect(game.Version()).To(Equal(int64(9)))
		Expect(game.TotalTurns()).To(Equal(int64(0)))
		Expect(game.CurrentTurn()).To(Equal(domain.Blue))
		Expect(game.NextTurnColor()).To(Equal(domain.Green))

		_, err := game.Player(domain.Red)
		Expect(err).To(BeNil())

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
			Expect(game.BuildSettlement(
				game.CurrentTurn(),
				grid.IntersectionCoord{R: 3, C: 1, D: grid.R},
				domain.NewSettlement(game.CurrentTurn()),
				startGameCommandOccurred,
			)).To(BeNil())
		})

		It("board should have the settlement with right color", func() {
			intersection, exists := game.Board().Intersection(grid.IntersectionCoord{R: 3, C: 1, D: grid.R})
			Expect(exists).To(Equal(true))
			Expect(intersection.Building()).To(Equal(domain.NewSettlement(domain.Blue)))
		})

		When("first player builds a legal road", func() {
			BeforeEach(func() {
				err = game.BuildRoad(
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
				Expect(game.CurrentTurn()).To(Equal(domain.Green))
			})
		})

		When("first player tries to build an illegal road", func() {
			It("should receive an error", func() {
				Expect(game.BuildRoad(
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
				Expect(game.BuildSettlement(
					game.CurrentTurn(),
					grid.IntersectionCoord{R: 3, C: 3, D: grid.R},
					domain.NewSettlement(game.CurrentTurn()),
					startGameCommandOccurred,
				)).To(Equal(domain.CommandIsForbiddenErr))
			})
		})
	})

	When("first player tries to build road before building", func() {
		It("should receive an error", func() {
			Expect(game.BuildRoad(
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

	When("first player tries to build city", func() {
		It("should receive an error", func() {

		})
	})

	When("not first player tries to build settlement", func() {
		It("should receive an error", func() {
			err := game.BuildSettlement(
				game.NextTurnColor(),
				grid.IntersectionCoord{R: 3, C: 3, D: grid.R},
				domain.NewSettlement(game.NextTurnColor()),
				startGameCommandOccurred,
			)
			Expect(err).To(Equal(domain.WrongTurnErr))
		})
	})
})
