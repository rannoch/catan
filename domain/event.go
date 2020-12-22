package domain

import (
	"reflect"
	"time"

	"github.com/rannoch/catan/grid"
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
type GameStartedEvent struct {
	Occurred time.Time
}

type BoardGeneratedEvent struct {
	newBoard Board
}

type PlayersShuffledEvent struct {
	playersInOrder []Color
}

type InitialSetupPhaseStartedEvent struct {
	occurred time.Time
}

type PlayPhaseStartedEvent struct {
	occurred time.Time
}

type PlayerRolledDiceEvent struct {
	roll roll
}

type PlayerPickedResourcesEvent struct {
	playerColor     Color
	pickedResources []resource
}

type PlayerWasRobbedByRobberEvent struct {
	robbedPlayerColor Color
	dumpedResources   []resource
}

type PlayerWasRobbedByPlayerEvent struct {
	robbingPlayerColor Color
	robbedPlayerColor  Color
	dumpedResources    []resource
}

type PlayerFinishedHisTurnEvent struct {
	playerColor Color
	occurred    time.Time
}

type PlayerStartedHisTurnEvent struct {
	playerColor Color
	occurred    time.Time
}

type PlayerBuiltSettlementEvent struct {
	playerColor       Color
	intersectionCoord grid.IntersectionCoord
	settlement        Settlement
}

type PlayerBuiltRoadEvent struct {
	playerColor Color
	pathCoord   grid.PathCoord
	road        Road
	occurred    time.Time
}
