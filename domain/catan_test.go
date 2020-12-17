package domain_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rannoch/catan/domain"
	"github.com/rannoch/catan/grid"
	"time"
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

var _ = Describe("Catan", func() {

	var (
		game    domain.Game
		err     error
		players = []domain.Player{
			domain.NewPlayer(domain.Blue, "baska"),
			domain.NewPlayer(domain.Red, "masha"),
			domain.NewPlayer(domain.Yellow, "vasya"),
			domain.NewPlayer(domain.Green, "bot"),
		}
	)

	BeforeEach(func() {
		game = domain.NewGame("test_id")
	})

	Describe("enter to initial set-up", func() {
		startGameCommandOccurred := time.Now()

		BeforeEach(func() {
			startGameCommand := domain.NewStartGameCommand(
				startGameCommandOccurred,
				players,
				domain.NewRandomPlayersShuffler(),
				testBoardGenerator{},
			)

			events, err := startGameCommand.Process(game)
			Expect(err).To(BeNil())

			game = game.WithAppliedEvents(events...)
		})

		It("game should enter set-up phase", func() {
			Expect(err).To(BeNil())
			Expect(game.Id()).To(Equal(domain.GameId("test_id")))
			Expect(game.IsStarted()).To(BeTrue())
			Expect(game.Status()).To(Equal(domain.GameStatusInitialSetup))
			Expect(game.TurnOrder()).To(Equal([]domain.Color{
				domain.Blue,
				domain.Red,
				domain.Yellow,
				domain.Green,
				domain.Green,
				domain.Yellow,
				domain.Red,
				domain.Blue,
			}))
			Expect(game.Version()).To(Equal(int64(5)))
			Expect(game.TotalTurns()).To(Equal(int64(0)))
			Expect(game.CurrentTurn()).To(Equal(domain.Blue))
			Expect(game.NextTurnColor()).To(Equal(domain.Red))

			playerRed, err := game.Player(domain.Red)
			Expect(err).To(BeNil())

			Expect(playerRed.AvailableSettlements()).To(Equal(int64(5)))
			Expect(playerRed.AvailableCities()).To(Equal(int64(4)))
			Expect(playerRed.AvailableRoads()).To(Equal(int64(15)))

			//Expect(game.Players).To(Equal(players))
		})

		It("intersection (3,1,R) should exist and for example (3213,1,R) not ", func() {
			intersection, exists := game.Board().Intersection(grid.IntersectionCoord{R: 3, C: 1, D: grid.R})
			Expect(exists).To(Equal(true))
			Expect(intersection.Building).To(BeNil())

			intersection, exists = game.Board().Intersection(grid.IntersectionCoord{R: 3213, C: 1, D: grid.R})
			Expect(exists).To(Equal(false))
			Expect(intersection.Building).To(BeNil())
		})

		When("first player builds initial settlement", func() {
			BeforeEach(func() {
				buildSettlementCommand := domain.NewBuildSettlementCommand(
					startGameCommandOccurred,
					game.CurrentTurn(),
					grid.IntersectionCoord{R: 3, C: 1, D: grid.R},
					domain.NewSettlement(game.CurrentTurn()),
				)

				events, err := buildSettlementCommand.Process(game)
				Expect(err).To(BeNil())

				game = game.WithAppliedEvents(events...)
			})

			It("board should have the settlement with right color", func() {
				intersection, exists := game.Board().Intersection(grid.IntersectionCoord{R: 3, C: 1, D: grid.R})
				Expect(exists).To(Equal(true))

				Expect(intersection.Building).To(Equal(domain.NewSettlement(domain.Blue)))
			})

			When("first player builds a legal road", func() {
				It("board should have the road with right color ", func() {

				})
			})
			When("first player tries to build an illegal road", func() {
				It("should receive the error", func() {

				})
			})

			When("not first player tries to build a settlement", func() {
				It("should receive the error", func() {

				})
			})
		})

		When("not first player tries to build settlement", func() {
			It("should return the errur", func() {

			})
		})
	})

	Describe("play phase", func() {

	})
})
