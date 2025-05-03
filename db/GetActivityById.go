package db

import (
	"OnTrek/utils"
	"database/sql"
)

func GetActivityById(db *sql.DB, userID string, activityID int) (utils.Activity, error) {
	var activity utils.Activity

	err := db.QueryRow("SELECT * FROM activities WHERE user_id = ? AND id = ?", userID, activityID).Scan(
		&activity.ID,
		&activity.UserID,
		&activity.Title,
		&activity.Description,
		&activity.StartTime,
		&activity.EndTime,
		&activity.CreatedAt,
		&activity.Distance,
		&activity.TotalAscent,
		&activity.TotalDescent,
		&activity.StartingElevation,
		&activity.MaximumElevation,
		&activity.AverageSpeed)
	if err != nil {
		return activity, err
	}

	return activity, nil
}
