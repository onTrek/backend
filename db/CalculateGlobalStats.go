package db

import (
	"OnTrek/utils"
	"database/sql"
	"fmt"
)

func CalculateGlobalStats(db *sql.DB, userID string) (utils.GlobalStats, error) {
	var globalStats utils.GlobalStats

	query := `
		SELECT 
			COUNT(DISTINCT a.id) AS total_activities,
			SUM(a.distance) AS total_distance,
			SUM(a.total_ascent) AS total_ascent,
			SUM(a.total_descent) AS total_descent,
			SUM(a.end_time - a.start_time) AS total_time
		FROM activities a
		WHERE a.user_id = ?
	`

	err := db.QueryRow(query, userID).Scan(
		&globalStats.TotalActivities,
		&globalStats.TotalDistance,
		&globalStats.TotalAscent,
		&globalStats.TotalDescent,
		&globalStats.TotalTime,
	)

	if err != nil {
		return globalStats, fmt.Errorf("failed to calculate global stats: %v", err)
	}

	return globalStats, nil
}
