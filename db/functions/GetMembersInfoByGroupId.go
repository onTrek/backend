package functions

import (
	"OnTrek/utils"
	"database/sql"
)

func GetMembersInfoByGroupId(db *sql.DB, groupId int) ([]utils.MemberInfo, error) {

	var members []utils.MemberInfo
	var GoingTo *string

	query := `SELECT u.id, u.username, gm.latitude, gm.longitude, gm.altitude, gm.accuracy, gm.help_request, gm.going_to, gm.timestamp FROM users u JOIN group_members gm ON u.id = gm.user_id WHERE gm.group_id = ?`
	rows, err := db.Query(query, groupId)
	if err != nil {
		return []utils.MemberInfo{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var member utils.MemberInfo
		err := rows.Scan(&member.User.ID, &member.User.Username, &member.Latitude, &member.Longitude, &member.Altitude, &member.Accuracy, &member.HelpRequested, &GoingTo, &member.TimeStamp)
		if err != nil {
			return []utils.MemberInfo{}, err
		}

		if GoingTo != nil {
			member.GoingTo = *GoingTo
		} else {
			member.GoingTo = ""
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return []utils.MemberInfo{}, err
	}

	return members, nil

}
