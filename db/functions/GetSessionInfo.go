package functions

import (
	"OnTrek/utils"
	"database/sql"
	"errors"
)

func GetSessionInfo(db *sql.DB, sessionId int) (utils.SessionInfoResponse, error) {
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

	// Get the members of the session
	query = `SELECT u.id, u.username, sm.color FROM users u JOIN session_members sm ON u.id = sm.user_id WHERE sm.session_id = ?`
	rows, err := db.Query(query, sessionId)
	if err != nil {
		return utils.SessionInfoResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var member utils.SessionMember
		err := rows.Scan(&member.ID, &member.Username, &member.Color)
		if err != nil {
			return utils.SessionInfoResponse{}, err
		}
		sessionInfo.Members = append(sessionInfo.Members, member)
	}

	return sessionInfo, nil
}
