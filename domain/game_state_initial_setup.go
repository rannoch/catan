package domain

import (
	"github.com/rannoch/catan/grid"
	"time"
)

type GameStateInitialSetup struct {
	GameStateDefault
	game *Game
}

func NewGameStateInitialSetup(game *Game) *GameStateInitialSetup {
	return &GameStateInitialSetup{game: game}
}

var _ GameState = (*GameStateInitialSetup)(nil)

func (gameStatusInitialSetup *GameStateInitialSetup) StartGame(time.Time, ) error {
	return GameAlreadyStartedErr
}

func (gameStatusInitialSetup *GameStateInitialSetup) EnterState(occurred time.Time) {
	playerStartedHisTurnEventMessage := NewEventDescriptor(
		string(gameStatusInitialSetup.game.Id()),
		PlayerStartedHisTurnEvent{
			playerColor: gameStatusInitialSetup.game.turnOrder[0],
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

	// todo check if the player already placed a building

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

func (gameStatusInitialSetup *GameStateInitialSetup) BuildRoad(playerColor Color, pathCoord grid.PathCoord, road road, occurred time.Time) error {
	game := gameStatusInitialSetup.game

	if game.CurrentTurn() != playerColor {
		return WrongTurnErr
	}

	// todo check if the player already placed a building

	// check if the board allowing to build
	if err := game.Board().CanBuildRoad(pathCoord, road); err != nil {
		return err
	}

	game.Apply(
		NewEventDescriptor(string(game.Id()), PlayerBuiltRoadEvent{
			playerColor: playerColor,
			pathCoord:   pathCoord,
			road:        road,
		}, nil, game.version, occurred,
		),
		true,
	)

	// todo make aftereffect
	game.Apply(
		NewEventDescriptor(string(game.Id()), PlayerFinishedHisTurnEvent{
			playerColor: playerColor,
		}, nil, game.version, occurred,
		),
		true,
	)

	return nil
}

func (gameStatusInitialSetup *GameStateInitialSetup) EndTurn(playerColor Color, occurred time.Time) error {
	// todo check if the player placed a building and a road

	gameStatusInitialSetup.game.Apply(
		NewEventDescriptor(
			string(gameStatusInitialSetup.game.Id()),
			PlayerFinishedHisTurnEvent{
				playerColor: playerColor,
			},
			nil,
			gameStatusInitialSetup.game.version,
			occurred,
		),
		true,
	)
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

		game.setBoard(game.Board().BuildRoad(event.pathCoord, event.road))

		// todo start play phase after last road has been built
	case PlayPhaseStartedEvent:
		game.setState(NewGameStatePlay(game), event.occurred)

	}
}
