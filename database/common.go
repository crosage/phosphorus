package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // MySQL 驱动
	_ "github.com/google/uuid"
	"github.com/rs/zerolog/log"
	_ "time"
)

var db *sql.DB

func InitDatabase() {
	var err error
	dsn := "admin:powerdealer213@tcp(47.109.51.81:3306)/powerdealer?charset=utf8&parseTime=True&loc=Local&serverTimezone=UTC"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to connect to database")
	}
	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to ping database")
	}
	createTables()
}

func createTables() {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(100) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		avatar_url VARCHAR(255),
		status ENUM('在线', '离线', '忙碌', '隐身') DEFAULT '离线',
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to create users table")
	}

	_, err = db.Exec(`
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
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to create private_messages table")
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS group_messages (
		id VARCHAR(36) PRIMARY KEY,
		sender_id INT NOT NULL,
		group_id INT NOT NULL,
		message TEXT,
		message_type ENUM('文本', '图片', '视频', '语音', '文件') DEFAULT '文本',
		file_url VARCHAR(255) NULL,
		reply_to VARCHAR(36) NULL,
		sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
		FOREIGN KEY (reply_to) REFERENCES group_messages(id) ON DELETE SET NULL
	);`)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to create group_messages table")
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS groups (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		creator_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
	);`)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to create groups table")
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS group_members (
		group_id INT NOT NULL,
		user_id INT NOT NULL,
		role ENUM('成员', '管理员', '群主') DEFAULT '成员',
		joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (group_id, user_id),
		FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to create group_members table")
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS contacts (
		user_id INT NOT NULL,
		contact_id INT NOT NULL,
		status ENUM('好友', '黑名单', '待确认') DEFAULT '好友',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, contact_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (contact_id) REFERENCES users(id) ON DELETE CASCADE
	);`)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to create contacts table")
	}
}
