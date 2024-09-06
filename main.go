package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

// RoomData represents the structure of the data we return
type RoomData struct {
	OccupancyRate    float64 `json:"occupancy_rate"`
	AverageNightRate float64 `json:"average_night_rate"`
	HighestNightRate float64 `json:"highest_night_rate"`
	LowestNightRate  float64 `json:"lowest_night_rate"`
}

// getRoomData mocks the data retrieval for the example
func getRoomData(roomID string) (*RoomData, error) {
	return &RoomData{
		OccupancyRate:    85.5,
		AverageNightRate: 120.00,
		HighestNightRate: 200.00,
		LowestNightRate:  90.00,
	}, nil
}

// roomDataHandler handles requests to get room data
func roomDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomID"]

	if roomID == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}

	data, err := getRoomData(roomID)
	if err != nil {
		http.Error(w, "Error fetching room data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

var limiter = rate.NewLimiter(rate.Every(time.Minute), 5)

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Exported function that Vercel will call
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	r.Use(rateLimitMiddleware)
	r.HandleFunc("/api/room/{roomID:[0-9]+}", roomDataHandler).Methods("GET")
}
