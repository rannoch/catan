package domain

import "time"

type UserId string

type User struct {
	id         UserId
	name       string
	totalGames int64
	winRate    float64
	createdAt  time.Time
	// etc
}
