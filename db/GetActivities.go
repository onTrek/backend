package db

import (
	"OnTrek/utils"
	"database/sql"
	"fmt"
)

func GetActivities(db *sql.DB, userID string) ([]utils.Activity, error) {
	var activities []utils.Activity

	rows, err := db.Query("SELECT * FROM activities WHERE user_id = ?", userID)
	if err != nil {
		fmt.Println("Error querying database:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var activity utils.Activity
		err := rows.Scan(
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
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}
