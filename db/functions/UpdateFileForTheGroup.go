package functions

import "database/sql"

func UpdateFileForTheGroup(db *sql.DB, groupId int, fileId int) error {
	// Prepare the SQL statement to update the file ID for the group
	query := `UPDATE groups SET file_id = ? WHERE id = ?`
	_, err := db.Exec(query, fileId, groupId)
	if err != nil {
		return err
	}
	return nil
}
