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
			COALESCE(SUM(a.distance), 0) AS total_distance,
			COALESCE(SUM(a.total_ascent), 0) AS total_ascent,
			COALESCE(SUM(a.total_descent), 0) AS total_descent,
			COALESCE(SUM(a.end_time - a.start_time), 0) AS total_time
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
