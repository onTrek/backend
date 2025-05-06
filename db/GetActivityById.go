package db

import (
	"OnTrek/utils"
	"database/sql"
)

func GetActivityByID(db *sql.DB, activityID int) (utils.Activity, error) {
	var activity utils.Activity

	err := db.QueryRow("SELECT * FROM activities WHERE id = ?", activityID).Scan(
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
