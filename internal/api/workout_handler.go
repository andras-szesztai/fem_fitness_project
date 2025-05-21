package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/andras-szesztai/fem_fitness_project/internal/store"
	"github.com/andras-szesztai/fem_fitness_project/internal/utils"
)

type WorkoutHandler struct {
	store  store.WorkoutStore
	logger *log.Logger
}

func NewWorkoutHandler(store store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{store: store, logger: logger}
}

func (wh *WorkoutHandler) HandleGetWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %s", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	workout, err := wh.store.GetWorkout(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkout: %s", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to get workout"})
		return
	}

	if workout == nil {
		wh.logger.Printf("ERROR: getWorkout: %s", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout not found"})
		return
	}

	wh.logger.Printf("INFO: getWorkout: %d", workoutID)
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": workout})

}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: decodeCreateWorkoutBody: %s", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request body"})
		return
	}

	createdWorkout, err := wh.store.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: createWorkout: %s", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to create workout"})
		return
	}

	wh.logger.Printf("INFO: createWorkout: %d", createdWorkout.ID)
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"data": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %s", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	existingWorkout, err := wh.store.GetWorkout(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkout: %s", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to get workout"})
		return
	}
	if existingWorkout == nil {
		wh.logger.Printf("ERROR: getWorkout: %s", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout not found"})
		return
	}

	var updatedWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updatedWorkoutRequest)
	if err != nil {
		wh.logger.Printf("ERROR: decodeUpdateWorkoutBody: %s", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request body"})
		return
	}

	if updatedWorkoutRequest.Title != nil {
		existingWorkout.Title = *updatedWorkoutRequest.Title
	}

	if updatedWorkoutRequest.Description != nil {
		existingWorkout.Description = *updatedWorkoutRequest.Description
	}

	if updatedWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updatedWorkoutRequest.DurationMinutes
	}

	if updatedWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updatedWorkoutRequest.CaloriesBurned
	}

	if updatedWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updatedWorkoutRequest.Entries
	}

	err = wh.store.UpdateWorkout(existingWorkout)
	if err != nil {
		wh.logger.Printf("ERROR: updateWorkout: %s", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to update workout"})
		return
	}

	wh.logger.Printf("INFO: updateWorkout: %d", existingWorkout.ID)
	utils.WriteJSON(w, http.StatusNoContent, utils.Envelope{})
}

func (wh *WorkoutHandler) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %s", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout ID"})
		return
	}

	err = wh.store.DeleteWorkout(workoutID)
	if err != nil {
		if err == sql.ErrNoRows {
			wh.logger.Printf("ERROR: getWorkout: %s", err)
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout not found"})
			return
		}
		wh.logger.Printf("ERROR: deleteWorkout: %s", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to delete workout"})
		return
	}

	wh.logger.Printf("INFO: deleteWorkout: %d", workoutID)
	utils.WriteJSON(w, http.StatusNoContent, utils.Envelope{})
}
