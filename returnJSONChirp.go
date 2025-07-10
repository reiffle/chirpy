package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/reiffle/chirpy/internal/database"
)

type returnVals struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

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
