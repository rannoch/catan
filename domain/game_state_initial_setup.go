package domain

import (
	"time"

	"github.com/rannoch/catan/grid"
)

type GameStateInitialSetup struct {
	GameStateDefault
	game *Game
}

func NewGameStateInitialSetup(game *Game) *GameStateInitialSetup {
	return &GameStateInitialSetup{game: game}
}

var _ GameState = (*GameStateInitialSetup)(nil)

func (gameStatusInitialSetup *GameStateInitialSetup) StartGame(time.Time) error {
	return GameAlreadyStartedErr
}

func (gameStatusInitialSetup *GameStateInitialSetup) EnterState(occurred time.Time) {
	playerStartedHisTurnEventMessage := NewEventDescriptor(
		gameStatusInitialSetup.game.Id(),
		PlayerStartedHisTurnEvent{
			PlayerColor: gameStatusInitialSetup.game.turnOrder[0],
		},
		nil,
		gameStatusInitialSetup.game.Version(),
		occurred,
	)

	gameStatusInitialSetup.game.Apply(playerStartedHisTurnEventMessage, true)
}

func (gameStatusInitialSetup *GameStateInitialSetup) BuildSettlement(playerColor Color, intersectionCoord grid.IntersectionCoord, settlement Settlement, occurred time.Time) error {
	game := gameStatusInitialSetup.game

	if game.CurrentTurn() != playerColor {
		return WrongTurnErr
	}

	player, err := game.Player(playerColor)
	if err != nil {
		return err
	}

	// check if the player already placed a building
	if player.availableSettlements != player.availableRoads-10 {
		return CommandIsForbiddenErr
	}

	if err := game.Board().CanBuildSettlementOrCity(intersectionCoord, settlement); err != nil {
		return err
	}

	playerBuiltSettlementEventMessage := NewEventDescriptor(
		game.Id(),
		PlayerBuiltSettlementEvent{
			PlayerColor:       playerColor,
			IntersectionCoord: intersectionCoord,
			Settlement:        settlement,
		},
		nil,
		game.Version(),
		occurred,
	)

	game.Apply(playerBuiltSettlementEventMessage, true)

	player, err = game.Player(playerColor)
	if err != nil {
		return err
	}

	if player.HasPlacedInitialBuildings() {
		playerPickedResourcesEventMessage := NewEventDescriptor(
			game.Id(),
			PlayerPickedResourcesEvent{
				PlayerColor:     playerColor,
				PickedResources: game.board.IntersectionInitialResources(intersectionCoord),
			},
			nil,
			game.version,
			occurred,
		)

		game.Apply(playerPickedResourcesEventMessage, true)
	}

	return nil
}

func (gameStatusInitialSetup *GameStateInitialSetup) BuildRoad(playerColor Color, pathCoord grid.PathCoord, road Road, occurred time.Time) error {
	game := gameStatusInitialSetup.game

	if game.CurrentTurn() != playerColor {
		return WrongTurnErr
	}

	player, err := game.Player(playerColor)
	if err != nil {
		return err
	}

	// check if the player already placed a building
	if player.availableRoads-player.availableSettlements != 11 {
		return CommandIsForbiddenErr
	}

	// check if the board allowing to build
	if err := game.Board().CanBuildRoad(pathCoord, road); err != nil {
		return err
	}

	game.Apply(
		NewEventDescriptor(
			game.Id(),
			PlayerBuiltRoadEvent{
				PlayerColor: playerColor,
				PathCoord:   pathCoord,
				Road:        road,
			},
			nil,
			game.version,
			occurred,
		),
		true,
	)

	if err := gameStatusInitialSetup.EndTurn(playerColor, occurred); err != nil {
		return err
	}

	return nil
}

