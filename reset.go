package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(403)
		w.Write([]byte("Forbidden"))
		return
	}
	err := cfg.DB.ResetUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to Process Request"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All Users Deleted"))
}
