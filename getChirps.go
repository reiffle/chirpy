package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	payload, err := cfg.DB.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	var cleanedChirps []returnVals
	for _, item := range payload {
		cleanedChirps = append(cleanedChirps, cleanChirp(item))
	}
	respondWithJSON(w, http.StatusOK, cleanedChirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("chirpID")
	parsedUUID, err := uuid.Parse(idString) // Use uuid.Parse() for parsing
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID format", err)
		return
	}
	payload, err := cfg.DB.GetChirp(r.Context(), parsedUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp", err)
		return
	}
	cleanedChirp := cleanChirp(payload)
	respondWithJSON(w, http.StatusOK, cleanedChirp)
}
