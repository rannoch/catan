package domain

import "errors"

type (
	// aggregate id
	GameId     string
	GameStatus string
)

const (
	GameStatusNew          GameStatus = "new"
	GameStatusStarted      GameStatus = "started"
	GameStatusInitialSetup GameStatus = "initial_setup"
	GameStatusPlay         GameStatus = "play"
	GameFinished           GameStatus = "finished"
)

var (
	GameAlreadyStartedErr  = errors.New("game already started")
	GameAlreadyFinishedErr = errors.New("game already finished")
	PlayerNotExistsErr     = errors.New("player does not exist")
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

	status GameStatus // todo

	turnOrder         []Color
	currentTurn       Color
	totalTurns        int64
	rollHistory       []roll
	availableCommands []Command // todo check

	// set-up phase
	setUpPhaseTurnOrder []Color

	version int64
	// todo trades
	// todo turn
}

func (game Game) Version() int64 {
	return game.version
}

func (game Game) SetUpPhaseTurnOrder() []Color {
	return game.setUpPhaseTurnOrder
}

func (game Game) AvailableCommands() []Command {
	return game.availableCommands
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

func NewGame(id GameId) Game {
	return Game{
		id:      id,
		status:  GameStatusNew,
		players: make(map[Color]Player),
	}
}

func (game Game) Id() GameId {
	return game.id
}

func (game Game) AddPlayers(players ...Player) Game {
	for _, player := range players {
		game.players[player.color] = player
	}

	return game
}

func (game Game) Player(color Color) (Player, error) {
	player, exists := game.players[color]
	if !exists {
		return Player{}, PlayerNotExistsErr
	}

	return player, nil
}

func (game Game) WithUpdatedPlayer(player Player) (Game, error) {
	_, err := game.Player(player.color)
	if err != nil {
		return Game{}, err
	}

	game.players[player.color] = player

	return game, nil
}

func (game Game) Board() Board {
	return game.board
}

func (game Game) WithBoard(board Board) Game {
	game.board = board
	return game
}

func (game Game) InStatus(status GameStatus) bool {
	return game.status == status
}

func (game Game) Status() GameStatus {
	return game.status
}

// todo bad method, need to split to different
func (game Game) WithStatus(status GameStatus) Game {
	game.status = status
	return game
}

func (game Game) IsStarted() bool {
	return game.status != GameStatusNew
}

func (game Game) NextTurnColor() Color {
	return game.TurnOrder()[int(game.totalTurns+1)%len(game.TurnOrder())]
}

func (game Game) TurnOrder() []Color {
	if game.InStatus(GameStatusPlay) {
		return game.turnOrder
	}

	if game.InStatus(GameStatusInitialSetup) {
		var turnOrderReversed = make([]Color, len(game.turnOrder))
		copy(turnOrderReversed, game.turnOrder)

		for i := len(turnOrderReversed)/2 - 1; i >= 0; i-- {
			opp := len(turnOrderReversed) - 1 - i
			turnOrderReversed[i], turnOrderReversed[opp] = turnOrderReversed[opp], turnOrderReversed[i]
		}

		return append(game.turnOrder, turnOrderReversed...)
	}

	return nil
}

func (game Game) WithTurnOrder(colors []Color) Game {
	game.turnOrder = colors
	return game
}

func (game Game) WithAppliedEvents(events ...Event) Game {
	for _, event := range events {
		game = event.Apply(game)
	}

	return game
}
