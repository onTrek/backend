package functions

import (
	"database/sql"
	"fmt"
	"time"
)

func JoinGroup(tx *sql.Tx, userId string, groupId int) error {

	// Enable foreign key enforcement
	_, err := tx.Exec("PRAGMA foreign_keys = ON") // Enable foreign key enforcement
	if err != nil {
		return fmt.Errorf("error enabling foreign key enforcement: %v", err)
	}

	// Prepare the SQL statement
	stmt, err := tx.Prepare(`
		INSERT INTO group_members (group_id, user_id, timestamp, color)
		VALUES (?, ?, ?, (
			SELECT color
			FROM (
				SELECT '#e6194b' AS color UNION ALL
				SELECT '#3cb44b' UNION ALL
				SELECT '#ffe119' UNION ALL
				SELECT '#4363d8' UNION ALL
				SELECT '#f58231' UNION ALL
				SELECT '#911eb4' UNION ALL
				SELECT '#46f0f0' UNION ALL
				SELECT '#f032e6' UNION ALL
				SELECT '#bcf60c' UNION ALL
				SELECT '#fabebe' UNION ALL
				SELECT '#008080' UNION ALL
				SELECT '#e6beff' UNION ALL
				SELECT '#9a6324' UNION ALL
				SELECT '#fffac8' UNION ALL
				SELECT '#800000' UNION ALL
				SELECT '#aaffc3' UNION ALL
				SELECT '#808000' UNION ALL
				SELECT '#ffd8b1' UNION ALL
				SELECT '#000075' UNION ALL
				SELECT '#808080'
			) AS color_pool
			WHERE color NOT IN (
				SELECT color
				FROM group_members
				WHERE group_id = ?
			)
			LIMIT 1
		));
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now().Format(time.RFC3339)

	_, err = stmt.Exec(groupId, userId, now, groupId)
	if err != nil {
		return err
	}

	return nil
}
