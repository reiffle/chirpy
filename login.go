package main

import (
	"encoding/json"
	"net/http"

	"github.com/reiffle/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string
		Password string
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
	err = auth.CheckPassword(params.Password, full_user.HashedPasswords.String)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Username and Password do not match", err)
		return
	}

	user := User{
		ID:        full_user.ID,
		CreatedAt: full_user.CreatedAt,
		UpdatedAt: full_user.UpdatedAt,
		Email:     full_user.Email,
	}
	respondWithJSON(w, http.StatusOK, user)
}
