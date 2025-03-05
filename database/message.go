package database

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

/*
CREATE TABLE IF NOT EXISTS private_messages (

	id VARCHAR(36) PRIMARY KEY,
	sender_id INT NOT NULL,
	receiver_id INT NOT NULL,
	message TEXT,
	message_type ENUM('文本', '图片', '视频', '语音', '文件') DEFAULT '文本',
	file_url VARCHAR(255) NULL,
	reply_to VARCHAR(36) NULL,
	status ENUM('已发送', '已接收', '已读') DEFAULT '已发送',
	sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (reply_to) REFERENCES private_messages(id) ON DELETE SET NULL

);`)
*/
func SendMessage(senderID, receiverID int, message, msgType, fileURL, replyTo string) error {
	messageID := uuid.New().String()
	sentAt := time.Now().Unix()
	query := "INSERT INTO private_messages (id, sender_id, receiver_id, message, message_type, file_url, sent_at, reply_to) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, messageID, senderID, receiverID, message, msgType, fileURL, sentAt, replyTo)
	return err
}

func GetMessages(receiverID int) ([]map[string]interface{}, error) {
	query := "SELECT id, sender_id, receiver_id, message, message_type, file_url, reply_to, status, sent_at FROM private_messages WHERE receiver_id = ? ORDER BY sent_at DESC"

	rows, err := db.Query(query, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []map[string]interface{}

	for rows.Next() {
		var (
			id          string
			senderID    int
			receiverID  int
			message     sql.NullString
			messageType string
			fileURL     sql.NullString
			replyTo     sql.NullString
			status      string
			sentAt      int64
		)

		err := rows.Scan(&id, &senderID, &receiverID, &message, &messageType, &fileURL, &replyTo, &status, &sentAt)
		if err != nil {
			return nil, err
		}

		msg := map[string]interface{}{
			"id":           id,
			"sender_id":    senderID,
			"receiver_id":  receiverID,
			"message":      message.String,
			"message_type": messageType,
			"file_url":     fileURL.String,
			"reply_to":     replyTo.String,
			"status":       status,
			"sent_at":      sentAt,
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func GetConversation(user1ID, user2ID int) ([]map[string]interface{}, error) {
	query := `
		SELECT id, sender_id, receiver_id, message, message_type, file_url, reply_to, status, sent_at 
		FROM private_messages 
		WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) 
		ORDER BY sent_at ASC
	`

	rows, err := db.Query(query, user1ID, user2ID, user2ID, user1ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []map[string]interface{}

	for rows.Next() {
		var (
			id          string
			senderID    int
			receiverID  int
			message     sql.NullString
			messageType string
			fileURL     sql.NullString
			replyTo     sql.NullString
			status      string
			sentAt      int64
		)

		err := rows.Scan(&id, &senderID, &receiverID, &message, &messageType, &fileURL, &replyTo, &status, &sentAt)
		if err != nil {
			return nil, err
		}

		msg := map[string]interface{}{
			"id":           id,
			"sender_id":    senderID,
			"receiver_id":  receiverID,
			"message":      message.String,
			"message_type": messageType,
			"file_url":     fileURL.String,
			"reply_to":     replyTo.String,
			"status":       status,
			"sent_at":      sentAt,
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
