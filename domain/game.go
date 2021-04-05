package domain

import (
	"errors"
	"reflect"
	"time"
)

type (
	// aggregate id
	GameId = string
)

var (
	// GameAlreadyStartedErr game is already started error
	GameAlreadyStartedErr = errors.New("game is already started")
	// GameAlreadyFinishedErr game is already finished
	GameAlreadyFinishedErr = errors.New("game is already finished")
	// PlayerNotExistsErr player with Color does not exist
	PlayerNotExistsErr = errors.New("player does not exist")
	// WrongTurnErr is occurred when player tries to do something during not his turn
	WrongTurnErr = errors.New("wrong turn")
)

// Game aggregate
type Game struct {
	id GameId

	players map[Color]Player
	board   Board

	stateNew          GameState
	stateStarted      GameState
	stateInitialSetup GameState
	statePlay         GameState

	currentState GameState

	turnOrder   []Color
	currentTurn Color
	totalTurns  int64
	rollHistory []Roll

	// set-up phase

	version int64
	changes []EventMessage

	boardGenerator  BoardGenerator
	playersShuffler PlayersShuffler
	diceRoller      DiceRoller

	availableResources map[ResourceCard][]ResourceCard // todo properly
	// todo trades
	// todo turn
}

func NewGame(id GameId, occurred time.Time) *Game {
	game := &Game{
		players: make(map[Color]Player),
	}

	game.Apply(NewEventDescriptor(id, GameCreated{
		GameId: id,
	}, nil, game.Version(), occurred), true)

	return game
}

func (game *Game) AddPlayer(player Player, occurred time.Time) error {
	return game.currentState.AddPlayer(player, occurred)
}

func (game *Game) SetBoardGenerator(boardGenerator BoardGenerator, occurred time.Time) error {
	return game.currentState.SetBoardGenerator(boardGenerator, occurred)
}

func (game *Game) SetPlayersShuffler(playersShuffler PlayersShuffler, occurred time.Time) error {
	return game.currentState.SetPlayersShuffler(playersShuffler, occurred)
}

func (game *Game) SetDiceRoller(diceRoller DiceRoller, occurred time.Time) error {
	return game.currentState.SetDiceRoller(diceRoller, occurred)
}

func (game *Game) GenerateBoard(occurred time.Time) error {
	return game.currentState.GenerateBoard(occurred)
}

func (game *Game) ShufflePlayers(occurred time.Time) error {
	return game.currentState.ShufflePlayers(occurred)
}

func (game *Game) StartGame(
	occurred time.Time,
) error {
	return game.currentState.StartGame(occurred)
}

func (game *Game) ChangeState(newState GameState, occurred time.Time) {
	game.Apply(
		NewEventDescriptor(
			game.Id(),
			GameEnteredState{NewState: newState},
			nil,
			game.version,
			occurred,
		),
		true,
	)

	newState.EnterState(occurred)
}

func (game *Game) BuyRoad(playerColor Color, occurred time.Time) error {
	return game.currentState.BuyRoad(playerColor, occurred)
}

func (game *Game) BuySettlement(playerColor Color, occurred time.Time) error {
	return game.currentState.BuySettlement(playerColor, occurred)
}

func (game *Game) BuyCity(playerColor Color, occurred time.Time) error {
	return game.currentState.BuyCity(playerColor, occurred)
}

func (game *Game) BuyDevelopmentCard(playerColor Color) error {
	return game.currentState.BuyDevelopmentCard(playerColor)
}

func (game *Game) PlaceSettlement(playerColor Color, settlement Settlement, occurred time.Time) error {
	return game.currentState.PlaceSettlement(playerColor, settlement, occurred)
}

func (game *Game) PlaceRoad(playerColor Color, road Road, occurred time.Time) error {
	return game.currentState.PlaceRoad(playerColor, road, occurred)
}

func (game *Game) RollDice(playerColor Color, occurred time.Time) error {
	return game.currentState.RollDice(playerColor, occurred)
}

func (game *Game) EndTurn(playerColor Color, occurred time.Time) error {
	return game.currentState.EndTurn(playerColor, occurred)
}

func (game *Game) Apply(eventMessage EventMessage, isNew bool) {
	game.incrementVersion()

	if isNew {
		game.trackChange(eventMessage)
	}

	switch event := eventMessage.Event().(type) {
	case GameEnteredState:
		game.setState(event.NewState)
	case GameCreated:
		gameStatePlayerIsToPlaceSettlement := NewGameStatePlayerIsToPlaceSettlement(game)
		gameStatePlayerIsToPlaceRoad := NewGameStatePlayerIsToPlaceRoad(game)
		gameStatePlayerIsRollingDice := NewGameStatePlayerIsRollingDice(game)

		game.id = event.GameId
		game.stateNew = NewGameStateNew(game)
		game.stateStarted = NewGameStateStarted(game)
		game.stateInitialSetup = NewGameStateInitialSetup(game, gameStatePlayerIsToPlaceSettlement, gameStatePlayerIsToPlaceRoad)
		game.statePlay = NewGameStatePlay(game, gameStatePlayerIsRollingDice, gameStatePlayerIsToPlaceSettlement, gameStatePlayerIsToPlaceRoad)

		game.setState(game.stateNew)
	case PlayerPlacedRoadEvent:
		game.trackChangeAndIncrementVersion(eventMessage)

		player, err := game.Player(event.PlayerColor)
		if err != nil {
			panic(err)
		}

		player.availableRoads--

		err = game.updatePlayer(player)
		if err != nil {
			panic(err)
		}

		err = game.placeRoad(event.Road)
		if err != nil {
			panic(err)
		}
	case PlayerPlacedInitialRoadEvent: // todo remove duplicate
		game.trackChangeAndIncrementVersion(eventMessage)

		player, err := game.Player(event.PlayerColor)
		if err != nil {
			panic(err)
		}

		player.availableRoads--

		err = game.updatePlayer(player)
		if err != nil {
			panic(err)
		}

		err = game.placeRoad(event.Road)
		if err != nil {
			panic(err)
		}
	case PlayerPlacedSettlementEvent:
		game.trackChangeAndIncrementVersion(eventMessage)

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

		err = game.placeSettlement(event.Settlement)
		if err != nil {
			panic(err)
		}
	case PlayerPlacedInitialSettlementEvent: // todo remove duplicate
		game.trackChangeAndIncrementVersion(eventMessage)

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

		err = game.placeSettlement(event.Settlement)
		if err != nil {
			panic(err)
		}
	case PlayPhaseStartedEvent:
		game.setState(game.statePlay)
	}
}

