package main

import (
	"encoding/json"
	"net/http"

	"github.com/reiffle/chirpy/internal/auth"
	"github.com/reiffle/chirpy/internal/database"
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
	err = auth.CheckPassword(params.Password, full_user.HashedPassword)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Username and Password do not match", err)
		return
	}

	token, err := auth.MakeJWT(full_user.ID, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create tokan", err)
		return
	}

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token string", err)
	}

	tokenParams := database.CreateRefreshTokenParams{Token: refreshTokenString, UserID: full_user.ID}
	refresh_token, err := cfg.DB.CreateRefreshToken(r.Context(), tokenParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token string", err)
	}
	user := User{
		ID:           full_user.ID,
		CreatedAt:    full_user.CreatedAt,
		UpdatedAt:    full_user.UpdatedAt,
		Email:        full_user.Email,
		Token:        token,
		RefreshToken: refresh_token.Token,
	}

	respondWithJSON(w, http.StatusOK, user)
}
