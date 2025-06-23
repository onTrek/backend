package functions

import (
	"OnTrek/utils"
	"database/sql"
	"fmt"
)

func GetFriendRequestsByUserId(db *sql.DB, userId string) ([]utils.UserEssentials, error) {
	var friendRequests []utils.UserEssentials

	rows, err := db.Query("SELECT u.id, u.username FROM users u JOIN friends f ON u.id = f.user_id1 WHERE f.user_id2 = ? AND f.pending = TRUE", userId)
	if err != nil {
		fmt.Println("Error querying database:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var request utils.UserEssentials
		err := rows.Scan(&request.ID, &request.Username)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		friendRequests = append(friendRequests, request)
	}

	return friendRequests, nil
}
