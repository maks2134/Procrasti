package handlers

import (
	"encoding/json"
	"net/http"
	"procrastigo/internal/storage"
)

type StatsHandler struct {
	storage storage.Storage
}

func NewStatsHandler(storage storage.Storage) *StatsHandler {
	return &StatsHandler{storage: storage}
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.storage.GetStats()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
