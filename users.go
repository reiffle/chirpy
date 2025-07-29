package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerIsChirpyRed(w http.ResponseWriter, r *http.Request) {
	type userID struct {
		UserID string `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  userID `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	stringID := params.Data.UserID
	uid, err := uuid.Parse(stringID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse UUID", err)
		return
	}

	err = cfg.DB.MakeChirpyRed(r.Context(), uid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
