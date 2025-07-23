package main

import (
	"net/http"

	"github.com/reiffle/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	refToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't retrieve token", err)
		return
	}

	userID, err := cfg.DB.GetUserFromRefreshToken(r.Context(), refToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find user", err)
		return
	}

	accessToken, err := auth.MakeJWT(userID, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find user", err)
		return
	}

	type tokenStruct struct {
		Token string `json:"token"`
	}

	jsonToken := tokenStruct{Token: accessToken}
	respondWithJSON(w, http.StatusOK, jsonToken)
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	refToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't retrieve token", err)
		return
	}

	err = cfg.DB.RevokeRefreshToken(r.Context(), refToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
