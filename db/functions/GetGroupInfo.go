package functions

import (
	"OnTrek/utils"
	"database/sql"
	"errors"
)

func GetGroupInfo(db *sql.DB, groupId int) (utils.GroupInfoResponse, error) {
	var groupInfo utils.GroupInfoResponse
	query := ` SELECT u.id, u.username, g.description, g.created_at FROM users u JOIN groups g ON u.id = g.created_by WHERE g.id = ?`
	row := db.QueryRow(query, groupId)
	err := row.Scan(&groupInfo.CreatedBy.ID, &groupInfo.CreatedBy.Username, &groupInfo.Description, &groupInfo.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return groupInfo, nil
		}
		return groupInfo, err
	}

	query = `SELECT u.id, u.username, gm.color FROM users u JOIN group_members gm ON u.id = gm.user_id WHERE gm.group_id = ?`
	rows, err := db.Query(query, groupId)
	if err != nil {
		return utils.GroupInfoResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var member utils.GroupMember
		err := rows.Scan(&member.ID, &member.Username, &member.Color)
		if err != nil {
			return utils.GroupInfoResponse{}, err
		}
		groupInfo.Members = append(groupInfo.Members, member)
	}

	return groupInfo, nil
}
