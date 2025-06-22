package db

import (
	"OnTrek/utils"
	"database/sql"
)

func GetMembersInfoBySessionId(db *sql.DB, sessionId int) ([]utils.MemberInfo, error) {

	var members []utils.MemberInfo

	query := `SELECT u.id, u.username, sm.latitude, sm.longitude, sm.altitude, sm.accuracy, sm.help_request, sm.going_to, sm.timestamp FROM users u JOIN session_members sm ON u.id = sm.user_id WHERE sm.session_id = ?`
	rows, err := db.Query(query, sessionId)
	if err != nil {
		return []utils.MemberInfo{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var member utils.MemberInfo
		err := rows.Scan(&member.User.ID, &member.User.Username, &member.Latitude, &member.Longitude, &member.Altitude, &member.Accuracy, &member.HelpRequested, &member.GoingTo, &member.TimeStamp)
		if err != nil {
			return []utils.MemberInfo{}, err
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return []utils.MemberInfo{}, err
	}

	return members, nil

}
