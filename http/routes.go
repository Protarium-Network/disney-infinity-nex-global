package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", handleRoot)
	mux.HandleFunc("/health", handleHealth)

	return mux
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Server is running!"))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]any{
		"success":   true,
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	status := http.StatusOK

	if DB == nil {
		response["database"] = "not configured"
	} else if err := DB.PingContext(r.Context()); err != nil {
		response["database"] = "disconnected"
		response["success"] = false
		response["status"] = "degraded"
		status = http.StatusServiceUnavailable
	} else {
		response["database"] = "connected"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