func (game *Game) trackChangeAndIncrementVersion(eventMessage EventMessage) {
	game.incrementVersion()
	game.trackChange(eventMessage)
}

// todo
func (game *Game) ProcessCommand(command interface{}) error {
	panic("todo")
}

func (game Game) Id() GameId {
	return game.id
}

func (game Game) Player(color Color) (Player, error) {
	player, exists := game.players[color]
	if !exists {
		return Player{}, PlayerNotExistsErr
	}

	return player, nil
}

func (game Game) Players() []Player {
	var players []Player

	for _, player := range game.players {
		players = append(players, player)
	}

	return players
}

func (game Game) Board() Board {
	return game.board
}

func (game Game) AvailableCommands() []string {

	return []string{}
}

func (game Game) Changes() []EventMessage {
	return game.changes
}

func (game Game) LastEvent() interface{} {
	return game.Changes()[len(game.Changes())-1].Event()
}

func (game Game) RollHistory() []Roll {
	return game.rollHistory
}

func (game Game) TotalTurns() int64 {
	return game.totalTurns
}

func (game Game) CurrentTurn() Color {
	return game.currentTurn
}

func (game Game) NextTurnColor() Color {
	turnOrder := game.State().TurnOrder()
	color := turnOrder[int(game.totalTurns+1)%len(turnOrder)]
	return color
}

func (game *Game) TurnOrder() []Color {
	return game.currentState.TurnOrder()
}

func (game Game) State() GameState {
	return game.currentState
}

func (game Game) InState(state GameState) bool {
	return reflect.TypeOf(state) == reflect.TypeOf(game.currentState)
}

func (game *Game) PlayersShuffler() PlayersShuffler {
	return game.playersShuffler
}

func (game *Game) BoardGenerator() BoardGenerator {
	return game.boardGenerator
}

func (game Game) Version() int64 {
	return game.version
}

func (game *Game) incrementVersion() {
	game.version++
}

func (game *Game) placeSettlement(settlement Settlement) error {
	intersection, exists := game.Board().Intersection(settlement.IntersectionCoord())
	if !exists {
		return BadIntersectionCoordErr
	}

	intersection.SetBuilding(settlement)

	return game.Board().UpdateIntersection(settlement.IntersectionCoord(), intersection)
}

func (game *Game) placeRoad(road Road) error {
	path, exists := game.Board().Path(road.PathCoord())
	if !exists {
		panic("todo") // todo
	}

	path.road = &road

	return game.Board().UpdatePath(road.PathCoord(), path)
}

func (game *Game) addPlayer(player Player) {
	if player.Color() == None {
		// set first available color
		for _, color := range allColors {
			_, err := game.Player(color)
			if err == PlayerNotExistsErr {
				player.SetColor(color)
				break
			}
		}
	}

	game.turnOrder = append(game.turnOrder, player.color)

	if game.players == nil {
		game.players = make(map[Color]Player)
	}

	game.players[player.Color()] = player
}

func (game *Game) removePlayer(player Player) {
	delete(game.players, player.Color())

	for i, color := range game.turnOrder {
		if player.Color() == color {
			copy(game.turnOrder[i:], game.turnOrder[i+1:])
			break
		}
	}
}

func (game *Game) updatePlayer(player Player) error {
	_, err := game.Player(player.color)
	if err != nil {
		return err
	}

	game.players[player.color] = player

	return nil
}

func (game *Game) setBoard(board Board) {
	game.board = board
}

func (game *Game) setState(state GameState) {
	game.currentState = state
}

func (game *Game) setTurnOrder(turnOrder []Color) {
	game.turnOrder = turnOrder
}

func (game *Game) trackChange(event EventMessage) {
	game.changes = append(game.changes, event)
}

func (game *Game) setPlayersShuffler(playersShuffler PlayersShuffler) {
	game.playersShuffler = playersShuffler
}

func (game *Game) setBoardGenerator(boardGenerator BoardGenerator) {
	game.boardGenerator = boardGenerator
}

func (game *Game) setDiceRoller(diceRoller DiceRoller) {
	game.diceRoller = diceRoller
}

func (game *Game) incrementTotalTurns() {
	game.totalTurns++
}

func (game *Game) setCurrentTurn(color Color) {
	game.currentTurn = color
}
