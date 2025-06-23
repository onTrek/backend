package functions

import (
	"OnTrek/utils"
	"database/sql"
	"time"
)

func UpdateGroup(db *sql.DB, userId string, group utils.GroupInfo) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("UPDATE group_members SET latitude = ?, longitude = ?, altitude = ?, accuracy = ?, help_request = ?, going_to = ?, timestamp = ? WHERE group_id = ? AND user_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement with the provided parameters
	_, err = stmt.Exec(group.Latitude, group.Longitude, group.Altitude, group.Accuracy, group.HelpRequest, group.GoingTo, time.Now().Format(time.RFC3339), group.GroupID, userId)
	if err != nil {
		return err
	}

	return nil
}
