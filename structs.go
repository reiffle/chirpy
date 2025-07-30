package main

import (
	"sync/atomic"
	"time"

	"github.com/reiffle/chirpy/internal/database"

	"github.com/google/uuid"
)

// Api config
type apiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
	platform       string
	secret         string
	polka_key      string
}

// User structs
type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type returnVals struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

// JSON chirp
func cleanChirp(item database.Chirp) returnVals {
	cleanedChirp := returnVals{
		ID:        item.ID,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
		Body:      item.Body,
		UserID:    item.UserID,
	}
	return cleanedChirp
}
