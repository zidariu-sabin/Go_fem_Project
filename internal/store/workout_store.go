package store

import (
	"database/sql"
	"fmt"
)

// adding this `json:` tag is a go feature that will allow us to assign a struct a json stucture aswell
type Workout struct {
	ID              int            `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

// reps, duration and weight are refferences because we want to check specifically if they are set to null exercise_mode comparison
type WorkoutEntry struct {
	ID              int      `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

// Store used for postgres database operations
type PostgresWorkoutStore struct {
	db *sql.DB
}

// postgres workout store constructore
func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

// We create an interface in order to separate database operations from store methods so the application store is not bound to postgress
type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutByID(id int64) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkout(id int64) error
}

func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	//transaction
	tx, err := pg.db.Begin()

	if err != nil {
		return nil, err
	}
	//rollback is defered so it only proceeds at the end if there are caught errors within the transaction object
	defer tx.Rollback()

	query :=
		`INSERT INTO  workouts(title, description, duration_minutes, calories_burned)
	VALUES ($1,$2,$3,$4)
	RETURNING id
	`

	err = tx.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)

	if err != nil {
		return nil, err
	}

	//inserting entries
	for i, entry := range workout.Entries {
		query := `INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index ) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
		`

		err = tx.QueryRow(query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&workout.Entries[i].ID)

		if err != nil {
			return nil, err
		}

	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return workout, nil

}

func (pg *PostgresWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {

	workout := &Workout{}

	queryWorkout := `SELECT id, title, description, duration_minutes, calories_burned 
	FROM workouts 
	WHERE id = $1`

	err := pg.db.QueryRow(queryWorkout, id).Scan(&workout.ID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	queryEntries := `SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index 
	FROM workout_entries 
	WHERE workout_id = $1
	ORDER BY order_index
	`

	rows, err := pg.db.Query(queryEntries, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var workout_entry WorkoutEntry
		err := rows.Scan(&workout_entry.ID,
			&workout_entry.ExerciseName,
			&workout_entry.Sets,
			&workout_entry.Reps,
			&workout_entry.DurationSeconds,
			&workout_entry.Weight,
			&workout_entry.Notes,
			&workout_entry.OrderIndex)
		if err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, workout_entry)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return workout, nil

}

// also used to delete/ update a workout entry
func (pg *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	fmt.Println("Updated workout with ID:", workout.ID)

	query := `
	UPDATE workouts
	SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4
	WHERE id = $5
	`
	result, err := tx.Exec(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned, workout.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	// we delete the existing entries for the workout to do a patch update to ease the data required for the client
	_, err = tx.Exec(`DELETE FROM workout_entries WHERE workout_id = $1`, workout.ID)

	if err != nil {
		return err
	}

	for _, entry := range workout.Entries {
		query := `
		INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index ) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		`

		_, err := tx.Exec(query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (pg *PostgresWorkoutStore) DeleteWorkout(id int64) error {
	query := `DELETE FROM workouts WHERE id = $1`

	result, err := pg.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
