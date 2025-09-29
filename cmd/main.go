package main

import (
	"log"
	"net/http"
	"procrastigo/internal/config"
	"procrastigo/internal/handlers"
	"procrastigo/internal/storage"
	"procrastigo/pkg/logger"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	logger.Init(cfg.LogLevel())

	store := storage.NewMemoryStorage()
	if err := store.LoadFromFile("data/excuses.json"); err != nil {
		log.Printf("Failed to load excuses: %v", err)
	}

	excuseHandler := handlers.NewExcuseHandler(store)
	statsHandler := handlers.NewStatsHandler(store)

	router := mux.NewRouter()

	v1 := router.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/excuses/random", excuseHandler.GetRandomExcuse).Methods("GET")
	v1.HandleFunc("/excuses", excuseHandler.GetExcuses).Methods("GET")
	v1.HandleFunc("/excuses", excuseHandler.CreateExcuse).Methods("POST")

	v1.HandleFunc("/stats", statsHandler.GetStats).Methods("GET")

	router.Use(handlers.LoggingMiddleware)
	router.Use(handlers.CORSMiddleware)

	log.Printf("🚀 Server starting on %s", cfg.ServerAddress())
	log.Fatal(http.ListenAndServe(cfg.ServerAddress(), router))
}
