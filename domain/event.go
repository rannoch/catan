package domain

import (
	"reflect"
	"time"
)

type EventMessage interface {
	AggregateId() string
	Event() interface{}
	EvenType() string
	Occurred() time.Time
}

type EventDescriptor struct {
	id       string
	event    interface{}
	headers  map[string]interface{}
	version  int64
	occurred time.Time
}

func NewEventDescriptor(
	id string,
	event interface{},
	headers map[string]interface{},
	version int64,
	occurred time.Time,
) *EventDescriptor {
	return &EventDescriptor{
		id:       id,
		event:    event,
		headers:  headers,
		version:  version,
		occurred: occurred,
	}
}

func (e EventDescriptor) AggregateId() string {
	return e.id
}

func (e EventDescriptor) Event() interface{} {
	return e.event
}

func (e EventDescriptor) EvenType() string {
	return reflect.TypeOf(e.event).Elem().Name()
}

func (e EventDescriptor) Occurred() time.Time {
	return e.occurred
}

type GameCreated struct {
	GameId GameId
}

// todo before game started events

type PlayerJoinedTheGameEvent struct {
	Player Player
}

type PlayerLeftTheGameEvent struct {
	Player Player
}

type BoardGeneratorSelectedEvent struct {
	boardGenerator BoardGenerator
}

type PlayersShufflerSelectedEvent struct {
	playersShuffler PlayersShuffler
}

/// In-game events
type GameStartedEvent struct{}

type BoardGeneratedEvent struct {
	NewBoard Board
}

type PlayersShuffledEvent struct {
	PlayersInOrder []Color
}

type InitialSetupPhaseStartedEvent struct {
}

type PlayPhaseStartedEvent struct {
}

type PlayerRolledDiceEvent struct {
	Roll Roll
}

type PlayerPickedResourcesEvent struct {
	PlayerColor     Color
	PickedResources []ResourceCard
}

type PlayerWasRobbedByRobberEvent struct {
	robbedPlayerColor Color
	dumpedResources   []ResourceCard
}

type PlayerWasRobbedByPlayerEvent struct {
	robbingPlayerColor Color
	robbedPlayerColor  Color
	dumpedResources    []ResourceCard
}

type PlayerStartedHisTurnEvent struct {
	PlayerColor Color
}

type PlayerFinishedHisTurnEvent struct {
	PlayerColor Color
}

type PlayerPlacedSettlementEvent struct {
	PlayerColor Color
	Settlement  Settlement
}

type PlayerPlacedRoadEvent struct {
	PlayerColor Color
	Road        Road
}
