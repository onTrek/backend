package functions

import (
	"OnTrek/utils"
	"database/sql"
)

func GetGroupsByUserId(db *sql.DB, userId string) ([]utils.Group, error) {

	var groups []utils.Group
	// Get the Groups for the user
	query := `SELECT g.id, g.description, g.created_by, g.created_at, gf.id, gf.filename FROM groups g JOIN group_members gm ON g.id = gm.group_id LEFT JOIN gpx_files gf ON g.file_id = gf.id WHERE gm.user_id = ? ORDER BY g.last_update DESC`
	rows, err := db.Query(query, userId)
	if err != nil {
		return []utils.Group{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var s utils.Group
		var fileID sql.NullInt64
        	var fileName sql.NullString
		
		err := rows.Scan(&s.ID, &s.Description, &s.CreatedBy, &s.CreatedAt, &fileID, &fileName)

	        if fileID.Valid {
	            s.File.ID = int(fileID.Int64)
	        } else {
	            s.File.ID = -1
	        }
	
	        if fileName.Valid {
	            s.File.Filename = fileName.String
	        } else {
	            s.File.Filename = ""
	        }
		
		if err != nil {
			return []utils.Group{}, err
		}
		groups = append(groups, s)
	}

	if err := rows.Err(); err != nil {
		return []utils.Group{}, err
	}

	if len(groups) == 0 {
		return []utils.Group{}, nil
	} else {
		return groups, nil
	}
}
