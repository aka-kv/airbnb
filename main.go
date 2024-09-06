package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

type RoomData struct {
	OccupancyRate    float64 `json:"occupancy_rate"`
	AverageNightRate float64 `json:"average_night_rate"`
	HighestNightRate float64 `json:"highest_night_rate"`
	LowestNightRate  float64 `json:"lowest_night_rate"`
}

// getRoomData generates room data
func getRoomData(roomID string) (*RoomData, error) {
	// You can modify this function to generate random values if needed
	return &RoomData{
		OccupancyRate:    82.5,
		AverageNightRate: 120.00,
		HighestNightRate: 200.00,
		LowestNightRate:  90.00,
	}, nil
}

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

func rateLimitMiddleware(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Every(time.Minute), 5)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CreateRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(rateLimitMiddleware)
	r.HandleFunc("/api/room/{roomID:[0-9]+}", roomDataHandler).Methods("GET")
	return r
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	router := CreateRouter()
	router.ServeHTTP(w, r)
}
