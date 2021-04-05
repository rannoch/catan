package domain_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rannoch/catan/domain"
	. "github.com/rannoch/catan/domain/games/catan_championship_premium_13_BUGGED_Semi_Final"
)

var _ = Describe("Catan state play", func() {
	var (
		game *domain.Game
		err  error
	)

	BeforeEach(func() {
		game = &domain.Game{}

		for _, event := range Events {
			game.Apply(domain.NewEventDescriptor(game.Id(), event, nil, game.Version(), time.Now()), true)
		}
	})

	It("aggregate id should match", func() {
		Expect(game.Id()).To(Equal("Catan Championship Premium #13 - BUGGED Semi-Final"))
	})

	Specify("turns checks", func() {
		Expect(game.CurrentTurn()).To(Equal(domain.Black))
		Expect(game.NextTurnColor()).To(Equal(domain.Orange))
		Expect(game.TotalTurns()).To(Equal(int64(8)))
	})

	It("game should enter play phase", func() {
		Expect(err).To(BeNil())

		Expect(game.InState(&domain.GameStatePlay{})).To(BeTrue())

		Expect(game.Version()).To(Equal(int64(51)))

		_, err := game.Player(domain.Orange)
		Expect(err).To(BeNil())
		Expect(len(game.Players())).To(Equal(4))
	})

	Specify("turn order is right", func() {
		Expect(game.TurnOrder()).To(Equal([]domain.Color{
			domain.Black,
			domain.Orange,
			domain.Blue,
			domain.Green,
		}))
	})

	Specify("players available buildings are correct", func() {
		for _, player := range game.Players() {
			Expect(player.AvailableSettlements()).To(Equal(int64(3)))
			Expect(player.AvailableCities()).To(Equal(int64(4)))
			Expect(player.AvailableRoads()).To(Equal(int64(13)))
		}
	})

	Specify("game state - black is rolling dice", func() {

	})

	When("black rolls a dice", func() {


		It("", func() {

		})
	})
})