func (gameStatusInitialSetup *GameStateInitialSetup) EndTurn(playerColor Color, occurred time.Time) error {
	game := gameStatusInitialSetup.game

	if game.CurrentTurn() != playerColor {
		return WrongTurnErr
	}

	player, err := game.Player(playerColor)
	if err != nil {
		return err
	}

	// check if the player placed a building and a road
	if !(player.availableSettlements == 4 && player.availableRoads == 14) &&
		!(player.availableSettlements == 3 && player.availableRoads == 13) {
		return CommandIsForbiddenErr
	}

	nextPlayerColor := game.NextTurnColor() // todo better approach

	gameStatusInitialSetup.game.Apply(
		NewEventDescriptor(
			gameStatusInitialSetup.game.Id(),
			PlayerFinishedHisTurnEvent{
				PlayerColor: playerColor,
			},
			nil,
			gameStatusInitialSetup.game.version,
			occurred,
		),
		true,
	)

	if gameStatusInitialSetup.checkAndMoveToPlayStateIfNeeded(occurred) {
		return nil
	}

	nextPlayerStartedHisTurnEventMessage := NewEventDescriptor(
		game.Id(),
		PlayerStartedHisTurnEvent{
			PlayerColor: nextPlayerColor,
		},
		nil,
		game.version,
		occurred,
	)

	game.Apply(nextPlayerStartedHisTurnEventMessage, true)

	return nil
}

func (gameStatusInitialSetup *GameStateInitialSetup) CurrentTurn() Color {
	return gameStatusInitialSetup.game.CurrentTurn()
}

func (gameStatusInitialSetup *GameStateInitialSetup) TurnOrder() []Color {
	var turnOrderReversed = make([]Color, len(gameStatusInitialSetup.game.turnOrder))
	copy(turnOrderReversed, gameStatusInitialSetup.game.turnOrder)

	for i := len(turnOrderReversed)/2 - 1; i >= 0; i-- {
		opp := len(turnOrderReversed) - 1 - i
		turnOrderReversed[i], turnOrderReversed[opp] = turnOrderReversed[opp], turnOrderReversed[i]
	}

	return append(gameStatusInitialSetup.game.turnOrder, turnOrderReversed...)
}

func (gameStatusInitialSetup *GameStateInitialSetup) Apply(eventMessage EventMessage, _ bool) {
	game := gameStatusInitialSetup.game

	switch event := eventMessage.Event().(type) {
	case PlayerStartedHisTurnEvent:
		game.setCurrentTurn(event.PlayerColor)
	case PlayerFinishedHisTurnEvent:
		game.incrementTotalTurns()
		game.setCurrentTurn(None)
	case PlayerBuiltSettlementEvent:
		player, err := game.Player(event.PlayerColor)
		if err != nil {
			panic(err)
		}

		player.victoryPoints += event.Settlement.VictoryPoints()
		player.availableSettlements--

		err = game.updatePlayer(player)
		if err != nil {
			panic(err)
		}

		game.setBoard(game.Board().BuildSettlementOrCity(event.IntersectionCoord, event.Settlement))
	case PlayerBuiltRoadEvent:
		player, err := game.Player(event.PlayerColor)
		if err != nil {
			panic(err)
		}

		player.availableRoads--

		err = game.updatePlayer(player)
		if err != nil {
			panic(err)
		}

		game.Board().BuildRoad(event.PathCoord, event.Road)
	case PlayerPickedResourcesEvent:
		player, err := game.Player(event.PlayerColor)
		if err != nil {
			panic(err) // todo
		}

		player.GainResources(event.PickedResources)

		err = game.updatePlayer(player)
		if err != nil {
			panic(err)
		}
	case PlayPhaseStartedEvent:
		game.setState(NewGameStatePlay(game))
	}
}

func (gameStatusInitialSetup *GameStateInitialSetup) checkAndMoveToPlayStateIfNeeded(occurred time.Time) bool {
	game := gameStatusInitialSetup.game

	for _, player := range game.Players() {
		if !player.HasPlacedInitialBuildingsAndRoads() {
			return false
		}
	}

	game.Apply(
		NewEventDescriptor(
			game.Id(),
			PlayPhaseStartedEvent{},
			nil,
			game.version,
			occurred,
		),
		true,
	)
	game.state.EnterState(occurred)

	return true
}
