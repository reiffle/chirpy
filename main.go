package main

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func handlerHealthZ(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Println("Error writing response: ", err)
	}
}

func (cfg *apiConfig) handlerNumberHits(w http.ResponseWriter, r *http.Request) {
	requests := cfg.fileserverHits.Load()
	requestString := fmt.Sprintf("Hits: %d", requests)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	_, err := w.Write([]byte(requestString))
	if err != nil {
		fmt.Println("Error writing response: ", err)
	}
}

func (cfg *apiConfig) handlerResetHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	requests := cfg.fileserverHits.Load()
	requestString := fmt.Sprintf("Hits: %d", requests)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	_, err := w.Write([]byte(requestString))
	if err != nil {
		fmt.Println("Error writing response: ", err)
	}
}

func main() {
	apiCFG := &apiConfig{}
	mux := http.NewServeMux() //Create a new traffic router
	mux.Handle("/app/", apiCFG.middlewareMetricsInc(http.StripPrefix("/app/", (http.FileServer(http.Dir("."))))))
	mux.HandleFunc("GET /healthz", handlerHealthZ)
	mux.HandleFunc("GET /metrics", apiCFG.handlerNumberHits)
	mux.HandleFunc("POST /reset", apiCFG.handlerResetHits)

	//Create a new server on port :8080 and put the router in it
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := s.ListenAndServe() // Start up the server

	if err != nil {
		fmt.Println("Couldn't start server")
		os.Exit(1)
	}
}
