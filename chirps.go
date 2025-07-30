package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/reiffle/chirpy/internal/auth"
	"github.com/reiffle/chirpy/internal/database"
)

// Create Chirps
func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	secret := cfg.secret
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	decoded_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token in make", err)
		return
	}

	uid, err := auth.ValidateJWT(decoded_token, secret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token in auth", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	body := getCleanedBody(params.Body, badWords)

	chirp_params := database.CreateChirpParams{
		Body:   body,
		UserID: uid,
	}
	new_chirp, err := cfg.DB.CreateChirp(context.Background(), chirp_params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", nil)
		return
	}
	cleaned_chirp := returnVals{
		ID:        new_chirp.ID,
		CreatedAt: new_chirp.CreatedAt,
		UpdatedAt: new_chirp.UpdatedAt,
		Body:      new_chirp.Body,
		UserID:    new_chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, cleaned_chirp)

}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

/*
func cleanMessage(message string) string {
	const replaceWord = "****"
	naughtyWords := []string{"kerfuffle", "sharbert", "fornax"}
	splitMessage := strings.Split(message, " ")
	for x, word := range splitMessage {
		for _, naughtyWord := range naughtyWords {
			if word == naughtyWord {
				splitMessage[x] = replaceWord
			}
		}
	}
	return strings.Join(splitMessage, " ")
}
*/

// Delete chirps
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

// Get chirps
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
