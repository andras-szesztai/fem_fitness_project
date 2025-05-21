package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andras-szesztai/fem_fitness_project/internal/store"
	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	store store.WorkoutStore
}

func NewWorkoutHandler(store store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{store: store}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.Error(w, "Workout ID is required", http.StatusBadRequest)
		return
	}

	workoutID, err := strconv.Atoi(paramsWorkoutID)
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	workout, err := wh.store.GetWorkoutByID(workoutID)
	if err != nil {
		http.Error(w, "Failed to get workout", http.StatusInternalServerError)
		return
	}

	if workout == nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)

}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdWorkout, err := wh.store.CreateWorkout(&workout)
	if err != nil {
		fmt.Println("Error creating workout:", err)
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdWorkout)
}
