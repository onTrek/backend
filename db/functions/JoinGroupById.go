package functions

import "database/sql"

func JoinGroupById(db *sql.DB, userId string, groupId int) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Rollback on error unless committed
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Insert into group_members
	err = JoinGroup(tx, userId, groupId)
	if err != nil {
		return err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
