package functions

import (
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"time"
)

func CreateGroup(db *sql.DB, user utils.User, info utils.GroupInfo) (utils.Group, error) {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return utils.Group{}, err
	}

	// Rollback on error unless committed
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Enable foreign key enforcement
	_, err = tx.Exec("PRAGMA foreign_keys = ON") // Enable foreign key enforcement
	if err != nil {
		return utils.Group{}, fmt.Errorf("error enabling foreign key enforcement: %v", err)
	}

	// Insert into sessions
	now := time.Now().Format(time.RFC3339)
	stmt1, err := tx.Prepare("INSERT INTO groups (created_by, description, created_at, file_id, last_update) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return utils.Group{}, err
	}
	defer stmt1.Close()

	var fileID interface{}
	if info.FileId == -1 {
		fileID = nil
	} else {
		fileID = info.FileId
	}
	res, err := stmt1.Exec(user.ID, info.Description, now, fileID, now)
	if err != nil {
		return utils.Group{}, err
	}

	groupID, err := res.LastInsertId()
	if err != nil {
		return utils.Group{}, err
	}

	err = JoinGroup(tx, user.ID, int(groupID))
	if err != nil {
		return utils.Group{}, fmt.Errorf("error joining group: %v", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return utils.Group{}, err
	}

	return utils.Group{ID: int(groupID), CreatedBy: user.ID}, nil
}
