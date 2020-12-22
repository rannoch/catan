package domain_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rannoch/catan/domain"
)

var _ = Describe("Catan play phase", func() {
	var (
		game *domain.Game
		//err     error
		players = []domain.Player{
			domain.NewPlayer(domain.Blue, "baska"),
			domain.NewPlayer(domain.Red, "masha"),
			domain.NewPlayer(domain.Yellow, "vasya"),
			domain.NewPlayer(domain.Green, "bot"),
		}
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
		// initial placing

	})
})
