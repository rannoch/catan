package domain_test

import (
	"time"

	"github.com/rannoch/catan/domain"
	"github.com/rannoch/catan/domain/fixtures"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Catan state play", func() {
	var (
		game *domain.Game
		err  error
	)

	BeforeEach(func() {
		game = &domain.Game{}

		for _, event := range fixtures.JustAfterInitialStateEvents {
			game.Apply(domain.NewEventDescriptor(game.Id(), event, nil, game.Version(), time.Now()), true)
		}
	})

	It("aggregate id should match", func() {
		Expect(game.Id()).To(Equal("test_id"))
	})

	Specify("turns checks", func() {
		Expect(game.CurrentTurn()).To(Equal(domain.Blue))
		Expect(game.NextTurnColor()).To(Equal(domain.White))
		Expect(game.TotalTurns()).To(Equal(int64(8)))
	})

	It("game should enter play phase", func() {
		Expect(err).To(BeNil())

		Expect(game.InState(&domain.GameStatePlay{})).To(BeTrue())

		Expect(game.Version()).To(Equal(int64(47)))

		_, err := game.Player(domain.Red)
		Expect(err).To(BeNil())
		Expect(len(game.Players())).To(Equal(4))
	})

	Specify("turn order is right", func() {
		Expect(game.TurnOrder()).To(Equal([]domain.Color{
			domain.Blue,
			domain.White,
			domain.Red,
			domain.Yellow,
		}))
	})

	Specify("players available buildings are correct", func() {
		for _, player := range game.Players() {
			Expect(player.AvailableSettlements()).To(Equal(int64(3)))
			Expect(player.AvailableCities()).To(Equal(int64(4)))
			Expect(player.AvailableRoads()).To(Equal(int64(13)))
		}
	})

})
