package db

import (
	"OnTrek/utils"
	"database/sql"
)

func GetSessionInfoMember(db *sql.DB, sessionId int, userID string) (utils.SessionInfoResponse, error) {
	var sessionInfo utils.SessionInfoResponse
	query := ` SELECT u.id, u.name, u.email, s.created_at, s.closed_at FROM users u	JOIN sessions s ON u.id = s.created_by WHERE s.id = ?`
	row := db.QueryRow(query, sessionId)
	err := row.Scan(&sessionInfo.CreatedBy.ID, &sessionInfo.CreatedBy.Name, &sessionInfo.CreatedBy.Email, &sessionInfo.CreatedAt, &sessionInfo.ClosedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return sessionInfo, nil
		}
		return sessionInfo, err
	}

	var members []utils.MemberInfo
	query = `SELECT u.id, u.name, u.email, sm.latitude, sm.longitude, sm.altitude, sm.accuracy, sm.timestamp FROM users u JOIN session_members sm ON u.id = sm.user_id WHERE sm.session_id = ?`
	rows, err := db.Query(query, sessionId)
	if err != nil {
		return sessionInfo, err
	}
	defer rows.Close()

	for rows.Next() {
		var member utils.MemberInfo
		err := rows.Scan(&member.User.ID, &member.User.Name, &member.User.Email, &member.SessionInfo.Latitude, &member.SessionInfo.Longitude, &member.SessionInfo.Altitude, &member.SessionInfo.Accuracy, &member.TimeStamp)
		if err != nil {
			return sessionInfo, err
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return sessionInfo, err
	}

	sessionInfo.Members = members

	return sessionInfo, nil
}
