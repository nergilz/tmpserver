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
	Private bool    `json:"private"` // for tow
	// Public  bool    `json:"public"`  // for many people
}

// SendChatRequestModel send msg in chat request
type SendChatRequestModel struct {
	ID      int64  `json:"chat_id"`
	Content string `json:"content"`
}

// InitChartStore ..
func InitChartStore(db *database.DB, log *logger.Logger) *ChatStore {
	cs := new(ChatStore)
	cs.db = db
	log.Service("init chat store")
	return cs
}

// CreateChat ...
func (cs *ChatStore) CreateChat(cm *ChatModel) error { // сделать проверку SQL запросом
	var id int64
	q := `INSERT INTO chats (user_id, private) VALUES($1, $2) RETURNING id`
	if err := cs.db.Conn().QueryRow(
		q,
		pq.Array(cm.UserIDs),
		cm.Private).Scan(&id); err != nil {
		return err
	}
	cm.ID = id
	return nil
}

// CheckChat ..
func (cs *ChatStore) CheckChat() (bool, error) {
	var id int64
	q := `SELECT id FROM chats ORDER BY id DESC LIMIT 1 WHERE id=$1 VALUES($1)`
	if err := cs.db.Conn().QueryRow(q).Scan(
		&id,
	); err != nil {
		return false, err
	}
	if id != 0 {
		return false, nil
	}
	return true, nil
}

// UpdateChat ..
func (cs *ChatStore) UpdateChat(cm *ChatModel) error {
	// q := ``
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
