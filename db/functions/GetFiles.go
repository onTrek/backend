package functions

import (
	"OnTrek/utils"
	"database/sql"
	"fmt"
)

func GetFiles(db *sql.DB, userID string) ([]utils.GpxInfo, error) {
	var files []utils.GpxInfo

	rows, err := db.Query(`SELECT id, filename, upload_date, title, km, ascent, descent, duration, max_altitude, min_altitude FROM gpx_files WHERE user_id = ?`, userID)
	if err != nil {
		fmt.Println("Error querying database:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var file utils.GpxInfo
		err := rows.Scan(
			&file.ID,
			&file.Filename,
			&file.UploadDate,
			&file.Title,
			&file.Stats.Km,
			&file.Stats.Ascent,
			&file.Stats.Descent,
			&file.Stats.Duration,
			&file.Stats.MaxAltitude,
			&file.Stats.MinAltitude,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}
