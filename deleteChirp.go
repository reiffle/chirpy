package main

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/reiffle/chirpy/internal/auth"
	"github.com/reiffle/chirpy/internal/database"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	//Authenticate credentials
	secret := cfg.secret

	decoded_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	userId, err := auth.ValidateJWT(decoded_token, secret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	//Authenticate delete request
	chirpIdString := r.PathValue("chirpID")
	parsedChirp, err := uuid.Parse(chirpIdString) // Use uuid.Parse() for parsing
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID format", err)
		return
	}
	chirp, err := cfg.DB.GetChirp(r.Context(), parsedChirp)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp", err)
		return
	}

	if chirp.UserID != userId {
		respondWithError(w, http.StatusForbidden, "Forbidden", errors.New("user did not author chirp"))
		return
	}

	deleteChirpStruct := database.DeleteChirpParams{
		ID:     chirp.ID,
		UserID: userId,
	}

	err = cfg.DB.DeleteChirp(r.Context(), deleteChirpStruct)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete chirp", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
