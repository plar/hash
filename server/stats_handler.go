package server

import (
	"encoding/json"
	"net/http"

	"github.com/plar/hash/service/stats"
)

type statsHandler struct {
	svc stats.Service
}

type statsResponse struct {
	Total   int64 `json:"total"`
	Average int64 `json:"average"`
}

func (h statsHandler) stats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	m := h.svc.Metric("Hasher.Create")
	if err := json.NewEncoder(w).Encode(statsResponse{
		Total:   m.Count(),
		Average: m.Average().Microseconds(),
	}); err != nil {
		http.Error(w, "Cannot encode stats", http.StatusInternalServerError)
	}
}
