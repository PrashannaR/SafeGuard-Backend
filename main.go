package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Coordinates struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

var (
	mutex            sync.RWMutex
	latestCoordinates Coordinates
)

func coordinatesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getCoordinates(w, r)
	case http.MethodPost:
		postCoordinates(w, r)
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func postCoordinates(w http.ResponseWriter, r *http.Request) {
	var coords Coordinates
	if err := json.NewDecoder(r.Body).Decode(&coords); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	latestCoordinates = coords
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Coordinates received"})
}

func getCoordinates(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	json.NewEncoder(w).Encode(latestCoordinates)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/coordinates", coordinatesHandler)
	mux.HandleFunc("/health", healthCheckHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Println("Server starting on port 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
