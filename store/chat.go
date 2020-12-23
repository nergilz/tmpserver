package store

import (
	"github.com/0LuigiCode0/Library/logger"
	"github.com/lib/pq"
	"github.com/nergilz/tmpserver/database"
)

// ChatStore chat store
type ChatStore struct {
	db *database.DB
}

// ChatModel chat model
type ChatModel struct {
	ID      int64   `json:"chat_id"`
	Name    string  `json:"chat_name"`
	UserIDs []int64 `json:"users_ids"`
	Private bool    `json:"private"`
}

// InitChartStore ..
func InitChartStore(db *database.DB, log *logger.Logger) *ChatStore {
	cs := new(ChatStore)
	cs.db = db
	log.Service("init chat store")
	return cs
}

/*
	сделать проверку SQL запросом
*/
// CreateChat ...
func (cs *ChatStore) CreateChat(cm *ChatModel) error {
	var id int64
	q := `INSERT INTO chats (user_id, private) VALUES ($1, $2) RETURNING id`
	if err := cs.db.Conn().QueryRow(
		q,
		pq.Array(cm.UserIDs),
		cm.Private).Scan(&id); err != nil {
		return err
	}
	cm.ID = id
	return nil
}

// UpdateChat ..
func (cs *ChatStore) UpdateChat(cm *ChatModel) error {

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
