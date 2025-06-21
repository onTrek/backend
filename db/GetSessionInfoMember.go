package db

import (
	"OnTrek/utils"
	"database/sql"
	"errors"
)

func GetSessionInfo(db *sql.DB, sessionId int, userID string) (utils.SessionInfoResponse, error) {
	var sessionInfo utils.SessionInfoResponse
	query := ` SELECT u.id, u.username, s.description, s.created_at, s.closed_at FROM users u JOIN sessions s ON u.id = s.created_by WHERE s.id = ?`
	row := db.QueryRow(query, sessionId)
	err := row.Scan(&sessionInfo.CreatedBy.ID, &sessionInfo.CreatedBy.Username, &sessionInfo.Description, &sessionInfo.CreatedAt, &sessionInfo.ClosedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sessionInfo, nil
		}
		return sessionInfo, err
	}

	return sessionInfo, nil
}
