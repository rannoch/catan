package domain

import (
	"time"

	"github.com/rannoch/catan/grid"
)

type GameStateInitialSetup struct {
	turnOrder []Color

	currentTurn int

	GameStateDefault

	game *Game
}

func NewGameStateInitialSetup(
	gameStateDefault GameStateDefault,
	game *Game,
) *GameStateInitialSetup {
	var turnOrderReversed = make([]Color, len(game.turnOrder))
	copy(turnOrderReversed, game.turnOrder)

	for i := len(turnOrderReversed)/2 - 1; i >= 0; i-- {
		opp := len(turnOrderReversed) - 1 - i
		turnOrderReversed[i], turnOrderReversed[opp] = turnOrderReversed[opp], turnOrderReversed[i]
	}

	initialSetupTurnOrder := append(game.turnOrder, turnOrderReversed...)

	return &GameStateInitialSetup{
		turnOrder:        initialSetupTurnOrder,
		GameStateDefault: gameStateDefault,
		game:             game,
	}
}

var _ GameState = (*GameStateInitialSetup)(nil)

func (gameStatusInitialSetup *GameStateInitialSetup) StartGame(time.Time) error {
	return GameAlreadyStartedErr
}

func (gameStatusInitialSetup *GameStateInitialSetup) EnterState(occurred time.Time) {
	// if all players finished their turns, head to play state
	if gameStatusInitialSetup.checkAndMoveToPlayStateIfNeeded(occurred) {
		return
	}

	// change state to next turn
	game := gameStatusInitialSetup.game

	gameStateInitialSetupPlayerTurn := NewGameStateInitialSetupPlayerTurn(game, gameStatusInitialSetup.CurrentTurn(), gameStatusInitialSetup)
	game.ChangeState(gameStateInitialSetupPlayerTurn, occurred)
}

func (gameStatusInitialSetup *GameStateInitialSetup) PlaceRoad(playerColor Color, road Road, occurred time.Time) error {
	if !gameStatusInitialSetup.isRoadAdjacentToLastSettlement(road.PathCoord()) {
		return CommandIsForbiddenErr
	}

	if err := gameStatusInitialSetup.currentSubState.PlaceRoad(playerColor, road, occurred); err != nil {
		return err
	}

	if err := gameStatusInitialSetup.endTurn(playerColor, occurred); err != nil {
		return err
	}

	return nil
}

func (gameStatusInitialSetup *GameStateInitialSetup) isRoadAdjacentToLastSettlement(pathCoord grid.PathCoord) bool {
	game := gameStatusInitialSetup.game

	if len(gameStatusInitialSetup.settlements) == 0 {
		return false
	}

	lastSettlement := gameStatusInitialSetup.settlements[len(gameStatusInitialSetup.settlements)-1]

	intersectionAdjacentPaths := game.board.IntersectionAdjacentPaths(lastSettlement.IntersectionCoord())

	for _, adjacentPathCoord := range intersectionAdjacentPaths {
		if adjacentPathCoord == pathCoord {
			return true
		}
	}

	return false
}

func (gameStatusInitialSetup *GameStateInitialSetup) endTurn(playerColor Color, occurred time.Time) error {
	game := gameStatusInitialSetup.game

	if game.CurrentTurn() != playerColor {
		return WrongTurnErr
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
	return gameStatusInitialSetup.turnOrder[gameStatusInitialSetup.currentTurn]
}

func (gameStatusInitialSetup *GameStateInitialSetup) TurnOrder() []Color {
	return gameStatusInitialSetup.turnOrder
}

func (gameStatusInitialSetup *GameStateInitialSetup) Apply(eventMessage EventMessage, isNew bool) {
	game := gameStatusInitialSetup.game

	switch event := eventMessage.Event().(type) {
	case PlayerStartedHisTurnEvent:
		game.setCurrentTurn(event.PlayerColor)
		gameStatusInitialSetup.currentSubState = gameStatusInitialSetup.statePlayerIsPlacingSettlement
	case PlayerFinishedHisTurnEvent:
		game.incrementTotalTurns()
		game.setCurrentTurn(None)
	case PlayerPlacedSettlementEvent:
		gameStatusInitialSetup.currentSubState.Apply(eventMessage, isNew)

		gameStatusInitialSetup.settlements = append(gameStatusInitialSetup.settlements, event.Settlement)
		gameStatusInitialSetup.currentSubState = gameStatusInitialSetup.statePlayerIsPlacingRoad
		gameStatusInitialSetup.playerIsToPlaceRoadAdjacentToBuilding = event.Settlement.IntersectionCoord()
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
	}
}

func (gameStatusInitialSetup *GameStateInitialSetup) checkAndMoveToPlayStateIfNeeded(occurred time.Time) bool {
	game := gameStatusInitialSetup.game

	for _, player := range game.Players() {
		if !player.HasPlacedInitialBuildingsAndRoads() {
			return false
		}
	}

	game.ChangeState(NewGameStatePlay(game), occurred)

	return true
}

func (gameStatusInitialSetup *GameStateInitialSetup) getInitialResources(intersectionCoord grid.IntersectionCoord) []ResourceCard {
	game := gameStatusInitialSetup.game

	hexCoords := game.Board().IntersectionAdjacentHexes(intersectionCoord)

	var resources []ResourceCard

	for _, hexCoord := range hexCoords {
		hex, exists := game.Board().Hex(hexCoord)
		if !exists {
			continue
		}

		resources = append(resources, hex.Resource.GetResourceCard(1)...)
	}

	return resources
}
