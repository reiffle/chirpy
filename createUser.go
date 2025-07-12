package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/reiffle/chirpy/internal/auth"
	"github.com/reiffle/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashWord, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't encrypt password", err)
		return
	}

	userParams := database.CreateUserParams{
		Email:           params.Email,
		HashedPasswords: sql.NullString{String: hashWord, Valid: true},
	}
	user, err := cfg.DB.CreateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	responseUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(responseUser)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(dat)
}
