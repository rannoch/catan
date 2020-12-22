package domain

import (
	"time"

	"github.com/rannoch/catan/grid"
)

type GameStatePlay struct {
	GameStateDefault
	game *Game
}

func NewGameStatePlay(game *Game) *GameStatePlay {
	return &GameStatePlay{game: game}
}

var _ GameState = (*GameStatePlay)(nil)

func (gameStatePlay *GameStatePlay) StartGame(time.Time) error {
	return GameAlreadyStartedErr
}

func (gameStatePlay *GameStatePlay) EnterState(occurred time.Time) {
	// todo spread initial resources
}

func (gameStatePlay *GameStatePlay) BuildSettlement(playerColor Color, intersectionCoord grid.IntersectionCoord, settlement Settlement, occurred time.Time) error {
	game := gameStatePlay.game

	player, err := game.Player(playerColor)
	if err != nil {
		return PlayerNotExistsErr
	}

	if game.CurrentTurn() != playerColor {
		return WrongTurnErr
	}

	if err := player.CanBuy(settlement); err != nil {
		return err
	}

	if err := player.CanBuildSettlement(); err != nil {
		return err
	}

	if err := game.Board().CanBuildSettlementOrCity(intersectionCoord, settlement); err != nil {
		return err
	}

	playerBuiltSettlementEventMessage := NewEventDescriptor(
		string(game.Id()),
		PlayerBuiltSettlementEvent{
			playerColor:       playerColor,
			intersectionCoord: intersectionCoord,
			settlement:        settlement,
		},
		nil,
		game.Version(),
		occurred,
	)

	game.Apply(playerBuiltSettlementEventMessage, true)

	return nil
}

func (gameStatePlay *GameStatePlay) BuildRoad(playerColor Color, pathCoord grid.PathCoord, road Road, occurred time.Time) error {
	game := gameStatePlay.game

	if game.CurrentTurn() != playerColor {
		return WrongTurnErr
	}

	player, err := game.Player(playerColor)
	if err != nil {
		return err
	}

	// check if the player has available road
	if err := player.HasAvailableRoad(); err != nil {
		return err
	}

	// check if the player has enough resources to buy a road
	if err := player.CanBuy(road); err != nil {
		return err
	}

	// check if the board allowing to build
	if err := game.Board().CanBuildRoad(pathCoord, road); err != nil {
		return err
	}

	gameStatePlay.Apply(
		NewEventDescriptor(string(game.Id()), PlayerBuiltRoadEvent{
			playerColor: playerColor,
			pathCoord:   pathCoord,
			road:        road,
		}, nil, game.version, occurred,
		),
		true,
	)

	return nil
}

func (gameStatePlay *GameStatePlay) BuyDevelopmentCard(playerColor Color, card DevelopmentCard) error {
	panic("implement me")
}

func (gameStatePlay *GameStatePlay) EndTurn(playerColor Color, occurred time.Time) error {
	panic("implement me")
	return nil
}

func (gameStatePlay *GameStatePlay) CurrentTurn() Color {
	panic("implement me")
}

func (gameStatePlay *GameStatePlay) TurnOrder() []Color {
	return gameStatePlay.game.turnOrder
}

func (gameStatePlay *GameStatePlay) Apply(eventMessage EventMessage, isNew bool) {
	game := gameStatePlay.game

	switch event := eventMessage.Event().(type) {
	case PlayerStartedHisTurnEvent:
		game.currentTurn = event.playerColor
	case PlayerFinishedHisTurnEvent:
		game.totalTurns++       // todo increment
		game.currentTurn = None // todo setter

		// todo invoke next player start his turn
	case PlayerBuiltSettlementEvent:
		player, err := game.Player(event.playerColor)
		if err != nil {
			panic(err)
		}

		player.victoryPoints += event.settlement.VictoryPoints()
		player.availableSettlements--

		err = game.updatePlayer(player)
		if err != nil {
			panic(err)
		}

		game.setBoard(game.Board().BuildSettlementOrCity(event.intersectionCoord, event.settlement))
	case PlayerBuiltRoadEvent:
		player, err := game.Player(event.playerColor)
		if err != nil {
			panic(err)
		}

		player.availableRoads--

		err = game.updatePlayer(player)
		if err != nil {
			panic(err)
		}

		game.Board().BuildRoad(event.pathCoord, event.road)
	case PlayerWasRobbedByRobberEvent:
		// todo
	case PlayerWasRobbedByPlayerEvent:
	case PlayerPickedResourcesEvent:
		player, err := game.Player(event.playerColor)
		if err != nil {
			panic(err) // todo
		}

		player.GainResources(event.pickedResources)
	case PlayerRolledDiceEvent:
		game.rollHistory = append(game.rollHistory, event.roll)
	}
}
