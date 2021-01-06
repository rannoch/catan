package domain

import (
	"errors"
	"reflect"
	"time"

	"github.com/rannoch/catan/grid"
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

type roll struct {
	Roller Color
	Roll   int64
}

// Game aggregate
type Game struct {
	id GameId

	players map[Color]Player
	board   Board

	state GameState

	turnOrder   []Color
	currentTurn Color
	totalTurns  int64
	rollHistory []roll

	// set-up phase

	version int64
	changes []EventMessage

	boardGenerator  BoardGenerator
	playersShuffler PlayersShuffler

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
	return game.state.AddPlayer(player, occurred)
}

func (game *Game) SetBoardGenerator(boardGenerator BoardGenerator) error {
	return game.state.SetBoardGenerator(boardGenerator)
}

func (game *Game) SetPlayersShuffler(playersShuffler PlayersShuffler) error {
	return game.state.SetPlayersShuffler(playersShuffler)
}

func (game *Game) GenerateBoard(occurred time.Time) error {
	return game.state.GenerateBoard(occurred)
}

func (game *Game) ShufflePlayers(occurred time.Time) error {
	return game.state.ShufflePlayers(occurred)
}

func (game *Game) StartGame(
	occurred time.Time,
) error {
	return game.state.StartGame(occurred)
}

func (game *Game) BuildSettlement(
	playerColor Color,
	intersectionCoord grid.IntersectionCoord,
	settlement Settlement,
	occurred time.Time,
) error {
	return game.state.BuildSettlement(
		playerColor, intersectionCoord, settlement, occurred,
	)
}

func (game *Game) BuildRoad(
	playerColor Color,
	pathCoord grid.PathCoord,
	road Road,
	occurred time.Time,
) error {
	return game.state.BuildRoad(
		playerColor, pathCoord, road, occurred,
	)
}

func (game *Game) EndTurn(playerColor Color, occurred time.Time) error {
	return game.state.EndTurn(playerColor, occurred)
}

func (game *Game) Apply(eventMessage EventMessage, isNew bool) {
	game.incrementVersion()

	if isNew {
		game.trackChange(eventMessage)
	}

	switch event := eventMessage.Event().(type) {
	case GameCreated:
		game.id = event.GameId
		game.setState(NewGameStateNew(game))
	default:
		game.state.Apply(eventMessage, isNew)
	}
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

func (game Game) RollHistory() []roll {
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
	return game.state.TurnOrder()
}

func (game Game) State() GameState {
	return game.state
}

func (game Game) InState(state GameState) bool {
	return reflect.TypeOf(state) == reflect.TypeOf(game.state)
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
	game.state = state
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

func (game *Game) incrementTotalTurns() {
	game.totalTurns++
}

func (game *Game) setCurrentTurn(color Color) {
	game.currentTurn = color
}
