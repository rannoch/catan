package domain

import (
	"time"

	"github.com/rannoch/catan/grid"
)

type GameStateInitialSetup struct {
	game *Game

	statePlayerIsToPlaceSettlement GameState
	statePlayerIsToPlaceRoad       GameState

	currentSubState                       GameState
	playerIsToPlaceRoadAdjacentToBuilding grid.IntersectionCoord

	settlements []Settlement

	GameStateDefault
}

func NewGameStateInitialSetup(
	game *Game,
	statePlayerIsToPlaceSettlement GameState,
	statePlayerIsToPlaceRoad GameState,
) *GameStateInitialSetup {
	return &GameStateInitialSetup{
		game:                           game,
		statePlayerIsToPlaceSettlement: statePlayerIsToPlaceSettlement,
		statePlayerIsToPlaceRoad:       statePlayerIsToPlaceRoad,
	}
}

var _ GameState = (*GameStateInitialSetup)(nil)

func (gameStatusInitialSetup *GameStateInitialSetup) StartGame(time.Time) error {
	return GameAlreadyStartedErr
}

func (gameStatusInitialSetup *GameStateInitialSetup) EnterState(occurred time.Time) {
	game := gameStatusInitialSetup.game

	playerStartedHisTurnEventMessage := NewEventDescriptor(
		game.Id(),
		PlayerStartedHisTurnEvent{
			PlayerColor: game.turnOrder[0],
		},
		nil,
		game.Version(),
		occurred,
	)

	game.Apply(playerStartedHisTurnEventMessage, true)
}

func (gameStatusInitialSetup *GameStateInitialSetup) PlaceSettlement(playerColor Color, settlement Settlement, occurred time.Time) error {
	game := gameStatusInitialSetup.game

	if err := gameStatusInitialSetup.currentSubState.PlaceSettlement(playerColor, settlement, occurred); err != nil {
		return err
	}

	player, err := gameStatusInitialSetup.game.Player(playerColor)
	if err != nil {
		return err
	}

	if player.HasPlacedInitialBuildings() {
		playerPickedResourcesEventMessage := NewEventDescriptor(
			game.Id(),
			PlayerPickedResourcesEvent{
				PlayerColor:     playerColor,
				PickedResources: gameStatusInitialSetup.getInitialResources(settlement.IntersectionCoord()),
			},
			nil,
			game.version,
			occurred,
		)

		game.Apply(playerPickedResourcesEventMessage, true)
	}

	return nil
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

func (gameStatusInitialSetup *GameStateInitialSetup) Apply(eventMessage EventMessage, isNew bool) {
	game := gameStatusInitialSetup.game

	switch event := eventMessage.Event().(type) {
	case PlayerStartedHisTurnEvent:
		game.setCurrentTurn(event.PlayerColor)
		gameStatusInitialSetup.currentSubState = gameStatusInitialSetup.statePlayerIsToPlaceSettlement
	case PlayerFinishedHisTurnEvent:
		game.incrementTotalTurns()
		game.setCurrentTurn(None)
	case PlayerPlacedSettlementEvent:
		gameStatusInitialSetup.currentSubState.Apply(eventMessage, isNew)

		gameStatusInitialSetup.settlements = append(gameStatusInitialSetup.settlements, event.Settlement)
		gameStatusInitialSetup.currentSubState = gameStatusInitialSetup.statePlayerIsToPlaceRoad
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
	case PlayPhaseStartedEvent:
		game.setState(game.statePlay)
	default:
		gameStatusInitialSetup.currentSubState.Apply(eventMessage, isNew)
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
	game.currentState.EnterState(occurred)

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
