package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/reiffle/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email              string
		Password           string
		Expires_in_seconds int
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	full_user, err := cfg.DB.FindUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Username and Password do not match", err)
		return
	}
	err = auth.CheckPassword(params.Password, full_user.HashedPassword)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Username and Password do not match", err)
		return
	}

	if params.Expires_in_seconds < 1 || params.Expires_in_seconds > 3600 {
		params.Expires_in_seconds = 3600
	}

	duration := time.Duration(params.Expires_in_seconds) * time.Second

	token, err := auth.MakeJWT(full_user.ID, cfg.secret, duration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create tokane", err)
		return
	}

	user := User{
		ID:        full_user.ID,
		CreatedAt: full_user.CreatedAt,
		UpdatedAt: full_user.UpdatedAt,
		Email:     full_user.Email,
		Token:     token,
	}
	respondWithJSON(w, http.StatusOK, user)
}
