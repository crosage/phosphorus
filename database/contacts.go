package database

import "database/sql"

func AddFriend(userID, contactID int) error {
	var status string
	queryCheck := `SELECT status FROM contacts WHERE user_id = ? AND contact_id = ?`
	err := db.QueryRow(queryCheck, userID, contactID).Scan(&status)

	if err == nil {
		if status == "好友" {
			return nil
		}
		return nil
	} else if err != sql.ErrNoRows {
		return err
	}

	queryInsert := `
		INSERT INTO contacts (user_id, contact_id, status) VALUES (?, ?, '待确认')
	`
	_, err = db.Exec(queryInsert, userID, contactID)
	return err
}

func ConfirmFriend(userID, contactID int) error {
	query := `
		UPDATE contacts 
		SET status = '好友'
		WHERE (user_id = ? AND contact_id = ?) OR (user_id = ? AND contact_id = ?)
	`
	_, err := db.Exec(query, userID, contactID, contactID, userID)
	return err
}

func GetPendingFriendRequests(userID int) ([]map[string]interface{}, error) {
	query := `
		SELECT u.id, u.username, c.created_at
		FROM contacts c
		JOIN users u ON c.user_id = u.id
		WHERE c.contact_id = ? AND c.status = '待确认'
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pendingRequests []map[string]interface{}

	for rows.Next() {
		var (
			friendID  int
			username  string
			createdAt string
		)

		err := rows.Scan(&friendID, &username, &createdAt)
		if err != nil {
			return nil, err
		}

		request := map[string]interface{}{
			"id":         friendID,
			"username":   username,
			"created_at": createdAt,
		}
		pendingRequests = append(pendingRequests, request)
	}

	return pendingRequests, nil
}
