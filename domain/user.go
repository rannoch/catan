package domain

import "time"

// UserId User aggregate id
type UserId = string

// User entity
type User struct {
	id         UserId
	name       string
	totalGames int64
	winRate    float64
	createdAt  time.Time
	// etc
}
