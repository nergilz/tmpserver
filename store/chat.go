package store

import (
	"github.com/0LuigiCode0/Library/logger"
	"github.com/nergilz/tmpserver/database"
)

// ChatStore chat store
type ChatStore struct {
	db *database.DB
}

// ChatModel chat model
type ChatModel struct {
	ID      int64   `json:"chat_id"`
	MsgIDs  []int64 `json:"messages_id"`
	UserIDs []int64 `json:"users_id"`
	Private bool    `json:"private"`
}

// InitChartStore ..
func InitChartStore(db *database.DB, log *logger.Logger) *ChatStore {
	cs := new(ChatStore)
	cs.db = db
	log.Service("init chat store")
	return cs
}

// CreateChat ...
func (cs *ChatStore) CreateChat(cm *ChatModel) error {
	var id int64
	q := `INSERT INTO chats (msg_id, user_id, private) VALUES ($1, $2, $3) RETURNING id`
	if err := cs.db.Conn().QueryRow(q, cm.MsgIDs, cm.UserIDs, cm.Private).Scan(&id); err != nil {
		return err
	}
	cm.ID = id
	return nil
}

// DeleteChat ...
func (cs *ChatStore) DeleteChat(chatID int64) error {
	q := `DELETE FROM chats WHERE id = $1`
	if err := cs.db.Conn().QueryRow(q, chatID).Err(); err != nil {
		return err
	}
	return nil
}

// DeleteUserInChat ...
func (cs *ChatStore) DeleteUserInChat(userID int64) error {

	return nil
}
