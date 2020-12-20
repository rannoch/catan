package domain

import (
	"errors"
	"github.com/rannoch/catan/grid"
	"reflect"
	"time"
)

type (
	// aggregate id
	GameId string
)

//const (
//	GameStatusNew          GameState = "new"
//	GameStatusStarted      GameState = "started"
//	GameStatusInitialSetup GameState = "initial_setup"
//	GameStatusPlay         GameState = "play"
//	GameFinished           GameState = "finished"
//)

var (
	GameAlreadyStartedErr  = errors.New("game already started")
	GameAlreadyFinishedErr = errors.New("game already finished")
	PlayerNotExistsErr     = errors.New("player does not exist")
	WrongTurnErr           = errors.New("wrong turn")
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

	state GameState // todo

	turnOrder   []Color
	currentTurn Color
	totalTurns  int64
	rollHistory []roll

	// set-up phase

	version int64
	changes []EventMessage

	boardGenerator  BoardGenerator
	playersShuffler PlayersShuffler
	// todo trades
	// todo turn
}

func (game Game) Changes() []EventMessage {
	return game.changes
}

func (game Game) Version() int64 {
	return game.version
}

func (game *Game) IncrementVersion() {
	game.version++
}

func (game Game) AvailableCommands() []string {

	return []string{}
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

func NewGame(id GameId, occurred time.Time) *Game {
	game := &Game{
		id:      id,
		players: make(map[Color]Player),
	}

	game.setState(NewGameStateNew(game), occurred)
	return game
}

func (game Game) Id() GameId {
	return game.id
}

func (game *Game) AddPlayer(player Player, occurred time.Time) error {
	return game.state.AddPlayer(player, occurred)
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

func (game *Game) updatePlayer(player Player) error {
	_, err := game.Player(player.color)
	if err != nil {
		return err
	}

	game.players[player.color] = player

	return nil
}

func (game Game) Board() Board {
	return game.board
}

func (game *Game) setBoard(board Board) {
	game.board = board
}

func (game Game) InStatus(status GameState) bool {
	return reflect.TypeOf(status) == reflect.TypeOf(game.state)
}

func (game Game) State() GameState {
	return game.state
}

func (game *Game) setState(state GameState, occurred time.Time) {
	game.state = state
	game.state.EnterState(occurred)
}

func (game Game) NextTurnColor() Color {
	return game.State().TurnOrder()[int(game.totalTurns+1)%len(game.State().TurnOrder())]
}

func (game *Game) TurnOrder() []Color {
	return game.state.TurnOrder()
}

func (game *Game) setTurnOrder(turnOrder []Color) {
	game.turnOrder = turnOrder
}

func (game *Game) TrackChange(event EventMessage) {
	game.changes = append(game.changes, event)
}

func (game *Game) PlayersShuffler() PlayersShuffler {
	return game.playersShuffler
}

func (game *Game) SetPlayersShuffler(playersShuffler PlayersShuffler) error {
	return game.state.SetPlayersShuffler(playersShuffler)
}

func (game *Game) setPlayersShuffler(playersShuffler PlayersShuffler) {
	game.playersShuffler = playersShuffler
}

func (game *Game) BoardGenerator() BoardGenerator {
	return game.boardGenerator
}

func (game *Game) SetBoardGenerator(boardGenerator BoardGenerator) error {
	return game.state.SetBoardGenerator(boardGenerator)
}

func (game *Game) setBoardGenerator(boardGenerator BoardGenerator) {
	game.boardGenerator = boardGenerator
}

func (game *Game) Apply(eventMessage EventMessage, isNew bool) {
	game.IncrementVersion()

	if isNew {
		game.TrackChange(eventMessage)
	}

	game.state.Apply(eventMessage, isNew)
}

// todo
func (game *Game) ProcessCommand(command interface{}) error {
	panic("todo")
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
	road road,
	occurred time.Time,
) error {
	return game.state.BuildRoad(
		playerColor, pathCoord, road, occurred,
	)
}
