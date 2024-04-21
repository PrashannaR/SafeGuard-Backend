package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Coordinates struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

var (
	// Protects access to the coordinates
	mutex sync.RWMutex
	// Holds the latest coordinates
	latestCoordinates Coordinates
)

func coordinatesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getCoordinates(w, r)
	case http.MethodPost:
		postCoordinates(w, r)
	default:
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
	}
}

func postCoordinates(w http.ResponseWriter, r *http.Request) {
	var coords Coordinates
	err := json.NewDecoder(r.Body).Decode(&coords)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	latestCoordinates = coords
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Coordinates received")
}

func getCoordinates(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	coords := latestCoordinates
	mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coords)
}

func main() {
	http.HandleFunc("/api/coordinates", coordinatesHandler)

	log.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
