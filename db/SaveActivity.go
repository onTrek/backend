package db

import (
	"OnTrek/utils"
	"database/sql"
	"time"
)

func SaveActivity(db *sql.DB, activity utils.Activity) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO activities (user_id, title, description, start_time, end_time, created_at, distance, total_ascent, total_descent, starting_elevation, maximum_altitude, average_speed) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(activity.UserID, activity.Title, activity.Description, activity.StartTime, activity.EndTime, time.Now().Format(time.RFC3339), activity.Distance, activity.TotalAscent, activity.TotalDescent, activity.StartingElevation, activity.MaximumElevation, activity.AverageSpeed)
	if err != nil {
		return err
	}

	return nil
}
