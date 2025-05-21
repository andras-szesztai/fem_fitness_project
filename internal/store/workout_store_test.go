package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost port=5433 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	t.Log("Successfully connected to test database")

	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	t.Log("Successfully ran migrations")

	_, err = db.Exec(`TRUNCATE TABLE workouts, workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("Failed to truncate test database: %v", err)
	}

	t.Log("Successfully truncated tables")

	return db
}

func teardownTestDB(db *sql.DB) {
	db.Close()
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	workoutStore := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{name: "Valid workout with reps", workout: &Workout{Title: "Test Workout", Description: "Test Description", DurationMinutes: 30, CaloriesBurned: 300, Entries: []WorkoutEntry{
			{
				ExerciseName: "Squats",
				Sets:         3,
				Reps:         &[]int{12}[0],
				Weight:       &[]float64{100.5}[0],
				Notes:        "Felt strong today",
				OrderIndex:   1,
			},
		}}, wantErr: false},
		{name: "Valid workout with duration", workout: &Workout{Title: "Test Workout", Description: "Test Description", DurationMinutes: 30, CaloriesBurned: 300, Entries: []WorkoutEntry{
			{
				ExerciseName:    "Plank",
				Sets:            1,
				DurationSeconds: &[]int{60}[0],
				Weight:          &[]float64{0}[0],
				Notes:           "Felt strong today",
				OrderIndex:      1,
			},
		}}, wantErr: false},
		{name: "Invalid workout - both reps and duration", workout: &Workout{Title: "Test Workout", Description: "Test Description", DurationMinutes: 30, CaloriesBurned: 300, Entries: []WorkoutEntry{
			{
				ExerciseName:    "Squats",
				Sets:            3,
				Reps:            &[]int{12}[0],
				DurationSeconds: &[]int{12}[0],
				Weight:          &[]float64{100.5}[0],
				Notes:           "Felt strong today",
				OrderIndex:      1,
			},
		}}, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			createdWorkout, err := workoutStore.CreateWorkout(test.workout)

			if test.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.workout.Title, createdWorkout.Title)
			assert.Equal(t, test.workout.Description, createdWorkout.Description)
			assert.Equal(t, test.workout.DurationMinutes, createdWorkout.DurationMinutes)
			assert.Equal(t, test.workout.CaloriesBurned, createdWorkout.CaloriesBurned)
			assert.Equal(t, len(test.workout.Entries), len(createdWorkout.Entries))
			assert.NotNil(t, createdWorkout.ID)

			retreived, err := workoutStore.GetWorkout(createdWorkout.ID)

			require.NoError(t, err)
			assert.Equal(t, createdWorkout.ID, retreived.ID)
			assert.Equal(t, createdWorkout.Title, retreived.Title)
			assert.Equal(t, createdWorkout.Description, retreived.Description)
			assert.Equal(t, createdWorkout.DurationMinutes, retreived.DurationMinutes)
			assert.Equal(t, createdWorkout.CaloriesBurned, retreived.CaloriesBurned)
			assert.Equal(t, len(createdWorkout.Entries), len(retreived.Entries))
			for i, entry := range createdWorkout.Entries {
				assert.Equal(t, entry.ExerciseName, retreived.Entries[i].ExerciseName)
				assert.Equal(t, entry.Sets, retreived.Entries[i].Sets)
				assert.Equal(t, entry.Reps, retreived.Entries[i].Reps)
				assert.Equal(t, entry.Weight, retreived.Entries[i].Weight)
				if entry.DurationSeconds != nil {
					assert.Equal(t, *entry.DurationSeconds, *retreived.Entries[i].DurationSeconds)
				}
			}
		})
	}
}
