package main

import (
	"encoding/json"
	"net/http"

	"github.com/reiffle/chirpy/internal/auth"
	"github.com/reiffle/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdatePassword(w http.ResponseWriter, r *http.Request) {
	secret := cfg.secret

	decoded_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	uid, err := auth.ValidateJWT(decoded_token, secret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}
	userParams := database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             uid,
	}

	returnedUser, err := cfg.DB.UpdateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't update user", err)
		return
	}
	CleanUser := User{
		ID:        returnedUser.ID,
		CreatedAt: returnedUser.CreatedAt,
		UpdatedAt: returnedUser.UpdatedAt,
		Email:     returnedUser.Email,
	}
	respondWithJSON(w, http.StatusOK, CleanUser)
}
