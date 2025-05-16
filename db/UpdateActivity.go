package db

import (
	"OnTrek/utils"
	"database/sql"
)

func UpdateActivity(db *sql.DB, activity utils.Activity) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("UPDATE activities SET title = ?, description = ?, start_time = ?, end_time = ?, distance = ?, total_ascent = ?, total_descent = ?, starting_elevation = ?, maximum_altitude = ?, average_speed = ?, average_heart_rate = ? WHERE id = ? AND user_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(activity.Title, activity.Description, activity.StartTime, activity.EndTime, activity.Distance, activity.TotalAscent, activity.TotalDescent, activity.StartingElevation, activity.MaximumElevation, activity.AverageSpeed, activity.AverageHeartRate, activity.ID, activity.UserID)
	if err != nil {
		return err
	}

	return nil
}
