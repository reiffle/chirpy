package main

import (
	"github.com/reiffle/chirpy/internal/database"
)

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
